// Business logic for `sync` subcommand
package sync

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Devourian/synchro/internal/logging"
)

var log = logging.GetLogger()

// Sync synchronizes files from sourceDirPath to targetDirPath.
//
// The decision to copy or delete files is based on the following rules:
//   - If a file exists in sourceDirPath but not in targetDirPath, it is copied to targetDirPath.
//   - If a file exists in both directories, and the last modified time of file from sourceDirPath is
//     newer than last modified time of file in targetDirPath and the size differs, it is copied to targetDirPath.
//   - If a file exists in both directories, and the size is the same, it is not copied.
//
// If deleteMissing is true, it deletes files in
// targetDirPath that are not present anymore in sourceDirPath.
func Sync(
	sourceDirPath string,
	targetDirPath string,
	deleteMissing bool,
) error {
	sourceFilesMap, err := getFilesMap(sourceDirPath)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	targetFilesMap, err := getFilesMap(targetDirPath)
	if err != nil {
		return fmt.Errorf("failed to read target directory: %w", err)
	}

	for sourceFileName, sourceFileEntry := range sourceFilesMap {
		sourceFilePath := filepath.Join(sourceDirPath, sourceFileName)
		targetFilePath := filepath.Join(targetDirPath, sourceFileName)

		if targetFileEntry, exists := targetFilesMap[sourceFileName]; exists {
			if !isOverwriteNeeded(sourceFileEntry, targetFileEntry) {
				continue
			}

			// File exists in both directories and needs to be overwritten
			log.Info("Overwriting file: %s with %s ...", targetFilePath, sourceFilePath)
			err := copyFile(sourceFilePath, targetFilePath)
			if err != nil {
				log.Error("Failed to overwrite file: %s", err)
				continue
			} else {
				log.Info("File overwritten: %s", targetFilePath)
			}
		} else {
			// File does not exist in targetDirPath, copy it
			log.Info("Copying file: %s to %s ...", sourceFilePath, targetFilePath)
			err = copyFile(sourceFilePath, targetFilePath)
			if err != nil {
				log.Error("Failed to copy file: %s", err)
				continue
			} else {
				log.Info("File copied: %s", targetFilePath)
			}
		}
	}

	if deleteMissing {
		err = deleteMissingFiles(sourceFilesMap, targetFilesMap, targetDirPath)
		if err != nil {
			log.Error("Failed to delete missing files: %s", err)
		}
	}

	return nil
}

// getFilesMap returns a map of files in the specified directory.
func getFilesMap(dirPath string) (map[string]os.DirEntry, error) {
	filesMap := make(map[string]os.DirEntry)

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return filesMap, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			log.Warn("Handling nested directories is not supported yet. Skipping %s", entry.Name())
			continue
		} else {
			filesMap[entry.Name()] = entry
		}
	}

	return filesMap, nil
}

// isOverwriteNeeded checks if the file in sourceDirPath should overwrite the file in targetDirPath.
func isOverwriteNeeded(
	sourceFileEntry os.DirEntry,
	targetFileEntry os.DirEntry,
) bool {
	if sourceFileEntry.Type() == targetFileEntry.Type() {
		sourceFileInfo, err := sourceFileEntry.Info()
		if err != nil {
			log.Error("Failed to get source file info: %s", err)
			return false
		}

		targetFileInfo, err := targetFileEntry.Info()
		if err != nil {
			log.Error("Failed to get target file info: %s", err)
			return false
		}

		if sourceFileInfo.ModTime().After(targetFileInfo.ModTime()) &&
			(sourceFileInfo.Size() != targetFileInfo.Size()) {
			return true
		}
	}

	return false
}

// copyFile copies the contents and permissions from src to dst.
// If dst does not exist, it will be created with the same file mode as src.
// If dst exists, it will be truncated.
func copyFile(src, dst string) error {
	// 1. Open source file for reading
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// 2. Stat source to get its file mode (permissions)
	fi, err := in.Stat()
	if err != nil {
		return err
	}

	// 3. Create (or truncate) the destination file
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fi.Mode())
	if err != nil {
		return err
	}

	// Ensure out is closed and capture any close error
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	// 4. Copy the actual file contents
	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	// 5. (Optionally) flush to stable storage
	if err = out.Sync(); err != nil {
		return err
	}

	return nil
}

// deleteMissingFiles deletes files in targetDirPath that are not present in sourceDirPath.
func deleteMissingFiles(
	sourceFilesMap map[string]os.DirEntry,
	targetFilesMap map[string]os.DirEntry,
	targetDirPath string,
) error {
	for targetFileName := range targetFilesMap {
		if _, exists := sourceFilesMap[targetFileName]; !exists {
			targetFilePath := filepath.Join(targetDirPath, targetFileName)

			log.Info("Deleting file: %s ...", targetFilePath)
			if err := os.Remove(targetFilePath); err != nil {
				return fmt.Errorf("failed to delete file: %w", err)
			}
			log.Info("File deleted: %s", targetFilePath)
		}
	}

	return nil
}
