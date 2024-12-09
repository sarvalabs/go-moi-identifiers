package identifiers

//// LogicID is a unique hex-encoded string identifier for logics in the MOI Protocol.
//// It encodes within itself key properties about the logic itself such as the address
//// at which it deployed, the nature of it state definitions and its edition.
////
//// The Logic ID Standard is an extensible standard and is compliant with the protocol specification
//// for it at https://sarvalabs.notion.site/Logic-ID-Standard-174a2cc6e3dc42e4bbf4dd708af0cd03?pvs=4
//type LogicID string
//
//// LogicIdentifier is a representation of LogicID that allows access to
//// the values encoded within it such as the address or identifier version
////
//// Each version of the LogicID standard is available
//// as its own type that implements this interface.
//type LogicIdentifier interface {
//	// Version returns the version of the LogicID standard for the LogicIdentifier
//	Version() int
//	// Address returns the 32-byte address associated with the LogicID
//	Address() Address
//	// LogicID returns the LogicIdentifier in its encoded representation as a LogicID
//	LogicID() LogicID
//
//	// Edition returns the edition of the logic from the LogicID
//	Edition() uint64
//	// AssetLogic returns whether the logic is associated with some Asset
//	AssetLogic() bool
//
//	// HasPersistentState returns whether the logic has some persistent state definition
//	HasPersistentState() bool
//	// HasEphemeralState returns whether the logic has some ephemeral state definition
//	HasEphemeralState() bool
//	// HasInteractableSites returns whether the logic has some interactable callsites
//	HasInteractableSites() bool
//}
//
//// NewLogicID generates a LogicID from some arbitrary string,
//// validating it in the process. It is version agnostic
//func NewLogicID(id string) (LogicID, error) {
//	logic := LogicID(id)
//
//	// Attempt to generate an identifier from the LogicID
//	// This will fail if the LogicID is invalid in any way
//	if _, err := logic.Identifier(); err != nil {
//		return "", err
//	}
//
//	return logic, nil
//}
//
//// Bytes returns the LogicID as a []byte after being
//// decoded from its hexadecimal string representation.
//// Panics if the LogicID is not a valid hex string.
//func (logic LogicID) Bytes() []byte {
//	return must(decodeHexString(string(logic)))
//}
//
//// String returns the LogicID as a string.
//// Implements the fmt.Stringer interface for LogicID.
//func (logic LogicID) String() string {
//	return "0x" + string(logic)
//}
//
//// Address returns the Address of the LogicID.
//// Returns NilAddress if the LogicID is invalid.
//// The LogicID standard expects the address to ALWAYS be the last 32 bytes
//func (logic LogicID) Address() Address {
//	// Error if length is too short
//	if len(logic) < 64 {
//		return NilAddress
//	}
//
//	// Trim the last 64 characters (32 bytes)
//	addr := string(logic[len(logic)-64:])
//	// Assertively decode into an Address
//	return must(NewAddressFromHex(addr))
//}
//
//// Identifier returns a LogicIdentifier for the LogicID.
////
//// This decodes the LogicID from its simple hex-encoded string
//// format into a representation appropriate for the version of
//// the LogicID to allow access to all the encoded fields within it.
////
//// It can also be used to verify the integrity of the LogicID
//func (logic LogicID) Identifier() (LogicIdentifier, error) {
//	id, err := decodeHexString(string(logic))
//	if err != nil {
//		return nil, errors.Wrap(err, "invalid logic ID")
//	}
//
//	// We verify that there is at least 1 byte, so that
//	// we can safely access the 0th index in the byte slice
//	if len(id) < 1 {
//		return nil, errors.New("invalid logic ID: missing version prefix")
//	}
//
//	// Determine the version of the LogicID and decode
//	switch version := int(id[0] & 0xF0); version {
//	case 0:
//		return decodeLogicIDv0(id)
//	default:
//		return nil, errors.Errorf("invalid logic ID: unsupported version: %v", version)
//	}
//}
//
//// MarshalText implements the encoding.TextMarshaler interface for LogicID
//func (logic LogicID) MarshalText() ([]byte, error) {
//	return []byte(logic.String()), nil
//}
//
//// UnmarshalText implements the encoding.TextUnmarshaler interface for LogicID
//func (logic *LogicID) UnmarshalText(text []byte) error {
//	// Assert that the 0x prefix exists
//	if !has0xPrefixBytes(text) {
//		return ErrMissing0xPrefix
//	}
//
//	// Trim the 0x prefix
//	text = trim0xPrefixBytes(text)
//	// Generate an identifier for the LogicID
//	if _, err := LogicID(text).Identifier(); err != nil {
//		return err
//	}
//
//	*logic = LogicID(text)
//
//	return nil
//}
//
//// MarshalJSON implements the json.Marshaler interface for LogicID
//func (logic LogicID) MarshalJSON() ([]byte, error) {
//	return json.Marshal(logic.String())
//}
//
//// UnmarshalJSON implements the json.Unmarshaler interface for LogicID
//func (logic *LogicID) UnmarshalJSON(data []byte) error {
//	var decoded string
//
//	// Decode the JSON data into a string
//	if err := json.Unmarshal(data, &decoded); err != nil {
//		return err
//	}
//
//	// Assert that the 0x prefix exists
//	if !has0xPrefixString(decoded) {
//		return ErrMissing0xPrefix
//	}
//
//	// Trim the 0x prefix
//	decoded = trim0xPrefixString(decoded)
//	// Generate an identifier for the AssetID
//	if _, err := LogicID(decoded).Identifier(); err != nil {
//		return err
//	}
//
//	*logic = LogicID(decoded)
//
//	return nil
//}

