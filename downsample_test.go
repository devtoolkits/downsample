package downsample

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"
)

func Test_Downsample(t *testing.T) {
	type args struct {
		step int
	}
	tests := []struct {
		name string
		rps  Points
		args args
		want Points
	}{
		{"only 1 point",
			[]Point{
				Point{0, 0.1},
			},
			args{15},
			[]Point{
				Point{0, 0.1},
			},
		},
		{"not enough points",
			[]Point{
				Point{0, 0.5},
				Point{10, 0.3},
			},
			args{15},
			[]Point{
				Point{0, 0.4},
			},
		},
		{"boundary",
			[]Point{
				Point{0, 0.5},
				Point{15, 0.3},
			},
			args{15},
			[]Point{
				Point{0, 0.5},
				Point{15, 0.3},
			},
		},
		{"normal",
			[]Point{
				Point{10, 0.5},
				Point{20, 0.3},
				Point{30, 0.4},
			},
			args{15},
			[]Point{
				Point{10, 0.4},
				Point{25, 0.4},
			},
		},
		{"NaN#1",
			[]Point{
				Point{10, 0.5},
				Point{20, math.NaN()},
				Point{30, 0.4},
			},
			args{20},
			[]Point{
				Point{10, 0.5},
				Point{30, 0.4},
			},
		},
		{"NaN#2",
			[]Point{
				Point{10, 0.5},
				Point{20, math.NaN()},
				Point{30, math.NaN()},
				Point{40, 0.6},
				Point{50, 0.8},
			},
			args{20},
			[]Point{
				Point{10, 0.5},
				Point{30, 0.6},
				Point{50, 0.8},
			},
		},
		{"should not be happen",
			[]Point{
				Point{10, 0.5},
				Point{20, 0.3},
				Point{30, 0.4},
			},
			args{5},
			[]Point{
				Point{10, 0.5},
				Point{15, math.NaN()},
				Point{20, 0.3},
				Point{25, math.NaN()},
				Point{30, 0.4},
			},
		},
		{"except#1",
			[]Point{
				Point{10, 0.5},
				Point{10, 0.3},
			},
			args{30},
			[]Point{
				Point{10, 0.5},
			},
		},
		{"unsorted",
			[]Point{
				Point{10, 0.5},
				Point{0, 0.3},
			},
			args{20},
			[]Point{
				Point{0, 0.4},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rps.Downsample(tt.args.step)
			// math.NaN() != math.NaN()
			// so reflect.DeepEqual is not work when value is math.NaN()
			gotBytes, _ := json.Marshal(got)
			wantBytes, _ := json.Marshal(tt.want)
			if !reflect.DeepEqual(gotBytes, wantBytes) {
				t.Errorf("Points.aggrByStep() = %v, want %v", got, tt.want)
			}
		})
	}
}
