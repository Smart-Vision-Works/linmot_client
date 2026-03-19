package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/Smart-Vision-Works/staged_robot/client"

	"primer/linmot"
	pb "primer/proto"
	"stage_primer_config"
)

type GRPCServer struct {
	pb.UnimplementedStagePrimerServer
	configPath string
	store      *config.ConfigStore
	server     *grpc.Server
	mockMode   bool
}

func NewGRPCServer(cfg config.Config) *GRPCServer {
	return &GRPCServer{
		store: config.NewConfigStore(cfg),
	}
}

// SetMockMode enables or disables mock mode for testing.
// When enabled, USB device discovery returns fake devices.
func (s *GRPCServer) SetMockMode(enabled bool) {
	s.mockMode = enabled
}

// SetConfigPath sets the file path used when writing config via SetConfig.
func (s *GRPCServer) SetConfigPath(configPath string) {
	s.configPath = strings.TrimSpace(configPath)
}

// SetConfigStore replaces the server's config store with a shared instance so
// that config written via SetConfig is immediately visible to other store
// readers (e.g. the fault monitor) without any filesystem round-trip.
func (s *GRPCServer) SetConfigStore(store *config.ConfigStore) {
	s.store = store
}

func (s *GRPCServer) loadConfig() (config.Config, error) {
	return s.store.Get()
}

func (s *GRPCServer) Start(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %v", port, err)
	}

	s.server = grpc.NewServer()
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(s.server, healthServer)
	pb.RegisterStagePrimerServer(s.server, s)
	if isDevelopmentEnv(os.Getenv("SERVER_ENVIRONMENT")) {
		reflection.Register(s.server)
	}

	fmt.Printf("Starting gRPC server on %s\n", port)
	return s.server.Serve(lis)
}

func (s *GRPCServer) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
	}
}

func isDevelopmentEnv(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "development", "dev", "local":
		return true
	default:
		return false
	}
}

func (s *GRPCServer) Jog(ctx context.Context, req *pb.JogRequest) (*pb.Empty, error) {
	fmt.Printf("[gRPC] Jog(robot=%d, stage=%d, z=%.3fmm)\n", req.RobotIndex, req.StageIndex, req.Z)
	currentConfig, err := s.loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	jogConfig := linmot.JogConfig{
		RobotIndex: int(req.RobotIndex),
		StageIndex: int(req.StageIndex),
		Config:     currentConfig,
		Position:   req.Z,
	}

	if err := linmot.Jog(ctx, jogConfig); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *GRPCServer) JogOffset(ctx context.Context, req *pb.JogOffsetRequest) (*pb.Empty, error) {
	fmt.Printf("[gRPC] JogOffset(robot=%d, stage=%d, offset=%.3fmm)\n", req.RobotIndex, req.StageIndex, req.ZOffset)
	currentConfig, err := s.loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	// Offset jog requires getting current position first
	posConfig := linmot.PositionConfig{
		RobotIndex: int(req.RobotIndex),
		StageIndex: int(req.StageIndex),
		Config:     currentConfig,
	}

	currentPos, err := linmot.GetPosition(ctx, posConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get current position for offset jog: %w", err)
	}

	jogConfig := linmot.JogConfig{
		RobotIndex: int(req.RobotIndex),
		StageIndex: int(req.StageIndex),
		Config:     currentConfig,
		Position:   currentPos + req.ZOffset,
	}

	if err := linmot.Jog(ctx, jogConfig); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *GRPCServer) GetPosition(ctx context.Context, req *pb.GetPositionRequest) (*pb.GetPositionResponse, error) {
	currentConfig, err := s.loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	posConfig := linmot.PositionConfig{
		RobotIndex: int(req.RobotIndex),
		StageIndex: int(req.StageIndex),
		Config:     currentConfig,
	}

	pos, err := linmot.GetPosition(ctx, posConfig)
	if err != nil {
		return nil, err
	}

	return &pb.GetPositionResponse{
		RobotIndex: req.RobotIndex,
		StageIndex: req.StageIndex,
		Z:          pos,
	}, nil
}

