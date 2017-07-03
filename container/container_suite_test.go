package container_test

import (
	"io/ioutil"
	"testing"

	"github.com/Microsoft/hcsshim"
	"github.com/Sirupsen/logrus"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

var _ = BeforeSuite(func() {
	logrus.SetOutput(ioutil.Discard)
})

func TestContainer(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Container Suite")
}

type ContainerManager interface {
	Create(spec *specs.Spec) error
	Delete() error
	State() (*specs.State, error)
	Exec(*specs.Process) (hcsshim.Process, error)
}
