package convert

import "github.com/oapi-codegen/nullable"

// MapNullable maps a nullable.Nullable[T] to *U using f.
// Returns nil when the field is unspecified OR explicitly null.
// Propagates Get() errors from the nullable package.
func MapNullable[T, U any](n nullable.Nullable[T], f func(T) U) (*U, error) {
	if !n.IsSpecified() || n.IsNull() {
		return nil, nil
	}
	v, err := n.Get()
	if err != nil {
		return nil, err
	}
	u := f(v)
	return &u, nil
}

// MapNullableSlice maps a nullable.Nullable[[]S] to *[]D using f for elements.
// Returns nil when the field is unspecified OR explicitly null.
func MapNullableSlice[S, D any](n nullable.Nullable[[]S], f func(S) D) (*[]D, error) {
	if !n.IsSpecified() || n.IsNull() {
		return nil, nil
	}
	in, err := n.Get()
	if err != nil {
		return nil, err
	}
	out := make([]D, len(in))
	for i := range in {
		out[i] = f(in[i])
	}
	return &out, nil
}

// CopyPtr returns a new pointer with the same value, or nil if input is nil.
func CopyPtr[T any](p *T) *T {
	if p == nil {
		return nil
	}
	v := *p
	return &v
}

// MapPtr applies f to *S and returns *D (nil-safe).
func MapPtr[S, D any](p *S, f func(S) D) *D {
	if p == nil {
		return nil
	}
	v := f(*p)
	return &v
}

// MapPtrSlice maps *[]S -> *[]D (nil-safe, deep copy).
func MapPtrSlice[S, D any](p *[]S, f func(S) D) *[]D {
	if p == nil {
		return nil
	}
	in := *p
	out := make([]D, len(in))
	for i := range in {
		out[i] = f(in[i])
	}
	return &out
}
