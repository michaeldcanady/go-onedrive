package drive

type Drive struct {
	ID       string
	Name     string
	Type     DriveType
	Owner    string
	ReadOnly bool
}
