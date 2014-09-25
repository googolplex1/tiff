// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"text/scanner"
)

type Type struct {
	TypeName string
	FileName string
	TypeList []string
	MapCode  string
}

func main() {
	var types = []Type{
		Type{
			TypeName: "CompressType",
			FileName: "compress_type.go",
		},
		Type{
			TypeName: "DataType",
			FileName: "data_type.go",
		},
		Type{
			TypeName: "TagType",
			FileName: "tag_type.go",
		},
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, `
// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by gen_helper.go,
// DO NOT EDIT!!!

package tiff

`[1:])

	for _, v := range types {
		v.Init()
		v.GenMapCode()

		fmt.Fprintf(&buf, "%s\n", v.MapCode)
	}

	data, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("types_table.go", data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Type) Init() {
	f, err := os.Open(p.FileName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var s scanner.Scanner
	var typeMap = make(map[string]bool)

	s.Init(f)
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		if tok&scanner.ScanIdents != 0 {
			if strings.HasPrefix(s.TokenText(), p.TypeName+"_") {
				typeMap[s.TokenText()] = true
			}
		}
	}

	for k, _ := range typeMap {
		p.TypeList = append(p.TypeList, k)
	}
	sort.Strings(p.TypeList)
}

func (p *Type) GenMapCode() {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "var _%sTable = map[%s]string {\n", p.TypeName, p.TypeName)
	for _, s := range p.TypeList {
		fmt.Fprintf(&buf, "\t%s: `%s`,\n", s, s)
	}
	fmt.Fprintf(&buf, "}\n")
	p.MapCode = buf.String()
}