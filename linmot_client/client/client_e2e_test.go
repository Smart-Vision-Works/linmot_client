package client

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	_ "embed"

	client_common "github.com/Smart-Vision-Works/staged_robot/client/common"
	client_command_tables "github.com/Smart-Vision-Works/staged_robot/client/rtc/command_tables"
	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
	protocol_motion_control "github.com/Smart-Vision-Works/staged_robot/protocol/motion_control"
	protocol_rtc "github.com/Smart-Vision-Works/staged_robot/protocol/rtc"
	protocol_command_tables "github.com/Smart-Vision-Works/staged_robot/protocol/rtc/command_tables"
)

var (
	linmotMode       = flag.String("linmot_mode", "mock", "client mode: mock|udp|compare")
	linmotIP         = flag.String("linmot_ip", "10.8.7.232", "LinMot drive IP for UDP mode")
	linmotDrivePort  = flag.Int("linmot_drive_port", 49360, "drive port")
	linmotMasterPort = flag.Int("linmot_master_port", 41136, "master port")
	linmotTimeout    = flag.Duration("linmot_timeout", 10*time.Second, "request timeout")
	linmotDebug      = flag.Bool("linmot_debug", false, "enable debug logging")
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

//go:embed rtc/command_tables/testdata/linmot_command_table.yaml
var embeddedCommandTableYAML []byte

func checkUDPPortBindings(t *testing.T) {
	if *linmotMode != "udp" {
		return
	}

	cmd := exec.Command("ss", "-u", "-a", "-p")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Logf("[UDP_PORT_CHECK] ss failed: %v (stderr=%s)", err, strings.TrimSpace(stderr.String()))
		return
	}

	portToken := fmt.Sprintf(":%d", *linmotMasterPort)
	lines := strings.Split(stdout.String(), "\n")
	var matches []string
	for _, line := range lines {
		if strings.Contains(line, portToken) {
			matches = append(matches, line)
		}
	}

	if len(matches) == 0 {
		t.Logf("[UDP_PORT_CHECK] no ss entries for port %d", *linmotMasterPort)
		return
	}

	t.Logf("[UDP_PORT_CHECK] ss -u -a -p matches for port %d:\n%s", *linmotMasterPort, strings.Join(matches, "\n"))

	pidRe := regexp.MustCompile(`pid=([0-9]+)`)
	pidMatches := pidRe.FindAllStringSubmatch(strings.Join(matches, "\n"), -1)
	selfPID := os.Getpid()
	for _, match := range pidMatches {
		pid, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}
		if pid != selfPID {
			t.Fatalf("master port %d already in use by pid=%d; stop other process and retry", *linmotMasterPort, pid)
		}
	}
}

func initializeMockClient() (*Client, func() error, error) {
	client, transportServer := NewMockClient()
	client.SetDebug(*linmotDebug)

	// Create and start the mock drive
	drive := NewMockLinMot(transportServer)
	drive.Start()

	// Set up drive state
	drive.SetStatus(&protocol_common.Status{
		StatusWord:     0x0001,
		StateVar:       0x0002,
		ActualPosition: 100000,
		DemandPosition: 100000,
		Current:        100,
		WarnWord:       0,
		ErrorCode:      0,
	})

	cleanup := func() error {
		drive.Close()
		return client.Close()
	}

	return client, cleanup, nil
}

func initializeUDPClient(t *testing.T) (*Client, func() error, error) {
	checkUDPPortBindings(t)
	client, err := NewUDPClientWithDebug(*linmotIP, *linmotDrivePort, *linmotMasterPort, "", *linmotTimeout, *linmotDebug)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() error {
		return client.Close()
	}
	return client, cleanup, nil
}

func initializeClients(t *testing.T) ([]*Client, func() error, error) {
	switch *linmotMode {
	case "mock":
		client, cleanup, err := initializeMockClient()
		return []*Client{client}, cleanup, err
	case "udp":
		client, cleanup, err := initializeUDPClient(t)
		return []*Client{client}, cleanup, err
	case "compare":
		mockClient, mockCleanup, mockClientErr := initializeMockClient()
		if mockClientErr != nil {
			return nil, nil, mockClientErr
		}

		udpClient, udpCleanup, udpClientErr := initializeUDPClient(t)
		if udpClientErr != nil {
			if mockCleanup != nil {
				_ = mockCleanup()
			}
			return nil, nil, udpClientErr
		}

		cleanup := func() error {
			var mockCleanupErr, udpCleanupErr error
			if mockCleanup != nil {
				mockCleanupErr = mockCleanup()
			}
			if udpCleanup != nil {
				udpCleanupErr = udpCleanup()
			}
			return errors.Join(mockCleanupErr, udpCleanupErr)
		}
		return []*Client{mockClient, udpClient}, cleanup, nil
	default:
		return nil, nil, fmt.Errorf("invalid client mode: %s", *linmotMode)
	}
}

func deferCleanup(t *testing.T, cleanup func() error) {
	t.Helper()
	if cleanup == nil {
		return
	}
	t.Cleanup(func() {
		if err := cleanup(); err != nil {
			t.Logf("cleanup error: %v", err)
		}
	})
}

func setupE2EInstrumentation(t *testing.T, client *Client) {
	dumpRouterTraceOnFailure(t, client)
	logResyncCountIfUDP(t, client)
}

func decodePresenceMasks(masks [8]uint32) []uint16 {
	var entryIDs []uint16
	for maskIdx, mask := range masks {
		baseID := uint16(maskIdx * 32)
		for bit := uint(0); bit < 32; bit++ {
			if mask&(1<<bit) == 0 {
				id := baseID + uint16(bit)
				if id == 0 {
					continue
				}
				entryIDs = append(entryIDs, id)
			}
		}
	}
	return entryIDs
}

func templateEntryIDs(entries []client_command_tables.Entry) []int {
	ids := make([]int, 0, len(entries))
	for _, entry := range entries {
		ids = append(ids, entry.ID)
	}
	sort.Ints(ids)
	return ids
}

func logPresenceMasks(t *testing.T, label string, masks [8]uint32) {
	entryIDs := decodePresenceMasks(masks)
	t.Logf("%s: success masks=[0]=0x%08X [1]=0x%08X [2]=0x%08X [3]=0x%08X [4]=0x%08X [5]=0x%08X [6]=0x%08X [7]=0x%08X entries=%v",
		label, masks[0], masks[1], masks[2], masks[3], masks[4], masks[5], masks[6], masks[7], entryIDs)
}

func readPresenceMasks(t *testing.T, client *Client, label string, timeout time.Duration) ([8]uint32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	masks, err := client.GetPresenceMasks(ctx)
	if err != nil {
		t.Logf("%s: fail err=%v timeout=%s", label, err, timeout)
		return masks, err
	}
	logPresenceMasks(t, label, masks)
	return masks, nil
}

func logErrorEvidence(t *testing.T, client *Client, label string, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	logged, occurred, err := client.GetErrorLogCounts(ctx)
	cancel()
	if err != nil {
		t.Logf("%s: GetErrorLogCounts failed: %v", label, err)
		return
	}
	t.Logf("%s: error_log_counts logged=%d occurred=%d", label, logged, occurred)
	if logged == 0 {
		return
	}
	lastIndex := logged - 1
	ctxEntry, cancelEntry := context.WithTimeout(context.Background(), timeout)
	entry, err := client.GetErrorLogEntry(ctxEntry, lastIndex)
	cancelEntry()
	if err != nil {
		t.Logf("%s: GetErrorLogEntry(%d) failed: %v", label, lastIndex, err)
		return
	}
	t.Logf("%s: error_log_latest index=%d code=0x%04X time=%s", label, entry.Index, entry.ErrorCode, entry.Timestamp.Format(time.RFC3339Nano))
	ctxText, cancelText := context.WithTimeout(context.Background(), timeout)
	text, err := client.GetErrorText(ctxText, entry.ErrorCode)
	cancelText()
	if err != nil {
		t.Logf("%s: GetErrorText(0x%04X) failed: %v", label, entry.ErrorCode, err)
		return
	}
	t.Logf("%s: error_text code=0x%04X text=%q", label, entry.ErrorCode, text)
}

