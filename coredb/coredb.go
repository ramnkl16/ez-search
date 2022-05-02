package coredb

import (
	"fmt"
	"path"

	"github.com/boltdb/bolt"
	"go.uber.org/zap/zapcore"

	"github.com/ramnkl16/ez-search/global"
	"github.com/ramnkl16/ez-search/utils/cache_utils"
)

const (
	coreDbName string = "core.db"
)

func setupDB() (*bolt.DB, error) {
	d, err := bolt.Open(path.Join(global.WorkingDir, coreDbName), 0600, nil)

	if err != nil {
		//logger.Error("could not open db", err)
		return nil, err
	}

	err = d.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("DB"))
		if err != nil {
			return fmt.Errorf("could not create root bucket: %v", err)
		}
		_, err = root.CreateBucketIfNotExists([]byte("catalog"))
		if err != nil {
			return fmt.Errorf("could not create %v bucket: %v", "catalog", err)
		}
		return nil
	})
	if err != nil {
		//logger.Error("could not set up buckets", err)
		return nil, err
	}
	fmt.Println("DB Setup is Done")
	return d, nil
}
func GetDb() *bolt.DB {
	var db *bolt.DB
	var err error
	//fmt.Println("Get db called", dbName, "Count ", cache_utils.Cache.Count())
	i, err := cache_utils.Cache.Get(coreDbName)
	//fmt.Println("after cache get")
	if err != nil {
		fmt.Println("getDb|Failed while cache_utils.cache", err, zapcore.Field{String: coreDbName, Key: "p1", Type: zapcore.StringType})
		i, err = setupDB()

		if err != nil {
			fmt.Sprintln("GetboltDb|Failed while open index", err, zapcore.Field{String: coreDbName, Key: "p1", Type: zapcore.StringType})
			return nil
		}
		//fmt.Println("GetboltDb|db has found under data folder", dbName, err)
		cache_utils.AddOrUpdateCache(coreDbName, i)
	}

	db = i.(*bolt.DB)
	//fmt.Println("before return ", db.Path())
	return db
}

func AddKey(key string, value []byte) error {
	db := GetDb()
	//fmt.Println("addkey", db.Path())
	err := db.Update(func(tx *bolt.Tx) error {
		//logger.Info(fmt.Sprintf("addkey:%s", key))
		err := tx.Bucket([]byte("DB")).Bucket([]byte("catalog")).Put([]byte(key), value)
		if err != nil {
			fmt.Println("Addkey|coredb|could not insert entry", err)
			return fmt.Errorf("could not insert entry: %v", err)
		}
		return nil
	})
	if err != nil {
		//logger.Error(fmt.Sprintf("Failed while add bolt db key:%s", key), err)
		//logger.Error("", err)
		fmt.Println(err)
		return err
	}
	//fmt.Println("Added Entry")
	//fmt.Println("AddKey|dbInfo", db.Info())
	return nil
}
func List() {
	db := GetDb()
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("DB")).Bucket([]byte("catalog")).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%d|%s\n", k, len(string(v)), v)
		}
		//b, _ := GetKey(db, "NA_B_916854_en-US")
		//fmt.Printf("key=NA_B_916854_en-US v=%s", b)
		return nil
	})
}
func GetKey(key string) (result []byte, err error) {
	db := GetDb()
	//fmt.Println("getkey|dbInfo", db.Info())
	db.View(func(tx *bolt.Tx) error {
		r := tx.Bucket([]byte("DB")).Bucket([]byte("catalog")).Get([]byte(key))
		//logger.Info(fmt.Sprintf("getkey [boltdb]:%s|%s", key, string(r)))
		if r != nil {
			result = make([]byte, len(r))
			copy(result, r)
		}
		return nil
	})
	return
}
func Delete(key string) error {
	db := GetDb()
	//fmt.Println("getkey|dbInfo", db.Info())
	db.View(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte("DB")).Bucket([]byte("catalog")).Delete([]byte(key))
		//logger.Info(fmt.Sprintf("getkey [boltdb]:%s|%s", key, string(r)))
		if err != nil {
			//logger.Error("Failed while delete key ", err)
			fmt.Println(err)
			return err
		}
		return nil
	})
	return nil
}
