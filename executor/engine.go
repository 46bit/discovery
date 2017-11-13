package main

import (
	"context"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"log"
	"syscall"
)

type TaskExitCode struct {
	GroupName  string
	TaskGUID   string
	TaskRemote string
	ExitCode   containerd.ExitStatus
}

type Executor struct {
	Namespace     string
	Client        *containerd.Client
	Ctx           context.Context
	Groups        map[string]Group
	Tasks         map[string]containerd.Task
	NewGroups     <-chan Group
	DeleteGroups  <-chan string
	TaskExitCodes chan TaskExitCode
}

func NewExecutor(namespace string, client *containerd.Client, newGroups <-chan Group, deleteGroups <-chan string) *Executor {
	ctx := namespaces.WithNamespace(context.Background(), namespace)
	taskExitCodes := make(chan TaskExitCode)
	return &Executor{
		Namespace:     namespace,
		Client:        client,
		Ctx:           ctx,
		Groups:        map[string]Group{},
		Tasks:         map[string]containerd.Task{},
		NewGroups:     newGroups,
		DeleteGroups:  deleteGroups,
		TaskExitCodes: taskExitCodes,
	}
}

func (e *Executor) createGroup(group Group) {
	e.Groups[group.Name] = group
	for _, machine := range group.Machines {
		e.createTask(group.Name, machine)
	}
}

func (e *Executor) createTask(groupName string, machine Machine) {
	task, err := runTask(machine, e.Namespace, e.Client)
	if err != nil {
		log.Fatalln(err)
	}
	exitStatusC, err := task.Wait(e.Ctx)
	if err != nil {
		log.Fatalln(err)
	}
	if err := task.Start(e.Ctx); err != nil {
		log.Fatalln(err)
	}
	e.Tasks[machine.Remote] = task
	go func(taskExitCodes chan TaskExitCode, exitStatusC <-chan containerd.ExitStatus) {
		exitStatus := <-exitStatusC
		taskExitCodes <- TaskExitCode{
			GroupName:  groupName,
			TaskGUID:   machine.GUID,
			TaskRemote: machine.Remote,
			ExitCode:   exitStatus,
		}
	}(e.TaskExitCodes, exitStatusC)
}

func (e *Executor) deleteGroup(groupName string) {
	for remote := range e.Groups[groupName].Machines {
		err := e.Tasks[remote].Kill(e.Ctx, syscall.SIGTERM)
		if err != nil {
			log.Fatalln(err)
		}
	}
	delete(e.Groups, groupName)
}

func (e *Executor) run() {
	for {
		select {
		case newGroup := <-e.NewGroups:
			e.createGroup(newGroup)
		case groupName := <-e.DeleteGroups:
			e.deleteGroup(groupName)
		case taskExitCode := <-e.TaskExitCodes:
			fmt.Println(taskExitCode)

			container, err := e.Client.LoadContainer(e.Ctx, taskExitCode.TaskGUID)
			if err != nil {
				log.Fatalln(err)
			}

			err = container.Delete(e.Ctx, containerd.WithSnapshotCleanup)
			if err != nil {
				log.Fatalln(err)
			}

			if _, ok := e.Groups[taskExitCode.GroupName]; ok {
				e.createTask(taskExitCode.GroupName, e.Groups[taskExitCode.GroupName].Machines[taskExitCode.TaskRemote])
			}
		}
	}
}
