package file

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Folder represents a file system directory as a string.
// It provides methods for working with directories, such as creating,
// removing, expanding paths, and checking existence.
type Folder string

// NewFolder creates a new Folder object from the provided path segments by joining them.
func NewFolder(paths ...string) Folder {
	return Folder(filepath.Join(paths...))
}

// IsSet checks whether the Folder has been set to a non-empty value.
func (f Folder) IsSet() bool {
	return f != ""
}

// Set assigns a new path to the Folder.
func (f *Folder) Set(path string) {
	*f = Folder(path)
}

// IsParentOf determines if the Folder is a parent directory of the given 'other' Folder.
func (f *Folder) IsParentOf(other Folder) bool {
	return strings.HasPrefix(other.Path(), f.Path())
}

// Expand expands a Folder path that begins with "~" to the user's home directory.
func (f *Folder) Expand() error {
	if strings.HasPrefix(f.Path(), "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("getting user home directory: %w", err)
		}

		f.Set(filepath.Join(homeDir, f.Path()[2:]))
	}

	return nil
}

// String returns the Folder as a string.
func (f Folder) String() string {
	return string(f)
}

// Path returns the Folder path as a string.
func (f Folder) Path() string {
	return f.String()
}

// Exists checks if the Folder exists in the file system.
func (f Folder) Exists() bool {
	_, err := os.Stat(f.Path())

	return err == nil
}

// Create creates the Folder and all necessary parent directories with 0755 permissions.
func (f Folder) Create() error {
	const perm = 0o755

	if err := os.MkdirAll(f.Path(), perm); err != nil {
		return fmt.Errorf("creating directory %s: %w", f.Path(), err)
	}

	return nil
}

// Name returns the base name (last element) of the Folder's path.
func (f Folder) Name() string {
	return filepath.Base(f.Path())
}

// CreateRandomInDir creates a new random directory inside the given directory
// and assigns the generated path to the Folder.
func (f *Folder) CreateRandomInDir(dir string) error {
	name, err := os.MkdirTemp(dir, "godyl-*")
	if err != nil {
		return fmt.Errorf("creating temporary directory in %s: %w", dir, err)
	}

	f.Set(name)

	return nil
}

// CreateRandomInTempDir creates a new random directory inside the system's temporary directory
// and assigns the generated path to the Folder.
func (f *Folder) CreateRandomInTempDir() error {
	name, err := os.MkdirTemp("", "godyl-*")
	if err != nil {
		return fmt.Errorf("creating temporary directory: %w", err)
	}

	f.Set(name)

	return nil
}

// CreateInTempDir creates a directory inside the system's temporary directory
// using the Folder's name and assigns the path to the Folder.
func (f *Folder) CreateInTempDir() error {
	name := filepath.Join(os.TempDir(), f.Name())

	const perm = 0o755

	if err := os.Mkdir(name, perm); err != nil {
		return fmt.Errorf("creating directory in temporary directory: %w", err)
	}

	f.Set(name)

	return nil
}

// Remove deletes the Folder and all of its contents.
func (f Folder) Remove() error {
	if err := os.RemoveAll(f.Path()); err != nil {
		return fmt.Errorf("removing directory %s: %w", f.Path(), err)
	}

	return nil
}

// CriteriaFunc defines a function type for filtering files during search operations.
type CriteriaFunc func(File) (bool, error)

// ErrNotFound is returned when a file is not found during a search operation.
var ErrNotFound = errors.New("file not found")

// FindFile searches for a file in the Folder that matches the provided criteria.
// It returns the first file found or an error if none are found.
func (f Folder) FindFile(criteria ...CriteriaFunc) (File, error) {
	var file File
	var found bool

	err := filepath.Walk(f.Path(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if File(path).IsDir() {
			return nil // Skip directories
		}

		path, err = filepath.Rel(f.Path(), path)
		if err != nil {
			return err
		}

		// Check if the file matches all criteria
		for _, criterion := range criteria {
			matches, err := criterion(File(path))
			if err != nil {
				return err
			}

			if !matches {
				return nil // Skip this file if it doesn't match all criteria
			}
		}

		file = NewFile(f.Path(), path)

		// If we've reached here, the file matches all criteria
		found = true

		return filepath.SkipAll // Stop the walk, we've found a match
	})
	if err != nil {
		return File(""), fmt.Errorf("error walking folder %s: %w", f.Path(), err)
	}

	if !found {
		return File(""), fmt.Errorf("%w: no file found matching all criteria in folder %s", ErrNotFound, f.Path())
	}

	return file, nil
}

// ListFolders returns a slice of Folders representing all subdirectories
// within the current Folder.
func (f Folder) ListFolders() ([]Folder, error) {
	var folders []Folder

	entries, err := os.ReadDir(f.Path())
	if err != nil {
		return nil, fmt.Errorf("reading directory %s: %w", f.Path(), err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			subFolder := NewFolder(f.Path(), entry.Name())
			folders = append(folders, subFolder)
		}
	}

	return folders, nil
}

// ListFiles returns a slice of Files representing all files
// within the current Folder.
func (f Folder) ListFiles() (Files, error) {
	var files Files

	entries, err := os.ReadDir(f.Path())
	if err != nil {
		return nil, fmt.Errorf("reading directory %s: %w", f.Path(), err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			file := NewFile(f.Path(), entry.Name())
			files = append(files, file)
		}
	}

	return files, nil
}
