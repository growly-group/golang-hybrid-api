# Authentication Service (auth-svc) Specification

## Overview
The Authentication Service provides secure user authentication, authorization, and session management for the hybrid API ecosystem. It supports both in-process and HTTP-based communication patterns.

## Port Assignment
- **Default Port**: 8081
- **Health Check**: 8081/health

## Core Features

### 1. User Authentication
- Email/password authentication
- JWT token-based authentication
- Password hashing with bcrypt
- Account lockout after failed attempts
- Password reset functionality

### 2. Authorization
- Role-based access control (RBAC)
- Permission-based authorization
- Resource-level access control
- API key management

### 3. Session Management
- JWT token generation and validation
- Token refresh mechanism
- Session expiration and cleanup
- Blacklist for revoked tokens

## API Endpoints

### Authentication Endpoints

#### POST /auth/register
Register a new user account.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securePassword123",
  "firstName": "John",
  "lastName": "Doe",
  "role": "user"
}
```

**Response (201):**
```json
{
  "id": "user-uuid",
  "email": "user@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "role": "user",
  "createdAt": "2025-09-22T10:00:00Z"
}
```

#### POST /auth/login
Authenticate user and return JWT token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Response (200):**
```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "refresh-token-here",
  "expiresIn": 3600,
  "user": {
    "id": "user-uuid",
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "role": "user"
  }
}
```

#### POST /auth/refresh
Refresh access token using refresh token.

**Request Body:**
```json
{
  "refreshToken": "refresh-token-here"
}
```

#### POST /auth/logout
Invalidate user session and blacklist tokens.

**Headers:**
```
Authorization: Bearer <access-token>
```

#### POST /auth/forgot-password
Initiate password reset process.

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

#### POST /auth/reset-password
Reset password using reset token.

**Request Body:**
```json
{
  "resetToken": "reset-token-here",
  "newPassword": "newSecurePassword123"
}
```

### Authorization Endpoints

#### GET /auth/verify
Verify JWT token and return user information.

**Headers:**
```
Authorization: Bearer <access-token>
```

**Response (200):**
```json
{
  "valid": true,
  "user": {
    "id": "user-uuid",
    "email": "user@example.com",
    "role": "user",
    "permissions": ["read:profile", "write:profile"]
  }
}
```

#### POST /auth/check-permission
Check if user has specific permission.

**Request Body:**
```json
{
  "userId": "user-uuid",
  "permission": "read:admin-panel",
  "resource": "admin-dashboard"
}
```

### User Management Endpoints

#### GET /auth/profile
Get current user profile.

#### PUT /auth/profile
Update user profile.

#### PUT /auth/change-password
Change user password.

### Admin Endpoints

#### GET /auth/users
List all users (admin only).

#### PUT /auth/users/:id/role
Update user role (admin only).

#### DELETE /auth/users/:id
Deactivate user account (admin only).

## Data Models

### User Model
```go
type User struct {
    ID           string    `json:"id" db:"id"`
    Email        string    `json:"email" db:"email"`
    PasswordHash string    `json:"-" db:"password_hash"`
    FirstName    string    `json:"firstName" db:"first_name"`
    LastName     string    `json:"lastName" db:"last_name"`
    Role         string    `json:"role" db:"role"`
    IsActive     bool      `json:"isActive" db:"is_active"`
    CreatedAt    time.Time `json:"createdAt" db:"created_at"`
    UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
    LastLoginAt  *time.Time `json:"lastLoginAt" db:"last_login_at"`
    FailedLogins int       `json:"-" db:"failed_logins"`
    LockedUntil  *time.Time `json:"-" db:"locked_until"`
}
```

### Session Model
```go
type Session struct {
    ID           string    `json:"id" db:"id"`
    UserID       string    `json:"userId" db:"user_id"`
    AccessToken  string    `json:"-" db:"access_token"`
    RefreshToken string    `json:"-" db:"refresh_token"`
    ExpiresAt    time.Time `json:"expiresAt" db:"expires_at"`
    CreatedAt    time.Time `json:"createdAt" db:"created_at"`
    IsRevoked    bool      `json:"isRevoked" db:"is_revoked"`
}
```

### Permission Model
```go
type Permission struct {
    ID          string `json:"id" db:"id"`
    Name        string `json:"name" db:"name"`
    Description string `json:"description" db:"description"`
}

