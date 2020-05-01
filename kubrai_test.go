package main

import (
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/ruslanbes/kubrai/fileutils"
	"github.com/ruslanbes/kubrai/property"
)

func setUpTestProperties(props map[string]string) {
	testPropertyDir := "./test/data/properties"
	os.RemoveAll(testPropertyDir)
	os.Mkdir(testPropertyDir, 0777)

	property.PropertiesPath = testPropertyDir
	property.SetProperties(props)
}

func Test_findExactVerb(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "empty",
			args: args{[]string{}},
			want: "",
		},
		{
			name: "play",
			args: args{[]string{"play"}},
			want: "play",
		},
		{
			name: "wrong",
			args: args{[]string{"blablabla"}},
			want: "",
		},
		{
			name: "wrong and right",
			args: args{[]string{"blablabla", "add"}},
			want: "add",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findExactVerb(tt.args.args); got != tt.want {
				t.Errorf("findExactVerb() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_guessVerb(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "play",
			args: args{[]string{}},
			want: "play",
		},
		{
			name: "solve",
			args: args{[]string{"КУБ_РАЯ"}},
			want: "solve",
		},
		{
			name: "view",
			args: args{[]string{"КУБ"}},
			want: "view",
		},
		{
			name: "wrongCmd",
			args: args{[]string{"vieww", "КУБ"}},
			want: "",
		},
		{
			name: "nothing",
			args: args{[]string{"КУБ_РАЯ", "ШАР"}},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := guessVerb(tt.args.args); got != tt.want {
				t.Errorf("guessVerb() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseVerb(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Solve",
			args: args{[]string{"ЛОДКА_В_Я"}},
			want: "solve",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseVerb(tt.args.args); got != tt.want {
				t.Errorf("parseVerb() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractArgs(t *testing.T) {
	setUpTestProperties(map[string]string{propArgsAutoLowercase: "ON"})

	type args struct {
		word  string
		words []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "One",
			args: args{
				word:  "one",
				words: []string{"One", "two"},
			},
			want: []string{"two"},
		},
		{
			name: "Many",
			args: args{
				word:  "one",
				words: []string{"one", "Two", "one"},
			},
			want: []string{"two"},
		},
		{
			name: "None",
			args: args{
				word:  "one",
				words: []string{"five", "two", "boom"},
			},
			want: []string{"five", "two", "boom"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractArgs(tt.args.word, tt.args.words); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractArgs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_backupName(t *testing.T) {
	type args struct {
		file      string
		backupNum int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Default",
			args: args{"test/file.txt", 5},
			want: "test/file.txt.5.bak",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := backupName(tt.args.file, tt.args.backupNum); got != tt.want {
				t.Errorf("backupName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_canonize(t *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Default",
			args: args{"bLa "},
			want: "BLA",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := canonize(tt.args.word); got != tt.want {
				t.Errorf("canonize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_saveAndLoadAssoc(t *testing.T) {
	setUpTestProperties(map[string]string{
		propPlaybooksDir:          "./test/data/playbooks",
		propPlaybookCurrent:       "default",
		propAssocFileKeySeparator: ":",
		propAssocFileValSeparator: ",",
	})

	assocFile := getCurrentPlaybookDir() + "/associations/associationsTest.txt"
	assoc := make(map[string][]string)
	assoc["aaa"] = []string{"bbb", "ccc"}

	saveAssoc(assocFile, assoc)
	got := loadAssoc(assocFile)

	if !reflect.DeepEqual(assoc, got) {
		t.Errorf("loadAssoc() = %v, want %v", got, assoc)
	}

	//cleanup
	err := os.Remove(assocFile)
	checkError(err)
}

func Test_addBeforeFirstLonger(t *testing.T) {
	type args struct {
		s   string
		slc []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Begin",
			args: args{"test", []string{"test1", "test2"}},
			want: []string{"test", "test1", "test2"},
		},
		{
			name: "Mid",
			args: args{"test11", []string{"test1", "test111"}},
			want: []string{"test1", "test11", "test111"},
		},
		{
			name: "End",
			args: args{"test", []string{"te", "t"}},
			want: []string{"te", "t", "test"},
		},
		{
			name: "SameLength",
			args: args{"test1", []string{"test9", "test0"}},
			want: []string{"test9", "test0", "test1"},
		},
		{
			name: "Empty",
			args: args{"test1", []string{}},
			want: []string{"test1"},
		},
		{
			name: "AlreadyThere",
			args: args{"test1", []string{"test1", "test111"}},
			want: []string{"test1", "test111"},
		},
		{
			name: "Cyrillic noncense",
			args: args{"БУЛАВА", []string{"ГОДНОСТЬ", "ЛЕТО"}},
			want: []string{"БУЛАВА", "ГОДНОСТЬ", "ЛЕТО"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addBeforeFirstLonger(tt.args.s, tt.args.slc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addBeforeFirstLonger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_runAdd(t *testing.T) {
	setUpTestProperties(map[string]string{
		propAddValMayEqualKey:     "ON",
		propAssocFileKeySeparator: ":",
		propAssocFileValSeparator: ",",
		propPlaybookCurrent:       "default",
		propPlaybooksDir:          "./test/data/playbooks",
	})

	saveDefaultAssoc(map[string][]string{})

	type args struct {
		a string
		b string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Test01",
			args: args{"boy", "girl"},
			want: []string{"girl"},
		},
		{
			name: "Test02",
			args: args{"boy", "man"},
			want: []string{"man", "girl"},
		},
		{
			name: "Test03",
			args: args{"girl", "woman"},
			want: []string{"woman"},
		},
		{
			name: "Test04",
			args: args{"boy", "child"},
			want: []string{"man", "girl", "child"},
		},
		{
			name: "Test05",
			args: args{"math", "science"},
			want: []string{"science"},
		},
		{
			name: "Test06",
			args: args{"math", "philosophy"},
			want: []string{"science", "philosophy"},
		},
		{
			name: "Test06Again",
			args: args{"math", "philosophy"},
			want: []string{"science", "philosophy"},
		},
		{
			name: "Test07",
			args: args{"math", "math"},
			want: []string{"math", "science", "philosophy"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runAdd(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("runAdd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_runAddBoth(t *testing.T) {
	setUpTestProperties(map[string]string{
		propAssocFileKeySeparator: ":",
		propAssocFileValSeparator: ",",
		propPlaybookCurrent:       "default",
		propPlaybooksDir:          "./test/data/playbooks",
	})

	saveDefaultAssoc(map[string][]string{})

	type args struct {
		a string
		b string
	}
	tests := []struct {
		name string
		args args
		want [2][]string
	}{
		{
			name: "Test01",
			args: args{"boy", "girl"},
			want: [2][]string{{"girl"}, {"boy"}},
		},
		{
			name: "Test02",
			args: args{"girl", "woman"},
			want: [2][]string{{"boy", "woman"}, {"girl"}},
		},
		{
			name: "Test03",
			args: args{"woman", "man"},
			want: [2][]string{{"man", "girl"}, {"woman"}},
		},
		{
			name: "Test04",
			args: args{"man", "boy"},
			want: [2][]string{{"boy", "woman"}, {"man", "girl"}},
		},
		{
			name: "Test05",
			args: args{"boy", "child"},
			want: [2][]string{{"man", "girl", "child"}, {"boy"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runAddBoth(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("runAddBoth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_runSmartAdd(t *testing.T) {
	setUpTestProperties(map[string]string{
		propAddAutoBothMaxlen:     "3",
		propAssocFileKeySeparator: ":",
		propAssocFileValSeparator: ",",
		propPlaybookCurrent:       "default",
		propPlaybooksDir:          "./test/data/playbooks",
	})

	saveDefaultAssoc(map[string][]string{})

	type args struct {
		a string
		b string
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{
			name: "Test01",
			args: args{"1", "12"},
			want: map[string][]string{"1": {"12"}, "12": {"1"}},
		},
		{
			name: "Test02",
			args: args{"12", "123"},
			want: map[string][]string{"12": {"1", "123"}, "123": {"12"}},
		},
		{
			name: "Test03",
			args: args{"123", "1234"},
			want: map[string][]string{"123": {"12", "1234"}, "1234": {"123"}},
		},
		{
			name: "Test04",
			args: args{"1234", "12345"},
			want: map[string][]string{"1234": {"123", "12345"}},
		},
		{
			name: "Test05",
			args: args{"12345", "1"},
			want: map[string][]string{"12345": {"1"}},
		},
		{
			name: "Test01again",
			args: args{"1", "12"},
			want: map[string][]string{"1": {"12"}, "12": {"1", "123"}},
		},
		{
			name: "Test01cyrillic",
			args: args{"шаг", "па"},
			want: map[string][]string{"шаг": {"па"}, "па": {"шаг"}},
		},
		{
			name: "Test02cyrillic",
			args: args{"па", "ап"},
			want: map[string][]string{"па": {"ап", "шаг"}, "ап": {"па"}},
		},
		{
			name: "Test03cyrillic",
			args: args{"ватта", "ома"},
			want: map[string][]string{"ватта": {"ома"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runSmartAdd(tt.args.a, tt.args.b)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("runSmartAdd() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeByValue(t *testing.T) {
	type args struct {
		s   string
		slc []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Begin",
			args: args{"aa", []string{"aa", "bb", "cc"}},
			want: []string{"bb", "cc"},
		},
		{
			name: "Mid",
			args: args{"bb", []string{"aa", "bb", "cc"}},
			want: []string{"aa", "cc"},
		},
		{
			name: "End",
			args: args{"cc", []string{"aa", "bb", "cc"}},
			want: []string{"aa", "bb"},
		},
		{
			name: "None",
			args: args{"dd", []string{"aa", "bb", "cc"}},
			want: []string{"aa", "bb", "cc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeByValue(tt.args.s, tt.args.slc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeByValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_runRemove(t *testing.T) {
	setUpTestProperties(map[string]string{
		propAssocFileKeySeparator: ":",
		propAssocFileValSeparator: ",",
		propPlaybookCurrent:       "default",
		propPlaybooksDir:          "./test/data/playbooks",
	})

	assoc := make(map[string][]string)
	assoc["boy"] = []string{"girl", "man", "child"}
	assoc["girl"] = []string{"woman"}

	saveDefaultAssoc(assoc)

	type args struct {
		a string
		b string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "NonExistentKey",
			args: args{"nonexistent", "anything"},
			want: []string{},
		},
		{
			name: "NonExistentVal",
			args: args{"boy", "magnet"},
			want: []string{"girl", "man", "child"},
		},
		{
			name: "Mid",
			args: args{"boy", "man"},
			want: []string{"girl", "child"},
		},
		{
			name: "Begin",
			args: args{"boy", "girl"},
			want: []string{"child"},
		},
		{
			name: "End",
			args: args{"boy", "child"},
			want: []string{},
		},
		{
			name: "NonExistentValAgain",
			args: args{"boy", "magnet"},
			want: []string{},
		},
		{
			name: "Again",
			args: args{"girl", "woman"},
			want: []string{},
		},
		{
			name: "AndAgain",
			args: args{"girl", "woman"},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runRemove(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("runRemove() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_runRemoveBoth(t *testing.T) {
	setUpTestProperties(map[string]string{
		propAssocFileKeySeparator: ":",
		propAssocFileValSeparator: ",",
		propPlaybookCurrent:       "default",
		propPlaybooksDir:          "./test/data/playbooks",
	})

	assoc := make(map[string][]string)
	assoc["from"] = []string{"to", "subject", "cc"}
	assoc["to"] = []string{"from"}
	assoc["subject"] = []string{"body"}

	saveDefaultAssoc(assoc)

	type args struct {
		a string
		b string
	}
	tests := []struct {
		name string
		args args
		want [2][]string
	}{
		{
			name: "Test0RemoveWhatIsntThere",
			args: args{"foo", "bar"},
			want: [2][]string{{}, {}},
		},
		{
			name: "Test1",
			args: args{"from", "to"},
			want: [2][]string{{"subject", "cc"}, {}},
		},
		{
			name: "Test2",
			args: args{"from", "subject"},
			want: [2][]string{{"cc"}, {"body"}},
		},
		{
			name: "Test3",
			args: args{"from", "cc"},
			want: [2][]string{{}, {}},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runRemoveBoth(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("runRemoveBoth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_runView(t *testing.T) {
	setUpTestProperties(map[string]string{
		propAssocFileKeySeparator: ":",
		propAssocFileValSeparator: ",",
		propPlaybookCurrent:       "default",
		propPlaybooksDir:          "./test/data/playbooks",
	})

	assoc := make(map[string][]string)
	assoc["boy"] = []string{"girl", "man", "child"}
	assoc["girl"] = []string{"woman"}

	saveDefaultAssoc(assoc)

	type args struct {
		a string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "NonExistentKey",
			args: args{"nonexistent"},
			want: []string{},
		},
		{
			name: "Test01",
			args: args{"boy"},
			want: []string{"girl", "man", "child"},
		},
		{
			name: "Test02",
			args: args{"girl"},
			want: []string{"woman"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runView(tt.args.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("runView() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadDicts(t *testing.T) {
	setUpTestProperties(map[string]string{
		propPlaybookCurrent: "default",
		propPlaybooksDir:    "./test/data/playbooks",
		propDictsExt:        ".test",
	})

	dictsDir := getFullDictsDir()
	fileutils.FilePutContents(dictsDir+"/"+"dict.test", "aa\nbb")
	defer fileutils.FileRemove(dictsDir + "/" + "dict.test")

	tests := []struct {
		name string
		want map[string][]string
	}{
		{
			name: "Test01",
			want: map[string][]string{"dict.test": {"aa", "bb"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loadDicts(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadDicts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_nextMultiDimValue(t *testing.T) {
	type args struct {
		counter    []int
		maxCounter []int
	}
	tests := []struct {
		name  string
		args  args
		want  []int
		want1 bool
	}{
		{
			name:  "Test000",
			args:  args{[]int{0, 0, 0}, []int{1, 2, 3}},
			want:  []int{1, 0, 0},
			want1: true,
		},
		{
			name:  "Test100",
			args:  args{[]int{1, 0, 0}, []int{1, 2, 3}},
			want:  []int{0, 1, 0},
			want1: true,
		},
		{
			name:  "Test010",
			args:  args{[]int{0, 1, 0}, []int{1, 2, 3}},
			want:  []int{1, 1, 0},
			want1: true,
		},
		{
			name:  "Test110",
			args:  args{[]int{1, 1, 0}, []int{1, 2, 3}},
			want:  []int{0, 2, 0},
			want1: true,
		},
		{
			name:  "Test020",
			args:  args{[]int{0, 2, 0}, []int{1, 2, 3}},
			want:  []int{1, 2, 0},
			want1: true,
		},
		{
			name:  "Test120",
			args:  args{[]int{1, 2, 0}, []int{1, 2, 3}},
			want:  []int{0, 0, 1},
			want1: true,
		},
		{
			name:  "Test003",
			args:  args{[]int{0, 0, 3}, []int{1, 2, 3}},
			want:  []int{1, 0, 3},
			want1: true,
		},
		{
			name:  "Test123",
			args:  args{[]int{1, 2, 3}, []int{1, 2, 3}},
			want:  []int{0, 0, 0},
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := nextMultiDimValue(tt.args.counter, tt.args.maxCounter)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("nextMultiDimValue() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("nextMultiDimValue() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_combinations(t *testing.T) {
	type args struct {
		items [][]string
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "Test1",
			args: args{[][]string{{"a", "b", "c"}, {"1", "2"}, {"TRUE"}, {"?"}}},
			want: [][]string{
				{"a", "1", "TRUE", "?"},
				{"b", "1", "TRUE", "?"},
				{"c", "1", "TRUE", "?"},
				{"a", "2", "TRUE", "?"},
				{"b", "2", "TRUE", "?"},
				{"c", "2", "TRUE", "?"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := combinations(tt.args.items); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("combinations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_runSolve(t *testing.T) {
	setUpTestProperties(map[string]string{
		propAssocFileKeySeparator: ":",
		propAssocFileValSeparator: ",",
		propPlaybookCurrent:       "default",
		propPlaybooksDir:          "./test/data/playbooks",
		propDictsExt:              ".test",
		propSolveMaxResults:       "3",
	})

	dictsDir := getFullDictsDir()
	fileutils.FilePutContents(dictsDir+"/"+"dict.test", strings.Join(
		[]string{
			"boycott",
			"copy",
			"proximity",
			"terminator",
			"mm",
			"am",
			"mi",
			"mn",
			"ai",
		}, "\n"))
	defer fileutils.FileRemove(dictsDir + "/" + "dict.test")

	fileutils.FilePutContents(dictsDir+"/"+"smallDict.test", strings.Join(
		[]string{
			"boycott",
			"copy",
		}, "\n"))
	defer fileutils.FileRemove(dictsDir + "/" + "smallDict.test")

	assoc := make(map[string][]string)
	assoc["policeman"] = []string{"cop", "thief"}
	assoc["why"] = []string{"y", "not", "ask"}
	assoc["girl"] = []string{"woman", "boy"}
	assoc["bed"] = []string{"sleep", "cot", "zzz"}
	assoc["tea"] = []string{"t", "coffee"}
	assoc["period"] = []string{"term"}
	assoc["out"] = []string{"in", "fall"}
	assoc["&t"] = []string{"at"}
	assoc["and"] = []string{"&", "or"}
	assoc["amateur"] = []string{"pro"}
	assoc["psi"] = []string{"xi", "storm"}
	assoc["6"] = []string{"7", "six", "mi"}
	assoc["thanks"] = []string{"yw", "ty"}
	assoc["max"] = []string{"m", "a", "x"}
	assoc["min"] = []string{"m", "i", "n"}

	saveDefaultAssoc(assoc)

	type args struct {
		kubraya string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 bool
	}{
		{
			name:  "Test01",
			args:  args{"policeman_why"},
			want:  []string{"copy"},
			want1: true,
		},
		{
			name:  "Test02",
			args:  args{"girl_bed_tea"},
			want:  []string{"boycott"},
			want1: true,
		},
		{
			name:  "Test03",
			args:  args{"period_out_&t_and"},
			want:  []string{"terminator"},
			want1: true,
		},
		{
			name:  "Test04",
			args:  args{"amateur_psi_6_thanks"},
			want:  []string{"proximity"},
			want1: true,
		},
		{
			name:  "Test05",
			args:  args{"amateur_psi_nonexistentword_thanks"},
			want:  []string{},
			want1: false,
		},
		{
			name:  "Test05",
			args:  args{"amateur_psi_why_thanks"},
			want:  []string{},
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := runSolve(tt.args.kubraya)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("runSolve() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("runSolve() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	t.Run("limit", func(t *testing.T) {
		got, got1 := runSolve("max_min")
		if !got1 {
			t.Errorf("runSolve() got1 = %v, want %v", got1, true)
		}

		sliceSize := 3
		if len(got) != sliceSize {
			t.Errorf("runSolve() got = %v, want slice size %d", got, sliceSize)
		}
	})
}

func Test_runSearchDict(t *testing.T) {
	setUpTestProperties(map[string]string{
		propPlaybookCurrent: "default",
		propPlaybooksDir:    "./test/data/playbooks",
		propDictsExt:        ".test",
	})

	dictsDir := getFullDictsDir()
	fileutils.FilePutContents(dictsDir+"/"+"dict.test", strings.Join(
		[]string{
			"worda",
			"wordb",
			"wordc",
			"wordd",
		}, "\n"))
	defer fileutils.FileRemove(dictsDir + "/" + "dict.test")

	fileutils.FilePutContents(dictsDir+"/"+"oddDict.test", strings.Join(
		[]string{
			"wordb",
			"wordd",
		}, "\n"))
	defer fileutils.FileRemove(dictsDir + "/" + "oddDict.test")

	type args struct {
		word       string
		maxResults int
	}
	tests := []struct {
		name string
		args args
		want map[string]int
	}{
		{
			name: "Test0a",
			args: args{"worda", 0},
			want: map[string]int{},
		},
		{
			name: "Test1a",
			args: args{"worda", 1},
			want: map[string]int{"dict.test": 0},
		},
		{
			name: "Test2b",
			args: args{"wordb", 2},
			want: map[string]int{"dict.test": 1, "oddDict.test": 0},
		},
		{
			name: "Test3d",
			args: args{"wordd", 3},
			want: map[string]int{"dict.test": 3, "oddDict.test": 1},
		},
		{
			name: "Test4e",
			args: args{"worde", 4},
			want: map[string]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runSearchDict(tt.args.word, tt.args.maxResults); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("runSearchDict() = %v, want %v", got, tt.want)
			}
		})
	}

	// test limit. Result is random
	got := runSearchDict("wordb", 1)
	if val, ok := got["dict.test"]; ok && val != 1 {
		t.Errorf("runSearchDict() = %v, test limit failed", got)
	} else if val, ok := got["oddDict.test"]; ok && val != 0 {
		t.Errorf("runSearchDict() = %v, test limit failed", got)
	}

}

func Test_allowUnknowns(t *testing.T) {
	setUpTestProperties(map[string]string{
		propGuessUnknownMarker: "???",
	})

	type args struct {
		kubAssoc [][]string
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "Test1",
			args: args{[][]string{{"one", "two"}, {}, {"a", "b", "c"}}},
			want: [][]string{{"one", "two", "???"}, {"???"}, {"a", "b", "c", "???"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := allowUnknowns(tt.args.kubAssoc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("allowUnknowns() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_searchDictByRegexpGetWords(t *testing.T) {
	setUpTestProperties(map[string]string{
		propPlaybookCurrent: "default",
		propPlaybooksDir:    "./test/data/playbooks",
		propDictsExt:        ".test",
	})

	dictsDir := getFullDictsDir()
	fileutils.FilePutContents(dictsDir+"/"+"dict.test", strings.Join(
		[]string{
			"aabbcc",
			"abc",
			"bcd",
			"def",
			"папа",
			"мама",
			"брат",
		}, "\n"))
	defer fileutils.FileRemove(dictsDir + "/" + "dict.test")

	type args struct {
		re         *regexp.Regexp
		maxResults int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Test0",
			args: args{regexp.MustCompile("aa.+cc"), 0},
			want: []string{},
		},
		{
			name: "Test1",
			args: args{regexp.MustCompile("aa.+cc"), 50},
			want: []string{"aabbcc"},
		},
		{
			name: "Test2",
			args: args{regexp.MustCompile("a.+.+"), 50},
			want: []string{"aabbcc", "abc"},
		},
		{
			name: "TestCyrillic1",
			args: args{regexp.MustCompile("па.+"), 50},
			want: []string{"папа"},
		},
		{
			name: "TestCyrillic2",
			args: args{regexp.MustCompile(".+па.+"), 50},
			want: []string{},
		},
		{
			name: "TestCyrillic3",
			args: args{regexp.MustCompile(".+а.+"), 50},
			want: []string{"папа", "мама", "брат"},
		},
		{
			name: "TestCyrillic4",
			args: args{regexp.MustCompile(".+а.+"), 2},
			want: []string{"папа", "мама"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := searchDictByRegexpGetWords(tt.args.re, tt.args.maxResults); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("searchDictByRegexpGetWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_runGuess(t *testing.T) {
	setUpTestProperties(map[string]string{
		propAssocFileKeySeparator: ":",
		propAssocFileValSeparator: ",",
		propGuessExplainResults:   "OFF",
		propGuessMaxResults:       "50",
		propGuessUnknownsLimit:    "2",
		propGuessUnknownMarker:    "???",
		propPlaybookCurrent:       "default",
		propPlaybooksDir:          "./test/data/playbooks",
		propDictsExt:              ".test",
		propSolveMaxResults:       "2",
	})

	dictsDir := getFullDictsDir()
	fileutils.FilePutContents(dictsDir+"/"+"dict.test", strings.Join(
		[]string{
			"boycott",
			"copy",
			"proximity",
			"terminator",
			"boyscout",
			"cowboy",
		}, "\n"))
	defer fileutils.FileRemove(dictsDir + "/" + "dict.test")

	assoc := make(map[string][]string)
	assoc["policeman"] = []string{"cop", "thief"}
	assoc["girl"] = []string{"boy"}
	assoc["tea"] = []string{"t", "coffee"}
	assoc["period"] = []string{"term"}
	assoc["out"] = []string{"in", "fall"}
	assoc["&t"] = []string{"at"}
	assoc["and"] = []string{"&", "or"}
	assoc["amateur"] = []string{"pro"}
	assoc["6"] = []string{"7", "six", "mi"}
	assoc["thanks"] = []string{"yw", "ty"}

	saveDefaultAssoc(assoc)

	type args struct {
		kubraya string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 bool
	}{
		{
			name:  "Test01",
			args:  args{"policeman_why"},
			want:  []string{"copy"},
			want1: true,
		},
		{
			name:  "Test03",
			args:  args{"period_out_&t_and"},
			want:  []string{"terminator"},
			want1: true,
		},
		{
			name:  "Test04",
			args:  args{"amateur_psi_6_thanks"},
			want:  []string{"proximity"},
			want1: true,
		},
		{
			name:  "Test05a",
			args:  args{"amateur_psi_nonexistentword_thanks"},
			want:  []string{"proximity"},
			want1: true,
		},
		{
			name:  "Test05b",
			args:  args{"amateur_psi_thanks_nonexistentword"},
			want:  []string{},
			want1: false,
		},
		{
			name:  "Test06",
			args:  args{"amateur_psi_tea_thanks"},
			want:  []string{"proximity"},
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := runGuess(tt.args.kubraya)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("runGuess() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("runGuess() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	// test multival. need adjustment for (more results exist)
	// got, got1 := runGuess("girl_bed_tea")
	// if got1 {
	// 	if !reflect.DeepEqual(got, []string{"boycott", "boyscout"}) && !reflect.DeepEqual(got, []string{"boyscout", "boycott"}) {
	// 		t.Errorf("runGuess() got = %v, want %v", got, "[boycott, boyscout] or vice-versa")
	// 	}

	// } else {
	// 	t.Errorf("runGuess() got1 = %v, want %v", got1, true)
	// }

}

func Test_countUnknownsInWord(t *testing.T) {
	setUpTestProperties(map[string]string{
		propGuessUnknownMarker: "???",
	})

	type args struct {
		word string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Test01",
			args: args{"??????ан"},
			want: 2,
		},
		{
			name: "Test02",
			args: args{"???триант"},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countUnknownsInWord(tt.args.word); got != tt.want {
				t.Errorf("countUnknownsInWord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortByUnknownsInWord(t *testing.T) {
	setUpTestProperties(map[string]string{
		propGuessUnknownMarker: "???",
	})

	type args struct {
		words []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Test01",
			args: args{[]string{"??????ан -> аргирофан", "??????ан -> аркан", "???триант -> репатриант", "??????ан -> Атлантический океан", "???триант -> экспатриант"}},
			want: []string{"???триант -> репатриант", "???триант -> экспатриант", "??????ан -> аргирофан", "??????ан -> аркан", "??????ан -> Атлантический океан"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortByUnknownsInWord(tt.args.words); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortByUnknownsInWord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPossibleVerbs(t *testing.T) {
	got := getPossibleVerbs()

	foundSolve := false
	re := regexp.MustCompile("^[A-Za-z]+$")
	for _, v := range got {
		if v == "solve" {
			foundSolve = true
		}

		if !re.MatchString(v) {
			t.Errorf("getPossibleVerbs() verb should avoid using special chars: %s", v)
		}
	}

	if !foundSolve {
		t.Errorf("getPossibleVerbs() doesn't contain solve")
	}
}
