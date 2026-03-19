package client_parameters

import (
	"context"
	"errors"
	"fmt"

	client_common "github.com/Smart-Vision-Works/staged_robot/client/common"
	protocol_common "github.com/Smart-Vision-Works/staged_robot/protocol/common"
	protocol_rtc "github.com/Smart-Vision-Works/staged_robot/protocol/rtc"
	protocol_parameters "github.com/Smart-Vision-Works/staged_robot/protocol/rtc/parameters"
)

type ParameterManager struct {
	requestManager *client_common.RequestManager
}

func NewParameterManager(requestManager *client_common.RequestManager) *ParameterManager {
	return &ParameterManager{
		requestManager: requestManager,
	}
}

// ParameterInfo contains UPID and its address usage information.
type ParameterInfo struct {
	UPID         uint16
	AddressUsage uint16 // Bit 0=RAM, Bit 1=ROM
}

var (
	ErrEndOfUPIDList    = errors.New("end of UPID list reached")
	ErrIteratorFinished = errors.New("iterator finished")
)

// UPIDIterator iterates over UPIDs in the drive's parameter list.
type UPIDIterator struct {
	manager   *ParameterManager
	ctx       context.Context
	started   bool
	finished  bool
	startUPID uint16
}

// ListParameters starts iteration over all UPIDs, starting from the given UPID.
func (m *ParameterManager) ListParameters(ctx context.Context, startUPID uint16) *UPIDIterator {
	return &UPIDIterator{
		manager:   m,
		ctx:       ctx,
		started:   false,
		startUPID: startUPID,
	}
}

// Next retrieves the next UPID in the list.
// Returns (upid, addressUsage, nil) on success.
// Returns (0, 0, io.EOF) when the list is exhausted.
// Returns (0, 0, error) on other errors.
func (iter *UPIDIterator) Next() (upid uint16, addressUsage uint16, err error) {
	if iter.finished {
		return 0, 0, ErrIteratorFinished
	}

	var nextResponse *protocol_parameters.GetNextUPIDListItemResponse

	if !iter.started {
		// First call - start the list
		startReq := protocol_parameters.NewStartGettingUPIDListRequest(iter.startUPID)
		_, err := client_common.SendRequestAndReceive[*protocol_parameters.StartGettingUPIDListResponse](iter.manager.requestManager, iter.ctx, startReq)
		if err != nil {
			iter.finished = true
			return 0, 0, err
		}
		iter.started = true

		// Get first item
		nextReq := protocol_parameters.NewGetNextUPIDListItemRequest()
		nextResponse, err = client_common.SendRequestAndReceive[*protocol_parameters.GetNextUPIDListItemResponse](iter.manager.requestManager, iter.ctx, nextReq)
		if err != nil {
			iter.finished = true
			return 0, 0, err
		}
	} else {
		// Subsequent calls - get next item
		nextReq := protocol_parameters.NewGetNextUPIDListItemRequest()
		var err2 error
		nextResponse, err2 = client_common.SendRequestAndReceive[*protocol_parameters.GetNextUPIDListItemResponse](iter.manager.requestManager, iter.ctx, nextReq)
		if err2 != nil {
			iter.finished = true
			return 0, 0, err2
		}
	}

	// Extract UPID and addressUsage from typed response
	foundUPID := nextResponse.FoundUPID()
	addressUsage = nextResponse.AddressUsage()

	// Check RTC status for end-of-list marker (0xC6)
	if nextResponse.RTCStatus() == 0xC6 {
		iter.finished = true
		return 0, 0, ErrEndOfUPIDList
	}

	// Check for other errors
	if protocol_rtc.IsStatusError(nextResponse.RTCStatus()) {
		iter.finished = true
		return 0, 0, fmt.Errorf("error reading UPID list: status 0x%02X", nextResponse.RTCStatus())
	}

	return foundUPID, addressUsage, nil
}

// GetParameterBounds retrieves the min, max, and default values for a parameter.
func (m *ParameterManager) GetParameterBounds(ctx context.Context, upid uint16) (min, max, def int32, err error) {
	// Get min value
	minReq := protocol_parameters.NewGetMinValueRequest(upid)
	minResp, err := client_common.SendRequestAndReceive[*protocol_parameters.GetMinValueResponse](m.requestManager, ctx, minReq)
	if err != nil {
		return 0, 0, 0, err
	}

	// Get max value
	maxReq := protocol_parameters.NewGetMaxValueRequest(upid)
	maxResp, err := client_common.SendRequestAndReceive[*protocol_parameters.GetMaxValueResponse](m.requestManager, ctx, maxReq)
	if err != nil {
		return 0, 0, 0, err
	}

	// Get default value
	defReq := protocol_parameters.NewGetDefaultValueRequest(upid)
	defResp, err := client_common.SendRequestAndReceive[*protocol_parameters.GetDefaultValueResponse](m.requestManager, ctx, defReq)
	if err != nil {
		return 0, 0, 0, err
	}

	return minResp.MinValue(), maxResp.MaxValue(), defResp.DefaultValue(), nil
}

