package state

type Repository interface {
	Load() (State, error)
	Save(state State) error
}