func (s *GRPCServer) DeployCommandTable(ctx context.Context, req *pb.DeployCommandTableRequest) (*pb.Empty, error) {
	fmt.Printf("[gRPC] DeployCommandTable(robot=%d, stage=%d, Z=%.2fmm, speed=%.1f%%, accel=%.1f%%, pickTime=%.3fs, inspect=%v, linmot_ip=%q)\n",
		req.RobotIndex, req.StageIndex, req.ZDistance, req.DefaultSpeed, req.DefaultAcceleration, req.PickTime, req.InspectMode, req.LinmotIp)

	// Load config as fallback for LinMot IP resolution (only needed if linmot_ip not provided)
	currentConfig, err := s.loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	deployConfig := linmot.DeployConfig{
		RobotIndex:          int(req.RobotIndex),
		StageIndex:          int(req.StageIndex),
		Config:              currentConfig,
		LinMotIP:            req.LinmotIp,
		ZDistance:            req.ZDistance,
		DefaultSpeed:        req.DefaultSpeed,
		DefaultAcceleration: req.DefaultAcceleration,
		PickTime:            req.PickTime,
	}

	// Deploy appropriate command table based on mode
	var deployErr error
	if req.InspectMode {
		deployErr = linmot.DeployInspectCommandTable(ctx, deployConfig)
	} else {
		deployErr = linmot.DeployCommandTable(ctx, deployConfig)
	}

	if deployErr != nil {
		return nil, deployErr
	}

	return &pb.Empty{}, nil
}

