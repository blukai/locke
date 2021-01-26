package locke

import "sync"

// Txn interface.
type Txn interface {
	sync.Locker
	CanLock() bool
}

type txn struct {
	l    *locke
	keys []interface{}
}

func (t *txn) Lock() {
	t.l.lock(t)
}

func (t *txn) CanLock() bool {
	return t.l.canLock(t) == nil
}

func (t *txn) Unlock() {
	t.l.unlock(t)
}
