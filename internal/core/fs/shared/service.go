package shared

// Service combines the Reader, Writer, and Manager interfaces to provide a full-featured filesystem service.
type Service interface {
	Reader
	Writer
	Manager
}
