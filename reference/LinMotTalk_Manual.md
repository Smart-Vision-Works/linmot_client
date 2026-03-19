# LinMot-Talk 6 Configuration Software Manual

*August 2023 — Doc.: 0185-1059-E_6V20_MA_LinMotTalk*

© 2023 NTI AG. This work is protected by copyright. Under the copyright laws, this publication may not be
reproduced or transmitted in any form, electronic or mechanical, including photocopying, recording,
microfilm, storing in an information retrieval system, not even for didactical use, or translating, in whole or
in part, without the prior written consent of NTI AG.

LinMot® is a registered trademark of NTI AG.

> The information in this documentation reflects the stage of development at the time of press and is
> therefore without obligation. NTI AG reserves itself the right to make changes at any time and without
> notice to reflect further technical advance or product improvement.

**NTI AG / LinMot**
Bodenaeckerstrasse 2, CH-8957 Spreitenbach
Tel.: +41 56 419 91 91 | Email: office@LinMot.com | www.LinMot.com

---

## 1 Introduction

The LinMot-Talk 6 software is a PC based tool which helps the user install firmware on the drive, set up the
drive's configuration, define and program motion profiles, emulate the PLC, watch variables, and read
messages and errors. LinMot-Talk 6 works with drive series A1100, B1100, C1100, E1100, C1200, E1200,
E1400, and B8050. It replaces the LinMot-Talk1100 software.

This manual covers LinMot-Talk version **6.10**. Features may differ in other versions.

### 1.1 System Generation (SG)

LinMot drive families are based on different hardware platforms called *system generations* (SG). The
following table shows which drive family belongs to which SG:

| SG  | Drives |
|-----|--------|
| SG1 | Families E400, E4000 V1 (not supported by LinMot-Talk 6) |
| SG2 | Families E400, E4000 V2 (not supported by LinMot-Talk 6) |
| SG3 | Family E1100 (GP, CO, DN, DP) (LC/HC/XC) |
| SG4 | Family B1100 (VF, PP, GP, ML) (LC/HC/XC) |
| SG5 | Family E1200 (GP, DP, DS, EC, IP, LU, PL, PN, SC, SE); Family E1400 (GP, DP, DS, EC, IP, PD, PL, PN, SC, SE) (0S/1S); Family B8000-ML (GP, EC, IP, PL, PN, SC) |
| SG6 | Family C1250 (CC, CM, DS, EC, IP, LU, PD, PL, PN, SC, SE) (0S/1S); Family E1400V2 (GP, DP, DS, EC, IP, LU, PD, PL, PN, SC, SE) (0S/1S) |
| SG7 | Family A1100; Family C1100 (GP, DS, EC, PD, PN, SE) (0S/1S) |

### 1.2 UPID (Unique Parameter ID)

All parameters have an assigned identification number called a **UPID** (Unique Parameter ID). All
parameters are accessed on the drive over this identification.

### 1.3 PnP (Plug and Play)

Drive families A1100, C1100, C1200, E1200, and E1400 support **Plug and Play** functionality. When a
motor is connected, it is automatically detected and parameters are set accordingly — the drive can then
control the motor without any further configuration. All PnP-capable components are marked on the type
label with "PnP".

> **Note:** All parameters set by the previous PnP motor that do not exist in the new motor will be set to
> default values before the new motor's parameters are loaded.

---

## 2 Overview

The most used functions after starting LinMot-Talk are **Install Firmware** and **Login to a drive**.

The main UI areas are:
- **Tool button bar** — configuration/setup tools, drive selection, shortcuts
- **Control Panel** — Control/Status Window, IO Panel, Motion Command Window, Monitoring Window
- **Menu** — additional functions and setup options

### 2.1 Tool Button Bar

The tool button bar is always present. Buttons from left:

| Button | Function |
|--------|----------|
| Show/Hide Tree | Shows or hides the project tree window |
| Up | Sets focus in project tree to parent of selection |
| Toggle | Toggles between the last two displayed tree branches |
| Import Configuration | Imports configurations to drives |
| Export Configuration | Exports configurations (parameters, variables, oscilloscope, curves) |
| Print | Prints items like curves, parameter configurations, etc. |
| Install Firmware | Starts firmware installation wizard |
| Open Login | Logs in to all drives in the selected workspace |
| Save Login | Saves the actual workspace |
| Reboot | Restarts the firmware on the drive |
| Stop | Stops firmware on the drive (for downloading/configuring) |
| Blink | Sends a blink LED command to the selected drive |
| Default | Defaults parameters by instance (OS, MC, INTF, APPL) |
| Go Offline | Logs out from active drive |
| Start Motor Wizard | Starts motor configuration wizard |
| Show Control Panel | Switches to control panel |
| Show Parameters | Switches to parameter view |
| Show Variables | Switches to variables view |
| Show Oscilloscope | Switches to oscilloscope |
| Show Messages | Switches to message viewer |
| Show Errors | Switches to error viewer |
| Show Curves | Switches to curve tool |
| Show Command Table | Switches to command table editor |
| Show Object Inspector | Displays help info for the selected object |
| +/- decimal | Show one more/fewer decimal place (only when "Round the decimal places" is active) |
| Information Window | Activates the Information Window |

### 2.2 Menu

#### 2.2.1 File

| Command | Description |
|---------|-------------|
| Login / Open Offline… | Opens the window to log in to a drive, or generate an offline drive from a config file |
| Create Offline… | Creates an offline drive with default configuration |
| Scanning (with CANusb) | Scans the CANusb board for a drive (requires CANusb board) |
| Scanning (via Ethernet) | Scans Ethernet and lists drives |
| Logout | Logs out from the active drive |
| Import… | Imports a configuration from a `.lmc` file |
| Export | Saves the configuration of drives |
| Save All | Fast save of entire configuration |
| Save Login | Saves the ports LinMot-Talk is currently logged into |
| Open Login | Opens a `.lws` file and tries to log in to all drives |
| Print | Prints parameters, variables, etc. |
| Install Firmware | Starts firmware installation |
| New | Opens a new LinMot-Talk window |
| Exit | Closes the active window (or shuts down LinMot-Talk if last window) |

#### 2.2.2 Search

| Command | Description |
|---------|-------------|
| Find with UPID… | Finds parameter or variable by UPID |
| Find with Caption… | Finds parameter or variable by caption |

#### 2.2.3 Drive

