package finesmith

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"log"

	yaml "gopkg.in/yaml.v2"
)

// TemplateProcessor is reading and processing template files supporting prismic.io and pongo2
type TemplateProcessor struct {
	prismicChannel chan *PrismicPageJob
	prismicWorker  *PrismicWorker
}

//NewTemplateProcessor creates a new instance with given config
func NewTemplateProcessor(prismicChannel chan *PrismicPageJob, prismicWorker *PrismicWorker) *TemplateProcessor {
	return &TemplateProcessor{
		prismicChannel,
		prismicWorker,
	}
}

// Fix conversion with yaml
func (t *TemplateProcessor) convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = t.convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = t.convert(v)
		}
	}
	return i
}

//ProcessFile reads templates file and processes it accordingly. Method is supposed to be used with io.Walk
func (t *TemplateProcessor) ProcessFile(fpath string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() || path.Ext(fpath) != ".html" {
		return nil
	}

	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		log.Fatalln("Cannot read file:", fpath)
		return err
	}

	var tc interface{}
	if err := yaml.Unmarshal(data, &tc); err != nil {
		log.Fatalln("Invalid template yaml file", fpath, err)
		return err
	}

	templateContent := t.convert(tc).(map[string]interface{})

	layout, hasLayout := templateContent["layout"]
	if !hasLayout {
		return errors.New("No layout defined")
	}

	if prismic, usesPrismic := templateContent["prismic"]; usesPrismic {
		country, hasCountry := templateContent["country"]
		if !hasCountry {
			return fmt.Errorf("No country defined in %s", fpath)
		}

		prismicSettings, ok := prismic.(map[string]interface{})
		if !ok {
			return fmt.Errorf("Prismic config not a map %s", prismic)
		}

		t.processPrismicPage(layout.(string), country.(string), prismicSettings)
		return nil
	}

	//TODO: Process normal template

	return nil
}

func mapToPrismicQuery(key string, page map[string]interface{}) (query PrismicQuery) {
	b, _ := json.Marshal(page)
	json.Unmarshal(b, &query)
	query.QueryKey = key
	return
}

func (t *TemplateProcessor) processPrismicPage(layout string, country string, prismicData map[string]interface{}) error {
	site, hasSite := prismicData["site"].(map[string]interface{})
	if !hasSite {
		return errors.New("No site content available")
	}

	siteResults, err := t.prismicWorker.Query(mapToPrismicQuery("site", site))
	if err != nil {
		log.Fatalln("Cannot fetch site results", err)
		return err
	}
	if len(siteResults) <= 0 {
		return errors.New("No site results, stopping")
	}

	siteResult := siteResults[0]
	if len(siteResults) != 1 {
		log.Fatalln("Multiple site results available, selecting first one. Number results:", len(siteResults))
	}

	content, hasContent := prismicData["content"].(map[string]interface{})
	if !hasContent {
		return errors.New("No content query available")
	}

	contentResult, err := t.prismicWorker.Query(mapToPrismicQuery("content", content))
	if err != nil {
		log.Fatalln("Cannot fetch content", err)
		return err
	}

	for _, documentData := range contentResult {
		job := PrismicPageJob{
			Layout:      layout,
			SiteData:    siteResult,
			ContentData: documentData,
			Country:     country,
		}
		t.prismicChannel <- &job
	}

	return nil
}
