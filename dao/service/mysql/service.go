package mysql

import (
	"database/sql"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/defineiot/keyauth/dao/models"
	"github.com/defineiot/keyauth/internal/exception"
	"github.com/defineiot/keyauth/internal/tools"
)

func (s *store) CreateService(svr *models.Service) error {
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
func (s *store) ListServices() ([]*models.Service, error) {
	rows, err := s.stmts[FindAllServices].Query()
	if err != nil {
		return nil, exception.NewInternalServerError("query service list error, %s", err)
	}
	defer rows.Close()

	svrs := []*models.Service{}
	for rows.Next() {
		svr := new(models.Service)
		if err := rows.Scan(&svr.ID, &svr.Type, &svr.Name, &svr.Description, &svr.Enabled, &svr.Status,
			&svr.StatusUpdateAt, &svr.CurrentVersion, &svr.UpgradeVersion, &svr.DowngradeVersion,
			&svr.CreateAt, &svr.UpdateAt, &svr.ClientID, &svr.ClientSecret, &svr.TokenExpireTime); err != nil {
			return nil, exception.NewInternalServerError("scan service record error, %s", err)
		}
		svrs = append(svrs, svr)
	}

	return svrs, nil
}

func (s *store) GetServiceByID(id string) (*models.Service, error) {
	svr := models.Service{}
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

func (s *store) GetServiceByName(name string) (*models.Service, error) {
	svr := models.Service{}
	if err := s.stmts[FindServiceByName].QueryRow(name).Scan(&svr.ID, &svr.Type, &svr.Name, &svr.Description,
		&svr.Enabled, &svr.Status, &svr.StatusUpdateAt, &svr.CurrentVersion, &svr.UpgradeVersion,
		&svr.DowngradeVersion, &svr.CreateAt, &svr.UpdateAt, &svr.ClientID, &svr.ClientSecret,
		&svr.TokenExpireTime); err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.NewNotFound("service %s not find", name)
		}

		return nil, exception.NewInternalServerError("query single service client error, %s", err)
	}

	return &svr, nil
}

