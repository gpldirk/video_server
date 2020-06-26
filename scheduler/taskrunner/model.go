package taskrunner

const (
	READY_TO_DISPATCH = "D"
	READY_TO_EXECUTE = "E"
	CLOSE = "C"
)

type controlChan chan string

type dataChan chan interface{}

type fn func(dc dataChan) error
