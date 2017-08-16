package sandbox

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	specs "github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
)

const CREATE_ATTEMPTS = 5
const DESTROY_ATTEMPTS = 10

type ImageSpec struct {
	RootFs string `json:"rootfs,omitempty"`
	specs.Spec
}

type DiskUsage struct {
	TotalBytesUsed     uint64 `json:"total_bytes_used"`
	ExclusiveBytesUsed uint64 `json:"exclusive_bytes_used"`
}

type ImageStats struct {
	Disk DiskUsage `json:"disk_usage"`
}

//go:generate counterfeiter . Limiter
type Limiter interface {
	SetDiskLimit(volumePath string, size uint64) error
}

//go:generate counterfeiter . Statser
type Statser interface {
	GetCurrentDiskUsage(string) (uint64, error)
}

//go:generate counterfeiter . LayerManager
type LayerManager interface {
	CreateLayer(string, string, []string) (string, error)
	DestroyLayer(string) error
	LayerExists(string) (bool, error)
	GetLayerMountPath(string) (string, error)
	Retryable(error) bool
	HomeDir() string
}

type Manager struct {
	layerManager LayerManager
	limiter      Limiter
	stats        Statser
	id           string
}

func NewManager(layerManager LayerManager, limiter Limiter, statser Statser, containerId string) *Manager {
	return &Manager{
		layerManager: layerManager,
		limiter:      limiter,
		stats:        statser,
		id:           containerId,
	}
}

func (s *Manager) Create(rootfs string, diskLimit uint64) (*ImageSpec, error) {
	_, err := os.Stat(rootfs)
	if os.IsNotExist(err) {
		return nil, &MissingRootfsError{Msg: rootfs}
	} else if err != nil {
		return nil, err
	}

	parentLayerChain, err := ioutil.ReadFile(filepath.Join(rootfs, "layerchain.json"))
	if err != nil {
		return nil, &MissingRootfsLayerChainError{Msg: rootfs}
	}

	exists, err := s.layerManager.LayerExists(s.id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, &LayerExistsError{Id: s.id}
	}

	parentLayers := []string{}
	if err := json.Unmarshal(parentLayerChain, &parentLayers); err != nil {
		return nil, &InvalidRootfsLayerChainError{Msg: rootfs}
	}
	sandboxLayers := append([]string{rootfs}, parentLayers...)

	var volumePath string
	var createErr error
	for i := 0; i < CREATE_ATTEMPTS; i++ {
		volumePath, createErr = s.layerManager.CreateLayer(s.id, rootfs, sandboxLayers)
		if createErr == nil || !s.layerManager.Retryable(createErr) {
			break
		}
	}
	if createErr != nil {
		_ = s.Delete()
		return nil, createErr
	}

	if err := s.limiter.SetDiskLimit(volumePath, diskLimit); err != nil {
		_ = s.Delete()
		return nil, err
	}

	volumeSize, err := s.stats.GetCurrentDiskUsage(volumePath)
	if err != nil {
		_ = s.Delete()
		return nil, err
	}

	err = ioutil.WriteFile(filepath.Join(s.layerManager.HomeDir(), s.id, "image_info"), []byte(strconv.FormatUint(volumeSize, 10)), 0644)
	if err != nil {
		_ = s.Delete()
		return nil, err
	}

	return &ImageSpec{
		RootFs: volumePath,
		Spec: specs.Spec{
			Root: &specs.Root{
				Path: volumePath,
			},
			Windows: &specs.Windows{
				LayerFolders: sandboxLayers,
			},
		},
	}, nil
}

func (s *Manager) Delete() error {
	exists, err := s.layerManager.LayerExists(s.id)
	if err != nil {
		return err
	}
	if !exists {
		logrus.Warningf("Layer `%s` not found. Skipping delete.", s.id)
		return nil
	}

	var destroyErr error
	for i := 0; i < DESTROY_ATTEMPTS; i++ {
		destroyErr = s.layerManager.DestroyLayer(s.id)
		if destroyErr == nil || !s.layerManager.Retryable(destroyErr) {
			break
		}
	}

	return destroyErr
}

func (s *Manager) Stats() (*ImageStats, error) {
	volumePath, err := s.layerManager.GetLayerMountPath(s.id)
	if err != nil {
		return nil, err
	}

	totalUsed, err := s.stats.GetCurrentDiskUsage(volumePath)
	if err != nil {
		return nil, err
	}

	vs, err := ioutil.ReadFile(filepath.Join(s.layerManager.HomeDir(), s.id, "image_info"))
	if err != nil {
		return nil, err
	}

	volumeSize, err := strconv.ParseUint(string(vs), 10, 64)
	if err != nil {
		return nil, err
	}

	return &ImageStats{
		Disk: DiskUsage{
			TotalBytesUsed:     totalUsed,
			ExclusiveBytesUsed: totalUsed - volumeSize,
		},
	}, nil
}
