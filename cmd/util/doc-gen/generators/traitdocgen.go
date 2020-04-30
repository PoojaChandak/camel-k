/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package generators

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"regexp"
	"sort"
	"strings"

	v1 "github.com/apache/camel-k/pkg/apis/camel/v1"
	"github.com/apache/camel-k/pkg/trait"
	"k8s.io/gengo/args"
	"k8s.io/gengo/generator"
	"k8s.io/gengo/types"
)

const (
	tagTrait = "+camel-k:trait"

	adocCommonMarkerStart = "// Start of autogenerated code - DO NOT EDIT!"
	adocCommonMarkerEnd   = "// End of autogenerated code - DO NOT EDIT!"

	adocDescriptionMarkerStart = adocCommonMarkerStart + " (description)"
	adocDescriptionMarkerEnd   = adocCommonMarkerEnd + " (description)"

	adocConfigurationMarkerStart = adocCommonMarkerStart + " (configuration)"
	adocConfigurationMarkerEnd   = adocCommonMarkerEnd + " (configuration)"

	adocNavMarkerStart = adocCommonMarkerStart + " (trait-nav)"
	adocNavMarkerEnd   = adocCommonMarkerEnd + " (trait-nav)"

	adocListMarkerStart = adocCommonMarkerStart + " (trait-list)"
	adocListMarkerEnd   = adocCommonMarkerEnd + " (trait-list)"
)

var (
	tagTraitID = regexp.MustCompile(fmt.Sprintf("%s=([a-z0-9-]+)", regexp.QuoteMeta(tagTrait)))
)

// traitDocGen produces documentation about traits
type traitDocGen struct {
	generator.DefaultGen
	arguments           *args.GeneratorArgs
	generatedTraitFiles []string
}

// NewTraitDocGen --
func NewTraitDocGen(arguments *args.GeneratorArgs) generator.Generator {
	return &traitDocGen{
		DefaultGen: generator.DefaultGen{},
		arguments:  arguments,
	}
}

func (g *traitDocGen) Filename() string {
	return "zz_generated_doc.go"
}

func (g *traitDocGen) Filter(context *generator.Context, t *types.Type) bool {
	for _, c := range t.CommentLines {
		if strings.Contains(c, tagTrait) {
			return true
		}
	}
	return false
}

