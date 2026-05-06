package tailscale

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
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

// CanManage returns true if the process can perform Tailscale management operations.
// Requires root or membership in the 'tailscale' group (which owns the socket).
func CanManage() bool {
	if os.Getuid() == 0 {
		return true
	}
	g, err := user.LookupGroup("tailscale")
	if err != nil {
		return false
	}
	cur, err := user.Current()
	if err != nil {
		return false
	}
	gids, err := cur.GroupIds()
	if err != nil {
		return false
	}
	for _, gid := range gids {
		if gid == g.Gid {
			return true
		}
	}
	return false
}
