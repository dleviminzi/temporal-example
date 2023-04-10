package farewell

type LeavingSignal struct {
	Message   string
	IsLeaving bool
}

// GreetInputSignal is the input required to execute the first child workflow: GreetWorkflow
type GreetInputSignal struct {
	Name      string `json:"name"`
	GreetDone bool   `json:"greet_done"`
}

// FarewellInputSignal is the input required to execute the second child workflow: FarewellWorkflow
type FarewellInputSignal struct {
	FarewellDone bool `json:"farewell_done"`
}

// WorkflowInputSignals can be used to agolmerate all of the received inputs so that future
// workflows can use inputs from earlier workflows
type WorkflowInputSignals struct {
	GreetInputSignal
	FarewellInputSignal
}
