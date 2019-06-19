package dag

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"testing"
)

func TestIsDAG(t *testing.T) {
	cases := []struct {
		name string
		dag  *DAG
		err  string
	}{
		{
			"loop",
			&DAG{
				map[string]*Task{
					"g1": {
						DependsOn: []string{
							"g2",
						},
					},
					"g2": {
						DependsOn: []string{
							"g3",
						},
					},
					"g3": {
						DependsOn: []string{
							"g1",
						},
					},
				},
			},
			"This request is not DAG",
		},
		{
			"liner",
			&DAG{
				map[string]*Task{
					"g1": {
						DependsOn: []string{
							"g2",
						},
					},
					"g2": {
						DependsOn: []string{
							"g3",
						},
					},
					"g3": {},
				},
			},
			"",
		},
	}

	for _, c := range cases {
		err := c.dag.Check()
		if (c.err == "") == (err != nil) {
			t.Errorf("[%s] wrong existence error: have=%+v, want=%+v", c.name, err != nil, c.err == "")
		}
		if c.err != "" && err.Error() != c.err {
			t.Errorf("[%s] got wrong err: have=%s, want='%s'", c.name, err.Error(), c.err)
		}
	}
}

func TestDoDAG(t *testing.T) {
	cases := []struct {
		name        string
		dag         *DAG
		retPatterns []string
		err         error
	}{
		{
			"forward",
			&DAG{
				map[string]*Task{
					"g1": {
						task: func(ctx context.Context, out io.Writer) error {
							fmt.Fprintf(out, "first")
							return nil
						},
					},
					"g2": {
						task: func(ctx context.Context, out io.Writer) error {
							fmt.Fprintf(out, "second")
							return nil
						},
						DependsOn: []string{
							"g1",
						},
					},
					"g3": {
						task: func(ctx context.Context, out io.Writer) error {
							fmt.Fprintf(out, "third")
							return nil
						},
						DependsOn: []string{
							"g2",
						},
					},
				},
			},
			[]string{
				"firstsecondthird",
			},
			nil,
		},
		{
			"triangle",
			&DAG{
				map[string]*Task{
					"g1": {
						task: func(ctx context.Context, out io.Writer) error {
							fmt.Fprintf(out, "first")
							return nil
						},
					},
					"g2-1": {
						task: func(ctx context.Context, out io.Writer) error {
							fmt.Fprintf(out, "second")
							return nil
						},
						DependsOn: []string{
							"g1",
						},
					},
					"g2-2": {
						task: func(ctx context.Context, out io.Writer) error {
							fmt.Fprintf(out, "third")
							return nil
						},
						DependsOn: []string{
							"g1",
						},
					},
				},
			},
			[]string{
				"firstsecondthird",
				"firstthirdsecond",
			},
			nil,
		},
		{
			"rollback",
			&DAG{
				map[string]*Task{
					"g1": {
						task: func(ctx context.Context, out io.Writer) error {
							fmt.Fprintf(out, "first")
							return nil
						},
						rollback: func(ctx context.Context, out io.Writer) error {
							fmt.Fprintf(out, "first")
							return nil
						},
					},
					"g2": {
						task: func(ctx context.Context, out io.Writer) error {
							fmt.Fprintf(out, "second")
							return nil
						},
						rollback: func(ctx context.Context, out io.Writer) error {
							fmt.Fprintf(out, "second")
							return nil
						},
						DependsOn: []string{
							"g1",
						},
					},
					"g3": {
						task: func(ctx context.Context, out io.Writer) error {
							fmt.Fprintf(out, "third")
							return fmt.Errorf("error")
						},
						DependsOn: []string{
							"g2",
						},
					},
				},
			},
			[]string{
				"firstsecondthirdsecondfirst",
			},
			&TaskErrors{
				RunErrors: map[string]error{
					"g3": errors.New("error"),
				},
			},
		},
	}

	ctx := context.Background()
	for _, c := range cases {
		out := &bytes.Buffer{}

		err := c.dag.Do(ctx, out)
		if err != nil && c.err != nil {
			if err.Error() != c.err.Error() {
				t.Errorf("[%s] DAG.Do() got wrong error:\n  have=%s\n  want=%s", c.name, err.Error(), c.err.Error())
			}
		} else {
			if err != c.err {
				t.Errorf("[%s] DAG.Do() got error: err=%s", c.name, err.Error())
			}
		}

		f := false
		for _, p := range c.retPatterns {
			f = f || (out.String() == p)
		}
		if !f {
			t.Errorf("[%s] DAG.Do() output wrong string: have=%s, want=%+v", c.name, out.String(), c.retPatterns)
		}
	}
}
