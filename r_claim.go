package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/livereload/api.livereload.com/licensecode"
	"github.com/livereload/api.livereload.com/model"
)

func claimLicense(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Added %d licenses, existed %d, timed out %d.\n\n", added, existed, timedout)
	if last != nil {
		fmt.Fprintf(w, "Last license processed: %s.\n\n", last.Code)
	}
	for _, o := range outcomes {
		fmt.Fprintf(w, "%-7s %s\n", o.result, o.license.Code)
	}
}
