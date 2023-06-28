# Learning from udemy course.
## Simple gRPC Server and Client for CRUD requests. Also, a Json Server available. Using  SQLite as a database ( list.db).  

## Generate Server and Client go code from .proto: 
```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/mail.proto
```
## APIs Runners: 
```
go run .
```

## Client Requests: 
```
go run ./client/
```
