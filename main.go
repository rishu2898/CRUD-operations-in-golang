package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Employee struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
	Role   int    `json:"role"`
}

type EmployeeHandler struct {
	DB *sql.DB
}

// function to create the table into database
func (empHandler *EmployeeHandler) CreateTable() {
	query := `CREATE TABLE article (
		id int,
		title varchar(50),
		desc varchar(50),
		Content varchar(50)
	)`
	_, err := empHandler.DB.Query(query)
	if err != nil {
		fmt.Println("tale not created")
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

// function to store the data from database into slice
func (empHandler *EmployeeHandler) storeRecord() []Employee {
	var emp []Employee
	dis, err := empHandler.DB.Query("select id, name, age, gender, role from employee")
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
func (empHandler *EmployeeHandler) returnAllEmployees(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllArticles")
	emp := empHandler.storeRecord()
	data, err := json.Marshal(emp)
	if err != nil {
		panic(err)
	}
	w.Write(data)
}

// function for return particular single record
func (empHandler *EmployeeHandler) returnSingleEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		vars := mux.Vars(r)
		key := vars["id"]
		res, _ := empHandler.DB.Query(fmt.Sprintf("SELECT * FROM employee WHERE id=%v", key))

		var emp []Employee
		for res.Next() {
			var row Employee
			err := res.Scan(&row.Id, &row.Name, &row.Age, &row.Gender, &row.Role)
			if err != nil {
				panic(err.Error())
			}
			emp = append(emp, row)
		}
		if len(emp) == 0 {
			fmt.Fprintf(w, "%v not found", http.StatusNoContent)
			return
		}
		json.NewEncoder(w).Encode(emp)
	}
}

// function for insert record
func (empHandler *EmployeeHandler) InsertRecord(w http.ResponseWriter, r *http.Request) {
	var emp Employee
	if r.Method == "POST" {
		err := json.NewDecoder(r.Body).Decode(&emp)
		if err != nil {
			panic(err)
		}
		res, err := empHandler.DB.Query(fmt.Sprintf("INSERT INTO employee(name, age, gender, role) VALUES('%v', %v, '%v', '%v')", emp.Name, emp.Age, emp.Gender, emp.Role))
		if err != nil {
			panic(err.Error())
		}
		data, err := json.Marshal(emp)
		w.Write(data)
		res.Close()
	}
}

// function for update single record
func (empHandler *EmployeeHandler) UpdateSingleRecord(w http.ResponseWriter, r *http.Request) {
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
		res, _ := empHandler.DB.Exec(fmt.Sprintf("UPDATE employee SET name = '%v', age = %v, gender = '%v', role = '%v' where id = %v", emp.Name, emp.Age, emp.Gender, emp.Role, key))
		// return the number of rows affected by query
		count, _ := res.RowsAffected()
		if count == 0 {
			fmt.Fprintf(w, "%v not found", http.StatusNoContent)
			return
		}
		fmt.Fprintf(w, "record having id: %v updated successfully", key)
	}
}

func (empHandler *EmployeeHandler) DeleteSingleRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		vars := mux.Vars(r)
		key := vars["id"]
		fmt.Println(key)
		res, _ := empHandler.DB.Exec(fmt.Sprintf("DELETE FROM employee WHERE id = %v", key))
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
	// connecting to the database
	db, err := sql.Open("mysql", "rishabh:Rishu2898@@(127.0.0.1)/company")
	if err != nil {
		panic(err.Error())
	}
	empHandler := &EmployeeHandler{DB: db}

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/employee", empHandler.returnAllEmployees).Methods("GET")
	myRouter.HandleFunc("/employee/{id}", empHandler.returnSingleEmployee).Methods("GET")
	myRouter.HandleFunc("/employee", empHandler.InsertRecord).Methods("POST")
	myRouter.HandleFunc("/employee/{id}", empHandler.UpdateSingleRecord).Methods("PUT")
	myRouter.HandleFunc("/employee/{id}", empHandler.DeleteSingleRecord).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", myRouter))
}
func main() {
	handleRequests()
}
