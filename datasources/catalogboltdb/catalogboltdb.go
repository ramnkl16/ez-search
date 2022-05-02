package catalogboltdb

import (
	"fmt"
	"path"

	"github.com/boltdb/bolt"

	"github.com/ramnkl16/ez-search/global"
	"github.com/ramnkl16/ez-search/logger"
	"github.com/ramnkl16/ez-search/utils/cache_utils"
)

var (
	//NaDb   *bolt.DB
	JsonDb *bolt.DB
)

// func InitializeBoltDB() {
// 	logger.Info("Initalize bolt db")
// 	logger.Info(config.NaCatalogDbName)
// 	var err error
// 	// NaDb, err = setupDB(config.NaCatalogDbName, config.BuketName)
// 	// if err != nil {
// 	// 	logger.Error(fmt.Sprintf("Failed bolt na db initalize %s", config.NaCatalogDbName), err)
// 	// }
// 	JsonDb, err = setupDB(config.EanzCatalogDbName, config.BuketName)
// 	if err != nil {
// 		logger.Error(fmt.Sprintf("Failed bolt na db initalize %s", config.EanzCatalogDbName), err)
// 	}
// }

func setupDB(dbName string, bucket string) (*bolt.DB, error) {
	d, err := bolt.Open(path.Join(global.WorkingDir, dbName), 0600, nil)

	if err != nil {
		logger.Error("could not open db", err)
		return nil, err
	}

	err = d.Update(func(tx *bolt.Tx) error {
		root, err := tx.CreateBucketIfNotExists([]byte("DB"))
		if err != nil {
			return fmt.Errorf("could not create root bucket: %v", err)
		}
		_, err = root.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return fmt.Errorf("could not create %v bucket: %v", bucket, err)
		}
		return nil
	})
	if err != nil {
		logger.Error("could not set up buckets", err)
		return nil, err
	}
	fmt.Println("DB Setup is Done")
	return d, nil
}

func AddKey(dbName, key string, value []byte) error {
	db := GetboltDb(dbName)
	//fmt.Println("addkey", db.Path())
	err := db.Update(func(tx *bolt.Tx) error {
		//logger.Info(fmt.Sprintf("addkey:%s", key))
		err := tx.Bucket([]byte("DB")).Bucket([]byte(config.BuketName)).Put([]byte(key), value)
		if err != nil {
			return fmt.Errorf("could not insert entry: %v", err)
		}
		return nil
	})
	if err != nil {
		logger.Error(fmt.Sprintf("Failed while add bolt db key:%s", key), err)
		return err
	}
	//fmt.Println("Added Entry")
	//fmt.Println("AddKey|dbInfo", db.Info())
	return nil
}
func List(dbName string, bucket string) {
	db := GetboltDb(dbName)
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("DB")).Bucket([]byte(config.BuketName)).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%d|%s\n", k, len(string(v)), v)
		}
		//b, _ := GetKey(db, "NA_B_916854_en-US")
		//fmt.Printf("key=NA_B_916854_en-US v=%s", b)
		return nil
	})
}
func GetKey(dbName string, key string) (result []byte, err error) {
	db := GetboltDb(dbName)
	//fmt.Println("getkey|dbInfo", db.Info())
	db.View(func(tx *bolt.Tx) error {
		r := tx.Bucket([]byte("DB")).Bucket([]byte(config.BuketName)).Get([]byte(key))
		//logger.Info(fmt.Sprintf("getkey [boltdb]:%s|%s", key, string(r)))
		if r != nil {
			result = make([]byte, len(r))
			copy(result, r)
		}
		return nil
	})
	return
}
func Delete(dbName string, key string) error {
	db := GetboltDb(dbName)
	//fmt.Println("getkey|dbInfo", db.Info())
	db.View(func(tx *bolt.Tx) error {
		err := tx.Bucket([]byte("DB")).Bucket([]byte(config.BuketName)).Delete([]byte(key))
		//logger.Info(fmt.Sprintf("getkey [boltdb]:%s|%s", key, string(r)))
		if err != nil {
			fmt.Printf("Failed while delete key %s", err.Error())
			return err
		}
		return nil
	})
	return nil
}

func GetboltDb(dbName string) *bolt.DB {
	var db *bolt.DB
	var err error
	//fmt.Println("Get db called", dbName, "Count ", cache_utils.Cache.Count())
	i, err := cache_utils.Cache.Get(dbName)
	//fmt.Println("after cache get")
	if err != nil {
		fmt.Println("GetboltDb|Failed while cache_utils.cache", dbName, err)
		i, err = setupDB(dbName, "catalog")

		if err != nil {
			fmt.Println("GetboltDb|Failed while open index", dbName, err)
			return nil
		}
		//fmt.Println("GetboltDb|db has found under data folder", dbName, err)
		cache_utils.AddOrUpdateCache(dbName, i)
	}

	db = i.(*bolt.DB)
	//fmt.Println("before return ", db.Path())
	return db
}
