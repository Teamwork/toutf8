// Package toutf8 converts text to UTF-8.
//
// This is a friendlier API for the golang.org/x/text collection of packages,
// making it easier to call for arbitrary input in an iconv-like manner.
package toutf8

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/ianaindex"
	"golang.org/x/text/encoding/unicode"
)

// ErrUnknownCharset is used when the passed character set is unknown to
// this library.
type ErrUnknownCharset string

func (s ErrUnknownCharset) Error() string { return string(s) }

// More aliases.
// TODO: this list can be expanded; Look at iconv?
var encMaps = map[string]encoding.Encoding{
	"ascii":    charmap.ISO8859_1,
	"us-ascii": charmap.ISO8859_1,
	"utf8":     unicode.UTF8,
}

var names = []*ianaindex.Index{ianaindex.MIME, ianaindex.IANA, ianaindex.MIB}

// FindEncoding attempts to find an encoding by name.
func FindEncoding(charset string) (encoding.Encoding, bool) {
	for i := range names {
		e, err := names[i].Encoding(charset)
		if err == nil && e != nil {
			return e, true
		}
	}

	// TODO: perhaps use some cleverness to normalize names? e.g. "shift jis" to
	// "shift-jis", "utf8" to "utf-8", etc.

	e, ok := encMaps[strings.ToLower(charset)]
	return e, ok
}

// Reader converts the bytes from the input reader to UTF-8.
func Reader(charset string, input io.Reader) (io.Reader, error) {
	enc, ok := FindEncoding(charset)
	if !ok {
		return nil, ErrUnknownCharset(fmt.Sprintf("unknown character set: %v", charset))
	}

	r := enc.NewDecoder().Reader(input)
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}

// Byte converts the bytes from the input to UTF-8.
func Byte(charset string, input []byte) ([]byte, error) {
	out, err := Reader(charset, bytes.NewReader(input))
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(out)
	if err != nil {
		return nil, fmt.Errorf("could not read from output: %v", err)
	}

	return b, nil
}
