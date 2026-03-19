package client_curves

import (
	"context"
	"fmt"

	client_common "github.com/Smart-Vision-Works/linmot_client/client/common"
	protocol_rtc "github.com/Smart-Vision-Works/linmot_client/protocol/rtc"
	protocol_curves "github.com/Smart-Vision-Works/linmot_client/protocol/rtc/curves"
)

type CurveManager struct {
	requestManager *client_common.RequestManager
}

func NewCurveManager(requestManager *client_common.RequestManager) *CurveManager {
	return &CurveManager{
		requestManager: requestManager,
	}
}

func packCurveChunk(chunk []byte) (uint16, uint16) {
	if len(chunk) > 4 {
		chunk = chunk[:4]
	}
	var low, high uint16
	if len(chunk) > 0 {
		low = uint16(chunk[0])
	}
	if len(chunk) > 1 {
		low |= uint16(chunk[1]) << 8
	}
	if len(chunk) > 2 {
		high = uint16(chunk[2])
	}
	if len(chunk) > 3 {
		high |= uint16(chunk[3]) << 8
	}
	return low, high
}

// SaveAllCurves saves all curves from RAM to Flash.
func (m *CurveManager) SaveAllCurves(ctx context.Context) error {
	request := protocol_curves.NewSaveAllCurvesRequest()
	_, err := client_common.SendRequestAndReceive[*protocol_curves.SaveAllCurvesResponse](m.requestManager, ctx, request)
	return err
}

// DeleteAllCurves deletes all curves from RAM.
func (m *CurveManager) DeleteAllCurves(ctx context.Context) error {
	request := protocol_curves.NewDeleteAllCurvesRequest()
	_, err := client_common.SendRequestAndReceive[*protocol_curves.DeleteAllCurvesResponse](m.requestManager, ctx, request)
	return err
}

// UploadCurve uploads a complete curve to the drive RAM.
// infoBlock and dataBlock are the raw curve data bytes.
func (m *CurveManager) UploadCurve(ctx context.Context, curveID uint16, infoBlock, dataBlock []byte) error {
	// Start adding curve
	startReq := protocol_curves.NewStartAddingCurveRequest(curveID)
	_, err := client_common.SendRequestAndReceive[*protocol_rtc.RTCSetParamResponse](m.requestManager, ctx, startReq)
	if err != nil {
		return err
	}

	// Send info block in 4-byte chunks
	for i := 0; i < len(infoBlock); i += 4 {
		dataLow, dataHigh := packCurveChunk(infoBlock[i:])

		req := protocol_curves.NewAddCurveInfoBlockRequest(curveID, dataLow, dataHigh)
		_, err := client_common.SendRequestAndReceive[*protocol_rtc.RTCSetParamResponse](m.requestManager, ctx, req)
		if err != nil {
			return err
		}
	}

	// Send data block in 4-byte chunks
	for i := 0; i < len(dataBlock); i += 4 {
		dataLow, dataHigh := packCurveChunk(dataBlock[i:])

		req := protocol_curves.NewAddCurveDataRequest(curveID, dataLow, dataHigh)
		resp, err := client_common.SendRequestAndReceive[*protocol_rtc.RTCSetParamResponse](m.requestManager, ctx, req)
		if err != nil {
			return err
		}

		// Check status - 0x04 means more data needed, 0x00 means complete
		if resp.RTCStatus() == 0x00 {
			// Upload complete
			break
		} else if resp.RTCStatus() != 0x04 {
			// Error status
			return fmt.Errorf("curve upload error: status 0x%02X", resp.RTCStatus())
		}
	}

	return nil
}

// DownloadCurve downloads a complete curve from the drive RAM.
// Returns the info block and data block as separate byte slices.
func (m *CurveManager) DownloadCurve(ctx context.Context, curveID uint16) (infoBlock, dataBlock []byte, err error) {
	// Start getting curve
	startReq := protocol_curves.NewStartGettingCurveRequest(curveID)
	_, err = client_common.SendRequestAndReceive[*protocol_rtc.RTCSetParamResponse](m.requestManager, ctx, startReq)
	if err != nil {
		return nil, nil, err
	}

	// Get info block
	getInfoReq := protocol_curves.NewGetCurveInfoBlockRequest(curveID)
	infoResp, err := client_common.SendRequestAndReceive[*protocol_rtc.RTCSetParamResponse](m.requestManager, ctx, getInfoReq)
	if err != nil {
		return nil, nil, err
	}

	// Extract info block size and data
	infoBlockSize := uint16(infoResp.Value() & 0xFFFF)
	infoBlock = make([]byte, 0, infoBlockSize)

	// For simplicity, assume info block is returned in value field
	// In practice, you'd need multiple requests for large info blocks
	val := uint32(infoResp.Value())
	infoBlock = append(infoBlock, byte(val&0xFF), byte((val>>8)&0xFF), byte((val>>16)&0xFF), byte((val>>24)&0xFF))

	// Get data block
	getDataReq := protocol_curves.NewGetCurveDataRequest(curveID)
	dataResp, err := client_common.SendRequestAndReceive[*protocol_rtc.RTCSetParamResponse](m.requestManager, ctx, getDataReq)
	if err != nil {
		return nil, nil, err
	}

	// Extract data block
	dataBlockSize := uint16(dataResp.Value() & 0xFFFF)
	dataBlock = make([]byte, 0, dataBlockSize)

	val = uint32(dataResp.Value())
	dataBlock = append(dataBlock, byte(val&0xFF), byte((val>>8)&0xFF), byte((val>>16)&0xFF), byte((val>>24)&0xFF))

	return infoBlock, dataBlock, nil
}

// ModifyCurve modifies an existing curve in the drive RAM.
// infoBlock and dataBlock are the raw curve data bytes.
func (m *CurveManager) ModifyCurve(ctx context.Context, curveID uint16, infoBlock, dataBlock []byte) error {
	// Start modifying curve
	startReq := protocol_curves.NewStartModifyingCurveRequest(curveID)
	_, err := client_common.SendRequestAndReceive[*protocol_rtc.RTCSetParamResponse](m.requestManager, ctx, startReq)
	if err != nil {
		return err
	}

	// Send info block in 4-byte chunks
	for i := 0; i < len(infoBlock); i += 4 {
		dataLow, dataHigh := packCurveChunk(infoBlock[i:])

		req := protocol_curves.NewModifyCurveInfoBlockRequest(curveID, dataLow, dataHigh)
		_, err := client_common.SendRequestAndReceive[*protocol_rtc.RTCSetParamResponse](m.requestManager, ctx, req)
		if err != nil {
			return err
		}
	}

	// Send data block in 4-byte chunks
	for i := 0; i < len(dataBlock); i += 4 {
		dataLow, dataHigh := packCurveChunk(dataBlock[i:])

		req := protocol_curves.NewModifyCurveDataRequest(curveID, dataLow, dataHigh)
		resp, err := client_common.SendRequestAndReceive[*protocol_rtc.RTCSetParamResponse](m.requestManager, ctx, req)
		if err != nil {
			return err
		}

		// Check status - 0x04 means more data needed, 0x00 means complete
		if resp.RTCStatus() == 0x00 {
			// Modify complete
			break
		} else if resp.RTCStatus() != 0x04 {
			// Error status
			return fmt.Errorf("curve modify error: status 0x%02X", resp.RTCStatus())
		}
	}

	return nil
}
