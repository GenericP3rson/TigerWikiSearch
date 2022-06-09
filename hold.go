package main 

import (
	"encoding/json"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"github.com/GenericP3rson/TigerGo"
	"sort"
)

func main_hold() {

	conn := TigerGo.TigerGraphConnection{ // Create Connection
		Token:     "",
		Host:      "https://bleve.i.tgcloud.io",
		GraphName: "NotMyGraph",
		Username:  "tigergraph",
		Password:  "tigergraph",
	}

	conn.Token, _ = conn.GetToken() // Generate Token

	res, err := conn.RunInstalledQuery("tg_pagerank", map[string]interface{}{"v_type": "Doc", "e_type": "LINKS_TO"}) // Runs modified Pagerank
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

	type Values struct {
		Vertex_ID string 
		Content string
		Score float32
	}
	
	type QueryResults struct {
		Top_scores_heap []Values `json:"@@top_scores_heap"`
	}
	
	type Res struct {
		Version struct {
			Edition string 
			Api string
			Schema int
		} 
		Error bool 
		Message string 
		Results []QueryResults
	}

	var resp Res	
	json.Unmarshal([]byte(res), &resp) // Maps response to structure
	fmt.Println(resp.Results[0].Top_scores_heap)

	// mapping := bleve.NewIndexMapping() // Create a new mapping
	// mapping.TypeField = "doc" // Type of document is "doc"
	// index, err := bleve.New("wiki_graph.bleve", mapping) // Create a new document with the mapping
	// blogMapping := bleve.NewDocumentMapping() // Create a document mapping
	// mapping.AddDocumentMapping("doc", blogMapping)
	// contentFieldMapping := bleve.NewTextFieldMapping() // Based on a certain field
	// contentFieldMapping.Analyzer = "en"
	// blogMapping.AddFieldMappingsAt("content", contentFieldMapping) // Add field mapping at "content"
	index, err  := bleve.Open("wiki_graph.bleve")
	if err != nil {
		fmt.Println(err)
	}

	// for _, element := range resp.Results[0].Top_scores_heap { // Iterates through and loads
	// 	fmt.Println("Uploading...")
	// 	index.Index(element.Vertex_ID, element)
	// 	fmt.Println(element.Vertex_ID, element)
	// 	fmt.Println(index.DocCount())
	// }

	// pb := &Values{ 
	// 	Vertex_ID: "0",
	// 	Content:  "contact",
	// 	Score: -8.0,
	// }	
	// index.Index("0", pb)

	query := bleve.NewMatchQuery("language")
	search := bleve.NewSearchRequest(query)
	search.Size = 50
	search.SortBy([]string{"-_score", "-Score"})
	search.Fields = []string{"Score"}
	searchResults, err := index.Search(search)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(searchResults)
	fmt.Println("Score?", search.Score)
	
	type Result struct {
		Id string
		Score float64
	}

	results := []Result{}

	weight_1, weight_2 := 2.0, 0.5

	for i, hit := range searchResults.Hits {
		fmt.Println(i, hit.ID, hit.Score, hit.Fields, hit.Score + hit.Fields["Score"].(float64))
		results = append(results, Result{Id: hit.ID, Score: hit.Score * weight_1 + hit.Fields["Score"].(float64) * weight_2})
	}
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	fmt.Println(results)

	fmt.Println("Modified Results")
	for i, res := range results {
		fmt.Println(i, res.Id, res.Score)
	}
	
	// sb := `{"version":{"edition":"enterprise","api":"v2","schema":0},"error":false,"message":"", "result":[{"v_id":"1000000379","v_type":"Patient","attributes":{"patient_id":"1000000379","global_num":9368,"birth_year":1981,"infection_case":"overseas inflow","contact_number":0,"symptom_onset_date":"1970-01-01 00:00:00","confirmed_date":"2020-03-27 00:00:00","released_date":"1970-01-01 00:00:00","deceased_date":"1970-01-01 00:00:00","state":"isolated","disease":"","sex":"male"}},{"v_id":"6022000013","v_type":"Patient","attributes":{"patient_id":"6022000013","global_num":0,"birth_year":1954,"infection_case":"","contact_number":0,"symptom_onset_date":"1970-01-01 00:00:00","confirmed_date":"2020-02-25 00:00:00","released_date":"2020-03-31 00:00:00","deceased_date":"1970-01-01 00:00:00","state":"released","disease":"","sex":"male"}},{"v_id":"2000000166","v_type":"Patient","attributes":{"patient_id":"2000000166","global_num":7665,"birth_year":1993,"infection_case":"contact with patient","contact_number":0,"symptom_onset_date":"1970-01-01 00:00:00","confirmed_date":"2020-03-10 00:00:00","released_date":"1970-01-01 00:00:00","deceased_date":"1970-01-01 00:00:00","state":"isolated","disease":"","sex":"female"}}]}`

	// type Attributes struct {
	// 	Patient_id string 
	// 	Global_num int 
	// 	Birth_year int
	// 	Infection_case string 
	// 	Contact_number int 
	// 	Symptom_onset_date string
	// 	Confirmed_date string
	// 	Released_date string
	// 	Deceased_date string
	// 	State string
	// 	Disease string
	// 	Sex string
	// }

	// type Res struct { // Result structure
	// 	Version struct {
	// 		Edition string 
	// 		Api string
	// 		Schema int
	// 	}
	// 	Error bool 
	// 	Message string
	// 	Result []struct {
	// 		V_id string 
	// 		V_type string 
	// 		Attributes Attributes
	// 	}
	// }

	// var resp Res	
	// json.Unmarshal([]byte(sb), &resp) // Maps response to structure

	// fmt.Println(resp)

}
