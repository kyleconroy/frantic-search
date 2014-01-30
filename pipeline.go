package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"sort"
	"sync"
	"time"
)

type Deckbox struct {
	Cards []Card
}

func (d *Deckbox) UnmarshalJSON(blob []byte) error {
	return json.Unmarshal(blob, &d.Cards)
}

func (d *Deckbox) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Cards)
}

func (d *Deckbox) Len() int {
	return len(d.Cards)
}

func (d *Deckbox) Swap(i, j int) {
	d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
}

func (d *Deckbox) Less(i, j int) bool {
	return d.Cards[i].Name < d.Cards[j].Name
}

func (d *Deckbox) Flush(path string) error {
	blob, err := json.Marshal(d)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, blob, 0644)

	if err != nil {
		return err
	}
	return nil
}

// Return a map of all card ids
func (d *Deckbox) IdSet() map[string]bool {
	set := map[string]bool{}

	for _, card := range d.Cards {
		set[card.Id] = true
	}

	return set
}

// Return a map of all Multiverse ids
func (d *Deckbox) MultiverseSet() map[int]bool {
	set := map[int]bool{}

	for _, card := range d.Cards {
		for _, edition := range card.Editions {
			set[edition.MultiverseId] = true
		}
	}

	return set
}

// Add a card to the deckbox
func (d *Deckbox) Add(newCard Card) error {
	if len(newCard.Editions) == 0 {
		return fmt.Errorf("%s has no editions", newCard.Name)
	}

	for i, c := range d.Cards {
		if c.Id == newCard.Id {
			for _, newe := range newCard.Editions {
				found := false
				for j, olde := range c.Editions {
					if olde.MultiverseId == newe.MultiverseId {
						found = true
						if olde.Set == "" && newe.Set != "" {
							d.Cards[i].Editions[j] = newe
						}
					}
				}

				if !found {
					d.Cards[i].Editions = append(newCard.Editions, newe)
				}
			}
			return nil
		}
	}

	d.Cards = append(d.Cards, newCard)
	return nil
}

func loadDeckBox(path string) (Deckbox, error) {
	blob, err := ioutil.ReadFile(path)

	if err != nil {
		log.Printf("WARNING: Couldn't open %s, creating new deckbox", path)
		blob = []byte(`[]`)
	}

	var box Deckbox

	err = json.Unmarshal(blob, &box)

	if err != nil {
		return box, err
	}

	return box, nil
}

func processSearchResults(seenCards map[string]bool, pageChan chan int, multiverseChan chan int) {
	var fetchGroup sync.WaitGroup

	log.Printf("Determining total number of pages")

	pages := TotalPages()

	log.Printf("Found %d pages", pages)

	for j := 0; j < pages; j++ {
		pageChan <- j
	}

	close(pageChan)

	// End Refactor
	log.Printf("Processing Gatherer search with concurrency 20")

	for j := 0; j < 20; j++ {
		fetchGroup.Add(1)
		go func(page int) {
			defer fetchGroup.Done()

			for {
				page, ok := <-pageChan

				if !ok {
					return
				}

				results, _, err := FetchSearch(page)

				if err != nil {
					log.Fatal(err)
				}

				toProcess := 0

				for _, result := range results {

					if !seenCards[result.Id] {
						toProcess += 1
						multiverseChan <- result.MultiverseId
					}
				}

				log.Printf("Found %d total cards on page %d, %d new", len(results), page, toProcess)
			}
		}(j)
	}

	go func() {
		fetchGroup.Wait()
		log.Printf("Close multiverse ID Channel")
		close(multiverseChan)
	}()
}

func processCards(multiverseChan chan int, cardChan chan Card) {
	// Start N go routines to go fetch and parse cards
	var parseGroup sync.WaitGroup

	log.Printf("Processing cards with concurrency of 100")

	for j := 0; j < 100; j++ {
		parseGroup.Add(1)
		go func() {
			defer parseGroup.Done()

			for {
				id, ok := <-multiverseChan

				if !ok {
					return
				}

				cards, err := FetchCards(id)

				if err != nil {
					log.Printf("ERROR Couldn't parse %d: %s", id, err)
					continue
				}

				for _, card := range cards {

					if card.Name == "" {
						log.Printf("ERROR No name found for %d", id)
						continue
					}

					cardChan <- card
				}
			}
		}()
	}

	go func() {
		parseGroup.Wait()
		close(cardChan)
	}()
}

// One go rotine pulls cards off the channel, adds them to the database
// And flushes it to memory
func saveCards(path string, box *Deckbox, cardChan chan Card) {
	count := 0
	for {
		card, ok := <-cardChan

		if !ok {
			log.Printf("FINISHED")
			sort.Sort(box)
			err := box.Flush(path)

			if err != nil {
				log.Fatal(err)
			}

			return
		}

		err := box.Add(card)

		if err != nil {
			log.Printf("ERROR Couldn't add card to database %s", card.Name)
			continue
		}

		count += 1

		if count >= 1000 {
			log.Printf("Added 1000 cards to the database")

			err := box.Flush(path)

			if err != nil {
				log.Fatal(err)
			}
			count = 0
		}
	}
}

func findEmptyEditions(box *Deckbox, multiverseChan chan int) {
	log.Printf("%d", len(box.Cards))
	count := 0
	for _, card := range box.Cards {
		for _, edition := range card.Editions {
			if edition.Set == "" {
				count += 1
				multiverseChan <- edition.MultiverseId
			}
		}
	}

	log.Printf("Found %d editions that need to be fetched", count)

	close(multiverseChan)
}

func processEditions(multiverseChan chan int, cardChan chan Card) {
	// Start N go routines to go fetch and parse cards
	var parseGroup sync.WaitGroup

	log.Printf("Processing editions with concurrency of 50")

	for j := 0; j < 50; j++ {
		parseGroup.Add(1)
		go func() {
			defer parseGroup.Done()

			for {
				id, ok := <-multiverseChan

				if !ok {
					return
				}

				cards, err := FetchCards(id)

				if err != nil {
					log.Printf("ERROR Couldn't parse %d: %s", id, err)
					continue
				}

				for _, card := range cards {

					if card.Name == "" {
						log.Printf("ERROR No name found for %d", id)
						continue
					}

					cardChan <- card
				}
			}
		}()
	}

	go func() {
		parseGroup.Wait()
		close(cardChan)
	}()
}

func saveEditions(path string, box Deckbox, editionChan chan Card) {
	for {
		time.Sleep(1000)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	path := flag.Arg(0)

	box, err := loadDeckBox(path)

	if err != nil {
		log.Fatal(err)
	}

	cardChannel := make(chan Card)
	editionChannel := make(chan Card)
	multiverseCardChannel := make(chan int, 15000)
	multiverseEditionChannel := make(chan int, 15000)
	pageChannel := make(chan int, 200)

	// Fetch all the cards
	go processSearchResults(box.IdSet(), pageChannel, multiverseCardChannel)
	go processCards(multiverseCardChannel, cardChannel)
	saveCards(path, &box, cardChannel)

	// Fetch all the editions
	go findEmptyEditions(&box, multiverseEditionChannel)
	go processEditions(multiverseEditionChannel, editionChannel)
	saveCards(path, &box, editionChannel)
}