func (s *GRPCServer) SetVacuum(ctx context.Context, req *pb.SetVacuumRequest) (*pb.Empty, error) {
	fmt.Printf("[gRPC] SetVacuum(robot=%d, stage=%d, action=%v)\n", req.RobotIndex, req.StageIndex, req.Action)
	currentConfig, err := s.loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	vacuumConfig := linmot.VacuumConfig{
		RobotIndex: int(req.RobotIndex),
		StageIndex: int(req.StageIndex),
		Config:     currentConfig,
		Action:     linmot.VacuumAction(req.Action),
	}

	if err := linmot.SetVacuum(ctx, vacuumConfig); err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (s *GRPCServer) GetUSBDevices(ctx context.Context, req *pb.Empty) (*pb.GetUSBDevicesResponse, error) {
	var devices []USBDevice
	var err error

	if s.mockMode {
		devices = getMockUSBDevices()
	} else {
		devices, err = getUSBDevicesFunc()
		if err != nil {
			return nil, fmt.Errorf("failed to list USB devices: %w", err)
		}
	}

	pbDevices := make([]*pb.USBDevice, 0, len(devices))
	for _, d := range devices {
		pbDevices = append(pbDevices, &pb.USBDevice{
			Bus:          d.Bus,
			Device:       d.Device,
			IdVendor:     d.IDVendor,
			IdProduct:    d.IDProduct,
			Manufacturer: d.Manufacturer,
			Product:      d.Product,
			Serial:       d.Serial,
		})
	}

	return &pb.GetUSBDevicesResponse{
		Devices: pbDevices,
	}, nil
}

func (s *GRPCServer) GetConfig(ctx context.Context, req *pb.Empty) (*pb.GetConfigResponse, error) {
	currentConfig, err := s.loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &pb.GetConfigResponse{
		Clearcores: mapConfigToProto(currentConfig),
	}, nil
}

func (s *GRPCServer) SetConfig(ctx context.Context, req *pb.SetConfigRequest) (*pb.Empty, error) {
	fmt.Printf("[gRPC] SetConfig(%d ClearCore(s))\n", len(req.Clearcores))
	newConfig := mapProtoToConfig(req.Clearcores)

	// Validate the configuration
	if err := validateConfig(&newConfig); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	// Detect new LinMot IPs that need Setup (RunMode, trigger config, IO config).
	// Setup writes to ROM so it only needs to happen once per LinMot, but we must
	// ensure it happens for any newly-added LinMot before command tables are deployed.
	oldConfig, _ := s.loadConfig()
	newLinMotIPs := findNewLinMotIPs(oldConfig, newConfig)

	// Write the validated config to the persistent file first.
	if err := saveConfigToFile(s.configPath, &newConfig); err != nil {
		return nil, fmt.Errorf("failed to save config: %w", err)
	}

	// Update the in-memory store so all readers see the change immediately,
	// without any further filesystem access.
	s.store.Set(newConfig)

	// Run Setup for any new LinMot IPs (non-blocking — log but don't fail SetConfig)
	// Run Setup for any new LinMot IPs (non-blocking — log but don't fail SetConfig)
	if len(newLinMotIPs) > 0 {
		fmt.Printf("[gRPC] SetConfig: %d new LinMot IP(s) detected, running Setup: %v\n", len(newLinMotIPs), newLinMotIPs)
		setupConfigs := make([]linmot.SetupConfig, 0, len(newLinMotIPs))
		for _, ip := range newLinMotIPs {
			setupConfigs = append(setupConfigs, linmot.SetupConfig{IP: ip})
		}
		if err := linmot.SetupAll(ctx, setupConfigs); err != nil {
			fmt.Printf("[gRPC] SetConfig: WARNING — LinMot setup failed for new IPs: %v\n", err)
		}
	}

	return &pb.Empty{}, nil
}

// findNewLinMotIPs returns LinMot IPs that are in newConfig but not in oldConfig.
func findNewLinMotIPs(oldConfig, newConfig config.Config) []string {
	oldIPs := make(map[string]struct{})
	for _, cc := range oldConfig.ClearCores {
		for _, lm := range cc.LinMots {
			if lm.IP != "" {
				oldIPs[lm.IP] = struct{}{}
			}
		}
	}

	var newIPs []string
	for _, cc := range newConfig.ClearCores {
		for _, lm := range cc.LinMots {
			if lm.IP == "" {
				continue
			}
			if _, exists := oldIPs[lm.IP]; !exists {
				newIPs = append(newIPs, lm.IP)
			}
		}
	}
	return newIPs
}

func mapConfigToProto(cfg config.Config) []*pb.ClearCoreConfig {
	pbCCs := make([]*pb.ClearCoreConfig, 0, len(cfg.ClearCores))
	for _, cc := range cfg.ClearCores {
		pbLinMots := make([]*pb.LinMotConfig, 0, len(cc.LinMots))
		for _, lm := range cc.LinMots {
			pbLinMots = append(pbLinMots, &pb.LinMotConfig{
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
			Linmots:               pbLinMots,
		})
	}
	return pbCCs
}

func mapProtoToConfig(pbCCs []*pb.ClearCoreConfig) config.Config {
	cfg := config.Config{
		ClearCores: make([]config.ClearCoreConfig, 0, len(pbCCs)),
	}
	for _, pbCC := range pbCCs {
		linmots := make([]config.LinMotConfig, 0, len(pbCC.Linmots))
		for _, pbLM := range pbCC.Linmots {
			linmots = append(linmots, config.LinMotConfig{
				IP: pbLM.Ip,
			})
		}

		cfg.ClearCores = append(cfg.ClearCores, config.ClearCoreConfig{
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
	return cfg
}

func (s *GRPCServer) MonitorFaults(req *pb.Empty, stream pb.StagePrimer_MonitorFaultsServer) error {
	faultChan := make(chan struct {
		ip  string
		err error
	}, 10)

	removeListener := linmot.AddFaultListener(func(ip string, err error) {
		select {
		case faultChan <- struct {
			ip  string
			err error
		}{ip, err}:
		default:
			// Buffer full, drop message or handle overflow
		}
	})
	defer removeListener()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case f := <-faultChan:
			notification := &pb.FaultNotification{
				Ip: f.ip,
			}

			var fault *client.DriveFaultError
			if errors.As(f.err, &fault) {
				notification.StatusWord = int32(fault.StatusWord)
				notification.StateVar = int32(fault.StateVar)
				notification.ErrorCode = int32(fault.ErrorCode)
				notification.ErrorText = fault.ErrorText
				notification.WarningWord = int32(fault.WarningWord)
				notification.WarningText = fault.WarningText

				var probeErr *client.DriveFaultProbeError
				if errors.As(f.err, &probeErr) && probeErr.ProbeErr != nil {
					notification.ProbeError = probeErr.ProbeErr.Error()
				}
			} else {
				notification.ErrorText = f.err.Error()
			}

			if err := stream.Send(notification); err != nil {
				return err
			}
		}
	}
}
