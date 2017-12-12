package application

// Application is oauth2's client: https://tools.ietf.org/html/rfc6749#section-2
type Application struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Name         string `json:"name"`
	Website      string `json:"website"`
	LogoImage    string `json:"logo_image"`
	Description  string `json:"description"`
	RedirectURI  string `json:"redirect_uri"`
	CreateAt     int64  `json:"create_at"`
	ClientType   string `json:"client_type"` // 1. confidential 2.public  https://tools.ietf.org/html/rfc6749#section-2.1
}

// Storage appliction stroage
type Storage interface {
	Registration(userID, name, redirectURI, clientType, description, website string) (*Application, error)
	Unregistration(userID, clientID string) error
	GetUserApps(userID string) ([]*Application, error)
}
