package downloader

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	_config "github.com/sijunda/govman/internal/config"
	_golang "github.com/sijunda/govman/internal/golang"
	_progress "github.com/sijunda/govman/internal/progress"
)

type Downloader struct {
	config *_config.Config
	client *http.Client
}

func New(cfg *_config.Config) *Downloader {
	return &Downloader{
		config: cfg,
		client: &http.Client{
			Timeout: cfg.Download.Timeout,
		},
	}
}

// Download downloads and installs a Go version
func (d *Downloader) Download(url, installDir, version string) error {
	// Get file info for verification
	fileInfo, err := _golang.GetFileInfo(version)
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Download file
	archivePath, err := d.downloadFile(url, fileInfo)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer os.Remove(archivePath) // Clean up downloaded file

	// Verify checksum
	if err := d.verifyChecksum(archivePath, fileInfo.Sha256); err != nil {
		return fmt.Errorf("checksum verification failed: %w", err)
	}

	// Extract archive
	if err := d.extractArchive(archivePath, installDir); err != nil {
		return fmt.Errorf("failed to extract archive: %w", err)
	}

	return nil
}

func (d *Downloader) downloadFile(url string, fileInfo *_golang.File) (string, error) {
	// Determine cache file path
	filename := filepath.Base(url)
	cachePath := filepath.Join(d.config.CacheDir, filename)

	// Check if file already exists and is complete
	if stat, err := os.Stat(cachePath); err == nil {
		if stat.Size() == fileInfo.Size {
			fmt.Printf("✅ Using cached file: %s\n", filename)
			return cachePath, nil
		}
		fmt.Printf("📦 Resuming download: %s\n", filename)
	} else {
		fmt.Printf("📦 Downloading: %s\n", filename)
	}

	// Open cache file for writing (append mode for resume)
	file, err := os.OpenFile(cachePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to create cache file: %w", err)
	}
	defer file.Close()

	// Get current file size for resume
	stat, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to stat cache file: %w", err)
	}
	currentSize := stat.Size()

	// Create HTTP request with range header for resume
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	if currentSize > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", currentSize))
	}

	// Execute request with retries
	var resp *http.Response
	for attempt := 0; attempt < d.config.Download.RetryCount; attempt++ {
		resp, err = d.client.Do(req)
		if err != nil {
			if attempt < d.config.Download.RetryCount-1 {
				fmt.Printf("⚠️  Download failed, retrying in 5 seconds... (%d/%d)\n",
					attempt+1, d.config.Download.RetryCount)
				time.Sleep(5 * time.Second)
				continue
			}
			return "", fmt.Errorf("failed to download after %d attempts: %w",
				attempt+1, err)
		}
		break
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return "", fmt.Errorf("download failed with status: %s", resp.Status)
	}

	// Create progress bar
	totalSize := fileInfo.Size
	if resp.StatusCode == http.StatusPartialContent {
		totalSize = currentSize + resp.ContentLength
	}

	progressBar := _progress.New(totalSize, fmt.Sprintf("Downloading %s", filename))
	progressBar.Set(currentSize) // Set current progress for resume

	// Download with progress
	reader := io.TeeReader(resp.Body, progressBar)

	if _, err := io.Copy(file, reader); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	progressBar.Finish()
	return cachePath, nil
}

func (d *Downloader) verifyChecksum(filePath, expectedSHA256 string) error {
	fmt.Printf("🔍 Verifying checksum...\n")

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return fmt.Errorf("failed to calculate checksum: %w", err)
	}

	actualSHA256 := fmt.Sprintf("%x", hasher.Sum(nil))
	if actualSHA256 != expectedSHA256 {
		return fmt.Errorf("checksum mismatch: expected %s, got %s",
			expectedSHA256, actualSHA256)
	}

	fmt.Printf("✅ Checksum verified\n")
	return nil
}

func (d *Downloader) extractArchive(archivePath, installDir string) error {
	fmt.Printf("📂 Extracting archive...\n")

	// Create install directory
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}

	// Determine archive type and extract
	if strings.HasSuffix(archivePath, ".tar.gz") {
		return d.extractTarGz(archivePath, installDir)
	} else if strings.HasSuffix(archivePath, ".zip") {
		return d.extractZip(archivePath, installDir)
	}

	return fmt.Errorf("unsupported archive format")
}

func (d *Downloader) extractTarGz(archivePath, installDir string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		// Skip the top-level "go" directory
		path := strings.TrimPrefix(header.Name, "go/")
		if path == "" {
			continue // Skip empty paths
		}

		targetPath := filepath.Join(installDir, path)
		// Create parent directory only if it doesn't exist
		parentDir := filepath.Dir(targetPath)
		if _, err := os.Stat(parentDir); os.IsNotExist(err) {
			if err := os.MkdirAll(parentDir, 0755); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}
		}

		// Create directory
		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}
			continue
		}

		// Create file
		if header.Typeflag == tar.TypeReg {
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}

			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", targetPath, err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("failed to write file %s: %w", targetPath, err)
			}
			outFile.Close()
		}
	}

	return nil
}

func (d *Downloader) extractZip(archivePath, installDir string) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open zip archive: %w", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		// Skip the top-level "go" directory
		path := file.Name
		if strings.HasPrefix(path, "go/") || strings.HasPrefix(path, "go\\") {
			path = path[3:] // Remove "go/" or "go\\" prefix
		}

		if path == "" {
			continue // Skip empty paths
		}

		targetPath := filepath.Join(installDir, path)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(targetPath, file.FileInfo().Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}
			continue
		}

		// Create parent directory only if it doesn't exist
		parentDir := filepath.Dir(targetPath)
		if _, err := os.Stat(parentDir); os.IsNotExist(err) {
			if err := os.MkdirAll(parentDir, 0755); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}
		}

		// Extract file
		srcFile, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in archive: %w", err)
		}

		dstFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, file.FileInfo().Mode())
		if err != nil {
			srcFile.Close()
			return fmt.Errorf("failed to create file %s: %w", targetPath, err)
		}

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			srcFile.Close()
			dstFile.Close()
			return fmt.Errorf("failed to write file %s: %w", targetPath, err)
		}

		srcFile.Close()
		dstFile.Close()
	}

	return nil
}
