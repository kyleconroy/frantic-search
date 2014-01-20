package main

import (
	"code.google.com/p/go.net/html"
	"crypto/md5"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	prefix = "#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_"
)

type Card struct {
	Name          string    `json:"name"`
	Id            string    `json:"id"`
	Types         []string  `json:"types"`
	Subtypes      []string  `json:"subtypes"`
	ConvertedCost int       `json:"converted_cost"`
	ManaCost      string    `json:"mana_cost"`
	Special       string    `json:"special"` //'flip', 'double-faced', 'split'
	PartnerCard   string    `json:"partner_card"`
	RulesText     []string  `json:"rules_text"`
	Power         string    `json:"power"`
	Toughness     string    `json:"toughness"`
	Editions      []Edition `json:"editions"`
}

type Edition struct {
	Set          string `json:"set"`
	Rarity       string `json:"rarity"`
	Artist       string `json:"artist"`
	MultiverseId int    `json:"multiverse"`
	Mark         string `json:"mark"`
	FlavorText   string `json:"flavor_text"`
	Number       string `json:"number"`
}

func (c Card) ImageURl() string {
	return "http://gatherer.wizards.com/Handlers/Image.ashx?multiverseid="
}

func cardName(n *html.Node) string {
	return ""
}

func extractString(n *html.Node, pattern string) string {
	if div, found := Find(n, pattern); found {
		return strings.TrimSpace(Flatten(div))
	} else {
		return ""
	}
}

func manaSymbol(alt string) string {
	switch alt {
	case "0", "1", "2", "3", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15":
		return alt
	case "Green":
		return "G"
	case "Red":
		return "R"
	case "Blue":
		return "U"
	case "Black":
		return "B"
	case "White":
		return "W"
	}
	return ""
}

func FlattenWithSymbols(n *html.Node) string {
	text := ""
	if n.Type == html.TextNode {
		text += n.Data
	}

	if n.Type == html.ElementNode && n.Data == "img" {
		text += "{" + manaSymbol(Attr(n, "alt")) + "}"
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += FlattenWithSymbols(c)
	}

	return text
}

func extractManaCost(n *html.Node) string {
	cost := ""
	for _, a := range FindAll(n, prefix+"manaRow .value img") {
		cost += manaSymbol(Attr(a, "alt"))
	}
	return cost
}

func extractPT(n *html.Node) (string, string) {
	div, found := Find(n, prefix+"ptRow .value")

	if !found {
		return "", ""
	}

	values := strings.Split(strings.TrimSpace(Flatten(div)), "/")

	if len(values) != 2 {
		return "", ""
	}

	return strings.TrimSpace(values[0]), strings.TrimSpace(values[1])
}

func SplitTrimSpace(source, pattern string) []string {
	result := []string{}

	for _, val := range strings.Split(strings.TrimSpace(source), pattern) {
		result = append(result, strings.TrimSpace(val))
	}

	return result
}

func extractTypes(n *html.Node) ([]string, []string) {
	div, found := Find(n, prefix+"typeRow .value")

	if !found {
		return []string{}, []string{}
	}

	ts := strings.Split(strings.ToLower(strings.TrimSpace(Flatten(div))), "—")

	var types []string
	var subtypes []string

	if len(ts) == 2 {
		types = SplitTrimSpace(ts[0], " ")
		subtypes = SplitTrimSpace(ts[1], " ")
	} else {
		types = SplitTrimSpace(ts[0], " ")
		subtypes = []string{}
	}

	return types, subtypes
}

func extractRarity(n *html.Node) string {
	if span, found := Find(n, prefix+"rarityRow .value span"); found {
		return Attr(span, "class")
	} else {
		return ""
	}
}

func extractInt(n *html.Node, pattern string) int {
	div, found := Find(n, pattern)

	if !found {
		return 0
	}

	number, err := strconv.Atoi(strings.TrimSpace(Flatten(div)))

	if err != nil {
		return 0
	}

	return number
}

func extractRulesText(n *html.Node) []string {
	rules := []string{}
	for _, node := range FindAll(n, prefix+"textRow .value .cardtextbox") {
		rules = append(rules, strings.TrimSpace(FlattenWithSymbols(node)))
	}
	return rules
}

func hash(in string) string {
	h := md5.New()
	io.WriteString(h, in)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func ParseCard(page io.Reader) (Card, error) {
	doc, err := html.Parse(page)

	card := Card{}
	edition := Edition{}

	if err != nil {
		return card, err
	}

	card.Name = extractString(doc, prefix+"nameRow .value")
	card.Id = hash(card.Name)
	card.ManaCost = extractManaCost(doc)
	card.ConvertedCost = extractInt(doc, prefix+"cmcRow .value")
	card.RulesText = extractRulesText(doc)
	card.Types, card.Subtypes = extractTypes(doc)
	card.Power, card.Toughness = extractPT(doc)

	edition.Number = extractString(doc, prefix+"numberRow .value")
	edition.Artist = extractString(doc, prefix+"artistRow .value")
	edition.Set = extractString(doc, prefix+"setRow .value")
	edition.FlavorText = extractString(doc, prefix+"FlavorText")
	edition.Rarity = extractRarity(doc)

	card.Editions = append(card.Editions, edition)

	return card, nil
}

func main() {
	fmt.Println("Foo")
}
