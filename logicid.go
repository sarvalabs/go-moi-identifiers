package identifiers

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// LogicID is a unique hex-encoded string identifier for logics in the MOI Protocol.
// It encodes within itself key properties about the logic itself such as the address
// at which it deployed, the nature of it state definitions and its edition.
//
// The Logic ID Standard is an extensible standard and is compliant with the protocol specification
// for it at https://sarvalabs.notion.site/Logic-ID-Standard-174a2cc6e3dc42e4bbf4dd708af0cd03?pvs=4
type LogicID string

// LogicIdentifier is a representation of LogicID that allows access to
// the values encoded within it such as the address or identifier version
//
// Each version of the LogicID standard is available
// as its own type that implements this interface.
type LogicIdentifier interface {
	// Version returns the version of the LogicID standard for the LogicIdentifier
	Version() int
	// Address returns the 32-byte address associated with the LogicID
	Address() Address
	// LogicID returns the LogicIdentifier in its encoded representation as a LogicID
	LogicID() LogicID

	// Edition returns the edition of the logic from the LogicID
	Edition() uint64
	// AssetLogic returns whether the logic is associated with some Asset
	AssetLogic() bool

	// HasPersistentState returns whether the logic has some persistent state definition
	HasPersistentState() bool
	// HasEphemeralState returns whether the logic has some ephemeral state definition
	HasEphemeralState() bool
	// HasInteractableSites returns whether the logic has some interactable callsites
	HasInteractableSites() bool
}

// NewLogicID generates a LogicID from some arbitrary string,
// validating it in the process. It is version agnostic
func NewLogicID(id string) (LogicID, error) {
	logic := LogicID(id)

	// Attempt to generate an identifier from the LogicID
	// This will fail if the LogicID is invalid in any way
	if _, err := logic.Identifier(); err != nil {
		return "", err
	}

	return logic, nil
}

// Bytes returns the LogicID as a []byte after being
// decoded from its hexadecimal string representation.
// Panics if the LogicID is not a valid hex string.
func (logic LogicID) Bytes() []byte {
	return must(decodeHexString(string(logic)))
}

// String returns the LogicID as a string.
// Implements the fmt.Stringer interface for LogicID.
func (logic LogicID) String() string {
	return "0x" + string(logic)
}

// Address returns the Address of the LogicID.
// Returns NilAddress if the LogicID is invalid.
// The LogicID standard expects the address to ALWAYS be the last 32 bytes
func (logic LogicID) Address() Address {
	// Error if length is too short
	if len(logic) < 64 {
		return NilAddress
	}

	// Trim the last 64 characters (32 bytes)
	addr := string(logic[len(logic)-64:])
	// Assertively decode into an Address
	return must(NewAddressFromHex(addr))
}

// Identifier returns a LogicIdentifier for the LogicID.
//
// This decodes the LogicID from its simple hex-encoded string
// format into a representation appropriate for the version of
// the LogicID to allow access to all the encoded fields within it.
//
// It can also be used to verify the integrity of the LogicID
func (logic LogicID) Identifier() (LogicIdentifier, error) {
	id, err := decodeHexString(string(logic))
	if err != nil {
		return nil, errors.Wrap(err, "invalid logic ID")
	}

	// We verify that there is at least 1 byte, so that
	// we can safely access the 0th index in the byte slice
	if len(id) < 1 {
		return nil, errors.New("invalid logic ID: missing version prefix")
	}

	// Determine the version of the LogicID and decode
	switch version := int(id[0] & 0xF0); version {
	case 0:
		return decodeLogicIDv0(id)
	default:
		return nil, errors.Errorf("invalid logic ID: unsupported version: %v", version)
	}
}

// MarshalText implements the encoding.TextMarshaler interface for LogicID
func (logic LogicID) MarshalText() ([]byte, error) {
	return []byte(logic.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for LogicID
func (logic *LogicID) UnmarshalText(text []byte) error {
	// Assert that the 0x prefix exists
	if !has0xPrefixBytes(text) {
		return ErrMissing0xPrefix
	}

	// Trim the 0x prefix
	text = trim0xPrefixBytes(text)
	// Generate an identifier for the LogicID
	if _, err := LogicID(text).Identifier(); err != nil {
		return err
	}

	*logic = LogicID(text)

	return nil
}

// MarshalJSON implements the json.Marshaler interface for LogicID
func (logic LogicID) MarshalJSON() ([]byte, error) {
	return json.Marshal(logic.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for LogicID
func (logic *LogicID) UnmarshalJSON(data []byte) error {
	var decoded string

	// Decode the JSON data into a string
	if err := json.Unmarshal(data, &decoded); err != nil {
		return err
	}

	// Assert that the 0x prefix exists
	if !has0xPrefixString(decoded) {
		return ErrMissing0xPrefix
	}

	// Trim the 0x prefix
	decoded = trim0xPrefixString(decoded)
	// Generate an identifier for the AssetID
	if _, err := LogicID(decoded).Identifier(); err != nil {
		return err
	}

	*logic = LogicID(decoded)

	return nil
}
