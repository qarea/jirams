package store

import (
	"testing"

	"github.com/boltdb/bolt"
	"github.com/qarea/jirams/entities"
	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	db, err := bolt.Open("var/bolt/store.db", 0666, nil)
	if err != nil {
		t.FailNow()
	}
	defer db.Close()
	store := New(db)

	// no key in empty DB
	key, err := store.GetKey(1, 1)
	assert.Nil(t, err)
	assert.Equal(t, "", key)

	// add key, get id
	id, err := store.GetID(1, "TEST")
	assert.Nil(t, err)
	assert.Equal(t, entities.UserID(1), id)

	// add another key, add incremented id
	id, err = store.GetID(1, "TEST2")
	assert.Nil(t, err)
	assert.Equal(t, entities.UserID(2), id)

	// read key back
	key, err = store.GetKey(1, 1)
	assert.Nil(t, err)
	assert.Equal(t, "TEST", key)

	// read other key back
	key, err = store.GetKey(1, 2)
	assert.Nil(t, err)
	assert.Equal(t, "TEST2", key)

	// try reading from other tracker
	key, err = store.GetKey(2, 1)
	assert.Nil(t, err)
	assert.Equal(t, "", key)
}
