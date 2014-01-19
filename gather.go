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
	Types         []string `json:"type"`
	Subtypes      []string `json:"type"`
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
	card.ConvertedCost = extractInt(doc, prefix+"cmcRow .value")
	card.Number = extractString(doc, prefix+"numberRow .value")
	card.RulesText = extractString(doc, prefix+"textRow .value")
	card.Artist = extractString(doc, prefix+"artistRow .value")
	card.Set = extractString(doc, prefix+"setRow .value")
	card.FlavorText = extractString(doc, prefix+"FlavorText")

	return card, nil
}

func main() {
	fmt.Println("Foo")
}
