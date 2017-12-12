package instance_test

import (
	"github.com/46bit/discovery/rainbow/instance"
	"github.com/containerd/containerd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Instance", func() {
	const (
		namespace = "default"
	)

	var (
		client    *containerd.Client
		instance1 *instance.Instance
	)

	BeforeEach(func() {
		var err error
		client, err = containerd.New("/run/containerd/containerd.sock")
		Expect(err).ToNot(HaveOccurred())

		instance1 = instance.NewInstance(namespace, "instance_test", "docker.io/46bit/hello-world:latest")
		Expect(instance1.Status()).To(Equal(instance.Described))
	})

	AfterEach(func() {
		if instance1.Status() == instance.Created {
			instance1.Delete()
		}
		client.Close()
	})

	Context("Create", func() {
		It("succeeds", func() {
			err := instance1.Create(client)
			Expect(err).ToNot(HaveOccurred())
			Expect(instance1.Status()).To(Equal(instance.Created))
		})

		It("fails if providing an invalid remote", func() {
			i := instance.NewInstance(namespace, "instance_test", "docker.io/46bit/does-not-exist:latest")
			Expect(i.Status()).To(Equal(instance.Described))
			err := i.Create(client)
			Expect(err).To(HaveOccurred())
			Expect(i.Status()).To(Equal(instance.Described))
		})
	})
})
