package utils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRemoveDuplicates(t *testing.T) {
	type args struct {
		slc []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "empty", args: args{[]string{}}, want: nil},
		{name: "blank", args: args{[]string{""}}, want: []string{""}},
		{name: "blank with comma", args: args{[]string{"", ""}}, want: []string{""}},
		{name: "duplicate elem", args: args{[]string{"1", "4", "3", "4", "1"}}, want: []string{"1", "4", "3"}},
		{name: "duplicate elem all", args: args{[]string{"1", "1", "1", "1", "1"}}, want: []string{"1"}},
		{name: "no duplicate elem", args: args{[]string{"1", "2", "3", "4", "5"}}, want: []string{"1", "2", "3", "4", "5"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveDuplicates(tt.args.slc); !reflect.DeepEqual(got, tt.want) {
				fmt.Printf("got,  %+v, %T, %p, %v\n", got, got, got, len(got))
				fmt.Printf("want, %+v, %T, %p, %v\n", tt.want, tt.want, tt.want, len(tt.want))
				t.Errorf("RemoveDuplicates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveDuplicatesByLoop(t *testing.T) {
	type args struct {
		slc []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "empty", args: args{[]string{}}, want: nil},
		{name: "blank", args: args{[]string{""}}, want: []string{""}},
		{name: "blank with comma", args: args{[]string{"", ""}}, want: []string{""}},
		{name: "duplicate elem", args: args{[]string{"1", "4", "3", "4", "1"}}, want: []string{"1", "4", "3"}},
		{name: "duplicate elem all", args: args{[]string{"1", "1", "1", "1", "1"}}, want: []string{"1"}},
		{name: "no duplicate elem", args: args{[]string{"1", "2", "3", "4", "5"}}, want: []string{"1", "2", "3", "4", "5"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveDuplicatesByLoop(tt.args.slc); !reflect.DeepEqual(got, tt.want) {
				fmt.Printf("got,  %+v, %T, %p, %v\n", got, got, got, len(got))
				fmt.Printf("want, %+v, %T, %p, %v\n", tt.want, tt.want, tt.want, len(tt.want))
				t.Errorf("RemoveDuplicatesByLoop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveDuplicatesByMap(t *testing.T) {
	type args struct {
		slc []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "empty", args: args{[]string{}}, want: nil},
		{name: "blank", args: args{[]string{""}}, want: []string{""}},
		{name: "blank with comma", args: args{[]string{"", ""}}, want: []string{""}},
		{name: "duplicate elem", args: args{[]string{"1", "4", "3", "4", "1"}}, want: []string{"1", "4", "3"}},
		{name: "duplicate elem all", args: args{[]string{"1", "1", "1", "1", "1"}}, want: []string{"1"}},
		{name: "no duplicate elem", args: args{[]string{"1", "2", "3", "4", "5"}}, want: []string{"1", "2", "3", "4", "5"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveDuplicatesByMap(tt.args.slc); !reflect.DeepEqual(got, tt.want) {
				fmt.Printf("got,  %+v, %T, %p, %v\n", got, got, got, len(got))
				fmt.Printf("want, %+v, %T, %p, %v\n", tt.want, tt.want, tt.want, len(tt.want))
				t.Errorf("RemoveDuplicatesByMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveDuplicatesWithSort(t *testing.T) {
	type args struct {
		slc []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "empty", args: args{[]string{}}, want: nil},
		{name: "blank", args: args{[]string{""}}, want: []string{""}},
		{name: "blank with comma", args: args{[]string{"", ""}}, want: []string{""}},
		{name: "duplicate elem", args: args{[]string{"1", "4", "3", "4", "1"}}, want: []string{"1", "3", "4"}},
		{name: "duplicate elem all", args: args{[]string{"1", "1", "1", "1", "1"}}, want: []string{"1"}},
		{name: "no duplicate elem", args: args{[]string{"1", "2", "3", "4", "5"}}, want: []string{"1", "2", "3", "4", "5"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveDuplicatesWithSort(tt.args.slc); !reflect.DeepEqual(got, tt.want) {
				fmt.Printf("got,  %+v, %T, %p, %v\n", got, got, got, len(got))
				fmt.Printf("want, %+v, %T, %p, %v\n", tt.want, tt.want, tt.want, len(tt.want))
				t.Errorf("RemoveDuplicatesWithSort() = %v, want %v", got, tt.want)
			}
		})
	}
}
