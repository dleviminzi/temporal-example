package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"worker/farewell"

	"go.temporal.io/sdk/client"
)

type Server struct {
	tc    client.Client
	wfID  string
	runID string
}

func NewServer(tc client.Client) Server {
	return Server{
		tc: tc,
	}
}

func (s Server) greetInput(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	switch r.Method {
	case "POST":
		defer r.Body.Close()

		var input farewell.GreetInputSignal
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		// send user greet input for first child workflow
		err := s.tc.SignalWorkflow(context.Background(), s.wfID, s.runID, "greet-input-signal", input)
		if err != nil {
			log.Fatalln("Error sending the Signal", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\"name\": \"%s\"}", input.Name)
	}
}

func (s Server) farewellInput(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	switch r.Method {
	case "POST":
		defer r.Body.Close()

		var input farewell.FarewellInputSignal
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		// send user farewell input for second child workflow
		err := s.tc.SignalWorkflow(context.Background(), s.wfID, s.runID, "farewell-input-signal", input)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalln("Error sending the Signal", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\"farewell_done\": %t\n}", input.FarewellDone)
	}
}

func (s Server) result(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	switch r.Method {
	case "GET":
		// get result of the workflows
		var result string
		wf := s.tc.GetWorkflow(context.Background(), s.wfID, s.runID)
		err := wf.Get(context.Background(), &result)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalln("Unable get parent child workflow result", err)
		}
		log.Println("Parent child workflow result:", result)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "{\"result\": \"%s\"}", result)
	}
}

func (s Server) serve() {
	http.HandleFunc("/greet", s.greetInput)
	http.HandleFunc("/farewell", s.farewellInput)
	http.HandleFunc("/result", s.result)

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

	// start parent workflow execution
	// ordinarily would also be triggered by user input
	optionsPC := client.StartWorkflowOptions{
		ID:        "parent-child-greetfarewell-workflow",
		TaskQueue: "greetfarewell-tasks",
	}
	wPC, err := c.ExecuteWorkflow(context.Background(), optionsPC, farewell.GreetFarewell)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", wPC.GetID(), "RunID", wPC.GetRunID())

	s := NewServer(c)
	s.wfID = wPC.GetID()
	s.runID = wPC.GetRunID()
	s.serve()
}
