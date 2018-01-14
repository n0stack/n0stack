package memcache

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/n0stack/n0core/message"
	"github.com/n0stack/n0core/model"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
)

type Memcache struct {
	MemcachedClient *memcache.Client
}

func (m Memcache) DigModel(id uuid.UUID, event string, depth uint) (model.AbstractModel, error) {
	i, err := m.MemcachedClient.Get(id.String())
	if err != nil {
		return nil, fmt.Errorf("")
	}

	mo := model.Model{}
	err = yaml.Unmarshal(i.Value, &mo)
	if err != nil {
		return nil, fmt.Errorf("")
	}

	pm, err := model.ParseYAMLModel(i.Value, mo.Type)
	if err != nil {
		return nil, err
	}

	if depth <= 0 {
		return pm, nil
	}

	for j, d := range pm.ToModel().Dependencies {
		pm.ToModel().Dependencies[j].Model, err = m.DigModel(d.Model.ToModel().ID, event, depth-1)
		if err != nil {
			return nil, err
		}
	}

	return pm, nil
}

func (m Memcache) StoreNotification(n *message.Notification) bool {
	y, err := yaml.Marshal(n.Model)
	if err != nil {
		return false
	}

	v := n.Model.ToModel()
	i := memcache.Item{Key: v.ID.String(), Value: y}
	m.MemcachedClient.Set(&i)

	return true
}
