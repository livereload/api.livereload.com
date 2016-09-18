package model

import (
	"github.com/lib/pq"
)

// https://www.postgresql.org/docs/9.3/static/errcodes-appendix.html
const errCodeUniqueViolation pq.ErrorCode = "23505"
