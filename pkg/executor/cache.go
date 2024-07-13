package executor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func makeCacheDir(cacheDirPath string) (string, error) {
	cacheDirPath = strings.TrimSpace(cacheDirPath)

	if cacheDirPath == "" {
		userCacheDir, err := os.UserCacheDir()
		if err != nil {
			return "", fmt.Errorf("failed to get the user cache directory: %w", err)
		}

		cacheDirPath = filepath.Join(userCacheDir, extensionName)
	} else {
		cacheDirPath = filepath.Clean(cacheDirPath)
		cacheDirPath = filepath.Join(cacheDirPath, extensionName)
	}

	if err := os.MkdirAll(cacheDirPath, 0750); err != nil {
		return "", fmt.Errorf("failed to create a cache directory: %w", err)
	}

	return cacheDirPath, nil
}