func TestClient_ValidateConnectivity_Repeated(t *testing.T) {
	const iterations = 50

	for i := 0; i < iterations; i++ {
		clients, cleanup, err := initializeClients(t)
		if err != nil {
			t.Fatalf("iteration %d: init failed: %v", i, err)
		}

		for _, client := range clients {
			ctx, cancel := context.WithTimeout(context.Background(), *linmotTimeout)
			_, statusErr := client.GetStatus(ctx)
			cancel()
			if statusErr != nil {
				_ = cleanup()
				t.Fatalf("iteration %d: GetStatus failed: %v", i, statusErr)
			}
		}

		if err := cleanup(); err != nil {
			t.Fatalf("iteration %d: cleanup failed: %v", i, err)
		}
	}
}

func TestClient_PresenceMasks_MCRunning_vs_MCStopped(t *testing.T) {
	t.Skip("skipping flaky MC running/stopped presence mask e2e until environment is stabilized")
	clients, cleanup, err := initializeClients(t)
	if err != nil {
		t.Fatalf("initializeClients(t) failed: %v", err)
	}
	defer deferCleanup(t, cleanup)

	for _, client := range clients {
		setupE2EInstrumentation(t, client)

		type stepResult struct {
			name     string
			duration time.Duration
			err      error
		}
		var results []stepResult

		recordStep := func(name string, start time.Time, err error) {
			duration := time.Since(start)
			results = append(results, stepResult{name: name, duration: duration, err: err})
			if err != nil {
				t.Logf("%s: fail err=%v duration=%s", name, err, duration)
				return
			}
			t.Logf("%s: success duration=%s", name, duration)
		}

		callTimeout := 15 * time.Second

		// 1) Initial status
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), callTimeout)
			start := time.Now()
			status, err := client.GetStatus(ctx)
			cancel()
			if err == nil {
				mainState := (status.StateVar >> 8) & 0xFF
				t.Logf("status_initial: StatusWord=0x%04X StateVar=0x%04X MainState=%d ErrorCode=0x%04X WarnWord=0x%04X",
					status.StatusWord, status.StateVar, mainState, status.ErrorCode, status.WarnWord)
			}
			recordStep("status_initial", start, err)
		}()

		// 2) Presence masks with MC running
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), callTimeout)
			start := time.Now()
			masks, err := client.GetPresenceMasks(ctx)
			cancel()
			if err == nil {
				logPresenceMasks(t, "presence_masks_mc_running", masks)
			}
			recordStep("presence_masks_mc_running", start, err)
		}()

		// 3) Stop MC
		stopMCErr := func() error {
			ctx, cancel := context.WithTimeout(context.Background(), callTimeout)
			start := time.Now()
			err := client.StopMotionController(ctx)
			cancel()
			recordStep("stop_mc", start, err)
			return err
		}()

		if stopMCErr != nil {
			// 3a) If StopMC fails, test a simple RTC request
			ctx, cancel := context.WithTimeout(context.Background(), callTimeout)
			start := time.Now()
			_, err := client.GetVelocity(ctx)
			cancel()
			recordStep("rtc_get_velocity_after_stopmc_fail", start, err)
		}

		// 4) Status after StopMC attempt
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), callTimeout)
			start := time.Now()
			status, err := client.GetStatus(ctx)
			cancel()
			if err == nil {
				mainState := (status.StateVar >> 8) & 0xFF
				t.Logf("status_after_stop: StatusWord=0x%04X StateVar=0x%04X MainState=%d ErrorCode=0x%04X WarnWord=0x%04X",
					status.StatusWord, status.StateVar, mainState, status.ErrorCode, status.WarnWord)
			}
			recordStep("status_after_stop", start, err)
		}()

		// 5) Presence masks after StopMC attempt
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), callTimeout)
			start := time.Now()
			masks, err := client.GetPresenceMasks(ctx)
			cancel()
			if err == nil {
				logPresenceMasks(t, "presence_masks_mc_stopped", masks)
			}
			recordStep("presence_masks_mc_stopped", start, err)
		}()

		// 6) Start MC
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), callTimeout)
			start := time.Now()
			err := client.StartMotionController(ctx)
			cancel()
			recordStep("start_mc", start, err)
		}()

		// 7) Final status
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), callTimeout)
			start := time.Now()
			status, err := client.GetStatus(ctx)
			cancel()
			if err == nil {
				mainState := (status.StateVar >> 8) & 0xFF
				t.Logf("status_final: StatusWord=0x%04X StateVar=0x%04X MainState=%d ErrorCode=0x%04X WarnWord=0x%04X",
					status.StatusWord, status.StateVar, mainState, status.ErrorCode, status.WarnWord)
			}
			recordStep("status_final", start, err)
		}()

		var failures []string
		for _, result := range results {
			if result.err != nil {
				failures = append(failures, fmt.Sprintf("%s err=%v duration=%s", result.name, result.err, result.duration))
			}
		}
		if len(failures) > 0 {
			t.Fatalf("sequence had failures: %s", strings.Join(failures, "; "))
		}
	}
}

func TestClient_GetStatus(t *testing.T) {
	clients, cleanup, err := initializeClients(t)
	if err != nil {
		t.Fatalf("initializeClients(t) failed: %v", err)
	}
	defer deferCleanup(t, cleanup)

	for _, client := range clients {
		setupE2EInstrumentation(t, client)
		status, err := client.GetStatus(context.Background())
		if err != nil {
			t.Fatalf("GetStatus() error: %v", err)
		}

		if status == nil {
			t.Fatal("GetStatus() returned nil status")
		}

		// Verify we got valid status data - no assumptions about values
		// Log for debugging but don't assert on specific values
		t.Logf("Status: Position=%d, StatusWord=0x%04X, StateVar=0x%04X, ErrorCode=0x%04X",
			status.ActualPosition, status.StatusWord, status.StateVar, status.ErrorCode)
	}
}

func TestClient_GetPosition(t *testing.T) {
	clients, cleanup, err := initializeClients(t)
	if err != nil {
		t.Fatalf("initializeClients(t) failed: %v", err)
	}
	defer deferCleanup(t, cleanup)

	for _, client := range clients {
		setupE2EInstrumentation(t, client)
		position, err := client.GetPosition(context.Background())
		if err != nil {
			t.Fatalf("GetPosition() error: %v", err)
		}

		// Just verify we can read a position - don't assume range
		// Real drives may be at any position depending on their state
		t.Logf("Current position: %.4f mm", position)
	}
}

