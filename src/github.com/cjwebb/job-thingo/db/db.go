package db

import (
	"log"
	"time"
	"encoding/json"

	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/dynamodb"
)

type Database struct {
	server *dynamodb.Server
}

type Job struct {
	Id            string
	Title         string
	Description   string
	ContactEmail  string
	Rate          string
	JobConsList   []JobRef // persisted as a json string 
}

type JobRef struct {
	Id	string `json:"id"`
	Email	string `json:"email"`
}

func NewDatabase() Database {
	// This assumes you have ENV vars:
	// AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY
	auth, err := aws.EnvAuth()
	if err != nil {
		log.Fatal(err.Error())
	}
	ddbs := dynamodb.New(auth, aws.EUWest)

	return Database{ddbs}
}

func (db Database) GetJob(id string) (Job, error) {
	primaryKey := dynamodb.PrimaryKey{dynamodb.NewStringAttribute("id", ""), nil}
	t := dynamodb.Table{db.server, "jobthing-jobs", primaryKey}
	r, err := t.GetItem(&dynamodb.Key{id, ""})

	if err != nil {
		return Job{}, err // todo(cjwebb) - better error handling
	}

	job := Job{
		get("id", r),
		get("title", r),
		get("description", r),
		get("contact_email", r),
		get("rate", r),
		stringToSlice(get("job_cons_list", r)),
	}
	return job, nil
}

func get(name string, m map[string]*dynamodb.Attribute) string {
	if val, ok := m[name]; ok {
		return val.Value
	} else {
		return ""
	}
}

func stringToSlice(s string) []JobRef {
	var dat []JobRef
	err := json.Unmarshal([]byte(s), &dat)
	if err != nil {
		return nil
	}
	return dat
}

func (db Database) PutJob(job Job) error {
	primaryKey := dynamodb.PrimaryKey{dynamodb.NewStringAttribute("id", ""), nil}
	t := dynamodb.Table{db.server, "jobthing-jobs", primaryKey}
	_, err := t.PutItem(job.Id, "", []dynamodb.Attribute{
		*dynamodb.NewStringAttribute("title", job.Title),
		*dynamodb.NewStringAttribute("description", job.Description),
		*dynamodb.NewStringAttribute("contact_email", job.ContactEmail),
		*dynamodb.NewStringAttribute("rate", job.Rate),
		*dynamodb.NewStringAttribute("job_cons_list", sliceToString(job.JobConsList)),
		*dynamodb.NewStringAttribute("date_created", time.Now().Format(time.RFC3339)),
	})
	return err
}

func sliceToString(arr []JobRef) string {
	s, err := json.Marshal(arr)
	if err != nil {
		return ""
	}
	return string(s)
}

