package client

import (
	"context"
	"encoding/binary"
	"testing"
	"time"

	client_common "github.com/Smart-Vision-Works/staged_robot/client/common"
	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
	protocol_rtc "github.com/Smart-Vision-Works/staged_robot/protocol/rtc"
)

func TestDiagnostic_MotionAcceptanceMatrix(t *testing.T) {
	if *linmotMode != "udp" {
		t.Skip("Skipping diagnostic in mock mode - requires real hardware")
	}

	client, err := newUDPClientNoValidation(*linmotIP, *linmotDrivePort, *linmotMasterPort, *linmotTimeout, *linmotDebug)
	if err != nil {
		t.Fatalf("Failed to create UDP client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	readUPID := func(upid uint16) (int32, error) {
		req := protocol_rtc.NewRTCGetParamRequest(upid, protocol_rtc.CommandCode.ReadRAM)
		resp, err := client_common.SendRequestAndReceive[*protocol_rtc.RTCGetParamResponse](client.requestManager, ctx, req)
		if err != nil {
			return 0, err
		}
		return resp.Value(), nil
	}

	logStatus := func(label string) *protocol_common.Status {
		st, err := client.GetStatus(ctx)
		if err != nil {
			t.Fatalf("%s: GetStatus failed: %v", label, err)
		}
		mainState := (st.StateVar >> 8) & 0xFF
		countNibble := st.StateVar & 0x000F
		eventActive := (st.StateVar & 0x0010) != 0
		motionActive := (st.StateVar & 0x0020) != 0
		inTarget := (st.StateVar & 0x0040) != 0
		t.Logf("%s: SW=0x%04X SV=0x%04X Main=%d Cnt=%d Event=%v Motion=%v InTarget=%v Err=0x%04X Act=%.3f Dem=%.3f",
			label, st.StatusWord, st.StateVar, mainState, countNibble, eventActive, motionActive, inTarget, st.ErrorCode, st.ActualPositionMM(), st.DemandPositionMM())
		return st
	}

	sendRawVAIGoToPosFromActual := func(counter uint8, targetMM float64) {
		// ReqDef bit1 (MotionControl), RepDef bits 0/1/3 (StatusWord/StateVar/DemandPosition).
		packet := make([]byte, 40)
		binary.LittleEndian.PutUint32(packet[0:4], protocol_common.RequestFlags.MotionControl)
		binary.LittleEndian.PutUint32(packet[4:8], protocol_common.RespBitStatusWord|protocol_common.RespBitStateVar|protocol_common.RespBitDemandPosition)

		// Header: master=0x01(VAI), sub=0x03(Go To Pos From Act Pos And Act Vel), counter=low nibble.
		packet[8] = (0x03 << 4) | (counter & 0x0F)
		packet[9] = 0x01

		// Parameters (int32/uint32 LE): target, velocity, accel, decel.
		// Cast the signed position units to uint32 so the wire packet preserves the raw two's-complement bit pattern.
		binary.LittleEndian.PutUint32(packet[10:14], uint32(protocol_common.ToProtocolPosition(targetMM)))
		binary.LittleEndian.PutUint32(packet[14:18], uint32(protocol_common.ToProtocolVelocity(0.2)))
		binary.LittleEndian.PutUint32(packet[18:22], uint32(protocol_common.ToProtocolAcceleration(2.0)))
		binary.LittleEndian.PutUint32(packet[22:26], uint32(protocol_common.ToProtocolAcceleration(2.0)))

		sendCtx, cancel := context.WithTimeout(ctx, 250*time.Millisecond)
		err := client.requestManager.TransportClient().SendPacket(sendCtx, packet)
		cancel()
		if err != nil {
			t.Fatalf("raw MC send failed (counter=%d): %v", counter, err)
		}
	}

	runMode, err := readUPID(uint16(protocol_common.PUID.RunMode))
	if err != nil {
		t.Logf("Read UPID 0x1450 (RunMode) failed: %v", err)
	} else {
		t.Logf("Read UPID 0x1450 (RunMode): 0x%04X (%d)", uint16(runMode), runMode)
		t.Cleanup(func() {
			restoreCtx, restoreCancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer restoreCancel()

			if err := client.SetRunMode(restoreCtx, protocol_common.RunMode(runMode), protocol_common.ParameterStorage.RAM); err != nil {
				t.Logf("Restore RunMode to 0x%04X failed: %v", uint16(runMode), err)
				return
			}
			t.Logf("Restored RunMode RAM to 0x%04X", uint16(runMode))
		})
	}

	triggerMode, err := readUPID(uint16(protocol_common.PUID.TriggerMode))
	if err != nil {
		t.Logf("Read UPID 0x170C (TriggerMode) failed: %v", err)
	} else {
		t.Logf("Read UPID 0x170C (TriggerMode): 0x%04X (%d)", uint16(triggerMode), triggerMode)
	}

	st0 := logStatus("Baseline")
	if st0.ErrorCode != 0 || (st0.StatusWord&0x0008) != 0 {
		t.Logf("Fault present; trying acknowledge")
		if _, err := client.AcknowledgeError(ctx); err != nil {
			t.Logf("AcknowledgeError failed: %v", err)
		}
	}
	if _, err := client.EnableDrive(ctx); err != nil {
		t.Logf("EnableDrive failed: %v", err)
	}

	modeSet := []protocol_common.RunMode{
		protocol_common.RunModes.MotionCommandInterface,
		protocol_common.RunModes.PCMotionCommandInterface,
	}
	for _, mode := range modeSet {
		t.Logf("Setting RunMode RAM to 0x%04X", uint16(mode))
		if err := client.SetRunMode(ctx, mode, protocol_common.ParameterStorage.RAM); err != nil {
			t.Logf("SetRunMode(0x%04X) failed: %v", uint16(mode), err)
			continue
		}

		time.Sleep(150 * time.Millisecond)
		base := logStatus("After SetRunMode")
		target := base.ActualPositionMM() + 10.0

		sendRawVAIGoToPosFromActual(1, target)
		time.Sleep(250 * time.Millisecond)
		a := logStatus("After raw MC counter=1")

		sendRawVAIGoToPosFromActual(2, target+5.0)
		time.Sleep(250 * time.Millisecond)
		b := logStatus("After raw MC counter=2")

		t.Logf("Mode 0x%04X summary: demand_delta_1=%.3f demand_delta_2=%.3f count_nibble=%d->%d->%d",
			uint16(mode),
			a.DemandPositionMM()-base.DemandPositionMM(),
			b.DemandPositionMM()-a.DemandPositionMM(),
			base.StateVar&0x000F, a.StateVar&0x000F, b.StateVar&0x000F)
	}
}
