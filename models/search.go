package models

import (
	"fmt"
	//"errors"
	//"fmt"
	"os"
	"strconv"

	"github.com/blevesearch/bleve"
	"github.com/jinzhu/copier"
	//"github.com/mitchellh/mapstructure"
	"github.com/will835559313/apiman/pkg/setting"
	//"github.com/yanyiwu/gojieba"
	//_ "github.com/yanyiwu/gojieba/bleve"
)

var (
	BleveIndex bleve.Index
	IndexCount = 1
)

func SearchInit() {
	// get config
	sec := setting.Cfg.Section("search")
	indexDir := sec.Key("index_dir").String()

	indexMapping := bleve.NewIndexMapping()
	os.RemoveAll(indexDir)
	// clean index when example finished
	defer os.RemoveAll(indexDir)

	//	err := indexMapping.AddCustomTokenizer("gojieba",
	//		map[string]interface{}{
	//			"dictpath":     gojieba.DICT_PATH,
	//			"hmmpath":      gojieba.HMM_PATH,
	//			"userdictpath": gojieba.USER_DICT_PATH,
	//			"idf":          gojieba.IDF_PATH,
	//			"stop_words":   gojieba.STOP_WORDS_PATH,
	//			"type":         "gojieba",
	//		},
	//	)
	//	if err != nil {
	//		panic(err)
	//	}
	//	err = indexMapping.AddCustomAnalyzer("gojieba",
	//		map[string]interface{}{
	//			"type":      "gojieba",
	//			"tokenizer": "gojieba",
	//		},
	//	)
	//	if err != nil {
	//		panic(err)
	//	}

	//	indexMapping.DefaultAnalyzer = "gojieba"

	BleveIndex, err = bleve.New(indexDir, indexMapping)
	if err != nil {
		panic(err)
	}

	IndexData()

	//ids, _ := Search("will", "user")
	//fmt.Println(ids)
}

type UserIndex struct {
	ID         string
	Name       string
	Nickname   string
	SearchType string
}

type TeamIndex struct {
	ID          string
	Name        string
	Description string
	SearchType  string
}

type ProjectIndex struct {
	ID          string
	Name        string
	Description string
	SearchType  string
}

type ApiGroupIndex struct {
	ID          string
	Name        string
	Description string
	SearchType  string
}

type ApiIndex struct {
	ID          string
	Name        string
	Description string
	SearchType  string
	URI         string
}

func IndexData() {
	// users
	users := make([]*User, 0)
	if err := db.Table("users").Select("id, name, nickname").Find(&users).Error; err == nil {
		for _, u := range users {
			d := new(UserIndex)
			copier.Copy(d, u)
			d.SearchType = "user"
			d.ID = strconv.Itoa(int(u.ID))
			//fmt.Printf("-----data %v\n", d)
			BleveIndex.Index("user:"+d.ID, d)
			IndexCount++
		}
	}

	// teams
	teams := make([]*Team, 0)
	if err := db.Table("teams").Select("id, name, description").Find(&teams).Error; err == nil {
		for _, t := range teams {
			d := new(TeamIndex)
			copier.Copy(d, t)
			d.SearchType = "team"
			d.ID = strconv.Itoa(int(t.ID))
			//fmt.Printf("-----data %v\n", d)
			BleveIndex.Index("team:"+d.ID, d)
			IndexCount++
		}
	}

	// projects
	projects := make([]*Project, 0)
	if err := db.Table("projects").Select("id, name, description").Find(&projects).Error; err == nil {
		for _, p := range projects {
			d := new(ProjectIndex)
			copier.Copy(d, p)
			d.SearchType = "project"
			d.ID = strconv.Itoa(int(p.ID))
			//fmt.Printf("-----data %v\n", d)
			BleveIndex.Index("project:"+d.ID, d)
			IndexCount++
		}
	}

	// apigroups
	api_groups := make([]*ApiGroup, 0)
	if err := db.Table("api_groups").Select("id, name, description").Find(&api_groups).Error; err == nil {
		for _, ag := range api_groups {
			d := new(ApiGroupIndex)
			copier.Copy(d, ag)
			d.SearchType = "api_group"
			d.ID = strconv.Itoa(int(ag.ID))
			//fmt.Printf("-----data %v\n", d)
			BleveIndex.Index("api_group:"+d.ID, d)
			IndexCount++
		}
	}

	// apis
	apis := make([]*Api, 0)
	if err := db.Table("apis").Select("id, name, description, uri").Find(&apis).Error; err == nil {
		for _, api := range apis {
			d := new(ApiIndex)
			copier.Copy(d, api)
			d.ID = strconv.Itoa(int(api.ID))
			d.SearchType = "api"
			//fmt.Printf("-----data %v\n", d)
			BleveIndex.Index("api:"+d.ID, d)
			IndexCount++
		}
	}
}

func Search(q, searchType, sort, order string) ([]uint, error) {
	//func Search(q, searchType string) {
	qstring := ""
	switch searchType {
	case "user":
		qstring = "+SearchType:user +" + q
	case "project":
		qstring = "+SearchType:project +" + q
	case "team":
		qstring = "+SearchType:team +" + q
	case "api":
		qstring = "+SearchType:api +" + q
	case "api_group":
		qstring = "+SearchType:api_group +" + q

	default:
		qstring = q
	}

	sq := bleve.NewQueryStringQuery(qstring)
	req := bleve.NewSearchRequest(sq)
	req.Fields = []string{"*"}

	if order == "desc" {
		req.SortBy([]string{"-_score", "-" + sort})
	} else {
		req.SortBy([]string{"-_score", sort})
	}

	searchResults, err := BleveIndex.Search(req)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}

	//fmt.Println(qstring)

	if searchResults.Total == 0 {
		fmt.Println("noting match")
		return nil, nil
	}

	ids := make([]uint, 0)
	for _, hit := range searchResults.Hits {
		if v, err := hit.Fields["ID"].(string); err {
			id, _ := strconv.Atoi(v)
			ids = append(ids, uint(id))
		}
	}

	fmt.Println(searchResults)

	return ids, nil
}
