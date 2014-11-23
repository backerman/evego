/*
Copyright © 2014 Brad Ackerman.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

// Package main is just a test driver - currently pointless.
package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/backerman/evego/pkg/dbaccess"
)

func readItemFile(filename string) []string {
	inFile, err := os.Open(filename)
	defer inFile.Close()
	if err != nil {
		log.Fatalf("Unable to open input file: %v", err)
	}
	inReader := csv.NewReader(inFile)
	inReader.Comma = '\t'
	recs, err := inReader.ReadAll()
	var items []string
	if err != nil {
		log.Fatalf("Unable to read input file: %v", err)
	}
	for _, rec := range recs {
		items = append(items, rec[0])
	}
	log.Printf("%v items read.", len(items))

	return items
}

func main() {
	fmt.Println("This is a test.")
	db := dbaccess.SQLiteDatabase("/Users/bsa3/Downloads/sqlite-latest.sqlite")
	defer db.Close()
	triBistot, _ := db.ItemForName("Triclinic Bistot")
	catTree, _ := db.MarketGroupForItem(triBistot)
	log.Printf("%v", catTree)
	glitter, _ := db.ItemForName("Dark Glitter")
	catTree, _ = db.MarketGroupForItem(glitter)
	log.Printf("%v", catTree)

}
