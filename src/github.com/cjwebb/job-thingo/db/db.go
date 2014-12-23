package db

import (
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/dynamodb"
	"log"
)

type Database struct {
	server *dynamodb.Server
}

type Job struct {
	uuid          string
	Title         string
	Link          string
	Description   string
	Contact_email string
	Rate          string
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

func (s Database) GetJob(id string) (Job, error) {
	primaryKey := dynamodb.PrimaryKey{dynamodb.NewStringAttribute("uuid", ""), nil}
	t := dynamodb.Table{s.server, "jt-jobs", primaryKey}
	r, err := t.GetItem(&dynamodb.Key{id, ""})

	if err != nil {
		return Job{}, err // todo(cjwebb) - better error handling
	}

	job := Job{
		get("uuid", r),
		get("title", r),
		get("link", r),
		get("description", r),
		get("contact_email", r),
		get("rate", r),
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

