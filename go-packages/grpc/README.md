# grpc

Reusable gRPC server wrapper with graceful shutdown support.

## Installation

```go
import grpcserver "polycoin/grpc"
```

## Usage

### Basic

```go
server := grpcserver.NewServer(nil) // uses default logger
pb.RegisterYourServiceServer(server.GrpcServer(), yourHandler)
server.EnableReflection() // optional, for grpcurl debugging

if err := server.Run("8080"); err != nil {
    log.Fatal(err)
}

// After Run returns (on SIGINT/SIGTERM), handle cleanup:
server.GracefulStop()
db.Close()
// ... other cleanup
```

### With Custom Logger

Provide a logger that implements the `Logger` interface:

```go
type Logger interface {
    Info(msg string, context ...any)
}
```

Example:

```go
server := grpcserver.NewServer(&grpcserver.ServerDeps{
    Log: appLogger,
})
```

## API

| Method | Description |
|--------|-------------|
| `NewServer(deps)` | Create server, deps is optional |
| `GrpcServer()` | Returns underlying `*grpc.Server` for service registration |
| `EnableReflection()` | Enable gRPC reflection for debugging tools like grpcurl |
| `Run(port)` | Start server, block until SIGINT/SIGTERM |
| `GracefulStop()` | Stop server gracefully |

## Full Example

```go
package main

import (
    "log"

    grpcserver "polycoin/grpc"
    pb "polycoin/grpc/pb/user"
    "polycoin/logger"
    "yourapp/internal/config"
    "yourapp/internal/handlers"
)

func main() {
    config.Load()

    appLogger := logger.NewZapMust(logger.ZapOptions{
        ServiceName: "my-service",
    })

    userHandler := handlers.NewUserHandler()

    server := grpcserver.NewServer(&grpcserver.ServerDeps{
        Log: appLogger,
    })
    pb.RegisterUserServiceServer(server.GrpcServer(), userHandler)
    server.EnableReflection()

    if err := server.Run(config.AppConfig.Port); err != nil {
        log.Fatal("server error:", err)
    }

    // Graceful shutdown sequence
    appLogger.Info("Shutting down...")
    server.GracefulStop()
    // db.Close()
    // cache.Close()
    // etc.

    appLogger.Flush()
}
```
