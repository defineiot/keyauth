/* 资源服务(resource service) 需要将实例注册到etcd中,
注册格式如下:
key: /registry/instances/<service_id>/<instance_id>
value:
*/

package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/coreos/etcd/clientv3"

	"github.com/defineiot/keyauth/dao/instance"
	"github.com/defineiot/keyauth/internal/exception"
)

func (s *store) ListServiceInstances(serviceID string) ([]*instance.Instance, error) {
	key := fmt.Sprintf("%s/%s/", s.prefix, serviceID)
	instances := []*instance.Instance{}

	ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
	resp, err := s.client.Get(ctx, key, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	defer cancel()
	if err != nil {
		return nil, err
	}

	for _, ev := range resp.Kvs {
		keyS := strings.Split(string(ev.Key), "/")
		if len(keyS) != 5 && len(keyS) != 6 {
			return nil, fmt.Errorf("key parse error, %v", keyS)
		}

		ins := new(instance.Instance)
		if err := json.Unmarshal(ev.Value, ins); err != nil {
			return nil, fmt.Errorf("unmarshal resource value from etcd error, %s", err)
		}
		ins.ID = keyS[4]
		ins.ServiceID = keyS[3]

		instances = append(instances, ins)
	}

	return instances, nil
}

func (s *store) RegistryServiceInstance(ins *instance.Instance) error {
	if err := ins.Validate(); err != nil {
		return err
	}

	key := fmt.Sprintf("%s/%s/%s", s.prefix, ins.ServiceID, ins.ID)
	// 放入etcd时, 清除多余的数据
	insCopy := *ins
	insCopy.ServiceID = ""
	insCopy.ID = ""

	value, err := json.Marshal(insCopy)
	if err != nil {
		return exception.NewInternalServerError("marshal service instance to json error when save to etcd, %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
	_, err = s.client.Put(ctx, key, string(value))
	defer cancel()
	if err != nil {
		return err
	}

	return nil
}

func (s *store) UnRegisteServiceInstance(serviceID, instanceID string) error {
	key := fmt.Sprintf("%s/%s/%s", s.prefix, serviceID, instanceID)

	ctx, cancel := context.WithTimeout(context.Background(), s.requestTimeout)
	defer cancel()
	if _, err := s.client.Delete(ctx, key); err != nil {
		return err
	}

	return nil
}
