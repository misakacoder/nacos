package collection

func Distinct[E comparable](slice []E) (result []E) {
	if slice != nil {
		dict := map[E]struct{}{}
		for _, v := range slice {
			if _, ok := dict[v]; !ok {
				result = append(result, v)
				dict[v] = struct{}{}
			}
		}
	}
	return
}
