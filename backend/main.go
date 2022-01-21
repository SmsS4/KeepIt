package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/SmsS4/KeepIt/backend/cache_api"
	"github.com/SmsS4/KeepIt/cache/utils"
	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/jwt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type GatewayConfig struct {
	Port               string
	RateLimitPerMinute int
}

type Config struct {
	CacheApi      cache_api.CacheConfig
	GatewayConfig GatewayConfig
}

var (
	mysupersecretpassword = "unicornsAreAwesome"
)

type UserInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type NewNoteInput struct {
	Note string `json:"note" binding:"required"`
}

type UpdateNoteInput struct {
	Note_id string `json:"note_id" binding:"required"`
	Note    string `json:"note" binding:"required"`
}

func GenerateId() string {
	res := strconv.Itoa(rand.Intn(99999999-10000000) + 10000000)
	return res
}

func ParseToken(tokenStr string) string {

	//creating the token
	token, _ := jwt_lib.Parse(tokenStr, func(token *jwt_lib.Token) (interface{}, error) {
		return mysupersecretpassword, nil
	})
	claims, _ := token.Claims.(jwt_lib.MapClaims)
	username := claims["username"].(string)
	return username
}

func getGatewayConfig(data map[string]string) GatewayConfig {
	port, err := strconv.Atoi(data["ratelimit-per-minute"])
	utils.CheckError(err)
	return GatewayConfig{
		Port:               data["port"],
		RateLimitPerMinute: port,
	}
}

func getConfig(configPath string) Config {
	log.Print("Getting config")
	configFile, err := ioutil.ReadFile(configPath)
	utils.CheckError(err)
	configMap := make(map[string]map[string]string)
	utils.CheckError(yaml.Unmarshal(configFile, &configMap))
	return Config{
		CacheApi:      cache_api.GetCacheConfig(configMap["cache"]),
		GatewayConfig: getGatewayConfig(configMap["gateway"]),
	}
}

var rateLimit = make(map[string][]int64)
var config = getConfig(os.Args[1])

func relaxRatelimit(requests []int64) {
	timestamp := time.Now().Unix()
	for len(requests) > 0 {
		if (timestamp - requests[0]) > 60 {
			requests = requests[1:]
		}
	}
}

func checkRatelimit(username string) bool {
	if _, ok := rateLimit[username]; !ok {
		rateLimit[username] = make([]int64, 0)
	}
	relaxRatelimit(rateLimit[username])
	if len(rateLimit[username]) > config.GatewayConfig.RateLimitPerMinute {
		return false
	}
	rateLimit[username] = append(rateLimit[username], time.Now().Unix())
	return true
}

func main() {
	kash := cache_api.CreateApi(config.CacheApi)
	r := gin.Default()
	public := r.Group("/pub")

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

		c.JSON(200, gin.H{"message": "user registered succssessfully"})

	})

	public.GET("/login", func(c *gin.Context) {
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
			c.JSON(500, gin.H{"message": "Could not generate token"})
		}
		c.JSON(200, gin.H{"token": tokenString})
	})

	private := r.Group("/private")
	private.Use(jwt.Auth(mysupersecretpassword))

	private.POST("/new_note", func(c *gin.Context) {
		var input NewNoteInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.String(400, "please register and login first")
			c.Abort()
			return
		}

		if len(input.Note) == 0 || len(input.Note) > 2048 {
			c.JSON(400, gin.H{"error": "note length error"})
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		username := ParseToken(token)
		if !checkRatelimit(username) {
			c.JSON(429, gin.H{"error": "ratelimit exceeded"})
			return
		}
		note_id := username + "$" + GenerateId()
		kash.Put(note_id, input.Note)
		c.JSON(500, gin.H{"note_id": note_id})
	})

	private.PUT("/update_note", func(c *gin.Context) {
		var input UpdateNoteInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.String(400, "please register and login first")
			c.Abort()
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		username := ParseToken(token)
		if !checkRatelimit(username) {
			c.JSON(429, gin.H{"error": "ratelimit exceeded"})
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
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.String(400, "please register and login first")
			c.Abort()
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		username := ParseToken(token)
		if !checkRatelimit(username) {
			c.JSON(429, gin.H{"error": "ratelimit exceeded"})
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
		c.JSON(200, gin.H{"note": kash.Get(noteId)})
	})

	private.DELETE("/delete_note", func(c *gin.Context) {
		note_id := c.Query("note_id")
		auth := c.Request.Header.Get("Authorization")
		if auth == "" {
			c.String(400, "please register and login first")
			c.Abort()
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		username := ParseToken(token)
		if !checkRatelimit(username) {
			c.JSON(429, gin.H{"error": "ratelimit exceeded"})
			return
		}
		if username != strings.Split(note_id, "$")[0] && username != "admin" {
			c.JSON(400, gin.H{"error": "you are not authorized to delete note"})
			return
		}
		if kash.Get(note_id) == "" {
			c.JSON(400, gin.H{"error": "note doesn't exist"})
			return
		}
		kash.Put(note_id, "")
		c.JSON(200, gin.H{"message": "note deleted successfully"})
	})

	r.Run("localhost:" + config.GatewayConfig.Port)
}
