package svgparser

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// Compose convert SVG from element
func (e *Element) Compose(w io.Writer) error {
	enc := xml.NewEncoder(w)

	err := EncodeXML(e, enc, []string{})
	if err != nil {
		return err
	}
	return nil
}

func (e *Element) ComposeExcludes(w io.Writer, uuids []string) error {
	enc := xml.NewEncoder(w)

	err := EncodeXML(e, enc, uuids)
	if err != nil {
		return err
	}
	return nil
}

// EncodeXML encode XML elements recursively
func EncodeXML(r *Element, e *xml.Encoder, excludes []string) (err error) {
	for _, uuid := range excludes {
		if r.UUID == uuid {
			return nil
		}
	}
	if val, ok := r.Attributes["xlink:href"]; ok {
		r.Attributes["xlink:href"] = strings.Replace(val, "\n", "", -1)
	}

	start := r.XMLElement()
	err = e.EncodeToken(start)
	if err != nil {
		return fmt.Errorf("failed to encode start element: %v", err)
	}
	for _, c := range r.Children {
		EncodeXML(c, e, excludes)
	}

	err = e.EncodeToken(xml.EndElement{start.Name})
	if err != nil {
		return fmt.Errorf("failed to encode end element: %v", err)
	}

	err = e.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush encoder: %v", err)
	}

	return nil
}