// ReadRAM reads the current RAM value of a parameter by UPID.
func (m *ParameterManager) ReadRAM(ctx context.Context, upid uint16) (int32, error) {
	request := protocol_rtc.NewRTCGetParamRequest(upid, protocol_rtc.CommandCode.ReadRAM)
	response, err := client_common.SendRequestAndReceive[*protocol_rtc.RTCGetParamResponse](m.requestManager, ctx, request)
	if err != nil {
		return 0, err
	}
	return response.Value(), nil
}

// WriteRAMAndROM writes both the RAM and ROM value of a parameter.
func (m *ParameterManager) WriteRAMAndROM(ctx context.Context, upid uint16, value int32) error {
	request := protocol_parameters.NewWriteRAMAndROMRequest(upid, value)
	_, err := client_common.SendRequestAndReceive[*protocol_rtc.RTCSetParamResponse](m.requestManager, ctx, request)
	return err
}

// GetMinValue gets the minimum allowed value for a parameter.
func (m *ParameterManager) GetMinValue(ctx context.Context, upid uint16) (int32, error) {
	request := protocol_parameters.NewGetMinValueRequest(upid)
	response, err := client_common.SendRequestAndReceive[*protocol_rtc.RTCGetParamResponse](m.requestManager, ctx, request)
	if err != nil {
		return 0, err
	}
	return response.Value(), nil
}

// GetMaxValue gets the maximum allowed value for a parameter.
func (m *ParameterManager) GetMaxValue(ctx context.Context, upid uint16) (int32, error) {
	request := protocol_parameters.NewGetMaxValueRequest(upid)
	response, err := client_common.SendRequestAndReceive[*protocol_rtc.RTCGetParamResponse](m.requestManager, ctx, request)
	if err != nil {
		return 0, err
	}
	return response.Value(), nil
}

// GetDefaultValue gets the default value for a parameter.
func (m *ParameterManager) GetDefaultValue(ctx context.Context, upid uint16) (int32, error) {
	request := protocol_parameters.NewGetDefaultValueRequest(upid)
	response, err := client_common.SendRequestAndReceive[*protocol_rtc.RTCGetParamResponse](m.requestManager, ctx, request)
	if err != nil {
		return 0, err
	}
	return response.Value(), nil
}

// GetAllUPIDs returns all available parameter UPIDs.
func (m *ParameterManager) GetAllUPIDs(ctx context.Context) ([]uint16, error) {
	iter := m.ListParameters(ctx, 0x0000)
	var upids []uint16
	for {
		upid, _, err := iter.Next()
		if err != nil {
			// Check for end of list
			if errors.Is(err, ErrEndOfUPIDList) || errors.Is(err, ErrIteratorFinished) {
				break
			}
			return nil, err
		}
		upids = append(upids, upid)
	}
	return upids, nil
}

// GetAllParameters returns all available parameters with address usage information.
func (m *ParameterManager) GetAllParameters(ctx context.Context) ([]ParameterInfo, error) {
	iter := m.ListParameters(ctx, 0x0000)
	var params []ParameterInfo
	for {
		upid, addressUsage, err := iter.Next()
		if err != nil {
			// Check for end of list
			if errors.Is(err, ErrEndOfUPIDList) || errors.Is(err, ErrIteratorFinished) {
				break
			}
			return nil, err
		}
		params = append(params, ParameterInfo{
			UPID:         upid,
			AddressUsage: addressUsage,
		})
	}
	return params, nil
}

// GetModifiedUPIDs returns all modified parameter UPIDs.
func (m *ParameterManager) GetModifiedUPIDs(ctx context.Context) ([]uint16, error) {
	// Start getting modified UPID list
	startReq := protocol_parameters.NewStartGettingModifiedUPIDListRequest(0x0000)
	_, err := client_common.SendRequestAndReceive[*protocol_parameters.StartGettingModifiedUPIDListResponse](m.requestManager, ctx, startReq)
	if err != nil {
		return nil, err
	}

	var upids []uint16
	for {
		// Get next modified UPID
		nextReq := protocol_parameters.NewGetNextModifiedUPIDListItemRequest()
		response, err := client_common.SendRequestAndReceive[*protocol_parameters.GetNextModifiedUPIDListItemResponse](m.requestManager, ctx, nextReq)
		if err != nil {
			return nil, err
		}

		// Check for end of list (status 0xC6)
		if response.RTCStatus() == 0xC6 {
			break
		}

		// Check for other errors
		if protocol_rtc.IsStatusError(response.RTCStatus()) {
			return nil, fmt.Errorf("error reading modified UPID list: status 0x%02X", response.RTCStatus())
		}

		foundUPID := response.FoundUPID()
		upids = append(upids, foundUPID)
	}

	return upids, nil
}

