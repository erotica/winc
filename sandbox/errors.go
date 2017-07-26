package sandbox

import "fmt"

type MissingRootfsError struct {
	Msg string
}

func (e *MissingRootfsError) Error() string {
	return fmt.Sprintf("rootfs layer does not exist: %s", e.Msg)
}

type MissingRootfsLayerChainError struct {
	Msg string
}

func (e *MissingRootfsLayerChainError) Error() string {
	return fmt.Sprintf("rootfs does not contain a layerchain.json: %s", e.Msg)
}

type InvalidRootfsLayerChainError struct {
	Path string
}

func (e *InvalidRootfsLayerChainError) Error() string {
	return fmt.Sprintf("rootfs contains an invalid layerchain.json: %s", e.Path)
}

type MissingBundlePathError struct {
	Msg string
}

func (e *MissingBundlePathError) Error() string {
	return fmt.Sprintf("bundlePath does not exist: %s", e.Msg)
}

type UnableToDestroyLayerError struct {
	Msg string
}

func (e *UnableToDestroyLayerError) Error() string {
	return fmt.Sprintf("unable to destroy layer file: %s", e.Msg)
}
