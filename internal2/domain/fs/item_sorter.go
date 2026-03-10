package fs

type ItemSorter interface {
	Sort(items []Item) ([]Item, error)
}
