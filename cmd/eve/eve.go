/*
Copyright © 2014–5 Brad Ackerman.

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

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"text/template"

	"github.com/backerman/evego"
	"github.com/backerman/evego/pkg/cache"
	"github.com/backerman/evego/pkg/dbaccess"
	"github.com/backerman/evego/pkg/eveapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// Register SQLite3 driver for static database export
	_ "github.com/mattn/go-sqlite3"
)

// Commands
var (
	rootCmd = &cobra.Command{
		Use:   "evego",
		Short: "Demo application for evego",
		Long:  "This is a demonstration for the evego library.",
	}
	accountCmd = &cobra.Command{
		Use:   "account",
		Short: "Account commands",
		Run:   callHelp,
	}
	accountListCharsCmd = &cobra.Command{
		Use:   "characters",
		Short: "List characters on an account",
		Run:   listCharacters,
	}
	charCmd = &cobra.Command{
		Use:   "character",
		Short: "Character commands",
		Run:   callHelp,
	}
	charSheetCmd = &cobra.Command{
		Use:   "sheet",
		Short: "Get the character sheet",
		Run:   characterInfo,
	}
	charAssetsCmd = &cobra.Command{
		Use:   "assets",
		Short: "Get the character's assets",
		Run:   characterAssets,
	}
)

var ts *httptest.Server

// Get the XML API key; fail if it hasn't been set.
func getAPIKey() *evego.XMLKey {
	keyID := viper.GetInt("keyID")
	vcode := viper.GetString("vcode")
	if keyID == 0 || vcode == "" {
		log.Fatalf("Error: API key ID and verification code must be set.")
	}
	return &evego.XMLKey{
		KeyID:            keyID,
		VerificationCode: vcode,
	}
}

func getXMLAPI(sdeDatabase evego.Database) evego.XMLAPI {
	var xmlEndpoint string
	// If faking the XML API, set up the fake server.
	fakeResultsFile := viper.GetString("localxml")
	if fakeResultsFile != "" {
		ts = httptest.NewServer(
			http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					respFile, _ := os.Open(fakeResultsFile)
					responseBytes, _ := ioutil.ReadAll(respFile)
					responseBuf := bytes.NewBuffer(responseBytes)
					responseBuf.WriteTo(w)
				}))
		xmlEndpoint = ts.URL
	} else {
		xmlEndpoint = viper.GetString("xmlapi")
	}
	// Using a NilCache just for demo purposes; do not use this in real life.
	return eveapi.XML(xmlEndpoint, sdeDatabase, cache.NilCache())
}

// Get a Database object for the SDE.
func getSDE() evego.Database {
	sdePath := viper.GetString("sdepath")
	if sdePath == "" {
		log.Fatalf("Error: You must specify the SDE file's path.")
	}
	db := dbaccess.SQLDatabase("sqlite3", sdePath)
	return db
}

// Call the help command.
func callHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func listCharacters(cmd *cobra.Command, args []string) {
	xmlKey := getAPIKey()
	sde := getSDE()
	xmlapi := getXMLAPI(sde)
	chars, err :=
		xmlapi.AccountCharacters(xmlKey)
	if err != nil {
		log.Fatalf("Unable to call API: %v", err)
	}
	fmt.Printf("Characters accessible using key %d:\n", xmlKey.KeyID)
	for _, ch := range chars {
		fmt.Printf("  %s (%d)\n", ch.Name, ch.ID)
		fmt.Printf("    Corporation: %s (%d)\n", ch.Corporation, ch.CorporationID)
		fmt.Printf("    Alliance:    %s (%d)\n", ch.Alliance, ch.AllianceID)
	}
}

type charInfoDisplay struct {
	*evego.CharacterSheet
	Skills [][]evego.Skill
}

// partitionSkills accepts a slice of Skills sorted by group and name,
// and returns a slice of slices, each containing one skill group's skills.
func partitionSkills(skills []evego.Skill) [][]evego.Skill {
	result := make([][]evego.Skill, 0, 5)
	if len(skills) == 0 {
		return result
	}
	lastGroup := skills[0].Group
	groupStart := 0
	for i, sk := range skills {
		if sk.Group != lastGroup {
			result = append(result, skills[groupStart:i])
			lastGroup = sk.Group
			groupStart = i
		}
	}
	return append(result, skills[groupStart:len(skills)])
}

func characterInfo(cmd *cobra.Command, args []string) {
	xmlKey := getAPIKey()
	sde := getSDE()
	xmlapi := getXMLAPI(sde)
	charID := viper.GetInt("charid")
	charSheet, err := xmlapi.CharacterSheet(xmlKey, charID)
	if err != nil {
		log.Fatalf("Unable to call API: %v", err)
	}
	// Set up templating.
	funcMap := template.FuncMap{
		"roman": romanNumerals,
	}
	tmpl, err := template.New("charsheet").Funcs(funcMap).Parse(charsheetTmpl)
	if err != nil {
		log.Fatalf("Unable to parse template: %v", err)
	}
	evego.SortSkills(charSheet.Skills)
	sheet := charInfoDisplay{
		CharacterSheet: charSheet,
		Skills:         partitionSkills(charSheet.Skills),
	}
	err = tmpl.Execute(os.Stdout, sheet)
	if err != nil {
		log.Fatalf("Unable to execute template: %v", err)
	}
}

func printAssetsInt(sde evego.Database, assets []evego.InventoryItem, indentLevel int) {
	for _, asset := range assets {
		thisItem, err := sde.ItemForID(asset.TypeID)
		if err != nil {
			log.Fatalf("Bad item ID: %v", asset.TypeID)
		}
		for i := 0; i < indentLevel; i++ {
			fmt.Print(" ")
		}
		var packaged string
		if asset.Unpackaged {
			packaged = "(unpackaged)"
		} else {
			packaged = ""
		}
		fmt.Printf("%v x %v in %v, %v %v (%v)\n", asset.Quantity, thisItem.Name,
			asset.LocationID, packaged, asset.BlueprintType, asset.Flag)
		if len(asset.Contents) > 0 {
			printAssetsInt(sde, asset.Contents, indentLevel+2)
		}

	}
}

func printAssets(sde evego.Database, assets []evego.InventoryItem) {
	fmt.Printf("Assets:\n")
	printAssetsInt(sde, assets, 2)
}

func characterAssets(cmd *cobra.Command, args []string) {
	xmlKey := getAPIKey()
	sde := getSDE()
	xmlapi := getXMLAPI(sde)
	charID := viper.GetInt("charid")
	assets, err := xmlapi.Assets(xmlKey, charID)
	if err != nil {
		log.Fatalf("Unable to call API: %v", err)
	}
	printAssets(sde, assets)
}

func main() {
	rootCmd.PersistentFlags().Int("keyid", 0, "The key ID to use for accessing the account.")
	rootCmd.PersistentFlags().String("vcode", "", "The API key's verification code.")
	rootCmd.PersistentFlags().String("xmlapi", "https://api.eveonline.com", "The XML API server endpoint.")
	rootCmd.PersistentFlags().String("sdepath", "", "The current SDE dump in SQLite format.")
	rootCmd.PersistentFlags().String("localxml", "", "A local XML file to read instead of making the actual API call. (Optional)")
	rootCmd.AddCommand(accountCmd)
	flagNames := []string{"keyid", "vcode", "xmlapi", "sdepath", "localxml"}
	for _, fname := range flagNames {
		viper.BindPFlag(fname, rootCmd.PersistentFlags().Lookup(fname))
	}
	accountCmd.AddCommand(accountListCharsCmd)
	rootCmd.AddCommand(charCmd)
	charCmd.AddCommand(charSheetCmd)
	charCmd.AddCommand(charAssetsCmd)
	charCmd.PersistentFlags().Int("charid", 0, "The character ID of the toon to get information on.")
	flagNames = []string{"charid"}
	for _, fname := range flagNames {
		viper.BindPFlag(fname, charCmd.PersistentFlags().Lookup(fname))
	}

	viper.SetEnvPrefix("EVE")
	viper.AutomaticEnv()

	rootCmd.Execute()
}
