// Code generated by counterfeiter. DO NOT EDIT.
package imagefakes

import (
	"sync"

	"code.cloudfoundry.org/winc/image"
)

type FakeLimiter struct {
	SetDiskLimitStub        func(volumePath string, size uint64) error
	setDiskLimitMutex       sync.RWMutex
	setDiskLimitArgsForCall []struct {
		volumePath string
		size       uint64
	}
	setDiskLimitReturns struct {
		result1 error
	}
	setDiskLimitReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeLimiter) SetDiskLimit(volumePath string, size uint64) error {
	fake.setDiskLimitMutex.Lock()
	ret, specificReturn := fake.setDiskLimitReturnsOnCall[len(fake.setDiskLimitArgsForCall)]
	fake.setDiskLimitArgsForCall = append(fake.setDiskLimitArgsForCall, struct {
		volumePath string
		size       uint64
	}{volumePath, size})
	fake.recordInvocation("SetDiskLimit", []interface{}{volumePath, size})
	fake.setDiskLimitMutex.Unlock()
	if fake.SetDiskLimitStub != nil {
		return fake.SetDiskLimitStub(volumePath, size)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.setDiskLimitReturns.result1
}

func (fake *FakeLimiter) SetDiskLimitCallCount() int {
	fake.setDiskLimitMutex.RLock()
	defer fake.setDiskLimitMutex.RUnlock()
	return len(fake.setDiskLimitArgsForCall)
}

func (fake *FakeLimiter) SetDiskLimitArgsForCall(i int) (string, uint64) {
	fake.setDiskLimitMutex.RLock()
	defer fake.setDiskLimitMutex.RUnlock()
	return fake.setDiskLimitArgsForCall[i].volumePath, fake.setDiskLimitArgsForCall[i].size
}

func (fake *FakeLimiter) SetDiskLimitReturns(result1 error) {
	fake.SetDiskLimitStub = nil
	fake.setDiskLimitReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeLimiter) SetDiskLimitReturnsOnCall(i int, result1 error) {
	fake.SetDiskLimitStub = nil
	if fake.setDiskLimitReturnsOnCall == nil {
		fake.setDiskLimitReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setDiskLimitReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeLimiter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.setDiskLimitMutex.RLock()
	defer fake.setDiskLimitMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeLimiter) recordInvocation(key string, args []interface{}) {
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

var _ image.Limiter = new(FakeLimiter)
