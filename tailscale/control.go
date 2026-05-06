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
// On distros that use a 'tailscale' group (Debian/Ubuntu), checks group membership.
// If the group doesn't exist, access is assumed to be open (Arch, NixOS, etc.).
func CanManage() bool {
	if os.Getuid() == 0 {
		return true
	}
	g, err := user.LookupGroup("tailscale")
	if err != nil {
		// No tailscale group — socket access is not group-restricted on this distro.
		return true
	}
	cur, err := user.Current()
	if err != nil {
		return true
	}
	gids, err := cur.GroupIds()
	if err != nil {
		return true
	}
	for _, gid := range gids {
		if gid == g.Gid {
			return true
		}
	}
	return false
}
