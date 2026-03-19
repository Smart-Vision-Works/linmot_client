package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "primer/proto"
	"primer/server"
)

type fakeStagePrimer struct {
	pb.UnimplementedStagePrimerServer
}

func (fakeStagePrimer) GetPosition(ctx context.Context, req *pb.GetPositionRequest) (*pb.GetPositionResponse, error) {
	return &pb.GetPositionResponse{
		RobotIndex: req.RobotIndex,
		StageIndex: req.StageIndex,
		Z:          12.34,
	}, nil
}

type errorStagePrimer struct {
	pb.UnimplementedStagePrimerServer
}

func (errorStagePrimer) GetPosition(context.Context, *pb.GetPositionRequest) (*pb.GetPositionResponse, error) {
	return nil, status.Error(codes.NotFound, "missing")
}

func TestGatewayGetPosition(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer lis.Close()

	grpcServer := grpc.NewServer()
	pb.RegisterStagePrimerServer(grpcServer, fakeStagePrimer{})
	go func() {
		_ = grpcServer.Serve(lis)
	}()
	defer grpcServer.Stop()

	gatewayMux, err := server.NewGatewayMux(context.Background(), lis.Addr().String())
	if err != nil {
		t.Fatalf("register gateway: %v", err)
	}

	httpServer, err := server.NewServer()
	if err != nil {
		t.Fatalf("new server: %v", err)
	}
	httpServer.MountGateway("/grpc", http.StripPrefix("/grpc", gatewayMux))

	ts := httptest.NewServer(httpServer.Handler())
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/grpc/stage_primer.StagePrimer/GetPosition", "application/json",
		bytes.NewBufferString(`{"robotIndex":1,"stageIndex":2}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}

	var payload struct {
		RobotIndex int32   `json:"robotIndex"`
		StageIndex int32   `json:"stageIndex"`
		Z          float64 `json:"z"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload.RobotIndex != 1 || payload.StageIndex != 2 || payload.Z != 12.34 {
		t.Fatalf("unexpected response: %+v", payload)
	}
}

func TestGatewayErrorShape(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer lis.Close()

	grpcServer := grpc.NewServer()
	pb.RegisterStagePrimerServer(grpcServer, errorStagePrimer{})
	go func() {
		_ = grpcServer.Serve(lis)
	}()
	defer grpcServer.Stop()

	gatewayMux, err := server.NewGatewayMux(context.Background(), lis.Addr().String())
	if err != nil {
		t.Fatalf("register gateway: %v", err)
	}

	httpServer, err := server.NewServer()
	if err != nil {
		t.Fatalf("new server: %v", err)
	}
	httpServer.MountGateway("/grpc", http.StripPrefix("/grpc", gatewayMux))

	ts := httptest.NewServer(httpServer.Handler())
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/grpc/stage_primer.StagePrimer/GetPosition", "application/json",
		bytes.NewBufferString(`{"robotIndex":1,"stageIndex":2}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("unexpected status: %d", resp.StatusCode)
	}

	var payload struct {
		Error string `json:"error"`
		Code  string `json:"code"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if payload.Error != "missing" || payload.Code != "NotFound" {
		t.Fatalf("unexpected error payload: %+v", payload)
	}
}
