package img

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

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
