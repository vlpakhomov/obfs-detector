package null

type Null[T any] struct {
	V   T
	Set bool
}

func New[T any](v T) Null[T] {
	return NewExplicit(v, true)
}

func NewPtr[T any](v T) *Null[T] {
	n := New(v)
	return &n
}

func NewExplicit[T any](v T, set bool) Null[T] {
	return Null[T]{
		V:   v,
		Set: set,
	}
}

func NewPtrExplicit[T any](v T, set bool) *Null[T] {
	n := NewExplicit(v, set)
	return &n
}

func NewFromPtr[T any](v *T) Null[T] {
	var value T
	if v != nil {
		value = *v
	}
	return NewExplicit(value, v != nil)
}

func NewPtrFromPtr[T any](v *T) *Null[T] {
	n := NewFromPtr(v)
	return &n
}

func (x Null[T]) ValuePtr() *T {
	if !x.Set {
		return nil
	}
	return &x.V
}
