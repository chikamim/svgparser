package svgparser_test

import (
	"testing"

	"github.com/chikamim/svgparser"
)

func testElement() *svgparser.Element {
	svg := `
		<svg width="1000" height="600">
			<g id="first">
				<rect width="5" height="3" id="inFirst"/>
				<rect width="5" height="2" id="inFirst"/>
			</g>
			<g id="second">
				<path d="M50 50 Q50 100 100 100"/>
				<rect width="5" height="1"/>
			</g>
		</svg>
	`
	element, _ := parse(svg, false)
	return element
}

func equals(t *testing.T, name string, expected, actual *svgparser.Element) {
	if !(expected == actual || expected.Compare(actual)) {
		t.Errorf("%s: expected %v, actual %v\n", name, expected, actual)
	}
}

func equalSlices(t *testing.T, name string, expected, actual []*svgparser.Element) {
	if len(expected) != len(actual) {
		t.Errorf("%s: expected %v, actual %v\n", name, expected, actual)
		return
	}

	for i, r := range actual {
		equals(t, name, expected[i], r)
	}
}

func TestFindAll(t *testing.T) {
	svgElement := testElement()

	equalSlices(t, "Find", []*svgparser.Element{
		element("rect", map[string]string{"width": "5", "height": "3", "id": "inFirst"}),
		element("rect", map[string]string{"width": "5", "height": "2", "id": "inFirst"}),
		element("rect", map[string]string{"width": "5", "height": "1"}),
	}, svgElement.FindAll("rect"))

	equalSlices(t, "Find", []*svgparser.Element{}, svgElement.FindAll("circle"))
}

func TestFindID(t *testing.T) {
	svgElement := testElement()

	equals(t, "Find", &svgparser.Element{
		Name:       "g",
		Attributes: map[string]string{"id": "second"},
		Children: []*svgparser.Element{
			element("path", map[string]string{"d": "M50 50 Q50 100 100 100"}),
			element("rect", map[string]string{"width": "5", "height": "1"}),
		},
	}, svgElement.FindID("second"))

	equals(t, "Find",
		element("rect", map[string]string{"width": "5", "height": "3", "id": "inFirst"}),
		svgElement.FindID("inFirst"),
	)

	equals(t, "Find", nil, svgElement.FindID("missing"))
}

func TestFindUUID(t *testing.T) {
	element := &svgparser.Element{
		Name:       "g",
		UUID:       "1",
		Attributes: map[string]string{"id": "first"},
		Children: []*svgparser.Element{
			&svgparser.Element{
				Name:       "rect",
				UUID:       "2",
				Attributes: map[string]string{"id": "second"},
				Children:   []*svgparser.Element{},
			},
		},
	}
	if element.FindUUID("2").Attributes["id"] != "second" {
		t.Error("FindUUID failed")
	}
}

func TestSelectByUUIDs(t *testing.T) {
	element := &svgparser.Element{
		Name:       "g",
		UUID:       "1",
		Attributes: map[string]string{"id": "first"},
		Children: []*svgparser.Element{
			&svgparser.Element{
				Name:       "rect",
				UUID:       "2",
				Attributes: map[string]string{"id": "second"},
				Children: []*svgparser.Element{&svgparser.Element{
					Name:       "rect",
					UUID:       "4",
					Attributes: map[string]string{"id": "forth"},
					Children:   []*svgparser.Element{},
				}},
			},
			&svgparser.Element{
				Name:       "rect",
				UUID:       "3",
				Attributes: map[string]string{"id": "third"},
				Children:   []*svgparser.Element{},
			},
		},
	}

	equals(t, "SelectByUUIDs", &svgparser.Element{
		UUID:       "1",
		Name:       "g",
		Attributes: map[string]string{"id": "first"},
		Children: []*svgparser.Element{
			&svgparser.Element{
				UUID:       "2",
				Name:       "rect",
				Attributes: map[string]string{"id": "second"},
				Children:   []*svgparser.Element{},
			},
		},
	}, element.SelectByUUIDs([]string{"1", "2"}))
}

func testLinkElement() *svgparser.Element {
	svg := `
		<svg width="1000" height="600" id="svg">
			<g id="grandfather">
				<g id="father">
					<g id="me">
						<use id="room" xlink:href="#bycicle"/>
					</g>
				</g>
				<g id="uncle">
					<rect id="garage">
						<path id="car"/>
						<path id="bycicle"/>
						</rect>
					<g id="cousin">
						<use xlink:href="#garage"/>
					</g>
				</g>
			</g>
		</svg>
	`
	element, _ := parse(svg, false)
	return element
}

func TestFindLinkedIDs(t *testing.T) {
	element := testLinkElement().FindID("room")
	expected := []string{"bycicle", "room"}
	for _, id := range element.FindLinkedIDs() {
		found := false
		for _, e := range expected {
			if id == e {
				found = true
			}
		}
		if !found {
			t.Errorf("FindLinkedIDs id %v not found", id)
		}
	}
}

func TestFindAllLinkedIDs(t *testing.T) {
	expected := []string{"garage", "cousin", "car", "bycicle"}
	for _, id := range svgparser.FindAllLinkedIDs(testLinkElement(), "cousin") {
		found := false
		for _, e := range expected {
			if id == e {
				found = true
			}
		}
		if !found {
			t.Errorf("FindAllLinkedIDs id %v not found", id)
		}
	}
}
