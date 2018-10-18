package svgparser

import (
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"path"

	"github.com/eknkc/basex"
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
	element.UUID = element.Hash()

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

// Ancestors returns ancestors' elements
func (e *Element) Ancestors() []*Element {
	ee := []*Element{e.Parent}
	p := e
	for p.Parent != nil {
		p = p.Parent
		ee = append(ee, p)
	}
	return ee
}

// Descendants returns descendants' elements
func (e *Element) Descendants() []*Element {
	ee := []*Element{}
	for _, child := range e.Children {
		ee = append(ee, child)
		ee = append(ee, child.Descendants()...)
	}
	return ee
}

// Generations returns all generations
func (e *Element) Generations() []*Element {
	ee := []*Element{}
	ee = append(ee, e.Ancestors()...)
	ee = append(ee, e)
	ee = append(ee, e.Descendants()...)
	return ee
}

// Hash returns element's unique hash string
func (e *Element) Hash() string {
	hasher := md5.New()
	for _, g := range e.Generations() {
		if g != nil {
			hasher.Write([]byte(fmt.Sprintf("%v:%v", g.Name, e.Attributes)))
		}
	}
	hasher.Write([]byte(fmt.Sprintf("%v:%v", e.Name, e.Attributes)))
	b64, _ := basex.NewEncoding("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	return b64.Encode(hasher.Sum(nil))
}