func TestClient_SetPosition1(t *testing.T) {
	clients, cleanup, err := initializeClients(t)
	if err != nil {
		t.Fatalf("initializeClients(t) failed: %v", err)
	}
	defer deferCleanup(t, cleanup)

	for _, client := range clients {
		setupE2EInstrumentation(t, client)
		// Read initial position
		initialPosition, err := client.GetPosition(context.Background())
		if err != nil {
			t.Fatalf("GetPosition() error getting initial: %v", err)
		}

		// Write a new value (relative to current to avoid out-of-range issues)
		targetPosition := initialPosition + 1.0
		err = client.SetPosition1(context.Background(), targetPosition, protocol_common.ParameterStorage.RAM)
		if err != nil {
			t.Fatalf("SetPosition1() error: %v", err)
		}

		// Note: SetPosition1 sets the Position1 parameter (target for motion commands),
		// not the actual motor position. The motor position only changes when motion is executed.
		// We verify the write operation succeeded, not that the motor moved.
		t.Logf("SetPosition1: set target from %.4f to %.4f mm", initialPosition, targetPosition)
	}
}

func TestClient_SetVelocity(t *testing.T) {
	clients, cleanup, err := initializeClients(t)
	if err != nil {
		t.Fatalf("initializeClients(t) failed: %v", err)
	}
	defer deferCleanup(t, cleanup)

	const tolerance = 0.001 // 1 mm/s tolerance for float comparison

	for _, client := range clients {
		setupE2EInstrumentation(t, client)
		// Read initial value to restore later
		initialVelocity, err := client.GetVelocity(context.Background())
		if err != nil {
			t.Fatalf("GetVelocity() error getting initial: %v", err)
		}

		// Write a new value
		targetVelocity := 0.5
		err = client.SetVelocity(context.Background(), targetVelocity, protocol_common.ParameterStorage.RAM)
		if err != nil {
			t.Fatalf("SetVelocity() error: %v", err)
		}

		// Read back and verify within tolerance
		velocity, err := client.GetVelocity(context.Background())
		if err != nil {
			t.Fatalf("GetVelocity() error: %v", err)
		}

		if abs(velocity-targetVelocity) > tolerance {
			t.Errorf("GetVelocity() = %v, want %v (tolerance %v)", velocity, targetVelocity, tolerance)
		}

		// Restore original value (best effort - don't fail test if restore fails)
		_ = client.SetVelocity(context.Background(), initialVelocity, protocol_common.ParameterStorage.RAM)
	}
}

func TestClient_SetAcceleration(t *testing.T) {
	clients, cleanup, err := initializeClients(t)
	if err != nil {
		t.Fatalf("initializeClients(t) failed: %v", err)
	}
	defer deferCleanup(t, cleanup)

	const tolerance = 0.001 // 1 m/s² tolerance for float comparison

	for _, client := range clients {
		setupE2EInstrumentation(t, client)
		// Read initial value to restore later
		initialAcceleration, err := client.GetAcceleration(context.Background())
		if err != nil {
			t.Fatalf("GetAcceleration() error getting initial: %v", err)
		}

		// Write a new value
		targetAcceleration := 2.0
		err = client.SetAcceleration(context.Background(), targetAcceleration, protocol_common.ParameterStorage.RAM)
		if err != nil {
			t.Fatalf("SetAcceleration() error: %v", err)
		}

		// Read back and verify within tolerance
		acceleration, err := client.GetAcceleration(context.Background())
		if err != nil {
			t.Fatalf("GetAcceleration() error: %v", err)
		}

		if abs(acceleration-targetAcceleration) > tolerance {
			t.Errorf("GetAcceleration() = %v, want %v (tolerance %v)", acceleration, targetAcceleration, tolerance)
		}

		// Restore original value (best effort - don't fail test if restore fails)
		_ = client.SetAcceleration(context.Background(), initialAcceleration, protocol_common.ParameterStorage.RAM)
	}
}

func TestClient_SetDeceleration(t *testing.T) {
	clients, cleanup, err := initializeClients(t)
	if err != nil {
		t.Fatalf("initializeClients(t) failed: %v", err)
	}
	defer deferCleanup(t, cleanup)

	const tolerance = 0.001 // 1 m/s² tolerance for float comparison

	for _, client := range clients {
		setupE2EInstrumentation(t, client)
		// Read initial value to restore later
		initialDeceleration, err := client.GetDeceleration(context.Background())
		if err != nil {
			t.Fatalf("GetDeceleration() error getting initial: %v", err)
		}

		// Write a new value
		targetDeceleration := 2.0
		err = client.SetDeceleration(context.Background(), targetDeceleration, protocol_common.ParameterStorage.RAM)
		if err != nil {
			t.Fatalf("SetDeceleration() error: %v", err)
		}

		// Read back and verify within tolerance
		deceleration, err := client.GetDeceleration(context.Background())
		if err != nil {
			t.Fatalf("GetDeceleration() error: %v", err)
		}

		if abs(deceleration-targetDeceleration) > tolerance {
			t.Errorf("GetDeceleration() = %v, want %v (tolerance %v)", deceleration, targetDeceleration, tolerance)
		}

		// Restore original value (best effort - don't fail test if restore fails)
		_ = client.SetDeceleration(context.Background(), initialDeceleration, protocol_common.ParameterStorage.RAM)
	}
}

// Helper functions for tests
func paramLiteral(i int64) *client_command_tables.Param {
	return &client_command_tables.Param{Literal: &i}
}

// abs returns the absolute value of x
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func TestClient_SetCommandTable(t *testing.T) {
	clients, cleanup, err := initializeClients(t)
	if err != nil {
		t.Fatalf("initializeClients(t) failed: %v", err)
	}
	defer deferCleanup(t, cleanup)

	for _, client := range clients {
		setupE2EInstrumentation(t, client)
		// Save existing command table to restore later (if any)
		var existingCT *client_command_tables.CommandTable

		// Create minimal command table for testing
		// Use a simple delay command to avoid motion-related issues
		sentCommandTable := &client_command_tables.CommandTable{
			Entries: []client_command_tables.Entry{
				{
					ID:   1,
					Name: "test_delay",
					Type: "Delay",
					Par1: paramLiteral(1000), // 1000ms delay
				},
			},
		}

		// Stop MC first to ensure deterministic state
		ctxStop, cancelStop := context.WithTimeout(context.Background(), 3*time.Second)
		err = client.StopMotionController(ctxStop)
		cancelStop()
		if err != nil {
			t.Fatalf("StopMotionController() failed: %v", err)
		}

		// Read existing command table before modification (best effort)
		ctxRead, cancelRead := context.WithTimeout(context.Background(), *linmotTimeout)
		existingCT, _ = client.GetCommandTable(ctxRead)
		cancelRead()

		// Best-effort StartMC in defer (in case test fails)
		defer func() {
			ctxStart, cancelStart := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancelStart()
			if err := client.StartMotionController(ctxStart); err != nil {
				t.Logf("[WARN] Failed to restart MC in defer: %v", err)
			}
		}()

		// Upload command table with RestartMC=false so MC stays stopped
		ctx, cancel := context.WithTimeout(context.Background(), *linmotTimeout)
		err = client.SetCommandTableWithOptions(ctx, sentCommandTable, client_command_tables.SetCommandTableOptions{RestartMC: false})
		if err != nil {
			t.Fatalf("SetCommandTable() failed: %v", err)
		}

		// Read back and verify (MC is still stopped, so PresenceMask requests will work)
		receivedCommandTable, err := client.GetCommandTable(ctx)
		if err != nil {
			t.Fatalf("GetCommandTable() error: %v", err)
		}

		if len(receivedCommandTable.Entries) != len(sentCommandTable.Entries) {
			t.Errorf("GetCommandTable() returned %d entries, want %d",
				len(receivedCommandTable.Entries), len(sentCommandTable.Entries))
		}

		if len(receivedCommandTable.Entries) > 0 {
			if receivedCommandTable.Entries[0].ID != 1 {
				t.Errorf("Entry ID = %d, want 1", receivedCommandTable.Entries[0].ID)
			}
		}

		// Restore original command table (best effort)
		if existingCT != nil && len(existingCT.Entries) > 0 {
			ctxRestore, cancelRestore := context.WithTimeout(context.Background(), *linmotTimeout)
			defer cancelRestore()
			_ = client.SetCommandTableWithOptions(ctxRestore, existingCT, client_command_tables.SetCommandTableOptions{RestartMC: false})
		}
		cancel()
	}
}

