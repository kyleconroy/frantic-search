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
}

func TestSearchResults(t *testing.T) {
	blob, _ := os.Open("fixtures/search.html")

	results, err := ParseSearch(blob)

	if err != nil {
		t.Fatal(err)
	}

	if len(results) != 100 {
		t.Fatalf("Search results should have 1, not %d", len(results))
	}

	card := results[0]

	expected := SearchResult{Name:"Academy at Tolaria West", MultiverseId:198073, Id:"d4d9a0bf18e13b326a3ca19c8283d28b"}

	if !reflect.DeepEqual(card, expected) {
		t.Errorf("Cards did not match: Got: \n%+v\ninstead of\n%+v", card, expected)
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

		path = fmt.Sprintf("fixtures/%d.json", id)
		blob, err := ioutil.ReadFile(path)

		if err != nil {
			t.Errorf("%5d: Can't open %s", id, path)
			continue
		}

		var expected Card

		err = json.Unmarshal(blob, &expected)

		if err != nil {
			t.Errorf("%5d: Couldn't load %s", id, path)
		}

		cards, err := ParseCards(file, id)

		if !reflect.DeepEqual(cards[0], expected) {

			card_json, _ := json.Marshal(cards[0])
			expect_json, _ := json.Marshal(expected)

			t.Errorf("%5d: Cards did not match: Got: \n%s\ninstead of\n%s", id, string(card_json), string(expect_json))

		}
	}
}
