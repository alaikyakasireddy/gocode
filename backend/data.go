package main

type User struct {
	Email    string
	Username string `datastore:"-"`
	Password string
	About string
	Image string
	JoinDate string
	OwnerStories []Story `datastore:"-"`
	ContributingStories []Story `datastore:"-"`
	IsMe bool `datastore:"-"`
}

type SessionData struct {
	User
	LoggedIn  bool
	LoginFail bool
	Debugging string
	ViewingUser User
	ViewingStory Story
	ViewingScenes []Scene
	Stories []Story
}

type Scene struct {
	Content string
	Name string
	CreatedDate string
	Author string
}

type Story struct {
	Title string
	Link string
	Owner string
	Description string
	Plot string
	CreatedDate string
	Scenes []Scene `datastore:"-"`
}
