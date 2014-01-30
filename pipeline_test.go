package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestDeckboxAdd(t *testing.T) {
	oldCard := Card{Id: "f", Editions: []Edition{Edition{MultiverseId: 1}}}
	sameCard := Card{Id: "f", Editions: []Edition{Edition{MultiverseId: 1}}}
	updatedEdition := Card{Id: "f", Editions: []Edition{Edition{MultiverseId: 1, Set: "Foo"}}}
	newEdition := Card{Id: "f", Editions: []Edition{Edition{MultiverseId: 2, Set: "Bar"}}}
	newCard := Card{Id: "g", Editions: []Edition{Edition{MultiverseId: 1, Set: "Foo"}}}

	box := Deckbox{Cards: []Card{oldCard}}

	box.Add(sameCard)

	card := box.Cards[0]

	if len(card.Editions) != 1 {
		t.Fatalf("Card should only have one edition, not %d", len(card.Editions))
	}

	box.Add(updatedEdition)

	card = box.Cards[0]
	set := card.Editions[0].Set

	if len(card.Editions) != 1 {
		t.Fatalf("Card should only have one edition, not %d", len(card.Editions))
	}

	if set != "Foo" {
		t.Fatalf("Card edition should be 'Foo', not '%s'", set)
	}

	box.Add(newEdition)

	card = box.Cards[0]

	if len(card.Editions) != 2 {
		t.Fatalf("Card should have two editions, not %d", len(card.Editions))
	}

	box.Add(newCard)

	if len(box.Cards) != 2 {
		t.Fatalf("Box should have two cards, not %d", len(box.Cards))
	}

}

func TestDeckboxJSON(t *testing.T) {
	path := "fixtures/testbox.json"
	box := Deckbox{Cards: []Card{Card{Id: "FOO"}}}
	err := box.Flush(path)

	if err != nil {
		t.Fatal(err)
	}

	blob, err := ioutil.ReadFile(path)

	if err != nil {
		t.Fatal(err)
	}

	var cart Deckbox

	err = json.Unmarshal(blob, &cart)

	if err != nil {
		t.Fatal(err)
	}

	if len(cart.Cards) != 1 {
		t.Fatalf("Box was supposed to have one card, had %d", len(cart.Cards))
	}

	if cart.Cards[0].Id != "FOO" {
		t.Fatalf("Loaded an empty card??")
	}
}
