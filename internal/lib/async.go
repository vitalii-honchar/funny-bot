package lib

func Async[V any](f func(chan V)) <-chan V {
	c := make(chan V, 1)
	go func() {
		defer close(c)
		f(c)
	}()
	return c
}
