package validation

import "errors"

// All returns a policy that requires all provided policies to pass.
func All[T any](policies ...Policy[T]) Policy[T] {
	return allPolicy[T]{policies: policies}
}

type allPolicy[T any] struct {
	policies []Policy[T]
}

func (p allPolicy[T]) Evaluate(candidate T) error {
	var errs []error
	for _, policy := range p.policies {
		if policy == nil {
			continue
		}
		if err := policy.Evaluate(candidate); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// Any returns a policy that requires at least one of the provided policies to pass.
func Any[T any](policies ...Policy[T]) Policy[T] {
	return anyPolicy[T]{policies: policies}
}

type anyPolicy[T any] struct {
	policies []Policy[T]
}

func (p anyPolicy[T]) Evaluate(candidate T) error {
	var errs []error
	for _, policy := range p.policies {
		if policy == nil {
			continue
		}
		err := policy.Evaluate(candidate)
		if err == nil {
			return nil
		}
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

// Each returns a policy that applies another policy to each element of a slice.
func Each[T any, E any](getter func(T) []E, policy Policy[E]) Policy[T] {
	return EachPolicy[T, E]{getter: getter, policy: policy}
}

type EachPolicy[T any, E any] struct {
	getter func(T) []E
	policy Policy[E]
}

func (p EachPolicy[T, E]) Evaluate(candidate T) error {
	elements := p.getter(candidate)
	var errs []error
	for _, element := range elements {
		if err := p.policy.Evaluate(element); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}