//// LogicIDV0Length is the length of the v0 specification of the LogicID Standard
//const LogicIDV0Length = 35
//
//// LogicIdentifierV0 is an implementation of v0 specification
//// of the LogicID Standard and implements the LogicIdentifier
//type LogicIdentifierV0 [LogicIDV0Length]byte
//
//// NewLogicIDv0 generates a new LogicID with the v0 specification. The LogicID v0 Form is defined as follows:
//// [version(4bits)|persistent(1bit)|ephemeral(1bit)|interactable(1bit)|asset(1bit)][edition(16bits)][address(256bits)]
//func NewLogicIDv0(persistent, ephemeral, interactable, assetlogic bool, edition uint16, addr Address) LogicID {
//	// The 4 MSB bits of the head are set the
//	// version of the Logic ID Form (v0)
//	var head uint8 = 0x00 << 4
//
//	// If persistent stateful flag is on, the 5th MSB is set
//	if persistent {
//		head |= 0x8
//	}
//
//	// If ephemeral stateful flag is on, the 6th MSB is set
//	if ephemeral {
//		head |= 0x4
//	}
//
//	// If interactable flag is on, the 7th MSB is set
//	if interactable {
//		head |= 0x2
//	}
//
//	// If asset logic flag is on, the 8th MSB is set
//	if assetlogic {
//		head |= 0x1
//	}
//
//	// Encode the 16-bit edition into its BigEndian bytes
//	editionBuf := make([]byte, 2)
//	binary.BigEndian.PutUint16(editionBuf, edition)
//
//	// Order the logic ID buffer [head][edition][address]
//	buf := make([]byte, 0, 35)
//	buf = append(buf, head)
//	buf = append(buf, editionBuf...)
//	buf = append(buf, addr[:]...)
//
//	return LogicID(hex.EncodeToString(buf))
//}
//
//// decodeLogicIDv0 can be used to decode some data into a LogicIdentifierV0.
//func decodeLogicIDv0(data []byte) (LogicIdentifierV0, error) {
//	// Check if data is the correct length for v0
//	if len(data) != LogicIDV0Length {
//		return LogicIdentifierV0{}, errors.New("invalid logic ID: insufficient length for v0")
//	}
//
//	// Create an LogicIdentifierV0 and copy the data into it
//	identifier := LogicIdentifierV0{}
//	copy(identifier[:], data)
//
//	return identifier, nil
//}
//
//// LogicID returns the LogicIdentifierV0 as a LogicID
//func (logic LogicIdentifierV0) LogicID() LogicID {
//	return LogicID(hex.EncodeToString(logic[:]))
//}
//
//// Version returns the version of the LogicIdentifierV0.
//func (logic LogicIdentifierV0) Version() int { return 0 }
//
//// HasPersistentState returns whether the persistent state flag is set for the LogicIdentifierV0.
//func (logic LogicIdentifierV0) HasPersistentState() bool {
//	// Determine the 5th LSB of the first byte (v0)
//	bit := (logic[0] >> 3) & 0x1
//	// Return true if bit is set
//	return bit != 0
//}
//
//// HasEphemeralState returns whether the ephemeral state flag is set for the LogicIdentifierV0.
//func (logic LogicIdentifierV0) HasEphemeralState() bool {
//	// Determine the 6th LSB of the first byte (v0)
//	bit := (logic[0] >> 2) & 0x1
//	// Return true if bit is set
//	return bit != 0
//}
//
//// HasInteractableSites returns whether the interactable flag is set for the LogicIdentifierV0.
//func (logic LogicIdentifierV0) HasInteractableSites() bool {
//	// Determine the 7th LSB of the first byte (v0)
//	bit := (logic[0] >> 1) & 0x1
//	// Return true if bit is set
//	return bit != 0
//}
//
//// AssetLogic returns whether the asset logic flag is set for the LogicIdentifierV0.
//func (logic LogicIdentifierV0) AssetLogic() bool {
//	// Determine the 8th LSB of the first byte (v0)
//	bit := logic[0] & 0x1
//	// Return true if bit is set
//	return bit != 0
//}
//
//// Edition returns the edition number of the LogicIdentifierV0.
//func (logic LogicIdentifierV0) Edition() uint64 {
//	// Decode the edition data from the second and third byte of
//	// the LogicID (v0). We decode it as 16-bit number and convert
//	edition := binary.BigEndian.Uint16(logic[1:3])
//
//	return uint64(edition)
//}
//
//// Address returns the Logic Address of the LogicIdentifierV0.
//func (logic LogicIdentifierV0) Address() Address {
//	// Address data is everything after the third byte (v0)
//	// We know it will be 32 bytes, because of the validity check
//	address := logic[3:]
//	// Convert address data into an Address and return
//	return NewAddressFromBytes(address)
//}
//

