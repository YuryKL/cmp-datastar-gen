package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

const (
	INPUT  = "input.json"
	OUTPUT = "data.lua"
)

type SnippetMap map[string]Snippet
type Snippet struct {
	Prefix      string      `json:"prefix"`
	Body        string      `json:"body"`
	Description string      `json:"description"`
	References  []Reference `json:"references"`
}
type Reference struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func main() {
	input, err := os.ReadFile(INPUT)
	if err != nil {
		panic(err)
	}

	var snippets SnippetMap
	if err := json.Unmarshal(input, &snippets); err != nil {
		panic(err)
	}

	keys := make([]string, 0, len(snippets))
	for k := range snippets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	out, err := os.Create(OUTPUT)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	if err := generateLua(out, snippets, keys); err != nil {
		panic(err)
	}
	fmt.Println("✅ data.lua generated successfully.")

}

func generateLua(out *os.File, snippets SnippetMap, keys []string) error {
	fmt.Fprintln(out, "local items = {")

	if len(keys) == 0 {
		fmt.Fprintln(out, "}")
		fmt.Fprintln(out, "return items")
		return nil
	}

	for _, key := range keys {
		s := snippets[key]
		url := ""
		if len(s.References) > 0 {
			url = s.References[0].URL
		}

		fmt.Fprintf(out, "  {\n")
		fmt.Fprintf(out, "    label = '%s',\n", s.Prefix)
		fmt.Fprintf(out, "    insertText = '%s',\n", s.Body)
		fmt.Fprintf(out, "    filterText = '%s',\n", s.Prefix)
		fmt.Fprintf(out, "    documentation = '%s',\n", s.Description)
		fmt.Fprintf(out, "    kind = 12,\n")
		fmt.Fprintf(out, "    detail = '%s',\n", url)
		fmt.Fprintf(out, "  },\n")
	}

	fmt.Fprintln(out, "}")
	fmt.Fprintln(out, "return items")
	return nil

}
