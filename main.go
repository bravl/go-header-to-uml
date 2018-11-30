package main

import (
	"./header-grapher"
	"flag"
)

func main() {
	pg := new(header_grapher.ParserGrapher)
	inputFile := flag.String("in", "none", "The file that will be parsed")
	outFile := flag.String("out", "none", "Graph text file")
	grapherTool := flag.String("tool", "plantuml", "The tool used to graph the diagram")
	flag.Parse()

	pg.RunParser(*inputFile)
	pg.RunGrapher(*outFile, *grapherTool)
}
