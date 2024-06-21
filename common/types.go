package common

type Pagination[T any] struct {
	Total   int64
	Count   int64
	Page    int64
	Data    []T
	HasNext bool
}
