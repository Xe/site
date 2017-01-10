package asar // import "layeh.com/asar"

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	errMalformed = errors.New("asar: malformed archive")
)

// Decode decodes the ASAR archive in ra.
//
// Returns the root element and nil on success. nil and an error is returned on
// failure.
func Decode(ra io.ReaderAt) (*Entry, error) {
	headerSize := uint32(0)
	headerStringSize := uint32(0)

	// [pickle object header (4 bytes) == 4]
	// [pickle uint32 = $header_object_size]
	{
		var buff [8]byte
		if n, _ := ra.ReadAt(buff[:], 0); n != 8 {
			return nil, errMalformed
		}

		dataSize := binary.LittleEndian.Uint32(buff[:4])
		if dataSize != 4 {
			return nil, errMalformed
		}
		headerSize = binary.LittleEndian.Uint32(buff[4:8])
	}

	// [pickle object header (4 bytes)]
	// [pickle data header (4 bytes) == $string_size]
	// [pickle string ($string_size bytes)]
	{
		var buff [8]byte
		if n, _ := ra.ReadAt(buff[:], 8); n != 8 {
			return nil, errMalformed
		}

		headerObjectSize := binary.LittleEndian.Uint32(buff[:4])
		if headerObjectSize != headerSize-4 {
			return nil, errMalformed
		}

		headerStringSize = binary.LittleEndian.Uint32(buff[4:8])
	}

	// read header string
	headerSection := io.NewSectionReader(ra, 8+8, int64(headerStringSize))
	baseOffset := 8 + 8 + int64(headerStringSize)
	baseOffset += baseOffset % 4 // pickle objects are uint32 aligned

	root, err := decodeHeader(ra, headerSection, baseOffset)
	if err != nil {
		return nil, err
	}

	return root, nil
}
