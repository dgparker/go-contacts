# go-contacts

#### This project is meant to fulfill the following user story:

As a user I need an online address book exposed as a REST API.  I need the data set to include the following data fields: 
First Name, Last Name, Email Address, and Phone Number

I need the api to follow standard rest semantics to support listing entries, showing a specific single entry, and adding, modifying, and deleting entries.  

The code for the address book should include regular go test files that demonstrate how to exercise all operations of the service.  

Finally I need the service to provide endpoints that can export and import the address book data in a CSV format.

### Endpoints
`GET /entry`<br>
`GET /entry/:id`<br>
`POST /entry`<br>
`PUT /entry`<br>
`DELETE /entry/:id`<br>
`GET /csv/entry`<br>
`POST /csv/entry`<br>


### Dependencies

Libraries - <br>
* https://github.com/julienschmidt/httprouter <br>
* https://github.com/gocarina/gocsv <br>
* https://gopkg.in/mgo.v2 <br>

Database - <br>

* https://www.mongodb.com/ <br>

--

### Setup

#### Clone the repository:<br>
`git clone https://github.com/dgparker/go-contacts`

#### Install dependencies:<br>
If you are running the official Go `dep` tool (https://github.com/golang/dep)
you can simply `cd` into the project directory and run `dep ensure`<br>
##### otherwise run the following command:<br>
`go get github.com/julienschmidt/httprouter github.com/gocarina/gocsv gopkg.in/mgo.v2`

#### Setup mongodb:<br>
If you don't currently have mongodb setup please use the following instructions for your relevant system. https://docs.mongodb.com/manual/administration/install-community/<br>

#### Build:<br>
In an effort to make building ambiguous you must pass your database connection info as well as a defined port for the server to listen on as part of a build flag.
It might look something similar to this for localhost and default settings listening on port 9000: <br>
`go build -ldflags "-X main.dbURI=mongodb://localhost:27017/devetst -X main.dbname=devtest -X main.dbcoll=entries -X main.port=:9000"`

#### Testing:<br>
If you are running the server you can use your preferred method for executing http requests i.e. `curl` or using Postman<br>
If you are using postman you can access my shared collection here https://www.getpostman.com/collections/3978ff85204cce216bda<br>

##### Note: Go test files are included within the `handler pkg` 
you do NOT need a running db to run these test files. They do however expect it's library dependencies to be met