| Command | Description |
|---------|-------------|
| Reboot… | Restarts firmware on the drive |
| Stop… | Stops firmware on the drive |
| Blink | Sends a blink LED command to the active drive |
| Download → Software | Downloads individual software parts manually |
| Download → Configuration | Downloads a `.gr3` file to default parameters |
| Export Raw Data… | Exports raw configuration data to a `.pvl` file |
| Create UPID List… | (B1100 only) Generates a list of UPIDs and Master UPIDs |
| Set Access Code… | Enters access code to activate special features |
| Save Config To SD-Card | Saves configuration to SD Card (drives with SD Card slot only) |
| Motor Wizard… | Opens the Motor Wizard |
| Compare Parameter | Compares settings between different drives |
| Advanced Save | Saves configuration in a format compatible with LinMot-Talk 6.4 and older |
| Reluctance and Friction Table Wizard | Starts the Correction Table Wizard |
| Install Interface/Application | Changes the interface or application without defaulting other instances |

#### 2.2.4 Services

| Command | Description |
|---------|-------------|
| Show Control Panel | Opens Control Panel view |
| Show Parameters | Opens Parameter view |
| Show Variables | Opens Variables view |
| Show Oscilloscope | Opens Oscilloscope view |
| Show Messages | Opens Messages view |
| Show Errors | Opens Error view |
| Show Curves | Opens Curve view |
| Show Command Table | Opens Command Table view |
| Show FS Par Validation | Opens Functional Safety Parameter Validation (2S drives only) |

#### 2.2.5 Options

| Command | Description |
|---------|-------------|
| Language | Sets language (English, German, Italian) |
| UPID Display Mode | Sets UPID display to Hexadecimal or Decimal |
| Raw Data Display Mode | Sets raw data display to Hexadecimal or Decimal |
| Exit Warning | Controls warning when closing with active devices |
| Set Login Timeout | Sets communication timeout (values < 2 s affect login only; > 2 s affect all communication) |
| Save Debug Window Data | Enables/disables saving communication in a ring buffer |
| Round the decimal places | Rounds numerical values to a defined number of decimal places |
| Group Moduleparts | Groups drives connected to a motor module in the tree view |
| Motor Data and Firmware Directory | Changes the storage directory for motor and firmware files (requires restart) |
| Multi Read Option | When disabled, reads each variable/parameter individually (much slower) |

#### 2.2.6 Window

Lists all open LinMot-Talk windows. Click an item to make that window active.

#### 2.2.7 Tools

| Command | Description |
|---------|-------------|
| LinRS Test Tool | Opens a window for testing LinRS communication |
| CANTalk Manager | Opens the CANTalk Manager for CAN settings |
| RSTalk Debug Window | Opens the communication debug window |
| Read Drive | Reads correction tables from the drive |
| Process Monitoring | Configures process monitoring (requires Process Monitoring application) |

#### 2.2.8 Manuals

| Command | Description |
|---------|-------------|
| Parameter and Variables | Generates HTML page with all parameters/variables for the selected part |
| Errors | Generates HTML page with all errors for the selected part |
| Motion Commands | Generates HTML page with all motion commands for the installed MC software |
| Relevant Documents | Lists relevant PDF documents for the active drive |
| All Documents | Lists all PDF documents supplied with LinMot-Talk |

#### 2.2.9 Help

| Command | Description |
|---------|-------------|
| Home Page | Opens www.linmot.com |
| Update Functions | See subsection below |
| Default LinMot-Talk Settings | Resets all LinMot-Talk settings to defaults |
| About LinMot-Talk 6.10 | Shows version and build information |

##### Update Functions

| Command | Description |
|---------|-------------|
| Check for Updates | Checks the home page for a newer version of LinMot-Talk |
| Check for Software Updates at Program Startup | Runs update check at every startup |
| Download and show News | Downloads and shows LinMot news |
| Check for new Motor Files | Checks server/local path for new motor files |
| Download Older Releases | Downloads older firmware releases |
| Generate Portable App | Copies LinMot-Talk files into a portable folder |
| Update Option | Switches update source between LinMot server and a local path |

### 2.3 Control Panel

The Control Panel provides direct access to the control and status word of the MC Software, allowing the
PC to command the drive without a PLC — useful for first commissioning.

- **Control Word** — Directly write the MC software's control word. Enable **Manual Override**, then set
  each flag with the **Override Value** checkbox. Additional flags can be set via the override mask at
  `\Parameters\Motion Control SW\State Machine Setup\Control Word\Ctrl, Word Parameter Force Mask`.
- **Status Word** — Shows the actual drive MC software status word; updated automatically.
- **General Monitoring** — Displays actual motor and drive information.
- **Additional Variables** — Choose variables to display and auto-update.
- **IO Panel** — Control the X4 IOs (E1100) or X14 IOs (B1100) for commissioning.
- **Motion Command Interface** — Directly access the MC software's motion command interface. When
  **Enable Manual Override** is set, MC commands can be selected, parametrized, and sent.

### 2.4 Messages

Reads out and displays all messages logged on the drive in chronological order. Not available on B1100
series (does not support message logging). Select a message and press **F1** for details in the Object
Inspector.

### 2.5 Errors

Reads out and displays all errors logged on the drive in chronological order. Firmware installations are also
logged. Select an error and press **F1** for details. Generate a complete error list via **Manuals → Errors**.

### 2.6 Oscilloscope

The drive's built-in oscilloscope records up to **eight channels** in real time.

During login, oscilloscope settings and data are read from the drive. If a shot is running or ready to read
out, a "Read out" item appears; otherwise a default item is generated.

**Oscilloscope controls:**

| Button | Function |
|--------|----------|
| Start/Abort | Starts or aborts an oscilloscope shot |
| Fit View | Scales all channels to fit the window |
| Fit View (same unit same fit) | Channels with the same unit share the same scale |
| Save Display | Stores current zoom/scale/offset settings |
| Recall Display | Restores previously saved display settings |
| Export Data | Exports last recorded shot to a CSV file |
| Oscilloscope Settings | Opens channel/trigger/time/mode setup |
| Display Settings | Sets scale, offset, and color per channel |
| Show/Hide | Shows or hides individual channels |
| Show/Hide Cursor | Displays two time cursors for measuring |
| Statistics Value | Shows statistics for each channel (between cursors, or whole shot if cursors disabled) |

### 2.7 Curves

The curve tool creates, joins, uploads, downloads, and saves motor motion profiles.

