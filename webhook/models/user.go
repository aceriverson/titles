package models

type UserPlan string

const (
	UserPlanNone UserPlan = "none"
	UserPlanFree UserPlan = "free"
	UserPlanPro  UserPlan = "pro"
)

type User struct {
	ID            int64    `json:"id"`
	Pic           string   `json:"pic"`
	Name          string   `json:"name"`
	AccessToken   string   `json:"access_token"`
	RefreshToken  string   `json:"refresh_token"`
	ExpiresAt     int64    `json:"expires_at"`
	Plan          UserPlan `json:"plan"`
	TermsAccepted bool     `json:"terms_accepted"`
}

func (u *User) UpdateToken(r RefreshResponse) {
	u.AccessToken = r.AccessToken
	u.RefreshToken = r.RefreshToken
	u.ExpiresAt = r.ExpiresAt
}
