# LinUDP V2 Interface

**Documentation of the LinUDP V2 Interface**

**Supported Drives:**
- E1250-LU-UC
- C1250-LU-XC
- E1450-LU-QN V2

**Document:** 0185-1108-E_1V9_MA_LinUDP_V2  
**Date:** November 2019

---

## Table of Contents

1. [Overview](#1-overview)
2. [Installation on Servo Drive](#2-installation-on-servo-drive)
3. [Connecting LinUDP V2](#3-connecting-linudp-v2)
   - 3.1 [Pin Assignment of Connectors X17-X18](#31-pin-assignment-of-connectors-x17-x18)
   - 3.2 [Default IP Address Settings](#32-default-ip-address-settings)
4. [LinUDP V2 Telegram](#4-linudp-v2-telegram)
   - 4.1 [DHCP Header](#41-dhcp-header)
   - 4.2 [IPv4 Header](#42-ipv4-header)
   - 4.3 [LinUDP V2 Header](#43-linudp-v2-header)
   - 4.4 [LinUDP V2 Data](#44-linudp-v2-data)
5. [LinUDP V2 Parameters](#5-linudp-v2-parameters)
6. [LinUDP V2 Modules](#6-linudp-v2-modules)
   - 6.1 [Master to Drive Modules](#61-master-to-drive-modules)
   - 6.2 [Drive to Master Modules](#62-drive-to-master-modules)
7. [Real Time Config Module](#7-real-time-config-module)

---

## 1. Overview

The **LinUDP V2 protocol** is an easy way for communication with a LinMot servo drive over Ethernet. 

**Key characteristics:**
- No checks are done to ensure messages reach their destination or are correctly received
- The drive has no active function - it only responds to requests with appropriate answers
- LinUDP V2 is the second version of LinUDP
- Functions on drives with "LU" in the name (e.g., C1250-LU-XC)
- The old LinUDP version does not run on these drives

---

## 2. Installation on Servo Drive

To install LinUDP V2 firmware:

1. Start the **LinMot-Talk** software
2. Press the **Install Firmware** button
3. Choose the file `Firmware_Buildxxxxxxxx.sct`
4. Press **Open**
5. Follow the wizard through installation

The installation will install LinUDP V2 on drives with "LU" in the name:
- C1250-LU-XC-xS-000
- E1250-LU-UC
- E1450-LU-QN-xS
- C1450-LU-VS-xS-000

---

## 3. Connecting LinUDP V2

### 3.1 Pin Assignment of Connectors X17-X18

The Ethernet/IP connector is a standard **RJ45 female connector** with pin assignment as defined by **EIA/TIA T568B**.

| Pin | Wire Color Code | Assignment 100BASE-TX |
|-----|-----------------|----------------------|
| 1   | WHT/ORG        | Rx+                  |
| 2   | ORG            | Rx-                  |
| 3   | WHT/GRN        | Tx+                  |
| 4   | BLU            | -                    |
| 5   | WHT/BLU        | -                    |
| 6   | GRN            | Tx-                  |
| 7   | WHT/BRN        | -                    |
| 8   | BRN            | -                    |

**Cable specification:** Use standard patch cables (twisted pair, S/UTP, AWG26), typically referred to as "Cat5e-Cable".

### 3.2 Default IP Address Settings

The default IP address is **192.168.001.xxx**, where the last byte `xxx` is defined via two hex switches **S1** and **S2**.

| Switch | Description |
|--------|-------------|
| S1 (bits 5-8) | Bus ID High (0…F). Bit 5 is LSB, bit 8 is MSB |
| S2 (bits 1-4) | Bus ID Low (0…F). Bit 1 is LSB, bit 4 is MSB |

⚠️ **Important:**
- Setting both ID high & low to `0xFF` resets the drive to manufacturer settings
- The switch value `S1 = S2 = 0` (factory default) is a special configuration that acquires the IP address via **DHCP**

---

## 4. LinUDP V2 Telegram

In LinUDP V2, there are two telegrams: one for the **request from the master** and one for the **response from the drive**.

### Telegram Structure

| Name | Size (Bytes) |
|------|--------------|
| DHCP Header | 14 |
| IPv4 Header | 20 |
| LinUDP Header | 8 |
| LinUDP Data | message dependent |

### 4.1 DHCP Header

Standard Ethernet frame header:

| Offset | Field |
|--------|-------|
| +0 | Destination MAC ID (6 bytes) |
| +6 | Source MAC ID (6 bytes) |
| +12 | Protocol Type (0x0800) |

### 4.2 IPv4 Header

The IPv4 header is described in **RFC 0791** chapter 3.1.  
Available at: www.ietf.org/rfc/rfc0791.txt

**Note:** The sections "options" and "padding" are not used.

### 4.3 LinUDP V2 Header

The LinUDP V2 header consists of standard UDP fields:

| Name | Size (Bytes) |
|------|--------------|
| Source Port | 2 |
| Destination Port | 2 |
| Length of UDP Telegram | 2 |
| UDP Checksum | 2 |

**Fixed Port Assignments:**
- **Master Port:** 41136 (0xA0B0 hex)
- **Drive Port:** 49360 (0xC0D0 hex)

### 4.4 LinUDP V2 Data

The LinUDP V2 data part has the same construction for requests and responses; only the source and destination are switched.

**Byte Order:** Low Byte First (Little Endian)

#### 4.4.1 Request from the Master

The first 32 bits define the **request**, and the following 32 bits define the format of the **response**.

**Request Definition:**

| Bit | Name | Data Size (Bytes) |
|-----|------|-------------------|
| 0 | Control Word | 2 |
| 1 | Motion Control | 32 |
| 2 | Realtime Configuration | 8 |
| 3-31 | Reserved for future expansions | - |

**Response Definition:**

| Bit | Name | Data Size (Bytes) |
|-----|------|-------------------|
| 0 | Status Word | 2 |
| 1 | State Var | 2 |
| 2 | Actual Position | 4 |
| 3 | Demand Position | 4 |
| 4 | Current | 2 |
| 5 | Warn Word | 2 |
| 6 | Error Code | 2 |
| 7 | Monitoring Channel | 16 |
| 8 | Realtime Configuration | 8 |
| 9-31 | Reserved for future expansions | - |

**How it works:**
- Each definition bit indicates whether the corresponding parameter is part of the communication
- The order of requested data parts matches the order of definition bits
- When a definition bit is not set, that data part is not transferred

**Full Request Frame (when all bits set):**

| Offset | Field |
|--------|-------|
| +0 | Request Definition (4 bytes) |
| +4 | Response Definition (4 bytes) |
| +8 | Control Word (2 bytes) |
| +10 | Motion Control (32 bytes) |
| +42 | Realtime Configuration (8 bytes) |

#### 4.4.2 Response from the Drive

The response data part has the same construction as the request.

**Exception:** When Realtime Configuration is activated:
- Bit 2 is set in the request definition
- Bit 8 is set in the response definition

**Full Response Frame (when all bits set):**

| Offset | Field |
|--------|-------|
| +0 | Request Definition (4 bytes) |
| +4 | Response Definition (4 bytes) |
| +8 | Status Word (2 bytes) |
| +10 | State Var (2 bytes) |
| +12 | Actual Position (4 bytes) |
| +16 | Demand Position (4 bytes) |
| +20 | Current (2 bytes) |
| +22 | Warn Word (2 bytes) |
| +24 | Error Code (2 bytes) |
| +26 | Monitoring Channel (16 bytes) |
| +42 | Realtime Configuration (8 bytes) |

**Minimum Frame Size:** If the response frame is shorter than 64 bytes, the drive fills it with zeros until the length is 64 bytes.

---

## 5. LinUDP V2 Parameters

LinUDP servo drives have an additional parameter tree branch called **"LinUDP Intf"**, which can be configured with LinMot-Talk software.

**Available at:** http://www.linmot.com (Download → Software & Manuals)

### Parameters

#### Dis-/Enable
Turn the interface on or off.

#### Ethernet Configuration
Choose the connection type.

#### Monitoring Channels
Define 4 UPIDs whose values are transmitted cyclically in the response when the monitoring channel bit is active.

| Channel | Description | Parameter UPID |
|---------|-------------|----------------|
| Channel 1 UPID | Source UPID for Monitoring Channel 1 | 0x20A8 |
| Channel 2 UPID | Source UPID for Monitoring Channel 2 | 0x20A9 |
| Channel 3 UPID | Source UPID for Monitoring Channel 3 | 0x20AA |
| Channel 4 UPID | Source UPID for Monitoring Channel 4 | 0x20AB |

#### Master Configuration
For communication safety. Three options:

| Option | Description |
|--------|-------------|
| **No Filter** | The drive does no control (default) |
| **Single Master** | Drive accepts IP address from sender of first LinUDP V2 telegram, then only responds to telegrams from that address |
| **Single Master with fixed IP** | Drive only responds to telegrams from the fixed IP address defined in Master IP Address parameter |

#### Master IP Address
The fixed IP address when using "Single Master with fixed IP" mode.

#### LinUDP V1 Mode
Changes the LinUDP behavior to old LinUDP V1 behavior. If activated, refer to LinUDP (V1) manual Art-Nr. 0185-1083.

---

## 6. LinUDP V2 Modules

LinUDP V2 has:
- **3 modules** for master to drive communication
- **8 modules** for drive to master communication

### 6.1 Master to Drive Modules

#### Control Word
Access the main state machine of the drive.

- Refer to "User Manual Motion Control Software" for control word details
- Transfer with **Low Byte First** order

#### MC Cmd Interface
Maps the MC command interface of the drive.

- Refer to MC software documentation
- Each part transferred in **Low Byte First** byte order
- **Example:** VAI Go To Pos 10mm order:
  ```
  00 01 A0 86 01 00 40 42 0F 00 40 42 0F 00 40 42 0F 00
  ```

#### Real Time Configuration
Allows accessing:
- Parameters
- Variables
- Curves
- Error log
- Command table

**Capabilities:**
- Restart, start, and stop of the drive
- Works independently from MC command interface
- Changing a parameter and sending a motion command can be done in parallel
- Influences both telegram directions

**Transfer:** Low Byte First order

See [Chapter 7: Real Time Config Module](#7-real-time-config-module) for details.

### 6.2 Drive to Master Modules

#### Status Word
16-bit status word.
- Refer to "User Manual Motion Control Software" for bit meanings

#### State Var
Consists of **MainState** and **SubState**.
- Refer to "User Manual Motion Control Software", Chapter 3, State Var table
- Contains all relevant flags and information for clean handshaking within one word
- Can replace modules "Get MC Header Echo" and "Get Error Code"

#### Actual Position
Returns the actual position of the motor.
- **Format:** 32-bit integer value
- **Resolution:** 0.1 µm

#### Demand Position
Returns the demand position of the motor.
- **Format:** 32-bit integer value
- **Resolution:** 0.1 µm

#### Current
Returns the set current of the motor.
- **Format:** 16-bit integer value
- **Resolution:** 1 mA

#### Warn Word
Returns the warn word.
- Refer to "User Manual Motion Control Software"

#### Error Code
Returns the error code.
- Refer to "User Manual Motion Control Software" for error code meanings

#### Monitoring Channel
Transmits cyclically the value of the variable defined by the monitoring channel parameter.
- See [LinUDP V2 Parameters](#5-linudp-v2-parameters)

---

## 7. Real Time Config Module

The Real Time Config (RTC) module structure uses Data Output (DO) and Data Input (DI) from the **Master's point of view**.

### RTC Structure

| Word # | DO (Master → Drive) | DI (Drive → Master) |
|--------|---------------------|---------------------|
| 1 | Parameter Channel Control | Parameter Channel Status |
| 2 | Argument (depends on Cmd ID) | Argument (depends on Cmd ID) |
| 3 | Argument (depends on Cmd ID) | Argument (depends on Cmd ID) |
| 4 | Argument (depends on Cmd ID) | Argument (depends on Cmd ID) |

### Real Time Config Control

**Parameter Channel Control** (16 bits):

| Bits | Field | Description |
|------|-------|-------------|
| 15-8 | Parameter Command ID | Command to be executed (see Command ID table) |
| 7-4 | Reserved | - |
| 3-0 | Command Count | Counter for command execution |

### Real Time Config Status

**Parameter Channel Status** (16 bits):

| Bits | Field | Description |
|------|-------|-------------|
| 15-8 | Parameter Status | Status of command execution (see Parameter Status table) |
| 7-4 | Reserved | - |
| 3-0 | Command Count Response | Echo of command count |

### Command Count

A new command is only evaluated if the value of the command count changes.

**Simplest approach:** Toggle bit 0.

---

## Command IDs

### Parameter Access (0x10-0x17)

| Cmd ID | Description |
|--------|-------------|
| 0x00 | No Operation |
| 0x10 | Read ROM Value of Parameter by UPID |
| 0x11 | Read RAM Value of Parameter by UPID |
| 0x12 | Write ROM Value of Parameter by UPID |
| 0x13 | Write RAM Value of Parameter by UPID |
| 0x14 | Write RAM and ROM Value of Parameter by UPID |
| 0x15 | Get minimal Value of Parameter by UPID |
| 0x16 | Get maximal Value of Parameter by UPID |
| 0x17 | Get default Value of Parameter by UPID |

### Parameter (UPID) List (0x20-0x23)

| Cmd ID | Description |
|--------|-------------|
| 0x20 | Start Getting UPID List |
| 0x21 | Get next UPID List item |
| 0x22 | Start Getting Modified UPID List |
| 0x23 | Get next Modified UPID List item |

### Stop / Start / Default (0x30-0x36)

| Cmd ID | Description |
|--------|-------------|
| 0x30 | Restart Drive |
| 0x31 | Set parameter ROM values to default (OS SW) |
| 0x32 | Set parameter ROM values to default (MC SW) |
| 0x33 | Set parameter ROM values to default (Interface SW) |
| 0x34 | Set parameter ROM values to default (Application SW) |
| 0x35 | Stop MC and Application Software (for Flash access) |
| 0x36 | Start MC and Application Software |

### Curve Service (0x40-0x62)

| Cmd ID | Description |
|--------|-------------|
| 0x40 | Save all Curves from RAM to Flash |
| 0x41 | Delete all Curves (RAM) |
| 0x50 | Start Adding Curve (RAM) |
| 0x51 | Add Curve Info Block (RAM) |
| 0x52 | Add Curve Data (RAM) |
| 0x53 | Start Modifying Curve (RAM) |
| 0x54 | Modify Curve Info Block (RAM) |
| 0x55 | Modify Curve Data (RAM) |
| 0x60 | Start Getting Curve (RAM) |
| 0x61 | Get Curve Info Block (RAM) |
| 0x62 | Get Curve Data (RAM) |

### Error Log (0x70-0x74)

| Cmd ID | Description |
|--------|-------------|
| 0x70 | Get Error Log Entry Counter |
| 0x71 | Get Error Log Entry Error Code |
| 0x72 | Get Error Log Entry Time low |
| 0x73 | Get Error Log Entry Time high |
| 0x74 | Get Error Code Text Stringlet |

### Command Table (0x80-0x8E)

| Cmd ID | Description |
|--------|-------------|
| 0x80 | Command Table: Save to Flash |
| 0x81 | Command Table: Delete All Entries (RAM) |
| 0x82 | Command Table: Delete Entry |
| 0x83 | Command Table: Write Entry |
| 0x84 | Command Table: Write Entry Data |
| 0x85 | Command Table: Get Entry |
| 0x86 | Command Table: Get Entry Data |
| 0x87 | Get Presence List of Entries 0..31 from RAM |
| 0x88 | Get Presence List of Entries 32..63 from RAM |
| 0x89 | Get Presence List of Entries 64..95 from RAM |
| 0x8A | Get Presence List of Entries 96..127 from RAM |
| 0x8B | Get Presence List of Entries 128..159 from RAM |
| 0x8C | Get Presence List of Entries 160..191 from RAM |
| 0x8D | Get Presence List of Entries 192..223 from RAM |
| 0x8E | Get Presence List of Entries 224..255 from RAM |

---

## Parameter Status Codes

| Status | Description |
|--------|-------------|
| 0x00 | OK, done |
| 0x02 | Command Running / Busy |
| 0x04 | Block not finished (Curve Service) |
| 0x05 | Busy |
| 0xC0 | UPID Error |
| 0xC1 | Parameter Type Error |
| 0xC2 | Range Error |
| 0xC3 | Address Usage Error |
| 0xC5 | Error: Command 0x21 "Get next UPID List item" was executed without prior execution of "Start Getting UPID List" |
| 0xC6 | End of UPID List reached (no next UPID List item found) |
| 0xD0 | Odd Address |
| 0xD1 | Size Error (Curve Service) |
| 0xD4 | Curve already defined / Curve not present (Curve Service) |

---

## RTC Command Details

### Parameter Access Commands

**Word Layout:**

| Word # | DO | DI |
|--------|----|----|
| 1 | Parameter Channel Control | Parameter Channel Status |
| 2 | Parameter UPID | Parameter UPID |
| 3 | Parameter Value Low | Parameter Value Low |
| 4 | Parameter Value High | Parameter Value High |

### Curve Access Commands

**Word Layout:**

| Word # | DO | DI |
|--------|----|----|
| 1 | Parameter Channel Control | Parameter Channel Status |
| 2 | Curve Number | Curve Number |
| 3 | Data Value Low / Info Block size | Data Value Low / Info Block size |
| 4 | Data Value High / Data Block size | Data Value High / Info Block size |

### Start Getting UPID List (0x20)

**Request:**

| Word # | DO | DI |
|--------|----|----|
| 1 | Parameter Channel Control | Parameter Channel Status |
| 2 | Start UPID (search from this UPID) | - |
| 3 | - | - |
| 4 | - | - |

**Response (Get next with 0x21):**

| Word # | DO | DI |
|--------|----|----|
| 1 | Parameter Channel Control | Parameter Channel Status |
| 2 | - | UPID found |
| 3 | - | Address Usage |
| 4 | - | - |

### Start Getting Modified UPID List (0x22)

**Request:**

| Word # | DO | DI |
|--------|----|----|
| 1 | Parameter Channel Control | Parameter Channel Status |
| 2 | Start UPID (search from this UPID) | - |
| 3 | - | - |
| 4 | - | - |

**Response (Get next with 0x23):**

| Word # | DO | DI |
|--------|----|----|
| 1 | Parameter Channel Control | Parameter Channel Status |
| 2 | - | UPID found |
| 3 | - | Data Value Low |
| 4 | - | Data Value High |

### Get Error Log Entry Counter (0x70)

| Word # | DO | DI |
|--------|----|----|
| 1 | Parameter Channel Control | Parameter Channel Status |
| 2 | - | - |
| 3 | - | Number of Logged Errors |
| 4 | - | Number of Occurred Errors |

### Get Error Log Entry Error Code (0x71)

| Word # | DO | DI |
|--------|----|----|
| 1 | Parameter Channel Control | Parameter Channel Status |
| 2 | Entry Number (0..20) | Entry Number |
| 3 | - | Logged Error Code |
| 4 | - | - |

### Get Error Log Entry Time Low (0x72)

| Word # | DO | DI |
|--------|----|----|
| 1 | Parameter Channel Control | Parameter Channel Status |
| 2 | Entry Number (0..20) | Entry Number |
| 3 | - | Entry Time Low Word |
| 4 | - | Entry Time Mid Low Word |

### Get Error Log Entry Time High (0x73)

| Word # | DO | DI |
|--------|----|----|
| 1 | Parameter Channel Control | Parameter Channel Status |
| 2 | Entry Number (0..20) | Entry Number |
| 3 | - | Entry Time Mid High Word |
| 4 | - | Entry Time High Word |

**Note:** The Error Log Entry Time consists of 32-bit hours (Time High) and 32-bit milliseconds (Time Low).

### Get Error Code Text Stringlet (0x74)

| Word # | DO | DI |
|--------|----|----|
| 1 | Parameter Channel Control | Parameter Channel Status |
| 2 | Error Code | Error code |
| 3 | Stringlet Number (0..7) | Stringlet Byte 0 and 1 |
| 4 | - | Stringlet Byte 2 and 3 |

### Command Table: Save to Flash (0x80)

| Word # | DO | DI |
|--------|----|----|
| 1 | Parameter Channel Control | Parameter Channel Status |
| 2 | - | - |
| 3 | - | - |
| 4 | - | - |

⚠️ **Requirement:** The MC software must be stopped (with command 0x35: Stop MC and Application Software). The LinUDP V2 Interface will stay active while the MC software is stopped.

### Command Table: Delete All Entries (RAM) (0x81)

| Word # | DO | DI |
|--------|----|----|
| 1 | Parameter Channel Control | Parameter Channel Status |
| 2 | - | - |
| 3 | - | - |
| 4 | - | - |

### Command Table: Delete Entry (0x82)

| Word # | DO | DI |
|--------|----|----|
| 1 | Parameter Channel Control | Parameter Channel Status |
| 2 | Entry Number | Entry Number |
| 3 | - | - |
| 4 | - | - |

### Command Table: Write Entry (0x83)

| Word # | DO | DI |
|--------|----|----|
| 1 | Parameter Channel Control | Parameter Channel Status |
| 2 | Entry Number | Entry Number |
| 3 | Block Size (even number of bytes) | Block Size |
| 4 | - | - |

### Command Table: Write Entry Data (0x84)

| Word # | DO | DI |
|--------|----|----|
| 1 | Parameter Channel Control | Parameter Channel Status |
| 2 | Entry Number | Entry Number |
| 3 | Data | Data |
| 4 | Data | Data |

### Command Table: Get Entry (0x85)

| Word # | DO | DI |
|--------|----|----|
| 1 | Parameter Channel Control | Parameter Channel Status |
| 2 | Entry Number | Entry Number |
| 3 | - | Block Size |
| 4 | - | - |

### Command Table: Get Entry Data (0x86)

| Word # | DO | DI |
|--------|----|----|
| 1 | Parameter Channel Control | Parameter Channel Status |
| 2 | Entry Number | Entry Number |
| 3 | - | Data |
| 4 | - | Data |

### Command Table: Get Entry List (0x87..0x8E)

| Word # | DO | DI |
|--------|----|----|
| 1 | Parameter Channel Control | Parameter Channel Status |
| 2 | - | Offset in bytes |
| 3 | - | Bit field (Bit set = undefined / Bit cleared = used) |
| 4 | - | Bit field (Bit set = undefined / Bit cleared = used) |

---

## Contact & Support

### SWITZERLAND

**NTI AG**  
Bodenaeckerstrasse 2  
CH-8957 Spreitenbach

- **Sales and Administration:** +41 56 419 91 91 | office@linmot.com
- **Tech. Support:** +41 56 544 71 00 | support@linmot.com | http://www.linmot.com/support
- **Tech. Support (Skype):** support.linmot
- **Fax:** +41 56 419 91 92
- **Web:** http://www.linmot.com

### USA

**LinMot USA Inc.**  
N1922 State Road 120, Unit 1  
Lake Geneva, WI 53147  
USA

- **Phone:** 262-743-2555
- **E-Mail:** usasales@linmot.com
- **Web:** http://www.linmot-usa.com/

**Find local distribution:** http://www.linmot.com/contact

---

## Copyright Notice

© 2019 NTI AG

This work is protected by copyright. Under the copyright laws, this publication may not be reproduced or transmitted in any form, electronic or mechanical, including photocopying, recording, microfilm, storing in an information retrieval system, not even for didactical use, or translating, in whole or in part, without the prior written consent of NTI AG.

**LinMot®** is a registered trademark of NTI AG.

**Note:** The information in this documentation reflects the stage of development at the time of press and is therefore without obligation. NTI AG reserves the right to make changes at any time and without notice to reflect further technical advance or product improvement.

