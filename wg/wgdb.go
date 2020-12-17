package wg

import (
	"encoding/json"
	"log"
	"os"

	"github.com/tidwall/buntdb"
)

type wgdb struct {
	filepath string
	DB       *buntdb.DB
}

func (wdb *wgdb) openDB(dbpath string) error {
	InitIndexes := true
	if _, err := os.Stat(dbpath); os.IsNotExist(err) {
		InitIndexes = false
	}
	ddb, err := buntdb.Open(dbpath)
	if err != nil {
		return err
	}
	if InitIndexes {
		// ddb.CreateIndex("ClientIPCIDR", "*", buntdb.IndexString)
	}
	wdb.filepath = dbpath
	wdb.DB = ddb
	return nil
}

//GetClient TODO: fix this gay shit
func (wdb *wgdb) GetClient(clientIPCIDR string) (*WGClient, error) {
	var wgc WGClient
	var found bool
	err := wdb.DB.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("", func(key, value string) bool {
			jerr := json.Unmarshal([]byte(value), &wgc)
			if jerr != nil {
				log.Println(jerr)
				return false
			}
			if wgc.ClientIPCIDR == clientIPCIDR {
				found = true
				return false //value found, stop accending
			}
			return true
		})
		return err
	})
	if err != nil || !found {
		return nil, err
	}
	return &wgc, nil
}

func (wdb *wgdb) InsertUpdateClient(client *WGClient) error {
	//marsh the client
	jsondata, err := json.MarshalIndent(client, "", "  ")
	if err != nil {
		return err
	}
	err = wdb.DB.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(client.ClientIPCIDR, string(jsondata), nil)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}
func (wdb *wgdb) DeleteClient(client *WGClient) error {
	err := wdb.DB.View(func(tx *buntdb.Tx) error {

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
