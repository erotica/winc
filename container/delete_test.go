package container_test

import (
	"errors"

	"code.cloudfoundry.org/winc/container"
	"code.cloudfoundry.org/winc/hcsclient/hcsclientfakes"
	"code.cloudfoundry.org/winc/network/networkfakes"
	"code.cloudfoundry.org/winc/sandbox/sandboxfakes"
	"github.com/Microsoft/hcsshim"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Delete", func() {
	const (
		expectedContainerId        = "containerid"
		expectedContainerBundleDir = "C:\\bundle"
	)

	var (
		hcsClient        *hcsclientfakes.FakeClient
		sandboxManager   *sandboxfakes.FakeSandboxManager
		fakeContainer    *hcsclientfakes.FakeContainer
		networkManager   *networkfakes.FakeNetworkManager
		containerManager container.ContainerManager
	)

	BeforeEach(func() {
		hcsClient = &hcsclientfakes.FakeClient{}
		sandboxManager = &sandboxfakes.FakeSandboxManager{}
		fakeContainer = &hcsclientfakes.FakeContainer{}
		networkManager = &networkfakes.FakeNetworkManager{}
		containerManager = container.NewManager(hcsClient, sandboxManager, networkManager, expectedContainerId)
	})

	Context("when the specified container is not running", func() {
		var pid int
		BeforeEach(func() {
			pid = 42
			fakeContainer.ProcessListReturns([]hcsshim.ProcessListItem{
				{ProcessId: uint32(pid), ImageName: "wininit.exe"},
			}, nil)
			hcsClient.OpenContainerReturns(fakeContainer, nil)
		})

		It("deletes it", func() {
			Expect(containerManager.Delete()).To(Succeed())

			Expect(sandboxManager.UnmountCallCount()).To(Equal(1))
			Expect(sandboxManager.UnmountArgsForCall(0)).To(Equal(pid))

			Expect(hcsClient.OpenContainerCallCount()).To(Equal(2))
			Expect(hcsClient.OpenContainerArgsForCall(0)).To(Equal(expectedContainerId))

			Expect(networkManager.DeleteContainerEndpointsCallCount()).To(Equal(1))
			container, containerID := networkManager.DeleteContainerEndpointsArgsForCall(0)
			Expect(container).To(Equal(fakeContainer))
			Expect(containerID).To(Equal(expectedContainerId))

			Expect(fakeContainer.ShutdownCallCount()).To(Equal(1))

			Expect(sandboxManager.DeleteCallCount()).To(Equal(1))
		})

		Context("when unmounting the sandbox fails", func() {
			BeforeEach(func() {
				sandboxManager.UnmountReturns(errors.New("unmounting failed"))
			})

			It("continues deleting the container and returns an error", func() {
				Expect(containerManager.Delete()).NotTo(Succeed())

				Expect(hcsClient.OpenContainerCallCount()).To(Equal(2))
				Expect(hcsClient.OpenContainerArgsForCall(0)).To(Equal(expectedContainerId))

				Expect(networkManager.DeleteContainerEndpointsCallCount()).To(Equal(1))
				container, containerID := networkManager.DeleteContainerEndpointsArgsForCall(0)
				Expect(container).To(Equal(fakeContainer))
				Expect(containerID).To(Equal(expectedContainerId))

				Expect(fakeContainer.ShutdownCallCount()).To(Equal(1))

				Expect(sandboxManager.DeleteCallCount()).To(Equal(1))
			})
		})

		Context("when shutting down the container does not immediately succeed", func() {
			var shutdownContainerError = errors.New("shutdown container failed")

			BeforeEach(func() {
				hcsClient.OpenContainerReturns(fakeContainer, nil)
				fakeContainer.ShutdownReturns(shutdownContainerError)
				hcsClient.IsPendingReturns(false)
			})

			It("calls terminate", func() {
				Expect(containerManager.Delete()).To(Succeed())
				Expect(fakeContainer.TerminateCallCount()).To(Equal(1))
				Expect(sandboxManager.DeleteCallCount()).To(Equal(1))
			})

			Context("when shutdown is pending", func() {
				BeforeEach(func() {
					hcsClient.IsPendingReturnsOnCall(0, true)
				})

				It("waits for shutdown to finish", func() {
					Expect(containerManager.Delete()).To(Succeed())
					Expect(fakeContainer.TerminateCallCount()).To(Equal(0))
					Expect(sandboxManager.DeleteCallCount()).To(Equal(1))
				})

				Context("when shutdown does not finish before the timeout", func() {
					var shutdownWaitError = errors.New("waiting for shutdown failed")

					BeforeEach(func() {
						fakeContainer.WaitTimeoutReturnsOnCall(0, shutdownWaitError)
					})

					It("it calls terminate", func() {
						Expect(containerManager.Delete()).To(Succeed())
						Expect(fakeContainer.TerminateCallCount()).To(Equal(1))
						Expect(sandboxManager.DeleteCallCount()).To(Equal(1))
					})

					Context("when terminate does not immediately succeed", func() {
						var terminateContainerError = errors.New("terminate container failed")

						BeforeEach(func() {
							fakeContainer.TerminateReturns(terminateContainerError)
						})

						It("errors", func() {
							Expect(containerManager.Delete()).To(Equal(terminateContainerError))
						})

						Context("when terminate is pending", func() {
							BeforeEach(func() {
								hcsClient.IsPendingReturnsOnCall(1, true)
							})

							It("waits for terminate to finish", func() {
								Expect(containerManager.Delete()).To(Succeed())
								Expect(sandboxManager.DeleteCallCount()).To(Equal(1))
							})

							Context("when terminate does not finish before the timeout", func() {
								var terminateWaitError = errors.New("waiting for terminate failed")

								BeforeEach(func() {
									fakeContainer.WaitTimeoutReturnsOnCall(1, terminateWaitError)
								})

								It("errors", func() {
									Expect(containerManager.Delete()).To(Equal(terminateWaitError))
									Expect(sandboxManager.DeleteCallCount()).To(Equal(0))
								})
							})
						})
					})
				})
			})
		})

		Context("when the sandbox delete fails", func() {
			var deleteSandboxError = errors.New("delete sandbox failed")

			BeforeEach(func() {
				sandboxManager.DeleteReturns(deleteSandboxError)
			})

			It("errors", func() {
				Expect(containerManager.Delete()).To(Equal(deleteSandboxError))
			})
		})
	})

	Context("when the container does not exist", func() {
		var openContainerError = errors.New("open container failed")

		BeforeEach(func() {
			hcsClient.OpenContainerReturns(nil, openContainerError)
		})

		It("errors", func() {
			Expect(containerManager.Delete()).To(Equal(openContainerError))
		})
	})
})
