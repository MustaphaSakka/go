package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type ConfigJson struct {
	Database struct {
		Host string `json:"host"`
		Port string `json:"port"`
	}
	User struct {
		Username string `json:"login"`
		Password string `json:"password"`
	}
}

type ConfigYaml struct {
	Database struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"database"`
	User struct {
		Username string `yaml:"login"`
		Password string `yaml:"password"`
	} `yaml:"user"`
}

type ConfigToml struct {
	Database struct {
		Host string
		Port string
	}
	User struct {
		Login    string
		Password string
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func eventSourcingService(client *redis.Client, subChannel, pubChannel, task string, notify bool) {
	// listen to an event
	sub := client.Subscribe(subChannel)
	msg, err := sub.ReceiveMessage()
	if err != nil {
		fmt.Println(msg)
	}
	// Do something
	fmt.Printf("[%s] > %s : %s \n", msg.Channel, task, msg.Payload)
	// notify if needed
	if notify {
		time.Sleep(2 * time.Second)
		client.Publish(pubChannel, msg.Payload)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "Hello from go server [%s]", r.URL)
}

func ArticlesHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-type", "application/json")
	json.NewEncoder(rw).Encode(articles)
}

func ArticleHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]
	//json.NewEncoder(rw).Encode(a)
	fmt.Printf("Article Num: %s", id)
}

func getArticleHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:@/api")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	// Get the 'id' parameter from the URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert 'id' to an integer
	articleID, err := strconv.Atoi(idStr)

	// Call the GetArticle function to fetch the article data from the database
	article, err := GetArticle(db, articleID)
	if err != nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	// Convert the article object to JSON and send it in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

func withLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(rw, req)
		end := time.Since(start)
		fmt.Printf("%s %s processing time %s\n", req.Method, req.URL, end)
	})
}

func main() {
	/**
	*********************************************
	API
	*********************************************
	**/
	// router := mux.NewRouter()
	// router.HandleFunc("/home", homeHandler)
	// router.HandleFunc("/articles", ArticlesHandler).Methods("GET")
	// router.HandleFunc("/article/{id:[0-9]}", getArticleHandler).Methods("GET")
	// http.Handle("/", router)

	// http.ListenAndServe(":3000", nil)
	// fmt.Println("Starting API on port 3000")

	/**
	*********************************************
	REDIS
	*********************************************
	**/
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	client.Set("monapi:/accueil", 0, time.Second).Result()
	client.Incr("monapi:/accueil").Result()
	client.Incr("monapi:/contact").Result()
	client.IncrBy("monapi:/contact", 20).Result()
	res, err := client.Get("monapi:/contact").Result()
	checkError(err)
	fmt.Println(res)

	// client.HSet("api")
	// client.HIncrBy("api", "/home", 5)
	// client.HIncrBy("api", "/help", 1)
	// client.HIncrBy("api", "/contact", 10)
	// pages, errHash := client.HGetAll("api").Result()
	// checkError(errHash)
	// fmt.Println(pages)

	// for page, NbVisit := range pages {
	// 	fmt.Println(page, NbVisit)
	// }

	// // PUB SUB
	// //client.Publish("api:notifications:security", "System under attack!").Result()
	// sub := client.PSubscribe("api:notifications:*")

	// msg, err := sub.ReceiveMessage()
	// checkError(err)
	// //PUBLISH "api:notifications:sport" "live matchs"
	// fmt.Println(msg.Channel, msg.Payload)

	//EVENT SOURCING
	go eventSourcingService(client, "user:signup", "user:confirm-mail", "Sending a confirmation mail", true)
	go eventSourcingService(client, "user:confirm-mail", "user:welcome-mail", "Sending a Welcome mail", true)
	go eventSourcingService(client, "user:welcome-mail", "user:activation", "Activating user account", true)

	fmt.Scanln()
	/**
	*********************************************
	Configuration
	*********************************************
	**/
	// JSON
	// fmt.Println("JSON config file.")
	// var configJson ConfigJson
	// file, err := os.Open("config/config.json")
	// checkError(err)
	// defer file.Close()

	// jsonDecoder := json.NewDecoder(file)
	// jsonDecoder.Decode(&configJson)

	// fmt.Println(configJson.Database.Host)

	// // YAML
	// fmt.Println("YAML config file.")
	// var configYaml ConfigYaml
	// fileYaml, errYaml := os.Open("config/config.yaml")
	// checkError(errYaml)
	// defer fileYaml.Close()

	// content, errYaml := io.ReadAll(fileYaml)
	// checkError(errYaml)

	// errYaml = yaml.Unmarshal(content, &configYaml)
	// checkError(errYaml)

	// fmt.Println(configYaml)

	// //TOML
	// fmt.Println("TOML config file.")
	// var configToml ConfigToml

	// fileToml, errToml := os.Open("Config/config.toml")
	// checkError(errToml)
	// defer fileToml.Close()

	// contentToml, errToml := io.ReadAll(fileToml)
	// checkError(errToml)

	// errToml = toml.Unmarshal(contentToml, &configToml)
	// checkError(errToml)

	// fmt.Println(configToml)
}
