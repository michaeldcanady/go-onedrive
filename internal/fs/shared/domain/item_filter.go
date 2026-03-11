package domain

type ItemFilterer interface {
	Filter([]Item) ([]Item, error)
}
