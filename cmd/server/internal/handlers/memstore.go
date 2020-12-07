package handlers

import (
	"errors"
	"sync"
)

// We have to store data, so expanding InMemoryStore with Stream storage. The storage is literaterally the same,
// except type of stored data. We can't simply add more fields to the structure,
// cause the different methods would share Counter and Mutex fields. Also, now adding new data types to store is much easier.

type buffMemStore struct {
	mu        *sync.RWMutex
	IdCounter uint64
	buffs     map[uint64]*Buff
}

type streamMemStore struct {
	mu        *sync.RWMutex
	IdCounter uint64
	streams   map[uint64]*Stream
}

type inMemStore struct {
	buffStore   *buffMemStore
	streamStore *streamMemStore
}

func (i *inMemStore) GetBuff(id uint64) (*Buff, error) {
	i.buffStore.mu.RLock()
	defer i.buffStore.mu.RUnlock()

	b, ok := i.buffStore.buffs[id]
	if !ok {
		return nil, errors.New("buff not found")
	}
	return b, nil
}

func (i *inMemStore) SetBuff(b *Buff) (uint64, error) {
	i.buffStore.mu.Lock()
	defer i.buffStore.mu.Unlock()

	i.buffStore.IdCounter++
	b.ID = i.buffStore.IdCounter
	i.buffStore.buffs[i.buffStore.IdCounter] = b
	return i.buffStore.IdCounter, nil
}

func (i *inMemStore) GetStream(id uint64) (*Stream, error) {
	i.streamStore.mu.RLock()
	defer i.streamStore.mu.RUnlock()

	s, ok := i.streamStore.streams[id]
	if !ok {
		return nil, errors.New("stream not found")
	}
	return s, nil

}

func (i *inMemStore) SetStream(s *Stream) (uint64, error) {
	i.streamStore.mu.Lock()
	defer i.streamStore.mu.Unlock()

	i.streamStore.IdCounter++
	s.ID = i.streamStore.IdCounter
	i.streamStore.streams[i.streamStore.IdCounter] = s
	return i.streamStore.IdCounter, nil

}

func (i *inMemStore) Count() (uint64, error) {
	return i.streamStore.IdCounter, nil

}

func (i *inMemStore) Streams() (map[uint64]*Stream, error) {
	return i.streamStore.streams, nil

}

func NewInMemStore() Store {
	return &inMemStore{
		&buffMemStore{
			mu:        &sync.RWMutex{},
			IdCounter: 0,
			buffs:     make(map[uint64]*Buff),
		},
		&streamMemStore{
			mu:        &sync.RWMutex{},
			IdCounter: 0,
			streams:   make(map[uint64]*Stream),
		},
	}
}
