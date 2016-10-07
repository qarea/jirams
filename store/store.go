package store

import (
	"encoding/binary"

	"github.com/boltdb/bolt"
	"github.com/powerman/narada-go/narada"
	"github.com/qarea/jirams/entities"
)

const (
	userBucket    = "Users"
	userKeyBucket = "UserKeys"
)

// UserKeyMapper interface defines a storage for user key string to user id number mapping
type UserKeyMapper interface {
	Init()
	GetID(trackerID entities.TrackerID, key entities.UserKey) (res entities.UserID, err error)
	GetKey(trackerID entities.TrackerID, userID entities.UserID) (res string, err error)
}

// Store implements BoldDB backed UserKeyMapper
type Store struct {
	DB *bolt.DB
}

// New creates an instance of Store
func New(db *bolt.DB) *Store {
	res := Store{DB: db}
	res.Init()
	return &res
}

var log = narada.NewLog("user store: ")

// Init initializes BoltDB storage if needed
func (store *Store) Init() {
	_ = store.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(userBucket))
		if err != nil {
			log.Fatal("Failed to create Jira Users store")
		}
		_, err = tx.CreateBucketIfNotExists([]byte(userKeyBucket))
		if err != nil {
			log.Fatal("Failed to create Jira Users store")
		}
		return nil
	})
}

// GetID looks for provided user key in the mapping, stores it if it is not present and returns associated numeric user id
func (store *Store) GetID(trackerID entities.TrackerID, key entities.UserKey) (res entities.UserID, err error) {

	err = store.DB.Update(func(tx *bolt.Tx) error {
		userKeys := tx.Bucket([]byte(userKeyBucket))
		found := userKeys.Get([]byte(key))
		// if key already known, return it's ID
		if found != nil {
			rawRes := userKeys.Get([]byte(key))
			if rawRes != nil {
				res = entities.UserID(btoi(rawRes))
			}
			return nil
		}

		// otherwise generate new ID and try to adding key to store
		users := tx.Bucket([]byte(userBucket))
		id, _ := users.NextSequence()
		storeKey := makeKey(trackerID, entities.UserID(id))

		if err = users.Put(storeKey, []byte(key)); err != nil {
			return err
		}

		if err = userKeys.Put([]byte(key), itob(id)); err != nil {
			return err
		}

		res = entities.UserID(id)
		return nil
	})
	return
}

// GetKey looks for provided user ID in the mapping and returns original user key
func (store *Store) GetKey(trackerID entities.TrackerID, userID entities.UserID) (res string, err error) {
	err = store.DB.View(func(tx *bolt.Tx) error {
		users := tx.Bucket([]byte(userBucket))
		storeKey := makeKey(trackerID, userID)
		rawRes := users.Get(storeKey)
		if rawRes != nil && len(rawRes) > 0 {
			res = string(rawRes)
		}
		return nil
	})
	return
}

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func btoi(v []byte) uint64 {
	return binary.BigEndian.Uint64(v)
}

func makeKey(trackerID entities.TrackerID, userID entities.UserID) []byte {
	return append(itob(uint64(trackerID)), itob(uint64(userID))...)
}
