package property

import (
	"testing"

	"github.com/ruslanbes/kubrai/fileutils"
)

func TestAsInt(t *testing.T) {
	fileutils.FilePutContents(PropertiesPath+"/testIntNeg", "-42")
	defer fileutils.FileRemove(PropertiesPath + "/testIntNeg")

	fileutils.FilePutContents(PropertiesPath+"/testIntZero", "0")
	defer fileutils.FileRemove(PropertiesPath + "/testIntZero")

	fileutils.FilePutContents(PropertiesPath+"/testIntPos", "42")
	defer fileutils.FileRemove(PropertiesPath + "/testIntPos")

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
	fileutils.FilePutContents(PropertiesPath+"/testBoolOn", "ON")
	defer fileutils.FileRemove(PropertiesPath + "/testBoolOn")

	fileutils.FilePutContents(PropertiesPath+"/testBoolOff", "OFF")
	defer fileutils.FileRemove(PropertiesPath + "/testBoolOff")

	fileutils.FilePutContents(PropertiesPath+"/testBoolWrong", "Wrong")
	defer fileutils.FileRemove(PropertiesPath + "/testBoolWrong")

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
