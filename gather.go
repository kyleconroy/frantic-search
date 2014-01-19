package main

import (
	"code.google.com/p/go.net/html"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	prefix = "#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_"
)

type Card struct {
	Artist        string   `json:"artist"`
	Name          string   `json:"name"`
	Number        string   `json:"number"`
	Rarity        string   `json:"rarity"`
	Types         []string `json:"types"`
	Subtypes      []string `json:"subtypes"`
	Set           string   `json:"set"`
	MultiverseId  int      `json:"multiverse"`
	ConvertedCost int      `json:"converted_cost"`
	ManaCost      string   `json:"mana_cost"`
	Special       string   `json:"special"` //'flip', 'double-faced', 'split'
	PartnerCard   int      `json:"partner_card"`
	Mark          string   `json:"mark"`
	RulesText     string   `json:"rules_text"`
	FlavorText    string   `json:"flavor_text"`
	Power         string   `json:"power"`
	Toughness     string   `json:"toughness"`
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

func extractManaCost(n *html.Node) string {
	cost := ""
	for _, a := range FindAll(n, prefix+"manaRow .value img") {
		cost += manaSymbol(Attr(a, "alt"))
	}
	return cost
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
		types = strings.Split(ts[0], " ")
		subtypes = strings.Split(ts[1], " ")
	} else {
		types = strings.Split(ts[0], " ")
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

func ParseCard(page io.Reader) (Card, error) {
	doc, err := html.Parse(page)

	card := Card{}

	if err != nil {
		return card, err
	}

	card.Rarity = extractRarity(doc)
	card.Name = extractString(doc, prefix+"nameRow .value")
	card.ManaCost = extractManaCost(doc)
	card.ConvertedCost = extractInt(doc, prefix+"cmcRow .value")
	card.Number = extractString(doc, prefix+"numberRow .value")
	card.RulesText = extractString(doc, prefix+"textRow .value")
	card.Artist = extractString(doc, prefix+"artistRow .value")
	card.Set = extractString(doc, prefix+"setRow .value")
	card.FlavorText = extractString(doc, prefix+"FlavorText")
	card.Types, card.Subtypes = extractTypes(doc)

	return card, nil
}

func main() {
	fmt.Println("Foo")
}
