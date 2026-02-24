package server

type CreateServerRequest struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	SSHKey   string `json:"sshKey"`
}

type TestConnectionRequest struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	SSHKey   string `json:"sshKey"`
	ServerID string `json:"serverId,omitempty"` // Optional: test existing server
}

type ConnectionTestResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Latency int64  `json:"latency"` // Milliseconds
}

type ServerResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
}
