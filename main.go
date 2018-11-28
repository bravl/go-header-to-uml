package main

import (
	"./header-grapher"
	"flag"
)

func main() {
	inputFile := flag.String("in", "none", "The file that will be parsed")
	grapherTool := flag.String("tool", "graphviz", "The tool used to graph the diagram")
	flag.Parse()

	header_grapher.RunGrapher(*inputFile, *grapherTool)
}
