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
	SaveFeature           = "save-feature"
	UpdateService         = "update-service"
	DeleteService         = "delete-service"
	DeleteServiceFeatures = "delete-services-features"
	FindAll               = "find-all-service"
	FindOneByID           = "find-one-service"
	FindOneByClient       = "find-one-service-by-client-id"
	FindAllFeatures       = "find-one-service-features"
	FindFullAllFeatures   = "find-all-featrues"
	CheckServiceExist     = "check-service-exist"
	CheckFeatureExist     = "check-service-feature-exist"
	CheckFeatureIDExist   = "check-feature-exist"
	FindRoleFeatures      = "find-role-features"
)

// NewServiceStore use to create domain storage service
func NewServiceStore(db *sql.DB, log log.IOTAuthLogger) (service.Store, error) {
	unprepared := map[string]string{
		SaveService: `
			INSERT INTO services (name, description, enabled, status, status_update_at, version, create_at, client_id) 
			VALUES (?,?,?,?,?,?,?,?);
		`,
		SaveFeature: `
			INSERT INTO features (name, method, endpoint, description, is_deleted, when_deleted_version, is_added, when_added_version, service_name) 
			VALUES (?,?,?,?,?,?,?,?,?)
		`,
		FindAll: `
			SELECT name, description, enabled, status, status_update_at, version, create_at, client_id 
			FROM services 
			ORDER BY create_at
			DESC;
		`,
		FindAllFeatures: `
			SELECT id, name, method, endpoint, description, is_deleted, when_deleted_version, is_added, when_added_version, service_name
			FROM features
			WHERE service_name = ? 
			ORDER BY method
			DESC;
		`,
		FindRoleFeatures: `
		    SELECT f.id, f.name, f.method, f.endpoint, f.description, f.is_deleted, f.when_deleted_version, f.is_added, f.when_added_version, service_name
			FROM features f
			LEFT JOIN roles_features_mapping m
			ON f.id = m.feature_id
		    WHERE m.role_name = ? 
		    ORDER BY f.method
		    DESC;
	    `,
		FindFullAllFeatures: `
		    SELECT id, name, method, endpoint, description, is_deleted, when_deleted_version, is_added, when_added_version, service_name
		    FROM features
		    ORDER BY method
		    DESC;
		`,
		FindOneByID: `
			SELECT s.name, s.description, s.enabled, s.status, s.status_update_at, s.version, s.create_at, client_id
			FROM services s
			WHERE s.name = ?;
		`,
		FindOneByClient: `
		    SELECT s.name, s.description, s.enabled, s.status, s.status_update_at, s.version, s.create_at, client_id
		    FROM services s
		    WHERE s.client_id = ?;
		`,
		DeleteService: `
			DELETE FROM services 
			WHERE name =?;
		`,
		DeleteServiceFeatures: `
		    DELETE FROM features 
			WHERE service_name = ?;
		`,
		CheckServiceExist: `
		    SELECT name
		    FROM services
		    WHERE name = ?;
	    `,
		CheckFeatureExist: `
		    SELECT name
		    FROM features
			WHERE name = ? 
			AND service_name = ?;
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
