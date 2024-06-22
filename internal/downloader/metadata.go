package downloader

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Post represents an Instagram post.
type Post struct {
	DisplayURL  string `json:"display_url"`
	MediaType   string `json:"__typename"`
	ID          string `json:"id"`
	VideoURL    string `json:"video_url,omitempty"`
	Description string `json:"edge_media_to_caption.edges.node.text"`
}

// fetchPostMetadata fetches the post metadata using the shortcode.
func (D *Downloader) fetchPostMetadata(shortcode string) ([]Post, error) {
	queryHash := "2b0673e0dc4580674a88d426fe00ea90"
	variables := map[string]string{"shortcode": shortcode}
	variablesJSON, err := json.Marshal(variables)
	if err != nil {
		return nil, err
	}

	graphqlURL := "https://www.instagram.com/graphql/query/"
	params := fmt.Sprintf("?query_hash=%s&variables=%s", queryHash, url.QueryEscape(string(variablesJSON)))

	response, err := getJSON(graphqlURL + params)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			ShortcodeMedia struct {
				EdgeSidecarToChildren struct {
					Edges []struct {
						Node Post `json:"node"`
					} `json:"edges"`
				} `json:"edge_sidecar_to_children"`
				DisplayURL         string `json:"display_url"`
				MediaType          string `json:"__typename"`
				ID                 string `json:"id"`
				VideoURL           string `json:"video_url"`
				EdgeMediaToCaption struct {
					Edges []struct {
						Node struct {
							Text string `json:"text"`
						} `json:"node"`
					} `json:"edges"`
				} `json:"edge_media_to_caption"`
			} `json:"shortcode_media"`
		} `json:"data"`
	}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, err
	}

	posts := make([]Post, 0)
	if len(result.Data.ShortcodeMedia.EdgeSidecarToChildren.Edges) > 0 {
		for _, edge := range result.Data.ShortcodeMedia.EdgeSidecarToChildren.Edges {
			post := edge.Node
			if len(result.Data.ShortcodeMedia.EdgeMediaToCaption.Edges) > 0 {
				post.Description = result.Data.ShortcodeMedia.EdgeMediaToCaption.Edges[0].Node.Text
			}
			posts = append(posts, post)
		}
	} else {
		media := result.Data.ShortcodeMedia
		post := Post{
			DisplayURL: media.DisplayURL,
			MediaType:  media.MediaType,
			ID:         media.ID,
			VideoURL:   media.VideoURL,
		}
		if len(media.EdgeMediaToCaption.Edges) > 0 {
			post.Description = media.EdgeMediaToCaption.Edges[0].Node.Text
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// getJSON performs an HTTP GET request and returns the response body as bytes.
func getJSON(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}
