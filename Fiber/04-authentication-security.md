# Authentication & Security

## Overview
This feature provides JWT-based authentication with Supabase integration, request validation, and comprehensive security middleware for the Motion-Index API built with Go Fiber.

## Current Python Implementation Analysis

### Key Components (from API analysis):
- **`src/middleware/auth.py`**: JWT authentication with Supabase integration
- **Security in `server.py`**: CORS, file validation, path traversal protection

### Current Features:
- JWT token verification with multiple algorithms (HS256, RS256)
- Supabase user management integration
- Flexible authentication supporting both secrets and public keys
- CORS configuration for multiple domains
- File serving with security checks

### Protected Endpoints:
- `POST /update-metadata` - Requires authentication

## Go Package Design

### Package Structure:
```
pkg/
├── auth/
│   ├── jwt/                 # JWT token handling
│   │   ├── validator.go     # JWT validation logic
│   │   ├── claims.go        # JWT claims parsing
│   │   ├── supabase.go      # Supabase-specific JWT handling
│   │   └── middleware.go    # Fiber JWT middleware
│   ├── user/                # User management
│   │   ├── service.go       # User service interface
│   │   ├── supabase.go      # Supabase user operations
│   │   └── models.go        # User data models
│   ├── security/            # Security utilities
│   │   ├── validation.go    # Input validation
│   │   ├── sanitizer.go     # Request sanitization
│   │   ├── rate_limit.go    # Rate limiting
│   │   └── headers.go       # Security headers
│   └── middleware/          # Fiber middleware
│       ├── cors.go          # CORS configuration
│       ├── auth.go          # Authentication middleware
│       ├── security.go      # Security headers middleware
│       └── validation.go    # Request validation middleware
```

### Core Interfaces:

```go
// JWTValidator interface for token validation
type JWTValidator interface {
    ValidateToken(ctx context.Context, tokenString string) (*Claims, error)
    RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
    GetUserFromToken(ctx context.Context, token string) (*User, error)
}

// UserService interface for user management
type UserService interface {
    GetUser(ctx context.Context, userID string) (*User, error)
    GetUserByEmail(ctx context.Context, email string) (*User, error)
    UpdateUser(ctx context.Context, userID string, updates map[string]interface{}) error
    ValidateUser(ctx context.Context, userID string) error
}

// SecurityService interface for security operations
type SecurityService interface {
    ValidateInput(input interface{}) error
    SanitizeString(input string) string
    CheckRateLimit(ctx context.Context, key string, limit int, window time.Duration) error
    GenerateCSRFToken(ctx context.Context) (string, error)
    ValidateCSRFToken(ctx context.Context, token string) error
}
```

### Data Models:

```go
type User struct {
    ID            string                 `json:"id"`
    Email         string                 `json:"email"`
    Role          string                 `json:"role"`
    Name          string                 `json:"name,omitempty"`
    Organization  string                 `json:"organization,omitempty"`
    Permissions   []string               `json:"permissions,omitempty"`
    Metadata      map[string]interface{} `json:"metadata,omitempty"`
    CreatedAt     time.Time              `json:"created_at"`
    UpdatedAt     time.Time              `json:"updated_at"`
    LastLoginAt   *time.Time             `json:"last_login_at,omitempty"`
    EmailVerified bool                   `json:"email_verified"`
    Active        bool                   `json:"active"`
}

type Claims struct {
    UserID      string                 `json:"sub"`
    Email       string                 `json:"email"`
    Role        string                 `json:"role"`
    Permissions []string               `json:"permissions,omitempty"`
    Metadata    map[string]interface{} `json:"user_metadata,omitempty"`
    jwt.RegisteredClaims
}

type TokenPair struct {
    AccessToken  string    `json:"access_token"`
    RefreshToken string    `json:"refresh_token,omitempty"`
    TokenType    string    `json:"token_type"`
    ExpiresIn    int64     `json:"expires_in"`
    ExpiresAt    time.Time `json:"expires_at"`
}

type AuthConfig struct {
    SupabaseURL        string        `env:"SUPABASE_URL" required:"true"`
    SupabaseAnonKey    string        `env:"SUPABASE_ANON_KEY" required:"true"`
    SupabaseServiceKey string        `env:"SUPABASE_SERVICE_KEY" required:"true"`
    JWTSecret          string        `env:"SUPABASE_JWT_SECRET"`
    JWTPublicKey       string        `env:"SUPABASE_PUBLIC_KEY"`
    TokenExpiry        time.Duration `env:"JWT_TOKEN_EXPIRY" default:"1h"`
    RefreshExpiry      time.Duration `env:"JWT_REFRESH_EXPIRY" default:"24h"`
}
```

