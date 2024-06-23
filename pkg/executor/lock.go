package executor

import "sync"

type pathLock struct {
	mu    sync.Mutex
	locks map[string]*sync.RWMutex
}

func newPathLock() *pathLock {
	return &pathLock{
		locks: make(map[string]*sync.RWMutex),
	}
}

func (pl *pathLock) Lock(path string) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	lock, ok := pl.locks[path]
	if !ok {
		lock = &sync.RWMutex{}
		pl.locks[path] = lock
	}

	lock.Lock()
}

func (pl *pathLock) RLock(path string) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	lock, ok := pl.locks[path]
	if !ok {
		lock = &sync.RWMutex{}
		pl.locks[path] = lock
	}

	lock.RLock()
}

func (pl *pathLock) Unlock(path string) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	lock, ok := pl.locks[path]
	if ok {
		lock.Unlock()
	}
}

func (pl *pathLock) RUnlock(path string) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	lock, ok := pl.locks[path]
	if ok {
		lock.RUnlock()
	}
}

var plock = newPathLock()
