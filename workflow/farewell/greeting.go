package farewell

import (
	"encoding/json"
	"time"

	"go.temporal.io/sdk/workflow"
)

// WorkflowState holds the current state of the workflow.
// This includes:
// Stage =  where we are in the workflow
// ReceivedInput = what we received in signals from the front end
// ChildOneResult = the result of the first child workflows execution
// ChildTwoResult = the result of the second child workflows execution
type WorkflowState struct {
	Stage          string `json:"stage"`
	ReceivedInput  WorkflowInputSignals
	ChildOneResult string `json:"child_one_result"`
	ChildTwoResult string `json:"child_two_result"`
}

// this is the query type that will be used in the client
const QueryTypeWorkflowStatus = "query-type-workflow-status"

// ---------------------------------------------------------------------------
// example of parent-children workflows

// parent workflow
func GreetFarewell(ctx workflow.Context) (string, error) {
	wfState := WorkflowState{Stage: "started"}

	// establish a query handler for this type of workflow that will return the
	// the current WorkflowState marshalled into json
	if err := workflow.SetQueryHandler(ctx, QueryTypeWorkflowStatus, func() (string, error) {
		currState, err := json.Marshal(wfState)
		if err != nil {
			return "", err
		}
		return string(currState), nil
	}); err != nil {
		return "", err
	}

	ctx = workflow.WithChildOptions(ctx, workflow.ChildWorkflowOptions{})

	// block while waiting for user input needed to complete the first child workflow
	greetInputChannel := workflow.GetSignalChannel(ctx, "greet-input-signal")
	for !wfState.ReceivedInput.GreetDone {
		greetInputChannel.Receive(ctx, &wfState.ReceivedInput.GreetInputSignal)
	}

	// execute the first child workflow now that we have received user input
	err := workflow.ExecuteChildWorkflow(ctx, GreetWorkflow, wfState.ReceivedInput.GreetInputSignal).Get(ctx, &wfState.ChildOneResult)
	if err != nil {
		return "", err
	}

	// update the WorkflowState so that it is accessible via our query
	wfState.Stage = "post-child-one"

	// block while waiting for user input needed to complete the second child workflow
	farewellInputChannel := workflow.GetSignalChannel(ctx, "farewell-input-signal")
	for !wfState.ReceivedInput.FarewellDone {
		farewellInputChannel.Receive(ctx, &wfState.ReceivedInput.FarewellInputSignal)
	}

	// execute the second child workflow now that we have received user input
	err = workflow.ExecuteChildWorkflow(ctx, FarewellWorkflow, wfState.ReceivedInput).Get(ctx, &wfState.ChildTwoResult)
	if err != nil {
		return "", err
	}

	// update the WorkflowState so that it is accessible via our query
	wfState.Stage = "post-child-two--finished"

	// return the result of the workflow (in this case the concatenation of the child workflow results)
	sumResponse := wfState.ChildOneResult + " " + wfState.ChildTwoResult
	return sumResponse, nil
}

// GreetWorkflow is the first of two child workflows for GreetFarewell workflow.
// This workflow will execute an activity to get a spanish greeting and then it
// will immediately request its result.
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

// FarewellWorkflow is the second and final workflow for the GreetFarewell workflow.
// It will execute an activity to get a spanish farewell message and then it will
// immediately request its result.
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
