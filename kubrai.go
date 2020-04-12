package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// verbs
const (
	vAdd        = "add"
	vAddBoth    = "addboth"
	vGuess      = "guess"
	vHint       = "hint"
	vInteract   = "interact"
	vRemove     = "remove"
	vRemoveBoth = "removeboth"
	vSearchDict = "searchdict"
	vSolve      = "solve"
	vUndo       = "undo"
	vView       = "view"
)

const propertyDir = "./properties"

// properties
const (
	propAddAutoBothMaxlen           = "AddAutoBothMaxlen"
	propAddAutoLowercase            = "AddAutoLowercase"
	propAddValMayEqualKey           = "AddValMayEqualKey"
	propAssocFileKeySeparator       = "AssocFileKeySeparator"
	propAssocFileLocation           = "AssocFileLocation"
	propAssocFileValSeparator       = "AssocFileValSeparator"
	propGuessExplainResults         = "GuessExplainResults"
	propGuessMaxResults             = "GuessMaxResults"
	propGuessUnknownMarker          = "GuessUnknownMarker"
	propGuessUnknownsLimit          = "GuessUnknownsLimit"
	propSearchDictDefaultMaxResults = "SearchDictDefaultMaxResults"
	propSolveAutoGuess              = "SolveAutoGuess"
	propSolveMaxResults             = "SolveMaxResults"
	propSolveKubrayaSeparator       = "SolveKubrayaSeparator"
	propDictsDir                    = "DictsDir"
	propDictsExt                    = "DictsExt"
)

func findExactVerb(args []string) string {
	possibleVerbs := getPossibleVerbs()

	for _, arg := range args {
		for _, possibleVerb := range possibleVerbs {
			if arg == possibleVerb {
				return arg
			}
		}
	}

	return ""
}

func getPossibleVerbs() [11]string {
	return [...]string{vAdd, vAddBoth, vRemove, vRemoveBoth, vView, vSearchDict, vSolve, vGuess, vHint, vInteract, vUndo}
}

func guessVerb(args []string) string {
	var verb string

	verb = guessInteractVerb(args)
	if verb != "" {
		return verb
	}

	verb = guessSolveVerb(args)
	if verb != "" {
		return verb
	}

	verb = guessViewVerb(args)
	if verb != "" {
		return verb
	}

	return ""
}

func guessInteractVerb(args []string) string {
	if len(args) == 0 {
		return vInteract
	}

	return ""
}

func guessSolveVerb(args []string) string {
	if len(args) == 1 && isKubraya(args[0]) {
		return vSolve
	}

	return ""
}

func guessViewVerb(args []string) string {
	if len(args) == 1 && !isKubraya(args[0]) {
		return vView
	}

	return ""
}

func isKubraya(str string) bool {
	return strings.Contains(str, "_")
}

func parseVerb(args []string) string {
	verb := findExactVerb(args)

	if verb == "" {
		verb = guessVerb(args)
	}

	return verb
}

func filterOut(word string, words []string) []string {
	n := 0
	for _, w := range words {
		if w != word {
			words[n] = w
			n++
		}
	}

	return words[:n]
}

func extractNouns(verb string, args []string) []string {
	return filterOut(verb, args)
}

func warnFileReadError(err error) {
	log.Println(fmt.Errorf("WARN : File read error: %s", err))
}

func getStringProperty(name string) string {
	bs, err := ioutil.ReadFile(propertyDir + "/" + name)
	if err != nil {
		warnFileReadError(err)
		return ""
	}
	str := string(bs)

	return str
}

func warnConvertToIntError(err error) {
	log.Println(fmt.Errorf("WARN : Convert to int error: %s", err))
}

func getByteProperty(name string) byte {
	return byte(getIntProperty(name))
}

func getIntProperty(name string) int {
	num, err := strconv.ParseInt(getStringProperty(name), 10, 64)
	if err != nil {
		warnConvertToIntError(err)
		return 0
	}
	return int(num)
}

func warnNonBool(b string) {
	log.Println(fmt.Errorf("WARN : Non-bool: %s", b))
}

func getBoolProperty(name string) bool {
	b := strings.Trim(getStringProperty(name), " ")
	if b == "ON" {
		return true
	} else if b == "OFF" {
		return false
	} else {
		warnNonBool(b)
		return false
	}
}

func checkError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func loadAssoc(assocFile string) map[string][]string {
	f, err := os.Open(assocFile)
	checkError(err)
	defer f.Close()

	assoc := make(map[string][]string)

	keySep := getStringProperty(propAssocFileKeySeparator)
	valSep := getStringProperty(propAssocFileValSeparator)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		keyVals := strings.Split(scanner.Text(), keySep)
		key := keyVals[0]
		vals := strings.Split(keyVals[1], valSep)
		assoc[key] = vals
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return assoc
}

