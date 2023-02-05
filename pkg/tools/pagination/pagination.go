package pagination

import (
	"math"
	"strconv"
)

type Pagination struct {
	Page  string
	Pages string
	Items string
}

func GetData(totalItems int, page, pageSize int32) Pagination {
	fPages := float64(totalItems) / float64(pageSize)
	pages := int(math.Ceil(fPages))

	return Pagination{
		Page:  strconv.Itoa(int(page)),
		Pages: strconv.Itoa(pages),
		Items: strconv.Itoa(totalItems),
	}
}
