# PX4 SITL Setup

The easiest way to try out Flightpath is to run one of the examples against a simulated version of the drone. This document describes how to run a headless PX4 SITL (Software in the Loop) in a docker container. See [Gazebo Classic Simulation](https://docs.px4.io/main/en/sim_gazebo_classic/) for the details on this SITL.

Instead of Docker Desktop, we will use the [Colima](https://github.com/abiosoft/colima) container runtime on MacOS. It is a lightweight equivalent of Docker Desktop with a very easy setup. Also, the latest Docker Desktop install on my Mac was interfering with my virus protection software, so I decided not to use it.

Follow the steps below to install and run the PX4 SITL.

### Install Homebrew

See latest install docs [here](https://brew.sh/).

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

### Install Colima

```bash
brew install colima
```

### Install Docker Runtime

On initial startup, Colima initiates with a user specified runtime that defaults to Docker. Hence the docker runtime is required. Install it using the following command:

```bash
brew install docker
```

### Start Colima

```bash
colima start
```

### Run PX4 SITL in Docker

Find the latest version of the PX4 SITL on [Docker Hub](https://hub.docker.com/) – search for "jonasvautherin/px4-gazebo-headless". As of this writing the latest version is 1.16.0.

Run the image:

```bash
# --rm: Automatically remove the container when it exits (cleanup)
# -it: Combined flags
#   -i: interactive: keep STDIN open
#   -t: allocate a pseudo-TTY for terminal interaction
docker run --rm -it jonasvautherin/px4-gazebo-headless:1.16.0
```

This will start the PX4 SITL with an interactive terminal:

```
______  __   __    ___
| ___ \ \ \ / /   /   |
| |_/ /  \ V /   / /| |
|  __/   /   \  / /_| |
| |     / /^\ \ \___  |
\_|     \/   \/     |_/

px4 starting.

INFO  [px4] startup script: /bin/sh etc/init.d-posix/rcS 0
Warning [parser.cc:833] XML Attribute[version] in element[sdf] not defined in SDF, ignoring.
INFO  [init] found model autostart file as SYS_AUTOSTART=10015
INFO  [param] selected parameter default file parameters.bson
INFO  [param] selected parameter backup file parameters_backup.bson
  SYS_AUTOCONFIG: curr: 0 -> new: 1
  SYS_AUTOSTART: curr: 0 -> new: 10015
  CAL_ACC0_ID: curr: 0 -> new: 1310988
  CAL_GYRO0_ID: curr: 0 -> new: 1310988
  CAL_ACC1_ID: curr: 0 -> new: 1310996
  CAL_GYRO1_ID: curr: 0 -> new: 1310996
  CAL_ACC2_ID: curr: 0 -> new: 1311004
  CAL_GYRO2_ID: curr: 0 -> new: 1311004
  CAL_MAG0_ID: curr: 0 -> new: 197388
  CAL_MAG0_PRIO: curr: -1 -> new: 50
  CAL_MAG1_ID: curr: 0 -> new: 197644
  CAL_MAG1_PRIO: curr: -1 -> new: 50
  SENS_BOARD_X_OFF: curr: 0.0000 -> new: 0.0000
  SENS_DPRES_OFF: curr: 0.0000 -> new: 0.0010
INFO  [dataman] data manager file './dataman' size is 1208528 bytes
INFO  [init] PX4_SIM_HOSTNAME: localhost
INFO  [simulator_mavlink] Waiting for simulator to accept connection on TCP port 4560
INFO  [simulator_mavlink] Simulator connected on TCP port 4560.
INFO  [lockstep_scheduler] setting initial absolute time to 196000 us
INFO  [commander] LED: open /dev/led0 failed (22)
INFO  [uxrce_dds_client] init UDP agent IP:127.0.0.1, port:8888
INFO  [mavlink] mode: Normal, data rate: 4000000 B/s on udp port 18570 remote port 14550
INFO  [mavlink] mode: Onboard, data rate: 4000000 B/s on udp port 14580 remote port 14540
INFO  [mavlink] mode: Onboard, data rate: 4000 B/s on udp port 14280 remote port 14030
INFO  [mavlink] mode: Gimbal, data rate: 400000 B/s on udp port 13030 remote port 13280
INFO  [logger] logger started (mode=all)
INFO  [logger] Start file log (type: full)
INFO  [logger] [logger] ./log/2025-12-09/20_06_19.ulg
INFO  [logger] Opened full log file: ./log/2025-12-09/20_06_19.ulg
INFO  [mavlink] MAVLink only on localhost (set param MAV_{i}_BROADCAST = 1 to enable network)
INFO  [mavlink] MAVLink only on localhost (set param MAV_{i}_BROADCAST = 1 to enable network)
INFO  [px4] Startup script returned successfully
pxh> INFO  [tone_alarm] home set
INFO  [commander] Ready for takeoff!

pxh>
```

You can enter commands on the terminal and the SITL will execute them, e.g.

```
pxh> commander takeoff
pxh> commander land
pxh> help             <--- list of all available commands
pxh> commander help   <--- usage info for commands in the commander module
```

Note that `commander` is one of the many modules available in PX4. The full list can be found in the [PX4 Modules & Commands Reference](https://docs.px4.io/main/en/modules/modules_main).

## SITL Simulation Environment
Reference: https://docs.px4.io/main/en/simulation/#sitl-simulation-environment

From the log above:

```
INFO  [mavlink] mode: Normal, data rate: 4000000 B/s on udp port 18570 remote port 14550
```

Communication flow:

```
     ┌──────────────┐                     ┌──────────────┐
     │     PX4      │                     │     GCS      │
     │   on SITL    │                     │ (Flightpath) │
     └──────┬───────┘                     └──────┬───────┘
            │                                    │
            │     broadcast messages from PX4    │
            │  ────────────────────────────────> │  Port 14550
            │                                    │
            │          commands from GCS         │
 Port 18570 |  <───────────────────────────────- │
            │                                    │
            │                                    │
```

## Useful Colima Commands

```bash
# Starts Colima with default settings, typically launching a Docker runtime
colima start

# Stops the running Colima instance
colima stop

# Displays information about all Colima instances, including their status, architecture, and resource allocation
colima list

# Shows the current status of the default Colima instance
colima status
```

## Useful Docker Commands

```bash
# Show all images on this machine
docker images -a

# Show all running containers
docker ps -a

# Stop a container
docker stop <hash>

# Remove the specified container from this machine
docker rm <hash>

# Remove the specified image from this machine
docker rmi <image-id>
```
