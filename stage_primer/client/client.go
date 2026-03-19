package client

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "primer/proto"
)

// Global map of gRPC connections
var (
	connectionsMu sync.Mutex
	connections   = make(map[string]*grpc.ClientConn)
	stageClients  = make(map[string]pb.StagePrimerClient)
)

const (
	// Client timeouts for different operations.
	// Deploy cycle: write (~1s) + flash save (~75s observed) + recovery (~60s) + homing (~30s) = ~166s.
	// Use 330s to match the primer's 300s DeploymentTimeout with gRPC transport overhead.
	deployClientTimeout      = 330 * time.Second
	jogPositionClientTimeout = 15 * time.Second
	configClientTimeout      = 10 * time.Second
	vacuumClientTimeout      = 10 * time.Second
)

func getGrpcClient(ip string) (pb.StagePrimerClient, error) {
	connectionsMu.Lock()
	defer connectionsMu.Unlock()

	// Normalize IP: if no port, assume 50051 (gRPC port default)
	target := ip
	if !strings.Contains(ip, ":") {
		target = ip + ":50051"
	}

	// Check cache
	if client, ok := stageClients[target]; ok {
		return client, nil
	}

	// Connect
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to dial gRPC target %s", target)
	}

	client := pb.NewStagePrimerClient(conn)
	connections[target] = conn
	stageClients[target] = client

	return client, nil
}

// MonitorFaults connects to the fault stream and calls handler for every fault.
// This blocks until the stream is closed or an error occurs.
func MonitorFaults(ctx context.Context, primerIP string, handler func(ip string, err error)) error {
	client, err := getGrpcClient(primerIP)
	if err != nil {
		return err
	}

	stream, err := client.MonitorFaults(ctx, &pb.Empty{})
	if err != nil {
		return errors.Wrap(err, "failed to start fault monitor stream")
	}

	for {
		notification, err := stream.Recv()
		if err != nil {
			return err
		}

		text := formatFaultNotification(notification)
		sourceIP := strings.TrimSpace(notification.Ip)
		if sourceIP == "" {
			sourceIP = primerIP
		}
		handler(sourceIP, errors.New(text))
	}
}

func formatFaultNotification(notification *pb.FaultNotification) string {
	if notification == nil {
		return "Unknown fault"
	}

	parts := []string{}

	errorText := strings.TrimSpace(notification.ErrorText)
	if errorText != "" {
		parts = append(parts, fmt.Sprintf("error=%q", errorText))
	} else if notification.ErrorCode != 0 {
		parts = append(parts, fmt.Sprintf("error_code=0x%s", strings.ToUpper(strconv.FormatInt(int64(notification.ErrorCode), 16))))
	}

	warningText := strings.TrimSpace(notification.WarningText)
	if warningText != "" {
		parts = append(parts, fmt.Sprintf("warning=%q", warningText))
	}
	if notification.WarningWord != 0 {
		parts = append(parts, fmt.Sprintf("warning_word=0x%s", strings.ToUpper(strconv.FormatInt(int64(notification.WarningWord), 16))))
	}

	probeError := strings.TrimSpace(notification.ProbeError)
	if probeError != "" {
		parts = append(parts, fmt.Sprintf("probe_error=%q", probeError))
	}

	if notification.StatusWord != 0 {
		parts = append(parts, fmt.Sprintf("status_word=0x%s", strings.ToUpper(strconv.FormatInt(int64(notification.StatusWord), 16))))
	}
	if notification.StateVar != 0 {
		parts = append(parts, fmt.Sprintf("state_var=0x%s", strings.ToUpper(strconv.FormatInt(int64(notification.StateVar), 16))))
	}

	if len(parts) == 0 {
		return "Unknown fault"
	}
	return strings.Join(parts, " ")
}

// GetStagePosition retrieves the current Z position of a stage via stage primer.
func GetStagePosition(primerIP string, robotIndex, stageIndex int) (float64, error) {
	client, err := getGrpcClient(primerIP)
	if err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), jogPositionClientTimeout)
	defer cancel()

	resp, err := client.GetPosition(ctx, &pb.GetPositionRequest{
		RobotIndex: int32(robotIndex),
		StageIndex: int32(stageIndex),
	})
	if err != nil {
		return 0, errors.Wrap(err, "gRPC GetPosition failed")
	}

	return resp.Z, nil
}

