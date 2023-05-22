package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/RichardKnop/machinery/example/tracers"
	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/log"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/urfave/cli"

	"toolbox/example/machinery/mq/exampletasks"
)

var (
	app *cli.App
)

func init() {
	app = cli.NewApp()
	app.Name = "machinery"
	app.Usage = "machinery worker and send example tasks with machinery send"
	app.Version = "0.0.1"
}

func main() {
	app.Commands = []cli.Command{
		{
			Name:  "send",
			Usage: "send example tasks",
			Action: func(c *cli.Context) error {
				fmt.Println("发送任务例子")
				if err := send(); err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				return nil
			},
		},
		{
			Name:  "worker",
			Usage: "launch machinery worker",
			Action: func(c *cli.Context) error {
				fmt.Println("运行工作者")
				if err := worker(); err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				return nil
			},
		},
	}

	app.Run(os.Args)
}

func StartServer() (*machinery.Server, error) {
	cnf := &config.Config{
		Broker:                  "amqp://admin:admin@localhost:5672/",
		Lock:                    "",
		MultipleBrokerSeparator: "",
		DefaultQueue:            "machinery_tasks",
		ResultBackend:           "amqp://admin:admin@localhost:5672/",
		ResultsExpireIn:         3600,
		AMQP: &config.AMQPConfig{
			Exchange:         "machinery_exchange",
			ExchangeType:     "direct",
			QueueDeclareArgs: nil,
			QueueBindingArgs: nil,
			BindingKey:       "machinery_task",
			PrefetchCount:    3,
			AutoDelete:       false,
		},
		SQS:           nil,
		Redis:         nil,
		GCPPubSub:     nil,
		MongoDB:       nil,
		TLSConfig:     nil,
		NoUnixSignals: false,
		DynamoDB:      nil,
	}

	server, err := machinery.NewServer(cnf)
	if err != nil {
		return nil, err
	}

	t := map[string]interface{}{
		"add":               exampletasks.Add,
		"multiply":          exampletasks.Multiply,
		"sum_ints":          exampletasks.SumInts,
		"sum_floats":        exampletasks.SumFloats,
		"concat":            exampletasks.Concat,
		"split":             exampletasks.Split,
		"panic_task":        exampletasks.PanicTask,
		"long_running_task": exampletasks.LongRunningTask,
	}

	return server, server.RegisterTasks(t)
}

func send() error {
	cleanup, err := tracers.SetupTracer("sender")
	if err != nil {
		log.FATAL.Fatalln("Unable to instantiate a tracer:", err)
	}
	defer cleanup()

	server, err := StartServer()
	if err != nil {
		return err
	}

	var (
		addTask0 tasks.Signature
	)

	var initTasks = func() {
		addTask0 = tasks.Signature{
			Name: "add",
			Args: []tasks.Arg{
				{
					Type:  "int64",
					Value: 1,
				},
				{
					Type:  "int64",
					Value: 10,
				},
			},
		}
	}

	initTasks()

	log.INFO.Println("Single task:")

	ctx := context.Background()
	asyncResult, err := server.SendTaskWithContext(ctx, &addTask0)
	if err != nil {
		return fmt.Errorf("Clould`t send task:%s", err.Error())
	}

	results, err := asyncResult.Get(time.Millisecond * 5)
	if err != nil {
		return fmt.Errorf("Getting task result failed with error:%s", err.Error())
	}
	log.INFO.Printf("Long running task returned = %v\n", tasks.HumanReadableResults(results))

	return nil
}

func worker() error {
	consumerTag := "machinery_worker"

	cleanup, err := tracers.SetupTracer(consumerTag)
	if err != nil {
		log.FATAL.Fatalln("Unable to instantiate a tracer:", err)
	}
	defer cleanup()

	server, err := StartServer()
	if err != nil {
		return err
	}

	worker := server.NewWorker(consumerTag, 0)

	errorHandler := func(err error) {
		log.ERROR.Println("I am an error handler:", err)
	}

	preTaskHandler := func(signature *tasks.Signature) {
		log.INFO.Println("I am a start of task handler for:", signature.Name)
	}

	postTaskHandler := func(signature *tasks.Signature) {
		log.INFO.Println("I am an end of task handler for:", signature.Name)
	}

	worker.SetPostTaskHandler(postTaskHandler)
	worker.SetErrorHandler(errorHandler)
	worker.SetPreTaskHandler(preTaskHandler)

	return worker.Launch()
}
