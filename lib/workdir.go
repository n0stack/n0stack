package lib

import (
	"os"
	"path/filepath"

	"github.com/satori/go.uuid"
)

func GetWorkDir(modelType string, id uuid.UUID) (string, error) {
	const basedir = "/var/lib/n0core"

	w := filepath.Join(basedir, modelType)
	w = filepath.Join(w, id.String())
	err := os.MkdirAll(w, os.ModePerm)

	return w, err
}
