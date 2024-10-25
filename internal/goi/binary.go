package goi

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/go-resty/resty/v2"

	"github.com/idelchi/godyl/pkg/download"
	"github.com/idelchi/godyl/pkg/file"
)

// Binary represents a Go binary, including its associated file, directory, and environment variables.
type Binary struct {
	File file.File   // File holds the file information for the Go binary.
	Dir  file.Folder // Dir refers to the directory where the binary is stored.
	Env  Env         // Env contains the environment variables for running the binary.

	noVerifySSL bool
}

var mu sync.Mutex

// New creates a new Binary instance, setting up the directory, downloading the latest release if necessary,
// and initializing environment variables. It ensures thread-safe execution by using a mutex lock.
func New(noVerifySSL bool) (binary Binary, err error) {
	mu.Lock()
	defer mu.Unlock()

	binary.noVerifySSL = noVerifySSL

	dir := file.NewFolder(".godyl-go")
	if err := dir.CreateInTempDir(); err != nil && !errors.Is(err, os.ErrExist) {
		return binary, fmt.Errorf("creating temp dir: %w", err)
	}

	if file, err := binary.Find(dir.Path()); err == nil {
		binary.File = file
		if dir.IsParentOf(file.Dir()) {
			binary.Dir = dir
			binary.Env.Default(binary.Dir.Path())
		} else {
			binary.Env = Env{}
			binary.Dir = file.Dir()
		}
		return binary, nil
	} else {
		binary.Dir = dir
	}

	release, err := binary.Latest()
	if err != nil {
		return binary, err
	}

	targets := Targets{}
	for _, file := range release.Files {
		if file.IsArchive() {
			targets.Files = append(targets.Files, file)
		}
	}

	path, err := targets.Match()
	if err != nil {
		return binary, err
	}

	err = binary.Download(path[0].Asset.Name)
	if err != nil {
		return binary, err
	}

	binary.Env.Default(binary.Dir.Path())

	return binary, nil
}

// Find searches for the Go binary in the given paths or system path, returning the file if found.
func (b *Binary) Find(paths ...string) (file.File, error) {
	binary, err := exec.LookPath("go")
	if err != nil {
		for _, path := range paths {
			file := file.NewFile(path, "go", "bin", "go")

			if file.Exists() {
				return file, nil
			}
		}

		return file.File(""), fmt.Errorf("go binary not found: %w", err)
	}

	return file.NewFile(binary), nil
}

// Download downloads the Go binary from the provided path and saves it to the directory.
// It returns an error if the download or file validation fails.
func (b *Binary) Download(path string) error {
	url := fmt.Sprintf("https://go.dev/dl/%s", path)

	downloader := download.New()
	downloader.InsecureSkipVerify = b.noVerifySSL

	destination, err := downloader.Download(url, b.Dir.Path())
	if err != nil {
		return fmt.Errorf("downloading %q: %w", url, err)
	}

	b.File = file.NewFile(destination.String(), "go", "bin", "go")

	return nil
}

// CleanUp removes the temporary directory associated with the binary.
// It returns an error if the directory removal fails.
func (b *Binary) CleanUp() error {
	if !b.Dir.IsSet() {
		return nil
	}

	if err := b.Dir.Remove(); err != nil {
		return fmt.Errorf("cleaning up: %w", err)
	}

	b.Dir = ""

	return nil
}

// Latest fetches the latest Go release information from the official Go download page.
// It returns the most recent release or an error if the process fails.
func (b Binary) Latest() (Release, error) {
	client := resty.New()
	resp, err := client.R().Get("https://go.dev/dl/?mode=json")
	if err != nil {
		return Release{}, err
	}

	var releases []Release

	if err := json.Unmarshal(resp.Body(), &releases); err != nil {
		return Release{}, err
	}

	if len(releases) > 0 {
		return releases[0], nil
	}

	return Release{}, fmt.Errorf("no versions found")
}