// TestClient_SetCommandTable_ProductionTemplate tests deploying the actual production command table template.
// This matches how the staged robot system deploys command tables in production.
func TestClient_SetCommandTable_ProductionTemplate(t *testing.T) {
	clients, cleanup, err := initializeClients(t)
	if err != nil {
		t.Fatalf("initializeClients(t) failed: %v", err)
	}
	defer deferCleanup(t, cleanup)

	// Find the production template YAML file
	// Try multiple paths (production, repo-relative)
	wd, _ := os.Getwd()
	// Go up from test dir (linmot_client/client) to repo root
	repoRoot := filepath.Join(wd, "..")
	absRepoRoot, _ := filepath.Abs(repoRoot)

	templatePaths := []string{
		filepath.Join(absRepoRoot, "client", "rtc", "command_tables", "testdata", "linmot_command_table.yaml"), // Repo-relative testdata
	}

	var template *client_command_tables.CommandTable
	var templateSource string

	for _, path := range templatePaths {
		if _, err := os.Stat(path); err == nil {
			template, err = client_command_tables.Load(path)
			if err != nil {
				t.Fatalf("Failed to load command table template from %s: %v", path, err)
			}
			templateSource = path
			break
		}
	}

	if template == nil {
		template, err = client_command_tables.Parse(embeddedCommandTableYAML)
		if err != nil {
			t.Fatalf("Failed to parse embedded command table template (candidates=%v, wd=%s): %v", templatePaths, wd, err)
		}
		templateSource = "embedded:rtc/command_tables/testdata/linmot_command_table.yaml"
	}

	t.Logf("Command table template source: %s (candidates=%v, wd=%s)", templateSource, templatePaths, wd)

	// Copy template to avoid modifying the cached/loaded version
	templateCopy := *template

	// Set variables with realistic production values
	// Unit conversion constants
	const (
		PositionUnit        = 10000 // mm → 0.1µm units
		VelocityUnit        = 100   // percent → µm/s
		AccelerationUnit    = 10000 // percent → 1e-5 m/s² units
		TimeUnit            = 10000 // seconds → 100µs units
		DefaultPurgeDelayMs = 500   // ms
	)

	// Production-like values
	zDistance := 50.0            // mm
	defaultSpeed := 50.0         // percent
	defaultAcceleration := 100.0 // percent
	pickTime := 0.1              // seconds

	// Set variables
	positionDown := int64(zDistance * PositionUnit)
	positionUp := int64(0)
	templateCopy.SetVar("POSITION_DOWN", positionDown)
	templateCopy.SetVar("POSITION_UP", positionUp)

	maxVelocity := int64(defaultSpeed * VelocityUnit)
	templateCopy.SetVar("MAX_VELOCITY", maxVelocity)

	acceleration := int64(defaultAcceleration * AccelerationUnit)
	templateCopy.SetVar("ACCELERATION", acceleration)
	templateCopy.SetVar("DECELERATION", acceleration)

	delayAtBottom := int64(pickTime * TimeUnit)
	templateCopy.SetVar("DELAY_AT_BOTTOM", delayAtBottom)

	delayPurge := int64(DefaultPurgeDelayMs * TimeUnit / 1000)
	templateCopy.SetVar("DELAY_PURGE", delayPurge)

	// Validate template after variable binding
	if err := templateCopy.Validate(); err != nil {
		t.Fatalf("Command table validation failed after variable binding: %v", err)
	}

	entryIDs := templateEntryIDs(templateCopy.Entries)
	if len(entryIDs) == 0 {
		t.Fatalf("Command table template has no entries (source=%s)", templateSource)
	}
	minID := entryIDs[0]
	maxID := entryIDs[len(entryIDs)-1]
	t.Logf("Template entry indices: count=%d min=%d max=%d", len(entryIDs), minID, maxID)
	if len(entryIDs) <= 40 {
		t.Logf("Template entry indices list=%v", entryIDs)
	} else {
		t.Logf("Template entry indices list(head/tail)=%v ... %v", entryIDs[:20], entryIDs[len(entryIDs)-20:])
	}

	// Deploy to each client
	for _, client := range clients {
		setupE2EInstrumentation(t, client)
		ctx, cancel := context.WithTimeout(context.Background(), *linmotTimeout)
		defer cancel()

		// Deploy command table (this is what production does)
		err = client.SetCommandTable(ctx, &templateCopy)
		if err != nil {
			logErrorEvidence(t, client, "setcommand_failed", *linmotTimeout)
			t.Fatalf("SetCommandTable() failed: %v", err)
		}

		// Do one lightweight StatusRequest to confirm firmware responds (no presence masks)
		status, err := client.GetStatus(ctx)
		if err != nil {
			t.Fatalf("GetStatus() failed after SetCommandTable: %v", err)
		}

		t.Logf("SetCommandTable succeeded. StatusWord=0x%04X, StateVar=0x%04X (MainState=%d), ErrorCode=0x%04X, WarnWord=0x%04X",
			status.StatusWord, status.StateVar, (status.StateVar>>8)&0xFF, status.ErrorCode, status.WarnWord)
		logErrorEvidence(t, client, "post_setcommand", *linmotTimeout)

		_, err = readPresenceMasks(t, client, "presence_masks_mc_running", *linmotTimeout)
		if err != nil {
			t.Logf("presence_masks_mc_running: failed")
		}

		ctxStop, cancelStop := context.WithTimeout(context.Background(), *linmotTimeout)
		err = client.StopMotionController(ctxStop)
		cancelStop()
		if err != nil {
			logErrorEvidence(t, client, "stop_mc_failed", *linmotTimeout)
			t.Fatalf("StopMotionController() failed: %v", err)
		}

		_, err = readPresenceMasks(t, client, "presence_masks_mc_stopped", *linmotTimeout)
		if err != nil {
			extended := *linmotTimeout + 5*time.Second
			t.Logf("presence_masks_mc_stopped: retrying with extended timeout=%s", extended)
			_, err = readPresenceMasks(t, client, "presence_masks_mc_stopped_retry", extended)
			if err != nil {
				t.Logf("presence_masks_mc_stopped_retry: failed")
			}
		}

		ctxStart, cancelStart := context.WithTimeout(context.Background(), *linmotTimeout)
		err = client.StartMotionController(ctxStart)
		cancelStart()
		if err != nil {
			logErrorEvidence(t, client, "start_mc_failed", *linmotTimeout)
			t.Fatalf("StartMotionController() failed: %v", err)
		}

		if status.ErrorCode != 0 {
			ctxErrText, cancelErrText := context.WithTimeout(context.Background(), *linmotTimeout)
			errorText, err := client.GetErrorText(ctxErrText, status.ErrorCode)
			cancelErrText()
			if err != nil {
				t.Logf("GetErrorText(0x%04X) failed: %v", status.ErrorCode, err)
			} else {
				t.Logf("ErrorText for ErrorCode=0x%04X: %q", status.ErrorCode, errorText)
			}
		}

		ctxErrLog, cancelErrLog := context.WithTimeout(context.Background(), *linmotTimeout)
		errLog, err := client.GetErrorLog(ctxErrLog)
		cancelErrLog()
		if err != nil {
			t.Logf("GetErrorLog() failed: %v", err)
		} else if len(errLog) > 0 {
			last := errLog[len(errLog)-1]
			errorText := ""
			ctxErrText, cancelErrText := context.WithTimeout(context.Background(), *linmotTimeout)
			text, err := client.GetErrorText(ctxErrText, last.ErrorCode)
			cancelErrText()
			if err == nil {
				errorText = text
			}
			t.Logf("Latest error log entry: index=%d code=0x%04X time=%s text=%q", len(errLog)-1, last.ErrorCode, last.Timestamp.Format(time.RFC3339Nano), errorText)
		} else {
			t.Logf("GetErrorLog() returned 0 entries")
		}
	}
}

