package main

import (
	"fmt"
	"github.com/fumiama/go-docx"
	"os"
)

type Entry struct {
	question string
	answers  []Answer
}
type Answer struct {
	text      string
	isCorrect bool
}

func CheckColor(p *docx.Paragraph) bool {
	for _, v := range p.Children {
		switch v.(type) {
		case *docx.Run:
			run := v.(*docx.Run)
			runProp := run.RunProperties
			if runProp != nil {
				if runProp.Color != nil {
					return true
				}
			}
		}
	}
	pP := p.Properties
	if pP == nil {
		return false
	}
	pRP := pP.RunProperties
	if pRP == nil {
		return false
	}
	color := pRP.Color
	if color != nil {
		if color.Val == "FF0000" {
			return true
		}

	}
	return false
}

func main() {
	readFile, err := os.Open("test/baza_AI.docx")
	if err != nil {
		panic(err)
	}
	fileInfo, err := readFile.Stat()
	if err != nil {
		panic(err)
	}
	size := fileInfo.Size()
	doc, err := docx.Parse(readFile, size)
	if err != nil {
		panic(err)
	}
	fmt.Println("Plain text:")
	entries := make([]Entry, 0, 5)
	for _, it := range doc.Document.Body.Items {
		switch it.(type) {
		case *docx.Table:
			table := it.(*docx.Table)
			var rows = table.TableRows
			for _, r := range rows {
				if len(r.TableCells) != 3 {
					continue
				}
				var entry Entry
				for i, c := range r.TableCells {
					if i == 0 {
						continue
					}
					if i == 1 {
						entry.question = c.Paragraphs[0].String()
						continue
					}

					//fmt.Println(i)
					for _, p := range c.Paragraphs {
						if len(p.String()) == 0 {
							continue
						}
						fmt.Print(p.String())
						var answer Answer
						answer.text = p.String()
						answer.isCorrect = CheckColor(p)
						entry.answers = append(entry.answers, answer)
						fmt.Println()
					}
				}
				entries = append(entries, entry)
				//}

			}
		}
	}
	for i, v := range entries {
		fmt.Println(i, v.question)
		for j, q := range v.answers {
			fmt.Println(j, q.text, q.isCorrect)
		}
	}
}
