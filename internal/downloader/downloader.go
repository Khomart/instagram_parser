package downloader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Downloader struct {
}

func NewDownloader() *Downloader {
	return &Downloader{}
}

func (d *Downloader) FetchPost(instagramURL string) string {
	shortcode, err := d.extractShortcode(instagramURL)
	if err != nil {
		log.Fatalf("Error extracting shortcode: %v\n", err)
	}

	fmt.Printf("Extracted shortcode: %s\n", shortcode)

	// Create download folder
	downloadFolder := filepath.Join("downloads", shortcode)
	if err := os.MkdirAll(downloadFolder, os.ModePerm); err != nil {
		log.Fatalf("Error creating download folder: %v\n", err)
	}

	// Fetch post metadata
	posts, err := d.fetchPostMetadata(shortcode)
	if err != nil {
		log.Fatalf("Error fetching post metadata: %v\n", err)
	}

	if len(posts) > 0 {
		// Save description if present
		if err := d.saveDescription(downloadFolder, posts[0].Description); err != nil {
			log.Fatalf("Error saving description: %v\n", err)
		}
	}

	for idx, post := range posts {
		fmt.Printf("Fetched post metadata: %+v\n", post)

		if err := d.downloadPost(downloadFolder, post, idx); err != nil {
			log.Fatalf("Error downloading post: %v\n", err)
		}
	}
	return downloadFolder
}

func (d *Downloader) downloadMedia(folder, filename, url string, mtime time.Time, filenameSuffix *string) (bool, error) {
	if filenameSuffix != nil {
		filename += "_" + *filenameSuffix
	}

	// Determine file extension from URL
	urlMatch := regexp.MustCompile(`\.[a-z0-9]*\?`).FindString(url)
	fileExtension := ""
	if urlMatch != "" {
		fileExtension = urlMatch[1 : len(urlMatch)-1]
	} else {
		fileExtension = url[len(url)-3:]
	}
	nominalFilename := filepath.Join(folder, filename+"."+fileExtension)

	// Check if the file already exists
	if _, err := os.Stat(nominalFilename); err == nil {
		fmt.Printf("%s exists\n", nominalFilename)
		return false, nil
	}

	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		return false, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	// Determine file extension from Content-Type header
	if contentType := resp.Header.Get("Content-Type"); contentType != "" {
		headerExtension := "." + strings.ToLower(strings.Split(strings.Split(contentType, ";")[0], "/")[1])
		if headerExtension == ".jpeg" {
			headerExtension = ".jpg"
		}
		filename += headerExtension
	} else {
		filename = nominalFilename
	}

	// Check if the file with the updated extension already exists
	if filename != nominalFilename {
		if _, err := os.Stat(filename); err == nil {
			fmt.Printf("%s exists\n", filename)
			return false, nil
		}
	}

	// Save the file
	out, err := os.Create(nominalFilename)
	if err != nil {
		return false, fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return false, fmt.Errorf("failed to save file: %w", err)
	}

	// Set the modification time
	if err := os.Chtimes(nominalFilename, time.Now(), mtime); err != nil {
		return false, fmt.Errorf("failed to set file modification time: %w", err)
	}

	return true, nil
}

// saveDescription saves the post description to a text file.
func (d *Downloader) saveDescription(folder, description string) error {
	if description == "" {
		return nil
	}

	filename := filepath.Join(folder, "description.txt")
	return os.WriteFile(filename, []byte(description), 0644)
}

// DownloadPost processes the post metadata and downloads the associated media.
func (d *Downloader) downloadPost(folder string, post Post, idx int) error {
	var mediaURL string
	if post.MediaType == "GraphVideo" {
		mediaURL = post.VideoURL
	} else {
		mediaURL = post.DisplayURL
	}

	// Create filename
	filename := fmt.Sprintf("%s_%d", post.ID, idx)

	// Download media
	mtime := time.Now() // Set modification time as current time for now
	if _, err := d.downloadMedia(folder, filename, mediaURL, mtime, nil); err != nil {
		return fmt.Errorf("error downloading media: %w", err)
	}

	fmt.Println("Media downloaded successfully:", filename)
	return nil
}
