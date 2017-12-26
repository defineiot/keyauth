package application

// Application is oauth2's client: https://tools.ietf.org/html/rfc6749#section-2
type Application struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Name        string  `json:"name"`
	Website     string  `json:"website,omitempty"`
	LogoImage   string  `json:"logo_image,omitempty"`
	Description string  `json:"description"`
	CreateAt    int64   `json:"create_at"`
	Client      *Client `json:"client"`
}

// Client is oauth2 client
type Client struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
	ClientType   string `json:"client_type"` // 1. confidential 2.public  https://tools.ietf.org/html/rfc6749#section-2.1
}


// Store application storage
type Store interface {
	StoreReader
	StoreWriter
	Close() error

}

// StoreReader use to read application information from store
type StoreReader interface {
	GetUserApps(userID string) ([]*Application, error)
	CheckAPPIsExistByID(appID string) (bool, error)
	CheckAPPIsExistByName(userID, name string) (bool, error)
	GetClient(clientID string) (*Client, error)
}

// StoreWriter use to write application information from store
type StoreWriter interface {
	Registration(userID, name, redirectURI, clientType, description, website string) (*Application, error)
	Unregistration(id string) error
}
