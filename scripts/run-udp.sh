#!/bin/bash

# Run the Flightpath server with UDP configuration
export FLIGHTPATH_MAVLINK_ENDPOINT_TYPE=udp-server
export FLIGHTPATH_MAVLINK_UDP_ADDRESS=0.0.0.0:14550
go run cmd/server/main.go
