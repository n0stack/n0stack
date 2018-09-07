// +build ignore

package n0app

import (
	"fmt"
	"io"

	"github.com/n0stack/n0core/application/n0app/operator"
	"github.com/n0stack/n0core/application/n0app/parser"
)

type Runner struct {
	operator operator.Operations
}

func (r Runner) Up(source []byte, stdout, stderr io.Writer) error {
	p, err := parser.Parse(source)
	if err != nil {
		fmt.Fprintf(stderr, "Failed to parse, err:'%s'", err)

		return fmt.Errorf("Failed to parse, err:'%s'", err)
	}

	for i, s := range p.Sentences {
		fmt.Fprintf(stdout, "Step %d/%d : %s\n", i, len(p.Sentences), string(s.Line))

		var c chan operator.OutputLine
		switch string(s.Opecode.Value) {
		case "FROM":
			c, err = r.operator.FROM(s.Operands)
			if err != nil {
				return err
			}

		case "RUN":
			c, err = r.operator.RUN(s.Operands)
			if err != nil {
				return err
			}

		case "COPY":
			c, err = r.operator.COPY(s.Operands)
			if err != nil {
				return err
			}

		case "DAEMON":
			c, err = r.operator.DAEMON(s.Operands)
			if err != nil {
				return err
			}
		}

		for {
			o, ok := <-c
			if !ok {
				break
			}
			if o.Err != nil {
				fmt.Fprintf(stderr, "Failed step, err:'%s'", o.Err.Error())

				return o.Err
			}

			switch o.OutputType {
			case operator.Stdout:
				fmt.Fprintf(stdout, "%s", o.Line)

			case operator.Stderr:
				fmt.Fprintf(stderr, "%s", o.Line)
			}
		}
	}

	return nil
}
