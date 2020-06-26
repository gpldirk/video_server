package main

// bucket token -> token 1 : N
// use connection -> assign token from bucket
// lose connection -> release token into bucket
// 因为bucket token有限，可以建立的connection也是有限的，进行流控
// bucket必须可以保证线程安全 -> channel

import (
	"log"
)

type ConnLimiter struct {
	concurrentConn int
	bucket chan int
}

// constructor
func NewConnLimiter(limit int) *ConnLimiter {
	return &ConnLimiter{
		concurrentConn: limit,
		bucket:         make(chan int, limit),
	}
}

// get token
func (connLimiter *ConnLimiter) GetConn() bool {
	if len(connLimiter.bucket) >= connLimiter.concurrentConn {
		log.Printf("Reached the max rate")
		return false
	}

	connLimiter.bucket <- 1
	return true
}

// release token
func (connLimiter *ConnLimiter) ReleaseConn() {
	c := <- connLimiter.bucket
	log.Printf("New connection coming: %d", c)
}



