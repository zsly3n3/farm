package datastruct

import (
	"sync"
)

type OnlinePlayerData struct {
	OnlineTime int64
	WillDelete bool
}

/*只统计执行过匹配的在线玩家们*/
type OnlinePlayers struct {
	Lock *sync.RWMutex               //读写互斥量
	Bm   map[string]OnlinePlayerData //根据token保存
}

// NewOnlinePlayers return new OnlinePlayers
func NewOnlinePlayers() *OnlinePlayers {
	return &OnlinePlayers{
		Lock: new(sync.RWMutex),
		Bm:   make(map[string]OnlinePlayerData),
	}
}

// Get from maps return the k's value
func (m *OnlinePlayers) Get(k string) (*OnlinePlayerData, bool) {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	val, ok := m.Bm[k]
	if ok {
		return &val, ok
	}
	return nil, ok
}

// Set Maps the given key and value. Returns false
// if the key is already in the map and changes nothing.
func (m *OnlinePlayers) Set(k string, v *OnlinePlayerData) bool {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	if val, ok := m.Bm[k]; !ok {
		m.Bm[k] = *v
	} else if val != *v {
		m.Bm[k] = *v
	} else {
		return false
	}
	return true
}
