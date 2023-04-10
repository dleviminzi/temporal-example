package main

import (
	"context"
	"log"
	"os"
	"time"

	"worker/farewell"

	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "greeting-workflow",
		TaskQueue: "greeting-tasks",
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, farewell.GreetSomeone, os.Args[1])
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	// wait a bit and then send leaving signal
	i := 0
	s := farewell.LeavingSignal{
		IsLeaving: false,
		Message:   "still here",
	}

	for i < 5 {

		if i == 4 {
			s.IsLeaving = true
			s.Message = "headed out"
		}

		err = c.SignalWorkflow(context.Background(), we.GetID(), we.GetRunID(), "leaving-signal", s)
		if err != nil {
			log.Fatalln("Error sending the Signal", err)
			return
		}

		time.Sleep(time.Second)
		i++
	}

	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}
	log.Println("Workflow result:", result)
}
