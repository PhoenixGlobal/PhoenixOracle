package models


import (
	"fmt"
	"github.com/asdine/storm"
	homedir "github.com/mitchellh/go-homedir"
	"log"
	"os"
	"path"
)

var db *storm.DB

func InitDB() {
	db = initializeDatabase("production")
	migrate()
}

func InitDBTest() {
	os.Remove(dbpath("test"))
	db = initializeDatabase("test")
	migrate()
}

func getDB() *storm.DB {
	return db
}

func Save(value interface{}) error {
	return getDB().Save(value)
}

func CloseDB() {
	getDB().Close()
}

func initializeDatabase(env string) *storm.DB {
	db, err := storm.Open(dbpath(env))
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func dbpath(env string) string {
	dir, err := homedir.Expand("~/.phoenix")
	if err != nil {
		log.Fatal(err)
	}

	os.MkdirAll(dir, os.FileMode(0700))
	var directory = path.Join(dir, "db."+env+".bolt")
	fmt.Println("directory  "+directory)
	return directory
}

func Find(field string, value interface{}, instance interface{}) error {
	return db.One(field, value, instance)
}


func Where(field string, value interface{}, instance interface{}) error {
	return db.Find(field, value, instance)
}
