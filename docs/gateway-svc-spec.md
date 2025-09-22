# API Gateway Service (gateway-svc) Specification

## Overview
The API Gateway Service serves as the central entry point for all client requests in the hybrid API ecosystem. It provides routing, load balancing, authentication middleware, rate limiting, and request/response transformation across all microservices.

## Port Assignment
- **Default Port**: 8082
- **Health Check**: 8082/health
- **Metrics**: 8082/metrics

## Core Features

### 1. Request Routing
- Dynamic service discovery
- Path-based routing
- Header-based routing
- Load balancing algorithms
- Circuit breaker pattern

### 2. Authentication & Authorization
- JWT token validation
- API key authentication
- Rate limiting per user/API key
- Permission-based access control

### 3. Request/Response Processing
- Request/response transformation
- Header manipulation
- CORS handling
- Request/response logging

### 4. Monitoring & Analytics
- Request metrics collection
- Performance monitoring
- Error rate tracking
- API usage analytics

## API Endpoints

### Gateway Management

#### GET /gateway/health
Health check endpoint for the gateway service.

**Response (200):**
```json
{
  "status": "healthy",
  "timestamp": "2025-09-22T10:00:00Z",
  "services": {
    "auth-svc": "healthy",
    "calculator-svc": "healthy",
    "pdf-svc": "healthy",
    "temperature-svc": "healthy"
  },
  "version": "1.0.0"
}
```

#### GET /gateway/metrics
Prometheus-compatible metrics endpoint.

#### GET /gateway/routes
List all configured routes (admin only).

**Response (200):**
```json
{
  "routes": [
    {
      "id": "auth-routes",
      "path": "/api/v1/auth/*",
      "service": "auth-svc",
      "target": "http://localhost:8081",
      "methods": ["GET", "POST", "PUT", "DELETE"],
      "middleware": ["cors", "rate-limit"],
      "enabled": true
    },
    {
      "id": "calculator-routes",
      "path": "/api/v1/calculator/*",
      "service": "calculator-svc",
      "target": "http://localhost:8080",
      "methods": ["POST"],
      "middleware": ["cors", "auth", "rate-limit"],
      "enabled": true
    }
  ]
}
```

### Service Proxy Endpoints

#### /api/v1/auth/*
Proxy all authentication-related requests to auth-svc.

**Example:**
- `POST /api/v1/auth/login` → `POST http://auth-svc:8081/auth/login`
- `GET /api/v1/auth/profile` → `GET http://auth-svc:8081/auth/profile`

#### /api/v1/calculator/*
Proxy all calculator requests to calculator-svc.

**Example:**
- `POST /api/v1/calculator/calculate` → `POST http://calculator-svc:8080/calculator`

#### /api/v1/pdf/*
Proxy all PDF-related requests to pdf-svc.

#### /api/v1/temperature/*
Proxy all temperature conversion requests to temperature-svc.

### Administrative Endpoints

#### POST /gateway/routes
Add new route configuration (admin only).

**Request Body:**
```json
{
  "id": "new-service-routes",
  "path": "/api/v1/newservice/*",
  "service": "new-svc",
  "target": "http://localhost:8085",
  "methods": ["GET", "POST"],
  "middleware": ["cors", "auth"],
  "healthCheck": "/health",
  "timeout": 30,
  "retries": 3
}
```

#### PUT /gateway/routes/:id
Update existing route configuration.

#### DELETE /gateway/routes/:id
Remove route configuration.

#### POST /gateway/reload
Reload gateway configuration without restart.

## Routing Configuration

### Route Definition
```go
type Route struct {
    ID          string            `json:"id" yaml:"id"`
    Path        string            `json:"path" yaml:"path"`
    Service     string            `json:"service" yaml:"service"`
    Target      string            `json:"target" yaml:"target"`
    Methods     []string          `json:"methods" yaml:"methods"`
    Middleware  []string          `json:"middleware" yaml:"middleware"`
    Headers     map[string]string `json:"headers" yaml:"headers"`
    HealthCheck string            `json:"healthCheck" yaml:"healthCheck"`
    Timeout     int               `json:"timeout" yaml:"timeout"`
    Retries     int               `json:"retries" yaml:"retries"`
    Enabled     bool              `json:"enabled" yaml:"enabled"`
    LoadBalancer LoadBalancerConfig `json:"loadBalancer" yaml:"loadBalancer"`
}

type LoadBalancerConfig struct {
    Algorithm string   `json:"algorithm" yaml:"algorithm"` // round-robin, least-connections, random
    Targets   []string `json:"targets" yaml:"targets"`
    HealthChecks bool  `json:"healthChecks" yaml:"healthChecks"`
}
```

