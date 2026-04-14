package sorting

// Direction represents the sort order: ascending, descending, or unspecified.
type Direction int

const (
	// DirectionAscending identifies an A-Z, smallest-to-largest, or oldest-to-newest sort order.
	DirectionAscending Direction = iota
	// DirectionDescending identifies a Z-A, largest-to-smallest, or newest-to-oldest sort order.
	DirectionDescending
)
