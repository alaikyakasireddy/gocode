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

func browse(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	sd := sessionInfo (req)
	sd.Stories = getAllStories(req)
	tpl.ExecuteTemplate(res, "browse.html", &sd)
}

func newStory(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	sd := sessionInfo (req)
	tpl.ExecuteTemplate(res, "newStory.html", &sd)
}

func newStoryProcess(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	sd := sessionInfo (req)
	
	t := time.Now()
	// y, m, d := t.Date()
	time := t.Local()
	s := fmt.Sprintf ("%v", time)
	
	story := Story{
		Owner: sd.Username,
		Title: req.FormValue("story"),
		Description: req.FormValue("description"),
		Plot: req.FormValue("plot"),
		Link: strings.Replace(req.FormValue("story"), " ", "-", -1),
		CreatedDate: s,
	}
	
	userkey := datastore.NewKey(ctx, "Users", sd.Username, 0, nil)
	key := datastore.NewKey(ctx, "Stories", story.Title, 0, userkey) //owner is ancestor - eliminates need for owner-story table
	key, err := datastore.Put(ctx, key, &story)
	if err != nil {
		log.Errorf(ctx, "error adding todo: %v", err)
		http.Error(res, err.Error(), 500)
		return
	}

	// redirect
	http.Redirect(res, req, "/browse", 302)
}

func getOwnerStories(user User, req *http.Request) []Story {
	ctx := appengine.NewContext(req)
	userkey := datastore.NewKey(ctx, "Users", user.Username, 0, nil)
	q := datastore.NewQuery("Stories").Ancestor(userkey)
	var stories []Story
	_, err := q.GetAll(ctx, &stories)
	if err != nil {
		panic(err)
	}
	return stories
}

func getAllStories(req *http.Request) []Story {
	ctx := appengine.NewContext(req)
	//Get 20 most recent stories, ordered by creation (newest first)
	q := datastore.NewQuery("Stories").Order("-CreatedDate").Limit(20)
	var stories []Story
	_, err := q.GetAll(ctx, &stories)
	if err != nil {
		panic(err)
	}
	return stories
}

func viewStory(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	ctx := appengine.NewContext(req)
	sd := sessionInfo(req)
	
	link := ps.ByName("story")
	title := strings.Replace(link, "-", " ", -1) //translate link to title
	owner := ps.ByName("owner")
	
	userkey := datastore.NewKey(ctx, "Users", owner, 0, nil)
	key := datastore.NewKey(ctx, "Stories", title, 0, userkey) //owner is ancestor - eliminates need for owner-story table
	var story Story
	err := datastore.Get(ctx, key, &story)
	if err != nil {
		panic(err)
	}
	var user User
	user.Username = owner
	err = datastore.Get(ctx, userkey, &user)
	if err != nil {
		panic(err)
	}
	
	scenes := getStoryScenes(story, req)
	sd.ViewingScenes = scenes
	sd.ViewingStory = story
	sd.ViewingUser = user
	tpl.ExecuteTemplate(res, "view.html", &sd)
}
