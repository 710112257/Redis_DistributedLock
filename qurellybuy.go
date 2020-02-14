package main

import (
	"sync"
	"sync/atomic"
)

type Food struct {
	name   string
	number int32
	sync.Mutex
	deal chan int
}

func NewFoods() Food {
	return Food{
		"apple",
		1000,
		sync.Mutex{},
		make(chan int, 1),
	}
}
func (this *Food) get() int32 {
	var getnumber int32 = 1
	if !this.surple() {
		return 0
	}
	this.Lock()
	defer this.Unlock()

	atomic.AddInt32(&this.number, int32(-1))
	return getnumber
}
func (this *Food) surple() bool {
	if this.number > 0 {
		return true
	}
	return false
}
func (this *Food) total() int32 {
	return this.number
}
