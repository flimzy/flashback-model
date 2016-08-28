package fb

import (
	"encoding/json"
	"errors"

	"github.com/pborman/uuid"
)

func isValidID(id string) bool {
	return uuid.Parse(id) != nil
}

// User repressents a user of Flashback
type User struct {
	ID       DbID
	uuid     uuid.UUID
	Rev      *string
	Username string
	Password string
	Salt     string
	UserType string
	FullName *string
	Email    *string
}

type userDoc struct {
	Type     string  `json:"type"`
	ID       DbID    `json:"_id"`
	Rev      *string `json:"_rev,omitempty"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	Salt     string  `json:"salt"`
	UserType string  `json:"userType"`
	FullName *string `json:"fullname,omitempty"`
	Email    *string `json:"email,omitempty"`
}

// func CreateUser(username string) (*User, error) {
// 	u := &User{}
// 	u.ID = NewID("user", uuid.NewRandom())
// 	return u, nil
// }

// NewUser returns a new User object, based on the provided UUID and username.
func NewUser(id uuid.UUID, username string) (*User, error) {
	uid, err := NewDbID("user", id)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:       uid,
		uuid:     id,
		Username: username,
	}, nil
}

// NewUserStub returns a new stub (bare bones) User object.
func NewUserStub(id string) (*User, error) {
	userID, err := ParseDbID("user", id)
	if err != nil {
		return nil, err
	}
	userUUID := uuid.UUID(userID.RawID())
	return NewUser(userUUID, "")
}

// RandomUser creates a new random user. For testing. Please don't use this!
// FIXME: Remove this after testing.
func RandomUser() *User {
	u, _ := NewUser(uuid.NewRandom(), "randomuser")
	return u
}

// UUID returns the UUID underlying the User object.
func (u *User) UUID() uuid.UUID {
	return u.uuid
}

// MarshalJSON implements the json.Marshaler interface for the User type.
func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(userDoc{
		Type:     "user",
		ID:       u.ID,
		Rev:      u.Rev,
		Username: u.Username,
		Password: u.Password,
		Salt:     u.Salt,
		UserType: u.UserType,
		FullName: u.FullName,
		Email:    u.Email,
	})
}

// UnmarshalJSON implements the json.Unmarshaler interface for the User type.
func (u *User) UnmarshalJSON(data []byte) error {
	doc := userDoc{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return err
	}
	if doc.Type != "user" {
		return errors.New("Invalid document type for user")
	}
	// 	id, err := b64encoder.DecodeString(doc.ID.Identity())
	// 	if err != nil {
	// 		return errors.New(doc.ID.Identity() + " is not a valid UUID")
	// 	}
	u.ID = doc.ID
	u.uuid = u.ID.RawID()
	u.Rev = doc.Rev
	u.Username = doc.Username
	u.Salt = doc.Salt
	u.FullName = doc.FullName
	u.Email = doc.Email

	return nil
}

// Fleshened returns true if the user object has been fleshened
func (u *User) Fleshened() bool {
	return u.Username != ""
}

// Equal returns true if the two users are equal.
func (u *User) Equal(id uuid.UUID) bool {
	return uuid.Equal(u.uuid, id)
}
