package policy

import (
	"sync"

	"github.com/intelops/vault-cred/internal/client"
)

type vaultConfigData struct {
	configMap map[string]client.ConfigMapData
	mutex     sync.Mutex
}

func newVaultConfigMapCache() vaultConfigData {
	return vaultConfigData{configMap: make(map[string]client.ConfigMapData)}
}

func (c *vaultConfigData) Put(key string, data client.ConfigMapData) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.configMap[key] = data
}

func (c *vaultConfigData) Get(key string) (client.ConfigMapData, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	data, ok := c.configMap[key]
	return data, ok
}
