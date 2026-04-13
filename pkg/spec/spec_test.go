package spec

type mockSpec struct {
	result bool
}

func (s mockSpec) IsSatisfiedBy(candidate any) bool {
	return s.result
}
