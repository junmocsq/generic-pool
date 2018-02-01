package pool

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"testing"
	"time"
)

type DemoCloser struct {
	Name string
}

func (p *DemoCloser) Close() error {
	fmt.Println(p.Name, "closed")
	return nil
}

func TestNewGenericPool(t *testing.T) {
	_, err := NewGenericPool(0, 10, time.Minute*10, func() (io.Closer, error) {
		time.Sleep(time.Second)
		return &DemoCloser{Name: "test"}, nil
	})
	if err != nil {
		t.Error(err)
	}
}

func TestGenericPool_Acquire(t *testing.T) {
	pool, err := NewGenericPool(0, 5, time.Minute*10, func() (io.Closer, error) {
		time.Sleep(time.Second)
		name := strconv.FormatInt(time.Now().Unix(), 10)
		log.Printf("%s created", name)
		return &DemoCloser{Name: name}, nil
	})
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < 10; i++ {
		s, err := pool.Acquire()
		if err != nil {
			t.Error(err)
			return
		}
		pool.Release(s)
	}
}

func TestGenericPool_Shutdown(t *testing.T) {
	pool, err := NewGenericPool(0, 10, time.Minute*10, func() (io.Closer, error) {
		time.Sleep(time.Second)
		return &DemoCloser{Name: "test"}, nil
	})
	if err != nil {
		t.Error(err)
		return
	}
	if err := pool.Shutdown(); err != nil {
		t.Error(err)
		return
	}
	if _, err := pool.Acquire(); err != ErrPoolClosed {
		t.Error(err)
	}
}