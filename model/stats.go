package model

import (
	"database/sql"
)

func CountUnclaimedLicenses(db *sql.DB, product, version, typ string) (int, error) {
	rows, err := db.Query("SELECT COUNT(*) FROM licenses WHERE NOT claimed AND product_code = $1 AND product_version = $2 AND license_type = $3", product, version, typ)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	if !rows.Next() {
		return 0, nil
	}

	var count int
	err = rows.Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}
