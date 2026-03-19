package client

import (
	"context"
	stderrors "errors"
	"fmt"
	"os"
	"time"

	"github.com/pkg/errors"

	client_common "github.com/Smart-Vision-Works/linmot_client/client/common"
	client_control_word "github.com/Smart-Vision-Works/linmot_client/client/control_word"
	client_monitoring "github.com/Smart-Vision-Works/linmot_client/client/monitoring"
	client_motion_control "github.com/Smart-Vision-Works/linmot_client/client/motion_control"
	client_rtc "github.com/Smart-Vision-Works/linmot_client/client/rtc"
	client_command_tables "github.com/Smart-Vision-Works/linmot_client/client/rtc/command_tables"
	client_errors "github.com/Smart-Vision-Works/linmot_client/client/rtc/errors"
	client_parameters "github.com/Smart-Vision-Works/linmot_client/client/rtc/parameters"
	protocol_common "github.com/Smart-Vision-Works/linmot_client/protocol/common"
	transport "github.com/Smart-Vision-Works/linmot_client/transport"
)

const (
	defaultTimerCycle          = 5 * time.Millisecond
	defaultConnectivityTimeout = 5 * time.Second // Timeout for initial connectivity validation
	envSkipConnectivityCheck   = "LINMOT_SKIP_VALIDATE_CONNECTIVITY"
)

// Client communicates with LinMot drives via the LinUDP V2 protocol.
//
// Client uses composition with specialized managers to organize functionality:
//   - controlWordManager: control word and status-related helpers
//   - rtcManager: real-time RTC, command table, and parameter helpers
//   - monitoringManager: monitoring channel configuration and data retrieval
//   - mcInterfaceManager: motion control command helpers (VAI/V motion commands)
type Client struct {
	driveIP            string
	requestManager     *client_common.RequestManager
	controlWordManager *client_control_word.ControlWordManager
	rtcManager         *client_rtc.RtcManager
	monitoringManager  *client_monitoring.MonitoringManager
	mcInterfaceManager *client_motion_control.MotionControlManager
}

// NewUDPClient creates a new UDP client for communicating with a LinMot drive.
// It validates connectivity by sending a status request before returning.
// Returns an error if the drive cannot be reached or does not respond.
// The underlying request manager starts its TX/RX goroutines immediately, so callers must
// call Close() once the client is no longer needed to stop those workers.
func NewUDPClient(driveIP string, drivePort, masterPort int, bindAddress string, timeout time.Duration) (*Client, error) {
	return newUDPClient(driveIP, drivePort, masterPort, bindAddress, timeout, false)
}

// NewUDPClientWithDebug creates a new UDP client and enables debug logging before validation.
func NewUDPClientWithDebug(driveIP string, drivePort, masterPort int, bindAddress string, timeout time.Duration, debug bool) (*Client, error) {
	return newUDPClient(driveIP, drivePort, masterPort, bindAddress, timeout, debug)
}

func newUDPClient(driveIP string, drivePort, masterPort int, bindAddress string, timeout time.Duration, debug bool) (*Client, error) {
	transportClient, err := transport.NewUDPTransportClient(driveIP, drivePort, masterPort, bindAddress, timeout)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to create UDP transport")
	}

	client := newClient(transportClient)
	client.driveIP = driveIP
	if debug {
		client.SetDebug(true)
		drained := drainUDPPackets(transportClient, 5, 2*time.Millisecond)
		if client.requestManager.DebugEnabled() {
			fmt.Printf("[UDP_DRAIN] drained=%d\n", drained)
		}
	}

	// Validate connectivity by sending a status request
	if err := client.validateConnectivity(defaultConnectivityTimeout); err != nil {
		// Clean up the client if validation fails
		_ = client.Close()
		return nil, errors.WithMessagef(err, "failed to validate connectivity to drive at %s:%d", driveIP, drivePort)
	}

	return client, nil
}

func drainUDPPackets(transportClient transport.Client, maxDrain int, perAttempt time.Duration) int {
	drained := 0
	for i := 0; i < maxDrain; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), perAttempt)
		_, err := transportClient.RecvPacket(ctx)
		cancel()
		if err != nil {
			break
		}
		drained++
	}
	return drained
}

// NewClientWithTransport creates a new client using an existing transport.Client.
// This is useful when using a shared transport (e.g., SharedUDPTransport) to avoid
// port binding conflicts when multiple drives share the same master port.
//
// The transport should already be initialized and ready to send/receive.
// It validates connectivity by sending a status request before returning.
// Returns an error if the drive cannot be reached or does not respond.
//
// The underlying request manager starts its TX/RX goroutines immediately, so callers must
// call Close() once the client is no longer needed to stop those workers.
func NewClientWithTransport(driveIP string, transportClient transport.Client) (*Client, error) {
	client := newClient(transportClient)
	client.driveIP = driveIP

	// Validate connectivity by sending a status request
	// The transport should already be active, so this validates the drive responds
	if os.Getenv(envSkipConnectivityCheck) == "" {
		if err := client.validateConnectivity(defaultConnectivityTimeout); err != nil {
			// Clean up the client if validation fails
			return nil, errors.WithMessagef(err, "failed to validate connectivity to drive at %s", driveIP)
		}
	}

	return client, nil
}

