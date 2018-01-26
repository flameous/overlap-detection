package overlap

import (
	"github.com/boltdb/bolt"
	"os"
	"log"
	"encoding/json"
)

type DB struct {
	*bolt.DB
}

var (
	proceededFiles = []byte("proceeded_files")
	overlapsData   = []byte("overlaps_data")
	exists         = []byte{'1'}
)

func NewDB(path string) *DB {
	db, err := bolt.Open(path, os.ModePerm, nil)
	if err != nil {
		log.Panicf("ошибка при создании файла БД: %v", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(proceededFiles)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(overlapsData)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Panicf("ошибка при инициализации БД: %v", err)
	}
	return &DB{db}
}

func (d *DB) IsFileProcessed(filename string) bool {
	ok := false
	d.View(func(tx *bolt.Tx) error {
		ok = tx.Bucket(proceededFiles).Get([]byte(filename)) != nil
		return nil
	})
	return ok
}
func (d *DB) SaveFileData(filename string, overlaps Overlaps) error {
	return d.Update(func(tx *bolt.Tx) error {
		err := tx.Bucket(proceededFiles).Put([]byte(filename), exists)
		if err != nil {
			return err
		}
		b, err := overlaps.dump()
		if err != nil {
			return err
		}
		return tx.Bucket(overlapsData).Put([]byte(filename), b)
	})
}

func (d *DB) GetProcessedFileList() ([]byte, error) {
	var data [][]byte
	var err error

	err = d.View(func(tx *bolt.Tx) error {
		return tx.Bucket(proceededFiles).ForEach(func(k, v []byte) error {
			data = append(data, k)
			return nil
		})
	})

	stringSlice := make([]string, 0, len(data))
	for _, v := range data {
		stringSlice = append(stringSlice, string(v))
	}
	b, err := json.MarshalIndent(stringSlice, "", "    ")
	return b, err
}

func (d *DB) GetFileData(filename string) ([]byte, error) {
	var b []byte
	err := d.View(func(tx *bolt.Tx) error {
		b = tx.Bucket(overlapsData).Get([]byte(filename))
		return nil
	})
	return b, err
}
