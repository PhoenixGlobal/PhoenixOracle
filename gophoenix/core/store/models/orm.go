package models

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"log"
	"path"
	"reflect"
)

type ORM struct {
	*storm.DB
}

func NewORM(dir string) *ORM {
	path := path.Join(dir, "db.bolt")
	orm := &ORM{initializeDatabase(path)}
	orm.migrate()
	return orm
}

func initializeDatabase(path string) *storm.DB {
	db, err := storm.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func (self *ORM) InitBucket(model interface{}) error {
	return self.Init(model)
}


func (self *ORM) Where(field string, value interface{}, instance interface{}) error {
	err:= self.Find(field, value, instance)
	if err == storm.ErrNotFound{
		emptySlice(instance)
		return nil
	}
	return err
}
func (self *ORM) JobRunsFor(job Job) ([]JobRun, error) {
	var runs []JobRun
	err := self.Where("JobID", job.ID, &runs)
	return runs, err
}

func emptySlice(to interface{}) {
	ref := reflect.ValueOf(to)
	results := reflect.MakeSlice(reflect.Indirect(ref).Type(),0,0)
	reflect.Indirect(ref).Set(results)
}


func (self *ORM) JobsWithCron() ([]Job, error) {
	initrs := []Initiator{}
	self.Where("Type", "cron", &initrs)
	jobIDs := []string{}
	for _, initr := range initrs {
		jobIDs = append(jobIDs, initr.JobID)
	}
	jobs := []Job{}
	err := self.Select(q.In("ID",jobIDs)).Find(&jobs)
	if err == storm.ErrNotFound {
		return jobs, nil
	}
	return jobs, err
}

func (self *ORM) SaveJob(job Job) error {
	tx, err := self.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := tx.Save(&job); err != nil {
		return err
	}
	for _, initr := range job.Initiators {
		initr.JobID = job.ID
		if err := tx.Save(&initr); err != nil {
			return err
		}
	}
	return tx.Commit()
}
