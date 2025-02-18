package connection

import (
	"github.com/kattana-io/mexc-golang-sdk/websocket/types"
	"sync"
)

type Subscribes struct {
	m map[string]mexcwstypes.OnReceive
	sync.RWMutex
}

func NewSubs() *Subscribes {
	return &Subscribes{
		m: map[string]mexcwstypes.OnReceive{},
	}
}

func (s *Subscribes) Add(req string, listener mexcwstypes.OnReceive) {
	s.Lock()
	defer s.Unlock()

	s.m[req] = listener
}

func (s *Subscribes) Remove(req string) {
	s.Lock()
	defer s.Unlock()

	delete(s.m, req)
}

func (s *Subscribes) Load(req string) (mexcwstypes.OnReceive, bool) {
	s.RLock()
	defer s.RUnlock()

	v, ok := s.m[req]

	return v, ok
}

func (s *Subscribes) Len() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.m)
}

func (s *Subscribes) GetAllChannels() []string {
	s.RLock()
	defer s.RUnlock()

	channels := make([]string, 0)
	for ch := range s.m {
		channels = append(channels, ch)
	}
	return channels
}
