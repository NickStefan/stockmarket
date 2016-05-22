package main

import (
	"gopkg.in/redsync.v1"
	"sync"
)

type Locker struct {
	redLock *redsync.Mutex
	mutLock *sync.Mutex
}

func (l *Locker) Lock() error {
	l.mutLock.Lock()
	err := l.redLock.Lock()
	if nil != err {
		l.mutLock.Unlock()
		return err
	}
	return nil
}

func (l *Locker) Unlock() {
	l.redLock.Unlock()
	l.mutLock.Unlock()
}
