package finesmith

import (
	"io/ioutil"
	"os"

	"path"

	"strings"

	"github.com/flosch/pongo2"
	"github.com/yosssi/gohtml"
)

func processData(ctx pongo2.Context, layoutPath string, outputFile string) {
	/*
		//outputs json next to index.html
		b, err := json.MarshalIndent(ctx, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		os.MkdirAll(path.Dir(outputFile), 0755)
		ioutil.WriteFile(fmt.Sprintf("%s.json", outputFile), b, 0755)
	*/
	pctx, err := pongo2.FromFile(layoutPath)
	if err != nil {
		panic(err)
	}

	os.MkdirAll(path.Dir(outputFile), 0755)
	bht, err := pctx.ExecuteBytes(ctx)
	if err != nil {
		panic(err)
	}
	bb := gohtml.FormatBytes(bht)
	ioutil.WriteFile(outputFile, bb, 0644)
}

func createPage(job PrismicPageJob, layoutBaseDir string, outputDir string) {
	uid := job.ContentData["uid"].(string)
	fpath := strings.Replace(uid, ".", "/", -1)

	ctx := make(map[string]interface{})
	ctx["site"] = job.SiteData
	ctx["content"] = job.ContentData

	fctx := make(map[string]interface{})
	fctx["prismic"] = ctx
	fctx["basePath"] = "/" + job.Country
	fctx["layout"] = job.Layout

	outFilename := "index.html"
	if strings.Contains(job.Layout, "404") {
		outFilename = "404.html"
	}
	outputFilepath := path.Join(outputDir, fpath, outFilename)
	if job.PathPrefix != "" {
		outputFilepath = path.Join(outputDir, job.PathPrefix, fpath, outFilename)
	}
	processData(fctx, path.Join(layoutBaseDir, job.Layout), outputFilepath)
}

func CreatePageWorker(channel chan *PrismicPageJob, layoutBaseDir string, outputDir string, done chan bool) {
	for work := range channel {
		createPage(*work, layoutBaseDir, outputDir)
	}
	done <- true
}
