package taskrunner

import (
	"testing"
	"log"
	"time"
	"errors"
)

func TestRunner(t *testing.T) {
	d := func(dc dataChan) error {
		for i := 0; i < 30; i++ {
			dc <- i
			log.Printf("dispatcher sent %v", i)
		}
		return nil
	}

	e := func(dc dataChan) error {
		forloop:
			for {
				select {
				case data := <- dc:
					log.Printf("Executor received: %v", data)
				default:
					break forloop // dispatcher 发送数据到dataChan，executor接收所有数据之后就可以退出了
				}
			}
		return errors.New("Executor finished")
	}

	runner := NewRunner(30, false, d, e)
	go runner.StartAll()
	time.Sleep(3 * time.Second)
}
