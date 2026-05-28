package handlers


import (
	"net/http"
	"sort"
	"stage/models"
	"strconv"
	"strings"
	"stage/utils"
	"fmt"

)



func ArtistsHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("query")

	filteredArtist, err := HandleSearch(w, r)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	if r.Method != http.MethodGet {
		utils.RenderError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	filter := r.URL.Query().Get("sort")
	filter = strings.TrimSpace(filter)

	page := r.URL.Query().Get("page")
	pageNum, _ := strconv.Atoi(page)
	if pageNum < 1{
		pageNum = 1
	}

	limit := 16
	totalPageNum := len(filteredArtist) / limit
	if len(filteredArtist) % limit != 0{
		totalPageNum = totalPageNum + 1
	}
	start := (pageNum - 1) * limit
	if page == ""{
		start = 0
	}
	end := limit + start

	if end > len(filteredArtist){
		end = len(filteredArtist)
	}

	pageSlice := filteredArtist[start:end]

	pageNumbers := make([]int, totalPageNum)

	for i := 0; i < totalPageNum; i++{
		pageNumbers[i] = i + 1
	}

	

	switch filter {
	case "ascending":
		sort.Slice(pageSlice, func(i, j int) bool {
			return pageSlice[i].Name < pageSlice[j].Name
		})
	case "descending":
		sort.Slice(pageSlice, func(i, j int) bool {
			return pageSlice[i].Name > pageSlice[j].Name
		})
	case "oldest":
		sort.Slice(pageSlice, func(i, j int) bool {
			return pageSlice[i].CreationDate < pageSlice[j].CreationDate
		})
	case "newest":
		sort.Slice(pageSlice, func(i, j int) bool {
			return pageSlice[i].CreationDate > pageSlice[j].CreationDate
		})
	case "default":
		// No sorting, keep the original order
		 filteredArtist=pageSlice
	}
    
   

	data := models.ArtistsPageData{
		Artists : pageSlice,
		PageNumbers: pageNumbers,
		PageNo:pageNum,
		Query:query,
		Sort: filter,
	}

	err = utils.Tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		utils.RenderError(w, http.StatusInternalServerError, "Template Error")
		return
	}
}

