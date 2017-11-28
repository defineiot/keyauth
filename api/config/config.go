package config

// Configer use to get conf
type Configer interface {
	GetConf() (*Config, error)
}

// Config is service conf
type Config struct {
	APP     *appConf
	MongoDB *mongoConf
}

type appConf struct {
	Host string
	Port string
}

type mongoConf struct {
	Host string
	Port string
	User string
	Pass string
	DB   string
}

// Validate use to check the service config
func (c *Config) Validate() error {
	if err := c.validateAPP(); err != nil {
		return err
	}

	if err := c.validateMySQL(); err != nil {
		return err
	}

	return nil
}

func (c *Config) validateAPP() error {
	if c.APP == nil {
		c.APP = &appConf{}
	}

	if c.APP.Host == "" {
		c.APP.Host = "0.0.0.0"
	}
	if c.APP.Port == "" {
		c.APP.Port = "8080"
	}

	return nil
}

func (c *Config) validateMySQL() error {
	if c.MongoDB == nil {
		c.MongoDB = &mongoConf{}
	}

	if c.MongoDB.Host == "" {
		c.MongoDB.Host = "127.0.0.1"
	}
	if c.MongoDB.Port == "" {
		c.MongoDB.Port = "3306"
	}

	return nil
}