// NewMockClient creates a new mock client for testing.
// Returns both the client and the transport server for setting up test scenarios.
// The request manager is started immediately; close the client to stop background goroutines.
func NewMockClient() (*Client, transport.Server) {
	transportClient, transportServer := transport.NewMockTransportClientWithServer()

	client := newClient(transportClient)
	client.driveIP = ""
	return client, transportServer
}

func newClient(transportClient transport.Client) *Client {
	requestManager := client_common.NewRequestManager(transportClient, defaultTimerCycle)
	rtcManager := client_rtc.NewRtcManager(requestManager)
	client := &Client{
		driveIP:            "",
		requestManager:     requestManager,
		controlWordManager: client_control_word.NewControlWordManager(requestManager),
		rtcManager:         rtcManager,
		monitoringManager:  client_monitoring.NewMonitoringManager(requestManager),
		mcInterfaceManager: client_motion_control.NewMotionControlManager(requestManager),
	}
	client.requestManager.Start()
	return client
}

// Close closes the transport connection and stops background workers.
func (client *Client) Close() error {
	return client.requestManager.Stop()
}

// SetDebug enables or disables debug logging for client operations.
func (client *Client) SetDebug(enabled bool) {
	client.requestManager.SetDebug(enabled)
	client.rtcManager.SetDebug(enabled)
}

// validateConnectivity validates that the client can communicate with the drive/server
// by attempting two strategies:
// 1. Send status requests (reqDef=0x1FF) to elicit proper status response
// 2. If status requests fail, fallback to C#-style cyclic telegram probes (reqDef=6/7)
//
// The second strategy handles drives that may ignore status-only requests but
// respond to continuous cyclic commands, as per the C# library behavior.
// Returns an error if both strategies fail or timeout.
func (client *Client) validateConnectivity(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	maxAttempts := 100 // Safety limit for status request loop

	if client.requestManager.DebugEnabled() {
		fmt.Printf("[VALIDATE_CONN] timeout=%v strategy=status_first\n", timeout)
	}

	// Strategy 1: Try status requests (0x7F then 0x1FF)
	statusResult, err := client.runStatusConnectivityStrategy(ctx, maxAttempts)
	if err != nil {
		return err
	}
	if statusResult.success {
		return nil
	}

	// Strategy 2: Fallback to C#-style cyclic telegram probe
	if client.requestManager.DebugEnabled() {
		fmt.Printf("[VALIDATE_CONN] status requests exhausted, attempting cyclic telegram probe\n")
	}

	// Check if we can still make attempts within the remaining context timeout
	if err := buildStatusTimeoutError(ctx, statusResult.invalidCount, statusResult.lastReason); err != nil {
		return err
	}

	// Send cyclic telegram probe via transport
	probeResult, err := client.sendCyclicTelegramProbe(ctx)
	if err != nil {
		return errors.WithMessage(err, "cyclic telegram probe failed")
	}

	if client.requestManager.DebugEnabled() {
		for _, line := range probeResult.RawOutput {
			fmt.Printf("[CYCLIC_PROBE] %s\n", line)
		}
	}

	if probeResult.Success {
		if client.requestManager.DebugEnabled() {
			fmt.Printf("[VALIDATE_CONN] SUCCESS: cyclic probe elicited %d response(s)\n", probeResult.Responses)
		}
		// Cyclic probe succeeded; treat as connectivity valid
		// Note: We skip assertDriveOperational here because cyclic responses are not full status frames
		return nil
	}

	// Both strategies failed
	message := buildConnectivityFailureMessage(statusResult.attempt, statusResult.invalidCount, statusResult.lastReason, probeResult)
	if client.requestManager.DebugEnabled() {
		fmt.Printf("[VALIDATE_HINT] Connectivity failed. Possible causes:\n")
		fmt.Printf("[VALIDATE_HINT]   1. Drive's UDP telegram support is not enabled in firmware\n")
		fmt.Printf("[VALIDATE_HINT]   2. Drive has Master Filter mode enabled (IP whitelist) - check LinMot-Talk for filter config\n")
		fmt.Printf("[VALIDATE_HINT]   3. Network return path blocked (asymmetric routing or firewall)\n")
		fmt.Printf("[VALIDATE_HINT]   4. Another master (e.g., LinMot-Talk) is connected - closes drive to other masters\n")
	}
	return errors.New(message)
}

type statusConnectivityResult struct {
	success      bool
	attempt      int
	invalidCount int
	lastReason   string
}

func (client *Client) runStatusConnectivityStrategy(ctx context.Context, maxAttempts int) (statusConnectivityResult, error) {
	result := statusConnectivityResult{}
	hintLogged := false
	useExtendedRepBits := false

	for result.attempt < maxAttempts {
		result.attempt++
		if client.requestManager.DebugEnabled() {
			repBitsStr := "0x7F"
			if useExtendedRepBits {
				repBitsStr = "0x1FF"
			}
			fmt.Printf("[VALIDATE_ATTEMPT] attempt=%d/%d ts=%s repBits=%s\n", result.attempt, maxAttempts, time.Now().Format(time.RFC3339Nano), repBitsStr)
		}

		var status *protocol_common.Status
		var err error
		if useExtendedRepBits {
			// Use extended repBits (0x1FF) for connectivity probe.
			status, err = client.getStatusWithRepBits(ctx, protocol_common.ResponseFlags.RTCReply)
		} else {
			// Try standard repBits (0x7F) first.
			status, err = client.GetStatus(ctx)
			// If it fails (timeout or any error), switch to extended repBits on next attempt.
			if err != nil {
				useExtendedRepBits = true
				continue
			}
		}

		if err != nil {
			if stderrors.Is(err, client_common.ErrInvalidStatusTelegram) {
				result.invalidCount++
				if snapshot, ok := client.requestManager.LastInvalidStatusSnapshot(); ok {
					result.lastReason = snapshot.Reason
					if client.requestManager.DebugEnabled() && snapshot.Reason == "all_zero_status" && !hintLogged {
						fmt.Printf("[STATUS_HINT] hint=\"possible multi-master or invalid telegram (e.g., LinMot-Talk connected)\"\n")
						hintLogged = true
					}
				}
				if ctx.Err() != nil {
					break
				}
				select {
				case <-ctx.Done():
					break
				case <-time.After(20 * time.Millisecond):
				}
				continue
			}
			// Other errors: log and fall through to cyclic probe fallback.
			if client.requestManager.DebugEnabled() {
				fmt.Printf("[VALIDATE_ATTEMPT] status request failed: %v, falling back to cyclic telegram probe\n", err)
			}
			break
		}

		if err := assertDriveOperational(status); err != nil {
			return result, err
		}
		result.success = true
		return result, nil
	}

	return result, nil
}

