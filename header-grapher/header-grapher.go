package header_grapher

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const nodeGraphDataGraphviz string = `"PLACEHOLDER1" [
		label = "PLACEHOLDER2" 
		shape = "record"
	];`
const nodeIndexStringGraphviz string = "<fX> "
const nodeIndexString2Graphviz string = ":fX"

const nodeGraphDataPlantUML string = `class PLACEHOLDER1 {
	PLACEHOLDER2
}`
const nodeGraphNotePlantUML string = `note bottom of PLACEHOLDER1 : XD array \n PLACEHOLDER2`

type grapher_node struct {
	index           int
	arrayDepth      int
	arrayParams     []string
	variableType    string
	variableName    string
	variableComment string
	isStruct        bool
	isArray         bool
	leafs           []*grapher_node
}

var multiLineDefine = false
var multiLineComment = false
var currentNode *grapher_node
var nodes []*grapher_node

func isComment(line string) bool {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) < 2 {
		return true
	}
	multiLineComment = false
	if trimmed[:2] == "/*" && !strings.Contains(trimmed, "*/") {
		multiLineComment = true
		return true
	}
	if trimmed[:2] == "/*" || trimmed[:2] == "//" {
		return true
	}
	return false
}

func isDefine(line string) bool {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) == 0 {
		return false
	}
	if trimmed[0] == '#' || multiLineDefine == true {
		multiLineDefine = false
		if trimmed[len(trimmed)-1] == '\\' {
			multiLineDefine = true
		}
		return true
	}
	multiLineDefine = false
	return false
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func retrieveTopLevelStruct(line string) *grapher_node {
	node := new(grapher_node)
	line = strings.Replace(line, "\t", " ", -1)
	line = strings.TrimSpace(line)
	line = standardizeSpaces(line)
	slices := strings.Split(line, " ")
	if slices[0] == "struct" {
		node.variableType = slices[1]
		node.isStruct = true
	} else {
		return nil
	}
	return node
}

func retrieveInfoFromLine(line string) *grapher_node {
	node := new(grapher_node)
	line = strings.Replace(line, "\t", " ", -1)
	line = strings.TrimSpace(line)
	line = standardizeSpaces(line)
	if isComment(line) {
		return nil
	}
	slices := strings.Split(line, " ")
	if slices[0] == "struct" {
		node.variableType = slices[1]
		node.variableName = slices[2]
		node.isArray, node.arrayDepth, node.arrayParams = isArray(node.variableName)
		node.variableName = strings.Split(node.variableName, "[")[0]
		node.variableName = strings.Replace(node.variableName, ";", "", -1)
		node.isStruct = true

		if len(slices) >= 4 {
			if isComment(slices[3]) {
				node.variableComment = strings.Join(slices[3:], " ")
			}
		}
	} else {
		node.variableType = slices[0]
		node.variableName = slices[1]
		node.isArray, node.arrayDepth, node.arrayParams = isArray(node.variableName)
		node.variableName = strings.Split(node.variableName, "[")[0]
		node.variableName = strings.Replace(node.variableName, ";", "", -1)
		if len(slices) >= 3 {
			if isComment(slices[2]) {
				node.variableComment = strings.Join(slices[2:], " ")
			}
		}
		node.isStruct = false
	}
	return node
}

func isArray(line string) (bool, int, []string) {
	if strings.Contains(line, "[") {
		count := strings.Count(line, "[")
		var re = regexp.MustCompile(`(?m)(\[.*?\])`)
		matches := re.FindAllString(line, -1)
		for i, match := range matches {
			match = strings.Replace(match, "[", "", -1)
			match = strings.Replace(match, "]", "", -1)
			matches[i] = match
		}
		return true, count, matches
	}
	return false, 0, nil
}

func isStruct(line string) bool {
	trimmed := strings.TrimSpace(line)
	if strings.Contains(trimmed, "struct") && trimmed[len(trimmed)-1] == '{' {
		return true
	}
	return false
}

func isEnum(line string) bool {
	trimmed := strings.TrimSpace(line)
	if strings.Contains(trimmed, "enum") && strings.Contains(trimmed, "{") {
		return true
	}
	return false
}

func isEndOfComment(line string) bool {
	trimmed := strings.TrimSpace(line)
	if strings.Contains(trimmed, "*/") {
		return true
	}
	return false
}

