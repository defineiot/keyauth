package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/satori/go.uuid"

	"github.com/defineiot/keyauth/dao/service"
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/defineiot/keyauth/internal/tools"
)

func (s *store) CreateService(svr *service.Service) error {
	if err := svr.Validate(); err != nil {
		return err
	}

	svr.CreateAt = time.Now().Unix()
	svr.ID = uuid.NewV4().String()
	svr.ClientID = tools.MakeUUID(24)
	svr.ClientSecret = tools.MakeUUID(32)

	if _, err := s.stmts[SaveService].Exec(svr.ID, string(svr.Type), svr.Name, svr.Description,
		svr.Enabled, svr.CreateAt, svr.ClientID, svr.ClientSecret, svr.TokenExpireTime); err != nil {
		return exception.NewInternalServerError("insert service exec sql err, %s", err)
	}

	return nil
}
func (s *store) ListServices() ([]*service.Service, error) {
	rows, err := s.stmts[FindAllServices].Query()
	if err != nil {
		return nil, exception.NewInternalServerError("query service list error, %s", err)
	}
	defer rows.Close()

	svrs := []*service.Service{}
	for rows.Next() {
		svr := new(service.Service)
		if err := rows.Scan(&svr.ID, &svr.Type, &svr.Name, &svr.Description, &svr.Enabled, &svr.Status,
			&svr.StatusUpdateAt, &svr.CurrentVersion, &svr.UpgradeVersion, &svr.DowngradeVersion,
			&svr.CreateAt, &svr.UpdateAt, &svr.ClientID, &svr.ClientSecret, &svr.TokenExpireTime); err != nil {
			return nil, exception.NewInternalServerError("scan service record error, %s", err)
		}
		svrs = append(svrs, svr)
	}

	return svrs, nil
}

func (s *store) GetServiceByID(id string) (*service.Service, error) {
	svr := service.Service{}
	if err := s.stmts[FindServiceByID].QueryRow(id).Scan(&svr.ID, &svr.Type, &svr.Name, &svr.Description,
		&svr.Enabled, &svr.Status, &svr.StatusUpdateAt, &svr.CurrentVersion, &svr.UpgradeVersion,
		&svr.DowngradeVersion, &svr.CreateAt, &svr.UpdateAt, &svr.ClientID, &svr.ClientSecret,
		&svr.TokenExpireTime); err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("service %s not find", id)
		}

		return nil, exception.NewInternalServerError("query single service client error, %s", err)
	}

	return &svr, nil
}

func (s *store) GetServiceByClientID(clientID string) (*service.Service, error) {
	svr := new(service.Service)
	if err := s.stmts[FindServiceByClient].QueryRow(clientID).Scan(&svr.ID, &svr.Type, &svr.Name, &svr.Description,
		&svr.Enabled, &svr.Status, &svr.StatusUpdateAt, &svr.CurrentVersion, &svr.UpgradeVersion,
		&svr.DowngradeVersion, &svr.CreateAt, &svr.UpdateAt, &svr.ClientID, &svr.ClientSecret,
		&svr.TokenExpireTime); err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("client %s service not find", clientID)
		}

		return nil, exception.NewInternalServerError("query single service client error, %s", err)
	}

	return svr, nil
}

func (s *store) DeleteService(id string) error {
	// 清除服务的功能列表
	_, err := s.stmts[DeleteServiceFeatures].Exec(id)
	if err != nil {
		return exception.NewInternalServerError("delete service features exec sql error, %s", err)
	}

	// 清除服务
	ret, err := s.stmts[DeleteService].Exec(id)
	if err != nil {
		return exception.NewInternalServerError("delete service exec sql error, %s", err)
	}
	count, err := ret.RowsAffected()
	if err != nil {
		return exception.NewInternalServerError("get delete service row affected error, %s", err)
	}
	if count == 0 {
		return exception.NewBadRequest("service %s not exist", id)
	}

	return nil
}

func (s *store) ListServiceFeatures(serviceID string) ([]*service.Feature, error) {
	rows, err := s.stmts[FindAllFeatures].Query(serviceID)
	if err != nil {
		return nil, exception.NewInternalServerError("query service feature list error, %s", err)
	}
	defer rows.Close()

	features := []*service.Feature{}
	for rows.Next() {
		f := new(service.Feature)
		if err := rows.Scan(&f.ID, &f.Name, &f.Tag, &f.HTTPEndpoint, &f.Description, &f.IsDeleted,
			&f.WhenDeletedVersion, &f.IsAdded, &f.WhenAddedVersion, &f.ServiceID); err != nil {
			return nil, exception.NewInternalServerError("scan service feature record error, %s", err)
		}
		features = append(features, f)
	}

	return features, nil
}

