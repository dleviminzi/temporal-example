package farewell

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

// ---------------------------------------------------------------------------
// example of parent-children workflows

// parent workflow
func GreetFarewell(ctx workflow.Context) (string, error) {
	childWorkflowOptions := workflow.ChildWorkflowOptions{}
	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)

	// parent blocks while waiting for user input needed to complete first workflow
	var greetInput GreetInputSignal
	greetInputChannel := workflow.GetSignalChannel(ctx, "greet-input-signal")
	for !greetInput.GreetDone {
		greetInputChannel.Receive(ctx, &greetInput)
	}

	var greetingResponse string
	err := workflow.ExecuteChildWorkflow(ctx, GreetWorkflow, greetInput).Get(ctx, &greetingResponse)
	if err != nil {
		return "", err
	}

	// create signal channel which will block until signal is received
	var leaveSignal LeavingSignal
	leaveInputChannel := workflow.GetSignalChannel(ctx, "leaving-signal")
	for !leaveSignal.IsLeaving {
		leaveInputChannel.Receive(ctx, &leaveSignal)
	}

	// parent blocks while waiting for user input needed to complete second workflow
	var farewellInput FarewellInputSignal
	farewellInputChannel := workflow.GetSignalChannel(ctx, "farewell-input-signal")
	for !farewellInput.FarewellDone {
		farewellInputChannel.Receive(ctx, &farewellInput)
	}

	var farewellResponse string
	err = workflow.ExecuteChildWorkflow(ctx, FarewellWorkflow, WorkflowInputSignals{greetInput, farewellInput}).Get(ctx, &farewellResponse)
	if err != nil {
		return "", err
	}

	sumResponse := "\n" + greetingResponse + "\n" + farewellResponse
	return sumResponse, nil
}

// children workflows
func GreetWorkflow(ctx workflow.Context, i GreetInputSignal) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	var spanishGreeting string
	err := workflow.ExecuteActivity(ctx, GreetInSpanish, i.Name).Get(ctx, &spanishGreeting)
	if err != nil {
		return "", err
	}

	return spanishGreeting, nil
}

func FarewellWorkflow(ctx workflow.Context, i WorkflowInputSignals) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	var spanishFarewell string
	err := workflow.ExecuteActivity(ctx, FarewellInSpanish, i.GreetInputSignal.Name).Get(ctx, &spanishFarewell)
	if err != nil {
		return "", err
	}

	return spanishFarewell, nil
}
