package main

import (
	"fmt"
	"time"
)

type RingBuffer struct {
	buff   []int            // слайс, буфер
	size   int              // его размер
	getCh  chan *int        // канал для получения значения (*int чтобы менять значение напрямую)
	tmp    chan int         // канал для ответа по таймеру
	putCh  chan int         // для записи в буфер
	wInd   int              // индекс записи, нужен для кольцевого буфера
	rInd   int              // индекс чтения, нужен для кольцевого буфера
	done   chan struct{}    // "сигнал" остановки loop цикла-селектора каналов буфера
	waitCh chan struct{}    // без этой штуки идет рассинхрон в запись переменной, нужно будет перепроверить еще раз, как без нее возможно
	timer  <-chan time.Time // канал тикания, время тика устанавливается в конструкторе.
}

func NewRingBuffer(size int, seconds time.Duration) *RingBuffer {
	return &RingBuffer{
		buff:   make([]int, size),
		size:   size,
		getCh:  make(chan *int),
		tmp:    make(chan int),
		putCh:  make(chan int),
		wInd:   0,
		rInd:   -1,
		done:   make(chan struct{}),
		waitCh: make(chan struct{}),
		timer:  time.Tick(seconds),
	}
}

// Get - функция чтения из буфера. Потокобезопасна, тк используется небуферизированный канал, для доступа: r.getCh
func (r *RingBuffer) Get() (int, error) {
	// fmt.Println("before get",r.buff,"r.rInd",r.rInd,r.wInd)
	if len(r.buff) == 0 || r.rInd == -1 {
		return 0, fmt.Errorf("buffer is empty1")
	}

	if r.rInd == r.wInd && r.wInd == -1 {
		return 0, fmt.Errorf("buffer is empty2")
	}

	x := r.buff[r.rInd%r.size]
	r.rInd = (r.rInd + 1) % r.size
	if r.rInd == r.wInd {
		r.rInd = -1
	}
	// fmt.Println("after get",r.buff,"r.rInd",r.rInd,r.wInd,"x",x)
	return x, nil
}

// Put - функция записи в буфер.
func (r *RingBuffer) Put(x int) {
	// fmt.Println("before put",r.buff,"r.rInd",r.rInd,r.wInd)
	if len(r.buff) == 0 {
		return
	}

	r.buff[r.wInd] = x

	if r.wInd == r.rInd {
		r.rInd = (r.rInd + 1) % r.size
	}

	if r.rInd == -1 {
		r.rInd = r.wInd
	}
	r.wInd = (r.wInd + 1) % r.size
	// fmt.Println("after put",r.buff,"r.rInd",r.rInd,r.wInd)
}

// loop() - селектор каналов буфера
func (r *RingBuffer) loop() {
	for {
		select {
		case x := <-r.getCh: // канал принимает указатели, чтобы изменить значение переменной которая поступает
			y, err := r.Get()
			if err == nil {
				*x = y
			}
			r.waitCh <- struct{}{} // без этой блокировки, синхронизации, переключается контекст горутин, и запись, присваивание x нового значения не отрабатывает корректно.
		case x := <-r.putCh:
			// fmt.Println("puting",x)
			r.Put(x)
		case <-r.done:
			return
		case <-r.timer:
			x, err := r.Get()
			if err != nil {
				fmt.Println("err",err)
			} else {
				// fmt.Println("Send x",x)
				r.tmp <- x // если нет ошибок, например канал не пустой, то будет передано сообщение по каналу
			}
		}
	}

}
