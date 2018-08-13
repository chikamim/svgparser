package svgparser

import (
	"encoding/xml"
	"path"

	"github.com/renstrom/shortuuid"
)

// Element is a representation of an SVG element.
type Element struct {
	UUID       string
	Name       string
	Attributes map[string]string
	Parent     *Element
	Children   []*Element
	Content    string
}

// NewElement creates element from decoder token.
func NewElement(token xml.StartElement) *Element {
	element := &Element{}
	element.UUID = shortuuid.New()
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
		attr = append(attr, xml.Attr{name, v})
	}
	return xml.StartElement{xml.Name{"", e.Name}, attr}
}

// Ancestors returns ancestors elements
func (e *Element) Ancestors() []*Element {
	ee := []*Element{e.Parent}
	p := e
	for p.Parent != nil {
		p = p.Parent
		ee = append(ee, p)
	}
	return ee
}

// Descendant returns descendant elements
func (e *Element) Descendant() []*Element {
	ee := []*Element{}
	for _, child := range e.Children {
		ee = append(ee, child)
		ee = append(ee, child.Descendant()...)
	}
	return ee
}

// Generations returns all generations
func (e *Element) Generations() []*Element {
	ee := []*Element{}
	ee = append(ee, e.Ancestors()...)
	ee = append(ee, e)
	ee = append(ee, e.Descendant()...)
	return ee
}
