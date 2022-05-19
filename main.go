package main

import (
    "fmt"
    // "html"
		"strings"
		"bytes"
		"io"
    "log"
		"io/ioutil"
    "net/http"
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

func main() {
		var url = "https://sentencedict.com/Fed%20up.html"
		data , err := getHtmlPage(url)
		// res , err := http.Get(url)
		// fmt.Println(data)
		if err != nil {
			log.Fatal(err)
		}
		doc, err := html.Parse(strings.NewReader(data))

    if err != nil {
        log.Fatal(err)
    }

    tag := getElementById(doc, "all")
    output := renderNode(tag)

    fmt.Println(output)

}