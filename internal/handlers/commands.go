package handlers

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"
)

// CommandRequest represents the command execution request
// @Description Command execution request payload
type CommandRequest struct {
	// Command to execute
	// required: true
	Command string `json:"command" example:"ls -la"`
	// Arguments for the command
	// required: false
	Args []string `json:"args" example:"['-l', '-a']"`
}

// CommandResponse represents the command execution response
// @Description Command execution response payload
type CommandResponse struct {
	// Command output
	// required: true
	Output string `json:"output" example:"total 1234\ndrwxr-xr-x  ..."`
	// Exit code
	// required: true
	ExitCode int `json:"exit_code" example:"0"`
	// Error message if any
	// required: false
	Error string `json:"error,omitempty" example:"command not found"`
}

// HandleCommand implements a vulnerable command execution endpoint
// @Summary Execute system command
// @Description Execute system command (intentionally vulnerable to command injection)
// @Tags commands
// @Accept json
// @Produce json
// @Param request body CommandRequest true "Command to execute"
// @Security BearerAuth
// @Success 200 {object} CommandResponse
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Server error"
// @Router /commands [post]
func HandleCommand(w http.ResponseWriter, r *http.Request) {
	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Vulnerable: Direct command execution without sanitization
	// Vulnerable: Command injection through string concatenation
	// Vulnerable: No command whitelisting
	// Vulnerable: No input validation
	cmd := exec.Command(req.Command, req.Args...)
	output, err := cmd.CombinedOutput()

	response := CommandResponse{
		Output:   string(output),
		ExitCode: cmd.ProcessState.ExitCode(),
	}

	if err != nil {
		response.Error = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleShellCommand implements a vulnerable shell command execution endpoint
// @Summary Execute shell command
// @Description Execute shell command (intentionally vulnerable to shell injection)
// @Tags commands
// @Accept json
// @Produce json
// @Param request body CommandRequest true "Shell command to execute"
// @Security BearerAuth
// @Success 200 {object} CommandResponse
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Server error"
// @Router /commands/shell [post]
func HandleShellCommand(w http.ResponseWriter, r *http.Request) {
	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Vulnerable: Shell command injection through string concatenation
	// Vulnerable: No command sanitization
	// Vulnerable: No input validation
	// Vulnerable: Shell metacharacter injection
	cmd := exec.Command("sh", "-c", req.Command)
	output, err := cmd.CombinedOutput()

	response := CommandResponse{
		Output:   string(output),
		ExitCode: cmd.ProcessState.ExitCode(),
	}

	if err != nil {
		response.Error = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleCommandWithFilter implements a vulnerable command execution with filtering
// @Summary Execute filtered command
// @Description Execute command with basic filtering (intentionally vulnerable)
// @Tags commands
// @Accept json
// @Produce json
// @Param request body CommandRequest true "Command to execute"
// @Security BearerAuth
// @Success 200 {object} CommandResponse
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Server error"
// @Router /commands/filter [post]
func HandleCommandWithFilter(w http.ResponseWriter, r *http.Request) {
	var req CommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Vulnerable: Ineffective command filtering
	// Vulnerable: Command injection through filter bypass
	// Vulnerable: No proper command validation
	// Vulnerable: Incomplete blacklist
	blacklist := []string{"rm", "mkfs", "dd", "format"}
	for _, cmd := range blacklist {
		if strings.Contains(req.Command, cmd) {
			http.Error(w, "Command not allowed", http.StatusBadRequest)
			return
		}
	}

	cmd := exec.Command("sh", "-c", req.Command)
	output, err := cmd.CombinedOutput()

	response := CommandResponse{
		Output:   string(output),
		ExitCode: cmd.ProcessState.ExitCode(),
	}

	if err != nil {
		response.Error = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
