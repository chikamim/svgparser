package svgparser

import (
	"log"
	"regexp"
)

// FindID finds the first child with the specified ID.
func (e *Element) FindID(id string) *Element {
	for _, child := range e.Children {
		if childID, ok := child.Attributes["id"]; ok && childID == id {
			return child
		}
		if element := child.FindID(id); element != nil {
			return element
		}
	}
	return nil
}

// FindUUIID finds the first child with the specified UUID.
func (e *Element) FindUUID(uuid string) *Element {
	for _, child := range e.Children {
		if child.UUID == uuid {
			return child
		}
		if element := child.FindUUID(uuid); element != nil {
			return element
		}
	}
	return nil
}

// FindAll finds all children with the given name.
func (e *Element) FindAll(name string) []*Element {
	var elements []*Element
	for _, child := range e.Children {
		if child.Name == name {
			elements = append(elements, child)
		}
		elements = append(elements, child.FindAll(name)...)
	}
	return elements
}

// FilterIDs filter all children with the given ids.
func (e *Element) SelectByUUIDs(uuids []string) *Element {
	c := &Element{}

	u := make(map[string]bool)
	uniq := []string{}
	for _, uuid := range uuids {
		if !u[uuid] {
			u[uuid] = true
			uniq = append(uniq, uuid)
		}
	}

	for _, uuid := range uniq {
		if e.UUID == uuid {
			*c = *e
			c.Children = []*Element{}
			break
		}
	}

	c.Children = e.selectChildrenByUUIDs(uniq)

	return c
}

func (e *Element) selectChildrenByUUIDs(uuids []string) []*Element {
	var elements []*Element
	for _, child := range e.Children {
		for _, uuid := range uuids {
			if child.UUID == uuid {
				c := &Element{}
				*c = *child
				c.Children = []*Element{}
				c.Children = child.selectChildrenByUUIDs(uuids)
				elements = append(elements, c)
			}
		}
	}
	return elements
}

// FindAllLinkedIDs finds related linked ID recursively
func FindAllLinkedIDs(r *Element, id string) []string {
	ids := []string{}
	f := r.FindID(id)
	if f == nil {
		return ids
	}

	c := f
	for c.Children != nil {
		for _, h := range c.Children {
			for _, id := range h.FindLinkedIDs() {
				ids = append(ids, id)
			}
			c = h
		}
	}

	ids = append(ids, f.FindLinkedIDs()...)
	for _, fid := range ids {
		if fid == id {
			continue
		}
		ids = append(ids, FindAllLinkedIDs(r, fid)...)
	}
	log.Println(ids)
	return ids
}

// FindLinkedIDs finds related linked ID
func (e *Element) FindLinkedIDs() []string {
	ids := []string{}
	if id, found := e.Attributes["id"]; found {
		ids = append(ids, id)
	}
	r := regexp.MustCompile("#(\\w+)")
	for _, v := range e.Attributes {
		if m := r.FindStringSubmatch(v); len(m) > 0 {
			ids = append(ids, m[1])
		}
	}
	return ids
}
