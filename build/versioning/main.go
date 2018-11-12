package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli"
)

var version = "undefined"

const VERSION_FILE = "VERSION"

func IsExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func getCurrentVersion(filename string) (uint64, error) {
	var v uint64

	if IsExists(filename) {
		f, err := os.OpenFile(filename, os.O_RDONLY, 0666)
		if err != nil {
			return 0, fmt.Errorf("Failed to open %s for read: err=%s", filename, err.Error())
		}
		defer f.Close()

		raw, err := ioutil.ReadAll(f)
		if err != nil {
			return 0, fmt.Errorf("Failed to read %s: err=%s", filename, err.Error())
		}

		v, err = strconv.ParseUint(string(raw), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("Failed to parse string %s as uint64: err=%s", string(raw), err.Error())
		}
	}

	return v, nil
}

func Print(ctx *cli.Context) error {
	v, err := getCurrentVersion(VERSION_FILE)
	if err != nil {
		return err
	}

	fmt.Println(v)

	return nil
}
func Commit(ctx *cli.Context) error {
	write := ctx.Bool("write")

	v, err := getCurrentVersion(VERSION_FILE)
	if err != nil {
		return err
	}

	v++
	fmt.Println(v)

	if write {
		f, err := os.OpenFile(VERSION_FILE, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return fmt.Errorf("Failed to open %s for write: err=%s", VERSION_FILE, err.Error())
		}
		defer f.Close()

		s := strconv.FormatUint(v, 10)
		if _, err := f.Write([]byte(s)); err != nil {
			return fmt.Errorf("Failed to write %s: err=%s", VERSION_FILE, err.Error())
		}
	}

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "versioning"
	app.Version = version
	// app.Usage = ""

	app.Commands = []cli.Command{
		{
			Name:   "commit",
			Usage:  "versioning commit -write",
			Action: Commit,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "write",
				},
			},
		},
		{
			Name:   "print",
			Usage:  "versioning print -write",
			Action: Print,
		},
	}

	log.SetFlags(log.Lshortfile)
	// log.SetOutput(ioutil.Discard)

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to command: %v\n", err.Error())
		os.Exit(1)
	}
}
