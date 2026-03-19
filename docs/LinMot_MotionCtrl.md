# LinMot Motion Control Software

**User Manual for Drive Series**
- A1100 / C1100 / C1200 / E1200 / E1400

**Document Version:** 4.3.4  
**Date:** September 2014

**© 2014 NTI AG**

This work is protected by copyright. Under the copyright laws, this publication may not be reproduced or transmitted in any form, electronic or mechanical, including photocopying, recording, microfilm, storing in an information retrieval system, not even for didactic use, or translating, in whole or in part, without the prior written consent of NTI AG.

**LinMot®** is a registered trademark of NTI AG.

**Note:** The information in this documentation reflects the stage of development at the time of press and is therefore without obligation. NTI AG reserves itself the right to make changes at any time and without notice to reflect further technical advance or product improvement.

---

## Table of Contents

1. [System Overview](#1-system-overview)
   - 1.1 [References](#11-references)
   - 1.2 [Definitions, Items, Shortcuts](#12-definitions-items-shortcuts)
   - 1.3 [Data Types](#13-data-types)
2. [Motion Control Interfaces](#2-motion-control-interfaces)
3. [State Machine](#3-state-machine)
   - 3.1 [State 0: Not Ready To Switch On](#31-state-0-not-ready-to-switch-on)
   - 3.2 [State 1: Switch On Disabled](#32-state-1-switch-on-disabled)
   - 3.3 [State 2: Ready To Switch On](#33-state-2-ready-to-switch-on)
   - 3.4 [State 3: Setup Error State](#34-state-3-setup-error-state)
   - 3.5 [State 4: Error State](#35-state-4-error-state)
   - 3.6 [State 5: HW Test](#36-state-5-hw-test)
   - 3.7 [State 6: Ready To Operate](#37-state-6-ready-to-operate)
   - 3.8 [State 8: Operation Enabled](#38-state-8-operation-enabled)
   - 3.9 [State 9: Homing](#39-state-9-homing)
   - 3.10 [State 10: Clearance Check](#310-state-10-clearance-check)
   - 3.11 [State 11: Going To Initial Position](#311-state-11-going-to-initial-position)
   - 3.12 [State 12: Aborting](#312-state-12-aborting)
   - 3.13 [State 13: Freezing](#313-state-13-freezing)
   - 3.14 [State 14: Error Behaviour Quick Stop](#314-state-14-error-behaviour-quick-stop)
   - 3.15 [State 15: Going To Position](#315-state-15-going-to-position)
   - 3.16 [State 16: Jogging +](#316-state-16-jogging-)
   - 3.17 [State 17: Jogging -](#317-state-17-jogging--)
   - 3.18 [State 18: Linearizing](#318-state-18-linearizing)
   - 3.19 [State 19: Phase Searching](#319-state-19-phase-searching)
   - 3.20 [State 20: Special Mode](#320-state-20-special-mode)
   - 3.21 [Building the Control Word](#321-building-the-control-word)
   - 3.22 [Control Word](#322-control-word)
   - 3.23 [Status Word](#323-status-word)
   - 3.24 [Warn Word](#324-warn-word)
4. [Motion Command Interface](#4-motion-command-interface)
   - 4.1 [Motion Command Interface](#41-motion-command-interface)
   - 4.2 [Overview Motion Commands](#42-overview-motion-commands)
   - 4.3 [Detailed Motion Command Description](#43-detailed-motion-command-description)
5. [Setpoint Generation](#5-setpoint-generation)
   - 5.1 [VA-Interpolator](#51-va-interpolator)
   - 5.2 [Sine VA Motion](#52-sine-va-motion)
   - 5.3 [Bestehorn VAJ Motion](#53-bestehorn-vaj-motion)
   - 5.4 [P(V)-Stream](#54-pv-stream)
   - 5.5 [CAM Motions](#55-cam-motions)
6. [Command Table](#6-command-table)
7. [Drive Configuration](#7-drive-configuration)
8. [Motor Configuration](#8-motor-configuration)
9. [State Machine Setup](#9-state-machine-setup)
10. [Error Code List](#10-error-code-list)
11. [Contact Addresses](#11-contact-addresses)

---

## 1. System Overview

This user manual describes the Motion Control SW functionality of the LinMot E1200 / E1400 drives.

### 1.1 References

| Ref | Title | Source |
|-----|-------|--------|
| 1 | Installation_Guide_E1200.pdf | www.linmot.com |
| 2 | Installation_Guide_E1400.pdf | www.linmot.com |
| 3 | Usermanual_LinMot-Talk_4.pdf | www.linmot.com |

The documentation is distributed with the LinMot-Talk configuration software or can be downloaded from the Internet from the download section of our homepage.

### 1.2 Definitions, Items, Shortcuts

| Shortcut | Meaning |
|----------|---------|
| LM | LinMot linear motor |
| OS | Operating system (Software) |
| MC (SW) | Motion Control (Software) |
| Intf | Interface (Software) |
| Appl | Application (Software) |
| VAI | VA-Interpolator (Max velocity limited acceleration position interpolator) |
| Pos | Position |
| Vel | Velocity |
| Acc | Acceleration |
| Dec | Deceleration |
| UPID | Unique Parameter ID (16 bit) |

### 1.3 Data Types

| Type | Range/Format | Num of Bytes |
|------|--------------|--------------|
| Bool | Boolean, False/True | 1/8 |
| Byte | 0..255 | 1 |
| Char | ASCII | 1 |
| String | Array of char last char = 00h | X |
| SInt16 | -32768..32767 | 2 |
| UInt16 | 0..65535 | 2 |
| SInt32 | -2147483648..2147483647 | 4 |
| UInt32 | 0..4294967295 | 4 |
| Float | - | 4 |

---

## 2. Motion Control Interfaces

For controlling the behavior of the motion control SW, two different interfaces are available. For controlling the main state machine, a bit coded control word can be used. For controlling the motion functionality a memory mapped motion command interface can be used. These two instances are mapped via an interface SW to an upper control system (PLC, IPC, PC, ..). The interfacing is done with digital I/Os or a serial link like Profibus DP, CAN bus (CANopen), RS485, RS422 or RS232 (LinRS protocol), Ethernet (POWERLINK, EtherCAT, Ethernet/IP).

With LinMot-Talk the control over the control word can be taken bit by bit, for testing and debugging. Unused control word bits can be forced by parameter value. Also the control of the motion command interpreter can be switched to the control panel of the LinMot-Talk software for testing.

⚠️ **Warning:** All this can be done while the system is running, so be careful using this features on a running machine!

---

## 3. State Machine

The main behavior of the axles is controlled with the control word, as shown in the following state diagram.

### State Diagram

```
Not Ready to Switch On (0)
         ↓
         Bit 7
         ↓
         Bit 0=0
Switch On Disabled (1)
         ↓
Control Word xxxx xxxx xxxx x110
         ↓
Error (4) ← Setup Error (3)
         ↓
Ready to Switch On (2)
         ↓
         Bit 0=1
         ↓
HW Tests (5)
         ↓
Ready to Operate (6)
         ↓
         Bit 3=1
         ↓
Operation Enabled (8)
         ↓
         Power On
         ↓
Homing (9) → Clearance Checking (10) → Going To Initial Position (11)
         ↓
Aborting (12) → Freezing (13) → Error Behaviour Quick Stop (14)
         ↓
Going To Position (15) → Jogging + (16) → Jogging - (17)
         ↓
Linearizing (18) → Phase Searching (19) → Special Mode (20)
```

The state machine can be followed in the PLCs with fieldbus using the StateVar. This response word can be configured for any supported fieldbus.

### State Var Structure

The State Var is divided into two sections: the **Main State** section (high byte) contains directly the number of the state machine, the content of the **Sub State** (low byte) is state depending.

| Bits | Field | Description |
|------|-------|-------------|
| 15-8 | Main State | State machine number (0-20) |
| 7-0 | Sub State | State-dependent information |

### State Var Main State Values

| Main State | Description |
|------------|-------------|
| 00 | Not Ready To Switch On |
| 01 | Switch On Disabled |
| 02 | Ready To Switch On |
| 03 | Setup Error |
| 04 | Error |
| 05 | HW Tests |
| 06 | Ready To Operate |
| 07 | - |
| 08 | Operation Enabled |
| 09 | Homing |
| 10 | Clearance Check |
| 11 | Going To Initial Position |
| 12 | Aborting |
| 13 | Freezing |
| 14 | Quick Stop (Error Behaviour) |
| 15 | Going To Position |
| 16 | Jogging + |
| 17 | Jogging - |
| 18 | Linearizing |
| 19 | Phase Search |
| 20 | Special Mode |

### State Var Sub State Values

| Main State | Sub State | Description |
|------------|-----------|-------------|
| 0-2 | 0 | - |
| 3 | 0 | Error Code which will be logged |
| 4 | 0 | Logged Error Code |
| 5-6 | 0 | (Not yet defined) |
| 8 | Bits 0..3 | Motion Command Count |
| 8 | Bit 4 | Event Handler Active |
| 8 | Bit 5 | Motion Active |
| 8 | Bit 6 | In Target Position |
| 8 | Bit 7 | Homed |
| 9 | 0Fh | Homing Finished |
| 10 | 0Fh | Clearance Check Finished |
| 11 | 0Fh | Going To Initial Position Finished |
| 12-13 | - | (Not yet defined) |
| 15 | 0Fh | Going To Position Finished |
| 16 | 01h | Moving positive |
| 16 | 0Fh | Jogging + Finished |
| 17 | 01h | Moving negative |
| 17 | 0Fh | Jogging - Finished |
| 18-20 | - | (Not yet defined) |

### 3.1 State 0: Not Ready To Switch On

In this state the release of control word bit 0 switch on is awaited. As soon as this bit is cleared a change to state 1 is performed. This behavior avoids self starting if all necessary bits for a start are set correctly in the control word.

### 3.2 State 1: Switch On Disabled

The state machine rests in this state as long as the bits 1 or 2 of the control word are cleared.

### 3.3 State 2: Ready To Switch On

The state machine rests in this state as long as the bit 0 is cleared.

### 3.4 State 3: Setup Error State

The state machine rests in this state as long the bit 0 is cleared.

### 3.5 State 4: Error State

The error state can be acknowledged with a rising edge of the control word bit 7 'Error Acknowledge'. If the error is fatal, bit 12 'Fatal Error' in the status word is set, no error acknowledgment is possible.

In the case of a fatal error, the error has to be checked, and the problem has to be solved before a reset or power cycle is done for resetting the error.

### 3.6 State 5: HW Test

The HW Test state is an intermediate state before turning on the power stage of the drive. If everything seems to be ok the servo changes to state 6 without any user action. The test takes about 300ms.

### 3.7 State 6: Ready to Operate

In this state the motor is either position controlled or with demand current = 0 and under voltage, but no motion commands are accepted. The mode is configurable with UPID 6300h.

Sending motion commands in this state will generate the error 'Motion command sent in wrong state' and a state change to the error state will be performed.

Clearing the control word bit 3 'Enable Operation' in state 8 or higher will stop immediately the set point generation and a state transition to 6 is performed. Clearing the bit while a motion is in execution a following error might be generated.

### 3.8 State 8: Operation Enabled

This is the state of the normal operation in which the motion commands are executed. It is strongly recommended to use the State Var for the motion command synchronization with any fieldbus system.

**State Var Structure in State 8:**

| Bit | Field | Description |
|-----|-------|-------------|
| 15-8 | Main State | = 8 (Operation Enabled) |
| 7 | Homed | Status word bit 'Homed' |
| 6 | In Target Position | Status word bit 'In Target Position' |
| 5 | Motion Active | Status word bit 'Motion Active' |
| 4 | Event Handler Active | Status word bit 'Event Handler Active' |
| 3-0 | Motion Command Count | Actual interpreted 'Motion Command Count' |

In the high byte stands the number of the main state = 8. In the low byte stands in the lowest 4 bits the actual interpreted 'Motion Command Count', bit 4 indicates if the event handler is active, in bit 5 stands the status word bit 'Motion Active', in bit 6 the status word bit 'In Target Position' and in bit 7 the status word bit 'Homed'. Because the 'Motion Command Count' echo and this status word bits are located in the same byte no data consistency problem is possible with any fieldbus.

A new motion command can be setup when the Motion Command Count has changed to the last sent and the 'Motion Active' bit is 0 or the 'In Target Position' bit is 1 if an exact positioning is required.

### 3.9 State 9: Homing

The homing state is used to define the position of the system according a mechanical reference, a home switch or an index.

For LinMot motors the slider home position at this home position is taken to compensate edge effects.

In the home sequence a position check of two positions and the motion to an initial position can be added.

⚠️ **Hint:** If a mechanical stop homing mode is chosen, the initial position should be a little apart from this mechanical stop to avoid overheating of the motor.

### 3.10 State 10: Clearance Check

Setting the Clearance Check bit in the Control Word, two positions are moved to, to check if the whole motion range is free. Normally this action is added to the homing sequence to ensure that the homing was done correctly.

### 3.11 State 11: Going To Initial Position

Setting the Go To Initial Position bit in the control word, the servo moves to the initial position, normally used to move away from the mechanical stop after homing, to protect the motor from overheating at the mechanical stop. After an error it is also recommended to move to a defined position again.

### 3.12 State 12: Aborting

Clearing the /Abort bit in the control word initiates a quick stop. After the motion has stopped the servo rests position controlled. Setting the bit again the drive rests in position until a new motion command is executed.

### 3.13 State 13: Freezing

Clearing the /Freeze bit in the control word initiates a quick stop. After the motion is stopped the servo rests position controlled. Setting the bit again the drive will finish the frozen motion (e.g. if it was a VAI command). Curve motion can be frozen but not restarted by releasing this bit, setting the bit again the motor moves at the target position of the last VAI command, if never used a VAI command it will go to the initial position.

### 3.14 State 14: Error Behaviour Quick Stop

Most of the errors, which can occur during an active motion, cause a quick stop behavior to stop the motion. After the quick stop is finished the motor is no longer position controlled.

### 3.15 State 15: Going To Position

Setting the Go To Position bit in the control word, the servo moves to the defined position, recommendable for example, after an error, to move to a defined position again.

### 3.16 State 16: Jogging +

Setting the Jog Move + bit in the control word, the servo moves either a defined position increment or to the maximal position with a limited speed. Releasing the bit will stop the motion.

### 3.17 State 17: Jogging -

Setting the Jog Move - bit in the control word, the servo moves either a defined position decrement or to the minimal position with a limited speed. Releasing the bit will stop the motion.

### 3.18 State 18: Linearizing

The linearizing state is used to correct position feedback parameters, to improve the linearity of the position feedback.

### 3.19 State 19: Phase Searching

The phase search is only defined for three phase EC motors with hall switches and ABZ-sensors to find the commutation offset for to the sensor. It cannot be guaranteed that this feature will work for all kinds of EC motors. The found offset can be found in the variable section Calculated Commutation Offset (UPID: 1C1Bh), and has to be set manually to the parameter Phase Angle (UPID 11F2h).

### 3.20 State 20: Special Mode

The Special Mode is available only on the B1100 drives. In this state the current command mode over the analog input is available.

### 3.21 Building the Control Word

The Control Word can be accessed bit by bit from different sources with different priorities. The highest priorities have the bits that are forced by parameters. The second highest priority has the control panel of the LinMot-Talk software, if logged in with the SW. The next lower priorities have the bits that are defined on the X4 IOs as control word input bits. The lowest priority have bits which are set over the interface (normally a serial fieldbus connection), so in the Ctrl Word Interface Copy mask all bits can be selected, without causing any problems, but bits which should not be accessed through the interface can be masked out.

### 3.22 Control Word

With the Control Word (16-bit) the main state machine of the drive can be accessed. Following table shows the meaning of each bit:

| Bit | Name | Value | Meaning | Remark |
|-----|------|-------|---------|--------|
| 0 | Switch On | 0 | OFF1 | A-Stop, → Current = 0, power switches disabled |
| | | 1 | ON | State change from switch on disabled to ready to switch on |
| 1 | Voltage Enable | 0 | OFF2 | Power switches disabled without microcontroller action |
| | | 1 | Operation | - |
| 2 | /Quick Stop | 0 | OFF3 | Quick Stop → Current = 0 → H-Bridges disabled |
| | | 1 | Operation | - |
| 3 | Enable Operation | 0 | Operation disabled | Position controller active Motion Commands disabled |
| | | 1 | Operation enable | Position controller active Motion Commands enabled |
| 4 | /Abort | 0 | Abort | Quick Stop position control rests active, motion command is cleared |
| | | 1 | Operation | - |
| 5 | /Freeze | 0 | Freeze motion | Quick Stop position control rests active, Target position not cleared, curves motions are aborted |
| | | 1 | Operation | Rising edge will reactivate motion command |
| 6 | Go To Position | 0 | - | Go to fixed parameterized Position. Wait for release of signal |
| | | 1 | Go To Position | - |
| 7 | Error Acknowledge | 0 | - | Rising edge of signal acknowledges error |
| | | 1 | Error Acknowledge | - |
| 8 | Jog Move + | 0 | - | - |
| | | 1 | Jog Move + | - |
| 9 | Jog Move - | 0 | - | - |
| | | 1 | Jog Move - | - |
| 10 | Special Mode | 0 | Special Mode | - |
| | | 1 | - | - |
| 11 | Home | 0 | Stop Homing | At startup bit 11 Status word is cleared, until procedure is finished |
| | | 1 | Homing | - |
| 12 | Clearance Check | 0 | Stop Clearance Check | - |
| | | 1 | Clearance Check | Enable Clearance Check Movements |
| 13 | Go To Initial Position | 0 | - | Rising edge will start go to initial position |
| | | 1 | Go To initial Position | - |
| 14 | Reserved | 0 | Reserved | - |
| | | 1 | Reserved | - |
| 15 | Phase Search | 0 | Stop Phase Search | - |
| | | 1 | Phase Search | Enable Phase Search Movements |

### 3.23 Status Word

Following table shows the meaning of the single bits:

| Bit | Name | Value | Meaning | Remark |
|-----|------|-------|---------|--------|
| 0 | Operation Enabled | 0 | State Nr < 8 | - |
| | | 1 | Operation Enabled | State Nr 8 or higher (copied to drive EN LED) |
| 1 | Switch On Active | 0 | Switch On Disabled | Control Word Bit 0 |
| | | 1 | Switch On Enabled | - |
| 2 | Enable Operation | 0 | Operation Disabled | Control Word Bit 3 |
| | | 1 | Operation | - |
| 3 | Error | 0 | No Error | - |
| | | 1 | Error | Acknowledge with Control word Bit 7 (Reset Error) |
| 4 | Voltage Enable | 0 | Power Bridge Off | Control Word Bit 1 |
| | | 1 | Operation | - |
| 5 | /Quick Stop | 0 | Active | Control Word Bit 2 |
| | | 1 | Operation | - |
| 6 | Switch On Locked | 0 | Not Locked | - |
| | | 1 | Switch On Locked | Release with 0 of Control word bit 0 (Switch On) |
| 7 | Warning | 0 | Warning not active | No bit is set in the Warn Word |
| | | 1 | Warning active | One or more bits in the Warn Word are set |
| 8 | Event Handler Active | 0 | Event Handler Inactive | Event Handler cleared or disabled |
| | | 1 | Event Handler Active | Event Handler setup |
| 9 | Special Motion Active | 0 | Normal Operation | - |
| | | 1 | Special Command runs | Special motion commands (Homing, ..) runs |
| 10 | In Target Position | 0 | Not In Pos | Motion active or actual position out of window |
| | | 1 | In Pos | Actual position after motion in window |
| 11 | Homed | 0 | Motor not homed | Incremental sensor not homed (referenced) |
| | | 1 | Motor homed | Position sensor system valid |
| 12 | Fatal Error | 0 | - | - |
| | | 1 | Fatal Error | A fatal error can not be acknowledged! |
| 13 | Motion Active | 0 | No Motion | Setpoint generation inactive |
| | | 1 | Motion active | Setpoint generation (VAI, curve) active |
| 14 | Range Indicator 1 | 0 | Not In Range 1 | Defined UPID is not in Range 1 |
| | | 1 | In Range1 | Defined UPID is in Range 1 |
| 15 | Range Indicator 2 | 0 | Not In Range 2 | Defined UPID is not in Range 2 |
| | | 1 | In Range2 | Defined UPID is in Range 2 |

### 3.24 Warn Word

Following table shows the meaning of the single bits of the Warn Word:

| Bit | Name | Value | Meaning |
|-----|------|-------|---------|
| 0 | Motor Hot Sensor | 0 | Normal Operation |
| | | 1 | Motor Temperature Sensor On |
| 1 | Motor Short Time Overload I²t | 0 | Normal Operation |
| | | 1 | Calculated Motor Temperature Reached Warn Limit |
| 2 | Motor Supply Voltage Low | 0 | Normal Operation |
| | | 1 | Motor Supply Voltage Reached Low Warn Limit |
| 3 | Motor Supply Voltage High | 0 | Normal Operation |
| | | 1 | Motor Supply Voltage Reached High Warn Limit |
| 4 | Position Lag Always | 0 | Normal Operation |
| | | 1 | Position Error during Moving Reached Warn Limit |
| 5 | Reserved | 0 | - |
| | | 1 | - |
| 6 | Drive Hot | 0 | Normal Operation |
| | | 1 | Temperature on Drive High |
| 7 | Motor Not Homed | 0 | Normal Operation |
| | | 1 | Warning Motor Not Homed Yet |
| 8 | PTC Sensor 1 Hot | 0 | Normal Operation |
| | | 1 | PTC Temperature Sensor 1 On |
| 9 | Reserved PTC 2 | 0 | Normal Operation |
| | | 1 | PTC Temperature Sensor 2 On |
| 10 | RR Hot Calculated | 0 | Normal Operation |
| | | 1 | Regenerative Resistor Temperature Hot Calculated |
| 11 | Reserved | 0 | - |
| | | 1 | - |
| 12 | Reserved | 0 | - |
| | | 1 | - |
| 13 | Reserved | 0 | - |
| | | 1 | - |
| 14 | Interface Warn Flag | 0 | Normal Operation |
| | | 1 | Warn Flag Of Interface SW layer |
| 15 | Application Warn Flag | 0 | Normal Operation |
| | | 1 | Warn Flag Of Application SW layer |

Normally the warn word bits are used to react in conditions before the drive goes into the error state. E.g. a typical reaction on the warning 'Motor Temperature Sensor' would be a stop of the machine, before the drive goes into the error state and the motor goes out of control to avoid crashes.

---

## 4. Motion Command Interface

### 4.1 Motion Command Interface

The motion command interface consists of one word that contains the command ID, and up to 16 command parameter words.

**Example:** 'VA-Interpolator 16 bit Go To Absolute Position'

| Word | Description | Example of command |
|------|-------------|-------------------|
| 1. | Command Header with ID | Go To Absolute Position Immediate |
| 2. | 1. Command Parameter | Position |
| 3. | 2. Command Parameter | Maximal Speed |
| 4. | 3. Command Parameter | Acceleration |
| 5. | 4. Command Parameter | Deceleration |
| 6.-16. | 5. - Command Parameter | Not used |

#### 4.1.1 Command Header

The header of the Motion command is split into three parts:

- **Master ID** (bits 15-8): Specifies the command group
- **Sub ID** (bits 7-4): Used to identify different commands from the same command group
- **Command Count** (bits 3-0): A new command will only be executed, if the value of the command count has changed. In the easiest way bit 0 can be toggled.

**Command Header Structure:**

```
Bits:  15 14 13 12 11 10  9  8 |  7  6  5  4 |  3  2  1  0
       └───── Master ID ──────┘ └─ Sub ID ─┘ └─ Command Count ─┘
```

##### 4.1.1.1 Master ID

The master ID specifies the command group.

##### 4.1.1.2 Sub ID

The sub ID is used to identify different commands from the same command group.

##### 4.1.1.3 Command Count

A new command will only be executed, if the value of the command count has changed. In the easiest way bit 0 can be toggled.

### 4.2 Overview Motion Commands

The motion commands are organized by Master ID and Sub ID. For a complete list of all motion commands, see the detailed descriptions in section 4.3.

**Key Command Groups:**
- **00h**: Basic operations (No Operation, Write Interface Control Word, etc.)
- **01h**: VAI commands (Go To Pos, Increment Dem Pos, etc.)
- **02h**: Predef VAI commands
- **03h**: Streaming commands (P Stream, PV Stream, PVA Stream)
- **04h**: Time Curve commands
- **05h**: Curve modification commands
- **06h**: Encoder CAM commands
- **07h**: Encoder position indexing
- **09h**: VAI 16-bit commands
- **0Ah**: Predef VAI 16-bit commands
- **0Bh**: VAI Predef Acc commands
- **0Ch**: VAI Dec=Acc commands
- **0Dh**: Command table variable commands
- **0Eh**: Sin VA commands
- **0Fh**: Bestehorn VAJ commands
- **10h**: Encoder CAM enable/disable
- **1yh**: Encoder CAM y commands (y = 1 or 2)
- **20h**: Command Table commands
- **22h**: Wait commands
- **24h**: Command Table variable operations
- **25h**: IF condition commands
- **30h**: Encoder winding commands
- **38h**: Force Control commands
- **39h**: Current Command Mode

### 4.3 Detailed Motion Command Description

This section provides detailed descriptions of all motion commands. Each command includes:
- Command ID (hex format)
- Parameter layout
- Byte offsets
- Data types
- Units
- Description

#### 4.3.1 No Operation (000xh)

| Field | Byte Offset | Description | Type | Unit |
|-------|-------------|-------------|------|------|
| Header | 0 | No Operation (000xh) | UInt16 | - |

This command does nothing. It can be sent in any operational state.

#### 4.3.2 Write Interface Control Word (001xh)

| Field | Byte Offset | Description | Type | Unit |
|-------|-------------|-------------|------|------|
| Header | 0 | Write Interface Control Word (001xh) | UInt16 | - |
| 1. Par | 2 | Interface Control Word | UInt16 | - |

This command allows writing the control word through the motion command interface. The fieldbus interfaces (CANOpen, DeviceNet, Profibus, LinRS, POWERLINK, EtherCAT) offer other ways to access the control word directly. Mostly a direct access is more comfortable than the way over the motion command interface.

#### 4.3.3 Write Live Parameter (002xh)

| Field | Byte Offset | Description | Type | Unit |
|-------|-------------|-------------|------|------|
| Header | 0 | Write Live Parameter (002xh) | UInt16 | - |
| 1. Par | 2 | UPID (Unique Parameter ID) | UInt16 | - |
| 2. Par | 4 | Parameter Value | UInt32 | Depends on Parameter |

This command allows writing any live parameter's RAM value through the motion command interface. The parameter has to be specified by its UPID (Unique Parameter ID). In order to keep the interface as simple as possible any parameter can be accessed as 32-bit integer value. The drive's operating system will filter out the relevant number of bits for parameters with smaller data size (e.g. only the lowest bit is considered for Boolean parameters).

The fieldbus interfaces (CANOpen, DeviceNet, Profibus, LinRS, POWERLINK, EtherCAT) offer other ways to read and write parameter values directly. Mostly a direct access is more comfortable than the way over the motion command interface.

#### 4.3.4 Write X4 Intf Outputs with Mask (003xh)

| Field | Byte Offset | Description | Type | Unit |
|-------|-------------|-------------|------|------|
| Header | 0 | Write X4 Intf Outputs with Mask (003xh) | UInt16 | - |
| 1. Par | 2 | Bit Mask; Bit 0 = X4.3 Bit 1 = X4.4… | UInt16 | - |
| 2. Par | 4 | Bit Value; Bit 0 = X4.3, Bit 1 X4.4… | UInt16 | - |

This command allows writing the configured X4 interface outputs with a write mask through the motion command interface. To write an output, the corresponding bit in the mask must be set. Bit 0 is mapped to output X4.3, bit 1 to output X4.4 etc.

#### 4.3.5 Select Position Controller Set (005xh)

| Field | Byte Offset | Description | Type | Unit |
|-------|-------------|-------------|------|------|
| Header | 0 | Select Position Controller Set (005xh) | UInt16 | - |
| 1. Par | 2 | Controller Set Selection (0 = Set A, 1 = Set B) | UInt16 | - |

This command selects the active position controller set (A/B) UPID 0x1393. For set A the ID is 0 and for Set B the ID is 1.

#### 4.3.6 Clear Event Evaluation (008xh)

| Field | Byte Offset | Description | Type | Unit |
|-------|-------------|-------------|------|------|
| Header | 0 | Clear Event Evaluation (008xh) | UInt16 | - |

This command resets the event handler. The event handler becomes active, if a motion command has been sent, that does not immediately start, but waits with its execution until other conditions are fulfilled (e.g. command 'VAI Go To Pos On Rising Trigger Event'). The bit 8 of the status word shows, if the event handler is active.

Once the event handler becomes active, it remains active, until it is deactivated with this clear command. As long the event handler is active, the command to be executed on the event situation will be restarted each time the event condition is fulfilled.

#### 4.3.7 Master Homing (009xh)

| Field | Byte Offset | Description | Type | Unit |
|-------|-------------|-------------|------|------|
| Header | 0 | Master Homing (009xh) | UInt16 | - |
| 1. Par | 2 | Home Position | SInt32 | 0.1 µm |

This command can be used, if the master system knows the home position without going to the home state in the state machine. The passed value of the home position is stored in the RAM value of the parameter Home Position (UPID 13C7h), then the corresponding value of the parameter Slider Home Position (UPID 13CAh) is calculated and stored in the RAM value. Then a homing at actual position is done without going into the homing state.

#### 4.3.8 Reset (00Fxh)

| Field | Byte Offset | Description | Type | Unit |
|-------|-------------|-------------|------|------|
| Header | 0 | Reset (00Fxh) | UInt16 | - |

This command resets all firmware instances of the drive. Use this command with count = 0, otherwise the drive reboots cyclic!

#### 4.3.9 VAI Go To Pos (010xh)

| Field | Byte Offset | Description | Type | Unit |
|-------|-------------|-------------|------|------|
| Header | 0 | VAI Go To Pos (010xh) | UInt16 | - |
| 1. Par | 2 | Target Position | SInt32 | 0.1 µm |
| 2. Par | 6 | Maximal Velocity | UInt32 | 1E-6 m/s |
| 3. Par | 10 | Acceleration | UInt32 | 1E-5 m/s² |
| 4. Par | 14 | Deceleration | UInt32 | 1E-5 m/s² |

This command sets a new target position and defines the maximal velocity, acceleration and deceleration for going there. The command execution starts immediately when the command has been sent.

---

*Note: Due to the extensive nature of this document (over 200 motion commands), the remaining command descriptions follow the same format. For complete documentation of all commands, refer to the original source material or LinMot-Talk software.*

---

## 5. Setpoint Generation

### 5.1 VA-Interpolator

The VA-Interpolator (Velocity-Acceleration Interpolator) is a setpoint generation algorithm that creates smooth motion profiles with limited velocity and acceleration.

**Parameters and Output:**
- Input: Target position, maximum velocity, acceleration, deceleration
- Output: Position setpoint, velocity setpoint
- Characteristics: Trapezoidal or triangular velocity profile

### 5.2 Sine VA Motion

Sine VA motion generates sinusoidal velocity profiles for smooth motion.

**Parameters and Output:**
- Similar to VAI but with sinusoidal velocity transitions
- Provides smoother acceleration/deceleration curves

### 5.3 Bestehorn VAJ Motion

Bestehorn VAJ motion includes jerk limitation in addition to velocity and acceleration limits.

**Parameters and Output:**
- Input: Target position, maximum velocity, acceleration, jerk
- Output: Position setpoint, velocity setpoint, acceleration setpoint
- Characteristics: S-curve motion profiles with jerk limitation

### 5.4 P(V)-Stream

P(V)-Stream allows streaming position (and optionally velocity) setpoints in real-time.

**Characteristics:**
- Real-time streaming of position/velocity data
- Requires cyclic communication
- Supports slave-generated timestamps or configured period times

### 5.5 CAM Motions

CAM motions synchronize slave motion to a master encoder.

#### 5.5.1 Triggered Cam Motions

CAM motions can be triggered by encoder events with configurable delay counts.

#### 5.5.2 Repeated Cam Motions with the Modulo CamMode

CAM motions can be repeated continuously using modulo arithmetic relative to the master encoder length.

---

## 6. Command Table

⚠️ **Note:** The command table in the B1100 drives is limited to 31 entries, and is flash only, for this reason it is not possible to write or modify the table over a serial interface!

The command table functionality can be used for programming sequences directly in the drive. Command table entries can contain:
- Motion commands
- Wait conditions
- Variable operations
- Conditional branching (IF commands)
- Loops

**Example Command Table Sequence:**

1. Set Command Table Var 1 = 0
2. GoTo 50mm
3. Wait Until In Target Position
4. GoTo 0mm
5. Wait until Motion Done
6. Increment Command Table Var 1
7. If Command Table Var 1 < 5 Then GoTo step 2 Else End
8. No Operation (End of Sequence)

---

## 7. Drive Configuration

The parameter configuration is normally done with LinMot-Talk software. The UPIDs, over which the parameter can be accessed, are the same for E1100, E1200 and E1400 drives, but are different for the B1100 drives. In this documentation the E1100/E1200/E1400 UPIDs are used. If a UPID for a B1100 drive is needed, a conversion list can be generated with the LinMot-Talk software.

### 7.1 Power Bridge

The E1100/B1100 drives are divided into three different power classes:
- **Normal drives**: Maximal current of 8A
- **High Current (HC)**: Maximal current of 15A
- **Extreme Current (XC)**: Maximal current of 25A

The E1200 series is only available as ultra current drives (name extension UC), with a current maximum of 32A.

### 7.2 X4 I/O Definitions

The functionality of most IOs can be programmed as a control word input bit or status word output bit, or they can be used as interface IO and read out or written over a serial bus interface. Apart from this general functionality a few IOs have a special functionality.

| IO | General Purpose IO | Special Functions |
|----|-------------------|-------------------|
| X4.3 | Yes | Brake control |
| X4.4 | Yes | - |
| X4.5 | Yes | - |
| X4.6 | Yes | Trigger input |
| X4.7 | Yes | - |
| X4.8 | Yes | Limit switch |
| X4.9 | Yes | Limit switch |
| X4.10 | Yes | PTC 1 sensor |
| X4.11 | Yes | PTC 2 sensor |
| X4.12 | No | SVE (Safety Voltage Enable) |

#### 7.2.1 X4.3 Brake

The X4.3 brake control has different behaviors depending on the operation state:
- **Operation Enabled Behavior**: Brake is released when operation is enabled
- **Operation /Abort Behavior**: Brake is applied on abort
- **Operation Quick Stop Behavior**: Brake behavior during quick stop

#### 7.2.2 X4.6 Trigger

The X4.6 trigger input can be configured in different modes:
- **Direct Trigger Mode**: Immediate response to trigger events
- **Inhibited Trigger Mode**: Trigger is inhibited until conditions are met
- **Delayed Trigger Mode**: Trigger events are delayed by a configurable count
- **Inhibited & Delayed Trigger Mode**: Combination of both modes

#### 7.2.3 X4.8 and X4.9 Limit Switches

These inputs can be configured as limit switches to prevent motion beyond safe positions.

#### 7.2.4 X4.10 and X4.11 PTC 1 and PTC 2

PTC (Positive Temperature Coefficient) sensors for thermal protection.

#### 7.2.5 X4.12 SVE (Safety Voltage Enable)

Safety voltage enable input for safety-related functions.

### 7.3 Master Encoder

Master encoder configuration for CAM motions and synchronization.

### 7.4 Monitoring

#### 7.4.1 Logic Supply Voltage

Monitoring of the 24V logic supply voltage with configurable warning and error thresholds.

#### 7.4.2 Motor Supply Voltage

Monitoring of the motor power supply voltage.

##### 7.4.2.1 Phase Switch On Test

The parameters in the Phase Switch On Test section are used in the HW Tests State (State 5) before enabling the power stage. If the motor power supply is on and there is no ground path in the motor (inclusive cabling) the phase voltage is approx. 6.5V when in power off state.

| Parameter Name | UPID | Description |
|----------------|------|-------------|
| Phase Voltage Low Level | 102Ch | If one of the phase voltages is below this limit before powering up an error will be generated |
| Phase Voltage High Level | 102Dh | If one of the phase voltages is above this limit before powering up an error will be generated |
| Phase Test Max Incurrent | 102Eh | If the current rises above this limit if one edge of a phase is set to a voltage an error will be generated |

#### 7.4.3 Regeneration Resistor

The regeneration resistor terminals on X1 can be used for energy dissipation, when the motor is decelerating.

**Configuration Parameters:**
- Enable/Disable (UPID 101Dh)
- Turn On Voltage (UPID 101Eh)
- Turn Off Voltage (UPID 101Fh)
- RR Resistance (UPID 1022h)
- Warning Temp (UPID 1024h)
- Error Temp (UPID 1025h)

⚠️ **Important:** The turn on voltage has to be at minimum 0.5V higher than the turn off voltage. Ensure that the idle motor supply voltage is lower than the turn off voltage!

#### 7.4.4 Temperature Monitoring

The E1100 drive hardware contains eight absolute temperature sensors for thermal protection. On the B1100 drive is one sensor placed.

| Parameter Name | UPID | Description |
|----------------|------|-------------|
| Temp Sens Warn Level | 1040h | If the maximal board temperature rises above this level, a warning is generated (bit 6 in Warn Word is set) |
| Temp Sens Error Level | 1041h | If the maximal board temperature rises above this level, the error is generated (error codes 10h..17h) |

### 7.5 PosCtrlStructure

Position controller structure configuration parameters.

---

## 8. Motor Configuration

The motor usually is set up with the motor wizard, which sets all needed parameters. Therefore a detailed description of the parameters will follow in the future.

### 8.1 Generic Motor Temperature Calculated

For third parties motors a generic calculated motor temperature model is used to adapt the winding resistance and to detect excess temperature.

**Temperature Model Parameters:**
- C Winding (UPID 120Ch): Heat capacity of the motor winding
- R Winding-Housing (UPID 1210h): Thermal resistance between winding and housing
- C Housing (UPID 1211h): Heat capacity of the motor housing
- R Housing-Mounting (UPID 1212h): Thermal resistance between housing and mounting
- C Mounting (UPID 1213h): Heat capacity of the mounting
- R Mounting-Environment (UPID 1214h): Thermal resistance between mounting and environment

The sum of all R defines the static power losses (consider also TW and TE). With the capacitance the thermal time constant can be influenced. The bigger the thermal capacitance the slower the temperature will rise.

---

## 9. State Machine Setup

In the state machine setup sections the parameters to influence the behavior of the single states can be defined.

---

## 10. Error Code List

| Code | Description | Actions to Take |
|------|-------------|------------------|
| 0000h | No Error | No error is pending |
| 0001h | Err: X4 Logic Supply Too Low | The logic supply voltage has been too low. Check your 24V logic power supply |
| 0002h | Err: X4 Logic Supply Too High | The logic supply voltage has been too high. Check your 24V logic power supply |
| 0003h | Err: X1 Pwr Voltage Too Low | The motor power supply voltage has been too low. Check your motor power supply, check the wiring, check the sizing of the power supply, add a capacitor to enforce your DC link |
| 0004h | Err: X1 Pwr Voltage Too High | The motor power supply voltage has been too high. Check your motor power supply, check the wiring, check the sizing of the power supply, use a regeneration resistor for power dissipation, add a capacitor to enforce your DC link |
| 0005h | Err: X1 RR Not Connected | A regeneration resistor is configured but not connected. Connect the regeneration resistor to X1 |
| 0006h | Err: PTC 1 Sensor Too Hot | The PTC 1 sensor on X4.10 is hot or not connected. Check the temperature, check the wiring |
| 0007h | Err: Min Pos Undershot | The motor position has been below the minimal position. Check the configuration, check the PLC program |
| 0008h | Err: Max Pos Overshot | The motor position has been above the maximal position. Check the configuration, check the PLC program |
| 0009h | Err: Ext-Int Sensor Diff Err | The position difference between sensor feedback on X3 and sensor feedback on X12 has been too big. Check sensor wiring, check sensor configuration (count direction, etc.), check parameter 1266h |
| 000Ah | Fatal Err: X12 Signals Missing | The external sensor is not connected to X12 or the wiring is not ok. Check the wiring |
| 000Bh | Err: Pos Lag Always Too Big | The motor was not able to follow the demand position. Check the motor load, check the motor stroke range for possible collisions, check the position controller setup, check the setpoint generation (unreachable speed/acceleration values?), check the motor sizing |
| 000Ch | Err: Pos Lag Standing Too Big (Not on B1100) | The motor was not able to reach the target position or was not able to stay at the target position. Check the motor load, check the motor stroke range for possible collisions, check the position controller setup, check the motor sizing |
| 000Dh | Fatal Err: X1 Pwr Over Current | Over current on X1 detected. Check motor wiring, check motor configuration, for P01-48 type motors: set parameter 11F4h to value 0001h |
| 000Eh | Err: Supply Dig Out Missing | Drive board defective. Contact support for repair |
| 000Fh | Err: PTC 2 Sensor Too Hot | The PTC 2 sensor on X4.11 is hot or not connected. Check the temperature, check the wiring |
| 0010h | Err: Drive Ph1+ Too Hot | Drive power bridge phase 1+ too hot. Check motor wiring |
| 0011h | Err: Drive Ph1- Too Hot | Drive power bridge phase 1- too hot. Check motor wiring |
| 0012h | Err: Drive Ph2+ Too Hot | Drive power bridge phase 2+ too hot. Check motor wiring |
| 0013h | Err: Drive Ph2- Too Hot | Drive power bridge phase 2- too hot. Check motor wiring |
| 0014h | Err: Drive Pwr Too Hot | DC link temp sensor has detected over temperature. Check wiring |
| 0015h | Err: Drive RR Hot Calc | Regeneration resistor switch hot. Check RR configuration (Turn On level, Resistance, etc.), check RR sizing |
| 0016h | Err: Drive X3 Too Hot | Temp sensor on X3 has detected over temperature. Check motor wiring |
| 0017h | Err: Drive Core Too Hot | Temp sensor on drive's PCB board reports core being hot. Drive power bridge phase 1+ may be defective. Contact support |
| 0018h | Err: Power Bridge Ph1+ Defective | Drive power bridge phase 1+ may be defective. Contact support |
| 0019h | Err: Power Bridge Ph1- Defective | Drive power bridge phase 1- may be defective. Contact support |
| 001Ah | Err: Power Bridge Ph2+ Defective | Drive power bridge phase 2+ may be defective. Contact support |
| 001Bh | Err: Power Bridge Ph2- Defective | Drive power bridge phase 2- may be defective. Contact support |
| 001Ch | Err: Supply DigOut X6 Fuse Blown | Supply fuse for digital outputs on X6 blown. Check X6 wiring, contact support for repair |
| 001Dh | Err: Supply X3.3 5V Fuse Blown | Supply X3.3 5V fuse blown. Motor or and/or wiring defective. Contact support for drive repair, check motor and wiring, replace motor and motor cables |
| 001Eh | Err: Supply X3.8 AGND Fuse Blown | Supply X3.8 analog ground fuse blown. Contact support for drive repair, check motor and wiring, replace motor and motor cables |
| 0020h | Err: Motor Hot Sensor | Temp sensor reports hot motor. Wait until motor has cooled down (until corresponding warning disappears), check load, check the motor configuration, check the setpoint generation (unreachable speed/acceleration values?), check the motor sizing |
| 0021h | Fatal Err: X3 Hall Sig Missing | Motor hall signals not connected to X3 or motor defective. Power down the drive and all power supplies, then reconnect motor, check motor and wiring, check parameter 1221h |
| 0022h | Fatal Err: Motor Slider Missing | Motor hall sensors cannot see magnetic field of the slider. The motor position was outside the allowed range defined through the motors ZP and Max Stroke data. Check stroke range, check slider orientation |
| 0023h | Err: Motor Short Time Overload | Short time motor overload detected. Check if motor is blocked, check motor sizing |
| 0024h | Err: RR Hot Calculated | Regeneration resistor hot calculated. Check RR configuration (Turn On level, Resistance, etc.), check RR sizing |
| 0025h | Err: Sensor Alarm | Sensor Alarm On X12 Occurred. Check sensor mounting, band contamination or motion speed |
| 0028h | Err: Ph1+ Short Circuit To GND | Short circuit between phase 1+ and ground detected. Check motor wiring, check motor |
| 0029h | Err: Ph1- Short Circuit To GND | Short circuit between phase 1- and ground detected. Check motor wiring, check motor |
| 002Ah | Err: Ph2+ Short Circuit To GND | Short circuit between phase 2+ and ground detected. Check motor wiring, check motor |
| 002Bh | Err: Ph2- Short Circuit To GND | Short circuit between phase 2- and ground detected. Check motor wiring, check motor |
| 002Ch | Err: Ph1 Short Circuit To Ph2 | Short circuit between motor phase 1 and phase 2 detected. Check motor wiring, check motor |
| 0030h | Err: Ph1+ Wired To Ph2+ | Motor phase 1+ has contact to phase 2+. Check motor wiring, check motor |
| 0031h | Err: Ph1+ Wired To Ph2- | Motor phase 1+ has contact to phase 2-. Check motor wiring, check motor |
| 0032h | Err: Ph1+ Not Wired To Ph1- | Motor phase 1+ has no connection to phase 1-. Check motor wiring, check motor |
| 0033h | Err: Ph2+ Wired To Ph1+ | Motor phase 2+ has contact to phase 1+. Check motor wiring, check motor |
| 0034h | Err: Ph2+ Wired To Ph1- | Motor phase 2+ has contact to phase 1-. Check motor wiring, check motor |
| 0035h | Err: Ph2+ Not Wired To Ph2- | Motor phase 2+ has no connection to phase 2-. Check motor wiring, check motor |
| 0036h | Err: Ph1 Short Circuit To Ph2+ | Short circuit between motor phase 1 and phase 2+ detected. Check motor wiring, check motor |
| 0037h | Err: Ph1 Short Circuit To Ph2- | Short circuit between motor phase 1 and phase 2- detected. Check motor wiring, check motor |
| 0038h | Err: Ph2 Short Circuit To Ph1+ | Short circuit between motor phase 2 and phase 1+ detected. Check motor wiring, check motor |
| 0039h | Err: Ph2 Short Circuit To Ph1- | Short circuit between motor phase 2 and phase 1- detected. Check motor wiring, check motor |
| 003Ah | Err: Phase U Broken | Motor phase U broken. Check motor wiring, check motor |
| 003Bh | Err: Phase V Broken | Motor phase V broken. Check motor wiring, check motor |
| 003Ch | Err: Phase W Broken | Motor phase W broken. Check motor wiring, check motor |
| 0040h | Err: X4.3 Brake Driver Error | X4.3 brake driver reports error. Check for short circuit on X4.3 |
| 0041h | Err: Dig Out X4.4..X4.11 Status | X4.3..X4.11 output driver reports error. Check for short circuit on outputs X4.4..X4.11 or output configurations |
| 0042h | Err: Dig Out X6 Status | X6 output driver reports error. Check for short circuit on outputs X6 |
| 0044h | Err: X4 Dig Out GND Fuse Blown | Ground fuse for digital outputs on X4 blown. Check X4 wiring, contact support for repair |
| 0045h | Fatal Err: Motor Comm Lost | Motor communication lost. Power down and check motor wiring and motor, replace cable and/or motor |
| 0046h | Err: PTC 1 Broken | PTC 1 on X4.10 broken or not connected. Power down and check PTC 1 wiring and resistance |
| 0047h | Err: PTC 1 Short To 24V | PTC 1 on X4.10 short to 24V. Power down and check PTC 1 wiring and resistance |
| 0050h | Setup Err: HW Not Supported | Setup error, hardware is not supported by the software. Download correct firmware, contact support |
| 0051h | Setup Err: SW Key Missing | Software key and access code for special functionality is missing. Order the SW key with your support together with the serial number of your HW |
| 0058h | Runtime Err: ROM write error | Runtime error, MC SW was not able to change parameter value in ROM. Verify PLC is not configuring during this action, contact support |
| 0060h | Cfg Err: RR Voltage Set Too Low | Configuration error: regeneration resistor turn on/off voltage parameter value is too low. Check parameters 101Eh and 101Fh |
| 0061h | Cfg Err: RR Hysteresis < 0.5V | Configuration error: regeneration resistor turn on/off voltage parameter values too close to each other. Check parameters 101Eh and 101Fh |
| 0062h | Cfg Err: Curve Not Defined | Configuration error. Software tried to start a curve that is not defined yet. Define the curve using the curves service, check if curves were downloaded to drive, check the curve IDs, check the configuration, check the PLC program |
| 0063h | Cfg Err: Pos Ctrl Max Curr High | Configuration error: Invalid max current setting in control parameters. Check parameters 13A6h and 13BAh, check PLC program |
| 0064h | Cfg Err (Fatal): No Motor Defined | Configuration error: No motor has been configured yet. Use the motor wizard to configure the motor |
| 0065h | Cfg Err (Fatal): No Trigger Mode Defined | Configuration error: Digital input X4.6 is configured for trigger input function, but the trigger mode is not defined yet. Configure parameter 170Ch |
| 0067h | Cfg Err (Fatal): Wrong Stator Type | Configuration error: The configured motor type does not match with the connected motor. Configure correct motor type by using the motor wizard, connect an appropriate motor |
| 0068h | Cfg Err (Fatal): No Motor Communication | Configuration error: The drive was not able to establish the communication to the microcontroller on the motor. Older P01 motors don't support motor communication. Check motor wiring, check motor, check the motor configuration, disable communication by using parameter 11FBh if you have an old P01 motor |
| 0069h | Cfg Err (Fatal): Wrong Slider | Configuration error: A wrong slider has been configured or slider home position has an invalid value. Reconfigure the motor by using the motor wizard |
| 0080h | User Err: Lin: Not Homed | User error: The PLC program tried to start an action that requires the motor to be already homed, but the motor was not homed. Check the PLC program, do a homing of the motor first |
| 0081h | User Err: Unknown Motion Cmd | User error: The PLC program sent an unknown motion command ID. Check PLC program, check firmware version |
| 0082h | User Err: PVT Buffer Overflow | User error: The PLC program has sent the stream position commands too fast, the buffer had an overflow. Streaming has to be strictly cyclic! Check PLC program, check the fieldbus by using bus monitor tools |
| 0083h | User Err: PVT Buffer Underflow | User error: The PLC program has sent the stream position commands too slowly, the buffer had an underflow. Streaming has to be strictly cyclic! Check PLC program, check the fieldbus by using bus monitor tools |
| 0084h | User Err: PVT Master Too Fast | User error: The PLC program has begun to send PVT streaming command. The commands were too close to each other. The drive expects new streaming commands every 2ms to 5ms. Check PLC program, check the fieldbus by using bus monitor tools |
| 0085h | User Err: PVT Master Too Slow | User error: The PLC program has begun to send PVT streaming command. The cycle time between the streaming commands has been too long. The drive expects new streaming commands every 2ms to 5ms. Check PLC program, check the fieldbus by using bus monitor tools |
| 0086h | User Err: Motion Cmd In Wrong St | User error: The PLC program has sent a motion command while the drive was not in an appropriate operational state. Most of the motion commands are accepted only in operational state 8 (Operation Enabled). Check the PLC program |
| 0087h | User Err: Limit Switch In High | User error: The motor moved into the Limit Switch In while it was still in the stroke range. Check the PLC program or check homing |
| 0088h | User Err: Limit Switch Out High | User error: The motor moved into the Limit Switch Out while it was still in the stroke range. Check the PLC program or check homing |
| 0089h | User Err: Curve Amp Scale Error | User error: The automatic calculated amplitude scale is out of range -2000percent to 2000percent. Check the PLC program or use other curve |
| 008Ah | User Err: Cmd Tab Entry Not Def | Called command Table entry is not defined. Check the PLC program or define Command Table Entry |

---

## 11. Contact Addresses

### SWITZERLAND

**NTI AG**  
Haerdlistr. 15  
CH-8957 Spreitenbach

- **Sales and Administration:** +41-(0)56-419 91 91 | office@linmot.com
- **Tech. Support:** +41-(0)56-544 71 00 | support@linmot.com
- **Tech. Support (Skype):** skype:support.linmot
- **Fax:** +41-(0)56-419 91 92
- **Web:** http://www.linmot.com/

### USA

**LinMot, Inc.**  
204 E Morrissey Dr.  
Elkhorn, WI 53121

- **Sales and Administration:** 877-546-3270 / 262-743-2555
- **Tech. Support:** 877-804-0718 / 262-743-1284
- **Fax:** 800-463-8708 / 262-723-6688
- **E-Mail:** us-sales@linmot.com
- **Web:** http://www.linmot-usa.com/

**Find local distribution:** http://www.linmot.com/contact

---

## Copyright Notice

© 2014 NTI AG

This work is protected by copyright. Under the copyright laws, this publication may not be reproduced or transmitted in any form, electronic or mechanical, including photocopying, recording, microfilm, storing in an information retrieval system, not even for didactic use, or translating, in whole or in part, without the prior written consent of NTI AG.

**LinMot®** is a registered trademark of NTI AG.

**Note:** The information in this documentation reflects the stage of development at the time of press and is therefore without obligation. NTI AG reserves itself the right to make changes at any time and without notice to reflect further technical advance or product improvement.
