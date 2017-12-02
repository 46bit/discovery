package executor

func transition(api cdApi, c *container) error {
	var err error
	switch c.task.state {
	case unspecified, deleted:
		err = c.task.create(api, c.container)
	case created:
		err = c.task.start(api)
	case stopped:
		err = c.task.delete(api)
	case started:
	}
	return err
}
