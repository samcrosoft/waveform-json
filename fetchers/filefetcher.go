package fetchers
import (
    "io/ioutil"
    "bytes"
    "io"
)

type FILEFetcher struct  {}

// this will fetch the io reader from a local file
// the argument should be the a string that represents the path to the image
// it should return the io-reader if the file is read correctly
func (f *FILEFetcher) FetchImage(sFileNameOrUrl string) (io.Reader, error) {
    // fetch the image and return it
    return fetchFromFile(sFileNameOrUrl)
}

// this method will perform the actual reading of the file
func fetchFromFile(sFileName string) (io.Reader, error) {
    var oReader *bytes.Reader
    if oImageData, err := ioutil.ReadFile(sFileName); err != nil {
        return oReader,err
    }else{
        oReader = bytes.NewReader(oImageData)
        return oReader, nil
    }
}