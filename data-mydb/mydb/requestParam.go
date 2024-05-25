package mydb

import "time"

type RequestParam struct {
	Start          time.Time
	End            time.Time
	IsTimeFieldSet bool
	From           int
	Size           int
	SearchField    string
	Keyword        string
	SortField      string
	SortType       string // desc / asc
}
