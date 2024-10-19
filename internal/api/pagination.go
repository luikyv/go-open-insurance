package api

import (
	"net/url"
	"strconv"
)

type Page[T any] struct {
	// Records are the records found for the page requested.
	Records []T
	// TotalRecords is the total number of records available.
	TotalRecords int
	// TotalPages is the total number of pages based on Size and TotalRecords.
	TotalPages int
	Pagination
}

type Pagination struct {
	// Number is the page number requested.
	Number int
	// Size is the page size requested.
	Size int
}

func NewPagination(p *PageNumber, ps *PageSize) Pagination {
	pagination := Pagination{
		Number: 1,
		Size:   25,
	}
	if p != nil {
		pagination.Number = int(*p)
	}

	if ps != nil && *ps <= 1000 {
		pagination.Size = int(*ps)
	}

	return pagination
}

// Paginate slices a list of records into a specific page of data based on the
// provided pagination parameters.
func Paginate[T any](records []T, pagination Pagination) Page[T] {
	numberOfRecords := len(records)
	page := Page[T]{
		TotalRecords: numberOfRecords,
		// Calculate the total number of pages using integer division.
		// Adding (pagination.Size - 1) ensures correct rounding up for partial pages.
		TotalPages: (numberOfRecords + pagination.Size - 1) / pagination.Size,
		Pagination: pagination,
	}

	// Subtracting 1 from the page number to convert it to a zero-based index,
	// as pages start at 1.
	start := (page.Number - 1) * page.Size
	if start >= numberOfRecords {
		return page
	}

	end := start + page.Size
	if end > numberOfRecords {
		end = numberOfRecords
	}

	page.Records = records[start:end]
	return page
}

// PaginatedLinks generates pagination links (self, first, prev, next, last)
// based on the current page information and the requested URL.
// T is a generic type parameter to make the function compatible with any Page type.
func PaginatedLinks[T any](requestedURL string, page Page[T]) Links {
	// Helper function to construct a URL with query parameters for pagination.
	buildURL := func(pageNumber int) string {
		u, _ := url.Parse(requestedURL)
		query := u.Query()
		query.Set("page", strconv.Itoa(pageNumber))
		query.Set("page-size", strconv.Itoa(page.Size))
		u.RawQuery = query.Encode()
		return u.String()
	}

	// Populate the Links struct.
	links := Links{
		Self: requestedURL,
	}

	// If the current page is not the first, generate the "first" and "previous"
	// links.
	if page.Number > 1 {
		first := buildURL(1)
		prev := buildURL(page.Number - 1)
		links.First = &first
		links.Prev = &prev
	}

	// If the current page is not the last, generate the "next" and "last" links.
	if page.Number < page.TotalPages {
		next := buildURL(page.Number + 1)
		last := buildURL(page.TotalPages)
		links.Next = &next
		links.Last = &last
	}

	return links
}
