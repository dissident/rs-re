package tg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatMessageStripsTags(t *testing.T) {
	message := "Abu Dhabi's BlueFive co-leads $3B <b>round</b> for China's Kling AI"
	expected := "Abu Dhabi's BlueFive co-leads $3B round for China's Kling AI"
	assert.Equal(t, expected, formatMessage(message))
}

func TestFormatMessageUnescapesEntities(t *testing.T) {
	message := "$3B &lt;b&gt;round&lt;/b&gt; &amp; more"
	assert.Equal(t, "$3B round & more", formatMessage(message))
}

func TestFormatMessageReplacesBrWithNewline(t *testing.T) {
	assert.Equal(t, "line one\nline two", formatMessage("line one<br />line two"))
	assert.Equal(t, "line one\nline two", formatMessage("line one<br>line two"))
}

func TestFormatMessagePlainText(t *testing.T) {
	message := "Just a plain title"
	assert.Equal(t, message, formatMessage(message))
}
