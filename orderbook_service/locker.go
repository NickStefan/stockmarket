package main

import (
	"fmt"
	"gopkg.in/redsync.v1"
	"sync"
)

type Locker struct {
	name    string
	redLock *redsync.Mutex
	mutLock *sync.Mutex
}

func (l *Locker) Lock() error {
	l.mutLock.Lock()
	err := l.redLock.Lock()
	if nil != err {
		fmt.Println("lock err", l.name)
		l.mutLock.Unlock()
		return err
	}
	return nil
}

func (l *Locker) Unlock() {
	l.redLock.Unlock()
	l.mutLock.Unlock()
}
