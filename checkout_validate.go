package acp

import (
	"errors"
	"fmt"
)

// Validate ensures CheckoutSessionCreateRequest satisfies required schema constraints.
func (r CheckoutSessionCreateRequest) Validate() error {
	if len(r.Items) == 0 {
		return errors.New("items must contain at least one entry")
	}
	for i, item := range r.Items {
		if item.ID == "" {
			return fmt.Errorf("items[%d]: id is required", i)
		}
		if item.Quantity <= 0 {
			return fmt.Errorf("items[%d]: quantity must be positive", i)
		}
	}
	if r.Buyer != nil {
		if r.Buyer.FirstName == "" || r.Buyer.LastName == "" || string(r.Buyer.Email) == "" {
			return errors.New("buyer requires first_name, last_name, and email")
		}
	}
	return nil
}

// Validate ensures CheckoutSessionUpdateRequest maintains schema constraints.
func (r CheckoutSessionUpdateRequest) Validate() error {
	if r.Items != nil {
		for i, item := range *r.Items {
			if item.ID == "" {
				return fmt.Errorf("items[%d]: id is required", i)
			}
			if item.Quantity <= 0 {
				return fmt.Errorf("items[%d]: quantity must be positive", i)
			}
		}
	}
	if r.Buyer != nil {
		if r.Buyer.FirstName == "" || r.Buyer.LastName == "" || string(r.Buyer.Email) == "" {
			return errors.New("buyer requires first_name, last_name, and email")
		}
	}
	return nil
}

// Validate ensures CheckoutSessionCompleteRequest satisfies payment requirements.
func (r CheckoutSessionCompleteRequest) Validate() error {
	if r.PaymentData.Token == "" {
		return errors.New("payment_data.token is required")
	}
	if r.PaymentData.Provider == "" {
		return errors.New("payment_data.provider is required")
	}
	return nil
}
