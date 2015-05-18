package waveformjson

import (
	"image"
	"image/png"
	"math"
	"encoding/json"
	"github.com/samcrosoft/waveform-json/fetchers"
	"errors"
	"github.com/samcrosoft/magic"
	"github.com/samcrosoft/magic/types"
	"bytes"
	"io"
)

const (
	IMAGE_TYPE = "png"
)

func init() {
	// Important :- Set the image format type to PNG
	image.RegisterFormat(IMAGE_TYPE, IMAGE_TYPE, png.Decode, png.DecodeConfig)
}

// this struct holds the source of the waveform and also the fetcher type that is to be used
type Waveform struct  {
	Source string
	Fetcher fetchers.IMAGEFetcher
}

// return a new pointer to the waveform type
func New(sSource string, oFetcher fetchers.IMAGEFetcher) *Waveform{
	return &Waveform{
		Source:sSource,
		Fetcher:oFetcher,
	}
}

// this will return the image filename/url to be converted to json points
func (w *Waveform) getSource() (string, error){
	if sSource := w.Source; sSource == ""{
		return "", (errors.New("The Image Source Must Be Set"))
	}
	return w.Source, nil
}

// this will return the fetcher to be used to fetch the image
func (w *Waveform) retrieveFetcher() fetchers.IMAGEFetcher {
	return w.Fetcher;
}

// this method will get the waveform from the source image
// it would then detect the float points as a float array that would be encoded and returned as json
// which is in effect a byte array
func (w *Waveform) GetWaveformJson() ([]byte, error) {

	var (
		sPath string
		oErr error
		oImage image.Image
		err error
		aReturn []byte
	)
	if sPath, oErr = w.getSource(); oErr != nil{
		return nil, oErr
	}
	imgReader,_ := w.retrieveFetcher().FetchImage(sPath)
	_ = imgReader
	// validate the image to be sure it has png type
	iPNGMinHeaderSize :=  int(types.PNGType{}.HeaderSize())
	aBytes := make([]byte,iPNGMinHeaderSize)
	imgReader.Read(aBytes)
	if bValidPng, pngError := magic.IsBytesContentAValidType(aBytes, &types.PNGType{}); pngError != nil || bValidPng == false{
		return aReturn, errors.New("Not a valid png")
	}

	// glue the bytes read for the png validation back to the reader
	oImageMultiReader := io.MultiReader(bytes.NewReader(aBytes), imgReader)

	if oImage, _, err = image.Decode(oImageMultiReader); err != nil{
		return nil, err
	}

	// get the image bounds and do the calculations
	bounds := oImage.Bounds()
	iWidth := bounds.Dx()
	iHeight := bounds.Dy()
	iHalf := int(math.Ceil(float64(iHeight/2)))

	// count up from top to half, find the point
	var fValPerPx float64
	fValPerPx = (1.0 / float64(iHalf))

	aPoints := make([]float64, iWidth)

	for i:=0;i<iWidth ;i++ {
		fPoint := calculateFirstAlphaPoint(&oImage, i, iHeight, iWidth, fValPerPx)
		aPoints[i] = fPoint
	}

	// write to json
	return computedValuesToJson(aPoints)
}

// this method will encode the computed points to json
func computedValuesToJson(aComputedValue []float64) ([]byte, error){
	if oJson, err := json.Marshal(aComputedValue); err != nil{
		return nil, err
	}else {
		return oJson, nil
	}
}


// This method will get the first point across the height of the image that is transparent
func calculateFirstAlphaPoint(img *image.Image, iXPos int, iHeight int, iWidth int, fValPerPx float64)float64{

	oImg := *img
	iHalf := int(math.Ceil(float64(iHeight/2)))
	var iFloatVal float64

	var iCheck  = iHeight
	for iCheck >= iHalf{
		iYPos := iCheck
		oPixel := oImg.At(iXPos, iYPos)
		_,_,_, a := oPixel.RGBA()
		if a == 0 && iYPos < iHeight{
			iFloatVal = fValPerPx * float64((iCheck - iHalf) + 1)
			break
		}
		iCheck--
	}

	return iFloatVal
}


