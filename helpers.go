package paginate

import (
	"net/url"
	"strconv"
)

// PageLink Returns URL for a given page index.
func (p *Paginator) PageLink(page int) string {
	link, _ := url.ParseRequestURI(p.Request().URL.String())
	values := link.Query()
	if page == 1 {
		values.Del(PageParam)
	} else {
		values.Set(PageParam, strconv.Itoa(page))
	}
	link.RawQuery = values.Encode()
	return link.String()
}

// PageLinkPrev Returns URL to the previous page.
func (p *Paginator) PageLinkPrev() (link string) {
	if p.HasPrev() {
		link = p.PageLink(p.Page() - 1)
	}
	return
}

// PageLinkNext Returns URL to the next page.
func (p *Paginator) PageLinkNext() (link string) {
	if p.HasNext() {
		link = p.PageLink(p.Page() + 1)
	}
	return
}

// PageLinkFirst Returns URL to the first page.
func (p *Paginator) PageLinkFirst() (link string) {
	return p.PageLink(1)
}

// PageLinkLast Returns URL to the last page.
func (p *Paginator) PageLinkLast() (link string) {
	return p.PageLink(p.PageNums())
}

// HasPrev Returns true if the current page has a predecessor.
func (p *Paginator) HasPrev() bool {
	return p.Page() > 1
}

// HasNext Returns true if the current page has a successor.
func (p *Paginator) HasNext() bool {
	return p.Page() < p.PageNums()
}

// IsActive Returns true if the given page index points to the current page.
func (p *Paginator) IsActive(page int) bool {
	return p.Page() == page
}

// getParam retrieve GET query param or POST form param
func (p *Paginator) getParam(param string) (val string) {
	val = p.request.URL.Query().Get(param)
	return
}

func (p *Paginator) SetLinks() {
	if p.HasPrev() {
		p.Prev = p.PageLinkPrev()
	}
	if p.HasNext() {
		p.Next = p.PageLinkNext()
	}
	if p.PageNums() > 0 {
		p.First = p.PageLinkFirst()
		p.Last = p.PageLinkLast()
	}
}
