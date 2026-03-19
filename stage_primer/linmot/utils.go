package linmot

import (
	"fmt"

	"stage_primer_config"
)

// findLinMotIP finds the LinMot IP address for a given robotIndex and stageIndex from config
func findLinMotIP(robotIndex int, stageIndex int, config config.Config) (string, error) {
	// Validate robotIndex bounds
	if robotIndex >= len(config.ClearCores) {
		return "", fmt.Errorf("robot index %d out of range (has %d robots)",
			robotIndex, len(config.ClearCores))
	}

	// Get the ClearCore at the specified robotIndex
	ccConfig := config.ClearCores[robotIndex]

	// Validate stageIndex bounds
	if stageIndex >= len(ccConfig.LinMots) {
		return "", fmt.Errorf("stage index %d out of range for robot %d (has %d stages)",
			stageIndex, robotIndex, len(ccConfig.LinMots))
	}

	return ccConfig.LinMots[stageIndex].IP, nil
}
