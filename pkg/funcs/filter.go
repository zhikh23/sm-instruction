package funcs

type Predicate[T any] func(T) bool

func Filter[T any](collection []T, predicate Predicate[T]) []T {
	res := make([]T, 0, len(collection))
	for _, t := range collection {
		if predicate(t) {
			res = append(res, t)
		}
	}
	return res
}
