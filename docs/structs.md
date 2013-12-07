Datastore: Profile
--------------------

	type User struct {
		Name			string
		Tagline			string
		Chef			bool
		RestaurantIds	[]string
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
		Comment			string
		FromID, ToID	string
		Time			time.Time[]
	}

	type Likes struct {
		FromID, ToID	string
	}