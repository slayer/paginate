package paginate

import (
	"math"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
)

const (
	defaultPerPage      = 25
	defaultPageParam    = "page"
	defaultPerPageParam = "per_page"
)

var (
	// PageParam is a name of URL query param "page"
	PageParam = defaultPageParam
	// PerPageParam is a name of URL query param "per_page"
	PerPageParam = defaultPerPageParam
)

// Totaler is an interface to calculate total number of items via `Total()` function
type Totaler interface {
	// Total returns total number of items
	Total() int
}

// Paginator within the state of a http request.
type Paginator struct {
	request  *http.Request
	PerPage  int `json:"per_page"`
	LastPage int `json:"last_page,omitempty"`

	// Links
	First string `json:"first,omitempty"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Last  string `json:"last,omitempty"`

	totalItems int64
	pageRange  []int
	pageNums   int
	page       int
}

func (p *Paginator) Request() *http.Request {
	return p.request
}

// PageNums Returns the total number of pages.
// it works in "page"
func (p *Paginator) PageNums() int {
	if p.pageNums != 0 {
		return p.pageNums
	}

	pageNums := math.Ceil(float64(p.totalItems) / float64(p.PerPage))
	if p.LastPage > 0 {
		pageNums = math.Min(pageNums, float64(p.LastPage))
	}
	p.pageNums = int(pageNums)
	return p.pageNums
}

// Nums Returns the total number of items (e.g. from doing SQL count).
func (p *Paginator) Nums() int64 {
	return p.totalItems
}

// SetTotal Sets the total number of items.
func (p *Paginator) SetTotal(nums interface{}) {
	if totaler, ok := nums.(Totaler); ok {
		p.totalItems = int64(totaler.Total())
	} else if query, ok := nums.(*gorm.DB); ok {
		query.Count(&p.totalItems)
	} else {
		p.totalItems, _ = toInt64(nums)
	}
	LastPage := math.Ceil(float64(p.totalItems) / float64(p.PerPage))
	p.LastPage = int(LastPage)
}

// SetPerPage set limit items per page
func (p *Paginator) SetPerPage(n int) {
	p.PerPage = n
	if p.PerPage > 0 {
		return
	}
	p.PerPage, _ = strconv.Atoi(p.getParam(PerPageParam))

	if p.PerPage > 0 {
		return
	}
	p.PerPage = defaultPerPage
}

// Page Returns the current page.
func (p *Paginator) Page() int {
	if p.page != 0 {
		return p.page
	}

	p.page, _ = strconv.Atoi(p.getParam(PageParam))
	if p.page > p.PageNums() {
		p.page = p.PageNums()
	}
	if p.page <= 0 {
		p.page = 1
	}
	return p.page
}

// Offset Returns the current OFFSET for SQL query
func (p *Paginator) Offset() int {
	return (p.Page() - 1) * p.PerPage
}

// Limit Returns the LIMIT for SQL query
func (p *Paginator) Limit() int {
	return p.PerPage
}

// HasPages Returns true if there is more than one page.
func (p *Paginator) HasPages() bool {
	return p.PageNums() > 1
}

func (p *Paginator) Paginate(query *gorm.DB) *gorm.DB {
	return query.Offset(p.Offset()).Limit(p.PerPage)
}

// NewPaginator Instantiates a paginator struct for the current http request.
// total can be a number (uints, ints) or gorm.DB to auto query COUNT()
func NewPaginator(r *http.Request, total interface{}) *Paginator {
	p := Paginator{
		request: r,
	}
	p.SetPerPage(0) // 0 - use "per_page" param or default value
	p.SetTotal(total)
	p.SetLinks()
	return &p
}
