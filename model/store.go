package model

import (
	"time"

	bolt "go.etcd.io/bbolt"
)

var DbNow *DB

type DB struct {
	db *bolt.DB
}

func (d *DB) Expire(key string, seconds int) error {
	timer := time.NewTimer(time.Duration(seconds) * time.Second)
	go func() {
		<-timer.C
		d.Del(key)
	}()
	return nil
}

var _ Store = (*DB)(nil)
var tableName = []byte("sms")

type Store interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	Del(key string) error
	Expire(key string, seconds int) error
	Len() (int, error)
	Close() error
}

func InitDb(dbName string) {
	DbNow = &DB{
		db: openDataBase(dbName),
	}
}

func openDataBase(dbName string) *bolt.DB {
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		return nil
	}

	tables := [...][]byte{
		tableName,
	}

	db.Update(func(tx *bolt.Tx) error {
		for _, table := range tables {
			_, err2 := tx.CreateBucketIfNotExists(table)
			if err2 != nil {
				return err2
			}
		}
		return nil
	})
	return db
}

func (d *DB) Del(key string) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(tableName)
		return b.Delete([]byte(key))
	})
}

func (d *DB) Set(key string, value []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(tableName)
		return b.Put([]byte(key), value)
	})
}

func (d *DB) Get(key string) ([]byte, error) {
	var value []byte
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(tableName)
		value = b.Get([]byte(key))
		return nil
	})
	return value, err
}

func (d *DB) Len() (int, error) {
	var count int
	err := d.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(tableName)
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			count++
		}
		return nil
	})
	return count, err
}
func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) Flush() {
	d.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(tableName)
	})
	d.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(tableName)
		return err
	})
}
