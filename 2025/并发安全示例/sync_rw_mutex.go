package main

import (
	"fmt"
	"sync"
)

type DataStore struct {
	data map[string]string
	rwmu sync.RWMutex
}

func (ds *DataStore) SetData(key, value string) {
	ds.rwmu.Lock()
	defer ds.rwmu.Unlock()
	ds.data[key] = value
}

func (ds *DataStore) GetData(key string) string {
	ds.rwmu.RLock()
	defer ds.rwmu.RUnlock()
	return ds.data[key]
}

func main() {
	dataStore := DataStore{
		data: make(map[string]string),
	}
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			dataStore.SetData(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
			defer wg.Done()
		}(i)
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			fmt.Println(dataStore.GetData(fmt.Sprintf("key%d", i)))
			wg.Done()
		}(i)
	}

	wg.Wait()

}
