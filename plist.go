package dawg

import "encoding/xml"

type PList PDict

const (
	PListHeader  = `<?xml version="1.0" encoding="UTF-8"?>`
	PListDoctype = `<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">`
)

func (p PList) PListWithHeader() ([]byte, error) {
	serialized, err := xml.MarshalIndent(p, "", "  ")
	if err != nil {
		return nil, err
	}
	data := []byte{}
	data = append(data, []byte(PListHeader)...)
	data = append(data, []byte("\n")...)
	data = append(data, []byte(PListDoctype)...)
	data = append(data, []byte("\n")...)
	data = append(data, serialized...)
	return data, nil
}

func (p PList) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "plist"}
	start.Attr = []xml.Attr{xml.Attr{Name: xml.Name{Local: "version"}, Value: "1.0"}}
	e.EncodeToken(start)
	PDict(p).MarshalXML(e, xml.StartElement{})
	e.EncodeToken(start.End())
	return nil
}

type PArray []xml.Marshaler

func (a PArray) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "array"}
	e.EncodeToken(start)
	for _, el := range a {
		el.(xml.Marshaler).MarshalXML(e, start)
	}
	e.EncodeToken(start.End())
	return nil
}

type PInteger int

func (i PInteger) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	e.EncodeElement(int(i), xml.StartElement{Name: xml.Name{Local: "integer"}})
	return nil
}

type PReal int

func (i PReal) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	e.EncodeElement(float64(i), xml.StartElement{Name: xml.Name{Local: "real"}})
	return nil
}

type PString string

func (s PString) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	e.EncodeElement(string(s), xml.StartElement{Name: xml.Name{Local: "string"}})
	return nil
}

type PBool bool

func (b PBool) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var val string
	if b {
		val = "true"
	} else {
		val = "false"
	}
	e.EncodeElement("", xml.StartElement{Name: xml.Name{Local: val}})
	return nil
}

type PDict map[string]xml.Marshaler

func (d PDict) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "dict"}
	e.EncodeToken(start)
	for k, v := range d {
		e.EncodeElement(k, xml.StartElement{Name: xml.Name{Local: "key"}})
		v.MarshalXML(e, start)
	}
	e.EncodeToken(start.End())
	return nil
}
