#region Assembly LinUDP, Version=2.1.1.0, Culture=neutral, PublicKeyToken=5cdb96f0bb61122b
// location unknown
// Decompiled with ICSharpCode.Decompiler 8.2.0.7535
#endregion

using System;
using System.Collections;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Net;
using System.Net.NetworkInformation;
using System.Net.Sockets;
using System.Reflection;
using System.Runtime.CompilerServices;
using System.Runtime.InteropServices;
using System.Text;
using System.Threading;
using System.Xml.Serialization;
using LinUDP.My;
using Microsoft.VisualBasic;
using Microsoft.VisualBasic.CompilerServices;

namespace LinUDP;

[Guid("E677D1A6-FE43-45F6-8CCB-6E5C1318BD90")]
[ClassInterface(ClassInterfaceType.None)]
[ComClass("E677D1A6-FE43-45F6-8CCB-6E5C1318BD90", "D79AE6E7-58C7-45A0-BDBD-EAD3F627E252", "F3E87F99-9003-43CD-902B-93746216C38C")]
public class ACI : ACI._ACI, IDisposable
{
    public struct TimestampData
    {
        public double value;

        public long Timestamp;
    }

    public struct TimestampMonitoring
    {
        public long value;

        public long Timestamp;
    }

    public struct TimestampDataUTC
    {
        public double value;

        public DateTime Timestamp;
    }

    public struct TimestampMonitoringUTC
    {
        public long value;

        public DateTime Timestamp;
    }

    public struct PseudoScopeEntry
    {
        public double ActualPos;

        public double DemandPos;

        public double Current;

        public long TimeStamp;

        public long MonitoringChannel2;

        public long MonitoringChannel3;

        public long MonitoringChannel4;
    }

    [Serializable]
    public struct UPID_List
    {
        public ushort UPID;

        public long RawDataValue;
    }

    private struct ReturnValueUPIDList
    {
        public long value;

        public ushort UPID;

        public byte State;
    }

    [Serializable]
    public struct CTEntry
    {
        public List<byte> data;
    }

    [Serializable]
    public struct CommandTableStructure
    {
        public List<CTEntry> Entries;

        public List<byte> CTPresenceList;
    }

    [Serializable]
    private struct CTEntryDataDefinition
    {
        public byte[] Data;

        public ushort Status;
    }

    [Serializable]
    public struct CurveDataDefinition
    {
        public ushort CurveID;

        public ushort SetpointCount;

        public string CurveName;

        public byte Xcode;

        public byte Ycode;

        public uint XLength;

        public ushort XDim;

        public ushort YDim;

        public int[] Setpoints;
    }

    [Serializable]
    public struct DriveData
    {
        public uint HashVaue;

        public List<UPID_List> Parameters;

        public List<CurveDataDefinition> Curves;

        public CommandTableStructure CommandTable;
    }

    internal enum DriveType : ushort
    {
        Unknown = 0,
        C1250LUXC1S = 5025,
        C1250IPXC1S = 4913,
        C1250LUXC0S = 4257,
        C1250IPXC0S = 4145,
        C1450LUVS1S000 = 7585,
        C1450LUVS0S000 = 7073,
        C1450IPVS1S000 = 7473,
        C1450IPVS0S000 = 6961,
        E1250IPUC = 2337,
        E1250LUUC = 2081,
        E1450LUQN1SV2 = 6817,
        E1450IPQN1SV2 = 6705,
        E1450LUQN0SV2 = 6561,
        E1450IPQN0SV2 = 6449
    }

    internal enum FirmwareVersion : ushort
    {
        Unknown = 0,
        LT66 = 1542,
        LT67 = 1543,
        LT68 = 1544,
        LT69 = 1545,
        LT610 = 1546
    }

    public enum StateMachineStates : ushort
    {
        NotReadyToSwitchOn = 0,
        SwitchOnDisabled = 1,
        ReadyToSwitchOn = 2,
        SetupError = 3,
        GeneralError = 4,
        HWTests = 5,
        ReadyToOperate = 6,
        OperationEnabled = 8,
        Homing = 9,
        HomingFinished = 90,
        ClearanceCheck = 10,
        ClearanceCheckFinished = 100,
        GoingToInitialPosition = 11,
        GoingToInitialPositionFinished = 110,
        Aborting = 12,
        Freezing = 13,
        QuickStopOnError = 14,
        GoingToPosition = 15,
        GoingToPositionFinished = 150,
        JoggingPlus = 16,
        Jogging_MovingPositive = 161,
        JoggingPlusFinished = 162,
        JoggingMinus = 17,
        Jogging_MovingNegative = 171,
        JoggingMinusFinished = 172,
        Linearizing = 18,
        PhaseSearch = 19,
        SpecialMode = 20,
        BrakeDelay = 21
    }

    private struct MCI
    {
        public byte CommandHeaderMasterID;

        public byte CommandHeaderSubID;

        public byte CountNibble;

        public byte CommandParameter1Low;

        public byte CommandParameter1High;

        public byte CommandParameter1LowLow;

        public byte CommandParameter1HighHigh;

        public byte CommandParameter2HighHigh;

        public byte CommandParameter2Low;

        public byte CommandParameter2High;

        public byte CommandParameter2LowLow;

        public byte CommandParameter3HighHigh;

        public byte CommandParameter3Low;

        public byte CommandParameter3High;

        public byte CommandParameter3LowLow;

        public byte CommandParameter4HighHigh;

        public byte CommandParameter4Low;

        public byte CommandParameter4High;

        public byte CommandParameter4LowLow;

        public byte CommandParameter5HighHigh;

        public byte CommandParameter5Low;

        public byte CommandParameter5High;

        public byte CommandParameter5LowLow;

        public byte CommandParameter6HighHigh;

        public byte CommandParameter6Low;

        public byte CommandParameter6High;

        public byte CommandParameter6LowLow;

        public byte CommandParameter7HighHigh;

        public byte CommandParameter7Low;

        public byte CommandParameter7High;

        public byte CommandParameter7LowLow;
    }

    private struct DriveFeedbackData
    {
        public DriveType DriveHardware;

        public byte ControlWordLow;

        public byte ControlWordHigh;

        public byte StatusWordHigh;

        public byte StatusWordLow;

        public byte StateVarHigh;

        public byte StateVarLow;

        public byte StateVarHighOld;

        public byte StateVarLowOld;

        public double ActualPos;

        public double DemandPos;

        public double Current;

        public byte WarnWordHigh;

        public byte WarnWordLow;

        public byte ErrorCodeHigh;

        public byte ErrorCodeLow;

        public long MonitoringChannel1;

        public long MonitoringChannel2;

        public long MonitoringChannel3;

        public long MonitoringChannel4;

        public bool ToggleJogP;

        public bool ToggleJogM;

        public bool Done;

        public bool FatalError;

        public MCI MotionCommandInterface;

        public bool UpdateTimerRunning;

        public byte RealTimeConfigID;

        public byte RealTimeConfigCommandCount;

        public byte RealTimeConfigArgs1High;

        public byte RealTimeConfigArgs1Low;

        public byte RealTimeConfigArgs2High;

        public byte RealTimeConfigArgs2Low;

        public byte RealTimeConfigArgs3High;

        public byte RealTimeConfigArgs3Low;

        public byte RealTimeConfigIDStatus;

        public byte RealTimeConfigStatusCommandCount;

        public byte RealTimeConfigStatusArgs1High;

        public byte RealTimeConfigStatusArgs1Low;

        public byte RealTimeConfigStatusArgs2High;

        public byte RealTimeConfigStatusArgs2Low;

        public byte RealTimeConfigStatusArgs3High;

        public byte RealTimeConfigStatusArgs3Low;

        public int RealTimeConfigCurveStatus;

        public string DLLErrorText;

        public int DLLError;

        public bool isRespondActual;

        public string SlaveIP;

        public string SlavePort;

        public byte[] SendData;

        public byte[] RecieveData;

        public bool isConnected;

        public long TimeStampSent;

        public CurveDataDefinition UploadedCurve;

        public long TimeStampReceive;

        public long TimeStampDifference;

        public List<PseudoScopeEntry> PseudoScopeData;

        public bool setTimeoutObservation;

        public int setSkipAmountResponsePackets;

        public DateTime OperationHours;

        public long OperationSubHoursOld;

        public uint SendCounter;

        public uint RecieveCounter;

        public FirmwareVersion FirmwareVersion;
    }

    private class CurveCall
    {
        public string TargetIP;

        public int CurveID;

        public int SetpointCount;

        public string CurveName;

        public uint XLength;

        public int XDim;

        public int YDim;

        public int[] Setpoints;

        public int Mode;

        public byte Xcode;

        public byte Ycode;
    }

    private class debugLogEntry
    {
        public string line;
    }

    private class DriveEntry
    {
        public string DriveIP;

        public string DriveMAC;

        public string DriveName;

        public string InstalledFW;

        public string Hardware;

        public string SerialNo;
    }

    private class HostAdapters
    {
        public string HostIP;

        public string HostMAC;

        public string NetMask;
    }

    [Guid("D79AE6E7-58C7-45A0-BDBD-EAD3F627E252")]
    [ComVisible(true)]
    public interface _ACI
    {
        [DispId(43)]
        int PseudoScopeSamples
        {
            [DispId(43)]
            get;
            [DispId(43)]
            set;
        }

        [DispId(1)]
        void Dispose();

        [DispId(2)]
        void Free();

        [DispId(3)]
        uint LongRAWDataToUINT32(long value);

        [DispId(4)]
        int LongRAWDataToINT32(long value);

        [DispId(5)]
        ushort LongRAWDataToUINT16(long value);

        [DispId(6)]
        short LongRAWDataToINT16(long value);

        [DispId(7)]
        void setHostMAC(string MAC, string Port);

        [DispId(8)]
        string getHostIP();

        [DispId(9)]
        string getVersion();

        [DispId(10)]
        string getDriveIP_byMAC(string MACAdress, string HostIP, string HostMetMask, string TargetPort = "49360");

        [DispId(11)]
        void CreateTargetAddressList();

        [DispId(12)]
        bool SetTargetAddressList(string TargetIP, string TargetPort);

        [DispId(13)]
        string SetTargetAddressListByMAC(string TargetMAC, string TargetPort);

        [DispId(14)]
        void ClearTargetAddressList();

        [DispId(15)]
        double getTimerCycle();

        [DispId(16)]
        bool setTimerCycle(uint Cycle);

        [DispId(17)]
        bool ActivateConnection(string LocalIP, string LocalPort);

        [DispId(18)]
        bool ActivateConnectionByMAC(string LocalMAC, string LocalPort);

        [DispId(19)]
        bool CloseConnection();

        [DispId(20)]
        bool isNetworkRunning();

        [DispId(21)]
        string getDLLError();

        [DispId(22)]
        bool clearDLLErrors();

        [DispId(23)]
        string getDriveType(string TargetIP);

        [DispId(24)]
        bool isConnected(string TargetIP);

        [DispId(25)]
        bool isResponseUpToDate(string TargetIP);

        [DispId(26)]
        bool isRealtimeConfigUpToDate(string TargetIP);

        [DispId(27)]
        long getDatagramCycleTime(string TargetIP);

        [DispId(28)]
        double getActualPos(string TargetIP);

        [DispId(29)]
        double getCurrent(string TargetIP);

        [DispId(30)]
        double getDemandPos(string TargetIP);

        [DispId(31)]
        long getMonitoringChannel1(string TargetIP);

        [DispId(32)]
        long getMonitoringChannel2(string TargetIP);

        [DispId(33)]
        long getMonitoringChannel3(string TargetIP);

        [DispId(34)]
        long getMonitoringChannel4(string TargetIP);

        [DispId(35)]
        TimestampData getActualPosWithTimestamp(string TargetIP);

        [DispId(36)]
        TimestampDataUTC getActualPosWithTimestampUTC(string TargetIP);

        [DispId(37)]
        TimestampData getDemandPosWithTimestamp(string TargetIP);

        [DispId(38)]
        TimestampDataUTC getDemandPosWithTimestampUTC(string TargetIP);

        [DispId(39)]
        TimestampData getDemandCurrentWithTimestamp(string TargetIP);

        [DispId(40)]
        TimestampDataUTC getDemandCurrentWithTimestampUTC(string TargetIP);

        [DispId(41)]
        TimestampMonitoring getMonitoringChannelWithTimestamp(string TargetIP, int Channel);

        [DispId(42)]
        TimestampMonitoringUTC getMonitoringChannelWithTimestampUTC(string TargetIP, int Channel);

        [DispId(44)]
        bool enablePseudoScopeTrace(string TargetIP);

        [DispId(45)]
        bool isPseudoScopeSampling();

        [DispId(46)]
        List<PseudoScopeEntry> getPseudoScopeSamples(string TargetIP);

        [DispId(47)]
        int getPseudoScopeProgress(string TargetIP);

        [DispId(48)]
        StateMachineStates getStateMachineState(string TargetIP);

        [DispId(49)]
        bool isSwitchOnActive(string TargetIP);

        [DispId(50)]
        bool isEventHandlerActive(string TargetIP);

        [DispId(51)]
        bool isSpecialMotionActive(string TargetIP);

        [DispId(52)]
        bool isInTargetPosition(string TargetIP);

        [DispId(53)]
        bool isHomed(string TargetIP);

        [DispId(54)]
        bool isFatalError(string TargetIP);

        [DispId(55)]
        bool isMotionActive(string TargetIP);

        [DispId(56)]
        bool isRangeIndicator1(string TargetIP);

        [DispId(57)]
        bool isRangeIndicator2(string TargetIP);

        [DispId(58)]
        bool isOperationEnable(string TargetIP);

        [DispId(59)]
        bool isEnableOperation(string TargetIP);

        [DispId(60)]
        bool isError(string TargetIP);

        [DispId(61)]
        bool isSafeVoltageEnable(string TargetIP);

        [DispId(62)]
        bool isQuickStop(string TargetIP);

        [DispId(63)]
        bool isSwitchOnLocked(string TargetIP);

        [DispId(64)]
        bool isWarning(string TargetIP);

        [DispId(65)]
        bool isNotReadyToSwitchOnSM(string TargetIP);

        [DispId(66)]
        bool isSwitchOnDisabledSM(string TargetIP);

        [DispId(67)]
        bool isReadyToSwitchOnSM(string TargetIP);

        [DispId(68)]
        bool isSetupErrorSM(string TargetIP);

        [DispId(69)]
        bool isErrorSM(string TargetIP);

        [DispId(70)]
        bool isHWTestsSM(string TargetIP);

        [DispId(71)]
        bool isReadyToOperateSM(string TargetIP);

        [DispId(72)]
        bool isOperationEnabledSM(string TargetIP);

        [DispId(73)]
        bool isHomingSM(string TargetIP);

        [DispId(74)]
        bool Active(string TargetIP);

        [DispId(75)]
        bool SwitchOn(string TargetIP);

        [DispId(76)]
        bool setSwitchOnBit(string TargetIP, bool State);

        [DispId(77)]
        bool Homing(string TargetIP);

        [DispId(78)]
        bool setHomingBit(string TargetIP, bool State);

        [DispId(79)]
        bool AckErrors(string TargetIP);

        [DispId(80)]
        bool SetErrorAcknowledgeBit(string TargetIP, bool State);

        [DispId(81)]
        bool JogPlus(string TargetIP);

        [DispId(82)]
        bool JogMinus(string TargetIP);

        [DispId(83)]
        bool setJogPlus(string TargetIP, bool state);

        [DispId(84)]
        bool setJogMinus(string TargetIP, bool state);

        [DispId(85)]
        bool setBit0(string TargetIP, bool State);

        [DispId(86)]
        bool setBit1(string TargetIP, bool State);

        [DispId(87)]
        bool setBit2(string TargetIP, bool State);

        [DispId(88)]
        bool setBit3(string TargetIP, bool State);

        [DispId(89)]
        bool setBit4(string TargetIP, bool State);

        [DispId(90)]
        bool setBit5(string TargetIP, bool State);

        [DispId(91)]
        bool setBit6(string TargetIP, bool State);

        [DispId(92)]
        bool setBit7(string TargetIP, bool State);

        [DispId(93)]
        bool setBit8(string TargetIP, bool State);

        [DispId(94)]
        bool setBit9(string TargetIP, bool State);

        [DispId(95)]
        bool setBit10(string TargetIP, bool State);

        [DispId(96)]
        bool setBit11(string TargetIP, bool State);

        [DispId(97)]
        bool setBit12(string TargetIP, bool State);

        [DispId(98)]
        bool setBit13(string TargetIP, bool State);

        [DispId(99)]
        bool setBit14(string TargetIP, bool State);

        [DispId(100)]
        bool setBit15(string TargetIP, bool State);

        [DispId(101)]
        long LMcf_StartStopDefault(string TargetIP, int Mode);

        [DispId(102)]
        long getROM_ByUPID(string TargetIP, uint UPID);

        [DispId(103)]
        long getRAM_ByUPID(string TargetIP, uint UPID);

        [DispId(104)]
        long getMinVal_ByUPID(string TargetIP, uint UPID);

        [DispId(105)]
        long getMaxVal_ByUPID(string TargetIP, uint UPID);

        [DispId(106)]
        long getDefault_ByUPID(string TargetIP, uint UPID);

        [DispId(107)]
        long SetRAM_ByUPID(string TargetIP, uint UPID, long Value);

        [DispId(108)]
        long SetROM_ByUPID(string TargetIP, uint UPID, long Value);

        [DispId(109)]
        long SetRAM_ROM_ByUPID(string TargetIP, uint UPID, long Value);

        [DispId(110)]
        bool LMmt_GoToPosFromActPosAndActVel(string TargetIP, float Pos1, float MaxVel1, float Acc1, float Dec1);

        [DispId(111)]
        bool LMmt_MoveAbs(string TargetIP, float Pos1, float MaxVel1, float Acc1, float Dec1);

        [DispId(112)]
        bool LMmt_MoveRel(string TargetIP, float Pos1, float MaxVel1, float Acc1, float Dec1);

        [DispId(113)]
        bool LMmt_Stop(string TargetIP, float Decceleration);

        [DispId(114)]
        bool LMmt_WriteLivePar(string TargetIP, uint UPID, int UPIDValue);

        [DispId(115)]
        bool LMmt_VAJIGoToPos(string TargetIP, float Pos, float MaxVel, float Acc, float Dec, float Jerk);

        [DispId(116)]
        bool LMmt_IncrementActPosStartingWithDemVel0ResetI(string TargetIP, float Pos1, float MaxVel1, float Acc1, float Dec1);

        [DispId(117)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, short MCParaWord0);

        [DispId(118)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, short MCParaWord0, ushort MCParaWord1, ushort MCParaWord2);

        [DispId(119)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, short MCParaWord0, int MCParaWord1);

        [DispId(120)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, short MCParaWord0, ushort MCParaWord1, ushort MCParaWord2, ushort MCParaWord3);

        [DispId(121)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, ushort MCParaWord1);

        [DispId(122)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, short MCParaWord1);

        [DispId(123)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, uint MCParaWord1);

        [DispId(124)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, int MCParaWord1);

        [DispId(125)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0);

        [DispId(126)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, int MCParaWord1, int MCParaWord2);

        [DispId(127)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, ushort MCParaWord1, short MCParaWord2);

        [DispId(128)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, ushort MCParaWord1, int MCParaWord2);

        [DispId(129)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, int MCParaWord1, ushort MCParaWord2, short MCParaWord3);

        [DispId(130)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, int MCParaWord1, uint MCParaWord2, short MCParaWord3);

        [DispId(131)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, int MCParaWord1, int MCParaWord2, int MCParaWord3);

        [DispId(132)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, int MCParaWord1, short MCParaWord2, int MCParaWord3);

        [DispId(133)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, ushort MCParaWord1, int MCParaWord2, int MCParaWord3);

        [DispId(134)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, ushort MCParaWord1, ushort MCParaWord2, ushort MCParaWord3);

        [DispId(135)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, int MCParaWord1, int MCParaWord2, int MCParaWord3, int MCParaWord4);

        [DispId(136)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0);

        [DispId(137)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, int MCParaWord1, int MCParaWord2);

        [DispId(138)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, uint MCParaWord1, uint MCParaWord2);

        [DispId(139)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, ushort MCParaWord1, ushort MCParaWord2);

        [DispId(140)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, ushort MCParaWord1);

        [DispId(141)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, uint MCParaWord1);

        [DispId(142)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, int MCParaWord1);

        [DispId(143)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, uint MCParaWord1, uint MCParaWord2, uint MCParaWord3);

        [DispId(144)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, uint MCParaWord1, uint MCParaWord2, ushort MCParaWord3);

        [DispId(145)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, uint MCParaWord1, uint MCParaWord2, short MCParaWord3);

        [DispId(146)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, uint MCParaWord1, uint MCParaWord2, uint MCParaWord3, uint MCParaWord4);

        [DispId(147)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, uint MCParaWord1, uint MCParaWord2, short MCParaWord3, short MCParaWord4);

        [DispId(148)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, uint MCParaWord0, uint MCParaWord1);

        [DispId(149)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, uint MCParaWord0, uint MCParaWord1, uint MCParaWord2);

        [DispId(150)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, uint MCParaWord0);

        [DispId(151)]
        bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, ushort MCParaWord1, int MCParaWord2, int MCParaWord3, int MCParaWord4);

        [DispId(152)]
        bool LMav_Mod16BitCTPar(string TargetIP, int CTEntryID, int ParaOffset, int ParaValue);

        [DispId(153)]
        bool LMav_Mod32BitCTPar(string TargetIP, int CTEntryID, int ParaOffset, int ParaValue);

        [DispId(154)]
        bool LMmt_StartCTCommand(string TargetIP, uint CTEntryID);

        [DispId(155)]
        bool LMmt_ClearEventEvaluation(string TargetIP);

        [DispId(156)]
        bool LMav_RunCurve(string TargetIP, int CurveID, int CurveOffset, int TimeScale, int AmplitudeScale);

        [DispId(157)]
        bool LMav_MoveBestehorn(string TargetIP, float Position, float Velocity, float Acceleration, float Jerk);

        [DispId(158)]
        bool LMav_MoveBestehornRelative(string TargetIP, float Position, float Velocity, float Acceleration, float Jerk);

        [DispId(159)]
        bool LMav_MoveSin(string TargetIP, float Position, float Velocity, float Acceleration);

        [DispId(160)]
        bool LMav_MoveSinRelative(string TargetIP, float Position, float Velocity, float Acceleration);

        [DispId(161)]
        bool LMfc_ChangeTargetForce(string TargetIP, float TargetForce);

        [DispId(162)]
        bool LMfc_GoToPosForceCtrlHighLim(string TargetIP, float Position, float Velocity, float Acceleration, float ForceLimit, float TargetForce);

        [DispId(163)]
        bool LMfc_GoToPosForceCtrlLowLim(string TargetIP, float Position, float Velocity, float Acceleration, float ForceLimit, float TargetForce);

        [DispId(164)]
        bool LMfc_GoToPosRstForceCtrl(string TargetIP, float Position, float Velocity, float Acceleration, float Deceleration);

        [DispId(165)]
        bool LMfc_GoToPosRstForceCtrlSetI(string TargetIP, float Position, float Velocity, float Acceleration, float Deceleration);

        [DispId(166)]
        bool LMfc_IncrementActPosAndResetForceControlSetI(string TargetIP, float Position, float Velocity, float Acceleration, float Deceleration);

        [DispId(167)]
        bool LMfc_IncrementActPosWithHigherForceCtrlLimitAndTargetForce(string TargetIP, float PositionIncrement, float MaxVelocity, float Acceleration, float ForceLimit, float TargetForce);

        [DispId(168)]
        bool LMav_SetCurrentCommandMode(string TargetIP, int Current);

        [DispId(169)]
        bool LMav_ResetCurrentCommandMode(string TargetIP);

        [DispId(170)]
        int LMcf_getCurveProgress(string TargetIP);

        [DispId(171)]
        bool LMcf_LoadCurve(string TargetIP, int CurveID, int SetpointCount, string CurveName, byte Xcode, byte Ycode, uint XLength, int XDim, int YDim, int Mode, int[] Setpoints);

        [DispId(172)]
        bool LMcf_LoadCurve(string TargetIP, int Mode, CurveDataDefinition CurveData);

        [DispId(173)]
        bool LMcf_isCurveLoading(string TargetIP);

        [DispId(174)]
        bool LMcf_isCurveOnDrive(string TargetIP, ushort CurveID);

        [DispId(175)]
        int[] LMcf_getAllCurveID(string TargetIP);

        [DispId(176)]
        bool LMcf_StartUploadCurve(string TargetIP, ushort CurveID);

        [DispId(177)]
        bool LMcf_StartDownloadCurve(string TargetIP);

        [DispId(178)]
        bool LMcf_setDownloadCurveData(string TargetIP, CurveDataDefinition CurveData);

        [DispId(179)]
        CurveDataDefinition LMcf_getUploadedCurveData(string TargetIP);

        [DispId(180)]
        bool LMcf_DeleteAllCurvesInRAM(string TargetIP);

        [DispId(181)]
        bool Save_CommandTable(string TargetIP, string FilePath);

        [DispId(182)]
        bool Save_CommandTable(CommandTableStructure CT, string FilePath);

        [DispId(183)]
        bool Load_CommandTable(string TargetIP, string FilePath);

        [DispId(184)]
        CommandTableStructure Load_CommandTable(string FilePath);

        [DispId(185)]
        bool Save_DriveParameters(string TargetIP, string FilePath);

        [DispId(186)]
        bool Save_DriveParameters(List<UPID_List> Data, string FilePath);

        [DispId(187)]
        bool Load_DriveParameters(string TargetIP, string FilePath);

        [DispId(188)]
        List<UPID_List> Load_DriveParameters(string FilePAth);

        [DispId(189)]
        bool Save_Curves(string TargetIP, string FilePath);

        [DispId(190)]
        bool Load_Curves(string TargetIP, string FilePath);

        [DispId(191)]
        bool Save_DriveConfiguration(string TargetIP, string FilePath);

        [DispId(192)]
        bool Load_DriveConfiguration(string TargetIP, string FilePath);

        [DispId(193)]
        uint getDriveHash(string TargetIP);

        [DispId(194)]
        bool isDriveConfigurationSame(string TargetIP, string FilePath);

        [DispId(195)]
        string LMcf_GetErrorTxt(string TargetIP);

        [DispId(196)]
        int LMcf_GetErrorCode(string TargetIP);

        [DispId(197)]
        int LMcf_GetWarningCode(string TargetIP);

        [DispId(198)]
        string LMcf_GetWarningTxt(string TargetIP);

        [DispId(199)]
        bool isMastmerSlaveOperationEnabled(string TargetIP);

        [DispId(200)]
        bool MasterSlaveHoming(string TargetIP);

        [DispId(201)]
        List<UPID_List> LMcf_getUPIDList(string TargetIP, ushort StartUPID, ushort StopUPID);

        [DispId(202)]
        List<UPID_List> LMcf_getModified_UPIDList(string TargetIP, ushort StartUPID, ushort StopUPID);

        [DispId(203)]
        bool LMcf_setUPIDList(string TargetIP, List<UPID_List> UPIDList);

        [DispId(204)]
        CommandTableStructure LMcf_getCommandTableContent(string TargetIP);

        [DispId(205)]
        ushort LMcf_DeleteCommandTable_RAM(string TargetIP);

        [DispId(206)]
        ushort LMcf_WriteCommandTableToFLASH(string TargetIP);

        [DispId(207)]
        bool LMcf_setCommandTableContent(string TargetIP, CommandTableStructure CT);
    }

    public const string ClassId = "E677D1A6-FE43-45F6-8CCB-6E5C1318BD90";

    public const string InterfaceId = "D79AE6E7-58C7-45A0-BDBD-EAD3F627E252";

    public const string EventsId = "F3E87F99-9003-43CD-902B-93746216C38C";

    internal const double MinimumDriveResponse = 3.0;

    private static DriveFeedbackData[] Axis = new DriveFeedbackData[2];

    private UdpClient udp;

    private static string ACIError;

    private static int TargetIPListCount;

    private int SendTimerTick;

    private bool CurveThreadStart;

    private double TimerCycle;

    private bool ConnectionActive;

    private static List<DriveEntry> MAC_List;

    private static List<string> TargetList;

    private static List<string> TargetPortList;

    private static List<HostAdapters> HostIP;

    private static string HostIPAdress = "";

    private static string HostNetMask = "";

    private static string HostPort = "";

    private static string DrivePort = "";

    private CurveCall CurveCallData;

    private object TaskLock;

    private bool FirstStartDLL;

    private HPTimer TxTimer;

    private Thread ReadData_Thread;

    private bool ReadData_Taskrunning;

    private int PseudoScopePoints;

    private bool EnablePseudoScope;

    private DateTime LastARPScan;

    private Thread CurveT;

    public CommandTableStructure CTContent;

    private bool RunLogActive;

    private bool RunLogRunning;

    private Thread LogIt;

    private Thread LogFileWriter;

    private int LogInterval;

    private string LogFileName;

    private Queue LogLine;

    private string LogIP;

    private long oldTimestampLog;

    private long newTimestampLog;

    private long SendTimeStamp;

    private long ResendTimeCycle;

    private bool SendEnable;

    private long WaitSentTime;

    private string StreamIP;

    private bool SkipDebugFileContent;

    public int PseudoScopeSamples
    {
        get
        {
            return PseudoScopePoints;
        }
        set
        {
            PseudoScopePoints = value;
        }
    }

    protected virtual void Dispose(bool disposing)
    {
        if (disposing)
        {
            if (udp != null)
            {
                Thread.Sleep(250);
                udp.Close();
            }

            if (ConnectionActive)
            {
                CloseConnection();
            }

            GC.Collect();
            base.Finalize();
        }
    }

    //
    // Summary:
    //     close all internal threads and dispose the axis object, use if application should
    //     close
    public void Dispose()
    {
        Dispose(disposing: true);
        GC.SuppressFinalize(this);
    }

    void IDisposable.Dispose()
    {
        //ILSpy generated this explicit interface implementation from .override directive in Dispose
        this.Dispose();
    }

    void _ACI.Dispose()
    {
        //ILSpy generated this explicit interface implementation from .override directive in Dispose
        this.Dispose();
    }

    //
    // Summary:
    //     Create the axis object for the assigned IP and Port
    public ACI()
    {
        SendTimerTick = 2;
        TimerCycle = 5.0;
        ConnectionActive = false;
        CurveCallData = new CurveCall();
        TaskLock = RuntimeHelpers.GetObjectValue(new object());
        FirstStartDLL = true;
        PseudoScopePoints = 100;
        EnablePseudoScope = false;
        LastARPScan = default(DateTime);
        RunLogActive = false;
        RunLogRunning = false;
        LogInterval = 5;
        LogLine = new Queue();
        oldTimestampLog = 0L;
        newTimestampLog = 0L;
        SendTimeStamp = 0L;
        ResendTimeCycle = 1L;
        SendEnable = true;
        WaitSentTime = 0L;
        SkipDebugFileContent = false;
        MAC_List = new List<DriveEntry>();
        TargetList = new List<string>();
        TargetPortList = new List<string>();
        HostIP = new List<HostAdapters>();
        Thread.CurrentThread.Priority = ThreadPriority.Highest;
    }

    //
    // Summary:
    //     Depricated, please use Dispose()!
    public void Free()
    {
        Dispose();
    }

    void _ACI.Free()
    {
        //ILSpy generated this explicit interface implementation from .override directive in Free
        this.Free();
    }

    private void getDriveHardware()
    {
        ushort num = 0;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                long rAM_ByUPID = getRAM_ByUPID(Axis[i].SlaveIP, 8u);
                num = (ushort)rAM_ByUPID;
                if (num > 0)
                {
                    if (!Enum.TryParse<DriveType>(num.ToString(), out Axis[i].DriveHardware))
                    {
                        Axis[i].DriveHardware = DriveType.Unknown;
                    }
                }
                else
                {
                    Axis[i].DriveHardware = DriveType.Unknown;
                }

                long rAM_ByUPID2 = getRAM_ByUPID(Axis[i].SlaveIP, 7500u);
                byte[] bytes = BitConverter.GetBytes(rAM_ByUPID2);
                Axis[i].ControlWordLow = bytes[0];
                Axis[i].ControlWordHigh = bytes[1];
            }

            FirstStartDLL = false;
        }
    }

    private void getDriveFirmware()
    {
        ushort num = 0;
        ushort num2 = 0;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                long rAM_ByUPID = getRAM_ByUPID(Axis[i].SlaveIP, 52u);
                long rAM_ByUPID2 = getRAM_ByUPID(Axis[i].SlaveIP, 53u);
                num = (ushort)rAM_ByUPID;
                num2 = (ushort)rAM_ByUPID2;
                num = unchecked((ushort)(num << 8));
                num = (ushort)unchecked((uint)(num + num2));
                if (num > 0)
                {
                    if (!Enum.TryParse<FirmwareVersion>(num.ToString(), out Axis[i].FirmwareVersion))
                    {
                        Axis[i].FirmwareVersion = FirmwareVersion.Unknown;
                    }
                }
                else
                {
                    Axis[i].FirmwareVersion = FirmwareVersion.Unknown;
                }
            }
        }
    }

    private bool isTimeOut(long time)
    {
        long num = DateAndTime.Now.ToFileTimeUtc();
        long num2 = checked(time - num);
        if (time < num)
        {
            return true;
        }

        return false;
    }

    private long getTimeOutTime(long delay)
    {
        return checked(DateAndTime.Now.ToFileTimeUtc() + delay * 10000);
    }

    private string getTimestamp()
    {
        return DateAndTime.Now.Hour + ":" + DateAndTime.Now.Minute + ":" + DateAndTime.Now.Second + ":" + DateAndTime.Now.Millisecond;
    }

    private static string BytesToString(ref byte[] Data, int max)
    {
        checked
        {
            Data = (byte[])Utils.CopyArray(Data, new byte[(short)max + 1]);
            string text = "";
            short num = (short)max;
            for (short num2 = 0; num2 <= num; num2 = (short)unchecked(num2 + 1))
            {
                text = text + " " + Conversions.ToString(Strings.Chr(Data[num2]));
            }

            return text;
        }
    }

    private static string BytesToString2(byte[] bytes_Input)
    {
        StringBuilder stringBuilder = new StringBuilder(checked(bytes_Input.Length * 2));
        foreach (byte number in bytes_Input)
        {
            stringBuilder.Append(Conversion.Hex(number));
            stringBuilder.Append("-");
        }

        if (stringBuilder[0].Equals(" "))
        {
            stringBuilder.Remove(0, 1);
        }

        return stringBuilder.ToString();
    }

    private static string BytesToStringChar(byte[] bytes_Input)
    {
        StringBuilder stringBuilder = new StringBuilder(checked(bytes_Input.Length * 2));
        foreach (byte b in bytes_Input)
        {
            if (b >= 32)
            {
                stringBuilder.Append(Strings.ChrW(b));
            }
        }

        return stringBuilder.ToString();
    }

    [DllImport("iphlpapi.dll", CharSet = CharSet.Ansi, ExactSpelling = true, SetLastError = true)]
    private static extern int SendARP(uint DestIP, uint SrcIP, byte[] pMacAddr, ref int PhyAddrLen);

    private string GetMAC(string IPAddress)
    {
        IPAddress iPAddress = System.Net.IPAddress.Parse(IPAddress);
        byte[] array = new byte[7];
        int PhyAddrLen = array.Length;
        SendARP(BitConverter.ToUInt32(iPAddress.GetAddressBytes(), 0), 0u, array, ref PhyAddrLen);
        return BitConverter.ToString(array, 0, PhyAddrLen);
    }

    private string getHostNetMask(string IPAddr)
    {
        IPAddress comparand = IPAddress.Parse(IPAddr);
        string result = "";
        NetworkInterface[] allNetworkInterfaces = NetworkInterface.GetAllNetworkInterfaces();
        foreach (NetworkInterface networkInterface in allNetworkInterfaces)
        {
            IPInterfaceProperties iPProperties = networkInterface.GetIPProperties();
            foreach (UnicastIPAddressInformation unicastAddress in iPProperties.UnicastAddresses)
            {
                if (unicastAddress.Address.Equals(comparand))
                {
                    result = unicastAddress.IPv4Mask.ToString();
                }
            }
        }

        return result;
    }

    public uint LongRAWDataToUINT32(long value)
    {
        byte[] bytes = BitConverter.GetBytes(value);
        return BitConverter.ToUInt32(bytes, 0);
    }

    uint _ACI.LongRAWDataToUINT32(long value)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LongRAWDataToUINT32
        return this.LongRAWDataToUINT32(value);
    }

    public int LongRAWDataToINT32(long value)
    {
        byte[] bytes = BitConverter.GetBytes(value);
        return BitConverter.ToInt32(bytes, 0);
    }

    int _ACI.LongRAWDataToINT32(long value)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LongRAWDataToINT32
        return this.LongRAWDataToINT32(value);
    }

    public ushort LongRAWDataToUINT16(long value)
    {
        byte[] bytes = BitConverter.GetBytes(value);
        return BitConverter.ToUInt16(bytes, 0);
    }

    ushort _ACI.LongRAWDataToUINT16(long value)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LongRAWDataToUINT16
        return this.LongRAWDataToUINT16(value);
    }

    public short LongRAWDataToINT16(long value)
    {
        byte[] bytes = BitConverter.GetBytes(value);
        return BitConverter.ToInt16(bytes, 0);
    }

    short _ACI.LongRAWDataToINT16(long value)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LongRAWDataToINT16
        return this.LongRAWDataToINT16(value);
    }

    //
    // Summary:
    //     Request relevant data for all installed network adapters with IPv4 addresses
    //
    //
    // Returns:
    //     Returns number of installed network adapters
    private int GetHostIPAdapters()
    {
        string hostName = Dns.GetHostName();
        HostIP.Clear();
        IPAddress[] addressList = Dns.GetHostEntry(hostName).AddressList;
        foreach (IPAddress iPAddress in addressList)
        {
            if (!(iPAddress.IsIPv6LinkLocal | iPAddress.IsIPv6Multicast | iPAddress.IsIPv6SiteLocal | iPAddress.IsIPv6Teredo))
            {
                HostAdapters hostAdapters = new HostAdapters();
                hostAdapters.HostIP = iPAddress.ToString();
                byte[] addressBytes = iPAddress.GetAddressBytes();
                hostAdapters.HostMAC = GetMAC(hostAdapters.HostIP);
                hostAdapters.NetMask = getHostNetMask(hostAdapters.HostIP);
                HostIP.Add(hostAdapters);
            }
        }

        return checked(HostIP.Count - 1);
    }

    //
    // Summary:
    //     Get the IP address of the host network adapter by given MAC address
    //
    // Parameters:
    //   MAC:
    //     MAC address in form of "xx-xx-xx-xx-xx-xx"
    //
    // Returns:
    //     IP address used by adapter with given MAC. Empty, if MAC is unknown!
    private string getHostIPbyHostMAC(string MAC)
    {
        string result = "";
        GetHostIPAdapters();
        foreach (HostAdapters item in HostIP)
        {
            if (item.HostMAC.Equals(MAC.ToUpper()))
            {
                result = item.HostIP;
                HostNetMask = item.NetMask;
            }
        }

        return result;
    }

    //
    // Summary:
    //     Set the Host Computer MAC address for searching IP
    //
    // Parameters:
    //   MAC:
    //     MAC address as String like "C8-D3-A3-02-F7-7B"
    public void setHostMAC(string MAC, string Port)
    {
        HostIPAdress = getHostIPbyHostMAC(MAC);
        HostPort = Port;
    }

    void _ACI.setHostMAC(string MAC, string Port)
    {
        //ILSpy generated this explicit interface implementation from .override directive in setHostMAC
        this.setHostMAC(MAC, Port);
    }

    //
    // Summary:
    //     Get Host IP, if empty, no Host MAC was set!
    //
    // Returns:
    //     Host IP address as string
    public string getHostIP()
    {
        return HostIPAdress;
    }

    string _ACI.getHostIP()
    {
        //ILSpy generated this explicit interface implementation from .override directive in getHostIP
        return this.getHostIP();
    }

    //
    // Summary:
    //     Get the version number of the DLL
    //
    // Returns:
    //     Version as String
    public string getVersion()
    {
        Assembly.GetExecutingAssembly().GetName().Version.ToString();
        return Assembly.GetExecutingAssembly().GetName().Version.ToString();
    }

    string _ACI.getVersion()
    {
        //ILSpy generated this explicit interface implementation from .override directive in getVersion
        return this.getVersion();
    }

    private long LongIntegerFromIP(string p_strIP)
    {
        string[] array = Strings.Split(p_strIP, ".");
        int num = Information.UBound(array);
        checked
        {
            long num2 = default(long);
            for (int i = 0; i <= num; i++)
            {
                num2 = (long)Math.Round((double)num2 + (double)Conversions.ToLong(array[i]) * Math.Pow(256.0, 3 - i));
            }

            return num2;
        }
    }

    private void ScanForDriveIP(string TargetPort)
    {
        if (Operators.CompareString(HostIPAdress, "", TextCompare: false) == 0)
        {
            return;
        }

        byte[] array = new byte[10] { 1, 0, 0, 0, 1, 0, 0, 0, 0, 0 };
        string[] array2 = HostNetMask.Split('.');
        byte b = Conversions.ToByte(array2[0]);
        byte b2 = Conversions.ToByte(array2[1]);
        byte b3 = Conversions.ToByte(array2[2]);
        byte b4 = Conversions.ToByte(array2[3]);
        Console.WriteLine(BitConverter.ToUInt32(new byte[4] { b4, b3, b2, b }, 0));
        checked
        {
            long num = 4294967295L - unchecked((long)BitConverter.ToUInt32(new byte[4] { b4, b3, b2, b }, 0));
            if (Operators.CompareString(HostPort, "", TextCompare: false) == 0)
            {
                HostPort = "41136";
            }

            IPAddress iPAddress = IPAddress.Parse(HostIPAdress);
            IPEndPoint localEP = new IPEndPoint(iPAddress, Conversions.ToInteger(HostPort));
            UdpClient udpClient = new UdpClient(localEP);
            IPAddress iPAddress2 = new IPAddress(BitConverter.ToUInt32(iPAddress.GetAddressBytes(), 0) & BitConverter.ToUInt32(new byte[4] { b, b2, b3, b4 }, 0));
            long num2 = num;
            for (long num3 = 0L; num3 <= num2; num3++)
            {
                if (Operators.CompareString(TargetPort, "", TextCompare: false) == 0)
                {
                    TargetPort = "49360";
                }

                byte[] bytes = BitConverter.GetBytes(num3);
                IPAddress iPAddress3 = new IPAddress(new byte[4]
                {
                    bytes[3],
                    bytes[2],
                    bytes[1],
                    bytes[0]
                });
                IPAddress iPAddress4 = new IPAddress(BitConverter.ToUInt32(iPAddress2.GetAddressBytes(), 0) + BitConverter.ToUInt32(iPAddress3.GetAddressBytes(), 0));
                string text = iPAddress4.ToString();
                string[] array3 = text.Split('.');
                if (!(array3[3].Equals("0") | array3[3].Equals("255")))
                {
                    udpClient.Send(array, array.Length, text, Conversions.ToInteger(TargetPort));
                    Thread.Sleep(10);
                }
            }

            udpClient.Close();
        }
    }

    private void ScanForDriveIP1(string TargetPort, string HostIP, string HostMetMask)
    {
        if (Operators.CompareString(HostIP, "", TextCompare: false) == 0)
        {
            return;
        }

        byte[] array = new byte[10] { 1, 0, 0, 0, 1, 0, 0, 0, 0, 0 };
        string[] array2 = HostMetMask.Split('.');
        byte b = Conversions.ToByte(array2[0]);
        byte b2 = Conversions.ToByte(array2[1]);
        byte b3 = Conversions.ToByte(array2[2]);
        byte b4 = Conversions.ToByte(array2[3]);
        checked
        {
            long num = 4294967295L - unchecked((long)BitConverter.ToUInt32(new byte[4] { b4, b3, b2, b }, 0));
            if (Operators.CompareString(HostPort, "", TextCompare: false) == 0)
            {
                HostPort = "41136";
            }

            IPAddress iPAddress = IPAddress.Parse(HostIP);
            IPEndPoint localEP = new IPEndPoint(iPAddress, Conversions.ToInteger(HostPort));
            UdpClient udpClient = new UdpClient(localEP);
            IPAddress iPAddress2 = new IPAddress(BitConverter.ToUInt32(iPAddress.GetAddressBytes(), 0) & BitConverter.ToUInt32(new byte[4] { b, b2, b3, b4 }, 0));
            long num2 = num;
            for (long num3 = 0L; num3 <= num2; num3++)
            {
                if (Operators.CompareString(TargetPort, "", TextCompare: false) == 0)
                {
                    TargetPort = "49360";
                }

                byte[] bytes = BitConverter.GetBytes(num3);
                IPAddress iPAddress3 = new IPAddress(new byte[4]
                {
                    bytes[3],
                    bytes[2],
                    bytes[1],
                    bytes[0]
                });
                IPAddress iPAddress4 = new IPAddress(BitConverter.ToUInt32(iPAddress2.GetAddressBytes(), 0) + BitConverter.ToUInt32(iPAddress3.GetAddressBytes(), 0));
                string text = iPAddress4.ToString();
                string[] array3 = text.Split('.');
                if (!(array3[3].Equals("0") | array3[3].Equals("255")))
                {
                    udpClient.Send(array, array.Length, text, Conversions.ToInteger(TargetPort));
                    Thread.Sleep(10);
                }
            }

            udpClient.Close();
        }
    }

    //
    // Parameters:
    //   MACAdress:
    //     MAC Address with - delimeters
    //
    //   TargetPort:
    //     Target port, default is 49360
    //
    //   HostIP:
    //     IP Address of the host network card
    //
    //   HostMetMask:
    //     Subnet mask from the host network card
    //
    // Returns:
    //     IP address of the drive with the given MAC
    public string getDriveIP_byMAC(string MACAdress, string HostIP, string HostMetMask, string TargetPort = "49360")
    {
        string result = "";
        if (!isNetworkRunning())
        {
            if (checked(LastARPScan.Ticks + new TimeSpan(0, 1, 0).Ticks) < DateTime.Now.Ticks)
            {
                ScanForDriveIP1(TargetPort, HostIP, HostMetMask);
                Thread.Sleep(1000);
                LastARPScan = DateTime.Now;
            }

            string fileName = "arp";
            string arguments = "-a";
            ProcessStartInfo processStartInfo = new ProcessStartInfo(fileName, arguments);
            string text = "";
            string text2 = "";
            processStartInfo.UseShellExecute = false;
            processStartInfo.RedirectStandardOutput = true;
            processStartInfo.CreateNoWindow = true;
            Process process = Process.Start(processStartInfo);
            string expression = process.StandardOutput.ReadToEnd();
            string[] array = Strings.Split(expression, "\r\n");
            string[] array2 = array;
            foreach (string text3 in array2)
            {
                if (Operators.CompareString(text3, null, TextCompare: false) == 0)
                {
                    continue;
                }

                string[] array3 = text3.Split(new char[1] { ' ' }, StringSplitOptions.RemoveEmptyEntries);
                if (Information.UBound(array3) == 2)
                {
                    text = array3[0].PadRight(20).ToUpper();
                    text2 = array3[1].PadRight(20).ToUpper();
                    if (Operators.CompareString(text2.Trim(' '), MACAdress.ToUpper(), TextCompare: false) == 0)
                    {
                        result = text.Trim(' ');
                        break;
                    }
                }
            }
        }

        return result;
    }

    string _ACI.getDriveIP_byMAC(string MACAdress, string HostIP, string HostMetMask, string TargetPort = "49360")
    {
        //ILSpy generated this explicit interface implementation from .override directive in getDriveIP_byMAC
        return this.getDriveIP_byMAC(MACAdress, HostIP, HostMetMask, TargetPort);
    }

    private string getDriveIP_byMAC1(string MACAdress, string TargetPort)
    {
        ScanForDriveIP(TargetPort);
        Thread.Sleep(1000);
        string result = "";
        string fileName = "arp";
        string arguments = "-a";
        ProcessStartInfo processStartInfo = new ProcessStartInfo(fileName, arguments);
        string text = "";
        string text2 = "";
        processStartInfo.UseShellExecute = false;
        processStartInfo.RedirectStandardOutput = true;
        processStartInfo.CreateNoWindow = true;
        Process process = Process.Start(processStartInfo);
        string expression = process.StandardOutput.ReadToEnd();
        string[] array = Strings.Split(expression, "\r\n");
        string[] array2 = array;
        foreach (string text3 in array2)
        {
            if (Operators.CompareString(text3, null, TextCompare: false) == 0)
            {
                continue;
            }

            string[] array3 = text3.Split(new char[1] { ' ' }, StringSplitOptions.RemoveEmptyEntries);
            if (Information.UBound(array3) == 2)
            {
                text = array3[0].PadRight(20).ToUpper();
                text2 = array3[1].PadRight(20).ToUpper();
                if (Operators.CompareString(text2.Trim(' '), MACAdress.ToUpper(), TextCompare: false) == 0)
                {
                    result = text.Trim(' ');
                    break;
                }
            }
        }

        return result;
    }

    private void ReadDataRadio()
    {
        if (ReadData_Taskrunning)
        {
            try
            {
                ReadData_Thread = new Thread(ReadResponseData);
                ReadData_Thread.Name = "CheckForDatagram";
                ReadData_Thread.SetApartmentState(ApartmentState.MTA);
                ReadData_Thread.IsBackground = true;
                ReadData_Thread.Start();
            }
            catch (Exception ex)
            {
                ProjectData.SetProjectError(ex);
                Exception ex2 = ex;
                Console.WriteLine("WaitOnResponse" + ex2.Message);
                ProjectData.ClearProjectError();
            }
        }
    }

    private void ReadResponseData()
    {
        checked
        {
            while (ReadData_Taskrunning)
            {
                if (udp == null)
                {
                    continue;
                }

                try
                {
                    IPEndPoint remoteEP = null;
                    byte[] recieveData = udp.Receive(ref remoteEP);
                    int targetIPListCount = TargetIPListCount;
                    for (int i = 0; i <= targetIPListCount; i++)
                    {
                        if ((Operators.CompareString(Axis[i].SlaveIP, remoteEP.Address.ToString(), TextCompare: false) == 0) & (Axis[i].setSkipAmountResponsePackets < 10))
                        {
                            object taskLock = TaskLock;
                            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                            bool lockTaken = false;
                            try
                            {
                                Monitor.Enter(taskLock, ref lockTaken);
                                Axis[i].RecieveData = recieveData;
                                byte[] recieveData2 = Axis[i].RecieveData;
                                Axis[i].StatusWordHigh = recieveData2[9];
                                Axis[i].StatusWordLow = recieveData2[8];
                                Axis[i].StateVarHigh = recieveData2[11];
                                Axis[i].StateVarLow = recieveData2[10];
                                byte[] value = new byte[4]
                                {
                                    recieveData2[12],
                                    recieveData2[13],
                                    recieveData2[14],
                                    recieveData2[15]
                                };
                                Axis[i].ActualPos = (double)BitConverter.ToInt32(value, 0) * 0.0001;
                                value = new byte[4]
                                {
                                    recieveData2[16],
                                    recieveData2[17],
                                    recieveData2[18],
                                    recieveData2[19]
                                };
                                Axis[i].DemandPos = (double)BitConverter.ToInt32(value, 0) * 0.0001;
                                value = new byte[2]
                                {
                                    recieveData2[20],
                                    recieveData2[21]
                                };
                                Axis[i].Current = (double)BitConverter.ToInt16(value, 0) * 0.001;
                                Axis[i].WarnWordHigh = recieveData2[22];
                                Axis[i].WarnWordLow = recieveData2[23];
                                Axis[i].ErrorCodeHigh = recieveData2[24];
                                Axis[i].ErrorCodeLow = recieveData2[25];
                                value = new byte[4]
                                {
                                    recieveData2[26],
                                    recieveData2[27],
                                    recieveData2[28],
                                    recieveData2[29]
                                };
                                Axis[i].MonitoringChannel1 = BitConverter.ToInt32(value, 0);
                                value = new byte[4]
                                {
                                    recieveData2[30],
                                    recieveData2[31],
                                    recieveData2[32],
                                    recieveData2[33]
                                };
                                Axis[i].MonitoringChannel2 = BitConverter.ToInt32(value, 0);
                                value = new byte[4]
                                {
                                    recieveData2[34],
                                    recieveData2[35],
                                    recieveData2[36],
                                    recieveData2[37]
                                };
                                Axis[i].MonitoringChannel3 = BitConverter.ToInt32(value, 0);
                                value = new byte[4]
                                {
                                    recieveData2[38],
                                    recieveData2[39],
                                    recieveData2[40],
                                    recieveData2[41]
                                };
                                Axis[i].MonitoringChannel4 = BitConverter.ToInt32(value, 0);
                                Axis[i].RealTimeConfigIDStatus = recieveData2[43];
                                Axis[i].RealTimeConfigStatusCommandCount = recieveData2[42];
                                Axis[i].RealTimeConfigStatusArgs1High = recieveData2[45];
                                Axis[i].RealTimeConfigStatusArgs1Low = recieveData2[44];
                                Axis[i].RealTimeConfigStatusArgs2High = recieveData2[47];
                                Axis[i].RealTimeConfigStatusArgs2Low = recieveData2[46];
                                Axis[i].RealTimeConfigStatusArgs3High = recieveData2[49];
                                Axis[i].RealTimeConfigStatusArgs3Low = recieveData2[48];
                                Axis[i].TimeStampReceive = DateAndTime.Now.ToFileTimeUtc();
                                Axis[i].TimeStampDifference = Axis[i].TimeStampReceive - Axis[i].TimeStampSent;
                                Axis[i].RecieveCounter = (uint)(unchecked((long)Axis[i].RecieveCounter) + 1L);
                                if (unchecked((long)Axis[i].RecieveCounter) > 1000L)
                                {
                                    Axis[i].RecieveCounter = 0u;
                                }

                                Axis[i].isRespondActual = true;
                                if (EnablePseudoScope)
                                {
                                    if (Axis[i].PseudoScopeData.Count != PseudoScopePoints)
                                    {
                                        PseudoScopeEntry item = default(PseudoScopeEntry);
                                        item.ActualPos = Axis[i].ActualPos;
                                        item.Current = Axis[i].Current;
                                        item.DemandPos = Axis[i].DemandPos;
                                        item.TimeStamp = Axis[i].MonitoringChannel1;
                                        item.MonitoringChannel2 = Axis[i].MonitoringChannel2;
                                        item.MonitoringChannel3 = Axis[i].MonitoringChannel3;
                                        item.MonitoringChannel4 = Axis[i].MonitoringChannel4;
                                        Axis[i].PseudoScopeData.Add(item);
                                    }
                                    else
                                    {
                                        EnablePseudoScope = false;
                                    }
                                }
                            }
                            finally
                            {
                                if (lockTaken)
                                {
                                    Monitor.Exit(taskLock);
                                }
                            }

                            continue;
                        }

                        object taskLock2 = TaskLock;
                        ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                        bool lockTaken2 = false;
                        try
                        {
                            Monitor.Enter(taskLock2, ref lockTaken2);
                            Axis[i].setSkipAmountResponsePackets = Axis[i].setSkipAmountResponsePackets - 1;
                            if (Axis[i].setSkipAmountResponsePackets < 0)
                            {
                                Axis[i].setSkipAmountResponsePackets = 0;
                            }
                        }
                        finally
                        {
                            if (lockTaken2)
                            {
                                Monitor.Exit(taskLock2);
                            }
                        }
                    }
                }
                catch (Exception ex)
                {
                    ProjectData.SetProjectError(ex);
                    Exception ex2 = ex;
                    ACIError = ex2.Message;
                    Console.WriteLine(ex2.Message + " " + ex2.Source);
                    ProjectData.ClearProjectError();
                }
            }
        }
    }

    private void SetTimer()
    {
        TxTimer = new HPTimer(TimerCycle);
        TxTimer.TimerTriggered += TxData;
        TxTimer.startTimer();
    }

    private void RemoveTimer()
    {
        TxTimer.stopTimer();
        TxTimer.dispose();
    }

    private void ChangeTimerInterval(double deltaT)
    {
        TxTimer.interval = deltaT;
    }

    private void TxData()
    {
        object taskLock = TaskLock;
        ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
        bool lockTaken = false;
        checked
        {
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                int targetIPListCount = TargetIPListCount;
                for (int i = 0; i <= targetIPListCount; i++)
                {
                    try
                    {
                        bool flag = false;
                        if (unchecked((int)Axis[i].FirmwareVersion) > 1545)
                        {
                            flag = true;
                        }

                        long num = Math.Abs(Axis[i].TimeStampReceive - Axis[i].TimeStampSent);
                        long num2 = (long)Math.Round(3000000.0 + TimerCycle * (double)(TargetIPListCount + 2) * 10000.0);
                        if (!flag)
                        {
                            num2 = 2 * num2;
                        }

                        if (CurveThreadStart)
                        {
                            num2 = (long)Math.Round(3000000.0 + TimerCycle * (double)(TargetIPListCount + 2) * 10000.0 + 100000000.0);
                        }

                        if (Axis[i].setTimeoutObservation | (Axis[i].setSkipAmountResponsePackets > 0))
                        {
                            continue;
                        }

                        if (flag)
                        {
                            if (num >= num2)
                            {
                                if (Axis[i].SendCounter == Axis[i].RecieveCounter)
                                {
                                    Axis[i].isConnected = false;
                                    Axis[i].DLLErrorText = "1: No response from Drive: " + Axis[i].SlaveIP + " after " + (double)num / 100000.0 + "ms , timeout limit:  " + (double)num2 / 100000.0 + "ms";
                                }
                                else if (Axis[i].SendCounter > unchecked((long)Axis[i].RecieveCounter) + 5L)
                                {
                                    Axis[i].isConnected = false;
                                    Axis[i].DLLErrorText = "1: No response from Drive: " + Axis[i].SlaveIP + " after 5 retries, difference: " + (double)num / 100000.0 + "ms , timeout limit:  " + (double)num2 / 100000.0 + "ms";
                                }
                            }
                            else
                            {
                                Axis[i].RecieveCounter = Axis[i].SendCounter;
                                Axis[i].isConnected = true;
                            }
                        }
                        else if (num >= num2)
                        {
                            Axis[i].isConnected = false;
                            Axis[i].DLLErrorText = "1: No response from Drive: " + Axis[i].SlaveIP + " after " + (double)num / 100000.0 + "ms , timeout limit:  " + (double)num2 / 100000.0 + "ms";
                        }
                        else
                        {
                            Axis[i].isConnected = true;
                        }
                    }
                    catch (Exception ex)
                    {
                        ProjectData.SetProjectError(ex);
                        Exception ex2 = ex;
                        Console.WriteLine(ex2.Message.ToString() + " " + ex2.Source.ToString());
                        ProjectData.ClearProjectError();
                    }
                }
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }

            int targetIPListCount2 = TargetIPListCount;
            for (int j = 0; j <= targetIPListCount2; j++)
            {
                object taskLock2 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                bool lockTaken2 = false;
                try
                {
                    Monitor.Enter(taskLock2, ref lockTaken2);
                    byte[] sendData = ((!FirstStartDLL) ? new byte[50]
                    {
                        7,
                        0,
                        0,
                        0,
                        255,
                        1,
                        0,
                        0,
                        Axis[j].ControlWordLow,
                        Axis[j].ControlWordHigh,
                        Axis[j].MotionCommandInterface.CommandHeaderSubID,
                        Axis[j].MotionCommandInterface.CommandHeaderMasterID,
                        Axis[j].MotionCommandInterface.CommandParameter1LowLow,
                        Axis[j].MotionCommandInterface.CommandParameter1Low,
                        Axis[j].MotionCommandInterface.CommandParameter1High,
                        Axis[j].MotionCommandInterface.CommandParameter1HighHigh,
                        Axis[j].MotionCommandInterface.CommandParameter2LowLow,
                        Axis[j].MotionCommandInterface.CommandParameter2Low,
                        Axis[j].MotionCommandInterface.CommandParameter2High,
                        Axis[j].MotionCommandInterface.CommandParameter2HighHigh,
                        Axis[j].MotionCommandInterface.CommandParameter3LowLow,
                        Axis[j].MotionCommandInterface.CommandParameter3Low,
                        Axis[j].MotionCommandInterface.CommandParameter3High,
                        Axis[j].MotionCommandInterface.CommandParameter3HighHigh,
                        Axis[j].MotionCommandInterface.CommandParameter4LowLow,
                        Axis[j].MotionCommandInterface.CommandParameter4Low,
                        Axis[j].MotionCommandInterface.CommandParameter4High,
                        Axis[j].MotionCommandInterface.CommandParameter4HighHigh,
                        Axis[j].MotionCommandInterface.CommandParameter5LowLow,
                        Axis[j].MotionCommandInterface.CommandParameter5Low,
                        Axis[j].MotionCommandInterface.CommandParameter5High,
                        Axis[j].MotionCommandInterface.CommandParameter5HighHigh,
                        Axis[j].MotionCommandInterface.CommandParameter6LowLow,
                        Axis[j].MotionCommandInterface.CommandParameter6Low,
                        Axis[j].MotionCommandInterface.CommandParameter6High,
                        Axis[j].MotionCommandInterface.CommandParameter6HighHigh,
                        Axis[j].MotionCommandInterface.CommandParameter7LowLow,
                        Axis[j].MotionCommandInterface.CommandParameter7Low,
                        Axis[j].MotionCommandInterface.CommandParameter7High,
                        Axis[j].MotionCommandInterface.CommandParameter7HighHigh,
                        0,
                        0,
                        Axis[j].RealTimeConfigCommandCount,
                        Axis[j].RealTimeConfigID,
                        Axis[j].RealTimeConfigArgs1Low,
                        Axis[j].RealTimeConfigArgs1High,
                        Axis[j].RealTimeConfigArgs2Low,
                        Axis[j].RealTimeConfigArgs2High,
                        Axis[j].RealTimeConfigArgs3Low,
                        Axis[j].RealTimeConfigArgs3High
                    } : new byte[48]
                    {
                        6,
                        0,
                        0,
                        0,
                        255,
                        1,
                        0,
                        0,
                        Axis[j].MotionCommandInterface.CommandHeaderSubID,
                        Axis[j].MotionCommandInterface.CommandHeaderMasterID,
                        Axis[j].MotionCommandInterface.CommandParameter1LowLow,
                        Axis[j].MotionCommandInterface.CommandParameter1Low,
                        Axis[j].MotionCommandInterface.CommandParameter1High,
                        Axis[j].MotionCommandInterface.CommandParameter1HighHigh,
                        Axis[j].MotionCommandInterface.CommandParameter2LowLow,
                        Axis[j].MotionCommandInterface.CommandParameter2Low,
                        Axis[j].MotionCommandInterface.CommandParameter2High,
                        Axis[j].MotionCommandInterface.CommandParameter2HighHigh,
                        Axis[j].MotionCommandInterface.CommandParameter3LowLow,
                        Axis[j].MotionCommandInterface.CommandParameter3Low,
                        Axis[j].MotionCommandInterface.CommandParameter3High,
                        Axis[j].MotionCommandInterface.CommandParameter3HighHigh,
                        Axis[j].MotionCommandInterface.CommandParameter4LowLow,
                        Axis[j].MotionCommandInterface.CommandParameter4Low,
                        Axis[j].MotionCommandInterface.CommandParameter4High,
                        Axis[j].MotionCommandInterface.CommandParameter4HighHigh,
                        Axis[j].MotionCommandInterface.CommandParameter5LowLow,
                        Axis[j].MotionCommandInterface.CommandParameter5Low,
                        Axis[j].MotionCommandInterface.CommandParameter5High,
                        Axis[j].MotionCommandInterface.CommandParameter5HighHigh,
                        Axis[j].MotionCommandInterface.CommandParameter6LowLow,
                        Axis[j].MotionCommandInterface.CommandParameter6Low,
                        Axis[j].MotionCommandInterface.CommandParameter6High,
                        Axis[j].MotionCommandInterface.CommandParameter6HighHigh,
                        Axis[j].MotionCommandInterface.CommandParameter7LowLow,
                        Axis[j].MotionCommandInterface.CommandParameter7Low,
                        Axis[j].MotionCommandInterface.CommandParameter7High,
                        Axis[j].MotionCommandInterface.CommandParameter7HighHigh,
                        0,
                        0,
                        Axis[j].RealTimeConfigCommandCount,
                        Axis[j].RealTimeConfigID,
                        Axis[j].RealTimeConfigArgs1Low,
                        Axis[j].RealTimeConfigArgs1High,
                        Axis[j].RealTimeConfigArgs2Low,
                        Axis[j].RealTimeConfigArgs2High,
                        Axis[j].RealTimeConfigArgs3Low,
                        Axis[j].RealTimeConfigArgs3High
                    });
                    Axis[j].SendData = sendData;
                    Axis[j].StateVarHighOld = Axis[j].StateVarHigh;
                    Axis[j].StateVarLowOld = Axis[j].StateVarLowOld;
                    Axis[j].TimeStampSent = DateAndTime.Now.ToFileTimeUtc();
                    Axis[j].SendCounter = (uint)(unchecked((long)Axis[j].SendCounter) + 1L);
                    if (unchecked((long)Axis[j].SendCounter) > 1000L)
                    {
                        Axis[j].SendCounter = 0u;
                    }

                    Axis[j].isRespondActual = false;
                    try
                    {
                        udp.Send(Axis[j].SendData, Axis[j].SendData.Length, Axis[j].SlaveIP, int.Parse(Axis[j].SlavePort));
                    }
                    catch (Exception ex3)
                    {
                        ProjectData.SetProjectError(ex3);
                        Exception ex4 = ex3;
                        Axis[j].DLLErrorText = "2: Transmission to Axis " + Axis[j].SlaveIP + " failed!";
                        ProjectData.ClearProjectError();
                    }
                    finally
                    {
                    }
                }
                finally
                {
                    if (lockTaken2)
                    {
                        Monitor.Exit(taskLock2);
                    }
                }
            }
        }
    }

    //
    // Summary:
    //     Obsolete, only for compatibility to old applications! Do not use anymore!
    public void CreateTargetAddressList()
    {
    }

    void _ACI.CreateTargetAddressList()
    {
        //ILSpy generated this explicit interface implementation from .override directive in CreateTargetAddressList
        this.CreateTargetAddressList();
    }

    //
    // Summary:
    //     Set the drive into axis communication
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the drive, which will be used.
    //
    //   TargetPort:
    //     Port value for use with LinUDP. If empty, default port 49360 for LinUDP is used
    //
    //
    // Returns:
    //     Returns true, if successful.
    public bool SetTargetAddressList(string TargetIP, string TargetPort)
    {
        bool result = false;
        if (Operators.CompareString(TargetPort, "", TextCompare: false) == 0)
        {
            TargetPort = "49360";
        }

        if (Operators.CompareString(TargetIP, "", TextCompare: false) != 0)
        {
            TargetList.Add(TargetIP);
            TargetPortList.Add(TargetPort);
            result = true;
        }

        return result;
    }

    bool _ACI.SetTargetAddressList(string TargetIP, string TargetPort)
    {
        //ILSpy generated this explicit interface implementation from .override directive in SetTargetAddressList
        return this.SetTargetAddressList(TargetIP, TargetPort);
    }

    //
    // Summary:
    //     Set the drive into axis communication
    //
    // Parameters:
    //   TargetMAC:
    //     MAC address of the drive, which will be used.
    //
    //   TargetPort:
    //     Port value for use with LinUDP. If empty, default port 49360 for LinUDP is used
    //
    //
    // Returns:
    //     Returns corresponding IP for use with the ACI methods later on. Return value
    //     is empty, if MAC does not exist!
    public string SetTargetAddressListByMAC(string TargetMAC, string TargetPort)
    {
        string result = "";
        string driveIP_byMAC = getDriveIP_byMAC1(TargetMAC, TargetPort);
        if (Operators.CompareString(TargetPort, "", TextCompare: false) == 0)
        {
            TargetPort = "49360";
        }

        if (Operators.CompareString(driveIP_byMAC, "", TextCompare: false) != 0)
        {
            result = driveIP_byMAC;
            TargetList.Add(driveIP_byMAC);
            TargetPortList.Add(TargetPort);
        }

        return result;
    }

    string _ACI.SetTargetAddressListByMAC(string TargetMAC, string TargetPort)
    {
        //ILSpy generated this explicit interface implementation from .override directive in SetTargetAddressListByMAC
        return this.SetTargetAddressListByMAC(TargetMAC, TargetPort);
    }

    //
    // Summary:
    //     Clear the internal drive list. Can be used, if not connected to set new IP addresses
    //     for reconnection
    public void ClearTargetAddressList()
    {
        TargetList.Clear();
        TargetPortList.Clear();
    }

    void _ACI.ClearTargetAddressList()
    {
        //ILSpy generated this explicit interface implementation from .override directive in ClearTargetAddressList
        this.ClearTargetAddressList();
    }

    //
    // Summary:
    //     Get timer value which is used for cyclic status data request
    //
    // Returns:
    //     Timer value in [ms]
    public double getTimerCycle()
    {
        return TimerCycle;
    }

    double _ACI.getTimerCycle()
    {
        //ILSpy generated this explicit interface implementation from .override directive in getTimerCycle
        return this.getTimerCycle();
    }

    //
    // Summary:
    //     Set the time span between cyclic polling. May be used, to reduce network traffic
    //     by increasing time value which can be read with "getTimerCycle". Recommended
    //     is 2ms
    //
    // Parameters:
    //   Cycle:
    //     New cycle time in [ms]
    //
    // Returns:
    //     True, if successful
    public bool setTimerCycle(uint Cycle)
    {
        bool flag = false;
        if (TxTimer != null)
        {
            if (!((double)Cycle < 3.0))
            {
                TimerCycle = Cycle;
                TxTimer.interval = TimerCycle;
                return true;
            }

            ACIError = "Error: Time span for data excahnge cycle below " + 3.0 + "!";
            return false;
        }

        ACIError = "Error: Value can changed only, after communication is established!";
        return false;
    }

    bool _ACI.setTimerCycle(uint Cycle)
    {
        //ILSpy generated this explicit interface implementation from .override directive in setTimerCycle
        return this.setTimerCycle(Cycle);
    }

    //
    // Summary:
    //     Start cyclic data exchange to drives
    //
    // Parameters:
    //   LocalIP:
    //     IP of Host, which should be used for communication
    //
    //   LocalPort:
    //     Host port which should be used, empty for default!
    //
    // Returns:
    //     Return true, if activation was successful
    public bool ActivateConnection(string LocalIP, string LocalPort)
    {
        bool result = false;
        if (Operators.CompareString(LocalPort, "", TextCompare: false) == 0)
        {
            LocalPort = "41136";
        }

        if (Operators.CompareString(LocalIP, "", TextCompare: false) == 0)
        {
            LocalIP = "127.0.0.1";
        }

        HostIPAdress = LocalIP;
        HostPort = LocalPort;
        checked
        {
            if (!ConnectionActive & (TargetList.Count - 1 > -1))
            {
                try
                {
                    IPAddress address = IPAddress.Parse(LocalIP);
                    IPEndPoint localEP = new IPEndPoint(address, int.Parse(LocalPort));
                    udp = new UdpClient(localEP);
                    TargetIPListCount = TargetList.Count - 1;
                    Axis = new DriveFeedbackData[TargetList.Count + 1];
                    int num = TargetList.Count - 1;
                    for (int i = 0; i <= num; i++)
                    {
                        DriveFeedbackData driveFeedbackData = default(DriveFeedbackData);
                        Axis[i] = default(DriveFeedbackData);
                        Axis[i].PseudoScopeData = new List<PseudoScopeEntry>();
                        Axis[i].SlaveIP = TargetList.ElementAt(i);
                        Axis[i].SlavePort = TargetPortList.ElementAt(i);
                        Axis[i].isConnected = true;
                        Axis[i].TimeStampReceive = DateAndTime.Now.ToFileTimeUtc();
                        Axis[i].TimeStampSent = DateAndTime.Now.ToFileTimeUtc();
                        Axis[i].OperationSubHoursOld = 0L;
                        Axis[i].OperationHours = DateTime.MinValue;
                        Axis[i].RealTimeConfigCommandCount = 1;
                        Axis[i].UploadedCurve = default(CurveDataDefinition);
                    }

                    double num2 = 3.0 - (double)(TargetList.Count + 1) * 0.2;
                    if (num2 > 0.0)
                    {
                        num2 = 3.0;
                    }

                    if (num2 < 0.0)
                    {
                        num2 = 0.0;
                    }

                    TimerCycle = (double)(TargetList.Count + 1) * 0.2 + num2;
                    if (TimerCycle < 3.0)
                    {
                        TimerCycle = 3.0;
                    }

                    WaitSentTime = (long)Math.Round(Math.Round(TimerCycle * 10000.0 / (double)Axis.Count()));
                    if (!ReadData_Taskrunning)
                    {
                        ReadData_Taskrunning = true;
                        ReadDataRadio();
                    }

                    SetTimer();
                    result = true;
                    ConnectionActive = true;
                    getDriveHardware();
                    getDriveFirmware();
                }
                catch (Exception ex)
                {
                    ProjectData.SetProjectError(ex);
                    Exception ex2 = ex;
                    ex2.Message.ToString();
                    ACIError = "Port is activated, please close before open!";
                    ProjectData.ClearProjectError();
                }
            }
            else
            {
                ACIError = "Port is activated, please close before open!";
            }

            return result;
        }
    }

    bool _ACI.ActivateConnection(string LocalIP, string LocalPort)
    {
        //ILSpy generated this explicit interface implementation from .override directive in ActivateConnection
        return this.ActivateConnection(LocalIP, LocalPort);
    }

    //
    // Summary:
    //     Start cyclic data exchange to drives
    //
    // Parameters:
    //   LocalMAC:
    //     MAC of Host, which should be used for communication
    //
    //   LocalPort:
    //     Host port which should be used, empty for default!
    //
    // Returns:
    //     Return true, if activation was successful
    public bool ActivateConnectionByMAC(string LocalMAC, string LocalPort)
    {
        bool result = false;
        if (Operators.CompareString(LocalPort, "", TextCompare: false) == 0)
        {
            LocalPort = "41136";
        }

        string hostIPbyHostMAC = getHostIPbyHostMAC(LocalMAC);
        if (Operators.CompareString(hostIPbyHostMAC, "", TextCompare: false) == 0)
        {
            return result;
        }

        HostIPAdress = hostIPbyHostMAC;
        HostPort = LocalPort;
        checked
        {
            if (!ConnectionActive & (TargetList.Count - 1 > -1))
            {
                try
                {
                    IPAddress address = IPAddress.Parse(hostIPbyHostMAC);
                    IPEndPoint localEP = new IPEndPoint(address, int.Parse(LocalPort));
                    udp = new UdpClient(localEP);
                    TargetIPListCount = TargetList.Count - 1;
                    Axis = new DriveFeedbackData[TargetList.Count + 1];
                    int num = TargetList.Count - 1;
                    for (int i = 0; i <= num; i++)
                    {
                        DriveFeedbackData driveFeedbackData = default(DriveFeedbackData);
                        Axis[i] = default(DriveFeedbackData);
                        Axis[i].PseudoScopeData = new List<PseudoScopeEntry>();
                        Axis[i].SlaveIP = TargetList.ElementAt(i);
                        Axis[i].SlavePort = TargetPortList.ElementAt(i);
                        Axis[i].isConnected = true;
                        Axis[i].TimeStampReceive = DateAndTime.Now.ToFileTimeUtc();
                        Axis[i].TimeStampSent = DateAndTime.Now.ToFileTimeUtc();
                        Axis[i].OperationSubHoursOld = 0L;
                        Axis[i].OperationHours = DateTime.MinValue;
                        Axis[i].RealTimeConfigCommandCount = 1;
                        Axis[i].UploadedCurve = default(CurveDataDefinition);
                    }

                    double num2 = 3.0 - (double)(TargetList.Count + 1) * 0.2;
                    if (num2 > 0.0)
                    {
                        num2 = 3.0;
                    }

                    if (num2 < 0.0)
                    {
                        num2 = 0.0;
                    }

                    TimerCycle = (double)(TargetList.Count + 1) * 0.2 + num2;
                    if (TimerCycle < 3.0)
                    {
                        TimerCycle = 3.0;
                    }

                    WaitSentTime = (long)Math.Round(Math.Round(TimerCycle * 10000.0 / (double)Axis.Count()));
                    if (!ReadData_Taskrunning)
                    {
                        ReadData_Taskrunning = true;
                        ReadDataRadio();
                    }

                    SetTimer();
                    result = true;
                    ConnectionActive = true;
                    getDriveHardware();
                    getDriveFirmware();
                }
                catch (Exception ex)
                {
                    ProjectData.SetProjectError(ex);
                    Exception ex2 = ex;
                    ex2.Message.ToString();
                    ACIError = "Port is activated, please close before open!";
                    ProjectData.ClearProjectError();
                }
            }
            else
            {
                ACIError = "Port is activated, please close before open!";
            }

            return result;
        }
    }

    bool _ACI.ActivateConnectionByMAC(string LocalMAC, string LocalPort)
    {
        //ILSpy generated this explicit interface implementation from .override directive in ActivateConnectionByMAC
        return this.ActivateConnectionByMAC(LocalMAC, LocalPort);
    }

    //
    // Summary:
    //     Stop cyclic data exchange to drives.
    //
    // Returns:
    //     Returns trus, if data exchange is stopped successful
    public bool CloseConnection()
    {
        bool result = false;
        if (ConnectionActive)
        {
            ReadData_Taskrunning = false;
            udp.Send(new byte[1], 1, HostIPAdress, int.Parse(HostPort));
            while (ReadData_Thread.IsAlive)
            {
            }

            Thread.Sleep(250);
            RemoveTimer();
            ConnectionActive = false;
            udp.Close();
            result = true;
        }

        return result;
    }

    bool _ACI.CloseConnection()
    {
        //ILSpy generated this explicit interface implementation from .override directive in CloseConnection
        return this.CloseConnection();
    }

    //
    // Summary:
    //     Check, if cyclic data exchange is running
    //
    // Returns:
    //     True, if cyclic data exchange is running
    public bool isNetworkRunning()
    {
        return ConnectionActive;
    }

    bool _ACI.isNetworkRunning()
    {
        //ILSpy generated this explicit interface implementation from .override directive in isNetworkRunning
        return this.isNetworkRunning();
    }

    //
    // Summary:
    //     Get Errors from DLL operation. This function deliver no drive related errors!
    //
    //
    // Returns:
    //     Return Error message. It is empty, if no error occurs
    public string getDLLError()
    {
        string text = "";
        text = "";
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].DLLErrorText, "", TextCompare: false) != 0)
            {
                text = text + " " + Axis[i].DLLErrorText;
            }
        }

        if (Operators.CompareString(ACIError, "", TextCompare: false) != 0)
        {
            text = text + " " + ACIError;
        }

        return text;
    }

    string _ACI.getDLLError()
    {
        //ILSpy generated this explicit interface implementation from .override directive in getDLLError
        return this.getDLLError();
    }

    public bool clearDLLErrors()
    {
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            Axis[i].DLLErrorText = "";
        }

        ACIError = "";
        return true;
    }

    bool _ACI.clearDLLErrors()
    {
        //ILSpy generated this explicit interface implementation from .override directive in clearDLLErrors
        return this.clearDLLErrors();
    }

    public string getDriveType(string TargetIP)
    {
        string result = "";
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = Axis[i].DriveHardware.ToString();
                break;
            }

            result = "";
        }

        return result;
    }

    string _ACI.getDriveType(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getDriveType
        return this.getDriveType(TargetIP);
    }

    //
    // Summary:
    //     check, if the drive is registered for cyclic data exchange
    //
    // Parameters:
    //   TargetIP:
    //     IP of drive, which should be registered
    //
    // Returns:
    //     True, if drive is registered for cyclic data exchange
    public bool isConnected(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = Axis[i].isConnected;
                break;
            }

            result = false;
        }

        return result;
    }

    bool _ACI.isConnected(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isConnected
        return this.isConnected(TargetIP);
    }

    //
    // Summary:
    //     Check the requested drive state is valid
    //
    // Parameters:
    //   TargetIP:
    //     IP of drive, which state should be up to date
    //
    // Returns:
    //     True, if drive state is up to date
    public bool isResponseUpToDate(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if ((Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0) & (Axis[i].MotionCommandInterface.CountNibble == (Axis[i].StateVarLow & 0xF)))
            {
                result = true;
                break;
            }

            result = false;
        }

        return result;
    }

    bool _ACI.isResponseUpToDate(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isResponseUpToDate
        return this.isResponseUpToDate(TargetIP);
    }

    //
    // Summary:
    //     Check the requested drive if parameter channel access is done
    //
    // Parameters:
    //   TargetIP:
    //     IP of drive
    //
    // Returns:
    //     True, if parameter execution has done
    public bool isRealtimeConfigUpToDate(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if ((Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0) & (Axis[i].RealTimeConfigCommandCount == Axis[i].RealTimeConfigStatusCommandCount))
            {
                result = true;
                break;
            }

            result = false;
        }

        return result;
    }

    bool _ACI.isRealtimeConfigUpToDate(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isRealtimeConfigUpToDate
        return this.isRealtimeConfigUpToDate(TargetIP);
    }

    //
    // Summary:
    //     Get the actual time between sending and recieving a datagram (DLL time)
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Difference time value [0.1 us]
    public long getDatagramCycleTime(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        long result = default(long);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = Axis[i].TimeStampDifference;
                break;
            }

            result = 0L;
        }

        return result;
    }

    long _ACI.getDatagramCycleTime(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getDatagramCycleTime
        return this.getDatagramCycleTime(TargetIP);
    }

    //
    // Summary:
    //     Get the actual slider position
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Position value [mm]
    public double getActualPos(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        double result = default(double);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = Axis[i].ActualPos;
                break;
            }

            result = 0.0;
        }

        return result;
    }

    double _ACI.getActualPos(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getActualPos
        return this.getActualPos(TargetIP);
    }

    //
    // Summary:
    //     Get the actual motor current
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Actual current [A]
    public double getCurrent(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        double result = default(double);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = Axis[i].Current;
                break;
            }

            result = 0.0;
        }

        return result;
    }

    double _ACI.getCurrent(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getCurrent
        return this.getCurrent(TargetIP);
    }

    //
    // Summary:
    //     Get the demand position
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Demand position [mm]
    public double getDemandPos(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        double result = default(double);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = Axis[i].DemandPos;
                break;
            }

            result = 0.0;
        }

        return result;
    }

    double _ACI.getDemandPos(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getDemandPos
        return this.getDemandPos(TargetIP);
    }

    //
    // Summary:
    //     Get the raw data value of monitoring channel 1, which must be configured with
    //     LinMot Talk
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Raw data value, must be scaled according to UPID scale (see LinMot Talk -> UPID)
    public long getMonitoringChannel1(string TargetIP)
    {
        long result = 0L;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = Axis[i].MonitoringChannel1;
            }
        }

        return result;
    }

    long _ACI.getMonitoringChannel1(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getMonitoringChannel1
        return this.getMonitoringChannel1(TargetIP);
    }

    //
    // Summary:
    //     Get the raw data value of monitoring channel 2, which must be configured with
    //     LinMot Talk
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Raw data value, must be scaled according to UPID scale (see LinMot Talk -> UPID)
    public long getMonitoringChannel2(string TargetIP)
    {
        long result = 0L;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = Axis[i].MonitoringChannel2;
            }
        }

        return result;
    }

    long _ACI.getMonitoringChannel2(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getMonitoringChannel2
        return this.getMonitoringChannel2(TargetIP);
    }

    //
    // Summary:
    //     Get the raw data value of monitoring channel 3, which must be configured with
    //     LinMot Talk
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Raw data value, must be scaled according to UPID scale (see LinMot Talk -> UPID)
    public long getMonitoringChannel3(string TargetIP)
    {
        long result = 0L;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = Axis[i].MonitoringChannel3;
            }
        }

        return result;
    }

    long _ACI.getMonitoringChannel3(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getMonitoringChannel3
        return this.getMonitoringChannel3(TargetIP);
    }

    //
    // Summary:
    //     Get the raw data value of monitoring channel 4, which must be configured with
    //     LinMot Talk
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Raw data value, must be scaled according to UPID scale (see LinMot Talk -> UPID)
    public long getMonitoringChannel4(string TargetIP)
    {
        long result = 0L;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = Axis[i].MonitoringChannel4;
            }
        }

        return result;
    }

    long _ACI.getMonitoringChannel4(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getMonitoringChannel4
        return this.getMonitoringChannel4(TargetIP);
    }

    //
    // Summary:
    //     Get actual position with drive timestamp, Monitoring Channel 1 must be configured
    //     to "Operating Sub Hours"¨!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Data structure contains actual position and timestamp
    public TimestampData getActualPosWithTimestamp(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        TimestampData result = default(TimestampData);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result.value = Axis[i].ActualPos;
                result.Timestamp = Axis[i].MonitoringChannel1;
                break;
            }
        }

        return result;
    }

    TimestampData _ACI.getActualPosWithTimestamp(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getActualPosWithTimestamp
        return this.getActualPosWithTimestamp(TargetIP);
    }

    //
    // Summary:
    //     Get actual position with drive timestamp in Coordinated Universal Time, Monitoring
    //     Channel 1 must be configured to "Operating Sub Hours"¨!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Data structure contains actual position and timestamp
    public TimestampDataUTC getActualPosWithTimestampUTC(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        checked
        {
            TimestampDataUTC result = default(TimestampDataUTC);
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                result.value = Axis[i].ActualPos;
                DateTime dateTime = default(DateTime);
                DateTime dateTime2 = Axis[i].OperationHours;
                long num = 0L;
                if (Axis[i].OperationHours.Ticks == 0)
                {
                    dateTime2 = DateTime.Now;
                }

                num = getMonitoringChannel1(TargetIP);
                long num2 = num - Axis[i].OperationSubHoursOld;
                if (num2 < 0)
                {
                    num2 += 3600000;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].OperationHours = dateTime2.AddMilliseconds(num2);
                    Axis[i].OperationSubHoursOld = num;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                dateTime = Axis[i].OperationHours;
                result.Timestamp = dateTime.ToUniversalTime();
                break;
            }

            return result;
        }
    }

    TimestampDataUTC _ACI.getActualPosWithTimestampUTC(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getActualPosWithTimestampUTC
        return this.getActualPosWithTimestampUTC(TargetIP);
    }

    //
    // Summary:
    //     Get demand position with drive timestamp, Monitoring Channel 1 must be configured
    //     to "Operating Sub Hours"¨!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Data structure contains demand position and timestamp
    public TimestampData getDemandPosWithTimestamp(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        TimestampData result = default(TimestampData);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result.value = Axis[i].DemandPos;
                result.Timestamp = Axis[i].MonitoringChannel1;
                break;
            }
        }

        return result;
    }

    TimestampData _ACI.getDemandPosWithTimestamp(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getDemandPosWithTimestamp
        return this.getDemandPosWithTimestamp(TargetIP);
    }

    //
    // Summary:
    //     Get demand position with drive timestamp in Coordinated Universal Time, Monitoring
    //     Channel 1 must be configured to "Operating Sub Hours"¨!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Data structure contains demand position and timestamp
    public TimestampDataUTC getDemandPosWithTimestampUTC(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        checked
        {
            TimestampDataUTC result = default(TimestampDataUTC);
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                result.value = Axis[i].DemandPos;
                DateTime dateTime = default(DateTime);
                DateTime dateTime2 = Axis[i].OperationHours;
                long num = 0L;
                if (Axis[i].OperationHours.Ticks == 0)
                {
                    dateTime2 = DateTime.Now;
                }

                num = getMonitoringChannel1(TargetIP);
                long num2 = num - Axis[i].OperationSubHoursOld;
                if (num2 < 0)
                {
                    num2 += 3600000;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].OperationHours = dateTime2.AddMilliseconds(num2);
                    Axis[i].OperationSubHoursOld = num;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                dateTime = Axis[i].OperationHours;
                result.Timestamp = dateTime.ToUniversalTime();
                break;
            }

            return result;
        }
    }

    TimestampDataUTC _ACI.getDemandPosWithTimestampUTC(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getDemandPosWithTimestampUTC
        return this.getDemandPosWithTimestampUTC(TargetIP);
    }

    //
    // Summary:
    //     Get current with drive timestamp, Monitoring Channel 1 must be configured to
    //     "Operating Sub Hours"¨!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Data structure contains current and timestamp
    public TimestampData getDemandCurrentWithTimestamp(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        TimestampData result = default(TimestampData);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result.value = Axis[i].Current;
                result.Timestamp = Axis[i].MonitoringChannel1;
                break;
            }
        }

        return result;
    }

    TimestampData _ACI.getDemandCurrentWithTimestamp(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getDemandCurrentWithTimestamp
        return this.getDemandCurrentWithTimestamp(TargetIP);
    }

    //
    // Summary:
    //     Get current with drive timestamp in Coordinated Universal Time, Monitoring Channel
    //     1 must be configured to "Operating Sub Hours"¨!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Data structure contains current and timestamp
    public TimestampDataUTC getDemandCurrentWithTimestampUTC(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        checked
        {
            TimestampDataUTC result = default(TimestampDataUTC);
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                result.value = Axis[i].Current;
                DateTime dateTime = default(DateTime);
                DateTime dateTime2 = Axis[i].OperationHours;
                long num = 0L;
                if (Axis[i].OperationHours.Ticks == 0)
                {
                    dateTime2 = DateTime.Now;
                }

                num = getMonitoringChannel1(TargetIP);
                long num2 = num - Axis[i].OperationSubHoursOld;
                if (num2 < 0)
                {
                    num2 += 3600000;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].OperationHours = dateTime2.AddMilliseconds(num2);
                    Axis[i].OperationSubHoursOld = num;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                dateTime = Axis[i].OperationHours;
                result.Timestamp = dateTime.ToUniversalTime();
                break;
            }

            return result;
        }
    }

    TimestampDataUTC _ACI.getDemandCurrentWithTimestampUTC(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getDemandCurrentWithTimestampUTC
        return this.getDemandCurrentWithTimestampUTC(TargetIP);
    }

    //
    // Summary:
    //     Get current with drive timestamp, Monitoring Channel 1 must be configured to
    //     "Operating Sub Hours"¨!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Data structure contains current and timestamp
    public TimestampMonitoring getMonitoringChannelWithTimestamp(string TargetIP, int Channel)
    {
        if (Channel < 1)
        {
            Channel = 2;
        }

        if (Channel > 4)
        {
            Channel = 2;
        }

        int targetIPListCount = TargetIPListCount;
        TimestampMonitoring result = default(TimestampMonitoring);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result.Timestamp = Axis[i].MonitoringChannel1;
                switch (Channel)
                {
                    case 2:
                        result.value = Axis[i].MonitoringChannel2;
                        break;
                    case 3:
                        result.value = Axis[i].MonitoringChannel3;
                        break;
                    case 4:
                        result.value = Axis[i].MonitoringChannel4;
                        break;
                    default:
                        continue;
                }

                break;
            }
        }

        return result;
    }

    TimestampMonitoring _ACI.getMonitoringChannelWithTimestamp(string TargetIP, int Channel)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getMonitoringChannelWithTimestamp
        return this.getMonitoringChannelWithTimestamp(TargetIP, Channel);
    }

    //
    // Summary:
    //     Get current with drive timestamp in Coordinated Universal Time, Monitoring Channel
    //     1 must be configured to "Operating Sub Hours"¨!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   Channel:
    //     Channel number of requested Monitoring channel from 2 ... 4
    //
    // Returns:
    //     Data structure contains current and timestamp
    public TimestampMonitoringUTC getMonitoringChannelWithTimestampUTC(string TargetIP, int Channel)
    {
        if (Channel < 1)
        {
            Channel = 2;
        }

        if (Channel > 4)
        {
            Channel = 2;
        }

        int targetIPListCount = TargetIPListCount;
        checked
        {
            TimestampMonitoringUTC result = default(TimestampMonitoringUTC);
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                DateTime dateTime = default(DateTime);
                DateTime dateTime2 = Axis[i].OperationHours;
                long num = 0L;
                if (Axis[i].OperationHours.Ticks == 0)
                {
                    dateTime2 = DateTime.Now;
                }

                num = getMonitoringChannel1(TargetIP);
                long num2 = num - Axis[i].OperationSubHoursOld;
                if (num2 < 0)
                {
                    num2 += 3600000;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].OperationHours = dateTime2.AddMilliseconds(num2);
                    Axis[i].OperationSubHoursOld = num;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                dateTime = Axis[i].OperationHours;
                result.Timestamp = dateTime.ToUniversalTime();
                switch (Channel)
                {
                    case 2:
                        result.value = Axis[i].MonitoringChannel2;
                        break;
                    case 3:
                        result.value = Axis[i].MonitoringChannel3;
                        break;
                    case 4:
                        result.value = Axis[i].MonitoringChannel4;
                        break;
                    default:
                        continue;
                }

                break;
            }

            return result;
        }
    }

    TimestampMonitoringUTC _ACI.getMonitoringChannelWithTimestampUTC(string TargetIP, int Channel)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getMonitoringChannelWithTimestampUTC
        return this.getMonitoringChannelWithTimestampUTC(TargetIP, Channel);
    }

    public bool enablePseudoScopeTrace(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            Axis[i].PseudoScopeData.Clear();
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0 && !EnablePseudoScope)
            {
                EnablePseudoScope = true;
            }
        }

        return EnablePseudoScope;
    }

    bool _ACI.enablePseudoScopeTrace(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in enablePseudoScopeTrace
        return this.enablePseudoScopeTrace(TargetIP);
    }

    public bool isPseudoScopeSampling()
    {
        return EnablePseudoScope;
    }

    bool _ACI.isPseudoScopeSampling()
    {
        //ILSpy generated this explicit interface implementation from .override directive in isPseudoScopeSampling
        return this.isPseudoScopeSampling();
    }

    public List<PseudoScopeEntry> getPseudoScopeSamples(string TargetIP)
    {
        List<PseudoScopeEntry> result = new List<PseudoScopeEntry>();
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = Axis[i].PseudoScopeData;
            }
        }

        return result;
    }

    List<PseudoScopeEntry> _ACI.getPseudoScopeSamples(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getPseudoScopeSamples
        return this.getPseudoScopeSamples(TargetIP);
    }

    public int getPseudoScopeProgress(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        checked
        {
            int result = default(int);
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
                {
                    result = (int)Math.Round(Math.Round((double)Axis[i].PseudoScopeData.Count / (double)PseudoScopePoints * 100.0));
                }
            }

            return result;
        }
    }

    int _ACI.getPseudoScopeProgress(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getPseudoScopeProgress
        return this.getPseudoScopeProgress(TargetIP);
    }

    //
    // Summary:
    //     This function gets the State Machine Main State and some corresponding sub states
    //
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Actual state as StateMAchineStates
    public StateMachineStates getStateMachineState(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                return (StateMachineStates)Axis[i].StateVarHigh;
            }
        }

        StateMachineStates result = default(StateMachineStates);
        return result;
    }

    StateMachineStates _ACI.getStateMachineState(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getStateMachineState
        return this.getStateMachineState(TargetIP);
    }

    public bool isSwitchOnActive(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = (Axis[i].StatusWordLow & 2) == 2;
            }
        }

        return result;
    }

    bool _ACI.isSwitchOnActive(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isSwitchOnActive
        return this.isSwitchOnActive(TargetIP);
    }

    //
    // Summary:
    //     Check, if status word bit 8 or eventHandler bit in state var low "Event Handler
    //     Active" is true
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if bit is active
    public bool isEventHandlerActive(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = (((((Axis[i].StateVarLow & 0x10) == 16) & (Axis[i].StateVarHigh == 8)) | ((Axis[i].StatusWordHigh & 1) == 1)) ? true : false);
            }
        }

        return result;
    }

    bool _ACI.isEventHandlerActive(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isEventHandlerActive
        return this.isEventHandlerActive(TargetIP);
    }

    //
    // Summary:
    //     Check, if status word bit 9 "Special Motion Active" is true
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if bit is active
    public bool isSpecialMotionActive(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = (Axis[i].StatusWordHigh & 2) == 2;
            }
        }

        return result;
    }

    bool _ACI.isSpecialMotionActive(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isSpecialMotionActive
        return this.isSpecialMotionActive(TargetIP);
    }

    //
    // Summary:
    //     Check, if drive is on target position. Bit is reset, while starting new motion
    //
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if bit is active
    public bool isInTargetPosition(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if ((Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0) & (Axis[i].MotionCommandInterface.CountNibble == (Axis[i].StateVarLow & 0xF)))
            {
                result = (Axis[i].StateVarLow & 0x40) == 64;
            }
        }

        return result;
    }

    bool _ACI.isInTargetPosition(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isInTargetPosition
        return this.isInTargetPosition(TargetIP);
    }

    //
    // Summary:
    //     Check, if axis is referenced
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if bit is active
    public bool isHomed(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = (Axis[i].StatusWordHigh & 8) == 8;
            }
        }

        return result;
    }

    bool _ACI.isHomed(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isHomed
        return this.isHomed(TargetIP);
    }

    //
    // Summary:
    //     Check, if a fatal error occured. Fatal error need reboot of the drive!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if bit is active
    public bool isFatalError(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = (Axis[i].StatusWordHigh & 0x10) == 16;
            }
        }

        return result;
    }

    bool _ACI.isFatalError(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isFatalError
        return this.isFatalError(TargetIP);
    }

    //
    // Summary:
    //     Check, if there is still a motion ongoing
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if bit is active
    public bool isMotionActive(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if ((Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0) & (Axis[i].MotionCommandInterface.CountNibble == (Axis[i].StateVarLow & 0xF)))
            {
                result = (Axis[i].StateVarLow & 0x20) == 32;
            }
        }

        return result;
    }

    bool _ACI.isMotionActive(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isMotionActive
        return this.isMotionActive(TargetIP);
    }

    //
    // Summary:
    //     Deliver state of range indicator 1. This range indicator must be configured on
    //     the drive
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if bit is active
    public bool isRangeIndicator1(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = (Axis[i].StatusWordHigh & 0x40) == 64;
            }
        }

        return result;
    }

    bool _ACI.isRangeIndicator1(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isRangeIndicator1
        return this.isRangeIndicator1(TargetIP);
    }

    //
    // Summary:
    //     Deliver state of range indicator 2. This range indicator must be configured on
    //     the drive
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if bit is active
    public bool isRangeIndicator2(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = (Axis[i].StatusWordHigh & 0x80) == 128;
            }
        }

        return result;
    }

    bool _ACI.isRangeIndicator2(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isRangeIndicator2
        return this.isRangeIndicator2(TargetIP);
    }

    //
    // Summary:
    //     Check, if status word bit 0 "Enable Operation" is active -> drive is ready to
    //     do movements
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if bit is active
    public bool isOperationEnable(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = (Axis[i].StatusWordLow & 1) == 1;
            }
        }

        return result;
    }

    bool _ACI.isOperationEnable(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isOperationEnable
        return this.isOperationEnable(TargetIP);
    }

    //
    // Summary:
    //     Check, if status word bit 2 "Enable Operation" is active
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if bit is active
    public bool isEnableOperation(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = (Axis[i].StatusWordLow & 4) == 4;
            }
        }

        return result;
    }

    bool _ACI.isEnableOperation(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isEnableOperation
        return this.isEnableOperation(TargetIP);
    }

    //
    // Summary:
    //     Check, if status word bit 3 "Error" is active -> active for any error state
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if bit is active
    public bool isError(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = (Axis[i].StatusWordLow & 8) == 8;
            }
        }

        return result;
    }

    bool _ACI.isError(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isError
        return this.isError(TargetIP);
    }

    //
    // Summary:
    //     Check, if status word bit 4 "Safety Voltage Enable" is active -> safety input
    //     from emergency switch
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if bit is active
    public bool isSafeVoltageEnable(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = (Axis[i].StatusWordLow & 0x10) == 16;
            }
        }

        return result;
    }

    bool _ACI.isSafeVoltageEnable(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isSafeVoltageEnable
        return this.isSafeVoltageEnable(TargetIP);
    }

    //
    // Summary:
    //     Check, if status word bit 5 "/Quick Stop" is active (inverted logic!)
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if bit is active
    public bool isQuickStop(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = (Axis[i].StatusWordLow & 0x20) == 32;
            }
        }

        return result;
    }

    bool _ACI.isQuickStop(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isQuickStop
        return this.isQuickStop(TargetIP);
    }

    //
    // Summary:
    //     Check, if status word bit 6 "Switch On Locked" is active
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if bit is active
    public bool isSwitchOnLocked(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = (Axis[i].StatusWordLow & 0x40) == 64;
            }
        }

        return result;
    }

    bool _ACI.isSwitchOnLocked(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isSwitchOnLocked
        return this.isSwitchOnLocked(TargetIP);
    }

    //
    // Summary:
    //     Check, if status word bit 7 "Warning" is active
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if bit is active
    public bool isWarning(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = (Axis[i].StatusWordLow & 0x80) == 128;
            }
        }

        return result;
    }

    bool _ACI.isWarning(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isWarning
        return this.isWarning(TargetIP);
    }

    //
    // Summary:
    //     Check, if drive state machine signals "not ready to swich on". State can be changed
    //     with function "Switch On"
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if state machine signals "not ready to switch on"
    public bool isNotReadyToSwitchOnSM(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = Axis[i].StateVarHigh == 0;
            }
        }

        return result;
    }

    bool _ACI.isNotReadyToSwitchOnSM(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isNotReadyToSwitchOnSM
        return this.isNotReadyToSwitchOnSM(TargetIP);
    }

    //
    // Summary:
    //     Check, if drive state machine signals "Switch On Disabled". State can be changed
    //     with function "Switch On"
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if state machine signals "switch on disabled"
    public bool isSwitchOnDisabledSM(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = Axis[i].StateVarHigh == 1;
            }
        }

        return result;
    }

    bool _ACI.isSwitchOnDisabledSM(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isSwitchOnDisabledSM
        return this.isSwitchOnDisabledSM(TargetIP);
    }

    //
    // Summary:
    //     Check, if drive state machine signals "Ready To Switch On". State can be changed
    //     with function "Switch On"
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if state machine signals "ready to switch on"
    public bool isReadyToSwitchOnSM(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = Axis[i].StateVarHigh == 2;
            }
        }

        return result;
    }

    bool _ACI.isReadyToSwitchOnSM(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isReadyToSwitchOnSM
        return this.isReadyToSwitchOnSM(TargetIP);
    }

    //
    // Summary:
    //     Check, if drive state machine signals "Setup Error" -> Drive setup incorrect
    //
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if state machine signals "Setup error"
    public bool isSetupErrorSM(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = Axis[i].StateVarHigh == 3;
            }
        }

        return result;
    }

    bool _ACI.isSetupErrorSM(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isSetupErrorSM
        return this.isSetupErrorSM(TargetIP);
    }

    //
    // Summary:
    //     Check, if drive state machine signals "Drive in error state"
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if state machine signals "Drive in error state"
    public bool isErrorSM(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = Axis[i].StateVarHigh == 4;
            }
        }

        return result;
    }

    bool _ACI.isErrorSM(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isErrorSM
        return this.isErrorSM(TargetIP);
    }

    //
    // Summary:
    //     Check, if drive state machine signals "Drive performs harware tests"
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if state machine signals "Drive perform hardware test"
    public bool isHWTestsSM(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = Axis[i].StateVarHigh == 5;
            }
        }

        return result;
    }

    bool _ACI.isHWTestsSM(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isHWTestsSM
        return this.isHWTestsSM(TargetIP);
    }

    //
    // Summary:
    //     Check, if drive state machine signals "Ready to Operate"
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if state machine signals "Ready to operate"
    public bool isReadyToOperateSM(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = Axis[i].StateVarHigh == 6;
            }
        }

        return result;
    }

    bool _ACI.isReadyToOperateSM(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isReadyToOperateSM
        return this.isReadyToOperateSM(TargetIP);
    }

    //
    // Summary:
    //     Check, if drive state machine signals "Operation Enabled"
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if state machine signals "Operation enabled"
    public bool isOperationEnabledSM(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = Axis[i].StateVarHigh == 8;
            }
        }

        return result;
    }

    bool _ACI.isOperationEnabledSM(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isOperationEnabledSM
        return this.isOperationEnabledSM(TargetIP);
    }

    //
    // Summary:
    //     Check, if drive state machine signals "Homing is ongoing"
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if state machine signals "Drive is homing"
    public bool isHomingSM(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        bool result = default(bool);
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = Axis[i].StateVarHigh == 9;
            }
        }

        return result;
    }

    bool _ACI.isHomingSM(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isHomingSM
        return this.isHomingSM(TargetIP);
    }

    //
    // Summary:
    //     Function is obsolete
    //
    // Parameters:
    //   TargetIP:
    public bool Active(string TargetIP)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                result = true;
                break;
            }

            ACIError = "TargetIP not valid!";
            result = false;
        }

        return result;
    }

    bool _ACI.Active(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in Active
        return this.Active(TargetIP);
    }

    //
    // Summary:
    //     Change the state of control word bit 0 "Switch on". Used to handle drive state
    //     machine
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if successful
    public bool SwitchOn(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        checked
        {
            bool result = default(bool);
            for (int i = 0; i <= targetIPListCount; i++)
            {
                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
                    {
                        if (!isSwitchOnActive(TargetIP))
                        {
                            Axis[i].ControlWordLow = (byte)(Axis[i].ControlWordLow | 1);
                            result = true;
                        }
                        else
                        {
                            Axis[i].ControlWordLow = (byte)(Axis[i].ControlWordLow & 0x3E);
                        }
                    }
                    else
                    {
                        result = false;
                    }
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.SwitchOn(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in SwitchOn
        return this.SwitchOn(TargetIP);
    }

    //
    // Summary:
    //     Control the drive Switch On Bit in Control Word
    //
    // Parameters:
    //   TargetIP:
    //     Number for accessing the drive with DLL functions (1..n)
    //
    //   State:
    //     True set bit 0, false reset bit 0 in Control Word
    //
    // Returns:
    //     Actual defined state of Bit 0 in Control Word
    public bool setSwitchOnBit(string TargetIP, bool State)
    {
        int targetIPListCount = TargetIPListCount;
        checked
        {
            bool result = default(bool);
            for (int i = 0; i <= targetIPListCount; i++)
            {
                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
                    {
                        if (State)
                        {
                            Axis[i].ControlWordLow = (byte)(Axis[i].ControlWordLow | 1);
                            result = true;
                        }
                        else
                        {
                            Axis[i].ControlWordLow = (byte)(Axis[i].ControlWordLow & 0x3E);
                        }
                    }
                    else
                    {
                        result = false;
                    }
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.setSwitchOnBit(string TargetIP, bool State)
    {
        //ILSpy generated this explicit interface implementation from .override directive in setSwitchOnBit
        return this.setSwitchOnBit(TargetIP, State);
    }

    //
    // Summary:
    //     Start homing procedure by handling control word bit 11
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if successful
    public bool Homing(string TargetIP)
    {
        long num = 0L;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0 || !isOperationEnable(TargetIP))
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh | 8);
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                Thread.Sleep(400);
                num = getTimeOutTime(60000L);
                while ((Axis[i].StatusWordHigh & 8) == 0)
                {
                    Thread.Sleep(5);
                    if (!isTimeOut(num))
                    {
                        continue;
                    }

                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        Axis[i].DLLErrorText = "3: Timeout during homing";
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    break;
                }

                object taskLock3 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                bool lockTaken3 = false;
                try
                {
                    Monitor.Enter(taskLock3, ref lockTaken3);
                    Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh & 0xF7);
                }
                finally
                {
                    if (lockTaken3)
                    {
                        Monitor.Exit(taskLock3);
                    }
                }

                num = getTimeOutTime(60000L);
                while ((Axis[i].StatusWordHigh & 8) == 1)
                {
                    Thread.Sleep(5);
                    if (!isTimeOut(num))
                    {
                        continue;
                    }

                    object taskLock4 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock4);
                    bool lockTaken4 = false;
                    try
                    {
                        Monitor.Enter(taskLock4, ref lockTaken4);
                        Axis[i].DLLErrorText = "3: Timeout during homing";
                    }
                    finally
                    {
                        if (lockTaken4)
                        {
                            Monitor.Exit(taskLock4);
                        }
                    }

                    break;
                }
            }

            return true;
        }
    }

    bool _ACI.Homing(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in Homing
        return this.Homing(TargetIP);
    }

    //
    // Summary:
    //     Set bit 11 for homing in Control Word
    //
    // Parameters:
    //   TargetIP:
    //     Number for accessing the drive with DLL functions (1..n)
    //
    //   State:
    //     True set bit 11, false reset bit 11 in Control Word
    //
    // Returns:
    //     Actual defined state of Bit 0 in Control Word
    public bool setHomingBit(string TargetIP, bool State)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
                    {
                        if (State)
                        {
                            Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh | 8);
                            result = true;
                        }
                        else
                        {
                            Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh & 0xF7);
                            result = false;
                        }
                    }
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.setHomingBit(string TargetIP, bool State)
    {
        //ILSpy generated this explicit interface implementation from .override directive in setHomingBit
        return this.setHomingBit(TargetIP, State);
    }

    //
    // Summary:
    //     Acknowledge error by control word bit 7
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if successful
    public bool AckErrors(string TargetIP)
    {
        long num = 0L;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                if (!isFatalError(TargetIP))
                {
                    if (isSwitchOnActive(TargetIP))
                    {
                        SwitchOn(TargetIP);
                    }

                    Axis[i].Done = false;
                    if (isError(TargetIP))
                    {
                        object taskLock = TaskLock;
                        ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                        bool lockTaken = false;
                        try
                        {
                            Monitor.Enter(taskLock, ref lockTaken);
                            Axis[i].ControlWordLow = (byte)(Axis[i].ControlWordLow | 0x80);
                        }
                        finally
                        {
                            if (lockTaken)
                            {
                                Monitor.Exit(taskLock);
                            }
                        }

                        num = getTimeOutTime(1000L);
                        while (isError(TargetIP))
                        {
                            Thread.Sleep(5);
                            if (!isTimeOut(num))
                            {
                                continue;
                            }

                            object taskLock2 = TaskLock;
                            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                            bool lockTaken2 = false;
                            try
                            {
                                Monitor.Enter(taskLock2, ref lockTaken2);
                                Axis[i].DLLErrorText = "4: Error during error acknowledge bit set to 1";
                            }
                            finally
                            {
                                if (lockTaken2)
                                {
                                    Monitor.Exit(taskLock2);
                                }
                            }

                            break;
                        }

                        object taskLock3 = TaskLock;
                        ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                        bool lockTaken3 = false;
                        try
                        {
                            Monitor.Enter(taskLock3, ref lockTaken3);
                            Axis[i].ControlWordLow = (byte)(Axis[i].ControlWordLow & 0x7F);
                        }
                        finally
                        {
                            if (lockTaken3)
                            {
                                Monitor.Exit(taskLock3);
                            }
                        }

                        num = getTimeOutTime(1000L);
                        while (isError(TargetIP))
                        {
                            Thread.Sleep(5);
                            if (!isTimeOut(num))
                            {
                                continue;
                            }

                            object taskLock4 = TaskLock;
                            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock4);
                            bool lockTaken4 = false;
                            try
                            {
                                Monitor.Enter(taskLock4, ref lockTaken4);
                                Axis[i].DLLErrorText = "4: Error during error acknowledge bit set to 0";
                            }
                            finally
                            {
                                if (lockTaken4)
                                {
                                    Monitor.Exit(taskLock4);
                                }
                            }

                            break;
                        }
                    }

                    object taskLock5 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock5);
                    bool lockTaken5 = false;
                    try
                    {
                        Monitor.Enter(taskLock5, ref lockTaken5);
                        Axis[i].FatalError = false;
                    }
                    finally
                    {
                        if (lockTaken5)
                        {
                            Monitor.Exit(taskLock5);
                        }
                    }

                    continue;
                }

                object taskLock6 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock6);
                bool lockTaken6 = false;
                try
                {
                    Monitor.Enter(taskLock6, ref lockTaken6);
                    Axis[i].FatalError = true;
                    Axis[i].Done = true;
                }
                finally
                {
                    if (lockTaken6)
                    {
                        Monitor.Exit(taskLock6);
                    }
                }
            }

            return true;
        }
    }

    bool _ACI.AckErrors(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in AckErrors
        return this.AckErrors(TargetIP);
    }

    //
    // Summary:
    //     Set the error acknowledge bit in the control word bit 7
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   State:
    //     Bit state, which should be set
    //
    // Returns:
    //     Returns true, if successful
    public bool SetErrorAcknowledgeBit(string TargetIP, bool State)
    {
        bool flag = false;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
                    {
                        if (State)
                        {
                            Axis[i].ControlWordLow = (byte)(Axis[i].ControlWordLow | 0x80);
                            flag = true;
                        }
                        else
                        {
                            Axis[i].ControlWordLow = (byte)(Axis[i].ControlWordLow & 0x7F);
                            flag = false;
                        }
                    }
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return true;
        }
    }

    bool _ACI.SetErrorAcknowledgeBit(string TargetIP, bool State)
    {
        //ILSpy generated this explicit interface implementation from .override directive in SetErrorAcknowledgeBit
        return this.SetErrorAcknowledgeBit(TargetIP, State);
    }

    //
    // Summary:
    //     Set/Reset control word bit 8 "Jog Move +"
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if successful
    public bool JogPlus(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                    {
                        continue;
                    }

                    if (isSwitchOnActive(TargetIP))
                    {
                        if (Axis[i].ToggleJogP)
                        {
                            Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh | 1);
                            Axis[i].MotionCommandInterface.CommandHeaderMasterID = 0;
                            Axis[i].MotionCommandInterface.CommandHeaderSubID = 0;
                        }
                        else
                        {
                            Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh & 0xFE);
                        }

                        Axis[i].Done = true;
                    }
                    else
                    {
                        Axis[i].DLLError = 1;
                        Axis[i].DLLErrorText = "5: Drive not switched on!";
                    }

                    Axis[i].ToggleJogP = !Axis[i].ToggleJogP;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return true;
        }
    }

    bool _ACI.JogPlus(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in JogPlus
        return this.JogPlus(TargetIP);
    }

    //
    // Summary:
    //     Set/Reset control word bit 9 "Jog Move -"
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if successful
    public bool JogMinus(string TargetIP)
    {
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                    {
                        continue;
                    }

                    if (isSwitchOnActive(TargetIP))
                    {
                        if (Axis[i].ToggleJogM)
                        {
                            Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh | 2);
                            Axis[i].MotionCommandInterface.CommandHeaderMasterID = 0;
                            Axis[i].MotionCommandInterface.CommandHeaderSubID = 0;
                        }
                        else
                        {
                            Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh & 0xFD);
                        }

                        Axis[i].Done = true;
                    }
                    else
                    {
                        Axis[i].DLLError = 1;
                        Axis[i].DLLErrorText = "5: Drive not switched on!";
                    }

                    Axis[i].ToggleJogM = !Axis[i].ToggleJogM;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return true;
        }
    }

    bool _ACI.JogMinus(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in JogMinus
        return this.JogMinus(TargetIP);
    }

    //
    // Summary:
    //     Set the bit state of control Word Jog + Bit and return true, if bit value sucessfully
    //     set
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   state:
    //     Bit state, which should be set
    //
    // Returns:
    //     Returns true, if bit value successfully written
    public bool setJogPlus(string TargetIP, bool state)
    {
        return setBit8(TargetIP, state);
    }

    bool _ACI.setJogPlus(string TargetIP, bool state)
    {
        //ILSpy generated this explicit interface implementation from .override directive in setJogPlus
        return this.setJogPlus(TargetIP, state);
    }

    //
    // Summary:
    //     Set the bit state of control Word Jog - Bit and return true, if bit value sucessfully
    //     set
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   state:
    //     Bit state, which should be set
    //
    // Returns:
    //     Returns true, if bit value successfully written
    public bool setJogMinus(string TargetIP, bool state)
    {
        return setBit9(TargetIP, state);
    }

    bool _ACI.setJogMinus(string TargetIP, bool state)
    {
        //ILSpy generated this explicit interface implementation from .override directive in setJogMinus
        return this.setJogMinus(TargetIP, state);
    }

    //
    // Summary:
    //     Control the drive Bit0 in Control Word
    //
    // Parameters:
    //   TargetIP:
    //     Number for accessing the drive with DLL functions (1..n)
    //
    //   State:
    //     True set bit 0, false reset bit 0 in Control Word
    //
    // Returns:
    //     Actual defined state of Bit 0 in Control Word
    public bool setBit0(string TargetIP, bool State)
    {
        int targetIPListCount = TargetIPListCount;
        checked
        {
            bool result = default(bool);
            for (int i = 0; i <= targetIPListCount; i++)
            {
                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
                    {
                        if (State)
                        {
                            Axis[i].ControlWordLow = (byte)(Axis[i].ControlWordLow | 1);
                            result = true;
                        }
                        else
                        {
                            Axis[i].ControlWordLow = (byte)(Axis[i].ControlWordLow & 0xFE);
                        }
                    }
                    else
                    {
                        result = false;
                    }
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.setBit0(string TargetIP, bool State)
    {
        //ILSpy generated this explicit interface implementation from .override directive in setBit0
        return this.setBit0(TargetIP, State);
    }

    //
    // Summary:
    //     Control the drive Bit1 in Control Word
    //
    // Parameters:
    //   TargetIP:
    //     Number for accessing the drive with DLL functions (1..n)
    //
    //   State:
    //     True set bit 1, false reset bit 1 in Control Word
    //
    // Returns:
    //     Actual defined state of Bit 1 in Control Word
    public bool setBit1(string TargetIP, bool State)
    {
        int targetIPListCount = TargetIPListCount;
        checked
        {
            bool result = default(bool);
            for (int i = 0; i <= targetIPListCount; i++)
            {
                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
                    {
                        if (State)
                        {
                            Axis[i].ControlWordLow = (byte)(Axis[i].ControlWordLow | 2);
                            result = true;
                        }
                        else
                        {
                            Axis[i].ControlWordLow = (byte)(Axis[i].ControlWordLow & 0xFD);
                        }
                    }
                    else
                    {
                        result = false;
                    }
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.setBit1(string TargetIP, bool State)
    {
        //ILSpy generated this explicit interface implementation from .override directive in setBit1
        return this.setBit1(TargetIP, State);
    }

    //
    // Summary:
    //     Control the drive Bit2 in Control Word
    //
    // Parameters:
    //   TargetIP:
    //     Number for accessing the drive with DLL functions (1..n)
    //
    //   State:
    //     True set bit 2, false reset bit 2 in Control Word
    //
    // Returns:
    //     Actual defined state of Bit 2 in Control Word
    public bool setBit2(string TargetIP, bool State)
    {
        int targetIPListCount = TargetIPListCount;
        checked
        {
            bool result = default(bool);
            for (int i = 0; i <= targetIPListCount; i++)
            {
                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
                    {
                        if (State)
                        {
                            Axis[i].ControlWordLow = (byte)(Axis[i].ControlWordLow | 4);
                            result = true;
                        }
                        else
                        {
                            Axis[i].ControlWordLow = (byte)(Axis[i].ControlWordLow & 0xFB);
                        }
                    }
                    else
                    {
                        result = false;
                    }
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.setBit2(string TargetIP, bool State)
    {
        //ILSpy generated this explicit interface implementation from .override directive in setBit2
        return this.setBit2(TargetIP, State);
    }

    //
    // Summary:
    //     Control the drive Bit3 in Control Word
    //
    // Parameters:
    //   TargetIP:
    //     Number for accessing the drive with DLL functions (1..n)
    //
    //   State:
    //     True set bit 3, false reset bit 3 in Control Word
    //
    // Returns:
    //     Actual defined state of Bit 3 in Control Word
    public bool setBit3(string TargetIP, bool State)
    {
        int targetIPListCount = TargetIPListCount;
        checked
        {
            bool result = default(bool);
            for (int i = 0; i <= targetIPListCount; i++)
            {
                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
                    {
                        if (State)
                        {
                            Axis[i].ControlWordLow = (byte)(Axis[i].ControlWordLow | 8);
                            result = true;
                        }
                        else
                        {
                            Axis[i].ControlWordLow = (byte)(Axis[i].ControlWordLow & 0xF7);
                        }
                    }
                    else
                    {
                        result = false;
                    }
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.setBit3(string TargetIP, bool State)
    {
        //ILSpy generated this explicit interface implementation from .override directive in setBit3
        return this.setBit3(TargetIP, State);
    }

    //
    // Summary:
    //     Control the drive Bit4 in Control Word
    //
    // Parameters:
    //   TargetIP:
    //     Number for accessing the drive with DLL functions (1..n)
    //
    //   State:
    //     True set bit 4, false reset bit 4 in Control Word
    //
    // Returns:
    //     Actual defined state of Bit 4 in Control Word
    public bool setBit4(string TargetIP, bool State)
    {
        int targetIPListCount = TargetIPListCount;
        checked
        {
            bool result = default(bool);
            for (int i = 0; i <= targetIPListCount; i++)
            {
                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
                    {
                        if (State)
                        {
                            Axis[i].ControlWordLow = (byte)(Axis[i].ControlWordLow | 0x10);
                            result = true;
                        }
                        else
                        {
                            Axis[i].ControlWordLow = (byte)(Axis[i].ControlWordLow & 0xEF);
                        }
                    }
                    else
                    {
                        result = false;
                    }
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.setBit4(string TargetIP, bool State)
    {
        //ILSpy generated this explicit interface implementation from .override directive in setBit4
        return this.setBit4(TargetIP, State);
    }

    //
    // Summary:
    //     Control the drive Bit5 in Control Word
    //
    // Parameters:
    //   TargetIP:
    //     Number for accessing the drive with DLL functions (1..n)
    //
    //   State:
    //     True set bit 5, false reset bit 5 in Control Word
    //
    // Returns:
    //     Actual defined state of Bit 5 in Control Word
    public bool setBit5(string TargetIP, bool State)
    {
        int targetIPListCount = TargetIPListCount;
        checked
        {
            bool result = default(bool);
            for (int i = 0; i <= targetIPListCount; i++)
            {
                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
                    {
                        if (State)
                        {
                            Axis[i].ControlWordLow = (byte)(Axis[i].ControlWordLow | 0x20);
                            result = true;
                        }
                        else
                        {
                            Axis[i].ControlWordLow = (byte)(Axis[i].ControlWordLow & 0xDF);
                        }
                    }
                    else
                    {
                        result = false;
                    }
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.setBit5(string TargetIP, bool State)
    {
        //ILSpy generated this explicit interface implementation from .override directive in setBit5
        return this.setBit5(TargetIP, State);
    }

    //
    // Summary:
    //     Control the drive Bit6 in Control Word
    //
    // Parameters:
    //   TargetIP:
    //     Number for accessing the drive with DLL functions (1..n)
    //
    //   State:
    //     True set bit 6, false reset bit 6 in Control Word
    //
    // Returns:
    //     Actual defined state of Bit 6 in Control Word
    public bool setBit6(string TargetIP, bool State)
    {
        int targetIPListCount = TargetIPListCount;
        checked
        {
            bool result = default(bool);
            for (int i = 0; i <= targetIPListCount; i++)
            {
                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
                    {
                        if (State)
                        {
                            Axis[i].ControlWordLow = (byte)(Axis[i].ControlWordLow | 0x40);
                            result = true;
                        }
                        else
                        {
                            Axis[i].ControlWordLow = (byte)(Axis[i].ControlWordLow & 0xBF);
                        }
                    }
                    else
                    {
                        result = false;
                    }
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.setBit6(string TargetIP, bool State)
    {
        //ILSpy generated this explicit interface implementation from .override directive in setBit6
        return this.setBit6(TargetIP, State);
    }

    //
    // Summary:
    //     Control the drive Bit7 in Control Word
    //
    // Parameters:
    //   TargetIP:
    //     Number for accessing the drive with DLL functions (1..n)
    //
    //   State:
    //     True set bit 7, false reset bit 7 in Control Word
    //
    // Returns:
    //     Actual defined state of Bit 7 in Control Word
    public bool setBit7(string TargetIP, bool State)
    {
        int targetIPListCount = TargetIPListCount;
        checked
        {
            bool result = default(bool);
            for (int i = 0; i <= targetIPListCount; i++)
            {
                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
                    {
                        if (State)
                        {
                            Axis[i].ControlWordLow = (byte)(Axis[i].ControlWordLow | 0x80);
                            result = true;
                        }
                        else
                        {
                            Axis[i].ControlWordLow = (byte)(Axis[i].ControlWordLow & 0x7F);
                        }
                    }
                    else
                    {
                        result = false;
                    }
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.setBit7(string TargetIP, bool State)
    {
        //ILSpy generated this explicit interface implementation from .override directive in setBit7
        return this.setBit7(TargetIP, State);
    }

    //
    // Summary:
    //     Control the drive Bit8 in Control Word
    //
    // Parameters:
    //   TargetIP:
    //     Number for accessing the drive with DLL functions (1..n)
    //
    //   State:
    //     True set bit 8, false reset bit 8 in Control Word
    //
    // Returns:
    //     Actual defined state of Bit 8 in Control Word
    public bool setBit8(string TargetIP, bool State)
    {
        int targetIPListCount = TargetIPListCount;
        checked
        {
            bool result = default(bool);
            for (int i = 0; i <= targetIPListCount; i++)
            {
                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
                    {
                        if (State)
                        {
                            Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh | 1);
                            result = true;
                        }
                        else
                        {
                            Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh & 0xFE);
                        }
                    }
                    else
                    {
                        result = false;
                    }
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.setBit8(string TargetIP, bool State)
    {
        //ILSpy generated this explicit interface implementation from .override directive in setBit8
        return this.setBit8(TargetIP, State);
    }

    //
    // Summary:
    //     Control the drive Bit9 in Control Word
    //
    // Parameters:
    //   TargetIP:
    //     Number for accessing the drive with DLL functions (1..n)
    //
    //   State:
    //     True set bit 9, false reset bit 9 in Control Word
    //
    // Returns:
    //     Actual defined state of Bit 9 in Control Word
    public bool setBit9(string TargetIP, bool State)
    {
        int targetIPListCount = TargetIPListCount;
        checked
        {
            bool result = default(bool);
            for (int i = 0; i <= targetIPListCount; i++)
            {
                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
                    {
                        if (State)
                        {
                            Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh | 2);
                            result = true;
                        }
                        else
                        {
                            Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh & 0xFD);
                        }
                    }
                    else
                    {
                        result = false;
                    }
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.setBit9(string TargetIP, bool State)
    {
        //ILSpy generated this explicit interface implementation from .override directive in setBit9
        return this.setBit9(TargetIP, State);
    }

    //
    // Summary:
    //     Control the drive Bit10 in Control Word
    //
    // Parameters:
    //   TargetIP:
    //     Number for accessing the drive with DLL functions (1..n)
    //
    //   State:
    //     True set bit 10, false reset bit 10 in Control Word
    //
    // Returns:
    //     Actual defined state of Bit 10 in Control Word
    public bool setBit10(string TargetIP, bool State)
    {
        int targetIPListCount = TargetIPListCount;
        checked
        {
            bool result = default(bool);
            for (int i = 0; i <= targetIPListCount; i++)
            {
                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
                    {
                        if (State)
                        {
                            Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh | 4);
                            result = true;
                        }
                        else
                        {
                            Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh & 0xFB);
                        }
                    }
                    else
                    {
                        result = false;
                    }
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.setBit10(string TargetIP, bool State)
    {
        //ILSpy generated this explicit interface implementation from .override directive in setBit10
        return this.setBit10(TargetIP, State);
    }

    //
    // Summary:
    //     Control the drive Bit11 in Control Word
    //
    // Parameters:
    //   TargetIP:
    //     Number for accessing the drive with DLL functions (1..n)
    //
    //   State:
    //     True set bit 11, false reset bit 11 in Control Word
    //
    // Returns:
    //     Actual defined state of Bit 11 in Control Word
    public bool setBit11(string TargetIP, bool State)
    {
        int targetIPListCount = TargetIPListCount;
        checked
        {
            bool result = default(bool);
            for (int i = 0; i <= targetIPListCount; i++)
            {
                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
                    {
                        if (State)
                        {
                            Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh | 8);
                            result = true;
                        }
                        else
                        {
                            Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh & 0xF7);
                        }
                    }
                    else
                    {
                        result = false;
                    }
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.setBit11(string TargetIP, bool State)
    {
        //ILSpy generated this explicit interface implementation from .override directive in setBit11
        return this.setBit11(TargetIP, State);
    }

    //
    // Summary:
    //     Control the drive Bit12 in Control Word
    //
    // Parameters:
    //   TargetIP:
    //     Number for accessing the drive with DLL functions (1..n)
    //
    //   State:
    //     True set bit 12, false reset bit 12 in Control Word
    //
    // Returns:
    //     Actual defined state of Bit 12 in Control Word
    public bool setBit12(string TargetIP, bool State)
    {
        int targetIPListCount = TargetIPListCount;
        checked
        {
            bool result = default(bool);
            for (int i = 0; i <= targetIPListCount; i++)
            {
                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
                    {
                        if (State)
                        {
                            Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh | 0x10);
                            result = true;
                        }
                        else
                        {
                            Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh & 0xEF);
                        }
                    }
                    else
                    {
                        result = false;
                    }
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.setBit12(string TargetIP, bool State)
    {
        //ILSpy generated this explicit interface implementation from .override directive in setBit12
        return this.setBit12(TargetIP, State);
    }

    //
    // Summary:
    //     Control the drive Bit13 in Control Word
    //
    // Parameters:
    //   TargetIP:
    //     Number for accessing the drive with DLL functions (1..n)
    //
    //   State:
    //     True set bit 13, false reset bit 13 in Control Word
    //
    // Returns:
    //     Actual defined state of Bit 13 in Control Word
    public bool setBit13(string TargetIP, bool State)
    {
        int targetIPListCount = TargetIPListCount;
        checked
        {
            bool result = default(bool);
            for (int i = 0; i <= targetIPListCount; i++)
            {
                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
                    {
                        if (State)
                        {
                            Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh | 0x20);
                            result = true;
                        }
                        else
                        {
                            Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh & 0xDF);
                        }
                    }
                    else
                    {
                        result = false;
                    }
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.setBit13(string TargetIP, bool State)
    {
        //ILSpy generated this explicit interface implementation from .override directive in setBit13
        return this.setBit13(TargetIP, State);
    }

    //
    // Summary:
    //     Control the drive Bit14 in Control Word
    //
    // Parameters:
    //   TargetIP:
    //     Number for accessing the drive with DLL functions (1..n)
    //
    //   State:
    //     True set bit 14, false reset bit 14 in Control Word
    //
    // Returns:
    //     Actual defined state of Bit 14 in Control Word
    public bool setBit14(string TargetIP, bool State)
    {
        int targetIPListCount = TargetIPListCount;
        checked
        {
            bool result = default(bool);
            for (int i = 0; i <= targetIPListCount; i++)
            {
                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
                    {
                        if (State)
                        {
                            Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh | 0x40);
                            result = true;
                        }
                        else
                        {
                            Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh & 0xBF);
                        }
                    }
                    else
                    {
                        result = false;
                    }
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.setBit14(string TargetIP, bool State)
    {
        //ILSpy generated this explicit interface implementation from .override directive in setBit14
        return this.setBit14(TargetIP, State);
    }

    //
    // Summary:
    //     Control the drive Bit15 in Control Word
    //
    // Parameters:
    //   TargetIP:
    //     Number for accessing the drive with DLL functions (1..n)
    //
    //   State:
    //     True set bit 15, false reset bit 15 in Control Word
    //
    // Returns:
    //     Actual defined state of Bit 15 in Control Word
    public bool setBit15(string TargetIP, bool State)
    {
        int targetIPListCount = TargetIPListCount;
        checked
        {
            bool result = default(bool);
            for (int i = 0; i <= targetIPListCount; i++)
            {
                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
                    {
                        if (State)
                        {
                            Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh | 0x80);
                            result = true;
                        }
                        else
                        {
                            Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh & 0x7F);
                        }
                    }
                    else
                    {
                        result = false;
                    }
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.setBit15(string TargetIP, bool State)
    {
        //ILSpy generated this explicit interface implementation from .override directive in setBit15
        return this.setBit15(TargetIP, State);
    }

    //
    // Summary:
    //     Parameter channel function: Reboot drive, set defaults, stop/start MC software
    //
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   Mode:
    //     0: reboot drive, 1: set OS ROM parameters to default, 2: set MC ROM parameters
    //     to default, 3: set interface ROM parameters to default, 4: set ROM application
    //     parameters to defualt, 5: Stop MC software, 6: Start MC software
    //
    // Returns:
    //     Returns true, if successful
    public long LMcf_StartStopDefault(string TargetIP, int Mode)
    {
        long result = 0L;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                switch (Mode)
                {
                    case 0:
                        Axis[i].RealTimeConfigID = 48;
                        break;
                    case 1:
                        Axis[i].RealTimeConfigID = 49;
                        break;
                    case 2:
                        Axis[i].RealTimeConfigID = 50;
                        break;
                    case 3:
                        Axis[i].RealTimeConfigID = 51;
                        break;
                    case 4:
                        Axis[i].RealTimeConfigID = 52;
                        break;
                    case 5:
                        Axis[i].RealTimeConfigID = 53;
                        break;
                    case 6:
                        Axis[i].RealTimeConfigID = 54;
                        break;
                }

                if (Axis[i].RealTimeConfigCommandCount > 14)
                {
                    Axis[i].RealTimeConfigCommandCount = 1;
                }

                Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(checked(Axis[i].RealTimeConfigCommandCount + 1)));
                if (Axis[i].RealTimeConfigCommandCount > 0 && Mode == 0)
                {
                    Axis[i].RealTimeConfigCommandCount = 0;
                }

                Axis[i].RealTimeConfigArgs1High = 0;
                Axis[i].RealTimeConfigArgs1Low = 0;
                Axis[i].RealTimeConfigArgs2High = 0;
                Axis[i].RealTimeConfigArgs2Low = 0;
                Axis[i].RealTimeConfigArgs3High = 0;
                Axis[i].RealTimeConfigArgs3Low = 0;
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = 0;
                Axis[i].MotionCommandInterface.CommandHeaderSubID = 0;
                Axis[i].MotionCommandInterface.CountNibble = 0;
                Axis[i].setTimeoutObservation = true;
                Axis[i].setSkipAmountResponsePackets = 20;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }

            Thread.Sleep(1000);
            long timeOutTime = getTimeOutTime(45000L);
            bool flag = false;
            if ((int)Axis[i].FirmwareVersion > 1545)
            {
                flag = true;
            }

            while (!((Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount) & (Axis[i].RealTimeConfigIDStatus == 0)))
            {
                if (!isTimeOut(timeOutTime))
                {
                    continue;
                }

                object taskLock2 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                bool lockTaken2 = false;
                try
                {
                    Monitor.Enter(taskLock2, ref lockTaken2);
                    Axis[i].DLLErrorText = "6: Config Channel Timeout by performing Start/Stop/Default settings. Error code: " + Conversion.Hex(Axis[i].RealTimeConfigIDStatus).ToString();
                    ACIError = "Config Channel Error";
                    Axis[i].RealTimeConfigStatusArgs2Low = 9;
                    Axis[i].RealTimeConfigStatusArgs2High = 0;
                    Axis[i].RealTimeConfigStatusArgs3Low = 9;
                    Axis[i].RealTimeConfigStatusArgs3High = 0;
                }
                finally
                {
                    if (lockTaken2)
                    {
                        Monitor.Exit(taskLock2);
                    }
                }

                break;
            }

            if (!flag)
            {
                Thread.Sleep(500);
            }

            object taskLock3 = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
            bool lockTaken3 = false;
            try
            {
                Monitor.Enter(taskLock3, ref lockTaken3);
                byte[] value = new byte[8]
                {
                    Axis[i].RealTimeConfigStatusArgs2Low,
                    Axis[i].RealTimeConfigStatusArgs2High,
                    Axis[i].RealTimeConfigStatusArgs3Low,
                    Axis[i].RealTimeConfigStatusArgs3High,
                    0,
                    0,
                    0,
                    0
                };
                Axis[i].RealTimeConfigID = 0;
                if (Axis[i].RealTimeConfigCommandCount == 0)
                {
                    Axis[i].RealTimeConfigCommandCount = 1;
                }

                Axis[i].DLLErrorText = "";
                Axis[i].DLLError = 0;
                if (!flag && Mode == 6)
                {
                    Axis[i].setSkipAmountResponsePackets = 50;
                    Axis[i].TimeStampReceive = 0L;
                    Axis[i].TimeStampSent = 0L;
                    Axis[i].DLLErrorText = "";
                    Axis[i].DLLError = 0;
                }

                Axis[i].setTimeoutObservation = false;
                result = checked((long)BitConverter.ToUInt64(value, 0));
            }
            finally
            {
                if (lockTaken3)
                {
                    Monitor.Exit(taskLock3);
                }
            }
        }

        return result;
    }

    long _ACI.LMcf_StartStopDefault(string TargetIP, int Mode)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMcf_StartStopDefault
        return this.LMcf_StartStopDefault(TargetIP, Mode);
    }

    //
    // Summary:
    //     Parameter channel function: getROM value by UPID, UPID and scale for raw output
    //     data can be found in LinMot Talk -> Variables
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   UPID:
    //     UPID address of drive variable
    //
    // Returns:
    //     Raw data of requested UPID
    public long getROM_ByUPID(string TargetIP, uint UPID)
    {
        long result = 0L;
        byte[] bytes = BitConverter.GetBytes(UPID);
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].RealTimeConfigID = 16;
                    if (Axis[i].RealTimeConfigCommandCount > 14)
                    {
                        Axis[i].RealTimeConfigCommandCount = 1;
                    }

                    Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[i].RealTimeConfigCommandCount + 1));
                    Axis[i].RealTimeConfigCommandCount = Axis[i].RealTimeConfigCommandCount;
                    Axis[i].RealTimeConfigArgs1High = bytes[1];
                    Axis[i].RealTimeConfigArgs1Low = bytes[0];
                    Axis[i].RealTimeConfigArgs2High = 0;
                    Axis[i].RealTimeConfigArgs2Low = 0;
                    Axis[i].RealTimeConfigArgs3High = 0;
                    Axis[i].RealTimeConfigArgs3Low = 0;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                long timeOutTime = getTimeOutTime(1000L);
                while (!((Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount) & (Axis[i].RealTimeConfigIDStatus == 0)))
                {
                    if (!isTimeOut(timeOutTime))
                    {
                        continue;
                    }

                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        Axis[i].DLLErrorText = "7: Config Channel Timeout by requesting ROM value by UPID. Error code: " + Conversion.Hex(Axis[i].RealTimeConfigIDStatus).ToString();
                        ACIError = "Config Channel Error";
                        Axis[i].RealTimeConfigStatusArgs2Low = 9;
                        Axis[i].RealTimeConfigStatusArgs2High = 0;
                        Axis[i].RealTimeConfigStatusArgs3Low = 9;
                        Axis[i].RealTimeConfigStatusArgs3High = 0;
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    break;
                }

                object taskLock3 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                bool lockTaken3 = false;
                try
                {
                    Monitor.Enter(taskLock3, ref lockTaken3);
                    byte[] value = new byte[8]
                    {
                        Axis[i].RealTimeConfigStatusArgs2Low,
                        Axis[i].RealTimeConfigStatusArgs2High,
                        Axis[i].RealTimeConfigStatusArgs3Low,
                        Axis[i].RealTimeConfigStatusArgs3High,
                        0,
                        0,
                        0,
                        0
                    };
                    Axis[i].RealTimeConfigID = 0;
                    result = (long)BitConverter.ToUInt64(value, 0);
                }
                finally
                {
                    if (lockTaken3)
                    {
                        Monitor.Exit(taskLock3);
                    }
                }
            }

            return result;
        }
    }

    long _ACI.getROM_ByUPID(string TargetIP, uint UPID)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getROM_ByUPID
        return this.getROM_ByUPID(TargetIP, UPID);
    }

    //
    // Summary:
    //     Parameter channel function: getRAM value by UPID, UPID and scale for raw output
    //     data can be found in LinMot Talk -> Variables
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   UPID:
    //     UPID address of drive variable
    //
    // Returns:
    //     Raw data of requested UPID
    public long getRAM_ByUPID(string TargetIP, uint UPID)
    {
        long result = 0L;
        byte[] bytes = BitConverter.GetBytes(UPID);
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].RealTimeConfigID = 17;
                    if (Axis[i].RealTimeConfigCommandCount > 14)
                    {
                        Axis[i].RealTimeConfigCommandCount = 1;
                    }

                    Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[i].RealTimeConfigCommandCount + 1));
                    Axis[i].RealTimeConfigArgs1High = bytes[1];
                    Axis[i].RealTimeConfigArgs1Low = bytes[0];
                    Axis[i].RealTimeConfigArgs2High = 0;
                    Axis[i].RealTimeConfigArgs2Low = 0;
                    Axis[i].RealTimeConfigArgs3High = 0;
                    Axis[i].RealTimeConfigArgs3Low = 0;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                long timeOutTime = getTimeOutTime(1000L);
                while (!((Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount) & (Axis[i].RealTimeConfigIDStatus == 0)))
                {
                    if (!isTimeOut(timeOutTime))
                    {
                        continue;
                    }

                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        Axis[i].DLLErrorText = "8: Config Channel Timeout by requesting RAM value by UPID. Error code: " + Conversion.Hex(Axis[i].RealTimeConfigIDStatus).ToString();
                        ACIError = "Config Channel Error";
                        Axis[i].RealTimeConfigStatusArgs2Low = 9;
                        Axis[i].RealTimeConfigStatusArgs2High = 0;
                        Axis[i].RealTimeConfigStatusArgs3Low = 9;
                        Axis[i].RealTimeConfigStatusArgs3High = 0;
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    break;
                }

                object taskLock3 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                bool lockTaken3 = false;
                try
                {
                    Monitor.Enter(taskLock3, ref lockTaken3);
                    byte[] value = new byte[8]
                    {
                        Axis[i].RealTimeConfigStatusArgs2Low,
                        Axis[i].RealTimeConfigStatusArgs2High,
                        Axis[i].RealTimeConfigStatusArgs3Low,
                        Axis[i].RealTimeConfigStatusArgs3High,
                        0,
                        0,
                        0,
                        0
                    };
                    Axis[i].RealTimeConfigID = 0;
                    result = (long)BitConverter.ToUInt64(value, 0);
                }
                finally
                {
                    if (lockTaken3)
                    {
                        Monitor.Exit(taskLock3);
                    }
                }
            }

            return result;
        }
    }

    long _ACI.getRAM_ByUPID(string TargetIP, uint UPID)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getRAM_ByUPID
        return this.getRAM_ByUPID(TargetIP, UPID);
    }

    //
    // Summary:
    //     Parameter channel function: get minimum value by UPID, UPID and scale for raw
    //     output data can be found in LinMot Talk -> Variables
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   UPID:
    //     UPID address of drive variable
    //
    // Returns:
    //     Raw data of requested UPID
    public long getMinVal_ByUPID(string TargetIP, uint UPID)
    {
        long result = 0L;
        byte[] bytes = BitConverter.GetBytes(UPID);
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].RealTimeConfigID = 21;
                    if (Axis[i].RealTimeConfigCommandCount > 14)
                    {
                        Axis[i].RealTimeConfigCommandCount = 1;
                    }

                    Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[i].RealTimeConfigCommandCount + 1));
                    Axis[i].RealTimeConfigCommandCount = Axis[i].RealTimeConfigCommandCount;
                    Axis[i].RealTimeConfigArgs1High = bytes[1];
                    Axis[i].RealTimeConfigArgs1Low = bytes[0];
                    Axis[i].RealTimeConfigArgs2High = 0;
                    Axis[i].RealTimeConfigArgs2Low = 0;
                    Axis[i].RealTimeConfigArgs3High = 0;
                    Axis[i].RealTimeConfigArgs3Low = 0;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                long timeOutTime = getTimeOutTime(1000L);
                while (!((Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount) & (Axis[i].RealTimeConfigIDStatus == 0)))
                {
                    if (!isTimeOut(timeOutTime))
                    {
                        continue;
                    }

                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        Axis[i].DLLErrorText = "9: Config Channel Timeout by requesting minimum value by UPID. Error code: " + Conversion.Hex(Axis[i].RealTimeConfigIDStatus).ToString();
                        Axis[i].RealTimeConfigStatusArgs2Low = 9;
                        Axis[i].RealTimeConfigStatusArgs2High = 0;
                        Axis[i].RealTimeConfigStatusArgs3Low = 9;
                        Axis[i].RealTimeConfigStatusArgs3High = 0;
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    break;
                }

                object taskLock3 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                bool lockTaken3 = false;
                try
                {
                    Monitor.Enter(taskLock3, ref lockTaken3);
                    byte[] value = new byte[8]
                    {
                        Axis[i].RealTimeConfigStatusArgs2Low,
                        Axis[i].RealTimeConfigStatusArgs2High,
                        Axis[i].RealTimeConfigStatusArgs3Low,
                        Axis[i].RealTimeConfigStatusArgs3High,
                        0,
                        0,
                        0,
                        0
                    };
                    Axis[i].RealTimeConfigID = 0;
                    result = (long)BitConverter.ToUInt64(value, 0);
                }
                finally
                {
                    if (lockTaken3)
                    {
                        Monitor.Exit(taskLock3);
                    }
                }
            }

            return result;
        }
    }

    long _ACI.getMinVal_ByUPID(string TargetIP, uint UPID)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getMinVal_ByUPID
        return this.getMinVal_ByUPID(TargetIP, UPID);
    }

    //
    // Summary:
    //     Parameter channel function: get maximum value by UPID, UPID and scale for raw
    //     output data can be found in LinMot Talk -> Variables
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   UPID:
    //     UPID address of drive variable
    //
    // Returns:
    //     Raw data of requested UPID
    public long getMaxVal_ByUPID(string TargetIP, uint UPID)
    {
        long result = 0L;
        byte[] bytes = BitConverter.GetBytes(UPID);
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].RealTimeConfigID = 22;
                    if (Axis[i].RealTimeConfigCommandCount > 14)
                    {
                        Axis[i].RealTimeConfigCommandCount = 1;
                    }

                    Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[i].RealTimeConfigCommandCount + 1));
                    Axis[i].RealTimeConfigCommandCount = Axis[i].RealTimeConfigCommandCount;
                    Axis[i].RealTimeConfigArgs1High = bytes[1];
                    Axis[i].RealTimeConfigArgs1Low = bytes[0];
                    Axis[i].RealTimeConfigArgs2High = 0;
                    Axis[i].RealTimeConfigArgs2Low = 0;
                    Axis[i].RealTimeConfigArgs3High = 0;
                    Axis[i].RealTimeConfigArgs3Low = 0;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                long timeOutTime = getTimeOutTime(1000L);
                while (!((Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount) & (Axis[i].RealTimeConfigIDStatus == 0)))
                {
                    if (!isTimeOut(timeOutTime))
                    {
                        continue;
                    }

                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        Axis[i].DLLErrorText = "10: Config Channel Timeout by requesting maximum value by UPID. Error code: " + Conversion.Hex(Axis[i].RealTimeConfigIDStatus).ToString();
                        ACIError = "Config Channel Error";
                        Axis[i].RealTimeConfigStatusArgs2Low = 9;
                        Axis[i].RealTimeConfigStatusArgs2High = 0;
                        Axis[i].RealTimeConfigStatusArgs3Low = 9;
                        Axis[i].RealTimeConfigStatusArgs3High = 0;
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    break;
                }

                object taskLock3 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                bool lockTaken3 = false;
                try
                {
                    Monitor.Enter(taskLock3, ref lockTaken3);
                    byte[] value = new byte[8]
                    {
                        Axis[i].RealTimeConfigStatusArgs2Low,
                        Axis[i].RealTimeConfigStatusArgs2High,
                        Axis[i].RealTimeConfigStatusArgs3Low,
                        Axis[i].RealTimeConfigStatusArgs3High,
                        0,
                        0,
                        0,
                        0
                    };
                    Axis[i].RealTimeConfigID = 0;
                    result = (long)BitConverter.ToUInt64(value, 0);
                }
                finally
                {
                    if (lockTaken3)
                    {
                        Monitor.Exit(taskLock3);
                    }
                }
            }

            return result;
        }
    }

    long _ACI.getMaxVal_ByUPID(string TargetIP, uint UPID)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getMaxVal_ByUPID
        return this.getMaxVal_ByUPID(TargetIP, UPID);
    }

    //
    // Summary:
    //     Parameter channel function: get default value by UPID, UPID and scale for raw
    //     output data can be found in LinMot Talk -> Variables
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   UPID:
    //     UPID address of drive variable
    //
    // Returns:
    //     Raw data of requested UPID
    public long getDefault_ByUPID(string TargetIP, uint UPID)
    {
        long result = 0L;
        byte[] bytes = BitConverter.GetBytes(UPID);
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].RealTimeConfigID = 23;
                    if (Axis[i].RealTimeConfigCommandCount > 14)
                    {
                        Axis[i].RealTimeConfigCommandCount = 1;
                    }

                    Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[i].RealTimeConfigCommandCount + 1));
                    Axis[i].RealTimeConfigCommandCount = Axis[i].RealTimeConfigCommandCount;
                    Axis[i].RealTimeConfigArgs1High = bytes[1];
                    Axis[i].RealTimeConfigArgs1Low = bytes[0];
                    Axis[i].RealTimeConfigArgs2High = 0;
                    Axis[i].RealTimeConfigArgs2Low = 0;
                    Axis[i].RealTimeConfigArgs3High = 0;
                    Axis[i].RealTimeConfigArgs3Low = 0;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                long timeOutTime = getTimeOutTime(1000L);
                while (!((Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount) & (Axis[i].RealTimeConfigIDStatus == 0)))
                {
                    if (!isTimeOut(timeOutTime))
                    {
                        continue;
                    }

                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        Axis[i].DLLErrorText = "11: Config Channel Timeout by requesting default value by UPID. Error code: " + Conversion.Hex(Axis[i].RealTimeConfigIDStatus).ToString();
                        ACIError = "Config Channel Error";
                        Axis[i].RealTimeConfigStatusArgs2Low = 9;
                        Axis[i].RealTimeConfigStatusArgs2High = 0;
                        Axis[i].RealTimeConfigStatusArgs3Low = 9;
                        Axis[i].RealTimeConfigStatusArgs3High = 0;
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    break;
                }

                object taskLock3 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                bool lockTaken3 = false;
                try
                {
                    Monitor.Enter(taskLock3, ref lockTaken3);
                    byte[] value = new byte[8]
                    {
                        Axis[i].RealTimeConfigStatusArgs2Low,
                        Axis[i].RealTimeConfigStatusArgs2High,
                        Axis[i].RealTimeConfigStatusArgs3Low,
                        Axis[i].RealTimeConfigStatusArgs3High,
                        0,
                        0,
                        0,
                        0
                    };
                    Axis[i].RealTimeConfigID = 0;
                    result = (long)BitConverter.ToUInt64(value, 0);
                }
                finally
                {
                    if (lockTaken3)
                    {
                        Monitor.Exit(taskLock3);
                    }
                }
            }

            return result;
        }
    }

    long _ACI.getDefault_ByUPID(string TargetIP, uint UPID)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getDefault_ByUPID
        return this.getDefault_ByUPID(TargetIP, UPID);
    }

    //
    // Summary:
    //     Parameter channel function: set RAM value by UPID, UPID and scale for raw input
    //     data can be found in LinMot Talk -> Variables
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   UPID:
    //     UPID address of drive variable
    //
    //   Value:
    //     Raw data value to send (must be scaled for UPID)
    //
    // Returns:
    //     Return send raw data, 0 if failed
    public long SetRAM_ByUPID(string TargetIP, uint UPID, long Value)
    {
        long result = 0L;
        byte[] bytes = BitConverter.GetBytes(Value);
        byte[] bytes2 = BitConverter.GetBytes(UPID);
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].RealTimeConfigID = 19;
                    if (Axis[i].RealTimeConfigCommandCount > 14)
                    {
                        Axis[i].RealTimeConfigCommandCount = 1;
                    }

                    Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[i].RealTimeConfigCommandCount + 1));
                    Axis[i].RealTimeConfigCommandCount = Axis[i].RealTimeConfigCommandCount;
                    Axis[i].RealTimeConfigArgs1High = bytes2[1];
                    Axis[i].RealTimeConfigArgs1Low = bytes2[0];
                    Axis[i].RealTimeConfigArgs2High = bytes[1];
                    Axis[i].RealTimeConfigArgs2Low = bytes[0];
                    Axis[i].RealTimeConfigArgs3High = bytes[3];
                    Axis[i].RealTimeConfigArgs3Low = bytes[2];
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                long timeOutTime = getTimeOutTime(1000L);
                while (!((Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount) & (Axis[i].RealTimeConfigIDStatus == 0)))
                {
                    if (!isTimeOut(timeOutTime))
                    {
                        continue;
                    }

                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        Axis[i].DLLErrorText = "12: Config Channel Timeout by writing RAM value by UPID. Error code: " + Conversion.Hex(Axis[i].RealTimeConfigIDStatus).ToString();
                        ACIError = "Config Channel Error";
                        Axis[i].RealTimeConfigStatusArgs2Low = 9;
                        Axis[i].RealTimeConfigStatusArgs2High = 0;
                        Axis[i].RealTimeConfigStatusArgs3Low = 9;
                        Axis[i].RealTimeConfigStatusArgs3High = 0;
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    break;
                }

                object taskLock3 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                bool lockTaken3 = false;
                try
                {
                    Monitor.Enter(taskLock3, ref lockTaken3);
                    byte[] value = new byte[8]
                    {
                        Axis[i].RealTimeConfigStatusArgs2Low,
                        Axis[i].RealTimeConfigStatusArgs2High,
                        Axis[i].RealTimeConfigStatusArgs3Low,
                        Axis[i].RealTimeConfigStatusArgs3High,
                        0,
                        0,
                        0,
                        0
                    };
                    Axis[i].RealTimeConfigID = 0;
                    result = (long)BitConverter.ToUInt64(value, 0);
                }
                finally
                {
                    if (lockTaken3)
                    {
                        Monitor.Exit(taskLock3);
                    }
                }
            }

            return result;
        }
    }

    long _ACI.SetRAM_ByUPID(string TargetIP, uint UPID, long Value)
    {
        //ILSpy generated this explicit interface implementation from .override directive in SetRAM_ByUPID
        return this.SetRAM_ByUPID(TargetIP, UPID, Value);
    }

    //
    // Summary:
    //     Parameter channel function: set ROM value by UPID, UPID and scale for raw input
    //     data can be found in LinMot Talk -> Variables
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   UPID:
    //     UPID address of drive variable
    //
    //   Value:
    //     Raw data value to send (must be scaled for UPID)
    //
    // Returns:
    //     Return send raw data, 0 if failed
    public long SetROM_ByUPID(string TargetIP, uint UPID, long Value)
    {
        long result = 0L;
        byte[] bytes = BitConverter.GetBytes(Value);
        byte[] bytes2 = BitConverter.GetBytes(UPID);
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].RealTimeConfigID = 18;
                    if (Axis[i].RealTimeConfigCommandCount > 14)
                    {
                        Axis[i].RealTimeConfigCommandCount = 1;
                    }

                    Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[i].RealTimeConfigCommandCount + 1));
                    Axis[i].RealTimeConfigCommandCount = Axis[i].RealTimeConfigCommandCount;
                    Axis[i].RealTimeConfigArgs1High = bytes2[1];
                    Axis[i].RealTimeConfigArgs1Low = bytes2[0];
                    Axis[i].RealTimeConfigArgs2High = bytes[1];
                    Axis[i].RealTimeConfigArgs2Low = bytes[0];
                    Axis[i].RealTimeConfigArgs3High = bytes[3];
                    Axis[i].RealTimeConfigArgs3Low = bytes[2];
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                long timeOutTime = getTimeOutTime(1000L);
                while (!((Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount) & (Axis[i].RealTimeConfigIDStatus == 0)))
                {
                    if (!isTimeOut(timeOutTime))
                    {
                        continue;
                    }

                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        Axis[i].DLLErrorText = "13: Config Channel Timeout by writing ROM value by UPID. Error code: " + Conversion.Hex(Axis[i].RealTimeConfigIDStatus).ToString();
                        ACIError = "Config Channel Error";
                        Axis[i].RealTimeConfigStatusArgs2Low = 9;
                        Axis[i].RealTimeConfigStatusArgs2High = 0;
                        Axis[i].RealTimeConfigStatusArgs3Low = 9;
                        Axis[i].RealTimeConfigStatusArgs3High = 0;
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    break;
                }

                object taskLock3 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                bool lockTaken3 = false;
                try
                {
                    Monitor.Enter(taskLock3, ref lockTaken3);
                    byte[] value = new byte[8]
                    {
                        Axis[i].RealTimeConfigStatusArgs2Low,
                        Axis[i].RealTimeConfigStatusArgs2High,
                        Axis[i].RealTimeConfigStatusArgs3Low,
                        Axis[i].RealTimeConfigStatusArgs3High,
                        0,
                        0,
                        0,
                        0
                    };
                    Axis[i].RealTimeConfigID = 0;
                    result = (long)BitConverter.ToUInt64(value, 0);
                }
                finally
                {
                    if (lockTaken3)
                    {
                        Monitor.Exit(taskLock3);
                    }
                }
            }

            return result;
        }
    }

    long _ACI.SetROM_ByUPID(string TargetIP, uint UPID, long Value)
    {
        //ILSpy generated this explicit interface implementation from .override directive in SetROM_ByUPID
        return this.SetROM_ByUPID(TargetIP, UPID, Value);
    }

    //
    // Summary:
    //     Parameter channel function: set RAM and ROM value by UPID, UPID and scale for
    //     raw input data can be found in LinMot Talk -> Variables
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   UPID:
    //     UPID address of drive variable
    //
    //   Value:
    //     Raw data value to send (must be scaled for UPID)
    //
    // Returns:
    //     Return send raw data, 0 if failed
    public long SetRAM_ROM_ByUPID(string TargetIP, uint UPID, long Value)
    {
        long result = 0L;
        byte[] bytes = BitConverter.GetBytes(Value);
        byte[] bytes2 = BitConverter.GetBytes(UPID);
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].RealTimeConfigID = 20;
                    if (Axis[i].RealTimeConfigCommandCount > 14)
                    {
                        Axis[i].RealTimeConfigCommandCount = 1;
                    }

                    Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[i].RealTimeConfigCommandCount + 1));
                    Axis[i].RealTimeConfigCommandCount = Axis[i].RealTimeConfigCommandCount;
                    Axis[i].RealTimeConfigArgs1High = bytes2[1];
                    Axis[i].RealTimeConfigArgs1Low = bytes2[0];
                    Axis[i].RealTimeConfigArgs2High = bytes[1];
                    Axis[i].RealTimeConfigArgs2Low = bytes[0];
                    Axis[i].RealTimeConfigArgs3High = bytes[3];
                    Axis[i].RealTimeConfigArgs3Low = bytes[2];
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                long timeOutTime = getTimeOutTime(1000L);
                while (!((Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount) & (Axis[i].RealTimeConfigIDStatus == 0)))
                {
                    if (!isTimeOut(timeOutTime))
                    {
                        continue;
                    }

                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        Axis[i].DLLErrorText = "14: Config Channel Timeout by writing RAM and ROM value by UPID. Error code: " + Conversion.Hex(Axis[i].RealTimeConfigIDStatus).ToString();
                        ACIError = "Config Channel Error";
                        Axis[i].RealTimeConfigStatusArgs2Low = 9;
                        Axis[i].RealTimeConfigStatusArgs2High = 0;
                        Axis[i].RealTimeConfigStatusArgs3Low = 9;
                        Axis[i].RealTimeConfigStatusArgs3High = 0;
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    break;
                }

                object taskLock3 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                bool lockTaken3 = false;
                try
                {
                    Monitor.Enter(taskLock3, ref lockTaken3);
                    byte[] value = new byte[8]
                    {
                        Axis[i].RealTimeConfigStatusArgs2Low,
                        Axis[i].RealTimeConfigStatusArgs2High,
                        Axis[i].RealTimeConfigStatusArgs3Low,
                        Axis[i].RealTimeConfigStatusArgs3High,
                        0,
                        0,
                        0,
                        0
                    };
                    Axis[i].RealTimeConfigID = 0;
                    result = (long)BitConverter.ToUInt64(value, 0);
                }
                finally
                {
                    if (lockTaken3)
                    {
                        Monitor.Exit(taskLock3);
                    }
                }
            }

            return result;
        }
    }

    long _ACI.SetRAM_ROM_ByUPID(string TargetIP, uint UPID, long Value)
    {
        //ILSpy generated this explicit interface implementation from .override directive in SetRAM_ROM_ByUPID
        return this.SetRAM_ROM_ByUPID(TargetIP, UPID, Value);
    }

    //
    // Summary:
    //     Make movement to target position from actual position with actual velocity
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   Pos1:
    //     Target position [mm]
    //
    //   MaxVel1:
    //     Maximum velocity [m/s]
    //
    //   Acc1:
    //     Acceleration [m/s2]
    //
    //   Dec1:
    //     Deceleration [m/s2]
    //
    // Returns:
    //     Returns true
    public bool LMmt_GoToPosFromActPosAndActVel(string TargetIP, float Pos1, float MaxVel1, float Acc1, float Dec1)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    int value = (int)Math.Round(Pos1 * 10000f);
                    uint value2 = (uint)Math.Round(MaxVel1 * 1000000f);
                    uint value3 = (uint)Math.Round(Acc1 * 100000f);
                    uint value4 = (uint)Math.Round(Dec1 * 100000f);
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 1;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(3 | Axis[i].MotionCommandInterface.CountNibble);
                    byte[] bytes = BitConverter.GetBytes(value);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value2);
                    Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value3);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value4);
                    Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter4High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter4Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes[0];
                    result = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMmt_GoToPosFromActPosAndActVel(string TargetIP, float Pos1, float MaxVel1, float Acc1, float Dec1)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GoToPosFromActPosAndActVel
        return this.LMmt_GoToPosFromActPosAndActVel(TargetIP, Pos1, MaxVel1, Acc1, Dec1);
    }

    //
    // Summary:
    //     Make absolute movement to target position
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   Pos1:
    //     Target position [mm]
    //
    //   MaxVel1:
    //     Maximum velocity [m/s]
    //
    //   Acc1:
    //     Acceleration [m/s2]
    //
    //   Dec1:
    //     Deceleration [m/s2]
    //
    // Returns:
    //     Returns true
    public bool LMmt_MoveAbs(string TargetIP, float Pos1, float MaxVel1, float Acc1, float Dec1)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    int value = (int)Math.Round(Pos1 * 10000f);
                    uint value2 = (uint)Math.Round(MaxVel1 * 1000000f);
                    uint value3 = (uint)Math.Round(Acc1 * 100000f);
                    uint value4 = (uint)Math.Round(Dec1 * 100000f);
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 1;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = Axis[i].MotionCommandInterface.CountNibble;
                    byte[] bytes = BitConverter.GetBytes(value);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value2);
                    Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value3);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value4);
                    Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter4High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter4Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes[0];
                    result = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMmt_MoveAbs(string TargetIP, float Pos1, float MaxVel1, float Acc1, float Dec1)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_MoveAbs
        return this.LMmt_MoveAbs(TargetIP, Pos1, MaxVel1, Acc1, Dec1);
    }

    //
    // Summary:
    //     Make relative motion from actual demand position
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   Pos1:
    //     Target position [mm]
    //
    //   MaxVel1:
    //     Maximum velocity [m/s]
    //
    //   Acc1:
    //     Acceleration [m/s2]
    //
    //   Dec1:
    //     Deceleration [m/s2]
    //
    // Returns:
    //     Returns true
    public bool LMmt_MoveRel(string TargetIP, float Pos1, float MaxVel1, float Acc1, float Dec1)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    int value = (int)Math.Round(Pos1 * 10000f);
                    int value2 = (int)Math.Round(MaxVel1 * 1000000f);
                    int value3 = (int)Math.Round(Acc1 * 100000f);
                    int value4 = (int)Math.Round(Dec1 * 100000f);
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 1;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(0x10 | Axis[i].MotionCommandInterface.CountNibble);
                    byte[] bytes = BitConverter.GetBytes(value);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value2);
                    Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value3);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value4);
                    Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter4High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter4Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes[0];
                    result = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMmt_MoveRel(string TargetIP, float Pos1, float MaxVel1, float Acc1, float Dec1)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_MoveRel
        return this.LMmt_MoveRel(TargetIP, Pos1, MaxVel1, Acc1, Dec1);
    }

    //
    // Summary:
    //     Stop actual motion by using defined deceleration
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   Decceleration:
    //     Deceleration [m/s2]
    //
    // Returns:
    //     Returns true
    public bool LMmt_Stop(string TargetIP, float Decceleration)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    int value = (int)Math.Round(Decceleration * 100000f);
                    int value2 = 0;
                    int value3 = 0;
                    int value4 = 0;
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 1;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(0x70 | Axis[i].MotionCommandInterface.CountNibble);
                    byte[] bytes = BitConverter.GetBytes(value);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value2);
                    Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value3);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value4);
                    Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter4High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter4Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes[0];
                    result = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMmt_Stop(string TargetIP, float Decceleration)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_Stop
        return this.LMmt_Stop(TargetIP, Decceleration);
    }

    //
    // Summary:
    //     Write live parameter over MC interface
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   UPID:
    //     UPID address which should be accessed
    //
    //   UPIDValue:
    //     Raw data value, must be scaled according to the UPID
    //
    // Returns:
    //     Returns true
    public bool LMmt_WriteLivePar(string TargetIP, uint UPID, int UPIDValue)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    int value = (int)UPID;
                    int value2 = 0;
                    int value3 = 0;
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 0;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(0x20 | Axis[i].MotionCommandInterface.CountNibble);
                    byte[] bytes = BitConverter.GetBytes(value);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(UPIDValue);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[0];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[2];
                    bytes = BitConverter.GetBytes(value2);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value3);
                    Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter4High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter4Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes[0];
                    result = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMmt_WriteLivePar(string TargetIP, uint UPID, int UPIDValue)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_WriteLivePar
        return this.LMmt_WriteLivePar(TargetIP, UPID, UPIDValue);
    }

    //
    // Summary:
    //     Make absolute movement to target position with Jerk (3A0xh)
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   Pos:
    //     Target position [mm]
    //
    //   MaxVel:
    //     Maximum velocity [m/s]
    //
    //   Acc:
    //     Acceleration [m/s2]
    //
    //   Dec:
    //     Deceleration [m/s2]
    //
    // Returns:
    //     Returns true
    public bool LMmt_VAJIGoToPos(string TargetIP, float Pos, float MaxVel, float Acc, float Dec, float Jerk)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    int value = (int)Math.Round(Pos * 10000f);
                    uint value2 = (uint)Math.Round(MaxVel * 1000000f);
                    uint value3 = (uint)Math.Round(Acc * 100000f);
                    uint value4 = (uint)Math.Round(Dec * 100000f);
                    uint value5 = (uint)Math.Round(Jerk * 10000f);
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 58;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(0x10 | Axis[i].MotionCommandInterface.CountNibble);
                    byte[] bytes = BitConverter.GetBytes(value);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value2);
                    Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value3);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value4);
                    Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter4High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter4Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value5);
                    Axis[i].MotionCommandInterface.CommandParameter5HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter5High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter5Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter5LowLow = bytes[0];
                    result = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMmt_VAJIGoToPos(string TargetIP, float Pos, float MaxVel, float Acc, float Dec, float Jerk)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_VAJIGoToPos
        return this.LMmt_VAJIGoToPos(TargetIP, Pos, MaxVel, Acc, Dec, Jerk);
    }

    //
    // Summary:
    //     Make movement to target position from actual position with actual velocity
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   Pos1:
    //     Target position [mm]
    //
    //   MaxVel1:
    //     Maximum velocity [m/s]
    //
    //   Acc1:
    //     Acceleration [m/s2]
    //
    //   Dec1:
    //     Deceleration [m/s2]
    //
    // Returns:
    //     Returns true
    public bool LMmt_IncrementActPosStartingWithDemVel0ResetI(string TargetIP, float Pos1, float MaxVel1, float Acc1, float Dec1)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    int value = (int)Math.Round(Pos1 * 10000f);
                    uint value2 = (uint)Math.Round(MaxVel1 * 1000000f);
                    uint value3 = (uint)Math.Round(Acc1 * 100000f);
                    uint value4 = (uint)Math.Round(Dec1 * 100000f);
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 13;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(0x90 | Axis[i].MotionCommandInterface.CountNibble);
                    byte[] bytes = BitConverter.GetBytes(value);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value2);
                    Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value3);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value4);
                    Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter4High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter4Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes[0];
                    result = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMmt_IncrementActPosStartingWithDemVel0ResetI(string TargetIP, float Pos1, float MaxVel1, float Acc1, float Dec1)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_IncrementActPosStartingWithDemVel0ResetI
        return this.LMmt_IncrementActPosStartingWithDemVel0ResetI(TargetIP, Pos1, MaxVel1, Acc1, Dec1);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as Int16
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, short MCParaWord0)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter1High = 0;
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter2High = 0;
                Axis[i].MotionCommandInterface.CommandParameter2Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter3High = 0;
                Axis[i].MotionCommandInterface.CommandParameter3Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, short MCParaWord0)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as SInt16
    //
    //   MCParaWord1:
    //     Parameter value as UInt16
    //
    //   MCParaWord2:
    //     Parameter value as UInt16
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, short MCParaWord0, ushort MCParaWord1, ushort MCParaWord2)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                byte[] bytes3 = BitConverter.GetBytes(MCParaWord2);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter2High = 0;
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes3[1];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes3[0];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter3High = 0;
                Axis[i].MotionCommandInterface.CommandParameter3Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, short MCParaWord0, ushort MCParaWord1, ushort MCParaWord2)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1, MCParaWord2);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as SInt16
    //
    //   MCParaWord1:
    //     Parameter value as SInt32
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, short MCParaWord0, int MCParaWord1)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter2High = 0;
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes2[3];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes2[2];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter3High = 0;
                Axis[i].MotionCommandInterface.CommandParameter3Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, short MCParaWord0, int MCParaWord1)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as SInt16
    //
    //   MCParaWord1:
    //     Parameter value as UInt16
    //
    //   MCParaWord2:
    //     Parameter value as UInt16
    //
    //   MCParaWord3:
    //     Parameter value as UInt16
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, short MCParaWord0, ushort MCParaWord1, ushort MCParaWord2, ushort MCParaWord3)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                byte[] bytes3 = BitConverter.GetBytes(MCParaWord2);
                byte[] bytes4 = BitConverter.GetBytes(MCParaWord3);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes4[1];
                Axis[i].MotionCommandInterface.CommandParameter2High = bytes4[0];
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes3[1];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes3[0];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter3High = 0;
                Axis[i].MotionCommandInterface.CommandParameter3Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, short MCParaWord0, ushort MCParaWord1, ushort MCParaWord2, ushort MCParaWord3)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1, MCParaWord2, MCParaWord3);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as UInt16
    //
    //   MCParaWord1:
    //     Parameter value as UInt16
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, ushort MCParaWord1)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter2High = 0;
                Axis[i].MotionCommandInterface.CommandParameter2Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter3High = 0;
                Axis[i].MotionCommandInterface.CommandParameter3Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, ushort MCParaWord1)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as UInt16
    //
    //   MCParaWord1:
    //     Parameter value as SInt16
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, short MCParaWord1)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter2High = 0;
                Axis[i].MotionCommandInterface.CommandParameter2Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter3High = 0;
                Axis[i].MotionCommandInterface.CommandParameter3Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, short MCParaWord1)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as UInt16
    //
    //   MCParaWord1:
    //     Parameter value as UInt32
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, uint MCParaWord1)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter2High = 0;
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes2[3];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes2[2];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter3High = 0;
                Axis[i].MotionCommandInterface.CommandParameter3Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, uint MCParaWord1)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as UInt16
    //
    //   MCParaWord1:
    //     Parameter value as Int32
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, int MCParaWord1)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter2High = 0;
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes2[3];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes2[2];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter3High = 0;
                Axis[i].MotionCommandInterface.CommandParameter3Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, int MCParaWord1)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as UInt16
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter1High = 0;
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter2High = 0;
                Axis[i].MotionCommandInterface.CommandParameter2Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter3High = 0;
                Axis[i].MotionCommandInterface.CommandParameter3Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as UInt16
    //
    //   MCParaWord1:
    //     Parameter value as SInt32
    //
    //   MCParaWord2:
    //     Parameter value as SInt32
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, int MCParaWord1, int MCParaWord2)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                byte[] bytes3 = BitConverter.GetBytes(MCParaWord2);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes3[1];
                Axis[i].MotionCommandInterface.CommandParameter2High = bytes3[0];
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes2[3];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes2[2];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter3High = 0;
                Axis[i].MotionCommandInterface.CommandParameter3Low = bytes3[3];
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes3[2];
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, int MCParaWord1, int MCParaWord2)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1, MCParaWord2);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as UInt16
    //
    //   MCParaWord1:
    //     Parameter value as UInt16
    //
    //   MCParaWord2:
    //     Parameter value as SInt16
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, ushort MCParaWord1, short MCParaWord2)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                byte[] bytes3 = BitConverter.GetBytes(MCParaWord2);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter2High = 0;
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes3[1];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes3[0];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter3High = 0;
                Axis[i].MotionCommandInterface.CommandParameter3Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, ushort MCParaWord1, short MCParaWord2)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1, MCParaWord2);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as UInt16
    //
    //   MCParaWord1:
    //     Parameter value as UInt16
    //
    //   MCParaWord2:
    //     Parameter value as SInt32
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, ushort MCParaWord1, int MCParaWord2)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                byte[] bytes3 = BitConverter.GetBytes(MCParaWord2);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes3[3];
                Axis[i].MotionCommandInterface.CommandParameter2High = bytes3[2];
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes3[1];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes3[0];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter3High = 0;
                Axis[i].MotionCommandInterface.CommandParameter3Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, ushort MCParaWord1, int MCParaWord2)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1, MCParaWord2);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as UInt16
    //
    //   MCParaWord1:
    //     Parameter value as SInt32
    //
    //   MCParaWord2:
    //     Parameter value as UInt16
    //
    //   MCParaWord3:
    //     Parameter value as SInt16
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, int MCParaWord1, ushort MCParaWord2, short MCParaWord3)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                byte[] bytes3 = BitConverter.GetBytes(MCParaWord2);
                byte[] bytes4 = BitConverter.GetBytes(MCParaWord3);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes3[1];
                Axis[i].MotionCommandInterface.CommandParameter2High = bytes3[0];
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes2[3];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes2[2];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter3High = 0;
                Axis[i].MotionCommandInterface.CommandParameter3Low = bytes4[1];
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes4[0];
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, int MCParaWord1, ushort MCParaWord2, short MCParaWord3)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1, MCParaWord2, MCParaWord3);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as UInt16
    //
    //   MCParaWord1:
    //     Parameter value as SInt32
    //
    //   MCParaWord2:
    //     Parameter value as SInt32
    //
    //   MCParaWord3:
    //     Parameter value as SInt16
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, int MCParaWord1, uint MCParaWord2, short MCParaWord3)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                byte[] bytes3 = BitConverter.GetBytes(MCParaWord2);
                byte[] bytes4 = BitConverter.GetBytes(MCParaWord3);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes3[1];
                Axis[i].MotionCommandInterface.CommandParameter2High = bytes3[0];
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes2[3];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes2[2];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes4[1];
                Axis[i].MotionCommandInterface.CommandParameter3High = bytes4[0];
                Axis[i].MotionCommandInterface.CommandParameter3Low = bytes3[3];
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes3[2];
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, int MCParaWord1, uint MCParaWord2, short MCParaWord3)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1, MCParaWord2, MCParaWord3);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as UInt16
    //
    //   MCParaWord1:
    //     Parameter value as SInt32
    //
    //   MCParaWord2:
    //     Parameter value as SInt32
    //
    //   MCParaWord3:
    //     Parameter value as SInt32
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, int MCParaWord1, int MCParaWord2, int MCParaWord3)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                byte[] bytes3 = BitConverter.GetBytes(MCParaWord2);
                byte[] bytes4 = BitConverter.GetBytes(MCParaWord3);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes3[1];
                Axis[i].MotionCommandInterface.CommandParameter2High = bytes3[0];
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes2[3];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes2[2];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes4[1];
                Axis[i].MotionCommandInterface.CommandParameter3High = bytes4[0];
                Axis[i].MotionCommandInterface.CommandParameter3Low = bytes3[3];
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes3[2];
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = bytes4[3];
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes4[2];
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, int MCParaWord1, int MCParaWord2, int MCParaWord3)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1, MCParaWord2, MCParaWord3);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as UInt16
    //
    //   MCParaWord1:
    //     Parameter value as SInt32
    //
    //   MCParaWord2:
    //     Parameter value as SInt16
    //
    //   MCParaWord3:
    //     Parameter value as SInt32
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, int MCParaWord1, short MCParaWord2, int MCParaWord3)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                byte[] bytes3 = BitConverter.GetBytes(MCParaWord2);
                byte[] bytes4 = BitConverter.GetBytes(MCParaWord3);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes3[1];
                Axis[i].MotionCommandInterface.CommandParameter2High = bytes3[0];
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes2[3];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes2[2];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes4[3];
                Axis[i].MotionCommandInterface.CommandParameter3High = bytes4[2];
                Axis[i].MotionCommandInterface.CommandParameter3Low = bytes4[1];
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes4[0];
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, int MCParaWord1, short MCParaWord2, int MCParaWord3)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1, MCParaWord2, MCParaWord3);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as UInt16
    //
    //   MCParaWord1:
    //     Parameter value as UInt16
    //
    //   MCParaWord2:
    //     Parameter value as SInt32
    //
    //   MCParaWord3:
    //     Parameter value as SInt32
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, ushort MCParaWord1, int MCParaWord2, int MCParaWord3)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                byte[] bytes3 = BitConverter.GetBytes(MCParaWord2);
                byte[] bytes4 = BitConverter.GetBytes(MCParaWord3);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes3[3];
                Axis[i].MotionCommandInterface.CommandParameter2High = bytes3[2];
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes3[1];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes3[0];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes4[3];
                Axis[i].MotionCommandInterface.CommandParameter3High = bytes4[2];
                Axis[i].MotionCommandInterface.CommandParameter3Low = bytes4[1];
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes4[0];
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, ushort MCParaWord1, int MCParaWord2, int MCParaWord3)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1, MCParaWord2, MCParaWord3);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as UInt16
    //
    //   MCParaWord1:
    //     Parameter value as UInt16
    //
    //   MCParaWord2:
    //     Parameter value as UInt16
    //
    //   MCParaWord3:
    //     Parameter value as UInt16
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, ushort MCParaWord1, ushort MCParaWord2, ushort MCParaWord3)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                byte[] bytes3 = BitConverter.GetBytes(MCParaWord2);
                byte[] bytes4 = BitConverter.GetBytes(MCParaWord3);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes4[1];
                Axis[i].MotionCommandInterface.CommandParameter2High = bytes4[0];
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes3[1];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes3[0];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter3High = 0;
                Axis[i].MotionCommandInterface.CommandParameter3Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, ushort MCParaWord1, ushort MCParaWord2, ushort MCParaWord3)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1, MCParaWord2, MCParaWord3);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as UInt16
    //
    //   MCParaWord1:
    //     Parameter value as SInt32
    //
    //   MCParaWord2:
    //     Parameter value as SInt32
    //
    //   MCParaWord3:
    //     Parameter value as SInt32
    //
    //   MCParaWord4:
    //     Parameter value as SInt32
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, int MCParaWord1, int MCParaWord2, int MCParaWord3, int MCParaWord4)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                byte[] bytes3 = BitConverter.GetBytes(MCParaWord2);
                byte[] bytes4 = BitConverter.GetBytes(MCParaWord3);
                byte[] bytes5 = BitConverter.GetBytes(MCParaWord4);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes3[1];
                Axis[i].MotionCommandInterface.CommandParameter2High = bytes3[0];
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes2[3];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes2[2];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes4[1];
                Axis[i].MotionCommandInterface.CommandParameter3High = bytes4[0];
                Axis[i].MotionCommandInterface.CommandParameter3Low = bytes3[3];
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes3[2];
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes5[1];
                Axis[i].MotionCommandInterface.CommandParameter4High = bytes5[0];
                Axis[i].MotionCommandInterface.CommandParameter4Low = bytes4[3];
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes4[2];
                Axis[i].MotionCommandInterface.CommandParameter5HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter5High = 0;
                Axis[i].MotionCommandInterface.CommandParameter5Low = bytes5[3];
                Axis[i].MotionCommandInterface.CommandParameter5LowLow = bytes5[2];
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, int MCParaWord1, int MCParaWord2, int MCParaWord3, int MCParaWord4)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1, MCParaWord2, MCParaWord3, MCParaWord4);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as Int32
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter2High = 0;
                Axis[i].MotionCommandInterface.CommandParameter2Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter3High = 0;
                Axis[i].MotionCommandInterface.CommandParameter3Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as SInt32
    //
    //   MCParaWord1:
    //     Parameter value as SInt32
    //
    //   MCParaWord2:
    //     Parameter value as SInt32
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, int MCParaWord1, int MCParaWord2)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                byte[] bytes3 = BitConverter.GetBytes(MCParaWord2);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes2[3];
                Axis[i].MotionCommandInterface.CommandParameter2High = bytes2[2];
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes3[3];
                Axis[i].MotionCommandInterface.CommandParameter3High = bytes3[2];
                Axis[i].MotionCommandInterface.CommandParameter3Low = bytes3[1];
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes3[0];
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, int MCParaWord1, int MCParaWord2)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1, MCParaWord2);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as SInt32
    //
    //   MCParaWord1:
    //     Parameter value as UInt32
    //
    //   MCParaWord2:
    //     Parameter value as UInt32
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, uint MCParaWord1, uint MCParaWord2)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                byte[] bytes3 = BitConverter.GetBytes(MCParaWord2);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes2[3];
                Axis[i].MotionCommandInterface.CommandParameter2High = bytes2[2];
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes3[3];
                Axis[i].MotionCommandInterface.CommandParameter3High = bytes3[2];
                Axis[i].MotionCommandInterface.CommandParameter3Low = bytes3[1];
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes3[0];
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, uint MCParaWord1, uint MCParaWord2)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1, MCParaWord2);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as SInt32
    //
    //   MCParaWord1:
    //     Parameter value as UInt16
    //
    //   MCParaWord2:
    //     Parameter value as UInt16
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, ushort MCParaWord1, ushort MCParaWord2)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                byte[] bytes3 = BitConverter.GetBytes(MCParaWord2);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes3[1];
                Axis[i].MotionCommandInterface.CommandParameter2High = bytes3[0];
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter3High = 0;
                Axis[i].MotionCommandInterface.CommandParameter3Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, ushort MCParaWord1, ushort MCParaWord2)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1, MCParaWord2);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as SInt32
    //
    //   MCParaWord1:
    //     Parameter value as UInt16
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, ushort MCParaWord1)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter2High = 0;
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter3High = 0;
                Axis[i].MotionCommandInterface.CommandParameter3Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, ushort MCParaWord1)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as SInt32
    //
    //   MCParaWord1:
    //     Parameter value as UInt32
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, uint MCParaWord1)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes2[3];
                Axis[i].MotionCommandInterface.CommandParameter2High = bytes2[2];
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter3High = 0;
                Axis[i].MotionCommandInterface.CommandParameter3Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, uint MCParaWord1)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as SInt32
    //
    //   MCParaWord1:
    //     Parameter value as SInt32
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, int MCParaWord1)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes2[3];
                Axis[i].MotionCommandInterface.CommandParameter2High = bytes2[2];
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter3High = 0;
                Axis[i].MotionCommandInterface.CommandParameter3Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, int MCParaWord1)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as SInt32
    //
    //   MCParaWord1:
    //     Parameter value as UInt32
    //
    //   MCParaWord2:
    //     Parameter value as UInt32
    //
    //   MCParaWord3:
    //     Parameter value as UInt32
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, uint MCParaWord1, uint MCParaWord2, uint MCParaWord3)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                byte[] bytes3 = BitConverter.GetBytes(MCParaWord2);
                byte[] bytes4 = BitConverter.GetBytes(MCParaWord3);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes2[3];
                Axis[i].MotionCommandInterface.CommandParameter2High = bytes2[2];
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes3[3];
                Axis[i].MotionCommandInterface.CommandParameter3High = bytes3[2];
                Axis[i].MotionCommandInterface.CommandParameter3Low = bytes3[1];
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes3[0];
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes4[3];
                Axis[i].MotionCommandInterface.CommandParameter4High = bytes4[2];
                Axis[i].MotionCommandInterface.CommandParameter4Low = bytes4[1];
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes4[0];
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, uint MCParaWord1, uint MCParaWord2, uint MCParaWord3)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1, MCParaWord2, MCParaWord3);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as SInt32
    //
    //   MCParaWord1:
    //     Parameter value as UInt32
    //
    //   MCParaWord2:
    //     Parameter value as UInt32
    //
    //   MCParaWord3:
    //     Parameter value as UInt16
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, uint MCParaWord1, uint MCParaWord2, ushort MCParaWord3)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                byte[] bytes3 = BitConverter.GetBytes(MCParaWord2);
                byte[] bytes4 = BitConverter.GetBytes(MCParaWord3);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes2[3];
                Axis[i].MotionCommandInterface.CommandParameter2High = bytes2[2];
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes3[3];
                Axis[i].MotionCommandInterface.CommandParameter3High = bytes3[2];
                Axis[i].MotionCommandInterface.CommandParameter3Low = bytes3[1];
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes3[0];
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = bytes4[1];
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes4[0];
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, uint MCParaWord1, uint MCParaWord2, ushort MCParaWord3)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1, MCParaWord2, MCParaWord3);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as SInt32
    //
    //   MCParaWord1:
    //     Parameter value as UInt32
    //
    //   MCParaWord2:
    //     Parameter value as UInt32
    //
    //   MCParaWord3:
    //     Parameter value as Int16
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, uint MCParaWord1, uint MCParaWord2, short MCParaWord3)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                byte[] bytes3 = BitConverter.GetBytes(MCParaWord2);
                byte[] bytes4 = BitConverter.GetBytes(MCParaWord3);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes2[3];
                Axis[i].MotionCommandInterface.CommandParameter2High = bytes2[2];
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes3[3];
                Axis[i].MotionCommandInterface.CommandParameter3High = bytes3[2];
                Axis[i].MotionCommandInterface.CommandParameter3Low = bytes3[1];
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes3[0];
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = bytes4[1];
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes4[0];
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, uint MCParaWord1, uint MCParaWord2, short MCParaWord3)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1, MCParaWord2, MCParaWord3);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as SInt32
    //
    //   MCParaWord1:
    //     Parameter value as UInt32
    //
    //   MCParaWord2:
    //     Parameter value as UInt32
    //
    //   MCParaWord3:
    //     Parameter value as UInt32
    //
    //   MCParaWord4:
    //     Parameter value as UInt32
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, uint MCParaWord1, uint MCParaWord2, uint MCParaWord3, uint MCParaWord4)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                byte[] bytes3 = BitConverter.GetBytes(MCParaWord2);
                byte[] bytes4 = BitConverter.GetBytes(MCParaWord3);
                byte[] bytes5 = BitConverter.GetBytes(MCParaWord4);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes2[3];
                Axis[i].MotionCommandInterface.CommandParameter2High = bytes2[2];
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes3[3];
                Axis[i].MotionCommandInterface.CommandParameter3High = bytes3[2];
                Axis[i].MotionCommandInterface.CommandParameter3Low = bytes3[1];
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes3[0];
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes4[3];
                Axis[i].MotionCommandInterface.CommandParameter4High = bytes4[2];
                Axis[i].MotionCommandInterface.CommandParameter4Low = bytes4[1];
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes4[0];
                Axis[i].MotionCommandInterface.CommandParameter5HighHigh = bytes5[3];
                Axis[i].MotionCommandInterface.CommandParameter5High = bytes5[2];
                Axis[i].MotionCommandInterface.CommandParameter5Low = bytes5[1];
                Axis[i].MotionCommandInterface.CommandParameter5LowLow = bytes5[0];
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, uint MCParaWord1, uint MCParaWord2, uint MCParaWord3, uint MCParaWord4)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1, MCParaWord2, MCParaWord3, MCParaWord4);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as SInt32
    //
    //   MCParaWord1:
    //     Parameter value as UInt32
    //
    //   MCParaWord2:
    //     Parameter value as UInt32
    //
    //   MCParaWord3:
    //     Parameter value as SInt16
    //
    //   MCParaWord4:
    //     Parameter value as SInt16
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, uint MCParaWord1, uint MCParaWord2, short MCParaWord3, short MCParaWord4)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                byte[] bytes3 = BitConverter.GetBytes(MCParaWord2);
                byte[] bytes4 = BitConverter.GetBytes(MCParaWord3);
                byte[] bytes5 = BitConverter.GetBytes(MCParaWord4);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes2[3];
                Axis[i].MotionCommandInterface.CommandParameter2High = bytes2[2];
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes3[3];
                Axis[i].MotionCommandInterface.CommandParameter3High = bytes3[2];
                Axis[i].MotionCommandInterface.CommandParameter3Low = bytes3[1];
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes3[0];
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes5[1];
                Axis[i].MotionCommandInterface.CommandParameter4High = bytes5[0];
                Axis[i].MotionCommandInterface.CommandParameter4Low = bytes4[1];
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes4[0];
                Axis[i].MotionCommandInterface.CommandParameter5HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter5High = 0;
                Axis[i].MotionCommandInterface.CommandParameter5Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter5LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, int MCParaWord0, uint MCParaWord1, uint MCParaWord2, short MCParaWord3, short MCParaWord4)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1, MCParaWord2, MCParaWord3, MCParaWord4);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as UInt32
    //
    //   MCParaWord1:
    //     Parameter value as UInt32
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, uint MCParaWord0, uint MCParaWord1)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes2[3];
                Axis[i].MotionCommandInterface.CommandParameter2High = bytes2[2];
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter3High = 0;
                Axis[i].MotionCommandInterface.CommandParameter3Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, uint MCParaWord0, uint MCParaWord1)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as UInt32
    //
    //   MCParaWord1:
    //     Parameter value as UInt32
    //
    //   MCParaWord2:
    //     Parameter value as UInt32
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, uint MCParaWord0, uint MCParaWord1, uint MCParaWord2)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                byte[] bytes3 = BitConverter.GetBytes(MCParaWord2);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes2[2];
                Axis[i].MotionCommandInterface.CommandParameter2High = bytes2[2];
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes3[3];
                Axis[i].MotionCommandInterface.CommandParameter3High = bytes3[2];
                Axis[i].MotionCommandInterface.CommandParameter3Low = bytes3[1];
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes3[0];
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, uint MCParaWord0, uint MCParaWord1, uint MCParaWord2)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1, MCParaWord2);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as UInt32
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, uint MCParaWord0)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter2High = 0;
                Axis[i].MotionCommandInterface.CommandParameter2Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter3High = 0;
                Axis[i].MotionCommandInterface.CommandParameter3Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = 0;
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = 0;
                Axis[i].MotionCommandInterface.CommandParameter4High = 0;
                Axis[i].MotionCommandInterface.CommandParameter4Low = 0;
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = 0;
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, uint MCParaWord0)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0);
    }

    //
    // Summary:
    //     Generic function, which allow use of all supported MC commands. See Motion Control
    //     Software manual for MC commands and parameters!
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   MCHeader:
    //     Motion command value (count nibble is handled automatical, set it to 0)
    //
    //   MCParaWord0:
    //     Parameter value as UInt32
    //
    // Returns:
    //     Returns true
    public bool LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, ushort MCParaWord1, int MCParaWord2, int MCParaWord3, int MCParaWord4)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                Axis[i].MotionCommandInterface.CountNibble = checked((byte)((Axis[i].StateVarLow & 0xF) + 1));
                if (Axis[i].MotionCommandInterface.CountNibble > 4)
                {
                    Axis[i].MotionCommandInterface.CountNibble = 1;
                }

                byte[] bytes = BitConverter.GetBytes(MCHeader);
                Axis[i].MotionCommandInterface.CommandHeaderMasterID = bytes[1];
                Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(bytes[0] | Axis[i].MotionCommandInterface.CountNibble);
                bytes = BitConverter.GetBytes(MCParaWord0);
                byte[] bytes2 = BitConverter.GetBytes(MCParaWord1);
                byte[] bytes3 = BitConverter.GetBytes(MCParaWord2);
                byte[] bytes4 = BitConverter.GetBytes(MCParaWord3);
                byte[] bytes5 = BitConverter.GetBytes(MCParaWord4);
                Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes2[1];
                Axis[i].MotionCommandInterface.CommandParameter1High = bytes2[0];
                Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes3[3];
                Axis[i].MotionCommandInterface.CommandParameter2High = bytes3[2];
                Axis[i].MotionCommandInterface.CommandParameter2Low = bytes3[1];
                Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes3[0];
                Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes4[3];
                Axis[i].MotionCommandInterface.CommandParameter3High = bytes4[2];
                Axis[i].MotionCommandInterface.CommandParameter3Low = bytes4[1];
                Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes4[0];
                Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes5[3];
                Axis[i].MotionCommandInterface.CommandParameter4High = bytes5[2];
                Axis[i].MotionCommandInterface.CommandParameter4Low = bytes5[1];
                Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes5[0];
                result = true;
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    bool _ACI.LMmt_GenericMC(string TargetIP, ushort MCHeader, ushort MCParaWord0, ushort MCParaWord1, int MCParaWord2, int MCParaWord3, int MCParaWord4)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_GenericMC
        return this.LMmt_GenericMC(TargetIP, MCHeader, MCParaWord0, MCParaWord1, MCParaWord2, MCParaWord3, MCParaWord4);
    }

    //
    // Summary:
    //     Modify 16 bit command table parameter (RAM only)
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   CTEntryID:
    //     Entry ID, line number where change should be done
    //
    //   ParaOffset:
    //     parameter offset, see DLL manual and MC software manual
    //
    //   ParaValue:
    //     Raw data of value to write into command table RAM
    //
    // Returns:
    //     Returns true
    public bool LMav_Mod16BitCTPar(string TargetIP, int CTEntryID, int ParaOffset, int ParaValue)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 32;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(0x80 | Axis[i].MotionCommandInterface.CountNibble);
                    byte[] bytes = BitConverter.GetBytes((uint)CTEntryID);
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes((uint)ParaOffset);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[0];
                    bytes = BitConverter.GetBytes(ParaValue);
                    Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(0);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(0);
                    Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter4High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter4Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes[0];
                    result = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMav_Mod16BitCTPar(string TargetIP, int CTEntryID, int ParaOffset, int ParaValue)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMav_Mod16BitCTPar
        return this.LMav_Mod16BitCTPar(TargetIP, CTEntryID, ParaOffset, ParaValue);
    }

    //
    // Summary:
    //     Modify 32 bit command table parameter (RAM only)
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   CTEntryID:
    //     Entry ID, line number where change should be done
    //
    //   ParaOffset:
    //     parameter offset, see DLL manual and MC software manual
    //
    //   ParaValue:
    //     Raw data of value to write into command table RAM
    //
    // Returns:
    //     Returns true
    public bool LMav_Mod32BitCTPar(string TargetIP, int CTEntryID, int ParaOffset, int ParaValue)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 32;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(0x90 | Axis[i].MotionCommandInterface.CountNibble);
                    byte[] bytes = BitConverter.GetBytes((uint)CTEntryID);
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes((uint)ParaOffset);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[0];
                    bytes = BitConverter.GetBytes(ParaValue);
                    Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(0);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(0);
                    Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter4High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter4Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes[0];
                    result = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMav_Mod32BitCTPar(string TargetIP, int CTEntryID, int ParaOffset, int ParaValue)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMav_Mod32BitCTPar
        return this.LMav_Mod32BitCTPar(TargetIP, CTEntryID, ParaOffset, ParaValue);
    }

    //
    // Summary:
    //     Start command table by defined entry line (A running command table can be checkd
    //     by isEventHandler!)
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   CTEntryID:
    //     >Entry ID, line number where command table should start
    //
    // Returns:
    //     Returns true
    public bool LMmt_StartCTCommand(string TargetIP, uint CTEntryID)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    int value = (int)CTEntryID;
                    int value2 = 0;
                    int value3 = 0;
                    int value4 = 0;
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 32;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(0 | Axis[i].MotionCommandInterface.CountNibble);
                    byte[] bytes = BitConverter.GetBytes(value);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value2);
                    Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value3);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value4);
                    Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter4High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter4Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes[0];
                    result = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMmt_StartCTCommand(string TargetIP, uint CTEntryID)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_StartCTCommand
        return this.LMmt_StartCTCommand(TargetIP, CTEntryID);
    }

    //
    // Summary:
    //     Stop a activated Trigger Event or a running Command Table
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true
    public bool LMmt_ClearEventEvaluation(string TargetIP)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    int value = 0;
                    int value2 = 0;
                    int value3 = 0;
                    int value4 = 0;
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 0;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(0x80 | Axis[i].MotionCommandInterface.CountNibble);
                    byte[] bytes = BitConverter.GetBytes(value);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value2);
                    Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value3);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value4);
                    Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter4High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter4Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes[0];
                    result = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMmt_ClearEventEvaluation(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMmt_ClearEventEvaluation
        return this.LMmt_ClearEventEvaluation(TargetIP);
    }

    //
    // Summary:
    //     Run a stored curve on the drive
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   CurveID:
    //     ID number of the stored curve
    //
    //   CurveOffset:
    //     Offset of curve start point in 0.1um
    //
    //   TimeScale:
    //     Scale factor between 0...200% of curve time base
    //
    //   AmplitudeScale:
    //     Scale factor for curve amplitude between -2000 ... 2000%
    //
    // Returns:
    //     Returns true
    public bool LMav_RunCurve(string TargetIP, int CurveID, int CurveOffset, int TimeScale, int AmplitudeScale)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 4;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(0x40 | Axis[i].MotionCommandInterface.CountNibble);
                    byte[] bytes = BitConverter.GetBytes(CurveID);
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(CurveOffset);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[0];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[2];
                    bytes = BitConverter.GetBytes(TimeScale);
                    Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter2High = bytes[0];
                    bytes = BitConverter.GetBytes(AmplitudeScale);
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = 0;
                    Axis[i].MotionCommandInterface.CommandParameter3High = 0;
                    result = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMav_RunCurve(string TargetIP, int CurveID, int CurveOffset, int TimeScale, int AmplitudeScale)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMav_RunCurve
        return this.LMav_RunCurve(TargetIP, CurveID, CurveOffset, TimeScale, AmplitudeScale);
    }

    //
    // Summary:
    //     Perform a jerk limited move
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   Position:
    //     Target position [mm]
    //
    //   Velocity:
    //     Maximum velocity [m/s]
    //
    //   Acceleration:
    //     Acceleration [m/s2]
    //
    //   Jerk:
    //     Jerk [m/s3]
    //
    // Returns:
    //     Returns true
    public bool LMav_MoveBestehorn(string TargetIP, float Position, float Velocity, float Acceleration, float Jerk)
    {
        checked
        {
            int value = (int)Math.Round(Position * 10000f);
            uint value2 = (uint)Math.Round(Velocity * 1000000f);
            uint value3 = (uint)Math.Round(Acceleration * 10000f);
            uint value4 = (uint)Math.Round(Jerk * 1000f);
            bool result = false;
            int targetIPListCount = TargetIPListCount;
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 15;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(0 | Axis[i].MotionCommandInterface.CountNibble);
                    byte[] bytes = BitConverter.GetBytes(value);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value2);
                    Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value3);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value4);
                    Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter4High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter4Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes[0];
                    result = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMav_MoveBestehorn(string TargetIP, float Position, float Velocity, float Acceleration, float Jerk)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMav_MoveBestehorn
        return this.LMav_MoveBestehorn(TargetIP, Position, Velocity, Acceleration, Jerk);
    }

    //
    // Summary:
    //     Perform a jerk limited relative move
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   Position:
    //     Target position [mm]
    //
    //   Velocity:
    //     Maximum velocity [m/s]
    //
    //   Acceleration:
    //     Acceleration [m/s2]
    //
    //   Jerk:
    //     Jerk [m/s3]
    //
    // Returns:
    //     Returns true
    public bool LMav_MoveBestehornRelative(string TargetIP, float Position, float Velocity, float Acceleration, float Jerk)
    {
        checked
        {
            int value = (int)Math.Round(Position * 10000f);
            uint value2 = (uint)Math.Round(Velocity * 1000000f);
            uint value3 = (uint)Math.Round(Acceleration * 10000f);
            uint value4 = (uint)Math.Round(Jerk * 1000f);
            bool result = false;
            int targetIPListCount = TargetIPListCount;
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 15;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(0x10 | Axis[i].MotionCommandInterface.CountNibble);
                    byte[] bytes = BitConverter.GetBytes(value);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value2);
                    Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value3);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value4);
                    Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter4High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter4Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes[0];
                    result = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMav_MoveBestehornRelative(string TargetIP, float Position, float Velocity, float Acceleration, float Jerk)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMav_MoveBestehornRelative
        return this.LMav_MoveBestehornRelative(TargetIP, Position, Velocity, Acceleration, Jerk);
    }

    //
    // Summary:
    //     Perform a sinusodial move
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   Position:
    //     Target position [mm]
    //
    //   Velocity:
    //     Maximum velocity [m/s]
    //
    //   Acceleration:
    //     Acceleration [m/s2]
    //
    // Returns:
    //     Returns true
    public bool LMav_MoveSin(string TargetIP, float Position, float Velocity, float Acceleration)
    {
        checked
        {
            int value = (int)Math.Round(Position * 10000f);
            uint value2 = (uint)Math.Round(Velocity * 1000000f);
            uint value3 = (uint)Math.Round(Acceleration * 100000f);
            bool result = false;
            int targetIPListCount = TargetIPListCount;
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 14;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(0 | Axis[i].MotionCommandInterface.CountNibble);
                    byte[] bytes = BitConverter.GetBytes(value);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value2);
                    Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value3);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(0);
                    Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter4High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter4Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes[0];
                    result = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMav_MoveSin(string TargetIP, float Position, float Velocity, float Acceleration)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMav_MoveSin
        return this.LMav_MoveSin(TargetIP, Position, Velocity, Acceleration);
    }

    //
    // Summary:
    //     Perform a sinusodial relative move
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   Position:
    //     Target position [mm]
    //
    //   Velocity:
    //     Maximum velocity [m/s]
    //
    //   Acceleration:
    //     Acceleration [m/s2]
    //
    // Returns:
    //     Returns true
    public bool LMav_MoveSinRelative(string TargetIP, float Position, float Velocity, float Acceleration)
    {
        checked
        {
            int value = (int)Math.Round(Position * 10000f);
            uint value2 = (uint)Math.Round(Velocity * 1000000f);
            uint value3 = (uint)Math.Round(Acceleration * 100000f);
            bool result = false;
            int targetIPListCount = TargetIPListCount;
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 14;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(0x10 | Axis[i].MotionCommandInterface.CountNibble);
                    byte[] bytes = BitConverter.GetBytes(value);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value2);
                    Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value3);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(0);
                    Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter4High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter4Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes[0];
                    result = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMav_MoveSinRelative(string TargetIP, float Position, float Velocity, float Acceleration)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMav_MoveSinRelative
        return this.LMav_MoveSinRelative(TargetIP, Position, Velocity, Acceleration);
    }

    //
    // Summary:
    //     Change the target force while force control is active
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   TargetForce:
    //     New target force value as setpoint
    //
    // Returns:
    //     Returns true
    public bool LMfc_ChangeTargetForce(string TargetIP, float TargetForce)
    {
        checked
        {
            int value = (int)Math.Round(TargetForce * 10f);
            bool result = false;
            int targetIPListCount = TargetIPListCount;
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 56;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(0x20 | Axis[i].MotionCommandInterface.CountNibble);
                    byte[] bytes = BitConverter.GetBytes(value);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[0];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[0];
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(0);
                    Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(0);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(0);
                    Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter4High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter4Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes[0];
                    result = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMfc_ChangeTargetForce(string TargetIP, float TargetForce)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMfc_ChangeTargetForce
        return this.LMfc_ChangeTargetForce(TargetIP, TargetForce);
    }

    //
    // Summary:
    //     Start force control mode by positive force values (push), if force limit is detected.
    //     If drive has changed to force mode, target force will be used as setpoint.
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   Position:
    //     Target position [mm]
    //
    //   Velocity:
    //     Maximum velocity [m/s]
    //
    //   Acceleration:
    //     Acceleration [m/s2]
    //
    //   ForceLimit:
    //     Force value, which must be reached, before drive change from position to force
    //     control
    //
    //   TargetForce:
    //     New target force value as setpoint
    //
    // Returns:
    //     Returns true
    public bool LMfc_GoToPosForceCtrlHighLim(string TargetIP, float Position, float Velocity, float Acceleration, float ForceLimit, float TargetForce)
    {
        bool result = false;
        checked
        {
            int value = (int)Math.Round(ForceLimit * 10f);
            int value2 = (int)Math.Round(TargetForce * 10f);
            int value3 = (int)Math.Round(Position * 10000f);
            uint value4 = (uint)Math.Round(Velocity * 1000000f);
            uint value5 = (uint)Math.Round(Acceleration * 100000f);
            int targetIPListCount = TargetIPListCount;
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 56;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(0x30 | Axis[i].MotionCommandInterface.CountNibble);
                    byte[] bytes = BitConverter.GetBytes(value3);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value4);
                    Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value5);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value);
                    Axis[i].MotionCommandInterface.CommandParameter4Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value2);
                    Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4High = bytes[0];
                    bytes = BitConverter.GetBytes(0);
                    Axis[i].MotionCommandInterface.CommandParameter5HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter5High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter5Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter5LowLow = bytes[0];
                    result = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMfc_GoToPosForceCtrlHighLim(string TargetIP, float Position, float Velocity, float Acceleration, float ForceLimit, float TargetForce)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMfc_GoToPosForceCtrlHighLim
        return this.LMfc_GoToPosForceCtrlHighLim(TargetIP, Position, Velocity, Acceleration, ForceLimit, TargetForce);
    }

    //
    // Summary:
    //     Start force control mode by negative force values (pull), if force limit is detected.
    //     If drive has changed to force mode, target force will be used as setpoint.
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   Position:
    //     Target position [mm]
    //
    //   Velocity:
    //     Maximum velocity [m/s]
    //
    //   Acceleration:
    //     Acceleration [m/s2]
    //
    //   ForceLimit:
    //     Force value, which must be reached, before drive change from position to force
    //     control
    //
    //   TargetForce:
    //     New target force value as setpoint
    //
    // Returns:
    //     Returns true
    public bool LMfc_GoToPosForceCtrlLowLim(string TargetIP, float Position, float Velocity, float Acceleration, float ForceLimit, float TargetForce)
    {
        bool result = false;
        checked
        {
            int value = (int)Math.Round(Position * 10000f);
            uint value2 = (uint)Math.Round(Velocity * 1000000f);
            uint value3 = (uint)Math.Round(Acceleration * 100000f);
            int value4 = (int)Math.Round(ForceLimit * 10f);
            int value5 = (int)Math.Round(TargetForce * 10f);
            int targetIPListCount = TargetIPListCount;
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 56;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(0x50 | Axis[i].MotionCommandInterface.CountNibble);
                    byte[] bytes = BitConverter.GetBytes(value);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value2);
                    Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value3);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value4);
                    Axis[i].MotionCommandInterface.CommandParameter4Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value5);
                    Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4High = bytes[0];
                    bytes = BitConverter.GetBytes(0);
                    Axis[i].MotionCommandInterface.CommandParameter5HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter5High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter5Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter5LowLow = bytes[0];
                    result = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMfc_GoToPosForceCtrlLowLim(string TargetIP, float Position, float Velocity, float Acceleration, float ForceLimit, float TargetForce)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMfc_GoToPosForceCtrlLowLim
        return this.LMfc_GoToPosForceCtrlLowLim(TargetIP, Position, Velocity, Acceleration, ForceLimit, TargetForce);
    }

    //
    // Summary:
    //     Reset force control mode, drive change to position control mode and move to the
    //     specified position
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   Position:
    //     Target position [mm]
    //
    //   Velocity:
    //     Maximum velocity [m/s]
    //
    //   Acceleration:
    //     Acceleration [m/s2]
    //
    //   Deceleration:
    //     Deceleration [m/s2]
    //
    // Returns:
    //     Returns true
    public bool LMfc_GoToPosRstForceCtrl(string TargetIP, float Position, float Velocity, float Acceleration, float Deceleration)
    {
        bool result = false;
        checked
        {
            int value = (int)Math.Round(Position * 10000f);
            uint value2 = (uint)Math.Round(Velocity * 1000000f);
            uint value3 = (uint)Math.Round(Acceleration * 100000f);
            uint value4 = (uint)Math.Round(Deceleration * 100000f);
            int targetIPListCount = TargetIPListCount;
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 56;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(0x10 | Axis[i].MotionCommandInterface.CountNibble);
                    byte[] bytes = BitConverter.GetBytes(value);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value2);
                    Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value3);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value4);
                    Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter4High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter4Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes[0];
                    result = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMfc_GoToPosRstForceCtrl(string TargetIP, float Position, float Velocity, float Acceleration, float Deceleration)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMfc_GoToPosRstForceCtrl
        return this.LMfc_GoToPosRstForceCtrl(TargetIP, Position, Velocity, Acceleration, Deceleration);
    }

    //
    // Summary:
    //     (386xh) Reset force control mode, drive change to position control mode and move
    //     to the specified position. The I part of the position controller is set to the
    //     last force control current
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   Position:
    //     Target position [mm]
    //
    //   Velocity:
    //     Maximum velocity [m/s]
    //
    //   Acceleration:
    //     Acceleration [m/s2]
    //
    //   Deceleration:
    //     Deceleration [m/s2]
    //
    // Returns:
    //     Returns true
    public bool LMfc_GoToPosRstForceCtrlSetI(string TargetIP, float Position, float Velocity, float Acceleration, float Deceleration)
    {
        bool result = false;
        checked
        {
            int value = (int)Math.Round(Position * 10000f);
            uint value2 = (uint)Math.Round(Velocity * 1000000f);
            uint value3 = (uint)Math.Round(Acceleration * 100000f);
            uint value4 = (uint)Math.Round(Deceleration * 100000f);
            int targetIPListCount = TargetIPListCount;
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 56;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(0x60 | Axis[i].MotionCommandInterface.CountNibble);
                    byte[] bytes = BitConverter.GetBytes(value);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value2);
                    Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value3);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value4);
                    Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter4High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter4Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes[0];
                    result = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMfc_GoToPosRstForceCtrlSetI(string TargetIP, float Position, float Velocity, float Acceleration, float Deceleration)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMfc_GoToPosRstForceCtrlSetI
        return this.LMfc_GoToPosRstForceCtrlSetI(TargetIP, Position, Velocity, Acceleration, Deceleration);
    }

    //
    // Summary:
    //     (387xh) Reinstalls the position control mode and increment actual position. The
    //     I part of the position controller is set to the last force control current.
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   Position:
    //     Target position [mm]
    //
    //   Velocity:
    //     Maximum velocity [m/s]
    //
    //   Acceleration:
    //     Acceleration [m/s2]
    //
    //   Deceleration:
    //     Deceleration [m/s2]
    //
    // Returns:
    //     Returns true
    public bool LMfc_IncrementActPosAndResetForceControlSetI(string TargetIP, float Position, float Velocity, float Acceleration, float Deceleration)
    {
        bool result = false;
        checked
        {
            int value = (int)Math.Round(Position * 10000f);
            uint value2 = (uint)Math.Round(Velocity * 1000000f);
            uint value3 = (uint)Math.Round(Acceleration * 100000f);
            uint value4 = (uint)Math.Round(Deceleration * 100000f);
            int targetIPListCount = TargetIPListCount;
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 56;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(0x70 | Axis[i].MotionCommandInterface.CountNibble);
                    byte[] bytes = BitConverter.GetBytes(value);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value2);
                    Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value3);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value4);
                    Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter4High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter4Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes[0];
                    result = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMfc_IncrementActPosAndResetForceControlSetI(string TargetIP, float Position, float Velocity, float Acceleration, float Deceleration)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMfc_IncrementActPosAndResetForceControlSetI
        return this.LMfc_IncrementActPosAndResetForceControlSetI(TargetIP, Position, Velocity, Acceleration, Deceleration);
    }

    //
    // Summary:
    //     (388xh) VAI increment actual position, if the measured force reaches the defined
    //     value, the drive switches to the force control mode with Target Force = Target
    //     Force.
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   PositionIncrement:
    //     Position Increment [mm]
    //
    //   MaxVelocity:
    //     Maximum velocity [m/s]
    //
    //   Acceleration:
    //     Acceleration [m/s2]
    //
    //   ForceLimit:
    //     Force Limit [N]
    //
    //   TargetForce:
    //     Target Force [N]
    //
    // Returns:
    //     Returns true
    public bool LMfc_IncrementActPosWithHigherForceCtrlLimitAndTargetForce(string TargetIP, float PositionIncrement, float MaxVelocity, float Acceleration, float ForceLimit, float TargetForce)
    {
        bool result = false;
        checked
        {
            int value = (int)Math.Round(PositionIncrement * 10000f);
            uint value2 = (uint)Math.Round(MaxVelocity * 1000000f);
            uint value3 = (uint)Math.Round(Acceleration * 100000f);
            int value4 = (int)Math.Round(ForceLimit * 10f);
            int value5 = (int)Math.Round(TargetForce * 10f);
            int targetIPListCount = TargetIPListCount;
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 56;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(0x80 | Axis[i].MotionCommandInterface.CountNibble);
                    byte[] bytes = BitConverter.GetBytes(value);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value2);
                    Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value3);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value4);
                    Axis[i].MotionCommandInterface.CommandParameter4Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(value5);
                    Axis[i].MotionCommandInterface.CommandParameter4HighHigh = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter4High = bytes[0];
                    bytes = BitConverter.GetBytes(0);
                    Axis[i].MotionCommandInterface.CommandParameter5HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter5High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter5Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter5LowLow = bytes[0];
                    result = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMfc_IncrementActPosWithHigherForceCtrlLimitAndTargetForce(string TargetIP, float PositionIncrement, float MaxVelocity, float Acceleration, float ForceLimit, float TargetForce)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMfc_IncrementActPosWithHigherForceCtrlLimitAndTargetForce
        return this.LMfc_IncrementActPosWithHigherForceCtrlLimitAndTargetForce(TargetIP, PositionIncrement, MaxVelocity, Acceleration, ForceLimit, TargetForce);
    }

    //
    // Summary:
    //     Enter current command mode. No position control is working, the send current
    //     value will push or pull the slider with this current. Can be used fo open loop
    //     force control.
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   Current:
    //     Current value [mA], which is set to push or pull the slider
    //
    // Returns:
    //     Returns true
    public bool LMav_SetCurrentCommandMode(string TargetIP, int Current)
    {
        bool flag = false;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 57;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(0 | Axis[i].MotionCommandInterface.CountNibble);
                    byte[] bytes = BitConverter.GetBytes(Current);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(0);
                    Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(0);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(0);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    flag = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return true;
        }
    }

    bool _ACI.LMav_SetCurrentCommandMode(string TargetIP, int Current)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMav_SetCurrentCommandMode
        return this.LMav_SetCurrentCommandMode(TargetIP, Current);
    }

    //
    // Summary:
    //     Leave current command mode and go back to position controlled mode
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true
    public bool LMav_ResetCurrentCommandMode(string TargetIP)
    {
        bool flag = false;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].MotionCommandInterface.CountNibble = (byte)((Axis[i].StateVarLow & 0xF) + 1);
                    if (Axis[i].MotionCommandInterface.CountNibble > 4)
                    {
                        Axis[i].MotionCommandInterface.CountNibble = 1;
                    }

                    Axis[i].MotionCommandInterface.CommandHeaderMasterID = 57;
                    Axis[i].MotionCommandInterface.CommandHeaderSubID = (byte)(0xF0 | Axis[i].MotionCommandInterface.CountNibble);
                    byte[] bytes = BitConverter.GetBytes(0);
                    Axis[i].MotionCommandInterface.CommandParameter1HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter1High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter1Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter1LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(0);
                    Axis[i].MotionCommandInterface.CommandParameter2HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter2High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter2Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter2LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(0);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    bytes = BitConverter.GetBytes(0);
                    Axis[i].MotionCommandInterface.CommandParameter3HighHigh = bytes[3];
                    Axis[i].MotionCommandInterface.CommandParameter3High = bytes[2];
                    Axis[i].MotionCommandInterface.CommandParameter3Low = bytes[1];
                    Axis[i].MotionCommandInterface.CommandParameter3LowLow = bytes[0];
                    flag = true;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }
            }

            return true;
        }
    }

    bool _ACI.LMav_ResetCurrentCommandMode(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMav_ResetCurrentCommandMode
        return this.LMav_ResetCurrentCommandMode(TargetIP);
    }

    //
    // Summary:
    //     Curve function: get curve download/save progress, value is in percent of data
    //     points
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true
    public int LMcf_getCurveProgress(string TargetIP)
    {
        int result = 0;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
            {
                continue;
            }

            object taskLock = TaskLock;
            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
            bool lockTaken = false;
            try
            {
                Monitor.Enter(taskLock, ref lockTaken);
                if (Axis[i].RealTimeConfigCurveStatus > 100)
                {
                    Axis[i].RealTimeConfigCurveStatus = 100;
                }

                result = Axis[i].RealTimeConfigCurveStatus;
                if (!CurveThreadStart)
                {
                    result = 0;
                }
            }
            finally
            {
                if (lockTaken)
                {
                    Monitor.Exit(taskLock);
                }
            }
        }

        return result;
    }

    int _ACI.LMcf_getCurveProgress(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMcf_getCurveProgress
        return this.LMcf_getCurveProgress(TargetIP);
    }

    //
    // Summary:
    //     Curve function: Delete, load or store a curve to drive.
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   CurveID:
    //     ID number of the curve
    //
    //   SetpointCount:
    //     length of array (setpoint count)
    //
    //   CurveName:
    //     Optional curve name, which can be defined
    //
    //   Xcode:
    //     See DLL manual for details (0: Time, 1: Encoder position, 2: Position
    //
    //   Ycode:
    //     See DLL manual for details (0: Position, 1: Velocity, 2: Current, 3: Encoder
    //     Position (Inc), 4: Encoder Speed, 5: MicroSteps
    //
    //   XLength:
    //     Defines length of x axis of curve, one step is 10us or one increment
    //
    //   XDim:
    //     Defines axis scale - 5: 0.1um, 26: 0.01ms, 27: one increment
    //
    //   YDim:
    //     Defines axis scale - 5: 0.1um, 26: 0.01ms, 27: one increment
    //
    //   Mode:
    //     Operation mode - 0: Load curve to RAM, 1: delete all stored curves (RAM and ROM),
    //     2: Store all curves from RAM to ROM
    //
    //   Setpoints:
    //     Array of the defined stepoints to be downloaded
    //
    // Returns:
    //     Returns true, if command has started
    public bool LMcf_LoadCurve(string TargetIP, int CurveID, int SetpointCount, string CurveName, byte Xcode, byte Ycode, uint XLength, int XDim, int YDim, int Mode, int[] Setpoints)
    {
        if (!CurveThreadStart)
        {
            CurveT = new Thread(CurveThread);
            CurveCallData.TargetIP = TargetIP;
            CurveCallData.CurveID = CurveID;
            CurveCallData.SetpointCount = SetpointCount;
            CurveCallData.CurveName = CurveName;
            CurveCallData.XLength = XLength;
            CurveCallData.XDim = XDim;
            CurveCallData.YDim = YDim;
            ref int[] setpoints = ref CurveCallData.Setpoints;
            setpoints = (int[])Utils.CopyArray(setpoints, new int[checked(SetpointCount + 2 + 1)]);
            CurveCallData.Setpoints = Setpoints;
            CurveCallData.Mode = Mode;
            CurveCallData.Xcode = Xcode;
            CurveCallData.Ycode = Ycode;
            CurveT.IsBackground = true;
            CurveT.Start();
            CurveT.Name = "CurveServiceRunning";
            CurveThreadStart = true;
        }

        return CurveThreadStart;
    }

    bool _ACI.LMcf_LoadCurve(string TargetIP, int CurveID, int SetpointCount, string CurveName, byte Xcode, byte Ycode, uint XLength, int XDim, int YDim, int Mode, int[] Setpoints)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMcf_LoadCurve
        return this.LMcf_LoadCurve(TargetIP, CurveID, SetpointCount, CurveName, Xcode, Ycode, XLength, XDim, YDim, Mode, Setpoints);
    }

    //
    // Summary:
    //     Curve function: Delete, load or store a curve to drive.
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   Mode:
    //     Operation mode - 0: Load curve to RAM, 1: delete all stored curves (RAM and ROM),
    //     2: Store all curves from RAM to ROM
    //
    //   CurveData:
    //     Curve data package
    public bool LMcf_LoadCurve(string TargetIP, int Mode, CurveDataDefinition CurveData)
    {
        if (!CurveThreadStart)
        {
            Thread thread = new Thread(CurveThread);
            CurveCallData.TargetIP = TargetIP;
            CurveCallData.CurveID = CurveData.CurveID;
            CurveCallData.SetpointCount = CurveData.SetpointCount;
            CurveCallData.CurveName = CurveData.CurveName;
            CurveCallData.XLength = CurveData.XLength;
            CurveCallData.XDim = CurveData.XDim;
            CurveCallData.YDim = CurveData.YDim;
            ref int[] setpoints = ref CurveCallData.Setpoints;
            setpoints = (int[])Utils.CopyArray(setpoints, new int[checked(CurveData.SetpointCount + 2 + 1)]);
            CurveCallData.Setpoints = CurveData.Setpoints;
            CurveCallData.Mode = Mode;
            CurveCallData.Xcode = CurveData.Xcode;
            CurveCallData.Ycode = CurveData.Ycode;
            thread.IsBackground = true;
            thread.Start();
            thread.Name = "CurveServiceRunning";
            CurveThreadStart = true;
        }

        return CurveThreadStart;
    }

    bool _ACI.LMcf_LoadCurve(string TargetIP, int Mode, CurveDataDefinition CurveData)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMcf_LoadCurve
        return this.LMcf_LoadCurve(TargetIP, Mode, CurveData);
    }

    //
    // Summary:
    //     Curve function: indicates, curve download/save is in progress
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Returns true, if curve loading process still active
    public bool LMcf_isCurveLoading(string TargetIP)
    {
        bool flag = CurveThreadStart;
        if (flag)
        {
            if (CurveT == null)
            {
                flag = false;
            }

            if (!CurveT.IsAlive)
            {
                flag = false;
            }
        }

        return flag;
    }

    bool _ACI.LMcf_isCurveLoading(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMcf_isCurveLoading
        return this.LMcf_isCurveLoading(TargetIP);
    }

    private void CurveThread()
    {
        int num = CurveCallData.SetpointCount;
        int curveID = CurveCallData.CurveID;
        uint xLength = CurveCallData.XLength;
        int xDim = CurveCallData.XDim;
        int yDim = CurveCallData.YDim;
        string targetIP = CurveCallData.TargetIP;
        string text = CurveCallData.CurveName;
        int[] setpoints = CurveCallData.Setpoints;
        byte ycode = CurveCallData.Ycode;
        byte xcode = CurveCallData.Xcode;
        byte[] bytes = BitConverter.GetBytes(num);
        byte[] bytes2 = BitConverter.GetBytes(curveID);
        byte[] bytes3 = BitConverter.GetBytes(xLength);
        byte[] bytes4 = BitConverter.GetBytes(xDim);
        byte[] bytes5 = BitConverter.GetBytes(yDim);
        int num2 = 0;
        checked
        {
            if (text.Length < 22)
            {
                int num3 = 22 - text.Length;
                for (int i = 1; i <= num3; i++)
                {
                    text += " ";
                }
            }

            ycode = unchecked((byte)(ycode << 4));
            ycode &= 0xF0;
            byte b = unchecked((byte)(ycode | xcode));
            byte[] array = new byte[73]
            {
                b,
                3,
                0,
                70,
                0,
                4,
                bytes[1],
                bytes[0],
                (byte)unchecked((int)text[3]),
                (byte)unchecked((int)text[2]),
                (byte)unchecked((int)text[1]),
                (byte)unchecked((int)text[0]),
                (byte)unchecked((int)text[7]),
                (byte)unchecked((int)text[6]),
                (byte)unchecked((int)text[5]),
                (byte)unchecked((int)text[4]),
                (byte)unchecked((int)text[11]),
                (byte)unchecked((int)text[10]),
                (byte)unchecked((int)text[9]),
                (byte)unchecked((int)text[8]),
                (byte)unchecked((int)text[15]),
                (byte)unchecked((int)text[14]),
                (byte)unchecked((int)text[13]),
                (byte)unchecked((int)text[12]),
                (byte)unchecked((int)text[19]),
                (byte)unchecked((int)text[18]),
                (byte)unchecked((int)text[17]),
                (byte)unchecked((int)text[16]),
                bytes2[1],
                bytes2[0],
                (byte)unchecked((int)text[20]),
                (byte)unchecked((int)text[21]),
                bytes3[3],
                bytes3[2],
                bytes3[1],
                bytes3[0],
                bytes5[1],
                bytes5[0],
                bytes4[1],
                bytes4[0],
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0,
                0
            };
            int targetIPListCount = TargetIPListCount;
            for (int j = 0; j <= targetIPListCount; j++)
            {
                if (Operators.CompareString(Axis[j].SlaveIP, targetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                double interval = TxTimer.interval;
                TxTimer.interval = 15.0;
                int num4 = 18 + num + 2;
                int num5 = 0;
                if (CurveCallData.Mode == 1)
                {
                    object taskLock = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                    bool lockTaken = false;
                    try
                    {
                        Monitor.Enter(taskLock, ref lockTaken);
                        Axis[j].RealTimeConfigID = 65;
                        if (Axis[j].RealTimeConfigCommandCount > 14)
                        {
                            Axis[j].RealTimeConfigCommandCount = 1;
                        }

                        Axis[j].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[j].RealTimeConfigCommandCount + 1));
                        Axis[j].RealTimeConfigArgs1High = 0;
                        Axis[j].RealTimeConfigArgs1Low = 0;
                        Axis[j].RealTimeConfigArgs2High = 0;
                        Axis[j].RealTimeConfigArgs2Low = 0;
                        Axis[j].RealTimeConfigArgs3High = 0;
                        Axis[j].RealTimeConfigArgs3Low = 0;
                        Axis[j].isRespondActual = false;
                    }
                    finally
                    {
                        if (lockTaken)
                        {
                            Monitor.Exit(taskLock);
                        }
                    }

                    long timeOutTime = getTimeOutTime(20000L);
                    do
                    {
                        Thread.Sleep(1);
                        if (Axis[j].isRespondActual)
                        {
                            num2++;
                        }

                        if ((Axis[j].RealTimeConfigIDStatus == 2) | (Axis[j].RealTimeConfigIDStatus == 5))
                        {
                            num2 = 0;
                        }

                        if (!isTimeOut(timeOutTime))
                        {
                            continue;
                        }

                        object taskLock2 = TaskLock;
                        ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                        bool lockTaken2 = false;
                        try
                        {
                            Monitor.Enter(taskLock2, ref lockTaken2);
                            Axis[j].DLLErrorText = "Timeout while waiting for curve response during deleting curve";
                        }
                        finally
                        {
                            if (lockTaken2)
                            {
                                Monitor.Exit(taskLock2);
                            }
                        }

                        break;
                    }
                    while (Axis[j].RealTimeConfigStatusCommandCount != Axis[j].RealTimeConfigCommandCount);
                }

                if (CurveCallData.Mode == 0)
                {
                    num2 = 0;
                    num4 = 18 + num + 4;
                    num5 = 0;
                    object taskLock3 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                    bool lockTaken3 = false;
                    try
                    {
                        Monitor.Enter(taskLock3, ref lockTaken3);
                        Axis[j].RealTimeConfigID = 80;
                        if (Axis[j].RealTimeConfigCommandCount > 14)
                        {
                            Axis[j].RealTimeConfigCommandCount = 1;
                        }

                        Axis[j].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[j].RealTimeConfigCommandCount + 1));
                        Axis[j].RealTimeConfigArgs1High = bytes2[1];
                        Axis[j].RealTimeConfigArgs1Low = bytes2[0];
                        byte[] bytes6 = BitConverter.GetBytes(num * 4);
                        Axis[j].RealTimeConfigArgs3High = bytes6[1];
                        Axis[j].RealTimeConfigArgs3Low = bytes6[0];
                        Axis[j].RealTimeConfigArgs2High = 0;
                        Axis[j].RealTimeConfigArgs2Low = 70;
                        Axis[j].isRespondActual = false;
                    }
                    finally
                    {
                        if (lockTaken3)
                        {
                            Monitor.Exit(taskLock3);
                        }
                    }

                    long timeOutTime = getTimeOutTime(1500L);
                    do
                    {
                        Thread.Sleep(1);
                        if (Axis[j].isRespondActual)
                        {
                            num2++;
                        }

                        if ((Axis[j].RealTimeConfigIDStatus == 2) | (Axis[j].RealTimeConfigIDStatus == 5))
                        {
                            num2 = 0;
                        }

                        if (!isTimeOut(timeOutTime))
                        {
                            continue;
                        }

                        object taskLock4 = TaskLock;
                        ObjectFlowControl.CheckForSyncLockOnValueType(taskLock4);
                        bool lockTaken4 = false;
                        try
                        {
                            Monitor.Enter(taskLock4, ref lockTaken4);
                            Axis[j].DLLErrorText = "Timeout while waiting for curve response during start adding curve";
                        }
                        finally
                        {
                            if (lockTaken4)
                            {
                                Monitor.Exit(taskLock4);
                            }
                        }

                        break;
                    }
                    while (Axis[j].RealTimeConfigStatusCommandCount != Axis[j].RealTimeConfigCommandCount);
                    num2 = 0;
                    object taskLock5 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock5);
                    bool lockTaken5 = false;
                    try
                    {
                        Monitor.Enter(taskLock5, ref lockTaken5);
                        Axis[j].RealTimeConfigID = 81;
                    }
                    finally
                    {
                        if (lockTaken5)
                        {
                            Monitor.Exit(taskLock5);
                        }
                    }

                    int num6 = 0;
                    do
                    {
                        object taskLock6 = TaskLock;
                        ObjectFlowControl.CheckForSyncLockOnValueType(taskLock6);
                        bool lockTaken6 = false;
                        try
                        {
                            Monitor.Enter(taskLock6, ref lockTaken6);
                            if (Axis[j].RealTimeConfigCommandCount > 14)
                            {
                                Axis[j].RealTimeConfigCommandCount = 1;
                            }

                            Axis[j].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[j].RealTimeConfigCommandCount + 1));
                            Axis[j].RealTimeConfigArgs1High = bytes2[1];
                            Axis[j].RealTimeConfigArgs1Low = bytes2[0];
                            Axis[j].RealTimeConfigArgs3High = array[num6];
                            Axis[j].RealTimeConfigArgs3Low = array[num6 + 1];
                            Axis[j].RealTimeConfigArgs2High = array[num6 + 2];
                            Axis[j].RealTimeConfigArgs2Low = array[num6 + 3];
                            timeOutTime = getTimeOutTime(1500L);
                            Axis[j].RealTimeConfigIDStatus = 0;
                            num5++;
                            Axis[j].RealTimeConfigCurveStatus = (int)Math.Round(Math.Round((double)num5 / (double)num4 * 100.0));
                            if (Axis[j].RealTimeConfigCurveStatus > 100)
                            {
                                Axis[j].RealTimeConfigCurveStatus = 100;
                            }

                            Axis[j].isRespondActual = false;
                        }
                        finally
                        {
                            if (lockTaken6)
                            {
                                Monitor.Exit(taskLock6);
                            }
                        }

                        do
                        {
                            Thread.Sleep(1);
                            if (Axis[j].isRespondActual)
                            {
                                num2++;
                            }

                            if (!isTimeOut(timeOutTime))
                            {
                                continue;
                            }

                            object taskLock7 = TaskLock;
                            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock7);
                            bool lockTaken7 = false;
                            try
                            {
                                Monitor.Enter(taskLock7, ref lockTaken7);
                                Axis[j].DLLErrorText = "Timeout while waiting for curve response during adding curve info block";
                            }
                            finally
                            {
                                if (lockTaken7)
                                {
                                    Monitor.Exit(taskLock7);
                                }
                            }

                            break;
                        }
                        while (Axis[j].RealTimeConfigStatusCommandCount != Axis[j].RealTimeConfigCommandCount);
                        num2 = 0;
                        object taskLock8 = TaskLock;
                        ObjectFlowControl.CheckForSyncLockOnValueType(taskLock8);
                        bool lockTaken8 = false;
                        try
                        {
                            Monitor.Enter(taskLock8, ref lockTaken8);
                            Axis[j].RealTimeConfigCurveStatus = 0;
                        }
                        finally
                        {
                            if (lockTaken8)
                            {
                                Monitor.Exit(taskLock8);
                            }
                        }

                        num6 += 4;
                    }
                    while (num6 <= 68);
                    object taskLock9 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock9);
                    bool lockTaken9 = false;
                    try
                    {
                        Monitor.Enter(taskLock9, ref lockTaken9);
                        Axis[j].RealTimeConfigCurveStatus = 0;
                        num4 = 18 + num + 4;
                        num5 = 0;
                        Axis[j].RealTimeConfigID = 82;
                    }
                    finally
                    {
                        if (lockTaken9)
                        {
                            Monitor.Exit(taskLock9);
                        }
                    }

                    int num7 = num - 1;
                    for (int k = 0; k <= num7; k++)
                    {
                        object taskLock10 = TaskLock;
                        ObjectFlowControl.CheckForSyncLockOnValueType(taskLock10);
                        bool lockTaken10 = false;
                        try
                        {
                            Monitor.Enter(taskLock10, ref lockTaken10);
                            if (Axis[j].RealTimeConfigCommandCount > 14)
                            {
                                Axis[j].RealTimeConfigCommandCount = 1;
                            }

                            Axis[j].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[j].RealTimeConfigCommandCount + 1));
                            Axis[j].RealTimeConfigArgs1High = bytes2[1];
                            Axis[j].RealTimeConfigArgs1Low = bytes2[0];
                            byte[] bytes7 = BitConverter.GetBytes(setpoints[k]);
                            Axis[j].RealTimeConfigArgs3High = bytes7[3];
                            Axis[j].RealTimeConfigArgs3Low = bytes7[2];
                            Axis[j].RealTimeConfigArgs2High = bytes7[1];
                            Axis[j].RealTimeConfigArgs2Low = bytes7[0];
                            timeOutTime = getTimeOutTime(500L);
                            Axis[j].RealTimeConfigIDStatus = 0;
                            num5++;
                            Axis[j].RealTimeConfigCurveStatus = (int)Math.Round(Math.Round((double)num5 / (double)num4 * 100.0));
                            if (Axis[j].RealTimeConfigCurveStatus > 100)
                            {
                                Axis[j].RealTimeConfigCurveStatus = 100;
                            }

                            Axis[j].isRespondActual = false;
                        }
                        finally
                        {
                            if (lockTaken10)
                            {
                                Monitor.Exit(taskLock10);
                            }
                        }

                        do
                        {
                            Thread.Sleep(1);
                            if (Axis[j].isRespondActual)
                            {
                                num2++;
                            }

                            if (!isTimeOut(timeOutTime))
                            {
                                continue;
                            }

                            object taskLock11 = TaskLock;
                            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock11);
                            bool lockTaken11 = false;
                            try
                            {
                                Monitor.Enter(taskLock11, ref lockTaken11);
                                string text2 = "Status CMD: " + Axis[j].RealTimeConfigStatusCommandCount + " - Config CMD: " + Axis[j].RealTimeConfigCommandCount + " RespondActual: " + num2;
                                Axis[j].DLLErrorText = "Timeout while waiting for curve response during adding curve data blocks - " + text2;
                            }
                            finally
                            {
                                if (lockTaken11)
                                {
                                    Monitor.Exit(taskLock11);
                                }
                            }

                            break;
                        }
                        while (Axis[j].RealTimeConfigStatusCommandCount != Axis[j].RealTimeConfigCommandCount);
                        num2 = 0;
                    }

                    object taskLock12 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock12);
                    bool lockTaken12 = false;
                    try
                    {
                        Monitor.Enter(taskLock12, ref lockTaken12);
                        Axis[j].RealTimeConfigCurveStatus = 0;
                    }
                    finally
                    {
                        if (lockTaken12)
                        {
                            Monitor.Exit(taskLock12);
                        }
                    }
                }

                if (CurveCallData.Mode == 2)
                {
                    num4 = 18 + num + 4;
                    num5 = 0;
                    object taskLock13 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock13);
                    bool lockTaken13 = false;
                    try
                    {
                        Monitor.Enter(taskLock13, ref lockTaken13);
                        Axis[j].RealTimeConfigID = 64;
                        if (Axis[j].RealTimeConfigCommandCount > 14)
                        {
                            Axis[j].RealTimeConfigCommandCount = 1;
                        }

                        Axis[j].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[j].RealTimeConfigCommandCount + 1));
                        Axis[j].RealTimeConfigArgs1High = 0;
                        Axis[j].RealTimeConfigArgs1Low = 0;
                        Axis[j].RealTimeConfigArgs2High = 0;
                        Axis[j].RealTimeConfigArgs2Low = 0;
                        Axis[j].RealTimeConfigArgs3High = 0;
                        Axis[j].RealTimeConfigArgs3Low = 0;
                        Axis[j].isRespondActual = false;
                    }
                    finally
                    {
                        if (lockTaken13)
                        {
                            Monitor.Exit(taskLock13);
                        }
                    }

                    long timeOutTime = getTimeOutTime(20000L);
                    do
                    {
                        if (Axis[j].isRespondActual)
                        {
                            num2++;
                        }

                        if (!isTimeOut(timeOutTime))
                        {
                            continue;
                        }

                        object taskLock14 = TaskLock;
                        ObjectFlowControl.CheckForSyncLockOnValueType(taskLock14);
                        bool lockTaken14 = false;
                        try
                        {
                            Monitor.Enter(taskLock14, ref lockTaken14);
                            Axis[j].DLLErrorText = "Timeout while waiting for curve response during storing curve";
                        }
                        finally
                        {
                            if (lockTaken14)
                            {
                                Monitor.Exit(taskLock14);
                            }
                        }

                        break;
                    }
                    while (Axis[j].RealTimeConfigStatusCommandCount != Axis[j].RealTimeConfigCommandCount);
                    object taskLock15 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock15);
                    bool lockTaken15 = false;
                    try
                    {
                        Monitor.Enter(taskLock15, ref lockTaken15);
                        Axis[j].RealTimeConfigIDStatus = 0;
                        Axis[j].RealTimeConfigID = 0;
                        if (Axis[j].RealTimeConfigCommandCount > 14)
                        {
                            Axis[j].RealTimeConfigCommandCount = 1;
                        }

                        Axis[j].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[j].RealTimeConfigCommandCount + 1));
                        if (num == 0)
                        {
                            num = 100;
                        }
                    }
                    finally
                    {
                        if (lockTaken15)
                        {
                            Monitor.Exit(taskLock15);
                        }
                    }

                    num2 = 0;
                    timeOutTime = getTimeOutTime(20000L);
                    do
                    {
                        if (Axis[j].isRespondActual)
                        {
                            num2++;
                        }

                        if (!isTimeOut(timeOutTime))
                        {
                            continue;
                        }

                        object taskLock16 = TaskLock;
                        ObjectFlowControl.CheckForSyncLockOnValueType(taskLock16);
                        bool lockTaken16 = false;
                        try
                        {
                            Monitor.Enter(taskLock16, ref lockTaken16);
                            Axis[j].DLLErrorText = "Timeout while waiting for curve response during storing curve";
                        }
                        finally
                        {
                            if (lockTaken16)
                            {
                                Monitor.Exit(taskLock16);
                            }
                        }

                        break;
                    }
                    while (Axis[j].RealTimeConfigStatusCommandCount != Axis[j].RealTimeConfigCommandCount);
                    num2 = 0;
                    object taskLock17 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock17);
                    bool lockTaken17 = false;
                    try
                    {
                        Monitor.Enter(taskLock17, ref lockTaken17);
                        Axis[j].RealTimeConfigCurveStatus = 0;
                    }
                    finally
                    {
                        if (lockTaken17)
                        {
                            Monitor.Exit(taskLock17);
                        }
                    }
                }

                object taskLock18 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock18);
                bool lockTaken18 = false;
                try
                {
                    Monitor.Enter(taskLock18, ref lockTaken18);
                    Axis[j].RealTimeConfigID = 0;
                }
                finally
                {
                    if (lockTaken18)
                    {
                        Monitor.Exit(taskLock18);
                    }
                }

                TxTimer.interval = interval;
            }

            CurveThreadStart = false;
        }
    }

    //
    // Summary:
    //     This function check, if the curve ID is a valid curve on the drive
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   CurveID:
    //     ID number of the curve
    //
    // Returns:
    //     Returns true, if curve ID exist
    public bool LMcf_isCurveOnDrive(string TargetIP, ushort CurveID)
    {
        bool result = false;
        byte[] bytes = BitConverter.GetBytes(CurveID);
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].RealTimeConfigID = 96;
                    if (Axis[i].RealTimeConfigCommandCount > 14)
                    {
                        Axis[i].RealTimeConfigCommandCount = 1;
                    }

                    Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[i].RealTimeConfigCommandCount + 1));
                    Axis[i].RealTimeConfigArgs1High = bytes[1];
                    Axis[i].RealTimeConfigArgs1Low = bytes[0];
                    Axis[i].RealTimeConfigArgs2High = 0;
                    Axis[i].RealTimeConfigArgs2Low = 0;
                    Axis[i].RealTimeConfigArgs3High = 0;
                    Axis[i].RealTimeConfigArgs3Low = 0;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                long timeOutTime = getTimeOutTime(5000L);
                while (!((Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount) & ((Axis[i].RealTimeConfigIDStatus == 0) | (Axis[i].RealTimeConfigIDStatus == 212))))
                {
                    if (!isTimeOut(timeOutTime))
                    {
                        continue;
                    }

                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        Axis[i].DLLErrorText = "8: Config Channel Timeout by requesting Command Table Data.. Error code: " + Conversion.Hex(Axis[i].RealTimeConfigIDStatus).ToString();
                        ACIError = "Config Channel Error";
                        Axis[i].RealTimeConfigStatusArgs2Low = 9;
                        Axis[i].RealTimeConfigStatusArgs2High = 0;
                        Axis[i].RealTimeConfigStatusArgs3Low = 9;
                        Axis[i].RealTimeConfigStatusArgs3High = 0;
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    break;
                }

                object taskLock3 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                bool lockTaken3 = false;
                try
                {
                    Monitor.Enter(taskLock3, ref lockTaken3);
                    byte[] array = new byte[2]
                    {
                        Axis[i].RealTimeConfigStatusArgs1Low,
                        Axis[i].RealTimeConfigStatusArgs1High
                    };
                    byte[] array2 = new byte[2]
                    {
                        Axis[i].RealTimeConfigStatusArgs2Low,
                        Axis[i].RealTimeConfigStatusArgs2High
                    };
                    Axis[i].RealTimeConfigID = 0;
                    if (Axis[i].RealTimeConfigIDStatus == 0)
                    {
                        result = true;
                    }
                }
                finally
                {
                    if (lockTaken3)
                    {
                        Monitor.Exit(taskLock3);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMcf_isCurveOnDrive(string TargetIP, ushort CurveID)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMcf_isCurveOnDrive
        return this.LMcf_isCurveOnDrive(TargetIP, CurveID);
    }

    //
    // Summary:
    //     This function delivres a list of all curve ID's hich are present on the drive
    //
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Array of Integer containing the curve ID's
    public int[] LMcf_getAllCurveID(string TargetIP)
    {
        int[] array = new int[1];
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                int num = 1;
                do
                {
                    if (LMcf_isCurveOnDrive(TargetIP, (ushort)num))
                    {
                        array[array.Count() - 1] = num;
                        array = (int[])Utils.CopyArray(array, new int[array.Count() + 1]);
                    }

                    num++;
                }
                while (num <= 100);
                array = (int[])Utils.CopyArray(array, new int[array.Count() - 2 + 1]);
            }

            return array;
        }
    }

    int[] _ACI.LMcf_getAllCurveID(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMcf_getAllCurveID
        return this.LMcf_getAllCurveID(TargetIP);
    }

    //
    // Summary:
    //     this function uploads a cuve from the drive into a curve data package
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   CurveID:
    //     ID number of the curve
    //
    // Returns:
    //     Returns True, if command is executed sucessfully
    public bool LMcf_StartUploadCurve(string TargetIP, ushort CurveID)
    {
        if (!CurveThreadStart)
        {
            Thread thread = new Thread(UploadCurve);
            CurveCallData.TargetIP = TargetIP;
            CurveCallData.CurveID = CurveID;
            CurveCallData.SetpointCount = 0;
            CurveCallData.CurveName = "";
            CurveCallData.XLength = 0u;
            CurveCallData.XDim = 0;
            CurveCallData.YDim = 0;
            CurveCallData.Setpoints = new int[1];
            CurveCallData.Mode = 0;
            CurveCallData.Xcode = 0;
            CurveCallData.Ycode = 0;
            thread.IsBackground = true;
            thread.Start();
            thread.Name = "CurveServiceRunning";
            CurveThreadStart = true;
        }

        return CurveThreadStart;
    }

    bool _ACI.LMcf_StartUploadCurve(string TargetIP, ushort CurveID)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMcf_StartUploadCurve
        return this.LMcf_StartUploadCurve(TargetIP, CurveID);
    }

    //
    // Summary:
    //     This function starts downloading the curve data into the drive
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     True, if process is started
    public bool LMcf_StartDownloadCurve(string TargetIP)
    {
        if (!CurveThreadStart)
        {
            int targetIPListCount = TargetIPListCount;
            int num = 0;
            if (num <= targetIPListCount && Operators.CompareString(Axis[num].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                try
                {
                    if (!((Axis[num].UploadedCurve.CurveID == 0) | (Axis[num].UploadedCurve.SetpointCount == 0)))
                    {
                        Thread thread = new Thread(CurveThread);
                        CurveCallData.TargetIP = TargetIP;
                        CurveCallData.CurveID = Axis[num].UploadedCurve.CurveID;
                        CurveCallData.SetpointCount = Axis[num].UploadedCurve.SetpointCount;
                        CurveCallData.CurveName = Axis[num].UploadedCurve.CurveName;
                        CurveCallData.XLength = Axis[num].UploadedCurve.XLength;
                        CurveCallData.XDim = Axis[num].UploadedCurve.XDim;
                        CurveCallData.YDim = Axis[num].UploadedCurve.YDim;
                        CurveCallData.Setpoints = Axis[num].UploadedCurve.Setpoints;
                        CurveCallData.Mode = 0;
                        CurveCallData.Xcode = Axis[num].UploadedCurve.Xcode;
                        CurveCallData.Ycode = Axis[num].UploadedCurve.Ycode;
                        thread.IsBackground = true;
                        thread.Start();
                        thread.Name = "CurveServiceRunning";
                        CurveThreadStart = true;
                    }
                    else
                    {
                        object taskLock = TaskLock;
                        ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                        bool lockTaken = false;
                        try
                        {
                            Monitor.Enter(taskLock, ref lockTaken);
                            Axis[num].DLLErrorText = "Error: No Curve defined for Downloading on " + TargetIP;
                        }
                        finally
                        {
                            if (lockTaken)
                            {
                                Monitor.Exit(taskLock);
                            }
                        }
                    }
                }
                catch (Exception ex)
                {
                    ProjectData.SetProjectError(ex);
                    Exception ex2 = ex;
                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        object taskLock3 = TaskLock;
                        ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                        bool lockTaken3 = false;
                        try
                        {
                            Monitor.Enter(taskLock3, ref lockTaken3);
                            Axis[num].DLLErrorText = "Error: No Curve defined for Downloading on " + TargetIP;
                        }
                        finally
                        {
                            if (lockTaken3)
                            {
                                Monitor.Exit(taskLock3);
                            }
                        }
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    ProjectData.ClearProjectError();
                }
            }
        }

        return CurveThreadStart;
    }

    bool _ACI.LMcf_StartDownloadCurve(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMcf_StartDownloadCurve
        return this.LMcf_StartDownloadCurve(TargetIP);
    }

    //
    // Summary:
    //     This function set curve data into the DLL data structure for download to the
    //     drive.
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   CurveData:
    //     Curve data for download to the drive
    //
    // Returns:
    //     True, if data is sucessfully set for download to drive
    public bool LMcf_setDownloadCurveData(string TargetIP, CurveDataDefinition CurveData)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                Axis[i].UploadedCurve = CurveData;
                result = true;
                break;
            }
        }

        return result;
    }

    bool _ACI.LMcf_setDownloadCurveData(string TargetIP, CurveDataDefinition CurveData)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMcf_setDownloadCurveData
        return this.LMcf_setDownloadCurveData(TargetIP, CurveData);
    }

    //
    // Summary:
    //     This function get an uploaded curve data package
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Curve data package
    public CurveDataDefinition LMcf_getUploadedCurveData(string TargetIP)
    {
        CurveDataDefinition result = default(CurveDataDefinition);
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                return Axis[i].UploadedCurve;
            }
        }

        return result;
    }

    CurveDataDefinition _ACI.LMcf_getUploadedCurveData(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMcf_getUploadedCurveData
        return this.LMcf_getUploadedCurveData(TargetIP);
    }

    //
    // Summary:
    //     This function deletes all curves in RAM
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     True, if sucessfully deleted
    public bool LMcf_DeleteAllCurvesInRAM(string TargetIP)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].RealTimeConfigID = 65;
                    if (Axis[i].RealTimeConfigCommandCount > 14)
                    {
                        Axis[i].RealTimeConfigCommandCount = 1;
                    }

                    Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[i].RealTimeConfigCommandCount + 1));
                    Axis[i].RealTimeConfigCommandCount = Axis[i].RealTimeConfigCommandCount;
                    Axis[i].RealTimeConfigArgs1High = 0;
                    Axis[i].RealTimeConfigArgs1Low = 0;
                    Axis[i].RealTimeConfigArgs2High = 0;
                    Axis[i].RealTimeConfigArgs2Low = 0;
                    Axis[i].RealTimeConfigArgs3High = 0;
                    Axis[i].RealTimeConfigArgs3Low = 0;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                long timeOutTime = getTimeOutTime(1000L);
                while (!((Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount) & (Axis[i].RealTimeConfigIDStatus == 0)))
                {
                    if (!isTimeOut(timeOutTime))
                    {
                        continue;
                    }

                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        Axis[i].DLLErrorText = "Timeout by deleting all curves in RAM. Error code: " + Conversion.Hex(Axis[i].RealTimeConfigIDStatus).ToString();
                        ACIError = "Config Channel Error";
                        Axis[i].RealTimeConfigStatusArgs2Low = 9;
                        Axis[i].RealTimeConfigStatusArgs2High = 0;
                        Axis[i].RealTimeConfigStatusArgs3Low = 9;
                        Axis[i].RealTimeConfigStatusArgs3High = 0;
                        result = false;
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    break;
                }

                object taskLock3 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                bool lockTaken3 = false;
                try
                {
                    Monitor.Enter(taskLock3, ref lockTaken3);
                    byte[] array = new byte[8]
                    {
                        Axis[i].RealTimeConfigStatusArgs2Low,
                        Axis[i].RealTimeConfigStatusArgs2High,
                        Axis[i].RealTimeConfigStatusArgs3Low,
                        Axis[i].RealTimeConfigStatusArgs3High,
                        0,
                        0,
                        0,
                        0
                    };
                    Axis[i].RealTimeConfigID = 0;
                    result = true;
                }
                finally
                {
                    if (lockTaken3)
                    {
                        Monitor.Exit(taskLock3);
                    }
                }
            }

            return result;
        }
    }

    bool _ACI.LMcf_DeleteAllCurvesInRAM(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMcf_DeleteAllCurvesInRAM
        return this.LMcf_DeleteAllCurvesInRAM(TargetIP);
    }

    //
    // Summary:
    //     This function copies the actual curve RAM into drives FLASH memory
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     True, if sucessful
    private bool LMcf_SaveAllCurvesFromRAMToFLASH(string TargetIP)
    {
        bool result = false;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].RealTimeConfigID = 64;
                    if (Axis[i].RealTimeConfigCommandCount > 14)
                    {
                        Axis[i].RealTimeConfigCommandCount = 1;
                    }

                    Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[i].RealTimeConfigCommandCount + 1));
                    Axis[i].RealTimeConfigCommandCount = Axis[i].RealTimeConfigCommandCount;
                    Axis[i].RealTimeConfigArgs1High = 0;
                    Axis[i].RealTimeConfigArgs1Low = 0;
                    Axis[i].RealTimeConfigArgs2High = 0;
                    Axis[i].RealTimeConfigArgs2Low = 0;
                    Axis[i].RealTimeConfigArgs3High = 0;
                    Axis[i].RealTimeConfigArgs3Low = 0;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                long timeOutTime = getTimeOutTime(50000L);
                while (!((Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount) & (Axis[i].RealTimeConfigIDStatus == 0)))
                {
                    if (!isTimeOut(timeOutTime))
                    {
                        continue;
                    }

                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        Axis[i].DLLErrorText = "Timeout by transferring all curves from RAM to FLASH. Error code: " + Conversion.Hex(Axis[i].RealTimeConfigIDStatus).ToString();
                        ACIError = "Config Channel Error";
                        Axis[i].RealTimeConfigStatusArgs2Low = 9;
                        Axis[i].RealTimeConfigStatusArgs2High = 0;
                        Axis[i].RealTimeConfigStatusArgs3Low = 9;
                        Axis[i].RealTimeConfigStatusArgs3High = 0;
                        result = false;
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    break;
                }

                object taskLock3 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                bool lockTaken3 = false;
                try
                {
                    Monitor.Enter(taskLock3, ref lockTaken3);
                    byte[] array = new byte[8]
                    {
                        Axis[i].RealTimeConfigStatusArgs2Low,
                        Axis[i].RealTimeConfigStatusArgs2High,
                        Axis[i].RealTimeConfigStatusArgs3Low,
                        Axis[i].RealTimeConfigStatusArgs3High,
                        0,
                        0,
                        0,
                        0
                    };
                    Axis[i].RealTimeConfigID = 0;
                    result = true;
                }
                finally
                {
                    if (lockTaken3)
                    {
                        Monitor.Exit(taskLock3);
                    }
                }
            }

            return result;
        }
    }

    private void UploadCurve()
    {
        CurveDataDefinition uploadedCurve = default(CurveDataDefinition);
        uploadedCurve.Setpoints = new int[1];
        int num = 1;
        int num2 = 0;
        string targetIP = CurveCallData.TargetIP;
        checked
        {
            ushort num3 = (ushort)CurveCallData.CurveID;
            byte[] bytes = BitConverter.GetBytes(num3);
            byte[] array = new byte[72];
            if (LMcf_isCurveOnDrive(targetIP, num3))
            {
                int targetIPListCount = TargetIPListCount;
                for (int i = 0; i <= targetIPListCount; i++)
                {
                    if (Operators.CompareString(Axis[i].SlaveIP, targetIP, TextCompare: false) != 0)
                    {
                        continue;
                    }

                    num2 = 0;
                    int num4 = 0;
                    int num5 = 0;
                    do
                    {
                        object taskLock = TaskLock;
                        ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                        bool lockTaken = false;
                        try
                        {
                            Monitor.Enter(taskLock, ref lockTaken);
                            Axis[i].RealTimeConfigID = 97;
                            if (Axis[i].RealTimeConfigCommandCount > 14)
                            {
                                Axis[i].RealTimeConfigCommandCount = 1;
                            }

                            Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[i].RealTimeConfigCommandCount + 1));
                            Axis[i].RealTimeConfigArgs1High = bytes[1];
                            Axis[i].RealTimeConfigArgs1Low = bytes[0];
                            byte[] bytes2 = BitConverter.GetBytes(70);
                            Axis[i].RealTimeConfigArgs3High = bytes2[1];
                            Axis[i].RealTimeConfigArgs3Low = bytes2[0];
                            Axis[i].RealTimeConfigArgs2High = 0;
                            Axis[i].RealTimeConfigArgs2Low = 70;
                        }
                        finally
                        {
                            if (lockTaken)
                            {
                                Monitor.Exit(taskLock);
                            }
                        }

                        long timeOutTime = getTimeOutTime(1500L);
                        do
                        {
                            Thread.Sleep(1);
                            if (Axis[i].isRespondActual)
                            {
                                num2++;
                            }

                            if ((Axis[i].RealTimeConfigIDStatus == 2) | (Axis[i].RealTimeConfigIDStatus == 5))
                            {
                                num2 = 0;
                            }

                            if (!isTimeOut(timeOutTime))
                            {
                                continue;
                            }

                            object taskLock2 = TaskLock;
                            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                            bool lockTaken2 = false;
                            try
                            {
                                Monitor.Enter(taskLock2, ref lockTaken2);
                                Axis[i].DLLErrorText = "Timeout while waiting for curve response during start adding curve";
                            }
                            finally
                            {
                                if (lockTaken2)
                                {
                                    Monitor.Exit(taskLock2);
                                }
                            }

                            break;
                        }
                        while (!unchecked(Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount && num2 > 0));
                        array[num5] = Axis[i].RealTimeConfigStatusArgs2Low;
                        array[num5 + 1] = Axis[i].RealTimeConfigStatusArgs2High;
                        array[num5 + 2] = Axis[i].RealTimeConfigStatusArgs3Low;
                        array[num5 + 3] = Axis[i].RealTimeConfigStatusArgs3High;
                        num5 += 4;
                    }
                    while (num5 <= 70);
                    byte b = array[3];
                    uploadedCurve.Ycode = (byte)((b & 0xF0) >> 4);
                    uploadedCurve.Xcode = (byte)(b & 0xF);
                    uploadedCurve.SetpointCount = BitConverter.ToUInt16(new byte[2]
                    {
                        array[4],
                        array[5]
                    }, 0);
                    num = uploadedCurve.SetpointCount + 1;
                    uploadedCurve.CurveName = Conversions.ToString(Strings.Chr(array[8])) + Conversions.ToString(Strings.Chr(array[9])) + Conversions.ToString(Strings.Chr(array[10])) + Conversions.ToString(Strings.Chr(array[11])) + Conversions.ToString(Strings.Chr(array[12])) + Conversions.ToString(Strings.Chr(array[13])) + Conversions.ToString(Strings.Chr(array[14])) + Conversions.ToString(Strings.Chr(array[15])) + Conversions.ToString(Strings.Chr(array[16])) + Conversions.ToString(Strings.Chr(array[17])) + Conversions.ToString(Strings.Chr(array[18])) + Conversions.ToString(Strings.Chr(array[19])) + Conversions.ToString(Strings.Chr(array[20])) + Conversions.ToString(Strings.Chr(array[21])) + Conversions.ToString(Strings.Chr(array[22])) + Conversions.ToString(Strings.Chr(array[23])) + Conversions.ToString(Strings.Chr(array[24])) + Conversions.ToString(Strings.Chr(array[25])) + Conversions.ToString(Strings.Chr(array[26])) + Conversions.ToString(Strings.Chr(array[27])) + Conversions.ToString(Strings.Chr(array[28])) + Conversions.ToString(Strings.Chr(array[29]));
                    uploadedCurve.CurveName = uploadedCurve.CurveName.Replace("\0", "");
                    uploadedCurve.CurveID = BitConverter.ToUInt16(new byte[2]
                    {
                        array[30],
                        array[31]
                    }, 0);
                    uploadedCurve.XLength = BitConverter.ToUInt32(new byte[4]
                    {
                        array[32],
                        array[33],
                        array[34],
                        array[35]
                    }, 0);
                    uploadedCurve.YDim = BitConverter.ToUInt16(new byte[2]
                    {
                        array[36],
                        array[37]
                    }, 0);
                    uploadedCurve.XDim = BitConverter.ToUInt16(new byte[2]
                    {
                        array[38],
                        array[39]
                    }, 0);
                    int setpointCount = uploadedCurve.SetpointCount;
                    for (int j = 0; j <= setpointCount; j++)
                    {
                        object taskLock3 = TaskLock;
                        ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                        bool lockTaken3 = false;
                        try
                        {
                            Monitor.Enter(taskLock3, ref lockTaken3);
                            Axis[i].RealTimeConfigID = 98;
                            if (Axis[i].RealTimeConfigCommandCount > 14)
                            {
                                Axis[i].RealTimeConfigCommandCount = 1;
                            }

                            Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[i].RealTimeConfigCommandCount + 1));
                            Axis[i].RealTimeConfigArgs1High = bytes[1];
                            Axis[i].RealTimeConfigArgs1Low = bytes[0];
                            Axis[i].RealTimeConfigArgs3High = 0;
                            Axis[i].RealTimeConfigArgs3Low = 0;
                            Axis[i].RealTimeConfigArgs2High = 0;
                            Axis[i].RealTimeConfigArgs2Low = 0;
                        }
                        finally
                        {
                            if (lockTaken3)
                            {
                                Monitor.Exit(taskLock3);
                            }
                        }

                        long timeOutTime = getTimeOutTime(1500L);
                        do
                        {
                            Thread.Sleep(1);
                            if (Axis[i].isRespondActual)
                            {
                                num2++;
                            }

                            if ((Axis[i].RealTimeConfigIDStatus == 2) | (Axis[i].RealTimeConfigIDStatus == 5))
                            {
                                num2 = 0;
                            }

                            if (!isTimeOut(timeOutTime))
                            {
                                continue;
                            }

                            object taskLock4 = TaskLock;
                            ObjectFlowControl.CheckForSyncLockOnValueType(taskLock4);
                            bool lockTaken4 = false;
                            try
                            {
                                Monitor.Enter(taskLock4, ref lockTaken4);
                                Axis[i].DLLErrorText = "Timeout while waiting for curve response during start adding curve";
                            }
                            finally
                            {
                                if (lockTaken4)
                                {
                                    Monitor.Exit(taskLock4);
                                }
                            }

                            break;
                        }
                        while (!unchecked(Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount && num2 > 0));
                        byte[] value = new byte[4]
                        {
                            Axis[i].RealTimeConfigStatusArgs2Low,
                            Axis[i].RealTimeConfigStatusArgs2High,
                            Axis[i].RealTimeConfigStatusArgs3Low,
                            Axis[i].RealTimeConfigStatusArgs3High
                        };
                        uploadedCurve.Setpoints[uploadedCurve.Setpoints.Count() - 1] = BitConverter.ToInt32(value, 0);
                        ref int[] setpoints = ref uploadedCurve.Setpoints;
                        setpoints = (int[])Utils.CopyArray(setpoints, new int[uploadedCurve.Setpoints.Count() + 1]);
                        num4++;
                        object taskLock5 = TaskLock;
                        ObjectFlowControl.CheckForSyncLockOnValueType(taskLock5);
                        bool lockTaken5 = false;
                        try
                        {
                            Monitor.Enter(taskLock5, ref lockTaken5);
                            Axis[i].RealTimeConfigCurveStatus = (int)Math.Round(Math.Round((double)num4 / (double)num * 100.0));
                            if (Axis[i].RealTimeConfigCurveStatus > 100)
                            {
                                Axis[i].RealTimeConfigCurveStatus = 100;
                            }
                        }
                        finally
                        {
                            if (lockTaken5)
                            {
                                Monitor.Exit(taskLock5);
                            }
                        }
                    }

                    ref int[] setpoints2 = ref uploadedCurve.Setpoints;
                    setpoints2 = (int[])Utils.CopyArray(setpoints2, new int[uploadedCurve.Setpoints.Count() - 3 + 1]);
                    Axis[i].UploadedCurve = uploadedCurve;
                }
            }

            CurveThreadStart = false;
        }
    }

    private bool SaveDevice(object DeviceData, string FilePath)
    {
        FileStream fileStream = null;
        bool result;
        try
        {
            fileStream = new FileStream(FilePath, FileMode.Create);
            XmlSerializer xmlSerializer = new XmlSerializer(DeviceData.GetType());
            xmlSerializer.Serialize(fileStream, RuntimeHelpers.GetObjectValue(DeviceData));
            fileStream.Close();
            result = true;
        }
        catch (Exception ex)
        {
            ProjectData.SetProjectError(ex);
            Exception ex2 = ex;
            fileStream?.Close();
            result = false;
            ProjectData.ClearProjectError();
        }

        return result;
    }

    private object LoadDevice(string FilePath, object Data)
    {
        FileStream fileStream = null;
        try
        {
            fileStream = new FileStream(FilePath, FileMode.Open);
            XmlSerializer xmlSerializer = new XmlSerializer(Data.GetType());
            Data = RuntimeHelpers.GetObjectValue(xmlSerializer.Deserialize(fileStream));
            fileStream.Close();
        }
        catch (Exception ex)
        {
            ProjectData.SetProjectError(ex);
            Exception ex2 = ex;
            fileStream?.Close();
            ProjectData.ClearProjectError();
        }

        return Data;
    }

    //
    // Summary:
    //     This function save a command table directly from the addressed drive into a file.
    //
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   FilePath:
    //     Path of the configuration file
    //
    // Returns:
    //     True, if sucessful
    public bool Save_CommandTable(string TargetIP, string FilePath)
    {
        bool flag = true;
        CommandTableStructure commandTableStructure = default(CommandTableStructure);
        commandTableStructure = LMcf_getCommandTableContent(TargetIP);
        return SaveDevice(commandTableStructure, FilePath);
    }

    bool _ACI.Save_CommandTable(string TargetIP, string FilePath)
    {
        //ILSpy generated this explicit interface implementation from .override directive in Save_CommandTable
        return this.Save_CommandTable(TargetIP, FilePath);
    }

    //
    // Summary:
    //     This function saves a Command Table data structur into a file.
    //
    // Parameters:
    //   CT:
    //     Data structure CommandTableStructure containig the table
    //
    //   FilePath:
    //     Path of the configuration file
    //
    // Returns:
    //     True if successful
    public bool Save_CommandTable(CommandTableStructure CT, string FilePath)
    {
        bool flag = true;
        return SaveDevice(CT, FilePath);
    }

    bool _ACI.Save_CommandTable(CommandTableStructure CT, string FilePath)
    {
        //ILSpy generated this explicit interface implementation from .override directive in Save_CommandTable
        return this.Save_CommandTable(CT, FilePath);
    }

    //
    // Summary:
    //     This function loads a command table, which was saved with Save_CommandTable()
    //     directly into the addressed drive.
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   FilePath:
    //     Path of the configuration file
    //
    // Returns:
    //     True, if successful
    public bool Load_CommandTable(string TargetIP, string FilePath)
    {
        object obj = LoadDevice(FilePath, default(CommandTableStructure));
        CommandTableStructure cT = ((obj != null) ? ((CommandTableStructure)obj) : default(CommandTableStructure));
        return LMcf_setCommandTableContent(TargetIP, cT);
    }

    bool _ACI.Load_CommandTable(string TargetIP, string FilePath)
    {
        //ILSpy generated this explicit interface implementation from .override directive in Load_CommandTable
        return this.Load_CommandTable(TargetIP, FilePath);
    }

    //
    // Summary:
    //     This function loads a command table, which was saved with Save_CommandTable()
    //
    //
    // Parameters:
    //   FilePath:
    //
    // Returns:
    //     Data structure CommandTableStructure containig the table
    public CommandTableStructure Load_CommandTable(string FilePath)
    {
        object obj = LoadDevice(FilePath, default(CommandTableStructure));
        return (obj != null) ? ((CommandTableStructure)obj) : default(CommandTableStructure);
    }

    CommandTableStructure _ACI.Load_CommandTable(string FilePath)
    {
        //ILSpy generated this explicit interface implementation from .override directive in Load_CommandTable
        return this.Load_CommandTable(FilePath);
    }

    //
    // Summary:
    //     This function save all drive parameters, which are different to factory defaults
    //
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   FilePath:
    //     Path of the configuration file
    //
    // Returns:
    //     True if success
    public bool Save_DriveParameters(string TargetIP, string FilePath)
    {
        List<UPID_List> list = new List<UPID_List>();
        list = LMcf_getModified_UPIDList(TargetIP, 0, ushort.MaxValue);
        return SaveDevice(list, FilePath);
    }

    bool _ACI.Save_DriveParameters(string TargetIP, string FilePath)
    {
        //ILSpy generated this explicit interface implementation from .override directive in Save_DriveParameters
        return this.Save_DriveParameters(TargetIP, FilePath);
    }

    //
    // Summary:
    //     This function save a drive parameter list into a file
    //
    // Parameters:
    //   Data:
    //     Parameter list of UPID_List
    //
    //   FilePath:
    //     Path of the configuration file
    //
    // Returns:
    //     True, if success
    public bool Save_DriveParameters(List<UPID_List> Data, string FilePath)
    {
        return SaveDevice(Data, FilePath);
    }

    bool _ACI.Save_DriveParameters(List<UPID_List> Data, string FilePath)
    {
        //ILSpy generated this explicit interface implementation from .override directive in Save_DriveParameters
        return this.Save_DriveParameters(Data, FilePath);
    }

    //
    // Summary:
    //     This function loads drive parameters from a file directly into the drive, which
    //     was stored with Save_DriveParameters()
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   FilePath:
    //     Path of the configuration file, stored with Save_DrivePArameters()
    //
    // Returns:
    //     True, if successful
    public bool Load_DriveParameters(string TargetIP, string FilePath)
    {
        bool flag = false;
        List<UPID_List> data = new List<UPID_List>();
        data = (List<UPID_List>)LoadDevice(FilePath, data);
        LMcf_StartStopDefault(TargetIP, 5);
        LMcf_StartStopDefault(TargetIP, 1);
        LMcf_StartStopDefault(TargetIP, 2);
        LMcf_StartStopDefault(TargetIP, 3);
        LMcf_StartStopDefault(TargetIP, 4);
        flag = LMcf_setUPIDList(TargetIP, data);
        LMcf_StartStopDefault(TargetIP, 6);
        return flag;
    }

    bool _ACI.Load_DriveParameters(string TargetIP, string FilePath)
    {
        //ILSpy generated this explicit interface implementation from .override directive in Load_DriveParameters
        return this.Load_DriveParameters(TargetIP, FilePath);
    }

    //
    // Summary:
    //     This function loads drive parameters from a file, which was stored with Save_DriveParameters()
    //
    //
    // Parameters:
    //   FilePath:
    //     Path of the configuration file, stored with Save_DriveParameters()
    //
    // Returns:
    //     Drive parameters as list of UPID_List
    public List<UPID_List> Load_DriveParameters(string FilePAth)
    {
        List<UPID_List> data = new List<UPID_List>();
        return (List<UPID_List>)LoadDevice(FilePAth, data);
    }

    List<UPID_List> _ACI.Load_DriveParameters(string FilePAth)
    {
        //ILSpy generated this explicit interface implementation from .override directive in Load_DriveParameters
        return this.Load_DriveParameters(FilePAth);
    }

    //
    // Summary:
    //     This function saes all curves from the addressed drive into a file
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   FilePath:
    //     Path of the configuration file
    //
    // Returns:
    //     True if successfully saved
    public bool Save_Curves(string TargetIP, string FilePath)
    {
        int[] array = LMcf_getAllCurveID(TargetIP);
        List<CurveDataDefinition> list = new List<CurveDataDefinition>();
        int[] array2 = array;
        foreach (int num in array2)
        {
            CurveDataDefinition curveDataDefinition = default(CurveDataDefinition);
            LMcf_StartUploadCurve(TargetIP, checked((ushort)num));
            while (LMcf_isCurveLoading(TargetIP))
            {
            }

            curveDataDefinition = LMcf_getUploadedCurveData(TargetIP);
            list.Add(curveDataDefinition);
        }

        return SaveDevice(list, FilePath);
    }

    bool _ACI.Save_Curves(string TargetIP, string FilePath)
    {
        //ILSpy generated this explicit interface implementation from .override directive in Save_Curves
        return this.Save_Curves(TargetIP, FilePath);
    }

    //
    // Summary:
    //     This function loads all curves from a saved file into the addressed drive.
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   FilePath:
    //     >Path of the configuration file, stored with Save_Curves()
    //
    // Returns:
    //     True by default
    public bool Load_Curves(string TargetIP, string FilePath)
    {
        List<CurveDataDefinition> data = new List<CurveDataDefinition>();
        data = (List<CurveDataDefinition>)LoadDevice(FilePath, data);
        CurveDataDefinition curveData = default(CurveDataDefinition);
        curveData.CurveName = " ";
        LMcf_StartStopDefault(TargetIP, 5);
        LMcf_LoadCurve(TargetIP, 1, curveData);
        while (LMcf_isCurveLoading(TargetIP))
        {
        }

        foreach (CurveDataDefinition item in data)
        {
            LMcf_LoadCurve(TargetIP, 0, item);
            while (LMcf_isCurveLoading(TargetIP))
            {
            }
        }

        LMcf_LoadCurve(TargetIP, 2, curveData);
        while (LMcf_isCurveLoading(TargetIP))
        {
        }

        LMcf_StartStopDefault(TargetIP, 6);
        return true;
    }

    bool _ACI.Load_Curves(string TargetIP, string FilePath)
    {
        //ILSpy generated this explicit interface implementation from .override directive in Load_Curves
        return this.Load_Curves(TargetIP, FilePath);
    }

    //
    // Summary:
    //     This function saves the whole drive configuration, including curves and command
    //     table into a file.
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   FilePath:
    //     Path of the configuration file
    //
    // Returns:
    //     True, if successful
    public bool Save_DriveConfiguration(string TargetIP, string FilePath)
    {
        bool flag = false;
        DriveData driveData = default(DriveData);
        driveData.CommandTable = default(CommandTableStructure);
        driveData.Curves = new List<CurveDataDefinition>();
        driveData.Parameters = new List<UPID_List>();
        driveData.Parameters = LMcf_getModified_UPIDList(TargetIP, 0, ushort.MaxValue);
        driveData.CommandTable = LMcf_getCommandTableContent(TargetIP);
        int[] array = LMcf_getAllCurveID(TargetIP);
        List<CurveDataDefinition> list = new List<CurveDataDefinition>();
        int[] array2 = array;
        checked
        {
            foreach (int num in array2)
            {
                CurveDataDefinition curveDataDefinition = default(CurveDataDefinition);
                LMcf_StartUploadCurve(TargetIP, (ushort)num);
                while (LMcf_isCurveLoading(TargetIP))
                {
                }

                curveDataDefinition = LMcf_getUploadedCurveData(TargetIP);
                list.Add(curveDataDefinition);
            }

            driveData.Curves = list;
            uint hashVaue = (uint)getRAM_ByUPID(TargetIP, 161u);
            driveData.HashVaue = hashVaue;
            return SaveDevice(driveData, FilePath);
        }
    }

    bool _ACI.Save_DriveConfiguration(string TargetIP, string FilePath)
    {
        //ILSpy generated this explicit interface implementation from .override directive in Save_DriveConfiguration
        return this.Save_DriveConfiguration(TargetIP, FilePath);
    }

    //
    // Summary:
    //     This function loads a drive configuration into the addressed drive.
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   FilePath:
    //     Path of the configuration file, stored with Save_DriveConfiguration()
    //
    // Returns:
    //     True, if sucessfully loaded into the drive
    public bool Load_DriveConfiguration(string TargetIP, string FilePath)
    {
        bool flag = true;
        DriveData driveData = default(DriveData);
        driveData.CommandTable = default(CommandTableStructure);
        driveData.Curves = new List<CurveDataDefinition>();
        driveData.Parameters = new List<UPID_List>();
        object obj = LoadDevice(FilePath, driveData);
        driveData = ((obj != null) ? ((DriveData)obj) : default(DriveData));
        LMcf_StartStopDefault(TargetIP, 5);
        LMcf_StartStopDefault(TargetIP, 1);
        LMcf_StartStopDefault(TargetIP, 2);
        LMcf_StartStopDefault(TargetIP, 3);
        LMcf_StartStopDefault(TargetIP, 4);
        flag &= LMcf_setUPIDList(TargetIP, driveData.Parameters);
        LMcf_DeleteCommandTable_RAM(TargetIP);
        flag &= LMcf_setCommandTableContent(TargetIP, driveData.CommandTable);
        LMcf_WriteCommandTableToFLASH(TargetIP);
        flag &= LMcf_DeleteAllCurvesInRAM(TargetIP);
        foreach (CurveDataDefinition curf in driveData.Curves)
        {
            flag &= LMcf_setDownloadCurveData(TargetIP, curf);
            flag &= LMcf_StartDownloadCurve(TargetIP);
            DateTime dateTime = DateTime.Now.Add(new TimeSpan(0, 0, 59));
            while (LMcf_isCurveLoading(TargetIP) && dateTime.Ticks >= DateTime.Now.Ticks)
            {
            }
        }

        flag &= LMcf_SaveAllCurvesFromRAMToFLASH(TargetIP);
        LMcf_StartStopDefault(TargetIP, 6);
        return flag;
    }

    bool _ACI.Load_DriveConfiguration(string TargetIP, string FilePath)
    {
        //ILSpy generated this explicit interface implementation from .override directive in Load_DriveConfiguration
        return this.Load_DriveConfiguration(TargetIP, FilePath);
    }

    //
    // Summary:
    //     This functon delivers the drive HASH value
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     The HASH value of the drive
    public uint getDriveHash(string TargetIP)
    {
        return checked((uint)getRAM_ByUPID(TargetIP, 161u));
    }

    uint _ACI.getDriveHash(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in getDriveHash
        return this.getDriveHash(TargetIP);
    }

    //
    // Summary:
    //     This function compares the actual drive configuration with a stored configuration
    //     by use of the drive hash value.
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   FilePath:
    //     Path of the configuration file, stored with Save_DriveConfiguration()
    //
    // Returns:
    //     True, if configuration is same
    public bool isDriveConfigurationSame(string TargetIP, string FilePath)
    {
        DriveData driveData = default(DriveData);
        driveData.CommandTable = default(CommandTableStructure);
        driveData.Curves = new List<CurveDataDefinition>();
        driveData.Parameters = new List<UPID_List>();
        object obj = LoadDevice(FilePath, driveData);
        driveData = ((obj != null) ? ((DriveData)obj) : default(DriveData));
        uint num = checked((uint)getRAM_ByUPID(TargetIP, 161u));
        bool result = false;
        if (driveData.HashVaue == num)
        {
            result = true;
        }

        return result;
    }

    bool _ACI.isDriveConfigurationSame(string TargetIP, string FilePath)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isDriveConfigurationSame
        return this.isDriveConfigurationSame(TargetIP, FilePath);
    }

    //
    // Summary:
    //     This function delivers the error message of the actual error, which is stored
    //     on the drive.
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Error message text
    public string LMcf_GetErrorTxt(string TargetIP)
    {
        string result = "";
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                byte[] value = new byte[2]
                {
                    Axis[i].ErrorCodeHigh,
                    Axis[i].ErrorCodeLow
                };
                int num = BitConverter.ToInt16(value, 0);
                int num2 = num;
                if (num2 < 0)
                {
                    result = "";
                }
                else if (num2 == 0)
                {
                    result = "No error pending";
                }
                else if (num2 == 1)
                {
                    result = "X4 Logic supply too low";
                }
                else if (num2 == 2)
                {
                    result = "X4 Logic supply too high";
                }
                else if (num2 == 3)
                {
                    result = "X1 Power Voltage too low";
                }
                else if (num2 == 4)
                {
                    result = "X1 Power Voltage too high";
                }
                else if (num2 == 5)
                {
                    result = "X1 RR not connected";
                }
                else if (num2 == 6)
                {
                    result = "PTC 1 Semsor too hot";
                }
                else if (num2 == 7)
                {
                    result = "Min Pos undershot";
                }
                else if (num2 == 8)
                {
                    result = "Max Pos overshot";
                }
                else if (num2 == 9)
                {
                    result = "Ext-Int Sensor Diff Err";
                }
                else if (num2 == 10)
                {
                    result = "Fatal error X12 Signals missing";
                }
                else if (num2 == 11)
                {
                    result = "Pos Lag always too big";
                }
                else if (num2 == 12)
                {
                    result = "Pos Lag standing too big";
                }
                else if (num2 == 13)
                {
                    result = "X1 Power overcurrent";
                }
                else if (num2 == 14)
                {
                    result = "Supply Dig Out Missing";
                }
                else if (num2 == 15)
                {
                    result = "PTC 2 Sensor too hot";
                }
                else if (num2 == 16)
                {
                    result = "Drive Ph1+ too hot";
                }
                else if (num2 == 17)
                {
                    result = "Drive Ph1- too hot";
                }
                else if (num2 == 18)
                {
                    result = "Drive Ph2+ too hot";
                }
                else if (num2 == 19)
                {
                    result = "Drive Ph2- too hot";
                }
                else if (num2 == 20)
                {
                    result = "Drive Power too hot";
                }
                else if (num2 == 21)
                {
                    result = "Drive RR Hot Calc";
                }
                else if (num2 == 22)
                {
                    result = "Drive X3 too hot";
                }
                else if (num2 == 23)
                {
                    result = "Drive Core too hot";
                }
                else if (num2 == 24)
                {
                    result = "Power Bridge Ph1+ defective";
                }
                else if (num2 == 25)
                {
                    result = "Power Bridge Ph1- defective";
                }
                else if (num2 == 26)
                {
                    result = "Power Bridge Ph2+ defective";
                }
                else if (num2 == 27)
                {
                    result = "Power Bridge Ph2- defective";
                }
                else if (num2 == 28)
                {
                    result = "Supply DigOut X6 Fuse Blown";
                }
                else if (num2 == 29)
                {
                    result = "Supply X3.3 5V Fuse Blown";
                }
                else if (num2 == 30)
                {
                    result = "Supply X3.8 AGND Fuse Blown";
                }
                else if (num2 == 31)
                {
                    result = "";
                }
                else if (num2 == 32)
                {
                    result = "Motor Hot Sensor";
                }
                else if (num2 == 33)
                {
                    result = "X3 Hall Signal Missing";
                }
                else if (num2 == 34)
                {
                    result = "Motor Slider Missing";
                }
                else if (num2 == 35)
                {
                    result = "Motor Short Time Overload";
                }
                else if (num2 == 36)
                {
                    result = "RR Hot Calculated";
                }
                else if (num2 == 37)
                {
                    result = "Sensor Alarm";
                }
                else if (num2 == 38)
                {
                    result = "";
                }
                else if (num2 == 39)
                {
                    result = "";
                }
                else if (num2 == 40)
                {
                    result = "Ph1+ Short Circuit to GND";
                }
                else if (num2 == 41)
                {
                    result = "Ph1- Short Circuit to GND";
                }
                else if (num2 == 42)
                {
                    result = "Ph2+ Short Circuit to GND";
                }
                else if (num2 == 43)
                {
                    result = "Ph2- Short Circuit to GND";
                }
                else if (num2 == 44)
                {
                    result = "Ph1 Short Circuit to Ph2";
                }
                else if (num2 == 45)
                {
                    result = "";
                }
                else if (num2 == 46)
                {
                    result = "";
                }
                else if (num2 == 47)
                {
                    result = "";
                }
                else if (num2 == 48)
                {
                    result = "Ph1+ Wired to Ph2+";
                }
                else if (num2 == 49)
                {
                    result = "Ph1+ Wired to Ph2-";
                }
                else if (num2 == 50)
                {
                    result = "Ph1+ not wired to Ph1-";
                }
                else if (num2 == 51)
                {
                    result = "Ph1+ Wired to Ph1+";
                }
                else if (num2 == 52)
                {
                    result = "Ph2+ Wired to Ph1-";
                }
                else if (num2 == 53)
                {
                    result = "Ph2+ not wired to Ph2-";
                }
                else if (num2 == 54)
                {
                    result = "Ph1 short circuit to Ph2+";
                }
                else if (num2 == 55)
                {
                    result = "Ph1 short circuit to Ph2-";
                }
                else if (num2 == 56)
                {
                    result = "Ph2 short circuit to Ph1+";
                }
                else if (num2 == 57)
                {
                    result = "Ph2 short circuit to Ph1-";
                }
                else if (num2 == 58)
                {
                    result = "Phase U broken";
                }
                else if (num2 == 59)
                {
                    result = "Phase V broken";
                }
                else if (num2 == 60)
                {
                    result = "Phase W broken";
                }
                else if (num2 == 61)
                {
                    result = "";
                }
                else if (num2 == 62)
                {
                    result = "";
                }
                else if (num2 == 63)
                {
                    result = "";
                }
                else if (num2 == 64)
                {
                    result = "X4.3 Brake Driver Error";
                }
                else if (num2 == 65)
                {
                    result = "Dig Out X4.4..X4.11 Status";
                }
                else if (num2 == 66)
                {
                    result = "Dig Out X6 Status";
                }
                else if (num2 == 67)
                {
                    result = "";
                }
                else if (num2 == 68)
                {
                    result = "X4 Dig Out Defective";
                }
                else if (num2 == 69)
                {
                    result = "Fatal Err:Motor Comm Lost";
                }
                else if (num2 == 70)
                {
                    result = "PTC 1 Broken";
                }
                else if (num2 == 71)
                {
                    result = "PTC 1 Short To 24V";
                }
                else if (num2 >= 72 && num2 <= 79)
                {
                    result = "";
                }
                else if (num2 == 80)
                {
                    result = "HW Not Supported";
                }
                else if (num2 == 81)
                {
                    result = "SW Key Missing";
                }
                else if (num2 >= 82 && num2 <= 87)
                {
                    result = "";
                }
                else if (num2 == 88)
                {
                    result = "ROM write error";
                }
                else if (num2 >= 89 && num2 <= 95)
                {
                    result = "";
                }
                else if (num2 == 96)
                {
                    result = "RR Voltage Set Too Low";
                }
                else if (num2 == 97)
                {
                    result = "RR Hysteresis < 0.5V";
                }
                else if (num2 == 98)
                {
                    result = "Cfg Err:Curve Not Defined";
                }
                else if (num2 == 99)
                {
                    result = "Cfg Err:Pos Ctrl Max Curr High";
                }
                else if (num2 == 100)
                {
                    result = "Cfg Err (Fatal):No Motor Defined";
                }
                else if (num2 == 101)
                {
                    result = "Cfg Err (Fatal): Configuration error: No Trigger Mode Defined";
                }
                else if (num2 == 102)
                {
                    result = "";
                }
                else if (num2 == 103)
                {
                    result = "Cfg Err (Fatal):Wrong Stator Type";
                }
                else if (num2 == 104)
                {
                    result = "Cfg Err (Fatal):No Motor Communication";
                }
                else if (num2 == 105)
                {
                    result = "Cfg Err:Wrong Slider";
                }
                else if (num2 >= 106 && num2 <= 127)
                {
                    result = "";
                }
                else if (num2 == 128)
                {
                    result = "User Err:Lin: Not Homed";
                }
                else if (num2 == 129)
                {
                    result = "User Err:Unknown Motion Cmd";
                }
                else if (num2 == 130)
                {
                    result = "User Err:PVT Buffer Overflow";
                }
                else if (num2 == 131)
                {
                    result = "User Err:PVT Buffer Underflow";
                }
                else if (num2 == 132)
                {
                    result = "User Err:PVT Master Too Fast";
                }
                else if (num2 == 133)
                {
                    result = "User Err:PVT Master Too Slow";
                }
                else if (num2 == 134)
                {
                    result = "User Err: Motion Cmd In Wrong State";
                }
                else if (num2 == 135)
                {
                    result = "User Err:Limit Switch In High";
                }
                else if (num2 == 136)
                {
                    result = "User Err:Limit Switch Out High";
                }
                else if (num2 == 137)
                {
                    result = "User Err:Curve Amp Scale Error";
                }
                else if (num2 == 138)
                {
                    result = "User Err:Cmd Tab Entry Not Def";
                }
                else if (num2 > 138)
                {
                    result = "";
                }
            }
        }

        return result;
    }

    string _ACI.LMcf_GetErrorTxt(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMcf_GetErrorTxt
        return this.LMcf_GetErrorTxt(TargetIP);
    }

    //
    // Summary:
    //     Get the error code number, if an error is active
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Error code number
    public int LMcf_GetErrorCode(string TargetIP)
    {
        int result = 0;
        uint num = 0u;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if ((Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0) & isError(TargetIP))
            {
                byte[] value = new byte[4]
                {
                    Axis[i].ErrorCodeHigh,
                    Axis[i].ErrorCodeLow,
                    0,
                    0
                };
                result = BitConverter.ToInt32(value, 0);
            }
        }

        return result;
    }

    int _ACI.LMcf_GetErrorCode(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMcf_GetErrorCode
        return this.LMcf_GetErrorCode(TargetIP);
    }

    //
    // Summary:
    //     Get the error code number, if an error is active
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Error code number
    public int LMcf_GetWarningCode(string TargetIP)
    {
        int result = 0;
        uint num = 0u;
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if ((Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0) & isWarning(TargetIP))
            {
                byte[] value = new byte[4]
                {
                    Axis[i].WarnWordHigh,
                    Axis[i].WarnWordLow,
                    0,
                    0
                };
                result = BitConverter.ToInt32(value, 0);
            }
        }

        return result;
    }

    int _ACI.LMcf_GetWarningCode(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMcf_GetWarningCode
        return this.LMcf_GetWarningCode(TargetIP);
    }

    //
    // Summary:
    //     Get warning message text
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Warning Message
    public string LMcf_GetWarningTxt(string TargetIP)
    {
        string text = "";
        int targetIPListCount = TargetIPListCount;
        for (int i = 0; i <= targetIPListCount; i = checked(i + 1))
        {
            if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) == 0)
            {
                Console.WriteLine("Low: " + Axis[i].WarnWordLow);
                Console.WriteLine("High: " + Axis[i].WarnWordHigh);
                if ((Axis[i].WarnWordHigh & 1) == 1)
                {
                    text += " Motor Temperature Sensor On";
                }

                if ((Axis[i].WarnWordHigh & 2) == 2)
                {
                    text += " Calculated Motor Temperature Reached Warn Limit";
                }

                if ((Axis[i].WarnWordHigh & 4) == 4)
                {
                    text += " Motor Supply Voltage Reached Low Warn Limit";
                }

                if ((Axis[i].WarnWordHigh & 8) == 8)
                {
                    text += " Motor Supply Voltage Reached High Warn Limit";
                }

                if ((Axis[i].WarnWordHigh & 0x10) == 16)
                {
                    text += " Position Error during Moving Reached Warn Limit";
                }

                if ((Axis[i].WarnWordHigh & 0x40) == 64)
                {
                    text += " Temperature on Servo Controller High";
                }

                if ((Axis[i].WarnWordHigh & 0x80) == 128)
                {
                    text += " Warning Motor Not Homed Yet";
                }

                if ((Axis[i].WarnWordLow & 1) == 1)
                {
                    text += " PTC Temperature Sensor 1 On";
                }

                if ((Axis[i].WarnWordLow & 2) == 2)
                {
                    text += " PTC Temperature Sensor 2 On";
                }

                if ((Axis[i].WarnWordLow & 4) == 4)
                {
                    text += " Regeneration Resistor Temperature Hot Calculated";
                }

                if ((Axis[i].WarnWordLow & 0x40) == 64)
                {
                    text += " Warn Flag Of Interface SW Layer";
                }

                if ((Axis[i].WarnWordLow & 0x80) == 128)
                {
                    text += " Warn Flag Of Application SW Layer";
                }
            }
        }

        return text;
    }

    string _ACI.LMcf_GetWarningTxt(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMcf_GetWarningTxt
        return this.LMcf_GetWarningTxt(TargetIP);
    }

    //
    // Summary:
    //     Check, if the state machine of Master and all Slaves are in Operation Enabled
    //     state
    //
    // Parameters:
    //   TargetIP:
    //     IP Adress for accessing the drive with DLL functions (1..n)
    //
    // Returns:
    //     True, if all drives are ready to execute motion commands
    public bool isMastmerSlaveOperationEnabled(string TargetIP)
    {
        long rAM_ByUPID = getRAM_ByUPID(TargetIP, 12512u);
        int num = 0;
        checked
        {
            int num2 = (int)rAM_ByUPID - 1;
            for (int i = 0; i <= num2; i++)
            {
                if (i == 0)
                {
                    long rAM_ByUPID2 = getRAM_ByUPID(TargetIP, 15216u);
                    byte[] bytes = BitConverter.GetBytes(rAM_ByUPID2);
                    if (bytes[1] == 8)
                    {
                        num++;
                    }
                }

                if (i == 1)
                {
                    long rAM_ByUPID3 = getRAM_ByUPID(TargetIP, 15218u);
                    byte[] bytes2 = BitConverter.GetBytes(rAM_ByUPID3);
                    if (bytes2[1] == 8)
                    {
                        num++;
                    }
                }

                if (i == 2)
                {
                    long rAM_ByUPID4 = getRAM_ByUPID(TargetIP, 15220u);
                    byte[] bytes3 = BitConverter.GetBytes(rAM_ByUPID4);
                    if (bytes3[1] == 8)
                    {
                        num++;
                    }
                }
            }

            int targetIPListCount = TargetIPListCount;
            bool result = default(bool);
            for (int j = 0; j <= targetIPListCount; j++)
            {
                if (Operators.CompareString(Axis[j].SlaveIP, TargetIP, TextCompare: false) == 0)
                {
                    result = (unchecked(Axis[j].StateVarHigh == 8 && num == rAM_ByUPID) ? true : false);
                }
            }

            return result;
        }
    }

    bool _ACI.isMastmerSlaveOperationEnabled(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in isMastmerSlaveOperationEnabled
        return this.isMastmerSlaveOperationEnabled(TargetIP);
    }

    //
    // Summary:
    //     Start homing procedure for Master/Slave applications
    //
    // Parameters:
    //   TargetIP:
    //     Number for accessing the drive with DLL functions (1..n)
    //
    // Returns:
    //     Always True
    public bool MasterSlaveHoming(string TargetIP)
    {
        bool result = false;
        long num = 0L;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0 || !isOperationEnabledSM(TargetIP))
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh | 8);
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                Thread.Sleep(200);
                num = getTimeOutTime(300000L);
                while ((Axis[i].StatusWordHigh & 8) != 8)
                {
                    Thread.Sleep(15);
                    if (!isTimeOut(num))
                    {
                        continue;
                    }

                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        Axis[i].DLLErrorText = "Timeout (5min) during homing procedure!1";
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    break;
                }

                num = getTimeOutTime(30000L);
                long rAM_ByUPID = getRAM_ByUPID(TargetIP, 12512u);
                long num2 = 0L;
                while (num2 != rAM_ByUPID)
                {
                    Thread.Sleep(15);
                    if (isTimeOut(num))
                    {
                        object taskLock3 = TaskLock;
                        ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                        bool lockTaken3 = false;
                        try
                        {
                            Monitor.Enter(taskLock3, ref lockTaken3);
                            Axis[i].DLLErrorText = "Timeout (5min) during homing procedure!2";
                        }
                        finally
                        {
                            if (lockTaken3)
                            {
                                Monitor.Exit(taskLock3);
                            }
                        }

                        break;
                    }

                    int num3 = (int)(rAM_ByUPID - 1);
                    for (int j = 0; j <= num3; j++)
                    {
                        if (j == 0)
                        {
                            long rAM_ByUPID2 = getRAM_ByUPID(TargetIP, 15216u);
                            byte[] bytes = BitConverter.GetBytes(rAM_ByUPID2);
                            if ((bytes[0] == 15) & (bytes[1] == 9))
                            {
                                num2++;
                            }
                        }

                        if (j == 1)
                        {
                            long rAM_ByUPID3 = getRAM_ByUPID(TargetIP, 15218u);
                            byte[] bytes2 = BitConverter.GetBytes(rAM_ByUPID3);
                            if ((bytes2[0] == 15) & (bytes2[1] == 9))
                            {
                                num2++;
                            }
                        }

                        if (j == 2)
                        {
                            long rAM_ByUPID4 = getRAM_ByUPID(TargetIP, 15220u);
                            byte[] bytes3 = BitConverter.GetBytes(rAM_ByUPID4);
                            if ((bytes3[0] == 15) & (bytes3[1] == 9))
                            {
                                num2++;
                            }
                        }
                    }
                }

                object taskLock4 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock4);
                bool lockTaken4 = false;
                try
                {
                    Monitor.Enter(taskLock4, ref lockTaken4);
                    Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh & 0xF7);
                }
                finally
                {
                    if (lockTaken4)
                    {
                        Monitor.Exit(taskLock4);
                    }
                }

                num = getTimeOutTime(5000L);
                while (Axis[i].StateVarHigh == 9)
                {
                    object taskLock5 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock5);
                    bool lockTaken5 = false;
                    try
                    {
                        Monitor.Enter(taskLock5, ref lockTaken5);
                        Axis[i].ControlWordHigh = (byte)(Axis[i].ControlWordHigh & 0xF7);
                    }
                    finally
                    {
                        if (lockTaken5)
                        {
                            Monitor.Exit(taskLock5);
                        }
                    }

                    if (isTimeOut(num))
                    {
                        object taskLock6 = TaskLock;
                        ObjectFlowControl.CheckForSyncLockOnValueType(taskLock6);
                        bool lockTaken6 = false;
                        try
                        {
                            Monitor.Enter(taskLock6, ref lockTaken6);
                            Axis[i].DLLErrorText = "Timeout (5min) during homing procedure!3";
                        }
                        finally
                        {
                            if (lockTaken6)
                            {
                                Monitor.Exit(taskLock6);
                            }
                        }

                        break;
                    }

                    Thread.Sleep(15);
                }

                result = true;
            }

            return result;
        }
    }

    bool _ACI.MasterSlaveHoming(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in MasterSlaveHoming
        return this.MasterSlaveHoming(TargetIP);
    }

    private int LogMotionToCSV(string IPtoLog, string FileName, bool Run, uint LoggingInterval = 5u)
    {
        int num = 0;
        LogIP = IPtoLog;
        if ((long)LoggingInterval < 5L)
        {
            LoggingInterval = 5u;
        }

        LogInterval = checked((int)LoggingInterval);
        if (Operators.CompareString(FileName, "", TextCompare: false) == 0)
        {
            num = 3;
        }
        else
        {
            LogFileName = FileName;
            num = 0;
        }

        if (Run & !RunLogActive)
        {
            try
            {
                string text = "Timestamp;Actual Pos;Actual Current; Monitoring 2;Monitoring 3;Monitoring 4\r\n";
                MyProject.Computer.FileSystem.WriteAllText(LogFileName, text, append: true);
            }
            catch (Exception ex)
            {
                ProjectData.SetProjectError(ex);
                Exception ex2 = ex;
                ProjectData.ClearProjectError();
            }

            RunLogActive = true;
            LogItStart();
            num = 0;
        }

        if (!Run & RunLogActive)
        {
            RunLogActive = false;
            num = 2;
        }

        if (Run & RunLogActive)
        {
            num = 1;
        }

        return num;
    }

    private void LogItStart()
    {
        if (RunLogActive)
        {
            LogIt = new Thread([SpecialName][DebuggerHidden] () =>
            {
                RunLog();
            });
            LogIt.Name = "Request data for log";
            LogIt.Start();
        }
    }

    private int RunLog()
    {
        try
        {
            oldTimestampLog = newTimestampLog;
            newTimestampLog = getMonitoringChannel1(LogIP);
            if (oldTimestampLog != newTimestampLog)
            {
                string obj = getMonitoringChannel1(LogIP) + ";" + getActualPos(LogIP) + ";" + getCurrent(LogIP) + ";" + getMonitoringChannel2(LogIP) + ";" + getMonitoringChannel3(LogIP) + ";" + getMonitoringChannel4(LogIP) + "\r\n";
                LogLine.Enqueue(obj);
                LogFileWrite();
            }
        }
        catch (Exception ex)
        {
            ProjectData.SetProjectError(ex);
            Exception ex2 = ex;
            ProjectData.ClearProjectError();
        }

        Thread.Sleep(LogInterval);
        LogItStart();
        return 0;
    }

    private void LogFileWrite()
    {
        if (RunLogActive)
        {
            LogFileWriter = new Thread(LogFileWriteIt);
            LogFileWriter.Name = "Log File Writer Task";
            LogFileWriter.Start();
        }
    }

    private void LogFileWriteIt()
    {
        try
        {
            if (LogLine.Count > 0)
            {
                string text = LogLine.Dequeue().ToString();
                MyProject.Computer.FileSystem.WriteAllText(LogFileName, text, append: true);
            }

            Thread.Sleep(2);
            if (LogLine.Count > 0)
            {
                LogFileWrite();
            }
        }
        catch (Exception ex)
        {
            ProjectData.SetProjectError(ex);
            Exception ex2 = ex;
            ProjectData.ClearProjectError();
        }
    }

    private long StartGettingUPIDList(string TargetIP, ushort StartUPID)
    {
        long result = 0L;
        byte[] bytes = BitConverter.GetBytes(StartUPID);
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].RealTimeConfigID = 32;
                    if (Axis[i].RealTimeConfigCommandCount > 14)
                    {
                        Axis[i].RealTimeConfigCommandCount = 1;
                    }

                    Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[i].RealTimeConfigCommandCount + 1));
                    Axis[i].RealTimeConfigCommandCount = Axis[i].RealTimeConfigCommandCount;
                    Axis[i].RealTimeConfigArgs1High = bytes[1];
                    Axis[i].RealTimeConfigArgs1Low = bytes[0];
                    Axis[i].RealTimeConfigArgs2High = 0;
                    Axis[i].RealTimeConfigArgs2Low = 0;
                    Axis[i].RealTimeConfigArgs3High = 0;
                    Axis[i].RealTimeConfigArgs3Low = 0;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                long timeOutTime = getTimeOutTime(1000L);
                while (!((Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount) & (Axis[i].RealTimeConfigIDStatus == 0)))
                {
                    if (!isTimeOut(timeOutTime))
                    {
                        continue;
                    }

                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        Axis[i].DLLErrorText = "7: Config Channel Timeout by requesting ROM value by UPID. Error code: " + Conversion.Hex(Axis[i].RealTimeConfigIDStatus).ToString();
                        ACIError = "Config Channel Error";
                        Axis[i].RealTimeConfigStatusArgs2Low = 9;
                        Axis[i].RealTimeConfigStatusArgs2High = 0;
                        Axis[i].RealTimeConfigStatusArgs3Low = 9;
                        Axis[i].RealTimeConfigStatusArgs3High = 0;
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    break;
                }

                object taskLock3 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                bool lockTaken3 = false;
                try
                {
                    Monitor.Enter(taskLock3, ref lockTaken3);
                    byte[] value = new byte[8]
                    {
                        Axis[i].RealTimeConfigStatusArgs2Low,
                        Axis[i].RealTimeConfigStatusArgs2High,
                        Axis[i].RealTimeConfigStatusArgs3Low,
                        Axis[i].RealTimeConfigStatusArgs3High,
                        0,
                        0,
                        0,
                        0
                    };
                    Axis[i].RealTimeConfigID = 0;
                    result = (long)BitConverter.ToUInt64(value, 0);
                }
                finally
                {
                    if (lockTaken3)
                    {
                        Monitor.Exit(taskLock3);
                    }
                }
            }

            return result;
        }
    }

    private long StartGettingModifiedUPIDList(string TargetIP, ushort StartUPID)
    {
        long result = 0L;
        byte[] bytes = BitConverter.GetBytes(StartUPID);
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].RealTimeConfigID = 34;
                    if (Axis[i].RealTimeConfigCommandCount > 14)
                    {
                        Axis[i].RealTimeConfigCommandCount = 1;
                    }

                    Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[i].RealTimeConfigCommandCount + 1));
                    Axis[i].RealTimeConfigCommandCount = Axis[i].RealTimeConfigCommandCount;
                    Axis[i].RealTimeConfigArgs1High = bytes[1];
                    Axis[i].RealTimeConfigArgs1Low = bytes[0];
                    Axis[i].RealTimeConfigArgs2High = 0;
                    Axis[i].RealTimeConfigArgs2Low = 0;
                    Axis[i].RealTimeConfigArgs3High = 0;
                    Axis[i].RealTimeConfigArgs3Low = 0;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                long timeOutTime = getTimeOutTime(1000L);
                while (!((Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount) & (Axis[i].RealTimeConfigIDStatus == 0)))
                {
                    if (!isTimeOut(timeOutTime))
                    {
                        continue;
                    }

                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        Axis[i].DLLErrorText = "7: Config Channel Timeout by requesting ROM value by UPID. Error code: " + Conversion.Hex(Axis[i].RealTimeConfigIDStatus).ToString();
                        ACIError = "Config Channel Error";
                        Axis[i].RealTimeConfigStatusArgs2Low = 9;
                        Axis[i].RealTimeConfigStatusArgs2High = 0;
                        Axis[i].RealTimeConfigStatusArgs3Low = 9;
                        Axis[i].RealTimeConfigStatusArgs3High = 0;
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    break;
                }

                object taskLock3 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                bool lockTaken3 = false;
                try
                {
                    Monitor.Enter(taskLock3, ref lockTaken3);
                    byte[] value = new byte[8]
                    {
                        Axis[i].RealTimeConfigStatusArgs2Low,
                        Axis[i].RealTimeConfigStatusArgs2High,
                        Axis[i].RealTimeConfigStatusArgs3Low,
                        Axis[i].RealTimeConfigStatusArgs3High,
                        0,
                        0,
                        0,
                        0
                    };
                    Axis[i].RealTimeConfigID = 0;
                    result = (long)BitConverter.ToUInt64(value, 0);
                }
                finally
                {
                    if (lockTaken3)
                    {
                        Monitor.Exit(taskLock3);
                    }
                }
            }

            return result;
        }
    }

    private ReturnValueUPIDList getNextUPIDListItem(string TargetIP)
    {
        ReturnValueUPIDList result = default(ReturnValueUPIDList);
        result.value = 0L;
        result.State = 0;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].RealTimeConfigID = 33;
                    if (Axis[i].RealTimeConfigCommandCount > 14)
                    {
                        Axis[i].RealTimeConfigCommandCount = 1;
                    }

                    Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[i].RealTimeConfigCommandCount + 1));
                    Axis[i].RealTimeConfigArgs1High = 0;
                    Axis[i].RealTimeConfigArgs1Low = 0;
                    Axis[i].RealTimeConfigArgs2High = 0;
                    Axis[i].RealTimeConfigArgs2Low = 0;
                    Axis[i].RealTimeConfigArgs3High = 0;
                    Axis[i].RealTimeConfigArgs3Low = 0;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                long timeOutTime = getTimeOutTime(5000L);
                while (!((Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount) & (Axis[i].RealTimeConfigIDStatus == 0)) && Axis[i].RealTimeConfigIDStatus != 198)
                {
                    if (!isTimeOut(timeOutTime))
                    {
                        continue;
                    }

                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        Axis[i].DLLErrorText = "8: Config Channel Timeout by requesting RAM value by UPID. Error code: " + Conversion.Hex(Axis[i].RealTimeConfigIDStatus).ToString();
                        ACIError = "Config Channel Error";
                        Axis[i].RealTimeConfigStatusArgs2Low = 9;
                        Axis[i].RealTimeConfigStatusArgs2High = 0;
                        Axis[i].RealTimeConfigStatusArgs3Low = 9;
                        Axis[i].RealTimeConfigStatusArgs3High = 0;
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    break;
                }

                object taskLock3 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                bool lockTaken3 = false;
                try
                {
                    Monitor.Enter(taskLock3, ref lockTaken3);
                    byte[] value = new byte[2]
                    {
                        Axis[i].RealTimeConfigStatusArgs1Low,
                        Axis[i].RealTimeConfigStatusArgs1High
                    };
                    byte[] value2 = new byte[8]
                    {
                        Axis[i].RealTimeConfigStatusArgs2Low,
                        Axis[i].RealTimeConfigStatusArgs2High,
                        Axis[i].RealTimeConfigStatusArgs3Low,
                        Axis[i].RealTimeConfigStatusArgs3High,
                        0,
                        0,
                        0,
                        0
                    };
                    Axis[i].RealTimeConfigID = 0;
                    result.State = Axis[i].RealTimeConfigIDStatus;
                    result.UPID = BitConverter.ToUInt16(value, 0);
                    result.value = (long)BitConverter.ToUInt64(value2, 0);
                }
                finally
                {
                    if (lockTaken3)
                    {
                        Monitor.Exit(taskLock3);
                    }
                }
            }

            return result;
        }
    }

    private ReturnValueUPIDList getNextModifiedUPIDListItem(string TargetIP)
    {
        ReturnValueUPIDList result = default(ReturnValueUPIDList);
        result.value = 0L;
        result.State = 0;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].RealTimeConfigID = 35;
                    if (Axis[i].RealTimeConfigCommandCount > 14)
                    {
                        Axis[i].RealTimeConfigCommandCount = 1;
                    }

                    Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[i].RealTimeConfigCommandCount + 1));
                    Axis[i].RealTimeConfigArgs1High = 0;
                    Axis[i].RealTimeConfigArgs1Low = 0;
                    Axis[i].RealTimeConfigArgs2High = 0;
                    Axis[i].RealTimeConfigArgs2Low = 0;
                    Axis[i].RealTimeConfigArgs3High = 0;
                    Axis[i].RealTimeConfigArgs3Low = 0;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                long timeOutTime = getTimeOutTime(5000L);
                while (!((Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount) & (Axis[i].RealTimeConfigIDStatus == 0)) && Axis[i].RealTimeConfigIDStatus != 198)
                {
                    if (!isTimeOut(timeOutTime))
                    {
                        continue;
                    }

                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        Axis[i].DLLErrorText = "8: Config Channel Timeout by requesting RAM value by UPID. Error code: " + Conversion.Hex(Axis[i].RealTimeConfigIDStatus).ToString();
                        ACIError = "Config Channel Error";
                        Axis[i].RealTimeConfigStatusArgs2Low = 9;
                        Axis[i].RealTimeConfigStatusArgs2High = 0;
                        Axis[i].RealTimeConfigStatusArgs3Low = 9;
                        Axis[i].RealTimeConfigStatusArgs3High = 0;
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    break;
                }

                object taskLock3 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                bool lockTaken3 = false;
                try
                {
                    Monitor.Enter(taskLock3, ref lockTaken3);
                    byte[] value = new byte[2]
                    {
                        Axis[i].RealTimeConfigStatusArgs1Low,
                        Axis[i].RealTimeConfigStatusArgs1High
                    };
                    byte[] value2 = new byte[8]
                    {
                        Axis[i].RealTimeConfigStatusArgs2Low,
                        Axis[i].RealTimeConfigStatusArgs2High,
                        Axis[i].RealTimeConfigStatusArgs3Low,
                        Axis[i].RealTimeConfigStatusArgs3High,
                        0,
                        0,
                        0,
                        0
                    };
                    Axis[i].RealTimeConfigID = 0;
                    result.State = Axis[i].RealTimeConfigIDStatus;
                    result.UPID = BitConverter.ToUInt16(value, 0);
                    result.value = (long)BitConverter.ToUInt64(value2, 0);
                }
                finally
                {
                    if (lockTaken3)
                    {
                        Monitor.Exit(taskLock3);
                    }
                }
            }

            return result;
        }
    }

    //
    // Summary:
    //     This function delivers a list of the available UPID adresses with their address
    //     usage (see manual drive configuration over Fieldbus for details), starting from
    //     StartUPID
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   StartUPID:
    //     Start UPID address value
    //
    //   StopUPID:
    //     Stop UPID address value>
    //
    // Returns:
    //     UPID list containing adddress and it's usage
    public List<UPID_List> LMcf_getUPIDList(string TargetIP, ushort StartUPID, ushort StopUPID)
    {
        List<UPID_List> list = new List<UPID_List>();
        StartGettingUPIDList(TargetIP, StartUPID);
        byte b = 0;
        bool flag = false;
        while (b == 0)
        {
            ReturnValueUPIDList nextUPIDListItem = getNextUPIDListItem(TargetIP);
            b = nextUPIDListItem.State;
            if (b == 198)
            {
                flag = true;
                break;
            }

            UPID_List item = default(UPID_List);
            item.UPID = nextUPIDListItem.UPID;
            item.RawDataValue = nextUPIDListItem.value;
            if ((uint)item.UPID > (uint)StopUPID)
            {
                flag = true;
                break;
            }

            list.Add(item);
        }

        if (b != 198 && !flag)
        {
            UPID_List item2 = default(UPID_List);
            item2.UPID = ushort.MaxValue;
            item2.RawDataValue = -1L;
            list.Clear();
            list.Add(item2);
        }

        return list;
    }

    List<UPID_List> _ACI.LMcf_getUPIDList(string TargetIP, ushort StartUPID, ushort StopUPID)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMcf_getUPIDList
        return this.LMcf_getUPIDList(TargetIP, StartUPID, StopUPID);
    }

    //
    // Summary:
    //     This function delivers all UPID's and their raw data, which are not on factory
    //     default. This can be used to stroe drive parameters in the master application.
    //
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   StartUPID:
    //     Start UPID address value
    //
    //   StopUPID:
    //     Stop UPID address value
    //
    // Returns:
    //     UPID list containing UPID and the raw data value
    public List<UPID_List> LMcf_getModified_UPIDList(string TargetIP, ushort StartUPID, ushort StopUPID)
    {
        List<UPID_List> list = new List<UPID_List>();
        StartGettingModifiedUPIDList(TargetIP, StartUPID);
        byte b = 0;
        bool flag = false;
        while (b == 0)
        {
            ReturnValueUPIDList nextModifiedUPIDListItem = getNextModifiedUPIDListItem(TargetIP);
            b = nextModifiedUPIDListItem.State;
            if (b == 198)
            {
                flag = true;
                break;
            }

            UPID_List item = default(UPID_List);
            item.UPID = nextModifiedUPIDListItem.UPID;
            item.RawDataValue = nextModifiedUPIDListItem.value;
            if ((uint)item.UPID > (uint)StopUPID)
            {
                flag = true;
                break;
            }

            list.Add(item);
        }

        if (b != 198 && !flag)
        {
            UPID_List item2 = default(UPID_List);
            item2.UPID = ushort.MaxValue;
            item2.RawDataValue = -1L;
            list.Clear();
            list.Add(item2);
        }

        return list;
    }

    List<UPID_List> _ACI.LMcf_getModified_UPIDList(string TargetIP, ushort StartUPID, ushort StopUPID)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMcf_getModified_UPIDList
        return this.LMcf_getModified_UPIDList(TargetIP, StartUPID, StopUPID);
    }

    //
    // Summary:
    //     This function can be used to write stored drive parameters back to a new or defautled
    //     drive.
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   UPIDList:
    //     UPID list, containing UPID address and their corresponding raw data
    //
    // Returns:
    //     True, if successful
    public bool LMcf_setUPIDList(string TargetIP, List<UPID_List> UPIDList)
    {
        bool flag = true;
        foreach (UPID_List UPID in UPIDList)
        {
            if (SetROM_ByUPID(TargetIP, UPID.UPID, UPID.RawDataValue) != 0)
            {
                flag = false;
            }
        }

        return !flag;
    }

    bool _ACI.LMcf_setUPIDList(string TargetIP, List<UPID_List> UPIDList)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMcf_setUPIDList
        return this.LMcf_setUPIDList(TargetIP, UPIDList);
    }

    private byte[] getCT_PresenceListList(string TargetIP, byte GetEntry)
    {
        byte[] result = new byte[4];
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].RealTimeConfigID = GetEntry;
                    if (Axis[i].RealTimeConfigCommandCount > 14)
                    {
                        Axis[i].RealTimeConfigCommandCount = 1;
                    }

                    Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[i].RealTimeConfigCommandCount + 1));
                    Axis[i].RealTimeConfigCommandCount = Axis[i].RealTimeConfigCommandCount;
                    Axis[i].RealTimeConfigArgs1High = 0;
                    Axis[i].RealTimeConfigArgs1Low = 0;
                    Axis[i].RealTimeConfigArgs2High = 0;
                    Axis[i].RealTimeConfigArgs2Low = 0;
                    Axis[i].RealTimeConfigArgs3High = 0;
                    Axis[i].RealTimeConfigArgs3Low = 0;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                long timeOutTime = getTimeOutTime(1000L);
                while (!((Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount) & (Axis[i].RealTimeConfigIDStatus == 0)))
                {
                    if (!isTimeOut(timeOutTime))
                    {
                        continue;
                    }

                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        Axis[i].DLLErrorText = "7: Config Channel Timeout by requesting ROM value by UPID. Error code: " + Conversion.Hex(Axis[i].RealTimeConfigIDStatus).ToString();
                        ACIError = "Config Channel Error";
                        Axis[i].RealTimeConfigStatusArgs2Low = 9;
                        Axis[i].RealTimeConfigStatusArgs2High = 0;
                        Axis[i].RealTimeConfigStatusArgs3Low = 9;
                        Axis[i].RealTimeConfigStatusArgs3High = 0;
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    break;
                }

                object taskLock3 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                bool lockTaken3 = false;
                try
                {
                    Monitor.Enter(taskLock3, ref lockTaken3);
                    byte[] array = new byte[4]
                    {
                        Axis[i].RealTimeConfigStatusArgs2Low,
                        Axis[i].RealTimeConfigStatusArgs2High,
                        Axis[i].RealTimeConfigStatusArgs3Low,
                        Axis[i].RealTimeConfigStatusArgs3High
                    };
                    Axis[i].RealTimeConfigID = 0;
                    result = array;
                }
                finally
                {
                    if (lockTaken3)
                    {
                        Monitor.Exit(taskLock3);
                    }
                }
            }

            return result;
        }
    }

    private CTEntryDataDefinition getCTEntry(string TargetIP, ushort EntryNumber)
    {
        CTEntryDataDefinition result = default(CTEntryDataDefinition);
        byte[] bytes = BitConverter.GetBytes(EntryNumber);
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].RealTimeConfigID = 133;
                    if (Axis[i].RealTimeConfigCommandCount > 14)
                    {
                        Axis[i].RealTimeConfigCommandCount = 1;
                    }

                    Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[i].RealTimeConfigCommandCount + 1));
                    Axis[i].RealTimeConfigArgs1High = bytes[1];
                    Axis[i].RealTimeConfigArgs1Low = bytes[0];
                    Axis[i].RealTimeConfigArgs2High = 0;
                    Axis[i].RealTimeConfigArgs2Low = 0;
                    Axis[i].RealTimeConfigArgs3High = 0;
                    Axis[i].RealTimeConfigArgs3Low = 0;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                long timeOutTime = getTimeOutTime(5000L);
                while (!((Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount) & (Axis[i].RealTimeConfigIDStatus == 0)))
                {
                    if (!isTimeOut(timeOutTime))
                    {
                        continue;
                    }

                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        Axis[i].DLLErrorText = "8: Config Channel Timeout by requesting Command Table Data.. Error code: " + Conversion.Hex(Axis[i].RealTimeConfigIDStatus).ToString();
                        ACIError = "Config Channel Error";
                        Axis[i].RealTimeConfigStatusArgs2Low = 9;
                        Axis[i].RealTimeConfigStatusArgs2High = 0;
                        Axis[i].RealTimeConfigStatusArgs3Low = 9;
                        Axis[i].RealTimeConfigStatusArgs3High = 0;
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    break;
                }

                object taskLock3 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                bool lockTaken3 = false;
                try
                {
                    Monitor.Enter(taskLock3, ref lockTaken3);
                    byte[] array = new byte[2]
                    {
                        Axis[i].RealTimeConfigStatusArgs1Low,
                        Axis[i].RealTimeConfigStatusArgs1High
                    };
                    byte[] data = new byte[2]
                    {
                        Axis[i].RealTimeConfigStatusArgs2Low,
                        Axis[i].RealTimeConfigStatusArgs2High
                    };
                    Axis[i].RealTimeConfigID = 0;
                    result.Data = data;
                    result.Status = Axis[i].RealTimeConfigIDStatus;
                }
                finally
                {
                    if (lockTaken3)
                    {
                        Monitor.Exit(taskLock3);
                    }
                }
            }

            return result;
        }
    }

    private CTEntryDataDefinition getCTEntryData(string TargetIP, ushort EntryNumber)
    {
        CTEntryDataDefinition result = default(CTEntryDataDefinition);
        byte[] bytes = BitConverter.GetBytes(EntryNumber);
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].RealTimeConfigID = 134;
                    if (Axis[i].RealTimeConfigCommandCount > 14)
                    {
                        Axis[i].RealTimeConfigCommandCount = 1;
                    }

                    Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[i].RealTimeConfigCommandCount + 1));
                    Axis[i].RealTimeConfigArgs1High = bytes[1];
                    Axis[i].RealTimeConfigArgs1Low = bytes[0];
                    Axis[i].RealTimeConfigArgs2High = 0;
                    Axis[i].RealTimeConfigArgs2Low = 0;
                    Axis[i].RealTimeConfigArgs3High = 0;
                    Axis[i].RealTimeConfigArgs3Low = 0;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                long timeOutTime = getTimeOutTime(5000L);
                while (!((Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount) & ((Axis[i].RealTimeConfigIDStatus == 0) | (Axis[i].RealTimeConfigIDStatus == 4))))
                {
                    if (!isTimeOut(timeOutTime))
                    {
                        continue;
                    }

                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        Axis[i].DLLErrorText = "8: Config Channel Timeout by requesting Command Table Data. Error code: " + Conversion.Hex(Axis[i].RealTimeConfigIDStatus).ToString();
                        ACIError = "Config Channel Error";
                        Axis[i].RealTimeConfigStatusArgs2Low = 9;
                        Axis[i].RealTimeConfigStatusArgs2High = 0;
                        Axis[i].RealTimeConfigStatusArgs3Low = 9;
                        Axis[i].RealTimeConfigStatusArgs3High = 0;
                        Console.WriteLine(Axis[i].DLLErrorText);
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    break;
                }

                object taskLock3 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                bool lockTaken3 = false;
                try
                {
                    Monitor.Enter(taskLock3, ref lockTaken3);
                    byte[] array = new byte[2]
                    {
                        Axis[i].RealTimeConfigStatusArgs1Low,
                        Axis[i].RealTimeConfigStatusArgs1High
                    };
                    byte[] data = new byte[4]
                    {
                        Axis[i].RealTimeConfigStatusArgs2Low,
                        Axis[i].RealTimeConfigStatusArgs2High,
                        Axis[i].RealTimeConfigStatusArgs3Low,
                        Axis[i].RealTimeConfigStatusArgs3High
                    };
                    Axis[i].RealTimeConfigID = 0;
                    result.Data = data;
                    result.Status = Axis[i].RealTimeConfigIDStatus;
                }
                finally
                {
                    if (lockTaken3)
                    {
                        Monitor.Exit(taskLock3);
                    }
                }
            }

            return result;
        }
    }

    private CTEntryDataDefinition WriteCT_Entry(string TargetIP, ushort EntryNumber, ushort BlockSize)
    {
        CTEntryDataDefinition result = default(CTEntryDataDefinition);
        byte[] bytes = BitConverter.GetBytes(EntryNumber);
        byte[] bytes2 = BitConverter.GetBytes(BlockSize);
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].RealTimeConfigID = 131;
                    if (Axis[i].RealTimeConfigCommandCount > 14)
                    {
                        Axis[i].RealTimeConfigCommandCount = 1;
                    }

                    Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[i].RealTimeConfigCommandCount + 1));
                    Axis[i].RealTimeConfigArgs1High = bytes[1];
                    Axis[i].RealTimeConfigArgs1Low = bytes[0];
                    Axis[i].RealTimeConfigArgs2High = bytes2[1];
                    Axis[i].RealTimeConfigArgs2Low = bytes2[0];
                    Axis[i].RealTimeConfigArgs3High = 0;
                    Axis[i].RealTimeConfigArgs3Low = 0;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                long timeOutTime = getTimeOutTime(5000L);
                while (!((Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount) & (Axis[i].RealTimeConfigIDStatus == 0)))
                {
                    if (!isTimeOut(timeOutTime))
                    {
                        continue;
                    }

                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        Axis[i].DLLErrorText = "8: Config Channel Timeout by requesting Command Table Data.. Error code: " + Conversion.Hex(Axis[i].RealTimeConfigIDStatus).ToString();
                        ACIError = "Config Channel Error";
                        Axis[i].RealTimeConfigStatusArgs2Low = 9;
                        Axis[i].RealTimeConfigStatusArgs2High = 0;
                        Axis[i].RealTimeConfigStatusArgs3Low = 9;
                        Axis[i].RealTimeConfigStatusArgs3High = 0;
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    break;
                }

                object taskLock3 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                bool lockTaken3 = false;
                try
                {
                    Monitor.Enter(taskLock3, ref lockTaken3);
                    byte[] array = new byte[2]
                    {
                        Axis[i].RealTimeConfigStatusArgs1Low,
                        Axis[i].RealTimeConfigStatusArgs1High
                    };
                    byte[] array2 = new byte[2]
                    {
                        Axis[i].RealTimeConfigStatusArgs2Low,
                        Axis[i].RealTimeConfigStatusArgs2High
                    };
                    Axis[i].RealTimeConfigID = 0;
                    result.Status = Axis[i].RealTimeConfigIDStatus;
                }
                finally
                {
                    if (lockTaken3)
                    {
                        Monitor.Exit(taskLock3);
                    }
                }
            }

            return result;
        }
    }

    private CTEntryDataDefinition WriteCT_EntryData(string TargetIP, ushort EntryNumber, byte[] Data)
    {
        CTEntryDataDefinition result = default(CTEntryDataDefinition);
        byte[] bytes = BitConverter.GetBytes(EntryNumber);
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].RealTimeConfigID = 132;
                    if (Axis[i].RealTimeConfigCommandCount > 14)
                    {
                        Axis[i].RealTimeConfigCommandCount = 1;
                    }

                    Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[i].RealTimeConfigCommandCount + 1));
                    Axis[i].RealTimeConfigArgs1High = bytes[1];
                    Axis[i].RealTimeConfigArgs1Low = bytes[0];
                    Axis[i].RealTimeConfigArgs2High = Data[1];
                    Axis[i].RealTimeConfigArgs2Low = Data[0];
                    Axis[i].RealTimeConfigArgs3High = Data[3];
                    Axis[i].RealTimeConfigArgs3Low = Data[2];
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                long timeOutTime = getTimeOutTime(5000L);
                while (!((Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount) & ((Axis[i].RealTimeConfigIDStatus == 0) | (Axis[i].RealTimeConfigIDStatus == 4))))
                {
                    if (!isTimeOut(timeOutTime))
                    {
                        continue;
                    }

                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        Axis[i].DLLErrorText = "8: Config Channel Timeout by requesting Command Table Data.. Error code: " + Conversion.Hex(Axis[i].RealTimeConfigIDStatus).ToString();
                        ACIError = "Config Channel Error";
                        Axis[i].RealTimeConfigStatusArgs2Low = 9;
                        Axis[i].RealTimeConfigStatusArgs2High = 0;
                        Axis[i].RealTimeConfigStatusArgs3Low = 9;
                        Axis[i].RealTimeConfigStatusArgs3High = 0;
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    break;
                }

                object taskLock3 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                bool lockTaken3 = false;
                try
                {
                    Monitor.Enter(taskLock3, ref lockTaken3);
                    byte[] array = new byte[2]
                    {
                        Axis[i].RealTimeConfigStatusArgs1Low,
                        Axis[i].RealTimeConfigStatusArgs1High
                    };
                    byte[] array2 = new byte[2]
                    {
                        Axis[i].RealTimeConfigStatusArgs2Low,
                        Axis[i].RealTimeConfigStatusArgs2High
                    };
                    Axis[i].RealTimeConfigID = 0;
                    result.Status = Axis[i].RealTimeConfigIDStatus;
                }
                finally
                {
                    if (lockTaken3)
                    {
                        Monitor.Exit(taskLock3);
                    }
                }
            }

            return result;
        }
    }

    //
    // Summary:
    //     This function loads the Command Table content from the drive into it's data structure
    //
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Return the Command Table content in it's return data structure
    public CommandTableStructure LMcf_getCommandTableContent(string TargetIP)
    {
        CommandTableStructure result = default(CommandTableStructure);
        result.CTPresenceList = new List<byte>();
        result.Entries = new List<CTEntry>();
        byte[] cT_PresenceListList = getCT_PresenceListList(TargetIP, 135);
        byte[] array = cT_PresenceListList;
        foreach (byte item in array)
        {
            result.CTPresenceList.Add(item);
        }

        cT_PresenceListList = getCT_PresenceListList(TargetIP, 136);
        byte[] array2 = cT_PresenceListList;
        foreach (byte item2 in array2)
        {
            result.CTPresenceList.Add(item2);
        }

        cT_PresenceListList = getCT_PresenceListList(TargetIP, 137);
        byte[] array3 = cT_PresenceListList;
        foreach (byte item3 in array3)
        {
            result.CTPresenceList.Add(item3);
        }

        cT_PresenceListList = getCT_PresenceListList(TargetIP, 138);
        byte[] array4 = cT_PresenceListList;
        foreach (byte item4 in array4)
        {
            result.CTPresenceList.Add(item4);
        }

        cT_PresenceListList = getCT_PresenceListList(TargetIP, 139);
        byte[] array5 = cT_PresenceListList;
        foreach (byte item5 in array5)
        {
            result.CTPresenceList.Add(item5);
        }

        cT_PresenceListList = getCT_PresenceListList(TargetIP, 140);
        byte[] array6 = cT_PresenceListList;
        foreach (byte item6 in array6)
        {
            result.CTPresenceList.Add(item6);
        }

        cT_PresenceListList = getCT_PresenceListList(TargetIP, 141);
        byte[] array7 = cT_PresenceListList;
        foreach (byte item7 in array7)
        {
            result.CTPresenceList.Add(item7);
        }

        cT_PresenceListList = getCT_PresenceListList(TargetIP, 142);
        byte[] array8 = cT_PresenceListList;
        foreach (byte item8 in array8)
        {
            result.CTPresenceList.Add(item8);
        }

        checked
        {
            int num3 = result.CTPresenceList.Count - 1;
            for (int num4 = 0; num4 <= num3; num4++)
            {
                byte b = result.CTPresenceList[num4];
                int num5 = 0;
                do
                {
                    if ((unchecked((byte)(b >>> (num5 & 7))) & 1) == 0)
                    {
                        CTEntryDataDefinition cTEntry = getCTEntry(TargetIP, (ushort)(num4 * 8 + num5));
                        if (cTEntry.Status == 0 && BitConverter.ToInt16(cTEntry.Data, 0) == 64)
                        {
                            CTEntry item9 = default(CTEntry);
                            item9.data = new List<byte>();
                            CTEntryDataDefinition cTEntryData = getCTEntryData(TargetIP, (ushort)(num4 * 8 + num5));
                            while (cTEntryData.Status == 4)
                            {
                                byte[] data = cTEntryData.Data;
                                foreach (byte item10 in data)
                                {
                                    item9.data.Add(item10);
                                }

                                cTEntryData = getCTEntryData(TargetIP, (ushort)(num4 * 8 + num5));
                            }

                            result.Entries.Add(item9);
                        }
                    }

                    num5++;
                }
                while (num5 <= 7);
            }

            return result;
        }
    }

    CommandTableStructure _ACI.LMcf_getCommandTableContent(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMcf_getCommandTableContent
        return this.LMcf_getCommandTableContent(TargetIP);
    }

    //
    // Summary:
    //     This function delete the Command Table content in the drive's RAM. If the WriteCommandTabkeToFLASH
    //     command will be called after this command, the Command Table will be empty.
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Return status data (see drive onfiguration over fieldbus manual). OK if 0
    public ushort LMcf_DeleteCommandTable_RAM(string TargetIP)
    {
        ushort result = 0;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].RealTimeConfigID = 129;
                    if (Axis[i].RealTimeConfigCommandCount > 14)
                    {
                        Axis[i].RealTimeConfigCommandCount = 1;
                    }

                    Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[i].RealTimeConfigCommandCount + 1));
                    Axis[i].RealTimeConfigArgs1High = 0;
                    Axis[i].RealTimeConfigArgs1Low = 0;
                    Axis[i].RealTimeConfigArgs2High = 0;
                    Axis[i].RealTimeConfigArgs2Low = 0;
                    Axis[i].RealTimeConfigArgs3High = 0;
                    Axis[i].RealTimeConfigArgs3Low = 0;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                long timeOutTime = getTimeOutTime(5000L);
                while (!((Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount) & (Axis[i].RealTimeConfigIDStatus == 0)))
                {
                    if (!isTimeOut(timeOutTime))
                    {
                        continue;
                    }

                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        Axis[i].DLLErrorText = "8: Config Channel Timeout by requesting Command Table Data.. Error code: " + Conversion.Hex(Axis[i].RealTimeConfigIDStatus).ToString();
                        ACIError = "Config Channel Error";
                        Axis[i].RealTimeConfigStatusArgs2Low = 9;
                        Axis[i].RealTimeConfigStatusArgs2High = 0;
                        Axis[i].RealTimeConfigStatusArgs3Low = 9;
                        Axis[i].RealTimeConfigStatusArgs3High = 0;
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    break;
                }

                object taskLock3 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                bool lockTaken3 = false;
                try
                {
                    Monitor.Enter(taskLock3, ref lockTaken3);
                    byte[] array = new byte[2]
                    {
                        Axis[i].RealTimeConfigStatusArgs1Low,
                        Axis[i].RealTimeConfigStatusArgs1High
                    };
                    byte[] array2 = new byte[2]
                    {
                        Axis[i].RealTimeConfigStatusArgs2Low,
                        Axis[i].RealTimeConfigStatusArgs2High
                    };
                    Axis[i].RealTimeConfigID = 0;
                    result = Axis[i].RealTimeConfigIDStatus;
                }
                finally
                {
                    if (lockTaken3)
                    {
                        Monitor.Exit(taskLock3);
                    }
                }
            }

            return result;
        }
    }

    ushort _ACI.LMcf_DeleteCommandTable_RAM(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMcf_DeleteCommandTable_RAM
        return this.LMcf_DeleteCommandTable_RAM(TargetIP);
    }

    //
    // Summary:
    //     This command will store the COmmand Table from RAM to FLASH ROM
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    // Returns:
    //     Return status data (see drive onfiguration over fieldbus manual). OK if 0
    public ushort LMcf_WriteCommandTableToFLASH(string TargetIP)
    {
        ushort result = 0;
        int targetIPListCount = TargetIPListCount;
        checked
        {
            for (int i = 0; i <= targetIPListCount; i++)
            {
                if (Operators.CompareString(Axis[i].SlaveIP, TargetIP, TextCompare: false) != 0)
                {
                    continue;
                }

                object taskLock = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock);
                bool lockTaken = false;
                try
                {
                    Monitor.Enter(taskLock, ref lockTaken);
                    Axis[i].RealTimeConfigID = 128;
                    if (Axis[i].RealTimeConfigCommandCount > 14)
                    {
                        Axis[i].RealTimeConfigCommandCount = 1;
                    }

                    Axis[i].RealTimeConfigCommandCount = byte.Parse(Conversions.ToString(Axis[i].RealTimeConfigCommandCount + 1));
                    Axis[i].RealTimeConfigArgs1High = 0;
                    Axis[i].RealTimeConfigArgs1Low = 0;
                    Axis[i].RealTimeConfigArgs2High = 0;
                    Axis[i].RealTimeConfigArgs2Low = 0;
                    Axis[i].RealTimeConfigArgs3High = 0;
                    Axis[i].RealTimeConfigArgs3Low = 0;
                    Axis[i].setTimeoutObservation = true;
                    Axis[i].setSkipAmountResponsePackets = 20;
                }
                finally
                {
                    if (lockTaken)
                    {
                        Monitor.Exit(taskLock);
                    }
                }

                long timeOutTime = getTimeOutTime(15000L);
                while (!((Axis[i].RealTimeConfigStatusCommandCount == Axis[i].RealTimeConfigCommandCount) & (Axis[i].RealTimeConfigIDStatus == 0)))
                {
                    if (!isTimeOut(timeOutTime))
                    {
                        continue;
                    }

                    object taskLock2 = TaskLock;
                    ObjectFlowControl.CheckForSyncLockOnValueType(taskLock2);
                    bool lockTaken2 = false;
                    try
                    {
                        Monitor.Enter(taskLock2, ref lockTaken2);
                        Axis[i].DLLErrorText = "8: Config Channel Timeout by requesting Command Table Data.. Error code: " + Conversion.Hex(Axis[i].RealTimeConfigIDStatus).ToString();
                        ACIError = "Config Channel Error";
                        Axis[i].RealTimeConfigStatusArgs2Low = 9;
                        Axis[i].RealTimeConfigStatusArgs2High = 0;
                        Axis[i].RealTimeConfigStatusArgs3Low = 9;
                        Axis[i].RealTimeConfigStatusArgs3High = 0;
                    }
                    finally
                    {
                        if (lockTaken2)
                        {
                            Monitor.Exit(taskLock2);
                        }
                    }

                    break;
                }

                object taskLock3 = TaskLock;
                ObjectFlowControl.CheckForSyncLockOnValueType(taskLock3);
                bool lockTaken3 = false;
                try
                {
                    Monitor.Enter(taskLock3, ref lockTaken3);
                    Axis[i].setTimeoutObservation = false;
                    byte[] array = new byte[2]
                    {
                        Axis[i].RealTimeConfigStatusArgs1Low,
                        Axis[i].RealTimeConfigStatusArgs1High
                    };
                    byte[] array2 = new byte[2]
                    {
                        Axis[i].RealTimeConfigStatusArgs2Low,
                        Axis[i].RealTimeConfigStatusArgs2High
                    };
                    Axis[i].RealTimeConfigID = 0;
                    result = Axis[i].RealTimeConfigIDStatus;
                }
                finally
                {
                    if (lockTaken3)
                    {
                        Monitor.Exit(taskLock3);
                    }
                }
            }

            return result;
        }
    }

    ushort _ACI.LMcf_WriteCommandTableToFLASH(string TargetIP)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMcf_WriteCommandTableToFLASH
        return this.LMcf_WriteCommandTableToFLASH(TargetIP);
    }

    //
    // Summary:
    //     This command will write a Command Table content to the drive RAM
    //
    // Parameters:
    //   TargetIP:
    //     IP address of the requested axis
    //
    //   CT:
    //     Command Table data
    //
    // Returns:
    //     True, if successful
    public bool LMcf_setCommandTableContent(string TargetIP, CommandTableStructure CT)
    {
        int num = 0;
        bool result = true;
        checked
        {
            int num2 = CT.CTPresenceList.Count - 1;
            for (int i = 0; i <= num2; i++)
            {
                byte b = CT.CTPresenceList[i];
                Console.WriteLine(b.ToString());
                int num3 = 0;
                do
                {
                    if ((unchecked((byte)(b >>> (num3 & 7))) & 1) == 0)
                    {
                        if (WriteCT_Entry(TargetIP, (ushort)(i * 8 + num3), 64).Status == 0)
                        {
                            while (CT.Entries[num].data.Count < 64)
                            {
                                CT.Entries[num].data.Add(default(byte));
                            }

                            int num4 = 0;
                            do
                            {
                                CTEntryDataDefinition cTEntryDataDefinition = WriteCT_EntryData(TargetIP, (ushort)(i * 8 + num3), new byte[4]
                                {
                                    CT.Entries[num].data[num4],
                                    CT.Entries[num].data[num4 + 1],
                                    CT.Entries[num].data[num4 + 2],
                                    CT.Entries[num].data[num4 + 3]
                                });
                                if (!((cTEntryDataDefinition.Status == 0) | (cTEntryDataDefinition.Status == 4)))
                                {
                                    result = false;
                                }

                                if (cTEntryDataDefinition.Status == 0)
                                {
                                    break;
                                }

                                num4 += 4;
                            }
                            while (num4 <= 64);
                        }

                        num++;
                    }

                    num3++;
                }
                while (num3 <= 7);
            }

            return result;
        }
    }

    bool _ACI.LMcf_setCommandTableContent(string TargetIP, CommandTableStructure CT)
    {
        //ILSpy generated this explicit interface implementation from .override directive in LMcf_setCommandTableContent
        return this.LMcf_setCommandTableContent(TargetIP, CT);
    }
}
#if false // Decompilation log
'12' items in cache
------------------
Resolve: 'mscorlib, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089'
Found single assembly: 'mscorlib, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089'
Load from: 'C:\Program Files (x86)\Reference Assemblies\Microsoft\Framework\.NETFramework\v4.8\mscorlib.dll'
------------------
Resolve: 'System, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089'
Found single assembly: 'System, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089'
Load from: 'C:\Program Files (x86)\Reference Assemblies\Microsoft\Framework\.NETFramework\v4.8\System.dll'
------------------
Resolve: 'Microsoft.VisualBasic, Version=10.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a'
Found single assembly: 'Microsoft.VisualBasic, Version=10.0.0.0, Culture=neutral, PublicKeyToken=b03f5f7f11d50a3a'
Load from: 'C:\Program Files (x86)\Reference Assemblies\Microsoft\Framework\.NETFramework\v4.8\Microsoft.VisualBasic.dll'
------------------
Resolve: 'System.Xml, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089'
Found single assembly: 'System.Xml, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089'
Load from: 'C:\Program Files (x86)\Reference Assemblies\Microsoft\Framework\.NETFramework\v4.8\System.Xml.dll'
------------------
Resolve: 'System.Core, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089'
Found single assembly: 'System.Core, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089'
Load from: 'C:\Program Files (x86)\Reference Assemblies\Microsoft\Framework\.NETFramework\v4.8\System.Core.dll'
#endif