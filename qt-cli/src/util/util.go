// Copyright (C) 2024 The Qt Company Ltd.
// SPDX-License-Identifier: LicenseRef-Qt-Commercial OR LGPL-3.0-only

package util

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"syscall"
)

type StringAnyMap map[string]interface{}

func MergeMap(base StringAnyMap, other StringAnyMap) StringAnyMap {
	all := StringAnyMap{}

	for k, v := range base {
		all[k] = v
	}

	for k, v := range other {
		all[k] = v
	}

	return all
}

func ReadAllFromFS(targetFS fs.FS, path string) ([]byte, error) {
	stat, err := fs.Stat(targetFS, path)
	if err != nil {
		return []byte{}, fmt.Errorf(
			Msg("cannot read file info, given = '%v'"), path)
	}

	if !stat.Mode().IsRegular() {
		return []byte{}, fmt.Errorf(
			Msg("cannot read non-regular file, given = '%v'"), path)
	}

	file, err := targetFS.Open(path)
	if err != nil {
		return []byte{}, err
	}

	defer file.Close()
	return io.ReadAll(file)
}

func WriteAll(data []byte, destPath string) (int, error) {
	dir := path.Dir(destPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return 0, err
	}

	destFile, err := os.Create(destPath)
	if err != nil {
		return 0, err
	}

	defer destFile.Close()
	return destFile.Write(data)
}

func EntryExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func EntryExistsFS(targetFS fs.FS, path string) bool {
	_, err := fs.Stat(targetFS, path)
	return !os.IsNotExist(err)
}

func ToBool(value interface{}, defaultValue bool) bool {
	switch c := value.(type) {
	case bool:
		return c

	case string:
		{
			s := strings.TrimSpace(strings.ToLower(c))
			return s == "true" || s == "yes"
		}

	case int:
		return c != 0

	case nil:
		return false

	default:
		return defaultValue
	}
}

func ToFloat64(value interface{}, defaultValue float64) float64 {
	switch c := value.(type) {
	case string:
		v, err := strconv.ParseFloat(c, 64)
		if err != nil {
			return defaultValue
		}

		return v

	case int:
		return float64(c)

	case nil:
		return 0.0

	default:
		return defaultValue
	}
}

func Msg(s string) string {
	return s
}

func IsValidDirName(name string) bool {
	tempDir := os.TempDir()
	testPath := path.Join(tempDir, name)

	err := os.Mkdir(testPath, 0755)
	if err != nil {
		return false
	}

	os.Remove(testPath)
	return true
}

func IsValidFileName(name string) bool {
	tempDir := os.TempDir()
	testPath := path.Join(tempDir, name)

	file, err := os.Create(testPath)
	if err != nil {
		return false
	}

	file.Close()
	os.Remove(testPath)
	return true
}

func IsWindowsReservedName(name string) bool {
	name = strings.ToUpper(name)
	if name == "CON" || name == "PRN" || name == "AUX" || name == "NUL" {
		return true
	}
	match, _ := regexp.MatchString(`^(COM[1-9]|LPT[1-9])$`, name)
	return match
}

func IsValidFolderName(name string) bool {
	name = strings.TrimSpace(filepath.Base(name))
	if name == "" {
		return false
	}

	if runtime.GOOS == "windows" {
		invalidChars := regexp.MustCompile(`[<>:"/\\|?*]`)
		if invalidChars.MatchString(name) {
			return false
		}

		if IsWindowsReservedName(name) {
			return false
		}
	} else {
		if strings.ContainsRune(name, '/') || strings.ContainsRune(name, 0) {
			return false
		}
	}

	return true
}

func IsDriveExists(path string) bool {
	if runtime.GOOS != "windows" {
		return true
	}

	if len(path) < 2 || path[1] != ':' {
		return true
	}

	driveRoot := path[:2] + `\`
	info, err := os.Stat(driveRoot)
	return err == nil && info.IsDir()
}

func IsValidFullPath(path string) bool {
	// 'C:\Users' => valid? true
	// 'D:\NotExists' => valid? false
	// 'COM1' => valid? false
	// 'COM10' => valid? true
	// 'valid_folder' => valid? true
	// 'invalid:name' => valid? false
	// '/tmp/test' => valid? true
	// 'my/folder' => valid? false
	// '.hidden' => valid? true
	// ' ' => valid? false

	return IsDriveExists(path) && IsValidFolderName(path)
}

func IsAbsPath(path string) bool {
	return filepath.IsAbs(path)
}

func IsValidProjectName(name string) bool {
	if len(name) == 0 {
		return false
	}

	validNameRegex := regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_-]+$")
	if !validNameRegex.MatchString(name) {
		return false
	}

	if runtime.GOOS == "windows" {
		if strings.ContainsAny(name, `\/:*?"<>|`) {
			return false
		}
	} else {
		if strings.ContainsAny(name, "/") {
			return false
		}
	}

	return true
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil || os.IsNotExist(err) {
		return false
	}

	return info.IsDir()
}

func SendSigTermOrKill(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	if runtime.GOOS == "windows" {
		return process.Kill()
	} else {
		return process.Signal(syscall.SIGTERM)
	}
}
