package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

func showLicensingStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendErrorMessage(w, http.StatusMethodNotAllowed, "")
		return
	}

	var cells [][]string
	rows, err := db.Query("SELECT product_code, product_version, license_type, COUNT(*) FROM licenses WHERE NOT claimed GROUP BY product_code, product_version, license_type")
	if err != nil {
		sendError(w, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var code, version, typ string
		var count int
		err = rows.Scan(&code, &version, &typ, &count)
		if err != nil {
			sendError(w, err)
			return
		}

		prefix := fmt.Sprintf("%s%s%s", code, version, typ)

		cells = append(cells, []string{prefix, strconv.Itoa(count)})
	}

	w.Header().Set("Content-Type", "text/plain")
	table := tablewriter.NewWriter(w)
	table.SetHeader([]string{"License Type", "Count"})
	table.SetCenterSeparator("-")
	// table.SetBorder(false)
	table.AppendBulk(cells)
	table.Render()
}