type RolePermission struct {
    RoleID       string `json:"roleId" db:"role_id"`
    PermissionID string `json:"permissionId" db:"permission_id"`
}
```

## SDK Interface

### Direct Mode (In-Process)
```go
type AuthSdk struct {
    Register       func(req RegisterRequest) (*User, error)
    Login          func(req LoginRequest) (*LoginResponse, error)
    VerifyToken    func(token string) (*User, error)
    CheckPermission func(userID, permission, resource string) (bool, error)
    RefreshToken   func(refreshToken string) (*LoginResponse, error)
}

func NewAuthSdk(mode string) *AuthSdk
```

### HTTP Mode (Network)
When running as separate service, the SDK makes HTTP calls to the auth service endpoints.

## Security Requirements

### Password Security
- Minimum 8 characters
- Must contain uppercase, lowercase, number, and special character
- Bcrypt hashing with cost factor 12
- Password history to prevent reuse

### Token Security
- JWT with HS256 algorithm
- Access token expires in 15 minutes
- Refresh token expires in 7 days
- Secure token storage recommendations

### Rate Limiting
- Login attempts: 5 per minute per IP
- Registration: 3 per hour per IP
- Password reset: 3 per hour per email

### Account Security
- Account lockout after 5 failed login attempts
- Lockout duration: 15 minutes
- Email verification for new accounts
- Audit logging for security events

## Environment Variables

```env
# Database
AUTH_DB_HOST=localhost
AUTH_DB_PORT=5432
AUTH_DB_NAME=auth_db
AUTH_DB_USER=auth_user
AUTH_DB_PASSWORD=auth_password

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key
JWT_ACCESS_TOKEN_EXPIRES=15m
JWT_REFRESH_TOKEN_EXPIRES=168h

# Email Configuration (for password reset)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=your-app-password

# Security
BCRYPT_COST=12
MAX_LOGIN_ATTEMPTS=5
LOCKOUT_DURATION=15m
```

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'user',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP,
    failed_logins INTEGER DEFAULT 0,
    locked_until TIMESTAMP
);
```

### Sessions Table
```sql
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    access_token_hash VARCHAR(255) NOT NULL,
    refresh_token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_revoked BOOLEAN DEFAULT false
);
```

### Permissions Table
```sql
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT
);

CREATE TABLE role_permissions (
    role_id VARCHAR(50) NOT NULL,
    permission_id UUID REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);
```

## Dependencies

```go
// Add to go.mod
require (
    github.com/golang-jwt/jwt/v5 v5.0.0
    golang.org/x/crypto v0.14.0
    github.com/lib/pq v1.10.9
    github.com/google/uuid v1.3.1
    gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
)
```

## Error Handling

### HTTP Status Codes
- 200: Success
- 201: Created (registration)
- 400: Bad Request (validation errors)
- 401: Unauthorized (invalid credentials)
- 403: Forbidden (insufficient permissions)
- 409: Conflict (email already exists)
- 429: Too Many Requests (rate limited)
- 500: Internal Server Error

### Error Response Format
```json
{
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Invalid email or password",
    "details": {
      "field": "password",
      "reason": "Password does not match"
    }
  }
}
```

## Testing Requirements

### Unit Tests
- User registration and validation
- Password hashing and verification
- JWT token generation and validation
- Permission checking logic

### Integration Tests
- Complete authentication flow
- Token refresh mechanism
- Account lockout functionality
- Rate limiting

### Security Tests
- SQL injection protection
- XSS protection
- CSRF protection
- Brute force attack protection

## Performance Requirements

- Login response time: < 200ms
- Token verification: < 50ms
- Concurrent users: 1000+
- Database connection pooling
- Redis caching for sessions (optional)

## Monitoring and Logging

### Metrics to Track
- Login success/failure rates
- Token generation/validation counts
- Account lockout events
- API response times

### Security Audit Log
- All authentication attempts
- Permission changes
- Account lockouts
- Password resets

## Migration Path

### Phase 1: Basic Authentication
- User registration and login
- JWT token generation
- Basic password security

### Phase 2: Advanced Security
- Role-based permissions
- Account lockout
- Rate limiting

### Phase 3: Enterprise Features
- SSO integration
- Multi-factor authentication
- Advanced audit logging