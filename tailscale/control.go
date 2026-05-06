package tailscale

import (
	"fmt"
	"os/exec"
	"strings"
)

// SetExitNode connects to the given exit node IP, or clears the exit node if ip is empty.
func SetExitNode(ip string) error {
	out, err := exec.Command("tailscale", "set", "--exit-node="+ip).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}
