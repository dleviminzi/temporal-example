package farewell

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

// this is a workflow with two activities that blocks between them on a signal
func GreetSomeone(ctx workflow.Context, name string) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	var spanishGreeting string
	err := workflow.ExecuteActivity(ctx, GreetInSpanish, name).Get(ctx, &spanishGreeting)
	if err != nil {
		return "", err
	}

	// create signal channel which will block until signal is received
	var s LeavingSignal
	c := workflow.GetSignalChannel(ctx, "leaving-signal")
	for !s.IsLeaving {
		c.Receive(ctx, &s)
	}

	var spanishFarewell string
	err = workflow.ExecuteActivity(ctx, FarewellInSpanish, name).Get(ctx, &spanishFarewell)
	if err != nil {
		return "", err
	}

	helloGoodbye := "\n" + spanishGreeting + "\n" + spanishFarewell

	return helloGoodbye, nil
}
