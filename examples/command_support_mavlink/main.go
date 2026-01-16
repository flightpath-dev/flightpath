package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bluenviron/gomavlib/v3"
	"github.com/bluenviron/gomavlib/v3/pkg/dialects/common"
	"github.com/flightpath-dev/flightpath/internal/config"
	"github.com/flightpath-dev/flightpath/internal/mavlink"
)

// ------------------------------------------------------------------------------------------------
// Command Support Probe using MAVLink
// ------------------------------------------------------------------------------------------------
// This example probes the connected flight controller to discover which MAVLink commands
// are supported by sending each command while the vehicle is disarmed and analyzing the
// COMMAND_ACK responses.
//
// Safety: All commands are sent while the vehicle is disarmed. The flight controller will
// perform logic checks and respond with acceptance or denial, but will not actuate motors
// in a disarmed state.
//
// Configuration is loaded from environment variables with sensible defaults:
//   - Default: UDP server on port 14550 (standard PX4 SITL port)
//   - See config.Load() function for all available environment variables
//
// To run this example:
//  1. Start a PX4 SITL (see docs/px4-sitl-setup.md)
//
//  2. Test all commands using the default configuration (MAVLink running as a UDP server on port 14550)
//     go run examples/command_support_mavlink/main.go
//
//  3. Test a single command by ID:
//     go run examples/command_support_mavlink/main.go 400
//
//  4. Or configure a serial connection via environment variables:
//     export FLIGHTPATH_MAVLINK_ENDPOINT_TYPE=serial
//     export FLIGHTPATH_MAVLINK_SERIAL_DEVICE=/dev/cu.usbserial-D30JAXGS
//     export FLIGHTPATH_MAVLINK_SERIAL_BAUD=57600
//
//     go run examples/command_support_mavlink/main.go
// ------------------------------------------------------------------------------------------------

const (
	// CommandTimeout is the timeout for waiting for COMMAND_ACK response
	CommandTimeout = 2 * time.Second
	// HeartbeatTimeout is the timeout for waiting for initial heartbeat
	HeartbeatTimeout = 60 * time.Second
	// DelayBetweenCommands is the delay between sending commands to avoid overwhelming the FC
	DelayBetweenCommands = 500 * time.Millisecond
)

// CommandResult represents the result of testing a command
type CommandResult struct {
	CommandID common.MAV_CMD
	Result    common.MAV_RESULT
	HasAck    bool
}

