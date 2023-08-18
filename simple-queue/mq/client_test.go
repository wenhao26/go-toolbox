package mq

import (
	"fmt"
	"sync"
	"testing"
)

func TestClient(t *testing.T) {
	var wg sync.WaitGroup

	b := NewClient()
	b.SetConditions(100)

	for i := 0; i < 100; i++ {
		topic := fmt.Sprintf("test_topic_%d", i)
		payload := fmt.Sprintf("number_%d", i)

		ch, err := b.Subscribe(topic)
		if err != nil {
			t.Fatal(err)
		}

		wg.Add(1)
		go func() {
			e := b.GetPayload(ch)
			if e != payload {
				t.Fatalf("%s expected %s but get %s", topic, payload, e)
			}
			if err := b.Unsubscribe(topic, ch); err != nil {
				t.Fatal(err)
			}
			wg.Done()
		}()

		if err := b.Publish(topic, payload); err != nil {
			t.Fatal(err)
		}
	}
	wg.Wait()
}
