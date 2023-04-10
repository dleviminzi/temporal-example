package main

import (
	"log"

	"worker/farewell"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "greetfarewell-tasks", worker.Options{})

	// parent
	w.RegisterWorkflow(farewell.GreetFarewell)

	// children
	w.RegisterWorkflow(farewell.GreetWorkflow)
	w.RegisterWorkflow(farewell.FarewellWorkflow)

	// activities
	w.RegisterActivity(farewell.GreetInSpanish)
	w.RegisterActivity(farewell.FarewellInSpanish)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
