package identifiers

// AssetID is a unique hex-encoded string identifier for assets in the MOI Protocol.
// It encodes within itself key properties about the asset itself such as the address
// at which the asset is deployed, the nature of it logic control and others.
//
// The Asset ID Standard is an extensible standard is compliant with the protocol specification
// for it at https://sarvalabs.notion.site/Asset-ID-Standard-e4fcd9151e7d4e7eb2447f1d8edf4672?pvs=4
// type AssetID string

// AssetIdentifier is a representation of AssetID allows access to the
// values encoded within it such as the address or the asset standard.
//
// Each version of the AssetID standard is available
// as its own type that implements this interface
//type AssetIdentifier interface {
//	// Version returns the version of the AssetID standard for the AssetIdentifier
//	Version() int
//	// Address returns the 32-byte address associated with the AssetID
//	Address() Address
//	// AssetID returns the AssetIdentifier in its encoded representation as a AssetID
//	AssetID() AssetID
//
//	// Standard returns the asset standard number from the AssetID
//	Standard() uint64
//	// Dimension returns the asset dimension value from the AssetID
//	Dimension() uint8
//
//	// IsLogical returns whether the asset has some logic associated with it
//	IsLogical() bool
//	// IsStateful returns whether the asset has some stateful information such as its supply
//	IsStateful() bool
//}
//
//// NewAssetID generates an AssetID from some arbitrary string,
//// validating it in the process. It is version agnostic.
//func NewAssetID(id string) (AssetID, error) {
//	asset := AssetID(id)
//
//	// Attempt to generate an identifier from the AssetID
//	// This will fail if the AssetID is invalid in any way
//	if _, err := asset.Identifier(); err != nil {
//		return "", err
//	}
//
//	return asset, nil
//}
//
//// Bytes returns the AssetID as a []byte after being
//// decoded from its hexadecimal string representation
//// Panics if the AssetID is not a valid hex string
//func (asset AssetID) Bytes() []byte {
//	return must(decodeHexString(string(asset)))
//}
//
//// String returns the AssetID as a string.
//// Implements the fmt.Stringer interface for AssetID.
//func (asset AssetID) String() string {
//	return "0x" + string(asset)
//}
//
//// Address returns the Address of the AssetID.
//// Returns NilAddress if the AssetID is invalid.
//// The AssetID standard expects the address to ALWAYS be the last 32 bytes
//func (asset AssetID) Address() Address {
//	// Error if length is too short
//	if len(asset) < 64 {
//		return NilAddress
//	}
//
//	// Trim the last 64 characters (32 bytes)
//	addr := string(asset[len(asset)-64:])
//	// Assertively decode into an Address
//	return must(NewAddressFromHex(addr))
//}
//
//// Identifier returns a AssetIdentifier for the AssetID.
////
//// This decodes the AssetID from its simple hex-encoded string
//// format into a representation appropriate for the version of
//// the AssetID to allow access to all the encoded fields within it.
////
//// It can also be used to verify the integrity of the AssetID
//func (asset AssetID) Identifier() (AssetIdentifier, error) {
//	id, err := decodeHexString(string(asset))
//	if err != nil {
//		return nil, errors.Wrap(err, "invalid asset ID")
//	}
//
//	// We verify that there is at least 1 byte, so that
//	// we can safely access the 0th index in the byte slice
//	if len(id) < 1 {
//		return nil, errors.New("invalid asset ID: missing version prefix")
//	}
//
//	// Determine the version of the AssetID and check if there are enough bytes
//	switch version := int(id[0] & 0xF0); version {
//	case 0:
//		return decodeAssetIDv0(id)
//	default:
//		return nil, errors.Errorf("invalid asset ID: unsupported version: %v", version)
//	}
//}
//
//// MarshalText implements the encoding.TextMarshaler interface for AssetID
//func (asset AssetID) MarshalText() ([]byte, error) {
//	return []byte(asset.String()), nil
//}
//
//// UnmarshalText implements the encoding.TextUnmarshaler interface for AssetID
//func (asset *AssetID) UnmarshalText(text []byte) error {
//	// Assert that the 0x prefix exists
//	if !has0xPrefixBytes(text) {
//		return ErrMissing0xPrefix
//	}
//
//	// Trim the 0x prefix
//	text = trim0xPrefixBytes(text)
//	// Generate an identifier for the AssetID
//	if _, err := AssetID(text).Identifier(); err != nil {
//		return err
//	}
//
//	*asset = AssetID(text)
//
//	return nil
//}
//
//// MarshalJSON implements the json.Marshaler interface for AssetID
//func (asset AssetID) MarshalJSON() ([]byte, error) {
//	return json.Marshal(asset.String())
//}
//
//// UnmarshalJSON implements the json.Unmarshaler interface for AssetID
//func (asset *AssetID) UnmarshalJSON(data []byte) error {
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
//	if _, err := AssetID(decoded).Identifier(); err != nil {
//		return err
//	}
//
//	*asset = AssetID(decoded)
//
//	return nil
//}

