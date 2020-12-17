package wg

import (
	"fmt"
	"os"

	"github.com/tidwall/buntdb"
)

type wgdb struct {
	filepath string
	DB       *buntdb.DB
}

func (db *wgdb) openDB(dbpath string) error {
	InitIndexes := true
	if _, err := os.Stat(dbpath); os.IsNotExist(err) {
		InitIndexes = false
	}
	ddb, err := buntdb.Open(dbpath)
	if err != nil {
		return err
	}
	if InitIndexes {
		db.CreateIndex("ages", "user:*:age", buntdb.IndexInt)
	}
	db.filepath = dbpath
	db.DB = ddb
	return nil
}
func (db *wgdb) GetClient(key string) (string, error) {
	err := db.View(func(tx *buntdb.Tx) error {
		tx.Ascend("ages", func(key, val string) bool {
			fmt.Printf(buf, "%s %s\n", key, val)
			return true
		})
		return nil
	})
	return nil
}

func (db *wgdb) UpdateClient(key string) (string, error) {
	err := db.View(func(tx *buntdb.Tx) error {

		return nil
	})
	return nil
}
func (db *wgdb) InsertClient(key string) (string, error) {
	err := db.View(func(tx *buntdb.Tx) error {

		return nil
	})
	return nil
}
func (db *wgdb) DeleteClient(key string) (string, error) {
	err := db.View(func(tx *buntdb.Tx) error {

		return nil
	})
	return nil
}