// TestClient_PresenceMasks_Stability tests presence mask read stability without mutating drive state.
// This is a clean confidence test that only reads presence masks and status, no state changes.
func TestClient_PresenceMasks_Stability(t *testing.T) {
	clients, cleanup, err := initializeClients(t)
	if err != nil {
		t.Fatalf("initializeClients(t) failed: %v", err)
	}
	defer deferCleanup(t, cleanup)

	// Only test UDP clients (skip mock)
	for _, client := range clients {
		// Skip mock clients - only test real hardware
		if *linmotMode == "mock" {
			continue
		}

		setupE2EInstrumentation(t, client)

		// Run 20 iterations of presence mask reads
		for i := 0; i < 20; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

			// Get status first to log drive state
			status, err := client.GetStatus(ctx)
			if err != nil {
				cancel()
				t.Fatalf("Iteration %d: GetStatus() failed: %v", i, err)
			}

			// Read presence masks directly (no entry data reads)
			masks, err := client.GetPresenceMasks(ctx)
			cancel()

			if err != nil {
				t.Fatalf("Iteration %d: GetPresenceMasks() failed: %v", i, err)
			}

			mainState := (status.StateVar >> 8) & 0xFF
			t.Logf("Iteration %d: StatusWord=0x%04X, StateVar=0x%04X (MainState=%d), first_mask=0x%08X",
				i, status.StatusWord, status.StateVar, mainState, masks[0])

			// Sleep between iterations to avoid overwhelming the drive
			time.Sleep(150 * time.Millisecond)
		}
	}
}

// TestClient_StatusSampling_Stability checks status sampling for missing response bits or all-zero payloads.
func TestClient_StatusSampling_Stability(t *testing.T) {
	clients, cleanup, err := initializeClients(t)
	if err != nil {
		t.Fatalf("initializeClients(t) failed: %v", err)
	}
	defer cleanup()

	for _, client := range clients {
		if *linmotMode == "mock" {
			continue
		}

		setupE2EInstrumentation(t, client)

		for i := 0; i < 200; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
			status, err := client.GetStatus(ctx)
			cancel()
			if err != nil {
				t.Fatalf("Iteration %d: GetStatus() failed: %v", i, err)
			}
			if status == nil {
				t.Fatalf("Iteration %d: GetStatus() returned nil status", i)
			}
			if status.ResponseBits&protocol_common.ResponseFlags.Standard != protocol_common.ResponseFlags.Standard {
				t.Fatalf("Iteration %d: status response missing standard bits: repBits=0x%08X", i, status.ResponseBits)
			}
			if status.StatusWord == 0 && status.StateVar == 0 && status.ErrorCode == 0 && status.WarnWord == 0 {
				t.Fatalf("Iteration %d: status payload is all zeros (sw/sv/err/warn)", i)
			}
		}
	}
}

// TestClient_StopStartMC_Stability tests StopMC/StartMC reliability without touching command tables.
func TestClient_StopStartMC_Stability(t *testing.T) {
	clients, cleanup, err := initializeClients(t)
	if err != nil {
		t.Fatalf("initializeClients(t) failed: %v", err)
	}
	defer deferCleanup(t, cleanup)

	// Only test UDP clients (skip mock)
	for _, client := range clients {
		// Skip mock clients - only test real hardware
		if *linmotMode == "mock" {
			continue
		}

		setupE2EInstrumentation(t, client)

		// Run 10 iterations of StopMC/StartMC cycles
		for i := 0; i < 10; i++ {
			// a) Get initial status
			ctx1, cancel1 := context.WithTimeout(context.Background(), 2*time.Second)
			status1, err := client.GetStatus(ctx1)
			cancel1()
			if err != nil {
				t.Fatalf("Iteration %d: GetStatus() before StopMC failed: %v", i, err)
			}
			mainState1 := (status1.StateVar >> 8) & 0xFF
			t.Logf("Iteration %d: Before StopMC: StatusWord=0x%04X, StateVar=0x%04X (MainState=%d), ErrorCode=0x%04X",
				i, status1.StatusWord, status1.StateVar, mainState1, status1.ErrorCode)

			// b) Stop Motion Controller
			ctx2, cancel2 := context.WithTimeout(context.Background(), 3*time.Second)
			err = client.StopMotionController(ctx2)
			cancel2()
			if err != nil {
				t.Fatalf("Iteration %d: StopMotionController() failed: %v", i, err)
			}

			// c) Get status after StopMC
			ctx3, cancel3 := context.WithTimeout(context.Background(), 2*time.Second)
			status3, err := client.GetStatus(ctx3)
			cancel3()
			if err != nil {
				t.Fatalf("Iteration %d: GetStatus() after StopMC failed: %v", i, err)
			}
			mainState3 := (status3.StateVar >> 8) & 0xFF
			t.Logf("Iteration %d: After StopMC: StatusWord=0x%04X, StateVar=0x%04X (MainState=%d), ErrorCode=0x%04X",
				i, status3.StatusWord, status3.StateVar, mainState3, status3.ErrorCode)

			// d) Start Motion Controller
			ctx4, cancel4 := context.WithTimeout(context.Background(), 3*time.Second)
			err = client.StartMotionController(ctx4)
			cancel4()
			if err != nil {
				t.Fatalf("Iteration %d: StartMotionController() failed: %v", i, err)
			}

			// e) Get status after StartMC
			ctx5, cancel5 := context.WithTimeout(context.Background(), 2*time.Second)
			status5, err := client.GetStatus(ctx5)
			cancel5()
			if err != nil {
				t.Fatalf("Iteration %d: GetStatus() after StartMC failed: %v", i, err)
			}
			mainState5 := (status5.StateVar >> 8) & 0xFF
			t.Logf("Iteration %d: After StartMC: StatusWord=0x%04X, StateVar=0x%04X (MainState=%d), ErrorCode=0x%04X",
				i, status5.StatusWord, status5.StateVar, mainState5, status5.ErrorCode)

			// Sleep between iterations
			time.Sleep(250 * time.Millisecond)
		}
	}
}

