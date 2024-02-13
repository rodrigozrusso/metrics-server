package healthcheck

type HealthCheckStatus string

const (
	Ok   HealthCheckStatus = "up"
	Warm HealthCheckStatus = "warm"
	Fail HealthCheckStatus = "down"
)

type ResourceHealthCheckResponse struct {
	Name    string            `json:"name"`
	Status  HealthCheckStatus `json:"status"`
	Message string            `json:"message"`
}

type HealthCheckResponse struct {
	Status    HealthCheckStatus             `json:"status"`
	Resources []ResourceHealthCheckResponse `json:"resources"`
}
