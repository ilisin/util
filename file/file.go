package file

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

type Path string

func (p Path) Expand() string {
	return os.ExpandEnv(string(p))
}
func (p Path) Add(s string) string {
	return path.Join(p.Expand(), string(s))
}
func Exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
func Move(dst string, src string) error {
	fi, err := os.Stat(src)
	if err != nil {
		return err
	}
	fd, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, fi.Mode())
	if err != nil {
		return err
	}
	defer fd.Close()
	err = Copy(fd, src)
	if err != nil {
		return err
	}
	return os.Remove(src)
}
func Copy(w io.Writer, src string) error {
	fd, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = io.Copy(w, fd)
	return err
}
func Cat(w io.Writer, path string) error {
	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = io.Copy(w, fd)
	return err
}
func Md5(file string) (hash string, err error) {
	h := md5.New()
	fd, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer fd.Close()
	io.Copy(h, fd)
	return fmt.Sprintf("%X", h.Sum(nil)), err
}
func Md5Path(path string) (hash string) {
	h := md5.New()
	h.Write([]byte(path))
	return fmt.Sprintf("%X", h.Sum(nil))
}
func ReadFile(path string) (data []byte, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}