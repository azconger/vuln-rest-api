# Vulnerable Go REST API

This project implements a deliberately vulnerable REST API in Go, designed for educational and security testing purposes. It mirrors the vulnerabilities present in the OWASP Benchmark application but is implemented in Go instead of Java.

## ⚠️ Security Warning

**IMPORTANT**: This API is intentionally vulnerable and should NEVER be deployed in a production environment. It is designed for:
- Security research and education
- Penetration testing practice
- Understanding common web vulnerabilities
- Testing security tools and scanners

## Features

- RESTful API implemented in Go
- Comprehensive OpenAPI/Swagger documentation
- Authentication system with multiple vulnerable implementations:
  - Basic JSON authentication
  - OAuth 2.0 Resource Owner Password Credentials (ROPC) flow
  - JWT-based authorization with known vulnerabilities
- Intentionally vulnerable endpoints demonstrating:
  - SQL Injection
  - Cross-Site Scripting (XSS)
  - Command Injection
  - Path Traversal
  - Insecure Deserialization
  - Authentication Bypass
  - And more...

## Prerequisites

- Go 1.21 or higher
- Docker (optional, for containerized deployment)
- Make (optional, for using Makefile commands)

## Project Structure

```
.
├── api/            # OpenAPI/Swagger specifications
├── cmd/            # Application entry points
├── internal/       # Private application code
│   ├── handlers/   # HTTP request handlers
│   ├── models/     # Data models
│   ├── database/   # Database interactions
│   ├── auth/       # Authentication implementations
│   └── middleware/ # HTTP middleware
├── pkg/           # Public libraries
├── tests/         # Test cases and security tests
└── docs/          # Additional documentation
```

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/vuln-rest-api.git
   cd vuln-rest-api
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

The API will be available at `http://localhost:8080`

## Authentication

The API provides two authentication methods, both intentionally vulnerable:

### 1. JSON Authentication
```http
POST /api/v1/auth/login
Content-Type: application/json

{
    "username": "user",
    "password": "password"
}
```

### 2. OAuth 2.0 ROPC
```http
POST /api/v1/auth/token
Content-Type: application/x-www-form-urlencoded

grant_type=password&username=user&password=password
```

Both methods return a JWT token that should be included in subsequent requests:
```http
Authorization: Bearer <jwt_token>
```

**Note**: The authentication system is intentionally vulnerable with:
- Hardcoded credentials (username: `user`, password: `password`)
- Weak JWT implementation
- Insufficient token validation
- Predictable token generation

## API Documentation

The API is fully documented using OpenAPI/Swagger. You can access the documentation at:
- Swagger UI: `http://localhost:8080/swagger-ui`
- OpenAPI JSON: `http://localhost:8080/swagger.json`

## Vulnerabilities

This API intentionally includes the following vulnerabilities:

1. **Authentication & Authorization**
   - Hardcoded credentials
   - Weak JWT implementation
   - Insufficient token validation
   - Predictable token generation
   - Session fixation vulnerabilities
   - Missing rate limiting
   - Weak password policies

2. **SQL Injection**
   - Unsanitized user input in database queries
   - Multiple injection points with varying difficulty levels

3. **Cross-Site Scripting (XSS)**
   - Reflected XSS
   - Stored XSS
   - DOM-based XSS

4. **Command Injection**
   - Unsafe command execution
   - Path traversal vulnerabilities

5. **Insecure Deserialization**
   - Unsafe object deserialization
   - XML external entity (XXE) vulnerabilities

6. **Security Misconfiguration**
   - Debug endpoints
   - Verbose error messages
   - Default credentials

## Testing

The project includes a comprehensive test suite that verifies the presence of vulnerabilities:

```bash
go test ./tests/...
```

## Contributing

Contributions are welcome! Please read our [Contributing Guidelines](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by the [OWASP Benchmark Project](https://owasp.org/www-project-benchmark/)
- Built with [Go](https://golang.org/)
- API documentation powered by [Swagger](https://swagger.io/) 