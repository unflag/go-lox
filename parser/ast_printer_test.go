package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unflag/go-lox/scanner"
)

func Test_Printer(t *testing.T) {
	testCases := []struct {
		name  string
		input Expr
		want  string
	}{
		{
			name: "basic",
			input: &Binary{
				Left: &Unary{
					Operator: scanner.NewToken(scanner.MINUS, "-", nil, 1),
					Right: &Literal{
						Value: 123,
					},
				},
				Operator: scanner.NewToken(scanner.STAR, "*", nil, 1),
				Right: &Grouping{
					Expression: &Literal{
						Value: 45.67,
					},
				},
			},
			want: "(* (- 123) (group 45.67))",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := Printer[string]{}
			out := p.Print(tc.input)
			if assert.Equal(t, tc.want, out) {
				t.Logf("%s: '%s'", tc.name, out)
			}
		})
	}
}
