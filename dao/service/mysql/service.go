package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/defineiot/keyauth/dao/service"
	"github.com/defineiot/keyauth/internal/exception"
)

var (
	domainSkip = []string{"CreateService", "ListServices", "GetService", "DeleteService",
		"RegistryServiceFeatures", "ListServiceFeatures", "CreateRole", "DeleteRole",
		"AddFeaturesToRole", "RemoveFeaturesFromRole", "CreateDomain", "DeleteDomain",
		"IssueToken", "ValidateToken", "RevolkToken"}

	memberSkip = []string{"ListServices", "GetService", "ListServiceFeatures",
		"ListApplications", "GetApplication", "ValidateToken"}
)

func (s *store) CreateService(name, description, clientID string) (*service.Service, error) {
	ok, err := s.CheckServiceIsExist(name)
	if err != nil {
		return nil, err
	}
	if ok {
		return nil, exception.NewBadRequest("service %s existed", name)
	}

	svr := service.Service{Name: name, Description: description, CreateAt: time.Now().Unix(), Enabled: true}

	_, err = s.stmts[SaveService].Exec(svr.Name, svr.Description, svr.Enabled, svr.Status, svr.StatusUpdateAt, svr.Version, svr.CreateAt, clientID)
	if err != nil {
		return nil, exception.NewInternalServerError("insert client exec sql err, %s", err)
	}

	return &svr, nil
}

func (s *store) CheckServiceIsExist(name string) (bool, error) {
	var n string
	if err := s.stmts[CheckServiceExist].QueryRow(name).Scan(&n); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, exception.NewInternalServerError("query check service %s error, %s", name, err)
	}

	return true, nil
}

func (s *store) DeleteService(name string) error {
	// delete service
	ret, err := s.stmts[DeleteService].Exec(name)
	if err != nil {
		return exception.NewInternalServerError("delete service exec sql error, %s", err)
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return exception.NewInternalServerError("get delete service row affected error, %s", err)
	}
	if count == 0 {
		return exception.NewBadRequest("service %s not exist", name)
	}

	// delete service features
	_, err = s.stmts[DeleteServiceFeatures].Exec(name)
	if err != nil {
		return exception.NewInternalServerError("delete service features exec sql error, %s", err)
	}

	return nil
}

func (s *store) ListServices() ([]*service.Service, error) {
	rows, err := s.stmts[FindAll].Query()
	if err != nil {
		return nil, exception.NewInternalServerError("query service list error, %s", err)
	}
	defer rows.Close()

	svrs := []*service.Service{}
	for rows.Next() {
		svr := service.Service{}
		if err := rows.Scan(&svr.Name, &svr.Description, &svr.Enabled, &svr.Status, &svr.StatusUpdateAt, &svr.Version, &svr.CreateAt, &svr.ClientID); err != nil {
			return nil, exception.NewInternalServerError("scan service record error, %s", err)
		}
		svrs = append(svrs, &svr)
	}

	return svrs, nil
}

func (s *store) GetService(name string) (*service.Service, error) {
	svr := service.Service{}
	if err := s.stmts[FindOneByID].QueryRow(name).Scan(&svr.Name, &svr.Description, &svr.Enabled, &svr.Status, &svr.StatusUpdateAt, &svr.Version, &svr.CreateAt, &svr.ClientID); err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("service %s not find", name)
		}

		return nil, exception.NewInternalServerError("query single service client error, %s", err)
	}

	return &svr, nil
}

func (s *store) GetServiceByClientID(clientID string) (*service.Service, error) {
	svr := service.Service{}
	if err := s.stmts[FindOneByClient].QueryRow(clientID).Scan(&svr.Name, &svr.Description, &svr.Enabled, &svr.Status, &svr.StatusUpdateAt, &svr.Version, &svr.CreateAt, &svr.ClientID); err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("client %s service not find", clientID)
		}

		return nil, exception.NewInternalServerError("query single service client error, %s", err)
	}

	return &svr, nil
}

