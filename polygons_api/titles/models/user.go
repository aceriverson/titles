package models

type User struct {
	ID   int64  `json:"id"`
	Pic  string `json:"pic"`
	Name string `json:"name"`
}

type UserPlan string

const (
	None UserPlan = "none"
	Free UserPlan = "free"
	Pro  UserPlan = "pro"
)

type UserInternal struct {
	ID           int64    `json:"id"`
	Pic          string   `json:"pic"`
	Name         string   `json:"name"`
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresAt    int64    `json:"expires_at"`
	AI           bool     `json:"ai"`
	Plan         UserPlan `json:"plan"`
}

type TokenExchangeResponse struct {
	ExpiresAt    int64  `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	Athlete      struct {
		ID        int64  `json:"id"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Profile   string `json:"profile"`
	} `json:"athlete"`
}

func (u *TokenExchangeResponse) ToUserInternal() UserInternal {
	return UserInternal{
		ID:           u.Athlete.ID,
		Pic:          u.Athlete.Profile,
		Name:         u.Athlete.Firstname + " " + u.Athlete.Lastname,
		AccessToken:  u.AccessToken,
		RefreshToken: u.RefreshToken,
		ExpiresAt:    u.ExpiresAt,
	}
}

func (u *UserInternal) UpdateToken(r RefreshResponse) {
	u.AccessToken = r.AccessToken
	u.RefreshToken = r.RefreshToken
	u.ExpiresAt = r.ExpiresAt
}
