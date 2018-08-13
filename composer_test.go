package svgparser_test

import (
	"bufio"
	"bytes"
	"testing"
)

func TestCompose(t *testing.T) {
	svg := `
<svg height="400" width="450">
  <g stroke="black" stroke-width="3" fill="black">
    <path id="AB" d="M 100 350 L 150 -300" stroke="red" />
    <path id="BC" d="M 250 50 L 150 300" stroke="red" />
    <path d="M 175 200 L 150 0" stroke="green" />
  </g>
</svg>
`
	element, _ := parse(svg, true)

	buf := bytes.Buffer{}
	w := bufio.NewWriter(&buf)
	err := element.Compose(w)
	if err != nil {
		t.Errorf("Compose failed: %v\n", err)
	}

	actual, _ := parse(buf.String(), true)
	if !element.Compare(actual) {
		t.Errorf("Compose output is not the same actual %+v, expected %+v", actual, element)
	}
}
