package file

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func GetAbsolutePwd() string {
	path, _ := os.Getwd()
	return path
}

func GetExecDir() string {
	execFileRelativePath, _ := exec.LookPath(os.Args[0])
	execFileAbsPath, _ := filepath.Abs(execFileRelativePath)
	execFileDirPath := filepath.Dir(execFileAbsPath)
	return execFileDirPath
}

func GetFileByRelative(file string) string {
	dir := GetExecDir()
	return path.Join(dir, file)
}
