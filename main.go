package main

import (
    "database/sql"
    "fmt"
    "net/http"
    "strconv"
    "strings"
    "github.com/joho/godotenv"
    "os"

    "github.com/ShiraazMoollatjie/goluhn"
    _ "github.com/lib/pq"
)

var db *sql.DB

func main() {
     // Load environment variables from .env file
     if err := godotenv.Load(); err != nil {
        fmt.Println("Error loading .env file:", err)
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
        fmt.Println("Error connecting to database:", err)
        return
    }
    defer db.Close()

    // Print a success message if the connection is successful
    fmt.Println("Connected to database successfully")

    http.HandleFunc("/", serveIndex)
    http.HandleFunc("/payment_form", handlePaymentInfoForm)

    fmt.Println("Fake Payment Gateway is running on port 8080...")
    http.ListenAndServe(":8080", nil)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "index.html")
}

func handlePaymentInfoForm(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse form data
    err := r.ParseForm()
    if err != nil {
        fmt.Println("Failed to parse form data:", err)
        http.Error(w, "Failed to parse form data", http.StatusInternalServerError)
        return
    }

    // Extract form data
    name := strings.TrimSpace(r.Form.Get("name"))
    cardNumber := strings.TrimSpace(r.Form.Get("card_number"))
    paymentAmountStr := strings.TrimSpace(r.Form.Get("payment_amount"))
    currency := strings.TrimSpace(r.Form.Get("currency"))

    fmt.Println("Form data:")
    fmt.Println("Name:", name)
    fmt.Println("Card Number:", cardNumber)
    fmt.Println("Payment Amount:", paymentAmountStr)
    fmt.Println("Currency:", currency)

    // Validate credit card number using Luhn algorithm
    if err := goluhn.Validate(cardNumber); err != nil {
        fmt.Println("Error validating credit card number:", err)
        http.Error(w, "Error validating credit card number", http.StatusBadRequest)
        return
    }

    // Convert payment amount to float64
    paymentAmount, err := strconv.ParseFloat(paymentAmountStr, 64)
    if err != nil {
        fmt.Println("Invalid payment amount:", err)
        http.Error(w, "Invalid payment amount", http.StatusBadRequest)
        return
    }

    // Insert data into the database
    _, err = db.Exec("INSERT INTO payment_info (name, payment_amount, currency, valid) VALUES ($1, $2, $3, $4)", name, paymentAmount, currency, true)
    if err != nil {
        fmt.Println("Failed to insert payment information into database:", err)
        http.Error(w, "Failed to insert payment information into database", http.StatusInternalServerError)
        return
    }

    // Redirect to the payment information page or show a success message
    http.Redirect(w, r, "/payment_success.html", http.StatusSeeOther)
}
