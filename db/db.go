package db

import (
	"github.com/colinrs/pkgx/logger"
)

type DB struct {
	data map[string]string
}

func (db *DB) Get(key string) string {
	value := db.data[key]
	logger.Info("getting key:%s value:%s", key, value)
	return value
}

func (db *DB) Set(key string, value string) {
	logger.Info("setting %s to %s", key, value)
	db.data[key] = value
}

func (db *DB) Del(key string) {
	logger.Info("del %s", key)
	delete(db.data, key)
}

func (db *DB) Data() map[string]string{
	logger.Info("data %+v", db.data)
	return db.data
}

func NewDB() *DB {
	ndb := new(DB)
	ndb.data = make(map[string]string)
	return ndb
}
