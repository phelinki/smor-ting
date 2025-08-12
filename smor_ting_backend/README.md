# Smor-Ting Backend

A robust, production-ready backend API for the Smor-Ting platform built with Go, Fiber, and SQLite/PostgreSQL.

## Features

- 🔐 **Secure Authentication**: JWT-based authentication with bcrypt password hashing
- 🗄️ **Database Support**: SQLite (development/testing) and PostgreSQL (production)
- 📝 **Structured Logging**: Comprehensive logging with Zap logger
- 🔒 **Security**: CORS configuration, input validation, and secure defaults
- 📚 **API Documentation**: Swagger/OpenAPI documentation
- 🧪 **Testing Mode**: In-memory database for development and testing
- ⚙️ **Configuration**: Environment-based configuration management
- 🚀 **Production Ready**: Error handling, health checks, and observability

## Quick Start

### Prerequisites

- Go 1.21 or higher
- SQLite (for development) or PostgreSQL (for production)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd smor-ting-backend
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables (see Configuration section)

4. Run the application:
```bash
go run cmd/main.go
```

The server will start on `http://localhost:8080` by default.

## Configuration

The application uses environment variables for configuration. Create a `.env` file in the root directory:

### Development Environment
```env
# Server Configuration
PORT=8080
HOST=0.0.0.0
ENV=development

# Database Configuration (SQLite for development)
DB_DRIVER=sqlite
DB_IN_MEMORY=true
DB_NAME=smor_ting.db

# Authentication
JWT_SECRET=YOUR_JWT_SECRET_MIN_32_CHARS
JWT_EXPIRATION=24h
BCRYPT_COST=12

# CORS (Development - allows all origins)
CORS_ALLOW_ORIGINS=*
CORS_ALLOW_HEADERS=Origin,Content-Type,Accept,Authorization
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOW_CREDENTIALS=true

# Logging
LOG_LEVEL=debug
LOG_FORMAT=console
LOG_OUTPUT=stdout
```

### Production Environment
```env
# Server Configuration
PORT=8080
HOST=0.0.0.0
ENV=production

# Database Configuration (PostgreSQL for production)
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=your_username
DB_PASSWORD=YOUR_MONGODB_PASSWORD
DB_NAME=smor_ting
DB_SSL_MODE=require
DB_IN_MEMORY=false

# Authentication (Use a strong secret in production!)
JWT_SECRET=YOUR_JWT_SECRET_MIN_32_CHARS
JWT_EXPIRATION=24h
BCRYPT_COST=12

# CORS (Production - restrict origins)
CORS_ALLOW_ORIGINS=https://smor-ting.com,https://www.smor-ting.com
CORS_ALLOW_HEADERS=Origin,Content-Type,Accept,Authorization
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOW_CREDENTIALS=true

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
LOG_OUTPUT=stdout
```

## Testing Mode

The application includes a testing mode that uses an in-memory SQLite database. This is perfect for:

- **Development**: Fast startup and no persistent data
- **Testing**: Isolated test environments
- **CI/CD**: Automated testing without external dependencies

### Enabling Testing Mode

Set the following environment variables:
```env
DB_DRIVER=sqlite
DB_IN_MEMORY=true
```

**Note**: In-memory databases are for development/testing only. For production, use PostgreSQL with persistent storage.

## API Documentation

### Authentication Endpoints

#### Register User
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response:**
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Login User
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Validate Token
```http
POST /api/v1/auth/validate
Authorization: Bearer <token>
```

**Response:**
```json
{
  "id": 1,
  "email": "user@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Health Check

#### Get Health Status
```http
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "smor-ting-backend",
  "version": "1.0.0"
}
```

## Error Handling

The application provides structured error responses:

```json
{
  "error": "Validation failed",
  "message": "Email is required"
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (invalid credentials)
- `409` - Conflict (user already exists)
- `500` - Internal Server Error

## Logging

The application uses structured logging with different levels:

- **DEBUG**: Detailed debugging information
- **INFO**: General application information
- **WARN**: Warning messages
- **ERROR**: Error messages with stack traces

### Log Format

**Development (Console):**
```
2024-01-01T12:00:00.000Z INFO User registered successfully email=user@example.com user_id=1
```

**Production (JSON):**
```json
{
  "level": "info",
  "timestamp": "2024-01-01T12:00:00.000Z",
  "message": "User registered successfully",
  "email": "user@example.com",
  "user_id": 1
}
```

## Security

### Authentication
- JWT tokens with configurable expiration
- bcrypt password hashing with configurable cost
- Secure token generation with random IDs

### CORS
- Configurable origins for production
- Proper header and method restrictions
- Credential support for authenticated requests

### Database
- SQL injection protection through parameterized queries
- Connection pooling for performance
- SSL/TLS support for PostgreSQL

## Development

### Project Structure
```
smor_ting_backend/
├── cmd/
│   └── main.go              # Application entry point
├── configs/
│   └── config.go            # Configuration management
├── internal/
│   └── auth/                # Authentication package
│       ├── service.go       # Business logic
│       └── handler.go       # HTTP handlers
├── pkg/
│   ├── database/            # Database package
│   ├── logger/              # Logging package
│   └── middleware/          # Middleware package
├── migrations/              # Database migrations
├── docs/                    # API documentation
└── scripts/                 # Utility scripts
```

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o bin/server cmd/main.go
```

## Deployment

### Docker
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o server cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
CMD ["./server"]
```

### Environment Variables for Production
- Set `ENV=production`
- Use strong `JWT_SECRET`
- Configure PostgreSQL database
- Set appropriate CORS origins
- Use JSON logging format

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

[Add your license information here] 