func printNodeGraphviz(node *grapher_node) {
	if node == nil {
		return
	}
	tmp := strings.Replace(nodeIndexStringGraphviz, "X", strconv.Itoa(node.index), -1) + node.variableType + " | "
	for _, leaf := range node.leafs {
		tmp += strings.Replace(nodeIndexStringGraphviz, "X", strconv.Itoa(leaf.index), -1) + leaf.variableName + " | "
	}
	tmpstr := strings.Replace(nodeGraphDataGraphviz, "PLACEHOLDER1", node.variableType, -1)
	tmpstr = strings.Replace(tmpstr, "PLACEHOLDER2", tmp, -1)
	fmt.Println(tmpstr)
}

func printNodePlantUML(node *grapher_node) {
	if node == nil {
		return
	}
	if strings.Contains(node.variableType, "Enum") {
		return
	}
	var tmp string = ""
	for _, leaf := range node.leafs {
		tmp += leaf.variableType + " : " + leaf.variableName + " \n "
	}
	tmpstr := strings.Replace(nodeGraphDataPlantUML, "PLACEHOLDER1", node.variableType, -1)
	tmpstr = strings.Replace(tmpstr, "PLACEHOLDER2", tmp, -1)
	fmt.Println(tmpstr)
}

func linkNodesGraphviz(node *grapher_node) {
	if node == nil {
		return
	}
	for _, leaf := range node.leafs {
		if leaf.isStruct {
			fmt.Println("\"" + node.variableType + "\"" + strings.Replace(nodeIndexString2Graphviz, "X", strconv.Itoa(leaf.index), -1) + " -> " + "\"" + leaf.variableType + "\"" + strings.Replace(nodeIndexString2Graphviz, "X", strconv.Itoa(0), -1))
		}
	}
}

func linkNodesPlantUML(node *grapher_node) {
	if node == nil {
		return
	}
	for _, leaf := range node.leafs {
		if leaf.isStruct {
			fmt.Println(node.variableType + " --> " + leaf.variableType)
			if leaf.isArray {
				tmp := strings.Replace(nodeGraphNotePlantUML, "PLACEHOLDER1", leaf.variableType, -1)
				tmp = strings.Replace(tmp, "X", strconv.Itoa(leaf.arrayDepth), -1)
				tmp = strings.Replace(tmp, "PLACEHOLDER2", strings.Join(leaf.arrayParams, " "), -1)
				fmt.Println(tmp)
			}
		}
	}
}
func isEndOfStruct(line string) bool {
	trimmed := strings.TrimSpace(line)
	if strings.Contains(trimmed, "};") {
		//printNodeGraphviz(currentNode)
		nodes = append(nodes, currentNode)
		currentNode = nil
		return true
	}
	return false
}

func parseStruct(scanner *bufio.Scanner, isAnEnum bool) bool {
	var index int = 1
	for scanner.Scan() {
		line := scanner.Text()
		if isComment(line) {
			if !multiLineComment {
				continue
			}
			for scanner.Scan() {
				if isEndOfComment(scanner.Text()) {
					break
				}
			}
		}
		if isDefine(line) {
			continue
		}

		if isEndOfStruct(line) {
			break
		}
		if isEnum(line) {
			parseStruct(scanner, true)
		} else if !isAnEnum {
			node := retrieveInfoFromLine(line)
			if node == nil {
				continue
			}
			node.index = index
			index++
			if currentNode != nil {
				currentNode.leafs = append(currentNode.leafs, node)
			}
		}
	}
	return true
}

func RunGrapher(fileName string, grapherTool string) {
	log.Println("Running Grapher on", fileName)
	if fileName == "none" || fileName == "" {
		log.Fatal("No input file provided")
		return
	}
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("Could not open file")
		return
	}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if isComment(line) {
			continue
		}
		if isDefine(line) {
			continue
		}
		if isStruct(line) {
			node := retrieveTopLevelStruct(line)
			if node != nil {
				currentNode = node
			}
			parseStruct(scanner, false)
		}
		if isEnum(line) {
			parseStruct(scanner, true)
		}
	}
	if grapherTool == "graphviz" {
		// Create boxes (graphviz)
		for _, node := range nodes {
			printNodeGraphviz(node)
		}

		// Create links (graphviz)

		for _, node := range nodes {
			linkNodesGraphviz(node)
		}
	} else if grapherTool == "plantuml" {
		for _, node := range nodes {
			printNodePlantUML(node)
		}
		for _, node := range nodes {
			linkNodesPlantUML(node)
		}

	}
}
