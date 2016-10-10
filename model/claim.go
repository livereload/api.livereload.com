package model

import (
	"database/sql"

	"github.com/shopspring/decimal"
)

type Claim struct {
	Store string
	Qty   int

	Txn     string
	Ticket  string
	OrderID string

	Company   string
	FirstName string
	LastName  string
	FullName  string
	Email     string
	Country   string

	Message    string
	Notes      string
	Additional string

	Raw string

	Currency     string
	Price        decimal.Decimal
	SaleGross    decimal.Decimal
	SaleTax      decimal.Decimal
	ProcessorFee decimal.Decimal
	Earnings     decimal.Decimal

	Coupon        string
	CouponSavings decimal.Decimal
}

func ClaimLicense(db *sql.DB, product, version, typ string, claim *Claim) (string, error) {
	var licenseCode string
	err := db.QueryRow(`
        UPDATE licenses
        SET claimed = TRUE, claimed_at = NOW(),
            claim_store = $1, claim_qty = $2,
            claim_txn = $3, claim_ticket = $4,
            claim_company = $5, claim_first_name = $6, claim_last_name = $7, claim_full_name = $8, claim_email = $9, claim_country = $10,
            claim_message = $11, claim_notes = $12, claim_additional = $13,
            claim_raw = $14,
            claim_currency = $15, claim_price = $16, claim_sale_gross = $17, claim_sale_tax = $18, claim_processor_fee = $19, claim_earnings = $20,
            claim_coupon = $21, claim_coupon_savings = $22,
            claim_order_id = $23
        WHERE
            id = (SELECT id FROM licenses WHERE product_code = $24 AND product_version = $25 AND license_type = $26 AND NOT claimed LIMIT 1 FOR UPDATE) RETURNING license_code`,
		claim.Store, claim.Qty,
		claim.Txn, claim.Ticket,
		claim.Company, claim.FirstName, claim.LastName, claim.FullName, claim.Email, claim.Country,
		claim.Message, claim.Notes, claim.Additional,
		claim.Raw,
		claim.Currency, claim.Price.String(), claim.SaleGross.String(), claim.SaleTax.String(), claim.ProcessorFee.String(), claim.Earnings.String(),
		claim.Coupon, claim.CouponSavings.String(),
		claim.OrderID,
		product, version, typ,
	).Scan(&licenseCode)
	if err != nil {
		return "", err
	}
	return licenseCode, nil
}
