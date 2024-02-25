package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
					Operator: newToken(MINUS, "-", nil, 1),
					Right: &Literal{
						Value: 123,
					},
				},
				Operator: newToken(STAR, "*", nil, 1),
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
			out := newPrinter().Print(tc.input)
			if assert.Equal(t, tc.want, out) {
				t.Logf("%s: '%s'", tc.name, out)
			}
		})
	}
}
