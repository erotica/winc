package hcsclient

import (
	"code.cloudfoundry.org/winc/hcscontainer"
	"github.com/Microsoft/hcsshim"
)

type HCSClient struct{}

func (c *HCSClient) GetContainers(q hcsshim.ComputeSystemQuery) ([]hcsshim.ContainerProperties, error) {
	return hcsshim.GetContainers(q)
}

func (c *HCSClient) NameToGuid(name string) (hcsshim.GUID, error) {
	return hcsshim.NameToGuid(name)
}

func (c *HCSClient) GetLayerMountPath(info hcsshim.DriverInfo, id string) (string, error) {
	return hcsshim.GetLayerMountPath(info, id)
}

func (c *HCSClient) CreateContainer(id string, config *hcsshim.ContainerConfig) (hcscontainer.Container, error) {
	return hcsshim.CreateContainer(id, config)
}

func (c *HCSClient) OpenContainer(id string) (hcscontainer.Container, error) {
	return hcsshim.OpenContainer(id)
}

func (c *HCSClient) IsPending(err error) bool {
	return hcsshim.IsPending(err)
}

func (c *HCSClient) CreateSandboxLayer(info hcsshim.DriverInfo, layerId, parentId string, parentLayerPaths []string) error {
	return hcsshim.CreateSandboxLayer(info, layerId, parentId, parentLayerPaths)
}

func (c *HCSClient) ActivateLayer(info hcsshim.DriverInfo, id string) error {
	return hcsshim.ActivateLayer(info, id)
}

func (c *HCSClient) PrepareLayer(info hcsshim.DriverInfo, layerId string, parentLayerPaths []string) error {
	return hcsshim.PrepareLayer(info, layerId, parentLayerPaths)
}

func (c *HCSClient) UnprepareLayer(info hcsshim.DriverInfo, layerId string) error {
	return hcsshim.UnprepareLayer(info, layerId)
}

func (c *HCSClient) DeactivateLayer(info hcsshim.DriverInfo, id string) error {
	return hcsshim.DeactivateLayer(info, id)
}

func (c *HCSClient) DestroyLayer(info hcsshim.DriverInfo, id string) error {
	return hcsshim.DestroyLayer(info, id)
}

func (c *HCSClient) LayerExists(info hcsshim.DriverInfo, id string) (bool, error) {
	return hcsshim.LayerExists(info, id)
}

func (c *HCSClient) GetContainerProperties(id string) (hcsshim.ContainerProperties, error) {
	query := hcsshim.ComputeSystemQuery{
		IDs:    []string{id},
		Owners: []string{"winc"},
	}
	cps, err := c.GetContainers(query)
	if err != nil {
		return hcsshim.ContainerProperties{}, err
	}

	if len(cps) == 0 {
		return hcsshim.ContainerProperties{}, &NotFoundError{Id: id}
	}

	if len(cps) > 1 {
		return hcsshim.ContainerProperties{}, &DuplicateError{Id: id}
	}

	return cps[0], nil
}

func (c *HCSClient) CreateEndpoint(endpoint *hcsshim.HNSEndpoint) (*hcsshim.HNSEndpoint, error) {
	return endpoint.Create()
}

func (c *HCSClient) DeleteEndpoint(endpoint *hcsshim.HNSEndpoint) (*hcsshim.HNSEndpoint, error) {
	return endpoint.Delete()
}

func (c *HCSClient) CreateNetwork(network *hcsshim.HNSNetwork) (*hcsshim.HNSNetwork, error) {
	return network.Create()
}

func (c *HCSClient) DeleteNetwork(network *hcsshim.HNSNetwork) (*hcsshim.HNSNetwork, error) {
	return network.Delete()
}

func (c *HCSClient) HNSListNetworkRequest() ([]hcsshim.HNSNetwork, error) {
	return hcsshim.HNSListNetworkRequest("GET", "", "")
}

func (c *HCSClient) GetHNSEndpointByID(id string) (*hcsshim.HNSEndpoint, error) {
	return hcsshim.GetHNSEndpointByID(id)
}

func (c *HCSClient) GetHNSNetworkByName(name string) (*hcsshim.HNSNetwork, error) {
	return hcsshim.GetHNSNetworkByName(name)
}
