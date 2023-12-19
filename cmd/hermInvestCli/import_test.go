package main

import (
	"reflect"
	"testing"
)

func TestSwapColumn(t *testing.T) {
	type args struct {
		row     []string
		indexes string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{

		{
			name: "Valid indexes",
			args: args{
				row:     []string{"A", "B", "C", "D", "E"},
				indexes: "1",
			},
			want:    []string{"B", "B", "C", "D", "E"},
			wantErr: false,
		},
		{
			name: "Valid indexes",
			args: args{
				row:     []string{"A", "B", "C", "D", "E"},
				indexes: "0,1,2,3",
			},
			want:    []string{"A", "B", "C", "D", "E"},
			wantErr: false,
		},
		{
			name: "Valid indexes",
			args: args{
				row:     []string{"A", "B", "C", "D", "E"},
				indexes: "0,1,3,2",
			},
			want:    []string{"A", "B", "D", "C", "E"},
			wantErr: false,
		},
		{
			name: "Valid indexes",
			args: args{
				row:     []string{"A", "B", "C", "D", "E"},
				indexes: "0,1,3,2",
			},
			want:    []string{"A", "B", "D", "C", "E"},
			wantErr: false,
		},
		{
			name: "Valid indexes",
			args: args{
				row:     []string{"A", "B", "C", "D", "E"},
				indexes: "3,0,1,2",
			},
			want:    []string{"D", "A", "B", "C", "E"},
			wantErr: false,
		},
		{
			name: "Empty indexes",
			args: args{
				row:     []string{"A", "B", "C", "D", "E"},
				indexes: "",
			},
			want:    []string{"A", "B", "C", "D", "E"},
			wantErr: true,
		},
		{
			name: "Out of range",
			args: args{
				row:     []string{"A", "B", "C", "D", "E"},
				indexes: "3,0,1,6",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Exceed length",
			args: args{
				row:     []string{"A", "B", "C", "D", "E"},
				indexes: "3,0,1,2,4,1",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Parse error",
			args: args{
				row:     []string{"A", "B", "C", "D", "E"},
				indexes: "3,0,$,2,4",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		// tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()
			got, err := swapColumn(tt.args.row, tt.args.indexes)
			if (err != nil) != tt.wantErr {
				t.Errorf("swapColumn(%v, %v) error = %v, wantErr %v", tt.args.row, tt.args.indexes, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("swapColumn(%v, %v) = %v, want %v", tt.args.row, tt.args.indexes, got, tt.want)
			}
		})
	}
}
