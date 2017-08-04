// Code generated by counterfeiter. DO NOT EDIT.
package containerfakes

import (
	"sync"

	"code.cloudfoundry.org/winc/container"
	"code.cloudfoundry.org/winc/hcscontainer"
	"github.com/Microsoft/hcsshim"
)

type FakeHCSClient struct {
	GetContainerPropertiesStub        func(string) (hcsshim.ContainerProperties, error)
	getContainerPropertiesMutex       sync.RWMutex
	getContainerPropertiesArgsForCall []struct {
		arg1 string
	}
	getContainerPropertiesReturns struct {
		result1 hcsshim.ContainerProperties
		result2 error
	}
	getContainerPropertiesReturnsOnCall map[int]struct {
		result1 hcsshim.ContainerProperties
		result2 error
	}
	NameToGuidStub        func(string) (hcsshim.GUID, error)
	nameToGuidMutex       sync.RWMutex
	nameToGuidArgsForCall []struct {
		arg1 string
	}
	nameToGuidReturns struct {
		result1 hcsshim.GUID
		result2 error
	}
	nameToGuidReturnsOnCall map[int]struct {
		result1 hcsshim.GUID
		result2 error
	}
	CreateContainerStub        func(string, *hcsshim.ContainerConfig) (hcscontainer.Container, error)
	createContainerMutex       sync.RWMutex
	createContainerArgsForCall []struct {
		arg1 string
		arg2 *hcsshim.ContainerConfig
	}
	createContainerReturns struct {
		result1 hcscontainer.Container
		result2 error
	}
	createContainerReturnsOnCall map[int]struct {
		result1 hcscontainer.Container
		result2 error
	}
	OpenContainerStub        func(string) (hcscontainer.Container, error)
	openContainerMutex       sync.RWMutex
	openContainerArgsForCall []struct {
		arg1 string
	}
	openContainerReturns struct {
		result1 hcscontainer.Container
		result2 error
	}
	openContainerReturnsOnCall map[int]struct {
		result1 hcscontainer.Container
		result2 error
	}
	IsPendingStub        func(error) bool
	isPendingMutex       sync.RWMutex
	isPendingArgsForCall []struct {
		arg1 error
	}
	isPendingReturns struct {
		result1 bool
	}
	isPendingReturnsOnCall map[int]struct {
		result1 bool
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeHCSClient) GetContainerProperties(arg1 string) (hcsshim.ContainerProperties, error) {
	fake.getContainerPropertiesMutex.Lock()
	ret, specificReturn := fake.getContainerPropertiesReturnsOnCall[len(fake.getContainerPropertiesArgsForCall)]
	fake.getContainerPropertiesArgsForCall = append(fake.getContainerPropertiesArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GetContainerProperties", []interface{}{arg1})
	fake.getContainerPropertiesMutex.Unlock()
	if fake.GetContainerPropertiesStub != nil {
		return fake.GetContainerPropertiesStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getContainerPropertiesReturns.result1, fake.getContainerPropertiesReturns.result2
}

func (fake *FakeHCSClient) GetContainerPropertiesCallCount() int {
	fake.getContainerPropertiesMutex.RLock()
	defer fake.getContainerPropertiesMutex.RUnlock()
	return len(fake.getContainerPropertiesArgsForCall)
}

func (fake *FakeHCSClient) GetContainerPropertiesArgsForCall(i int) string {
	fake.getContainerPropertiesMutex.RLock()
	defer fake.getContainerPropertiesMutex.RUnlock()
	return fake.getContainerPropertiesArgsForCall[i].arg1
}

func (fake *FakeHCSClient) GetContainerPropertiesReturns(result1 hcsshim.ContainerProperties, result2 error) {
	fake.GetContainerPropertiesStub = nil
	fake.getContainerPropertiesReturns = struct {
		result1 hcsshim.ContainerProperties
		result2 error
	}{result1, result2}
}

func (fake *FakeHCSClient) GetContainerPropertiesReturnsOnCall(i int, result1 hcsshim.ContainerProperties, result2 error) {
	fake.GetContainerPropertiesStub = nil
	if fake.getContainerPropertiesReturnsOnCall == nil {
		fake.getContainerPropertiesReturnsOnCall = make(map[int]struct {
			result1 hcsshim.ContainerProperties
			result2 error
		})
	}
	fake.getContainerPropertiesReturnsOnCall[i] = struct {
		result1 hcsshim.ContainerProperties
		result2 error
	}{result1, result2}
}

func (fake *FakeHCSClient) NameToGuid(arg1 string) (hcsshim.GUID, error) {
	fake.nameToGuidMutex.Lock()
	ret, specificReturn := fake.nameToGuidReturnsOnCall[len(fake.nameToGuidArgsForCall)]
	fake.nameToGuidArgsForCall = append(fake.nameToGuidArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("NameToGuid", []interface{}{arg1})
	fake.nameToGuidMutex.Unlock()
	if fake.NameToGuidStub != nil {
		return fake.NameToGuidStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.nameToGuidReturns.result1, fake.nameToGuidReturns.result2
}

func (fake *FakeHCSClient) NameToGuidCallCount() int {
	fake.nameToGuidMutex.RLock()
	defer fake.nameToGuidMutex.RUnlock()
	return len(fake.nameToGuidArgsForCall)
}

func (fake *FakeHCSClient) NameToGuidArgsForCall(i int) string {
	fake.nameToGuidMutex.RLock()
	defer fake.nameToGuidMutex.RUnlock()
	return fake.nameToGuidArgsForCall[i].arg1
}

func (fake *FakeHCSClient) NameToGuidReturns(result1 hcsshim.GUID, result2 error) {
	fake.NameToGuidStub = nil
	fake.nameToGuidReturns = struct {
		result1 hcsshim.GUID
		result2 error
	}{result1, result2}
}

func (fake *FakeHCSClient) NameToGuidReturnsOnCall(i int, result1 hcsshim.GUID, result2 error) {
	fake.NameToGuidStub = nil
	if fake.nameToGuidReturnsOnCall == nil {
		fake.nameToGuidReturnsOnCall = make(map[int]struct {
			result1 hcsshim.GUID
			result2 error
		})
	}
	fake.nameToGuidReturnsOnCall[i] = struct {
		result1 hcsshim.GUID
		result2 error
	}{result1, result2}
}

func (fake *FakeHCSClient) CreateContainer(arg1 string, arg2 *hcsshim.ContainerConfig) (hcscontainer.Container, error) {
	fake.createContainerMutex.Lock()
	ret, specificReturn := fake.createContainerReturnsOnCall[len(fake.createContainerArgsForCall)]
	fake.createContainerArgsForCall = append(fake.createContainerArgsForCall, struct {
		arg1 string
		arg2 *hcsshim.ContainerConfig
	}{arg1, arg2})
	fake.recordInvocation("CreateContainer", []interface{}{arg1, arg2})
	fake.createContainerMutex.Unlock()
	if fake.CreateContainerStub != nil {
		return fake.CreateContainerStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.createContainerReturns.result1, fake.createContainerReturns.result2
}

func (fake *FakeHCSClient) CreateContainerCallCount() int {
	fake.createContainerMutex.RLock()
	defer fake.createContainerMutex.RUnlock()
	return len(fake.createContainerArgsForCall)
}

func (fake *FakeHCSClient) CreateContainerArgsForCall(i int) (string, *hcsshim.ContainerConfig) {
	fake.createContainerMutex.RLock()
	defer fake.createContainerMutex.RUnlock()
	return fake.createContainerArgsForCall[i].arg1, fake.createContainerArgsForCall[i].arg2
}

func (fake *FakeHCSClient) CreateContainerReturns(result1 hcscontainer.Container, result2 error) {
	fake.CreateContainerStub = nil
	fake.createContainerReturns = struct {
		result1 hcscontainer.Container
		result2 error
	}{result1, result2}
}

func (fake *FakeHCSClient) CreateContainerReturnsOnCall(i int, result1 hcscontainer.Container, result2 error) {
	fake.CreateContainerStub = nil
	if fake.createContainerReturnsOnCall == nil {
		fake.createContainerReturnsOnCall = make(map[int]struct {
			result1 hcscontainer.Container
			result2 error
		})
	}
	fake.createContainerReturnsOnCall[i] = struct {
		result1 hcscontainer.Container
		result2 error
	}{result1, result2}
}

func (fake *FakeHCSClient) OpenContainer(arg1 string) (hcscontainer.Container, error) {
	fake.openContainerMutex.Lock()
	ret, specificReturn := fake.openContainerReturnsOnCall[len(fake.openContainerArgsForCall)]
	fake.openContainerArgsForCall = append(fake.openContainerArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("OpenContainer", []interface{}{arg1})
	fake.openContainerMutex.Unlock()
	if fake.OpenContainerStub != nil {
		return fake.OpenContainerStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.openContainerReturns.result1, fake.openContainerReturns.result2
}

func (fake *FakeHCSClient) OpenContainerCallCount() int {
	fake.openContainerMutex.RLock()
	defer fake.openContainerMutex.RUnlock()
	return len(fake.openContainerArgsForCall)
}

func (fake *FakeHCSClient) OpenContainerArgsForCall(i int) string {
	fake.openContainerMutex.RLock()
	defer fake.openContainerMutex.RUnlock()
	return fake.openContainerArgsForCall[i].arg1
}

func (fake *FakeHCSClient) OpenContainerReturns(result1 hcscontainer.Container, result2 error) {
	fake.OpenContainerStub = nil
	fake.openContainerReturns = struct {
		result1 hcscontainer.Container
		result2 error
	}{result1, result2}
}

func (fake *FakeHCSClient) OpenContainerReturnsOnCall(i int, result1 hcscontainer.Container, result2 error) {
	fake.OpenContainerStub = nil
	if fake.openContainerReturnsOnCall == nil {
		fake.openContainerReturnsOnCall = make(map[int]struct {
			result1 hcscontainer.Container
			result2 error
		})
	}
	fake.openContainerReturnsOnCall[i] = struct {
		result1 hcscontainer.Container
		result2 error
	}{result1, result2}
}

func (fake *FakeHCSClient) IsPending(arg1 error) bool {
	fake.isPendingMutex.Lock()
	ret, specificReturn := fake.isPendingReturnsOnCall[len(fake.isPendingArgsForCall)]
	fake.isPendingArgsForCall = append(fake.isPendingArgsForCall, struct {
		arg1 error
	}{arg1})
	fake.recordInvocation("IsPending", []interface{}{arg1})
	fake.isPendingMutex.Unlock()
	if fake.IsPendingStub != nil {
		return fake.IsPendingStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.isPendingReturns.result1
}

func (fake *FakeHCSClient) IsPendingCallCount() int {
	fake.isPendingMutex.RLock()
	defer fake.isPendingMutex.RUnlock()
	return len(fake.isPendingArgsForCall)
}

func (fake *FakeHCSClient) IsPendingArgsForCall(i int) error {
	fake.isPendingMutex.RLock()
	defer fake.isPendingMutex.RUnlock()
	return fake.isPendingArgsForCall[i].arg1
}

func (fake *FakeHCSClient) IsPendingReturns(result1 bool) {
	fake.IsPendingStub = nil
	fake.isPendingReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeHCSClient) IsPendingReturnsOnCall(i int, result1 bool) {
	fake.IsPendingStub = nil
	if fake.isPendingReturnsOnCall == nil {
		fake.isPendingReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isPendingReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *FakeHCSClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getContainerPropertiesMutex.RLock()
	defer fake.getContainerPropertiesMutex.RUnlock()
	fake.nameToGuidMutex.RLock()
	defer fake.nameToGuidMutex.RUnlock()
	fake.createContainerMutex.RLock()
	defer fake.createContainerMutex.RUnlock()
	fake.openContainerMutex.RLock()
	defer fake.openContainerMutex.RUnlock()
	fake.isPendingMutex.RLock()
	defer fake.isPendingMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeHCSClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ container.HCSClient = new(FakeHCSClient)
