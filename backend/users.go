package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/nu7hatch/gouuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
	"io/ioutil"
	"net/http"
	"time"
)

func viewProfile(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	sd := sessionInfo(req)
	var user User
	user.Username = ps.ByName("username")
	if user.Username != sd.Username {
		// Get user to view
		ctx := appengine.NewContext(req)
		key := datastore.NewKey(ctx, "Users", user.Username, 0, nil)
		err := datastore.Get(ctx, key, &user)
		if err != nil {
			panic(err)
		}
	} else {
		user = sd.User
		user.IsMe = true
	}
	user.OwnerStories = getOwnerStories(user, req)
	sd.ViewingUser = user
	tpl.ExecuteTemplate(res, "profile.html", &sd)
}

func login(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	tpl.ExecuteTemplate(res, "login.html", nil)
}

func signup(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	tpl.ExecuteTemplate(res, "signup.html", nil)
}

func editProfile(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	sd := sessionInfo (req)
	tpl.ExecuteTemplate(res, "editProfile.html", &sd)
}


func checkUserName(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	bs, err := ioutil.ReadAll(req.Body)
	sbs := string(bs)
	log.Infof(ctx, "REQUEST BODY: %v", sbs)
	var user User
	key := datastore.NewKey(ctx, "Users", sbs, 0, nil)
	err = datastore.Get(ctx, key, &user)
	// if there is an err, there is NO user
	log.Infof(ctx, "ERR: %v", err)
	if err != nil {
		// there is an err, there is NO user
		fmt.Fprint(res, "false")
		return
	}
	
	//there is a user
	fmt.Fprint(res, "true")
}

func checkEmail(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	bs, err := ioutil.ReadAll(req.Body)
	sbs := string(bs)
	q := datastore.NewQuery("Users").Filter("Email =", sbs)
	var u []User
	_, err = q.GetAll(ctx, &u)
	if err != nil {
		panic(err)
	}
	if len(u) > 0 {
		// users, there is a user with that email
		fmt.Fprint(res, "true")
	} else {
		// there is NO user with that email
		fmt.Fprint(res, "false")
		return
	}
}

func createUser(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.FormValue("password1")), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf(ctx, "error creating password: %v", err)
		http.Error(res, err.Error(), 500)
		return
	}
	
	t := time.Now()
	y, m, d := t.Date()
	s := fmt.Sprintf ("%v %v, %v", m, d, y)
	
	user := User{
		Email: req.FormValue("email"),
		Username: req.FormValue("username"),
		About: req.FormValue("about"),
		Image: req.FormValue("image"),
		Password: string(hashedPass),
		JoinDate: s,
	}
	key := datastore.NewKey(ctx, "Users", user.Username, 0, nil)
	key, err = datastore.Put(ctx, key, &user)
	if err != nil {
		log.Errorf(ctx, "error adding todo: %v", err)
		http.Error(res, err.Error(), 500)
		return
	}

	createSession(res, req, user)
	// redirect
	http.Redirect(res, req, "/", 302)
}

func loginProcess(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	key := datastore.NewKey(ctx, "Users", req.FormValue("username"), 0, nil)
	var user User
	err := datastore.Get(ctx, key, &user)
	if err != nil || 
		bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.FormValue("password"))) != nil {
		// failure logging in
		var sd SessionData
		sd.LoginFail = true
		tpl.ExecuteTemplate(res, "login.html", sd)
		return
	}
	
	// success logging in
	user.Username = req.FormValue("username")
	user.OwnerStories = getOwnerStories(user, req)
	createSession(res, req, user)
	// redirect
	http.Redirect(res, req, "/", 302)
	
}

func createSession(res http.ResponseWriter, req *http.Request, user User) {
	ctx := appengine.NewContext(req)
	// SET COOKIE
	id, _ := uuid.NewV4()
	cookie := &http.Cookie{
		Name:  "session",
		Value: id.String(),
		Path:  "/",
		// twenty minute session:
		// MaxAge: 60 * 20,
		//		UNCOMMENT WHEN DEPLOYED:
		//		Secure: true,
		//		HttpOnly: true,
	}
	http.SetCookie(res, cookie)

	// SET MEMCACHE session data (sd)
	json, err := json.Marshal(user)
	if err != nil {
		log.Errorf(ctx, "error marshalling during user creation: %v", err)
		http.Error(res, err.Error(), 500)
		return
	}
	
	sd := memcache.Item{
		Key:   id.String(),
		Value: json,
	}
	memcache.Set(ctx, &sd)
}

func logout(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)

	cookie, err := req.Cookie("session")
	// cookie is not set
	if err != nil {
		http.Redirect(res, req, "/", 302)
		return
	}

	// clear memcache
	sd := memcache.Item{
		Key:        cookie.Value,
		Value:      []byte(""),
		Expiration: time.Duration(1 * time.Microsecond),
	}
	memcache.Set(ctx, &sd)

	// clear the cookie
	cookie.MaxAge = -1
	http.SetCookie(res, cookie)

	// redirect
	http.Redirect(res, req, "/login", 302)
}

func editProfileProcess(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	sd := sessionInfo (req)
	
	user := User{
		Username: sd.Username,
		Email: req.FormValue("email"),
		About: req.FormValue("about"),
		Image: req.FormValue("image"),
		Password: sd.Password,
		JoinDate: sd.JoinDate,
	}
	key := datastore.NewKey(ctx, "Users", sd.Username, 0, nil)
	key, err := datastore.Put(ctx, key, &user)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
		
	createSession(res, req, user)

	// redirect to profile page
	http.Redirect(res, req, "/user/" + sd.Username, 302)
}

func editPassword(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	ctx := appengine.NewContext(req)
	sd := sessionInfo (req)
	
	key := datastore.NewKey(ctx, "Users", sd.Username, 0, nil)
	var user User
	err := datastore.Get(ctx, key, &user)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.FormValue("password"))) != nil {
		// wrong current password
		http.Redirect(res, req, "/user/" + sd.Username, 302)
		return
	}
	
	// correct current password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.FormValue("password1")), bcrypt.DefaultCost)
	if err != nil {
		log.Errorf(ctx, "error creating password: %v", err)
		http.Error(res, err.Error(), 500)
		return
	}
	
	user.Password = string(hashedPass)
	key = datastore.NewKey(ctx, "Users", sd.Username, 0, nil)
	key, err = datastore.Put(ctx, key, &user)
	if err != nil {
		http.Error(res, err.Error(), 500)
		return
	}
	
	createSession(res, req, user)
	
	// redirect to profile page
	http.Redirect(res, req, "/user/" + sd.Username, 302)
}
