package main

import (
	"testing"
)

func TestDeckboxAdd(t *testing.T) {
	oldCard := Card{Id: "f", Editions: []Edition{Edition{MultiverseId: 1}}}
	updatedEdition := Card{Id: "f", Editions: []Edition{Edition{MultiverseId: 1, Set: "Foo"}}}
	newEdition := Card{Id: "f", Editions: []Edition{Edition{MultiverseId: 2, Set: "Bar"}}}
	newCard := Card{Id: "g", Editions: []Edition{Edition{MultiverseId: 1, Set: "Foo"}}}

	box := Deckbox{Cards: []Card{oldCard}}

	box.Add(updatedEdition)

	card := box.Cards[0]
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
