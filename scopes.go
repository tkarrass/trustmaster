package trustmaster

type Scope int

const (
	GetAgentName = 0
	GetEmailAddress = 1
	GetGoogleName = 2
	GetGoogleAvatar = 3
	GetTelegram = 4
	GetLocation = 5
	GetZello = 6
	GetLanguages = 7
)

var names = [...]string{
	"get-agent-name",
	"get-email-address",
	"get-google-name",
	"get-google-avatar",
	"get-telegram",
	"get-location",
	"get-zello",
	"get-languages",
}

func (s Scope) String() string {
	if s < GetAgentName || s > GetLanguages {
		return "unknown"
	}
	return names[s]
}