### Configuration File (gateway.yaml)
```yaml
server:
  port: 8082
  readTimeout: 30s
  writeTimeout: 30s
  maxHeaderBytes: 1048576

cors:
  allowOrigins: ["*"]
  allowMethods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
  allowHeaders: ["Content-Type", "Authorization"]
  maxAge: 86400

rateLimit:
  requests: 1000
  window: "1h"
  storage: "memory" # memory or redis

routes:
  - id: "auth-routes"
    path: "/api/v1/auth/*"
    service: "auth-svc"
    target: "http://localhost:8081"
    methods: ["GET", "POST", "PUT", "DELETE"]
    middleware: ["cors", "rate-limit"]
    healthCheck: "/health"
    timeout: 30
    retries: 3
    enabled: true

  - id: "calculator-routes"
    path: "/api/v1/calculator/*"
    service: "calculator-svc"
    target: "http://localhost:8080"
    methods: ["POST"]
    middleware: ["cors", "auth", "rate-limit"]
    healthCheck: "/health"
    timeout: 15
    retries: 3
    enabled: true

middleware:
  auth:
    endpoint: "http://localhost:8081/auth/verify"
    timeout: 5
    cacheTokens: true
    cacheTTL: "5m"

  rateLimit:
    defaultLimit: 100
    window: "1m"
    storage: "memory"
    keyFunc: "ip" # ip, user, api-key

serviceDiscovery:
  enabled: false
  provider: "consul" # consul, etcd, kubernetes
  config:
    address: "localhost:8500"
    interval: "30s"
```

## Middleware System

### Authentication Middleware
```go
type AuthMiddleware struct {
    AuthServiceURL string
    CacheEnabled   bool
    CacheTTL       time.Duration
    TokenCache     map[string]*CachedToken
}

type CachedToken struct {
    User      *User
    ExpiresAt time.Time
}
```

### Rate Limiting Middleware
```go
type RateLimitConfig struct {
    Requests   int           `json:"requests"`
    Window     time.Duration `json:"window"`
    KeyFunc    string        `json:"keyFunc"` // ip, user, api-key
    Storage    string        `json:"storage"` // memory, redis
    SkipPaths  []string      `json:"skipPaths"`
}
```

### CORS Middleware
```go
type CORSConfig struct {
    AllowOrigins     []string `json:"allowOrigins"`
    AllowMethods     []string `json:"allowMethods"`
    AllowHeaders     []string `json:"allowHeaders"`
    ExposeHeaders    []string `json:"exposeHeaders"`
    AllowCredentials bool     `json:"allowCredentials"`
    MaxAge           int      `json:"maxAge"`
}
```

### Circuit Breaker Middleware
```go
type CircuitBreakerConfig struct {
    MaxRequests      uint32        `json:"maxRequests"`
    Interval         time.Duration `json:"interval"`
    Timeout          time.Duration `json:"timeout"`
    FailureThreshold uint32        `json:"failureThreshold"`
}
```

## Load Balancing

### Supported Algorithms
1. **Round Robin**: Distributes requests evenly across all healthy targets
2. **Least Connections**: Routes to target with fewest active connections
3. **Random**: Randomly selects a healthy target
4. **Weighted Round Robin**: Routes based on assigned weights
5. **IP Hash**: Routes based on client IP hash for session affinity

### Health Checking
```go
type HealthChecker struct {
    Interval     time.Duration
    Timeout      time.Duration
    HealthyThreshold   int
    UnhealthyThreshold int
    Path         string
    ExpectedCode int
}
```

## Service Discovery

### Manual Configuration
Static configuration via YAML file or environment variables.

### Consul Integration
```go
type ConsulConfig struct {
    Address       string        `json:"address"`
    Scheme        string        `json:"scheme"`
    Datacenter    string        `json:"datacenter"`
    Token         string        `json:"token"`
    Interval      time.Duration `json:"interval"`
    Timeout       time.Duration `json:"timeout"`
}
```

### Kubernetes Integration
Automatic service discovery using Kubernetes API.

## SDK Interface

### Direct Mode (In-Process)
```go
type GatewaySdk struct {
    RouteRequest    func(req *http.Request) (*http.Response, error)
    AddRoute        func(route Route) error
    RemoveRoute     func(routeID string) error
    GetHealthStatus func() (*HealthStatus, error)
    GetMetrics      func() (*Metrics, error)
}

func NewGatewaySdk(mode string) *GatewaySdk
```

### HTTP Mode (Network)
Standard HTTP proxy functionality with load balancing and health checking.

## Monitoring & Metrics

