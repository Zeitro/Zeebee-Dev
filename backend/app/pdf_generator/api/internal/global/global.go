package global

import (
	"time"

	"zeebee-backend/app/pdf_generator/api/internal/pdf_generator"
)

var ChromePool *pdf_generator.ChromePool

// InitChromePool
func InitChromePool(maxConcurrent int, acquireTimeout time.Duration) error {
	pool, err := pdf_generator.NewChromePool(maxConcurrent, acquireTimeout)
	if err != nil {
		return err
	}
	ChromePool = pool
	return nil
}

func CloseChromePool() {
	if ChromePool != nil {
		_ = ChromePool.Close()
		ChromePool = nil
	}
}
