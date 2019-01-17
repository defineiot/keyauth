package mysql

import (
	"database/sql"

	"github.com/defineiot/keyauth/dao"
	"github.com/defineiot/keyauth/dao/service"
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/defineiot/keyauth/internal/logger"
	"github.com/defineiot/keyauth/internal/tools"
)

const (
	SaveService           = "save-service"
	FindAllServices       = "find-all-service"
	FindServiceByID       = "find-service-by-id"
	FindServiceByClient   = "find-service-by-client"
	DeleteService         = "delete-service"
	DeleteServiceFeatures = "delete-services-features"

	SaveFeature             = "save-feature"
	MarkDeleteFeature       = "mark-delete-feature"
	AssociateFeaturesToRole = "add-feature-to-role"
	UnlinkFeatureFromRole   = "delete-feature-to-role"
	UpdateService           = "update-service"

	FindServiceFeatures = "find-service-features"
	FindRoleFeatures    = "find-role-features"
)

// NewServiceStore use to create domain storage service
func NewServiceStore(opt *dao.Options) (service.Store, error) {
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
		SaveFeature: `
			INSERT INTO features (id, name, tag, endpoint, description, is_deleted, when_deleted_version, when_deleted_time, is_added, when_added_version, when_added_time, service_id) 
			VALUES (?,?,?,?,?,?,?,?,?,?,?,?)
		`,
		MarkDeleteFeature: `
			UPDATE features  
			SET is_deleted=?, when_deleted_version=?, when_deleted_time=?   
			WHERE name=? 
			AND service_id=?;
		`,
		DeleteServiceFeatures: `
			DELETE FROM features 
			WHERE service_id = ?;
		`,
		AssociateFeaturesToRole: `
			INSERT INTO role_feature_mappings (feature_id, role_id) 
			VALUES (?,?);
		`,
		UnlinkFeatureFromRole: `
			DELETE FROM role_feature_mappings 
			WHERE feature_id = ? 
			AND role_id = ?;
		`,
		FindServiceFeatures: `
			SELECT id, name, tag, endpoint, description, is_deleted, when_deleted_version, when_deleted_time, is_added, when_added_version, when_added_time, service_id
			FROM features
			WHERE service_id = ? 
			ORDER BY tag
			DESC;
		`,
		FindRoleFeatures: `
			SELECT f.id, f.name, f.tag, f.endpoint, f.description, f.is_deleted, f.when_deleted_version, f.when_deleted_time, f.is_added, f.when_added_version, f.when_added_time, service_id 
			FROM features f 
			LEFT JOIN role_feature_mappings m
			ON m.feature_id = f.id 
			WHERE m.role_id = ? 
			ORDER BY f.tag 
			DESC; 
		`,
	}

	// prepare all statements to verify syntax
	stmts, err := tools.PrepareStmts(opt.DB, unprepared)
	if err != nil {
		return nil, exception.NewInternalServerError("prepare service store query statment error, %s", err)
	}

	s := store{
		db:    opt.DB,
		stmts: stmts,
		sql:   unprepared,
	}
	s.Logger = opt.LOG

	return &s, nil
}

// DomainManager is use mongodb as storage
type store struct {
	logger.Logger

	db    *sql.DB
	stmts map[string]*sql.Stmt
	sql   map[string]string
}

// Close closes the database, releasing any open resources.
func (s *store) Close() error {
	return s.db.Close()
}

func init() {
	dao.Registe(NewServiceStore)
}
