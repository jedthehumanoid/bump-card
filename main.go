package main

import (
	"fmt"
	"github.com/jedthehumanoid/card-cabinet"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"
)

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
		fmt.Println("   -f --force      Bump even if no previous frontmatter")
		fmt.Println()
		fmt.Println("If no files given, list all files")
		fmt.Println()
		os.Exit(0)
	}

	files := []string{}
	if len(arguments) == 0 {
		if containsString(flags, "-r") || containsString(flags, "--recursive") {
			files = cardcabinet.FindFiles(".")
		} else {
			files = readDir(".")
		}
	} else {
		files = arguments
	}
	cards := cardcabinet.ReadCards(files)

	if containsString(flags, "-b") || containsString(flags, "--bump") {
		if len(arguments) == 0 {
			fmt.Println("Will not bump all files implicitly")
			os.Exit(0)
		}
		for _, card := range cards {
			fmt.Printf("Bumping %s\n", card.Title)

			time := time.Now().UTC().Format(time.RFC3339)

			if card.Properties == nil {
				if containsString(flags, "-f") || containsString(flags, "--force") {
					card.Properties = map[string]interface{}{}
				} else {
					fmt.Println("No frontmatter in file, will not bump unless forced")
					os.Exit(0)
				}
			}
			if card.Frontmatter == "" {
				card.Frontmatter = "toml"
			}

			card.Properties["bumped"] = time
			contents := card.Contents
			contents = card.MarshalFrontmatter(true) + "\n" + contents
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
