package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimer1second(t *testing.T) {

	sizeRingBuffer := 3
	rb := NewRingBuffer(sizeRingBuffer, time.Duration(time.Second*1))

	rb.Put(1)
	rb.Put(2)
	// rb.Put(3)

	go rb.loop()
	time.Sleep(time.Second * 15)

	// x, err := rb.Get()
	// if err != nil {
	// 	t.Log("err", err)

	// }

	// t.Log("x", x)
}

func TestPut4size3(t *testing.T) {

	sizeRingBuffer := 3
	rb := NewRingBuffer(sizeRingBuffer, time.Duration(time.Second*30))

	go func() {
		go rb.loop()
		time.Sleep(time.Second * 3)
		rb.done <- struct{}{}
	}()

	rb.putCh <- 1
	// <-rb.waitCh
	time.Sleep(time.Second)
	var x int
	rb.getCh <- &x

	<-rb.waitCh
	time.Sleep(time.Second)
	assert.Equal(t, 1, x)
}

func TestPut4size4(t *testing.T) {

	sizeRingBuffer := 4
	rb := NewRingBuffer(sizeRingBuffer, time.Duration(time.Second*30))

	rb.Put(1)
	rb.Put(2)
	rb.Put(3)
	rb.Put(4)

	for i := 1; i <= 4; i++ {
		out, err := rb.Get()
		if err != nil {
			t.Log("err", err)
		}
		assert.Equal(t, i, out)
	}

	a, err := rb.Get()
	assert.Equal(t,"buffer is empty1", err.Error() )
	assert.Equal(t, 0, a)
}

func TestPutsize3Get(t *testing.T) {

	sizeRingBuffer := 3
	rb := NewRingBuffer(sizeRingBuffer, time.Duration(time.Second*30))

	rb.Put(1)
	rb.Put(2)
	rb.Put(3)
	rb.Put(4)

	out, err := rb.Get()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 2, out)

	out, err = rb.Get()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 3, out)

	out, err = rb.Get()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 4, out)

	out, err = rb.Get() // ошибка неокуда читать достигнут wInd
	if err != nil {
		assert.Equal(t, "buffer is empty1", err.Error())
	}
	assert.Equal(t, 0, out)

	rb.Put(5) // 4 5 3

	out, err = rb.Get()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 5, out)

	out, err = rb.Get() // ошибка неокуда читать достигнут wInd
	if err != nil {
		assert.Equal(t, "buffer is empty1", err.Error())
	}
	assert.Equal(t, 0, out)

	rb.Put(6) // 4 5 6
	rb.Put(7) // 7 5 6

	out, err = rb.Get()
	if err != nil {
		assert.Equal(t, err, out)
	}
	assert.Equal(t, 6, out)

	out, err = rb.Get()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 7, out)

}

func TestPut3Get3Ch(t *testing.T) {

	sizeRingBuffer := 3
	rb := NewRingBuffer(sizeRingBuffer, time.Duration(time.Second*30))

	go func() {
		go rb.loop()
		time.Sleep(time.Second * 3)
		rb.done <- struct{}{}
	}()

	rb.putCh <- 1
	rb.putCh <- 2
	rb.putCh <- 3
	var x int
	rb.getCh <- &x
	<-rb.waitCh
	assert.Equal(t, 1, x)

	rb.getCh <- &x
	<-rb.waitCh
	assert.Equal(t, 2, x)

	rb.getCh <- &x
	<-rb.waitCh
	assert.Equal(t, 3, x)

}

func TestPut3FuncGet3Ch(t *testing.T) {

	sizeRingBuffer := 3
	rb := NewRingBuffer(sizeRingBuffer, time.Duration(time.Second*30))

	rb.Put(1)
	rb.Put(2)
	rb.Put(3)

	go rb.loop()

	var x int
	rb.getCh <- &x
	<-rb.waitCh
	assert.Equal(t, 1, x)

	rb.getCh <- &x
	<-rb.waitCh
	assert.Equal(t, 2, x)

	rb.getCh <- &x

	<-rb.waitCh
	assert.Equal(t, 3, x)
}
