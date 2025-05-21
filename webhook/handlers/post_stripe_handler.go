package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

var webhookSecret = os.Getenv("STRIPE_WEBHOOK_SECRET")

func (h *Handler) PostStripeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const MaxBodyBytes = int64(65536)
		r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
		payload, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		event := stripe.Event{}

		if err := json.Unmarshal(payload, &event); err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Webhook error while parsing basic request. %v\n", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		signatureHeader := r.Header.Get("Stripe-Signature")
		event, err = webhook.ConstructEvent(payload, signatureHeader, webhookSecret)
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠️  Webhook signature verification failed. %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Unmarshal the event data into an appropriate struct depending on its Type
		switch event.Type {
		case "checkout.session.completed":
			var session stripe.CheckoutSession
			if err := json.Unmarshal(event.Data.Raw, &session); err == nil {
				h.titles.CheckoutCompleted(session.ClientReferenceID, session.Customer.ID, session.Subscription.ID)
				h.titles.SubscriptionCreated(session.Customer.ID, session.Subscription.ID)
			}

		case "customer.subscription.created":
			var sub stripe.Subscription
			if err := json.Unmarshal(event.Data.Raw, &sub); err == nil {
				h.titles.SubscriptionCreated(sub.Customer.ID, sub.ID)
			}

		case "customer.subscription.deleted":
			var sub stripe.Subscription
			if err := json.Unmarshal(event.Data.Raw, &sub); err == nil {
				h.titles.SubscriptionDeleted(sub.Customer.ID, sub.ID)
			}

		default:
			fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
		}

		w.WriteHeader(http.StatusOK)
	}
}
