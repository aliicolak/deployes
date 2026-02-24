package deployment 

type CreateDeploymentRequest struct {
	ProjectID string `json:"projectId"`
	ServerID  string `json:"serverId"`
}

type DeploymentResponse struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	Logs       string `json:"logs"`
	ProjectID  string `json:"projectId"`
	ServerID   string `json:"serverId"`
	CommitHash string `json:"commitHash,omitempty"`
}

type StatsResponse struct {
	Total           int     `json:"total"`
	Successful      int     `json:"successful"`
	Failed          int     `json:"failed"`
	AverageDuration float64 `json:"averageDuration"`
	Last7Days       struct {
		Dates  []string `json:"dates"`
		Counts []int    `json:"counts"`
	} `json:"last7Days"`
}
