package shared

// Service combines the Reader, Writer, and Manager interfaces to provide a full-featured filesystem service.
type Service interface {
	Namer
	Reader
	Writer
	Manager
}

type Namer interface {
	Name() string
}
