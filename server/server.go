package server

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/kiwamunet/image-optim/compare"
	"github.com/kiwamunet/image-optim/operation"

	"code.cloudfoundry.org/bytefmt"
)

const (
	addr       = ":8080"
	SRC        = "src"
	CMD        = "cmd"
	INPUTFILE  = "inputfile"
	OUTPUTFILE = "outputfile"
)

type Opt struct {
	URL              string
	OriginalFileName string
	TargetFileName   string
	Commands         []string
}

func Serve() {
	s := &http.Server{
		Addr: addr,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("/image-optim", func(w http.ResponseWriter, r *http.Request) {
		imageOptim(w, r)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})
	s.Handler = mux
	s.ListenAndServe()
}

func imageOptim(w http.ResponseWriter, r *http.Request) {
	urls, err := url.QueryUnescape(r.URL.Scheme + r.URL.Host + r.RequestURI)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	o := &Opt{
		URL: urls,
	}

	// src url parse
	srcURL := r.URL.Query().Get(SRC)
	if srcURL == "" {
		http.Error(w, fmt.Sprintf("err %s", "src url is unavalable"), http.StatusInternalServerError)
		return
	}

	// filename set
	_, s := path.Split(srcURL)
	o.OriginalFileName = s
	o.TargetFileName = "tgt-" + o.OriginalFileName

	// file download
	err = operation.DownLoad(srcURL, o.OriginalFileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// copy
	err = operation.CopyFile(o.OriginalFileName, o.TargetFileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// command チェック
	commands := r.URL.Query().Get(CMD)
	for _, x := range strings.Split(commands, ",") {
		arg := strings.Split(x, " ")
		arg, cmdStr := renameFileName(arg, o.TargetFileName)
		o.Commands = append(o.Commands, cmdStr)
		log.Println(o.Commands)
		err := operation.ExeCmd(arg[0], arg[1:])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// exist
		if operation.Exist("_" + o.TargetFileName) {
			if err = operation.CopyFile("_"+o.TargetFileName, o.TargetFileName); err != nil {
				return
			}
			if err = os.Remove("_" + o.TargetFileName); err != nil {
				return
			}
		}
	}

	// get bynary
	original, err := ioutil.ReadFile(o.OriginalFileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	output, err := ioutil.ReadFile(o.TargetFileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ssim, psnr, err := compare.ImageComp(o.OriginalFileName, o.TargetFileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if t, err := ParseAssets("assets/index.tpl"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		t.Execute(w, struct {
			URL              string
			OriginalMimeType string
			DstMimeType      string
			OriginalImage    string
			DstImage         string
			OriginalSize     string
			OutputSize       string
			Commands         []string
			Compressibility  string
			SSIM             string
			PSNR             string
		}{
			URL:              "http://localhost:8080" + o.URL,
			OriginalMimeType: http.DetectContentType(original),
			DstMimeType:      http.DetectContentType(output),
			OriginalImage:    base64.StdEncoding.EncodeToString(original),
			DstImage:         base64.StdEncoding.EncodeToString(output),
			OriginalSize:     bytefmt.ByteSize(uint64(binary.Size(original))),
			OutputSize:       bytefmt.ByteSize(uint64(binary.Size(output))),
			Commands:         o.Commands,
			Compressibility:  "(" + strconv.FormatFloat(Round((float64(binary.Size(output))/float64(binary.Size(original))*100), .5, 2), 'f', 2, 64) + "%)",
			SSIM:             strconv.FormatFloat(ssim, 'f', 4, 64),
			PSNR:             strconv.FormatFloat(psnr, 'f', 4, 64),
		})
	}
}
func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
func ParseAssets(path string) (*template.Template, error) {
	ns := template.New("opt")
	src, err := Asset(path)
	if err != nil {
		return nil, err
	}
	return ns.Parse(string(src))
}

func renameFileName(args []string, targetFileName string) ([]string, string) {
	var cmdStr = ""
	for i, v := range args {
		switch v {
		case INPUTFILE:
			args[i] = targetFileName
		case OUTPUTFILE:
			args[i] = "_" + targetFileName
		}
		cmdStr += " " + args[i]
	}
	return args, cmdStr
}
