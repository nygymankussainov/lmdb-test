package main

import (
	"fmt"
	"os"

	"github.com/bmatsuo/lmdb-go/lmdb"
)

func main() {
	path := "./db-test"
	name := "DBI"
	env, err := lmdb.NewEnv()
	if err != nil {
		panic(err)
	}
	if err = env.SetMaxDBs(1); err != nil {
		panic(err)
	}
	if err = env.SetMapSize(1 << 30); err != nil {
		panic(err)
	}
	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err = os.Mkdir(path, 0755); err != nil {
			panic(err)
		}
	}
	if err = env.Open(path, 0, 0644); err != nil {
		panic(err)
	}
	err = env.Update(func(txn *lmdb.Txn) error {
		_, err := txn.OpenDBI(name, lmdb.Create)
		return err
	})
	if err != nil {
		panic(err)
	}
	err = env.Update(func(txn *lmdb.Txn) error {
		dbi, err := txn.OpenDBI(name, 0)
		if err != nil {
			return err
		}
		err = txn.Put(dbi, []byte("ping"), []byte("pong"), 0)
		return err
	})
	if err != nil {
		panic(err)
	}
	err = env.View(func(txn *lmdb.Txn) error {
		dbi, err := txn.OpenDBI(name, 0)
		if err != nil {
			return err
		}
		val, err := txn.Get(dbi, []byte("ping"))
		fmt.Println(string(val))
		return err
	})
	if err != nil {
		panic(err)
	}
}
