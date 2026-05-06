package tailscale

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"time"
)

const defaultSocketPath = "/var/run/tailscale/tailscaled.sock"

type Location struct {
	Country     string `json:"Country"`
	CountryCode string `json:"CountryCode"`
	City        string `json:"City"`
	CityCode    string `json:"CityCode"`
}

type PeerStatus struct {
	HostName       string    `json:"HostName"`
	DNSName        string    `json:"DNSName"`
	TailscaleIPs   []string  `json:"TailscaleIPs"`
	OS             string    `json:"OS"`
	Online         bool      `json:"Online"`
	LastSeen       time.Time `json:"LastSeen"`
	ExitNodeOption bool      `json:"ExitNodeOption"`
	ExitNode       bool      `json:"ExitNode"`
	Location       *Location `json:"Location"`
}

type SelfStatus struct {
	HostName     string   `json:"HostName"`
	TailscaleIPs []string `json:"TailscaleIPs"`
	OS           string   `json:"OS"`
	Online       bool     `json:"Online"`
}

type TailnetStatus struct {
	Name           string `json:"Name"`
	MagicDNSSuffix string `json:"MagicDNSSuffix"`
}

type ExitNodeStatus struct {
	Active bool `json:"Active"`
}

type Status struct {
	BackendState   string                `json:"BackendState"`
	Self           SelfStatus            `json:"Self"`
	CurrentTailnet TailnetStatus         `json:"CurrentTailnet"`
	Peer           map[string]PeerStatus `json:"Peer"`
	ExitNodeStatus *ExitNodeStatus       `json:"ExitNodeStatus"`
}

func socketPath() string {
	if s := os.Getenv("TAILSCALE_SOCKET"); s != "" {
		return s
	}
	return defaultSocketPath
}

func newClient() *http.Client {
	sock := socketPath()
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "unix", sock)
			},
		},
	}
}

// CheckSocket returns nil if the tailscaled socket is reachable.
func CheckSocket() error {
	conn, err := net.DialTimeout("unix", socketPath(), 2*time.Second)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}

// FetchStatus fetches current Tailscale status from the local API.
func FetchStatus() (*Status, error) {
	resp, err := newClient().Get("http://local-tailscaled.sock/localapi/v0/status")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var s Status
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		return nil, err
	}
	return &s, nil
}