func main() {
	// Step 1: Parse command-line arguments
	var singleCommandID *common.MAV_CMD
	if len(os.Args) > 1 {
		cmdIDStr := os.Args[1]
		cmdIDUint, err := strconv.ParseUint(cmdIDStr, 10, 64)
		if err != nil {
			log.Fatalf("Invalid command ID: %s (must be a number)", cmdIDStr)
		}
		cmdID := common.MAV_CMD(cmdIDUint)

		// Validate that this is a known command
		allCommands := mavlink.GetAllCommands()
		found := false
		for _, c := range allCommands {
			if c == cmdID {
				found = true
				break
			}
		}
		if !found {
			log.Fatalf("Command ID %d is not a known MAV_CMD in the common dialect", cmdID)
		}

		singleCommandID = &cmdID
	}

	// Step 2: Display safety checklist and wait for confirmation
	printSafetyChecklist()

	// Step 3: Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Step 4: Initialize MAVLink node
	log.Println("Initializing MAVLink node...")
	node := &gomavlib.Node{
		Endpoints:   []gomavlib.EndpointConf{cfg.MAVLink.Endpoint},
		Dialect:     common.Dialect,
		OutVersion:  gomavlib.V2,
		OutSystemID: 254,
	}
	err = node.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize node: %v", err)
	}
	defer node.Close()
	log.Println("Node initialized successfully")

	// Step 5: Wait for heartbeat to discover channel and system/component IDs
	heartbeatChan, heartbeatSystemID, heartbeatComponentID := detectChannel(node)

	// Step 6: Wait a moment for drone to recognize that GCS is connected
	log.Println("Waiting 2 seconds for drone to recognize GCS connection...")
	sleepWhileListening(node, 2*time.Second)
	log.Println("GCS connection established")
	fmt.Println()

	// Step 7: Test single or all commands
	if singleCommandID != nil {
		// Test single command
		log.Printf("Testing single command on system %d, component %d...\n", heartbeatSystemID, heartbeatComponentID)
		fmt.Println()

		result := testCommand(node, heartbeatChan, heartbeatSystemID, heartbeatComponentID, *singleCommandID)
		printSingleCommandResult(result)

		log.Println("Example completed successfully")
	} else {
		// Test all commands
		commands := mavlink.GetAllCommands()
		log.Printf("Testing commands on system %d, component %d...\n", heartbeatSystemID, heartbeatComponentID)
		log.Printf("Testing %d MAV_CMD commands from common dialect\n", len(commands))
		estimatedTime := time.Duration(len(commands)) * DelayBetweenCommands
		log.Printf("Estimated time: %v (at %v per command)\n", estimatedTime.Round(time.Second), DelayBetweenCommands)
		fmt.Println()

		results := []CommandResult{}
		totalCommands := len(commands)
		startTime := time.Now()
		lastProgressTime := startTime

		for i, cmdID := range commands {
			now := time.Now()

			// Show progress every 10 commands or every 10 seconds
			showProgress := i%10 == 0 || now.Sub(lastProgressTime) >= 10*time.Second

			if showProgress {
				elapsed := now.Sub(startTime)
				if i > 0 {
					rate := float64(i) / elapsed.Seconds()
					remaining := time.Duration(float64(totalCommands-i)/rate) * time.Second
					fmt.Printf("[%3d/%3d] Testing %s (%d) - Elapsed: %v, Remaining: ~%v\r",
						i+1, totalCommands, cmdID.String(), cmdID,
						elapsed.Round(time.Second), remaining.Round(time.Second))
				} else {
					fmt.Printf("[%3d/%3d] Testing %s (%d)...\r",
						i+1, totalCommands, cmdID.String(), cmdID)
				}
				lastProgressTime = now
			}

			result := testCommand(node, heartbeatChan, heartbeatSystemID, heartbeatComponentID, cmdID)
			results = append(results, result)

			// Small delay to avoid overwhelming the FC
			time.Sleep(DelayBetweenCommands)
		}

		fmt.Println() // Clear the progress line
		fmt.Println()

		// Print results
		printResults(results)

		log.Println("Example completed successfully")
	}
}

// printSafetyChecklist
// Displays the pre-test safety checklist and waits for user confirmation.
func printSafetyChecklist() {
	fmt.Println("==================================================================================")
	fmt.Println("                        PRE-TEST SAFETY CHECKLIST")
	fmt.Println("==================================================================================")
	fmt.Println()
	fmt.Println("Before proceeding, verify the following safety measures are in place:")
	fmt.Println()
	fmt.Println("  1. PROPELLERS REMOVED")
	fmt.Println("     This is the only 100% effective safeguard against injury or damage.")
	fmt.Println("     Ensure all propellers have been physically removed from the vehicle.")
	fmt.Println()
	fmt.Println("  2. HARDWARE SAFETY SWITCH DISENGAGED")
	fmt.Println("     If your flight controller has a physical safety switch (standard on")
	fmt.Println("     Pixhawk-style hardware), verify it remains disengaged. When disengaged,")
	fmt.Println("     motor output is inhibited regardless of software commands.")
	fmt.Println()
	fmt.Println("NOTICE: Testing will send MAVLink commands to the flight controller while")
	fmt.Println("disarmed. The FC will perform logic checks and respond with acceptance or")
	fmt.Println("denial, but will not actuate motors in a disarmed state.")
	fmt.Println()
	fmt.Println("==================================================================================")
	fmt.Println()
	fmt.Print("Press ENTER to confirm safety checklist and begin testing, or Ctrl+C to abort: ")

	reader := bufio.NewReader(os.Stdin)
	_, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}
	fmt.Println()
}

