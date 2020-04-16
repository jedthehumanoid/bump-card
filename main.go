package main

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Card is a text with properties, meant for displaying on a board.
type Card struct {
	Title       string                 `json:"title"`
	Contents    string                 `json:"contents"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
	Frontmatter string                 `json:"frontmatter,omitempty"`
}

func getArguments(args []string) ([]string, []string) {
	arguments := []string{}
	flags := []string{}
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
		} else {
			arguments = append(arguments, arg)
		}
	}
	return arguments, flags
}

func main() {
	arguments, flags := getArguments(os.Args[1:])

	if containsString(flags, "-h") || containsString(flags, "--help") {
		fmt.Printf("Usage: %s [options] [files...]\n", os.Args[0])
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("   -h --help       Print help")
		fmt.Println("   -b --bump       Bump files")
		fmt.Println("   -r --recursive  List subdirectories recursively")
		fmt.Println()
		fmt.Println("If no files given, list all files")
		fmt.Println()
		os.Exit(0)
	}

	files := []string{}
	if len(arguments) == 0 {
		if containsString(flags, "-r") || containsString(flags, "--recursive") {
			files = findFiles(".")
		} else {
			files = readDir(".")
		}
	} else {
		files = arguments
	}
	cards := readCards(files)

	if containsString(flags, "-b") || containsString(flags, "--bump") {
		if len(arguments) == 0 {
			fmt.Println("Will not bump all files implicitly")
			os.Exit(0)
		}
		for _, card := range cards {
			fmt.Printf("Bumping %s\n", card.Title)

			time := time.Now().UTC().Format(time.RFC3339)

			if card.Properties == nil {
				card.Properties = map[string]interface{}{}
			}
			if card.Frontmatter == "" {
				card.Frontmatter = "toml"
			}

			card.Properties["bumped"] = time
			contents := card.Contents
			contents = MarshalFrontmatter(card, true) + "\n" + contents
			err := ioutil.WriteFile(card.Title, []byte(contents), 0644)
			if err != nil {
				panic(err)
			}
		}

	} else {
		output := []string{}
		for _, card := range cards {
			bumped := card.Properties["bumped"]
			if bumped == nil {
				bumped = "        ----        "
			}
			output = append(output, fmt.Sprintf("%s %s", bumped, card.Title))
		}
		sort.Strings(output)
		sort.Sort(sort.Reverse(sort.StringSlice(output)))
		fmt.Println()
		for _, line := range output {
			fmt.Println(line)
		}
		fmt.Println()
	}
}

// ReadCardFile takes a file path, reading file in to a card.
func readCard(path string) (Card, error) {
	var card Card

	card.Title = path

	contents, err := ioutil.ReadFile(filepath.FromSlash(path))
	if err != nil {
		return card, err
	}
	frontmatter, raw, b := getFrontmatter(contents)
	card.Frontmatter = frontmatter

	card.Contents = strings.TrimPrefix(string(contents), string(raw))
	card.Contents = strings.TrimSpace(card.Contents)

	switch frontmatter {
	case "yaml":
		err = yaml.Unmarshal(b, &card.Properties)
		return card, err
	case "toml":
		_, err = toml.Decode(string(b), &card.Properties)
		return card, err
	}
	return card, nil
}

func isCard(file string) bool {
	return strings.HasSuffix(file, ".md")
}

func readCards(files []string) []Card {
	cards := []Card{}

	for _, file := range files {
		if !isCard(file) {
			continue
		}
		card, err := readCard(file)
		if err != nil {
			panic(err)
		}
		cards = append(cards, card)
	}
	return cards
}

func MarshalFrontmatter(card Card, fences bool) string {
	ret := ""

	switch card.Frontmatter {
	case "yaml":
		b, _ := yaml.Marshal(card.Properties)
		frontmatter := strings.TrimSpace(string(b))
		if frontmatter != "{}" {
			ret = frontmatter
			if fences {
				ret = "---\n" + ret + "\n---"
			}
		}
	case "toml":
		buf := new(bytes.Buffer)
		toml.NewEncoder(buf).Encode(card.Properties)
		frontmatter := strings.TrimSpace(buf.String())
		if frontmatter != "" {
			ret = frontmatter
			if fences {
				ret = "+++\n" + ret + "\n+++"
			}
		}
	}

	return ret
}
