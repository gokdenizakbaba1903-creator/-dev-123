package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	// Setup SQLite
	var err error
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	db.Exec("CREATE TABLE users (id INTEGER, name TEXT, secret TEXT)")
	db.Exec("INSERT INTO users VALUES (1, 'Admin', 'SUPER_SECRET_TOKEN_123')")
	db.Exec("INSERT INTO users VALUES (2, 'User', 'HELLO_WORLD')")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<h1>Vulnerable Target App</h1><p>Welcome to the security lab target.</p>")
	})

	// VULNERABLE: Reflected XSS
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<h2>Search results for: %s</h2>", query)
	})

	// VULNERABLE: SQL Injection
	http.HandleFunc("/api/user", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		// The query is intentionally built using string fmt for SQLi
		query := fmt.Sprintf("SELECT name FROM users WHERE id = %s", id)
		
		var name string
		err := db.QueryRow(query).Scan(&name)
		if err != nil {
			// Leaking SQL error to client
			http.Error(w, "Database Error: "+err.Error(), 500)
			return
		}
		fmt.Fprintf(w, "User: %s", name)
	})

	// VULNERABLE: Directory Listing / File Exposure
	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		data, _ := os.ReadFile("config.json") // Simple enough
		w.Write(data)
	})

	fmt.Println("Vulnerable Target App running on :9090")
	fmt.Println("Try: http://localhost:9090/search?q=<script>alert(1)</script>")
	fmt.Println("Try: http://localhost:9090/api/user?id=1 OR 1=1")
	
	log.Fatal(http.ListenAndServe(":9090", nil))
}
