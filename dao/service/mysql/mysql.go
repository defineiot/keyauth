package mysql

import (
	"database/sql"

	"github.com/defineiot/keyauth/dao/service"
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/defineiot/keyauth/internal/log"
	"github.com/defineiot/keyauth/internal/tools"
)

const (
	SaveService           = "save-service"
	FindAllServices       = "find-all-service"
	FindServiceByID       = "find-service-by-id"
	FindServiceByClient   = "find-service-by-client"
	DeleteService         = "delete-service"
	DeleteServiceFeatures = "delete-services-features"
	CheckServiceExist     = "check-service-exist"

	SaveFeature   = "save-feature"
	UpdateService = "update-service"

	FindAllFeatures     = "find-one-service-features"
	FindFullAllFeatures = "find-all-featrues"
	CheckFeatureExist   = "check-service-feature-exist"
	CheckFeatureIDExist = "check-feature-exist"
	FindRoleFeatures    = "find-role-features"
)

// NewServiceStore use to create domain storage service
func NewServiceStore(db *sql.DB, log log.IOTAuthLogger) (service.Store, error) {
	unprepared := map[string]string{
		SaveService: `
			INSERT INTO services (id, type, name, description, enabled, create_at, client_id, client_secret, token_expire_time) 
			VALUES (?,?,?,?,?,?,?,?,?);
		`,
		FindAllServices: `
			SELECT id, type, name, description, enabled, status, status_update_at, current_version, upgrade_version, downgrade_version, create_at, update_at, client_id, client_secret, token_expire_time 
			FROM services;
		`,
		FindServiceByID: `
			SELECT id, type, name, description, enabled, status, status_update_at, current_version, upgrade_version, downgrade_version, create_at, update_at, client_id, client_secret, token_expire_time 
			FROM services 
			WHERE id = ?;
		`,
		FindServiceByClient: `
			SELECT id, type, name, description, enabled, status, status_update_at, current_version, upgrade_version, downgrade_version, create_at, update_at, client_id, client_secret, token_expire_time 
			FROM services 
			WHERE client_id = ?;
		`,
		DeleteService: `
    		DELETE FROM services 
	    	WHERE id =?;
		`,
		DeleteServiceFeatures: `
			DELETE FROM features 
			WHERE service_id = ?;
		`,
		FindAllFeatures: `
			SELECT id, name, tag, endpoint, description, is_deleted, when_deleted_version, is_added, when_added_version, service_id
			FROM features
			WHERE service_id = ? 
			ORDER BY tag
			DESC;
		`,
		FindRoleFeatures: `
		    SELECT f.id, f.name, f.tag, f.endpoint, f.description, f.is_deleted, f.when_deleted_version, f.is_added, f.when_added_version, service_id
			FROM features f
			LEFT JOIN role_feature_mappings m
			ON f.id = m.feature_id
		    WHERE m.role_id = ? 
		    ORDER BY f.tag
		    DESC;
	    `,
		FindFullAllFeatures: `
		    SELECT id, name, tag, endpoint, description, is_deleted, when_deleted_version, is_added, when_added_version, service_id
		    FROM features
		    ORDER BY tag
		    DESC;
		`,
		CheckServiceExist: `
		    SELECT id
		    FROM services
		    WHERE id = ?;
	    `,
		CheckFeatureExist: `
		    SELECT name
		    FROM features
			WHERE name = ? 
			AND service_id = ?;
		`,
		CheckFeatureIDExist: `
		    SELECT id
		    FROM features
		    WHERE id = ?;
		`,
	}

	// prepare all statements to verify syntax
	stmts, err := tools.PrepareStmts(db, unprepared)
	if err != nil {
		return nil, exception.NewInternalServerError("prepare service store query statment error, %s", err)
	}

	s := store{
		db:    db,
		stmts: stmts,
		log:   log,
	}

	return &s, nil
}

// DomainManager is use mongodb as storage
type store struct {
	db    *sql.DB
	log   log.IOTAuthLogger
	stmts map[string]*sql.Stmt
}

// Close closes the database, releasing any open resources.
func (s *store) Close() error {
	return s.db.Close()
}
