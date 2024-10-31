package rusti

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/go-resty/resty/v2"

	"github.com/idelchi/godyl/pkg/download"
	"github.com/idelchi/godyl/pkg/file"
)

type Binary struct {
	File file.File
	Dir  file.Folder
	Env  Env
}

var mu sync.Mutex

func New() (binary Binary, err error) {
	mu.Lock()
	defer mu.Unlock()

	dir := file.NewFolder(".rusti")
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

	version, err := binary.Latest()
	if err != nil {
		return binary, err
	}

	target, err := binary.MatchTarget(version)
	if err != nil {
		return binary, err
	}

	err = binary.Download(target)
	if err != nil {
		return binary, err
	}

	binary.Env.Default(binary.Dir.Path())

	return binary, nil
}

func (b *Binary) Find(paths ...string) (file.File, error) {
	binary, err := exec.LookPath("rustc")
	if err != nil {
		for _, path := range paths {
			file := file.NewFile(path, "rust", "bin", "rustc")

			if file.Exists() {
				return file, nil
			}
		}

		return file.File(""), fmt.Errorf("rustc binary not found: %w", err)
	}

	return file.NewFile(binary), nil
}

func (b *Binary) Download(target string) error {
	url := fmt.Sprintf("https://static.rust-lang.org/dist/%s", target)

	downloader := download.New()

	destination, err := downloader.Download(url, b.Dir.Path())
	if err != nil {
		return fmt.Errorf("downloading %q: %w", url, err)
	}

	// Extract the tar.gz file
	cmd := exec.Command("tar", "xzf", destination.String(), "-C", b.Dir.Path())
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("extracting archive: %w", err)
	}

	// The Rust binary will be in the 'rust-[version]/bin' directory
	extractedDir := strings.TrimSuffix(filepath.Base(destination.String()), ".tar.gz")
	b.File = file.NewFile(b.Dir.Path(), extractedDir, "bin", "rustc")

	return nil
}

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

func (b Binary) Latest() (string, error) {
	client := resty.New()
	resp, err := client.R().Get("https://static.rust-lang.org/dist/channel-rust-stable.toml")
	if err != nil {
		return "", err
	}

	// Parse the TOML content to find the version
	lines := strings.Split(string(resp.Body()), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "version = ") {
			return strings.Trim(strings.TrimPrefix(line, "version = "), "\""), nil
		}
	}

	return "", fmt.Errorf("no version found")
}

func (b Binary) MatchTarget(version string) (string, error) {
	os := runtime.GOOS
	arch := runtime.GOARCH

	switch os {
	case "windows":
		arch = "pc-windows-msvc"
	case "darwin":
		os = "apple-darwin"
	case "linux":
		os = "unknown-linux-gnu"
	default:
		return "", fmt.Errorf("unsupported OS: %s", os)
	}

	return fmt.Sprintf("rust-%s-%s-%s.tar.gz", version, arch, os), nil
}
