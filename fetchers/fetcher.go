package fetchers
import (
    "io"
)

// create an interface that will handle the image fetching interface
type IMAGEFetcher interface  {
    FetchImage(sFileSource string) (io.Reader, error)
}