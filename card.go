package fb

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

/*
type CardQueue int

const (
	QueueNew CardQueue = iota
	QueueLearning
	QueueReview
)
*/

// Card represents a struct of card-related statistics and configuration.
type Card struct {
	bundleID   string
	noteID     string
	templateID uint32
	themeID    string
	modelID    uint32
	Rev        *string
	Created    time.Time
	Modified   time.Time
	Imported   *time.Time
	// 	Queue       CardQueue
	// 	Suspended   bool
	// 	Buried      bool
	// 	AutoBuried  bool
	Due        *Due
	Interval   *Interval
	EaseFactor float32
	// 	SRSFactor   float32
	ReviewCount int
	// 	LapseCount  int
}

type cardDoc struct {
	Type     string     `json:"type"`
	ID       string     `json:"_id"`
	Rev      *string    `json:"_rev,omitempty"`
	Created  time.Time  `json:"created"`
	Modified time.Time  `json:"modified"`
	Imported *time.Time `json:"imported,omitempty"`
	ModelID  string     `json:"model"`
	// 	Queue       CardQueue      `json:"state"`
	// 	Suspended   *bool          `json:"suspended,omitempty"`
	// 	Buried      *bool          `json:"buried,omitempty"`
	// 	AutoBuried  *bool          `json:"autoBuried,omitempty"`
	Due      *Due      `json:"due,omitempty"`
	Interval *Interval `json:"interval,omitempty,string"`
	// 	SRSFactor   *float32       `json:"srsFactor,omitempty"`
	// 	ReviewCount *int           `json:"reviewCount,omitempty"`
	// 	LapseCount  *int           `json:"lapseCount,omitempty"`
}

func parseID(id string) (string, string, uint32, error) {
	parts := strings.Split(strings.TrimPrefix(id, "card-"), ".")
	if len(parts) != 3 {
		return "", "", 0, errors.New("Invalid id format")
	}
	templateID, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", "", 0, errors.Wrap(err, "Invalid TemplateID")
	}
	return parts[0], parts[1], uint32(templateID), nil
}

// NewCard returns a new Card instance, with the requested id
func NewCard(theme string, model uint32, id string) (*Card, error) {
	bundleID, noteID, templateID, err := parseID(id)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing card ID")
	}
	return &Card{
		bundleID:   bundleID,
		noteID:     noteID,
		templateID: templateID,
		themeID:    theme,
		modelID:    model,
	}, nil
}

// MarshalJSON implements the json.Marshaler interface for the Card type.
func (c *Card) MarshalJSON() ([]byte, error) {
	doc := cardDoc{
		Type:     "card",
		ID:       c.DocID(),
		Rev:      c.Rev,
		Created:  c.Created,
		Modified: c.Modified,
		Imported: c.Imported,
		ModelID:  fmt.Sprintf("%s/%d", c.themeID, c.modelID),
		// 		Queue:    c.Queue,
		Due:      c.Due,
		Interval: c.Interval,
	}
	return json.Marshal(doc)
}

// UnmarshalJSON implements the json.Unmarshaler interface for the Card type.
func (c *Card) UnmarshalJSON(data []byte) error {
	doc := &cardDoc{}
	if err := json.Unmarshal(data, doc); err != nil {
		return err
	}
	if doc.Type != "card" {
		return errors.New("Invalid document type for card: " + doc.Type)
	}
	bundleID, noteID, templateID, err := parseID(doc.ID)
	if err != nil {
		return err
	}
	c.bundleID = bundleID
	c.noteID = noteID
	c.templateID = templateID
	c.Rev = doc.Rev
	c.Created = doc.Created
	c.Modified = doc.Modified
	c.Imported = doc.Imported
	model := strings.Split(doc.ModelID, "/")
	c.themeID = model[0]
	m, err := strconv.Atoi(model[1])
	if err != nil {
		return errors.Wrap(err, "invalid model ID")
	}
	c.modelID = uint32(m)
	// 	c.Queue = doc.Queue
	// 	if doc.Suspended != nil {
	// 		c.Suspended = *doc.Suspended
	// 	}
	// 	if doc.Buried != nil {
	// 		c.Buried = *doc.Buried
	// 	}
	// 	if doc.AutoBuried != nil {
	// 		c.AutoBuried = *doc.AutoBuried
	// 	}
	c.Due = doc.Due
	c.Interval = doc.Interval
	// 	if doc.SRSFactor != nil {
	// 		c.SRSFactor = *doc.SRSFactor
	// 	}
	// 	if doc.ReviewCount != nil {
	// 		c.ReviewCount = *doc.ReviewCount
	// 	}
	// 	if doc.LapseCount != nil {
	// 		c.LapseCount = *doc.LapseCount
	// 	}
	return nil
}

// Identity returns the identity of the card as a string.
func (c *Card) Identity() string {
	return fmt.Sprintf("%s.%s.%d", c.bundleID, c.noteID, c.templateID)
}

// SetRev sets the Card's _rev attribute
func (c *Card) SetRev(rev string) { c.Rev = &rev }

// DocID returns the Card's _id attribute
func (c *Card) DocID() string { return "card-" + c.Identity() }

// ImportedTime returns the Card's imported time, or nil
func (c *Card) ImportedTime() *time.Time { return c.Imported }

// ModifiedTime returns the Card's last modified time
func (c *Card) ModifiedTime() *time.Time { return &c.Modified }

// MergeImport attempts to merge i into c, returning true on success, or false
// if no merge was necessary.
func (c *Card) MergeImport(i interface{}) (bool, error) {
	existing := i.(*Card)
	if c.Identity() != existing.Identity() {
		return false, errors.New("IDs don't match")
	}
	if !c.Created.Equal(existing.Created) {
		return false, errors.New("Created timestamps don't match")
	}
	c.Rev = existing.Rev
	if c.Modified.After(existing.Modified) {
		// The new version is newer than the existing one, so update
		return true, nil
	}
	// The new version is older, so we need to use the version we just read
	c.Modified = existing.Modified
	c.Imported = existing.Imported
	return false, nil
}

// BundleID returns the card's BundleID
func (c *Card) BundleID() string {
	return "bundle-" + c.bundleID
}

// TemplateID returns the card's TemplateID
func (c *Card) TemplateID() uint32 {
	return c.templateID
}

// NoteID returns the card's NoteID
func (c *Card) NoteID() string {
	return "note-" + c.noteID
}
