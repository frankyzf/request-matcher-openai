package export

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"request-matcher-openai/data-mydb/mydb"
)

func GetTimeParam(c *gin.Context) (mydb.RequestParam, error) {
	var err error
	// now := time.Now().In(loc)
	// year, month, day := now.Date()
	// start := time.Date(year, month, day-7, 0, 0, 0, 0, loc)
	// end := time.Date(year, month, day, 23, 59, 59, 999, loc)
	start := time.Unix(0, 0)
	end, _ := time.Parse("2006-01-02", "2049-12-31")
	from := 0
	size := 20
	searchField := ""
	keyword := ""
	sortField := ""
	sortType := "" // desc / asc
	bIsTimeSet := false
	if c.Query("start") != "" {
		start, err = time.Parse("2006-01-02", c.Query("start"))
		start = mydb.GetBeginOfDay(start)
		bIsTimeSet = true
	}
	if err == nil {
		if c.Query("end") != "" {
			end, err = time.Parse("2006-01-02", c.Query("end"))
			end = mydb.GetEndOfDay(end)
			bIsTimeSet = true
		}
	}
	if err == nil {
		if c.Query("from") != "" {
			from, err = strconv.Atoi(c.Query("from"))
		}
	}
	if err == nil {
		if c.Query("size") != "" {
			size, err = strconv.Atoi(c.Query("size"))
		}
	}
	if err == nil {
		if c.Query("search_field") != "" {
			searchField = c.Query("search_field")
		}
	}
	if err == nil {
		if c.Query("keyword") != "" {
			keyword = c.Query("keyword")
		}
	}
	if err == nil {
		if c.Query("sort_field") != "" {
			sortField = c.Query("sort_field")
		}
	}
	if err == nil {
		if c.Query("sort_type") != "" {
			sortType = c.Query("sort_type")
		}
	}
	//no search key word, reset search field
	if keyword == "" {
		searchField = ""
	}
	if sortField == "" {
		sortType = ""
	}

	data := mydb.RequestParam{
		Start:          start,
		End:            end,
		IsTimeFieldSet: bIsTimeSet,
		From:           from,
		Size:           size,
		SearchField:    searchField,
		Keyword:        keyword,
		SortField:      sortField,
		SortType:       sortType,
	}
	return data, err
}
