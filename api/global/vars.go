package global

import (
	"github.com/defineiot/keyauth/internal/conf"
	"github.com/defineiot/keyauth/internal/logger"
	"github.com/defineiot/keyauth/store"
)

var (
	// Conf app config
	Conf *conf.Config
	// Log app log
	Log logger.Logger
	// Store db stroe
	Store *store.Store
)
