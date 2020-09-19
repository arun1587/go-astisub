package astisub_test

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/asticode/go-astisub"
	"github.com/stretchr/testify/assert"
)

func TestWebVTT(t *testing.T) {
	// Open
	s, err := astisub.OpenFile("./testdata/example-in.vtt")
	assert.NoError(t, err)
	assertSubtitleItems(t, s)
	// Comments
	assert.Equal(t, []string{"this a nice example", "of a VTT"}, s.Items[0].Comments)
	assert.Equal(t, []string{"This a comment inside the VTT", "and this is the second line"}, s.Items[1].Comments)
	// Regions
	assert.Equal(t, 2, len(s.Regions))
	assert.Equal(t, astisub.Region{ID: "fred", InlineStyle: &astisub.StyleAttributes{WebVTTLines: 3, WebVTTRegionAnchor: "0%,100%", WebVTTScroll: "up", WebVTTViewportAnchor: "10%,90%", WebVTTWidth: "40%"}}, *s.Regions["fred"])
	assert.Equal(t, astisub.Region{ID: "bill", InlineStyle: &astisub.StyleAttributes{WebVTTLines: 3, WebVTTRegionAnchor: "100%,100%", WebVTTScroll: "up", WebVTTViewportAnchor: "90%,90%", WebVTTWidth: "40%"}}, *s.Regions["bill"])
	assert.Equal(t, s.Regions["bill"], s.Items[0].Region)
	assert.Equal(t, s.Regions["fred"], s.Items[1].Region)
	// Styles
	assert.Equal(t, astisub.StyleAttributes{WebVTTAlign: "left", WebVTTPosition: "10%,start", WebVTTSize: "35%"}, *s.Items[1].InlineStyle)

	// No subtitles to write
	w := &bytes.Buffer{}
	err = astisub.Subtitles{}.WriteToWebVTT(w)
	assert.EqualError(t, err, astisub.ErrNoSubtitlesToWrite.Error())

	// Write
	c, err := ioutil.ReadFile("./testdata/example-out.vtt")
	assert.NoError(t, err)
	err = s.WriteToWebVTT(w)
	assert.NoError(t, err)
	assert.Equal(t, string(c), w.String())
}

func TestBroken1WebVTT(t *testing.T) {
	// Open bad, broken WebVTT file
	_, err := astisub.OpenFile("./testdata/broken-1-in.vtt")
	assert.Nil(t, err)
}

func TestWebVTTWithVoiceName(t *testing.T) {
	voiceName1 := "Roger Bingham"
	voiceName2 := "Bingham"
	voiceName3 := "Lee"
	testData := `WEBVTT

	NOTE this a example with voicename

	1
	00:02:34.00 --> 00:02:35.00
	<v.first.local Roger Bingham> I'm the fist speaker

	2
	00:02:34.00 --> 00:02:35.00
	<v Bingham> I'm the second speaker

	3
	00:00:04.000 --> 00:00:08.000
	<v Lee>What are you doing here?</v>

	4
	00:00:04.000 --> 00:00:08.000
	<v Bob>Incorrect tag?</vi>`

	s, err := astisub.ReadFromWebVTT(strings.NewReader(testData))
	assert.NoError(t, err)

	assert.Len(t, s.Items, 4)
	assert.Equal(t, voiceName1, s.Items[0].Lines[0].VoiceName)
	assert.Equal(t, voiceName2, s.Items[1].Lines[0].VoiceName)
	assert.Equal(t, voiceName3, s.Items[2].Lines[0].VoiceName)
	assert.Equal(t, "", s.Items[3].Lines[0].VoiceName)
}

func TestParseTextWebVTT(t *testing.T) {

	t.Run("When both voice tags are available", func(t *testing.T) {
		testData := `<v Bob>Correct tag</v>`

		s := astisub.ParseTextWebVTT(testData)
		assert.Equal(t, "Bob", s.VoiceName)
		assert.Equal(t, 1, len(s.Items))
		assert.Equal(t, "Correct tag", s.Items[0].Text)
	})

	t.Run("When there is no end tag", func(t *testing.T) {
		testData := `<v Bob> Text without end tag`

		s := astisub.ParseTextWebVTT(testData)
		assert.Equal(t, "Bob", s.VoiceName)
		assert.Equal(t, 1, len(s.Items))
		assert.Equal(t, "Text without end tag", s.Items[0].Text)
	})

	t.Run("When the end tag is correct", func(t *testing.T) {
		testData := `<v Bob>Incorrect end tag</vi>`

		s := astisub.ParseTextWebVTT(testData)
		assert.Equal(t, "", s.VoiceName)
		assert.Equal(t, 0, len(s.Items))
	})
}
