/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package transientstore

import (
	"errors"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/protos/ledger/rwset"

	"github.com/hyperledger/fabric/common/ledger/util/leveldbhelper"
	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/ledger"
	"github.com/hyperledger/fabric/core/ledger/pvtdatastorage"
	"github.com/syndtr/goleveldb/leveldb/iterator"
)

var emptyValue = []byte{}

// ErrStoreEmpty is used to indicate that there are no entries in transient store
var ErrStoreEmpty = errors.New("Transient store is empty")

//////////////////////////////////////////////
// Interfaces and data types
/////////////////////////////////////////////

// StoreProvider provides an instance of a TransientStore
type StoreProvider interface {
	OpenStore(ledgerID string) (Store, error)
	Close()
}

// RWSetScanner provides an iterator for EndorserPvtSimulationResults
type RWSetScanner interface {
	// Next returns the next EndorserPvtSimulationResults from the RWSetScanner.
	// It may return nil, nil when it has no further data, and also may return an error
	// on failure
	Next() (*EndorserPvtSimulationResults, error)
	// Close frees the resources associated with this RWSetScanner
	Close()
}

// Store manages the storage of private write sets for a ledgerId.
// Ideally, a ledger can remove the data from this storage when it is committed to
// the permanent storage or the pruning of some data items is enforced by the policy
type Store interface {
	// Persist stores the private write set of a transaction in the transient store
	// based on txid and the block height the private data was received at
	Persist(txid string, blockHeight uint64, privateSimulationResults *rwset.TxPvtReadWriteSet) error
	// GetTxPvtRWSetByTxid returns an iterator due to the fact that the txid may have multiple private
	// write sets persisted from different endorsers (via Gossip)
	GetTxPvtRWSetByTxid(txid string, filter ledger.PvtNsCollFilter) (RWSetScanner, error)
	// PurgeByTxids removes private write sets of a given set of transactions from the
	// transient store
	PurgeByTxids(txids []string) error
	// PurgeByHeight removes private write sets at block height lesser than
	// a given maxBlockNumToRetain. In other words, Purge only retains private write sets
	// that were persisted at block height of maxBlockNumToRetain or higher. Though the private
	// write sets stored in transient store is removed by coordinator using PurgebyTxids()
	// after successful block commit, PurgeByHeight() is still required to remove orphan entries (as
	// transaction that gets endorsed may not be submitted by the client for commit)
	PurgeByHeight(maxBlockNumToRetain uint64) error
	// GetMinTransientBlkHt returns the lowest block height remaining in transient store
	GetMinTransientBlkHt() (uint64, error)
	Shutdown()
}

// EndorserPvtSimulationResults captures the details of the simulation results specific to an endorser
type EndorserPvtSimulationResults struct {
	ReceivedAtBlockHeight uint64
	PvtSimulationResults  *rwset.TxPvtReadWriteSet
}

//////////////////////////////////////////////
// Implementation
/////////////////////////////////////////////

// storeProvider encapsulates a leveldb provider which is used to store
// private write sets of simulated transactions, and implements TransientStoreProvider
// interface.
type storeProvider struct {
	dbProvider *leveldbhelper.Provider
}

// store holds an instance of a levelDB.
type store struct {
	db       *leveldbhelper.DBHandle
	ledgerID string
}

type RwsetScanner struct {
	txid   string
	dbItr  iterator.Iterator
	filter ledger.PvtNsCollFilter
}

// NewStoreProvider instantiates TransientStoreProvider
func NewStoreProvider() StoreProvider {
	dbProvider := leveldbhelper.NewProvider(&leveldbhelper.Conf{DBPath: GetTransientStorePath()})
	return &storeProvider{dbProvider: dbProvider}
}

// OpenStore returns a handle to a ledgerId in Store
func (provider *storeProvider) OpenStore(ledgerID string) (Store, error) {
	dbHandle := provider.dbProvider.GetDBHandle(ledgerID)
	return &store{db: dbHandle, ledgerID: ledgerID}, nil
}

// Close closes the TransientStoreProvider
func (provider *storeProvider) Close() {
	provider.dbProvider.Close()
}

