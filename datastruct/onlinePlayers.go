package datastruct

import (
	"sync"
)

/*只统计执行过匹配的在线玩家们*/
type OnlinePlayers struct {
	Lock *sync.RWMutex   //读写互斥量
	Bm   map[string]bool //根据token保存
}

// NewOnlinePlayers return new OnlinePlayers
func NewOnlinePlayers() *OnlinePlayers {
	return &OnlinePlayers{
		Lock: new(sync.RWMutex),
		Bm:   make(map[string]bool),
	}
}

// Get from maps return the k's value
func (m *OnlinePlayers) Get(k string) (bool, bool) {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	val, ok := m.Bm[k]
	if ok {
		return val, ok
	}
	return false, ok
}

// Set Maps the given key and value. Returns false
// if the key is already in the map and changes nothing.
func (m *OnlinePlayers) Set(k string, v bool) bool {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	if val, ok := m.Bm[k]; !ok {
		m.Bm[k] = v
	} else if val != v {
		m.Bm[k] = v
	} else {
		return false
	}
	return true
}

// Check Returns true if k is exist in the map.
func (m *OnlinePlayers) Check(k string) (bool, bool) {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	v, ok := m.Bm[k]
	return v, ok
}

func (m *OnlinePlayers) IsExist(k string) bool {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	_, ok := m.Bm[k]
	return ok
}

// Delete the given key and value.
func (m *OnlinePlayers) Delete(k string) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	delete(m.Bm, k)
}

// Items returns all items in safemap.
func (m *OnlinePlayers) Items() map[string]bool {
	m.Lock.RLock()
	r := make(map[string]bool)
	for k, v := range m.Bm {
		r[k] = v
	}
	m.Lock.RUnlock()
	return r
}

// Count returns the number of items within the map.
func (m *OnlinePlayers) Count() int {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	return len(m.Bm)
}
