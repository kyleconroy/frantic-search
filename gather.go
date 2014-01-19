package main

import (
	"code.google.com/p/go.net/html"
	"fmt"
	"io"
	"strings"
	"strconv"
)

const (
	prefix = "ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_"
)

type Card struct {
	Artist        string   `json:"artist"`
	Name          string   `json:"name"`
	Number        int      `json:"number"`
	Rarity        string   `json:"rarity"`
	Types         []string `json:"type"`
	Subtypes      []string `json:"type"`
	Expansion     string   `json:"set"`
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

func hasId(n *html.Node, id string) bool {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if a.Key == "id" {
				return a.Val == id
			}
		}
	}
	return false
}

func find(n *html.Node, class string) *html.Node {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			for _, a := range c.Attr {
				if a.Key == "class" && a.Val == class {
					return c
				}
			}
		}

	}
	return nil
}

func text(n *html.Node) string {
	if n.FirstChild == nil {
		return ""
	}

	if n.FirstChild.Type == html.TextNode {
		return n.FirstChild.Data
	}

	return ""
}

func extractName(n *html.Node) string {
	div := find(n, "value")

	if div == nil {
		return ""
	}

	return strings.TrimSpace(text(div))
}

func extractConvertedCost(n *html.Node) int {
	div := find(n, "value")

	if div == nil {
		return 0
	}

	result := strings.TrimSpace(text(div))
	cost, err := strconv.Atoi(result)

	if err != nil {
		return 0
	}

	return cost
}


func ParseCard(page io.Reader) (Card, error) {
	doc, err := html.Parse(page)

	card := Card{}

	if err != nil {
		return card, err
	}

	var searchTheCity func(n *html.Node)

	searchTheCity = func(n *html.Node) {
		if hasId(n, prefix+"nameRow") {
			card.Name = extractName(n)
		}
		if hasId(n, prefix+"cmcRow") {
			card.ConvertedCost = extractConvertedCost(n)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			searchTheCity(c)
		}
	}

	searchTheCity(doc)

	return card, nil
}

func main() {
	fmt.Println("Foo")
}
