package main
import (
  "github.com/gorilla/mux"
  "database/sql"
  _"github.com/go-sql-driver/mysql"
  "net/http"
  "encoding/json"
  "fmt"
  "io/ioutil"
)
type Account struct { //stores data used for log-in, not public
	Pid      string `json:"pid,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}
var db *sql.DB
var err error
func main() {
db, err = sql.Open("mysql", "<user>:<password>@tcp(127.0.0.1:3306)/<dbname>")
  if err != nil {
    panic(err.Error())
  }
  defer db.Close()
  router := mux.NewRouter()
  router.HandleFunc("/posts", getProfiles).Methods("GET")
  router.HandleFunc("/posts", createProfile).Methods("POST")
  router.HandleFunc("/posts/{id}", getProfile).Methods("GET")
  router.HandleFunc("/posts/{id}", updateProfile).Methods("PUT")
  router.HandleFunc("/posts/{id}", deleteProfile).Methods("DELETE")
  http.ListenAndServe(":8000", router)
}
func getProfiles(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  var accounts []Account
  result, err := db.Query("SELECT pid, email from accounts")
  if err != nil {
    panic(err.Error())
  }
  defer result.Close()
  for result.Next() {
    var post Post
    err := result.Scan(&post.ID, &post.Title)
    if err != nil {
      panic(err.Error())
    }
    posts = append(posts, post)
  }
  json.NewEncoder(w).Encode(posts)
}
func createProfile(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  stmt, err := db.Prepare("INSERT INTO accounts(email) VALUES(?)")
  if err != nil {
    panic(err.Error())
  }
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    panic(err.Error())
  }
  keyVal := make(map[string]string)
  json.Unmarshal(body, &keyVal)
  title := keyVal["title"]
  _, err = stmt.Exec(title)
  if err != nil {
    panic(err.Error())
  }
  fmt.Fprintf(w, "New post was created")
}
func getProfile(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  params := mux.Vars(r)
  result, err := db.Query("SELECT id, title FROM posts WHERE id = ?", params["id"])
  if err != nil {
    panic(err.Error())
  }
  defer result.Close()
  var post Post
  for result.Next() {
    err := result.Scan(&post.ID, &post.Title)
    if err != nil {
      panic(err.Error())
    }
  }
  json.NewEncoder(w).Encode(post)
}
func updateProfile(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  params := mux.Vars(r)
  stmt, err := db.Prepare("UPDATE posts SET title = ? WHERE id = ?")
  if err != nil {
    panic(err.Error())
  }
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    panic(err.Error())
  }
  keyVal := make(map[string]string)
  json.Unmarshal(body, &keyVal)
  newTitle := keyVal["title"]
  _, err = stmt.Exec(newTitle, params["id"])
  if err != nil {
    panic(err.Error())
  }
  fmt.Fprintf(w, "Profile with ID = %s was updated", params["id"])
}
func deleteProfile(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  params := mux.Vars(r)
  stmt, err := db.Prepare("DELETE FROM posts WHERE id = ?")
  if err != nil {
    panic(err.Error())
  }
  _, err = stmt.Exec(params["id"])
  if err != nil {
    panic(err.Error())
  }
  fmt.Fprintf(w, "Profile with ID = %s was deleted", params["id"])
}
