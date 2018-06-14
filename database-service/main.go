package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "micropostgres.postgres.database.azure.com"
	port     = 5432
	user     = "user"
	password = "password"
	dbname   = "Customer"
)

func main() {
	log.Println("Before response")

	// jsonResponse := retrieveDatabaseRecords()
	// os.Stdout.Write(jsonResponse)

	http.HandleFunc("/get", getAllRecords)
	http.HandleFunc("/create", createRecord)

	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}

func getAllRecords(w http.ResponseWriter, r *http.Request) {
	jsonResponse := retrieveDatabaseRecords()
	w.Write(jsonResponse)
}

func createRecord(w http.ResponseWriter, r *http.Request) {

	name := r.URL.Query().Get("name")
	apikey := r.URL.Query().Get("apikey")

	if len(name) <= 0 || len(apikey) <= 0 {
		log.Fatal("No argument provided")
	}

	log.Println("Writing record")
	db := getConnection()
	_, err := db.Exec("insert into \"Customer\" values ($1, $2)", name, apikey)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Record created")
	w.Write([]byte("{message:\"Record successfully created\"}"))
}

func retrieveDatabaseRecords() (response []byte) {
	fmt.Println("Trying to connect to postgres...")

	db := getConnection()

	//rows, err := db.Query("select id, name from test")
	rows, err := db.Query("select apikey, name from \"Customer\"")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var records []Record

	for rows.Next() {
		var apikey sql.NullString
		var name sql.NullString
		err := rows.Scan(&apikey, &name)
		if err != nil {
			log.Fatal(err)
		}
		if name.Valid && apikey.Valid {
			// log.Println(name.String, apikey.String)
			records = append(records, Record{name.String, apikey.String})
		} else if name.Valid {
			// log.Println(name.String, ", null")
			records = append(records, Record{name.String, ", null"})
		} else if apikey.Valid {
			// log.Println("null, ", apikey.String)
			records = append(records, Record{"null, ", apikey.String})
		} else {
			// log.Println("null, null")
			records = append(records, Record{"null", "null"})
		}
	}

	fmt.Println("Successfully retrieved records!")
	json, _ := json.Marshal(records)

	return json
}

func getConnection() (dbout *sql.DB) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db
}

type Record struct {
	Name   string
	Apikey string
}
