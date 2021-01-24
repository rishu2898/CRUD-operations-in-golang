package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestReturnSingleEmployee(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Println("error creating mock database")
		return
	}
	defer db.Close()

	// create empHandler with mocked db, request and response to test
	empHandler := &EmployeeHandler{db}

	testCases := []Employee{
		{1, "rk", 22, "M", 2},
		{6, "seek", 15, "M", 2},
	}

	for _, tc := range testCases {
		url := "/employee/%v"
		req, err := http.NewRequest("GET", fmt.Sprintf(url, tc.Id), nil)
		if err != nil {
			t.Fatalf("an error '%s' was not expected while creating request", err)
		}
		req = mux.SetURLVars(req, map[string]string{
			"id": strconv.Itoa(tc.Id),
		})
		// returns an initialized ResponseRecorder
		w := httptest.NewRecorder()

		// before we actually execute our api function, we need to expect required DB actions
		rows := sqlmock.NewRows([]string{"id", "name", "age", "gender", "role"}).
			AddRow(tc.Id, tc.Name, tc.Age, tc.Gender, tc.Role)

		query := "SELECT id, name, age, gender, role FROM employee WHERE id = ?"
		mock.ExpectQuery(query).WithArgs(tc.Id).WillReturnRows(rows)
		empHandler.returnSingleEmployee(w, req)
		if w.Code != 200 {
			t.Fatalf("expected status code to be 500, but got: %d", w.Code)
		}

		data := []Employee{{tc.Id, tc.Name, tc.Age, tc.Gender, tc.Role}}
		// returns the json encoding of data
		// The Marshal() function can take anything, which in Go means the empty interface and return a slice of bytes and error.
		expected, err := json.Marshal(data)
		actual := w.Body.Bytes()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
		}
		check := bytes.Compare(expected, actual[0:len(actual)-1])
		if check != 0 {
			t.Errorf("the expected json: %s is different from actual %s", expected, actual)
		}
	}
}
