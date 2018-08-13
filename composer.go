package svgparser

import (
	"encoding/xml"
	"fmt"
	"io"
)

// Compose convert SVG from element
func (e *Element) Compose(w io.Writer) error {
	enc := xml.NewEncoder(w)

	err := EncodeXML(e, enc)
	if err != nil {
		return err
	}
	return nil
}

// EncodeXML encode XML elements recursively
func EncodeXML(r *Element, e *xml.Encoder) (err error) {
	start := r.XMLElement()
	err = e.EncodeToken(start)
	if err != nil {
		return fmt.Errorf("failed to encode start element: %v", err)
	}
	for _, c := range r.Children {
		EncodeXML(c, e)
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
