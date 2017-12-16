package rainbow

import (
	"encoding/json"
	"fmt"
	"github.com/46bit/discovery/rainbow/executor"
	"github.com/containerd/containerd"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type Server struct {
	exec     *executor.Executor
	operator *Operator
}

func NewServer(client *containerd.Client) *Server {
	exec := executor.NewExecutor(client)
	operator := NewOperator(exec.CmdChan, exec.EventChan)
	return &Server{
		exec:     exec,
		operator: operator,
	}
}

func (s *Server) Run(listenAddress string) {
	go s.exec.Run()
	go s.operator.Run()

	r := mux.NewRouter()
	r.HandleFunc("/deployments", s.getDeployments).Methods("GET")
	r.HandleFunc("/deployments", s.postDeployment).Methods("POST")
	r.HandleFunc("/deployments/{name}", s.getDeployment).Methods("GET")
	r.HandleFunc("/deployments/{name}", s.deleteDeployment).Methods("DELETE")
	http.Handle("/", r)
	http.ListenAndServe(listenAddress, nil)
}

func (s *Server) getDeployments(w http.ResponseWriter, r *http.Request) {
	deployments := []Deployment{}
	for _, deployment := range s.operator.Deployments {
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
	var deployment Deployment
	deployment.Jobs = []Job{}
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
	} else if _, ok := s.operator.Deployments[deployment.Name]; ok {
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
	s.operator.Add(deployment)
	http.Redirect(w, r, "/deployments/"+deployment.Name, http.StatusFound)
}

func (s *Server) getDeployment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	deployment, ok := s.operator.Deployments[name]
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
	if _, ok := s.operator.Deployments[name]; !ok {
		http.Error(w, "deployment was not found", http.StatusNotFound)
		return
	}
	s.operator.Remove(name)
	w.WriteHeader(http.StatusNoContent)
}
