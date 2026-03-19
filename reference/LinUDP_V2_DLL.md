# LinUDP Library — DLL Integration Guide

**Version:** 2.1.1 (en), 17/06/21
**Minimum firmware release:** 6.10 Build20210521
**Minimum .NET:** V4.0

Tested with:
- Microsoft Excel 2013 / Windows 7 32 & 64 bit / Windows 8
- LabView 2011
- Microsoft Visual Studio 2015

© 2021 NTI AG — LinMot® is a registered trademark of NTI AG.

---

## Use of this Document

The presented DLL provides function blocks to control LinMot drives over the LinUDP interface. The library is provided by NTI AG / LinMot free of charge with no warranty for updates. LinMot accepts no liability for damages that may be caused by using this library.

**Applicable controllers:** C1250-LU(-XX), E1250-LU(-XX), E1450-LU(-XX)

### Recommended Documentation

1. LinMot-Talk
2. Motion Control Software
3. Installation guide C1200 servo controllers
4. Installation guide E1200 servo controllers
5. Installation guide E1400 servo controllers
6. TF Force Control

---

## Release Notes

| DLL Version | Description |
|-------------|-------------|
| 2.0.0 | Initial version for use with LinUDP V2; all commands compatible with previous DLL V1.0.2 |
| 2.0.1 | Additional functions for Master-Slave applications: `setSwitchOn`, `setHoming`, `isMasterSlaveOperationEnabled`, `MasterSlaveHoming` |
| 2.0.2 | Network tools added for internal use |
| 2.0.3 | `LMmt_WriteLivePar` – correction of data order; cycle time of communication task improved |
| 2.0.4 | Special functions with process data combined with monitoring channel for drive timestamp; "Pseudo" oscilloscope implemented; `LMmt_GenericMC` reworked |
| 2.0.5 | `LMav_MoveBestehornRelative` and `LMav_MoveSinRelative` added |
| 2.0.6 | DLL adapts drive control word on first start; Start/Stop/Reboot function reworked; helper functions for data type conversion added |
| 2.0.7 | Error Acknowledge Bit function added; additional functions for setting Control Word bits 0–15; additional state machine status evaluation function |
| 2.0.8 | New force control functions: VAI Go To Pos From Act Pos And Reset Force Control Set I (386xh), VAI Increment Act Pos And Reset Force Control Set I (387xh), VAI Inc Act Pos With Higher Force Ctrl Limit and Target Force (388xh); deceleration issue on `VAI Stop` fixed |
| 2.1.0 | New curve loading function; load/save Command Table functions; save/restore drive configuration functions |
| 2.1.1 | Communication behavior adjusted for LT 6.10 LinUDP V2 bug fix |

---

## 1. Overview LinUDP

The LinUDP protocol is an easy way to communicate with a LinMot servo drive over Ethernet using UDP. Because UDP provides no delivery guarantees, the drive only responds to requests — it has no active role. LinUDP uses fixed ports (configurable in LinMot Talk):

- **Master (PC):** port 41136
- **Slaves (LinMot drives):** port 49360

> **Attention:** This DLL can only be used with servo controllers of type `X1X50-LU-XX`. If Ethernet/IP controllers are used, LinUDP V2 is not available — use the latest DLL version 1 instead.

> **Firewall:** Configure firewall port settings to allow both ports if a firewall is active.

### 1.1 Timing of LinUDP DLL

**LinUDP is not a real-time protocol.** The DLL and user application depend on system timing controlled by the operating system — changes in communication timing and command execution may occur.

Typical data transmission timing tested on Dell Latitude E5570 with Windows 7 64-bit:
- ~49.86% of datagrams handled in less than 5 ms
- ~50.74% above 5 ms cycle time

The actual cycle time varies randomly. No deterministic data exchange is possible. Readings of position or current data versus time are not recommended for high-speed movements.

### 1.2 Working Principle of LinUDP DLL

```
LinUDP DLL
  Task Timer
    Axis Memory 1  ←→  Network
    Axis Memory 2  ←→  Network
    Axis Memory …  ←→  Network
  DLL functions  ←→  User application
```

The DLL reserves data memory for each registered drive (Axis Memory). The Task Timer handles UDP packet exchange between the DLL and the individual drives. DLL functions exposed to the user application only access this data memory — they read data or set motion commands. All actual network I/O is managed by the Task Timer as a Windows task.

### 1.3 Installation of the DLL

