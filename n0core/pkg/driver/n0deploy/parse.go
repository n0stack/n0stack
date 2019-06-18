package n0deploy

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// ValidateCopy return source and destination to copy file
func ValidateCopy(line string) (string, string, error) {
	r := csv.NewReader(strings.NewReader(line))
	r.Comma = ' '
	args, err := r.Read()
	if err != nil {
		return "", "", errors.Wrap(err, "Failed to parse as csv, csv package is used for splitting a string at Space except inside quotation marks") // https://stackoverflow.com/questions/47489745/splitting-a-string-at-space-except-inside-quotation-marks-go
	}

	if len(args) != 3 {
		return "", "", fmt.Errorf("COPY instruction takes 2 arguments")
	}
	// srcがpwdから出る場合はエラーを返す
	// if args[1]

	if !filepath.IsAbs(args[2]) {
		return "", "", fmt.Errorf("Set an absolute path for COPY destination")
	}

	return args[1], args[2], nil
}

type N0deploy struct {
	Bootstrap []Instruction
	Deploy    []Instruction
}

type BlockType int

const (
	BlockType_BOOTSTRAP BlockType = iota
	BlockType_DEPLOY
)

type Instruction interface {
	Do(ctx context.Context, out io.Writer) error
	String() string
}

type Parser struct {
	NewRunInstruction  func(command string) (Instruction, error)
	NewCopyInstruction func(src, dst string) (Instruction, error)
}

func (p Parser) Parse(src string) (*N0deploy, error) {
	src = strings.Replace(src, "\\\n", "", -1)
	lines := strings.Split(src, "\n")

	ins := map[BlockType][]Instruction{
		BlockType_BOOTSTRAP: make([]Instruction, 0, len(lines)),
		BlockType_DEPLOY:    make([]Instruction, 0, len(lines)),
	}
	block := BlockType_BOOTSTRAP

	for _, l := range lines {
		l = strings.TrimSpace(l)
		l = strings.Trim(l, "\t")

		if len(l) == 0 {
			continue
		}

		switch {
		case strings.HasPrefix(l, "RUN"):
			i, err := p.NewRunInstruction(strings.TrimPrefix(l, "RUN "))
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to parse '%s'", l)
			}

			ins[block] = append(ins[block], i)

		case strings.HasPrefix(l, "COPY"):
			src, dst, err := ValidateCopy(l)
			if err != nil {
				return nil, err
			}

			i, err := p.NewCopyInstruction(src, dst)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to parse '%s'", l)
			}

			ins[block] = append(ins[block], i)

		case strings.HasPrefix(l, "BOOTSTRAP"):
			if len(ins[BlockType_BOOTSTRAP]) != 0 {
				return nil, errors.New("BOOTSTAP block is duplicated")
			}

			block = BlockType_BOOTSTRAP

		case strings.HasPrefix(l, "DEPLOY"):
			if len(ins[BlockType_DEPLOY]) != 0 {
				return nil, errors.New("DEPLOY block is duplicated")
			}

			block = BlockType_DEPLOY

		default:
			return nil, fmt.Errorf("Failed to parse '%s': the instruction does not exist", l)
		}
	}

	return &N0deploy{
		Bootstrap: ins[BlockType_BOOTSTRAP],
		Deploy:    ins[BlockType_DEPLOY],
	}, nil
}
