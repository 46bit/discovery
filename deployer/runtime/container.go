const (
  // An unused description of a Container
  Described ContainerState = iota
  // A Container whose containerd.Container has been created
  Created
  // A Container whose containerd.Task has been created
  Prepared
  // A Container whose containerd.Task has been started
  Started
  // A Container that has been SIGTERMed
  Sigtermed
  // A Container that has been SIGKILLed
  Sigkilled
  // A Container that has exited
  Exited
  // A Container whose containerd.Task has been deleted
  Deleted
  // A Container whose containerd.Container has been deleted
  Erased
)

type ContainerState uint

// RESTARTING BUT WITH A BRAND-NEW CONTAINER AND TASK:
//   Created -> Prepared -> Started -> Stopped -> Deleted -> Erased -> Created -> ...
// STOPPING:
//   Created -> Prepared -> Started -> Sigtermed -> Stopped -> Deleted -> Erased
// FORCE-STOPPING AFTER TIMEOUT:
//   Created -> Prepared -> Started -> Sigtermed -> Sigkilled -> Stopped -> Deleted -> Erased

type Container struct {
  ID        string
  Remote    string
  Namespace string
  State     ContainerState
}

func (c *Container) Create(client *containerd.Client, ctx *context.Context) (*containerd.Container, error) {
  containerdContainer, err := createContainer(client, ctx, c.ID, c.Remote)
  if err != nil {
    return nil, err
  }
  c.State = Created
  return containerdContainer, nil
}

func (c *Container) Prepare(ctx *context.Context, containerdContainer containerd.Container) (*containerd.Task, exitStatusC <-chan containerd.ExitStatus, error) {
  containerdTask, err := createTask(ctx, containerdContainer)
  if err != nil {
    return nil, err
  }
  c.State = Prepared
  return containerdTask, nil
}

func (c *Container) Start(ctx *context.Context, containerdTask containerd.Task) (exitStatusC <-chan containerd.ExitStatus, error) {
  exitStatusC, err := startTask(ctx, containerdTask)
  if err != nil {
    return nil, nil, err
  }
  c.State = Prepared
  return exitStatusC, nil
}

func (c *Container) Sigterm(ctx *context.Context, containerdTask containerd.Task) error {
  err := c.signal(ctx, containerdTask, syscall.SIGTERM)
  if err != nil {
    return err
  }
  c.State = Sigtermed
  return nil
}

func (c *Container) Sigkill(ctx *context.Context, containerdTask containerd.Task) error {
  err := c.signal(ctx, containerdTask, syscall.SIGKILL)
  if err != nil {
    return err
  }
  c.State = Sigkilled
  return nil
}

func (c *Container) signal(ctx *context.Context, containerdTask containerd.Task, signal syscall.Signal) error {
  status, err := containerdTask.Status(ctx)
  if err != nil {
    return err
  }
  switch status.Status {
  case containerd.Running:
    fallthrough
  case containerd.Paused:
    fallthrough
  case containerd.Pausing:
    err = containerdTask.Kill(ctx, signal, containerd.WithKillAll)
    if err != nil {
      return err
    }
  }
  return nil
}

func (c *Container) Delete(ctx *context.Context, containerdTask containerd.Task) error {
  _, err := containerdTask.Delete(ctx)
  if err != nil {
    return fmt.Errorf("Error deleting task %s: %s", id, err)
  }
  c.State = Deleted
  return nil
}

func (c *Container) Erase(client *containerd.Client, ctx *context.Context, containerdContainer containerd.Container) error {
  err = containerdContainer.Delete(ctx, containerd.WithSnapshotCleanup)
  if err != nil {
    return err
  }
  c.State = Erased
  return nil
}
