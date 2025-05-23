package rabbitmq_bak

import (
	"errors"

	"github.com/sirupsen/logrus"
)

type Clients struct {
	Url       string
	Ex        *Exchange
	MaxClient int
	MqClients chan *MyMQ
}

func NewClients(url string, exchange *Exchange, maxClient int) (*Clients, error) {
	clients := Clients{
		Url:       url,
		Ex:        exchange,
		MaxClient: maxClient,
		MqClients: make(chan *MyMQ, maxClient),
	}

	for i := 0; i < maxClient; i++ {
		client, err := newClient(url, exchange)
		if err != nil {
			logrus.Errorf("connect to %s mq err: %s", url, err.Error())
			return nil, err
		}
		clients.MqClients <- client
	}
	return &clients, nil
}

func newClient(url string, exchange *Exchange) (*MyMQ, error) {
	client, err := New(&Config{
		Addr:     url,
		Exchange: exchange,
	})
	return client, err
}

func (c *Clients) Push(q *Queue, msg *Message, exchanges ...*Exchange) error {
	client, err := c.Get()
	defer c.Put(client)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return client.Push(q, msg, exchanges...)
}

func (c *Clients) Pub(routingKey string, msg *Message, exchanges ...*Exchange) error {
	client, err := c.Get()
	defer c.Put(client)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return client.Pub(routingKey, msg, exchanges...)
}

func (c *Clients) Sub(q *Queue) (<-chan *Message, error) {
	client, err := c.Get()
	defer c.Put(client)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return client.Sub(q)
}

func (c *Clients) Get() (*MyMQ, error) {
	if c.MqClients == nil {
		return nil, errors.New("Rpc Clients is nil ")
	}
	return <-c.MqClients, nil
}

func (c *Clients) Put(client *MyMQ) {
	c.MqClients <- client
}

func (c *Clients) Close() error {
	for c.MaxClient > 0 {
		client, err := c.Get()
		if err != nil {
			logrus.Error(err)
			return err
		}
		client.Close()
		c.MaxClient--
	}
	return nil
}
