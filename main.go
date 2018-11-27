package main

import (
	"./header-grapher"
	"flag"
	"log"
)

func main() {
	inputFile := flag.String("in", "none", "The file that will be parsed")
	grapherTool := flag.String("tool", "graphviz", "The tool used to graph the diagram")
	flag.Parse()

	log.Println("Input file: ", *inputFile)
	header_grapher.RunGrapher(*inputFile, *grapherTool)
}
