package fb

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/pkg/errors"
)

// CardCollection represents a collection of cards, which make up a deck.
type CardCollection struct {
	col map[string]struct{}
}

// MarshalJSON fulfills the json.Marshaler interface for the CardCollection type.
func (cc *CardCollection) MarshalJSON() ([]byte, error) {
	ids := make([]string, 0, len(cc.col))
	for id := range cc.col {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return json.Marshal(ids)
}

// NewCardCollection returns a new, empty CardCollection.
func NewCardCollection() *CardCollection {
	return &CardCollection{
		col: make(map[string]struct{}),
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface for the CardCollection
// type.
func (cc *CardCollection) UnmarshalJSON(data []byte) error {
	var ids []string
	if err := json.Unmarshal(data, &ids); err != nil {
		return errors.Wrap(err, "failed to unmarshal CardCollection")
	}
	cc.col = make(map[string]struct{})
	for _, id := range ids {
		cc.col[id] = struct{}{}
	}
	return nil
}

// All returns an array containing the list of all card IDs in the deck
func (cc *CardCollection) All() []string {
	ids := make([]string, 0, len(cc.col))
	for id := range cc.col {
		ids = append(ids, id)
	}
	return ids
}

// Deck represents a Flashback Deck
type Deck struct {
	ID          DocID
	Rev         *string
	Created     time.Time
	Modified    time.Time
	Imported    *time.Time
	Name        *string
	Description *string
	Cards       *CardCollection
}

type deckDoc struct {
	Type        string          `json:"type"`
	ID          DocID           `json:"_id"`
	Rev         *string         `json:"_rev,omitempty"`
	Created     time.Time       `json:"created"`
	Modified    time.Time       `json:"modified"`
	Imported    *time.Time      `json:"imported,omitempty"`
	Name        *string         `json:"name,omitempty"`
	Description *string         `json:"description,omitempty"`
	Cards       *CardCollection `json:"cards,omitempty"`
}

/*
type Deck struct {
	ConfigID                ID                `json:"conf"`             // ID of option group from dconf in `col` table
}

type DeckConfig struct {
	ID               ID                `json:"id"`       // Deck ID
	Name             string            `json:"name"`     // Deck Name
	ReplayAudio      bool              `json:"replayq"`  // When answer shown, replay both question and answer audio
	ShowTimer        BoolInt           `json:"timer"`    // Show answer timer
	MaxAnswerSeconds int               `json:"maxTaken"` // Ignore answers that take longer than this many seconds
	Modified         *TimestampSeconds `json:"mod"`      // Modified timestamp
	AutoPlay         bool              `json:"autoplay"` // Automatically play audio
	Lapses           struct {
		LeechFails      int               `json:"leechFails"`  // Leech threshold
		MinimumInterval DurationDays      `json:"minInt"`      // Minimum interval in days
		LeechAction     LeechAction       `json:"leechAction"` // Leech action: Suspend or Tag Only
		Delays          []DurationMinutes `json:"delays"`      // Steps in minutes
		NewInterval     float32           `json:"mult"`        // New Interval Multiplier
	} `json:"lapse"`
	Reviews struct {
		PerDay           int          `json:"perDay"` // Maximum reviews per day
		Fuzz             float32      `json:"fuzz"`   // Apparently not used?
		IntervalModifier float32      `json:"ivlFct"` // Interval modifier (fraction)
		MaxInterval      DurationDays `json:"maxIvl"` // Maximum interval in days
		EasyBonus        float32      `json:ease4"`   // Easy bonus
		Bury             bool         `json:"bury"`   // Bury related reviews until next day
	} `json:"rev"`
	New struct {
		PerDay        int               `json:"perDay"`        // Maximum new cards per day
		Delays        []DurationMinutes `json:"delays"`        // Steps in minutes
		Bury          bool              `json:"bury"`          // Bury related cards until the next day
		Separate      bool              `json:"separate"`      // Unused??
		Intervals     [3]DurationDays   `json:"ints"`          // Intervals??
		InitialFactor float32           `json:"initialFactor"` // Starting Ease
		Order         NewCardOrder      `json:"order"`         // New card order: Random, or order added
	} `json:"new"`
}
*/

// NewDeck creates a new Deck with the provided id.
func NewDeck(id []byte) (*Deck, error) {
	did, err := NewDocID("deck", id)
	if err != nil {
		return nil, err
	}
	return &Deck{
		ID:    did,
		Cards: NewCardCollection(),
	}, nil
}

// MarshalJSON implements the json.Marshaler interface for the Deck type.
func (d *Deck) MarshalJSON() ([]byte, error) {
	return json.Marshal(deckDoc{
		Type:        "deck",
		ID:          d.ID,
		Rev:         d.Rev,
		Created:     d.Created,
		Modified:    d.Modified,
		Imported:    d.Imported,
		Name:        d.Name,
		Description: d.Description,
		Cards:       d.Cards,
	})
}

// AddCard adds the provided card to the deck.
func (d *Deck) AddCard(cardID string) error {
	d.Cards.col[cardID] = struct{}{}
	return nil
}

// UnmarshalJSON fulfills the json.Unmarshaler interface for the Deck type.
func (d *Deck) UnmarshalJSON(data []byte) error {
	doc := &deckDoc{}
	if err := json.Unmarshal(data, doc); err != nil {
		return errors.Wrap(err, "failed to unmarshal Deck")
	}
	if doc.Type != "deck" {
		return errors.New("Invalid document type for deck: " + doc.Type)
	}
	d.ID = doc.ID
	d.Rev = doc.Rev
	d.Created = doc.Created
	d.Modified = doc.Modified
	d.Imported = doc.Imported
	d.Name = doc.Name
	d.Description = doc.Description
	d.Cards = doc.Cards

	return nil
}

// SetRev sets the Deck's _rev attribute.
func (d *Deck) SetRev(rev string) { d.Rev = &rev }

// DocID returns the Deck's _id attribute.
func (d *Deck) DocID() string { return d.ID.String() }

// ImportedTime returns the time the Deck was imported, or nil.
func (d *Deck) ImportedTime() *time.Time { return d.Imported }

// ModifiedTime returns the time the Deck was last modified.
func (d *Deck) ModifiedTime() *time.Time { return &d.Modified }

// MergeImport attempts to merge i into d, returning true on success, or false
// if no merge was necessary.
func (d *Deck) MergeImport(i interface{}) (bool, error) {
	existing := i.(*Deck)
	if !d.ID.Equal(&existing.ID) {
		return false, errors.New("IDs don't match")
	}
	if !d.Created.Equal(existing.Created) {
		return false, errors.New("Created timestamps don't match")
	}
	d.Rev = existing.Rev
	if d.Modified.After(existing.Modified) {
		// The new version is newer than the existing one, so update
		return true, nil
	}
	// The new version is older, so we need to use the version we just read
	d.Modified = existing.Modified
	d.Imported = existing.Imported
	d.Name = existing.Name
	d.Description = existing.Description
	d.Cards = existing.Cards
	return false, nil
}
