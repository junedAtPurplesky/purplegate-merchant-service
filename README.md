# PurpleGate Merchant Service

A microservice for managing merchant operations in the PurpleGate platform.

## Overview

This service handles merchant-related operations including registration, management, and business logic for the PurpleGate ecosystem.

## Getting Started

### Prerequisites

- Go 1.24.4 or later
- Docker (optional)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/junedAtPurplesky/purplegate-merchant-service.git
cd purplegate-merchant-service
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the service:
```bash
go run cmd/main.go
```

## Project Structure

```
.
├── cmd/                 # Application entrypoints
├── internal/           # Private application code
├── proto/              # Protocol buffer definitions
├── configs/            # Configuration files
├── go.mod              # Go module file
└── README.md           # This file
```

## Development

### Building

```bash
go build -o bin/merchant-service cmd/main.go
```

### Testing

```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

