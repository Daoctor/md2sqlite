package main

import (
	db "./savedb"
	"bufio"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"regexp"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func pop(element []string) []string {
	return element[:len(element)-1]
}

func isEmpty(element []string) bool {
	if len(element) > 0 {
		return false
	}
	return true
}

func getResult(d string, r *regexp.Regexp) (string, error) {
	result := r.FindStringSubmatch(d)
	if len(result) > 1 {
		return result[1], nil
	}
	return "not found", errors.New("result not found")
}

func main() {
	// args 1 db_path 2 file_path
	args := os.Args
	args = args[1:]

	if len(args) != 2 {
		fmt.Print("Args not right!\n")
		fmt.Print("Usage:\n")
		fmt.Print("    ./savemd  db_file_path  md_file_path\n")
		return
	}
	dbPath := args[0]
	fileName := args[1]
	dataGet := false
	dataR, err := regexp.Compile(`date\:\s*([\d-]+)`)
	linkR, err := regexp.Compile(`https://github.com/(.*?)\"`)
	tagR, err := regexp.Compile(`\*\*(.*?)\*\*`)
	fh, err := os.Open(fileName)
	f := bufio.NewReader(fh)
	check(err)
	defer fh.Close()
	buf := make([]byte, 1024)
	stack := []string{}
	resultList := []db.TendingData{}
	var currentLanguage string
	var curData string
	for {
		buf, _, err = f.ReadLine()
		if err != nil {
			break
		}
		sData := string(buf)
		if !dataGet && strings.HasPrefix(sData, "date") {
			r, err := getResult(sData, dataR)
			if err != nil {
			} else {
				curData = r
			}
		}
		if strings.HasPrefix(sData, "###") {
			r, err := getResult(sData, linkR)
			if err != nil {
			} else {
				stack = append(stack, r)
			}
		}
		if strings.HasPrefix(sData, `# Language`) {
			r, err := getResult(sData, tagR)
			if err != nil {
			} else {
				if !isEmpty(stack) && r != currentLanguage {
					t := stack[:]
					s := db.TendingData{Language: currentLanguage, Date: curData, Links: t}
					resultList = append(resultList, s)
					stack = stack[:0]
				}
				currentLanguage = r
			}
		}
	}
	s := db.TendingData{Language: currentLanguage, Date: curData, Links: stack}
	resultList = append(resultList, s)
	fmt.Print("done\n")
	for _, d := range resultList {
		db.Save(dbPath, d)
	}
}
