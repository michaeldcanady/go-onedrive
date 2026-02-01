package sorting

type Direction int64

const (
	DirectionUnknown Direction = iota - 1
	DirectionAscending
	DirectionDescending
)

func (d Direction) String() string {
	str, ok := map[Direction]string{
		DirectionUnknown:    "unknown",
		DirectionAscending:  "ascending",
		DirectionDescending: "descending",
	}[d]
	if !ok {
		return DirectionUnknown.String()
	}
	return str
}
