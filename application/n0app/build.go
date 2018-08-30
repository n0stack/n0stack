package n0app

import (
	"fmt"
	"io"

	"github.com/n0stack/n0core/application/n0app/operator"
	"github.com/n0stack/n0core/application/n0app/parser"
)

type Builder struct {
	operator operator.Operations
}

func (b Builder) Build(source []byte, stdout, stderr io.Writer) error {
	p, err := parser.Parse(source)
	if err != nil {
		return fmt.Errorf("Failed to parse, err:'%s'", err)
	}

	for i, s := range p.Sentences {
		fmt.Fprintf(stdout, "Step %d/%d : %s\n", i, len(p.Sentences), string(s.Line))

		var c chan operator.OutputLine
		switch string(s.Opecode.Value) {
		case "FROM":
			c, err = b.operator.FROM(s.Operands)
			if err != nil {
				return err
			}

		case "RUN":
			c, err = b.operator.RUN(s.Operands)
			if err != nil {
				return err
			}

		case "COPY":
			c, err = b.operator.COPY(s.Operands)
			if err != nil {
				return err
			}

		case "DAEMON":
			c, err = b.operator.DAEMON(s.Operands)
			if err != nil {
				return err
			}
		}

		for {
			o, ok := <-c // 処理の途中でフェイルする可能性を考慮する必要がある
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
