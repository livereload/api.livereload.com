package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"unicode/utf8"

	"github.com/shopspring/decimal"

	"github.com/livereload/api.livereload.com/licensecode"
	"github.com/livereload/api.livereload.com/model"
)

type paddleClaimRequest struct {
	Token string `json:"token"`

	Txn     string `json:"txn"`
	Qty     string `json:"quantity"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`

	Passthrough string `json:"passthrough"`

	ProductID         string `json:"p_product_id"`
	OrderID           string `json:"p_order_id"`
	Country           string `json:"p_country"`
	Coupon            string `json:"p_coupon"`
	CouponSavings     string `json:"p_coupon_savings"`
	Currency          string `json:"p_currency"`
	Earnings          string `json:"p_earnings"`
	PaddleFee         string `json:"p_paddle_fee"`
	Price             string `json:"p_price"`
	Quantity          string `json:"p_quantity"`
	SaleGross         string `json:"p_sale_gross"`
	TaxAmount         string `json:"p_tax_amount"`
	UsedPriceOverride string `json:"p_used_price_override"`

	Signature string `json:"p_signature"`
}

func claimLicenseForPaddle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorMessage(w, http.StatusMethodNotAllowed, "")
		return
	}

	ctype, _ := parseRequestContentType(r)
	if ctype != "application/json" {
		sendErrorMessage(w, http.StatusBadRequest, "application/json Content-Type required")
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1024*1024))
	if err != nil {
		sendErrorMessage(w, http.StatusBadRequest, "Failed to read the POST payload")
		return
	}
	if !utf8.Valid(body) {
		sendErrorMessage(w, http.StatusBadRequest, "POST payload is not a valid UTF-8 string")
		return
	}
	log.Printf("body: %v", string(body))

	var rq paddleClaimRequest
	err = json.Unmarshal(body, &rq)
	if err != nil {
		sendErrorFmt(w, http.StatusBadRequest, "Failed to decode POST payload as JSON: %v", err)
		return
	}

	err = verifyToken(rq.Token, paddleToken)
	if err != nil {
		sendError(w, err)
		return
	}

	if rq.ProductID != "489469" {
		sendErrorMessage(w, http.StatusBadRequest, "Invalid p_product_id")
		return
	}

	earningsMap := map[string]string{}
	if rq.Earnings != "" {
		err = json.Unmarshal([]byte(rq.Earnings), &earningsMap)
		if err != nil {
			sendErrorFmt(w, http.StatusBadRequest, "Failed to decode payload.p_earnings as JSON: %v", err)
			return
		}
	}
	earningsString := earningsMap["128"]

	claim := &model.Claim{
		Store: "Paddle",
		Qty:   1,

		Txn:      rq.Txn,
		FullName: rq.Name,
		Email:    rq.Email,
		Message:  rq.Message,

		Raw: string(body),

		OrderID:  rq.OrderID,
		Country:  rq.Country,
		Currency: rq.Currency,
		Coupon:   rq.Coupon,
	}

	claim.Price, err = decimal.NewFromString(rq.Price)
	if err != nil {
		sendErrorFmt(w, http.StatusBadRequest, "p_price is not a decimal number: %v", rq.Price)
		return
	}
	claim.SaleGross, err = decimal.NewFromString(rq.SaleGross)
	if err != nil {
		sendErrorFmt(w, http.StatusBadRequest, "p_sale_gross is not a decimal number: %v", rq.SaleGross)
		return
	}
	claim.SaleTax, err = decimal.NewFromString(rq.TaxAmount)
	if err != nil {
		sendErrorFmt(w, http.StatusBadRequest, "p_tax_amount is not a decimal number: %v", rq.TaxAmount)
		return
	}
	claim.ProcessorFee, err = decimal.NewFromString(rq.PaddleFee)
	if err != nil {
		sendErrorFmt(w, http.StatusBadRequest, "p_paddle_fee is not a decimal number: %v", rq.PaddleFee)
		return
	}
	claim.Earnings, err = decimal.NewFromString(earningsString)
	if err != nil {
		sendErrorFmt(w, http.StatusBadRequest, `p_earnings["128"] is not a decimal number: %v`, earningsString)
		return
	}
	claim.CouponSavings, err = decimal.NewFromString(rq.CouponSavings)
	if err != nil {
		sendErrorFmt(w, http.StatusBadRequest, "p_coupon_savings is not a decimal number: %v", rq.CouponSavings)
		return
	}

	var code string
	if false {
		code = "LR2A-X92VM-MGI6H-KT5XM-KRD52-FIDXF-CQG3F-PUPE3"
	} else {
		code, err = model.ClaimLicense(db, "LR", "2", licensecode.TypeIndividual, claim)
		if err != nil {
			sendError(w, err)
			return
		}
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "%s\n", code)
}
