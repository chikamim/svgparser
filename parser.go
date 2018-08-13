package svgparser

import (
	"bytes"
	"encoding/xml"

	"io"
	"io/ioutil"
	"strings"

	"golang.org/x/net/html/charset"
)

// ValidationError contains errors which have occured when parsing svg input.
type ValidationError struct {
	msg string
}

func (err ValidationError) Error() string {
	return err.msg
}

// DecodeFirst creates the first element from the decoder.
func DecodeFirst(decoder *xml.Decoder) (*Element, error) {
	for {
		token, err := decoder.Token()
		if token == nil && err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		switch element := token.(type) {
		case xml.StartElement:
			return NewElement(element), nil
		}
	}
	return &Element{}, nil
}

// Decode decodes the child elements of element.
func (e *Element) Decode(decoder *xml.Decoder) error {
	for {
		token, err := decoder.Token()
		if token == nil && err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		switch element := token.(type) {
		case xml.StartElement:
			nextElement := NewElement(element)
			err := nextElement.Decode(decoder)
			if err != nil {
				return err
			}

			e.Children = append(e.Children, nextElement)

		case xml.CharData:
			data := strings.TrimSpace(string(element))
			if data != "" {
				e.Content = string(element)
			}

		case xml.EndElement:
			if element.Name.Local == e.Name {
				return nil
			}
		}
	}
	return nil
}

// Parse creates an Element instance from an SVG input.
func Parse(source io.Reader, validate bool) (*Element, error) {
	raw, err := ioutil.ReadAll(source)
	if err != nil {
		return nil, err
	}
	decoder := xml.NewDecoder(bytes.NewReader(raw))
	decoder.CharsetReader = charset.NewReaderLabel
	element, err := DecodeFirst(decoder)
	if err != nil {
		return nil, err
	}

	if err := element.Decode(decoder); err != nil && err != io.EOF {
		return nil, err
	}
	return element, nil
}
