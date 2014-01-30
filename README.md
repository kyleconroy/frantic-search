![Frantic Search](http://www.jeffmiracola.com/images/art/paintings/frantic_searchBIG.jpg)

Scrape Gatherer to build a local database of card information.

This project is a personal project. If you want card and image data, I suggest
[mtgjson.com](http://mtgjson.com) and [mtgimages.com](http://mtgimages.com).
The data sets at those website are more inclusive than this one.

## Usage

    ./frantic cards.json

This will create a JSON file with all metadata for all Magic: The Gathering
cards. It should take about 5 minutes. 

## Latest JSON

- [cards.json.zip (2.1mb)](https://github.com/kyleconroy/frantic-search/releases/download/BTG/cards.json.zip)

## Card Structure

Here is a sample card structure. A card's unique ID is generated by computing
the MD5 of concatenating the name and mana cost of the card.

```js
{
    "id": "bf52be96ef89803d43bc8df52ecece62",
    "name": "Bestial Menace",
    "mana_cost": "{3}{G}{G}",
    "converted_cost": 5,
    "rules_text": [
        "Put a 1/1 green Snake creature token, a 2/2 green Wolf creature token, and a 3/3 green Elephant creature token onto the battlefield."
    ],
    "types": [
        "sorcery"
    ]
    "editions": [
        {
            "artist": "Andrew Robinson",
            "flavor_text": [
                "\"My battle cry reaches ears far keener than yours.\"",
                "\u2014Saidah, Joraga hunter"
            ],
            "multiverse_id": 197843,
            "number": "97",
            "rarity": "uncommon",
            "set": "Worldwake"
        },
        {
            "artist": "Andrew Robinson",
            "flavor_text": [
                "\"My battle cry reaches ears far keener than yours.\"",
                "\u2014Saidah, Joraga hunter"
            ],
            "multiverse_id": 247535,
            "number": "144",
            "rarity": "uncommon",
            "set": "Magic: The Gathering-Commander"
        }
    ]
}
```

## Testing

    go test -v

## License and Copyright

Card names and text are all copyright Wizards of the Coast.

This website is not affiliated with Wizards of the Coast in any way.

I am providing the JSON files under the public domain license.

