package main

import "testing"

func TestHandoffMessages(t *testing.T) {
	tests := []struct {
		name                  string
		order                 Order
		wantSeller, wantBuyer string
		wantErr               bool
	}{
		{"complete", Order{"A-17", "+15550001", "+15550002", "oak table"}, "Order A-17: handoff ready for oak table", "Order A-17 is ready for handoff", false},
		{"missing asset", Order{"A-17", "+1", "+2", ""}, "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seller, buyer, err := HandoffMessages(tt.order)
			if (err != nil) != tt.wantErr || seller != tt.wantSeller || buyer != tt.wantBuyer {
				t.Fatalf("got %q %q %v", seller, buyer, err)
			}
		})
	}
}
