package common

// Pagination is a struct that holds the configuration for the Pagination
//
//   - Total is the total number of items
//   - Count is the number of items in the current page
//   - Page is the current page
//   - Data is the items in the current page
//   - HasNext is a boolean that indicates if there is a next page
type Pagination[T any] struct {
	Total   int64
	Count   int64
	Page    int64
	Data    []T
	HasNext bool
}
