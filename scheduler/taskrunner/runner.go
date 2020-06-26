package taskrunner

type Runner struct {
	Controller controlChan
	Error controlChan
	Data dataChan
	dataSize int
	longLived bool
	Dispatcher fn
	Executor fn
}

func NewRunner(dataSize int, longLived bool, dispatcher, executor fn) *Runner {
	return &Runner{
		Controller: make(chan string, 1), // 使用带缓冲区的chan，防止开始时就阻塞
		Error:      make(chan string, 1),
		Data:       make(chan interface{}, dataSize),
		dataSize:   dataSize,
		longLived:  longLived,
		Dispatcher: dispatcher,
		Executor:   executor,
	}
}

func (r *Runner) startDispatch() {
	defer func() {
		if !r.longLived { // 如果不是long live，则需要在dispatch退出时关闭所有通道
			close(r.Controller)
			close(r.Error)
			close(r.Data)
		}
	}()

	for {
		select { // non-blocking: 一旦有消息，则立即取出消息执行某个case
		case c := <- r.Controller:
			if c == READY_TO_DISPATCH {
				err := r.Dispatcher(r.Data)
				if err != nil {
					r.Error <- CLOSE
				} else {
					r.Controller <- READY_TO_EXECUTE
				}
			} else if c == READY_TO_EXECUTE {
				err := r.Executor(r.Data)
				if err != nil {
					r.Error <- CLOSE
				} else {
					r.Controller <- READY_TO_DISPATCH
				}
			}

		case e := <- r.Error:
			if e == CLOSE {
				return
			}

		default:
			
		}
	}
}

func (r *Runner) StartAll() {
	r.Controller <- READY_TO_DISPATCH
	r.startDispatch()
}
