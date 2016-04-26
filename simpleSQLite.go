package main

import (
	"database/sql"
	"net/http"
	"html/template"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type PaperStruct struct{
	Authors string
	Year int
	Title string
}

var db *sql.DB //global database connection

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./researchPapers.db")
	checkErr(err)
	defer db.Close()

	createTable()
	c := PaperStruct{authors: "john smith", year: 1996, title: "foo"}
	insertData(c)
	//x := retrieveData()
	//fmt.Printf("%+v\n", x)

	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":8080", nil)
}

func checkErr(err error) {
        if err != nil {
            log.Fatal(err)
        }
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	x := retrieveData()
	t,err := template.ParseFiles("template.html")
	checkErr(err)
	t.Execute(w, x)
}

//create the table if it doesn't exist
func createTable(){
	//create the table if it doesn't exist
	sqlStmt := `
	create table if not exists papers(
		id integer not null primary key, 
		title text, 
		year integer, 
		authors text);`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func insertData(c PaperStruct){
	_,err := db.Exec(fmt.Sprintf("insert into papers(title,year,authors) values('%s',%d,'%s')",c.title,c.year,c.authors))
	checkErr(err)
}

func retrieveData()[]PaperStruct{
	//make a slice of paper structs of length 0
	results := make([]PaperStruct,0)
	rows, err := db.Query("select title,year,authors from papers")
	checkErr(err)

	defer rows.Close()
	var table_data PaperStruct
	for rows.Next() {
		err = rows.Scan(&table_data.title,&table_data.year,&table_data.authors)
		checkErr(err)
		results = append(results, table_data)
	}
	err = rows.Err()
	checkErr(err)
	return results
}
