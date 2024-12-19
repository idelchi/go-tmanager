// Package env provides utilities for managing environment variables in Go.
// It offers functionality for loading, manipulating, and normalizing environment variables
// from various sources, such as the system environment or .env files. Additionally,
// it allows environment variable expansion, filtering, and conversion between different formats.
//
// Key features include:
// - Loading environment variables from the OS or .env files.
// - Adding, merging, and normalizing environment variables.
// - Expanding environment variable values based on the current system environment.
// - Converting environment variables to and from slices of `key=value` strings.
package env