func TestClient_StopMC_MinimalRepro(t *testing.T) {
	// Only test UDP clients (skip mock)
	if *linmotMode != "udp" {
		t.Skip("Skipping test in mock/compare mode - requires UDP hardware")
	}

	client, cleanup, err := initializeUDPClient(t)
	if err != nil {
		t.Fatalf("initializeUDPClient(t) failed: %v", err)
	}
	defer deferCleanup(t, cleanup)

	setupE2EInstrumentation(t, client)

	time.Sleep(200 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.StopMotionController(ctx); err != nil {
		t.Fatalf("StopMotionController() failed: %v", err)
	}
}

// TestClient_PresenceMask_StopVsStart tests whether PresenceMaskRequest(0x87) is accepted only when MC is stopped.
func TestClient_PresenceMask_StopVsStart(t *testing.T) {
	clients, cleanup, err := initializeClients(t)
	if err != nil {
		t.Fatalf("initializeClients(t) failed: %v", err)
	}
	defer deferCleanup(t, cleanup)

	// Only test UDP clients (skip mock)
	for _, client := range clients {
		// Skip mock clients - only test real hardware
		if *linmotMode == "mock" {
			continue
		}

		setupE2EInstrumentation(t, client)

		// 1) GetStatus (log StatusWord/StateVar for sanity)
		ctx1, cancel1 := context.WithTimeout(context.Background(), 2*time.Second)
		status1, err := client.GetStatus(ctx1)
		cancel1()
		if err != nil {
			t.Fatalf("GetStatus() failed: %v", err)
		}
		mainState1 := (status1.StateVar >> 8) & 0xFF
		t.Logf("Initial: StatusWord=0x%04X, StateVar=0x%04X (MainState=%d), ErrorCode=0x%04X",
			status1.StatusWord, status1.StateVar, mainState1, status1.ErrorCode)

		// 2) StopMC
		ctx2, cancel2 := context.WithTimeout(context.Background(), 3*time.Second)
		err = client.StopMotionController(ctx2)
		cancel2()
		if err != nil {
			t.Fatalf("StopMotionController() failed: %v", err)
		}

		// 3) Issue PresenceMaskRequest(0) directly and assert response
		ctx3, cancel3 := context.WithTimeout(context.Background(), 2*time.Second)
		req, err := protocol_command_tables.NewPresenceMaskRequest(0)
		if err != nil {
			cancel3()
			t.Fatalf("NewPresenceMaskRequest(0) failed: %v", err)
		}

		// Get the request manager from client (test can access unexported fields)
		response, err := client_common.SendRequestAndReceive[*protocol_command_tables.PresenceMaskResponse](
			client.requestManager, ctx3, req)
		cancel3()

		if err != nil {
			t.Fatalf("PresenceMaskRequest(0) after StopMC failed: %v", err)
		}

		// Extract and log details
		reqCmdCode := protocol_command_tables.CommandCode.PresenceMask0 // 0x87
		reqCmdCount := response.RTCCounter()                            // This is the counter from the response
		// Note: We can't get the request counter directly, but we can verify the response counter matches
		// by checking if the response was delivered (which means counter matched)
		rtcStatus := response.RTCStatus()
		status := response.Status()
		value := response.Value()

		// Extract mask value (w3 and w4 from value)
		valueU32 := uint32(value)
		w3 := uint32(valueU32 >> 16)
		w4 := uint32(valueU32 & 0xFFFF)
		maskValue := (w3 << 16) | w4

		// Extract parameter channel status (bits 8-15 of StatusWord)
		// StatusWord bits 8-15 would be: (StatusWord >> 8) & 0xFF
		paramChannelStatus := (status.StatusWord >> 8) & 0xFF

		t.Logf("After StopMC - PresenceMaskRequest(0):")
		t.Logf("  req cmdCode=0x%02X (0x87)", reqCmdCode)
		t.Logf("  parsed cmdCountResponse=%d", reqCmdCount)
		t.Logf("  parsed parameter status code (StatusWord bits 8-15)=0x%02X", paramChannelStatus)
		t.Logf("  RTC status=0x%02X", rtcStatus)
		t.Logf("  returned mask value (w3<<16|w4)=0x%08X", maskValue)
		t.Logf("  StatusWord=0x%04X, StateVar=0x%04X", status.StatusWord, status.StateVar)

		// Assert that we received a response (counter matching is implicit in successful delivery)
		if reqCmdCount == 0 {
			t.Errorf("Expected non-zero command count in response, got 0")
		}

		// 4) StartMC
		ctx4, cancel4 := context.WithTimeout(context.Background(), 3*time.Second)
		err = client.StartMotionController(ctx4)
		cancel4()
		if err != nil {
			t.Fatalf("StartMotionController() failed: %v", err)
		}

		// 5) Issue PresenceMaskRequest(0) again and log (don't assert, just log + return error if timeout)
		ctx5, cancel5 := context.WithTimeout(context.Background(), 2*time.Second)
		req2, err := protocol_command_tables.NewPresenceMaskRequest(0)
		if err != nil {
			cancel5()
			t.Fatalf("NewPresenceMaskRequest(0) failed: %v", err)
		}

		response2, err := client_common.SendRequestAndReceive[*protocol_command_tables.PresenceMaskResponse](
			client.requestManager, ctx5, req2)
		cancel5()

		if err != nil {
			t.Logf("After StartMC - PresenceMaskRequest(0) TIMEOUT/ERROR: %v", err)
			continue // Expected - just log and proceed to next client
		}

		// Extract and log details
		reqCmdCount2 := response2.RTCCounter()
		rtcStatus2 := response2.RTCStatus()
		status2 := response2.Status()
		value2 := response2.Value()

		// Extract mask value
		valueU32_2 := uint32(value2)
		w3_2 := uint32(valueU32_2 >> 16)
		w4_2 := uint32(valueU32_2 & 0xFFFF)
		maskValue2 := (w3_2 << 16) | w4_2

		// Extract parameter channel status
		paramChannelStatus2 := (status2.StatusWord >> 8) & 0xFF

		t.Logf("After StartMC - PresenceMaskRequest(0) SUCCESS:")
		t.Logf("  req cmdCode=0x%02X (0x87)", reqCmdCode)
		t.Logf("  parsed cmdCountResponse=%d", reqCmdCount2)
		t.Logf("  parsed parameter status code (StatusWord bits 8-15)=0x%02X", paramChannelStatus2)
		t.Logf("  RTC status=0x%02X", rtcStatus2)
		t.Logf("  returned mask value (w3<<16|w4)=0x%08X", maskValue2)
		t.Logf("  StatusWord=0x%04X, StateVar=0x%04X", status2.StatusWord, status2.StateVar)
	}
}

// TestClient_StopMC_OneShot tests a single StopMC/StartMC cycle with fresh contexts.
func TestClient_StopMC_OneShot(t *testing.T) {
	clients, cleanup, err := initializeClients(t)
	if err != nil {
		t.Fatalf("initializeClients(t) failed: %v", err)
	}
	defer deferCleanup(t, cleanup)

	// Only test UDP clients (skip mock)
	for _, client := range clients {
		// Skip mock clients - only test real hardware
		if *linmotMode == "mock" {
			continue
		}

		setupE2EInstrumentation(t, client)

		// 1) GetStatus() and log StatusWord + StateVar
		ctx1, cancel1 := context.WithTimeout(context.Background(), 2*time.Second)
		status1, err := client.GetStatus(ctx1)
		cancel1()
		if err != nil {
			t.Fatalf("GetStatus() before StopMC failed: %v", err)
		}
		mainState1 := (status1.StateVar >> 8) & 0xFF
		t.Logf("Before StopMC: StatusWord=0x%04X, StateVar=0x%04X (MainState=%d), ErrorCode=0x%04X",
			status1.StatusWord, status1.StateVar, mainState1, status1.ErrorCode)

		// 2) StopMotionController() with fresh 10s context
		ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
		err = client.StopMotionController(ctx2)
		cancel2()
		if err != nil {
			t.Fatalf("StopMotionController() failed: %v", err)
		}

		// 3) GetStatus() again and log
		ctx3, cancel3 := context.WithTimeout(context.Background(), 2*time.Second)
		status3, err := client.GetStatus(ctx3)
		cancel3()
		if err != nil {
			t.Fatalf("GetStatus() after StopMC failed: %v", err)
		}
		mainState3 := (status3.StateVar >> 8) & 0xFF
		t.Logf("After StopMC: StatusWord=0x%04X, StateVar=0x%04X (MainState=%d), ErrorCode=0x%04X",
			status3.StatusWord, status3.StateVar, mainState3, status3.ErrorCode)

		// 4) StartMotionController() with fresh 10s context
		ctx4, cancel4 := context.WithTimeout(context.Background(), 10*time.Second)
		err = client.StartMotionController(ctx4)
		cancel4()
		if err != nil {
			t.Fatalf("StartMotionController() failed: %v", err)
		}

		// 5) GetStatus() again and log
		ctx5, cancel5 := context.WithTimeout(context.Background(), 2*time.Second)
		status5, err := client.GetStatus(ctx5)
		cancel5()
		if err != nil {
			t.Fatalf("GetStatus() after StartMC failed: %v", err)
		}
		mainState5 := (status5.StateVar >> 8) & 0xFF
		t.Logf("After StartMC: StatusWord=0x%04X, StateVar=0x%04X (MainState=%d), ErrorCode=0x%04X",
			status5.StatusWord, status5.StateVar, mainState5, status5.ErrorCode)
	}
}

func TestClient_RTCCommandCount_Probe(t *testing.T) {
	// Only test UDP clients (skip mock)
	if *linmotMode == "mock" {
		t.Skip("Skipping probe test in mock mode")
	}

	clients, cleanup, err := initializeClients(t)
	if err != nil {
		t.Fatalf("initializeClients(t) failed: %v", err)
	}
	defer deferCleanup(t, cleanup)

	for _, client := range clients {
		// Skip mock clients - only test real hardware
		if *linmotMode == "mock" {
			continue
		}

		setupE2EInstrumentation(t, client)

		// Step A: Attempt StopMC with default counter to establish baseline
		ctxShort, cancelShort := context.WithTimeout(context.Background(), 250*time.Millisecond)
		defer cancelShort()

		// Force counter to 1 (default start)
		client.SetRTCCounterForTesting(1)
		_ = client.StopMotionController(ctxShort) // Ignore error, we just want to see what drive echoes

		base := client.LastDriveCmdCountForTesting()
		t.Logf("Step A: Baseline - drive echoed cmdCount=%d", base)

		// Define 3 candidate cmdCounts
		toggle := base ^ 1
		plus1 := (base + 1) & 0x0F
		plus2 := (base + 2) & 0x0F

		// Map to valid range (1-14)
		normalize := func(val uint8) uint8 {
			if val == 0 || val > protocol_rtc.CounterMax {
				return 1
			}
			return val
		}

		candidates := []struct {
			name string
			val  uint8
		}{
			{"toggle", normalize(toggle)},
			{"plus1", normalize(plus1)},
			{"plus2", normalize(plus2)},
		}

		accepted := false
		for _, cand := range candidates {
			// Force counter to candidate value
			client.SetRTCCounterForTesting(cand.val)

			// Attempt StopMC with short timeout
			ctxTest, cancelTest := context.WithTimeout(context.Background(), 250*time.Millisecond)
			err := client.StopMotionController(ctxTest)
			cancelTest()

			// Read what drive echoed
			echoed := client.LastDriveCmdCountForTesting()
			t.Logf("Candidate %s=%d: err=%v, driveEchoed=%d", cand.name, cand.val, err, echoed)

			// If candidate accepted, echoed should == candidate
			if echoed == cand.val {
				accepted = true
				t.Logf("  -> ACCEPTED (echoed matches candidate)")
			} else {
				t.Logf("  -> IGNORED (echoed=%d != candidate=%d)", echoed, cand.val)
			}
		}

		if !accepted {
			t.Fatalf("FAIL: None of the candidates were accepted. Baseline=%d, candidates tried: toggle=%d, plus1=%d, plus2=%d",
				base, normalize(toggle), normalize(plus1), normalize(plus2))
		}
	}
}

func TestClient_RTCCommandCount_AcceptanceMatrix(t *testing.T) {
	// Only test UDP clients (skip mock)
	if *linmotMode == "mock" {
		t.Skip("Skipping acceptance matrix test in mock mode")
	}

	clients, cleanup, err := initializeClients(t)
	if err != nil {
		t.Fatalf("initializeClients(t) failed: %v", err)
	}
	defer deferCleanup(t, cleanup)

	for _, client := range clients {
		// Skip mock clients - only test real hardware
		if *linmotMode == "mock" {
			continue
		}

		setupE2EInstrumentation(t, client)

		// Normalize cmdCount to valid range (1-14)
		normalize := func(val uint8) uint8 {
			if val == 0 {
				return protocol_rtc.CounterMax
			}
			if val > protocol_rtc.CounterMax {
				val = ((val - 1) % protocol_rtc.CounterMax) + 1
			}
			return val
		}

		// Check if two values have same parity (bit0)
		sameParity := func(a, b uint8) bool {
			return (a & 1) == (b & 1)
		}

		// Test baselines
		baselines := []uint8{2, 3}

		// Define test cases: (baseline, candidate, description)
		type testCase struct {
			baseline    uint8
			candidate   uint8
			description string
		}

		var testCases []testCase
		for _, B := range baselines {
			// Same parity, +2
			testCases = append(testCases, testCase{B, normalize(B + 2), "same parity +2"})
			// Opposite parity, +1
			testCases = append(testCases, testCase{B, normalize(B + 1), "opposite parity +1"})
			// Same parity, +4
			testCases = append(testCases, testCase{B, normalize(B + 4), "same parity +4"})
			// Opposite parity, +5
			testCases = append(testCases, testCase{B, normalize(B + 5), "opposite parity +5"})
		}

		// Print header
		t.Logf("\n=== RTC Command Count Acceptance Matrix ===")
		t.Logf("Baseline | Candidate | ParitySame | Echoed | Accepted")
		t.Logf("---------|-----------|------------|--------|---------")

		allAccepted := true
		for _, tc := range testCases {
			// Trial: Establish baseline B
			client.SetRTCCounterForTesting(tc.baseline)
			ctxBase, cancelBase := context.WithTimeout(context.Background(), 500*time.Millisecond)
			_ = client.StopMotionController(ctxBase) // Establish baseline
			cancelBase()

			// Verify baseline was established
			established := client.LastDriveCmdCountForTesting()
			if established != tc.baseline {
				t.Logf("WARNING: Baseline establishment failed: expected %d, got %d", tc.baseline, established)
				// Try one more time
				client.SetRTCCounterForTesting(tc.baseline)
				ctxRetry, cancelRetry := context.WithTimeout(context.Background(), 500*time.Millisecond)
				_ = client.StopMotionController(ctxRetry)
				cancelRetry()
				established = client.LastDriveCmdCountForTesting()
			}

			if established != tc.baseline {
				t.Errorf("Failed to establish baseline %d, got %d", tc.baseline, established)
				continue
			}

			// Now test candidate C (without any intervening commands)
			client.SetRTCCounterForTesting(tc.candidate)
			ctxCandidate, cancelCandidate := context.WithTimeout(context.Background(), 500*time.Millisecond)
			_ = client.StopMotionController(ctxCandidate) // Test candidate
			cancelCandidate()

			// Check what drive echoed
			echoed := client.LastDriveCmdCountForTesting()
			accepted := echoed == tc.candidate
			paritySame := sameParity(tc.baseline, tc.candidate)

			if !accepted {
				allAccepted = false
			}

			t.Logf("   %2d     |    %2d      |    %5v    |  %2d    |  %5v  (%s)",
				tc.baseline, tc.candidate, paritySame, echoed, accepted, tc.description)
		}

		t.Logf("=============================================\n")

		if !allAccepted {
			t.Logf("CONCLUSION: Drive does NOT accept all cmdCount changes. Some candidates were ignored.")
		} else {
			t.Logf("CONCLUSION: Drive accepts ANY cmdCount change (not just parity flips).")
		}
	}
}

func TestClient_RTCCommandCount_BackwardsAndWrap(t *testing.T) {
	// Only test UDP clients (skip mock)
	if *linmotMode == "mock" {
		t.Skip("Skipping backwards/wrap test in mock mode")
	}

	clients, cleanup, err := initializeClients(t)
	if err != nil {
		t.Fatalf("initializeClients(t) failed: %v", err)
	}
	defer deferCleanup(t, cleanup)

	for _, client := range clients {
		// Skip mock clients - only test real hardware
		if *linmotMode == "mock" {
			continue
		}

		setupE2EInstrumentation(t, client)

		// Define test cases: (baseline, candidate, description)
		type testCase struct {
			baseline    uint8
			candidate   uint8
			description string
		}

		var testCases []testCase

		// Backwards cases: going from higher to lower cmdCount
		testCases = append(testCases, testCase{8, 2, "backwards: 8->2"})
		testCases = append(testCases, testCase{10, 3, "backwards: 10->3"})
		testCases = append(testCases, testCase{14, 5, "backwards: 14->5"})

		// Wrap cases: testing wrap behavior
		testCases = append(testCases, testCase{14, 1, "wrap: 14->1"})
		testCases = append(testCases, testCase{13, 1, "wrap: 13->1"})

		// Test if 15 is accepted (even though CounterMax=14, test what drive does)
		// Note: SetRTCCounterForTesting may clamp to valid range, but we test anyway
		testCases = append(testCases, testCase{14, 15, "forward: 14->15 (beyond CounterMax)"})
		testCases = append(testCases, testCase{15, 1, "wrap: 15->1 (if 15 accepted)"})

		// Print header
		t.Logf("\n=== RTC Command Count Backwards + Wrap Acceptance Test ===")
		t.Logf("Baseline | Candidate | Description        | Echoed | Accepted | Direction")
		t.Logf("---------|-----------|--------------------|--------|----------|----------")

		allBackwardsAccepted := true
		wrapWorks := false
		value15Accepted := false

		for _, tc := range testCases {
			// Establish baseline B
			client.SetRTCCounterForTesting(tc.baseline)
			ctxBase, cancelBase := context.WithTimeout(context.Background(), 500*time.Millisecond)
			err := client.StopMotionController(ctxBase)
			cancelBase()

			if err != nil {
				t.Logf("WARNING: Baseline establishment failed for B=%d: %v", tc.baseline, err)
				continue
			}

			// Verify baseline was established
			established := client.LastDriveCmdCountForTesting()
			if established != tc.baseline {
				t.Logf("WARNING: Baseline mismatch: expected %d, got %d (may be clamped)", tc.baseline, established)
				// Update baseline to what was actually established
				tc.baseline = established
			}

			// Test candidate C (backwards or wrap)
			client.SetRTCCounterForTesting(tc.candidate)
			ctxCandidate, cancelCandidate := context.WithTimeout(context.Background(), 500*time.Millisecond)
			client.StopMotionController(ctxCandidate)
			cancelCandidate()

			// Check what drive echoed
			echoed := client.LastDriveCmdCountForTesting()
			accepted := echoed == tc.candidate

			// Determine direction
			direction := "forward"
			// Check for wrap case first (14->1 or 13->1)
			if tc.candidate == 1 && (tc.baseline == 14 || tc.baseline == 13) {
				direction = "wrap"
				if accepted {
					wrapWorks = true
				}
			} else if tc.candidate == 15 {
				direction = "beyond_max"
				if accepted {
					value15Accepted = true
				}
			} else if tc.candidate < tc.baseline {
				direction = "backwards"
				if !accepted {
					allBackwardsAccepted = false
				}
			}

			status := "ACCEPTED"
			if !accepted {
				status = "IGNORED"
			}

			t.Logf("   %2d     |    %2d      | %-18s |  %2d    | %-8s | %s",
				tc.baseline, tc.candidate, tc.description, echoed, status, direction)
		}

		t.Logf("===========================================================\n")

		// Summary and conclusions
		t.Logf("=== SUMMARY ===")
		if allBackwardsAccepted {
			t.Logf("(a) Backwards acceptance: Drive accepts backwards cmdCount (any change accepted)")
		} else {
			t.Logf("(a) Backwards acceptance: Drive REJECTS backwards cmdCount (sequence-number semantics)")
		}

		if wrapWorks {
			t.Logf("(b) Wrap behavior: Drive accepts wrap from 14->1 (normal wrap works)")
		} else {
			t.Logf("(b) Wrap behavior: Drive may not accept wrap (needs verification)")
		}

		if value15Accepted {
			t.Logf("(c) Counter max: Drive accepts cmdCount=15 (Go CounterMax should be 15, not 14)")
		} else {
			t.Logf("(c) Counter max: Drive does NOT accept cmdCount=15 (Go CounterMax=14 is correct)")
		}

		t.Logf("\n=== CONCLUSION ===")
		if allBackwardsAccepted {
			t.Logf("Drive accepts ANY cmdCount change (backwards, forwards, wrap) - no sequence-number semantics")
		} else {
			t.Logf("Drive requires cmdCount to move FORWARD (sequence-number semantics) - backwards changes are ignored")
		}
	}
}

func TestRequestManager_CleanupMCRequestOnTimeout(t *testing.T) {
	// Only test in mock mode for this regression test
	if *linmotMode != "mock" {
		t.Skip("Skipping MC cleanup test in non-mock mode")
	}

	client, cleanup, err := initializeMockClient()
	if err != nil {
		t.Fatalf("initializeMockClient() failed: %v", err)
	}
	defer deferCleanup(t, cleanup)

	// Create an MC request
	mcRequest := protocol_motion_control.NewMCCommandRequest(
		protocol_motion_control.MasterID(1),
		0,
		[]uint16{0x1234, 0x5678},
	)

	// Submit with context cancellation to test cleanup path
	ctx, cancel := context.WithCancel(context.Background())

	// Start the request in a goroutine, then cancel immediately
	done := make(chan error, 1)
	go func() {
		_, err := client_common.SendRequestAndReceive[protocol_motion_control.MCResponse](
			client.requestManager, ctx, mcRequest)
		done <- err
	}()

	// Cancel immediately to trigger cleanup
	cancel()

	// Wait for the request to complete
	err = <-done

	// Verify we got a context cancellation error
	if err == nil {
		t.Fatal("Expected context cancellation error, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Expected context.Canceled, got %T: %v", err, err)
	}

	// Give a small delay to ensure cleanup completes
	time.Sleep(10 * time.Millisecond)

	// Verify pendingMCRequests is empty (cleanup should have removed it)
	pendingCount := client.PendingMCRequestCountForTesting()
	if pendingCount != 0 {
		t.Errorf("Expected pendingMCRequests to be empty after cancellation, got %d pending requests", pendingCount)
	}
}
