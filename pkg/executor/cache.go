package executor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func makeCacheDir(dirPath string, dirPerm os.FileMode) (string, error) {
	dirPath = strings.TrimSpace(dirPath)

	if dirPath == "" {
		userCacheDir, err := os.UserCacheDir()
		if err != nil {
			return "", fmt.Errorf("failed to get the user cache directory: %w", err)
		}

		dirPath = filepath.Join(userCacheDir, extensionName)
	} else {
		dirPath = filepath.Join(dirPath, extensionName)
	}

	dirPath = filepath.Clean(dirPath)

	if err := os.MkdirAll(dirPath, dirPerm); err != nil {
		return "", fmt.Errorf("failed to create a cache directory: %w", err)
	}

	return dirPath, nil
}
