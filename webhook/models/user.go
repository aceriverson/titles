package models

type UserPlan string

const (
	None UserPlan = "none"
	Free UserPlan = "free"
	Pro  UserPlan = "pro"
)

type User struct {
	ID           int64    `json:"id"`
	Pic          string   `json:"pic"`
	Name         string   `json:"name"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresAt    int64    `json:"expires_at"`
	AI           bool     `json:"ai"`
	Plan         UserPlan `json:"plan"`
}

func (u *User) UpdateToken(r RefreshResponse) {
	u.AccessToken = r.AccessToken
	u.RefreshToken = r.RefreshToken
	u.ExpiresAt = r.ExpiresAt
}
