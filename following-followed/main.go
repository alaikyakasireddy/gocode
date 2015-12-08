package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"encoding/json"
	"github.com/dustin/go-humanize"
	// "google.golang.org/appengine/memcache"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

var tpl *template.Template

func init() {
	r := httprouter.New()
	http.Handle("/", r)
	r.GET("/", Home)
	r.GET("/user/:user", Home)
	r.GET("/profile/:name", profile)
	r.GET("/form/login", Login)
	r.GET("/form/signup", Signup)
	r.POST("/api/checkusername", checkUserName)
	r.POST("/api/createuser", createUser)
	r.POST("/api/login", loginProcess)
	r.POST("/api/tweet", tweetProcess)
	r.GET("/api/logout", logout)
	r.GET("/api/follow/:name", follow)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("public/"))))

	tpl = template.New("roottemplate")
	tpl = tpl.Funcs(template.FuncMap{
		"humanize_time": humanize.Time,
	})

	tpl = template.Must(tpl.ParseGlob("templates/html/*.html"))
}

func sessionInfo(req *http.Request) SessionData {
	// get session
	memItem, err := getSession(req)
	var sd SessionData
	if err == nil {
		// logged in
		json.Unmarshal(memItem.Value, &sd)
		sd.LoggedIn = true
	}
	return sd
}

func Home(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	//get tweets
	var err error
	var tweets []Tweet
	if len(ps) != 0 {
		user := User{UserName: ps.ByName("user")}
		tweets, err = getTweets(req, &user)
	} else {
		tweets, err = getTweets(req, nil)
	}
	if err != nil {
		// lala
	}
	sd := sessionInfo(req)
	sd.Tweets = tweets
	tpl.ExecuteTemplate(res, "home.html", &sd)
}

func stringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

func profile(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	sd := sessionInfo(req)
	var user User
	user.UserName = ps.ByName("name")
	
	ctx := appengine.NewContext(req)
	
	//Get ppl they're following
	q := datastore.NewQuery("Follow").Filter("Follower =", user.UserName)
	var f []Follow
	_, err := q.GetAll(ctx, &f)
	if err != nil {
		panic(err)
	}
	for _, val := range f {
		user.Following = append(user.Following, val.Following)
	}
	
	//Get ppl they're followed by
	q2 := datastore.NewQuery("Follow").Filter("Following =", user.UserName)
	var f2 []Follow
	_, err2 := q2.GetAll(ctx, &f2)
	if err2 != nil {
		panic(err)
	}
	for _, val := range f2 {
		user.FollowedBy = append(user.FollowedBy, val.Follower)
	}
	
	//Get ppl they're following
	q3 := datastore.NewQuery("Follow").Filter("Follower =", sd.UserName)
	var f3 []Follow
	_, err = q3.GetAll(ctx, &f3)
	if err != nil {
		panic(err)
	}
	for _, val := range f3 {
		sd.Following = append(sd.Following, val.Following)
	}
		
	sd.FollowingUser = stringInSlice(user.UserName, sd.Following)
	if (user.UserName == sd.UserName) {
		sd.FollowingUser = true
	}
	sd.ViewingUser = user
	
	tpl.ExecuteTemplate(res, "profile.html", &sd)
}

func Login(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	serveTemplate(res, req, "login.html")
}

func Signup(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	serveTemplate(res, req, "signup.html")
}

/*
TO DO:
session
-memcache templates
- uuid in a cookie
--- https while logged in? - depends upon security required
- encrypt password on datastore?
--- never store an unencrypted password, so, resoundingly, YES
--- sha-256 fast hash value
- user memcache?
- datastore / memcache
session interface change
- change login button to logout when user logged in
post tweets
follow people
see tweets for everyone
see tweets for individual user
*/
