

# Fake Payment Gateway

This project serves as a demonstration of card validation through a simple payment gateway. It does not include secure authentication mechanisms. To learn how to implement secure user authentication and sessions in Golang, please refer to one of my other projects, the full-scale [Ticket Raising Platform Backend](https://github.com/amy324/Ticket-Raising-Platform-Backend).

As a backend engineer, I have written this project as a purely backend application designed to handle payment transactions. As such, it does not include a user interface. To interact with the Fake Payment Gateway, you can test its endpoints using tools like Postman or curl. Alternatively, you can create your own frontend application to integrate with the provided backend services.

[Deployed Backend](https://payment-gateway-4gca.onrender.com/)


## Table of Contents

1. [Overview](#overview)
2. [Getting Started](#getting-started)
3. [Code Details](#code-details)
4. [Deployed App](#deployed-app)
5. [Usage](#usage)
6. [Examples](#examples)
7. [Contributions](#contributions)
7. [License](#license)

## Overview

The Fake Payment Gateway is designed to demonstrate how card validation works fundamentally without needing to sign up to third-party providers such as PayPal or Stripe to try writing this program yourself. It allows users to submit payment information, including name, card number, payment amount, and currency, via HTTP requests. The server validates the credit card number using the Luhn algorithm and inserts the payment information into a PostgreSQL database. Clients can also retrieve invoice details using the payment ID.

Please note that this program is solely for demonstration purposes. If you plan on trying out this code yourself, it should be for learning purposes only and should never be used to genuinely accept payments. In real-life scenarios, integrating with established payment processors like PayPal or Stripe is recommended for secure and reliable payment processing. You can also use third parties like Stripe or Paypal for testing and sandbox purposes with dummy data before attempting to implement any genuine payment processing solutions. As mentioned above, this program provides insight into payment gateway functionality without requiring integration  or registering with such services. The title of this repo is to emphaise that you should never use this program in its current state for actual payment processing!



## Getting Started

1. **Clone the repository:**

   ```bash
   git clone https://github.com/amy324/Fake-Payment-Gateway.git
   ```

2. **Install dependencies:**

   ```bash
   go mod tidy
   ```

3. **Set up the environment variables** by creating a `.env` file and specifying the following:

   ```
   DB_USER=your_database_user
   DB_PASSWORD=your_database_password
   DB_HOST=your_database_host
   DB_NAME=your_database_name
   ```

4. **Build and Run the Application:**

   - **Option 1: Build and Execute**

     ```bash
     go build
     ./payment.gateway.exe
     ```

   - **Option 2: Run Directly**

     ```bash
     go run .
     ```



## Code Details

This section provides insights into the functionality and structure of the code used in the Fake Payment Gateway project.

### Main Functionality

The `main` function serves as the entry point of the application. It establishes a connection to the PostgreSQL database, initializes a Gorilla Mux router to handle HTTP requests, and defines endpoints for submitting payment information (`/payment_form`) and retrieving invoice details (`/get_invoice/{payment_id}`). Additionally, it includes a route (`/`) to verify the server's live status.

### Payment Information Handling

The `handlePaymentInfoForm` function is responsible for processing payment information submitted via HTTP POST requests. It validates the credit card number using the Luhn algorithm, inserts the payment information into the PostgreSQL database, and returns a JSON response indicating the success or failure of the payment transaction.

### Invoice Retrieval

The `handleGetInvoice` function handles HTTP GET requests to retrieve invoice details based on the payment ID provided in the URL path parameters. It queries the PostgreSQL database for the specified payment ID and returns the invoice details in JSON format.

### Database Connectivity

The application interacts with a PostgreSQL database to store payment information securely. It utilizes the `database/sql` package along with the PostgreSQL driver (`github.com/lib/pq`) for database operations. However, it's important to note that the application does not store sensitive credit card details in the database. Instead, it stores essential payment details such as the name of the cardholder, payment amount, currency, and a boolean value indicating the validity status of the credit card number. If the credit card number passes validation, the corresponding entry in the database will have the validity status set to `true`. The application does not store entries for invalid credit card numbers, and the API response for such cases will be "Error validating credit card number". Environment variables specified in a `.env` file are used to provide database connection details such as the database username, password, host, and name, utilizing the `github.com/joho/godotenv` package.


### Dependencies

- **gorilla/mux**: A powerful HTTP router and URL matcher for building Go web servers.
- **ShiraazMoollatjie/goluhn**: A Go library for validating and generating Luhn numbers.
- **github.com/joho/godotenv**: A Go port of the Ruby dotenv library for loading environment variables from `.env` files.
- **github.com/lib/pq**: Go PostgreSQL driver for the database/sql package.

### PostgreSQL Database

The application uses a PostgreSQL database to store payment information securely. PostgreSQL is a robust, open-source relational database management system known for its reliability, extensibility, and compliance with SQL standards.



## Deployed App
A deployed version of this fake payment gateway microservice is available at [https://payment-gateway-4gca.onrender.com/](https://payment-gateway-4gca.onrender.com/)].Feel free to use it to test the API endpoints. You can verify if the service is live by visting the URL in your browser and seeing the text"The fake payment gateway microservice is live"

Feel free to use this URL to test the API endpoints; To use the fake payment gateway, you can send HTTP POST requests to `/payment_form` with payment information in the request body. You can also retrieve invoice details by sending HTTP GET requests to `/get_invoice/{payment_id}` where `{payment_id}` is the ID of the payment. See more details on how to do this [below](#usage)



### Note:

- **Do not use real credit card details:** Although this program does not store credit card details, it is strongly advised not to use any real card numbers. Use fake card numbers for testing purposes. For instructions on how to generate valid credit card numbers, visit the [goluhn documentation here](https://github.com/ShiraazMoollatjie/goluhn). Alternatively here are some examples to get you started:


Valid Card Numbers:

1. **Visa**: 4111111111111111
2. **Mastercard**: 5555555555554444
3. **American Express**: 378282246310005

Invalid Card Numbers:

1. **Invalid Length**: 12345678901234567890
2. **Invalid Checksum (Luhn Algorithm)**: 4111111111111112
3. **Invalid Prefix**: 1234567890123456

- **Preferably use fake names:** For privacy and security reasons, avoid using real names in the payment information submitted to the fake payment gateway.




## Usage

### Payment Form Endpoint

- **URL**: `/payment_form`
- **Method**: `POST`
- **Request Body**: JSON structure with the following fields:
  - `name`: Name of the cardholder
  - `card_number`: Credit card number
  - `payment_amount`: Payment amount
  - `currency`: Currency code

### Get Invoice Endpoint

- **URL**: `/get_invoice/{payment_id}`
- **Method**: `GET`
- **Path Parameter**:
  - `payment_id`: ID of the payment for which the invoice details are requested



## Examples

### Making a Payment:

**Endpoint:** `POST https://payment-gateway-4gca.onrender.com/payment_form`

**Request Body:**
```json
{
    "name": "John Doe",
    "card_number": "4111111111111111",
    "payment_amount": 20.99,
    "currency": "USD"
}
```

**Response:**
```json
{
    "message": "Payment Successful",
    "payment_id": 40
}
```

### Retrieving Invoice:

**Endpoint:** `GET https://payment-gateway-4gca.onrender.com/get_invoice/40`

**Response:**
```json
{
    "currency": "USD",
    "id": 40,
    "name": "John Doe",
    "payment_amount": 20.99,
    "valid": true
}
```

Feel free to try out your own examples in Postman or curl using the above strucure to interact with the fake payment gateway API. Remember, do not use real credit card details for testing purposes.

## Contributions

Contributions to this project are welcome! To contribute, follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix.
3. Make your changes and commit them to your branch.
4. Push your changes to your fork.
5. Submit a pull request to the main repository.


## License

This project is licensed under the MIT License

---

