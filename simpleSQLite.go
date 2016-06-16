package main

import (
	"database/sql"
	"html/template"
	"net/http"
	//"fmt"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

//A PaperStruct is a struct that describes the structure of a single research paper
type PaperStruct struct {
	ID      int
	Authors string
	Year    int
	Title   string
	Title2  string
	Tag1    string
	Tag2    string
	Tag3    string
	Tag4    string
	Notes   string
}

//A PageMulti is a struct that describes a single HTML Page with multiple papers
type PageMulti struct {
	Title  string
	Papers []PaperStruct
}

//A PageSingle is a struct that describes a single HTML Page with a single paper
type PageSingle struct {
	Title  string
	CSS    string
	Papers PaperStruct
}

var db *sql.DB //global database connection

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./researchPapers.db")
	checkErr(err)
	defer db.Close()

	createTable()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/add/", addHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/update/", updateHandler)
	http.ListenAndServe(":8080", nil)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	x := PageMulti{Title: "View Papers", Papers: retrieveMultiLineData()}
	tmplBytes, err := Asset("resources/view_template.html")
	checkErr(err)
	t, err := template.New("viewtmpl").Parse(string(tmplBytes))
	checkErr(err)
	err = t.Execute(w, x)
	checkErr(err)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	x := PageSingle{Title: "Add Record"}
	tmplBytes, err := Asset("resources/add_template.html")
	checkErr(err)
	t, err := template.New("addtmpl").Parse(string(tmplBytes))
	checkErr(err)
	err = t.Execute(w, x)
	checkErr(err)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/view/", http.StatusFound)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	yearInt, err := strconv.Atoi(r.FormValue("Year"))
	checkErr(err)
	c := PaperStruct{
		Authors: r.FormValue("Authors"),
		Year:    yearInt,
		Title:   r.FormValue("Title"),
		Title2:  r.FormValue("Title2"),
		Tag1:    r.FormValue("Tag1"),
		Tag2:    r.FormValue("Tag2"),
		Tag3:    r.FormValue("Tag3"),
		Tag4:    r.FormValue("Tag4"),
		Notes:   r.FormValue("Notes")}
	insertData(c)
	http.Redirect(w, r, "/", http.StatusFound)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	recordID, err := strconv.Atoi(r.URL.Path[6:]) //strip "/edit/" from the URL string and convert to int
	checkErr(err)
	x := PageSingle{Title: "Add Record", Papers: retrieveRecord(recordID)}
	tmplBytes, err := Asset("resources/edit_template.html")
	checkErr(err)
	t, err := template.New("edittmpl").Parse(string(tmplBytes))
	checkErr(err)
	err = t.Execute(w, x)
	checkErr(err)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	yearInt, err := strconv.Atoi(r.FormValue("Year"))
	checkErr(err)
	idInt, err := strconv.Atoi(r.FormValue("ID"))
	checkErr(err)
	c := PaperStruct{
		Authors: r.FormValue("Authors"),
		Year:    yearInt,
		Title:   r.FormValue("Title"),
		Title2:  r.FormValue("Title2"),
		Tag1:    r.FormValue("Tag1"),
		Tag2:    r.FormValue("Tag2"),
		Tag3:    r.FormValue("Tag3"),
		Tag4:    r.FormValue("Tag4"),
		Notes:   r.FormValue("Notes")}
	_, err = db.Exec("update papers set authors=?,year=?,title=?,title2=?,tag1=?,tag2=?,tag3=?,tag4=?,notes=? where id = ?", c.Authors, c.Year, c.Title, c.Title2, c.Tag1, c.Tag2, c.Tag3, c.Tag4, c.Notes, idInt)
	checkErr(err)
	http.Redirect(w, r, "/", http.StatusFound)
}

//create the table if it doesn't exist
func createTable() {
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

func insertData(c PaperStruct) {
	_, err := db.Exec("insert into papers(authors,year,title,title2,tag1,tag2,tag3,tag4,notes) values(?,?,?,?,?,?,?,?,?)", c.Authors, c.Year, c.Title, c.Title2, c.Tag1, c.Tag2, c.Tag3, c.Tag4, c.Notes)
	checkErr(err)
}

func retrieveRecord(id int) PaperStruct {
	row, err := db.Query("select * from papers where ID = ?", id)
	checkErr(err)
	defer row.Close()
	var tableData PaperStruct
	for row.Next() {
		err = row.Scan(&tableData.ID, &tableData.Authors, &tableData.Year, &tableData.Title, &tableData.Title2, &tableData.Tag1, &tableData.Tag2, &tableData.Tag3, &tableData.Tag4, &tableData.Notes)
		checkErr(err)
	}
	return tableData
}

func retrieveMultiLineData() []PaperStruct {
	//make a slice of paper structs of length 0
	var results []PaperStruct
	rows, err := db.Query("select * from papers")
	checkErr(err)

	defer rows.Close()
	var tableData PaperStruct
	for rows.Next() {
		err = rows.Scan(&tableData.ID, &tableData.Authors, &tableData.Year, &tableData.Title, &tableData.Title2, &tableData.Tag1, &tableData.Tag2, &tableData.Tag3, &tableData.Tag4, &tableData.Notes)
		checkErr(err)
		results = append(results, tableData)
	}
	err = rows.Err()
	checkErr(err)
	return results
}
