package filtering

type Filter interface {
	Filter(items any) error
}
