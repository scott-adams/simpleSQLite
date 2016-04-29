package main

import (
	"database/sql"
	"net/http"
	"html/template"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"log"
)

type PaperStruct struct{
	ID int
	Authors string
	Year int
	Title string
	Title2 string
	Tag1 string
	Tag2 string
	Tag3 string
	Tag4 string
	Notes string
}

type Page struct{
	Title string
}

var db *sql.DB //global database connection

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./researchPapers.db")
	checkErr(err)
	defer db.Close()

	createTable()

	http.HandleFunc("/", viewFunc)
	http.HandleFunc("/add/",addFunc)
	http.HandleFunc("/save/",saveFunc)
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	http.ListenAndServe(":8080", nil)
}

func checkErr(err error) {
        if err != nil {
            log.Fatal(err)
        }
}

func viewFunc(w http.ResponseWriter, r *http.Request) {
	x := retrieveData()
	t,err := template.ParseFiles("templates/view_template.html")
	checkErr(err)
	err = t.Execute(w,x)
	checkErr(err)
}

func addFunc(w http.ResponseWriter, r *http.Request) {
	x:= Page{Title: "Add Record"}
	t,err := template.ParseFiles("templates/add_template.html")
	checkErr(err)
	err = t.Execute(w,x)
	checkErr(err)
}

func saveFunc(w http.ResponseWriter, r *http.Request) {
	yearInt, err := strconv.Atoi(r.FormValue("Year"))
	checkErr(err)
	c := PaperStruct{
		Authors:r.FormValue("Authors"),
		Year:yearInt,
		Title:r.FormValue("Title"),
		Title2:r.FormValue("Title2"),
		Tag1:r.FormValue("Tag1"),
		Tag2:r.FormValue("Tag2"),
		Tag3:r.FormValue("Tag3"),
		Tag4:r.FormValue("Tag4"),
		Notes:r.FormValue("Notes")}
	insertData(c)
	http.Redirect(w, r, "/", http.StatusFound)
}

//create the table if it doesn't exist
func createTable(){
	//create the table if it doesn't exist
	sqlStmt := `
	create table if not exists papers(
		id integer not null primary key, 
		authors text, 
		year integer, 
		title text,
		title2 text,
		tag1 text,
		tag2 text,
		tag3 text,
		tag4 text,
		notes text);`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func insertData(c PaperStruct){
	_,err := db.Exec(fmt.Sprintf("insert into papers(authors,year,title,title2,tag1,tag2,tag3,tag4,notes) values('%s',%d,'%s','%s','%s','%s','%s','%s','%s')",c.Authors,c.Year,c.Title,c.Title2,c.Tag1,c.Tag2,c.Tag3,c.Tag4,c.Notes))
	checkErr(err)
}

func retrieveData()[]PaperStruct{
	//make a slice of paper structs of length 0
	results := make([]PaperStruct,0)
	rows, err := db.Query("select * from papers")
	checkErr(err)

	defer rows.Close()
	var table_data PaperStruct
	for rows.Next() {
		err = rows.Scan(&table_data.ID,&table_data.Authors,&table_data.Year,&table_data.Title,&table_data.Title2,&table_data.Tag1,&table_data.Tag2,&table_data.Tag3,&table_data.Tag4,&table_data.Notes)
		checkErr(err)
		results = append(results, table_data)
	}
	err = rows.Err()
	checkErr(err)
	return results
}
