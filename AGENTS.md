# Project Instructions

This project exposes a gRPC API to control a drone.

- It uses Protobuf to define the API and Go language to implement it.
- It uses the [MAVLink Protocol](https://mavlink.io/en/) to control the drone.

## Code Style

- Use Go for all new files
- Follow the project structure defined in README.md
- When writing comments for functions, use the following format:

```go
// FunctionName
// Describes what the function does, starting with a verb (e.g., "Converts", "Processes", "Validates").
// Additional sentences provide more details.
func FunctionName(param Type) ReturnType {
    // ...
}
```

Example:

```go
// HeartbeatMessageToMap
// Converts a HEARTBEAT message to a map with decoded fields for better readability.
// For example, the PX4 CustomMode is decoded into a human-readable format.
func HeartbeatMessageToMap(msg *common.MessageHeartbeat) (map[string]interface{}, error) {
    // ...
}
```

## Reference documentation
- MAVLink Protocol: https://mavlink.io/en/
- gomavlib – a library that implements the Mavlink protocol in Go: https://github.com/bluenviron/gomavlib