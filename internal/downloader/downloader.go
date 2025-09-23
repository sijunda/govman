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
	_logger "github.com/sijunda/govman/internal/logger"
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
	_logger.InternalProgress("Retrieving file information")
	timer := _logger.StartTimer("file info retrieval")
	fileInfo, err := _golang.GetFileInfoWithConfig(version,
		d.config.GoReleases.APIURL,
		d.config.GoReleases.CacheExpiry)
	if err != nil {
		_logger.StopTimer(timer)
		return fmt.Errorf("failed to get file info: %w", err)
	}
	_logger.StopTimer(timer)

	// Download file
	_logger.InternalProgress("Downloading file")
	archivePath, err := d.downloadFile(url, fileInfo)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer os.Remove(archivePath) // Clean up downloaded file

	// Verify checksum
	_logger.InternalProgress("Verifying checksum")
	timer = _logger.StartTimer("checksum verification")
	if err := d.verifyChecksum(archivePath, fileInfo.Sha256); err != nil {
		_logger.StopTimer(timer)
		return fmt.Errorf("checksum verification failed: %w", err)
	}
	_logger.StopTimer(timer)

	// Extract archive
	_logger.InternalProgress("Extracting archive")
	timer = _logger.StartTimer("archive extraction")
	if err := d.extractArchive(archivePath, installDir); err != nil {
		_logger.StopTimer(timer)
		return fmt.Errorf("failed to extract archive: %w", err)
	}
	_logger.StopTimer(timer)

	return nil
}

func (d *Downloader) downloadFile(url string, fileInfo *_golang.File) (string, error) {
	// Determine cache file path
	filename := filepath.Base(url)
	cachePath := filepath.Join(d.config.CacheDir, filename)

	// Check if file already exists and is complete
	if stat, err := os.Stat(cachePath); err == nil {
		if stat.Size() == fileInfo.Size {
			_logger.Success("Using cached file: %s", filename)
			return cachePath, nil
		}
		_logger.Download("Resuming download: %s", filename)
	} else {
		_logger.Download("Downloading: %s", filename)
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
				_logger.Warning("Download failed, retrying in 5 seconds... (%d/%d)",
					attempt+1, d.config.Download.RetryCount)
				time.Sleep(d.config.Download.RetryDelay)
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
	if progressBar != nil {
		progressBar.Set(currentSize) // Set current progress for resume
	}

	// Download with progress
	var reader io.Reader
	if progressBar != nil {
		reader = io.TeeReader(resp.Body, progressBar)
	} else {
		reader = resp.Body
	}

	// Note: Uncomment this to show the progress bar only when verbose mode is enabled
	// var progressBar *_progress.ProgressBar
	// if _logger.Get().Level() >= _logger.VerboseLevel {
	// 	progressBar = _progress.New(totalSize, fmt.Sprintf("Downloading %s", filename))
	// 	progressBar.Set(currentSize) // Set current progress for resume
	// }

	// var reader io.Reader
	// if progressBar != nil {
	// 	reader = io.TeeReader(resp.Body, progressBar)
	// } else {
	// 	reader = resp.Body
	// }

	if _, err := io.Copy(file, reader); err != nil {
		// Ensure the file is closed before returning
		file.Close()
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	if progressBar != nil {
		progressBar.Finish()
	}
	return cachePath, nil
}

func (d *Downloader) verifyChecksum(filePath, expectedSHA256 string) error {
	_logger.Verify("Verifying checksum...")

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

	_logger.Success("Checksum verified")
	return nil
}

func (d *Downloader) extractArchive(archivePath, installDir string) error {
	_logger.Extract("Extracting archive...")

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

		// Validate that the path is safe and doesn't contain traversal sequences
		if strings.Contains(path, "..") || filepath.IsAbs(path) {
			return fmt.Errorf("unsafe path in archive: %s", header.Name)
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

		// Validate that the path is safe and doesn't contain traversal sequences
		if strings.Contains(path, "..") || filepath.IsAbs(path) {
			return fmt.Errorf("unsafe path in archive: %s", file.Name)
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
