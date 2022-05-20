package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func getAttribute(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
			if attr.Key == key {
					return attr.Val, true
			}
	}

	return "", false
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)

	err := html.Render(w, n)

	if err != nil {
			return ""
	}

	return buf.String()
}

func checkId(n *html.Node, id string) bool {

	if n.Type == html.ElementNode {

	s, ok := getAttribute(n, "id")

			if ok && s == id {
					return true
			}
	}

	return false
}

func traverse(n *html.Node, id string) *html.Node {
	if checkId(n, id) {
			return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
			res := traverse(c, id)
			if res != nil {
					return res
			}
	}

	return nil
}

func getElementById(n *html.Node, id string) *html.Node {
	return traverse(n, id)
}

func getHtmlPage(webPage string) (string, error) {
	resp, err := http.Get(webPage)

	if err != nil {
			return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {

			return "", err
	}

	return string(body), nil
}

type sentence  struct { 
	id uint8
	shown bool
	text string
}
type sentencesList map[string][]sentence

func main() {
		var url = "https://sentencedict.com/Fed%20up.html"
		data , err := getHtmlPage(url)
		sentencesSlice := sentencesList{}
		sentencesSlice["fed up"] = append(sentencesSlice["fed up"], sentence{id:12,shown:false,text: "test"})

		if err != nil {
			log.Fatal(err)
		}
		doc, err := html.Parse(strings.NewReader(data))

    if err != nil {
        log.Fatal(err)
    }

    tag := getElementById(doc, "all")
    output := renderNode(tag)
		digitDot := regexp.MustCompile(`\d*\.`)
		tokenizer := html.NewTokenizer(strings.NewReader(output))
		previousStartTokenTest := tokenizer.Token()
		loopDomTest :
		for {
			ttt := ""
			tokenType := tokenizer.Next()
			switch {
				case tokenType == html.ErrorToken:
					break loopDomTest
				case tokenType == html.StartTagToken:
					previousStartTokenTest = tokenizer.Token()
				case tokenType == html.TextToken :
					mainWord := ""
					if previousStartTokenTest.Data == "script" {
						continue
					}
					if previousStartTokenTest.Data == "em" {
						mainWord = strings.TrimSpace(html.UnescapeString(string(tokenizer.Text())))
					}
					TxtContent := strings.TrimSpace(html.UnescapeString(string(tokenizer.Text())))
					TxtContent = digitDot.ReplaceAllString(TxtContent,"${1}")
					// fmt.Printf(previousStartTokenTest.Data)
					if len(mainWord) > 0 {
						TxtContent += mainWord
					}
					ttt += TxtContent
					if len(TxtContent) > 0 {
							// fmt.Printf("%s\n", TxtContent)
					}
					fmt.Printf("%s\n", ttt)
			}
		}

}