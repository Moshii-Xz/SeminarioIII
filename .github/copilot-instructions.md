# VenciTrack - AI Agent Instructions

## Project Overview
VenciTrack is a Go-based REST API for managing product expiration tracking with e-commerce capabilities. Built with **GORM**, **Gin** (framework implied by structure), and **PostgreSQL**.

**Module path**: `github.com/Moshii-Xz/SeminarioIII`

## Architecture: Layered Module Design

This project follows **Clean Architecture** with vertical module slicing:

```
cmd/api/           - Application entry point
internal/
  ├── config/      - App configuration loading
  ├── domain/      - Core business entities (User, Product, Order, Payment)
  ├── middleware/  - HTTP middleware (auth, cors, logging)
  ├── modules/     - Feature modules (catalog, orders, payments, reports, reviews, users)
  │   └── [module]/
  │       ├── dto.go        - Request/response DTOs
  │       ├── handler.go    - HTTP handlers (Gin)
  │       ├── service.go    - Business logic
  │       └── repository.go - Data access (GORM)
  ├── platform/    - Infrastructure abstractions
  │   ├── database/ - GORM/Postgres setup, transactions
  │   ├── http/     - HTTP response utilities
  │   └── logger/   - Logging configuration
  └── server/      - HTTP server setup (Gin router, server init)
pkg/               - Shared public utilities (pagination, validations)
migrations/        - Database migrations
```

### Key Architectural Decisions

1. **Vertical Module Slicing**: Each module (`catalog`, `orders`, `payments`, etc.) is self-contained with its own handler→service→repository layers. Cross-module communication should go through service interfaces.

2. **Domain-Driven Boundaries**: 
   - `internal/domain/` contains pure business entities
   - Modules reference domain entities but own their DTOs
   - Keep domain models free from framework tags (GORM, JSON) when possible

3. **Dependency Flow**: 
   - Handler → Service → Repository → Domain
   - Platform packages are infrastructure; inject DB connections at module initialization
   - Middleware is registered at router level (`internal/server/router.go`)

## Critical Conventions

### Module Pattern (Handler-Service-Repository)
Each module follows this structure:

**handler.go** - HTTP layer (Gin handlers):
```go
type Handler struct {
    service *Service
}

func (h *Handler) GetProduct(c *gin.Context) {
    // Parse request → call service → return response
}
```

**service.go** - Business logic:
```go
type Service struct {
    repo *Repository
}

func (s *Service) GetProduct(id uint) (*domain.Product, error) {
    // Business rules here
}
```

**repository.go** - Data access (GORM):
```go
type Repository struct {
    db *gorm.DB
}

func (r *Repository) FindByID(id uint) (*domain.Product, error) {
    // GORM queries
}
```

**dto.go** - Data Transfer Objects (request/response models):
```go
type CreateProductRequest struct {
    Name string `json:"name" binding:"required"`
    // Use Gin binding tags for validation
}
```

### Database Patterns

**Connection Setup** (`internal/platform/database/postgres.go`):
- Use the existing `database.Config` and `database.New()` pattern
- Database config example:
  ```go
  type Config struct {
      Host     string
      Port     int
      User     string
      Password string
      DBName   string
      SSLMode  string
      TimeZone string
  }
  ```
- Returns `*gorm.DB` with connection pooling configured

**Transactions** (`internal/platform/database/transaction.go`):
- Implement repository methods to accept `*gorm.DB` for transaction support
- Service layer orchestrates transactions across repositories

**Migrations**:
- Store SQL migrations in `migrations/` directory
- Use sequential numbering: `001_create_users.sql`, `002_create_products.sql`

### Response Handling (`internal/platform/http/response.go`)
Standardize API responses:
```go
// Success response
type SuccessResponse struct {
    Data interface{} `json:"data"`
}

// Error response
type ErrorResponse struct {
    Error   string `json:"error"`
    Code    string `json:"code,omitempty"`
    Details interface{} `json:"details,omitempty"`
}
```

### Middleware Registration
Register middleware in `internal/server/router.go`:
```go
router.Use(middleware.CORS())
router.Use(middleware.Logging())

// Protected routes
authorized := router.Group("/api/v1")
authorized.Use(middleware.Auth())
```

## Development Workflow

### Setup & Dependencies
```powershell
# Install dependencies
go mod download

# Add new dependency
go get github.com/example/package

# Update dependencies
go get -u ./...
```

### Running the Application
```powershell
# Run from project root
go run cmd/api/main.go

# Build binary
go build -o vencitrack.exe cmd/api/main.go

# Run with hot reload (if air is installed)
air
```

### Database Migrations
```powershell
# Apply migrations (implement migration runner in main.go or separate tool)
# Common pattern: use golang-migrate or custom SQL runner
```

### Testing
```powershell
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific module tests
go test ./internal/modules/users/...
```

## Module-Specific Notes

### Users Module (`internal/modules/users/`)
- Handles authentication, user CRUD
- Auth middleware validates JWT tokens from this module

### Catalog Module (`internal/modules/catalog/`)
- Product management with expiration tracking
- Core domain for the app's expiration monitoring features

### Orders Module (`internal/modules/orders/`)
- Order lifecycle management
- Coordinates with payments module for payment processing

### Payments Module (`internal/modules/payments/`)
- Payment gateway integration (`gateway.go` suggests external payment provider)
- Handle payment webhooks and status updates

### Reports Module (`internal/modules/reports/`)
- Analytics and reporting endpoints
- May aggregate data from multiple modules

### Reviews Module (`internal/modules/reviews/`)
- Product reviews and ratings
- Links to users and catalog modules

## Technology Stack

- **Go 1.24**: Latest Go features available
- **GORM v1.31.1**: ORM for PostgreSQL (`gorm.io/gorm`, `gorm.io/driver/postgres`)
- **PostgreSQL**: Primary database (via pgx driver v5)
- **Gin** (implied): HTTP framework - handlers follow Gin patterns (`gin.Context`)
- **JWT**: Authentication (implement in `internal/middleware/auth.go`)

## Common Tasks

### Adding a New Module
1. Create directory: `internal/modules/[name]/`
2. Create files: `dto.go`, `handler.go`, `service.go`, `repository.go`
3. Define domain entities in `internal/domain/[name].go`
4. Wire up handler in `internal/server/router.go`
5. Initialize repository with DB connection

### Adding a New Endpoint
1. Define DTO in `dto.go` (request/response structs)
2. Add service method in `service.go` (business logic)
3. Add repository method in `repository.go` (if DB access needed)
4. Create handler function in `handler.go` (HTTP layer)
5. Register route in `internal/server/router.go`

### Database Schema Changes
1. Create migration file in `migrations/` (e.g., `003_add_expiration_alerts.sql`)
2. Update domain models in `internal/domain/`
3. Update repository queries as needed
4. Run migration script

## Code Style & Patterns

- **Error Handling**: Return errors, don't panic (except in `main.go` for fatal setup errors)
- **Context**: Pass `context.Context` for cancellation support in long operations
- **Naming**: Use descriptive names; `GetProductByID`, not `Get`
- **Package Organization**: Keep packages focused; avoid circular dependencies
- **Dependency Injection**: Pass dependencies via constructors (`NewHandler(service *Service)`)

## Important Files to Reference

- `go.mod` - Module dependencies and Go version
- `internal/platform/database/postgres.go` - Database connection pattern
- `internal/server/router.go` - Route definitions and middleware setup
- `internal/domain/*.go` - Core business entities
- Individual module directories - Follow existing patterns when adding features
