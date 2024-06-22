package interfaces

type Downloader interface {
	FetchPost(url string) string
}

type Parser interface {
	SummarizeContent(folderPath string) (string, error)
}