func backupName(file string, backupNum int) string {
	return file + "." + strconv.Itoa(backupNum) + ".bak"
}

func rotateBackups(file string) {
	lastBakFile := backupName(file, 10)
	os.Remove(lastBakFile)

	for i := 9; i > 0; i-- {
		bacFile := backupName(file, i)
		nextBacFile := backupName(file, i+1)
		os.Rename(bacFile, nextBacFile)
	}
}

func backupFile(file string) {
	rotateBackups(file)
	os.Rename(file, backupName(file, 1))
}

func canonize(word string) string {
	return strings.Trim(strings.ToUpper(word), " ")
}

func buildAssocString(word string, assocSingle []string) string {
	return word + getStringProperty(propAssocFileKeySeparator) + strings.Join(assocSingle, getStringProperty(propAssocFileValSeparator))
}

func saveAssoc(assocFile string, assoc map[string][]string) {
	backupFile(assocFile)

	f, err := os.Create(assocFile)
	checkError(err)
	defer f.Close()

	for k, v := range assoc {
		f.WriteString(buildAssocString(k, v) + "\n")
	}
}

func saveDefaultAssoc(assoc map[string][]string) {
	saveAssoc(getStringProperty(propAssocFileLocation), assoc)
}

func loadDefaultAssoc() map[string][]string {
	return loadAssoc(getStringProperty(propAssocFileLocation))
}

func addBeforeFirstLonger(s string, slc []string) []string {
	longer := -1
	for i, w := range slc {
		if w == s {
			return slc
		}

		if len(w) > len(s) {
			longer = i
			break
		}
	}

	if longer == -1 {
		return append(slc, s)
	}

	res := make([]string, len(slc)+1)
	copy(res, slc[:longer])
	res[longer] = s
	copy(res[longer+1:], slc[longer:])
	return res
}

func runAdd(a string, b string) []string {
	assoc := loadDefaultAssoc()

	if getBoolProperty(propAddAutoLowercase) {
		a = strings.ToLower(a)
		b = strings.ToLower(b)
	}

	if a != b || a == b && getBoolProperty(propAddValMayEqualKey) {
		assoc[a] = addBeforeFirstLonger(b, assoc[a])
		saveDefaultAssoc(assoc)
	}

	return assoc[a]
}

func runAddBoth(a string, b string) [2][]string {
	var res [2][]string
	res[0] = runAdd(a, b)
	res[1] = runAdd(b, a)
	return res
}

func runSmartAdd(a string, b string) ([2][]string, uint8) {
	if len([]rune(a)) <= getIntProperty(propAddAutoBothMaxlen) {
		return runAddBoth(a, b), 2
	}
	var res [2][]string
	res[0] = runAdd(a, b)
	return res, 1
}

func removeByValue(s string, slc []string) []string {
	for k, v := range slc {
		if v == s {
			return append(slc[:k], slc[k+1:]...)
		}
	}

	return slc
}

func runRemove(a string, b string) []string {
	assoc := loadDefaultAssoc()
	if assoca, ok := assoc[a]; ok {
		assoca = removeByValue(b, assoca)
		if len(assoca) == 0 {
			delete(assoc, a)
		} else {
			assoc[a] = assoca
		}
		saveDefaultAssoc(assoc)
		return assoca
	}

	return []string{}
}

func runRemoveBoth(a string, b string) [2][]string {
	var res [2][]string
	res[0] = runRemove(a, b)
	res[1] = runRemove(b, a)
	return res
}

func runView(a string) []string {
	assoc := loadDefaultAssoc()
	if assoca, ok := assoc[a]; ok {
		return assoca
	}
	return []string{}
}

func readDirNames(dir string) []string {
	f, err := os.Open(dir)
	defer f.Close()
	checkError(err)

	list, err := f.Readdirnames(-1)
	checkError(err)

	return list
}