// JogStage moves a stage to an absolute Z position via stage primer.
func JogStage(primerIP string, robotIndex, stageIndex int, z float64) error {
	client, err := getGrpcClient(primerIP)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), jogPositionClientTimeout)
	defer cancel()

	_, err = client.Jog(ctx, &pb.JogRequest{
		RobotIndex: int32(robotIndex),
		StageIndex: int32(stageIndex),
		Z:          z,
	})
	if err != nil {
		return errors.Wrap(err, "gRPC Jog failed")
	}

	return nil
}

// JogStageOffset moves a stage by a relative Z offset via stage primer.
func JogStageOffset(primerIP string, robotIndex, stageIndex int, zOffset float64) error {
	client, err := getGrpcClient(primerIP)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), jogPositionClientTimeout)
	defer cancel()

	_, err = client.JogOffset(ctx, &pb.JogOffsetRequest{
		RobotIndex: int32(robotIndex),
		StageIndex: int32(stageIndex),
		ZOffset:    zOffset,
	})
	if err != nil {
		return errors.Wrap(err, "gRPC JogOffset failed")
	}

	return nil
}

// DeployStageCommandTable deploys command table via stage primer gRPC endpoint.
func DeployStageCommandTable(primerIP string, robotIndex, stageIndex int, params CommandTableParams) error {
	client, err := getGrpcClient(primerIP)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), deployClientTimeout)
	defer cancel()

	_, err = client.DeployCommandTable(ctx, &pb.DeployCommandTableRequest{
		RobotIndex:          int32(robotIndex),
		StageIndex:          int32(stageIndex),
		ZDistance:            params.ZDistance,
		DefaultSpeed:         params.DefaultSpeed,
		DefaultAcceleration:  params.DefaultAcceleration,
		PickTime:             params.PickTime,
		LinmotIp:             params.LinMotIP,
	})
	if err != nil {
		return errors.Wrap(err, "gRPC DeployCommandTable failed")
	}

	return nil
}

// DeployInspectCommandTable deploys an inspect-mode command table via gRPC endpoint.
func DeployInspectCommandTable(primerIP string, robotIndex, stageIndex int, params CommandTableParams) error {
	client, err := getGrpcClient(primerIP)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), deployClientTimeout)
	defer cancel()

	_, err = client.DeployCommandTable(ctx, &pb.DeployCommandTableRequest{
		RobotIndex:          int32(robotIndex),
		StageIndex:          int32(stageIndex),
		ZDistance:            params.ZDistance,
		DefaultSpeed:         params.DefaultSpeed,
		DefaultAcceleration:  params.DefaultAcceleration,
		PickTime:             params.PickTime,
		InspectMode:          true,
		LinmotIp:             params.LinMotIP,
	})
	if err != nil {
		return errors.Wrap(err, "gRPC DeployInspectCommandTable failed")
	}

	return nil
}

// SetVacuum controls vacuum for a stage via stage primer gRPC endpoint.
func SetVacuum(primerIP string, robotIndex, stageIndex int, action string) error {
	client, err := getGrpcClient(primerIP)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), vacuumClientTimeout)
	defer cancel()

	_, err = client.SetVacuum(ctx, &pb.SetVacuumRequest{
		RobotIndex: int32(robotIndex),
		StageIndex: int32(stageIndex),
		Action:     action,
	})
	if err != nil {
		return errors.Wrap(err, "gRPC SetVacuum failed")
	}

	return nil
}

