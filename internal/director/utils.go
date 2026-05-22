package director

// WithEndEvent invokes fn and calls onSuccess or onError based on the result.
func WithEndEvent(fn func() error, onSuccess func(), onError func(error)) error {
	err := fn()
	if err != nil {
		onError(err)
	} else {
		onSuccess()
	}
	return err
}

// WithUnaryEndEvent invokes fn and calls onSuccess or onError based on the result.
func WithUnaryEndEvent[T any](fn func() (T, error), onSuccess func(T), onError func(error)) (T, error) {
	t, err := fn()
	if err != nil {
		onError(err)
	} else {
		onSuccess(t)
	}
	return t, err
}
