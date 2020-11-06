package goutil

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func GetImportPathForDir(dir string) string {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return ""
	}

	modpath := findGoModPath(dir)
	if modpath != "" {
		modpkg := findGoModPkg(dir)
		if modpkg != "" {
			return modpkg + strings.TrimPrefix(dir, filepath.Dir(modpath))
		}
	}

	return ""
}

func findGoModPath(dir string) string {
	path, err := execGo([]string{"env", "GOMOD"}, dir)
	if err != nil {
		return ""
	}

	return path
}

func findGoModPkg(dir string) string {
	pkg, err := execGo([]string{"list", "-m"}, dir)
	if err != nil {
		return ""
	}

	return pkg
}

func execGo(args []string, dir string) (string, error) {
	cmd := exec.Command("go", args...)
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.WithStack(err)
	}

	out = bytes.TrimRight(out, "\n")
	if len(out) == 0 {
		return "", nil
	}

	return string(out), nil
}
