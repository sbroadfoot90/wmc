Public facing:
----------------------
* **/**						Home page, checks session var
* **/profile?id=$user**      	View user profile, comments, etc
* **/retaurant?name=$name**  View restaurant, current & previous
* **/edit/**					Edit logged in user profile
* **/search/?query=$query** Make a search query

Internal processes:
----------------------
* **/like?id=$user**			Submit a like
* **/comment?id=$user**		Submit a comment