package main 

import (
	"encoding/json"
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"github.com/GenericP3rson/TigerGo"
)

func process_data() { // This function processes the data from the graph database and uploads it to Bleve

	conn := TigerGo.TigerGraphConnection{ // Create Connection
		Token:     "",
		Host:      "https://SUBDOMAIN.i.tgcloud.io", // Replace with subdomain
		GraphName: "MyGraph", // Replace with graphname
		Username:  "tigergraph",
		Password:  "tigergraph", // Replace with password
	}
	conn.Token, _ = conn.GetToken() // Generate Token

	res, err := conn.RunInstalledQuery("tg_pagerank", map[string]interface{}{"v_type": "Doc", "e_type": "LINKS_TO"}) // Runs modified Pagerank
	if err != nil {
		fmt.Println(err)
	}

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

	mapping := bleve.NewIndexMapping() // Create a new mapping
	mapping.TypeField = "doc" // Type of document is "doc"
	index, err := bleve.New("wiki_graph.bleve", mapping) // Create a new document with the mapping
	blogMapping := bleve.NewDocumentMapping() // Create a document mapping
	mapping.AddDocumentMapping("doc", blogMapping)
	contentFieldMapping := bleve.NewTextFieldMapping() // Based on a certain field
	contentFieldMapping.Analyzer = "en"
	blogMapping.AddFieldMappingsAt("content", contentFieldMapping) // Add field mapping at "content"
	if err != nil {
		fmt.Println(err)
	}

	for _, element := range resp.Results[0].Top_scores_heap { // Iterates through and loads
		fmt.Println("Uploading...")
		index.Index(element.Vertex_ID, element)
	}

}
