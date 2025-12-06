package pagination

import (
	"net/url"
	"strconv"
)

// ===========================================
// HATEOAS LINK BUILDER
// ===========================================

// buildLinks generates HATEOAS navigation links
func buildLinks(req *AdvancedPaginationRequest, resp *AdvancedPaginationResponse) PaginationLinks {
	links := PaginationLinks{}

	if req.BaseURL == "" {
		return links
	}

	builder := newLinkBuilder(req)

	if req.Type == CursorPagination {
		return builder.buildCursorLinks(resp)
	}

	return builder.buildOffsetLinks(resp)
}

// linkBuilder helps construct pagination links
type linkBuilder struct {
	req *AdvancedPaginationRequest
}

func newLinkBuilder(req *AdvancedPaginationRequest) *linkBuilder {
	return &linkBuilder{req: req}
}

func (lb *linkBuilder) buildCursorLinks(resp *AdvancedPaginationResponse) PaginationLinks {
	links := PaginationLinks{
		Self: lb.buildURL(map[string]string{"cursor": lb.req.Cursor}),
	}

	if resp.HasNextPage && resp.EndCursor != "" {
		links.Next = lb.buildURL(map[string]string{"cursor": resp.EndCursor})
	}

	return links
}

func (lb *linkBuilder) buildOffsetLinks(resp *AdvancedPaginationResponse) PaginationLinks {
	links := PaginationLinks{
		Self:  lb.buildURL(map[string]string{"page": strconv.Itoa(lb.req.Page)}),
		First: lb.buildURL(map[string]string{"page": "1"}),
	}

	if resp.TotalPages > 0 {
		links.Last = lb.buildURL(map[string]string{"page": strconv.Itoa(resp.TotalPages)})
	}

	if resp.HasNextPage {
		links.Next = lb.buildURL(map[string]string{"page": strconv.Itoa(lb.req.Page + 1)})
	}

	if resp.HasPrevPage {
		links.Prev = lb.buildURL(map[string]string{"page": strconv.Itoa(lb.req.Page - 1)})
	}

	return links
}

func (lb *linkBuilder) buildURL(params map[string]string) string {
	u, err := url.Parse(lb.req.BaseURL)
	if err != nil {
		return lb.req.BaseURL
	}

	q := u.Query()

	// Preserve existing params
	lb.preserveParams(q)

	// Add page size
	q.Set("page_size", strconv.Itoa(lb.req.PageSize))

	// Add custom params
	for k, v := range params {
		if v != "" {
			q.Set(k, v)
		}
	}

	u.RawQuery = q.Encode()
	return u.String()
}

func (lb *linkBuilder) preserveParams(q url.Values) {
	if lb.req.SortRaw != "" {
		q.Set("sort", lb.req.SortRaw)
	}
	if lb.req.FilterRaw != "" {
		q.Set("filter", lb.req.FilterRaw)
	}
	if lb.req.Search != "" {
		q.Set("search", lb.req.Search)
	}
}
