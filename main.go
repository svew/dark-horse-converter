package main

import (
	"archive/tar"
	"archive/zip"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"mime"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var input = flag.String("i", "", "Path of a comic book file in tar format to convert (required)")
var output = flag.String("o", "", "Output directory (if not specified, result will be placed in the same directory as the input)")

type inputFile struct {
	header  *tar.Header
	content []byte
}

type outputFile struct {
	header  *zip.FileHeader
	content []byte
}

type page struct {
	SortOrder   int    `json:"sort_order"`
	SourceImage string `json:"src_image"`
	MimeType    string `json:"mime_type"`
}

type manifest struct {
	RightToLeft bool   `json:"is_rtl"`
	Pages       []page `json:"pages"`
}

func (p manifest) Len() int {
	return len(p.Pages)
}

func (p manifest) Swap(i, j int) {
	p.Pages[i], p.Pages[j] = p.Pages[j], p.Pages[i]
}

func (p manifest) Less(i, j int) bool {
	if p.RightToLeft {
		return p.Pages[i].SortOrder > p.Pages[j].SortOrder
	} else {
		return p.Pages[i].SortOrder < p.Pages[j].SortOrder
	}
}

func checkedClose(c io.Closer, err *error) {
	if cerr := c.Close(); cerr != nil && *err == nil {
		*err = cerr
	}
}

func read(path string) (map[string]inputFile, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	archive := tar.NewReader(file)
	files := make(map[string]inputFile)

	for {
		header, err := archive.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		content, err := ioutil.ReadAll(archive)

		if err != nil {
			return nil, err
		}

		files[header.Name] = inputFile{header: header, content: content}
	}

	return files, nil
}

func convert(inputFiles map[string]inputFile) ([]outputFile, error) {
	file, ok := inputFiles["manifest.json"]

	if !ok {
		return nil, errors.New("Manifest not found inside archive")
	}

	manifest := new(manifest)

	if err := json.Unmarshal(file.content, manifest); err != nil {
		return nil, err
	}

	if len(manifest.Pages) == 0 {
		return nil, errors.New("No pages found inside archive")
	}

	sort.Sort(manifest)

	outputFiles := make([]outputFile, len(manifest.Pages))
	numberOfDigits := int(math.Floor(math.Log10(float64(len(manifest.Pages)))) + 1)

	for i, page := range manifest.Pages {
		image, ok := inputFiles[page.SourceImage]

		if !ok {
			return nil, errors.New("Archive is missing one or more pages")
		}

		fileInfo := image.header.FileInfo()
		header, err := zip.FileInfoHeader(fileInfo)

		if err != nil {
			return nil, err
		}

		extensions, err := mime.ExtensionsByType(page.MimeType)

		if err != nil {
			return nil, err
		}

		if len(extensions) == 0 {
			return nil, errors.New("Unknown image type found inside archive")
		}

		ext := extensions[0]
		if ext == ".jfif" {
			ext = ".jpeg"
		}
		format := fmt.Sprintf("%%0%dd%s", numberOfDigits, ext)
		header.Name = fmt.Sprintf(format, i+1)
		outputFiles[i] = outputFile{header: header, content: image.content}
	}

	return outputFiles, nil
}

func write(name string, files []outputFile) (err error) {
	file, err := os.Create(name)

	if err != nil {
		return err
	}

	defer checkedClose(file, &err)

	zipWriter := zip.NewWriter(file)
	defer checkedClose(zipWriter, &err)

	for _, file := range files {
		writer, err := zipWriter.CreateHeader(file.header)

		if err != nil {
			return err
		}

		_, err = writer.Write(file.content)

		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	flag.Parse()

	if *input == "" {
		flag.Usage()
		os.Exit(1)
	}

	dir, file := filepath.Split(*input)

	if *output != "" {
		dir = *output
	}

	extension := filepath.Ext(file)
	name := strings.TrimSuffix(file, extension)
	*output = filepath.Join(dir, name+".cbz")

	inputFiles, err := read(*input)

	if err != nil {
		log.Fatal(err)
	}

	outputFiles, err := convert(inputFiles)

	if err != nil {
		log.Fatal(err)
	}

	if write(*output, outputFiles) != nil {
		log.Fatal(err)
	}

	log.Println("Comic book successfully converted")
}
