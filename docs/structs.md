Datastore: Profile
--------------------

type Foodie struct {
	Name	string
	Tagline	string

}

type Chef struct {
	Foodie
	Restaurant	string
}


Datastore: Restaurant
----------------------
type Restaurant struct {
	Name	string
	Address	string
}


Datastore: Interaction
----------------------
type Comments struct {
	UserFrom	ID
	UserTo		ID
	Comment 	string
}

type Likes struct {
	UserFrom	ID
	UserTo		ID
}