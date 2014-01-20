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

func TestCreatureCard(t *testing.T) {
	for _, id := range cards {
		path := fmt.Sprintf("test/fixtures/%d.html", id)
		file, err := os.Open(path)

		if err != nil {
			t.Errorf("%5d: Can't open %s", id, path)
			continue
		}

		path = fmt.Sprintf("test/fixtures/%d.json", id)
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

		card, err := ParseCard(file, id)

		if !reflect.DeepEqual(card, expected) {

			card_json, _ := json.Marshal(card)
			expect_json, _ := json.Marshal(expected)

			t.Errorf("%5d: Cards did not match: Got: \n%s\ninstead of\n%s", id, string(card_json), string(expect_json))

		}
	}
}