## Fiber Middleware Implementation

### JWT Authentication Middleware:
```go
func NewJWTMiddleware(validator JWTValidator) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Extract token from Authorization header
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return fiber.NewError(fiber.StatusUnauthorized, "Authorization header required")
        }
        
        // Parse Bearer token
        tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            return fiber.NewError(fiber.StatusUnauthorized, "Invalid authorization format")
        }
        
        tokenString := tokenParts[1]
        
        // Validate token
        claims, err := validator.ValidateToken(c.Context(), tokenString)
        if err != nil {
            return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired token")
        }
        
        // Store user information in context
        c.Locals("user_id", claims.UserID)
        c.Locals("user_email", claims.Email)
        c.Locals("user_role", claims.Role)
        c.Locals("user_permissions", claims.Permissions)
        c.Locals("user_claims", claims)
        
        return c.Next()
    }
}
```

### CORS Middleware:
```go
func NewCORSMiddleware() fiber.Handler {
    return cors.New(cors.Config{
        AllowOrigins: strings.Join([]string{
            "https://api.motionindex.techjusticelab.org",
            "http://localhost:5173",
            "https://localhost:5173",
            "http://localhost:3000",
            "https://localhost:3000",
            "https://motionindex.techjusticelab.org",
            "https://motion-index.vercel.app",
        }, ","),
        AllowCredentials: true,
        AllowMethods: strings.Join([]string{
            fiber.MethodGet,
            fiber.MethodPost,
            fiber.MethodPut,
            fiber.MethodDelete,
            fiber.MethodOptions,
        }, ","),
        AllowHeaders: strings.Join([]string{
            "Origin",
            "Content-Type",
            "Accept",
            "Authorization",
            "X-Requested-With",
            "X-CSRF-Token",
        }, ","),
        ExposeHeaders: strings.Join([]string{
            "X-Total-Count",
            "X-Rate-Limit-Remaining",
            "X-Rate-Limit-Reset",
        }, ","),
        MaxAge: int(12 * time.Hour / time.Second),
    })
}
```

### Security Headers Middleware:
```go
func NewSecurityMiddleware() fiber.Handler {
    return func(c *fiber.Ctx) error {
        // Security headers
        c.Set("X-Content-Type-Options", "nosniff")
        c.Set("X-Frame-Options", "DENY")
        c.Set("X-XSS-Protection", "1; mode=block")
        c.Set("Referrer-Policy", "strict-origin-when-cross-origin")
        c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
        
        // HSTS for HTTPS (conditional)
        if c.Protocol() == "https" {
            c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        }
        
        return c.Next()
    }
}
```

### Rate Limiting Middleware:
```go
func NewRateLimitMiddleware(redisClient *redis.Client) fiber.Handler {
    return limiter.New(limiter.Config{
        Max:        100,              // Maximum requests
        Expiration: 1 * time.Minute,  // Time window
        KeyGenerator: func(c *fiber.Ctx) string {
            // Use user ID if authenticated, otherwise IP
            if userID := c.Locals("user_id"); userID != nil {
                return fmt.Sprintf("user:%s", userID.(string))
            }
            return fmt.Sprintf("ip:%s", c.IP())
        },
        LimitReached: func(c *fiber.Ctx) error {
            return fiber.NewError(fiber.StatusTooManyRequests, "Rate limit exceeded")
        },
        Storage: limiter.ConfigDefault.Storage, // In-memory storage
    })
}
```

## JWT Validation Implementation

