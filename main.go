package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/ShiraazMoollatjie/goluhn"
	_ "github.com/lib/pq"
)

var db *sql.DB

// PaymentInfo struct represents the payment information received from the client
type PaymentInfo struct {
	Name          string  `json:"name"`
	CardNumber    string  `json:"card_number"`
	PaymentAmount float64 `json:"payment_amount"`
	Currency      string  `json:"currency"`
	PaymentID     int     `json:"payment_id"`
}

// Invoice struct represents the invoice details retrieved from the database
type Invoice struct {
	ID                          int    `json:"id"`
	Name                        string `json:"name"`
	PaymentAmountInSmallestUnit int    `json:"payment_amount_in_smallest_unit"`
	Currency                    string `json:"currency"`
	Valid                       bool   `json:"valid"`
}

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file:", err)
		return
	}

	// Retrieve database connection details from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	// Construct the database connection string
	dbURI := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=require", dbUser, dbPassword, dbHost, dbName)

	// Establish a connection to the PostgreSQL database
	var err error
	db, err = sql.Open("postgres", dbURI)
	if err != nil {
		log.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	// Print a success message if the connection is successful
	log.Println("Connected to database successfully")

	// Create a new Gorilla Mux router
	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/payment_form", handlePaymentInfoForm).Methods("POST")
	router.HandleFunc("/get_invoice/{payment_id}", handleGetInvoice).Methods("GET")

	log.Println("Fake Payment Gateway is running on port 8080...")
	http.ListenAndServe(":8080", router)
}

// PaymentRequestBody represents the JSON structure expected in the payment form request body
type PaymentRequestBody struct {
	Name          string  `json:"name"`
	CardNumber    string  `json:"card_number"`
	PaymentAmount float64 `json:"payment_amount"`
	Currency      string  `json:"currency"`
}

func handlePaymentInfoForm(w http.ResponseWriter, r *http.Request) {
	// Decode JSON data from the request body
	var paymentInfo PaymentInfo
	err := json.NewDecoder(r.Body).Decode(&paymentInfo)
	if err != nil {
		log.Println("Failed to decode JSON data:", err)
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	log.Println("Parsed payment info:", paymentInfo)

	// Validate credit card number using Luhn algorithm
	if err := goluhn.Validate(paymentInfo.CardNumber); err != nil {
		log.Println("Error validating credit card number:", err)
		http.Error(w, "Error validating credit card number", http.StatusBadRequest)
		return
	}

	log.Println("Credit card number validated successfully.")

	// Convert payment amount to the smallest unit (e.g., cents)
	paymentAmountInCents := paymentInfo.PaymentAmount * 100

	// Insert data into the database after validation
	_, err = db.Exec("INSERT INTO payment_info (name, payment_amount, currency, valid) VALUES ($1, $2, $3, $4)", paymentInfo.Name, paymentAmountInCents, paymentInfo.Currency, true)
	if err != nil {
		log.Println("Failed to insert payment information into database:", err)
		http.Error(w, "Failed to insert payment information into database", http.StatusInternalServerError)
		return
	}

	// Retrieve the ID of the inserted payment info based on name, currency, and payment amount
	row := db.QueryRow("SELECT id FROM payment_info WHERE name = $1 AND payment_amount = $2 AND currency = $3", paymentInfo.Name, paymentAmountInCents, paymentInfo.Currency)
	var paymentID int
	err = row.Scan(&paymentID)
	if err != nil {
		log.Println("Failed to retrieve payment ID:", err)
		http.Error(w, "Failed to retrieve payment ID", http.StatusInternalServerError)
		return
	}

	// Update the paymentInfo struct with the retrieved ID
	paymentInfo.PaymentID = paymentID

	log.Println("Payment information inserted into database successfully with ID:", paymentID)

	// Send a JSON response indicating success
	response := map[string]interface{}{
		"message":    "Payment Successful",
		"payment_id": paymentID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleGetInvoice(w http.ResponseWriter, r *http.Request) {
	// Extract the payment ID from the URL path parameters
	vars := mux.Vars(r)
	paymentIDStr := vars["payment_id"]
	if paymentIDStr == "" {
		http.Error(w, "Payment ID not found in URL", http.StatusBadRequest)
		return
	}

	// Convert the payment ID to an integer
	paymentID, err := strconv.Atoi(paymentIDStr)
	if err != nil {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	// Query the database to retrieve invoice details for the specific payment ID
	row := db.QueryRow("SELECT id, name, payment_amount, currency, valid FROM payment_info WHERE id = $1", paymentID)

	// Process the query result
	var invoiceDetails Invoice
	var paymentAmountInCents int
	if err := row.Scan(&invoiceDetails.ID, &invoiceDetails.Name, &paymentAmountInCents, &invoiceDetails.Currency, &invoiceDetails.Valid); err != nil {
		http.Error(w, "Failed to retrieve invoice details", http.StatusInternalServerError)
		return
	}

	// Convert payment amount from smallest unit (e.g., cents) to decimal
	paymentAmount := float64(paymentAmountInCents) / 100

	// Send a JSON response with the invoice details, including the payment amount in decimal format
	response := map[string]interface{}{
		"id":            invoiceDetails.ID,
		"name":          invoiceDetails.Name,
		"payment_amount": paymentAmount,
		"currency":      invoiceDetails.Currency,
		"valid":         invoiceDetails.Valid,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
