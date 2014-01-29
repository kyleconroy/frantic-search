package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var cards = []int{
	21382,
	189211,
	233056,
	212241,
}

func TestSearchResults(t *testing.T) {
	blob, _ := os.Open("fixtures/search.html")

	results, total, err := ParseSearch(blob)

	if err != nil {
		t.Fatal(err)
	}

	if total != 13993 {
		t.Errorf("Total search results should be 13993, not %d", total)
	}

	if len(results) != 100 {
		t.Errorf("Search results should have 1, not %d", len(results))
	}

	card := results[0]

	expected := SearchResult{Name: "Academy at Tolaria West", MultiverseId: 198073, Id: "d4d9a0bf18e13b326a3ca19c8283d28b"}

	if !reflect.DeepEqual(card, expected) {
		t.Errorf("Cards did not match: Got: \n%+v\ninstead of\n%+v", card, expected)
	}
}

func loadCard(id int) (Card, error) {
	path := fmt.Sprintf("fixtures/%d.json", id)
	blob, err := ioutil.ReadFile(path)

	if err != nil {
		return Card{}, err
	}

	var expected Card

	err = json.Unmarshal(blob, &expected)

	if err != nil {
		return Card{}, err
	}

	return expected, nil
}

func checkCards(t *testing.T, card, expected Card) {
	if !reflect.DeepEqual(card, expected) {
		card_json, _ := json.Marshal(card)
		expect_json, _ := json.Marshal(expected)

		t.Errorf("Cards did not match: Got: \n%s\ninstead of\n%s", string(card_json), string(expect_json))
	}
}

func TestCreatureCard(t *testing.T) {
	for _, id := range cards {
		path := fmt.Sprintf("fixtures/%d.html", id)
		file, err := os.Open(path)

		if err != nil {
			t.Errorf("%5d: Can't open %s", id, path)
			continue
		}

		expected, err := loadCard(id)

		if err != nil {
			t.Errorf("%5d: Couldn't load %d.json", id, id)
			continue
		}

		cards, err := ParseCards(file, id)
		checkCards(t, cards[0], expected)
	}
}

func TestSplitCards(t *testing.T) {
	split := func(path string, frontId int, backId int) {
		file, err := os.Open(path)

		if err != nil {
			t.Errorf("Can't open %s", path)
			return
		}

		front, err := loadCard(frontId)

		if err != nil {
			t.Errorf("Couldn't load %d.json", frontId)
			return
		}

		back, err := loadCard(backId)

		if err != nil {
			t.Errorf("Couldn't load %d.json", backId)
			return
		}

		cards, err := ParseCards(file, frontId)

		checkCards(t, cards[0], front)

		if len(cards) != 2 {
			t.Error("Split card page only returned one card")
			return
		}

		checkCards(t, cards[1], back)
	}

	split("fixtures/huntmaster.html", 262875, 262699)
	split("fixtures/standdeliver.html", 20574, 205740)
	split("fixtures/bushi.html", 78600, 78601)
}
