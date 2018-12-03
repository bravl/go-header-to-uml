package header_grapher

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
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
	variableType string
	vars         []*variable_node
}

type ParserGrapher struct {
	gStructs []*struct_node
}

var commentRegex = regexp.MustCompile(`(?ms)\/\*(.*?)\*\/|\/\/(.*?).?^`)
var structRegex = regexp.MustCompile(`(?ms)^ ?struct .*?\{(.*?)};`)
var structNameRegex = regexp.MustCompile(`(?ms)^ ?struct ([^\s]+).?{`)
var enumRegex = regexp.MustCompile(`(?ms)^ ?struct ([^\s]+) ?\{.?enum ([^\s]+) ?{(.*?)};.?};`)
var enumNameRegex = regexp.MustCompile(`(?ms)^ ?enum(.*?)\{`)
var variableRegex = regexp.MustCompile(`(?ms)^.?([^\s]+) ?([^\s]+) ?([^\s]+);`)
var bracketRegex = regexp.MustCompile(`(?ms)\[(.*?)\]`)
var isStructVarRegex = regexp.MustCompile(`(?ms)^.?struct ?([^\s]+) ?([^\s]+)[;\[]`)

const plantUMLClassText string = `class PLACEHOLDER1 {
PLACEHOLDER2
}`

const plantUMLVarLinkText string = `note bottom of PLACEHOLDER1 : XD array \n PLACEHOLDER2`

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

func (pg *ParserGrapher) matchStructs(fileContents string) {
	tFile := enumRegex.ReplaceAllLiteralString(fileContents, "")
	names := structNameRegex.FindAllString(tFile, -1)
	structs := structRegex.FindAllStringSubmatch(tFile, -1)
	for index, str := range names {
		str = strings.Replace(str, "\n", "", -1)
		str = strings.TrimSpace(str)

		tmpStruct := new(struct_node)
		tmpStruct.variableType = structNameRegex.FindStringSubmatch(str)[1]

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

			tmpVar.variableName = strings.Replace(tmpVar.variableName, ";", "", -1)
			tmpStruct.vars = append(tmpStruct.vars, tmpVar)
		}
		pg.gStructs = append(pg.gStructs, tmpStruct)
	}

}

func (pg *ParserGrapher) linkNodesPlantUML(file *os.File) {
	if file == nil {
		return
	}

	var tmpString string = ""

	for _, str := range pg.gStructs {
		for _, v := range str.vars {
			if v.isStruct {
				tmpString += str.variableType + " --> " + v.variableType + "\n"
				if v.arrayDepth != 0 {
					tmp := strings.Replace(plantUMLVarLinkText, "PLACEHOLDER1", v.variableType, -1)
					tmp = strings.Replace(tmp, "X", strconv.Itoa(v.arrayDepth), -1)
					tmp = strings.Replace(tmp, "PLACEHOLDER2", strings.Join(v.arrayDepends, " "), -1)
					tmpString += tmp + "\n"
				}
			}
		}
	}
	file.WriteString(tmpString + "\n")
}

func (pg *ParserGrapher) RunGrapher(outFile, tool string) bool {
	if outFile == "none" {
		return false
	}

	output, _ := os.Create(outFile)

	for _, str := range pg.gStructs {
		tmpString := strings.Replace(plantUMLClassText, "PLACEHOLDER1", str.variableType, -1)
		var tmp string = ""

		for _, v := range str.vars {
			tmp += v.variableName + " : " + v.variableType
			for i := 0; i < v.arrayDepth; i++ {
				tmp += "[]"
			}
			tmp += "\n"
		}
		tmpString = strings.Replace(tmpString, "PLACEHOLDER2", tmp, -1)
		output.WriteString(tmpString + "\n")
	}
	pg.linkNodesPlantUML(output)

	return true
}

func (pg *ParserGrapher) RunParser(inFile string) bool {
	if inFile == "none" {
		return false
	}
	fmt.Println("Running grapher")
	bFile, _ := ioutil.ReadFile(inFile)
	sFile := string(bFile)
	pg.matchStructs(prepareFile(sFile))

	return true
}
