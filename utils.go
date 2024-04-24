package languageserver

import "fmt"

func (h *Handler) makeErr(while string, err error) error {
	wrappedE := fmt.Errorf(while+": %w", err)
	h.Logger.Error(wrappedE.Error())
	return wrappedE
}