// Persist stores the private write set of a transaction in the transient store
// based on txid and the block height the private data was received at
func (s *store) Persist(txid string, blockHeight uint64,
	privateSimulationResults *rwset.TxPvtReadWriteSet) error {
	dbBatch := leveldbhelper.NewUpdateBatch()

	// Create compositeKey with appropriate prefix, txid, uuid and blockHeight
	// Due to the fact that the txid may have multiple private write sets persisted from different
	// endorsers (via Gossip), we postfix an uuid with the txid to avoid collision.
	uuid := util.GenerateUUID()
	compositeKeyPvtRWSet := createCompositeKeyForPvtRWSet(txid, uuid, blockHeight)
	privateSimulationResultsBytes, err := proto.Marshal(privateSimulationResults)
	if err != nil {
		return err
	}
	dbBatch.Put(compositeKeyPvtRWSet, privateSimulationResultsBytes)

	// Create two index: (i) by txid, and (ii) by height

	// Create compositeKey for purge index by height with appropriate prefix, blockHeight,
	// txid, uuid and store the compositeKey (purge index) with a null byte as value. Note that
	// the purge index is used to remove orphan entries in the transient store (which are not removed
	// by PurgeTxids()) using BTL policy by PurgeByHeight(). Note that orphan entries are due to transaction
	// that gets endorsed but not submitted by the client for commit)
	compositeKeyPurgeIndexByHeight := createCompositeKeyForPurgeIndexByHeight(blockHeight, txid, uuid)
	dbBatch.Put(compositeKeyPurgeIndexByHeight, emptyValue)

	// Create compositeKey for purge index by txid with appropriate prefix, txid, uuid,
	// blockHeight and store the compositeKey (purge index) with a null byte as value.
	// Though compositeKeyPvtRWSet itself can be used to purge private write set by txid,
	// we create a separate composite key with a null byte as value. The reason is that
	// if we use compositeKeyPvtRWSet, we unnecessarily read (potentially large) private write
	// set associated with the key from db. Note that this purge index is used to remove non-orphan
	// entries in the transient store and is used by PurgeTxids()
	// Note: We can create compositeKeyPurgeIndexByTxid by just replacing the prefix of compositeKeyPvtRWSet
	// with purgeIndexByTxidPrefix. For code readability and to be expressive, we use a
	// createCompositeKeyForPurgeIndexByTxid() instead.
	compositeKeyPurgeIndexByTxid := createCompositeKeyForPurgeIndexByTxid(txid, uuid, blockHeight)
	dbBatch.Put(compositeKeyPurgeIndexByTxid, emptyValue)

	return s.db.WriteBatch(dbBatch, true)
}

// GetTxPvtRWSetByTxid returns an iterator due to the fact that the txid may have multiple private
// write sets persisted from different endorsers.
func (s *store) GetTxPvtRWSetByTxid(txid string, filter ledger.PvtNsCollFilter) (RWSetScanner, error) {
	// Construct startKey and endKey to do an range query
	startKey := createTxidRangeStartKey(txid)
	endKey := createTxidRangeEndKey(txid)

	iter := s.db.GetIterator(startKey, endKey)
	return &RwsetScanner{txid, iter, filter}, nil
}

// PurgeByTxids removes private write sets of a given set of transactions from the
// transient store. PurgeByTxids() is expected to be called by coordinator after
// committing a block to ledger.
func (s *store) PurgeByTxids(txids []string) error {
	dbBatch := leveldbhelper.NewUpdateBatch()

	for _, txid := range txids {
		// Construct startKey and endKey to do an range query
		startKey := createPurgeIndexByTxidRangeStartKey(txid)
		endKey := createPurgeIndexByTxidRangeEndKey(txid)

		iter := s.db.GetIterator(startKey, endKey)

		// Get all txid and uuid from above result and remove it from transient store (both
		// write set and the corresponding indexes.
		for iter.Next() {
			// For each entry, remove the private read-write set and correponding indexes

			// Remove private write set
			compositeKeyPurgeIndexByTxid := iter.Key()
			// Note: We can create compositeKeyPvtRWSet by just replacing the prefix of compositeKeyPurgeIndexByTxid
			// with  prwsetPrefix. For code readability and to be expressive, we split and create again.
			uuid, blockHeight := splitCompositeKeyOfPurgeIndexByTxid(compositeKeyPurgeIndexByTxid)
			compositeKeyPvtRWSet := createCompositeKeyForPvtRWSet(txid, uuid, blockHeight)
			dbBatch.Delete(compositeKeyPvtRWSet)

			// Remove purge index -- purgeIndexByHeight
			compositeKeyPurgeIndexByHeight := createCompositeKeyForPurgeIndexByHeight(blockHeight, txid, uuid)
			dbBatch.Delete(compositeKeyPurgeIndexByHeight)

			// Remove purge index -- purgeIndexByTxid
			dbBatch.Delete(compositeKeyPurgeIndexByTxid)
		}
		iter.Release()
	}
	// If peer fails before/while writing the batch to golevelDB, these entries will be
	// removed as per BTL policy later by PurgeByHeight()
	return s.db.WriteBatch(dbBatch, true)
}

