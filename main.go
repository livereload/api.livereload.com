package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/keighl/postmark"
	_ "github.com/lib/pq"
)

var adminToken string
var paddleToken string
var db *sql.DB
var postmarkClient *postmark.Client

func index(w http.ResponseWriter, r *http.Request) {
	sendErrorMessage(w, http.StatusNotFound, "Nothing to see here, move along.")
}

func main() {
	var err error

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	adminToken = os.Getenv("ADMIN_TOKEN")
	if adminToken == "" {
		log.Fatal("Missing required ADMIN_TOKEN.")
	}

	paddleToken = os.Getenv("PADDLE_TOKEN")
	if paddleToken == "" {
		log.Fatal("Missing required PADDLE_TOKEN.")
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("Missing required DATABASE_URL.")
	}

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(4)

	postmarkEmailAPIKey := os.Getenv("POSTMARK_API_KEY")
	if postmarkEmailAPIKey == "" {
		log.Fatal("Missing required POSTMARK_API_KEY.")
	}
	postmarkClient = postmark.NewClient(postmarkEmailAPIKey, "")

	http.HandleFunc("/licensing/admin/import/", importLicenses)
	http.HandleFunc("/licensing/admin/stats/", showLicensingStats)
	http.HandleFunc("/licensing/callback/paddle", claimLicenseForPaddle)
	http.HandleFunc("/licensing/callback/paddle/", claimLicenseForPaddle)
	http.HandleFunc("/", index)
	log.Printf("Listening on port %s.", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
