// +build small

package qemu_img

import (
	"testing"

	"code.cloudfoundry.org/bytefmt"
)

// TODO: 現状は正常系のみ

func TestCreateImg(t *testing.T) {
	p := "test.qcow2"

	i, err := OpenQemuImg(p)
	if err != nil {
		t.Fatalf("Cannot open '%s': err='%s'", p, err.Error())
	}
	if i.IsExists() {
		t.Errorf("Test environment is invalid, image is already existing: path='%s'", p)
	}

	s, _ := bytefmt.ToBytes("1G")
	if err := i.Create(s); err != nil {
		t.Errorf("Failed to create image: err='%s'", err.Error())
	}
	if !i.IsExists() {
		t.Errorf("Failed to create image: image is not existing yet")
	}

	if err := i.Delete(); err != nil {
		t.Errorf("Failed to delete image: err='%s'", err.Error())
	}
	if i.IsExists() {
		t.Errorf("Failed to delete image: image is existing yet")
	}
}

func TestResize(t *testing.T) {
	p := "test.qcow2"
	src := "1G"
	dst := "2G"

	i, err := OpenQemuImg(p)
	if err != nil {
		t.Fatalf("Cannot open '%s': err='%s'", p, err.Error())
	}
	if i.IsExists() {
		t.Errorf("Test environment is invalid, image is already existing: path='%s'", p)
	}

	s, _ := bytefmt.ToBytes(src)
	if err := i.Create(s); err != nil {
		t.Errorf("Failed to create image: err='%s'", err.Error())
	}
	if bytefmt.ByteSize(i.Info.VirtualSize) != src {
		t.Errorf("Image size is mismatch: want='%s', have='%s'", src, bytefmt.ByteSize(i.Info.VirtualSize))
	}

	s, _ = bytefmt.ToBytes(dst)
	if err := i.Resize(s); err != nil {
		t.Errorf("Failed to resize image: err='%s'", err.Error())
	}
	if bytefmt.ByteSize(i.Info.VirtualSize) != dst {
		t.Errorf("Image size is mismatch: want='%s', have='%s'", dst, bytefmt.ByteSize(i.Info.VirtualSize))
	}

	if err := i.Delete(); err != nil {
		t.Errorf("Failed to delete image: err='%s'", err.Error())
	}
	if i.IsExists() {
		t.Errorf("Failed to delete image: image is existing yet")
	}
}

func TestBackingfile(t *testing.T) {
	p := "test.qcow2"
	pb := "test.diff.qcow2"

	i, err := OpenQemuImg(p)
	if err != nil {
		t.Fatalf("Cannot open '%s': err='%s'", p, err.Error())
	}
	if i.IsExists() {
		t.Errorf("Test environment is invalid, image is already existing: path='%s'", p)
	}

	s, _ := bytefmt.ToBytes("1G")
	if err := i.Create(s); err != nil {
		t.Errorf("Failed to create image: err='%s'", err.Error())
	}

	bi, err := i.CreateBackingFile(pb)
	if !bi.IsExists() {
		t.Errorf("Failed to create backing image: image is not existing yet")
	}
	if bi.Info.BackingFilename != i.Info.Filename {
		t.Errorf("Filename is mismatch: want='%s', have='%s'", i.Info.Filename, bi.Info.BackingFilename)
	}

	if err := bi.Delete(); err != nil {
		t.Errorf("Failed to delete image: err='%s'", err.Error())
	}
	if bi.IsExists() {
		t.Errorf("Failed to delete image: image is existing yet")
	}
	if err := i.Delete(); err != nil {
		t.Errorf("Failed to delete image: err='%s'", err.Error())
	}
	if i.IsExists() {
		t.Errorf("Failed to delete image: image is existing yet")
	}
}
