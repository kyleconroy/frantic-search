## Schema


class Artist:
	name = string

class Card:
	artist = key artist
	converted_cost =  int
	mana_cost =  string
	special = ['flip', 'double-faced', 'split']
	partner_card = key card
	name = string
	number = string
	multiverse_id = int
	rarity = ['common', 'uncommon', 'rare', 'mythic']
	types = array???
	subtypes = array???
	mark = ""
	rules_text = ""
	flavor_text = ""
	power = string
	toughness = string
	expansion = string 


## Printings Versus Cards

A card has some attributes that cannot change

- Name (key)
- ID (md5 of the English version of the card)
- Converted Cost
- Mana Cost
- Rules Text
- Power
- Toughness
- Partner Card
- Special
- Types
- Color Identity
- Subtypes

A printing has these attributes

- Number
- Artist
- Rarity
- FlavorText
- MultiverseId
- Mark
