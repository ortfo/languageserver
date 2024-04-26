package languageserver

import "fmt"

func (h *Handler) makeErr(while string, err error) error {
	wrappedE := fmt.Errorf(while+": %w", err)
	h.Logger.Error(wrappedE.Error())
	return wrappedE
}

func keys[K comparable, V any](m map[K]V) []K {
	result := make([]K, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	return result
}
