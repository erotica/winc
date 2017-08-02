package volume

type Limiter struct{}

func (l *Limiter) SetDiskLimit(volumePath string, size int64) error {
	return nil
}