func (s *store) ListRoleFeatures(roleID string) ([]*service.Feature, error) {
	rows, err := s.stmts[GetRoleFeatures].Query(roleID)
	if err != nil {
		return nil, exception.NewInternalServerError("query role features error, %s", err)
	}
	defer rows.Close()

	features := []*service.Feature{}
	for rows.Next() {
		f := new(service.Feature)
		if err := rows.Scan(&f.ID, &f.Name, &f.Tag, &f.HTTPEndpoint, &f.Description, &f.IsDeleted,
			&f.WhenDeletedVersion, &f.IsAdded, &f.WhenAddedVersion, &f.ServiceID); err != nil {
			return nil, exception.NewInternalServerError("scan role feature mapping record error, %s", err)
		}
		features = append(features, f)
	}

	return features, nil
}

func (s *store) RegistryServiceFeatures(serviceID string, features ...*service.Feature) error {
	s.Debugf("registry service :%s features: %v", serviceID, features)

	hasF, err := s.ListServiceFeatures(serviceID)
	if err != nil {
		return err
	}

	// 找出需要新增的功能
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
			added = append(added, features[i])
		}
	}

	// 找出需要删除的功能
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

	// 启动事物
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("start save features transaction error, %s", err)
	}
	// commit transaction
	defer func() {
		if err := tx.Commit(); err != nil {
			s.Errorf("feature transaction commit error, %s", err)
		}
	}()

	// 添加需要的新功能
	addFeaturePre, err := tx.Prepare(s.sql[SaveFeature])
	if err != nil {
		return exception.NewInternalServerError("prepare insert feature stmt error, name: %s, %s", serviceID, err)
	}
	defer addFeaturePre.Close()

	for _, f := range added {
		s.Infof("service: %s add feature: %s", serviceID, f.Name)
		if _, err := addFeaturePre.Exec(f.ID, f.Name, f.Tag, f.HTTPEndpoint, f.Description, f.IsDeleted,
			f.WhenDeletedVersion, f.IsAdded, f.WhenAddedVersion, serviceID); err != nil {
			if err := tx.Rollback(); err != nil {
				s.Errorf("insert feature transaction rollback error, %s", err)
			}
			return exception.NewInternalServerError("insert feature exec sql err, %s", err)
		}
	}

	// 标记需要删除的功能
	delFeaturePre, err := tx.Prepare(s.sql[MarkDeleteFeature])
	if err != nil {
		return exception.NewInternalServerError("prepare update delete mark feature stmt error, name: %s, %s", serviceID, err)
	}
	defer delFeaturePre.Close()

	for _, f := range deleted {
		s.Infof("service: %s del feature: %s", serviceID, f.Name)
		_, err := delFeaturePre.Exec(1, f.Name, serviceID)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				s.Errorf("update delete mark feature feature transaction rollback error, %s", err)
			}
			return exception.NewInternalServerError("update delete mark feature feature exec sql err, %s", err)
		}
	}

	return nil
}

func (s *store) AssociateFeaturesToRole(roleID string, features ...*service.Feature) error {
	// start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("start associate features to role transaction error, %s", err)
	}

	// prepare insert feature
	mappingPre, err := tx.Prepare(s.sql[AssociateFeaturesToRole])
	if err != nil {
		return exception.NewInternalServerError("prepare insert feature role mapping stmt error, name: %s, %s", roleID, err)
	}
	defer mappingPre.Close()

	for _, f := range features {
		_, err := mappingPre.Exec(f.ID, roleID)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				s.Errorf("insert feature role mapping transaction rollback error, %s", err)
			}
			return exception.NewInternalServerError("insert feature role mapping exec sql err, %s", err)
		}

	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		s.Errorf("insert feature transaction rollback error, %s", err)
		return exception.NewInternalServerError("insert feature transaction commit error, but rollback success, %s", err)
	}

	return nil
}

func (s *store) UnlinkFeatureFromRole(roleID string, features ...*service.Feature) error {
	// start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("start unlink features from role transaction error, %s", err)
	}

	// prepare insert feature
	mappingPre, err := tx.Prepare(s.sql[UnlinkFeatureFromRole])
	if err != nil {
		return exception.NewInternalServerError("prepare insert feature role mapping stmt error, name: %s, %s", roleID, err)
	}
	defer mappingPre.Close()

	for _, f := range features {
		_, err := mappingPre.Exec(f.ID, roleID)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				s.Errorf("unlik feature role mapping transaction rollback error, %s", err)
			}
			return exception.NewInternalServerError("unlik feature role mapping exec sql err, %s", err)
		}

	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		s.Errorf("unlik feature transaction rollback error, %s", err)
		return exception.NewInternalServerError("unlik feature transaction commit error, but rollback success, %s", err)
	}

	return nil
}
