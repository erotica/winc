// Code generated by counterfeiter. DO NOT EDIT.
package layerfakes

import (
	"sync"

	"code.cloudfoundry.org/winc/layer"
	"github.com/Microsoft/hcsshim"
)

type FakeHCSClient struct {
	CreateSandboxLayerStub        func(hcsshim.DriverInfo, string, string, []string) error
	createSandboxLayerMutex       sync.RWMutex
	createSandboxLayerArgsForCall []struct {
		arg1 hcsshim.DriverInfo
		arg2 string
		arg3 string
		arg4 []string
	}
	createSandboxLayerReturns struct {
		result1 error
	}
	createSandboxLayerReturnsOnCall map[int]struct {
		result1 error
	}
	ActivateLayerStub        func(hcsshim.DriverInfo, string) error
	activateLayerMutex       sync.RWMutex
	activateLayerArgsForCall []struct {
		arg1 hcsshim.DriverInfo
		arg2 string
	}
	activateLayerReturns struct {
		result1 error
	}
	activateLayerReturnsOnCall map[int]struct {
		result1 error
	}
	PrepareLayerStub        func(hcsshim.DriverInfo, string, []string) error
	prepareLayerMutex       sync.RWMutex
	prepareLayerArgsForCall []struct {
		arg1 hcsshim.DriverInfo
		arg2 string
		arg3 []string
	}
	prepareLayerReturns struct {
		result1 error
	}
	prepareLayerReturnsOnCall map[int]struct {
		result1 error
	}
	UnprepareLayerStub        func(hcsshim.DriverInfo, string) error
	unprepareLayerMutex       sync.RWMutex
	unprepareLayerArgsForCall []struct {
		arg1 hcsshim.DriverInfo
		arg2 string
	}
	unprepareLayerReturns struct {
		result1 error
	}
	unprepareLayerReturnsOnCall map[int]struct {
		result1 error
	}
	DeactivateLayerStub        func(hcsshim.DriverInfo, string) error
	deactivateLayerMutex       sync.RWMutex
	deactivateLayerArgsForCall []struct {
		arg1 hcsshim.DriverInfo
		arg2 string
	}
	deactivateLayerReturns struct {
		result1 error
	}
	deactivateLayerReturnsOnCall map[int]struct {
		result1 error
	}
	DestroyLayerStub        func(hcsshim.DriverInfo, string) error
	destroyLayerMutex       sync.RWMutex
	destroyLayerArgsForCall []struct {
		arg1 hcsshim.DriverInfo
		arg2 string
	}
	destroyLayerReturns struct {
		result1 error
	}
	destroyLayerReturnsOnCall map[int]struct {
		result1 error
	}
	LayerExistsStub        func(hcsshim.DriverInfo, string) (bool, error)
	layerExistsMutex       sync.RWMutex
	layerExistsArgsForCall []struct {
		arg1 hcsshim.DriverInfo
		arg2 string
	}
	layerExistsReturns struct {
		result1 bool
		result2 error
	}
	layerExistsReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	GetLayerMountPathStub        func(hcsshim.DriverInfo, string) (string, error)
	getLayerMountPathMutex       sync.RWMutex
	getLayerMountPathArgsForCall []struct {
		arg1 hcsshim.DriverInfo
		arg2 string
	}
	getLayerMountPathReturns struct {
		result1 string
		result2 error
	}
	getLayerMountPathReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeHCSClient) CreateSandboxLayer(arg1 hcsshim.DriverInfo, arg2 string, arg3 string, arg4 []string) error {
	var arg4Copy []string
	if arg4 != nil {
		arg4Copy = make([]string, len(arg4))
		copy(arg4Copy, arg4)
	}
	fake.createSandboxLayerMutex.Lock()
	ret, specificReturn := fake.createSandboxLayerReturnsOnCall[len(fake.createSandboxLayerArgsForCall)]
	fake.createSandboxLayerArgsForCall = append(fake.createSandboxLayerArgsForCall, struct {
		arg1 hcsshim.DriverInfo
		arg2 string
		arg3 string
		arg4 []string
	}{arg1, arg2, arg3, arg4Copy})
	fake.recordInvocation("CreateSandboxLayer", []interface{}{arg1, arg2, arg3, arg4Copy})
	fake.createSandboxLayerMutex.Unlock()
	if fake.CreateSandboxLayerStub != nil {
		return fake.CreateSandboxLayerStub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.createSandboxLayerReturns.result1
}

func (fake *FakeHCSClient) CreateSandboxLayerCallCount() int {
	fake.createSandboxLayerMutex.RLock()
	defer fake.createSandboxLayerMutex.RUnlock()
	return len(fake.createSandboxLayerArgsForCall)
}

func (fake *FakeHCSClient) CreateSandboxLayerArgsForCall(i int) (hcsshim.DriverInfo, string, string, []string) {
	fake.createSandboxLayerMutex.RLock()
	defer fake.createSandboxLayerMutex.RUnlock()
	return fake.createSandboxLayerArgsForCall[i].arg1, fake.createSandboxLayerArgsForCall[i].arg2, fake.createSandboxLayerArgsForCall[i].arg3, fake.createSandboxLayerArgsForCall[i].arg4
}

func (fake *FakeHCSClient) CreateSandboxLayerReturns(result1 error) {
	fake.CreateSandboxLayerStub = nil
	fake.createSandboxLayerReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeHCSClient) CreateSandboxLayerReturnsOnCall(i int, result1 error) {
	fake.CreateSandboxLayerStub = nil
	if fake.createSandboxLayerReturnsOnCall == nil {
		fake.createSandboxLayerReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.createSandboxLayerReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeHCSClient) ActivateLayer(arg1 hcsshim.DriverInfo, arg2 string) error {
	fake.activateLayerMutex.Lock()
	ret, specificReturn := fake.activateLayerReturnsOnCall[len(fake.activateLayerArgsForCall)]
	fake.activateLayerArgsForCall = append(fake.activateLayerArgsForCall, struct {
		arg1 hcsshim.DriverInfo
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("ActivateLayer", []interface{}{arg1, arg2})
	fake.activateLayerMutex.Unlock()
	if fake.ActivateLayerStub != nil {
		return fake.ActivateLayerStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.activateLayerReturns.result1
}

func (fake *FakeHCSClient) ActivateLayerCallCount() int {
	fake.activateLayerMutex.RLock()
	defer fake.activateLayerMutex.RUnlock()
	return len(fake.activateLayerArgsForCall)
}

func (fake *FakeHCSClient) ActivateLayerArgsForCall(i int) (hcsshim.DriverInfo, string) {
	fake.activateLayerMutex.RLock()
	defer fake.activateLayerMutex.RUnlock()
	return fake.activateLayerArgsForCall[i].arg1, fake.activateLayerArgsForCall[i].arg2
}

func (fake *FakeHCSClient) ActivateLayerReturns(result1 error) {
	fake.ActivateLayerStub = nil
	fake.activateLayerReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeHCSClient) ActivateLayerReturnsOnCall(i int, result1 error) {
	fake.ActivateLayerStub = nil
	if fake.activateLayerReturnsOnCall == nil {
		fake.activateLayerReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.activateLayerReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeHCSClient) PrepareLayer(arg1 hcsshim.DriverInfo, arg2 string, arg3 []string) error {
	var arg3Copy []string
	if arg3 != nil {
		arg3Copy = make([]string, len(arg3))
		copy(arg3Copy, arg3)
	}
	fake.prepareLayerMutex.Lock()
	ret, specificReturn := fake.prepareLayerReturnsOnCall[len(fake.prepareLayerArgsForCall)]
	fake.prepareLayerArgsForCall = append(fake.prepareLayerArgsForCall, struct {
		arg1 hcsshim.DriverInfo
		arg2 string
		arg3 []string
	}{arg1, arg2, arg3Copy})
	fake.recordInvocation("PrepareLayer", []interface{}{arg1, arg2, arg3Copy})
	fake.prepareLayerMutex.Unlock()
	if fake.PrepareLayerStub != nil {
		return fake.PrepareLayerStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.prepareLayerReturns.result1
}

func (fake *FakeHCSClient) PrepareLayerCallCount() int {
	fake.prepareLayerMutex.RLock()
	defer fake.prepareLayerMutex.RUnlock()
	return len(fake.prepareLayerArgsForCall)
}

func (fake *FakeHCSClient) PrepareLayerArgsForCall(i int) (hcsshim.DriverInfo, string, []string) {
	fake.prepareLayerMutex.RLock()
	defer fake.prepareLayerMutex.RUnlock()
	return fake.prepareLayerArgsForCall[i].arg1, fake.prepareLayerArgsForCall[i].arg2, fake.prepareLayerArgsForCall[i].arg3
}

func (fake *FakeHCSClient) PrepareLayerReturns(result1 error) {
	fake.PrepareLayerStub = nil
	fake.prepareLayerReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeHCSClient) PrepareLayerReturnsOnCall(i int, result1 error) {
	fake.PrepareLayerStub = nil
	if fake.prepareLayerReturnsOnCall == nil {
		fake.prepareLayerReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.prepareLayerReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeHCSClient) UnprepareLayer(arg1 hcsshim.DriverInfo, arg2 string) error {
	fake.unprepareLayerMutex.Lock()
	ret, specificReturn := fake.unprepareLayerReturnsOnCall[len(fake.unprepareLayerArgsForCall)]
	fake.unprepareLayerArgsForCall = append(fake.unprepareLayerArgsForCall, struct {
		arg1 hcsshim.DriverInfo
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("UnprepareLayer", []interface{}{arg1, arg2})
	fake.unprepareLayerMutex.Unlock()
	if fake.UnprepareLayerStub != nil {
		return fake.UnprepareLayerStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.unprepareLayerReturns.result1
}

func (fake *FakeHCSClient) UnprepareLayerCallCount() int {
	fake.unprepareLayerMutex.RLock()
	defer fake.unprepareLayerMutex.RUnlock()
	return len(fake.unprepareLayerArgsForCall)
}

func (fake *FakeHCSClient) UnprepareLayerArgsForCall(i int) (hcsshim.DriverInfo, string) {
	fake.unprepareLayerMutex.RLock()
	defer fake.unprepareLayerMutex.RUnlock()
	return fake.unprepareLayerArgsForCall[i].arg1, fake.unprepareLayerArgsForCall[i].arg2
}

func (fake *FakeHCSClient) UnprepareLayerReturns(result1 error) {
	fake.UnprepareLayerStub = nil
	fake.unprepareLayerReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeHCSClient) UnprepareLayerReturnsOnCall(i int, result1 error) {
	fake.UnprepareLayerStub = nil
	if fake.unprepareLayerReturnsOnCall == nil {
		fake.unprepareLayerReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.unprepareLayerReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeHCSClient) DeactivateLayer(arg1 hcsshim.DriverInfo, arg2 string) error {
	fake.deactivateLayerMutex.Lock()
	ret, specificReturn := fake.deactivateLayerReturnsOnCall[len(fake.deactivateLayerArgsForCall)]
	fake.deactivateLayerArgsForCall = append(fake.deactivateLayerArgsForCall, struct {
		arg1 hcsshim.DriverInfo
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("DeactivateLayer", []interface{}{arg1, arg2})
	fake.deactivateLayerMutex.Unlock()
	if fake.DeactivateLayerStub != nil {
		return fake.DeactivateLayerStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.deactivateLayerReturns.result1
}

func (fake *FakeHCSClient) DeactivateLayerCallCount() int {
	fake.deactivateLayerMutex.RLock()
	defer fake.deactivateLayerMutex.RUnlock()
	return len(fake.deactivateLayerArgsForCall)
}

func (fake *FakeHCSClient) DeactivateLayerArgsForCall(i int) (hcsshim.DriverInfo, string) {
	fake.deactivateLayerMutex.RLock()
	defer fake.deactivateLayerMutex.RUnlock()
	return fake.deactivateLayerArgsForCall[i].arg1, fake.deactivateLayerArgsForCall[i].arg2
}

func (fake *FakeHCSClient) DeactivateLayerReturns(result1 error) {
	fake.DeactivateLayerStub = nil
	fake.deactivateLayerReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeHCSClient) DeactivateLayerReturnsOnCall(i int, result1 error) {
	fake.DeactivateLayerStub = nil
	if fake.deactivateLayerReturnsOnCall == nil {
		fake.deactivateLayerReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deactivateLayerReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeHCSClient) DestroyLayer(arg1 hcsshim.DriverInfo, arg2 string) error {
	fake.destroyLayerMutex.Lock()
	ret, specificReturn := fake.destroyLayerReturnsOnCall[len(fake.destroyLayerArgsForCall)]
	fake.destroyLayerArgsForCall = append(fake.destroyLayerArgsForCall, struct {
		arg1 hcsshim.DriverInfo
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("DestroyLayer", []interface{}{arg1, arg2})
	fake.destroyLayerMutex.Unlock()
	if fake.DestroyLayerStub != nil {
		return fake.DestroyLayerStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.destroyLayerReturns.result1
}

func (fake *FakeHCSClient) DestroyLayerCallCount() int {
	fake.destroyLayerMutex.RLock()
	defer fake.destroyLayerMutex.RUnlock()
	return len(fake.destroyLayerArgsForCall)
}

func (fake *FakeHCSClient) DestroyLayerArgsForCall(i int) (hcsshim.DriverInfo, string) {
	fake.destroyLayerMutex.RLock()
	defer fake.destroyLayerMutex.RUnlock()
	return fake.destroyLayerArgsForCall[i].arg1, fake.destroyLayerArgsForCall[i].arg2
}

func (fake *FakeHCSClient) DestroyLayerReturns(result1 error) {
	fake.DestroyLayerStub = nil
	fake.destroyLayerReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeHCSClient) DestroyLayerReturnsOnCall(i int, result1 error) {
	fake.DestroyLayerStub = nil
	if fake.destroyLayerReturnsOnCall == nil {
		fake.destroyLayerReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.destroyLayerReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeHCSClient) LayerExists(arg1 hcsshim.DriverInfo, arg2 string) (bool, error) {
	fake.layerExistsMutex.Lock()
	ret, specificReturn := fake.layerExistsReturnsOnCall[len(fake.layerExistsArgsForCall)]
	fake.layerExistsArgsForCall = append(fake.layerExistsArgsForCall, struct {
		arg1 hcsshim.DriverInfo
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("LayerExists", []interface{}{arg1, arg2})
	fake.layerExistsMutex.Unlock()
	if fake.LayerExistsStub != nil {
		return fake.LayerExistsStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.layerExistsReturns.result1, fake.layerExistsReturns.result2
}

func (fake *FakeHCSClient) LayerExistsCallCount() int {
	fake.layerExistsMutex.RLock()
	defer fake.layerExistsMutex.RUnlock()
	return len(fake.layerExistsArgsForCall)
}

func (fake *FakeHCSClient) LayerExistsArgsForCall(i int) (hcsshim.DriverInfo, string) {
	fake.layerExistsMutex.RLock()
	defer fake.layerExistsMutex.RUnlock()
	return fake.layerExistsArgsForCall[i].arg1, fake.layerExistsArgsForCall[i].arg2
}

func (fake *FakeHCSClient) LayerExistsReturns(result1 bool, result2 error) {
	fake.LayerExistsStub = nil
	fake.layerExistsReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeHCSClient) LayerExistsReturnsOnCall(i int, result1 bool, result2 error) {
	fake.LayerExistsStub = nil
	if fake.layerExistsReturnsOnCall == nil {
		fake.layerExistsReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.layerExistsReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeHCSClient) GetLayerMountPath(arg1 hcsshim.DriverInfo, arg2 string) (string, error) {
	fake.getLayerMountPathMutex.Lock()
	ret, specificReturn := fake.getLayerMountPathReturnsOnCall[len(fake.getLayerMountPathArgsForCall)]
	fake.getLayerMountPathArgsForCall = append(fake.getLayerMountPathArgsForCall, struct {
		arg1 hcsshim.DriverInfo
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("GetLayerMountPath", []interface{}{arg1, arg2})
	fake.getLayerMountPathMutex.Unlock()
	if fake.GetLayerMountPathStub != nil {
		return fake.GetLayerMountPathStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getLayerMountPathReturns.result1, fake.getLayerMountPathReturns.result2
}

func (fake *FakeHCSClient) GetLayerMountPathCallCount() int {
	fake.getLayerMountPathMutex.RLock()
	defer fake.getLayerMountPathMutex.RUnlock()
	return len(fake.getLayerMountPathArgsForCall)
}

func (fake *FakeHCSClient) GetLayerMountPathArgsForCall(i int) (hcsshim.DriverInfo, string) {
	fake.getLayerMountPathMutex.RLock()
	defer fake.getLayerMountPathMutex.RUnlock()
	return fake.getLayerMountPathArgsForCall[i].arg1, fake.getLayerMountPathArgsForCall[i].arg2
}

func (fake *FakeHCSClient) GetLayerMountPathReturns(result1 string, result2 error) {
	fake.GetLayerMountPathStub = nil
	fake.getLayerMountPathReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeHCSClient) GetLayerMountPathReturnsOnCall(i int, result1 string, result2 error) {
	fake.GetLayerMountPathStub = nil
	if fake.getLayerMountPathReturnsOnCall == nil {
		fake.getLayerMountPathReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.getLayerMountPathReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeHCSClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createSandboxLayerMutex.RLock()
	defer fake.createSandboxLayerMutex.RUnlock()
	fake.activateLayerMutex.RLock()
	defer fake.activateLayerMutex.RUnlock()
	fake.prepareLayerMutex.RLock()
	defer fake.prepareLayerMutex.RUnlock()
	fake.unprepareLayerMutex.RLock()
	defer fake.unprepareLayerMutex.RUnlock()
	fake.deactivateLayerMutex.RLock()
	defer fake.deactivateLayerMutex.RUnlock()
	fake.destroyLayerMutex.RLock()
	defer fake.destroyLayerMutex.RUnlock()
	fake.layerExistsMutex.RLock()
	defer fake.layerExistsMutex.RUnlock()
	fake.getLayerMountPathMutex.RLock()
	defer fake.getLayerMountPathMutex.RUnlock()
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

var _ layer.HCSClient = new(FakeHCSClient)
