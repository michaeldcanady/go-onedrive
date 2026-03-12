package sorting

// Direction represents the sort order: ascending, descending, or unspecified.
type Direction int

const (
	// DirectionUnknown identifies an unspecified sort order.
	DirectionUnknown Direction = iota - 1
	// DirectionAscending identifies an A-Z, smallest-to-largest, or oldest-to-newest sort order.
	DirectionAscending
	// DirectionDescending identifies a Z-A, largest-to-smallest, or newest-to-oldest sort order.
	DirectionDescending
)

// String returns the string representation of the Direction.
func (d Direction) String() string {
	switch d {
	case DirectionAscending:
		return "ascending"
	case DirectionDescending:
		return "descending"
	default:
		return "unknown"
	}
}
