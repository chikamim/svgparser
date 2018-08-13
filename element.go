package svgparser

import (
	"encoding/xml"
	"path"
	"strings"
)

// Element is a representation of an SVG element.
type Element struct {
	Name       string
	Attributes map[string]string
	Children   []*Element
	Content    string
}

// NewElement creates element from decoder token.
func NewElement(token xml.StartElement) *Element {
	element := &Element{}
	attributes := make(map[string]string)
	for _, attr := range token.Attr {
		key := attr.Name.Local
		s := path.Base(attr.Name.Space)
		if s != "." {
			key = s + ":" + key
		}
		attributes[key] = attr.Value
	}
	element.Name = token.Name.Local
	element.Attributes = attributes
	return element
}

// Compare compares two elements.
func (e *Element) Compare(o *Element) bool {
	if e.Name != o.Name || e.Content != o.Content ||
		len(e.Attributes) != len(o.Attributes) ||
		len(e.Children) != len(o.Children) {
		return false
	}

	for k, v := range e.Attributes {
		if v != o.Attributes[k] {
			return false
		}
	}

	for i, child := range e.Children {
		if !child.Compare(o.Children[i]) {
			return false
		}
	}
	return true
}

// XMLElement convert to XMLElement
func (e *Element) XMLElement() xml.StartElement {
	attr := []xml.Attr{}
	for k, v := range e.Attributes {
		name := xml.Name{"", k}
		c := strings.Split(k, ":")
		if len(c) > 1 {
			name = xml.Name{c[0], c[1]}
		}
		attr = append(attr, xml.Attr{name, v})
	}
	return xml.StartElement{xml.Name{"", e.Name}, attr}
}
