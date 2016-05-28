package main

import (
	"fmt"
	"gopkg.in/redsync.v1"
	"sync"
)

type Locker struct {
	env     string
	name    string
	redLock *redsync.Mutex
	mutLock *sync.Mutex
}

func (l *Locker) Lock() error {
	l.mutLock.Lock()
	if "TESTING" == l.env {
		return nil
	}
	err := l.redLock.Lock()
	if nil != err {
		fmt.Println(err, l.name)
		l.mutLock.Unlock()
		return err
	}
	return nil
}

func (l *Locker) Unlock() {
	if "TESTING" != l.env {
		l.redLock.Unlock()
	}
	l.mutLock.Unlock()
}
