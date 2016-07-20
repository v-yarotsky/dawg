package dawg

import (
	"bytes"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func Test_generateLogoThumbnail(t *testing.T) {
	f, err := os.Open("./fixtures/test.png")
	must(t, err)
	defer f.Close()

	m, _, err := image.Decode(f)
	must(t, err)

	thumb, err := generateLogoThumbnail(m)
	must(t, err)

	expected, err := os.Open("./fixtures/test-resized.png")
	must(t, err)
	defer expected.Close()

	expectedBytes, err := ioutil.ReadAll(expected)
	must(t, err)

	var out bytes.Buffer
	err = png.Encode(&out, thumb)
	must(t, err)

	if !reflect.DeepEqual(out.Bytes(), expectedBytes) {
		t.Errorf("resized image contents did not match the expected result")
	}
}

func must(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
