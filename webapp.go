// Aaron Salazar
//
// change dsn string to connect to mysql db
// visit http://localhost:3000/home/ to view webpage
//

package main
  
import (
  "html/template"
	"log"
  "net/http"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
)

type DataTable struct{
  ItemName []string
}

var dataTable *DataTable
var homeTemplate = template.Must(template.ParseFiles("home.html"))
var name string
var ptrDB *sql.DB
var dsn string = "user:password@tcp(localhost:3306)/" 

// check() deals with error checking
func check(err error){
	if err != nil {
		log.Fatal(err)
	}
}

// setupDB() checks if required db/table exists and creates them if needed 
func setupDB(dbName string, tbName string){
	db, err := sql.Open("mysql",dsn) 
	check(err)
	defer db.Close()
	_, err = db.Query("create database if not exists " +dbName)
	check(err) 
	_, err = db.Query("create table if not exists "+dbName+"."+tbName+" (name varchar(50))")
	check(err) 
}

// loadData() transfers data from the db to a DataTable
func loadData(){
	dataTable = &DataTable{make([]string,0)}
	rows, err := ptrDB.Query("select name from item")
	check(err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name)
		check(err)
		dataTable.ItemName = append(dataTable.ItemName,name);
	}
}

// renderTemplate() shows webpage
func renderTemplate(w http.ResponseWriter, r *http.Request){
	err := homeTemplate.Execute(w,dataTable)
	check(err)
}

// saveHandler() updates the db with data from the textbox, then
// redirects the user to /home/
func saveHandler(w http.ResponseWriter, r *http.Request) {
	input := r.FormValue("input")
	_, err := ptrDB.Query("insert into item(name) values(?)",input)
	check(err)
	http.Redirect(w,r,"/home/",http.StatusFound)
}

// homeHandler() deals with requests to /home/
func homeHandler(w http.ResponseWriter, r *http.Request){
	loadData()
	renderTemplate(w,r)
}

func main(){
	setupDB("sphere","item")
	db, err := sql.Open("mysql",dsn+"sphere") 
	ptrDB = db
	check(err)
	defer db.Close()

  http.HandleFunc("/home/",homeHandler)
	http.HandleFunc("/save/",saveHandler)
  http.ListenAndServe(":3000",nil)
}
