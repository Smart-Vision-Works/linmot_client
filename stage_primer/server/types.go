package server

// JogRequest represents a jog request
type JogRequest struct {
	RobotIndex int     `json:"robotIndex" binding:min=0"`
	StageIndex int     `json:"stageIndex" binding:min=0"`
	Z          float64 `json:"z"`
}

// JogOffsetRequest represents a jog offset request
type JogOffsetRequest struct {
	RobotIndex int     `json:"robotIndex" binding:min=0"`
	StageIndex int     `json:"stageIndex" binding:min=0"`
	ZOffset    float64 `json:"zOffset"`
}

// PositionRequest represents a position query request
type PositionRequest struct {
	RobotIndex int `form:"robotIndex" binding:"min=0"`
	StageIndex int `form:"stageIndex" binding:"min=0"`
}

// PositionResponse represents a position response
type PositionResponse struct {
	RobotIndex int     `json:"robotIndex"`
	StageIndex int     `json:"stageIndex"`
	Z          float64 `json:"z"`
}

// CommandTableDeployRequest represents a command table deployment request
type CommandTableDeployRequest struct {
	RobotIndex          int     `json:"robotIndex" binding:"min=0"`
	StageIndex          int     `json:"stageIndex" binding:"min=0"`
	ZDistance           float64 `json:"zDistance"`
	DefaultSpeed        float64 `json:"defaultSpeed"`
	DefaultAcceleration float64 `json:"defaultAcceleration"`
	PickTime            float64 `json:"pickTime"`
	InspectMode         bool    `json:"inspectMode"` // If true, deploy motion-only table (no vacuum)
}

// USBDevice represents a USB device
type USBDevice struct {
	Bus          string `json:"bus"`
	Device       string `json:"device"`
	IDVendor     string `json:"idVendor"`
	IDProduct    string `json:"idProduct"`
	Manufacturer string `json:"manufacturer"`
	Product      string `json:"product"`
	Serial       string `json:"serial"`
}