//func TestLogicIDv0(t *testing.T) {
//	// Setup fuzzer
//	f := fuzz.New().NilChance(0)
//
//	type params struct {
//		Persistent  bool
//		Ephemeral   bool
//		Interactive bool
//		AssetLogic  bool
//
//		Edition uint16
//		Address Address
//	}
//
//	var x params
//
//	for i := 0; i < 1000; i++ {
//		// Fuzz Parameters
//		f.Fuzz(&x)
//
//		// Create new LogicID
//		logicID := NewLogicIDv0(x.Persistent, x.Ephemeral, x.Interactive, x.AssetLogic, x.Edition, x.Address)
//		require.Equal(t, x.Address, logicID.Address())
//
//		// Check type conversions
//		require.Equal(t, "0x"+string(logicID), logicID.String())
//		require.Equal(t, must(decodeHexString(string(logicID))), logicID.Bytes())
//
//		// Check identifier conversion
//		identifier, err := logicID.Identifier()
//		require.NoError(t, err)
//
//		// Check encoded parameter accessors
//		require.Equal(t, 0, identifier.Version())
//		require.Equal(t, logicID, identifier.LogicID())
//
//		require.Equal(t, x.Persistent, identifier.HasPersistentState())
//		require.Equal(t, x.Ephemeral, identifier.HasEphemeralState())
//		require.Equal(t, x.Interactive, identifier.HasInteractableSites())
//		require.Equal(t, x.AssetLogic, identifier.AssetLogic())
//
//		require.Equal(t, uint64(x.Edition), identifier.Edition())
//		require.Equal(t, x.Address, identifier.Address())
//
//		// Check serialization
//		encoded, err := json.Marshal(logicID)
//		require.NoError(t, err)
//		require.Equal(t, fmt.Sprintf(`"0x%v"`, string(logicID)), string(encoded))
//
//		// Check deserialization
//		var decoded LogicID
//		err = json.Unmarshal(encoded, &decoded)
//		require.NoError(t, err)
//		require.Equal(t, logicID, decoded)
//	}
//}
