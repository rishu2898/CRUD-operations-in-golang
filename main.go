package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Employee struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
	Role   int    `json:"role"`
}

var DB *sql.DB

// function to create the table into database
func CreateTable() {
	query := `CREATE TABLE article (
		id int,
		title varchar(50),
		desc varchar(50),
		Content varchar(50)
	)`
	_, err := DB.Query(query)
	if err != nil {
		fmt.Println("tale not created")
	}
}
// function to connect the database
func ConnectDB() {
	db, err := sql.Open("mysql", "rishabh:Rishu2898@@(127.0.0.1)/company")
	if err != nil {
		panic(err.Error())
	}
	DB = db
	//defer DB.Close()
	//CreateTable()
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}
// function to store the data from database into slice
func storeRecord() []Employee {
	var emp []Employee
	dis, err := DB.Query("select * from employee")
	if err != nil {
		panic(err.Error())
	}
	for dis.Next() {
		var row Employee
		err = dis.Scan(&row.Id, &row.Name, &row.Age, &row.Gender, &row.Role)
		if err != nil {
			panic(err)
		}
		emp = append(emp, row)
	}
	return emp
}
// function for find all record
func returnAllEmployees(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	emp := storeRecord()
	data, err := json.Marshal(emp)
	if err != nil {
		panic(err)
	}
	w.Write(data)
}
// function for return particular single record
func returnSingleEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	emp := storeRecord()
	for _, rec := range emp {
		if strconv.Itoa(rec.Id) == key {
			data, err := json.Marshal(rec)
			if err != nil {
				panic(err)
			}
			w.Write(data)
			return
		}
	}
	fmt.Fprintf(w, "%v page not found", http.StatusNotFound)
}
// function for insert record
func InsertRecord(w http.ResponseWriter, r *http.Request) {
	var emp Employee
	if r.Method == "POST" {
		err := json.NewDecoder(r.Body).Decode(&emp)
		if err != nil {
			panic(err)
		}
		res, err := DB.Query(fmt.Sprintf("INSERT INTO employee(id, name, age, gender, role) VALUES(%v, '%v', %v, '%v', '%v')", emp.Id, emp.Name, emp.Age, emp.Gender, emp.Role))
		if err != nil {
			panic(err.Error())
		}
		data, err := json.Marshal(emp)
		w.Write(data)
		res.Close()
	}
}
// function for update single record
func UpdateSingleRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		vars := mux.Vars(r)
		key := vars["id"]

		// A Decoder reads and decodes JSON values from an input stream.
		// NewDecoder returns a new decoder that reads from r.
		// Decode reads the next JSON-encoded value from its input and stores it in the value pointed to by v.
		var emp Employee
		err := json.NewDecoder(r.Body).Decode(&emp)
		if err != nil {
			panic(err)
		}
		res, _ := DB.Exec(fmt.Sprintf("UPDATE employee SET name = '%v', age = %v, gender = '%v', role = '%v' where id = %v", emp.Name, emp.Age, emp.Gender, emp.Role, key))
		// return the number of rows affected by query
		count, _ := res.RowsAffected()
		if count == 0 {
			fmt.Fprintf(w, "%v not found", http.StatusNoContent)
			return
		}
		fmt.Fprintf(w, "record having id: %v updated successfully", key)
	}
}

func DeleteSingleRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		vars := mux.Vars(r)
		key := vars["id"]

		res, _ := DB.Exec(fmt.Sprintf("DELETE FROM employee WHERE id = %v", key))
		// return the number of rows affected by query
		count, _ := res.RowsAffected()
		if count == 0 {
			fmt.Fprintf(w, "%v not found", http.StatusNoContent)
			return
		}
		fmt.Fprintf(w, "record having id: %v deleted successfully", key)
	}
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/employees", returnAllEmployees).Methods("GET")
	myRouter.HandleFunc("/employee/{id}", returnSingleEmployee).Methods("GET")
	myRouter.HandleFunc("/insert", InsertRecord).Methods("POST")
	myRouter.HandleFunc("/update/{id}", UpdateSingleRecord).Methods("PUT")
	myRouter.HandleFunc("/delete/{id}", DeleteSingleRecord).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", myRouter))
}
func main() {
	ConnectDB()
	handleRequests()
}
