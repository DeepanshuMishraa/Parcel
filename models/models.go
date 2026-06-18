package models

type Status string

const (
	QUEUED   Status = "queued"
	RUNNING  Status = "running"
	FINISHED Status = "finished"
	FAILED   Status = "failed"
)

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Job struct {
	JobID     string         `json:"jobId"`
	JobName   string         `json:"job_name"`
	JobStatus Status         `json:"status"`
	Payload   map[string]any `json:"payload"`
}
