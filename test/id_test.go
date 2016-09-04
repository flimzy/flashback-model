package test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/FlashbackSRS/flashback-model"
	. "github.com/FlashbackSRS/flashback-model/test/util"
)

var frozenDocID = []byte(`"note-VGVzdCBOb3Rl"`)
var frozenDbID = []byte(`"user-krsxg5bakvzwk4q"`)

func TestDocID(t *testing.T) {
	id, err := fb.NewDocID("note", []byte("Test Note"))
	if err != nil {
		t.Fatalf("Error creating B64 ID: %s\n", err)
	}
	if id.String() != "note-VGVzdCBOb3Rl" {
		t.Fatalf("Error stringifying note id. Got %s\n", id.String())
	}
	if id.Identity() != "VGVzdCBOb3Rl" {
		t.Fatalf("Unexpected identity for note id. Got %s\n", id.Identity())
	}
	JSONDeepEqual(t, "Create DocID", Marshal(t, "Create ID1", id), frozenDocID)

	id2 := fb.DocID{}
	if err := json.Unmarshal(frozenDocID, &id2); err != nil {
		t.Fatalf("Error thawing DocID: %s", err)
	}
	JSONDeepEqual(t, "Thawed DocID", Marshal(t, "Thaw DocID", id2), frozenDocID)

	if !reflect.DeepEqual(id, id2) {
		PrintDiff(id2, id)
		t.Fatalf("Thawed and created B64 IDs don't match")
	}
}

func TestDbID(t *testing.T) {
	id, err := fb.NewDbID("user", []byte("Test User"))
	if err != nil {
		t.Fatalf("Error creating Hex ID: %s\n", err)
	}
	if id.String() != "user-krsxg5bakvzwk4q" {
		t.Fatalf("Error stringifying note id. Got %s\n", id.String())
	}
	if id.Identity() != "krsxg5bakvzwk4q" {
		t.Fatalf("Unexpected identity for note id. Got %s\n", id.Identity())
	}
	JSONDeepEqual(t, "Create Hex ID", Marshal(t, "Create ID1", id), frozenDbID)

	id2 := fb.DbID{}
	if err := json.Unmarshal(frozenDbID, &id2); err != nil {
		t.Fatalf("Error thawing Hex ID: %s", err)
	}
	JSONDeepEqual(t, "Thawed Hex ID", Marshal(t, "Thaw Hex ID", id2), frozenDbID)

	if !reflect.DeepEqual(id, id2) {
		PrintDiff(id2, id)
		t.Fatalf("Thawed and created Hex IDs don't match")
	}
}

func TestID(t *testing.T) {
	id, err := fb.NewDbID("user", []byte("User Bob"))
	if err != nil {
		t.Fatalf("Error creating user: %s\n", err)
	}
	id2, err := fb.ParseDbID(id.String())
	if err != nil {
		t.Fatalf("We can't even parse the IDs we generate: %s", err.Error())
	}
	if id.Type() != id2.Type() {
		t.Errorf("Types: %s != %s\n", id.Type(), id2.Type())
	}
	if id.Identity() != id2.Identity() {
		t.Errorf("ID: %x != %x\n", id.Identity(), id2.Identity())
	}
}

func TestID2(t *testing.T) {
	id, err := fb.NewDbID("user", []byte("User Bob"))
	if err != nil {
		t.Fatalf("Error creating user: %s\n", err)
	}
	id2, err := fb.ParseDbID(id.String())
	if err != nil {
		t.Fatalf("We can't even parse the IDs we generate: %s", err.Error())
	}
	if id.Type() != id2.Type() {
		t.Errorf("Types: %s != %s\n", id.Type(), id2.Type())
	}
	if id.Identity() != id2.Identity() {
		t.Errorf("ID: %x != %x\n", id.Identity(), id2.Identity())
	}
}
