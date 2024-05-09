package main

import (
	"twittueur_api/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	// Initialiser echo
	e := echo.New()
	/* On doit donner le type de requete (Get, Post...), leur nom,
	et leur emplacement dans le programme.
	Une requete GET permet de demander des données à un serveur
	Une requete POST permet d'envyer des données à un serveur
	*/
	e.GET("/", routes.HelloWorld)

	// Enregistrement et connexion
	e.POST("/register", routes.Register)
	e.POST("/login", routes.Login)

	// Posts
	e.POST("/posts", routes.PostData)
	e.GET("/posts", routes.GetPosts)

	// postsLength
	e.GET("/globalPostsLength", routes.GlobalPostsLength)

	// Likes
	e.POST("/likes", routes.PostLikes)
	e.GET("/likes", routes.GetLikes)
	e.GET("/postsLikes", routes.GetLikesByPost)

	// Bookmarks
	e.POST("/bookmarks", routes.PostBookmarks)
	e.GET("/bookmarks", routes.GetBookmarks)
	e.GET("/postsBookmarks", routes.GetBookmarksByPost)

	// User
	e.GET("/user", routes.GetUser)

	// Laisser l'utilisateur accéder à db/avatars etdb/images
	e.Static("/avatars", "db/avatars")
	e.Static("/images", "db/images")

	// Lancer le serveur en local sur le port :1323
	e.Logger.Fatal(e.Start(":1323"))
}
