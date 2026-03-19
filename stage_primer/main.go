package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	config "stage_primer_config"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"gsail-go/logger"
	"primer/linmot"
	"primer/server"
)

const logDir = "/logs"

func printUsage() {
	fmt.Printf("Stage Primer Server Tool\n\n")
	fmt.Printf("Usage: %s <command> [options]\n\n", os.Args[0])
	fmt.Printf("Commands:\n")
	fmt.Printf("  startup     Server startup: REST API + LinMot setup\n")
	fmt.Printf("  server      Start REST API server for LinMot control\n")
	fmt.Printf("  linmot      LinMot drive management (setup, setpos, getpos)\n")
	fmt.Printf("  help        Show this help message\n\n")
	fmt.Printf("Most common usage: %s startup\n", os.Args[0])
	fmt.Printf("Use '%s <command> -h' for command-specific help\n", os.Args[0])
}

// startServer starts the REST API and gRPC servers with the given configuration
// This is the core server logic without command-line argument parsing
func startServer(httpPort, grpcPort, configFile string, sigChan chan os.Signal, mockMode bool) error {
	fmt.Printf("Stage Primer Server\n")
	fmt.Printf("HTTP Port: %s\n", httpPort)
	fmt.Printf("gRPC Port: %s\n", grpcPort)
	fmt.Printf("Config File: %s\n", configFile)
	if mockMode {
		fmt.Printf("Mock Mode: enabled\n")
	}
	fmt.Printf("\n")

	// Load configuration from file
	cfg, err := config.LoadConfigFromFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %v", err)
	}

	// Create an in-memory config store shared by the gRPC server and the fault
	// monitor. This eliminates per-tick and per-RPC filesystem reads: config is
	// read from disk once at startup and updated in-memory whenever SetConfig is
	// called via gRPC.
	store := config.NewConfigStore(cfg)

	// Create and start server
	srv, err := server.NewServer(configFile)
	if err != nil {
		return fmt.Errorf("failed to create server: %v", err)
	}

	// Enable mock mode if requested
	if mockMode {
		srv.SetMockMode(true)
	}

	grpcSrv := server.NewGRPCServer(cfg)
	grpcSrv.SetConfigPath(configFile)
	grpcSrv.SetConfigStore(store)
	rootCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcEndpoint := fmt.Sprintf("127.0.0.1:%s", grpcPort)
	gatewayMux, err := server.NewGatewayMux(rootCtx, grpcEndpoint)
	if err != nil {
		return fmt.Errorf("failed to register grpc-gateway: %w", err)
	}
	srv.MountGateway("/grpc", http.StripPrefix("/grpc", gatewayMux))
	srv.SetGrpcHealthTarget(grpcEndpoint)

	group, groupCtx := errgroup.WithContext(rootCtx)

	if sigChan != nil {
		group.Go(func() error {
			select {
			case <-sigChan:
				fmt.Printf("\nReceived signal, shutting down server gracefully...\n")
				cancel()
			case <-groupCtx.Done():
			}
			return nil
		})
	}

	group.Go(func() error {
		return linmot.MonitorFaultsWithConfigProvider(groupCtx, store.Get)
	})

	group.Go(func() error {
		addr := fmt.Sprintf(":%s", httpPort)
		fmt.Printf("Starting HTTP server on %s\n", addr)
		if err := srv.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	group.Go(func() error {
		addr := fmt.Sprintf(":%s", grpcPort)
		if err := grpcSrv.Start(addr); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			return err
		}
		return nil
	})

	group.Go(func() error {
		<-groupCtx.Done()
		grpcSrv.Stop()
		if err := srv.Shutdown(); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
		return nil
	})

	if err := group.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func runStartup() {
	serverCmd := flag.NewFlagSet("startup", flag.ExitOnError)
	configFile := serverCmd.String("config", config.DefaultConfigPath, "Path to JSON configuration file")
	port := serverCmd.String("port", "80", "Port for REST API server")

	serverCmd.Usage = func() {
		fmt.Printf("Server Startup — REST API + LinMot setup (no ClearCore hardware)\n")
		fmt.Printf("Designed to run as a standalone Docker service without USB access.\n\n")
		fmt.Printf("Usage: %s startup [options]\n\n", os.Args[0])
		fmt.Printf("Options:\n")
		serverCmd.PrintDefaults()
		fmt.Printf("\nExamples:\n")
		fmt.Printf("  %s startup                           # Start on port 80\n", os.Args[0])
		fmt.Printf("  %s startup -port 9090                # Custom port\n", os.Args[0])
		fmt.Printf("  %s startup -config /path/config.json # Custom config\n", os.Args[0])
	}

	serverCmd.Parse(os.Args[2:])

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start REST API and gRPC servers in background
	go func() {
		if err := startServer(*port, "50051", *configFile, nil, false); err != nil {
			log.Printf("Server failed to start: %v", err)
		}
	}()

	// Load configuration for LinMot setup
	cfg, err := config.LoadConfigFromFile(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup LinMot drives if configured
	linmot.SetupLinMots(cfg)

	// Wait for signal
	<-sigChan
	fmt.Printf("\nReceived signal, shutting down gracefully...\n")
	fmt.Printf("Server stopped. Goodbye!\n")
}

func runServer() {
	serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
	port := serverCmd.String("port", "80", "Port to run the HTTP server on")
	grpcPort := serverCmd.String("grpc-port", "50051", "Port to run the gRPC server on")
	configFile := serverCmd.String("config", config.DefaultConfigPath, "Path to JSON configuration file")
	mockLinmot := serverCmd.Bool("mock-linmot", false, "Use mock LinMot clients for testing (no real hardware)")

	serverCmd.Usage = func() {
		fmt.Printf("Start REST API and gRPC server for LinMot control\n")
		fmt.Printf("Provides HTTP endpoints and gRPC service for configuration and jogging control\n\n")
		fmt.Printf("Usage: %s server [options]\n\n", os.Args[0])
		fmt.Printf("Options:\n")
		serverCmd.PrintDefaults()
		fmt.Printf("\nExamples:\n")
		fmt.Printf("  %s server                      # Start server on default ports\n", os.Args[0])
		fmt.Printf("  %s server -port 9090          # Start HTTP server on port 9090\n", os.Args[0])
		fmt.Printf("  %s server -config /path/to/config.json  # Use custom config file\n", os.Args[0])
		fmt.Printf("  %s server -mock-linmot        # Start with mock LinMot for testing\n", os.Args[0])
	}

	serverCmd.Parse(os.Args[2:])

	// Setup mock LinMot factory if requested
	var mockFactory *linmot.MockClientFactory
	if *mockLinmot {
		log.Println("Starting in mock LinMot mode - no real hardware connections")
		mockFactory = linmot.NewMockClientFactory()
		linmot.SetClientFactory(mockFactory)
		defer func() {
			log.Println("Shutting down mock LinMot factory")
			mockFactory.Close()
			linmot.ResetClientFactory()
		}()
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the server with signal handling
	if err := startServer(*port, *grpcPort, *configFile, sigChan, *mockLinmot); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func main() {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}
	logger.InitGlobalLogger(filepath.Join(logDir, "clearcore.log"), nil, true)

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "startup":
		runStartup()
	case "linmot":
		linmot.RunLinmot()
	case "server":
		runServer()
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}
