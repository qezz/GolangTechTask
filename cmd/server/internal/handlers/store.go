package handlers

type BuffStore interface {
	GetBuff(id uint64) (*Buff, error)
	SetBuff(*Buff) (id uint64, err error)
}

type StreamStore interface {
	GetStream(id uint64) (*Stream, error)
	SetStream(s *Stream) (uint64, error)
	Count() (uint64, error)
	Streams() (map[uint64]*Stream, error)
}

type Store interface {
	BuffStore
	StreamStore
}
