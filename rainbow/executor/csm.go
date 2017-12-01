package executor

func transition(api cdApi, container *container) error {
	if container.task == nil {
		task, err := newTask(api, container.container)
		if err != nil {
			return err
		}
		container.task = task
	} else {
		if container.task.state == created {
			if err := container.task.start(api); err != nil {
				return err
			}
			container.task.state = started
		} else if container.task.state == stopped {
			if err := container.task.delete(api); err != nil {
				return err
			}
			container.task = nil
		} else {
			return nil
		}
	}
	return nil
}
