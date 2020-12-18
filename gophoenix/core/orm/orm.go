package orm


import (
	"fmt"
	"github.com/asdine/storm"
	homedir "github.com/mitchellh/go-homedir"
	"log"
	"os"
	"path"
)

var db *storm.DB

func Init() {
	db = initializeDatabase("production")
	migrate()
}

func InitTest() {
	os.Remove(dbpath("test"))
	db = initializeDatabase("test")
	migrate()
}

func GetDB() *storm.DB {
	return db
}

func Close() {
	db.Close()
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

func Save(data interface{}) error {
	return db.Save(data)
}
