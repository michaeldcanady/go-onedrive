package domain

type Repository interface {
	Load() (State, error)
	Save(state State) error
}
