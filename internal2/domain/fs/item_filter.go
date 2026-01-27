package fs

type ItemFilterer interface {
	Filter([]Item) ([]Item, error)
}
