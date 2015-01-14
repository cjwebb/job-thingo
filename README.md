job-thingo
==========
Website for JobThing, which enables tracking/sharing of job adverts via email addresses and job uuids.

## To Run (in development):
This project is written using [Golang](https://golang.org/doc/code.html) (initially v1.4) and uses [go-martini](https://github.com/go-martini/martini) as a web framework.

The database used is [DynamoDB](http://aws.amazon.com/dynamodb/), and as such, you will need to set AWS environmental variables prior to running the project.

```
export AWS_ACCESS_KEY_ID={YOUR_KEY}
export AWS_SECRET_ACCESS_KEY={YOUR_SECRET_ACCESS_KEY}
```

This key/secret pair will need access to DynamoDB. At present, you may also need to manually create a table with DynamoDB named `jobthing-jobs`, with a string hash key of `id`.

Once these keys have been set, and you have installed Golang, simply run:

```
go run main.go
```

## To Run (in Production):
This project currently has a hideously manual build/run process. Whilst this will be addressed in the future, here is how it is currently done:

**Step 1:** `rsync` recent changes onto servers. Usually like this:

```
rsync -av src/ {IP_ADDRESS}:./job-thingo/src/
```

**Step 2:** SSH into servers, and move to directory specified by previous `rsync` command.

**Step 3:** Set environmental variables

```
export AWS_ACCESS_KEY_ID={YOUR_KEY}
export AWS_SECRET_ACCESS_KEY={YOUR_SECRET_ACCESS_KEY}
export PORT={PORT_REFERENCED_BY_NGINX[1]}
export MARTINI_ENV=production
export BASE_URL=http://www.jobthing.org
```

[1] Note that Nginx is running on all servers, and configured as a forward proxy. The deploy process for that is detailed elsewhere for now.

**Step 4:** Build and restart

```
/usr/src/go/bin/go get
/usr/src/go/bin/go build
kill {PID}
./job-thingo > job-thingo.log &
```

**Step 5:** Check server is running okay.

**Step 6:** Fix this horrible deploy process, so you never have to do that again.

