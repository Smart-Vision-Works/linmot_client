# stage_primer

REST/gRPC service for controlling LinMot Z-axis actuators on the staged robot system. Uses the [`linmot_client`](../linmot_client/) library for all drive communication.

## Module Info

- **Module path:** `primer`
- **Go version:** 1.24
- **Depends on:** `github.com/Smart-Vision-Works/staged_robot` (via `replace => ../linmot_client`)

## What This Service Does

The stage_primer runs on the StagePrimer device (a BalenaOS Linux box). It:

1. **Configures LinMot drives** on startup (14 ROM parameters for triggered command table mode)
2. **Deploys command tables** when the operator enters RunMode (YAML template with Z-distance, speed, acceleration, pick time)
3. **Jogs drives** to specific Z positions for calibration
4. **Monitors faults** in the background and attempts auto-recovery
5. **Controls vacuum/purge** solenoids via LinMot digital outputs

All LinMot communication goes through the `linmot_client` library over LinUDP V2 (UDP port 49360).

## Directory Layout

```
stage_primer/
├── main.go                 # Entry point: startup, server, or linmot CLI
├── go.mod
│
├── linmot/                 # LinMot drive management (calls into linmot_client)
│   ├── client_factory.go   # Client pool setup, factory interface
│   ├── setup.go            # Drive commissioning (14 ROM parameters)
│   ├── command_table.go    # Command table deploy pipeline
│   ├── jog.go              # Z-axis positioning
│   ├── faults.go           # Fault monitoring + auto-recovery
│   ├── vacuum.go           # Vacuum/purge solenoid control
│   ├── commands.go         # CLI subcommands (setup, setpos, getpos)
│   ├── recovery_state.go   # Per-drive recovery state tracking
│   └── *.yaml              # Command table templates
│
├── server/                 # REST API + gRPC server
│   ├── server.go           # HTTP + gRPC server lifecycle
│   ├── grpc_server.go      # gRPC handlers (DeployCommandTable, Jog, etc.)
│   └── api_config.go       # REST config endpoints
│
├── client/                 # gRPC client (used by TensorPro to call stage_primer)
│   └── client.go
│
├── proto/                  # Protobuf definitions + generated code
│   └── stage_primer.proto
│
├── config/                 # Configuration (separate Go module: stage_primer_config)
│   └── config.go           # ClearCoreConfig, LinMotConfig structs
│
└── stubs/gsail-go/         # Minimal stub for gsail-go/logger (standalone build only)
```

---

## Where stage_primer Calls linmot_client

Every call from stage_primer into the linmot_client library is listed below, grouped by file. All linmot_client types are imported from `github.com/Smart-Vision-Works/staged_robot/...`.

### linmot/client_factory.go — Client Pool Management

Creates and manages the global client pool. One UDP socket (port 41136) shared across all drives.

| Line | Call                                     | Purpose                                 |
| ---- | ---------------------------------------- | --------------------------------------- |
| 75   | `client.NewClientPool()`                 | Create global client pool               |
| 83   | `globalClientPool.GetClient(linmotIP)`   | Get/create pooled client for a drive IP |
| 88   | `globalClientPool.Close()`               | Close all pooled clients                |
| 131  | `globalClientPool.EvictClient(linmotIP)` | Evict stale client after recovery       |

### linmot/setup.go — Drive Commissioning (14 ROM Parameters)

Reads current drive parameters and writes only those that differ. Handles the 30–90 second ROM write recovery.

| Line | Call                                            | Purpose                                     |
| ---- | ----------------------------------------------- | ------------------------------------------- |
| 85   | `linmotClient.ReadRAM(ctx, upid)`               | Read current parameter value                |
| 96   | `linmotClient.WriteRAMAndROM(ctx, upid, value)` | Write parameter to RAM + ROM                |
| 150  | `client.GetStatus(pollCtx)`                     | Poll drive status during ROM write recovery |

### linmot/command_table.go — Command Table Deploy Pipeline

Loads YAML template, binds variables (Z-distance, speed, acceleration, pick time), deploys to drive, optionally saves to flash.

| Line    | Call                                                                | Purpose                                            |
| ------- | ------------------------------------------------------------------- | -------------------------------------------------- |
| 174     | `linmot_command_tables.Load(templatePath)`                          | Load YAML command table template                   |
| 207     | `linmot_command_tables.Load(templatePath)`                          | Load inspect variant template                      |
| 418–442 | `templateCopy.SetVar(...)`                                          | Bind Z-distance, velocity, acceleration, pick time |
| 442     | `templateCopy.Validate()`                                           | Validate command table after variable binding      |
| 474     | `linmotClient.SetCommandTableWithOptions(ctx, &templateCopy, opts)` | Deploy command table to drive RAM (with verify)    |
| 287     | `flashClient.SaveCommandTableToFlash(flashCtx)`                     | Save command table to flash (fire-and-forget)      |
| 310     | `freshClient.StartMotionController(restartCtx)`                     | Restart motion controller after flash save         |
| 320     | `freshClient.Home(homeCtx)`                                         | Re-home drive after MC restart                     |
| 336     | `freshClient.GetStatus(pollCtx)`                                    | Poll status during homing                          |
| 580     | `linmotClient.SetCommandTable(deployCtx, &templateCopy)`            | Deploy inspect command table (no flash)            |

### linmot/jog.go — Z-Axis Positioning

Moves the drive to an absolute Z position. Handles state transitions (fault acknowledge, enable, mode switch).

