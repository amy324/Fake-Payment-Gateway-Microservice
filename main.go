package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ShiraazMoollatjie/goluhn"
)

func main() {
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/charge", handleCharge)

	fmt.Println("Fake Payment Gateway is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func handleCharge(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
		return
	}

	// Extract credit card number from form
	cardNumber := strings.TrimSpace(r.Form.Get("card_number"))

	// Validate credit card number using goluhn package
	if err := goluhn.Validate(cardNumber); err != nil {
		fmt.Println("Credit card number is invalid:", err)
		http.Error(w, "Invalid credit card number", http.StatusBadRequest)
		return
	}

	// Simulate charging the customer
	// You can customize this function to simulate different scenarios
	fmt.Println("Charge successful for credit card number:", cardNumber)
	fmt.Fprintf(w, "Charge successful for credit card number: %s\n", cardNumber)
}

func handleRefund(w http.ResponseWriter, r *http.Request) {
	// Simulate refunding a payment
	// You can customize this function to simulate different scenarios
	response := map[string]string{"message": "Refund successful"}
	jsonResponse(w, response)
}

func handlePaymentInfo(w http.ResponseWriter, r *http.Request) {
	// Simulate retrieving payment information
	// You can customize this function to return mock payment data
	response := map[string]interface{}{
		"customer_id":    "123456",
		"payment_amount": 50.00,
		"currency":       "USD",
	}
	jsonResponse(w, response)
}

func handleDefault(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not found", http.StatusNotFound)
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
