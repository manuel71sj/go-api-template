package file

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

// SelfPath 컴파일된 실행 파일 절대 경로를 가져옵니다.
func SelfPath() string {
	pt, _ := filepath.Abs(os.Args[0])

	return pt
}

// RealPath 빌드된 실행 파일을 기반으로 절대 파일 경로 가져오기
func RealPath(fp string) (string, error) {
	if path.IsAbs(fp) {
		return fp, nil
	}

	wd, err := os.Getwd()

	return path.Join(wd, fp), err
}

// SelfDir 컴파일된 실행 파일 디렉터리를 가져옵니다.
func SelfDir() string {
	return filepath.Dir(SelfPath())
}

// Basename 파일 경로 기본 이름 가져오기
func Basename(fp string) string {
	return path.Base(fp)
}

// Dir 파일 경로 디렉터리 이름 가져오기
func Dir(fp string) string {
	return path.Dir(fp)
}

func InsureDir(fp string) error {
	if IsExist(fp) {
		return nil
	}
	return os.MkdirAll(fp, os.ModePerm)
}

// EnsureDir 존재하지 않는 경우 디렉터리 만들기
func EnsureDir(fp string) error {
	return os.MkdirAll(fp, os.ModePerm)
}

// EnsureDirRW the dataDir and make sure it's re-able
func EnsureDirRW(dataDir string) error {
	err := EnsureDir(dataDir)
	if err != nil {
		return err
	}

	checkFile := fmt.Sprintf("%s/rw.%d", dataDir, time.Now().UnixNano())
	fd, err := Create(checkFile)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("open %s: rw permission denied", dataDir)
		}

		return err
	}

	if err := Close(fd); err != nil {
		return fmt.Errorf("close error: %s", err)
	}

	if err := Remove(checkFile); err != nil {
		return fmt.Errorf("remove error: %s", err)
	}

	return nil
}

// Create one file
func Create(name string) (*os.File, error) {
	return os.Create(name)
}

// Remove one file
func Remove(name string) error {
	return os.Remove(name)
}

// Close fd
func Close(fd *os.File) error {
	return fd.Close()
}

func Ext(fp string) string {
	return path.Ext(fp)
}

// Rename file name
func Rename(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

// Unlink delete file
func Unlink(fp string) error {
	return os.Remove(fp)
}

// IsFile 경로가 파일인지 확인하며, 디렉터리이거나 존재하지 않으면 false를 반환합니다.
func IsFile(fp string) bool {
	f, e := os.Stat(fp)
	if e != nil {
		return false
	}

	return !f.IsDir()
}

// IsExist 파일 또는 디렉터리의 존재 여부를 확인합니다. 파일이나 디렉터리가 존재하지 않으면 false를 반환합니다.
func IsExist(fp string) bool {
	_, err := os.Stat(fp)
	return err == nil || os.IsExist(err)
}

// SearchFile 경로에서 파일을 검색합니다. 이것은 종종 / etc ~/의 config file 검색에 사용됩니다.
func SearchFile(filename string, paths ...string) (fullPath string, err error) {
	for _, pt := range paths {
		if fullPath = filepath.Join(pt, filename); IsExist(fullPath) {
			return
		}
	}

	err = fmt.Errorf("%s not found in paths", fullPath)
	return
}

// FileMTime 파일 수정 시간 가져오기
func FileMTime(fp string) (int64, error) {
	f, e := os.Stat(fp)
	if e != nil {
		return 0, e
	}

	return f.ModTime().Unix(), nil
}

// FileSize 파일 크기를 바이트로 가져옵니다.
func FileSize(fp string) (int64, error) {
	f, e := os.Stat(fp)
	if e != nil {
		return 0, e
	}

	return f.Size(), nil
}

// DirsUnder dirPath 아래의 디렉터리 나열
func DirsUnder(dirPath string) ([]string, error) {
	if !IsExist(dirPath) {
		return []string{}, nil
	}

	fs, err := os.ReadDir(dirPath)
	if err != nil {
		return []string{}, err
	}

	sz := len(fs)
	if sz == 0 {
		return []string{}, nil
	}

	ret := make([]string, 0, sz)
	for i := 0; i < sz; i++ {
		if fs[i].IsDir() {
			name := fs[i].Name()
			if name != "." && name != ".." {
				ret = append(ret, name)
			}
		}
	}

	return ret, nil
}

// FilesUnder 디렉터리 경로 아래에 파일 나열
func FilesUnder(dirPath string) ([]string, error) {
	if !IsExist(dirPath) {
		return []string{}, nil
	}

	fs, err := os.ReadDir(dirPath)
	if err != nil {
		return []string{}, err
	}

	sz := len(fs)
	if sz == 0 {
		return []string{}, nil
	}

	ret := make([]string, 0, sz)
	for i := 0; i < sz; i++ {
		if !fs[i].IsDir() {
			ret = append(ret, fs[i].Name())
		}
	}

	return ret, nil
}

func MustOpenLogFile(fp string) *os.File {
	if strings.Contains(fp, "/") {
		dir := Dir(fp)
		err := EnsureDir(dir)
		if err != nil {
			log.Fatalf("mkdir -p %s occur error %v", dir, err)
		}
	}

	f, err := os.OpenFile(fp, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("open %s occur error %v", fp, err)
	}

	return f
}
