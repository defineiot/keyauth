package client

const (
	Confidential Type = "confidential"
	Public       Type = "public"
)

// Type the type of authorization request
type Type string

// Client is oauth2 client
type Client struct {
	ID          string `json:"client_id"`
	Secret      string `json:"client_secret"`
	RedirectURI string `json:"redirect_uri,omitempty"`
	CreateAt    int64  `json:"client_create_at,omitempty"`
	Type        Type   `json:"client_type"` // 1. confidential 2.public  https://tools.ietf.org/html/rfc6749#section-2.1
}

// Store service store interface
type Store interface {
	Reader
	Writer
	Close() error
}

// Reader read service information from store
type Reader interface {
	ListClients() ([]*Client, error)
	GetClient(id string) (*Client, error)
}

// Writer write service information to store
type Writer interface {
	CreateClient(clientType Type, redirectURI string) (*Client, error)
	DeleteClient(id string) error
}
