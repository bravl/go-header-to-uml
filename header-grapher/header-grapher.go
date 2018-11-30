package header_grapher

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

type variable_node struct {
	variableType string
	variableName string
	arrayDepends []string
	arrayDepth   int
	isArray      bool
}
type struct_node struct {
	varType string
	vars    []*variable_node
}

var gStructs []*struct_node

var commentRegex = regexp.MustCompile(`(?ms)\/\*(.*?)\*\/|\/\/(.*?).?^`)
var structRegex = regexp.MustCompile(`(?ms)^ ?struct .*?\{(.*?)};`)
var structNameRegex = regexp.MustCompile(`(?ms)^ ?struct ([^\s]+).?{`)
var enumRegex = regexp.MustCompile(`(?ms)^ ?struct ([^\s]+) ?\{.?enum ([^\s]+) ?{(.*?)};.?};`)
var enumNameRegex = regexp.MustCompile(`(?ms)^ ?enum(.*?)\{`)
var variableRegex = regexp.MustCompile(`(?ms)^ ?([^\s]+) ?([^\s]+);`)
var bracketRegex = regexp.MustCompile(`(?ms)\[(.*?)\]`)

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func prepareFile(fileContents string) string {
	tFile := commentRegex.ReplaceAllLiteralString(fileContents, "\n")
	var lines []string
	for _, match := range strings.Split(tFile, "\n") {
		match = standardizeSpaces(match)
		if len(match) == 0 {
			continue
		}
		lines = append(lines, match)
	}
	return strings.Join(lines, "\n")
}

func matchStructs(fileContents string) {
	tFile := enumRegex.ReplaceAllLiteralString(fileContents, "")
	names := structNameRegex.FindAllString(tFile, -1)
	structs := structRegex.FindAllStringSubmatch(tFile, -1)
	fmt.Println(len(names))
	fmt.Println(len(structs))
	for index, str := range names {
		str = strings.Replace(str, "\n", "", -1)
		str = strings.TrimSpace(str)
		fmt.Println(strings.Replace(str, "{", "", -1))

		tmpStruct := new(struct_node)

		for _, tmp := range strings.Split(structs[index][1], "\n") {
			arrayParams := bracketRegex.FindAllString(tmp, -1)
			fmt.Println(arrayParams)
		}
		gStructs = append(gStructs, tmpStruct)
	}

}

func RunGrapher(file, tool string) bool {
	fmt.Println("Running grapher")
	bFile, _ := ioutil.ReadFile(file)
	sFile := string(bFile)
	//fmt.Println(prepareFile(sFile))
	matchStructs(prepareFile(sFile))
	return true
}
