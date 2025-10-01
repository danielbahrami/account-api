# Account API
This project implements a bank account API that allows users to create accounts, list their accounts, and transfer funds between accounts. The API is built using Go, PostgreSQL, and Docker Compose.

## Technology Choices
- **Go:** Strong type safety, fast compilation, robust standard library, and good support for concurrency.
- **PostgreSQL:** Reliably stores structured data like users, accounts, transactions, and supports relationships and constraints.
- **Docker Compose:** Simplifies development setup and provides a consistent environment for the API and PostgreSQL.

## Initial Data
The database initializes two demo users:
- John Doe (email: `johndoe@example.com`)
- Jane Doe (email: `janedoe@example.com`)

Both have the password: `password123`. Each user has two accounts with demo balances for testing.

The API does not provide user registration, as creating user accounts is outside the scope of this project.

## API Endpoints

All endpoints (except `/login`) require JWT authentication via the `Authorization: Bearer <token>` header.

### 1. `POST /login`
Authenticates a user and provides a JWT token. This token must be included in the `Authorization` header as a Bearer token for all other endpoints. Without a valid JWT, requests to protected endpoints such as `/accounts` and `/transfer` will return an error.

**Request body:**
```json
{
    "email": "johndoe@example.com",
    "password": "password123"
}
```

### 2. `GET /accounts`
Retrieves all accounts belonging to the authenticated user.

### 3. `POST /accounts`
Create a new account for the authenticated user.

**Request body:**
```json
{
    "account_name": "Savings"
}
```

### 4. `POST /transfer`
Transfer funds between accounts owned by the authenticated user.

**Request body:**
```json
{
    "from_account_number": "1111111111",
    "to_account_number": "2222222222",
    "amount": 100.00
}
```

## Code Overview

### main.go
Starts the application, sets up the database connection, and registers the HTTP routes.
### routes.go
Registers all API endpoints (`/login`, `/accounts`, `/transfer`) to their corresponding handlers.
### auth.go
Handles user authentication, including login and JWT generation.
### middleware.go
Validates JWT's on incoming requests and injects the authenticated user ID into the request context.
### accounts.go
Implements the logic for listing user accounts and creating new accounts.
### transfer.go
Implements the logic for transferring funds between accounts and updating balances.
### init.sql
Sets up the database tables and inserts demo users and accounts for testing.
### compose.yaml
Orchestrates the application and PostgreSQL for easy deployment.