func (s *store) RegistryServiceFeatures(name string, features ...service.Feature) error {
	s.log.Debugf("registry service :%s features: %v", name, features)

	hasF, err := s.ListServiceFeatures(name)
	if err != nil {
		return err
	}

	// find need add
	added := []*service.Feature{}
	for i, in := range features {
		exist := false
		for _, has := range hasF {
			if in.Name == has.Name {
				exist = true
				break
			}
		}
		if !exist {
			added = append(added, &features[i])
		}
	}

	// find need delete
	deleted := []*service.Feature{}
	for _, has := range hasF {
		var exist bool
		for _, in := range features {
			if has.Name == in.Name {
				exist = true
				break
			}
		}

		if !exist {
			deleted = append(deleted, has)
		}
	}

	// start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("start save features transaction error, %s", err)
	}
	// commit transaction
	defer func() {
		if err := tx.Commit(); err != nil {
			s.log.Errorf("feature transaction commit error, %s", err)
		}
	}()

	// added new features
	addFeaturePre, err := tx.Prepare("INSERT INTO features (name, method, endpoint, description, is_deleted, when_deleted_version, is_added, when_added_version, service_name) VALUES (?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return exception.NewInternalServerError("prepare insert feature stmt error, name: %s, %s", name, err)
	}
	defer addFeaturePre.Close()

	for _, f := range added {
		s.log.Infof("service: %s add feature: %s", name, f.Name)
		// exec sql
		_, err := addFeaturePre.Exec(f.Name, f.Method, f.Endpoint, f.Description, f.IsDeleted, f.WhenDeletedVersion, f.IsAdded, f.WhenAddedVersion, name)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				s.log.Errorf("insert feature transaction rollback error, %s", err)
			}
			return exception.NewInternalServerError("insert feature exec sql err, %s", err)
		}
	}

	// mark delete features
	delFeaturePre, err := tx.Prepare("UPDATE features SET is_deleted=? WHERE name=? AND service_name=?")
	if err != nil {
		return exception.NewInternalServerError("prepare update delete mark feature stmt error, name: %s, %s", name, err)
	}
	defer delFeaturePre.Close()

	for _, f := range deleted {
		s.log.Infof("service: %s del feature: %s", name, f.Name)
		// exec sql
		_, err := delFeaturePre.Exec(1, f.Name, name)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				s.log.Errorf("update delete mark feature feature transaction rollback error, %s", err)
			}
			return exception.NewInternalServerError("update delete mark feature feature exec sql err, %s", err)
		}
	}
	return nil
}

func (s *store) ListServiceFeatures(name string) ([]*service.Feature, error) {
	rows, err := s.stmts[FindAllFeatures].Query(name)
	if err != nil {
		return nil, exception.NewInternalServerError("query service feature list error, %s", err)
	}
	defer rows.Close()

	features := []*service.Feature{}
	for rows.Next() {
		f := service.Feature{}
		if err := rows.Scan(&f.ID, &f.Name, &f.Method, &f.Endpoint, &f.Description, &f.IsDeleted, &f.WhenDeletedVersion, &f.IsAdded, &f.WhenAddedVersion, &f.ServiceName); err != nil {
			return nil, exception.NewInternalServerError("scan service feature record error, %s", err)
		}
		features = append(features, &f)
	}

	return features, nil
}

func (s *store) CheckServiceHasFeature(serviceName, featureName string) (bool, error) {
	var name string
	if err := s.stmts[CheckFeatureExist].QueryRow(featureName, serviceName).Scan(&name); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, exception.NewInternalServerError("check service feature exist error, %s", err)
	}

	return true, nil
}

func (s *store) ListDomainFeatures() ([]*service.Feature, error) {
	df := []*service.Feature{}
	features, err := s.listFeatures()
	if err != nil {
		return nil, err
	}

	for _, f := range features {
		isok := true
		for _, skip := range domainSkip {
			if f.Name == skip {
				isok = false
				break
			}
		}

		if isok {
			df = append(df, f)
		}
	}

	return df, nil

}

func (s *store) ListMemberFeatures() ([]*service.Feature, error) {
	df := []*service.Feature{}
	features, err := s.listFeatures()
	if err != nil {
		return nil, err
	}

	for _, f := range features {
		isok := true
		for _, skip := range memberSkip {
			if f.Name == skip {
				isok = false
				break
			}
		}

		if f.Name == "SetUserPassword" {
			df = append(df, f)
		}

		if isok && f.Method == "GET" {
			df = append(df, f)
		}
	}

	return df, nil
}

func (s *store) CheckFeatureIsExist(featureID int64) (bool, error) {
	var id string
	if err := s.stmts[CheckFeatureIDExist].QueryRow(featureID).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, exception.NewInternalServerError("query feature exist error, %s", err)
	}

	return true, nil
}

func (s *store) ListRoleFeatures(name string) ([]*service.Feature, error) {
	rows, err := s.stmts[FindRoleFeatures].Query(name)
	if err != nil {
		return nil, exception.NewInternalServerError("query all feature list error, %s", err)
	}
	defer rows.Close()

	features := []*service.Feature{}
	for rows.Next() {
		f := service.Feature{}
		if err := rows.Scan(&f.ID, &f.Name, &f.Method, &f.Endpoint, &f.Description, &f.IsDeleted, &f.WhenDeletedVersion, &f.IsAdded, &f.WhenAddedVersion, &f.ServiceName); err != nil {
			return nil, exception.NewInternalServerError("scan service feature record error, %s", err)
		}
		features = append(features, &f)
	}

	return features, nil

}

func (s *store) listFeatures() ([]*service.Feature, error) {
	rows, err := s.stmts[FindFullAllFeatures].Query()
	if err != nil {
		return nil, exception.NewInternalServerError("query all feature list error, %s", err)
	}
	defer rows.Close()

	features := []*service.Feature{}
	for rows.Next() {
		f := service.Feature{}
		if err := rows.Scan(&f.ID, &f.Name, &f.Method, &f.Endpoint, &f.Description, &f.IsDeleted, &f.WhenDeletedVersion, &f.IsAdded, &f.WhenAddedVersion, &f.ServiceName); err != nil {
			return nil, exception.NewInternalServerError("scan service feature record error, %s", err)
		}
		features = append(features, &f)
	}

	return features, nil
}
