package server

import (
	"encoding/json"
	"fmt"
	"github.com/46bit/discovery/rainbow"
	"github.com/46bit/discovery/rainbow/executor"
	"github.com/46bit/discovery/rainbow/operator"
	"github.com/containerd/containerd"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type Server struct {
	exec *executor.Executor
	op   *operator.Operator
}

func NewServer(client *containerd.Client) *Server {
	exec := executor.NewExecutor(client)
	op = operator.NewOperator(exec.CmdChan, exec.EventChan)
	return &Server{
		exec: exec,
		op:   op,
	}
}

func (s *Server) Run(listenAddress string) {
	go s.exec.Run()
	go s.op.Run()

	r := mux.NewRouter()
	r.HandleFunc("/deployments", s.getDeployments).Methods("GET")
	r.HandleFunc("/deployments", s.postDeployment).Methods("POST")
	r.HandleFunc("/deployments/{name}", s.getDeployment).Methods("GET")
	r.HandleFunc("/deployments/{name}", s.deleteDeployment).Methods("DELETE")
	http.Handle("/", r)
	http.ListenAndServe(listenAddress, nil)
}

func (s *Server) getDeployments(w http.ResponseWriter, r *http.Request) {
	deployments := []rainbow.Deployment{}
	for _, deployment := range op.Deployments {
		deployments = append(deployments, deployment)
	}
	responseBody, err := json.Marshal(deployments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(responseBody)
}

func (s *Server) postDeployment(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var deployment rainbow.Deployment
	deployment.Jobs = []rainbow.Job{}
	err = json.Unmarshal(body, &deployment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if deployment.Name == "" {
		http.Error(w, "deployment name was empty", http.StatusBadRequest)
		return
	} else if len(deployment.Name) < 2 {
		http.Error(w, fmt.Sprintf("deployment name '%s' was too short", deployment.Name), http.StatusBadRequest)
		return
	} else if _, ok := op.Deployments[deployment.Name]; ok {
		http.Error(w, "deployment name already in use", http.StatusBadRequest)
		return
	}
	for _, job := range deployment.Jobs {
		if job.Name == "" {
			http.Error(w, "job name was empty", http.StatusBadRequest)
			return
		} else if len(job.Name) < 2 {
			http.Error(w, fmt.Sprintf("job name '%s' was too short", job.Name), http.StatusBadRequest)
			return
		}
	}
	op.Add(deployment)
	http.Redirect(w, r, "/deployments/"+deployment.Name, http.StatusFound)
}

func (s *Server) getDeployment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	deployment, ok := op.Deployments[name]
	if !ok {
		http.Error(w, "deployment was not found", http.StatusNotFound)
		return
	}
	responseBody, err := json.Marshal(deployment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write(responseBody)
}

func (s *Server) deleteDeployment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	if _, ok := op.Deployments[name]; !ok {
		http.Error(w, "deployment was not found", http.StatusNotFound)
		return
	}
	op.Remove(name)
	w.WriteHeader(http.StatusNoContent)
}
