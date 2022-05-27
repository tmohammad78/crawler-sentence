package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
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

func generateURL(url string,word string) string{
	newUrl := url
	if(strings.Contains(word," ")){
		words := strings.Split(word," ")
		newUrl += "/" + words[0] + "%20" + words[1]+ ".html"
	}
	return newUrl
}

func generateSentense(tokenizer *html.Tokenizer) string {
	ttt := ""
	digitDot := regexp.MustCompile(`\d*\.`)
	previousStartTokenTest := tokenizer.Token()
	loopDomTest :
	for {
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
				// fmt.Printf("%s\n", TxtContent)
				// fmt.Printf(previousStartTokenTest.Data)
				if len(mainWord) > 0 {
					TxtContent += " " + mainWord + " "
				}
				ttt += TxtContent
				// if len(TxtContent) > 0 {
				// 		// fmt.Printf("%s\n", TxtContent)
				// }			
			}
		}
		return ttt
}

func generatedObject( wantedWord string,ttt string ) sentencesList {
	sentencesSlice := sentencesList{}
	for index,item:= range strings.Split(ttt,"."){
		sentencesSlice[wantedWord] = append(sentencesSlice[wantedWord], sentence{id:index,shown:false,text:item})
	}
	return sentencesSlice
}

func generateHTMLDOM(data string) *html.Node {
	doc, err := html.Parse(strings.NewReader(data))

	if err != nil {
			log.Fatal(err)
	}
	return doc
}

var (
	seq   = 1
)

type sentence  struct { 
	id int
	shown bool
	text string
}

type word struct{
	ID int `json:"id"`
	Word string `json:"word"`
}

type sentencesList map[string][]sentence

func main() {
		var url = "https://sentencedict.com"
		e := echo.New()
		result := sentencesList{}
		wantedWord :=""
		e.GET("/",func (c echo.Context) error {
			return c.String(http.StatusOK,"Hello")
		})
		e.POST("/word", func (c echo.Context) error{
			w := &word{
				ID: seq,
			}
			if err := c.Bind(w);err != nil {
				return err
			}
			wantedWord = c.FormValue("word")
			
			return c.JSON(http.StatusOK, "Done")
		})

		e.GET("/word",func (c echo.Context) error  {
			return c.JSON(http.StatusOK, result)
		})
		e.Logger.Fatal(e.Start(":1323"))
		data , err := getHtmlPage(generateURL(url,wantedWord))

		if err != nil {
			log.Fatal(err)
		}

		output := renderNode(getElementById(generateHTMLDOM(data), "all"))
		tokenizer := html.NewTokenizer(strings.NewReader(output))
		ttt:= generateSentense(tokenizer)
		result = generatedObject(wantedWord,ttt)
		file , err := os.Create("wtf"); 
		if err != nil {
			panic(err)
		}
		defer file.Close()
		file1, _ := json.MarshalIndent(result, "", " ")
		_ = ioutil.WriteFile(file.Name(), file1, 0644)

}