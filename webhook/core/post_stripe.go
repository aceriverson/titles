package titles

import (
	"errors"
	"strconv"
)

func (h *TitlesCore) CheckoutCompleted(userIDString, customer, session string) error {
	if userIDString == "" {
		return errors.New("userID is empty")
	}
	if customer == "" {
		return errors.New("customer is empty")
	}

	userID, err := strconv.ParseInt(userIDString, 10, 64)
	if err != nil {
		return errors.New("invalid userID")
	}

	// Check if the user exists
	_, err = h.DB.GetUserInternal(userID)
	if err != nil {
		return err
	}

	// Create new subscription
	h.DB.CreateSubscription(userID, customer, session)

	return nil
}

func (h *TitlesCore) SubscriptionCreated(customer, session string) error {
	if customer == "" {
		return errors.New("customer is empty")
	}

	h.DB.UpdateSubscription(customer, session, "pro")

	return nil
}

func (h *TitlesCore) SubscriptionDeleted(customer, session string) error {
	if customer == "" {
		return errors.New("customer is empty")
	}

	h.DB.UpdateSubscription(customer, session, "free")

	return nil
}