func (s *store) GetServiceByClientID(clientID string) (*models.Service, error) {
	svr := new(models.Service)
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

func (s *store) ListAllFeatures() ([]*models.Feature, error) {
	rows, err := s.stmts[FindAllFeatures].Query()
	if err != nil {
		return nil, exception.NewInternalServerError("query all feature list error, %s", err)
	}
	defer rows.Close()

	features := []*models.Feature{}
	for rows.Next() {
		f := new(models.Feature)
		if err := rows.Scan(&f.ID, &f.Name, &f.Tag, &f.HTTPEndpoint, &f.Description, &f.IsDeleted,
			&f.DeletedVersion, &f.DeleteAt, &f.IsAdded, &f.AddedVersion, &f.AddedAt, &f.ServiceID); err != nil {
			return nil, exception.NewInternalServerError("scan service feature record error, %s", err)
		}
		features = append(features, f)
	}

	return features, nil
}

func (s *store) ListServiceFeatures(serviceID string) ([]*models.Feature, error) {
	rows, err := s.stmts[FindServiceFeatures].Query(serviceID)
	if err != nil {
		return nil, exception.NewInternalServerError("query service feature list error, %s", err)
	}
	defer rows.Close()

	features := []*models.Feature{}
	for rows.Next() {
		f := new(models.Feature)
		if err := rows.Scan(&f.ID, &f.Name, &f.Tag, &f.HTTPEndpoint, &f.Description, &f.IsDeleted,
			&f.DeletedVersion, &f.DeleteAt, &f.IsAdded, &f.AddedVersion, &f.AddedAt, &f.ServiceID); err != nil {
			return nil, exception.NewInternalServerError("scan service feature record error, %s", err)
		}
		features = append(features, f)
	}

	return features, nil
}

func (s *store) ListRoleFeatures(roleID string) ([]*models.Feature, error) {
	rows, err := s.stmts[FindRoleFeatures].Query(roleID)
	if err != nil {
		return nil, exception.NewInternalServerError("query role features error, %s", err)
	}
	defer rows.Close()

	features := []*models.Feature{}
	for rows.Next() {
		f := new(models.Feature)
		if err := rows.Scan(&f.ID, &f.Name, &f.Tag, &f.HTTPEndpoint, &f.Description, &f.IsDeleted,
			&f.DeletedVersion, &f.DeleteAt, &f.IsAdded, &f.AddedVersion, &f.AddedAt, &f.ServiceID); err != nil {
			return nil, exception.NewInternalServerError("scan role feature mapping record error, %s", err)
		}
		features = append(features, f)
	}

	return features, nil
}

func (s *store) RegistryServiceFeatures(serviceID, version string, features ...*models.Feature) error {
	hasF, err := s.ListServiceFeatures(serviceID)
	if err != nil {
		return err
	}

	// 找出需要新增的功能(同一个服务下, 名称相同的功能)
	added := []*models.Feature{}
	for i := range features {
		exist := false
		for _, has := range hasF {
			if features[i].Name == has.Name && serviceID == has.ServiceID {
				exist = true
				break
			}
		}
		if !exist {
			// 处理新增加的功能
			newF := features[i]
			if newF.ID == "" {
				newF.ID = uuid.NewV4().String()
			}
			newF.IsAdded = true
			newF.AddedVersion = version
			newF.AddedAt = time.Now().Unix()
			added = append(added, newF)
		} else {

		}
	}
	s.Debug("added features: %s", added)

	// 找出需要删除的功能(之前存在, 而新注册的功能列表里面没有功能)
	deleted := []*models.Feature{}
	for _, has := range hasF {
		exist := false
		for i := range features {
			if has.Name == features[i].Name && has.ServiceID == serviceID {
				exist = true
				break
			}
		}

		if !exist {
			// 处理需要删除的功能
			has.IsDeleted = true
			has.DeletedVersion = version
			has.DeleteAt = time.Now().Unix()
			deleted = append(deleted, has)
		}
	}
	s.Debug("deleted features: %s", deleted)

	// 启动事物
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("start save features transaction error, %s", err)
	}
	// commit transaction
	defer func() {
		if err := tx.Commit(); err != nil {
			s.Error("feature transaction commit error, %s", err)
		}
	}()

	// 添加需要的新功能
	addFeaturePre, err := tx.Prepare(s.sql[SaveFeature])
	if err != nil {
		return exception.NewInternalServerError("prepare insert feature stmt error, name: %s, %s", serviceID, err)
	}
	defer addFeaturePre.Close()

	for _, f := range added {
		s.Info("service: %s add feature: %s", serviceID, f)
		if _, err := addFeaturePre.Exec(f.ID, f.Name, f.Tag, f.HTTPEndpoint, f.Description, f.IsDeleted,
			f.DeletedVersion, f.DeleteAt, f.IsAdded, f.AddedVersion, f.AddedAt, serviceID); err != nil {
			if err := tx.Rollback(); err != nil {
				s.Error("insert feature transaction rollback error, %s", err)
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
		s.Info("service: %s del feature: %s", serviceID, f.Name)
		_, err := delFeaturePre.Exec(f.IsAdded, f.DeletedVersion, f.DeleteAt, f.Name, f.ServiceID)
		if err != nil {
			if err := tx.Rollback(); err != nil {
				s.Error("update delete mark feature feature transaction rollback error, %s", err)
			}
			return exception.NewInternalServerError("update delete mark feature feature exec sql err, %s", err)
		}
	}

	return nil
}

func (s *store) AssociateFeaturesToRole(roleID string, features ...*models.Feature) error {
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
				s.Error("insert feature role mapping transaction rollback error, %s", err)
			}
			return exception.NewInternalServerError("insert feature role mapping exec sql err, %s", err)
		}

	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		s.Error("insert feature transaction rollback error, %s", err)
		return exception.NewInternalServerError("insert feature transaction commit error, but rollback success, %s", err)
	}

	return nil
}

func (s *store) UnlinkFeatureFromRole(roleID string, features ...*models.Feature) error {
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
				s.Error("unlik feature role mapping transaction rollback error, %s", err)
			}
			return exception.NewInternalServerError("unlik feature role mapping exec sql err, %s", err)
		}

	}

	// commit transaction
	if err := tx.Commit(); err != nil {
		s.Error("unlik feature transaction rollback error, %s", err)
		return exception.NewInternalServerError("unlik feature transaction commit error, but rollback success, %s", err)
	}

	return nil
}
