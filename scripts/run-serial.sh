#!/bin/bash

# Run the Flightpath server with serial port configuration
export FLIGHTPATH_MAVLINK_ENDPOINT_TYPE=serial
export FLIGHTPATH_MAVLINK_SERIAL_DEVICE=/dev/cu.usbserial-D30JAXGS
export FLIGHTPATH_MAVLINK_SERIAL_BAUD=57600
go run cmd/server/main.go
