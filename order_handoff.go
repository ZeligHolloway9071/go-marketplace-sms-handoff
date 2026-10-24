package main

import "fmt"

type Order struct{ ID, SellerPhone, BuyerPhone, Asset string }

// HandoffMessages makes the business decision explicit: both parties get the same order state.
func HandoffMessages(o Order) (seller, buyer string, err error) {
	if o.ID == "" || o.SellerPhone == "" || o.BuyerPhone == "" || o.Asset == "" {
		return "", "", fmt.Errorf("order requires id, phones, and asset")
	}
	return fmt.Sprintf("Order %s: handoff ready for %s", o.ID, o.Asset), fmt.Sprintf("Order %s is ready for handoff", o.ID), nil
}
