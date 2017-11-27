package containers

import (
	cd "github.com/containerd/containerd"
	"syscall"
	"time"
)

type task struct {
	state state
	task  *cd.Task
}

func newTask(api cdApi, cdContainer *cd.Container) (*task, error) {
	t := task{}
	if err := t.create(api, cdContainer); err != nil {
		return nil, err
	}
	return &t, nil
}

func (t *task) create(api cdApi, cdContainer *cd.Container) error {
	cdTask, err := (*cdContainer).NewTask(api.context, cd.Stdio)
	if err != nil {
		return err
	}
	t.task = &cdTask
	t.state = created
	return nil
}

func (t *task) start(api cdApi) error {
	err := (*t.task).Start(api.context)
	if err != nil {
		return err
	}
	t.state = started
	return nil
}

func (t *task) stop(api cdApi) error {
	exitStatusC, err := (*t.task).Wait(api.context)
	if err != nil {
		return err
	}
	err = (*t.task).Kill(api.context, syscall.SIGTERM, cd.WithKillAll)
	if err != nil {
		return err
	}
	select {
	case <-exitStatusC:
	case <-time.After(5 * time.Second):
		err = (*t.task).Kill(api.context, syscall.SIGKILL, cd.WithKillAll)
		if err != nil {
			return err
		}
	}
	<-exitStatusC
	t.state = stopped
	return nil
}

func (t *task) delete(api cdApi) error {
	_, err := (*t.task).Delete(api.context)
	if err != nil {
		return err
	}
	t.state = deleted
	return nil
}
