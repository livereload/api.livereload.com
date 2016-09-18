package licensecode

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidFormat = errors.New("Invalid license code")
var ErrLegacyFormat = errors.New("Legacy license codes are not supported")

var ErrUnknownProduct = errors.New("Unknown product")
var ErrUnknownVersion = errors.New("Unknown version")
var ErrUnknownType = errors.New("Unknown license type")

type License struct {
	Product string
	Version string
	Type    string
	Code    string
}

func (lic License) String() string {
	return fmt.Sprintf("%s:%s:%s-%s", lic.Product, lic.Version, lic.Type, lic.Code)
}

func Parse(code string) (*License, error) {
	if len(code) < 5 {
		return nil, ErrInvalidFormat
	}

	code = strings.ToUpper(code)
	if code[2] == '-' {
		return nil, ErrLegacyFormat
	}
	if code[4] != '-' {
		return nil, ErrInvalidFormat
	}

	product := code[0:2]
	version := code[2:3]
	typ := code[3:4]

	if product != "LR" {
		return nil, ErrUnknownProduct
	}
	if !(version == "2" || version == "3") {
		return nil, ErrUnknownVersion
	}
	if !IsValidType(typ) {
		return nil, ErrUnknownType
	}

	return &License{product, version, typ, code}, nil
}
