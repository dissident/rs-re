package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnwrapRedirectGoogleAlerts(t *testing.T) {
	link := "https://www.google.com/url?rct=j&sa=t&url=https://example.com/news/article&ct=ga&cd=CAIyGmVi&usg=AOvVaw0"
	assert.Equal(t, "https://example.com/news/article", unwrapRedirect(link))
}

func TestUnwrapRedirectGoogleWithoutURLParam(t *testing.T) {
	link := "https://www.google.com/url?rct=j"
	assert.Equal(t, link, unwrapRedirect(link))
}

func TestUnwrapRedirectPlainLink(t *testing.T) {
	link := "https://example.com/news/article"
	assert.Equal(t, link, unwrapRedirect(link))
}

func TestUnwrapRedirectOtherGooglePath(t *testing.T) {
	link := "https://www.google.com/search?q=test"
	assert.Equal(t, link, unwrapRedirect(link))
}

func TestUnwrapRedirectInvalidURL(t *testing.T) {
	link := "://not-a-url"
	assert.Equal(t, link, unwrapRedirect(link))
}

func TestMemTitlesContains(t *testing.T) {
	titles := []string{"one", "two"}
	assert.True(t, memTitlesContains(titles, "two"))
	assert.False(t, memTitlesContains(titles, "three"))
}