// detectChannel
// Detects the channel by waiting for a heartbeat message and returns the channel, system ID, and component ID.
func detectChannel(node *gomavlib.Node) (*gomavlib.Channel, uint8, uint8) {
	log.Println("Waiting for heartbeat from flight controller...")
	t := time.NewTimer(HeartbeatTimeout)
	defer t.Stop()

	for {
		select {
		case evt := <-node.Events():
			if evt, ok := evt.(*gomavlib.EventFrame); ok {
				if _, ok := evt.Message().(*common.MessageHeartbeat); ok {
					channel := evt.Channel
					systemID := evt.SystemID()
					componentID := evt.ComponentID()
					log.Printf("Received heartbeat from system %d, component %d\n", systemID, componentID)
					return channel, systemID, componentID
				}
			}
		case <-t.C:
			log.Fatalf("Timeout waiting for heartbeat after %v", HeartbeatTimeout)
		}
	}
}

// testCommand
// Sends a command and waits for ACK with timeout, returns the result.
func testCommand(
	node *gomavlib.Node,
	channel *gomavlib.Channel,
	targetSystemID, targetComponentID uint8,
	commandID common.MAV_CMD,
) CommandResult {
	cmd := &common.MessageCommandLong{
		TargetSystem:    targetSystemID,
		TargetComponent: targetComponentID,
		Command:         commandID,
		Confirmation:    0,
		Param1:          0,
		Param2:          0,
		Param3:          0,
		Param4:          0,
		Param5:          0,
		Param6:          0,
		Param7:          0,
	}

	err := node.WriteMessageTo(channel, cmd)
	if err != nil {
		log.Printf("Warning: Failed to write command %d: %v\n", commandID, err)
		return CommandResult{
			CommandID: commandID,
			Result:    common.MAV_RESULT_FAILED,
			HasAck:    false,
		}
	}

	t := time.NewTimer(CommandTimeout)
	defer t.Stop()

	for {
		select {
		case evt := <-node.Events():
			if evt, ok := evt.(*gomavlib.EventFrame); ok {
				if ack, ok2 := evt.Message().(*common.MessageCommandAck); ok2 {
					if ack.Command == commandID &&
						evt.SystemID() == targetSystemID &&
						evt.ComponentID() == targetComponentID {
						return CommandResult{
							CommandID: commandID,
							Result:    ack.Result,
							HasAck:    true,
						}
					}
				}
			}

		case <-t.C:
			return CommandResult{
				CommandID: commandID,
				Result:    common.MAV_RESULT_FAILED,
				HasAck:    false,
			}
		}
	}
}

// resultToString
// Converts MAV_RESULT to human-readable status string.
func resultToString(result CommandResult) string {
	if !result.HasAck {
		return "NO RESPONSE"
	}

	switch result.Result {
	case common.MAV_RESULT_ACCEPTED:
		return "ACCEPTED"
	case common.MAV_RESULT_TEMPORARILY_REJECTED:
		return "TEMPORARILY_REJECTED"
	case common.MAV_RESULT_DENIED:
		return "DENIED"
	case common.MAV_RESULT_UNSUPPORTED:
		return "UNSUPPORTED"
	case common.MAV_RESULT_FAILED:
		return "FAILED"
	case common.MAV_RESULT_IN_PROGRESS:
		return "IN_PROGRESS"
	case common.MAV_RESULT_CANCELLED:
		return "CANCELLED"
	case common.MAV_RESULT_COMMAND_LONG_ONLY:
		return "COMMAND_LONG_ONLY"
	case common.MAV_RESULT_COMMAND_INT_ONLY:
		return "COMMAND_INT_ONLY"
	case common.MAV_RESULT_COMMAND_UNSUPPORTED_MAV_FRAME:
		return "UNSUPPORTED_MAV_FRAME"
	case common.MAV_RESULT_NOT_IN_CONTROL:
		return "NOT_IN_CONTROL"
	default:
		return fmt.Sprintf("OTHER(%d)", result.Result)
	}
}

// isSupported
// Returns true if the result indicates the command is supported by the FC.
// Supported means the FC knows about the command (even if it can't execute it right now).
func isSupported(result common.MAV_RESULT) bool {
	switch result {
	case common.MAV_RESULT_ACCEPTED,
		common.MAV_RESULT_TEMPORARILY_REJECTED,
		common.MAV_RESULT_DENIED,
		common.MAV_RESULT_IN_PROGRESS,
		common.MAV_RESULT_CANCELLED,
		common.MAV_RESULT_COMMAND_LONG_ONLY,
		common.MAV_RESULT_COMMAND_INT_ONLY,
		common.MAV_RESULT_NOT_IN_CONTROL:
		return true
	default:
		return false
	}
}

