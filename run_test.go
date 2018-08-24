package main

import (
	"reflect"
	"testing"
)

func Test_getCompErrs(t *testing.T) {
	type args struct {
		output string
	}
	tests := []struct {
		name string
		args args
		want []CompErr
	}{
		{
			name: "1 line err",
			args: args{output: "../../../../../../../var/folders/5s/pxq8rc1d6wx8d5f5vsbz5vth0000gn/T/example010888711/main.go:47:2: expected operand, found 'return'"},
			want: []CompErr{
				CompErr{
					Row:     46,
					Column:  2,
					Text:    "expected operand, found 'return'",
					ErrType: "error",
				},
			},
		},
		{
			name: "multi line err",
			args: args{output: `../../../../../../../var/folders/5s/pxq8rc1d6wx8d5f5vsbz5vth0000gn/T/example565701415/main.go:38:19: expected ';', found ','
			../../../../../../../var/folders/5s/pxq8rc1d6wx8d5f5vsbz5vth0000gn/T/example565701415/main.go:46:44: expected ';', found '!'
			../../../../../../../var/folders/5s/pxq8rc1d6wx8d5f5vsbz5vth0000gn/T/example565701415/main.go:67:2: expected declaration, found 'IDENT' screen`},
			want: []CompErr{
				CompErr{
					Row:     37,
					Column:  19,
					Text:    "expected ';', found ','",
					ErrType: "error",
				},
				CompErr{
					Row:     45,
					Column:  44,
					Text:    "expected ';', found '!'",
					ErrType: "error",
				},
				CompErr{
					Row:     66,
					Column:  2,
					Text:    "expected declaration, found 'IDENT' screen",
					ErrType: "error",
				},
			},
		},
		

		{
			name: "Error message with colons in it",
			args: args{output: `../../../../../../../var/folders/5s/pxq8rc1d6wx8d5f5vsbz5vth0000gn/T/example586307890/main.go:46:9: invalid operation: r (variable of type *invalid type) has no field or method w`},
			want: []CompErr{
				CompErr{
					Row:     45,
					Column:  9,
					Text:    "invalid operation: r (variable of type *invalid type) has no field or method w",
					ErrType: "error",
				},
			},
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getCompErrs(tt.args.output); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCompErrs() = %v, want %v", got, tt.want)
			}
		})
	}
}
