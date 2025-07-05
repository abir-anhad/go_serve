// Test Code

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Name string `json:"name"`
	Age  int16  `json:"age"`
}

type HealthReport struct {
	Live      bool  `json:"live"`
	ErrCounts int16 `json:"errorcounts"`
}

type DefaultResponse struct {
	Res string `json:"res"`
}

// refactoring
// first lets write a methodguard

// HTTP method Guard - middleware
func methodGuard(allowedMethod string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != allowedMethod {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		handler(w, req)
	}
}

// writeJSON response - utility
func writeJSONResponse(w http.ResponseWriter, data interface{}) error {
	jsonResponse, err := json.MarshalIndent(data, "", " ")

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	// fmt.Fprint(w, string(jsonResponse))
	w.Write(jsonResponse)

	return nil
}

// handler for '/'
func rootHandler(w http.ResponseWriter, req *http.Request) {

	response := DefaultResponse{
		Res: "Welcome to root route",
	}
	err := writeJSONResponse(w, response)

	if err != nil {
		// debug print
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// handler for '/api/users'
func usersHandler(w http.ResponseWriter, req *http.Request) {

	users := []User{
		{Name: "Abir", Age: 33},
		{Name: "Babai", Age: 23},
	}

	err := writeJSONResponse(w, users)

	if err != nil {
		// debug print
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// handler for '/health'
func healthHandler(w http.ResponseWriter, req *http.Request) {

	reports := []HealthReport{
		{Live: true, ErrCounts: 20},
		{Live: false, ErrCounts: 1200},
	}

	err := writeJSONResponse(w, reports)

	if err != nil {
		// debug print
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// main func
func main() {
	http.HandleFunc("/", methodGuard(http.MethodGet, rootHandler))
	http.HandleFunc("/api/users", methodGuard(http.MethodGet, usersHandler))
	http.HandleFunc("/health", methodGuard(http.MethodGet, healthHandler))
	http.ListenAndServe(":8080", nil)
}
