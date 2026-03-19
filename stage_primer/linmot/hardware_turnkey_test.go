package linmot

import (
	"context"
	"fmt"
	"testing"
	"time"

	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"

	"stage_primer_config"
)

func prepareDriveForVAIMotion(ctx context.Context, t *testing.T, targetIP string) {
	t.Helper()

	linmotClient, err := globalClientFactory.CreateClient(targetIP)
	if err != nil {
		t.Fatalf("CreateClient failed: %v", err)
	}

	fmt.Printf("[PREP] Setting Run Mode = MotionCommandInterface\n")
	if err := linmotClient.SetRunMode(ctx, protocol_common.RunModes.MotionCommandInterface, protocol_common.ParameterStorage.Both); err != nil {
		t.Fatalf("SetRunMode(MotionCommandInterface) failed: %v", err)
	}

	fmt.Printf("[PREP] Clearing trigger/event-handler config on input 4.6\n")
	if err := linmotClient.SetTriggerMode(ctx, protocol_common.TriggerModeConfig.None, protocol_common.ParameterStorage.Both); err != nil {
		t.Fatalf("SetTriggerMode(None) failed: %v", err)
	}
	if err := linmotClient.SetIODefInputFunction(ctx, protocol_common.IOPin.Input46, protocol_common.InputFunction.None, protocol_common.ParameterStorage.Both); err != nil {
		t.Fatalf("SetIODefInputFunction(Input46=None) failed: %v", err)
	}
	if err := linmotClient.SetEasyStepsInputRisingEdgeFunction(ctx, protocol_common.IOPin.Input46, protocol_common.EasyStepsIOMotion.None, protocol_common.ParameterStorage.Both); err != nil {
		t.Fatalf("SetEasyStepsInputRisingEdgeFunction(Input46=None) failed: %v", err)
	}

	status, err := linmotClient.GetStatus(ctx)
	if err != nil {
		t.Fatalf("GetStatus after prep failed: %v", err)
	}
	fmt.Printf("[PREP] Status after prep: StatusWord=0x%04X StateVar=0x%04X ErrorCode=0x%04X\n", status.StatusWord, status.StateVar, status.ErrorCode)

	if (status.StatusWord&0x0008) != 0 || status.ErrorCode != 0 {
		fmt.Printf("[PREP] Fault present, acknowledging...\n")
		if _, err := linmotClient.AcknowledgeError(ctx); err != nil {
			t.Fatalf("AcknowledgeError during prep failed: %v", err)
		}
	}
}

