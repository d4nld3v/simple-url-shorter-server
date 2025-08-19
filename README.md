# 🔗 URL Shortener Go

A fast and secure URL shortening service built with Go, designed with clean architecture and security best practices.

## ✨ Features

- 🚀 **High performance** - Built with native Go without heavy frameworks
- 🛡️ **Robust security** - Protection against SSRF, XSS, SQL injection and brute force attacks
- 📊 **Rate limiting** - Request rate control per IP with automatic memory cleanup
- 🔒 **Strict validation** - Complete URL validation and prevention of private network access
- 📈 **Click tracking** - Click counter per shortened URL
- 🗄️ **SQLite database** - Persistent storage with connection pooling
- 🔐 **Security headers** - HTTP security headers implemented
- 📝 **Standardized responses** - API with consistent JSON format and clear error codes
- 📊 **Structured logging** - Complete request traceability with JSON logs

## 🏗️ Architecture

The project follows **Clean Architecture** principles:

```
cmd/api/v1/              # Application entry point
├── main.go              # Server configuration and startup

internal/                # Internal application code
├── config/              # Configuration and database connection
│   ├── config.go        # Configuration loading
│   └── db.go            # Singleton pattern for DB with connection pooling
├── handler/             # HTTP controllers (Presentation Layer)
│   ├── url.go           # Handlers for URL endpoints
│   └── response.go      # Response structures and helpers
├── services/            # Business logic (Business Layer)
│   ├── url.go           # URL shortening services
│   └── validator.go     # URL validation and security
├── repository/          # Data access (Data Layer)
│   └── url.go           # CRUD operations with database
└── server/              # HTTP server configuration
    └── server.go        # Route and middleware setup

pkg/                     # Reusable packages
├── logger/              # Structured logging system
│   └── logger.go        # JSON logger with multiple levels
└── middleware/          # Reusable middlewares
    ├── api.go           # Security headers
    ├── logging.go       # HTTP request logging
    └── ratelimit.go     # Rate limiting with token bucket algorithm

docs/                    # Documentation
├── LOGGING.md           # Logging system documentation
└── PATH_PARAMETERS.md   # Parameter documentation
```

## 🔒 Security Features

### **Protection against attacks:**

- **SSRF (Server-Side Request Forgery)** - Blocking private IPs and metadata servers
- **SQL Injection** - Prepared statements with placeholders
- **XSS** - Input validation and sanitization
- **Rate Limiting** - Request limits per IP with token bucket algorithm
- **Hash Collision** - SHA-256 + timestamp + random bytes instead of MD5

### **Implemented validations:**

- Malformed or dangerous URLs
- Private network access (127.0.0.1, 192.168.x.x, etc.)
- Control and potentially dangerous characters
- Length and data format limits
- Timeouts to prevent DoS attacks

## 📊 Logging System

### **Structured JSON Logging:**

```json
{
  "time": "2025-08-19T10:30:00.123Z",
  "level": "INFO",
  "msg": "Request completed successfully",
  "method": "POST",
  "path": "/shorten",
  "client_ip": "192.168.1.100",
  "status_code": 201,
  "duration_ms": 45,
  "request_id": "20250819103000-abc123"
}
```

### **Log Levels:**

- **DEBUG** - Detailed information for development
- **INFO** - General operational information
- **WARN** - Warning conditions that don't affect operation
- **ERROR** - Error conditions that affect operation
- **FATAL** - Critical errors that terminate the application

### **Automatic Metrics:**

- Response times in milliseconds
- Status codes for health monitoring
- Request rates per IP and endpoint
- Error rates classified by type

## 🚀 Quick Start

### **Prerequisites**

- Go 1.19 or higher
- SQLite (included in project)

### **Installation**

```bash
# Clone the repository
git clone https://github.com/your-username/url-shortener-go.git
cd url-shortener-go

# Install dependencies
go mod download

# Run the application
go run cmd/api/v1/main.go
```

The server will be available at `http://localhost:8080`

## 📚 API Reference

### **Shorten URL**

```http
POST /shorten
Content-Type: application/json

{
  "url": "https://example.com/very/long/url"
}
```

