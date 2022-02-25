package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/argon2"
)

var db *pgxpool.Pool

func main() {
	// Explicitly declaring err here to avoid using := syntax in the
	// db connection statement. Using := will create a new db variable
	// limited to this scope instead of initializing the global db var.
	var err error

	db, err = pgxpool.Connect(context.Background(), os.Getenv("DB_CONN"))
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/signup", handleSignup)

	http.HandleFunc("/login", handleLogin)

	log.Fatalln(http.ListenAndServe(":8080", nil))
}

func handleSignup(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "/signup only accept POST request.\n")
	} else {

		username := r.FormValue("username")
		password := r.FormValue("password")

		salt := make([]byte, 32)
		_, err = rand.Read(salt)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Internal server error: %v\n", err)
		} else {
			saltStr := fmt.Sprintf("%x", salt)

			passHash := argon2.IDKey([]byte(password), salt, 10, 64*1024, 8, 32)
			passHashStr := fmt.Sprintf("%x", passHash)

			_, err = db.Exec(context.Background(), "INSERT INTO signuplogin VALUES ($1, $2, $3)", username, passHashStr, saltStr)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Creating user account. Username must be unique: %v\n", err)
			} else {
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "User account created.\n")
			}
		}
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "/login only accept POST request.\n")
	} else {

		username := r.FormValue("username")
		password := r.FormValue("password")
		var saltStr string
		var passHashStr string

		err = db.QueryRow(context.Background(), "SELECT passwordhash, salt FROM signuplogin WHERE username=$1", username).Scan(&passHashStr, &saltStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "No user found: %v\n", err)
		} else {

			var salt []byte
			salt, err = hex.DecodeString(saltStr)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "Internal server error: %v\n", err)
			} else {

				newPassHash := argon2.IDKey([]byte(password), salt, 10, 64*1024, 8, 32)
				newPassHashStr := fmt.Sprintf("%x", newPassHash)
				if newPassHashStr == passHashStr {
					w.WriteHeader(http.StatusOK)
					fmt.Fprintf(w, "Login successful\n")
				} else {
					w.WriteHeader(http.StatusBadRequest)
					fmt.Fprintf(w, "Invalid password\n")
				}
			}
		}
	}
}
