job-thingo
==========
Website for JobThing, which enables tracking/sharing of job adverts via email addresses and job uuids.

## To Run:
This project is written using [Golang](https://golang.org/doc/code.html) (initially v1.4) and uses [go-martini](https://github.com/go-martini/martini) as a web framework.

The database used is [DynamoDB](http://aws.amazon.com/dynamodb/), and as such, you will need to set AWS environmental variables prior to running the project.

```
export AWS_ACCESS_KEY_ID=YOUR_KEY
export AWS_SECRET_ACCESS_KEY=YOUR_SECRET_ACCESS_KEY
```

This key/secret pair will need access to DynamoDB. At present, you may also need to manually create a table with DynamoDB named `jobthing-jobs`, with a string hash key of `id`.

Once these keys have been set, and you have installed Golang, simply run:

```
go run main.go
```

to start a the project running as a development server.

