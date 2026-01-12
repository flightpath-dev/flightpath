# Flightpath

Go platform exposing a gRPC API to control a drone.

## Architecture

![architecture](./assets/flightpath-architecture.png)

## Project Structure
```
flightpath/
├── cmd/                            # Main applications for this project
│   └── server/                     # Server entry point
│       └── main.go
├── examples/                       # API usage examples
├── gen/                            # gRPC code generated from protobufs
│   ├── go/
│   └── ts/
├── internal/                       # Private application and library code
│   ├── config/
│   │   ├── config.go               # Configuration structure
│   │   └── loader.go               # Configuration loader
│   ├── mavlink/                    # Drone side communication code
│   │   ├── command_dispatcher.go   # Sends commands to drone
│   │   └── message_receiver.go     # Receives messages from Drone
│   ├── middleware/                 # Handlers to process HTTP requests
│   │   ├── cors.go                 # CORS middleware
│   │   ├── logging.go              # Request logging middleware
│   │   └── recovery.go             # Panic recovery middleware
│   ├── server/                     # Server state and lifecycle management
│   │   └── server.go
│   └── services/                   # Core functionality offered by the application
│       ├── context.go              # Shared context for all services (config, logger, etc.)
│       └── mavlink_service.go      # Distributes MAVLink messages to gRPC subscribers
├── proto/                          # Protocol Buffers definitions
│   └── flightpath/
│       └── mavlink_service.proto   # Handles commands and messages from gRPC clients
├── go.mod
└── go.sum
```

## Quick Start

### Prerequisites

```bash
# 1. Clone repository
git clone https://github.com/flightpath-dev/flightpath
cd flightpath

# 2. Install dependencies
go mod tidy
```

### Run with a PX4 SITL

Start a PX4 SITL by following the instructions in [PX4 SITL Setup](docs/px4-sitl-setup.md).

```bash
# 1. Run server
go run cmd/server/main.go

# 2. Monitor messages from the SITL
go run examples/monitor_heartbeat_flightpath/main.go
```

### Run with a drone connected to a serial port

```bash
# 1. Turn on the drone

# 2. Run the server with a serial port configuration 
./scripts/run-serial.sh

# 3. Monitor messages from the drone
go run examples/monitor_heartbeat_flightpath/main.go
```

### Run with a drone connected over a UDP port

```bash
# 1. Turn on the drone

# 2. Run the server with a UDP configuration
./scripts/run-udp.sh

# 3. Monitor messages from the drone
go run examples/monitor_heartbeat_flightpath/main.go
```

## License

MIT
