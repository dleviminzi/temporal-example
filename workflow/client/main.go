package main

import (
	"context"
	"log"
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

	// start parent workflow execution
	var result string
	optionsPC := client.StartWorkflowOptions{
		ID:        "parent-child-greetfarewell-workflow",
		TaskQueue: "greetfarewell-tasks",
	}
	wPC, err := c.ExecuteWorkflow(context.Background(), optionsPC, farewell.GreetFarewell)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", wPC.GetID(), "RunID", wPC.GetRunID())

	// send user greet input for first child workflow
	greetInputSignal := farewell.GreetInputSignal{
		Name:      "daniel",
		GreetDone: true,
	}
	err = c.SignalWorkflow(context.Background(), wPC.GetID(), wPC.GetRunID(), "greet-input-signal", greetInputSignal)
	if err != nil {
		log.Fatalln("Error sending the Signal", err)
		return
	}

	// wait and send leaving signal
	time.Sleep(time.Second)
	leavingSignal := farewell.LeavingSignal{
		IsLeaving: true,
		Message:   "still here",
	}
	err = c.SignalWorkflow(context.Background(), wPC.GetID(), wPC.GetRunID(), "leaving-signal", leavingSignal)
	if err != nil {
		log.Fatalln("Error sending the Signal", err)
		return
	}

	// send user farewell input for second child workflow
	farewellInputSignal := farewell.FarewellInputSignal{
		FarewellDone: true,
	}
	err = c.SignalWorkflow(context.Background(), wPC.GetID(), wPC.GetRunID(), "farewell-input-signal", farewellInputSignal)
	if err != nil {
		log.Fatalln("Error sending the Signal", err)
		return
	}

	// get result of the workflows
	err = wPC.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable get parent child workflow result", err)
	}
	log.Println("Parent child workflow result:", result)
}
