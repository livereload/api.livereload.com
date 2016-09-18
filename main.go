package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/livereload/api.livereload.com/licensecode"
	"github.com/livereload/api.livereload.com/model"
)

var adminToken string
var db *sql.DB

type outcome struct {
	license *licensecode.License
	result  string
}

func importLicenses(w http.ResponseWriter, r *http.Request) {
	deadline := time.Now().Add(25 * time.Second)

	err := verifyToken(r, adminToken)
	if err != nil {
		sendError(w, err)
		return
	}

	if r.Method != http.MethodPost {
		sendErrorMessage(w, http.StatusMethodNotAllowed, "")
		return
	}

	var licenses []*licensecode.License

	scanner := bufio.NewScanner(io.LimitReader(r.Body, 10*1024*1024))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == '#' {
			continue
		}

		license, err := licensecode.Parse(line)
		if err != nil {
			sendErrorFmt(w, http.StatusBadRequest, "Invalid license: %s", line)
			return
		}

		licenses = append(licenses, license)
	}
	if err := scanner.Err(); err != nil {
		sendError(w, err)
		return
	}

	// rows, err db.Query
	// queryRow(ctx, 'UPDATE licenses SET claimed = TRUE, claimed_at = NOW(), claim_store = $1, claim_txn = $2, claim_qty = 1, claim_first_name = $3, claim_last_name = $4, claim_email = $5, claim_notes = $6, claim_full_name = $7, claim_message = $8, claim_raw = $9, claim_currency = $10, claim_price = $11, claim_coupon = $12, claim_coupon_savings = $13, claim_additional = $14, claim_country = $15, claim_sale_gross = $16, claim_sale_tax = $17, claim_processor_fee = $18, claim_earnings = $19 WHERE id IN (SELECT id FROM licenses WHERE product_code = \'LR\' AND license_type = \'A\' AND NOT claimed LIMIT 1) RETURNING id, license_code as "licenseCode";', params.store, params.txn, params.firstName || null, params.lastName || null, params.email, params.notes || '', params.fullName || null, params.message || null, params.raw || null, params.currency || null, params.price || null, params.coupon || null, params.savings || null, params.additional || null, params.country || null, params.gross || null, params.tax || null, params.fee || null, params.earnings || null)
	// queryValue(ctx, "SELECT COUNT(*) FROM licenses WHERE product_code = 'LR' AND license_type = 'A' AND NOT claimed")
	// queryValue(ctx, "SELECT COUNT(*) FROM licenses WHERE product_code = 'LR' AND license_type = 'A' AND NOT claimed")

	if len(licenses) == 0 {
		sendErrorMessage(w, http.StatusNoContent, "No licenses specified")
		return
	}

	var outcomes []outcome
	var added int
	var existed int
	var timedout int
	var count int
	var last *licensecode.License

	imp, err := model.NewImporter(db)
	if err != nil {
		sendError(w, err)
		return
	}
	for _, license := range licenses {
		var result string
		if time.Now().After(deadline) {
			result = "timeout"
			timedout++
		} else {
			isNew, err := imp.Import(license)
			if err != nil {
				sendError(w, err)
				return
			}
			if isNew {
				result = "added"
				added++
			} else {
				result = "existed"
				existed++
			}
			last = license
		}
		outcomes = append(outcomes, outcome{license, result})
		count++
	}

	err = imp.Commit()
	if err != nil {
		sendError(w, err)
		return
	}

	// components := strings.Split(r.URL.Path, "/")

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Added %d licenses, existed %d, timed out %d.\n\n", added, existed, timedout)
	if last != nil {
		fmt.Fprintf(w, "Last license processed: %s.\n\n", last.Code)
	}
	for _, o := range outcomes {
		fmt.Fprintf(w, "%-7s %s\n", o.result, o.license.Code)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
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

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("Missing required DATABASE_URL.")
	}

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(4)

	// age := 21
	// rows, err := db.Query("SELECT name FROM users WHERE age = $1", age)

	http.HandleFunc("/licensing/admin/import/", importLicenses)
	http.HandleFunc("/", index)
	log.Printf("Listening on port %s.", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