### Prometheus Metrics
```
# Request metrics
gateway_requests_total{method, path, service, status}
gateway_request_duration_seconds{method, path, service}
gateway_request_size_bytes{method, path, service}
gateway_response_size_bytes{method, path, service}

# Service health metrics
gateway_service_up{service}
gateway_service_response_time{service}

# Load balancer metrics
gateway_lb_targets_total{service}
gateway_lb_targets_healthy{service}
gateway_lb_requests_total{service, target}

# Rate limiting metrics
gateway_rate_limit_hits_total{key}
gateway_rate_limit_rejects_total{key}
```

### Request Logging
```json
{
  "timestamp": "2025-09-22T10:00:00Z",
  "requestId": "req-uuid",
  "method": "POST",
  "path": "/api/v1/calculator",
  "service": "calculator-svc",
  "target": "http://localhost:8080",
  "clientIP": "192.168.1.100",
  "userAgent": "Mozilla/5.0...",
  "userId": "user-uuid",
  "duration": 125,
  "statusCode": 200,
  "requestSize": 256,
  "responseSize": 128,
  "error": null
}
```

## Environment Variables

```env
# Server Configuration
GATEWAY_PORT=8082
GATEWAY_READ_TIMEOUT=30s
GATEWAY_WRITE_TIMEOUT=30s

# Service URLs (when not using service discovery)
AUTH_SERVICE_URL=http://localhost:8081
CALCULATOR_SERVICE_URL=http://localhost:8080
PDF_SERVICE_URL=http://localhost:8083
TEMPERATURE_SERVICE_URL=http://localhost:8084

# CORS Configuration
CORS_ALLOW_ORIGINS=*
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOW_HEADERS=Content-Type,Authorization

# Rate Limiting
RATE_LIMIT_REQUESTS=1000
RATE_LIMIT_WINDOW=1h
RATE_LIMIT_STORAGE=memory

# Authentication
AUTH_ENDPOINT=http://localhost:8081/auth/verify
AUTH_CACHE_TTL=5m

# Service Discovery
SERVICE_DISCOVERY_ENABLED=false
SERVICE_DISCOVERY_PROVIDER=consul
CONSUL_ADDRESS=localhost:8500

# Monitoring
METRICS_ENABLED=true
LOGGING_LEVEL=info
TRACING_ENABLED=false
```

## Error Handling

### Gateway Error Responses
```json
{
  "error": {
    "code": "SERVICE_UNAVAILABLE",
    "message": "The requested service is currently unavailable",
    "service": "calculator-svc",
    "requestId": "req-uuid",
    "timestamp": "2025-09-22T10:00:00Z",
    "retryAfter": 30
  }
}
```

### Error Codes
- `ROUTE_NOT_FOUND`: No route configured for the requested path
- `SERVICE_UNAVAILABLE`: Target service is down or unhealthy
- `AUTHENTICATION_REQUIRED`: Valid authentication token required
- `RATE_LIMIT_EXCEEDED`: Request rate limit exceeded
- `TIMEOUT`: Request timeout exceeded
- `CIRCUIT_BREAKER_OPEN`: Circuit breaker is open for the service

## Dependencies

```go
// Add to go.mod
require (
    github.com/gin-gonic/gin v1.10.1
    github.com/hashicorp/consul/api v1.24.0
    github.com/prometheus/client_golang v1.17.0
    github.com/redis/go-redis/v9 v9.2.1
    gopkg.in/yaml.v3 v3.0.1
    github.com/sony/gobreaker v0.5.0
    golang.org/x/time v0.3.0
)
```

## Performance Requirements

- Request routing latency: < 10ms
- Service discovery refresh: < 5s
- Health check frequency: 30s
- Concurrent connections: 10,000+
- Request throughput: 50,000 req/sec

## Security Considerations

### TLS/SSL
- HTTPS termination at gateway
- Certificate management
- TLS 1.3 support

### Security Headers
```go
// Security headers added to all responses
securityHeaders := map[string]string{
    "X-Frame-Options":           "DENY",
    "X-Content-Type-Options":    "nosniff",
    "X-XSS-Protection":          "1; mode=block",
    "Strict-Transport-Security": "max-age=31536000; includeSubDomains",
    "Content-Security-Policy":   "default-src 'self'",
}
```

### Request Validation
- Input sanitization
- Request size limits
- Header validation
- Path traversal protection

## Testing Requirements

### Unit Tests
- Route matching logic
- Middleware functionality
- Load balancer algorithms
- Health checker logic

### Integration Tests
- End-to-end request routing
- Service discovery integration
- Authentication middleware
- Rate limiting functionality

### Load Tests
- High concurrent request handling
- Performance under load
- Memory usage optimization
- Connection pooling

## Deployment Strategies

### Standalone Gateway
Single gateway instance handling all traffic.

### High Availability
Multiple gateway instances behind a load balancer.

### Edge Gateway
Deploy at network edge for geo-distributed services.

### Sidecar Pattern
Deploy as sidecar container with each service group.