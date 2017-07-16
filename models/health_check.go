package models

type Healthy struct {
	Healthy bool `json:"healthy"`
}

type HealthCheckStatus struct {
	Alive        *Healthy `json:"alive"`
	ClientAccess *Healthy `json:"client_access"`
	UserAccess   *Healthy `json:"user_access"`
	KvsAccess    *Healthy `json:"key_value_store_access"`
}

func NewHealthCheckStatus() *HealthCheckStatus {
	return &HealthCheckStatus{
		Alive:        &Healthy{false},
		ClientAccess: &Healthy{false},
		UserAccess:   &Healthy{false},
		KvsAccess:    &Healthy{false},
	}
}
