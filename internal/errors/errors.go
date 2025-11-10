package errors

import "errors"

var (
	NoDataFound error = errors.New("no data found for provided IDs")
)