func buildStatusTimeoutError(ctx context.Context, invalidCount int, lastReason string) error {
	if ctx.Err() == nil {
		return nil
	}
	if invalidCount > 0 {
		message := fmt.Sprintf("connectivity validation failed: invalid telemetry hits=%d", invalidCount)
		if lastReason != "" {
			message = fmt.Sprintf("%s last_reason=%s", message, lastReason)
		}
		return errors.WithMessage(client_common.ErrInvalidStatusTelegram, message)
	}
	return errors.WithMessage(ctx.Err(), "connectivity validation timeout (status strategy)")
}

func buildConnectivityFailureMessage(attempt int, invalidCount int, lastReason string, probeResult *transport.ProbeResult) string {
	message := fmt.Sprintf("connectivity validation failed: status_attempts=%d invalid_hits=%d cyclic_probe_responses=%d/%d",
		attempt, invalidCount, probeResult.Responses, probeResult.Attempts)
	if lastReason != "" {
		message = fmt.Sprintf("%s last_reason=%s", message, lastReason)
	}
	if probeResult.FirstError != "" {
		message = fmt.Sprintf("%s probe_error=%s", message, probeResult.FirstError)
	}
	return message
}

// sendCyclicTelegramProbe attempts to elicit a UDP response using C#-style cyclic telegrams.
// This is a fallback connectivity check when status requests don't work.
func (client *Client) sendCyclicTelegramProbe(ctx context.Context) (*transport.ProbeResult, error) {
	// Get transport client
	transportClient := client.requestManager.TransportClient()
	if transportClient == nil {
		return nil, errors.New("transport client unavailable")
	}

	// Use default probes (reqDef=7 first, then reqDef=6)
	probes := transport.DefaultCyclicProbes()

	// Send probes with 5 attempts per type, return on first success
	result, err := transport.SendCyclicTelegramProbe(ctx, transportClient, probes, 5, true)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetStatus gets the current drive status.
func (client *Client) GetStatus(ctx context.Context) (*protocol_common.Status, error) {
	request := protocol_common.NewStatusRequest()
	response, err := client_common.SendRequestAndReceive[*protocol_common.StatusResponse](client.requestManager, ctx, request)
	if err != nil {
		return nil, err
	}
	return response.Status(), nil
}

// getStatusWithRepBits sends a status request with custom repBits and returns the status.
// This is used for connectivity probes that request extended response data (e.g., 0x1FF).
func (client *Client) getStatusWithRepBits(ctx context.Context, repBits uint32) (*protocol_common.Status, error) {
	request := protocol_common.NewConnectivityProbeRequest(repBits)
	response, err := client_common.SendRequestAndReceive[*protocol_common.StatusResponse](client.requestManager, ctx, request)
	if err != nil {
		return nil, err
	}
	return response.Status(), nil
}

// SetRunMode sets the run mode with the specified storage type.
func (client *Client) SetRunMode(ctx context.Context, mode protocol_common.RunMode, storageType protocol_common.ParameterStorageType) error {
	return client.rtcManager.SetRunMode(ctx, mode, storageType)
}

// Motion Control

// GetPosition returns the current position in millimeters.
func (client *Client) GetPosition(ctx context.Context) (float64, error) {
	status, err := client.GetStatus(ctx)
	if err != nil {
		return 0, err
	}

	return status.ActualPositionMM(), nil
}

// SetPosition1 sets the target position for Position 1 (upid 0x145A) in millimeters.
// Position 1 has independent speed, acceleration, deceleration, and wait parameters.
func (client *Client) SetPosition1(ctx context.Context, positionMM float64, storageType protocol_common.ParameterStorageType) error {
	return client.rtcManager.SetPosition1(ctx, positionMM, storageType)
}

// SetPosition2 sets the target position for Position 2 (upid 0x145F) in millimeters.
// Position 2 has independent speed, acceleration, deceleration, and wait parameters.
func (client *Client) SetPosition2(ctx context.Context, positionMM float64, storageType protocol_common.ParameterStorageType) error {
	return client.rtcManager.SetPosition2(ctx, positionMM, storageType)
}

// GetVelocity returns the maximum velocity in meters per second.
func (client *Client) GetVelocity(ctx context.Context) (float64, error) {
	return client.rtcManager.GetVelocity(ctx)
}

// SetVelocity sets the maximum velocity in meters per second.
func (client *Client) SetVelocity(ctx context.Context, velocityMS float64, storageType protocol_common.ParameterStorageType) error {
	return client.rtcManager.SetVelocity(ctx, velocityMS, storageType)
}

// GetAcceleration returns the acceleration in meters per second squared.
func (client *Client) GetAcceleration(ctx context.Context) (float64, error) {
	return client.rtcManager.GetAcceleration(ctx)
}

// SetAcceleration sets the acceleration in meters per second squared.
func (client *Client) SetAcceleration(ctx context.Context, accelMS2 float64, storageType protocol_common.ParameterStorageType) error {
	return client.rtcManager.SetAcceleration(ctx, accelMS2, storageType)
}

// GetDeceleration returns the deceleration in meters per second squared.
func (client *Client) GetDeceleration(ctx context.Context) (float64, error) {
	return client.rtcManager.GetDeceleration(ctx)
}

// SetDeceleration sets the deceleration in meters per second squared.
func (client *Client) SetDeceleration(ctx context.Context, decelMS2 float64, storageType protocol_common.ParameterStorageType) error {
	return client.rtcManager.SetDeceleration(ctx, decelMS2, storageType)
}

// Configuration //

// SetEasyStepsAutoStart sets the Easy Steps auto start configuration.
// Use protocol_common.EasyStepsAutoStart.Enabled or .Disabled as the value.
func (client *Client) SetEasyStepsAutoStart(ctx context.Context, value int32, storageType protocol_common.ParameterStorageType) error {
	return client.rtcManager.SetEasyStepsAutoStart(ctx, value, storageType)
}

// SetEasyStepsAutoHome sets the Easy Steps auto home configuration.
// Use protocol_common.EasyStepsAutoHome.Enabled or .Disabled as the value.
func (client *Client) SetEasyStepsAutoHome(ctx context.Context, value int32, storageType protocol_common.ParameterStorageType) error {
	return client.rtcManager.SetEasyStepsAutoHome(ctx, value, storageType)
}

// SetEasyStepsInputRisingEdgeFunction sets the Easy Steps rising edge action for any input pin.
// inputNumber should be protocol_common.IOPin.Input45, .Input46, .Input47, or .Input48.
// Use protocol_common.EasyStepsIOMotion constants as the value.
func (client *Client) SetEasyStepsInputRisingEdgeFunction(ctx context.Context, inputNumber protocol_common.IOPinNumber, value int32, storageType protocol_common.ParameterStorageType) error {
	return client.rtcManager.SetEasyStepsRisingEdge(ctx, inputNumber, value, storageType)
}

// SetEasyStepsInputCurveCmdID sets the Easy Steps IO motion config curve/command table ID for any input pin.
// inputNumber should be protocol_common.IOPin.Input45, .Input46, .Input47, or .Input48.
func (client *Client) SetEasyStepsInputCurveCmdID(ctx context.Context, inputNumber protocol_common.IOPinNumber, curveCmdID int32, storageType protocol_common.ParameterStorageType) error {
	return client.rtcManager.SetEasyStepsIOMotionConfigCmd(ctx, inputNumber, curveCmdID, storageType)
}

// SetIODefOutputFunction sets an output pin function configuration for any output pin.
// outputNumber should be protocol_common.IOPin.Output36, .Output43, or .Output44.
// Use protocol_common.OutputConfig constants as the value.
func (client *Client) SetIODefOutputFunction(ctx context.Context, outputNumber protocol_common.IOPinNumber, value int32, storageType protocol_common.ParameterStorageType) error {
	return client.rtcManager.SetOutputFunction(ctx, outputNumber, value, storageType)
}

// SetIODefInputFunction sets an input pin function configuration for any input pin.
// inputNumber should be protocol_common.IOPin.Input45, .Input46, .Input47, or .Input48.
// Use protocol_common.InputFunction constants as the value.
func (client *Client) SetIODefInputFunction(ctx context.Context, inputNumber protocol_common.IOPinNumber, value int32, storageType protocol_common.ParameterStorageType) error {
	return client.rtcManager.SetInputFunction(ctx, inputNumber, value, storageType)
}

// SetTriggerMode sets the trigger mode configuration.
// Use protocol_common.TriggerModeConfig constants as the value.
func (client *Client) SetTriggerMode(ctx context.Context, value int32, storageType protocol_common.ParameterStorageType) error {
	return client.rtcManager.SetTriggerMode(ctx, value, storageType)
}

// Command Tables //

// GetCommandTable retrieves the current command table from the drive.
// Note that Version and DriveModel will be empty strings as these metadata fields are not stored on the drive.
func (client *Client) GetCommandTable(ctx context.Context) (*client_command_tables.CommandTable, error) {
	return client.rtcManager.GetCommandTable(ctx)
}

// GetPresenceMasks retrieves the presence mask values from the drive.
// Returns 8 uint32 masks indicating which entries exist in each range (mask 0 = entries 0-31, etc.).
func (client *Client) GetPresenceMasks(ctx context.Context) ([8]uint32, error) {
	return client.rtcManager.GetPresenceMasks(ctx)
}

// StopMotionController stops the Motion Controller software on the drive.
func (client *Client) StopMotionController(ctx context.Context) error {
	return client.rtcManager.StopMotionController(ctx)
}

// StartMotionController restarts the Motion Controller software on the drive.
func (client *Client) StartMotionController(ctx context.Context) error {
	return client.rtcManager.StartMotionController(ctx)
}

// SaveCommandTableToFlash sends the SaveCommandTable (0x80) RTC command to persist
// the current RAM command table to flash. MC must be stopped before calling this.
// The LinMot drive may not send a standard RTC response for this command, so callers
// should expect a timeout and verify completion via read-back after recovery.
func (client *Client) SaveCommandTableToFlash(ctx context.Context) error {
	return client.rtcManager.SaveCommandTableToFlash(ctx)
}

// SetCommandTable sets the command table on the drive from a CommandTable struct.
func (client *Client) SetCommandTable(ctx context.Context, ct *client_command_tables.CommandTable) error {
	return client.rtcManager.SetCommandTable(ctx, ct)
}

// SetCommandTableWithOptions sets the command table on the drive from a CommandTable struct with configurable options.
func (client *Client) SetCommandTableWithOptions(ctx context.Context, ct *client_command_tables.CommandTable, opts client_command_tables.SetCommandTableOptions) error {
	return client.rtcManager.SetCommandTableWithOptions(ctx, ct, opts)
}

// Extended Parameter Operations //

// ReadRAM reads the current RAM value of a parameter by UPID.
func (client *Client) ReadRAM(ctx context.Context, upid uint16) (int32, error) {
	return client.rtcManager.ReadRAM(ctx, upid)
}

// WriteRAMAndROM writes both RAM and ROM value of a parameter by UPID.
// This is equivalent to calling both WriteRAM and WriteROM with the same value.
func (client *Client) WriteRAMAndROM(ctx context.Context, upid uint16, value int32) error {
	return client.rtcManager.WriteRAMAndROM(ctx, upid, value)
}

// GetParameterMinValue gets the minimum allowed value for a parameter by UPID.
func (client *Client) GetParameterMinValue(ctx context.Context, upid uint16) (int32, error) {
	return client.rtcManager.GetParameterMinValue(ctx, upid)
}

// GetParameterMaxValue gets the maximum allowed value for a parameter by UPID.
func (client *Client) GetParameterMaxValue(ctx context.Context, upid uint16) (int32, error) {
	return client.rtcManager.GetParameterMaxValue(ctx, upid)
}

// GetParameterDefaultValue gets the default value for a parameter by UPID.
func (client *Client) GetParameterDefaultValue(ctx context.Context, upid uint16) (int32, error) {
	return client.rtcManager.GetParameterDefaultValue(ctx, upid)
}

// UPID List Operations //

// GetAllParameterIDs returns a list of all available parameter UPIDs on the drive.
// This iterates through the entire UPID list starting from 0x0000.
func (client *Client) GetAllParameterIDs(ctx context.Context) ([]uint16, error) {
	return client.rtcManager.GetAllParameterIDs(ctx)
}

// GetModifiedParameterIDs returns a list of all modified parameter UPIDs.
// Modified parameters are those whose RAM or ROM values differ from factory defaults.
func (client *Client) GetModifiedParameterIDs(ctx context.Context) ([]uint16, error) {
	return client.rtcManager.GetModifiedParameterIDs(ctx)
}

// GetAllParameters returns all available parameters with address usage information.
// This provides richer metadata than GetAllParameterIDs by including whether each
// parameter is stored in RAM, ROM, or both.
func (client *Client) GetAllParameters(ctx context.Context) ([]client_parameters.ParameterInfo, error) {
	return client.rtcManager.GetAllParameters(ctx)
}

// GetModifiedParameters returns all modified parameters with address usage information.
// Modified parameters are those whose RAM or ROM values differ from factory defaults.
func (client *Client) GetModifiedParameters(ctx context.Context) ([]client_parameters.ParameterInfo, error) {
	return client.rtcManager.GetModifiedParameters(ctx)
}

// Drive Operations //

// RestartDrive restarts the drive.
// WARNING: This is a dangerous operation that will restart the entire drive.
// The drive will go through its boot sequence and all runtime state will be lost.
func (client *Client) RestartDrive(ctx context.Context) error {
	return client.rtcManager.RestartDrive(ctx)
}

// ResetOSParametersToDefault resets all Operating System software parameters to factory defaults.
// WARNING: This will reset OS-level configuration parameters to their default values.
func (client *Client) ResetOSParametersToDefault(ctx context.Context) error {
	return client.rtcManager.ResetOSParametersToDefault(ctx)
}

// ResetMCParametersToDefault resets all Motion Control parameters to factory defaults.
// WARNING: This will reset all motion control configuration to default values.
func (client *Client) ResetMCParametersToDefault(ctx context.Context) error {
	return client.rtcManager.ResetMCParametersToDefault(ctx)
}

// ResetInterfaceParametersToDefault resets all Interface parameters to factory defaults.
// WARNING: This will reset interface configuration (e.g., I/O, communication) to default values.
func (client *Client) ResetInterfaceParametersToDefault(ctx context.Context) error {
	return client.rtcManager.ResetInterfaceParametersToDefault(ctx)
}

// ResetApplicationParametersToDefault resets all Application parameters to factory defaults.
// WARNING: This will reset application-level configuration to default values.
func (client *Client) ResetApplicationParametersToDefault(ctx context.Context) error {
	return client.rtcManager.ResetApplicationParametersToDefault(ctx)
}

// Curve Service //

// SaveAllCurves saves all curves from RAM to Flash memory.
// This makes the current curves persistent across power cycles.
func (client *Client) SaveAllCurves(ctx context.Context) error {
	return client.rtcManager.SaveAllCurves(ctx)
}

// DeleteAllCurves deletes all curves from RAM.
// WARNING: This will remove all motion curves. Call SaveAllCurves first if you want to preserve them.
func (client *Client) DeleteAllCurves(ctx context.Context) error {
	return client.rtcManager.DeleteAllCurves(ctx)
}

// UploadCurve uploads a complete curve to the drive RAM.
// curveID is the curve number (1-100), infoBlock and dataBlock are the raw curve data.
func (client *Client) UploadCurve(ctx context.Context, curveID uint16, infoBlock, dataBlock []byte) error {
	return client.rtcManager.UploadCurve(ctx, curveID, infoBlock, dataBlock)
}

// DownloadCurve downloads a complete curve from the drive RAM.
// Returns the info block and data block as separate byte slices.
func (client *Client) DownloadCurve(ctx context.Context, curveID uint16) ([]byte, []byte, error) {
	return client.rtcManager.DownloadCurve(ctx, curveID)
}

// ModifyCurve modifies an existing curve in the drive RAM.
// curveID is the curve number, infoBlock and dataBlock are the updated curve data.
func (client *Client) ModifyCurve(ctx context.Context, curveID uint16, infoBlock, dataBlock []byte) error {
	return client.rtcManager.ModifyCurve(ctx, curveID, infoBlock, dataBlock)
}

// Error Log //

// GetErrorLog retrieves the complete error log from the drive.
// Returns a slice of error log entries with error codes and timestamps.
func (client *Client) GetErrorLog(ctx context.Context) ([]client_errors.ErrorLogEntry, error) {
	return client.rtcManager.GetErrorLog(ctx)
}

// GetErrorLogCounts returns the count of logged and occurred errors.
func (client *Client) GetErrorLogCounts(ctx context.Context) (logged, occurred uint16, err error) {
	return client.rtcManager.GetErrorLogCounts(ctx)
}

// GetErrorLogEntry retrieves a single error log entry by index.
func (client *Client) GetErrorLogEntry(ctx context.Context, entryNum uint16) (*client_errors.ErrorLogEntry, error) {
	return client.rtcManager.GetErrorLogEntry(ctx, entryNum)
}

// GetErrorText retrieves the human-readable description for an error code.
// Error text is retrieved in stringlets (4-byte chunks) up to a maximum of 32 bytes.
func (client *Client) GetErrorText(ctx context.Context, errorCode uint16) (string, error) {
	return client.rtcManager.GetErrorText(ctx, errorCode)
}

// GetErrorLogWithText retrieves the complete error log from the drive with human-readable descriptions.
// This is more expensive than GetErrorLog as it requires additional requests to fetch text for each error.
func (client *Client) GetErrorLogWithText(ctx context.Context) ([]client_errors.ErrorLogEntry, error) {
	return client.rtcManager.GetErrorLogWithText(ctx)
}

// State Machine Control //

// EnableDrive transitions the drive to Operation Enabled state.
// Sets the control word to enable switch on, voltage, quick stop release, and operation.
// Returns the drive status after the command is sent.
func (client *Client) EnableDrive(ctx context.Context) (*protocol_common.Status, error) {
	return client.controlWordManager.EnableDrive(ctx)
}

// DisableDrive transitions the drive to Switch On Disabled state.
// Clears all control word bits to disable the drive.
func (client *Client) DisableDrive(ctx context.Context) (*protocol_common.Status, error) {
	return client.controlWordManager.DisableDrive(ctx)
}

// AcknowledgeError acknowledges and clears a drive error.
// Sends a rising edge followed by a falling edge on the error acknowledge bit.
// Fails if a fatal error is present (fatal errors cannot be acknowledged).
func (client *Client) AcknowledgeError(ctx context.Context) (*protocol_common.Status, error) {
	return client.controlWordManager.AcknowledgeError(ctx)
}

// QuickStop triggers emergency quick stop.
// Immediately stops motion by clearing the quick stop bit (inverted logic).
func (client *Client) QuickStop(ctx context.Context) (*protocol_common.Status, error) {
	return client.controlWordManager.QuickStop(ctx)
}

// Home initiates the homing sequence.
// Sets the home bit in the control word to start the homing procedure.
// This operation may take up to 30 seconds to complete.
func (client *Client) Home(ctx context.Context) (*protocol_common.Status, error) {
	return client.controlWordManager.Home(ctx)
}

// GetDriveStatus queries the current drive status without changing the control word.
// This is a read-only operation that returns the current state machine state, position, and error status.
func (client *Client) GetDriveStatus(ctx context.Context) (*protocol_common.Status, error) {
	return client.controlWordManager.GetStatus(ctx)
}

// SendControlWord sends a raw control word value to the drive.
// This is an advanced method for direct control word manipulation.
// Most users should use the higher-level methods like EnableDrive, DisableDrive, etc.
func (client *Client) SendControlWord(ctx context.Context, word uint16) (*protocol_common.Status, error) {
	return client.controlWordManager.SendControlWord(ctx, word)
}

// ============================================================================
// Monitoring Channel Methods
// ============================================================================

// ConfigureMonitoringChannel configures a single monitoring channel to monitor a specific UPID.
// channelNum must be 1-4.
// upid is the parameter to monitor (e.g., position, velocity, current).
func (client *Client) ConfigureMonitoringChannel(ctx context.Context, channelNum int, upid uint16) error {
	return client.monitoringManager.ConfigureChannel(ctx, channelNum, upid)
}

// ConfigureMonitoringChannels bulk configures all 4 monitoring channels.
// upids[0] configures channel 1, upids[1] configures channel 2, etc.
func (client *Client) ConfigureMonitoringChannels(ctx context.Context, upids [4]uint16) error {
	return client.monitoringManager.ConfigureChannels(ctx, upids)
}

// GetMonitoringChannelConfiguration reads which UPID is assigned to a monitoring channel.
// channelNum must be 1-4.
func (client *Client) GetMonitoringChannelConfiguration(ctx context.Context, channelNum int) (uint16, error) {
	return client.monitoringManager.GetChannelConfiguration(ctx, channelNum)
}

// GetAllMonitoringChannelConfigurations returns the UPIDs configured for all 4 monitoring channels.
func (client *Client) GetAllMonitoringChannelConfigurations(ctx context.Context) ([4]uint16, error) {
	return client.monitoringManager.GetAllChannelConfigurations(ctx)
}

// GetMonitoringData retrieves drive status with monitoring channel data.
// Returns a Status struct with MonitoringChannel field populated.
func (client *Client) GetMonitoringData(ctx context.Context) (*protocol_common.Status, error) {
	return client.monitoringManager.GetMonitoringData(ctx)
}

// GetMonitoringSnapshot retrieves monitoring data as a convenient snapshot struct.
// This method retrieves status and monitoring channel values in one call.
func (client *Client) GetMonitoringSnapshot(ctx context.Context) (*client_monitoring.MonitoringSnapshot, error) {
	return client.monitoringManager.GetMonitoringSnapshot(ctx)
}

// ============================================================================
// Motion Control Methods (VAI Commands)
// ============================================================================

// VAIGoToPosition sends a VAI Go To Position command via Motion Control.
//
// Parameters:
//   - positionMM: Target position in millimeters
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
func (client *Client) VAIGoToPosition(ctx context.Context, positionMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return client.mcInterfaceManager.VAIGoToPosition(ctx, positionMM, velocityMS, accelMS2, decelMS2)
}

// VAIIncrementDemandPosition sends a VAI Increment Demand Position command via Motion Control.
//
// Parameters:
//   - incrementMM: Position increment in millimeters
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
func (client *Client) VAIIncrementDemandPosition(ctx context.Context, incrementMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return client.mcInterfaceManager.VAIIncrementDemandPosition(ctx, incrementMM, velocityMS, accelMS2, decelMS2)
}

// VAIIncrementTargetPosition sends a VAI Increment Target Position command via Motion Control.
//
// Parameters:
//   - incrementMM: Position increment in millimeters
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
func (client *Client) VAIIncrementTargetPosition(ctx context.Context, incrementMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return client.mcInterfaceManager.VAIIncrementTargetPosition(ctx, incrementMM, velocityMS, accelMS2, decelMS2)
}

// VAIStop sends a VAI Stop command via Motion Control.
func (client *Client) VAIStop(ctx context.Context) (*protocol_common.Status, error) {
	return client.mcInterfaceManager.VAIStop(ctx)
}

// VAIGoToPositionFromActual sends a VAI Go To Position From Actual Position command via Motion Control.
//
// Parameters:
//   - positionMM: Target position in millimeters
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
func (client *Client) VAIGoToPositionFromActual(ctx context.Context, positionMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return client.mcInterfaceManager.VAIGoToPositionFromActual(ctx, positionMM, velocityMS, accelMS2, decelMS2)
}

// VAIGoToPositionFromActualDemVelZero sends a VAI Go To Position From Actual with start velocity = 0 via Motion Control.
//
// Parameters:
//   - positionMM: Target position in millimeters
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
func (client *Client) VAIGoToPositionFromActualDemVelZero(ctx context.Context, positionMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return client.mcInterfaceManager.VAIGoToPositionFromActualDemVelZero(ctx, positionMM, velocityMS, accelMS2, decelMS2)
}

// VAIIncrementActualPosition sends a VAI Increment Actual Position command via Motion Control.
//
// Parameters:
//   - incrementMM: Position increment in millimeters
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
func (client *Client) VAIIncrementActualPosition(ctx context.Context, incrementMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return client.mcInterfaceManager.VAIIncrementActualPosition(ctx, incrementMM, velocityMS, accelMS2, decelMS2)
}

// VAIIncrementActualPositionDemVelZero sends a VAI Increment Actual Position with start velocity = 0 via Motion Control.
//
// Parameters:
//   - incrementMM: Position increment in millimeters
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
func (client *Client) VAIIncrementActualPositionDemVelZero(ctx context.Context, incrementMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return client.mcInterfaceManager.VAIIncrementActualPositionDemVelZero(ctx, incrementMM, velocityMS, accelMS2, decelMS2)
}

// SetRTCCounterForTesting sets the RTC counter value for testing purposes only.
// This should only be used in test code.
func (client *Client) SetRTCCounterForTesting(value uint8) {
	client.requestManager.SetRTCCounterForTesting(value)
}

// LastDriveCmdCountForTesting returns the last cmdCount seen from the drive for testing purposes only.
// This should only be used in test code.
func (client *Client) LastDriveCmdCountForTesting() uint8 {
	return client.requestManager.LastDriveCmdCountForTesting()
}

// PendingMCRequestCountForTesting returns the number of pending MC requests for testing purposes only.
// This should only be used in test code.
func (client *Client) PendingMCRequestCountForTesting() int {
	return client.requestManager.PendingMCRequestCountForTesting()
}

// VAIGoToPositionAfterActualCommand sends a VAI Go To Position After Actual Command via Motion Control.
//
// Parameters:
//   - positionMM: Target position in millimeters
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
func (client *Client) VAIGoToPositionAfterActualCommand(ctx context.Context, positionMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return client.mcInterfaceManager.VAIGoToPositionAfterActualCommand(ctx, positionMM, velocityMS, accelMS2, decelMS2)
}

// VAIGoToAnalogPosition sends a VAI Go To Analog Position command via Motion Control.
// Target position comes from analog input.
//
// Parameters:
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
func (client *Client) VAIGoToAnalogPosition(ctx context.Context, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return client.mcInterfaceManager.VAIGoToAnalogPosition(ctx, velocityMS, accelMS2, decelMS2)
}

// VAIGoToPositionOnRisingTrigger sends a VAI Go To Position On Rising Trigger Event via Motion Control.
//
// Parameters:
//   - positionMM: Target position in millimeters
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
func (client *Client) VAIGoToPositionOnRisingTrigger(ctx context.Context, positionMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return client.mcInterfaceManager.VAIGoToPositionOnRisingTrigger(ctx, positionMM, velocityMS, accelMS2, decelMS2)
}

// VAIIncrementTargetPositionOnRisingTrigger sends a VAI Increment Target Position On Rising Trigger Event via Motion Control.
//
// Parameters:
//   - incrementMM: Position increment in millimeters
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
func (client *Client) VAIIncrementTargetPositionOnRisingTrigger(ctx context.Context, incrementMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return client.mcInterfaceManager.VAIIncrementTargetPositionOnRisingTrigger(ctx, incrementMM, velocityMS, accelMS2, decelMS2)
}

// VAIGoToPositionOnFallingTrigger sends a VAI Go To Position On Falling Trigger Event via Motion Control.
//
// Parameters:
//   - positionMM: Target position in millimeters
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
func (client *Client) VAIGoToPositionOnFallingTrigger(ctx context.Context, positionMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return client.mcInterfaceManager.VAIGoToPositionOnFallingTrigger(ctx, positionMM, velocityMS, accelMS2, decelMS2)
}

// VAIIncrementTargetPositionOnFallingTrigger sends a VAI Increment Target Position On Falling Trigger Event via Motion Control.
//
// Parameters:
//   - incrementMM: Position increment in millimeters
//   - velocityMS: Maximal velocity in meters per second
//   - accelMS2: Acceleration in meters per second squared
//   - decelMS2: Deceleration in meters per second squared
func (client *Client) VAIIncrementTargetPositionOnFallingTrigger(ctx context.Context, incrementMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return client.mcInterfaceManager.VAIIncrementTargetPositionOnFallingTrigger(ctx, incrementMM, velocityMS, accelMS2, decelMS2)
}

// VAIChangeMotionParamsOnPositiveTransition sends a VAI Change Motion Parameters On Positive Position Transition via Motion Control.
//
// Parameters:
//   - transitionPosMM: Position where parameters change in millimeters
//   - velocityMS: Maximal velocity after event in meters per second
//   - accelMS2: Acceleration after event in meters per second squared
//   - decelMS2: Deceleration after event in meters per second squared
func (client *Client) VAIChangeMotionParamsOnPositiveTransition(ctx context.Context, transitionPosMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return client.mcInterfaceManager.VAIChangeMotionParamsOnPositiveTransition(ctx, transitionPosMM, velocityMS, accelMS2, decelMS2)
}

// VAIChangeMotionParamsOnNegativeTransition sends a VAI Change Motion Parameters On Negative Position Transition via Motion Control.
//
// Parameters:
//   - transitionPosMM: Position where parameters change in millimeters
//   - velocityMS: Maximal velocity after event in meters per second
//   - accelMS2: Acceleration after event in meters per second squared
//   - decelMS2: Deceleration after event in meters per second squared
func (client *Client) VAIChangeMotionParamsOnNegativeTransition(ctx context.Context, transitionPosMM, velocityMS, accelMS2, decelMS2 float64) (*protocol_common.Status, error) {
	return client.mcInterfaceManager.VAIChangeMotionParamsOnNegativeTransition(ctx, transitionPosMM, velocityMS, accelMS2, decelMS2)
}

// ============================================================================
// Interface Control Methods
// ============================================================================

// WriteOutputsWithMask sets digital outputs on the X4 interface with a bitmask.
// This allows direct control of digital outputs like vacuum control valves.
//
// Parameters:
//   - bitMask: Which bits to modify (bit 0 = X4.3, bit 1 = X4.4, etc.)
//   - bitValue: Values to set for masked bits
//
// Common usage for vacuum control:
//   - Vacuum ON:  WriteOutputsWithMask(ctx, 0x0003, 0x0002)  // bit 1 high
//   - Vacuum OFF: WriteOutputsWithMask(ctx, 0x0003, 0x0001)  // bit 0 high (purge)
//   - All OFF:    WriteOutputsWithMask(ctx, 0x0003, 0x0000)
func (client *Client) WriteOutputsWithMask(ctx context.Context, bitMask, bitValue uint16) (*protocol_common.Status, error) {
	return client.mcInterfaceManager.WriteOutputsWithMask(ctx, bitMask, bitValue)
}
