const (
  // An unused description of a Container
  Described ContainerState = iota
  // A Container whose containerd.Container has been created
  Created
  // A Container whose containerd.Task has been created and started
  Started
  // A Container that has been SIGTERMed
  Sigtermed
  // A Container that has been SIGKILLed
  Sigkilled
  // A Container that has exited
  Stopped
  // A Container whose containerd.Task has been deleted
  Deleted
  // A Container whose containerd.Container has been deleted
  Erased
)

type ContainerState uint

// RESTARTING BUT WITH A BRAND-NEW CONTAINER AND TASK:
//   Created -> Started -> Stopped -> Deleted -> Erased -> Created -> ...
// STOPPING:
//   Created -> Started -> Sigtermed -> Stopped -> Deleted -> Erased
// FORCE-STOPPING AFTER TIMEOUT:
//   Created -> Started -> Sigtermed -> Sigkilled -> Stopped -> Deleted -> Erased

type Container struct {
  ID        string
  Remote    string
  Namespace string
  State     ContainerState
}
