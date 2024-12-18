package sliceH

func Distinct[T any](arr []T) (res []T) {
	m := map[any]struct{}{}
	for _, v := range arr {
		if _, ok := m[v]; !ok {
			res = append(res, v)
			m[v] = struct{}{}
		}
	}

	return
}
