package feed

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

func clearImage(filePath string) error {
	absolutePath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("Error converting to absolute path: %w", err)
	}

	if _, err := os.Stat(absolutePath); err == nil {
		if err := os.Remove(absolutePath); err != nil {
			return fmt.Errorf("Error removing file: %w", err)
		} else {
			slog.Info(fmt.Sprintf("%s removed successfully", absolutePath))
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("Error checking file existence: %w", err)
	}

	return nil
}