> **Note:** On B1100, the curve feature must be enabled with an access key.

**Edit window buttons:**

| Button | Function |
|--------|----------|
| New Curve | Starts the curve wizard |
| Edit Properties | Modifies properties of a selected curve (name, time, stroke) |
| Edit Curve Values | Manually edits curve points |
| Join Curves | Joins all selected curves with a wizard |

**Download window buttons:**

| Button | Function |
|--------|----------|
| Upload Curves from Drive | Uploads all curves stored on the drive |
| Download Curves to Drive | Synchronizes the drive's curve sector with the download window |
| Auto Numerate Curves | Sets curve IDs automatically (must be unique) |

**Maximum curve storage:**

| Series | Max Curves | Storage Limit |
|--------|-----------|---------------|
| B1100 | 16 | `#Curves × 70 B + #SamplePoints × 4 B ≤ 2016 B` |
| A1100/C1100 | 50 | `#Curves × 70 B + #SamplePoints × 4 B ≤ 32512 B` |
| All other series | 100 | `#Curves × 70 B + #SamplePoints × 4 B ≤ 65280 B` |

### 2.8 Parameters

Drive parameters are displayed in a tree view.

**Parameter controls:**

| Button | Function |
|--------|----------|
| Show/Hide Details | Shows/hides UPID, scaling, min/max per parameter |
| Show UPID Browser | Opens UPID browser (visible when editing UPID-type parameters) |
| OK / Enter | Confirms the input value |
| Cancel | Cancels the typed value |
| Read | Reads and refreshes all parameters from the drive |

Parameters marked with a red **L** icon are **live parameters** — they can be changed without stopping
firmware. All other parameters require firmware to be stopped first.

**Parameter table columns:**

| Column | Description |
|--------|-------------|
| Name | Name of the parameter |
| Value | ROM value — written to RAM on restart; editable in LinMot-Talk |
| Raw Data | Data as stored in ROM (no scaling/offset), in hexadecimal |
| Value (RAM) | Active value in RAM; not editable in LinMot-Talk (live parameters update immediately) |
| UPID | Unique Parameter IDentification |
| Type | Data type (e.g., SInt32, UInt32, String) |
| Scale | Scale factor from raw data to value |
| Offset | Offset added to raw data to produce value |
| Min | Minimum allowed value |
| Max | Maximum allowed value |
| Default | Value after defaulting the drive |
| Attr. | Access rights: R (read), W (write), RW (both) |

Select a parameter and press **F1** to open the Object Inspector with additional documentation. A blue
"more" link opens the full documentation for that parameter.

### 2.9 Variables

Drive variables are arranged in functional groups. The **MC SW overview** group contains the most
commonly used variables.

**Variable controls:**

| Button | Function |
|--------|----------|
| Show/Hide Details | Shows/hides UPID, scaling, min/max |
| Read Variable | Reads the selected variable once |
| Write Variable | Writes the selected variable to the drive |
| Read All Variables | Reads all variables in the section once |
| Read All Variables Cyclically | Reads all variables in the section cyclically |
| Remove (Del) | Removes the selected variable from the list |
| Edit Properties | Displays and changes parameter properties |
| New … Variable | Adds a new variable (dropdown with types) |
| New Bit Variable | Adds a bit-type variable |
| New String Variable | Adds a string-type variable |
| New Float32 | Adds a float32 variable |
| New With UPID | Adds a variable by UPID |

Under **User Defined**, any variables or parameters can be grouped together. Variables can be selected via
UPID or dragged and dropped from the parameter or variable section.

### 2.10 Command Table

The command table (CT) stores up to **255 motion commands** (31 for B1100GP/VF; not supported on
B1100PP). Commands include motion commands, conditions, sequence directives, and parameter access.
CT entries can be executed via digital inputs (X6) or via interface software.

**CT editing elements:**

| Field | Description |
|-------|-------------|
| Entry ID | The CT entry being edited |
| Entry Name | Descriptive string, max 16 characters |
| Motion Command Category | Groups commands for better overview |
| Motion Command Type | Specifies the command to execute |
| Auto execute new command on next cycle | On the next cycle, executes the entry specified in "ID of Sequenced Entry" |
| ID of Sequenced Entry | CT entry to execute on next cycle when auto-execute is active |
| Apply | Writes edited values into the entry |
| Upload from Drive | Reads the entire command table from the drive |
| Download to Drive | Writes the edited table to the drive |

### 2.11 Access Codes

Special features and customer-specific applications are protected by software keys, enabled by
drive-specific access codes (pinned to the serial number). Navigate to **Drive → Set Access Code**.

- Up to **four keys** can be set on the drive.
- **Active Keys** lists all valid installed keys (key value and access code).
- To add a key: select the key name, define the value and access code, click **Write**, then click
  **Activate** (the drive will reboot).
- Access codes are **drive-specific** — they cannot be copied between drives.

**Evaluation mode** (LinMot-Talk 6.10+): Technology functions can be enabled for a **4-hour trial**. After
4 hours, the drive enters an error state. Check remaining time at `Variables\OS SW Keys`. The evaluation
access code appears as `FFFFFFFFh` in the Active Keys list.

**Technology function availability:**

| Function | Curve | Force Control | Process Monitoring |
|----------|-------|---------------|--------------------|
| E1100 | Standard | Technology Function | N/A |
| B1100 | Technology Function | Technology Function | N/A |
| B1150ML | Technology Function | Technology Function | N/A |
| E1200 | Standard | Technology Function | N/A |
| E1400V2 | Standard | Technology Function | N/A |
| A1100 | Standard | Technology Function | N/A |
| C1100 | Standard | Technology Function | N/A |
| C1200 | Standard | TF (Eval Mode Supported) | TF (Eval Mode Supported) |
| C1400 | Standard | Technology Function | N/A |
| D1400 | Standard | N/A | N/A |

### 2.12 Information Window

The Information Window is only visible when it has a message and is activated (via the Information Window
button in the tool button bar).

| Source | Message | When Cleared |
|--------|---------|-------------|
| Motor Wizard | Motor was not configured with Motor Wizard; only PnP was used | When user completes the Motor Wizard |
| Oscilloscope | An oscilloscope has new data | When user navigates to the oscilloscope |

---

## 3 Quick Start Guide

### 3.1 Cabling E1100

