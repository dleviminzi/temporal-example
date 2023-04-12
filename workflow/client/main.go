package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"worker/farewell"

	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

// this is pretty ugly stuff, apologies for the complete lack of code organization

var wf = "parent-child-greetfarewell-workflow-%s"

// Server is bunch of endpoints that translate requests and commands into
// requests and commands that the temporal client can execute
type Server struct {
	tc client.Client
}

func NewServer(tc client.Client) Server {
	return Server{
		tc: tc,
	}
}

// greetInput receives the input for the first child workflow from the frontend.
// It will take that input, decode it into the input type that the workflow expects
// and it will send that as a signal. That signal will stop the parent from blocking
// and allow it to execute the first child workflow.
func (s Server) greetInput(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	id := r.URL.Query().Get("id")

	defer r.Body.Close()

	// decode body into input signal
	var input farewell.GreetInputSignal
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// send input signal
	err := s.tc.SignalWorkflow(context.Background(), fmt.Sprintf(wf, id), "", "greet-input-signal", input)
	if err != nil {
		log.Fatalln("Error sending the Signal", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(input)
}

// farewellInput receives the input for the secon child workflow from the frontend.
// It will take that, decode it into the second workflows input type.
func (s Server) farewellInput(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	id := r.URL.Query().Get("id")

	defer r.Body.Close()

	var input farewell.FarewellInputSignal
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// send user farewell input for second child workflow
	err := s.tc.SignalWorkflow(context.Background(), fmt.Sprintf(wf, id), "", "farewell-input-signal", input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln("Error sending the Signal", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(input)
}

// result is used to get the latest workflow result for a user id
// If you want historical results, you can use the queryWorkflows endpoint
func (s Server) result(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	id := r.URL.Query().Get("id")

	// get result of the workflows
	var result string
	wf := s.tc.GetWorkflow(context.Background(), fmt.Sprintf(wf, id), "")
	err := wf.Get(context.Background(), &result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatalln("Unable get parent child workflow result", err)
	}
	log.Println("Parent child workflow result:", result)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// queryWorkflows will request a list of all workflows in all states for a given id
func (s Server) queryWorkflows(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	id := r.URL.Query().Get("id")

	list, err := s.tc.ListWorkflow(r.Context(), &workflowservice.ListWorkflowExecutionsRequest{
		Namespace: "default",
		Query:     fmt.Sprintf("(WorkflowId = 'parent-child-greetfarewell-workflow-%s')", id),
	})
	if err != nil {
		log.Println("Unable get parent child workflow result", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	if list != nil {
		fmt.Println(list.Executions)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(list.Executions)
	}
}

// queryActiveWorkflow will return the current state of the active workflow run
// for an id
func (s Server) queryActiveWorkflow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	wfID := fmt.Sprintf(wf, r.URL.Query().Get("id"))

	// use the query hook in the workflow definition
	response, err := s.tc.QueryWorkflow(r.Context(), wfID, "", farewell.QueryTypeWorkflowStatus)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !response.HasValue() {
		log.Printf("no active state found for %s\n", wfID)
		return
	}

	var wfState string
	if err := response.Get(&wfState); err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(wfState)

	w.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprint(w, wfState); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// startWorkflow will start a workflow run for the given id or resume the active
// workflow run for the id
func (s Server) startWorkflow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// garbage im sorry
	type Id struct {
		Id int `json:"id"`
	}
	var id Id
	if err := json.NewDecoder(r.Body).Decode(&id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	// important! notice that we specify the workflow id here!
	optionsPC := client.StartWorkflowOptions{
		ID:        fmt.Sprintf(wf, strconv.Itoa(id.Id)),
		TaskQueue: "greetfarewell-tasks",
	}
	wPC, err := s.tc.ExecuteWorkflow(context.Background(), optionsPC, farewell.GreetFarewell)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", wPC.GetID(), "RunID", wPC.GetRunID())
}

func (s Server) serve() {
	http.HandleFunc("/greet", s.greetInput)
	http.HandleFunc("/farewell", s.farewellInput)
	http.HandleFunc("/result", s.result)
	http.HandleFunc("/wfs", s.queryWorkflows)
	http.HandleFunc("/startwf", s.startWorkflow)
	http.HandleFunc("/activewf", s.queryActiveWorkflow)

	if err := http.ListenAndServe(":8090", nil); err != nil {
		log.Fatal(err)
	}
}

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	s := NewServer(c)
	s.serve()
}
