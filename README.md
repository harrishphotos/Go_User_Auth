# Go User Authentication System

A robust, scalable user authentication and management service built with Go, featuring secure authentication, email verification, and authorization capabilities.

## 🚀 Features

- **User Authentication**

  - Registration with email verification
  - Secure login with PASETO tokens
  - Password reset functionality
  - Token refresh mechanism
  - Argon2id password hashing

- **Security**

  - JWT/PASETO token-based authentication
  - Role-based access control
  - Protection against common security vulnerabilities

- **Infrastructure**
  - Kafka integration for reliable event processing
  - PostgreSQL for data persistence
  - Docker containerization
  - Scalable architecture

## 🛠️ Tech Stack

- **Backend**: Go with Fiber web framework
- **Database**: PostgreSQL with GORM ORM
- **Security**: PASETO tokens, Argon2id password hashing
- **Messaging**: Kafka for asynchronous processing
- **Containerization**: Docker and Docker Compose

## 📋 Prerequisites

- Go 1.21+
- Docker and Docker Compose
- PostgreSQL (or use the Docker setup)
- Kafka (or use the Docker setup)

## 🚦 Getting Started

### Using Docker

1. Clone the repository

```bash
git clone https://github.com/your-username/Go_User_Auth.git
cd Go_User_Auth
```

2. Create a `.env` file based on provided example config

```bash
# Example .env content
PORT=3000
DB_HOST=postgres
DB_PORT=5432
DB_USER=harish
DB_PASSWORD=postgres
DB_NAME=book_crud
KAFKA_BROKERS=kafka:9094
TOKEN_SYMMETRIC_KEY=your-secret-key-here
```

3. Start the application using Docker Compose

```bash
docker-compose up -d
```

### Manual Setup

1. Install Go dependencies

```bash
go mod tidy
```

2. Set up PostgreSQL database
3. Configure Kafka
4. Run the application

```bash
go run main.go
```

## 🔑 API Endpoints

### Authentication

- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - User login
- `GET /api/auth/verify-email` - Email verification
- `POST /api/auth/refresh` - Refresh access token
- `POST /api/auth/forgot-password` - Request password reset
- `POST /api/auth/reset-password` - Reset password
- `POST /api/auth/logout` - Logout (requires authentication)

### Books API

- Book management endpoints (CRUD operations)

### Store API

- Store management endpoints

## 🔒 Authentication Flow

### 1. User Registration Flow

**Request:**

```json
POST /api/auth/register
{
  "username": "johndoe",
  "email": "john.doe@example.com",
  "password": "SecurePassword123!",
  "full_name": "John Doe",
  "phone": "+1234567890"
}
```

**Response:**

```json
{
  "status": "success",
  "message": "User registered successfully. Please check your email to verify your account.",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "johndoe",
    "email": "john.doe@example.com",
    "created_at": "2023-04-28T14:45:30Z"
  }
}
```

**Process:**

1. User submits registration information
2. System validates data and creates a new user account (unverified)
3. Verification token is generated and stored
4. Email with verification link is sent via Kafka message queue
5. User receives verification email with a link to confirm their account

### 2. Email Verification Flow

**Request:**

```
GET /api/auth/verify-email?token=abcdef123456
```

**Response:**

```json
{
  "status": "success",
  "message": "Email verified successfully. You can now log in."
}
```

**Process:**

1. User clicks verification link in email
2. System validates the token
3. User account is marked as verified
4. User is redirected to login page

### 3. Login Flow

**Request:**

```json
POST /api/auth/login
{
  "email": "john.doe@example.com",
  "password": "SecurePassword123!"
}
```

**Response:**

```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "access_token": "v4.local.payload_data_here",
    "refresh_token": "v4.local.payload_data_here",
    "user": {
      "user_id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "johndoe",
      "email": "john.doe@example.com"
    },
    "expires_at": "2023-04-28T15:45:30Z"
  }
}
```

**Process:**

1. User submits login credentials
2. System validates credentials
3. If valid, access and refresh tokens are generated
   - Access token: short-lived (15-60 minutes)
   - Refresh token: longer-lived (7-30 days)
4. Tokens are returned to the client

### 4. Token Refresh Flow

**Request:**

```json
POST /api/auth/refresh
{
  "refresh_token": "v4.local.payload_data_here"
}
```

**Response:**

```json
{
  "status": "success",
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "v4.local.new_payload_data_here",
    "refresh_token": "v4.local.new_payload_data_here",
    "expires_at": "2023-04-28T16:45:30Z"
  }
}
```

**Process:**

1. Client submits refresh token
2. System validates the refresh token
3. If valid, new access and refresh tokens are generated
4. New tokens are returned to the client

### 5. Password Reset Flow

**Step 1: Request Password Reset**

**Request:**

```json
POST /api/auth/forgot-password
{
  "email": "john.doe@example.com"
}
```

**Response:**

```json
{
  "status": "success",
  "message": "If the email exists, a password reset link has been sent."
}
```

**Process:**

1. User requests password reset with email
2. System generates a reset token (valid for limited time, e.g., 15 minutes)
3. Reset email with token is sent via Kafka message queue
4. User receives email with password reset link

**Step 2: Reset Password**

**Request:**

```json
POST /api/auth/reset-password
{
  "token": "abcdef123456",
  "new_password": "NewSecurePassword456!"
}
```

**Response:**

```json
{
  "status": "success",
  "message": "Password has been reset successfully. You can now log in with your new password."
}
```

**Process:**

1. User submits reset token and new password
2. System validates the token
3. If valid, password is updated with new Argon2id hash
4. User can log in with the new password

### 6. Logout Flow

**Request:**

```json
POST /api/auth/logout
Authorization: Bearer v4.local.payload_data_here
```

**Response:**

```json
{
  "status": "success",
  "message": "Logged out successfully"
}
```

**Process:**

1. Client submits request with valid access token
2. System invalidates the refresh token
3. Client deletes tokens from local storage

## 🐳 Docker Services

- **app**: Go application service
- **postgres**: PostgreSQL database
- **kafka**: Kafka message broker

##licence

##
