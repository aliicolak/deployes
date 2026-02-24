package deployment

import "time"

type Deployment struct {
	ID             string
	UserID         string
	ProjectID      string
	ServerID       string
	Status         string
	Logs           string
	CommitHash     string
	RollbackFromID string
	CreatedAt      time.Time
	StartedAt      *time.Time
	FinishedAt     *time.Time
}

const (
	StatusQueued  = "queued"
	StatusRunning = "running"
	StatusSuccess = "success"
	StatusFailed  = "failed"
)

type Stats struct {
	Total                  int
	Successful             int
	Failed                 int
	AverageDurationSeconds float64
	// Last 7 days, simple slice for chart
	Last7DaysCounts []int
	Last7DaysDates  []string
}