func readFileToSlice(file string, cap int) []string {
	f, err := os.Open(file)
	checkError(err)
	defer f.Close()

	slc := make([]string, 0, cap)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		slc = append(slc, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return slc
}

func loadDicts() map[string][]string {
	dictsDir := getStringProperty(propDictsDir)
	dictsExt := getStringProperty(propDictsExt)
	lenExt := len(dictsExt)
	list := readDirNames(dictsDir)

	dicts := make(map[string][]string)
	for _, n := range list {
		if n[len(n)-lenExt:] != dictsExt {
			continue
		}

		dicts[n] = readFileToSlice(dictsDir+"/"+n, 200000)
	}

	return dicts
}

func splitKubraya(kubraya string) []string {
	if getBoolProperty(propAddAutoLowercase) {
		kubraya = strings.ToLower(kubraya)
	}

	kubSep := getStringProperty(propSolveKubrayaSeparator)
	return strings.Split(kubraya, kubSep)
}

func nextMultiDimValue(counter []int, maxCounter []int) ([]int, bool) {
	if len(counter) != len(maxCounter) {
		log.Println(fmt.Errorf("WARN : Counter length is incorrect: %d != %d", len(counter), len(maxCounter)))
	}

	ok := false
	for i, c := range counter {
		if c != maxCounter[i] {
			counter[i] = c + 1
			ok = true
			break
		} else {
			counter[i] = 0
		}
	}

	return counter, ok
}

func combinations(items [][]string) [][]string {
	lenitems := len(items)

	maxCounter := make([]int, lenitems)
	for i, item := range items {
		maxCounter[i] = len(item) - 1
	}

	ok := true
	counter := make([]int, lenitems)
	combs := make([][]string, 0, lenitems*10)
	for true {
		comb := make([]string, lenitems)
		for i, j := range counter {
			comb[i] = items[i][j]
		}
		combs = append(combs, comb)

		counter, ok = nextMultiDimValue(counter, maxCounter)
		if !ok {
			break
		}
	}

	return combs
}

func runSearchDict(word string, maxResults int) map[string]int {
	// improve it with fuzzy search
	dicts := loadDicts()

	results := make(map[string]int)
	if maxResults == 0 {
		return results
	}
	for dictName, dict := range dicts {
		for line, dictWord := range dict {
			if word == dictWord {
				results[dictName] = line
				if len(results) == maxResults {
					return results
				}
			}
		}
	}

	return results
}

func searchDictByRegexpGetWords(re *regexp.Regexp, maxResults int) []string {
	dicts := loadDicts()

	results := make([]string, 0, maxResults)
	if maxResults == 0 {
		return results
	}
	for _, dict := range dicts {
		for _, dictWord := range dict {
			if re.MatchString(dictWord) {
				results = append(results, dictWord)
				if len(results) == maxResults {
					return results
				}
			}
		}
	}

	return results
}

func mapKeys(m map[string]bool) []string {
	keys := make([]string, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

func buildKubAssoc(kubraya string) ([][]string, bool) {
	kubParts := splitKubraya(kubraya)

	kubAssoc := make([][]string, len(kubParts))
	for i, part := range kubParts {
		kubAssoc[i] = runView(part)
		if len(kubAssoc[i]) == 0 {
			return [][]string{}, false
		}
	}

	return kubAssoc, true
}

func buildKubAssocComplete(kubraya string) ([][]string, bool) {
	kubParts := splitKubraya(kubraya)

	complete := true
	kubAssoc := make([][]string, len(kubParts))
	for i, part := range kubParts {
		kubAssoc[i] = runView(part)
		if len(kubAssoc[i]) == 0 {
			complete = false
			kubAssoc[i] = []string{}
		}
	}

	return kubAssoc, complete
}

func combsOfKubraya(kubraya string) ([][]string, bool) {
	if kubAssoc, ok := buildKubAssoc(kubraya); ok {
		combs := combinations(kubAssoc)
		return combs, true
	}

	return [][]string{}, false
}

func runSolve(kubraya string) ([]string, bool) {
	maxResults := getIntProperty(propSolveMaxResults)
	results := make(map[string]bool)

	combs, ok := combsOfKubraya(kubraya)
	if !ok {
		return []string{}, false
	}
	for _, comb := range combs {
		word := strings.Join(comb, "")

		res := runSearchDict(word, 1)
		if len(res) > 0 {
			if !results[word] {
				results[word] = true
			}

			if len(results) == maxResults {
				return mapKeys(results), true
			}
		}
	}
	if len(results) > 0 {
		return mapKeys(results), true
	}

	return []string{}, false
}

func allowUnknowns(kubAssoc [][]string) [][]string {
	guessUnknownMarker := getStringProperty(propGuessUnknownMarker)

	for i, list := range kubAssoc {
		if len(list) == 0 {
			kubAssoc[i] = []string{guessUnknownMarker}
		} else {
			kubAssoc[i] = append(kubAssoc[i], guessUnknownMarker)
		}
	}

	return kubAssoc
}

func wordGuessToRegexp(wordGuess string) string {
	guessUnknownMarker := getStringProperty(propGuessUnknownMarker)

	return "^" + strings.ReplaceAll(wordGuess, guessUnknownMarker, ".+") + "$"
}

func countValsInSlice(slc []string, val string) int {
	counter := 0
	for _, v := range slc {
		if v == val {
			counter++
		}
	}

	return counter
}

func isCombGuessable(comb []string) bool {
	guessUnknownMarker := getStringProperty(propGuessUnknownMarker)
	guessUnknownsLimit := getIntProperty(propGuessUnknownsLimit)

	unk := countValsInSlice(comb, guessUnknownMarker)

	// not really good. logic about guessability needs to be injected from params
	// also maybe need smarter criteria - 3 of out 5 unknown is fine
	// rely on percentage??
	return unk <= guessUnknownsLimit && unk < len(comb)
}

func filterGuessableCombs(combs [][]string) [][]string {
	res := [][]string{}
	for _, comb := range combs {
		if isCombGuessable(comb) {
			res = append(res, comb)
		}
	}

	return res
}

func countUnknowns(comb []string) int {
	guessUnknownMarker := getStringProperty(propGuessUnknownMarker)
	return countValsInSlice(comb, guessUnknownMarker)
}

func sortCombsByBestChances(combs [][]string) [][]string {
	sort.Slice(combs, func(i, j int) bool {
		return countUnknowns(combs[i]) < countUnknowns(combs[j])
	})

	return combs
}

func countUnknownsInWord(word string) int {
	guessUnknownMarker := getStringProperty(propGuessUnknownMarker)
	return strings.Count(word, guessUnknownMarker)
}

func sortByUnknownsInWord(words []string) []string {
	sort.Slice(words, func(i, j int) bool {
		return countUnknownsInWord(words[i]) < countUnknownsInWord(words[j])
	})
	return words
}

func runGuess(kubraya string) ([]string, bool) {
	kubAssoc, complete := buildKubAssocComplete(kubraya)
	if complete {
		if res, ok := runSolve(kubraya); ok {
			return res, true
		}
	}

	kubAssoc = allowUnknowns(kubAssoc)
	combs := combinations(kubAssoc)
	combs = filterGuessableCombs(combs)
	combs = sortCombsByBestChances(combs)

	maxResults := getIntProperty(propGuessMaxResults)
	results := make(map[string]bool)
	explains := make(map[string]string)
	guessExplainResults := getBoolProperty(propGuessExplainResults)
	for _, comb := range combs {
		wordGuess := strings.Join(comb, "")
		wordRegexp := wordGuessToRegexp(wordGuess)
		re := regexp.MustCompile(wordRegexp)

		words := searchDictByRegexpGetWords(re, maxResults)
		if len(words) > 0 {
			for _, word := range words {
				if !results[word] {
					if guessExplainResults {
						explains[word] = wordGuess
					}
					results[word] = true
				}
				if len(results) == maxResults {
					break
				}
			}
		}
		if len(results) == maxResults {
			break
		}
	}

	if len(results) > 0 {
		keys := mapKeys(results)
		if guessExplainResults {
			for i, key := range keys {
				keys[i] = explains[key] + " -> " + key
			}
			keys = sortByUnknownsInWord(keys)
		}
		if len(keys) >= maxResults {
			keys = append(keys, "(First "+strconv.Itoa(maxResults)+" shown, more exist)")
		}
		return keys, true
	}

	return []string{}, false
}

func runCommand(verb string, nouns []string) string {
	switch verb {
	case vAdd:
		res, num := runSmartAdd(nouns[0], nouns[1])
		if num == 2 {
			return buildAssocString(nouns[0], res[0]) + "\n" + buildAssocString(nouns[1], res[1])
		}
		return buildAssocString(nouns[0], res[0])
	case vAddBoth:
		res := runAddBoth(nouns[0], nouns[1])
		return buildAssocString(nouns[0], res[0]) + "\n" + buildAssocString(nouns[1], res[1])
	case vGuess:
		if res, ok := runGuess(nouns[0]); ok {
			return strings.Join(res, "\n")
		}
		return "404 NOT FOUND"
	case vRemove:
		res := runRemove(nouns[0], nouns[1])
		return buildAssocString(nouns[0], res)
	case vRemoveBoth:
		res := runRemoveBoth(nouns[0], nouns[1])
		return buildAssocString(nouns[0], res[0]) + "\n" + buildAssocString(nouns[1], res[1])
	case vSearchDict:
		res := runSearchDict(nouns[0], getIntProperty(propSearchDictDefaultMaxResults))
		if len(res) == 0 {
			return "404 NOT FOUND"
		}

		tmp := make([]string, 0, len(res))
		for k, v := range res {
			tmp = append(tmp, k+": "+strconv.Itoa(v))
		}
		return strings.Join(tmp, "\n")
	case vSolve:
		if res, ok := runSolve(nouns[0]); ok {
			return strings.Join(res, "\n")
		}
		return "404 NOT FOUND"
	case vView:
		res := runView(nouns[0])
		return buildAssocString(nouns[0], res)
	default:
		msg := fmt.Sprintf("501 NOT IMPLEMENTED\n%s", verb)
		return msg
	}
}

func main() {
	verb := parseVerb(os.Args[1:])
	if verb == "" {
		answer := "400 BAD REQUEST"
		fmt.Println(answer)
	} else {
		nouns := extractNouns(verb, os.Args[1:])
		answer := runCommand(verb, nouns)
		fmt.Println(answer)
	}
}