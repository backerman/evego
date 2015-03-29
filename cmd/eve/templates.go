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

import "fmt"

var numMap = [...]string{"0", "I", "II", "III", "IV", "V"}

func romanNumerals(n int) string {
	if n < 0 || n >= len(numMap) {
		return fmt.Sprintf("(out of range: %d)", n)
	}
	return numMap[n]
}

var (
	charsheetTmpl = `
{{.Name}} ({{.ID}})
Corporation: {{.Corporation}} ({{.CorporationID}})
Alliance:    {{.Alliance}} ({{.AllianceID}})
Skills:{{with .Skills}}{{range $i, $skillGroup := .}}{{range $j, $sk := .}}
{{if eq $j 0 }}  {{$sk.Group}}: ({{len $skillGroup}} skills)
{{end}}    {{$sk.Name}} {{roman $sk.Level}} ({{$sk.NumSkillpoints}} pts.){{end}}{{end}}{{end}}
`
)
