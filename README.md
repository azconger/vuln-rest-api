# Vulnerable REST API

A deliberately vulnerable REST API for testing and demonstration purposes, inspired by the OWASP Benchmark project.

## Features
- Multiple endpoints with common web vulnerabilities (SQLi, weak auth, etc.)
- OpenAPI (Swagger) documentation and UI
- Example authentication system with JWT (intentionally weak)
- PostgreSQL backend with seed data

## Getting Started

### Running with Docker Compose (Recommended)

1. **Build and start all services:**
   ```sh
   docker-compose up --build
   ```
   This will:
   - Start a PostgreSQL database with seed data
   - Build and run the Go app
   - Automatically generate the OpenAPI docs

2. **Access the API and documentation:**
   - **Swagger UI:** [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
   - **OpenAPI JSON:** [http://localhost:8080/swagger/doc.json](http://localhost:8080/swagger/doc.json)
   - **API Base Path:** [http://localhost:8080/api/v1/](http://localhost:8080/api/v1/)

3. **Stop all services:**
   ```sh
   docker-compose down
   ```

### Running Locally (Manual)
- Make sure you have Go, PostgreSQL, and the `swag` CLI installed.
- Start PostgreSQL and load `init.sql`.
- Run `swag init -g main.go` to generate docs.
- Run the app: `go run main.go`

## Test Credentials
Use these credentials to test authentication endpoints:

| Username | Password   | Role      |
|----------|------------|-----------|
| admin    | admin123   | admin     |
| user     | user123    | user      |

## References
- [OWASP Benchmark](https://owasp.org/www-project-benchmark/)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)

## Warning
**This project is intentionally vulnerable. Do NOT deploy in production or expose to the public internet.** 