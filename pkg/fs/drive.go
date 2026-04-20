package fs

// Drive represents an available OneDrive/Cloud drive.
type Drive struct {
	ID       string
	Name     string
	Type     string
	Owner    string
	ReadOnly bool
}