**Successful response (201 Created):**

```json
{
  "data": {
    "original_url": "https://example.com/very/long/url",
    "shortened_id": "abc123xy",
    "shortened_url": "http://localhost:8080/short/abc123xy",
    "created_at": "2025-08-19T10:30:00Z"
  },
  "message": "URL shortened successfully",
  "timestamp": "2025-08-19T10:30:00Z"
}
```

### **Redirect URL**

```http
GET /short/{shortID}
```

**Response:** HTTP 302 redirect to the original URL

### **Error Responses**

All error responses follow the standard format:

```json
{
  "error": {
    "code": "INVALID_URL",
    "message": "Invalid URL",
    "details": "only http and https schemes are allowed"
  },
  "timestamp": "2025-08-19T10:30:00Z",
  "path": "/shorten",
  "method": "POST"
}
```

### **Error Codes**

| Code                  | Description                 |
| --------------------- | --------------------------- |
| `INVALID_INPUT`       | Invalid input data          |
| `INVALID_URL`         | Incorrectly formatted URL   |
| `URL_TOO_LONG`        | URL exceeds character limit |
| `URL_NOT_FOUND`       | Short URL not found         |
| `RATE_LIMIT_EXCEEDED` | Request limit exceeded      |
| `METHOD_NOT_ALLOWED`  | HTTP method not allowed     |
| `INTERNAL_ERROR`      | Internal server error       |

## ⚙️ Configuration

The service is configured via environment variables:

```bash
# Server port
ADDR=:8080

# Rate limiting
RATE_LIMIT=100          # Requests per minute per IP
BURST_LIMIT=10          # Maximum burst allowed

# Database
DB_SOURCE=./url_shortener.db

# Logging
LOG_LEVEL=INFO          # DEBUG, INFO, WARN, ERROR, FATAL
LOG_FILE=logs/app.log   # Log file path (optional)
LOG_FORMAT=json         # Log format (json, text)
```

## 🧪 Testing

```bash
# Run tests
go test ./...

# Tests with coverage
go test -cover ./...

# Integration tests
go test -tags=integration ./...
```

## 📊 Usage Examples

### **Shorten a URL**

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com"}'
```

### **Use the shortened URL**

```bash
curl -L http://localhost:8080/short/abc123xy
```

## � Monitoring & Logs

### **Log Analysis**

```bash
# Find errors in the last 24 hours
grep '"level":"ERROR"' logs/app.log | tail -100

# Rate limiting events
grep '"status_code":429' logs/app.log

# Performance metrics (slow requests >1000ms)
grep '"duration_ms"' logs/app.log | awk -F'"duration_ms":' '{print $2}' | awk -F',' '{if($1>1000) print}' | wc -l
```

### **Health Monitoring**

- Error rate monitoring
- Response time percentiles
- Request rate per endpoint
- Database connection health
- Rate limiting effectiveness

## 🔧 Development

### **Commit Structure**

- `feat:` new features
- `fix:` bug fixes
- `docs:` documentation changes
- `refactor:` code refactoring
- `test:` add or modify tests

### **Contributing**

1. Fork the project
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📈 Performance

- **Response time**: < 50ms for shortening
- **Redirects**: < 10ms
- **Rate limit**: 100 requests/min per IP by default
- **Database**: Connection pooling with maximum 25 connections
- **Logging**: Structured JSON logs with request correlation

## 🛠️ Technologies Used

- **[Go](https://golang.org/)** - Programming language
- **[SQLite](https://sqlite.org/)** - Embedded database
- **[net/http](https://pkg.go.dev/net/http)** - Native HTTP server
- **[slog](https://pkg.go.dev/log/slog)** - Structured logging
- **Clean Architecture** - Layer separation and responsibilities

## 📝 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## 👥 Authors

- **Daniel Enriquez** - [@d4nld3v](https://github.com/d4nld3v)

## 🙏 Acknowledgments

- Inspired by OWASP security best practices
- Based on Clean Architecture principles
- Implementation of standard Go design patterns
- Security guidelines from the Go community
