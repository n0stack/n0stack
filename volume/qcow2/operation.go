package qcow2

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/n0stack/n0core/lib"
	"github.com/n0stack/n0core/notification"

	"code.cloudfoundry.org/bytefmt"
	pnotification "github.com/n0stack/go.proto/notification/v0"
	qcow2 "github.com/n0stack/go.proto/qcow2/v0"
	uuid "github.com/satori/go.uuid"
)

const fileName = "disk.qcow2"

type volume struct {
	qcow2.Status
	qcow2.Spec

	id      uuid.UUID
	workDir string
}

// qemu-img info $url
// TODO: なんかぐちゃぐちゃなので処理を分けたい
func (v *volume) getQcow2() *pnotification.Notification {
	var err error
	v.workDir, err = lib.GetWorkDir(modelType, v.id)
	if err != nil {
		return notification.MakeNotification("getQcow2.GetWorkDir", false, fmt.Sprintf("error message '%s', when creating work directory, '%s'", v.workDir, err.Error()))
	}

	u := &url.URL{
		Scheme: "file",
		Path:   filepath.Join(v.workDir, fileName), // もう少し汎用的にファイルを取得したい
	}

	v.Url = u.String()
	if _, err := os.Stat(u.Path); err != nil {
		return notification.MakeNotification("getQcow2", true, fmt.Sprintf("Not exists disk image: error message '%s'", err.Error()))
	}

	a := []string{"qemu-img", "info", "--output=json", v.Url}
	cmd := exec.Command(a[0], a[1:]...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return notification.MakeNotification("getQcow2.getImageInfo", false, fmt.Sprintf("error message '%s', std '%s'", err.Error(), out))
	}

	type imageInfo struct {
		Size uint64 `json:"virtual-size"`
		// ActualSize uint64 `json:"actual-size"`
		// Format     string `json:"format"`
	}
	info := &imageInfo{}
	if err = json.Unmarshal(out, &info); err != nil {
		return notification.MakeNotification("getQcow2.UnmarshalJson", false, fmt.Sprintf("error message '%s'", err.Error()))
	}

	v.Bytes = info.Size

	return notification.MakeNotification("getQcow2", true, "Already exists disk image")
}

// qemu-img create -f qcow2 $image $size
func (v *volume) createImage(size uint64) *pnotification.Notification {
	a := []string{"qemu-img", "create", "-f", "qcow2", v.Url, bytefmt.ByteSize(size)}
	cmd := exec.Command(a[0], a[1:]...)
	if o, err := cmd.CombinedOutput(); err != nil {
		return notification.MakeNotification("createImage.runProcess", false, fmt.Sprintf("error message '%s', args '%s', output '%s'", err.Error(), a, o))
	}

	return notification.MakeNotification("createImage", true, "")
}

// qemu-img resize $image $size
func (v *volume) resizeImage(size uint64) *pnotification.Notification {
	return notification.MakeNotification("resizeImage", true, "")
}

// rm $image
func (v *volume) deleteImage() *pnotification.Notification {
	u, err := url.Parse(v.Status.Url)
	if err != nil {
		return notification.MakeNotification("deleteImage.parseURL", false, fmt.Sprintf("error message '%s'", err.Error()))
	}

	if err := os.Remove(u.Path); err != nil {
		return notification.MakeNotification("deleteImage", false, fmt.Sprintf("error message '%s'", err.Error()))
	}

	return notification.MakeNotification("deleteImage", true, "")
}