// PurgeByHeight removes private write sets at block height lesser than
// a given maxBlockNumToRetain. In other words, Purge only retains private write sets
// that were persisted at block height of maxBlockNumToRetain or higher. Though the private
// write sets stored in transient store is removed by coordinator using PurgebyTxids()
// after successful block commit, PurgeByHeight() is still required to remove orphan entries (as
// transaction that gets endorsed may not be submitted by the client for commit)
func (s *store) PurgeByHeight(maxBlockNumToRetain uint64) error {
	// Do a range query with 0 as startKey and maxBlockNumToRetain-1 as endKey
	startKey := createPurgeIndexByHeightRangeStartKey(0)
	endKey := createPurgeIndexByHeightRangeEndKey(maxBlockNumToRetain - 1)
	iter := s.db.GetIterator(startKey, endKey)

	dbBatch := leveldbhelper.NewUpdateBatch()

	// Get all txid and uuid from above result and remove it from transient store (both
	// write set and the corresponding index.
	for iter.Next() {
		// For each entry, remove the private read-write set and correponding indexes

		// Remove private write set
		compositeKeyPurgeIndexByHeight := iter.Key()
		txid, uuid, blockHeight := splitCompositeKeyOfPurgeIndexByHeight(compositeKeyPurgeIndexByHeight)
		compositeKeyPvtRWSet := createCompositeKeyForPvtRWSet(txid, uuid, blockHeight)
		dbBatch.Delete(compositeKeyPvtRWSet)

		// Remove purge index -- purgeIndexByTxid
		compositeKeyPurgeIndexByTxid := createCompositeKeyForPurgeIndexByTxid(txid, uuid, blockHeight)
		dbBatch.Delete(compositeKeyPurgeIndexByTxid)

		// Remove purge index -- purgeIndexByHeight
		dbBatch.Delete(compositeKeyPurgeIndexByHeight)
	}
	iter.Release()

	return s.db.WriteBatch(dbBatch, true)
}

// GetMinTransientBlkHt returns the lowest block height remaining in transient store
func (s *store) GetMinTransientBlkHt() (uint64, error) {
	// Current approach performs a range query on purgeIndex with startKey
	// as 0 (i.e., blockHeight) and returns the first key which denotes
	// the lowest block height remaining in transient store. An alternative approach
	// is to explicitly store the minBlockHeight in the transientStore.
	startKey := createPurgeIndexByHeightRangeStartKey(0)
	iter := s.db.GetIterator(startKey, nil)
	// Fetch the minimum transient block height
	if iter.Next() {
		dbKey := iter.Key()
		_, _, blockHeight := splitCompositeKeyOfPurgeIndexByHeight(dbKey)
		return blockHeight, nil
	}
	iter.Release()
	// Returning an error may not be the right thing to do here. May be
	// return a bool. -1 is not possible due to unsigned int as first
	// return value
	return 0, ErrStoreEmpty
}

func (s *store) Shutdown() {
	// do nothing because shared db is used
}

// Next moves the iterator to the next key/value pair.
// It returns whether the iterator is exhausted.
func (scanner *RwsetScanner) Next() (*EndorserPvtSimulationResults, error) {
	if !scanner.dbItr.Next() {
		return nil, nil
	}
	dbKey := scanner.dbItr.Key()
	dbVal := scanner.dbItr.Value()
	_, blockHeight := splitCompositeKeyOfPvtRWSet(dbKey)

	txPvtRWSet := &rwset.TxPvtReadWriteSet{}
	if err := proto.Unmarshal(dbVal, txPvtRWSet); err != nil {
		return nil, err
	}
	filteredTxPvtRWSet := pvtdatastorage.TrimPvtWSet(txPvtRWSet, scanner.filter)

	return &EndorserPvtSimulationResults{
		ReceivedAtBlockHeight: blockHeight,
		PvtSimulationResults:  filteredTxPvtRWSet,
	}, nil
}

// Close releases resource held by the iterator
func (scanner *RwsetScanner) Close() {
	scanner.dbItr.Release()
}
