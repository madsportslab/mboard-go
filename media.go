package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/madsportslab/glbs"
)

const (

	EXT_DOT     = "."
	EXT_JPG			= ".jpg"
	EXT_MP4			= ".mp4"
	EXT_MOV     = ".mov"
	EXT_PNG			= ".png"

)

const (

	MediaCreate = "INSERT into media(format_id, key, meta) " +
	  "VALUES($1, $2, $3)"

  MediaDelete = "DELETE from media WHERE id=?"
	
	MediaGet = "SELECT " +
	  "id, format_id, key, meta, created, updated " +
		"FROM media " +
		"WHERE id=?"

	MediaGetAll = "SELECT " +
	  "id, format_id, key, meta, created, updated " + 
		"FROM media " +
		"ORDER BY created DESC"

	FormatsGetAll = "SELECT id, name FROM formats"

)

type MediaMeta struct {
	Size		int64			`json:"size"`
	Name    string    `json:"name"`
	Ext     string    `json:"ext"`
}

type Media struct {
	ID					int   				`json:"id"`
	FormatID 		int						`json:"formatId"`
	Key					string				`json:"key"`
	Meta  			*MediaMeta		`json:"meta"`
	Created     string				`json:"created"`
	Updated     string				`json:"updated"`
}

func getFormat(ext string) int  {
	
	rows, err := data.Query(FormatsGetAll)

	if err != nil {

		log.Println("getFormat(): ", err)
		return 0

	}

	defer rows.Close()

	for rows.Next() {

		id 		:= 0
		name	:= ""

		err := rows.Scan(&id, &name)

		if err == sql.ErrNoRows || err != nil {
			
			log.Println("getFormat(): ", err)
			return 0

		}

		if strings.ToUpper(ext) == EXT_DOT + name {
			return id
		}

	}

	return 0

} // getFormat

func getMeta(filename string, size int64) (*MediaMeta, error) {

	m := MediaMeta{}

	m.Size 	= size
	m.Name 	= filename
	m.Ext		= filepath.Ext(filename)

	return &m, nil
	
} // getMeta

func createMedia(key string, filename string, size int64) {

	meta, err := getMeta(filename, size)

	j, err := json.Marshal(meta)

	if err != nil {
		log.Println(err)
	} else {

		_, err := data.Exec(
			MediaCreate, getFormat(meta.Ext), key, j,
		)
	
		if err != nil {
			log.Println(err)
		}
	
	}

} // createMedia


func removeMedia() {

} // removeMedia


func getMediaList() []Media {

	rows, err := data.Query(MediaGetAll)
	
	if err != nil {

		log.Println("getMediaList(): ", err)
		return nil

	}

	defer rows.Close()

	all := []Media{}

	for rows.Next() {

		m := Media{}

		jstr := ""

		err := rows.Scan(&m.ID, &m.FormatID, &m.Key, &jstr, &m.Created,
		  &m.Updated)

		if err == sql.ErrNoRows || err != nil {
			
			log.Println("getMediaList(): ", err)
			return nil

		}

		mm := MediaMeta{}

		errJson := json.Unmarshal([]byte(jstr), &mm)

		if errJson != nil {
			log.Println("getMediaList(): ", errJson)
		} else {

			m.Meta = &mm
	
			all = append(all, m)
	
		}


	}

	return all

} // getMediaList


func mediaHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
	case http.MethodPost:

		err := r.ParseMultipartForm(200000)

		if err != nil {
			log.Println(err)
		}

    form := r.MultipartForm

		media := form.File["media"]
		
		for _, m := range media {

			file, err := m.Open()

			defer file.Close()

			if err != nil {
				log.Println(err)
			} else {

				glbs.SetNamespace("blobs")

				k := glbs.Put(file)

				log.Printf("%s uploaded successfully, %s", m.Filename , *k)

				createMedia(*k, m.Filename, m.Size)

			}

		}

	case http.MethodGet:
		
		// search database for all media

		all := getMediaList()

		j, err := json.Marshal(all)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.Write(j)
		}
		
  case http.MethodDelete:
	case http.MethodPut:
	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // mediaHandler
