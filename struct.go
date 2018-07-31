package trustmaster

import "time"

type AccessToken struct {
	Token string  `json:"access_token"`
	ExpiresIn uint64 `json:"expires_in"`
	Expires time.Time
	Type string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	generic bool
	user bool
}

func (token *AccessToken) Expired() bool {
	return time.Now().After(token.Expires)
}

type Version struct {
	Version uint `json:"version"`
}

type Profile struct {
	AgentName string `json:"agent_name"`
	Email string `json:"email"`
	AvatarUrl string `json:"google_avatar"`
	GoogleId string `json:"google_id"`
	GoogleName string `json:"google_name"`
	Languages []string `json:"languages"`
	Telegram *Telegram `json:"telegram"`
	Location *Location `json:"location"`
	Zello string `json:"zello_username"`
}

type Location struct {
	Lat string `json:"lat"`
	Lon string `json:"lng"`
	Privacy string `json:"privacy"`
}

type Telegram struct {
	Id string `json:"telegram_id"`
	Name string `json:"telegram_username"`
}

type Trust struct {
	Generation uint `json:"generation"`
	GoogleId string `json:"google_id"`
	NewestTrust *TrustDecision `json:"newest_trust"`
	Trust *TrustDecision `json:"trust"`
	Summary *TrustSummary `json:"summary"`
	SummaryUnverified *TrustSummary `json:"summary_unverified"`
}

type TrustDecision struct {
	Decision string `json:"decision"`
	Updated string `json:"updated_at"`
}

type TrustSummary struct {
	Admin uint `json:"admin"`
	// ... ????
}

