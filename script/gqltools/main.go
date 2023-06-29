// GqlInOne is a tool to merge all gqlgen schema files into one file
package main

import (
	"flag"
	"github.com/99designs/gqlgen/codegen/config"
	"log"
	"os"
	"path/filepath"
)

var (
	gqlgenFile = flag.String("gqlgen", "./codegen/gqlgen/gqlgen.yaml", "gqlgen config file")
	targetFile = flag.String("target", "./test/tmp/all.graphql", "target file")
)

func main() {
	flag.Parse()
	cfg, err := config.LoadConfig(*gqlgenFile)
	if err != nil {
		log.Fatal("failed to load config", err.Error())
	}
	if err := os.MkdirAll(filepath.Dir(*targetFile), os.ModePerm); err != nil {
		log.Fatal(err)
	}
	// clear file
	err = os.WriteFile(*targetFile, []byte(""), os.ModePerm)
	if err != nil {
		log.Fatal("failed to write file", err.Error())
	}

	f, err := os.OpenFile(*targetFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("failed to write file", err.Error())
	}
	defer f.Close()

	for _, source := range cfg.Sources {
		// write to file
		_, err = f.WriteString(source.Input)
		if err != nil {
			log.Fatal("failed to write file", err.Error())
		}
	}
}