// TestHardwareTurnkeyRecovery provides a rigorous, end-to-end verification of
// the auto-recovery system on real hardware.
//
// IT DOES THE FOLLOWING:
// 1. Connects to the LinMot at 10.8.7.232.
// 2. Starts a background monitor loop (The "Brain").
// 3. Initiates a series of JOG commands (The "Motion").
// 4. Waits for YOU to obstruct the motor to induce a fault.
// 5. Detects the fault and triggers the auto-recovery (The "Action").
// 6. Verifies the drive returns to a Healthy state (The "Proof").
//
// USAGE:
// go test -v ./stage_primer/primer/linmot -run TestHardwareTurnkeyRecovery -linmot_mode=udp -linmot_ip=10.8.7.232 -timeout 10m
func TestHardwareTurnkeyRecovery(t *testing.T) {
	if *linmotMode != "udp" {
		t.Skip("Skipping hardware test; run with -linmot_mode=udp")
	}

	targetIP := *linmotIP

	// --- SETUP ---
	SetClientFactory(&pooledClientFactory{})
	defer ResetClientFactory()
	ResetFaultBudget(targetIP)
	resetFaultLifecycleStateForTests()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Printf("\n========================================================\n")
	fmt.Printf("   LINMOT AUTO-RECOVERY TURNKEY VERIFICATION\n")
	fmt.Printf("========================================================\n")
	fmt.Printf("Target Hardware: %s\n", targetIP)

	// --- 1. START MONITORING ---
	fmt.Printf("[1/4] Preparing drive for direct VAI motion and starting monitor loop...\n")
	setupCtx, setupCancel := context.WithTimeout(ctx, 10*time.Second)
	prepareDriveForVAIMotion(setupCtx, t, targetIP)
	setupCancel()

	monitorErrCh := make(chan error, 1)
	go func() {
		cfg := config.Config{
			ClearCores: []config.ClearCoreConfig{
				{LinMots: []config.LinMotConfig{{IP: targetIP}}},
			},
		}
		if err := MonitorFaults(ctx, cfg); err != nil && ctx.Err() == nil {
			select {
			case monitorErrCh <- err:
			default:
			}
		}
	}()

	// --- 2. PREPARATION WINDOW ---
	fmt.Printf("[2/4] PREPARE TO OBSTRUCT: Motion starts in 5 seconds.\n")
	for i := 5; i > 0; i-- {
		fmt.Printf("      Starting in %d...\n", i)
		time.Sleep(1 * time.Second)
	}

	// --- 3. MOTION GENERATOR ---
	fmt.Printf("[3/4] MOTION STARTED: Motor is toggling between 10mm and 40mm.\n")
	fmt.Printf("      >>> ACTION: PLEASE OBSTRUCT THE MOTOR NOW! <<<\n")
	fmt.Printf("      (Stop it from reaching its target to trigger 'Position Lag')\n\n")

	motionCtx, stopMotion := context.WithCancel(ctx)
	defer stopMotion()
	jogErrCh := make(chan error, 1)

	go func() {
		posA, posB := 10.0, 40.0
		current := posA
		for {
			select {
			case <-motionCtx.Done():
				return
			default:
				fmt.Printf("[MOTION] Command: Move to %.1f mm\n", current)
				if err := Jog(motionCtx, JogConfig{Position: current, Config: config.Config{
					ClearCores: []config.ClearCoreConfig{{LinMots: []config.LinMotConfig{{IP: targetIP}}}},
				}}); err != nil && motionCtx.Err() == nil {
					select {
					case jogErrCh <- err:
					default:
					}
				}

				// Verify motion happened
				time.Sleep(1500 * time.Millisecond)
				pos, _ := GetPosition(motionCtx, PositionConfig{
					Config: config.Config{
						ClearCores: []config.ClearCoreConfig{{LinMots: []config.LinMotConfig{{IP: targetIP}}}},
					},
				})
				fmt.Printf("[MOTION] Verification: Actual Position = %.2f mm\n", pos)

				if current == posA {
					current = posB
				} else {
					current = posA
				}
				time.Sleep(2 * time.Second)
			}
		}
	}()

	// --- 4. VERIFICATION LOOP ---
	start := time.Now()
	timeout := 2 * time.Minute
	faultDetected := false
	recoveryObserved := false

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	deadline := time.After(timeout)
	for {
		select {
		case <-deadline:
			t.Fatalf("Test timed out after %v without observing recovery.", timeout)
		case err := <-monitorErrCh:
			t.Fatalf("MonitorFaults exited unexpectedly: %v", err)
		case err := <-jogErrCh:
			t.Logf("[MOTION] Jog error observed during recovery exercise: %v", err)
		case <-ticker.C:
			snapshot := getFaultLifecycleSnapshot(targetIP)

			// LOG STATE TRANSITIONS
			switch snapshot.State {
			case FaultLifecycleStateRecovering:
				if !faultDetected {
					fmt.Printf("\n[!] EVIDENCE: Fault Detected (Code: 0x%04X). System is now RECOVERING...\n", snapshot.LastFaultCode)
					faultDetected = true
				}
			case FaultLifecycleStateHealthy:
				if faultDetected && !recoveryObserved {
					fmt.Printf("[!] EVIDENCE: Auto-Recovery SUCCESSFUL! Drive is back to Healthy state.\n")
					fmt.Printf("    Recovery took: %v since detection.\n", time.Since(start))
					recoveryObserved = true

					fmt.Printf("\n========================================================\n")
					fmt.Printf("   VERIFICATION COMPLETE: AUTO-RECOVERY IS WORKING!\n")
					fmt.Printf("========================================================\n")
					return
				}
			case FaultLifecycleStateEscalated:
				fmt.Printf("[!] WARNING: Fault was ESCALATED (Budget exceeded or fatal). Check logs.\n")
				t.Errorf("Recovery failed and escalated.")
				return
			}

			// Heartbeat during wait
			if !faultDetected && time.Since(start).Seconds() > 5 {
				// Also try to get raw status to prove we have connectivity and motion is happening
				if linmotClient, err := globalClientFactory.CreateClient(targetIP); err == nil {
					statusCtx, statusCancel := context.WithTimeout(ctx, 2*time.Second)
					if status, statusErr := linmotClient.GetStatus(statusCtx); statusErr == nil {
						fmt.Printf("      ...waiting for obstruction (Primer State: %s | Drive Actual Pos: %.2f mm | Demand: %.2f mm)...\n",
							snapshot.State, status.ActualPositionMM(), status.DemandPositionMM())
					} else {
						fmt.Printf("      ...waiting for obstruction (Primer State: %s | Drive status check failed)...\n", snapshot.State)
					}
					statusCancel()
				} else {
					fmt.Printf("      ...waiting for obstruction (Primer State: %s)...\n", snapshot.State)
				}
			}
		}
	}
}

