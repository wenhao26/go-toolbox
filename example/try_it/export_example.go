package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Mock struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}

type Result struct {
	Err error
}

type InputCh chan *Mock
type OutputCh chan *Result

func readData() InputCh {
	rand.Seed(time.Now().UnixNano())

	result := make(InputCh)
	go func() {
		defer close(result)

		for i := 0; i < 100000; i++ {
			mock := genMock(i)
			//fmt.Println(mock)
			result <- &mock
			//time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
		}
	}()
	return result
}

func genMock(i int) Mock {
	return Mock{
		Id:      i,
		Message: time.Now().Format("2006-01-02 15:04:05.000000000"),
	}
}

func writeData(in InputCh) OutputCh {
	out := make(OutputCh)
	go func() {
		defer close(out)
		for d := range in {
			out <- &Result{Err: saveMock(d)}
		}
	}()
	return out
}

func saveMock(d *Mock) error {
	rand.Seed(time.Now().UnixNano())
	//time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)
	return nil
}

type readCmd func() InputCh
type writeCmd func(in InputCh) OutputCh

func pipeHandler(rc readCmd, wc ...writeCmd) OutputCh {
	wg := sync.WaitGroup{}

	data := rc()
	out := make(OutputCh)
	for _, c := range wc {
		output := c(data)
		wg.Add(1)
		go func(outData OutputCh) {
			defer wg.Done()
			for i := range outData {
				out <- i
			}
		}(output)
	}
	go func() {
		defer close(out)
		wg.Wait()
	}()
	return out
}

func main() {
	start := time.Now().UnixNano() / 1e6
	out := pipeHandler(readData, writeData, writeData)
	for res := range out {
		fmt.Printf("Result=%v\n", res.Err)
	}
	end := time.Now().UnixNano() / 1e6
	fmt.Printf("测试--用时:%d秒\r\n", (end-start)/1000)
}
