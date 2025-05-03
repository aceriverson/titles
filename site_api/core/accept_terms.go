package core

func (h *Core) AcceptTerms(userID int64) error {
	if err := h.DB.AcceptTerms(userID); err != nil {
		return err
	}

	return nil
}
