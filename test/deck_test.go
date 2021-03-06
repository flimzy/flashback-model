package test

import (
	"encoding/json"
	"testing"

	"github.com/flimzy/testify/require"

	"github.com/FlashbackSRS/flashback-model"
)

var frozenDeck = []byte(`
{
    "type": "deck",
    "_id": "deck-VGVzdCBEZWNr",
    "created": "2016-07-31T15:08:24.730156517Z",
    "modified": "2016-07-31T15:08:24.730156517Z",
    "imported": "2016-08-02T15:08:24.730156517Z",
    "name": "Test Deck",
    "description": "Deck for testing",
    "cards": []
}
`)

func TestDecks(t *testing.T) {
	require := require.New(t)
	d, err := fb.NewDeck([]byte("Test Deck"))
	require.Nil(err, "Error creating deck: %s", err)

	name := "Test Deck"
	d.Name = &name
	descr := "Deck for testing"
	d.Description = &descr
	d.Created = now
	d.Modified = now
	imp := now.AddDate(0, 0, 2)
	d.Imported = &imp
	require.MarshalsToJSON(frozenDeck, d, "Create Deck")

	d2 := &fb.Deck{}
	err = json.Unmarshal(frozenDeck, d2)
	require.Nil(err, "Error thawing deck: %s", err)
	require.MarshalsToJSON(frozenDeck, d2, "Thawed Deck")

	require.DeepEqual(d, d2, "Thawed vs. Created Decks")
}

var frozenExistingDeck = []byte(`
{
    "type": "deck",
    "_id": "deck-VGVzdCBEZWNr",
    "_rev": "1-6e1b6fb5352429cf3013eab5d692aac8",
    "created": "2016-07-31T15:08:24.730156517Z",
    "modified": "2016-07-15T15:07:24.730156517Z",
    "imported": "2016-08-01T15:08:24.730156517Z",
    "name": "Test deck",
    "description": "Deck for testing",
    "cards": []
}
`)

var frozenMergedDeck = []byte(`
{
    "type": "deck",
    "_id": "deck-VGVzdCBEZWNr",
    "_rev": "1-6e1b6fb5352429cf3013eab5d692aac8",
    "created": "2016-07-31T15:08:24.730156517Z",
    "modified": "2016-07-31T15:08:24.730156517Z",
    "imported": "2016-08-02T15:08:24.730156517Z",
    "name": "Test Deck",
    "description": "Deck for testing",
    "cards": []
}
`)

func TestDeckMergeImport(t *testing.T) {
	require := require.New(t)
	d := &fb.Deck{}
	err := json.Unmarshal(frozenDeck, d)
	require.Nil(err, "Error thawing deck: %s", err)

	e := &fb.Deck{}
	err = json.Unmarshal(frozenExistingDeck, e)
	require.Nil(err, "Error thawing Existing Deck: %s", err)

	changed, err := d.MergeImport(e)
	require.Nil(err, "Error merging Deck: %s", err)
	require.True(changed, "Deck changed during merge")
	require.MarshalsToJSON(frozenMergedDeck, d, "Merged Deck")
}
