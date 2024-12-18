package db_mdbx

import (
	"encoding/hex"
	"errors"
	"github.com/RevittConsulting/mdbx-viewer/pkg/utils"
	"github.com/RevittConsulting/mdbx-viewer/types"
	"github.com/erigontech/mdbx-go/mdbx"
	"log"
)

type MDBX struct {
	env *mdbx.Env
}

func New() *MDBX {
	env, err := mdbx.NewEnv()
	if err != nil {
		log.Fatal(err)
	}

	err = env.SetOption(mdbx.OptMaxDB, 200)
	if err != nil {
		log.Fatal(err)
	}

	return &MDBX{
		env: env,
	}
}

// Open - Open the environment.
func (m *MDBX) Open(path string) error {
	return m.env.Open(path, mdbx.NoTLS|mdbx.Readonly, 0444)
}

// Close - Close the environment.
func (m *MDBX) Close() {
	m.env.Close()
}

// List - Returns all dbi names.
func (m *MDBX) List() ([]string, error) {
	txn, err := m.env.BeginTxn(nil, mdbx.Readonly)
	if err != nil {
		return nil, err
	}
	defer txn.Abort()

	return txn.ListDBI()
}

// Entries - Number of data items.
func (m *MDBX) Entries(name string) (uint64, error) {
	txn, err := m.env.BeginTxn(nil, mdbx.Readonly)
	if err != nil {
		return 0, err
	}
	defer txn.Abort()

	dbi, err := txn.OpenDBISimple(name, 0)
	if err != nil {
		return 0, err
	}

	dbiStat, err := txn.StatDBI(dbi)
	if err != nil {
		return 0, err
	}

	return dbiStat.Entries, nil
}

// ValuesByKey - Returns all values with the given key.
func (m *MDBX) ValuesByKey(name string, key []byte) ([][]byte, error) {
	txn, err := m.env.BeginTxn(nil, mdbx.Readonly)
	if err != nil {
		return nil, err
	}
	defer txn.Abort()

	dbi, err := txn.OpenDBISimple(name, 0)
	if err != nil {
		return nil, err
	}
	defer m.env.CloseDBI(dbi)

	cursor, err := txn.OpenCursor(dbi)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var values [][]byte
	for k, v, err := cursor.Get(nil, nil, mdbx.First); err == nil; k, v, err = cursor.Get(nil, nil, mdbx.Next) {
		if utils.BytesEqual(k, key) {
			values = append(values, v)
		}
	}

	return values, nil
}

// KeysByValue - Returns all keys with the given value.
func (m *MDBX) KeysByValue(name string, value []byte) ([][]byte, error) {
	txn, err := m.env.BeginTxn(nil, mdbx.Readonly)
	if err != nil {
		return nil, err
	}
	defer txn.Abort()

	dbi, err := txn.OpenDBISimple(name, 0)
	if err != nil {
		return nil, err
	}
	defer m.env.CloseDBI(dbi)

	cursor, err := txn.OpenCursor(dbi)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	var keys [][]byte
	for k, v, err := cursor.Get(nil, nil, mdbx.First); err == nil; k, v, err = cursor.Get(nil, nil, mdbx.Next) {
		if utils.BytesEqual(v, value) {
			keys = append(keys, k)
		}
	}

	return keys, nil
}

// Read - Read data from the database.
func (m *MDBX) Read(name string, take, offset uint64) ([]types.KeyValuePair, error) {
	var data []types.KeyValuePair

	txn, err := m.env.BeginTxn(nil, mdbx.Readonly)
	if err != nil {
		return nil, err
	}
	defer txn.Abort()

	dbi, err := txn.OpenDBISimple(name, 0)
	if err != nil {
		return nil, err
	}
	defer m.env.CloseDBI(dbi)

	cursor, err := txn.OpenCursor(dbi)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()

	k, v, err := cursor.Get(nil, nil, mdbx.First)
	if err != nil {
		return nil, err
	}

	keyCount, err := m.Entries(name)
	if err != nil {
		return nil, err
	}

	if take > keyCount {
		take = keyCount
	}

	count := 0
	for ; err == nil; k, v, err = cursor.Get(nil, nil, mdbx.Next) {
		if uint64(count) >= offset && uint64(count) < (offset+take) {
			data = append(data, types.KeyValuePair{
				Key:   k,
				Value: v,
			})
		}
		count++
		if uint64(count) >= (offset + take) {
			break
		}
	}

	if err != nil && !errors.Is(err, mdbx.NotFound) {
		return nil, err
	}

	return data, nil
}

// CountKeysOfLength - Count and keys of a given length.
func (m *MDBX) CountKeysOfLength(name string, length uint64) (uint64, []string, error) {
	var count uint64
	var keys []string

	txn, err := m.env.BeginTxn(nil, mdbx.Readonly)
	if err != nil {
		return 0, nil, err
	}
	defer txn.Abort()

	dbi, err := txn.OpenDBISimple(name, 0)
	if err != nil {
		return 0, nil, err
	}
	defer m.env.CloseDBI(dbi)

	cursor, err := txn.OpenCursor(dbi)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close()

	limit := uint64(1000)
	for k, _, err := cursor.Get(nil, nil, mdbx.First); err == nil; k, _, err = cursor.Get(nil, nil, mdbx.Next) {
		if uint64(len(k)*2) == length {
			count++
			if count <= limit {
				keyHex := hex.EncodeToString(k)
				keys = append(keys, keyHex)
			}
		}
	}

	return count, keys, nil
}

// CountValuesOfLength - Count and values of a given length.
func (m *MDBX) CountValuesOfLength(bucketName string, length uint64) (uint64, []string, error) {
	var count uint64
	var values []string

	txn, err := m.env.BeginTxn(nil, mdbx.Readonly)
	if err != nil {
		return 0, nil, err
	}
	defer txn.Abort()

	dbi, err := txn.OpenDBISimple(bucketName, 0)
	if err != nil {
		return 0, nil, err
	}
	defer m.env.CloseDBI(dbi)

	cursor, err := txn.OpenCursor(dbi)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close()

	limit := uint64(1000)
	for _, v, err := cursor.Get(nil, nil, mdbx.First); err == nil; _, v, err = cursor.Get(nil, nil, mdbx.Next) {
		if uint64(len(v)*2) == length {
			count++
			if count <= limit {
				keyHex := hex.EncodeToString(v)
				values = append(values, keyHex)
			}
		}
	}

	return count, values, nil
}