// GetModifiedParameters returns all modified parameters with full data.
func (m *ParameterManager) GetModifiedParameters(ctx context.Context) ([]ParameterInfo, error) {
	// Start getting modified UPID list
	startReq := protocol_parameters.NewStartGettingModifiedUPIDListRequest(0x0000)
	_, err := client_common.SendRequestAndReceive[*protocol_parameters.StartGettingModifiedUPIDListResponse](m.requestManager, ctx, startReq)
	if err != nil {
		return nil, err
	}

	var params []ParameterInfo
	for {
		// Get next modified UPID
		nextReq := protocol_parameters.NewGetNextModifiedUPIDListItemRequest()
		response, err := client_common.SendRequestAndReceive[*protocol_parameters.GetNextModifiedUPIDListItemResponse](m.requestManager, ctx, nextReq)
		if err != nil {
			return nil, err
		}

		// Check for end of list (status 0xC6)
		if response.RTCStatus() == 0xC6 {
			break
		}

		// Check for other errors
		if protocol_rtc.IsStatusError(response.RTCStatus()) {
			return nil, fmt.Errorf("error reading modified UPID list: status 0x%02X", response.RTCStatus())
		}

		foundUPID := response.FoundUPID()
		// Modified UPID list returns the actual value in words 3 and 4
		// Use addressUsage from the value field
		addressUsage := uint16(response.Value() & 0xFFFF)

		params = append(params, ParameterInfo{
			UPID:         foundUPID,
			AddressUsage: addressUsage,
		})
	}

	return params, nil
}

// ============================================================================
// Motion Parameter Methods (RTC-based wrappers for motion-related UPIDs)
// ============================================================================

// SetPosition1 sets the target position for Position 1 (upid 0x145A) in millimeters.
func (m *ParameterManager) SetPosition1(ctx context.Context, positionMM float64, storageType protocol_common.ParameterStorageType) error {
	request, err := protocol_parameters.NewWritePosition1Request(positionMM, storageType)
	if err != nil {
		return err
	}
	return m.requestManager.SendRequest(ctx, request)
}

// SetPosition2 sets the target position for Position 2 (upid 0x145F) in millimeters.
func (m *ParameterManager) SetPosition2(ctx context.Context, positionMM float64, storageType protocol_common.ParameterStorageType) error {
	request, err := protocol_parameters.NewWritePosition2Request(positionMM, storageType)
	if err != nil {
		return err
	}
	return m.requestManager.SendRequest(ctx, request)
}

// GetVelocity returns the maximum velocity in meters per second.
func (m *ParameterManager) GetVelocity(ctx context.Context) (float64, error) {
	request := protocol_parameters.NewReadVelocityRequest()
	response, err := client_common.SendRequestAndReceive[*protocol_parameters.ReadVelocityResponse](m.requestManager, ctx, request)
	if err != nil {
		return 0, err
	}
	return response.VelocityMS(), nil
}

// SetVelocity sets the maximum velocity in meters per second.
func (m *ParameterManager) SetVelocity(ctx context.Context, velocityMS float64, storageType protocol_common.ParameterStorageType) error {
	request, err := protocol_parameters.NewWriteVelocityRequest(velocityMS, storageType)
	if err != nil {
		return err
	}
	return m.requestManager.SendRequest(ctx, request)
}

// GetAcceleration returns the acceleration in meters per second squared.
func (m *ParameterManager) GetAcceleration(ctx context.Context) (float64, error) {
	request := protocol_parameters.NewReadAccelerationRequest()
	response, err := client_common.SendRequestAndReceive[*protocol_parameters.ReadAccelerationResponse](m.requestManager, ctx, request)
	if err != nil {
		return 0, err
	}
	return response.AccelerationMS2(), nil
}

// SetAcceleration sets the acceleration in meters per second squared.
func (m *ParameterManager) SetAcceleration(ctx context.Context, accelMS2 float64, storageType protocol_common.ParameterStorageType) error {
	request, err := protocol_parameters.NewWriteAccelerationRequest(accelMS2, storageType)
	if err != nil {
		return err
	}
	return m.requestManager.SendRequest(ctx, request)
}

