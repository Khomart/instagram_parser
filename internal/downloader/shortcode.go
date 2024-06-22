package downloader

import (
	"fmt"
	"net/url"
	"strings"
)

func (d *Downloader) extractShortcode(instagramURL string) (string, error) {
	parsedURL, err := url.Parse(instagramURL)
	if err != nil {
		return "", err
	}

	segments := strings.Split(parsedURL.Path, "/")
	if len(segments) > 2 && segments[1] == "p" {
		return segments[2], nil
	}

	return "", fmt.Errorf("invalid Instagram post URL")
}
