package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/SmsS4/KeepIt/backend/cache_api"
	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/jwt"
	"github.com/gin-gonic/gin"
)

var (
	mysupersecretpassword = "unicornsAreAwesome"
)

var rateLimit = make(map[string][]int64)
var config = GetConfig(os.Args[1])

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Print("hey ")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func CheckAuth(c *gin.Context) (string, error) {
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		c.JSON(400, gin.H{"error": "please register and login first"})
		return "", gin.Error{}
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	username := ParseToken(token)
	if !CheckRatelimit(username) {
		c.JSON(429, gin.H{"error": "ratelimit exceeded"})
		return "", gin.Error{}
	}
	return username, nil
}

func main() {
	kash := cache_api.CreateApi(config.CacheApi)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	public := r.Group("/pub")

	kash.Put("admin", "admin")

	public.POST("/register", func(c *gin.Context) {
		var input UserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if len(input.Username) == 0 || len(input.Username) > 55 {
			c.JSON(400, gin.H{"error": "username length error"})
			return
		}

		if len(input.Password) == 0 {
			c.JSON(400, gin.H{"error": "password cant be empty"})
			return
		}
		if kash.Get(input.Username) != "" {
			c.JSON(400, gin.H{"error": "user already exists"})
			return
		}

		kash.Put(input.Username, input.Password)

		log.Printf("Sucessfuly added %s", input.Username)

		c.JSON(200, gin.H{"message": "user registered succssessfully"})

	})

	public.POST("/login", func(c *gin.Context) {
		var input UserInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		if kash.Get(input.Username) == "" || kash.Get(input.Username) != input.Password {
			c.JSON(404, gin.H{"error": "invalid username or password"})
			return
		}

		// Create the token
		token := jwt_lib.New(jwt_lib.GetSigningMethod("HS256"))
		// Set some claims
		token.Claims = jwt_lib.MapClaims{
			"username": input.Username,
			"exp":      time.Now().Add(time.Hour * 1).Unix(),
		}
		// Sign and get the complete encoded token as a string
		tokenString, err := token.SignedString([]byte(mysupersecretpassword))
		if err != nil {
			c.JSON(500, gin.H{"error": "Could not generate token"})
		}
		c.JSON(200, gin.H{"token": tokenString, "message": "Welcome :)"})
	})

	private := r.Group("/private")
	private.Use(jwt.Auth(mysupersecretpassword))

	private.POST("/list_note", func(c *gin.Context) {
		var input NotesListInput
		if err := c.ShouldBindJSON(&input); err != nil {
			log.Print(err.Error())
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		username, err := CheckAuth(c)
		if err != nil {
			return
		}
		if username != "admin" || len(input.Username) == 0 {
			notes := kash.Get(username + "$" + "notes")
			c.JSON(200, gin.H{"message": "Fetched", "notes": strings.Split(notes, ",")})
		} else {
			notes := kash.Get(input.Username + "$" + "notes")
			c.JSON(200, gin.H{"message": "Fetched", "notes": strings.Split(notes, ",")})
		}
	})

	private.POST("/new_note", func(c *gin.Context) {
		var input NewNoteInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		username, err := CheckAuth(c)
		if err != nil {
			return
		}
		note_id := username + "$" + GenerateId()
		notes := kash.Get(username + "$" + "notes")
		kash.Put(username+"$"+"notes", notes+","+note_id)
		kash.Put(note_id, input.Note)
		c.JSON(200, gin.H{"message": "Note added", "note_id": note_id})
	})

	private.PUT("/update_note", func(c *gin.Context) {
		var input UpdateNoteInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		username, err := CheckAuth(c)
		if err != nil {
			return
		}
		if username != strings.Split(input.Note_id, "$")[0] && username != "admin" {
			c.JSON(400, gin.H{"error": "you are not authorized to uppdate note"})
			return
		}

		if kash.Get(input.Note_id) == "" {
			c.JSON(400, gin.H{"error": "note doesnt exist"})
			return
		}

		if len(input.Note) == 0 || len(input.Note) > 2048 {
			c.JSON(400, gin.H{"error": "note length error"})
			return
		}

		kash.Put(input.Note_id, input.Note)
		c.JSON(200, gin.H{"message": "note changed succssessfully"})
	})

	private.GET("/get_note", func(c *gin.Context) {
		noteId := c.Query("note_id")
		username, err := CheckAuth(c)
		if err != nil {
			return
		}
		if username != strings.Split(noteId, "$")[0] && username != "admin" {
			c.JSON(400, gin.H{"error": "You are not authorized to read the note"})
			return
		}
		if kash.Get(noteId) == "" {
			c.JSON(400, gin.H{"error": "Note doesn't exist"})
			return
		}
		c.JSON(200, gin.H{"message": "Note fetched", "note": kash.Get(noteId)})
	})

	private.DELETE("/delete_note", func(c *gin.Context) {
		noteId := c.Query("note_id")
		username, err := CheckAuth(c)
		if err != nil {
			return
		}
		if username != strings.Split(noteId, "$")[0] && username != "admin" {
			c.JSON(400, gin.H{"error": "you are not authorized to delete note"})
			return
		}
		if kash.Get(noteId) == "" {
			c.JSON(400, gin.H{"error": "note doesn't exist"})
			return
		}
		kash.Put(noteId, "")
		notes := kash.Get(username + "$" + "notes")
		notes = strings.Replace(notes, noteId+",", "", -1)
		notes = strings.Replace(notes, ","+noteId, "", -1)
		notes = strings.Replace(notes, noteId, "", -1)
		kash.Put(username+"$"+"notes", notes)
		c.JSON(200, gin.H{"message": "note deleted successfully"})
	})

	r.Run("0.0.0.0:" + config.GatewayConfig.Port)
}