// GetDeceleration returns the deceleration in meters per second squared.
func (m *ParameterManager) GetDeceleration(ctx context.Context) (float64, error) {
	request := protocol_parameters.NewReadDecelerationRequest()
	response, err := client_common.SendRequestAndReceive[*protocol_parameters.ReadDecelerationResponse](m.requestManager, ctx, request)
	if err != nil {
		return 0, err
	}
	return response.DecelerationMS2(), nil
}

// SetDeceleration sets the deceleration in meters per second squared.
func (m *ParameterManager) SetDeceleration(ctx context.Context, decelMS2 float64, storageType protocol_common.ParameterStorageType) error {
	request, err := protocol_parameters.NewWriteDecelerationRequest(decelMS2, storageType)
	if err != nil {
		return err
	}
	return m.requestManager.SendRequest(ctx, request)
}

// SetEasyStepsAutoStart sets the Easy Steps auto start configuration (upid 0x30D4).
func (m *ParameterManager) SetEasyStepsAutoStart(ctx context.Context, value int32, storageType protocol_common.ParameterStorageType) error {
	request, err := protocol_parameters.NewWriteEasyStepsAutoStartRequest(value, storageType)
	if err != nil {
		return err
	}
	return m.requestManager.SendRequest(ctx, request)
}

// SetEasyStepsAutoHome sets the Easy Steps auto home configuration (upid 0x30D5).
func (m *ParameterManager) SetEasyStepsAutoHome(ctx context.Context, value int32, storageType protocol_common.ParameterStorageType) error {
	request, err := protocol_parameters.NewWriteEasyStepsAutoHomeRequest(value, storageType)
	if err != nil {
		return err
	}
	return m.requestManager.SendRequest(ctx, request)
}

// SetEasyStepsRisingEdge sets the Easy Steps rising edge action for an input pin.
func (m *ParameterManager) SetEasyStepsRisingEdge(ctx context.Context, inputNumber protocol_common.IOPinNumber, value int32, storageType protocol_common.ParameterStorageType) error {
	request, err := protocol_parameters.NewWriteEasyStepsRisingEdgeRequest(inputNumber, value, storageType)
	if err != nil {
		return err
	}
	return m.requestManager.SendRequest(ctx, request)
}

// SetEasyStepsIOMotionConfigCmd sets the Easy Steps IO motion config curve/CMD ID.
func (m *ParameterManager) SetEasyStepsIOMotionConfigCmd(ctx context.Context, inputNumber protocol_common.IOPinNumber, curveCmdID int32, storageType protocol_common.ParameterStorageType) error {
	request, err := protocol_parameters.NewWriteEasyStepsIOMotionConfigCmdRequest(inputNumber, curveCmdID, storageType)
	if err != nil {
		return err
	}
	return m.requestManager.SendRequest(ctx, request)
}

// SetOutputFunction sets the output pin function configuration.
func (m *ParameterManager) SetOutputFunction(ctx context.Context, outputNumber protocol_common.IOPinNumber, value int32, storageType protocol_common.ParameterStorageType) error {
	request, err := protocol_parameters.NewWriteOutputFunctionRequest(outputNumber, value, storageType)
	if err != nil {
		return err
	}
	return m.requestManager.SendRequest(ctx, request)
}

// SetInputFunction sets the input pin function configuration.
func (m *ParameterManager) SetInputFunction(ctx context.Context, inputNumber protocol_common.IOPinNumber, value int32, storageType protocol_common.ParameterStorageType) error {
	request, err := protocol_parameters.NewWriteInputFunctionRequest(inputNumber, value, storageType)
	if err != nil {
		return err
	}
	return m.requestManager.SendRequest(ctx, request)
}

// SetRunMode sets the run mode configuration (upid 0x1450).
func (m *ParameterManager) SetRunMode(ctx context.Context, mode protocol_common.RunMode, storageType protocol_common.ParameterStorageType) error {
	request, err := protocol_parameters.NewWriteRunModeRequest(mode, storageType)
	if err != nil {
		return err
	}
	return m.requestManager.SendRequest(ctx, request)
}

// SetTriggerMode sets the trigger mode configuration (upid 0x170C).
func (m *ParameterManager) SetTriggerMode(ctx context.Context, value int32, storageType protocol_common.ParameterStorageType) error {
	request, err := protocol_parameters.NewWriteTriggerModeRequest(value, storageType)
	if err != nil {
		return err
	}
	return m.requestManager.SendRequest(ctx, request)
}
