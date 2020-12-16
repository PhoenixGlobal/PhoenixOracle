package orm


import (
	"github.com/asdine/storm"
	homedir "github.com/mitchellh/go-homedir"
	"log"
	"os"
	"path"
)

type phoenixDB interface {
	GetDB() *storm.DB
	Close()
}

type TestDB struct {
	*storm.DB
}

type persistentDB struct {
	*storm.DB
}

var db phoenixDB

func Init() {
	db = persistentDB{initializeDatabase("production")}
	migrate()
}

func InitTest() {
	db = TestDB{initializeDatabase("test")}
	migrate()
}

func GetDB() *storm.DB {
	return db.GetDB()
}

func Close() {
	db.Close()
}

func (d TestDB) GetDB() *storm.DB {
	return d.DB
}

func (d persistentDB) GetDB() *storm.DB {
	return d.DB
}

func (d TestDB) Close() {
	d.GetDB().Close()
	os.Remove(dbpath("test"))
}

func (d persistentDB) Close() {
	d.GetDB().Close()
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
	return path.Join(dir, "db."+env+".bolt")
}