| Line  | Call                                                        | Purpose                        |
| ----- | ----------------------------------------------------------- | ------------------------------ |
| 60    | `linmotClient.GetStatus(ctx)`                               | Check drive state before jog   |
| 71    | `linmotClient.AcknowledgeError(ctx)`                        | Clear fault if in error state  |
| 76–91 | `linmotClient.EnableDrive(ctx)`                             | Enable drive (multiple paths)  |
| 105   | `linmotClient.SetPosition1(ctx, pos, RAM)`                  | Set target position 1          |
| 108   | `linmotClient.SetPosition2(ctx, pos, RAM)`                  | Set target position 2 (safety) |
| 114   | `linmotClient.SetRunMode(ctx, MotionCommandInterface, ...)` | Exit VAI2PosContinuous         |
| 122   | `linmotClient.SetRunMode(ctx, VAI2PosContinuous, ...)`      | Re-enter to trigger motion     |
| 145   | `linmotClient.GetPosition(ctx)`                             | Read actual position           |

### linmot/faults.go — Fault Monitoring + Auto-Recovery

Background goroutine polls drives every 1 second. Auto-acknowledges recoverable faults (up to 3 attempts per 60-second window).

| Line | Call                                    | Purpose                          |
| ---- | --------------------------------------- | -------------------------------- |
| 491  | `linmotClient.AcknowledgeError(ackCtx)` | Auto-acknowledge during recovery |
| 509  | `linmotClient.CheckDriveFault(pollCtx)` | Poll drive for fault status      |

### linmot/vacuum.go — Vacuum/Purge Solenoid Control

Controls X4.3 (purge) and X4.4 (vacuum) digital outputs by toggling pin functions.

| Line  | Call                                                           | Purpose                                   |
| ----- | -------------------------------------------------------------- | ----------------------------------------- |
| 62–92 | `linmotClient.SetIODefOutputFunction(ctx, pin, function, RAM)` | Set output pin function (AlwaysOn / None) |

Pin mapping: X4.3 = Purge, X4.4 = Vacuum.

### server/grpc_server.go — gRPC Entry Point

The gRPC server is the entry point for TensorPro requests. It creates `linmot_client.Client` instances via the factory.

| Line                 | Call                             | Purpose                                                   |
| -------------------- | -------------------------------- | --------------------------------------------------------- |
| (via linmot package) | `linmot.DeployCommandTable(...)` | Called from `DeployCommandTable` gRPC handler             |
| (via linmot package) | `linmot.Jog(...)`                | Called from `Jog` gRPC handler                            |
| (via linmot package) | `linmot.Setup(...)`              | Called from `SetConfig` handler when new LinMots detected |

### linmot/mock_factory.go — Test Infrastructure

Creates mock clients backed by `test.MockLinMot` (full drive simulator from linmot_client).

| Line | Call                                  | Purpose                     |
| ---- | ------------------------------------- | --------------------------- |
| 46   | `client.NewMockClient()`              | Create mock transport pair  |
| 48   | `test.NewMockLinMot(transportServer)` | Create mock drive simulator |

---

## Unit Conversion Constants

stage_primer converts from user-friendly units to LinMot hardware units in `command_table.go`:

| Parameter    | User Unit    | LinMot Unit  | Conversion |
| ------------ | ------------ | ------------ | ---------- |
| Z Distance   | mm           | 0.1 um       | x 10,000   |
| Speed        | % (0-100)    | um/s         | x 10,000   |
| Acceleration | % (0-200)    | 1e-5 m/s^2   | x 10,000   |
| Pick Time    | seconds      | 100 us units | x 10,000   |
| Purge Delay  | fixed 500 ms | 100 us units | x 10       |

## Command Table Template

The 11-entry pick cycle command table (`linmot/linmot_command_table.yaml`):

```
Entry  1: SetDO          Vacuum ON, Purge OFF
Entry  2: VAI_GoToPos    Move down to ${POSITION_DOWN}
Entry  3: WaitMotionDone Wait until motion finished
Entry  4: Delay          Wait ${DELAY_AT_BOTTOM} (pick dwell)
Entry  5: VAI_GoToPos    Move up to 0 (home)
Entry  6: WaitMotionDone Wait until motion finished
Entry  7: IfMaskedX4Eq   Check if trigger still HIGH
Entry  8: WaitFalling    Wait for trigger falling edge
Entry  9: ClearDO        Vacuum OFF, Purge ON
Entry 10: Delay          Wait 500ms (purge)
Entry 11: ClearDO        All outputs OFF
```

## 14 ROM Parameters (setup.go)

Written during drive commissioning. Compare-before-write minimizes flash wear.

| Parameter               | Value                   | Purpose                                  |
| ----------------------- | ----------------------- | ---------------------------------------- |
| Run Mode                | TriggeredCommandTable   | Execute command table on trigger         |
| Input 4.6 Function      | Trigger                 | ClearCore DIO pin = trigger input        |
| Trigger Mode            | Direct                  | Immediate trigger evaluation             |
| Output 4.3, 4.4         | Interface Output        | Vacuum/purge solenoid control            |
| Easy Steps Auto Start   | Enabled                 | Auto-start on power-up                   |
| Easy Steps Auto Home    | Enabled                 | Auto-home on power-up                    |
| Rising Edge 4.6         | EvalCommandTableCommand | Trigger HIGH = start command table       |
| Rising Edge CurveCmdID  | 1                       | Start at entry 1 (move down)             |
| Falling Edge CurveCmdID | 5                       | Trigger LOW = start at entry 5 (move up) |
| Go To Position          | 10mm                    | Default safe position                    |
| Initial Position        | 10mm                    | Default safe position                    |
| Smart Control Word      | 1                       | Operation enabled flag                   |
