package qcow2

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"code.cloudfoundry.org/bytefmt"
	n0stack "github.com/n0stack/go-proto"
	"github.com/n0stack/go-proto/device/volume"
	"github.com/n0stack/go-proto/resource/storage"
	"github.com/n0stack/n0core/lib"
	uuid "github.com/satori/go.uuid"
)

const fileName = "disk.qcow2"

type qcow2 struct {
	volume.Status
	volume.Spec

	id      uuid.UUID
	workDir string
}

// qemu-img info $url
func (q *qcow2) getStatus(model *n0stack.Model) *n0stack.Notification {
	q.Model = model

	var err error
	q.id, err = uuid.FromBytes(model.Id)
	if err != nil {
		return lib.MakeNotification("getQcow2.validateUUID", false, fmt.Sprintf("error message '%s'", err.Error()))
	}

	q.workDir, err = lib.GetWorkDir(modelType, q.id)
	if err != nil {
		return lib.MakeNotification("getQcow2.GetWorkDir", false, fmt.Sprintf("error message '%s', when creating work directory, '%s'", q.workDir, err.Error()))
	}

	u := &url.URL{
		Scheme: "file",
		Path:   filepath.Join(q.workDir, fileName), // もう少し汎用的にファイルを取得したい
	}

	q.Url = u.String()
	if _, err := os.Stat(u.Path); err != nil {
		return lib.MakeNotification("getQcow2", true, fmt.Sprintf("Not exists disk image: error message '%s'", err.Error()))
	}

	a := []string{"qemu-img", "info", "--output=json", q.Url}
	cmd := exec.Command(a[0], a[1:]...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return lib.MakeNotification("getQcow2.getImageInfo", false, fmt.Sprintf("error message '%s', std '%s'", err.Error(), out))
	}

	type imageInfo struct {
		Size       uint64 `json:"virtual-size"`
		ActualSize uint64 `json:"actual-size"`
		Format     string `json:"format"`
	}
	info := &imageInfo{}
	if err = json.Unmarshal(out, &info); err != nil {
		return lib.MakeNotification("getQcow2.UnmarshalJson", false, fmt.Sprintf("error message '%s'", err.Error()))
	}

	q.Storage = &storage.Spec{
		Bytes: info.Size,
	}

	return lib.MakeNotification("getQcow2", true, "Already exists disk image")
}

// qemu-img create -f qcow2 $image $size
func (q *qcow2) createImage(size uint64) *n0stack.Notification {
	a := []string{"qemu-img", "create", "-f", "qcow2", q.Url, bytefmt.ByteSize(size)}
	cmd := exec.Command(a[0], a[1:]...)
	if o, err := cmd.CombinedOutput(); err != nil {
		return lib.MakeNotification("createImage.runProcess", false, fmt.Sprintf("error message '%s', args '%s', output '%s'", err.Error(), a, o))
	}

	return lib.MakeNotification("createImage", true, "")
}

// qemu-img resize $image $size
func (q *qcow2) resizeImage(size uint64) *n0stack.Notification {
	return lib.MakeNotification("resizeImage", true, "")
}

// rm $image
func (q *qcow2) deleteImage() *n0stack.Notification {
	u, err := url.Parse(q.Status.Url)
	if err != nil {
		return lib.MakeNotification("deleteImage.parseURL", false, fmt.Sprintf("error message '%s'", err.Error()))
	}

	if err := os.Remove(u.Path); err != nil {
		return lib.MakeNotification("deleteImage", false, fmt.Sprintf("error message '%s'", err.Error()))
	}

	return lib.MakeNotification("deleteImage", true, "")
}
