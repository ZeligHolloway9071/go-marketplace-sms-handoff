package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	client, err := NewSMSClient()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/orders/handoff", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "POST required", http.StatusMethodNotAllowed)
			return
		}
		var order Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, "invalid order", 400)
			return
		}
		seller, buyer, err := HandoffMessages(order)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		first, err := client.Send(order.SellerPhone, seller, order.ID+"-seller")
		if err != nil {
			http.Error(w, err.Error(), 502)
			return
		}
		second, err := client.Send(order.BuyerPhone, buyer, order.ID+"-buyer")
		if err != nil {
			http.Error(w, err.Error(), 502)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"seller_message_id": first, "buyer_message_id": second})
	})
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
