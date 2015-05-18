package fetchers
import (
    "net/http"
    "io/ioutil"
    "bytes"
    "io"
)

type URLFetcher struct  {}

// this will fetch the io reader from a URL
// the argument should be the a string that represents the path to the image
// it should return the io-reader if the file is read correctly
func (u *URLFetcher) FetchImage(sURL string) (io.Reader, error){

    var oReader *bytes.Reader

    // Just a simple GET request to the image URL
    // We get back a *Response, and an error
    res, err := http.Get(sURL)

    // defer the closing of the body
    defer func() {
        res.Body.Close()
    }()

    if err != nil {
        return oReader,err
    }

    // We read all the bytes of the image
    // Types: data []byte
    var data []byte

    if data, err = ioutil.ReadAll(res.Body); err != nil {
        return oReader, err
    }

    oReader = bytes.NewReader(data)

    return oReader, nil
}