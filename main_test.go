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

	// create app with mocked db, request and response to test
	app := &EmployeeHandler{db}

	testCases := []struct {
		id     int
		name   string
		age    int
		gender string
		role   int
	}{
		{1, "rk", 22, "M", 2},
		{6, "seek", 15, "M", 2},
	}

	for _, tc := range testCases {
		url := "/employee/%v"
		req, err := http.NewRequest("GET", fmt.Sprintf(url, tc.id), nil)
		if err != nil {
			t.Fatalf("an error '%s' was not expected while creating request", err)
		}
		req = mux.SetURLVars(req, map[string]string{
			"id": strconv.Itoa(tc.id),
		})
		// returns an initialized ResponseRecorder
		w := httptest.NewRecorder()

		// before we actually execute our api function, we need to expect required DB actions
		rows := sqlmock.NewRows([]string{"id", "name", "age", "gender", "role"}).
			AddRow(tc.id, tc.name, tc.age, tc.gender, tc.role)

		query := "SELECT id, name, age, gender, role FROM employee WHERE id = ?"
		mock.ExpectQuery(query).WithArgs(tc.id).WillReturnRows(rows)
		app.returnSingleEmployee(w, req)
		if w.Code != 200 {
			t.Fatalf("expected status code to be 500, but got: %d", w.Code)
		}

		data := struct {
			employee []*Employee
		}{employee: []*Employee{
			{tc.id, tc.name, tc.age, tc.gender, tc.role},
		}}

		expected, err := json.Marshal(data)
		actual := w.Body.Bytes()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when marshaling expected json data", err)
		}

		if bytes.Compare(expected, actual) != 0 {
			t.Errorf("the expected json: %s is different from actual %s", expected, actual)
		}
	}
}