//// AssetIDV0Length is the length of the v0 specification of the AssetID Standard
//const AssetIDV0Length = 36
//
//// AssetIdentifierV0 is an implementation of v0 specification
//// of the AssetID Standard and implements the AssetIdentifier
//type AssetIdentifierV0 [AssetIDV0Length]byte
//
//// NewAssetIDv0 generates a new AssetID with the v0 specification. The AssetID v0 Form is defined as follows:
//// [version(4bits)|logical(1bit)|stateful(1bit)|reserved(2bits)][dimension(8bits)][standard(16bits)][address(256bits)]
//func NewAssetIDv0(logical, stateful bool, dimension uint8, standard uint16, addr Address) AssetID {
//	// The 4 MSB bits of the head are set the
//	// version of the Asset ID Form (v0)
//	var head uint8 = 0x00 << 4
//
//	// If logical flag is on, the 5th MSB is set
//	if logical {
//		head |= 0x8
//	}
//
//	// If stateful flag is on, the 6th MSB is set
//	if stateful {
//		head |= 0x4
//	}
//
//	// Encode the 16-bit standard into its BigEndian bytes
//	standardBuf := make([]byte, 2)
//	binary.BigEndian.PutUint16(standardBuf, standard)
//
//	// Order the asset ID buffer [head][dimension][standard][address]
//	buf := make([]byte, 0, 36)
//	buf = append(buf, head)
//	buf = append(buf, dimension)
//	buf = append(buf, standardBuf...)
//	buf = append(buf, addr[:]...)
//
//	return AssetID(hex.EncodeToString(buf))
//}
//
//// decodeAssetIDv0 can be used to decode some data into a AssetIdentifierV0.
//func decodeAssetIDv0(data []byte) (AssetIdentifierV0, error) {
//	// Check if data is the correct length for v0
//	if len(data) != AssetIDV0Length {
//		return AssetIdentifierV0{}, errors.New("invalid asset ID: insufficient length for v0")
//	}
//
//	// Create an AssetIdentifierV0 and copy the data into it
//	identifier := AssetIdentifierV0{}
//	copy(identifier[:], data)
//
//	return identifier, nil
//}
//
//// AssetID returns the AssetIdentifierV0 as an AssetID
//func (asset AssetIdentifierV0) AssetID() AssetID {
//	return AssetID(hex.EncodeToString(asset[:]))
//}
//
//// Version returns the version of the AssetIdentifierV0.
//func (asset AssetIdentifierV0) Version() int { return 0 }
//
//// IsLogical returns whether the logical flag is set for the AssetIdentifierV0.
//func (asset AssetIdentifierV0) IsLogical() bool {
//	// Determine the 5th LSB of the first byte (v0)
//	bit := (asset[0] >> 3) & 0x1
//	// Return true if bit is set
//	return bit != 0
//}
//
//// IsStateful returns whether the stateful flag is set for the AssetIdentifierV0.
//func (asset AssetIdentifierV0) IsStateful() bool {
//	// Determine the 6th LSB of the first byte (v0)
//	bit := (asset[0] >> 2) & 0x1
//	// Return true if bit is set
//	return bit != 0
//}
//
//// Dimension returns the dimension of the AssetIdentifierV0.
//func (asset AssetIdentifierV0) Dimension() uint8 {
//	// Dimension data is in the second byte of the AssetID (v0)
//	return asset[1]
//}
//
//// Standard returns the standard of the AssetIdentifierV0.
//func (asset AssetIdentifierV0) Standard() uint64 {
//	// Decode the edition data from the third and fourth byte of
//	// the AssetID (v0). We decode it as 16-bit number and convert
//	standard := binary.BigEndian.Uint16(asset[2:4])
//
//	return uint64(standard)
//}
//
//// Address returns the Asset Address of the AssetIdentifier.
//func (asset AssetIdentifierV0) Address() Address {
//	// Address data is everything after the fourth byte (v0)
//	// We know it will be 32 bytes, because of the validity check
//	address := asset[4:]
//	// Convert address data into an Address and return
//	return NewAddressFromBytes(address)
//}

//func TestAssetIDv0(t *testing.T) {
//	// Setup fuzzer
//	f := fuzz.New().NilChance(0)
//
//	type params struct {
//		Logical   bool
//		Stateful  bool
//		Dimension uint8
//		Standard  uint16
//		Address   Address
//	}
//
//	var x params
//
//	for i := 0; i < 1000; i++ {
//		// Fuzz Parameters
//		f.Fuzz(&x)
//
//		// Create new AssetID
//		assetID := NewAssetIDv0(x.Logical, x.Stateful, x.Dimension, x.Standard, x.Address)
//		require.Equal(t, x.Address, assetID.Address())
//
//		// Check type conversions
//		require.Equal(t, "0x"+string(assetID), assetID.String())
//		require.Equal(t, must(decodeHexString(string(assetID))), assetID.Bytes())
//
//		// Check identifier conversion
//		identifier, err := assetID.Identifier()
//		require.NoError(t, err)
//
//		// Check encoded parameter accessors
//		require.Equal(t, 0, identifier.Version())
//		require.Equal(t, assetID, identifier.AssetID())
//
//		require.Equal(t, x.Logical, identifier.IsLogical())
//		require.Equal(t, x.Stateful, identifier.IsStateful())
//
//		require.Equal(t, x.Dimension, identifier.Dimension())
//		require.Equal(t, uint64(x.Standard), identifier.Standard())
//		require.Equal(t, x.Address, identifier.Address())
//
//		// Check serialization
//		encoded, err := json.Marshal(assetID)
//		require.NoError(t, err)
//		require.Equal(t, fmt.Sprintf(`"0x%v"`, string(assetID)), string(encoded))
//
//		// Check deserialization
//		var decoded AssetID
//		err = json.Unmarshal(encoded, &decoded)
//		require.NoError(t, err)
//		require.Equal(t, assetID, decoded)
//	}
//}
