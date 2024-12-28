package models

type Map struct {
	Polyline string `json:"polyline"`
}

type Activity struct {
	ID          string `json:"id"`
	Map         Map    `json:"map"`
	Description string `json:"description"`
	SportType   string `json:"sport_type"`
}

type Webhook struct {
	ObjectType     string            `json:"object_type"`
	ObjectID       int64             `json:"object_id"`
	AspectType     string            `json:"aspect_type"`
	Updates        map[string]string `json:"updates"`
	OwnerID        int64             `json:"owner_id"`
	SubscriptionID int64             `json:"subscription_id"`
	EventTime      int64             `json:"event_time"`
}

type RefreshBody struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresAt    int64  `json:"expires_at"`
	RefreshToken string `json:"refresh_token"`
}

type Update struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
