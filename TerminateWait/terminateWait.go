package TerminateWait

import (
	"container/list"
	log "github.com/Sirupsen/logrus"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	forWait     list.List
	lock        sync.RWMutex
	quitingList list.List
	onceTer     sync.Once
)

type CanWait interface {
	Status() string
	Stop() <-chan struct{}
	Name() string
}

type QuitingDescriptor struct {
	CanWait
	Finished <-chan struct{}
}

func Add(cw CanWait) {
	lock.Lock()
	defer lock.Unlock()
	forWait.PushBack(cw)
}

func waitQuit() {
	for ele := quitingList.Front(); ele != nil; ele = ele.Next() {
		<-ele.Value.(*QuitingDescriptor).Finished
	}
	os.Exit(0)
}

func printQuitingStatus() {
	for {
		log.Warn("waiting for the below tasks quitting...")
		for ele := quitingList.Front(); ele != nil; ele = ele.Next() {
			qd := ele.Value.(*QuitingDescriptor)
			select {
			case <-qd.Finished:
				continue
			default:
				log.Warnf("%s : in status [%s]", qd.Name(), qd.Status())
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func Terminate() {
	lock.RLock()
	defer lock.RUnlock()

	for ele := forWait.Front(); ele != nil; ele = ele.Next() {
		quitingList.PushBack(&QuitingDescriptor{ele.Value.(CanWait), ele.Value.(CanWait).Stop()})
	}
	go waitQuit()
	go printQuitingStatus()
}

func init() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		onceTer.Do(Terminate)
	}()
}
