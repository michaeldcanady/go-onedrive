package domain

type ItemSorter interface {
	Sort(items []Item) ([]Item, error)
}