### Supabase JWT Validator:
```go
type SupabaseJWTValidator struct {
    jwtSecret    string
    publicKey    *rsa.PublicKey
    supabaseURL  string
    userService  UserService
}

func NewSupabaseJWTValidator(config *AuthConfig, userService UserService) (*SupabaseJWTValidator, error) {
    validator := &SupabaseJWTValidator{
        jwtSecret:   config.JWTSecret,
        supabaseURL: config.SupabaseURL,
        userService: userService,
    }
    
    // Parse public key if provided
    if config.JWTPublicKey != "" {
        publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(config.JWTPublicKey))
        if err != nil {
            return nil, fmt.Errorf("failed to parse JWT public key: %w", err)
        }
        validator.publicKey = publicKey
    }
    
    return validator, nil
}

func (v *SupabaseJWTValidator) ValidateToken(ctx context.Context, tokenString string) (*Claims, error) {
    // Parse token with claims
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        // Check signing method
        switch token.Method.(type) {
        case *jwt.SigningMethodHMAC:
            if v.jwtSecret == "" {
                return nil, fmt.Errorf("HMAC signing method requires JWT secret")
            }
            return []byte(v.jwtSecret), nil
        case *jwt.SigningMethodRSA:
            if v.publicKey == nil {
                return nil, fmt.Errorf("RSA signing method requires public key")
            }
            return v.publicKey, nil
        default:
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
    })
    
    if err != nil {
        return nil, fmt.Errorf("failed to parse token: %w", err)
    }
    
    if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }
    
    claims, ok := token.Claims.(*Claims)
    if !ok {
        return nil, fmt.Errorf("invalid token claims")
    }
    
    // Validate token expiration
    if claims.ExpiresAt.Time.Before(time.Now()) {
        return nil, fmt.Errorf("token has expired")
    }
    
    // Optional: Validate user still exists and is active
    if v.userService != nil {
        err = v.userService.ValidateUser(ctx, claims.UserID)
        if err != nil {
            return nil, fmt.Errorf("user validation failed: %w", err)
        }
    }
    
    return claims, nil
}
```

## Request Validation

### Input Validation Middleware:
```go
func NewValidationMiddleware() fiber.Handler {
    validate := validator.New()
    
    return func(c *fiber.Ctx) error {
        // Skip validation for GET requests and specific paths
        if c.Method() == fiber.MethodGet || 
           strings.HasPrefix(c.Path(), "/health") ||
           strings.HasPrefix(c.Path(), "/api/documents/") {
            return c.Next()
        }
        
        // Parse request body based on content type
        contentType := c.Get("Content-Type")
        
        if strings.Contains(contentType, "application/json") {
            var body map[string]interface{}
            if err := c.BodyParser(&body); err != nil {
                return fiber.NewError(fiber.StatusBadRequest, "Invalid JSON body")
            }
            
            // Validate required fields based on endpoint
            if err := validateEndpointRequirements(c.Path(), c.Method(), body, validate); err != nil {
                return fiber.NewError(fiber.StatusBadRequest, err.Error())
            }
        }
        
        return c.Next()
    }
}

func validateEndpointRequirements(path, method string, body map[string]interface{}, validate *validator.Validate) error {
    switch {
    case path == "/search" && method == "POST":
        return validateSearchRequest(body, validate)
    case path == "/update-metadata" && method == "POST":
        return validateMetadataUpdate(body, validate)
    case strings.HasPrefix(path, "/categorise") && method == "POST":
        return validateFileUpload(body, validate)
    default:
        return nil // No specific validation required
    }
}
```

## User Management

### Supabase User Service:
```go
type SupabaseUserService struct {
    client      *supabase.Client
    serviceKey  string
}

func NewSupabaseUserService(config *AuthConfig) *SupabaseUserService {
    client := supabase.CreateClient(config.SupabaseURL, config.SupabaseAnonKey)
    
    return &SupabaseUserService{
        client:     client,
        serviceKey: config.SupabaseServiceKey,
    }
}

func (s *SupabaseUserService) GetUser(ctx context.Context, userID string) (*User, error) {
    // Use service key for admin operations
    adminClient := supabase.CreateClient(s.client.BaseURL, s.serviceKey)
    
    var user User
    err := adminClient.DB.From("auth.users").
        Select("*").
        Eq("id", userID).
        Single().
        ExecuteTo(&user)
    
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return &user, nil
}

func (s *SupabaseUserService) ValidateUser(ctx context.Context, userID string) error {
    user, err := s.GetUser(ctx, userID)
    if err != nil {
        return err
    }
    
    if !user.Active {
        return fmt.Errorf("user account is inactive")
    }
    
    if !user.EmailVerified {
        return fmt.Errorf("user email is not verified")
    }
    
    return nil
}
```

## Security Utilities

