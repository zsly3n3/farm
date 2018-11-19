package datastruct

import (
	"sync"
)

type OnlinePlayer struct {
}

/*只统计执行过匹配的在线玩家们*/
type OnlinePlayers struct {
	Lock *sync.RWMutex           //读写互斥量
	Bm   map[string]OnlinePlayer //map[int]*Player 根据Id保存
}

// NewOnlinePlayers return new OnlinePlayers
func NewOnlinePlayers() *OnlinePlayers {
	return &OnlinePlayers{
		Lock: new(sync.RWMutex),
		Bm:   make(map[string]OnlinePlayer),
	}
}

// Get from maps return the k's value
func (m *OnlinePlayers) Get(k string) (OnlinePlayer, bool) {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	val, ok := m.Bm[k]
	if ok {
		return val, ok
	}
	return OnlinePlayer{}, ok
}

// func (m *OnlinePlayers) GetAndUpdateState(key []string,state PlayerEnterType,room_id string) []Player {
// 	m.lock.RLock()
// 	defer m.lock.RUnlock()
// 	if val, ok := m.bm[k]; ok {
// 		return val
// 	}
// 	return Player{}
// }

// func (m *OnlinePlayers) GetWithAddr(addr string) (string,*Player) {
// 	key:=NULLKEY
// 	var pl *Player = nil
// 	m.lock.RLock()
// 	for k, v := range m.bm {
// 		str:=v.Agent.RemoteAddr().String()
// 		if str == addr{
// 			key = k
// 			pl = v
// 			break
// 		}
// 	}
// 	m.lock.RUnlock()
// 	return key,pl
// }

// Set Maps the given key and value. Returns false
// if the key is already in the map and changes nothing.
// func (m *OnlinePlayers) Set(k string, v Player) bool {
// 	m.Lock.Lock()
// 	defer m.Lock.Unlock()
// 	if val, ok := m.Bm[k]; !ok {
// 		m.Bm[k] = v
// 	} else if val != v {
// 		m.Bm[k] = v
// 	} else {
// 		return false
// 	}
// 	return true
// }

// Check Returns true if k is exist in the map.
func (m *OnlinePlayers) Check(k string) (OnlinePlayer, bool) {
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
func (m *OnlinePlayers) Items() map[string]OnlinePlayer {
	m.Lock.RLock()
	r := make(map[string]OnlinePlayer)
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
