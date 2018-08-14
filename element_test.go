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

func testFamilyElement() *svgparser.Element {
	svg := `
		<svg width="1000" height="600" id="svg">
			<g id="grandfather">
				<g id="father">
					<g id="me">
						<g id="son">
							<rect id="toy" width="5" height="1"/>
						</g>
					</g>
				</g>
			</g>
		</svg>
	`
	element, _ := parse(svg, false)
	return element
}

func TestAncestors(t *testing.T) {
	element := testFamilyElement().FindID("me")
	actual := element.Ancestors()[0].Attributes["id"]
	expected := "father"
	if actual != expected {
		t.Errorf("Ancestors expected %v, got %v", expected, actual)
	}
}

func TestDescendants(t *testing.T) {
	element := testFamilyElement().FindID("me")
	actual := element.Descendants()[0].Attributes["id"]
	expected := "son"
	if actual != expected {
		t.Errorf("Ancestors expected %v, got %v", expected, actual)
	}
}

func TestGenerations(t *testing.T) {
	element := testFamilyElement().FindID("me")
	generations := []string{"svg", "grandfather", "father", "me", "son", "toy"}
	for _, e := range element.Generations() {
		found := false
		for _, g := range generations {
			if e.Attributes["id"] == g {
				found = true
			}
		}
		if !found {
			t.Error("Expected Generation element not found")
		}
	}
}