| Connector | Description |
|-----------|-------------|
| X1 | Motor Supply: 48–72 VDC (PWR+ to PGND) |
| X2 | Motor Phases (if present; otherwise connect motor to X3 only) |
| X3 | Motor Signals: DSUB-9 direct or via adapter |
| X4 | Commission with PC: wire Pin1 (GND), Pin2 (+24 VDC), and if present Pin12 SVE (+24 VDC) |
| X5 | **RS232**: DSUB-9 F/F 1:1 cable (X-modem). No COM port? Use USB-to-RS232 converter (LinMot art. 0150-2473) |

### 3.2 Cabling E1200

| Connector | Description |
|-----------|-------------|
| X1 | Motor Supply: 48–72 VDC (PWR+ to PGND) |
| X2 | Motor Phases |
| X3 | Motor Signals (motor phases are **not** on this connector — always wire phases to X2) |
| X4 | Commission with PC: wire Pin1 (GND), Pin2 (+24 VDC), and if present Pin12 SVE (+24 VDC) |
| X15/X16 | **Ethernet**: standard RJ45 patch cable to LAN |
| X19 | **RS232**: use RS232 PC configuration cable (LinMot art. 0150-2143). No COM port? Use USB-to-RS232 converter (art. 0150-2473) |

### 3.3 Cabling E1400

| Connector | Description |
|-----------|-------------|
| X2 | Motor Phases |
| X3 | Motor Encoder Signals |
| X4 | Commission with PC: wire Pin1 (GND) and Pin2 (+24 VDC) |
| X15/X16 | **Ethernet**: standard RJ45 patch cable to LAN |
| X19 | **RS232**: RS232 PC config cable (art. 0150-2143). No COM port? Use USB-to-RS232 (art. 0150-2473) |
| X30 | Motor Supply: 3×400 / 3×480 VAC 50/60 Hz |
| X33 | Safety Relays: separate +24 VDC supply; wire both Ksr+ (X33.4, X33.8) to +24 VDC and both Ksr− (X33.3, X33.7) to GND |

### 3.4 Cabling B1100

| Connector | Description |
|-----------|-------------|
| X1 | Motor Supply: 48–72 VDC (PWR+ to PGND) |
| X2 | Motor Phases |
| X3 | Motor Signals: DSUB-9 direct or via adapter |
| X5 | **RS232**: DSUB-9 F/F 1:1 cable. No COM port? Use USB-to-RS232 (art. 0150-2473) |
| X14 | Commission with PC: wire Pin13 (DGND) and Pin25 (+24 VDC) |

### 3.5 Cabling B8050-ML

| Connector | Description |
|-----------|-------------|
| X23 | **RS232**: DSUB-9 F/F 1:1 cable. No COM port? Use USB-to-RS232 (art. 0150-2473) |
| X24 | 24 V switched power supply |

### 3.6 Cabling A1100

| Connector | Description |
|-----------|-------------|
| X2 | Motor Phases |
| X3 | Motor Signals |
| X19 | **RS232**: RS232 PC config cable (art. 0150-3544). No COM port? Use USB-to-RS232 (art. 0150-2473) |
| X40 | Signal supply: Pin1 (GND) and Pin2 (+24 VDC). Motor supply: 48–72 VDC for PWR+ (Pin4) and PGND (Pin3). (LinMot art. 0150-3545 provides a connector with 1.5 m pre-crimped wires.) |

### 3.7 Cabling C1100

| Connector | Description |
|-----------|-------------|
| X1 | Motor Supply: 48–72 VDC (PWR+ to PGND) |
| X2 | Motor Phases |
| X3 | Motor Signals (motor phases are **not** on this connector — always wire phases to X2) |
| X4 | Commission with PC: wire Pin1 (GND) and Pin2 (+24 VDC) |
| X7-8 | **RS485**: USB-to-RS485 converter (art. 0150-3356). Switch S4.1 must be set to **ON** (FW ≥ 6.9). Only point-to-point connections supported. |
| X19 | **RS232**: RS232 PC config cable (art. 0150-2143). No COM port? USB-to-RS232 (art. 0150-2473). Switch S4.1 must be set to **OFF**. |
| X33 | Safety Relays (1S only): separate +24 VDC; Ksr+ (X33.4, X33.8) to +24 VDC, Ksr− (X33.3, X33.7) to GND |
| S4 | S4.1 selects communication channel: OFF (default) = RS232 on X19; ON = RS485 on X7/8. LinRS interface automatically uses the other channel. (FW ≥ 6.9) |

### 3.8 Cabling C1200

| Connector | Description |
|-----------|-------------|
| X1 | Motor Supply: 48–72 VDC (PWR+ to PGND) |
| X2 | Motor Phases |
| X3 | Motor Signals (motor phases are **not** on this connector — always wire phases to X2) |
| X4 | Commission with PC: wire Pin1 (GND) and Pin2 (+24 VDC) |
| X19 | **RS232**: RS232 PC config cable (art. 0150-2143). No COM port? USB-to-RS232 (art. 0150-2473) |
| X33 | Safety Relays (1S only): separate +24 VDC; Ksr+ (X33.4, X33.8) to +24 VDC, Ksr− (X33.3, X33.7) to GND |

### 3.9 Cabling M8000

| Connector | Description |
|-----------|-------------|
| X3 | Motor: single connector for both phases and signals |
| X19 | **RS232**: RS232 PC config cable (art. 0150-2143). No COM port? USB-to-RS232 (art. 0150-2473) |
| X33 | Safety Relays (1S only): separate +24 VDC; Ksr+ (X33.4, X33.8) to +24 VDC, Ksr− (X33.3, X33.7) to GND |
| X34 | Motor Supply: 48–72 VDC (PWR+ to PGND). Axes 1–4 and 5–8 are supplied separately. |
| X36 | Commission with PC: wire Pin1 (GND) and Pin2 (+24 VDC) |

### 3.10 Firmware Download

With cabling complete, power on the drive and start LinMot-Talk. Before first use, firmware must be
installed:

1. Click the **Install Firmware** button.
2. Choose the firmware file (e.g., `Firmware_Build20101126.sct`) and press **Open**.
3. The wizard guides through installation.

