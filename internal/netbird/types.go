package netbird

import "time"

type Peer struct {
	AccessiblePeersCount int       `json:"accessible_peers_count"`
	ApprovalRequired     bool      `json:"approval_required"`
	CityName             string    `json:"city_name"`
	Connected            bool      `json:"connected"`
	ConnectionIP         string    `json:"connection_ip"`
	CountryCode          string    `json:"country_code"`
	CreatedAt            time.Time `json:"created_at"`
	DNSLabel             string    `json:"dns_label"`
	Ephemeral            bool      `json:"ephemeral"`
	ExtraDNSLabels       []any     `json:"extra_dns_labels"`
	GeonameID            int       `json:"geoname_id"`
	Groups               []struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		PeersCount     int    `json:"peers_count"`
		ResourcesCount int    `json:"resources_count"`
	} `json:"groups"`
	Hostname                    string    `json:"hostname"`
	ID                          string    `json:"id"`
	InactivityExpirationEnabled bool      `json:"inactivity_expiration_enabled"`
	IP                          string    `json:"ip"`
	Ipv6                        string    `json:"ipv6"`
	KernelVersion               string    `json:"kernel_version"`
	LastLogin                   time.Time `json:"last_login"`
	LastSeen                    time.Time `json:"last_seen"`
	LocalFlags                  struct {
		BlockInbound          bool `json:"block_inbound"`
		BlockLanAccess        bool `json:"block_lan_access"`
		DisableClientRoutes   bool `json:"disable_client_routes"`
		DisableDNS            bool `json:"disable_dns"`
		DisableFirewall       bool `json:"disable_firewall"`
		DisableServerRoutes   bool `json:"disable_server_routes"`
		LazyConnectionEnabled bool `json:"lazy_connection_enabled"`
		RosenpassEnabled      bool `json:"rosenpass_enabled"`
		RosenpassPermissive   bool `json:"rosenpass_permissive"`
		ServerSSHAllowed      bool `json:"server_ssh_allowed"`
	} `json:"local_flags"`
	LoginExpirationEnabled bool   `json:"login_expiration_enabled"`
	LoginExpired           bool   `json:"login_expired"`
	Name                   string `json:"name"`
	Os                     string `json:"os"`
	SerialNumber           string `json:"serial_number"`
	SSHEnabled             bool   `json:"ssh_enabled"`
	UIVersion              string `json:"ui_version"`
	UserID                 string `json:"user_id"`
	Version                string `json:"version"`
}

type Key struct {
	Expires time.Time `json:"expiration_date"`
}
