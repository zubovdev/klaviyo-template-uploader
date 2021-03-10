package scanner

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"klavio-template/internal/app/config"
	"klavio-template/internal/app/email"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	ScannedFileName          = ".scanned"
	ScannedFileDelimiter     = "    "
	ScannedFileMode          = 0744
	DirectoryProcessDuration = time.Second * 5
)

// Scanner ...
type Scanner struct {
	workDir, apiKey string
	scanned         map[string]int64
	scannedFile     *os.File
}

// NewScanner ...
func NewScanner(cfg config.App) *Scanner {
	return &Scanner{
		workDir: filepath.Clean(cfg.TemplatePath),
		apiKey:  cfg.ApiKey,
		scanned: make(map[string]int64),
	}
}

// Start ...
func (s *Scanner) Start() error {
	if err := s.loadScanned(); err != nil {
		return err
	}

	s.scan()
	return nil
}

// loadScanned ...
func (s *Scanner) loadScanned() error {
	s.scanned = make(map[string]int64)
	scannedPath := path.Join(s.workDir, ScannedFileName)

	if _, err := os.Stat(scannedPath); err == nil {
		s.scannedFile, err = os.OpenFile(scannedPath, os.O_RDWR|os.O_APPEND, ScannedFileMode)
		if err != nil {
			return err
		}

		// Load upload history.
		scanner := bufio.NewScanner(s.scannedFile)
		for scanner.Scan() {
			val := strings.Split(scanner.Text(), ScannedFileDelimiter)
			if len(val) == 2 {
				s.scanned[val[0]], _ = strconv.ParseInt(val[1], 10, 64)
			}
		}
	} else if os.IsNotExist(err) {
		s.scannedFile, err = os.Create(scannedPath)
		if err != nil {
			return err
		}
		_ = s.scannedFile.Chmod(ScannedFileMode)
	} else {
		return err
	}

	return nil
}

// scan ...
func (s *Scanner) scan() {
	fmt.Printf("Scan directory \"%s\" started successfully.\n", s.workDir)
	defer s.scannedFile.Close()

	for {
		if err := filepath.WalkDir(s.workDir, s.processDir); err != nil {
			log.Printf("failed to proceed dir: %v", err)
		}

		time.Sleep(3 * time.Second)
	}
}

// processDir ...
func (s *Scanner) processDir(path string, dir fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if !dir.IsDir() || path == s.workDir {
		return nil
	}

	info, _ := dir.Info()
	modTime := info.ModTime().Unix()

	scanTime, ok := s.scanned[dir.Name()]
	if ok && scanTime == modTime {
		return nil
	}

	s.scanned[dir.Name()] = modTime

	go func() {
		time.Sleep(DirectoryProcessDuration)

		images, _ := filepath.Glob(filepath.Join(path, "*.png"))
		if err = s.uploadTemplate(dir.Name(), email.NewEmailTemplate(s.workDir, images).Render()); err != nil {
			fmt.Printf("Failed to upload %s: %v.\n", dir.Name(), err)
		}

		fmt.Printf("Template %s successfully uploaded.\n", dir.Name())
		_, _ = fmt.Fprintf(s.scannedFile, "%s%s%d\n", dir.Name(), ScannedFileDelimiter, modTime)
	}()

	return nil
}

// uploadTemplate ...
func (s *Scanner) uploadTemplate(name, html string) error {
	values := url.Values{}
	values.Set("api_key", s.apiKey)
	values.Set("name", name)
	values.Set("html", html)

	res, err := http.PostForm("https://a.klaviyo.com/api/v1/email-templates", values)
	if err != nil {
		return err
	} else if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}

	return nil
}
