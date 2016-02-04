package dawg

import (
	"fmt"
	"testing"
)

func TestPlist(t *testing.T) {
	data := PList{
		"bundleid": PString("val"),
		"category": PString("Productivity"),
		"connections": PDict{
			"UUID": PArray{
				PDict{
					"destinationuid": PString("UUID"),
					"modifiers":      PInteger(0),
					"modifersubtext": PString(""),
				},
			}}}

	out, err := data.PListWithHeader()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(out))
}
