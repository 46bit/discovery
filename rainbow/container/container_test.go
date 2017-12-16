package container_test

import (
	"github.com/46bit/discovery/rainbow/container"
	"github.com/containerd/containerd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Describe("Instance", func() {
	const (
		namespace = "default"
	)

	var (
		client     *containerd.Client
		container1 *container.Container
	)

	BeforeEach(func() {
		var err error
		client, err = containerd.New("/run/containerd/containerd.sock")
		Expect(err).ToNot(HaveOccurred())

		container1 = container.NewInstance(namespace, "container_test", "docker.io/46bit/hello-world:latest")
		Expect(container1.Status()).To(Equal(container.Described))
	})

	AfterEach(func() {
		if container1.Status() == container.Created {
			container1.Delete()
		}
		client.Close()
	})

	Context("Create", func() {
		It("succeeds", func() {
			err := container1.Create(client)
			Expect(err).ToNot(HaveOccurred())
			Expect(container1.Status()).To(Equal(container.Created))
		})

		It("fails if providing an invalid remote", func() {
			container1.Remote = "docker.io/46bit/does-not-exist:latest"
			err := container1.Create(client)
			Expect(err).To(HaveOccurred())
			Expect(errors.Cause(err).Error()).To(Equal("authorization failed"))
			Expect(container1.Status()).To(Equal(container.Described))
		})
	})

	Context("Delete", func() {
		BeforeEach(func() {
			err := container1.Create(client)
			Expect(err).ToNot(HaveOccurred())
			Expect(container1.Status()).To(Equal(container.Created))
		})

		It("succeeds", func() {
			err := container1.Delete()
			Expect(err).ToNot(HaveOccurred())
			Expect(container1.Status()).To(Equal(container.Deleted))
		})

		It("allows an container to be recreated", func() {
			container1.Delete()
			err := container1.Create(client)
			Expect(err).ToNot(HaveOccurred())
			Expect(container1.Status()).To(Equal(container.Created))
		})
	})
})
