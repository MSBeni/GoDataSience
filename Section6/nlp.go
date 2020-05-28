package main

import(
	"fmt"
	"gopkg.in/jdkato/prose.v2"
	"regexp"
	"strings"
)

func preProcessingNLP(text string) string{
	// Find all characters that are not alphabets
	reg := regexp.MustCompile("[^a-zA-Z]+")

	// Replace those chars with spaces
	text = reg.ReplaceAllString(text, " ")

	// Lower case
	text = strings.ToLower(text)
	return text
}

func main(){
	str := "On Wednesday, May 27, @AstroBehnken and @Astro_Doug will lift off aboard @SpaceX's Crew Dragon on the " +
		"company’s Falcon 9 rocket and fly to the @Space_Station. Weather remains 40% favorable: https://t.co/m2wtN8"
	str = preProcessingNLP(str)

	println("preProcessingNLP result: ", str)

	// Tokenization 1
	fmt.Println("#################################")
	fmt.Println("Tokenization 1: ")
	tokens := strings.Fields(str)
	for idx, token := range tokens{
		fmt.Println(idx, token, len(tokens))
	}

	// Create a new document with the default configuration:
	text :=  "This spacecraft will carry @AstroBehnken and @Astro_Doug to the @Space_Station when it launches atop a" +
		" @SpaceX Falcon 9 rocket on May 27, at 4:33 p.m. ET: https://t.co/yvfOCG4"
	doc, err := prose.NewDocument(text)
	if err != nil{
		panic(err)
	}
	// Tokenization 2
	fmt.Println("#################################")
	fmt.Println("Tokenization 2: ")
	for _, tok := range doc.Tokens(){
		fmt.Println(tok.Text, tok.Tag, tok.Label)
	}

	// Get named entities from document
	fmt.Println("#################################")
	fmt.Println("Get named entities from document: ")
	for _, ent := range doc.Entities(){
		fmt.Println(ent.Text, ent.Label)
	}

	// Get document sentences
	fmt.Println("#################################")
	fmt.Println("Get document sentences: ")
	for idx, sent := range doc.Sentences(){
		fmt.Println(idx, sent.Text)
	}
}