package db

import (
	"log"
	"time"
)

type SlowDB struct {
	data map[string]string
}

func (db *SlowDB) Get(key string) string {
	time.Sleep(time.Duration(300) * time.Millisecond)
	value := db.data[key]
	log.Printf("getting key:%s value:%s\n", key, value)
	return value
}

func (db *SlowDB) Set(key string, value string) {
	log.Printf("setting %s to %s\n", key, value)
	db.data[key] = value
}

func (db *SlowDB) Del(key string) {
	log.Printf("del %s\n", key)
	delete(db.data, key)
}

func (db *SlowDB) Data() map[string]string{
	log.Printf("data %+v\n", db.data)
	return db.data
}

func NewSlowDB() *SlowDB {
	ndb := new(SlowDB)
	ndb.data = make(map[string]string)
	return ndb
}