// TestHardwareMotionProof is a "Pre-Flight" check. It moves the motor twice
// and verifies the position changed, proving that our commands are reaching
// the physical hardware.
func TestHardwareMotionProof(t *testing.T) {
	if *linmotMode != "udp" {
		t.Skip("Skipping hardware test; run with -linmot_mode=udp")
	}

	targetIP := *linmotIP
	SetClientFactory(&pooledClientFactory{})
	defer ResetClientFactory()

	fmt.Printf("\n--- MOTION PROOF-OF-LIFE STARTED ---\n")

	// Prepare drive for direct VAI motion mode (avoid Triggered Command Table mode).
	fmt.Printf("Preparing drive for direct VAI motion...\n")
	setupCtx, setupCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer setupCancel()
	prepareDriveForVAIMotion(setupCtx, t, targetIP)

	moveAndVerify := func(target float64) {
		fmt.Printf("Moving to %.2f mm...\n", target)
		jogCtx, jogCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer jogCancel()
		err := Jog(jogCtx, JogConfig{
			Position: target,
			Config: config.Config{
				ClearCores: []config.ClearCoreConfig{{LinMots: []config.LinMotConfig{{IP: targetIP}}}},
			},
		})
		if err != nil {
			t.Fatalf("Jog failed: %v", err)
		}
		linmotClient, err := globalClientFactory.CreateClient(targetIP)
		if err == nil {
			statusCtx, statusCancel := context.WithTimeout(context.Background(), 3*time.Second)
			if status, statusErr := linmotClient.GetStatus(statusCtx); statusErr == nil {
				fmt.Printf("[MOTION] Post-command status: SW=0x%04X SV=0x%04X Err=0x%04X Act=%.2f Dem=%.2f\n",
					status.StatusWord, status.StateVar, status.ErrorCode, status.ActualPositionMM(), status.DemandPositionMM())
			}
			statusCancel()
		}

		// Wait for motion to settle
		time.Sleep(2 * time.Second)

		positionCtx, positionCancel := context.WithTimeout(context.Background(), 3*time.Second)
		pos, err := GetPosition(positionCtx, PositionConfig{
			Config: config.Config{
				ClearCores: []config.ClearCoreConfig{{LinMots: []config.LinMotConfig{{IP: targetIP}}}},
			},
		})
		positionCancel()
		if err != nil {
			t.Fatalf("GetPosition failed: %v", err)
		}

		// Allow 0.5mm tolerance
		if pos < target-0.5 || pos > target+0.5 {
			t.Fatalf("Position mismatch: Target %.2f, Actual %.2f", target, pos)
		}
		fmt.Printf("SUCCESS: Reached %.2f mm (Actual: %.2f)\n", target, pos)
	}

	moveAndVerify(10.0)
	moveAndVerify(40.0)
	fmt.Printf("--- MOTION PROOF-OF-LIFE PASSED ---\n\n")
}
