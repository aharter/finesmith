package finesmith

import (
	"encoding/json"
	"fmt"

	"github.com/SoCloz/goprismic"
	"github.com/SoCloz/goprismic/fragment"
	"github.com/SoCloz/goprismic/fragment/link"
	"github.com/SoCloz/goprismic/proxy"
)

const prismicDefaultPageSize = 100
const prismicMaxPageSize = 100

type PrismicPageJob struct {
	Layout      string
	ContentData map[string]interface{}
	SiteData    map[string]interface{}
	Country     string
	PathPrefix  string
}

type PrismicQuery struct {
	QueryKey   string  `json:"-"`
	Query      string  `json:"query,omitempty"`
	Bookmark   string  `json:"bookmark,omitempty"`
	DocumentID string  `json:"documentID,omitempty"`
	Orderings  string  `json:"orderings,omitempty"`
	FormName   string  `json:"formName,omitempty"`
	LinkDepth  int     `json:"linkDepth,omitempty"`
	Ref        *string `json:"ref,omitempty"`
	PageSize   *int    `json:"pageSize,omitempty"`
	AllPages   *bool   `json:"allPages,omitempty"`
}

type PrismicWorker struct {
	api      *proxy.Proxy
	worker   chan *PrismicPageJob
	resolver func(link.Link) string
}

func NewPrismicWorker(url string, token string, worker chan *PrismicPageJob, resolver func(link.Link) string) *PrismicWorker {
	if api, err := proxy.New(url, token, goprismic.DefaultConfig, proxy.Config{CacheSize: 1000}); err == nil {
		return &PrismicWorker{api, worker, resolver}
	}

	return nil
}

func (p *PrismicWorker) fetchSubdocuments(inter fragment.Interface, linkDepth int) fragment.Interface {
	if linkDepth <= 0 {
		return inter
	}

	switch t := inter.(type) {
	case *fragment.Link:
		if link, ok2 := t.Link.(*link.DocumentLink); linkDepth > 0 && ok2 {
			parentPage := PrismicQuery{
				DocumentID: link.Document.Id,
				LinkDepth:  linkDepth - 1,
			}
			parentDocuments, _ := p.prisimicLookup(parentPage)
			if len(parentDocuments) <= 0 {
				b, _ := json.MarshalIndent(inter, "", "  ")
				fmt.Println("Issues with:", string(b))
				return inter
			}
			// we only have 1 parent document
			parentDocument := parentDocuments[0]

			// Make sure URL is resolved
			parentDocument.ResolveLinks(p.resolver)
			// Replace link with goprismic.DocumentLink containing the data
			t.Link = parentDocument.AsDocumentLink()
		}
		break
	case *fragment.Group:
		for _, groupFragments := range *t {
			for groupFragmentName, groupFragment := range groupFragments {
				groupFragments[groupFragmentName] = p.fetchSubdocuments(groupFragment, linkDepth)
			}
		}
	}

	return inter
}

func (p *PrismicWorker) processFragmentList(fragments *fragment.Fragments, linkDepth int) {
	for _, currentFragment := range *fragments {
		for fragmentPieceIndex, fragmentPiece := range currentFragment {
			currentFragment[fragmentPieceIndex] = p.fetchSubdocuments(fragmentPiece, linkDepth)
		}
	}
}

func (p *PrismicWorker) prismicSearch(page PrismicQuery) (*goprismic.SearchResult, error) {
	var query string
	const QueryByID = "[[:d = at(document.id, \"%s\")]]"
	if page.Query != "" {
		query = page.Query
	} else if page.Bookmark != "" {
		query = fmt.Sprintf(QueryByID, p.api.Direct().Data.Bookmarks[page.Bookmark])
	} else if page.DocumentID != "" {
		query = fmt.Sprintf(QueryByID, page.DocumentID)
	}

	if page.FormName == "" {
		page.FormName = "everything"
	}

	pageSize := prismicDefaultPageSize
	if page.PageSize != nil {
		if *page.PageSize > prismicMaxPageSize {
			pageSize = prismicMaxPageSize
		} else {
			pageSize = *page.PageSize
		}
	}

	if page.Ref == nil {
		return p.api.Search().
			Form(page.FormName).
			PageSize(pageSize).
			Query(query).Submit()
	}

	return p.api.Direct().
		ForceRef(*page.Ref).
		Form(page.FormName).
		PageSize(pageSize).
		Query(query).Submit()
}

func (p *PrismicWorker) prisimicLookup(page PrismicQuery) ([]goprismic.Document, error) {
	var result *goprismic.SearchResult
	var err error
	prismicDocuments := make([]goprismic.Document, 0)
	if page.AllPages != nil && *page.AllPages {
		*page.AllPages = false
		page.PageSize = new(int)
		*page.PageSize = prismicMaxPageSize

		result, err = p.prismicSearch(page)
		for err == nil {
			prismicDocuments = append(prismicDocuments, result.Results...)
			if result.TotalPages != *page.PageSize {
				break
			}
			result, err = p.prismicSearch(page)
		}
	} else {
		// If 'page.AllPages' is not set
		result, err = p.prismicSearch(page)
		prismicDocuments = append(prismicDocuments, result.Results...)
	}

	if err != nil {
		return nil, err
	}

	for _, currentDocument := range prismicDocuments {
		currentDocument.ResolveLinks(p.resolver)
		for _, currentFragments := range currentDocument.Fragments {
			p.processFragmentList(&currentFragments, page.LinkDepth)
		}
	}

	return prismicDocuments, nil
}

func (p *PrismicWorker) Query(page PrismicQuery) ([]map[string]interface{}, error) {
	documents, err := p.prisimicLookup(page)
	if err != nil {
		return nil, err
	}

	results := make([]map[string]interface{}, len(documents))
	for index, currentDocument := range documents {
		// Move fragments into temporary variable to avoid unmarshaling it the fragments twice
		docFragments := currentDocument.Fragments
		currentDocument.Fragments = nil

		by, _ := json.Marshal(currentDocument)
		resultMap := make(map[string]interface{})
		json.Unmarshal(by, &resultMap)

		for _, fragment := range docFragments {
			by, _ := json.Marshal(fragment)
			dataMap := make(map[string]interface{})
			json.Unmarshal(by, &dataMap)
			resultMap["data"] = dataMap
		}
		results[index] = resultMap
	}

	return results, nil
}
