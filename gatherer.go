package main

import (
	"code.google.com/p/go.net/html"
	"crypto/md5"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

const (
	prefixSingle = "#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_"
	prefixFront  = "#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl07_"
	prefixBack   = "#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl08_"
	prefixLeft   = "#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl09_"
	prefixRight  = "#ctl00_ctl00_ctl00_MainContent_SubContent_SubContent_ctl10_"
	gathererUrl  = "http://gatherer.wizards.com/Pages/Card/Details.aspx?multiverseid=%d"
	searchUrl    = "http://gatherer.wizards.com/Pages/Search/Default.aspx?output=compact&action=advanced&special=true&cmc=|>%%3d[0]|<%%3d[0]&page=%d"
)

type Card struct {
	Name           string    `json:"name"`
	Id             string    `json:"id"`
	Types          []string  `json:"types"`
	Subtypes       []string  `json:"subtypes,omitempty"`
	ConvertedCost  int       `json:"converted_cost"`
	ManaCost       string    `json:"mana_cost"`
	Special        string    `json:"special,omitempty"` //'flip', 'double-faced', 'split'
	PartnerCard    string    `json:"partner_card,omitempty"`
	RulesText      []string  `json:"rules_text"`
	ColorIndicator []string  `json:"color_indicator,omitempty"`
	Power          string    `json:"power,omitempty"`
	Toughness      string    `json:"toughness,omitempty"`
	Loyalty        int       `json:"loyalty,omitempty"`
	Editions       []Edition `json:"editions"`
}

type Edition struct {
	Set          string   `json:"set,omitempty"`
	Watermark    string   `json:"watermark,omitempty"`
	Rarity       string   `json:"rarity,omitempty"`
	Artist       string   `json:"artist,omitempty"`
	MultiverseId int      `json:"multiverse_id"`
	FlavorText   []string `json:"flavor_text,omitempty"`
	Number       string   `json:"number,omitempty"`
}

func (c Card) ImageURl() string {
	return "http://gatherer.wizards.com/Handlers/Image.ashx?multiverseid="
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
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15":
		return "{" + alt + "}"
	case "Phyrexian":
		return "{P}"
	case "Phyrexian Green":
		return "{G/P}"
	case "Phyrexian Red":
		return "{R/P}"
	case "Phyrexian Blue":
		return "{U/P}"
	case "Phyrexian Black":
		return "{B/P}"
	case "Phyrexian White":
		return "{W/P}"
	case "White or Blue":
		return "{W/U}"
	case "White or Black":
		return "{W/B}"
	case "Blue or Black":
		return "{U/B}"
	case "Blue or Red":
		return "{U/R}"
	case "Black or Red":
		return "{B/R}"
	case "Black or Green":
		return "{B/G}"
	case "Red or Green":
		return "{R/G}"
	case "Red or White":
		return "{R/W}"
	case "Green or White":
		return "{G/W}"
	case "Green or Blue":
		return "{G/U}"
	case "Two or White":
		return "{2/W}"
	case "Two or Blue":
		return "{2/U}"
	case "Two or Black":
		return "{2/B}"
	case "Two or Red":
		return "{2/R}"
	case "Two or Green ":
		return "{2/G}"
	case "Variable Colorless":
		return "{X}"
	case "Snow":
		return "{S}"
	case "Green":
		return "{G}"
	case "Red":
		return "{R}"
	case "Blue":
		return "{U}"
	case "Black":
		return "{B}"
	case "White":
		return "{W}"
	case "Tap":
		return "{T}"
	case "Untap":
		return "{Q}"
	case "[chaos]":
		return "{C}"
	}
	return ""
}

func FlattenWithSymbols(n *html.Node) string {
	text := ""
	if n.Type == html.TextNode {
		text += n.Data
	}

	if n.Type == html.ElementNode && n.Data == "img" {
		text += manaSymbol(Attr(n, "alt"))
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += FlattenWithSymbols(c)
	}

	return text
}

func extractManaCost(n *html.Node, prefix string) string {
	cost := ""
	for _, a := range FindAll(n, prefix+"manaRow .value img") {
		cost += manaSymbol(Attr(a, "alt"))
	}
	return cost
}