- **LabVIEW only:** Use the DLL directly from the LabVIEW example folder — no installation required.
- **Visual Studio / Excel:** Install the DLL (requires Windows 7 or higher and administrator rights). Run the setup program in the `Setup` folder.
- **Windows XP:** Manual installation required — see [Appendix V](#appendix-v-manual-installation-of-linudp-dll-driver).

To uninstall, open System Manager → Software → select "LinUDP Interface Test" → Uninstall. (Uninstalling the test tool also removes the DLL.)

The installer registers the DLL per Microsoft .NET guidelines and installs a small test utility. See [Appendix III](#appendix-iii-linudp-testtool) for details.

### 1.4 Using LinUDP with Microsoft Excel

Install the DLL (see section 1.3), then open the example project in `\Examples\Excel\`. See [Appendix V.I](#vi-i-dll-installation-for-windows-7) for how to add the library reference in Excel.

### 1.5 Using LinUDP with LabVIEW

The DLL is a .NET assembly. For LabVIEW, create an additional configuration file in the folder containing the running `.exe`. For `LabVIEW.exe`, name it `LabVIEW.exe.conf` (some versions use `.config`):

```xml
<?xml version ="1.0"?>
<configuration>
  <startup useLegacyV2RuntimeActivationPolicy="true">
    <supportedRuntime version="v4.0.30319"/>
  </startup>
</configuration>
```

This is necessary for LabVIEW 2011, which uses .NET 2.0 by default. Refer to National Instruments documentation for higher LabVIEW versions. Example project is in `\Examples\LabView\`.

---

## 2. Running LinUDP by Using the LinMot DLL

### Data Types

| .NET Type | Description | Size | Range |
|-----------|-------------|------|-------|
| Boolean | True or False | depends | — |
| Byte | Unsigned byte | 1 byte | 0 to 255 |
| Double | Double-precision float | 8 bytes | ±1.8×10³⁰⁸ |
| Int16 | Short integer | 2 bytes | −32,768 to 32,767 |
| Int32 | Integer | 4 bytes | −2,147,483,648 to 2,147,483,647 |
| Int64 | Long integer | 8 bytes | −9,223,372,036,854,775,808 to 9,223,372,036,854,775,807 |
| SByte | Signed byte | 1 byte | −128 to 127 |
| Single | Single-precision float | 4 bytes | ±3.4×10³⁸ |
| String | Unicode string | variable | Up to 2 billion characters |
| UInt16 | Unsigned short | 2 bytes | 0 to 65,535 |
| UInt32 | Unsigned integer | 4 bytes | 0 to 4,294,967,295 |
| UInt64 | Unsigned long | 8 bytes | 0 to 18,446,744,073,709,551,615 |

### 2.1 Minimum Steps to Operate a LinMot Drive

To get started:

1. Create an `ACI` object (Axis Communication Interface)
2. Call `ClearTargetAddressList()`
3. Register each drive IP with `SetTargetAddressList(TargetIP, TargetPort)`
4. Call `ActivateConnection(HostIP, Port)` to open the UDP socket

> **Note:** Calling `ActivateConnection` a second time without first closing the connection causes a runtime error. Call `CloseConnection()` before re-initializing.

```vb
Dim ACI1 As LinUDP.ACI()

Private Sub EstablishConnection()
    ' Register drives IP address and establish connection
    ACI1.CreateTargetAddressList
    ACI1.ClearTargetAddressList
    ACI1.SetTargetAddressList(IPAddress1, TargetPort)
    ACI1.SetTargetAddressList(IPAddress2, TargetPort)
    X = ACI1.ActivateConnection(HostIP, Port)
End Sub

Private Sub CloseConnection()
    ' Close connection
    X = ACI1.CloseConnection
End Sub
```

> **Note:** Port number assignment is optional. Pass empty strings to use the default ports (`TargetPort` = 49360, `LocalPort` = 41136).

#### 2.1.1 Axis Communication Interface

The ACI object manages all UDP datagrams to/from connected drives.

```vb
' Declare the object
Public Shared ACI As LinUDP.ACI

' Create the object
ACI = New LinUDP.ACI

' Register drives
ACI.ClearTargetAddressList
ACI.SetTargetAddressList(IPAddress1, TargetPort)
ACI.SetTargetAddressList(IPAddress2, TargetPort)

' Activate (pass empty strings for defaults)
X = ACI.ActivateConnection("", "")
```

If you have multiple network adapters, specify the host IP (e.g., `"192.1.1.10"`). If not using the default port 41136, enter the configured port number.

#### 2.1.2 Axis Control by Using ACI Methods

All drive-related functions require the drive's IP address (`TargetIP`):

```vb
' Check if drive is switched on
IsActive = ACI.isSwitchOnActive(TargetIP)   ' TargetIP e.g. "192.1.3.5"

' Get actual position (mm)
ActPos = ACI.getActualPos(TargetIP)

' Send absolute move command
ACI.LMmt_MoveAbs(TargetIP, Pos, MaxVel, Acc, Dec)

' With boolean return (true = TargetIP found in drive list)
isDone = ACI.LMmt_MoveAbs(TargetIP, Pos, MaxVel, Acc, Dec)
```

#### 2.1.3 State Machine Handling with ACI Functions

All LinMot drives implement a state machine for basic operation. See the `Usermanual_MotionCtrlSW_e_recent.pdf` for details.

Example: switching on and homing a drive:

```vb
ACI.SwitchOn(TargetIP)
If ACI.isSwitchOnActive(TargetIP) And Not ACI.isHomed(TargetIP) Then
    ACI.Homing(TargetIP)
End If
```

---

### 2.2 Overview: ACI Methods and Functions (Class ACI)

#### 2.2.1 ACI-Specific Methods, Functions, and Data Access

```
New()
boolean getTimerCycle()
boolean SetTimerCycle(Cycle As UInt32)
long getDatagramCycleTime(TargetIP As String)
boolean CreateTargetAddressList()           [Obsolete]
boolean ClearTargetAddressList()
boolean SetTargetAddressList(TargetIP As String, TargetPort As String)
boolean ActivateConnection(LocalIP As String, LocalPort As String)
setHostMAC(MAC As String, Port As String)
string getHostIP()
string SetTargetAddressListByMAC(TargetMAC As String, TargetPort As String)
boolean ActivateConnectionByMAC(LocalMAC As String, LocalPort As String)
string getDriveIP_byMAC(MACAdress As String, HostIP As String, HostNetMask As String, TargetPort As String)
boolean CloseConnection()
Free()

DLL Status:
boolean isConnected(TargetIP As String)
boolean isResponseUpToDate(TargetIP As String)
boolean isRealtimeConfigUpToDate(TargetIP As String)
boolean isNetWorkRunning()
string getDLLError()
boolean clearDLLErrors()
String getVersion()
```

#### 2.2.2 State Machine–Related Functions, Methods, and Data Access

```
boolean Active(TargetIP As String)
boolean SwitchOn(TargetIP As String)
boolean Homing(TargetIP As String)
boolean AckErrors(TargetIP As String)
boolean JogPlus(TargetIP As String)
boolean JogMinus(TargetIP As String)
boolean isSwitchOnActive(TargetIP As String)
boolean isEventHandlerActive(TargetIP As String)
boolean isSpecialMotionActive(TargetIP As String)
boolean isInTargetPosition(TargetIP As String)
boolean isHomed(TargetIP As String)
boolean isFatalError(TargetIP As String)
boolean isMotionActive(TargetIP As String)
boolean isRangeIndicator1(TargetIP As String)
boolean isRangeIndicator2(TargetIP As String)
boolean isOperationEnable(TargetIP As String)
boolean isEnableOperation(TargetIP As String)
boolean isError(TargetIP As String)
boolean isSafeVoltageEnable(TargetIP As String)
boolean isQuickStop(TargetIP As String)
boolean isSwitchOnLocked(TargetIP As String)
boolean isWarning(TargetIP As String)

State Machine State Queries:
boolean isNotReadyToSwitchOnSM(TargetIP As String)
boolean isSwitchOnDisabledSM(TargetIP As String)
boolean isReadyToSwitchOnSM(TargetIP As String)
boolean isSetupErrorSM(TargetIP As String)
boolean isErrorSM(TargetIP As String)
boolean isHardwareTestsSM(TargetIP As String)
boolean isReadyToOperateSM(TargetIP As String)
boolean isOperationEnabledSM(TargetIP As String)
boolean isHomingSM(TargetIP As String)

Control Word Bit Setters:
Boolean setHomingBit(TargetIP As String, State As Boolean)
Boolean setSwitchOnBit(TargetIP As String, State As Boolean)
Boolean setErrorAcknowledgeBit(TargetIP As String, State As Boolean)
Boolean setBit0(TargetIP As String, State As Boolean)
...
Boolean setBit15(TargetIP As String, State As Boolean)
Boolean setJogPlus(TargetIP As String, State As Boolean)
Boolean setJogMinus(TargetIP As String, State As Boolean)

Data Readings:
double getActualPos(TargetIP As String)
double getCurrent(TargetIP As String)
double getDemandPos(TargetIP As String)
long getMonitoringChannel1(TargetIP As String)
long getMonitoringChannel2(TargetIP As String)
long getMonitoringChannel3(TargetIP As String)
long getMonitoringChannel4(TargetIP As String)
Boolean isMasterSlaveOperationEnabled(TargetIP As String)
Boolean MasterSlaveHoming(TargetIP As String)
StateMachineStates getStateMachineState(TargetIP As String)

Timestamp Functions:
TimeStampData getActualPosWithTimestamp(TargetIP As String)
TimeStampData getCurrentWithTimestamp(TargetIP As String)
TimeStampData getDemandPosWithTimestamp(TargetIP As String)
TimeStampMonitoring getMonitoringChannelWithTimestamp(TargetIP As String, Channel As Integer)
TimeStampDataUTC getActualPosWithTimestampUTC(TargetIP As String)
TimeStampDataUTC getCurrentWithTimestampUTC(TargetIP As String)
TimeStampDataUTC getDemandPosWithTimestampUTC(TargetIP As String)
TimeStampMonitoringUTC getMonitoringChannelWithTimestampUTC(TargetIP As String, Channel As Integer)

PseudoScope:
Integer PseudoScopeSamples()
Boolean EnablePseudoScopeTrace(TargetIP As String)
Boolean isPseudoScopeSampling()
List(of PseudoScopeEntry) getPseudoScopeSamples(TargetIP As String)
Integer getPseudoScopeProgress(TargetIP As String)

Helper Functions:
UInt32 LongRawDataToUInt32(value As Long)
Int32  LongRawDataToInt32(value As Long)
UInt16 LongRawDataToUInt16(value As Long)
Int16  LongRawDataToInt16(value As Long)
```

##### StateMachineStates Enumeration

```vb
Enum StateMachineStates As UShort
    NotReadyToSwitchOn          = 0
    SwitchOnDisabled            = 1
    ReadyToSwitchOn             = 2
    SetupError                  = 3
    GeneralError                = 4
    HWTests                     = 5
    ReadyToOperate              = 6
    OperationEnabled            = 8
    Homing                      = 9
    HomingFinished              = 90
    ClearanceCheck              = 10
    ClearanceCheckFinished      = 100
    GoingToInitialPosition      = 11
    GoingToInitialPositionFinished = 110
    Aborting                    = 12
    Freezing                    = 13
    QuickStopOnError            = 14
    GoingToPosition             = 15
    GoingToPositionFinished     = 150
    JoggingPlus                 = 16
    Jogging_MovingPositive      = 161
    JoggingPlusFinished         = 162
    JoggingMinus                = 17
    Jogging_MovingNegative      = 171
    JoggingMinusFinished        = 172
    Linearizing                 = 18
    PhaseSearch                 = 19
    SpecialMode                 = 20
    BrakeDelay                  = 21
End Enum
```

#### 2.2.3 Motion Commands, Parameter Access

```
boolean LMmt_MoveAbs(TargetIP, Pos1, MaxVel1, Acc1, Dec1)
boolean LMmt_MoveRel(TargetIP, Pos1, MaxVel1, Acc1, Dec1)
boolean LMmt_StartCTCommand(TargetIP, CTEntryID As UInteger)
boolean LMmt_ClearEventEvaluation(TargetIP)
boolean LMmt_Stop(TargetIP, Decceleration As Single)
boolean LMmt_WriteLivePar(TargetIP, UPID As UInteger, UPIDValue As Single)
boolean LMav_Mod16BitCTPar(TargetIP, CTEntryID, ParaOffset, ParaValue As Integer)
boolean LMav_Mod32BitCTPar(TargetIP, CTEntryID, ParaOffset, ParaValue As Integer)
boolean LMav_RunCurve(TargetIP, CurveID, CurveOffset, TimeScale, AmplitudeScale)
boolean LMav_MoveBestehorn(TargetIP, Position, Velocity, Acceleration, Jerk)
boolean LMav_MoveBestehornRelative(TargetIP, Position, Velocity, Acceleration, Jerk)
boolean LMav_MoveSin(TargetIP, Position, Velocity, Acceleration)
boolean LMav_MoveSinRelative(TargetIP, Position, Velocity, Acceleration)
boolean LMfc_ChangeTargetForce(TargetIP, TargetForce As Single)
boolean LMfc_GoToPosForceCtrlHighLim(TargetIP, Position, Velocity, Acceleration, ForceLimit, TargetForce)
boolean LMfc_GoToPosForceCtrlLowLim(TargetIP, Position, Velocity, Acceleration, ForceLimit, TargetForce)
boolean LMfc_GoToPosRstForceCtrl(TargetIP, Position, Velocity, Acceleration, Deceleration)
boolean LMfc_GoToPosRstForceCtrlSetI(TargetIP, Position, Velocity, Acceleration, Deceleration)
boolean LMfc_IncrementActPosAndResetForceControlSetI(TargetIP, Position, Velocity, Acceleration, Deceleration)
boolean LMfc_IncrementActPosWithHigherForceCtrlLimitAndTargetForce(TargetIP, PositionIncrement, MaxVelocity, Acceleration, ForceLimit, TargetForce)
string  LMcf_GetErrorTxt(TargetIP)
integer LMcf_GetErrorCode(TargetIP)
string  LMcf_GetWarningTxt(TargetIP)
integer LMcf_GetWarningCode(TargetIP)
long    LMcf_StartStopDefault(TargetIP, Mode As Integer)
long    getROM_ByUPID(TargetIP, UPID As UInt32)
long    getRAM_ByUPID(TargetIP, UPID As UInt32)
long    getMinVal_ByUPID(TargetIP, UPID As UInt32)
long    getMaxVal_ByUPID(TargetIP, UPID As UInt32)
long    getDefault_ByUPID(TargetIP, UPID As UInt32)
long    setRAM_ByUPID(TargetIP, UPID As UInt32, Value As Long)
long    setROM_ByUPID(TargetIP, UPID As UInt32, Value As Long)
long    setRAM_ROM_ByUPID(TargetIP, UPID As UInt32, Value As Long)
boolean LMav_SetCurrentCommandMode(TargetIP, Current As Single)
boolean LMav_ResetCurrentCommandMode(TargetIP)
boolean LMmt_GenericMC(TargetIP, MCHeader As UInt16, MCParaWord0..MCParaWordN)
```

#### 2.2.4 Curve Access

```
boolean LMcf_LoadCurve(TargetIP, CurveID, SetpointCount, CurveName, Xlength, Xdim, Ydim, Mode, Setpoints As Integer())
Boolean LMcf_LoadCurve(TargetIP, Mode, CurveData As CurveDataDefinition)
boolean LMcf_isCurveLoading(TargetIP)
integer LMcf_getCurveProgress(TargetIP)
CurveDataDefinition getUploadedCurveData(TargetIP)
Boolean LMcf_setDownloadCurveData(TargetIP, CurveData As CurveDataDefinition)
Boolean LMcf_StartUploadCurve(TargetIP, CurveID As UInt16)
Boolean LMcf_StartDownloadCurve(TargetIP)
Integer() LMcf_getAllCurveID(TargetIP)
Boolean LMcf_isCurveOnDrive(TargetIP, CurveID As UInt16)
Boolean LMcf_DeleteAllCurvesInRAM(TargetIP)
Boolean LMcf_SaveAllCurvesFromRAMToFLASH(TargetIP)
```

#### 2.2.5 Command Table Access

```
CommandTableStructure LMcf_getCommandTableContent(TargetIP)
UInt16 LMcf_DeleteCommandTable_RAM(TargetIP)
UInt16 LMcf_WriteCommandTableToFLASH(TargetIP)
Boolean LMcf_setCommandTableContent(TargetIP, CT As CommandTableStructure)
```

#### 2.2.6 Drive Parameter Access

```
List(Of UPID_List) LMcf_getUPIDList(TargetIP, StartUPID As UShort, StopUPID As UShort)
List(Of UPID_List) LMcf_getModified_UPIDList(TargetIP, StartUPID As UShort, StopUPID As UShort)
Boolean LMcf_setUPIDList(TargetIP, UPIDList As List(Of UPID_List))
```

#### 2.2.7 Load / Save Drive Configuration

```
Boolean Save_CommandTable(TargetIP, FilePath As String)
Boolean Save_CommandTable(CT As CommandTableStructure, FilePath As String)
Boolean Load_CommandTable(TargetIP, FilePath As String)
CommandTableStructure Load_CommandTable(FilePath As String)
Boolean Save_DriveParameters(TargetIP, FilePath As String)
Boolean Save_DriveParameters(Data As List(Of UPID_List), FilePath As String)
Boolean Load_DriveParameters(TargetIP, FilePath As String)
List(Of UPID_List) Load_DriveParameters(FilePath As String)
Boolean Save_Curves(TargetIP, FilePath As String)
Boolean Load_Curves(TargetIP, FilePath As String)
Boolean Save_DriveConfiguration(TargetIP, FilePath As String)
Boolean Load_DriveConfiguration(TargetIP, FilePath As String)
UInt32 getDriveHash(TargetIP)
Boolean isDriveConfigurationSame(TargetIP, FilePath As String)
```

---

### 2.3 Details

#### 2.3.1 ACI-Specific Methods, Functions, and Data Access in Detail

##### `boolean CreateTargetAddressList()`
**Obsolete.** Kept for compatibility. Do not use in new applications.

##### `boolean ClearTargetAddressList()`
Empties the registered target drive address list. Call this after `CreateTargetAddressList`. Any communication must be inactive. To add new controllers, call `CloseConnection()` first, then clear and rebuild the list.

##### `boolean SetTargetAddressList(TargetIP As String, TargetPort As String)`
Creates an axis module in the ACI object and sets the drive's IP and port address. Must be called for every drive before calling `ActivateConnection`. Pass an empty string for `TargetPort` to use the default (49360).

##### `boolean ActivateConnection(LocalIP As String, LocalPort As String)`
Establishes the data link between the DLL object and the target drives. Pass an empty string for `LocalIP` to use the first available network adapter. Pass an empty string for `LocalPort` to use the default (41136). All registered drives must use the same master port.

##### `boolean CloseConnection()`
Closes all UDP connections to registered drives. Call this before re-registering drives or when closing your application to release allocated resources.

##### `Free()` / `Dispose()`
`Free()` is replaced by `Dispose()`. Should be called when the application is closing to release all DLL system resources. After calling `Dispose()`, further function calls will cause an "Index out of range" error.

> **Note:** Since the new communication stack, these functions are largely no longer needed.

##### `boolean setTimerCycle(Cycle As UInt32)`
Changes the internal data exchange interval (default: 5 ms). Returns false if `Cycle <= 0`. Cycle time is in milliseconds.

> **Note:** Communication runs via Windows-specific modules and is not real-time. The configured cycle is a minimum — actual cycle times will vary.

##### `integer getTimerCycle()`
Returns the currently configured cycle time in ms.

##### `long getDatagramCycleTime(TargetIP As String)`
Returns the time span between transmitting a request and receiving the drive response, in milliseconds.

##### `string getDLLError()`
Returns active DLL error messages, or an empty string if none.

##### `boolean clearDLLErrors()`
Resets the last DLL error.

##### `boolean isConnected(TargetIP As String)`
Returns true if the addressed servo controller is active and data exchange is running.

##### `boolean isResponseUpToDate(TargetIP As String)`
Returns true if the last motion command has been executed successfully.

##### `boolean isRealtimeConfigUpToDate(TargetIP As String)`
Returns true if the last parameter channel command has been executed successfully.

##### `boolean isNetworkRunning()`
Returns true if the DLL is in cyclic data exchange. Only checks the computer-side communication.

##### `String getVersion()`
Returns the DLL version number.

##### MAC Address-Based Functions

> **Attention:** To identify possible LinMot drives, the DLL sends a request to every possible IP address based on the computer's subnet mask. This can generate heavy traffic and take significant time. Use only when necessary.

**`setHostMAC(MAC As String, Port As String)`**
Defines the network adapter by its MAC address. Must be called first to define the adapter for LinUDP. Empty port string uses default 41136.

**`string getHostIP()`**
Returns the IP address of the network card selected by `setHostMAC`.

**`string SetTargetAddressListByMAC(TargetMAC As String, TargetPort As String)`**
Registers a controller by its MAC address (printed on the controller label). Returns the corresponding IP address, which must be stored for later DLL function calls. Empty `TargetPort` uses default 49360.

**`boolean ActivateConnectionByMAC(LocalMAC As String, LocalPort As String)`**
Establishes the data link using the MAC address of the network card. Empty `LocalPort` uses default 41136.

> MAC address format: `XX-XX-XX-XX-XX`

**`string getDriveIP_byMAC(MACAdress As String, HostIP As String, HostNetMask As String, TargetPort As String)`**
Determines a drive's IP address from its MAC. The PC's IP and subnet mask must be provided.

> **Attention:** This process can take several minutes depending on subnet mask size. After a successful search, the participant list is fixed internally for 1 minute. Newly connected subscribers are not recognized until 1 minute after connection. This function only works when no DLL connection is established (`isNetworkRunning()` must be false).

---

#### 2.3.2 State Machine–Related Functions

##### `boolean isSwitchOnActive(TargetIP As String)`
Returns true if the drive is switch-on active.

##### `boolean isEventHandlerActive(TargetIP As String)`
Returns true if an event handler (e.g., a command table sequence) is active. Use this to check if a command table has finished before sending new motion commands.

##### `boolean isSpecialMotionActive(TargetIP As String)`
Returns true if a special motion (e.g., Current Command Mode) is active.

##### `boolean isInTargetPosition(TargetIP As String)`
Returns true if the stator is in the target position.

##### `boolean isHomed(TargetIP As String)`
Returns true if the drive has been homed.

##### `boolean isFatalError(TargetIP As String)`
Returns true if a fatal error has occurred. Fatal errors cannot be acknowledged — they require a drive reboot.

##### `boolean isMotionActive(TargetIP As String)`
Returns true if a motion is currently active.

##### `boolean isNotReadyToSwitchOnSM(TargetIP As String)`
Returns true if the controller is not ready to switch on. Toggle bit 0 with `SwitchOn` to change this state.

##### `boolean isSwitchOnDisabledSM(TargetIP As String)`
Returns true if bit 0 or bit 1 in the Control Word is zero (monitors switch-on and SVE/STO or Voltage Enable inputs, depending on controller type).

##### `boolean isReadyToSwitchOnSM(TargetIP As String)`
Returns true if the controller is ready to be switched on with `SwitchOn`.

##### `boolean isSetupErrorSM(TargetIP As String)`
Returns true if there is a configuration error.

##### `boolean isErrorSM(TargetIP As String)`
Returns true if an error has occurred.

##### `boolean isHWTestsSM(TargetIP As String)`
Returns true while the controller performs its power stage self-test during switch-on.

##### `boolean isReadyToOperateSM(TargetIP As String)`
Returns true when the controller is ready to operate.

##### `boolean isOperationEnabledSM(TargetIP As String)`
Returns true when the controller is in operation and will execute motion commands.

##### `boolean isHomingSM(TargetIP As String)`
Returns true while a homing procedure is running.

##### `boolean isRangeIndicator1/2(TargetIP As String)`
Returns true if range indicator 1 or 2 is active.

##### `boolean isOperationEnable(TargetIP As String)` / `isEnableOperation`
Returns true if operation enable / enable operation is set.

##### `boolean isError(TargetIP As String)`
Returns true if an error has occurred.

##### `boolean isSafeVoltageEnable(TargetIP As String)`
Returns true if the safety voltage enable (SVE or STO) input is set (only available on controllers with SVE/STO).

##### `boolean isQuickStop(TargetIP As String)`
Returns true if QuickStop is triggered. QuickStop uses negative logic — this bit must be set for the controller to operate.

##### `boolean isSwitchOnLocked(TargetIP As String)`
Returns true if QuickStop is set (same negative logic as above).

##### `boolean isWarning(TargetIP As String)`
Returns true if a warning has occurred.

##### `boolean isMasterSlaveOperationEnabled(TargetIP As String)`
Returns true if the master and all connected slaves are ready for operation.

> **Hint:** Only works if Master/Slave application is installed and at least one slave is connected to the master.

##### `double getActualPos(TargetIP As String)`
Returns the actual slider position in mm.

##### `double getCurrent(TargetIP As String)`
Returns the actual motor current in amperes.

##### `double getDemandPos(TargetIP As String)`
Returns the actual slider demand position in mm.

##### `long getMonitoringChannel1/2/3/4(TargetIP As String)`
Returns the monitoring channel value as a long. Raw data must be scaled according to the assigned UPID.

##### Timestamp Functions

> **Attention:** When using timestamp functions, monitoring channel 1 on all drives must be configured to UPID `82h` (Operating Sub Hours).

**`TimestampData getActualPosWithTimestamp(TargetIP As String)`**
Returns actual position (mm) and drive timestamp (ms, resets each hour).

**`TimestampData getCurrentWithTimestamp(TargetIP As String)`**
Returns actual current (A) and drive timestamp (ms, resets each hour).

**`TimestampData getDemandPosWithTimestamp(TargetIP As String)`**
Returns demand position (mm) and drive timestamp (ms, resets each hour).

**`TimestampMonitoring getMonitoringChannelWithTimestamp(TargetIP As String, Channel As Integer)`**
Returns raw monitoring channel value (channels 2–4) and timestamp. Values outside 2–4 are treated as channel 2.

Return data structures:
```vb
Structure TimestampData
    value     As Double
    Timestamp As Long
End Structure

Structure TimestampMonitoring
    value     As Long
    Timestamp As Long
End Structure
```

##### UTC Timestamp Functions

Same as the above timestamp functions, but the timestamp includes hour counting and is delivered as Coordinated Universal Time (UTC).

```vb
Structure TimestampDataUTC
    value     As Double
    Timestamp As DateTime
End Structure

Structure TimestampMonitoringUTC
    value     As Long
    Timestamp As DateTime
End Structure
```

##### Helper Functions

**`UInt32 LongRawDataToUInt32(value As Long)`** — Converts long monitoring channel value to UInt32.
**`Int32 LongRawDataToInt32(value As Long)`** — Converts to Int32.
**`UInt16 LongRawDataToUInt16(value As Long)`** — Converts to UInt16.
**`Int16 LongRawDataToInt16(value As Long)`** — Converts to Int16.

Use these to interpret raw monitoring channel data according to the UPID's data type as shown in LinMot Talk.

##### PseudoScope Feature

A simple data acquisition feature. **Not real-time** — the OS controls timing, so time intervals between readings are not fixed.

> **Attention:** When using PseudoScope, monitoring channel 1 on all drives must be configured to UPID `82h` (Operating Sub Hours).

**`Integer PseudoScopeSamples()`**
Gets or sets the number of data points to record. One data point is stored per network exchange cycle.

> **Caution:** Maximum available data points are not automatically managed — very large values may exhaust system memory. For each used drive, this amount of data points is reserved.

**`Boolean EnablePseudoScopeTrace(TargetIP As String)`**
Enables data acquisition for the specified IP. Starts immediately. Returns true while recording is active.

**`Boolean isPseudoScopeSampling()`**
Returns true while data acquisition is running; false when the defined sample count has been reached.

**`List(of PseudoScopeEntry) getPseudoScopeSamples(TargetIP As String)`**
Reads stored data from DLL memory as a list of `PseudoScopeEntry` structures.

```vb
Structure PseudoScopeEntry
    ActualPos         As Double
    DemandPos         As Double
    Current           As Double
    TimeStamp         As Long
    MonitoringChannel2 As Long
    MonitoringChannel3 As Long
    MonitoringChannel4 As Long
End Structure
```

**`Integer getPseudoScopeProgress(TargetIP As String)`**
Returns acquisition progress as a percentage (0–100%).

##### `boolean Active(TargetIP As String)`
**Deprecated.** Available for backward compatibility only — do not use in new applications.

##### `boolean SwitchOn(TargetIP As String)`
**Toggles** bit 0 in the Control Word of the specified controller (0→1 or 1→0). Returns true if bit was set to 1, false if set to 0.

##### `boolean Homing(TargetIP As String)`
Performs a homing procedure. Returns true on success. Times out after 60 seconds (timeout error written to DLL error).

##### `boolean AckErrors(TargetIP As String)`
Resets pending errors on the specified controller. Fatal errors cannot be reset.

##### `boolean JogPlus(TargetIP As String)` / `boolean JogMinus(TargetIP As String)`
Starts jog mode in the positive/negative direction. Call again to stop. Jog behavior must be configured in LinMot Talk.

##### `boolean setHomingBit(TargetIP As String, State As Boolean)`
Sets/resets bit 11 (homing) in the Control Word directly. The user application must manage the homing procedure — reset this bit when homing is complete, or the drive remains in homing state.

##### `boolean setSwitchOnBit(TargetIP As String, State As Boolean)`
Sets/resets bit 0 in the Control Word.

##### `boolean setErrorAcknowledgeBit(TargetIP As String, State As Boolean)`
Sets/resets bit 7 in the Control Word.

##### `boolean setBit0` … `setBit15(TargetIP As String, State As Boolean)`
Sets/resets the corresponding bit (0–15) in the Control Word.

##### `boolean setJogPlus(TargetIP As String, State As Boolean)` / `setJogMinus`
Sets/resets the Jog+ / Jog− bit in the Control Word.

##### `boolean MasterSlaveHoming(TargetIP As String)`
Performs homing on the master and observes homing state of connected slaves. Times out after 5 minutes.

> **Hint:** Only works if Master/Slave application is installed and at least one slave is connected.

##### `StateMachineStates getStateMachineState(TargetIP As String)`
Returns the current state machine state as a `StateMachineStates` enumeration value (see section 2.2.2).

---

#### 2.3.3 Motion Commands, Parameter Access in Detail

> **Note:** All motion commands return boolean: `true` = command handled successfully, `false` = malfunction. The count nibble is managed internally by the DLL — set it to zero in generic function calls.

##### `boolean LMmt_MoveAbs(TargetIP, Pos1, MaxVel1, Acc1, Dec1 As Single)`
Absolute movement to `Pos1` (relative to home position). Parameters in mm and m/s units.

##### `boolean LMmt_MoveRel(TargetIP, Pos1, MaxVel1, Acc1, Dec1 As Single)`
Relative movement by `Pos1` distance (relative to actual position).

##### `boolean LMmt_StartCTCommand(TargetIP, CTEntryID As UInteger)`
Starts the command table at the specified entry ID.

##### `boolean LMmt_ClearEventEvaluation(TargetIP)`
Clears an active trigger or cancels a running command table. Note: depending on command table execution, an "Unknown Motion Command" error may occur. A command table should normally terminate itself.

##### `boolean LMmt_Stop(TargetIP, Decceleration As Single)`
Stops the current movement with the specified deceleration (m/s²).

##### `boolean LMmt_WriteLivePar(TargetIP, UPID As UInteger, UPIDValue As Integer)`
Writes to the specified UPID on the drive. Only RAM values with runtime access can be set.

##### `boolean LMav_Mod16BitCTPar(TargetIP, CTEntryID, ParaOffset, ParaValue As Integer)`
Modifies a 16-bit parameter in the command table (RAM only) on E series drives.

| Offset | Meaning |
|--------|---------|
| 2 | Linked ID (following line ID) |
| 4 | Motion Command Header |
| 6 | First motion command parameter |

##### `boolean LMav_Mod32BitCTPar(TargetIP, CTEntryID, ParaOffset, ParaValue As Integer)`
Modifies a 32-bit parameter in the command table (RAM only) on E series drives.

| Offset | Meaning |
|--------|---------|
| 6 | First motion command parameter |
| 10 | Second motion command parameter |

##### `boolean LMav_RunCurve(TargetIP, CurveID, CurveOffset, TimeScale, AmplitudeScale)`
Starts a predefined curve stored on the controller. `TimeScale` is scaled at 0.01% (range 0–200%); `AmplitudeScale` at 0.1% (range −2000%–2000%).

##### `boolean LMav_MoveBestehorn(TargetIP, Position, Velocity, Acceleration, Jerk)`
Absolute movement with Bestehorn profile (limited jerk). Supported on X12x0 and E14x0 series. Units: Position in mm, Velocity in m/s, Acceleration in m/s², Jerk in m/s³.

##### `boolean LMav_MoveBestehornRelative(TargetIP, Position, Velocity, Acceleration, Jerk)`
Relative movement / demand position increment with Bestehorn profile. Same units and series support as above.

##### `boolean LMav_MoveSin(TargetIP, Position, Velocity, Acceleration)`
Absolute movement with sinoid profile. Supported on X12x0 and E14x0 series.

##### `boolean LMav_MoveSinRelative(TargetIP, Position, Velocity, Acceleration)`
Relative movement / demand position increment with sinoid profile.

##### `boolean LMmt_GenericMC(TargetIP, MCHeader As UInt16, MCParaWord0..N)`
Sends any motion command described in the Motion Control Software manual. Set `MCHeader` to the hex motion command header value (as decimal or with conversion). The count nibble is handled internally — enter `0` for the nibble placeholder in the header.

**Example — VAI Go To Pos (010xh):**

Set `x = 0` → MCHeader = `0100h` = 256 decimal.

**Curve run example:**

| Parameter | Type | Decimal | Hex |
|-----------|------|---------|-----|
| MCHeader | UInt16 | 1248 | 04E0 |
| Curve ID = 1 | UInt16 | 1 | 0001 |
| Position = 50mm | SInt32 | 500000 | 7A120 |
| Time = 750ms | SInt32 | 75000 | 124F8 |

```vb
.LMmt_GenericMC("192.168.0.1", CUInt16(1248), CUInt16(1), CInt(500000), CInt(75000))
```

**Position move with current limit example:**

| Parameter | Type | Decimal | Hex |
|-----------|------|---------|-----|
| MCHeader | UInt16 | 3152 | 0C50 |
| Target Position = 50mm | SInt32 | 500000 | 7A120 |
| Max Velocity = 2 m/s | UInt32 | 2000000 | 1E8480 |
| Acc = Dec = 3 m/s² | UInt32 | 300000 | 493E0 |
| Current Limit = 2A | UInt16 | 2000 | 7D0 |

```vb
.LMmt_GenericMC("192.168.0.1", CUInt16(3152), CInt(500000), CUInt32(2000000), CUInt32(300000), CUInt16(2000))
```

> **Note:** Raw data for all parameter words must be in the specified data type and scale per the Motion Control SW manual. Use explicit type casts when entering literal values.

##### Force Control Functions

All force control functions require the "Force Control" technology function (order no. 0150-2503).

**`boolean LMfc_ChangeTargetForce(TargetIP, TargetForce As Single)`**
Sets the target force (N) in force control mode.

**`boolean LMfc_GoToPosForceCtrlHighLim(TargetIP, Position, Velocity, Acceleration, ForceLimit, TargetForce)`**
Moves to position; switches to force control when force reaches `ForceLimit` (upper limit).

**`boolean LMfc_GoToPosForceCtrlLowLim(TargetIP, Position, Velocity, Acceleration, ForceLimit, TargetForce)`**
Moves to position; switches to force control when actual force falls below `ForceLimit` (lower limit).

**`boolean LMfc_GoToPosRstForceCtrl(TargetIP, Position, Velocity, Acceleration, Deceleration)`**
Resets force control mode to normal position control mode by moving to the specified position.

**`boolean LMfc_GoToPosRstForceCtrlSetI(TargetIP, Position, Velocity, Acceleration, Deceleration)`**
Resets force control mode to position control mode, setting the I-term to the current force value.

**`boolean LMfc_IncrementActPosAndResetForceControlSetI(TargetIP, Position, Velocity, Acceleration, Deceleration)`**
Same as above, but increments the current position rather than moving to an absolute position.

**`boolean LMfc_IncrementActPosWithHigherForceCtrlLimitAndTargetForce(TargetIP, PositionIncrement, MaxVelocity, Acceleration, ForceLimit, TargetForce)`**
Increments current position; switches to force control and adjusts to target force when `ForceLimit` is reached.

##### Error and Warning Functions

**`string LMcf_GetErrorTxt(TargetIP)`** — Returns pending error as string (see Motion Control Software Manual for codes).
**`integer LMcf_GetErrorCode(TargetIP)`** — Returns pending error code as integer.
**`string LMcf_GetWarningTxt(TargetIP)`** — Returns pending warning message.
**`integer LMcf_GetWarningCode(TargetIP)`** — Returns pending warning code.

##### `long LMcf_StartStopDefault(TargetIP, Mode As Integer)`

| Mode | Description |
|------|-------------|
| 0 | Reboot the controller |
| 1 | Set all ROM values to factory default (OS Software) |
| 2 | Set all ROM values to factory default (MC Software) |
| 3 | Set all ROM values to factory default (Interface Software) |
| 4 | Set all ROM values to factory default (Application Software) |
| 5 | Stop Motion Command and Application software |
| 6 | Start Motion Command and Application Software |

Returns `0` on success, `2` if still executing, `9090` on timeout (see also [Return Values During Parameter Access](#224-return-values-during-parameter-access)).

##### UPID Parameter Access Functions

All UPID addresses are entered as UInt32. Addresses shown in hex in LinMot Talk must be converted to integer.

**`long getROM_ByUPID(TargetIP, UPID)`** — Reads UPID value from drive ROM.
**`long getRAM_ByUPID(TargetIP, UPID)`** — Reads UPID value from drive RAM.
**`long getMinVal_ByUPID(TargetIP, UPID)`** — Reads minimum allowed value for UPID.
**`long getMaxVal_ByUPID(TargetIP, UPID)`** — Reads maximum allowed value for UPID.
**`long getDefault_ByUPID(TargetIP, UPID)`** — Reads factory default value for UPID.

**`long setRAM_ByUPID(TargetIP, UPID, Value As Long)`** — Writes value to drive RAM. Returns 0 (OK), 2 (running), or 9090 (timeout).

**`long setROM_ByUPID(TargetIP, UPID, Value As Long)`** — Writes value to drive ROM. Active after next reboot. Returns 0, 2, or 9090.

> **Attention:** ROM memory allows only ~100,000 write cycles. Use ROM writes for configuration only, never in cyclic process steps.

**`long setRAM_ROM_ByUPID(TargetIP, UPID, Value As Long)`** — Writes to both ROM and RAM. Active after next reboot.

##### Current Command Mode

**`boolean LMav_SetCurrentCommandMode(TargetIP, Current As Single)`**
Switches from position control to current command mode. Current is in milliamperes; positive = positive direction, negative = negative direction. Sets `isSpecialMotionActive` to true. No other motion commands execute while active.

> Force accuracy is approximately 10% at the same target position. For higher precision, use closed-loop force control.

**`boolean LMav_ResetCurrentCommandMode(TargetIP)`**
Returns from current command mode to position control mode. Clears `isSpecialMotionActive`.

---

#### 2.2.4 Return Values During Parameter Access

| Code (Dec) | Code (Hex) | Description |
|------------|------------|-------------|
| 0 | 0h | OK — command executed |
| 2 | 2h | Command is running |
| 4 | 4h | Data block is writing (Curve service) |
| 5 | 5h | Busy |
| 192 | C0h | UPID wrong |
| 193 | C1h | Parameter data type wrong |
| 194 | C2h | Range error |
| 195 | C3h | Address usage wrong |
| 196 | C5h | `Get next UPID List item` sent without prior `Start Getting UPID List` |
| 197 | C6h | End of parameter list reached |
| 208 | D0h | Odd address |
| 209 | D1h | Unit error (Curve service) |
| 212 | D4h | Curve already defined / not available |
| 9090 | 2382h | Timeout during parameter channel access |

---

#### 2.3.5 Curve Access

##### `boolean LMcf_LoadCurve(TargetIP, CurveID, SetpointCount, CurveName, Xcode, Ycode, Xlength, Xdim, Ydim, Mode, Setpoints As Integer())`

Loads a curve into drive curve memory. `CurveID` = 1–100. `SetpointCount` minimum = 2. `CurveName` up to 22 characters.

**Ycode / Xcode axis definitions:**

| Ycode | Y Axis | Xcode | X Axis |
|-------|--------|-------|--------|
| 0 | Position | 0 | Time |
| 1 | Velocity | 1 | Encoder Position [Inc] |
| 2 | Current | 2 | Position [mm] |
| 3 | Encoder Position | | |
| 4 | Encoder Speed | | |
| 5 | MicroSteps | | |

**XDim / YDim scale units:**

| Value | Scale Unit | Description |
|-------|-----------|-------------|
| 5 | 0.1 µm | Standard linear position unit |
| 26 | 0.01 ms | Standard curve timebase |
| 27 | 1 increment | Standard encoder position unit |

**Loading sequence:**

| Mode | Description |
|------|-------------|
| 0 | Load curve into RAM |
| 1 | Delete all curves in RAM and ROM |
| 2 | Store all curves from RAM to ROM |

Loading procedure:
1. Call with **Mode 1** to delete all curves (non-TargetIP params can be dummy values).
2. Wait until function returns `false` (deletion complete).
3. Call with **Mode 0** for each curve to load (returns `true` during load, `false` when done).
4. If curves need to persist after reboot, call with **Mode 2**.

> **Note:** The firmware must be stopped when saving to ROM. Run only one curve download at a time.

##### `boolean LMcf_LoadCurve(TargetIP, Mode, CurveData As CurveDataDefinition)`
Alternative overload using a structured data type:

```vb
Structure CurveDataDefinition
    CurveID       As UInt16
    SetpointCount As UInt16
    CurveName     As String
    Xcode         As Byte
    Ycode         As Byte
    XLength       As UInt32
    XDim          As UInt16
    YDim          As UInt16
    Setpoints     As Integer()
End Structure
```

##### `boolean LMcf_isCurveLoading(TargetIP)`
Returns true while a curve load/write process is ongoing.

##### `integer LMcf_getCurveProgress(TargetIP)`
Returns percentage progress (0–100%) of the current curve upload or download.

##### `boolean LMcf_StartUploadCurve(TargetIP, CurveID As UInt16)`
Initiates curve upload from drive to application (background service). Check progress with `LMcf_isCurveLoading()` and `LMcf_getCurveProgress()`. Result is read with `LMcf_getUploadedCurveData()`.

##### `boolean LMcf_StartDownloadCurve(TargetIP)`
Initiates curve download from application to drive (background service). Prepare data with `LMcf_setDownloadCurveData()` first.

##### `CurveDataDefinition LMcf_getUploadedCurveData(TargetIP)`
Returns curve data previously loaded with `LMcf_StartUploadCurve()`.

##### `Boolean LMcf_setDownloadCurveData(TargetIP, CurveData As CurveDataDefinition)`
Sets curve data for download. Then call `LMcf_StartDownloadCurve()`.

##### `integer() LMcf_getAllCurveID(TargetIP)`
Returns an integer array of all available curve IDs on the drive.

##### `boolean LMcf_isCurveOnDrive(TargetIP, CurveID As UInt16)`
Returns true if the specified curve ID exists on the drive.

##### `Boolean LMcf_DeleteAllCurvesInRAM(TargetIP)`
Deletes all curves in RAM.

##### `Boolean LMcf_SaveAllCurvesFromRAMToFLASH(TargetIP)`
Saves current RAM contents to FLASH memory.

> **Attention:** FLASH memory allows only ~100,000 write cycles. Use for configuration only, not in cyclic process steps.

---

#### 2.3.6 Command Table Access

For temporary changes, modifying the Command Table in RAM is sufficient. Saving to FLASH is only needed if the Command Table must survive drive reboots without prior loading.

For individual parameter changes, use: `LMav_Mod32BitCTPar()` and `LMav_Mod16BitCTPar()`

Data structures:
```vb
Structure CTEntry
    data As List(Of Byte)
End Structure

Structure CommandTableStructure
    Entries         As List(Of CTEntry)
    CTPresenceList  As List(Of Byte)
End Structure
```

The data structure content is defined per the manual *"Drive Configuration over Fieldbus SG5-SG7"*.

##### `CommandTableStructure LMcf_getCommandTableContent(TargetIP)`
Loads the Command Table from the drive into the `CommandTableStructure` data structure.

##### `Boolean LMcf_setCommandTableContent(TargetIP, CT As CommandTableStructure)`
Writes `CommandTableStructure` data to the Command Table RAM of the drive.

##### `UInt16 LMcf_DeleteCommandTable_RAM(TargetIP)`
Deletes the Command Table in the drive's RAM.

##### `UInt16 LMcf_WriteCommandTableToFLASH(TargetIP)`
Writes the RAM Command Table area to FLASH.

> **Attention:** FLASH memory allows only ~100,000 write cycles. Use for configuration only.

---

#### 2.3.7 Drive Parameter Access

##### `List(Of UPID_List) LMcf_getUPIDList(TargetIP, StartUPID, StopUPID As UShort)`
Generates a list of all UPIDs from StartUPID to StopUPID with their usage flags. Bit definitions of the returned data value:

| Bit | 15–13 | 12 | 11–6 | 5 | 4 | 3 | 2 | 1 | 0 |
|-----|-------|----|------|---|---|---|---|---|---|
| Meaning | — | Not used for HASH | — | Live Parameter | ROM Write | ROM Read | RAM Write | RAM Read | — |

##### `List(Of UPID_List) LMcf_getModified_UPIDList(TargetIP, StartUPID, StopUPID As UShort)`
Returns all UPIDs whose values differ from factory defaults. Useful for saving drive configurations. Address range: 0 to 0xFFFF.

See *"Drive Configuration over Fieldbus SG5-SG7"* for UPID group overview.

##### `Boolean LMcf_setUPIDList(TargetIP, UPIDList As List(Of UPID_List))`
Writes raw data values from a UPID list to the specified drive. To ensure defined states, first reset the relevant parameter range to factory defaults using `LMcf_StartStopReset()`.

---

#### 2.3.8 Load / Save Drive Configuration

All save functions use XML-based data format. File path and extension are freely definable.

##### Command Table

**`Boolean Save_CommandTable(TargetIP, FilePath As String)`** — Saves Command Table from drive directly to file.
**`Boolean Save_CommandTable(CT As CommandTableStructure, FilePath As String)`** — Saves a loaded Command Table structure to file.
**`Boolean Load_CommandTable(TargetIP, FilePath As String)`** — Loads Command Table from file to drive (RAM only). To persist after reboot, also call `LMcf_WriteCommandTableToFLASH()`.
**`CommandTableStructure Load_CommandTable(FilePath As String)`** — Loads stored Command Table into a structure.

##### Drive Parameters

**`Boolean Save_DriveParameters(TargetIP, FilePath As String)`** — Saves all parameters that differ from factory values.
**`Boolean Save_DriveParameters(Data As List(Of UPID_List), FilePath As String)`** — Saves a UPID list to file.
**`Boolean Load_DriveParameters(TargetIP, FilePath As String)`** — Loads parameters from file to drive.

> **Attention:** This function resets all drive parameters to factory settings first, then loads the parameters. A restart must be triggered afterward.

**`List(Of UPID_List) Load_DriveParameters(FilePath As String)`** — Loads parameters from file into a UPID list.

##### Curves

**`Boolean Save_Curves(TargetIP, FilePath As String)`** — Saves all curves stored on the drive to file.
**`Boolean Load_Curves(TargetIP, FilePath As String)`** — Loads saved curves from file directly into drive.

##### Full Drive Configuration

**`Boolean Save_DriveConfiguration(TargetIP, FilePath As String)`** — Saves the entire drive configuration (parameters, curves, Command Table) to file. This is a blocking operation and may take several minutes.

**`Boolean Load_DriveConfiguration(TargetIP, FilePath As String)`** — Loads a complete configuration to the drive. Resets drive to factory defaults first, then loads configuration. Blocking operation, may take several minutes. Drive must be restarted (with `LMcf_StartStopReset()`) to activate the loaded configuration.

> **Attention:** This function completely resets the drive to factory settings before loading configuration data.

##### Drive Hash / Configuration Integrity

**`UInt32 getDriveHash(TargetIP)`** — Returns the drive's hash value, calculated from relevant parameters during reboot. Bit 12 of UPID data indicates whether a UPID is included in the hash calculation.

**`Boolean isDriveConfigurationSame(TargetIP, FilePath As String)`** — Compares a previously saved configuration file against the current drive to detect parameter changes.

> **Note:** Any parameter change on the controller changes the HASH. The HASH is only recalculated during startup/reboot.

**Example scenario for automated configuration:**
1. After commissioning or parameter adjustments, trigger a drive restart via `LMcf_StartStopReset()` to recalculate the hash.
2. Load and save the new configuration with `Save_DriveConfiguration()`.
3. On application startup, use `isDriveConfigurationSame()` to check for unauthorized changes.
4. If needed, reload configuration with `Load_DriveConfiguration()`, then restart the drive.

---

## Appendix I: Commissioning LinUDP Interface

For use of LinUDP, the interface must be configured. After firmware installation, the interface is at factory defaults:

- **Dis-/Enable:** set to "Enable" (required for operation)
- **Ethernet Configuration → IP Configuration Mode:** set to S1 and S2 switches

**IP Address configuration options:**
- **DHCP:** for automatic addressing with a DHCP server
- **Static by IP Configuration:** enter the static IP in "IP Configuration"

**Monitoring channels:** All channels are unconfigured by default. Up to 4 channels can be assigned UPID addresses — the addressed values are transmitted in responses and accessible via DLL functions.

### Master Configuration (Advanced)

Configure only if needed:
- Change **Master Port** and **Drive Port** (slave) if defaults conflict with your network
- Configure IP filters for multi-master environments:
  - **Single Master:** drive filters to the first recognized IP address
  - **Single Master with fix IP:** specify the master's IP in "Master IP Address"; all other LinUDP masters are ignored

> **Attention:** When using modified port numbers, ensure the DLL connection uses the same ports.

---

## Appendix II: Realtime Connection

LinUDP is available on `-LU` servo controllers, which provide a two-port internal switch for the real-time interface. Standard network switches can be used.

**Connection points:**

| Controller | Connection Location |
|------------|---------------------|
| C1250-LU-XC-xS | See Installation Guide |
| E1250-LU-XC-xS | See Installation Guide |
| E1450-LU-XC-xS | See Installation Guide |

**Supported topologies:** Daisy chain and Star topology.

---

## Appendix III: LinUDP Testtool

The LinUDP test tool is installed with the DLL and provides a quick way to verify drive connection and interface configuration.

> **Note:** Uninstalling the test tool also removes the DLL.

### Network Settings Tab

| Setting | Description |
|---------|-------------|
| Computer IP address | Local IP of the PC |
| Port | Pre-set to the controller's default |
| Target IP | IP address of the controller |
| Target Port | Pre-set to the controller's default |

### Drive Tab

| Feature | Description |
|---------|-------------|
| Status bits | Overview of controller status |
| Actual position / current | Live readings |
| SwitchOn | Toggle bit 0 to change state machine |
| Homing | Perform homing operation |
| Ack Errors | Acknowledge drive errors |
| Move Absolute | Move to an absolute position |
| Establish / Close Connection | Connect and disconnect |
| DLL error display | Status bar shows DLL error messages |

**Menu:** Edit → change software language. Help → library description, version number, check for updates, DLL installation path.

---

## Appendix IV: Time Stamp Integration

Windows is non-deterministic, so time-based measurements are not recommended. For applications where this is acceptable, the drive internal clock (1 ms resolution) can be added to a monitoring channel via UPID `82h` (Operating Sub Hours).

After configuring monitoring channel 1 to UPID `82h`, use the following DLL functions to access timestamped data:

- `getActualPosWithTimestamp`
- `getDemandPosWithTimestamp`
- `getDemandCurrentWithTimestamp`
- `getMonitoringChannelWithTimestamp`
- PseudoScope functions

> **Configuration:** `Monitoring Channels → Channel 1 UPID` must be set to `82h`.

---

## Appendix V: Manual Installation of LinUDP DLL Driver

**Prerequisites:** Microsoft .NET 4.0 Framework or newer. Administrator rights required.

> This registration is required for use with Visual Studio or Excel. **Not required** for LabVIEW.

### VI.I DLL Installation for Windows 7

1. Open Control Panel → Programs → check for **.NET Framework 4.0** or higher.
   - If not installed, download from: http://www.microsoft.com/en-us/download/details.aspx?id=17718

2. Copy `LinUDP.dll` to `C:\Windows\System32\` (64-bit: also `System32`).

3. Run **Regedit** and navigate to:
   ```
   HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System
   ```

4. Modify these registry entries:
   - `ConsentPromptBehaviorAdmin` → set to `0`
   - `EnableLUA` → set to `0`
   
   (Note down original values first, then reboot.)

5. Open a command window and navigate to the .NET folder:
   - 32-bit: `C:\Windows\Microsoft.NET\Framework\V4.xxxxx`
   - 64-bit: `C:\Windows\Microsoft.NET\Framework64\V4.xxxxx`

6. Register the DLL:
   ```
   regasm c:\windows\system32\LinUDP.dll /tlb:LinUDP.dll /codebase
   ```

7. Restore the registry settings (`ConsentPromptBehaviorAdmin` and `EnableLUA`) to their original values, then reboot.

**Adding the library reference in Excel 2013:**
1. Activate the Developer tab in Excel.
2. Click **Visual Basic**.
3. In the VB editor, click **Tools → References**.
4. Select **LinUDP Driver** from the list and click **OK**.

> **Note:** The reference is per-document. New Excel files require the reference to be added again. The LinUDP driver must be separately installed on each computer.

### V.II DLL Installation for Windows XP or Older

1. Open Control Panel → Add or Remove Programs → check for **.NET Framework 4.0** or higher.
   - If not installed, download from: http://www.microsoft.com/en-us/download/details.aspx?id=17718

2. Copy `LinUDP.dll` to `C:\Windows\System32\`.

3. Open a command window and navigate to the .NET Framework folder.

4. Register the DLL:
   ```
   regasm c:\windows\system32\LinUDP.dll /tlb:LinUDP.tlb /codebase
   ```

**Adding the library reference in Excel (2000 or newer):**
1. Click **Extras → Macro → Visual Basic Editor**.
2. In the VB editor, click **Extras → References**.
3. Select **LinUDP Driver** and click **OK**.

---

## Contact

**NTI AG (Switzerland)**
Bodenaeckerstr. 2, CH-8957 Spreitenbach

| | |
|---|---|
| Sales & Administration | +41-(0)56-419 91 91 / office@linmot.com |
| Tech. Support | +41-(0)56-544 71 00 / support@linmot.com |
| Tech. Support (Skype) | skype:support.linmot |
| Fax | +41-(0)56-419 91 92 |
| Web | http://www.linmot.com/ |

**LinMot USA Inc.**
N1922 State Road 120, Unit 1, Lake Geneva, WI 53147 USA

| | |
|---|---|
| Sales & Administration | 877-546-3270 / 262-743-2555 |
| Tech. Support | 877-804-0718 / 262-743-1284 |
| Fax | 800-463-8708 / 262-723-6688 |
| E-Mail | usasales@linmot.com |
| Web | http://www.linmot-usa.com/ |

Visit http://www.linmot.com to find the nearest distributor.
