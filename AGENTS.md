# Project Instructions

This project exposes a gRPC API to control a drone.

- It uses Protobuf to define the API and Go language to implement it.
- It uses the [MAVLink Protocol](https://mavlink.io/en/) to control the drone.

## Reference documentation

- [MAVLink Protocol Specification](https://mavlink.io/en/): Used to control the drone.
- [MAVLINK Common Message Set](https://mavlink.io/en/messages/common.html): Part of the above specification that contains a set of common messages and commands that should be implemented by MAVLink-compatible systems.
- [MAVLINK Common Message Set in XML format](https://github.com/mavlink/mavlink/blob/master/message_definitions/v1.0/common.xml): The common message set in a structured XML format. This is easier for programmatic understanding. Note that this set includes [standard.xml](https://github.com/mavlink/mavlink/blob/master/message_definitions/v1.0/standard.xml), which itself includes [minimal.xml](https://github.com/mavlink/mavlink/blob/master/message_definitions/v1.0/minimal.xml). For a pull understanding of the protocol, it's important to understand all three.
- [PX4 implementation of MAVLink Specs](https://docs.px4.io/main/en/mavlink/standard_modes): This page from the PX4 Guide describes how PX4 implements the MAVLink protocol. Specifically look at [Other MAVLink Mode-changing Commands](https://docs.px4.io/main/en/mavlink/standard_modes#other-mavlink-mode-changing-commands) for a list of specific commands to change modes. These can be more convenient that just starting the mode, in particular when the message allows additional settings to be configured.
- [gomavlib](https://github.com/bluenviron/gomavlib): The Go library used to send MAVLink commands to the drone and receive MAVLink messages. Especially, look at the [examples](https://github.com/bluenviron/gomavlib/tree/main/examples) directory to understand MAVLink workflows to achieve specific goals.
- [Protobuf Style Guide](https://buf.build/docs/best-practices/style-guide/)

## Coding Style

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
