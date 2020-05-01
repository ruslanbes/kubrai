package property

import (
	"testing"

	"github.com/ruslanbes/kubrai/fileutils"
)

func TestAsInt(t *testing.T) {
	fileutils.FilePutContents(PropertyDir+"/testIntNeg", "-42")
	defer fileutils.FileRemove(PropertyDir + "/testIntNeg")

	fileutils.FilePutContents(PropertyDir+"/testIntZero", "0")
	defer fileutils.FileRemove(PropertyDir + "/testIntZero")

	fileutils.FilePutContents(PropertyDir+"/testIntPos", "42")
	defer fileutils.FileRemove(PropertyDir + "/testIntPos")

	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Neg",
			args: args{"testIntNeg"},
			want: -42,
		},
		{
			name: "Zero",
			args: args{"testIntZero"},
			want: 0,
		},
		{
			name: "Pos",
			args: args{"testIntPos"},
			want: 42,
		},
		{
			name: "NonExist",
			args: args{"testNonexistentFile"},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AsInt(tt.args.name); got != tt.want {
				t.Errorf("AsInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAsBool(t *testing.T) {
	fileutils.FilePutContents(PropertyDir+"/testBoolOn", "ON")
	defer fileutils.FileRemove(PropertyDir + "/testBoolOn")

	fileutils.FilePutContents(PropertyDir+"/testBoolOff", "OFF")
	defer fileutils.FileRemove(PropertyDir + "/testBoolOff")

	fileutils.FilePutContents(PropertyDir+"/testBoolWrong", "Wrong")
	defer fileutils.FileRemove(PropertyDir + "/testBoolWrong")

	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "On",
			args: args{"testBoolOn"},
			want: true,
		},
		{
			name: "Off",
			args: args{"testBoolOff"},
			want: false,
		},
		{
			name: "Wrong",
			args: args{"testBoolWrong"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AsBool(tt.args.name); got != tt.want {
				t.Errorf("AsBool() = %v, want %v", got, tt.want)
			}
		})
	}
}
