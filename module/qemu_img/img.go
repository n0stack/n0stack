package qemu_img

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"code.cloudfoundry.org/bytefmt"
)

// ImgInfo is response of `qemu-img info`.
type ImgInfo struct {
	VirtualSize     uint64                 `json:"virtual-size"`
	Filename        string                 `json:"filename"`
	ClusterSize     uint64                 `json:"cluster-size"`
	Format          string                 `json:"format"`
	ActuralSize     uint64                 `json:"actual-size"`
	FormatSpecific  map[string]interface{} `json:"format-specific"`
	BackingFilename string                 `json:"backing-filename"`
	DirtyFlag       bool                   `json:"dirty-flag"`
}

// とりあえず入れてるだけ

type QemuImg struct {
	Info *ImgInfo

	path string
}

func OpenQemuImg(path string) (*QemuImg, error) {
	p, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("") // TODO
	}
	// check permission

	q := &QemuImg{
		path: p,
	}
	q.updateInfo()

	return q, nil
}

func (q *QemuImg) Create(bytes uint64) error {
	if q.IsExists() {
		return fmt.Errorf("Already exists") // TODO
	}

	args := []string{"qemu-img", "create", "-f", "qcow2", q.path, bytefmt.ByteSize(bytes)}
	cmd := exec.Command(args[0], args[1:]...)
	if o, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Failed to create image: err='%s', args='%v', output='%s'", err.Error(), args, o)
	}

	q.updateInfo()

	return nil
}

func (q QemuImg) CreateBackingFile(path string) (*QemuImg, error) {
	if q.Info.Format != "qcow2" {
		return nil, fmt.Errorf("Cannot create backing file because base file is not qcow2: format='%s'", q.Info.Format)
	}

	args := []string{"qemu-img", "create", "-f", "qcow2", "-b", q.path, path}
	cmd := exec.Command(args[0], args[1:]...)
	if o, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("Failed to create image: err='%s', args='%v', output='%s'", err.Error(), args, o)
	}

	nq, err := OpenQemuImg(path)
	if err != nil {
		return nil, err
	}

	return nq, nil
}

func (q *QemuImg) Download(src *url.URL) error {
	if q.IsExists() {
		return fmt.Errorf("Already exists") // TODO
	}

	res, err := http.Get(src.String())
	if err != nil {
		return fmt.Errorf("Failed to get file from http: url='%s', err='%s'", src.String(), err.Error())
	}
	defer res.Body.Close()

	f, err := os.Create(q.path)
	if err != nil {
		return fmt.Errorf("Failed to open file: path='%s', err='%s'", q.path, err.Error())
	}
	defer f.Close()

	io.Copy(f, res.Body)

	q.updateInfo()

	return nil
}

func (q *QemuImg) Resize(bytes uint64) error {
	args := []string{"qemu-img", "resize", q.path, bytefmt.ByteSize(bytes)}
	cmd := exec.Command(args[0], args[1:]...)
	if o, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("Failed to create image: err='%s', args='%v', output='%s'", err.Error(), args, o)
	}

	q.updateInfo()

	return nil
}

func (q *QemuImg) Delete() error {
	if !q.IsExists() {
		return fmt.Errorf("Already deleted") // nilを返したほうがいいかも？
	}

	if err := os.Remove(q.path); err != nil {
		return fmt.Errorf("Failed to delete image: err='%s'", err.Error())
	}

	q.updateInfo()

	return nil
}

func (q QemuImg) IsExists() bool {
	if q.Info == nil {
		return false
	}

	return true
}

// ファイルが存在しない場合は正常
// -> 空のイメージを作成するかダウンロードをする
// イメージではないファイルが存在する場合にはエラー
func (q *QemuImg) updateInfo() error {
	if _, err := os.Stat(q.path); err != nil {
		q.Info = nil

		return nil
	}

	args := []string{"qemu-img", "info", "--output=json", q.path}
	cmd := exec.Command(args[0], args[1:]...)
	o, err := cmd.CombinedOutput()
	if err != nil {
		// TODO: 書いたほうがいいか？
		// q.Info = nil

		return fmt.Errorf("Failed to get info from image: err='%s', args='%v', output='%s'", err.Error(), args, o)
	}

	if q.Info == nil {
		q.Info = &ImgInfo{}
	}
	if err := json.Unmarshal(o, q.Info); err != nil {
		return fmt.Errorf("Failed to parse qemu-img info by json: err='%s'", err.Error())
	}

	return nil
}
