package internal

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"sort"
	"sync"
	"unicode"
)

var (
	//go:embed font_size.json
	fontSizeBytes     []byte
	onceSetupFontSize sync.Once
	fontSizeData      []fontSize
	emSize            float32 = 10.7
)

type fontSize struct {
	Char rune
	Size float32
}

func (f *fontSize) UnmarshalJSON(b []byte) error {
	var rawArray [3]float32
	if err := json.Unmarshal(b, &rawArray); err != nil {
		return err
	}

	if len(rawArray) != 3 {
		return errors.New("unexpected number of element in fontsize")
	}

	*f = fontSize{
		Char: rune(rawArray[0]),
		Size: rawArray[2],
	}

	return nil
}

const badgeSvgTmplString = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="{{.Width}}" height="20" role="img"
aria-label="{{.Title}}">
<title>{{.Title}}</title>
<linearGradient id="s" x2="0" y2="100%">
	<stop offset="0" stop-color="#bbb" stop-opacity=".1" />
	<stop offset="1" stop-opacity=".1" />
</linearGradient>
<clipPath id="r">
	<rect width="{{.Width}}" height="20" rx="3" fill="#fff" />
</clipPath>
<g clip-path="url(#r)">
	<rect width="{{.LeftWidth}}" height="20" fill="#555" />
	<rect x="{{.LeftWidth}}" width="{{.RightWidth}}" height="20" fill="#4c1" />
	<rect width="{{.Width}}" height="20" fill="url(#s)" />
</g>
<g fill="#fff" text-anchor="middle" font-family="Verdana,Geneva,DejaVu Sans,sans-serif"
	text-rendering="geometricPrecision" font-size="110">
	<text
		aria-hidden="true"
		x="{{.LabelTextX}}"
		y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)"
		textLength="{{.LabelTextWidth}}">{{.Label}}</text>
	<text
		x="{{.LabelTextX}}"
		y="140" transform="scale(.1)" fill="#fff" textLength="{{.LabelTextWidth}}">{{.Label}}</text>
	<text
		aria-hidden="true"
		x="{{.DescriptionTextX}}"
		y="150"
		fill="#010101"
		fill-opacity=".3"
		transform="scale(.1)"
		textLength="{{.DescriptionTextWidth}}"
	>{{.Description}}</text>
	<text
		x="{{.DescriptionTextX}}"
		y="140"
		transform="scale(.1)"
		fill="#fff"
		textLength="{{.DescriptionTextWidth}}"
	>{{.Description}}</text>
</g>
</svg>`

var badgeSvgTmpl = template.Must(template.New("badge").Parse(badgeSvgTmplString))

type badgeTmplData struct {
	Title                string
	Label                string
	Description          string
	LeftWidth            float32
	RightWidth           float32
	LabelTextWidth       float32
	LabelTextX           float32
	DescriptionTextX     float32
	DescriptionTextWidth float32
	Width                float32
}

func sizeOfRune(c rune) float32 {
	onceSetupFontSize.Do(func() {
		if err := json.Unmarshal(fontSizeBytes, &fontSizeData); err != nil {
			log.Printf("[WARN] error unmarshalling JSON: %v", err)
		}
		fontSizeBytes = nil
	})

	if unicode.IsControl(c) {
		return 0
	}

	log.Println(fontSizeData[0])

	i := sort.Search(len(fontSizeData), func(i int) bool { return fontSizeData[i].Char >= c })
	if i < len(fontSizeData) {
		return fontSizeData[i].Size
	}
	log.Println("not found code: ", c, " - ", i)
	return emSize
}

func sizeOfString(s string) float32 {
	var size float32
	for _, c := range s {
		size += sizeOfRune(c)
	}
	return size
}

func renderBadge(w io.Writer, label, description string) error {
	marginWidth := 5
	labelWidth := sizeOfString(label)
	descWidth := sizeOfString(description)
	leftWidth := labelWidth + float32(2*marginWidth)
	rightWidth := descWidth + float32(2*marginWidth)
	labelX := 0.5*labelWidth + float32(marginWidth)
	descriptionX := 0.5*descWidth + leftWidth + float32(marginWidth)
	return badgeSvgTmpl.Execute(w, badgeTmplData{
		Title:                fmt.Sprintf("%s: %s", label, description),
		Label:                label,
		Description:          description,
		LeftWidth:            leftWidth,
		RightWidth:           rightWidth,
		LabelTextWidth:       labelWidth / 0.1,
		LabelTextX:           labelX / 0.1,
		DescriptionTextWidth: descWidth / 0.1,
		DescriptionTextX:     descriptionX / 0.1,
		Width:                leftWidth + rightWidth,
	})
}
