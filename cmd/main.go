package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"url_shortener/pkg/url_shortener"

	_ "github.com/lib/pq"
)

var hostName string

func main() {
	if val := os.Getenv("hostname"); val == "" {
		panic("Need hostname env value")
	}
	hostName = os.Getenv("hostname")
	//TODO: do this for all envs

	db := connectDB()
	a := url_shortener.NewApp(hostName, db)
	err := a.Start()
	if err != nil {
		log.Fatal(err)
	}
}

func connectDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PG_HOST"), os.Getenv("PG_PORT"), os.Getenv("PG_USER"), os.Getenv("PG_PASSWORD"), os.Getenv("PG_DB_NAME"))

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	log.Println("Ping...")
	err = db.Ping()
	if err != nil {
		log.Println("Ping unsuccessfull")
		log.Fatalln(err)
	}
	log.Println("PONG")
	return db
}
