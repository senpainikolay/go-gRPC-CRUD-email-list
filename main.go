package main

import (
	"database/sql"
	"log"
	"simple_proj/db_sql"
	"simple_proj/grpcapi"
	"simple_proj/jsonapi"
	"sync"

	"github.com/alexflint/go-arg"
)

var args struct {
	DBpath   string `arg:"env:DBPATH"`
	JsonPort string `arg:"env:APIPORT"`
	GrpcPort string `arg:"env:gRPCPORT"`
}

func main() {
	arg.MustParse(&args)

	if args.DBpath == "" {
		args.DBpath = "list.db"
	}
	if args.JsonPort == "" {
		args.JsonPort = ":8080"
	}
	if args.GrpcPort == "" {
		args.GrpcPort = ":8081"
	}

	log.Printf("using database '%v'\n", args.DBpath)
	db, err := sql.Open("sqlite3", args.DBpath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db_sql.TryCreate(db)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		log.Printf("starting JSON API server...\n")
		jsonapi.Serve(db, args.JsonPort)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		log.Printf("starting gRPC API server...\n")
		grpcapi.Serve(db, args.GrpcPort)
		wg.Done()
	}()

	wg.Wait()
}