### Input Sanitization:
```go
func SanitizeString(input string) string {
    // Remove null bytes
    input = strings.ReplaceAll(input, "\x00", "")
    
    // Remove control characters except newlines and tabs
    var result strings.Builder
    for _, r := range input {
        if unicode.IsControl(r) && r != '\n' && r != '\t' && r != '\r' {
            continue
        }
        result.WriteRune(r)
    }
    
    return result.String()
}

func ValidateFilePath(path string) error {
    // Prevent path traversal
    if strings.Contains(path, "..") {
        return fmt.Errorf("path traversal detected")
    }
    
    // Prevent absolute paths
    if strings.HasPrefix(path, "/") {
        return fmt.Errorf("absolute paths not allowed")
    }
    
    // Check for null bytes
    if strings.Contains(path, "\x00") {
        return fmt.Errorf("null bytes not allowed in path")
    }
    
    return nil
}
```

## Test Strategy

### Unit Tests:
```go
func TestJWTValidator_ValidateToken(t *testing.T) {
    tests := []struct {
        name      string
        token     string
        secret    string
        wantError bool
    }{
        {
            name:      "valid token",
            token:     createValidToken("user123", "secret"),
            secret:    "secret",
            wantError: false,
        },
        {
            name:      "expired token",
            token:     createExpiredToken("user123", "secret"),
            secret:    "secret",
            wantError: true,
        },
        {
            name:      "invalid signature",
            token:     createValidToken("user123", "wrong-secret"),
            secret:    "secret",
            wantError: true,
        },
    }
    
    validator := NewSupabaseJWTValidator(&AuthConfig{JWTSecret: "secret"}, nil)
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := validator.ValidateToken(context.Background(), tt.token)
            if (err != nil) != tt.wantError {
                t.Errorf("ValidateToken() error = %v, wantError %v", err, tt.wantError)
            }
        })
    }
}
```

### Integration Tests:
- Supabase JWT token validation
- User management operations
- Rate limiting functionality
- CORS preflight requests
- Security header validation

## Implementation Priority

1. **JWT Validation** - Core authentication functionality
2. **Middleware Setup** - CORS, security headers, rate limiting
3. **User Management** - Supabase integration
4. **Input Validation** - Request sanitization and validation
5. **Advanced Security** - CSRF protection, additional headers
6. **Monitoring** - Authentication metrics and logging

## Dependencies

### External Libraries:
- `github.com/golang-jwt/jwt/v5` - JWT token handling
- `github.com/go-playground/validator/v10` - Input validation
- `github.com/gofiber/fiber/v2/middleware/cors` - CORS middleware
- `github.com/gofiber/fiber/v2/middleware/limiter` - Rate limiting
- `github.com/supabase-community/supabase-go` - Supabase client

### Configuration:
```go
type SecurityConfig struct {
    // JWT Configuration
    JWTSecret          string        `env:"JWT_SECRET"`
    JWTPublicKey       string        `env:"JWT_PUBLIC_KEY"`
    TokenExpiry        time.Duration `env:"TOKEN_EXPIRY" default:"1h"`
    
    // Supabase Configuration
    SupabaseURL        string `env:"SUPABASE_URL" required:"true"`
    SupabaseAnonKey    string `env:"SUPABASE_ANON_KEY" required:"true"`
    SupabaseServiceKey string `env:"SUPABASE_SERVICE_KEY" required:"true"`
    
    // Rate Limiting
    RateLimit          int           `env:"RATE_LIMIT" default:"100"`
    RateLimitWindow    time.Duration `env:"RATE_LIMIT_WINDOW" default:"1m"`
    
    // Security Headers
    EnableHSTS         bool   `env:"ENABLE_HSTS" default:"true"`
    CSPPolicy          string `env:"CSP_POLICY"`
    TrustedProxies     string `env:"TRUSTED_PROXIES"`
}
```

## Performance Considerations

- **Token Caching**: Cache validated tokens to reduce validation overhead
- **Connection Pooling**: Efficient Supabase client management
- **Rate Limiting**: Memory-efficient rate limiting with optional Redis backend
- **Middleware Order**: Optimize middleware execution order for performance

## Security Considerations

- **Token Security**: Secure JWT secret management
- **HTTPS Only**: Force HTTPS in production
- **Secure Headers**: Comprehensive security header implementation
- **Input Validation**: Strict input validation and sanitization
- **Rate Limiting**: Prevent brute force and DoS attacks
- **CORS Policy**: Restrictive CORS configuration