// printSingleCommandResult
// Prints the result for a single command test in a simple format.
func printSingleCommandResult(result CommandResult) {
	fmt.Println("==================================================================================")
	fmt.Println("                              TEST RESULT")
	fmt.Println("==================================================================================")
	fmt.Println()

	var category string
	if !result.HasAck {
		category = "UNKNOWN"
	} else if isSupported(result.Result) {
		category = "SUPPORTED"
	} else {
		category = "UNSUPPORTED"
	}

	fmt.Printf("Command: %s (%d)\n", result.CommandID.String(), result.CommandID)
	fmt.Printf("Status:  %s\n", category)
	if result.HasAck {
		fmt.Printf("Result:  %s\n", resultToString(result))
	} else {
		fmt.Printf("Result:  NO RESPONSE (timeout)\n")
	}

	fmt.Println()
	fmt.Println("==================================================================================")
}

// printResults
// Formats and prints the categorized results table with summary.
// Results are categorized as:
//   - Supported: FC responded with ACCEPTED, DENIED, TEMPORARILY_REJECTED, etc.
//   - Unsupported: FC responded with UNSUPPORTED, FAILED, or COMMAND_UNSUPPORTED_MAV_FRAME
//   - Unknown: No response (timeout)
func printResults(results []CommandResult) {
	// Categorize results
	supported := []CommandResult{}
	unsupported := []CommandResult{}
	unknown := []CommandResult{}

	for _, result := range results {
		if !result.HasAck {
			unknown = append(unknown, result)
		} else if isSupported(result.Result) {
			supported = append(supported, result)
		} else {
			unsupported = append(unsupported, result)
		}
	}

	// Sort by command ID for consistent output
	sort.Slice(supported, func(i, j int) bool {
		return supported[i].CommandID < supported[j].CommandID
	})
	sort.Slice(unsupported, func(i, j int) bool {
		return unsupported[i].CommandID < unsupported[j].CommandID
	})
	sort.Slice(unknown, func(i, j int) bool {
		return unknown[i].CommandID < unknown[j].CommandID
	})

	fmt.Println()
	fmt.Println("==================================================================================")
	fmt.Println("                              TEST RESULTS")
	fmt.Println("==================================================================================")
	fmt.Println()

	// Print supported commands
	fmt.Printf("SUPPORTED: %d commands\n", len(supported))
	fmt.Println(strings.Repeat("-", 80))
	if len(supported) > 0 {
		for _, result := range supported {
			fmt.Printf("  %s (%d) - %s\n", result.CommandID.String(), result.CommandID, resultToString(result))
		}
	} else {
		fmt.Println("  (none)")
	}
	fmt.Println()

	// Print unsupported commands
	fmt.Printf("UNSUPPORTED: %d commands\n", len(unsupported))
	fmt.Println(strings.Repeat("-", 80))
	if len(unsupported) > 0 {
		for _, result := range unsupported {
			fmt.Printf("  %s (%d) - %s\n", result.CommandID.String(), result.CommandID, resultToString(result))
		}
	} else {
		fmt.Println("  (none)")
	}
	fmt.Println()

	// Print unknown commands (no response)
	fmt.Printf("UNKNOWN (NO RESPONSE): %d commands\n", len(unknown))
	fmt.Println(strings.Repeat("-", 80))
	if len(unknown) > 0 {
		for _, result := range unknown {
			fmt.Printf("  %s (%d)\n", result.CommandID.String(), result.CommandID)
		}
	} else {
		fmt.Println("  (none)")
	}
	fmt.Println()

	// Print summary
	fmt.Println("==================================================================================")
	fmt.Printf("SUMMARY: %d supported, %d unsupported, %d unknown (%d total)\n",
		len(supported), len(unsupported), len(unknown), len(results))
	fmt.Println("==================================================================================")
}

// sleepWhileListening
// Sleeps for the specified duration while continuing to process node events.
func sleepWhileListening(node *gomavlib.Node, d time.Duration) {
	t := time.NewTimer(d)
	defer t.Stop()

	for {
		select {
		case <-node.Events():
		case <-t.C:
			return
		}
	}
}
