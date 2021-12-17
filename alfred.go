package dawg

import "encoding/xml"

// <?xml version="1.0"?>
// <items>
//   <item uid="#{q dashboard_name}" valid="YES" autocomplete="#{q dashboard_name}">
//     <title>#{dashboard_name}</title>
//     <arg>https://app.datadoghq.com/screen/#{dashboard_id}</arg>
//   </item>
// </items>
type AlfredOutput []AlfredOutputItem

func (a AlfredOutput) MakeXML() ([]byte, error) {
	serialized, err := xml.MarshalIndent(a, "", "  ")
	if err != nil {
		return nil, err
	}
	data := []byte{}
	data = append(data, []byte(xml.Header)...)
	data = append(data, serialized...)
	return data, nil
}

func (a AlfredOutput) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "items"}
	e.EncodeToken(start)
	for _, item := range a {
		e.EncodeElement(item, xml.StartElement{Name: xml.Name{Local: "item"}})
	}
	e.EncodeToken(start.End())
	return nil
}

type AlfredOutputItem struct {
	UID          string `xml:"uid,attr"`
	Autocomplete string `xml:"autocomplete,attr"`
	Title        string `xml:"title"`
	Icon         string `xml:"icon"`
	Arg          string `xml:"arg"`
}
