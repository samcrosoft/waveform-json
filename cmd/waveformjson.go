package main
import (
	"github.com/samcrosoft/waveform-json"
	"github.com/samcrosoft/waveform-json/fetchers"
	"io/ioutil"
	"os"
	"flag"
	"strings"
	"runtime"
)

func main() {

	// run on full throttle, use all CPU
	runtime.GOMAXPROCS(runtime.NumCPU())

	var sOutputFileName string
	var sSource string
	var oFetcher fetchers.IMAGEFetcher
	flag.StringVar(&sOutputFileName, "o", "output.json", "the outfile filename, default is output.json")
	flag.StringVar(&sSource, "s", "", "the source of the waveform, e.g ./corpus/sample.png")
	flag.Parse()

	// get the source
	if(sSource == ""){
		os.Stderr.WriteString("Source Of Waveform cannot be empty")
		os.Exit(0)
	}
	// determine the fetcher
	oFetcher = determineFetcherTypeFromSource(sSource)
	// initiate the waveform calculations
	oWave := waveformjson.New(sSource, oFetcher)

	// save to file
	if oJson, err := oWave.GetWaveformJson(); err == nil{
		if err2 := ioutil.WriteFile(sOutputFileName, oJson, os.FileMode(os.O_WRONLY)); err2 != nil{
			// do nothing
		}
	}else{
		os.Stderr.WriteString("Not A Valid Png")
	}

}


func determineFetcherTypeFromSource(sSource string) fetchers.IMAGEFetcher{
	if sSource != "" && strings.HasPrefix(sSource, "http"){
		return &fetchers.URLFetcher{}
	}else{
		return &fetchers.FILEFetcher{}
	}
}