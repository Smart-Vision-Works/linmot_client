package linmot

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"stage_primer_config"
)

// RunLinmot handles the linmot command and its subcommands
func RunLinmot() {
	if len(os.Args) < 3 {
		printLinmotUsage()
		os.Exit(1)
	}

	subcommand := os.Args[2]

	switch subcommand {
	case "setup":
		runLinmotSetup()
	case "setpos":
		runLinmotSetPos()
	case "getpos":
		runLinmotGetPos()
	default:
		fmt.Printf("Unknown linmot subcommand: %s\n\n", subcommand)
		printLinmotUsage()
		os.Exit(1)
	}
}

func printLinmotUsage() {
	fmt.Printf("LinMot Drive Management Tool\n\n")
	fmt.Printf("Usage: %s linmot <subcommand> [options]\n\n", os.Args[0])
	fmt.Printf("Subcommands:\n")
	fmt.Printf("  setup    Set up LinMot drives (configure for triggered operation)\n")
	fmt.Printf("  setpos   Set Z-axis position on a LinMot drive\n")
	fmt.Printf("  getpos   Get current Z-axis position from a LinMot drive\n\n")
	fmt.Printf("Use '%s linmot <subcommand> -h' for subcommand-specific help\n", os.Args[0])
}

func runLinmotSetup() {
	// Define setup-specific flags
	setupCmd := flag.NewFlagSet("linmot setup", flag.ExitOnError)
	configFile := setupCmd.String("config", config.DefaultConfigPath, "Path to JSON configuration file")

	setupCmd.Usage = func() {
		fmt.Printf("Set up LinMot drives\n")
		fmt.Printf("Configures LinMot drives for triggered command table operation\n\n")
		fmt.Printf("Usage: %s linmot setup [options]\n\n", os.Args[0])
		fmt.Printf("Options:\n")
		setupCmd.PrintDefaults()
		fmt.Printf("\nExamples:\n")
		fmt.Printf("  %s linmot setup                    # Set up all LinMot drives from default config\n", os.Args[0])
		fmt.Printf("  %s linmot setup -config /path/to/config.json  # Use custom config file\n", os.Args[0])
	}

	setupCmd.Parse(os.Args[3:])

	fmt.Printf("LinMot Setup Tool\n")
	fmt.Printf("Config File: %s\n\n", *configFile)

	// Load configuration from file
	cfg, err := config.LoadConfigFromFile(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup LinMot drives
	SetupLinMots(cfg)

	fmt.Printf("LinMot setup complete. Exiting.\n")
}

func runLinmotSetPos() {
	// Define setpos-specific flags
	setPosCmd := flag.NewFlagSet("linmot setpos", flag.ExitOnError)
	configFile := setPosCmd.String("config", config.DefaultConfigPath, "Path to JSON configuration file")
	robotIndex := setPosCmd.Int("robot", 0, "Robot index in configuration (required)")
	stageIndex := setPosCmd.Int("stage", 0, "Stage index within robot (required)")
	position := setPosCmd.Float64("position", 0, "Z-axis position in mm to set")

	setPosCmd.Usage = func() {
		fmt.Printf("Set Z-axis Position Tool\n")
		fmt.Printf("Send Z-axis positioning commands to LinMot drives using configuration\n\n")
		fmt.Printf("Usage: %s linmot setpos -robot <robotIndex> -stage <stageIndex> -position <position>\n\n", os.Args[0])
		setPosCmd.PrintDefaults()
		fmt.Printf("\nExamples:\n")
		fmt.Printf("  %s linmot setpos -robot 0 -stage 0 -position 50.0    # Set robot 0, stage 0 to 50.0mm\n", os.Args[0])
		fmt.Printf("  %s linmot setpos -robot 1 -stage 0 -position 25.5    # Set robot 1, stage 0 to 25.5mm\n", os.Args[0])
		fmt.Printf("  %s linmot setpos -robot 0 -stage 1 -position 0.0     # Set robot 0, stage 1 to home position\n", os.Args[0])
	}

	setPosCmd.Parse(os.Args[3:])

	fmt.Printf("LinMot Set Position Command\n")
	fmt.Printf("Config File: %s\n", *configFile)
	fmt.Printf("Robot Index: %d\n", *robotIndex)
	fmt.Printf("Stage Index: %d\n", *stageIndex)
	fmt.Printf("Position: %.2f mm\n\n", *position)

	// Load configuration from file
	cfg, err := config.LoadConfigFromFile(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Send set position command to LinMot
	jogConfig := JogConfig{
		RobotIndex: *robotIndex,
		StageIndex: *stageIndex,
		Config:     cfg,
		Position:   *position,
	}

	ctx := context.Background()
	if err := Jog(ctx, jogConfig); err != nil {
		log.Fatalf("Failed to set LinMot position: %v", err)
	}

	fmt.Printf("Position command completed successfully\n")
}

func runLinmotGetPos() {
	// Define getpos-specific flags
	getPosCmd := flag.NewFlagSet("linmot getpos", flag.ExitOnError)
	configFile := getPosCmd.String("config", config.DefaultConfigPath, "Path to JSON configuration file")
	robotIndex := getPosCmd.Int("robot", 0, "Robot index in configuration (required)")
	stageIndex := getPosCmd.Int("stage", 0, "Stage index within robot (required)")

	getPosCmd.Usage = func() {
		fmt.Printf("Get Z-axis Position Tool\n")
		fmt.Printf("Get current Z-axis position from LinMot drive using configuration\n\n")
		fmt.Printf("Usage: %s linmot getpos -robot <robotIndex> -stage <stageIndex>\n\n", os.Args[0])
		getPosCmd.PrintDefaults()
		fmt.Printf("\nExamples:\n")
		fmt.Printf("  %s linmot getpos -robot 0 -stage 0    # Get position from robot 0, stage 0\n", os.Args[0])
		fmt.Printf("  %s linmot getpos -robot 1 -stage 0    # Get position from robot 1, stage 0\n", os.Args[0])
	}

	getPosCmd.Parse(os.Args[3:])

	fmt.Printf("LinMot Get Position Command\n")
	fmt.Printf("Config File: %s\n", *configFile)
	fmt.Printf("Robot Index: %d\n", *robotIndex)
	fmt.Printf("Stage Index: %d\n\n", *stageIndex)

	// Load configuration from file
	cfg, err := config.LoadConfigFromFile(*configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Get position from LinMot
	posConfig := PositionConfig{
		RobotIndex: *robotIndex,
		StageIndex: *stageIndex,
		Config:     cfg,
	}

	ctx := context.Background()
	position, err := GetPosition(ctx, posConfig)
	if err != nil {
		log.Fatalf("Failed to get position from LinMot drive: %v", err)
	}

	fmt.Printf("Position: %.2f mm\n", position)
}

// setupLinMots sets up all LinMot drives configured in the config
// This function needs to be accessible from main.go for the startup routine
func SetupLinMots(cfg config.Config) {
	allLinMots := cfg.GetAllLinMots()
	if len(allLinMots) == 0 {
		fmt.Printf("No LinMot drives configured\n")
		return
	}

	separator := strings.Repeat("=", 50)
	fmt.Printf("\n%s\n", separator)
	fmt.Printf("Setting up %d LinMot drive(s)...\n", len(allLinMots))
	fmt.Printf("%s\n\n", separator)

	setupConfigs := make([]SetupConfig, 0, len(allLinMots))
	for _, lm := range allLinMots {
		if lm.IP != "" {
			setupConfigs = append(setupConfigs, SetupConfig{IP: lm.IP})
		}
	}

	if len(setupConfigs) == 0 {
		fmt.Printf("No valid LinMot IP addresses found\n")
		return
	}

	ctx := context.Background()
	if err := SetupAll(ctx, setupConfigs); err != nil {
		log.Printf("Warning: LinMot setup failed: %v", err)
		return
	}

	fmt.Printf("\n%s\n", separator)
	fmt.Printf("LinMot setup completed successfully\n")
	fmt.Printf("%s\n\n", separator)
}
