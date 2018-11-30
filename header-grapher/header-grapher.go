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
	isStruct     bool
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
var variableRegex = regexp.MustCompile(`(?ms)^.?([^\s]+) ?([^\s]+) ?([^\s]+);`)
var bracketRegex = regexp.MustCompile(`(?ms)\[(.*?)\]`)
var isStructVarRegex = regexp.MustCompile(`(?ms)^.?struct ?([^\s]+) ?([^\s]+)[;\[]`)

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

func printVar(v *variable_node) {
	fmt.Println(v.variableName + " : " + v.variableType)
	fmt.Println("Struct: ", v.isStruct)
	fmt.Println("Array Depth: ", v.arrayDepth)
	fmt.Println("Array Params: ", v.arrayDepends)
	fmt.Println()
}

func matchStructs(fileContents string) {
	tFile := enumRegex.ReplaceAllLiteralString(fileContents, "")
	names := structNameRegex.FindAllString(tFile, -1)
	structs := structRegex.FindAllStringSubmatch(tFile, -1)
	for index, str := range names {
		str = strings.Replace(str, "\n", "", -1)
		str = strings.TrimSpace(str)

		tmpStruct := new(struct_node)

		for _, tmp := range strings.Split(structs[index][1], "\n") {
			tmpVar := new(variable_node)

			arrayParams := bracketRegex.FindAllString(tmp, -1)
			tmp = bracketRegex.ReplaceAllLiteralString(tmp, "")

			variable := variableRegex.FindString(tmp)
			if len(variable) == 0 {
				continue
			}

			tmpVar.arrayDepends = arrayParams
			tmpVar.arrayDepth = len(arrayParams)

			tmpVar.isStruct = isStructVarRegex.MatchString(variable)
			variables := strings.Split(variable, " ")
			if tmpVar.isStruct {
				tmpVar.variableType = variables[1]
				tmpVar.variableName = variables[2]
			} else {
				tmpVar.variableType = variables[0]
				tmpVar.variableName = variables[1]
			}

			tmpStruct.vars = append(tmpStruct.vars, tmpVar)
		}
		gStructs = append(gStructs, tmpStruct)
	}

}

func RunGrapher(file, tool string) bool {
	fmt.Println("Running grapher")
	bFile, _ := ioutil.ReadFile(file)
	sFile := string(bFile)
	matchStructs(prepareFile(sFile))
	for _, str := range gStructs {
		for _, v := range str.vars {
			printVar(v)
		}
	}
	return true
}