// GetUSBDevices lists connected USB devices via stage primer gRPC endpoint.
func GetUSBDevices(primerIP string) ([]USBDevice, error) {
	client, err := getGrpcClient(primerIP)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), configClientTimeout)
	defer cancel()

	resp, err := client.GetUSBDevices(ctx, &pb.Empty{})
	if err != nil {
		return nil, errors.Wrap(err, "gRPC GetUSBDevices failed")
	}

	devices := make([]USBDevice, 0, len(resp.Devices))
	for _, d := range resp.Devices {
		devices = append(devices, USBDevice{
			Bus:          d.Bus,
			Device:       d.Device,
			IDVendor:     d.IdVendor,
			IDProduct:    d.IdProduct,
			Manufacturer: d.Manufacturer,
			Product:      d.Product,
			Serial:       d.Serial,
		})
	}

	return devices, nil
}

// GetConfig retrieves the current stage_primer configuration via stage primer gRPC endpoint.
func GetConfig(primerIP string) (*StagePrimerConfig, error) {
	client, err := getGrpcClient(primerIP)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), configClientTimeout)
	defer cancel()

	resp, err := client.GetConfig(ctx, &pb.Empty{})
	if err != nil {
		return nil, errors.Wrap(err, "gRPC GetConfig failed")
	}

	return &StagePrimerConfig{
		ClearCores: mapProtoToConfig(resp.Clearcores),
	}, nil
}

// SetConfig updates the stage_primer configuration via stage primer gRPC endpoint.
func SetConfig(primerIP string, config StagePrimerConfig) error {
	client, err := getGrpcClient(primerIP)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), configClientTimeout)
	defer cancel()

	_, err = client.SetConfig(ctx, &pb.SetConfigRequest{
		Clearcores: mapConfigToProto(config.ClearCores),
	})
	if err != nil {
		return errors.Wrap(err, "gRPC SetConfig failed")
	}

	return nil
}

func mapConfigToProto(clearcores []ClearCoreConfig) []*pb.ClearCoreConfig {
	pbCCs := make([]*pb.ClearCoreConfig, 0, len(clearcores))
	for _, cc := range clearcores {
		linmots := make([]*pb.LinMotConfig, 0, len(cc.LinMots))
		for _, lm := range cc.LinMots {
			linmots = append(linmots, &pb.LinMotConfig{
				Ip: lm.IP,
			})
		}
		pbCCs = append(pbCCs, &pb.ClearCoreConfig{
			UsbId:                 cc.USBID,
			Dhcp:                  cc.DHCP,
			IpAddress:             cc.IPAddress,
			Gateway:               cc.Gateway,
			Subnet:                cc.Subnet,
			Dns:                   cc.DNS,
			RetransmissionTimeout: uint32(cc.RetransmissionTimeout),
			RetransmissionCount:   uint32(cc.RetransmissionCount),
			Linmots:               linmots,
		})
	}
	return pbCCs
}

func mapProtoToConfig(pbCCs []*pb.ClearCoreConfig) []ClearCoreConfig {
	clearcores := make([]ClearCoreConfig, 0, len(pbCCs))
	for _, pbCC := range pbCCs {
		linmots := make([]LinMotConfig, 0, len(pbCC.Linmots))
		for _, pbLM := range pbCC.Linmots {
			linmots = append(linmots, LinMotConfig{
				IP: pbLM.Ip,
			})
		}
		clearcores = append(clearcores, ClearCoreConfig{
			USBID:                 pbCC.UsbId,
			DHCP:                  pbCC.Dhcp,
			IPAddress:             pbCC.IpAddress,
			Gateway:               pbCC.Gateway,
			Subnet:                pbCC.Subnet,
			DNS:                   pbCC.Dns,
			RetransmissionTimeout: uint8(pbCC.RetransmissionTimeout),
			RetransmissionCount:   uint8(pbCC.RetransmissionCount),
			LinMots:               linmots,
		})
	}
	return clearcores
}

// containsPort checks if an address string contains a port.
func containsPort(addr string) bool {
	// Simple check: if it contains a colon and something after it, assume it has a port.
	for i := len(addr) - 1; i >= 0; i-- {
		if addr[i] == ':' {
			return i < len(addr)-1
		}
		// Stop at the first non-numeric character (IPv4) or bracket (IPv6).
		if addr[i] == ']' || addr[i] == '.' {
			break
		}
	}
	return false
}
