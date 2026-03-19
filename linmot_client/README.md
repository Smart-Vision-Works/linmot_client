# linmot_client

> **Module Description:** Go implementation of the LinUDP V2 protocol for LinMot C1250 servo drives. Full high-level client API with multi-drive pooling, YAML-loadable command tables with template variables, RTC parameter access, motion control (VAI/streaming/predefined), 4-channel monitoring, and a complete mock drive for unit and E2E testing.

A pure-Go implementation of the **LinUDP V2** communication protocol for LinMot C1250-series servo drives. Covers everything from raw packet encoding to a high-level `Client` API that handles connectivity validation, state-machine transitions, motion commands, parameter reads/writes, command table management, and multi-drive connection pooling.

---

## Table of Contents

1. [Module Layout](#1-module-layout)
2. [Architecture](#2-architecture)
   - 2.1 [Layers at a Glance](#21-layers-at-a-glance)
   - 2.2 [Transport](#22-transport)
   - 2.3 [Protocol](#23-protocol)
   - 2.4 [Client](#24-client)
3. [Quick Start](#3-quick-start)
4. [Key Concepts](#4-key-concepts)
   - 4.1 [LinUDP V2 Protocol](#41-linudp-v2-protocol)
   - 4.2 [Engineering Units](#42-engineering-units)
   - 4.3 [Storage Types — RAM vs. ROM](#43-storage-types--ram-vs-rom)
   - 4.4 [UPID Parameters](#44-upid-parameters)
   - 4.5 [Command Tables](#45-command-tables)
5. [Client API Reference](#5-client-api-reference)
   - 5.1 [Creating a Client](#51-creating-a-client)
   - 5.2 [State Machine Control](#52-state-machine-control)
   - 5.3 [Motion Control](#53-motion-control)
   - 5.4 [RTC Parameters](#54-rtc-parameters)
   - 5.5 [Command Tables](#55-command-tables)
   - 5.6 [Monitoring Channels](#56-monitoring-channels)
   - 5.7 [Curves](#57-curves)
   - 5.8 [Error Log](#58-error-log)
   - 5.9 [Drive Operations](#59-drive-operations)
6. [Multi-Drive: ClientPool](#6-multi-drive-clientpool)
7. [Motion Control Interfaces](#7-motion-control-interfaces)
8. [Testing Infrastructure](#8-testing-infrastructure)
   - 8.1 [MockLinMot — Full Drive Simulator](#81-mocklinmot--full-drive-simulator)
   - 8.2 [Running Tests](#82-running-tests)
9. [Reference Documentation](#9-reference-documentation)
10. [Known Issues & Quirks](#10-known-issues--quirks)

---

## 1. Module Layout

```
linmot_client/
├── client.go                          # Top-level convenience: NewClient(driveIP)
│
├── client/                            # High-level Client API
│   ├── client.go                      # Client struct, NewUDPClient, NewClientWithTransport
│   ├── pool.go                        # ClientPool — multi-drive connection cache
│   ├── operational.go                 # ErrDriveNotOperational, state name lookup
│   ├── faults.go                      # DriveFaultError, CheckDriveFault
│   │
│   ├── common/
│   │   ├── request_manager.go         # Core TX/RX goroutine loop, request dispatch
│   │   ├── pending_request.go         # In-flight request with response channel
│   │   ├── rtc_counter.go             # RTC 8-bit cyclic counter management
│   │   └── errors.go                  # Sentinel errors (ErrInvalidStatusTelegram, etc.)
│   │
│   ├── control_word/
│   │   └── manager.go                 # EnableDrive, DisableDrive, Home, QuickStop, AcknowledgeError
│   │
│   ├── monitoring/
│   │   └── manager.go                 # ConfigureChannel(s), GetMonitoringData
│   │
│   ├── motion_control/
│   │   ├── manager.go                 # MotionControlManager — VAI, Predefined, VAI16, Streaming
│   │   ├── vai/manager.go             # GoToPosition, IncrementDemandPosition, Stop, etc.
│   │   ├── vai16/manager.go           # 16-bit VAI variant
│   │   ├── predefined/manager.go      # Execute pre-stored command table entries
│   │   ├── interface_control/manager.go  # Direct interface control commands
│   │   └── streaming/manager.go       # Streaming motion commands
│   │
│   └── rtc/
│       ├── manager.go                 # RtcManager — orchestrates all RTC sub-managers
│       ├── command_tables/
│       │   ├── manager.go             # SetCommandTable, GetCommandTable, Stop/StartMC, SaveToFlash
│       │   ├── command_table.go       # CommandTable & Entry types, YAML load/parse, Validate
│       │   └── testdata/              # YAML fixtures: linmot_command_table, pick_sequence, templates
│       ├── curves/manager.go          # UploadCurve, DownloadCurve, ModifyCurve, Save/DeleteAll
│       ├── errors/manager.go          # GetErrorLog, GetErrorText, GetErrorLogWithText
│       ├── operations/manager.go      # RestartDrive, SetOS/MC/Interface/AppROMToDefault
│       └── parameters/manager.go      # ReadRAM, WriteRAMAndROM, GetAll/ModifiedUPIDs, motion params
│
├── protocol/                          # Pure packet encoding/decoding — no I/O
│   ├── common/
│   │   ├── constants.go               # RequestFlags, ResponseFlags, packet sizes
│   │   ├── packet.go                  # Packet framing: header, request, response assembly
│   │   ├── requests.go                # NewStatusRequest, NewConnectivityProbeRequest
│   │   ├── responses.go               # StatusResponse — parse raw bytes into Status struct
│   │   ├── status.go                  # Status struct: StatusWord, StateVar, Position, WarnWord, etc.
│   │   ├── units.go                   # ToProtocol*/FromProtocol* unit converters
│   │   ├── offsets.go                 # Byte offset constants for packet fields
│   │   ├── writable.go                # Writable interface for encodable request payloads
│   │   └── timeouts.go                # Timeout constants
│   │
│   ├── control_word/
│   │   ├── builder.go                 # Control word bit construction
│   │   ├── constants.go               # MainState, control bit masks
│   │   ├── requests.go                # Control word request encoding
│   │   └── status.go                  # Control word status decoding
│   │
│   ├── motion_control/                # 15 MC interface sub-packages
│   │   ├── constants.go               # InterfaceID values
│   │   ├── header.go                  # MC frame header
│   │   ├── counter.go                 # MC sequence counter
│   │   ├── requests.go / responses.go # Generic MC request/response framing
│   │   ├── vai/                       # VAI commands (GoToPos, IncrementDemand, etc.)
│   │   ├── vai16/                     # 16-bit VAI variant
│   │   ├── predefined/                # Predefined command execution
│   │   ├── predef_vai16/              # Predefined with 16-bit params
│   │   ├── streaming/                 # Streaming motion
│   │   ├── interface_control/         # Interface control commands
│   │   ├── advanced/                  # Advanced commands
│   │   ├── bestehorn/                 # Bestehorn cam profile
│   │   ├── cams/                      # Cam commands
│   │   ├── encoder_cams/              # Encoder cam commands
│   │   ├── indexing/                  # Indexing commands
│   │   ├── vai_dec_acc/               # VAI with separate deceleration/acceleration
│   │   ├── vai_positioning/           # VAI positioning variant
│   │   └── vai_predef_acc/            # VAI predefined with acceleration
│   │
│   └── rtc/                           # RTC frame encoding/decoding
│       ├── constants.go               # RTC command codes
│       ├── requests.go / responses.go # Generic RTC request/response framing
│       ├── response_registry.go       # Maps RTC command codes to response parsers
│       ├── validation.go              # Response validation helpers
│       ├── command_tables/            # CT read/write/presence-mask protocol
│       ├── curves/                    # Curve upload/download protocol
│       ├── errors/                    # Error log + error text protocol
│       ├── operations/                # Restart/reset-to-default protocol
│       └── parameters/                # UPID read/write/enumerate protocol
│
├── transport/                         # UDP transport layer
│   ├── client.go                      # Client interface (SendPacket, RecvPacket, Close)
│   ├── server.go                      # Server interface (for mock drives)
│   ├── constants.go                   # DefaultDrivePort=49360, DefaultMasterPort=41136, DefaultTimeout=1s
│   ├── udp_client.go                  # Single-drive UDP transport
│   ├── shared_udp.go                  # SharedUDPTransport — one socket, many drives
│   ├── shared_client.go               # Per-drive client on shared transport
│   ├── cyclic_probe.go                # CyclicProbe — C#-style connectivity fallback
│   ├── debug.go                       # Debug helpers
│   ├── tracer.go                      # Packet trace logging
│   ├── mock_client.go                 # In-memory mock transport (client side)
│   ├── mock_server.go                 # In-memory mock transport (server side)
│   └── mock_config.go                 # MockTransport configuration
│
├── test/
│   └── mock_linmot.go                 # MockLinMot — full state-machine drive simulator
│
└── reference/                              # Official LinMot documentation
    ├── README.md                      # Doc index, key reference points, drive configuration table
    ├── C1250_MI_Installation_Guide.md # Hardware guide: connectors, DIP switches, LED blink codes
    ├── LinUDP_V2.md                   # Protocol manual: packet format, UPID list, monitoring channels
    ├── LinUDP_V2_DLL.md               # DLL integration guide: commissioning, static IP config path
    ├── LinMotTalk_Manual.md           # LinMot-Talk manual: factory reset (§4.1), parameter editing
    ├── LinMot_MotionCtrl.md           # Motion Control manual: Easy Steps, command tables, position controller
    └── decompiled_linudp_csharp_lib.cs  # Decompiled official C# SDK (canonical reference)
```

---

## 2. Architecture

### 2.1 Layers at a Glance

```
┌─────────────────────────────────────────────────────────────┐
│                   Your application               │
├─────────────────────────────────────────────────────────────┤
│  linmot.NewClient(ip)  ─────────────────────────────────────┤  ← top-level entry point
│  client.Client  (client.go)                                  │
│    ├── ControlWordManager  (enable, disable, home, ack)      │
│    ├── RtcManager  (params, command tables, curves, errors)  │
│    ├── MonitoringManager  (4-channel live variable monitor)  │
│    └── MotionControlManager  (VAI, predefined, streaming)    │
├─────────────────────────────────────────────────────────────┤
│  client/common.RequestManager                                │
│    TX goroutine → encodes & sends requests                   │
│    RX goroutine → receives responses, matches to pending     │
├─────────────────────────────────────────────────────────────┤
│  protocol/  (pure encoding/decoding, no I/O)                 │
│    common  ·  control_word  ·  motion_control  ·  rtc       │
├─────────────────────────────────────────────────────────────┤
│  transport/  (UDP I/O)                                        │
│    UDPTransportClient  ·  SharedUDPTransport  ·  MockClient  │
└─────────────────────────────────────────────────────────────┘
```

### 2.2 Transport

The transport layer owns the UDP socket(s) and exposes a simple two-method interface:

```go
type Client interface {
    SendPacket(ctx context.Context, packet []byte) error
    RecvPacket(ctx context.Context) ([]byte, error)
    Close() error
}
```

| Transport                                         | When to use                                                                                       |
| ------------------------------------------------- | ------------------------------------------------------------------------------------------------- |
| [`UDPTransportClient`](transport/udp_client.go)   | Single drive, dedicated socket per client                                                         |
| [`SharedUDPTransport`](transport/shared_udp.go)   | Multiple drives from **one** master port (41136) — required when running >1 client simultaneously |
| [`MockTransportClient`](transport/mock_client.go) | Unit tests — in-memory channels, no network                                                       |

**Why SharedUDPTransport?** LinMot drives filter incoming packets by source port. All clients must originate from port 41136. You cannot bind two sockets to the same port. `SharedUDPTransport` owns one socket and routes received packets to the correct per-drive channel by source IP.

### 2.3 Protocol

`protocol/` is entirely pure Go — no I/O, no goroutines. It handles:

- **Packet framing:** 8-byte header (4-byte request flags + 4-byte response flags), then optional ControlWord (2 bytes), MotionControl (32 bytes), and/or RTC command (8 bytes)
- **Response parsing:** Status fields (StatusWord, StateVar, ActualPosition, DemandPosition, Current, WarnWord, ErrorCode), plus optional MonitoringChannel and RTC reply data
- **Unit conversion:** see [§4.2 Engineering Units](#42-engineering-units)

### 2.4 Client

`client.Client` composes four specialized managers:

| Manager                | Responsibility                                                                              |
| ---------------------- | ------------------------------------------------------------------------------------------- |
| `ControlWordManager`   | State machine transitions: enable, disable, home, quick-stop, error-acknowledge             |
| `RtcManager`           | All RTC channel operations: parameters, command tables, curves, error log, drive operations |
| `MonitoringManager`    | Configure and read 4-channel live variable monitoring                                       |
| `MotionControlManager` | All MC interface commands: VAI, predefined, VAI16, streaming, interface control             |

All managers share a single `RequestManager`, which serialises requests to the drive via its TX/RX goroutine pair. The `RequestManager` starts on client creation and stops when `Close()` is called.

---

## 3. Quick Start

```go
import "github.com/Smart-Vision-Works/staged_robot"

// Single drive — uses default ports (drive:49360, master:41136)
c, err := linmot.NewClient("10.8.7.232")
if err != nil {
    log.Fatal(err) // includes connectivity validation
}
defer c.Close()

ctx := context.Background()

// Enable drive
if _, err := c.EnableDrive(ctx); err != nil {
    log.Fatal(err)
}

// Move to 10 mm at 0.1 m/s, 5 m/s² accel/decel
status, err := c.VAIGoToPosition(ctx, 10.0, 0.1, 5.0, 5.0)

// Read actual position
positionMM, err := c.GetPosition(ctx)

// Write a parameter to RAM + ROM
err = c.WriteRAMAndROM(ctx, 0x1450, 100000) // e.g., position in 0.1µm units
```

**Multiple drives (recommended for production):**

```go
pool := client.NewClientPool()
defer pool.Close()

c0, err := pool.GetClient("10.8.7.232")
c1, err := pool.GetClient("10.8.7.234")
// Both share one SharedUDPTransport on port 41136
```

---

## 4. Key Concepts

### 4.1 LinUDP V2 Protocol

LinUDP V2 is a UDP-based binary protocol defined by NTI AG (LinMot). Key facts:

- **Drive port:** 49360 (drive listens here)
- **Master port:** 41136 (all packets must originate from this port)
- **Packet structure:** 8-byte header + optional ControlWord (2B) + optional MotionControl (32B) + optional RTC command (8B)
- **Response structure:** 8-byte header + 20-byte status block + optional MonitoringChannel (8B) + optional RTC reply
- **Request flags** (bits in the first 4 bytes of the header):

  | Bit | Flag          | Payload size |
  | --- | ------------- | ------------ |
  | 0   | ControlWord   | 2 bytes      |
  | 1   | MotionControl | 32 bytes     |
  | 2   | RTCCommand    | 8 bytes      |

- **Response flags** (bits in bytes 4–7 of the header):

  | Bits | Flag              | Description                                        |
  | ---- | ----------------- | -------------------------------------------------- |
  | 0–6  | Standard          | Status fields (always present in normal responses) |
  | 7    | MonitoringChannel | Inline 4-channel monitoring data                   |
  | 8    | RTCReply          | RTC response payload                               |

Reference: [`protocol/common/constants.go`](protocol/common/constants.go), [`reference/LinUDP_V2.md`](reference/LinUDP_V2.md)

### 4.2 Engineering Units

All public API methods accept/return **engineering units** (mm, m/s, m/s²). The protocol uses integer units internally.

| Quantity     | API unit | Protocol unit | Conversion factor |
| ------------ | -------- | ------------- | ----------------- |
| Position     | `mm`     | 0.1 µm        | × 10,000          |
| Velocity     | `m/s`    | 1 µm/s        | × 1,000,000       |
| Acceleration | `m/s²`   | 10 µm/s²      | × 100,000         |
| Jerk         | `m/s³`   | 1 µm/s³       | × 1,000,000       |

Converters: [`protocol/common/units.go`](protocol/common/units.go)

```go
// Examples
protocol_common.ToProtocolPosition(10.5)   // → 105000 (10.5 mm)
protocol_common.ToProtocolVelocity(0.5)    // → 500000 (0.5 m/s)
protocol_common.ToProtocolAcceleration(2.5) // → 250000 (2.5 m/s²)
```

### 4.3 Storage Types — RAM vs. ROM

Every parameter write takes a `ParameterStorageType`:

| Type             | Behaviour               | Survives power cycle? |
| ---------------- | ----------------------- | --------------------- |
| `WriteRAM`       | Write to drive RAM only | ❌ Lost on reboot     |
| `WriteROM`       | Write to drive NV flash | ✅ Persistent         |
| `WriteRAMAndROM` | Write both atomically   | ✅ Persistent         |

> ⚠️ **Always use `WriteRAMAndROM`** (or call `WriteRAMAndROM` directly) for any configuration that must survive a power cycle. Pure RAM writes are wiped on every reboot.

### 4.4 UPID Parameters

Every drive parameter has a **UPID** (Universal Parameter ID), a 16-bit address used across all LinMot firmware. The client exposes named helpers for the most common UPIDs, plus raw access for anything else:

```go
// Named helpers
c.SetVelocity(ctx, 0.5, protocol_common.WriteRAMAndROM)
c.SetPosition1(ctx, 10.0, protocol_common.WriteRAM)

// Raw UPID access
c.ReadRAM(ctx, 0x1450)                  // read any UPID
c.WriteRAMAndROM(ctx, 0x1450, 100000)   // write any UPID to RAM+ROM

// Discovery
upids, _ := c.GetAllParameterIDs(ctx)      // all UPIDs on this drive
modified, _ := c.GetModifiedParameterIDs(ctx)  // only non-default values
```

UPID list reference: [`reference/LinUDP_V2.md`](reference/LinUDP_V2.md) §6

### 4.5 Command Tables

Command tables are sequences of up to 255 motion entries stored in drive flash. They are defined in YAML and support **template variables** (`${VAR_NAME}`) for runtime parameterization (e.g., configurable positions/velocities without re-defining the whole table).

**YAML format:**

```yaml
version: "1"
drive_model: "C1250-MI"
entries:
  - id: 1
    name: "move to pick"
    type: "VAI_GoToPos"
    par1: ${POSITION_DOWN} # Position in 0.1µm units
    par2: ${MAX_VELOCITY} # Velocity in µm/s
    par3: ${ACCELERATION} # Accel in 1e-5 m/s²
    par4: ${DECELERATION}
    sequenced_entry: 2 # Next entry to run (chaining)

  - id: 2
    name: "wait stop"
    type: "WaitDemandVelLT"
    par1: 0
```

**Supported entry types** include (non-exhaustive): `VAI_GoToPos`, `Delay`, `SetDO`, `ClearDO`, `WaitDemandVelLT`, `WaitFalling`, `WaitRising`, `IfDemandPosLT`, `IfDemandPosGT`, `IfMaskedX4Eq`, `IfMaskedStatusEq`, and many more — see the full list in [`reference/LinMot_MotionCtrl.md`](reference/LinMot_MotionCtrl.md).

**Loading and deploying a command table:**

```go
ct, err := client_command_tables.Load("path/to/table.yaml")
// Bind template variables
ct.SetVar("POSITION_DOWN", 500000)  // 50 mm
ct.SetVar("MAX_VELOCITY", 100000)   // 0.1 m/s
ct.SetVar("ACCELERATION", 50000)
ct.SetVar("DECELERATION", 50000)
ct.SetVar("DELAY_AT_BOTTOM", 1000)

// Validate and deploy (stops MC, writes all entries, flash-saves, restarts MC)
err = c.SetCommandTable(ctx, ct)
```

Example tables: [`client/rtc/command_tables/testdata/`](client/rtc/command_tables/testdata/)

---

## 5. Client API Reference

### 5.1 Creating a Client

| Function                                                            | Description                                                    |
| ------------------------------------------------------------------- | -------------------------------------------------------------- |
| `linmot.NewClient(ip)`                                              | Convenience wrapper — default ports, validates connectivity    |
| `client.NewUDPClient(ip, drivePort, masterPort, bindAddr, timeout)` | Full control over ports and bind address                       |
| `client.NewClientWithTransport(ip, transport)`                      | Inject a pre-built transport (e.g., from `SharedUDPTransport`) |
| `client.NewMockClient()`                                            | Testing only — returns `(*Client, transport.Server)`           |

On creation, the client **validates connectivity** by sending status requests. If the drive is not reachable or not in an operational state, construction fails with a descriptive error. Set `LINMOT_SKIP_VALIDATE_CONNECTIVITY=1` to bypass this check.

Always `defer c.Close()` — this stops the internal TX/RX goroutines.

### 5.2 State Machine Control

```go
c.EnableDrive(ctx)       // → Operation Enabled state
c.DisableDrive(ctx)      // → Switch On Disabled state
c.AcknowledgeError(ctx)  // rising + falling edge on error-ack bit
c.QuickStop(ctx)         // emergency stop (clears quick-stop bit)
c.Home(ctx)              // initiates homing sequence (up to 30s)
c.SendControlWord(ctx, word)  // raw control word (advanced use)
c.GetDriveStatus(ctx)    // read-only status poll
```

**State machine states** (decoded from `StateVar >> 8`):

| Code | Name               |
| ---- | ------------------ |
| 0    | NotReadyToSwitchOn |
| 1    | SwitchOnDisabled   |
| 2    | ReadyToSwitchOn    |
| 3    | SetupError         |
| 4    | GeneralError       |
| 5    | HWTests            |
| 6    | ReadyToOperate     |
| 8    | OperationEnabled   |
| 9    | Homing             |

Source: [`client/operational.go`](client/operational.go)

### 5.3 Motion Control

**VAI commands** (Velocity/Acceleration/Incremental):

```go
// All positions in mm, velocities in m/s, accel/decel in m/s²
c.VAIGoToPosition(ctx, positionMM, velocityMS, accelMS2, decelMS2)
c.VAIIncrementDemandPosition(ctx, incrementMM, velocityMS, accelMS2, decelMS2)
c.VAIIncrementTargetPosition(ctx, incrementMM, velocityMS, accelMS2, decelMS2)
c.VAIStop(ctx)
```

**Predefined / command table execution:**

```go
c.ExecutePredefinedCommand(ctx, commandID)   // trigger a command table entry by ID
```

**Streaming:**

```go
c.StreamingGoToPosition(ctx, positionMM, velocityMS)  // streaming motion command
```

Full interface list: see [§7 Motion Control Interfaces](#7-motion-control-interfaces).

### 5.4 RTC Parameters

```go
// Named helpers (most common motion parameters)
c.GetPosition(ctx)                                        // actual position in mm
c.SetPosition1(ctx, 10.0, protocol_common.WriteRAM)      // Easy Steps Position 1
c.SetPosition2(ctx, 50.0, protocol_common.WriteRAMAndROM)// Easy Steps Position 2
c.GetVelocity(ctx) / c.SetVelocity(ctx, 0.5, ...)
c.GetAcceleration(ctx) / c.SetAcceleration(ctx, 5.0, ...)
c.GetDeceleration(ctx) / c.SetDeceleration(ctx, 5.0, ...)

// Easy Steps configuration
c.SetEasyStepsAutoStart(ctx, protocol_common.EasyStepsAutoStart.Enabled, ...)
c.SetEasyStepsAutoHome(ctx, protocol_common.EasyStepsAutoHome.Enabled, ...)
c.SetEasyStepsInputRisingEdgeFunction(ctx, inputPin, value, ...)
c.SetEasyStepsInputCurveCmdID(ctx, inputPin, cmdID, ...)

// I/O configuration
c.SetIODefOutputFunction(ctx, outputPin, protocol_common.OutputConfig.InterfaceOutput, ...)
c.SetIODefInputFunction(ctx, inputPin, protocol_common.InputFunction.Trigger, ...)
c.SetTriggerMode(ctx, protocol_common.TriggerModeConfig.Direct, ...)
c.SetRunMode(ctx, protocol_common.RunMode.TriggeredCommandTable, ...)

// Raw UPID access
c.ReadRAM(ctx, upid)
c.WriteRAMAndROM(ctx, upid, value)
c.GetParameterMinValue(ctx, upid)
c.GetParameterMaxValue(ctx, upid)
c.GetParameterDefaultValue(ctx, upid)

// Parameter discovery
c.GetAllParameterIDs(ctx)       // []uint16
c.GetModifiedParameterIDs(ctx)  // []uint16
c.GetAllParameters(ctx)         // []ParameterInfo (with RAM/ROM usage)
c.GetModifiedParameters(ctx)    // []ParameterInfo
```

### 5.5 Command Tables

```go
// Deploy a command table (stop MC → write entries → flash save → restart MC)
c.SetCommandTable(ctx, ct)

// Deploy with options
c.SetCommandTableWithOptions(ctx, ct, client_command_tables.SetCommandTableOptions{
    RestartMC:    true,   // default true
    SkipFlashSave: false, // default false
})

// Read back current command table from drive
ct, err := c.GetCommandTable(ctx)

// Presence masks (256-entry bitmap)
masks, err := c.GetPresenceMasks(ctx) // [8]uint32

// MC lifecycle (needed before/after manual entry writes)
c.StopMotionController(ctx)
c.StartMotionController(ctx)
c.SaveCommandTableToFlash(ctx) // ⚠️ see Known Issues §10
```

> **`ErrCommandTableUnchanged`** is returned by `SetCommandTableWithOptions` when the on-drive entries already match the desired table byte-for-byte. No writes are performed. Callers should skip any post-write recovery/homing.

### 5.6 Monitoring Channels

The drive supports 4 live monitoring channels. Each channel can be assigned any UPID, and its value is returned inline with every response when the monitoring flag is set.

```go
// Assign UPIDs to channels
c.ConfigureMonitoringChannel(ctx, 1, 0x0020) // channel 1 = actual position
c.ConfigureMonitoringChannel(ctx, 2, 0x0028) // channel 2 = demand position
c.ConfigureMonitoringChannels(ctx, [4]uint16{0x0020, 0x0028, 0x0040, 0x0050})

// Read live data
status, err := c.GetMonitoringData(ctx)
// status.MonitoringChannel[0..3] = current values

// Inspect current channel assignments
c.GetMonitoringChannelConfiguration(ctx, 1)
c.GetAllMonitoringChannelConfigurations(ctx)  // [4]uint16
```

Monitoring channel UPIDs: `0x20A8`–`0x20AB`. Reference: [`reference/LinUDP_V2.md`](reference/LinUDP_V2.md) §5.

### 5.7 Curves

Curves are pre-stored motion profiles (e.g., cam profiles) saved in drive flash.

```go
c.UploadCurve(ctx, curveID, infoBlock, dataBlock)
c.DownloadCurve(ctx, curveID)   // returns (infoBlock, dataBlock, err)
c.ModifyCurve(ctx, curveID, infoBlock, dataBlock)
c.SaveAllCurves(ctx)            // RAM → flash
c.DeleteAllCurves(ctx)          // ⚠️ irreversible from RAM
```

### 5.8 Error Log

```go
c.GetErrorLog(ctx)                  // []ErrorLogEntry
c.GetErrorLogWithText(ctx)          // same + human-readable descriptions
c.GetErrorLogCounts(ctx)            // (logged, occurred uint16)
c.GetErrorLogEntry(ctx, entryNum)   // single entry
c.GetErrorText(ctx, errorCode)      // "Position Lag Error", etc.
c.CheckDriveFault(ctx)              // nil if healthy; DriveFaultError if fault active
```

### 5.9 Drive Operations

```go
c.RestartDrive(ctx)                        // ⚠️ full reboot
c.ResetOSParametersToDefault(ctx)          // ⚠️ wipes OS-level config
c.ResetMCParametersToDefault(ctx)          // ⚠️ wipes motion control config
c.ResetInterfaceParametersToDefault(ctx)   // ⚠️ wipes I/O, comms config
c.ResetApplicationParametersToDefault(ctx) // ⚠️ wipes application config
```

---

## 6. Multi-Drive: ClientPool

`ClientPool` manages a thread-safe cache of `Client` instances keyed by drive IP. All clients in a pool share one `SharedUDPTransport` (one UDP socket on port 41136).

```go
pool := client.NewClientPool()
defer pool.Close()

c0, err := pool.GetClient("10.8.7.232")  // creates on first call, cached thereafter
c1, err := pool.GetClient("10.8.7.234")

// After flash save (drive UDP stack reset) — evict and reconnect
pool.EvictClient("10.8.7.232")
c0, err = pool.GetClient("10.8.7.232")  // fresh client, new socket state
```

**Key behaviours:**

- First `GetClient` for a new IP blocks until the client is created and connectivity is validated
- Concurrent callers for the same IP wait on a shared future — only one creation attempt per IP at a time
- Failed entries are retried after 5 seconds (`failedEntryRetryInterval`)
- `EvictClient` closes the old client and removes it; next `GetClient` creates fresh (needed after `SaveCommandTableToFlash` — see [§10](#10-known-issues--quirks))

Source: [`client/pool.go`](client/pool.go)

---

## 7. Motion Control Interfaces

The MC channel supports 15 interface types. The `MotionControlManager` (and `Client`) expose the ones used in production. The full protocol-level implementations live in `protocol/motion_control/`:

| Interface         | Package              | Description                                   |
| ----------------- | -------------------- | --------------------------------------------- |
| VAI               | `vai/`               | Go-to position, increment demand/target, stop |
| VAI16             | `vai16/`             | 16-bit variant of VAI                         |
| Predefined        | `predefined/`        | Execute a stored command table entry          |
| PredefVAI16       | `predef_vai16/`      | Predefined with 16-bit params                 |
| Streaming         | `streaming/`         | Streaming position/velocity commands          |
| Interface Control | `interface_control/` | Direct MC interface control commands          |
| Advanced          | `advanced/`          | Advanced motion commands                      |
| Bestehorn         | `bestehorn/`         | Bestehorn cam profile execution               |
| Cams              | `cams/`              | Cam table execution                           |
| Encoder Cams      | `encoder_cams/`      | Encoder-synchronized cam execution            |
| Indexing          | `indexing/`          | Incremental indexing moves                    |
| VAI Dec/Acc       | `vai_dec_acc/`       | VAI with separate dec/acc profiles            |
| VAI Positioning   | `vai_positioning/`   | VAI positioning variant                       |
| VAI Predef Acc    | `vai_predef_acc/`    | VAI with predefined acceleration              |

---

## 8. Testing Infrastructure

### 8.1 MockLinMot — Full Drive Simulator

[`test/mock_linmot.go`](test/mock_linmot.go) implements a complete LinMot drive simulator in memory. It is used by both unit tests and E2E tests across the codebase.

**Simulated capabilities:**

- Full CiA 402 state machine (NotReadyToSwitchOn → SwitchOnDisabled → ReadyToOperate → OperationEnabled → Homing)
- RAM and ROM parameter storage (`map[uint16]int32` each)
- RTC command table: read/write/flash-save/presence-masks
- Motion simulation: VAI GoToPosition with demand position tracking
- Error and warning injection (`simulateError`, `persistentError`, `simulateWarning`)
- Error acknowledgment edge detection
- 4-channel monitoring (reads from `monitoredVariables`)
- MC lifecycle: stop/start, sequence counter deduplication
- Configurable delays for status and error text responses

```go
c, server := client.NewMockClient()
drive := test.NewMockLinMot(server)
go drive.Run()
defer drive.Close()

// Now use c like a real client — it talks to the in-memory drive
```

### 8.2 Running Tests

```bash
# From module root
go test ./linmot_client/...

# Specific package
go test -v ./linmot_client/client/...

# E2E tests (require a real drive, skipped by default)
LINMOT_E2E_HOST=10.8.7.232 go test -v ./linmot_client/client/ -run TestE2E

# Enable test debug logging
LINMOT_TEST_DEBUG=1 go test ./linmot_client/...
```

The module has **42 test files** and **148 Go source files**. Tests cover:

- `client/common/` — RequestManager parity tests, RTC counter, pending request lifecycle
- `client/` — operational state assertions, pool eviction, fault detection, diagnostic motion matrix
- `protocol/*/` — packet encoding/decoding roundtrips for every interface
- `transport/` — UDP client, shared UDP routing, packet loss, context deadline handling

---

## 9. Reference Documentation

All official LinMot documentation lives in [`reference/`](reference/). See [`reference/README.md`](reference/README.md) for the full index, key reference table, and our drive configuration values.

| Document                                                                       | Most Important Section                                                       |
| ------------------------------------------------------------------------------ | ---------------------------------------------------------------------------- |
| [`C1250_MI_Installation_Guide.md`](reference/C1250_MI_Installation_Guide.md)   | §9.9 DIP switch layout · §10 LED blink codes                                 |
| [`LinUDP_V2.md`](reference/LinUDP_V2.md)                                       | §3.2 default IP addressing · §4 packet format · §5 monitoring · §6 UPID list |
| [`LinUDP_V2_DLL.md`](reference/LinUDP_V2_DLL.md)                               | Appendix I — "Static by IP Configuration" commissioning path                 |
| [`LinMotTalk_Manual.md`](reference/LinMotTalk_Manual.md)                       | §4.1 — factory reset (6-step hardware procedure)                             |
| [`LinMot_MotionCtrl.md`](reference/LinMot_MotionCtrl.md)                       | Easy Steps auto start/home · command table types · Set Stiff vs. Set Soft    |
| [`decompiled_linudp_csharp_lib.cs`](reference/decompiled_linudp_csharp_lib.cs) | Canonical reference for packet construction and ACI function signatures      |

---

## 10. Known Issues & Quirks

### Flash Save Kills the LinUDP Stack

`SaveCommandTableToFlash` sends RTC command `0x80` to persist the command table. The drive internally calls `LMcf_StopDefault(0x35)` as part of the flash-save routine, which **tears down the LinUDP networking stack**. As a result:

- The drive never sends a standard RTC ACK for this command
- The Go client's RTC request times out waiting for a response
- After the flash save completes, the drive rebuilds its network stack — but the existing UDP socket state on both sides is stale

**Workaround:** After calling `SaveCommandTableToFlash`, always call `pool.EvictClient(driveIP)` to close the old client and force a fresh connection on the next `GetClient` call. `SetCommandTableWithOptions` handles this internally.

This behaviour is the subject of the open Karsten consultation ([`_AGENTS/karsten/LinMot_Consultation_Karsten.zip`](../../_AGENTS/karsten/LinMot_Consultation_Karsten.zip)).

### Multi-Master Lock

LinMot drives operate in single-master mode. If LinMot-Talk (the Windows configuration tool) is connected to a drive, it holds the master lock and the Go client will receive all-zero status telegrams (`ErrInvalidStatusTelegram`) or no response at all. Disconnect LinMot-Talk before running the Go client.

### RTCCommand Flag Bug (Fixed)

The `RTCCommand` request flag was historically set to `0x00000100` (response bit 8) instead of the correct `0x00000004` (request bit 2). This was a silent bug that caused drives to ignore RTC requests. The correct value is documented in [`protocol/common/constants.go`](protocol/common/constants.go) with a `FIXED` comment.

### Connectivity Validation Fallback

Some drives (or drives fresh off a factory reset) ignore status-only requests but respond to C#-style cyclic telegrams. `NewUDPClient` automatically falls back to a cyclic telegram probe if the normal status strategy exhausts its attempts. This behaviour mirrors the official C# SDK.
