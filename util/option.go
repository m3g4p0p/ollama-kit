package util

type Option[T any] func(*T)

func ApplyOptions[T any](v T, options []Option[T]) T {
	for _, opt := range options {
		opt(&v)
	}

	return v
}

func ApplyOptionsTo[T any](v *T, options []Option[T]) {
	for _, opt := range options {
		opt(v)
	}
}
