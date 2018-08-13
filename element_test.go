package svgparser_test

import (
	"encoding/xml"
	"testing"

	"github.com/chikamim/svgparser"
)

func TestNewElement(t *testing.T) {
	attr := []xml.Attr{}
	attr = append(attr, xml.Attr{xml.Name{"", "id"}, "1234"})
	start := xml.StartElement{xml.Name{"", "mask"}, attr}

	element := svgparser.NewElement(start)
	if element.Name != "mask" {
		t.Errorf("Element Name expected %v, actual %v\n", "mask", element.Name)
	}

	if val, ok := element.Attributes["id"]; ok {
		if val != "1234" {
			t.Errorf("Element Attibute key expected %v, actual %v\n", "1234", val)
		}
	} else {
		t.Error("Element Attibute key not found")
	}
}

func TestNewElementWithSpace(t *testing.T) {
	attr := []xml.Attr{}
	attr = append(attr, xml.Attr{xml.Name{"http://www.w3.org/1999/xlink", "href"}, "#1234"})
	start := xml.StartElement{xml.Name{"", "mask"}, attr}

	element := svgparser.NewElement(start)
	if element.Name != "mask" {
		t.Errorf("Element Name expected %v, actual %v\n", "mask", element.Name)
	}

	if val, ok := element.Attributes["xlink:href"]; ok {
		if val != "#1234" {
			t.Errorf("Element Attibute key expected %v, actual %v\n", "#1234", val)
		}
	} else {
		t.Error("Element Attibute key not found")
	}
}
