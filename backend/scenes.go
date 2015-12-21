package main

import (
	"github.com/julienschmidt/httprouter"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"net/http"
	"strings"
	"time"
	"fmt"
)

func newScene(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	sd := sessionInfo (req)
	sd.OwnerStories = getOwnerStories(sd.User, req)
	tpl.ExecuteTemplate(res, "newScene.html", &sd)
}

func newSceneProcess(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	sd := sessionInfo (req)
	
	t := time.Now()
	time := t.Local()
	s := fmt.Sprintf ("%v", time)
	
	scene := Scene{
		Author: sd.Username,
		Name: req.FormValue("name"),
		Content: req.FormValue("scene"),
		CreatedDate: s,
	}
	
	storykey := datastore.NewKey(ctx, "Stories", req.FormValue("story"), 0, nil)
	key := datastore.NewKey(ctx, "Scenes", req.FormValue("name"), 0, storykey) //owner is ancestor - eliminates need for owner-story table
	key, err := datastore.Put(ctx, key, &scene)
	if err != nil {
		log.Errorf(ctx, "error adding todo: %v", err)
		http.Error(res, err.Error(), 500)
		return
	}

	// redirect
	link := strings.Replace(req.FormValue("story"), " ", "-", -1)
	http.Redirect(res, req, "/view/" + link + "/" + sd.Username, 302)
}

func getStoryScenes(story Story, req *http.Request) []Scene {
	ctx := appengine.NewContext(req)
	storykey := datastore.NewKey(ctx, "Stories", story.Title, 0, nil)
	q := datastore.NewQuery("Scenes").Ancestor(storykey).Order("CreatedDate")
	var scenes []Scene
	_, err := q.GetAll(ctx, &scenes)
	if err != nil {
		panic(err)
	}
	//Format date nicely for viewing
	for key, val := range scenes {
		//starting form
		//2015-12-06 00:01:45.3242642 +0000 UTC
		//reference time
		//Mon Jan 2 15:04:05 -0700 MST 2006
		longForm := "2006-01-02 15:04:05 -0700 MST"
		t, _ := time.Parse(longForm, val.CreatedDate)
		layout := "Jan 2, 2006 at 3:04pm"
		scenes[key].CreatedDate = t.Format(layout)
	}
	return scenes
}