func extractPT(n *html.Node, prefix string) (string, string) {
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

func extractTypes(n *html.Node, prefix string) ([]string, []string) {
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

func extractRarity(n *html.Node, prefix string) string {
	if span, found := Find(n, prefix+"rarityRow .value span"); found {
		return Attr(span, "class")
	} else {
		return ""
	}
}

func extractColorIndicator(n *html.Node, pattern string) []string {
	div, found := Find(n, pattern+"colorIndicatorRow .value")

	if !found {
		return nil
	}

	colors := []string{}

	for _, color := range strings.Split(strings.TrimSpace(Flatten(div)), ",") {
		colors = append(colors, strings.ToLower(strings.TrimSpace(color)))
	}

	return colors
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

func extractId(n *html.Node, pattern string) int {
	img, found := Find(n, pattern)

	if !found {
		return 0
	}

	url, err := url.Parse(Attr(img, "src"))

	if err != nil {
		return 0
	}

	mid := url.Query().Get("multiverseid")
	multiverseid, err := strconv.Atoi(mid)

	if err != nil {
		return 0
	}

	return multiverseid
}

func extractText(n *html.Node, pattern string) []string {
	rules := []string{}
	for _, node := range FindAll(n, pattern) {
		rule := strings.TrimSpace(FlattenWithSymbols(node))

		if rule != "" {
			rules = append(rules, rule)
		}
	}
	return rules
}

func extractResultSize(n *html.Node) int {
	div, found := Find(n, "#ctl00_ctl00_ctl00_MainContent_SubContent_SubContentHeader_searchTermDisplay")

	if !found {
		return 0
	}

	parts := strings.Split(Flatten(div), "(")

	if len(parts) != 2 {
		return 0
	}

	result := strings.Replace(parts[1], ")", "", 1)

	count, err := strconv.Atoi(result)

	if err != nil {
		return 0
	}

	return count
}

func extractMultiverseIds(n *html.Node, id int, pattern string) []int {
	ids := []int{}
    found := false

	for _, a := range FindAll(n, pattern) {

		u, err := url.Parse(Attr(a, "href"))

		if err != nil {
			continue
		}

		mid := u.Query().Get("multiverseid")
		multiverseid, err := strconv.Atoi(mid)

		if err != nil {
			continue
		}

        if multiverseid == id {
                found = true
        }

		ids = append(ids, multiverseid)
	}

    if !found {
		ids = append(ids, id)
}

	sort.Ints(ids)
	return ids
}

func hash(in string) string {
	h := md5.New()
	io.WriteString(h, in)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func parseCard(doc *html.Node, prefix string) Card {
	card := Card{}
	card.Name = extractString(doc, prefix+"nameRow .value")
	card.ManaCost = extractManaCost(doc, prefix)
	card.Id = hash(card.Name + card.ManaCost)
	card.ConvertedCost = extractInt(doc, prefix+"cmcRow .value")
	card.RulesText = extractText(doc, prefix+"textRow .value .cardtextbox")
	card.Loyalty = extractInt(doc, prefix+"ptRow .value")
	card.ColorIndicator = extractColorIndicator(doc, prefix)
	card.Types, card.Subtypes = extractTypes(doc, prefix)
	card.Power, card.Toughness = extractPT(doc, prefix)

	edition := Edition{}
	edition.Number = extractString(doc, prefix+"numberRow .value")
	edition.Artist = extractString(doc, prefix+"artistRow .value")
	edition.Set = extractString(doc, prefix+"setRow .value")
	edition.FlavorText = extractText(doc, prefix+"flavorRow .value .cardtextbox")
	edition.Rarity = extractRarity(doc, prefix)
	edition.Watermark = extractString(doc, prefix+"markRow .value")
	edition.MultiverseId = extractId(doc, prefix+"cardImage")

	// Gross?
	ids := extractMultiverseIds(doc, edition.MultiverseId, prefix+"otherSetsValue a")

		for _, id := range ids {
			if edition.MultiverseId == id {
				card.Editions = append(card.Editions, edition)
			} else {
				card.Editions = append(card.Editions, Edition{MultiverseId: id})
			}
		}

	return card
}

func FetchCards(multiverseId int) ([]Card, error) {
	url := fmt.Sprintf(gathererUrl, multiverseId)
	resp, err := http.Get(url)

	if err != nil {
		return []Card{}, err
	}

	cards, err := ParseCards(resp.Body, multiverseId)

	if err != nil {
		return []Card{}, err
	}

	return cards, nil
}

func ParseCards(page io.Reader, multiverseid int) ([]Card, error) {
	doc, err := html.Parse(page)

	if err != nil {
		return []Card{Card{}}, err
	}

	var prefixA, prefixB, special string

	if _, found := Find(doc, prefixLeft+"cardImage"); found {
		prefixA = prefixLeft
		prefixB = prefixRight
		special = "split"
	} else if _, found := Find(doc, prefixBack+"colorIndicatorRow"); found {
		prefixA = prefixFront
		prefixB = prefixBack
		special = "double-faced"
	} else if _, found := Find(doc, prefixFront+"cardImage"); found {
		prefixA = prefixFront
		prefixB = prefixBack
		special = "flip"
	}

	if special != "" {
		a := parseCard(doc, prefixA)
		b := parseCard(doc, prefixB)

		a.PartnerCard = b.Id
		b.PartnerCard = a.Id

		a.Special = special
		b.Special = special
		return []Card{a, b}, nil
	} else {
		return []Card{parseCard(doc, prefixSingle)}, nil
	}
}

type SearchResult struct {
	Name         string
	Id           string
	MultiverseId int
}

func TotalPages() int {
	_, total, err := FetchSearch(0)

	if err != nil {
		return 0
	}
	return int(math.Ceil(float64(total) / float64(100)))
}

func FetchSearch(page int) ([]SearchResult, int, error) {
	url := fmt.Sprintf(searchUrl, 0)

	resp, err := http.Get(url)

	if err != nil {
		return []SearchResult{}, 0, err
	}

	results, total, err := ParseSearch(resp.Body)

	if err != nil {
		return []SearchResult{}, 0, err
	}

	return results, total, nil
}

func ParseSearch(page io.Reader) ([]SearchResult, int, error) {
	doc, err := html.Parse(page)

	results := []SearchResult{}

	if err != nil {
		return results, 0, err
	}

	for _, a := range FindAll(doc, ".cardItem .name a") {
		// Parse the link to get the id
		url, err := url.Parse(Attr(a, "href"))

		if err != nil {
			return results, 0, err
		}

		mid := url.Query().Get("multiverseid")

		multiverseid, err := strconv.Atoi(mid)

		if err != nil {
			return results, 0, err
		}

		name := strings.TrimSpace(Flatten(a))

		results = append(results, SearchResult{
			Name:         name,
			Id:           hash(name),
			MultiverseId: multiverseid,
		})
	}

	return results, extractResultSize(doc), nil
}
