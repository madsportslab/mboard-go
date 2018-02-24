package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/eknkc/amber"
	"github.com/gorilla/mux"
	"github.com/madsportslab/glbs"
)

const (
	OMX_PLAYER						= "omxplayer"
	OMX_HDMI							= "hdmi"
	OMX_OPTION_O  				= "-o"
	OMX_OPTION_D  				= "-d"
	OMX_DBUS_DIR  				= "/org/mpris/MediaPlayer2"
	OMX_DBUS_PLAYER				= "org.mpris.MediaPlayer2.omxplayer"
  OMX_DBUS_STOP     		= 15
)

var omx *exec.Cmd

func videoPlay(options map[string]string) error {

	if omx != nil {
		return errors.New("Video already playing")
	}

	glbs.SetNamespace("blobs")

	filename := fmt.Sprintf("/home/mboard/bin/%s", *glbs.GetPath(options["key"]))

	log.Println(filename)

	omx = exec.Command(OMX_PLAYER, OMX_OPTION_O, OMX_HDMI, filename)

	err := omx.Start()
 
	if err != nil {
		return err
	}

	err = omx.Wait()

	if err != nil {
		return err
	}

	return nil
	
} // videoPlay

func videoStop(options map[string]string) error {

	if omx == nil {
		return errors.New("No video playing")
	}

	err := thebus.Emit(OMX_DBUS_DIR, OMX_DBUS_PLAYER, OMX_DBUS_STOP)

	if err != nil {
		return err
	}

	return nil

} // videoStop

func videoHandler(w http.ResponseWriter, r *http.Request) {

  switch r.Method {
  case http.MethodGet:

		data := make(map[string]string)

		vars := mux.Vars(r)

		id := vars["id"]

		if id == "" {
			w.WriteHeader(http.StatusNotFound)
		} else {

			glbs.SetNamespace("blobs")
			data["video"] = "/" + *glbs.GetPath(id)

			compiler := amber.New()
	
			err := compiler.ParseFile("mboard-www/video.amber")
	
			if err != nil {
				
				log.Printf("[%s][Error] %s", version(), err)
				w.WriteHeader(http.StatusInternalServerError)
				return
	
			}
	
			template, err2 := compiler.Compile()
	
			if err2 != nil {
				
				log.Printf("[%s][Error] %s", version(), err2)
				w.WriteHeader(http.StatusInternalServerError)
				return
	
			}
	
			template.Execute(w, data)

		}

  case http.MethodPost:
  case http.MethodDelete:
	case http.MethodPut:
	default:
	  w.WriteHeader(http.StatusMethodNotAllowed)
	}

} // videoHandler
