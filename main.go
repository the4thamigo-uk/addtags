package main

import (
	"bytes"
	"fmt"
	"github.com/fatih/structtag"
	"github.com/spf13/pflag"
	"github.com/the4thamigo-uk/conflate"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
)

type structTags map[string]fieldTags
type fieldTags map[string]string

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {

	var dir string
	pflag.StringVarP(&dir, "dir", "d", ".", `directory containing go code where you want the tags to be applied`)
	pflag.Parse()

	var files []string
	pflag.StringArrayVarP(&files, "tags", "t", nil, `list of filenames (JSON, YAML or TOML) containing tags to add to struct fields`)
	pflag.Parse()

	if len(files) == 0 {
		return fmt.Errorf("no tags files specified")
	}

	newTags, err := loadConfig(files...)
	if err != nil {
		return err
	}

	err = appendTags(dir, *newTags)
	if err != nil {
		return err
	}

	return nil
}

func loadConfig(files ...string) (*structTags, error) {

	c, err := conflate.FromFiles(files...)
	if err != nil {
		return nil, err
	}

	var newTags structTags
	err = c.Unmarshal(&newTags)
	if err != nil {
		return nil, err
	}
	return &newTags, nil
}

func appendTags(dir string, newTags structTags) error {
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		for fileName, file := range pkg.Files {
			structs := findStructs(file)

			var modified bool
			for strName, str := range structs {
				strTags, ok := newTags[strName]
				if !ok {
					continue
				}
				for _, fld := range str.Fields.List {
					fldTag, ok := strTags[fld.Names[0].Name]
					if !ok {
						continue
					}
					newTag, err := addTag(unquoteTag(fld.Tag.Value), unquoteTag(fldTag))
					if err != nil {
						return err
					}
					fld.Tag.Value = quoteTag(newTag)
					modified = true
				}
			}
			if modified {
				var buf bytes.Buffer
				err = format.Node(&buf, fset, file)
				if err != nil {
					return err
				}
				err = ioutil.WriteFile(fileName, buf.Bytes(), 0)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func unquoteTag(tag string) string {
	if len(tag) > 2 && tag[0] == '`' && tag[len(tag)-1] == '`' {
		return tag[1 : len(tag)-1]
	}
	return tag
}

func quoteTag(tag string) string {
	if len(tag) > 2 && tag[0] == '`' && tag[len(tag)-1] == '`' {
		return tag
	}
	return "`" + tag + "`"
}

func addTag(tag string, newTag string) (string, error) {
	tags, err := structtag.Parse(tag)
	if err != nil {
		return "", err
	}
	newTags, err := structtag.Parse(newTag)
	if err != nil {
		return "", err
	}
	for _, t := range newTags.Tags() {
		tags.Set(t)
	}
	return tags.String(), nil
}

func findStructs(node ast.Node) map[string]*ast.StructType {
	structs := make(map[string]*ast.StructType, 0)
	findStructs := func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.TypeSpec:
			if s, ok := t.Type.(*ast.StructType); ok {
				if t.Name.Name == "CargoEvent" {
					structs[t.Name.Name] = s
				}
			}
		}
		return true
	}
	ast.Inspect(node, findStructs)
	return structs
}
