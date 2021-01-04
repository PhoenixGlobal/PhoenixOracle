package models


import (
	"fmt"
	"github.com/asdine/storm"
	homedir "github.com/mitchellh/go-homedir"
	"log"
	"os"
	"path"
	"reflect"
)

type ORM struct {
	*storm.DB
}

func InitORM(env string) ORM {
	orm := ORM{initializeDatabase(env)}
	orm.migrate()
	return orm
}

func initializeDatabase(env string) *storm.DB {
	db, err := storm.Open(DBpath(env))
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func DBpath(env string) string {
	dir, err := homedir.Expand("~/.phoenix")
	if err != nil {
		log.Fatal(err)
	}

	os.MkdirAll(dir, os.FileMode(0700))
	var directory = path.Join(dir, "db."+env+".bolt")
	fmt.Println("directory  "+directory)
	return directory
}

func (self ORM) InitBucket(model interface{}) error {
	return self.Init(model)
}


func (self ORM) Where(field string, value interface{}, instance interface{}) error {
	err:= self.Find(field, value, instance)
	if err == storm.ErrNotFound{
		emptySlice(instance)
		return nil
	}
	return err
}

func emptySlice(to interface{}) {
	ref := reflect.ValueOf(to)
	results := reflect.MakeSlice(reflect.Indirect(ref).Type(),0,0)
	reflect.Indirect(ref).Set(results)
}


func (self ORM) JobsWithCron() ([]Job, error) {
	jobs := []Job{}
	err := self.AllByIndex("Cron", &jobs)
	fmt.Println(jobs)
	return jobs, err
}


