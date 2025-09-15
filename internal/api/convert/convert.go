package convert

import "github.com/oapi-codegen/nullable"

// Ptr returns a pointer to a copy of v (for value -> *T).
func Ptr[T any](v T) *T { return &v }

// PtrIfNonZero returns a pointer to v if v != 0, otherwise nil.
func PtrIfNonZero[T ~int | ~int32 | ~int64](v T) *T {
	if v == 0 {
		return nil
	}
	return &v
}

// MapSlice maps []S -> []D using f.
func MapSlice[S any, D any](in []S, f func(S) D) []D {
	if in == nil {
		return nil
	}
	out := make([]D, len(in))
	for i := range in {
		out[i] = f(in[i])
	}
	return out
}

// PtrSlice returns *[]T with a copied backing array from a value slice.
func PtrSlice[T any](in []T) *[]T {
	if in == nil {
		return nil
	}
	out := make([]T, len(in))
	copy(out, in)
	return &out
}

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

// Deref returns the value pointed to by p, or def if p is nil.
func Deref[T any](p *T, def T) T {
	if p == nil {
		return def
	}
	return *p
}

// DerefSlice returns a copy of the slice pointed to by p, or nil if p is nil.
func DerefSlice[T any](p *[]T) []T {
	if p == nil {
		return nil
	}
	out := make([]T, len(*p))
	copy(out, *p)
	return out
}

// DerefSlicePtr maps *[]S -> []D
func DerefSlicePtr[S any, D any](p *[]S, f func(S) D) []D {
	if p == nil {
		return nil
	}
	out := make([]D, len(*p))
	for i, v := range *p {
		out[i] = f(v)
	}
	return out
}

// DerefMap pointer to value using mapper and default
func DerefMap[S any, D any](p *S, f func(S) D, def D) D {
	if p == nil {
		return def
	}
	return f(*p)
}

// DerefSlicePtrE maps *[]S -> []D using f(S) (D, error).
// Returns nil, nil if p is nil.
func DerefSlicePtrE[S any, D any](p *[]S, f func(S) (D, error)) ([]D, error) {
	if p == nil {
		return nil, nil
	}
	in := *p
	out := make([]D, len(in))
	for i := range in {
		d, err := f(in[i])
		if err != nil {
			return nil, err
		}
		out[i] = d
	}
	return out, nil
}

// DerefSlicePtrEP is the pointer-to-slice variant: *[]S -> *[]D.
func DerefSlicePtrEP[S any, D any](p *[]S, f func(S) (D, error)) (*[]D, error) {
	if p == nil {
		return nil, nil
	}
	in := *p
	out := make([]D, len(in))
	for i := range in {
		d, err := f(in[i])
		if err != nil {
			return nil, err
		}
		out[i] = d
	}
	return &out, nil
}
