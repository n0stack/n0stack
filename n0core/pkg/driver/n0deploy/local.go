package n0deploy

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func NewLocalParser() *Parser {
	return &Parser{
		NewCopyInstruction: NewLocalCopyInstruction,
		NewRunInstruction:  NewLocalRunInstruction,
	}
}

type LocalRunInstruction struct {
	command string
}

func NewLocalRunInstruction(line string) (Instruction, error) {
	return &LocalRunInstruction{
		command: strings.TrimPrefix(line, "RUN "),
	}, nil
}

func (i LocalRunInstruction) Do(ctx context.Context, out io.Writer) error {
	cmd := exec.CommandContext(ctx, "sh", "-c", i.command)
	cmd.Stdout = out
	cmd.Stderr = out

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (i LocalRunInstruction) String() string {
	return fmt.Sprintf("RUN %s", i.command)
}

type LocalCopyInstruction struct {
	src string
	dst string
}

func NewLocalCopyInstruction(src, dst string) (Instruction, error) {
	return &LocalCopyInstruction{
		src: src,
		dst: dst,
	}, nil
}

func (i LocalCopyInstruction) Do(ctx context.Context, out io.Writer) error {
	cmd := exec.CommandContext(ctx, "cp", i.src, i.dst)
	cmd.Stdout = out
	cmd.Stderr = out

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (i LocalCopyInstruction) String() string {
	return fmt.Sprintf("COPY %s %s", i.src, i.dst)
}
