# 📦 Golang Hybrid API

A flexible microservices architecture in Go that supports both **monolithic** and **distributed** deployment patterns. Each service can run independently as a standalone server or be composed together in a single application, providing maximum deployment flexibility without code changes.

## 🏗️ Architecture Overview

This hybrid API model allows you to:
- **Single Server Mode**: Run all services in one application process
- **Distributed Mode**: Deploy each service as an independent server
- **Mixed Mode**: Run some services together while others run separately
- **Seamless Communication**: Services can communicate in-process or via HTTP without code changes

## 🚀 Available Services

| Service | Description | Port | Status |
|---------|-------------|------|--------|
| `calculator-svc` | Mathematical operations (add, subtract, multiply, divide) | 8080 | ✅ Implemented |
| `auth-svc` | Authentication and authorization | TBD | 🚧 In Progress |
| `gateway-svc` | API Gateway and routing | TBD | 🚧 In Progress |
| `pdf-svc` | PDF generation and manipulation | TBD | 🚧 In Progress |
| `temperature-svc` | Temperature conversion utilities | TBD | 🚧 In Progress |

## 🔧 Quick Start

### Prerequisites
- Go 1.24.2 or higher
- Git

### Installation
```bash
git clone https://github.com/growly-group/golang-hybrid-api.git
cd golang-hybrid-api
go mod download
```

### Environment Configuration
Create a `.env` file in the root directory (copy from `.env.example`):
```bash
# Copy the example configuration
cp .env.example .env
```

**Example .env content:**
```env
# Required: Comma-separated list of services to run
TARGET_SERVICES=calculator-svc

# Optional: Service URLs for HTTP communication
CALCULADOR_SERVICE_URL=http://localhost:8080
AUTH_SERVICE_URL=http://localhost:8081
GATEWAY_SERVICE_URL=http://localhost:8082
```

## 🎯 Deployment Patterns

### Pattern 1: Single Service (Microservice)
Run only the calculator service:
```bash
# Set environment
export TARGET_SERVICES=calculator-svc
# Or on Windows
$env:TARGET_SERVICES="calculator-svc"

# Run
go run main.go
```

### Pattern 2: Multiple Services (Distributed Monolith)
Run multiple services in one process:
```bash
export TARGET_SERVICES=calculator-svc,auth-svc,gateway-svc
go run main.go
```

### Pattern 3: Distributed Deployment
Deploy each service separately:

**Terminal 1 - Calculator Service:**
```bash
export TARGET_SERVICES=calculator-svc
go run main.go
```

**Terminal 2 - Auth Service:**
```bash
export TARGET_SERVICES=auth-svc
go run main.go
```

### Pattern 4: Docker Deployment

**Single Container (Monolithic):**
```bash
# Build and run all services in one container
docker-compose --profile monolith up --build
```

**Multiple Containers (Distributed):**
```bash
# Build and run each service in separate containers
docker-compose --profile distributed up --build
```

**Custom Docker Build:**
```bash
# Build custom image
docker build -t golang-hybrid-api .

# Run calculator service only
docker run -e TARGET_SERVICES=calculator-svc -p 8080:8080 golang-hybrid-api

# Run multiple services
docker run -e TARGET_SERVICES=calculator-svc,auth-svc -p 8080:8080 golang-hybrid-api
```

## 🔌 Service Communication

Services can communicate in two modes without code changes:

### 1. In-Process Communication (Direct SDK)
When services run in the same process:
```go
// Automatic direct function calls
sdk := calculatorsvc.NewCalculatorSdk("direct")
result, err := sdk.Calculate(calculatorsvc.CalculationRequest{
    A: 10,
    B: 5,
    Operation: "add",
})
```

### 2. HTTP Communication (Network SDK)
When services run in separate processes:
```go
// Automatic HTTP calls
sdk := calculatorsvc.NewCalculatorSdk("http")
result, err := sdk.Calculate(calculatorsvc.CalculationRequest{
    A: 10,
    B: 5,
    Operation: "add",
})
```

## 📝 API Documentation

### Calculator Service

**Endpoint:** `POST /calculator`

**Request Body:**
```json
{
  "a": 10,
  "b": 5,
  "operation": "add"
}
```

**Supported Operations:**
- `add` - Addition
- `subtract` - Subtraction
- `multiply` - Multiplication
- `divide` - Division

**Response:**
```json
{
  "result": 15
}
```

**Error Response:**
```json
{
  "error": "Division by zero is not allowed"
}
```

**Example using curl:**
```bash
curl -X POST http://localhost:8080/calculator \
  -H "Content-Type: application/json" \
  -d '{"a": 10, "b": 5, "operation": "add"}'
```

## 🏢 Production Deployment Strategies

### Strategy 1: Full Microservices
- Deploy each service as a separate container/instance
- Use service discovery (Consul, etcd) or API Gateway
- Enable horizontal scaling per service
- Ideal for: Large teams, different scaling requirements

### Strategy 2: Service Groups
- Group related services together
- Deploy calculator + temperature services together
- Deploy auth + gateway services together
- Ideal for: Medium teams, logical service boundaries

### Strategy 3: Monolithic Start
- Deploy all services in one application initially
- Split services as the application grows
- Zero downtime migration path
- Ideal for: Startups, rapid prototyping

## 🔧 Development

### Adding a New Service

1. **Create service directory:**
```bash
mkdir -p src/my-service
```

2. **Implement required files:**
```
src/my-service/
├── entrypoint.go    # HTTP server setup
├── sdk.go          # SDK for in-process/HTTP calls
├── types.go        # Request/Response types
└── methods.go      # Business logic
```

3. **Register in main.go:**
```go
import myservice "github.com/growly-group/golang-hybrid-api/src/my-service"

entrypoints := map[string]entrypointFunc{
    // ... existing services
    "my-service": myservice.Entrypoint,
}
```

### Service Template

**entrypoint.go:**
```go
package myservice

import (
    "log"
    "github.com/gin-gonic/gin"
)

func Entrypoint() {
    r := gin.Default()
    r.POST("/my-endpoint", myHandler)
    
    log.Println("Starting my-service on :8081")
    if err := r.Run(":8081"); err != nil {
        log.Fatalf("Failed to run server: %v", err)
    }
}
```

## 🌟 Benefits

### Deployment Flexibility
- **Start Simple**: Begin with a monolith, scale to microservices
- **Cost Effective**: Reduce infrastructure costs in early stages
- **Team Independence**: Teams can deploy services independently when ready

### Development Experience
- **Code Reuse**: Same service code works in both modes
- **Easy Testing**: Test services in isolation or together
- **Gradual Migration**: Move to microservices incrementally

### Operational Benefits
- **Simplified Monitoring**: Fewer moving parts in monolithic mode
- **Service Resilience**: HTTP communication includes retry logic
- **Performance**: In-process calls eliminate network overhead

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-service`
3. Follow the service template for new services
4. Add tests for your service
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔮 Roadmap

- [ ] Complete implementation of all services
- [ ] Add service discovery mechanism
- [ ] Implement circuit breaker pattern
- [ ] Add comprehensive monitoring and metrics
- [ ] Create Kubernetes deployment manifests
- [ ] Add automated tests and CI/CD pipeline
