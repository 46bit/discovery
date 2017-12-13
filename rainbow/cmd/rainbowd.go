package main

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
	"math/rand"
	"net/http"
	"time"
)

var op *operator.Operator

func main() {
	rand.Seed(time.Now().UnixNano())

	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	exec := executor.NewExecutor(client)
	go exec.Run()

	op = operator.NewOperator(exec.CmdChan, exec.EventChan)
	go op.Run()

	r := mux.NewRouter()
	r.HandleFunc("/deployments", getDeployments).Methods("GET")
	r.HandleFunc("/deployments", postDeployment).Methods("POST")
	r.HandleFunc("/deployments/{name}", getDeployment).Methods("GET")
	r.HandleFunc("/deployments/{name}", deleteDeployment).Methods("DELETE")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func getDeployments(w http.ResponseWriter, r *http.Request) {
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

func postDeployment(w http.ResponseWriter, r *http.Request) {
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

func getDeployment(w http.ResponseWriter, r *http.Request) {
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

func deleteDeployment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	if _, ok := op.Deployments[name]; !ok {
		http.Error(w, "deployment was not found", http.StatusNotFound)
		return
	}
	op.Remove(name)
	w.WriteHeader(http.StatusNoContent)
}
