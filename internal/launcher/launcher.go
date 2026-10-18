package launcher

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const AgentArg = "-agentlib:zbNative"

type Config struct {
	MainClass string                            `json:"mainClass"`
	Classpath []string                          `json:"classpath"`
	VMArgs    []string                          `json:"vmArgs"`
	Windows   map[string]launcherWindowsVersion `json:"windows,omitempty"`
}

type launcherWindowsVersion struct {
	VMArgs []string `json:"vmArgs"`
}

// HasAgent reports whether vmArgs already includes the ZombieBuddy agent.
func HasAgent(args []string) bool {
	for _, arg := range args {
		if isAgentOption(arg) {
			return true
		}
	}
	return false
}

func isAgentOption(value string) bool {
	return value == AgentArg || strings.HasPrefix(value, AgentArg+"=")
}

// NeedsPatch reads the launcher JSON and reports if the agent arg is missing.
func NeedsPatch(path string) (bool, error) {
	cfg, found, err := readConfig(path)
	if err != nil {
		return false, err
	}
	if !found {
		return false, fmt.Errorf("launcher config not found: %s", path)
	}
	return !HasAgent(cfg.VMArgs), nil
}

// Patch prepends -agentlib:zbNative to vmArgs when missing.
func Patch(path string) (changed bool, err error) {
	cfg, found, err := readConfig(path)
	if err != nil {
		return false, err
	}
	if !found {
		return false, fmt.Errorf("launcher config not found: %s", path)
	}
	if HasAgent(cfg.VMArgs) {
		return false, nil
	}
	cfg.VMArgs = append([]string{AgentArg}, cfg.VMArgs...)
	output, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return false, err
	}
	output = append(output, '\n')
	return true, os.WriteFile(path, output, 0644)
}

func readConfig(path string) (Config, bool, error) {
	input, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, false, nil
		}
		return Config{}, false, err
	}
	var cfg Config
	if err := json.Unmarshal(input, &cfg); err != nil {
		return Config{}, false, err
	}
	return cfg, true, nil
}

// AddAgent returns vmArgs with the agent prepended (for tests).
func AddAgent(args []string) ([]string, bool) {
	if HasAgent(args) {
		return args, false
	}
	return append([]string{AgentArg}, args...), true
}
