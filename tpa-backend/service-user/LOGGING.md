# User Service with Zap Logging

This service now includes structured logging using Uber's Zap logger with environment-based configuration.

## Environment Configuration

The logging behavior is controlled by the `ENVIRONMENT` variable in your `.env` file:

### Development Mode (default)

```env
ENVIRONMENT=development
```

- Colored console output
- Debug level logging
- Human-readable format
- Includes development-specific features

### Production Mode

```env
ENVIRONMENT=production
```

- JSON structured logs
- Info level logging and above
- Optimized for performance
- Machine-readable format

## Logging Features

### Structured Logging

All logs include contextual information using structured fields:

```go
logger.Info("User registered successfully",
    zap.String("user_id", user.ID.String()),
    zap.String("username", req.Username),
    zap.String("email", req.Email))
```

### Log Levels

- **Debug**: Detailed information for debugging (development only)
- **Info**: General information about service operations
- **Warn**: Warning conditions that don't prevent operation
- **Error**: Error conditions that may affect functionality
- **Fatal**: Critical errors that cause service termination

### Key Operations Logged

#### User Registration

- Registration requests with username and email
- Validation failures
- Duplicate username/email attempts
- Successful registrations

#### User Authentication

- Login attempts
- Inactive/banned account access attempts
- User lookup operations

#### Follow/Unfollow Operations

- Follow requests with user IDs
- Unfollow operations
- Database operation results

#### Premium Requests

- Premium registration requests
- Duplicate request attempts
- Request processing results

#### File Uploads

- Image/video upload requests with file information
- Upload success/failure with file details
- Supabase integration operations

#### Database Operations

- Connection initialization
- Migration execution
- Seeding operations
- Query errors and warnings

## Usage Examples

### Logging in Handlers

```go
func (h *UserHandler) SomeOperation(ctx context.Context, req *pb.Request) (*pb.Response, error) {
    h.logger.Info("Operation started",
        zap.String("operation", "SomeOperation"),
        zap.String("user_id", req.UserId))

    // ... operation logic ...

    if err != nil {
        h.logger.Error("Operation failed",
            zap.String("user_id", req.UserId),
            zap.Error(err))
        return nil, err
    }

    h.logger.Info("Operation completed successfully",
        zap.String("user_id", req.UserId))
    return result, nil
}
```

### Using Logger Utilities

```go
import "github.com/Wenfuuu/tpa-web/service-user/utils"

// Get logger instance
logger := utils.GetLogger()

// Or use utility functions
utils.LogInfo("Service started")
utils.LogError("Database connection failed", zap.Error(err))
```

## Log Output Examples

### Development Mode

```
2025-01-16T10:30:45.123+0700    INFO    Starting user service
2025-01-16T10:30:45.124+0700    INFO    Database connection established successfully
2025-01-16T10:30:45.125+0700    INFO    User service started    {"port": "50051", "environment": "development"}
```

### Production Mode

```json
{"level":"info","ts":1642320645.123,"msg":"Starting user service"}
{"level":"info","ts":1642320645.124,"msg":"Database connection established successfully"}
{"level":"info","ts":1642320645.125,"msg":"User service started","port":"50051","environment":"production"}
```

## Configuration Notes

- Logger is initialized once and reused throughout the application
- Environment variable is read from `.env` file
- Default fallback is development mode if `ENVIRONMENT` is not set
- All log output goes to stdout/stderr for container compatibility
- Logger instance is thread-safe and can be used concurrently
