package types

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binidng:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email" binidng:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type CreateJobRequest struct {
	JobName string         `json:"job_name" binding:"required"`
	Status  string         `json:"status"`
	UserId  string         `json:"user_id" binding:"required"`
	Payload map[string]any `json:"payload" binding:"required"`
}

type CreateJobResponse struct {
	JobID   string `json:"job_id"`
	JobName string `json:"job_name"`
	Status  string `json:"status"`
}
