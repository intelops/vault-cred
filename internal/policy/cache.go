package policy
import (
	"sync"
	"time"
)

// Define a struct to hold the Cache data and synchronization Mutex
type ConfigMapTimestampCache struct {
	Cache map[string]time.Time
	Mutex sync.Mutex
}
//var configMapTimestampCache = Cache.NewLRUExpireCache(1000)

// Initialize the Cache
var configMapTimestampCache = &ConfigMapTimestampCache{
	Cache: make(map[string]time.Time),
}

// configMapTimestampCache := &ConfigMapTimestampCache{
// 	Cache: make(map[string]time.Time),
// }

// Function to add or update the Cache entry
func (c *ConfigMapTimestampCache) AddOrUpdate(key string, timestamp time.Time) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	c.Cache[key] = timestamp
}

// Function to check if a Cache entry exists
func (c *ConfigMapTimestampCache) Exists(key string) bool {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	_, ok := c.Cache[key]
	return ok
}

// Function to retrieve the timestamp from the Cache
func (c *ConfigMapTimestampCache) Get(key string) (time.Time, bool) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	timestamp, ok := c.Cache[key]
	return timestamp, ok
}
