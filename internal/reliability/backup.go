package reliability

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	// File permissions
	FileModeReadWrite = 0600 // rw------- for files
	FileModeReadOnly  = 0400 // r-------- for files
	DirModeReadWrite  = 0700 // rwx------ for directories
	DirModeReadOnly   = 0500 // r-x------ for directories
)

// BackupFile creates a backup of a file before modification
func BackupFile(filePath string) error {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // No backup needed for new files
	}

	// Create backup directory if it doesn't exist
	backupDir := filepath.Join(filepath.Dir(filePath), ".templater_backups")
	if err := os.MkdirAll(backupDir, DirModeReadWrite); err != nil {
		return fmt.Errorf("error creating backup directory: %w", err)
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("20060102_150405")
	backupName := fmt.Sprintf("%s_%s.bak", filepath.Base(filePath), timestamp)
	backupPath := filepath.Join(backupDir, backupName)

	// Read original file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file for backup: %w", err)
	}

	// Write backup file with restricted permissions
	if err := os.WriteFile(backupPath, content, FileModeReadWrite); err != nil {
		return fmt.Errorf("error writing backup file: %w", err)
	}

	return nil
}

// RestoreFromBackup restores a file from its most recent backup
func RestoreFromBackup(filePath string) error {
	backupDir := filepath.Join(filepath.Dir(filePath), ".templater_backups")

	// Find most recent backup
	backups, err := filepath.Glob(filepath.Join(backupDir, filepath.Base(filePath)+"_*.bak"))
	if err != nil {
		return fmt.Errorf("error finding backups: %w", err)
	}

	if len(backups) == 0 {
		return fmt.Errorf("no backups found for %s", filePath)
	}

	// Get most recent backup
	var mostRecent string
	var mostRecentTime time.Time
	for _, backup := range backups {
		info, err := os.Stat(backup)
		if err != nil {
			continue
		}
		if info.ModTime().After(mostRecentTime) {
			mostRecentTime = info.ModTime()
			mostRecent = backup
		}
	}

	if mostRecent == "" {
		return fmt.Errorf("no valid backups found for %s", filePath)
	}

	// Read backup file
	content, err := os.ReadFile(mostRecent)
	if err != nil {
		return fmt.Errorf("error reading backup file: %w", err)
	}

	// Restore original file with restricted permissions
	if err := os.WriteFile(filePath, content, FileModeReadWrite); err != nil {
		return fmt.Errorf("error restoring from backup: %w", err)
	}

	return nil
}

// VerifyFileIntegrity checks if a file's content matches its expected state
func VerifyFileIntegrity(filePath string, expectedContent []byte) error {
	actualContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file for integrity check: %w", err)
	}

	if string(actualContent) != string(expectedContent) {
		return fmt.Errorf("file integrity check failed: content mismatch")
	}

	return nil
}

// CleanupOldBackups removes backups older than the specified duration
func CleanupOldBackups(filePath string, maxAge time.Duration) error {
	backupDir := filepath.Join(filepath.Dir(filePath), ".templater_backups")

	// Find all backups
	backups, err := filepath.Glob(filepath.Join(backupDir, filepath.Base(filePath)+"_*.bak"))
	if err != nil {
		return fmt.Errorf("error finding backups: %w", err)
	}

	// Remove old backups
	now := time.Now()
	for _, backup := range backups {
		info, err := os.Stat(backup)
		if err != nil {
			continue
		}

		if now.Sub(info.ModTime()) > maxAge {
			if err := os.Remove(backup); err != nil {
				return fmt.Errorf("error removing old backup %s: %w", backup, err)
			}
		}
	}

	return nil
}