func (g *traitDocGen) GenerateType(context *generator.Context, t *types.Type, out io.Writer) error {
	docDir := g.arguments.CustomArgs.(*CustomArgs).DocDir
	traitPath := g.arguments.CustomArgs.(*CustomArgs).TraitPath
	traitID := getTraitID(t)
	traitFile := traitID + ".adoc"
	filename := path.Join(docDir, traitPath, traitFile)

	g.generatedTraitFiles = append(g.generatedTraitFiles, traitFile)

	file, content, err := readFile(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writeTitle(traitID, &content)
	writeDescription(t, traitID, &content)
	writeFields(t, traitID, &content)

	return writeFile(file, content)
}

func (g *traitDocGen) Finalize(c *generator.Context, w io.Writer) error {
	return g.FinalizeNav(c)
}

func (g *traitDocGen) FinalizeNav(*generator.Context) error {
	docDir := g.arguments.CustomArgs.(*CustomArgs).DocDir
	navPath := g.arguments.CustomArgs.(*CustomArgs).NavPath
	filename := path.Join(docDir, navPath)

	file, content, err := readFile(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	pre, post := split(content, adocNavMarkerStart, adocNavMarkerEnd)

	content = append([]string(nil), pre...)
	content = append(content, adocNavMarkerStart)
	sort.Strings(g.generatedTraitFiles)
	for _, t := range g.generatedTraitFiles {
		name := traitNameFromFile(t)
		content = append(content, "** xref:traits:"+t+"["+name+"]")
	}
	content = append(content, adocNavMarkerEnd)
	content = append(content, post...)

	return writeFile(file, content)
}

func traitNameFromFile(file string) string {
	name := strings.TrimSuffix(file, ".adoc")
	name = strings.ReplaceAll(name, "trait", "")
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.Trim(name, " ")
	name = strings.Title(name)
	return name
}

func writeTitle(traitID string, content *[]string) {
	res := append([]string(nil), *content...)
	for _, s := range res {
		if strings.HasPrefix(s, "= ") {
			// Already has a title
			return
		}
	}
	res = append([]string{"= " + strings.Title(strings.ReplaceAll(traitID, "-", " ")) + " Trait"}, res...)
	*content = res
}

func writeDescription(t *types.Type, traitID string, content *[]string) {
	pre, post := split(*content, adocDescriptionMarkerStart, adocDescriptionMarkerEnd)
	res := append([]string(nil), pre...)
	res = append(res, adocDescriptionMarkerStart)
	res = append(res, filterOutTagsAndComments(t.CommentLines)...)
	profiles := strings.Join(determineProfiles(traitID), ", ")
	res = append(res, "", fmt.Sprintf("This trait is available in the following profiles: **%s**.", profiles))
	if isPlatformTrait(traitID) {
		res = append(res, "", fmt.Sprintf("WARNING: The %s trait is a *platform trait*: disabling it may compromise the platform functionality.", traitID))
	}
	res = append(res, "", adocDescriptionMarkerEnd)
	res = append(res, post...)
	*content = res
}

func writeFields(t *types.Type, traitID string, content *[]string) {
	pre, post := split(*content, adocConfigurationMarkerStart, adocConfigurationMarkerEnd)
	res := append([]string(nil), pre...)
	res = append(res, adocConfigurationMarkerStart, "== Configuration", "")
	res = append(res, "Trait properties can be specified when running any integration with the CLI:")
	res = append(res, "```")
	if len(t.Members) > 1 {
		res = append(res, fmt.Sprintf("kamel run --trait %s.[key]=[value] --trait %s.[key2]=[value2] integration.groovy", traitID, traitID))
	} else {
		res = append(res, fmt.Sprintf("kamel run --trait %s.[key]=[value] integration.groovy", traitID))
	}
	res = append(res, "```")
	res = append(res, "The following configuration options are available:", "")
	res = append(res, "[cols=\"2,1,5a\"]", "|===")
	res = append(res, "|Property | Type | Description", "")
	writeMembers(t, traitID, &res)
	res = append(res, "|===", "", adocConfigurationMarkerEnd)
	res = append(res, post...)
	*content = res
}

func writeMembers(t *types.Type, traitID string, content *[]string) {
	res := append([]string(nil), *content...)
	for _, m := range t.Members {
		prop := reflect.StructTag(m.Tags).Get("property")
		if prop != "" {
			if strings.Contains(prop, "squash") {
				writeMembers(m.Type, traitID, &res)
			} else {
				res = append(res, "| "+traitID+"."+prop)
				res = append(res, "| "+strings.TrimPrefix(m.Type.Name.Name, "*"))
				first := true
				for _, l := range filterOutTagsAndComments(m.CommentLines) {
					if first {
						res = append(res, "| "+l)
						first = false
					} else {
						res = append(res, l)
					}
				}
				res = append(res, "")
			}
		}
	}
	*content = res
}

func getTraitID(t *types.Type) string {
	for _, s := range t.CommentLines {
		if strings.Contains(s, tagTrait) {
			matches := tagTraitID.FindStringSubmatch(s)
			if len(matches) < 2 {
				panic(fmt.Sprintf("unable to extract trait ID from tag line `%s`", s))
			}
			return matches[1]
		}
	}
	panic(fmt.Sprintf("trait ID not found in type %s", t.Name.Name))
}

func filterOutTagsAndComments(comments []string) []string {
	res := make([]string, 0, len(comments))
	for _, l := range comments {
		if !strings.HasPrefix(strings.TrimLeft(l, " \t"), "+") &&
			!strings.HasPrefix(strings.TrimLeft(l, " \t"), "TODO:") {
			res = append(res, l)
		}
	}
	return res
}

func split(doc []string, startMarker, endMarker string) (pre []string, post []string) {
	if len(doc) == 0 {
		return nil, nil
	}
	idx := len(doc)
	for i, s := range doc {
		if s == startMarker {
			idx = i
			break
		}
	}
	idy := len(doc)
	for j, s := range doc {
		if j > idx && s == endMarker {
			idy = j
			break
		}
	}
	pre = doc[0:idx]
	if idy < len(doc) {
		post = doc[idy+1:]
	}
	return pre, post
}

func readFile(filename string) (file *os.File, content []string, err error) {
	if file, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0777); err != nil {
		return file, content, err
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return file, content, err
	}
	content = strings.Split(string(bytes), "\n")
	return file, content, nil
}

func writeFile(file *os.File, content []string) error {
	if err := file.Truncate(0); err != nil {
		return err
	}
	max := 0
	for i, line := range content {
		if line != "" {
			max = i
		}
	}
	for i, line := range content {
		if i <= max {
			if _, err := file.WriteString(line + "\n"); err != nil {
				return err
			}
		}
	}
	return nil
}

func isPlatformTrait(traitID string) bool {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	catalog := trait.NewCatalog(ctx, nil)
	t := catalog.GetTrait(traitID)
	return t.IsPlatformTrait()
}

func determineProfiles(traitID string) (profiles []string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	catalog := trait.NewCatalog(ctx, nil)
	for _, p := range v1.AllTraitProfiles {
		traits := catalog.TraitsForProfile(p)
		for _, t := range traits {
			if string(t.ID()) == traitID {
				profiles = append(profiles, string(p))
			}
		}
	}
	return profiles
}
