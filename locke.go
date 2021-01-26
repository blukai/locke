package locke

import (
	"sync"
)

// Locke interface.
type Locke interface {
	NewTxn(keys ...interface{}) Txn
}

type locke struct {
	wait      chan struct{}
	mu        sync.Mutex
	resources map[interface{}]struct{}
}

// New Locke.
func New() Locke {
	return &locke{
		resources: make(map[interface{}]struct{}),
	}
}

func (l *locke) NewTxn(keys ...interface{}) Txn {
	return &txn{
		keys: keys,
		l:    l,
	}
}

func (l *locke) lock(t *txn) {
	for {
		if wait := l.canLock(t); wait != nil {
			<-wait
		}
		return
	}
}

func (l *locke) canLock(t *txn) <-chan struct{} {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, key := range t.keys {
		if _, ok := l.resources[key]; ok {
			return l.exposeWait()
		}
	}

	for _, key := range t.keys {
		l.resources[key] = struct{}{}
	}

	return nil
}

func (l *locke) unlock(t *txn) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, key := range t.keys {
		delete(l.resources, key)
	}

	l.notifyWait()
}

func (l *locke) exposeWait() <-chan struct{} {
	if l.wait == nil {
		l.wait = make(chan struct{}, 0)
	}
	return l.wait
}

func (l *locke) notifyWait() {
	if l.wait != nil {
		close(l.wait)
		l.wait = nil
	}
}