> **Note — Ethernet firmware installation:** A service password is required (for safety, to prevent
> accidentally flashing the wrong drive). By default no password is set. If the password is unknown,
> factory-reset parameters using the hex switches (see [Section 4.1](#41-setting-all-parameters-to-default-values)).
> Ethernet firmware install is only supported on drives with a separate Config Ethernet port (E1200, E1400).
> Installation over RT Ethernet is **not** possible.

Select the appropriate interface and application software for your drive type during the wizard.

### 3.11 Install only Interface or Application

**Drive → Install Interface/Application** installs only a different interface or application, in most cases
without defaulting the Motion Control Software and OS parameters.

- The first combo box selects the interface to install (current interface shown by default; unchanged if not
  modified).
- The second combo box selects the application to install (current application shown by default).
- Warnings in red rectangles indicate which parameters will be defaulted.

After clicking **Continue**, LinMot-Talk logs out, installs the selected firmware parts, then logs back in.

### 3.12 Login

After downloading firmware:

1. Go to **File → Login…** or double-click **Project** in the project tree.
2. Select the appropriate port and press **OK**.
3. A login info window shows login progress.

> If the firmware on the drive has a different version than LinMot-Talk, you can download compatible files
> to complete login — see [Section 3.13](#313-downloading-older-releases).

The Object Inspector window can be dragged away or closed; reopen with **F1**.

### 3.13 Downloading Older Releases

LinMot-Talk supports the current version and one version prior (e.g., LinMot-Talk 6.9 can log into drives
with firmware 6.9 or 6.8). To work with older versions, download the files first.

#### 3.13.1 Manual Downloading

1. Go to **Help → Update Functions → Download Older Releases**.
2. The window lists all downloadable versions. Grey = already installed; black = not yet installed.
3. Select the needed version and click **OK**.
4. LinMot-Talk downloads from the server (internet connection required).

Downloaded files are stored at:
```
C:\Users\USERNAME\AppData\Local\LinMot\LinMot-Talk X.Z - Build YYYYMMDD\Firmware\OlderReleases
```

#### 3.13.2 Downloading by Login

When logging into a drive whose firmware version is not yet on the computer, a dialog appears. Click
**Check for compatible version** to have LinMot-Talk find and download the correct files from the server.
After successful download, click **OK** and LinMot-Talk will offer to restart the login process.

### 3.14 Scanning CAN Bus

To automatically discover all drives on a CAN bus:

1. Go to **File → Scanning (with CANusb)** (requires a CANusb board).
2. A list of all present drives appears.
3. Click once to log in to all drives simultaneously.

### 3.15 Scanning Ethernet

To automatically discover drives connected via Ethernet:

1. Go to **File → Scanning (via Ethernet)**.
2. Select the network interface.
3. Optionally activate **Group Number** scanning to filter by a specific group number (parameter **Net Group**, UPID `0078h`).
4. A list of present drives appears with color-coded status:
   - **Green** — drive is ready to log in
   - **Grey** — already logged in to this drive
   - **Red** — another instance (other user or interface) is logged in

> **IP Address Assignment:** The default mode is DHCP. If no DHCP server responds, the drive uses
> IPv4 Link-Local (APIPA) addressing in the range `169.254.0.1` – `169.254.255.254`
> (subnet mask `255.255.0.0`). This process can take up to **one minute**.

### 3.16 Motor Wizard

Press the Motor Wizard button to open the wizard. Select **LinMot Linear Motors**, choose the stator
family (e.g., `PS0x-23x`) and stator subfamily (e.g., `PS01-23x160x`), select the actuator type, and press
**Open**.

If the motor has PnP functionality, LinMot-Talk opens the correct file and shows the Actuator Selection
page automatically. If the correct file cannot be found, update motor files first (see [Section 3.19](#319-update-motor-files)).

#### 3.16.1 Actuator Selection

Define the stator and slider. Derived settings show the complete motor type, article numbers, and key
technical data. Changing the positive moving direction is supported since release 6R7 and only for motors
with PnP version V3S2 or higher (V3S1 motors do not support this change).

#### 3.16.2 Drive Settings

Choose a drive name and, if applicable, a regeneration resistor.

#### 3.16.3 Extension Cable Setup

Define up to two cable segments for longer extension cables, which affect the motor's phase resistance.

#### 3.16.4 External Position Sensor System

Define an external position sensor if present:
- **E1100**: none, incremental AB(Z), or analog sine/cosine 1Vpp
- **B1100**: none, incremental AB(Z), or AB encoder simulation

#### 3.16.5 Force and Torque Sensor

Configure force or torque sensors if supported by the motor file.

| Parameter | Description |
|-----------|-------------|
| Installed Sensor | Usually set by PnP or motor file; can be overridden here |
| Speed Filter Time | Delta-t length for force speed calculation (derivative of measured force) |
| Acceleration Filter Time | Delta-t length for actual acceleration calculation |
| Speed Limit | Maximum velocity of speed limiter (0 = disabled); recommended when using force control |
| 0V/−10V Force | Force equivalent to the 0V/−10V analog input value (X4.10/X4.11) — SG6/SG7 only |
| 10V Force | Force equivalent to the 10V analog input value (X4.10/X4.11) — SG6/SG7 only |
| Minimum Force | Minimum force sensor measuring range value — SG8 only |
| Maximum Force | Maximum force sensor measuring range value — SG8 only |
| Minimum Torque | Minimum torque sensor measuring range value — SG8 only |
| Maximum Torque | Maximum torque sensor measuring range value — SG8 only |
| Sensor Direction | Selects the direction that produces positive force/torque |
| FIR Filter Mode | Sets the low-pass filter cutoff frequency on the sensor (only shown when sensor PnP works) |

#### 3.16.6 Feed Forward Parameters

Set feed forward parameters based on moving mass, additional load mass, friction, and orientation. Check
the **Derived Settings** to see the influence of your inputs.

#### 3.16.7 PID Position Controller

Configure the position controller parameters:

| Setting | Description |
|---------|-------------|
| Soft settings | Low PID values; quieter motor, less stiff position control |
| Stiff settings | Higher PID values; more noise and power, tighter position control |
| I Gain | Set to 0 by default (steady-state deviation may occur; enabling I Gain can cause oscillation) |
| Noise Filter | Reduces noise from position feedback sensor at standstill |

> **Tip:** Start with soft settings. For iterative PID tuning, refer to document 0185-1156 *Loop Tuning*.

#### 3.16.8 Homing 1

Define the homing procedure. The most common mode is **Mechanical Stop Negative Search** — the
slider moves toward the stator's front end (the end without the cable). Other modes support home
switches, limit switches, indexer inputs, or combinations.

#### 3.16.9 Homing 2

Define the **slider home position** — the most important value, which determines where the slider is
positioned relative to the stator at the home position, and thus how far the motor can move in each
direction.

#### 3.16.10 Homing 3

Define the user's coordinate system. Press **Finish** to write all parameters to the drive.

The Motor Wizard can be re-run at any time (e.g., to set up an external sensor, change load setup, or
change motor type). On re-run, a list of parameters that will change is shown before applying.

### 3.17 Correction Table Wizard

The Correction Table Wizard configures up to two correction tables stored in the drive:
- **Reluctance correction table**
- **Friction correction table**

Tables are measured by a bidirectional measuring movement. They correct disturbances caused by
reluctance forces and friction.

> **Warning:** The motor will physically move during the measuring movement. Ensure no obstacles
> exist between the configured minimal and maximal positions.

#### 3.17.1 Info Page

Contains information only — no configurable options. The key message is that the motor will move during
the wizard.

#### 3.17.2 Settings for the Measuring Movement

| Parameter | Default | Description |
|-----------|---------|-------------|
| Minimal Position | From Motor Wizard | Motor moves to this position first |
| Maximal Position | From Motor Wizard | Motor moves to this position next |
| Acceleration | 1 m/s² | Too fast = inaccurate measurements |
| Maximal Velocity | 0.01 m/s | Too fast = inaccurate measurements |
| Linearizing Mode (UPID 173Eh) | — | Must be set to **System Compensation** for correction tables to be used |

The measuring movement sequence: minimal position → maximal position → minimal position.

#### 3.17.3 Define the Usage of the Reluctance Table

| Option | Description |
|--------|-------------|
| Enable Usage Reluctance Table correction | Whether the table is used when available |
| Enable the Measurement of the Reluctance Table | Whether the table is measured during the movement |
| Load Reluctance Table after reboot | **Clear Table**: deleted on reboot; **Keep Table**: retained; **Load From Motor, clear if not found**: reads from motor, clears if absent; **Load From Motor, keep if not found**: reads from motor, keeps old table if absent |
| Table save in | **Save in RAM**: lost on reboot; **Save in ROM**: persists across reboots |

#### 3.17.4 Define the Usage of the Friction Table

**Friction Table Compensation Mode:**

| Mode | Description |
|------|-------------|
| Off | No friction table compensation |
| Faded In With Actual Velocity Filtered | Compensation uses actual filtered velocity; value multiplied by factor F (velocity-dependent) |
| Demand Velocity | Compensation uses demand velocity; **recommended default** |

| Option | Description |
|--------|-------------|
| Friction I Decrement Mode | When axis is at standstill and position deviation is below `1397h` (Friction I Decrement Start Limit), the integrator is decremented to `1398h` (Friction I Decrement Stop Limit). Noise deadband is automatically deactivated when enabled. |
| Enable Usage of Friction Table correction | Whether the friction table is used when available |
| Enable the Measurement of the Friction Table | Whether the friction table is measured during the movement |
| Load Friction Table after reboot | Same options as reluctance table |
| Table save in | Same options as reluctance table |

#### 3.17.5 General Table Settings

| Option | Description |
|--------|-------------|
| Compensation Type | **Linear** for linear motors; **Rotative** for rotary motors |
| Control Parameter Set Selection | Selects the control parameter set used during the measuring movement |

#### 3.17.6 Write Parameters and Execute the Measuring Movement

Button sequence:

1. **Stop Firmware and Write Parameter** — stops MC software and writes linearising parameters
2. **Start Firmware** — starts MC software (circulating dots appear while waiting)
3. **Make Measuring Move** — performs homing (if needed) then the measuring movement
4. **Stop** — emergency stop button
5. **Next** — advances to results (only available after measuring movement completes)

#### 3.17.7 Result Page

Shows the reluctance table (upper graph) and friction table (lower graph), both auto-scaled. Press
**Finish** to end the wizard and activate the correction tables.

### 3.18 Unit System

For LinMot rotary motors and the rotary part of PR01 motors, the Motor Wizard shows a unit system page.
The unit system is **active in LinMot-Talk only** — it has no effect on the drive.

| System | Position Unit | Notes |
|--------|--------------|-------|
| Linear | mm | "1 Revolution" parameter = mm per revolution |
| Rotary | ° | "1 Revolution" parameter = encoder ticks per revolution |

> **Important:**
> - For LinMot rotary motors (EC02): the "1 Revolution" value must be **divisible by 4**.
> - For PR01 rotary part: the "1 Revolution" value must be **divisible by 12**.
>
> Using a multiple of 360 gives well-rounded numbers. Alternatively, a value where position wraps back
> to 0 after each revolution can also be used.

### 3.19 Update Motor Files

1. Go to **Help → Update Functions → Check for new Motor Files**.
2. Click **Check for Updates** — LinMot-Talk finds updates since the last update date.
3. Use **Reset Date** to show all updates ever available.
4. Select the files to update and click **OK**.

Each motor file shows its creation date. Selecting **All** downloads every motor file from the server (usually
only possible after resetting the update date).

#### 3.19.1 Adding a Motor Repository

Enable the **Advanced** checkbox in the motor files update window, then click **Add Repository**.

| Field | Description |
|-------|-------------|
| Repository Name | Display name for this repository in the update window |
| URL | Internet address of the repository source |
| Username | Login username for the repository |
| Password | Login password for the repository |

> **Note:** NTI AG does not validate third-party repositories — the customer is responsible for their
> content.

#### 3.19.2 Problems with the Update Function

| Case | Solution |
|------|----------|
| No internet connection | Connect the computer to the internet |
| Firewall blocks TCP port 443 | LinMot-Talk uses HTTPS — allow outbound port 443 |
| Update server offline | Try again later |

**Manual download alternative:** Visit https://repo001.linmot.com/svn to download files as a GNU tarball.
Click the needed folder, then download the tarball.

**Local package alternative:**
1. Download the Motor and OlderReleases packages from the LinMot homepage (under LinMot-Talk
   downloads).
2. Unzip to a local folder (packages must be in the same location).
3. In LinMot-Talk, change **Update Option** from **Server** to **Local**.
4. LinMot-Talk will prompt for the package path.

> **Note:** The OlderReleases package is over **1 GB**. Local packages may be older than the server.
> When connecting to the server fails, a **"Change to local Package"** button appears in the error dialog.

### 3.20 Continuous Curve Mode

To run a curve cyclically:

1. In the parameter tree, navigate to:
   `\Motion Control SW\Motion Interface\Run Mode Settings\RunMode Selection\`
   and set it to **Continuous Curve**.
2. Navigate to:
   `\Motion Control SW\Motion Interface\Time Curve Settings\`
   and set **Curve ID** to the ID of the curve to run (e.g., `1`).

> **Note:** On B1100 drives, the curve feature must be enabled with an access key before use.

### 3.21 Defining Curves

Example: define two sine curves (50 mm out and 50 mm in) and join them.

1. Open the curve tool (**Show Curves** button).
2. Click **New Curve** to start the curve wizard.
3. Keep **position vs. time** mode (default) and press **Next**.
4. Set **Curve ID** to `2`, name it `SineOut`, set End Point to `50 mm`. Press **Next**, accept the
   proposed sample points, press **Finish**.
5. Start the wizard again: set **Curve ID** to `3`, name it `SineIn`, **Curve Length** `500 ms`,
   **Start Point** `50 mm`, **End Point** `0 mm`. Click **Next** twice, then **Finish**.
6. Select both curves, then press **Join Curves**.
7. In the join wizard, set **Curve Name** to `SineOutIn`, **Curve ID** to `1`. Press **Next** and
   **Finish**.
8. Move all curves from the edit window to the download window.
9. Press **Download Curves into Drive**, confirm the warning, and wait for the progress window to finish.

### 3.22 Control Status

After configuring parameters and curves, start the motor:

1. Switch to the control panel.
2. Press **Start** (starts the drive's firmware) and wait for the control status panel to update.
3. Enable the **Switch On** and **Home** flags.
4. Turn **Switch On** off, then on again (auto-start prevention).
5. The motor powers up and holds its current position.
6. Set the **Home** flag — the motor initializes against the inner hard stop.
7. Clear the **Home** flag — the motor begins running the curve continuously.

Detailed state diagram information is in the MC Software manual.

### 3.23 Oscilloscope

The LinMot-Talk default oscilloscope samples: actual position, demand position, position difference, and
demand current.

1. Click **Show Oscilloscope**.
2. Click **Start** — recorded data is read from the drive and displayed.
3. Press **Fit View** if channels need rescaling.

#### 3.23.1 Oscilloscope Settings

The settings window has three tabs: General, Trigger, and Advanced.

##### General

| Setting | Description |
|---------|-------------|
| Acquisition Mode | Single shot or continuous recording |
| Recording Time | Duration of one oscilloscope shot |
| Channel X checkbox | Activates/deactivates the channel |
| Is math channel | Enables mathematical operations on other channels |
| Group | Variable group for the channel |
| Variable | Variable recorded in the channel |

**Math channel functions:**

| Function | Description |
|----------|-------------|
| Addition | Sum of two selected channels |
| Subtraction | Difference between two selected channels |
| Product | Product of two selected channels |
| Ghost | Shows a channel from another oscilloscope (select the oscilloscope and channel) |

##### Trigger

Two trigger conditions (A and B) can be defined and combined with AND/OR logic.

| Event | Description |
|-------|-------------|
| Rising edge | Triggers when variable goes from below value to at/above value |
| Falling edge | Triggers when variable goes from above value to at/below value |
| Any edge | Triggers when variable passes through or reaches value |
| Greater than | Triggers when variable > value |
| Less than | Triggers when variable < value |
| Greater or equal | Triggers when variable ≥ value |
| Less or equal | Triggers when variable ≤ value |
| Equal | Triggers when variable = value |
| Not Equal | Triggers when variable ≠ value |
| Change | Triggers on any change (independent of defined value) |
| Difference greater or equal | Triggers if elevation between neighboring points ≥ value/ms |
| Difference less or equal | Triggers if elevation between neighboring points ≤ value/ms |
| ABS difference greater or equal | Triggers if absolute elevation between points ≥ value/ms |
| ABS difference smaller or equal | Triggers if absolute elevation between points ≤ value/ms |
| Masked Bits = False | Triggers if all set bits in value are false in variable |
| Masked Bits = True | Triggers if all set bits in value are true in variable |

##### Advanced

| Setting | Description |
|---------|-------------|
| Pretrigger | Defined in % of recording time |
| Delay | Delay after trigger event before recording starts (absolute time) |
| Sample period | Time between two neighboring measurement points |
| Number of samples | Measurement points per channel (maximum shown in the UI based on current settings) |
| Preview function | Draws an estimated graph during measurement (only works when recording time > 10 s) |

#### 3.23.2 Display Settings

**Scaling methods:**

1. **Fit Buttons** — Press while holding number keys to apply fit only to those channel numbers (e.g.,
   hold 2 and 3 while clicking Fit to fit channels 2 and 3 only).
2. **Mouse Wheel** — Scroll to scale. Hold number keys to target specific channels. Hold **X** for
   X-axis only; hold **Y** for Y-axis only.
3. **Display Settings Window** — Change scale, offset, and color per channel. Switch between
   Offset/Division and Min/Max for axis definition. Also allows changing the time scale.

**Curve representation button** cycles between: line only, points only, or line with points.

The **Print** tab allows adding UPIDs whose values are printed as comments when the oscilloscope window
is printed.

### 3.24 Continuous Two Point Mode

The simplest way to run a motor continuously is **VAI 2 Pos Continuous** mode:

1. Set the run mode to **VAI 2 Pos Continuous**.
2. Set positions under **Trig Fall Config\Position** and **Trig Rise Config\Position**.
3. Optionally configure speed, acceleration, and deceleration at the same location in the parameter tree.
4. Start the motor as described in [Section 3.22](#322-control-status).

The motor moves continuously between the two positions. Dwell times at each position are set under
**VAI 2 Pos Cont Settings**.

### 3.25 Export Configuration

Save the complete drive configuration via **File → Export…** or the Export button.

1. The **Save Config** window opens. Select drives to export.
2. Without **Advanced Options**: all variables (including unread ones) are read before saving.
3. With **Advanced Options**: select specific parts per drive (parameters, curves, command table, etc.).
4. Choose a filename and folder.

> **Tip for support requests:** Export without advanced options with all drives selected — this gives
> supporters all necessary information.

### 3.26 Import Configuration

Import a configuration via **File → Import…** or the Import button.

| Import Symbol | Meaning |
|---------------|---------|
| (none) | Not used — nothing happens |
| Open offline | Creates an offline device and loads the configuration |
| Same drive type | Import to a drive matching the configuration's drive type |
| Different drive type | Import to a different drive type — may cause inconsistent parameter trees |

In the green area on the right side, select which parts of the configuration to import (e.g., only curves or
command table). Only one configuration can be assigned to a drive at a time.

When importing to a drive, a compatibility list is shown before the import begins.

### 3.27 Open Offline Configuration

Open a configuration without a physical drive via **File → Login/Open Offline…**. Useful for support and
offline configuration work.

### 3.28 Create Offline Configuration

Create a default configuration for any supported drive without a physical connection via **File → Create
Offline…**:

1. Select the drive family.
2. Select the drive type.
3. Choose the interface and application software.
4. The offline configuration is created with all parameters at default values.

The configuration can be modified and saved normally.

### 3.29 Compare Parameters

**Drive → Compare Parameter** compares settings between online and offline drives. Options include
filtering by firmware instance (OS, MC, INTF, APPL) or parameter type (read-only or writable). The
result can be saved as a `.pvl` file (comma-separated text).

### 3.30 Portable App

Create a portable installation of LinMot-Talk (e.g., for a USB drive):

1. Go to **Help → Update Functions → Generate Portable App**.
2. Click **Browse** to choose a target folder (or type a path and check **Generate the path if it does not
   exist**).
3. Click **OK** — LinMot-Talk copies all necessary files to the folder.

> **Tip:** Update motor files and download all required older releases from the LinMot server **before**
> generating the portable app. Generation takes some time — wait for the mouse cursor to return to
> normal before using the portable app.

---

## 4 Troubleshooting

### 4.1 Setting all Parameters to Default Values

With LinMot-Talk connected, use the **DEF** button (see [Section 2.1](#21-tool-button-bar)). Without LinMot-Talk,
use the hardware procedures below.

#### E1100, E1200, E1400, B8050, MB8050, C1100-GP, and C1250 (SG3/SG5)

1. Power off the drive.
2. Set both ID switches to `0xFF`.
3. Power on the drive — the Error and Warn LEDs blink alternately at ~4 Hz.
4. Set both ID switches to `0x00`.
5. Wait until the Warn and EN LEDs flash together at ~2 Hz.
6. Power off, then power on again.

#### B1100 (SG4)

1. Set the parameter with UPID `0x6085` to `0x0001`.
2. Power off the drive.
3. Power on the drive.

The parameter `0x6085` is automatically cleared back to `0x0000`. The default image was stored during
firmware installation.

#### A1100

1. Power off the drive.
2. Set DIP switch **S5.2** to **ON**.
3. Power on the drive — the Error and Warn LEDs blink alternately at ~4 Hz.
4. Set DIP switch **S5.2** to **OFF**.
5. Wait until the Warn and EN LEDs flash together at ~2 Hz.
6. Power off, then power on again.

### 4.2 Interface Does Not Run

If interface software (DeviceNet, CANopen, Profibus, LinRS) is not communicating, check:

| Check | Action |
|-------|--------|
| Interface software installed? | Install the correct interface software |
| Switch S3.4 "Interface" | Must be **ON** (for LinRS: **OFF** when configuring over RS232; **ON** when running LinRS) |
| UPID `2008h` | Must not be set to disabled |
| Baud rate and Node ID | Verify correct values in parameters and on ID switches |

### 4.3 Stopping Firmware

When the same link is used for configuration (RS232) and the interface (e.g., LinRS), LinMot-Talk may
not be able to log in.

- **E1100**: Set interface switch **S3.4** to **OFF** and power cycle — the interface software will be
  deactivated and the configuration link freed.
- **B1100** (or if S3.4 doesn't help): Use **File → Open → StopFirmware.sct**, which continuously tries
  to stop the firmware during the first **2 seconds** after power-up.

### 4.4 Communication Debug Window

Open via **Tools → RSTalk Debug Window**. Enable the **Enable Debug Mode** checkbox to view live
communication between LinMot-Talk and the drives.

**Automatic logging:** Enable via **Options → Save Debug Window Data**. When active, communication
is logged to:
```
C:\Users\username\AppData\Local\LinMot\LinMot-Talk6.6-BuildXXXXXXXX\Communication\
```
LinMot-Talk maintains a **10-file ring buffer** (oldest file overwritten). This option is automatically
disabled on every restart of LinMot-Talk.

### 4.5 Measures Against Cyber Attacks

#### 4.5.1 Firmware on the Drives

Firmware is transmitted **encrypted** from LinMot-Talk to the drives. Each firmware part is secured with
a checksum — if the checksum is invalid, that firmware part will not start.

#### 4.5.2 Configuration of the Drives

Drive configuration is **not** protected against changes. Because a connected PLC can legitimately change
parameters, it is not possible to exclude parameter changes from correctly crafted fieldbus packets.
However, such packets would need to match very specifically (e.g., correct IP address, UDP port, and
command format for EtherNet/IP).

A **hash value** (UPID `00A1h`) is calculated over the configuration at startup. Record this value after a
fresh install — on subsequent startups, compare the hash to detect unauthorized configuration changes.

#### 4.5.3 Restoring the Correct State

**Compromised firmware:**
1. Re-download LinMot-Talk from the LinMot homepage (or restore from an uncontaminated backup).
2. Reinstall firmware using LinMot-Talk.

**Compromised configuration:**
1. Default all parameters (using LinMot-Talk DEF or the hardware procedure in [Section 4.1](#41-setting-all-parameters-to-default-values)).
2. Load a known-good saved configuration onto the drive.

> **Important:** Default all parameters **before** loading a saved configuration — a configuration file may
> not contain every parameter, and unloaded parameters would retain potentially compromised values.

---

## Contact & Support

**Switzerland:**
NTI AG, Bodenaeckerstrasse 2, CH-8957 Spreitenbach
- Sales: +41 56 419 91 91 | office@linmot.com
- Tech Support: +41 56 544 71 00 | support@linmot.com | http://www.linmot.com/support
- Skype: support.linmot
- Fax: +41 56 419 91 92 | Web: http://www.linmot.com

**USA:**
LinMot USA Inc., N1922 State Road 120, Unit 1, Lake Geneva, WI 53147
- Phone: 262-743-2555 | usasales@linmot.com | http://www.linmot-usa.com/

For worldwide distribution: http://www.linmot.com/contact
