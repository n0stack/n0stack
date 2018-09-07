// +build ignore

package parser

import (
	"fmt"
	"testing"
)

func TestCheckVersion(t *testing.T) {
	Cases := []struct {
		name    string
		source  []byte
		program *Program
		err     string
	}{
		{
			name:   "[Valid] simple <opecode> <operand_string>",
			source: []byte("FROM test"),
			program: &Program{
				Sentences: []*Sentence{
					&Sentence{
						Line:       []byte("FROM test"),
						LineNumber: 1,
						Opecode: &Token{
							Kind:  Opecode,
							Value: []byte("FROM"),
						},
						Operands: []*Token{
							&Token{
								Kind:  Operand_String,
								Value: []byte("test"),
							},
						},
					},
				},
			},
			err: "",
		},
		{
			name:   "[Valid] comment <opecode> <operand_string> <comment>",
			source: []byte("FROM test # comment"),
			program: &Program{
				Sentences: []*Sentence{
					&Sentence{
						Line:       []byte("FROM test # comment"),
						LineNumber: 1,
						Opecode: &Token{
							Kind:  Opecode,
							Value: []byte("FROM"),
						},
						Operands: []*Token{
							&Token{
								Kind:  Operand_String,
								Value: []byte("test"),
							},
							&Token{
								Kind:  Comment,
								Value: []byte("# comment"),
							},
						},
					},
				},
			},
			err: "",
		},
		// {
		// 	name:   "[Valid] comment line <comment>",
		// 	source: []byte("# comment"),
		// 	program: &Program{
		// 		Sentences: []*Sentence{
		// 			&Sentence{
		// 				Line:       []byte("# comment"),
		// 				LineNumber: 1,
		// 				Tokens: []*Token{
		// 					&Token{
		// 						Kind:  Comment,
		// 						Value: []byte("# comment"),
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// 	err: "",
		// },
		{
			name:   "[Valid] multiple operands <opecode> <operand_string> <operand_string>",
			source: []byte("FROM ope1 ope2"),
			program: &Program{
				Sentences: []*Sentence{
					&Sentence{
						Line:       []byte("FROM ope1 ope2"),
						LineNumber: 1,
						Opecode: &Token{
							Kind:  Opecode,
							Value: []byte("FROM"),
						},
						Operands: []*Token{
							&Token{
								Kind:  Operand_String,
								Value: []byte("ope1"),
							},
							&Token{
								Kind:  Operand_String,
								Value: []byte("ope2"),
							},
						},
					},
				},
			},
			err: "",
		},
		{
			name:   "[Valid] single quote with space <opecode> '<operand_string>'",
			source: []byte("FROM 'ope1 ope2'"),
			program: &Program{
				Sentences: []*Sentence{
					&Sentence{
						Line:       []byte("FROM 'ope1 ope2'"),
						LineNumber: 1,
						Opecode: &Token{
							Kind:  Opecode,
							Value: []byte("FROM"),
						},
						Operands: []*Token{
							&Token{
								Kind:  Operand_String,
								Value: []byte("ope1 ope2"),
							},
						},
					},
				},
			},
			err: "",
		},
		{
			name:   "[Valid] double quote with space <opecode> \"<operand_string>\"",
			source: []byte("FROM \"ope1 ope2\""),
			program: &Program{
				Sentences: []*Sentence{
					&Sentence{
						Line:       []byte("FROM \"ope1 ope2\""),
						LineNumber: 1,
						Opecode: &Token{
							Kind:  Opecode,
							Value: []byte("FROM"),
						},
						Operands: []*Token{
							&Token{
								Kind:  Operand_String,
								Value: []byte("ope1 ope2"),
							},
						},
					},
				},
			},
			err: "",
		},
		{
			name:   "[Valid] multiple lines <opecode> <operand_string>\\n<operand_string>",
			source: []byte(fmt.Sprintf("FROM ope1\\\nope2")), // TODO: 正しいかわからない
			program: &Program{
				Sentences: []*Sentence{
					&Sentence{
						Line:       []byte("FROM ope1 ope2"),
						LineNumber: 1,
						Opecode: &Token{
							Kind:  Opecode,
							Value: []byte("FROM"),
						},
						Operands: []*Token{
							&Token{
								Kind:  Operand_String,
								Value: []byte("ope1"),
							},
							&Token{
								Kind:  Operand_String,
								Value: []byte("ope2"),
							},
						},
					},
				},
			},
			err: "",
		},
		{
			name:   "[Valid] multiple lines in a operand <opecode> <operand_string>",
			source: []byte(fmt.Sprintf("FROM 'ope1\nope2'")),
			program: &Program{
				Sentences: []*Sentence{
					&Sentence{
						Line:       []byte("FROM 'ope1ope2'"),
						LineNumber: 1,
						Opecode: &Token{
							Kind:  Opecode,
							Value: []byte("FROM"),
						},
						Operands: []*Token{
							&Token{
								Kind:  Operand_String,
								Value: []byte("ope1"),
							},
						},
					},
				},
			},
			err: "",
		},
	}

	for _, c := range Cases {
		p, err := Parse(c.source)

		if p != c.program {
			t.Errorf("[%s] Wrong parsed\n\thave:%d\n\twant:%v", c.name, p, c.program)
		}

		if err == nil && c.err != "" {
			t.Errorf("[%s] Need error message\n\thave:nil\n\twant:%s", c.name, c.err)
		}

		if err != nil && err.Error() != c.err {
			t.Errorf("[%s] Wrong error message\n\thave:%s\n\twant:%s", c.name, err.Error(), c.err)
		}
	}
}
