package main
import (
   "fmt"
   "github.com/gin-gonic/gin"
   "github.com/gin-contrib/cors"
   "net/http"
   "time"
	"github.com/blevesearch/bleve/v2"
	"sort"
)
func main() {

	index, err  := bleve.Open("wiki_graph.bleve") // Opening the Bleve document
	if err != nil {
		fmt.Println(err)
	}

	r := gin.Default() // Default router

	r.Use(cors.New(cors.Config{ // CORS!
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

   r.POST("/search", func(c *gin.Context) { // Create endpoint /search
		term := c.Query("term") // Takes in a parameter of the query
		query := bleve.NewMatchQuery(term) // Searches for matches
		search := bleve.NewSearchRequest(query)
		search.Size = 50
		search.SortBy([]string{"-_score", "-Score"}) // First sort by their order then by Pageranke
		search.Fields = []string{"Score", "Content"} // Add Score and Content field
		searchResults, err := index.Search(search) // Initial score results
		if err != nil {
			fmt.Println(err)
		}
		
		type Result struct {
			Id string
			Score float64
			Content string
		}

		results := []Result{}
		weight_score, weight_pagerank := 2.0, 0.5 // Creating weighing values based on the score and pagerank

		for _, hit := range searchResults.Hits { // Modify results and re-sort
			results = append(results, Result{Id: hit.ID, Score: hit.Score * weight_score + hit.Fields["Score"].(float64) * weight_pagerank, Content: hit.Fields["Content"].(string)})
		}
		sort.SliceStable(results, func(i, j int) bool {
			return results[i].Score > results[j].Score
		})
		
		c.JSON(http.StatusOK, results) // Send the modified results
   })

   r.Run(":8080")
}