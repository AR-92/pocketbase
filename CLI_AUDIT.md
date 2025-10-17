# Detailed Audit of PocketBase Core Directories

## APIs Directory (`apis/`)

The `apis/` directory contains the HTTP API layer for PocketBase, providing REST endpoints for all backend operations.

### Key Components:
- **Core API Handlers**:
  - `record_crud.go`: CRUD operations for records
  - `record_auth.go`: Authentication endpoints (login, register, password reset, OAuth2)
  - `collection.go`: Collection management (create, update, delete)
  - `backup.go`: Backup creation, restoration, and management
  - `logs.go`: Log viewing and statistics
  - `settings.go`: Application configuration
  - `realtime.go`: WebSocket subscriptions
  - `batch.go`: Batch operations
  - `file.go`: File upload/download handling
  - `health.go`: Health check endpoints

- **Middlewares**:
  - `middlewares_cors.go`: CORS handling
  - `middlewares_rate_limit.go`: Rate limiting
  - `middlewares_body_limit.go`: Request body size limits
  - `middlewares_gzip.go`: Response compression

- **Authentication Flows**:
  - Email verification, password reset, OAuth2 redirects
  - MFA and OTP support
  - Impersonation capabilities

- **Testing**: Comprehensive test coverage for all endpoints

### Architecture:
- Uses custom router with middleware chaining
- JSON-based request/response format
- Error handling with structured responses
- Pagination and filtering support

## Core Directory (`core/`)

The `core/` directory contains the business logic and data models for PocketBase.

### Key Components:
- **Application Core**:
  - `base.go`: Main application struct and lifecycle
  - `app.go`: Application interface and implementations
  - `events.go`: Event system for hooks

- **Data Models**:
  - `collection_model.go`: Collection definitions and validation
  - `record_model.go`: Record operations and field resolution
  - `field_*.go`: Individual field type implementations
  - `auth_origin_model.go`: External authentication
  - `external_auth_model.go`: OAuth2 providers
  - `mfa_model.go`, `otp_model.go`: Multi-factor authentication

- **Database Layer**:
  - `db.go`: Database connection and operations
  - `db_tx.go`: Transaction management
  - `db_builder.go`: Query building
  - `migrations_runner.go`: Migration execution

- **Field System**:
  - Support for 15+ field types (text, number, bool, date, file, relation, etc.)
  - Validation and type conversion
  - Field resolution for queries

- **Authentication & Security**:
  - Password hashing and validation
  - JWT token management
  - MFA/OTP implementation
  - Security utilities

### Architecture:
- Clean separation between data models and business logic
- Event-driven architecture with hooks
- Transaction-safe operations
- Extensible field system

## Tools Directory (`tools/`)

The `tools/` directory contains utility packages that support the core functionality.

### Key Components:
- **Authentication Providers** (`auth/`): 30+ OAuth2 providers (Google, GitHub, Discord, etc.)
- **Cron System** (`cron/`): Scheduled job execution
- **Mailer** (`mailer/`): Email sending (SMTP, Sendmail)
- **Router** (`router/`): HTTP routing with middleware support
- **Search & Filtering** (`search/`): Query parsing and field resolution
- **Security** (`security/`): Encryption, JWT, random generation
- **Subscriptions** (`subscriptions/`): Real-time client management
- **Template Engine** (`template/`): Template rendering
- **Types** (`types/`): Custom data types (JSON, DateTime, GeoPoint)
- **Logger** (`logger/`): Structured logging
- **OS Utilities** (`osutils/`): Cross-platform OS operations

### Architecture:
- Modular design with clear interfaces
- Comprehensive test coverage
- Cross-platform compatibility
- Performance-optimized implementations

## Integration Points

### APIs ↔ Core:
- APIs call core methods for business logic
- Core provides data validation and processing
- Event system allows APIs to trigger core hooks

### Core ↔ Tools:
- Core uses tools for authentication, mailing, scheduling
- Tools provide low-level utilities for core operations
- Shared types and interfaces

### APIs ↔ Tools:
- APIs use router and security tools
- Tools provide middleware and utility functions
- Search tools power API filtering

## Testing Coverage
- **APIs**: 40+ test files with comprehensive endpoint coverage
- **Core**: 50+ test files covering models, fields, and operations
- **Tools**: 30+ test files for utilities and integrations

## Performance Characteristics
- Efficient database queries with prepared statements
- Connection pooling and transaction management
- Caching for frequently accessed data
- Optimized field resolution and validation

This architecture provides a robust, scalable backend foundation with clear separation of concerns and extensive testing.</content>
</xai:function_call"><xai:function_call name="todowrite">
<parameter name="todos">[{"content":"Audit apis/, core/, and tools/ directories","status":"completed","priority":"high","id":"audit_directories"}]