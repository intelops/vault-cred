package policy

import (
	"sync"

	"github.com/intelops/vault-cred/internal/k8s"
)

type vaultConfigData struct {
	configMap map[string]k8s.ConfigMapData
	mutex     sync.Mutex
}

func newVaultConfigMapCache() vaultConfigData {
	return vaultConfigData{configMap: make(map[string]k8s.ConfigMapData)}
}

func (c *vaultConfigData) Put(key string, data k8s.ConfigMapData) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.configMap[key] = data
}

func (c *vaultConfigData) Get(key string) (k8s.ConfigMapData, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	data, ok := c.configMap[key]
	return data, ok
}
