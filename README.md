# Auth Service

A secure authentication service with user management and comprehensive logging built with Go, PostgreSQL, and Kafka.

## Features

- User registration and login
- Password hashing and validation
- Role-based access control
- Comprehensive logging with Kafka integration
- Database persistence with PostgreSQL
- Containerized for easy deployment

## Prerequisites

- Docker and Docker Compose
- Git

## Quick Start

1. Clone the repository:
   ```
   git clone https://github.com/KfcEnjoyer/auth-service.git
   cd auth-service
   ```

2. Start the services:
   ```
   docker-compose up -d
   ```

3. The API will be available at `http://localhost:8080`

## API Endpoints

### Create User
```
POST /create
Content-Type: application/json

{
  "username": "testuser",
  "password": {
    "plain": "secure_password@123"
  },
  "role": "regular"
}
```

### Login
```
POST /login
Content-Type: application/json

{
  "username": "testuser",
  "password": {
    "plain": "secure_password@123"
  }
}
```

### Home (Welcome Page)
```
GET /home
```

### Health Check
```
GET /health
```

## Architecture

```
.
├── cmd/
│   ├── auth/         # Authentication service
│   └── logservice/   # Log processing service
├── configs/          # Configuration files
├── internal/         # Internal packages
│   ├── database/     # Database interaction
│   ├── services/     # Core business logic
│   └── user/         # User domain models
├── kafka/            # Kafka producers and consumers
├── pkg/              # Utility packages
│   ├── logger/       # Logging utilities
│   └── validator/    # Input validation
├── docker-compose.yml
└── Dockerfile
```

## Configuration

The service uses a YAML configuration file located at `configs/config.yaml` with the following structure:

```yaml
database:
  host: "postgres"  # Use "postgres" when running with Docker, "localhost" for local development
  port: 5432
  user: "postgres"
  password: "lolkek12"
  dbname: "users"
```

You can override configuration settings using environment variables:

- `CONFIG_PATH`: Path to the configuration file
- `KAFKA_BROKERS`: Comma-separated list of Kafka brokers
- `KAFKA_TOPIC`: Kafka topic for logs (default: "auth-logs")
- `KAFKA_GROUP_ID`: Consumer group ID (default: "log-consumer-group")
- `PORT`: HTTP server port (default: "8080")

## Monitoring and Logging

All authentication events and system logs are stored in both PostgreSQL and Kafka.

To view the logs:

1. Connect to the PostgreSQL database:
   ```
   docker-compose exec postgres psql -U postgres -d users
   ```

2. Execute a query:
   ```sql
   SELECT * FROM logs ORDER BY timestamp DESC LIMIT 10;
   ```

## Development

### Without Docker

1. Install dependencies:
   - Go 1.24+
   - PostgreSQL
   - Kafka

2. Set up the database:
   ```sql
   CREATE DATABASE users;
   ```
   
3. Run the SQL in `init-db.sql` to set up tables

4. Update `configs/config.yaml` with your local connection details

5. Start the services:
   ```
   go run cmd/auth/main.go
   go run cmd/logservice/main.go
   ```

### Testing

```
go test ./...
```

## Security Considerations

- Passwords are hashed using bcrypt
- Input validation is performed on all endpoints
- All authentication events are logged
- Suspicious login attempts are flagged

## License

MIT
