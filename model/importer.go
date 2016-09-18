package model

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/livereload/api.livereload.com/licensecode"
)

type Importer struct {
	stmt *sql.Stmt
}

func NewImporter(db *sql.DB) (*Importer, error) {
	stmt, err := db.Prepare("INSERT INTO licenses (product_code, product_version, license_type, license_code) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return nil, err
	}

	return &Importer{stmt}, nil
}

func (imp *Importer) Import(license *licensecode.License) (bool, error) {
	_, err := imp.stmt.Exec(license.Product, license.Version, license.Type, license.Code)
	if err != nil {
		if pgerr, ok := err.(*pq.Error); ok {
			if pgerr.Code == errCodeUniqueViolation {
				return false, nil
			}
		}
		return false, err
	}

	return true, nil
}

func (imp *Importer) Commit() error {
	return nil
}
