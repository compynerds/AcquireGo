package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"net/http"
	"io/ioutil"
	 "github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	rootRouter := mux.NewRouter()
	resourceRouter := mux.NewRouter()

	router.HandleFunc("/registration", registrationPage).Methods("GET")
	router.HandleFunc("/registration", createPlayer).Methods("POST")
	router.HandleFunc("/registration", updateProfile).Methods("UPDATE")//this won't be implemented just yet
	http.Handle("/registration", router)

	rootRouter.HandleFunc("/", redirect)
	http.Handle("/",rootRouter)

	resourceRouter.HandleFunc("/resources", getResources).Methods("GET")
	http.Handle("/resources",resourceRouter)

	http.ListenAndServe(":8000", nil)

}


/**
 * This will redirect the root '/' to '/registration'
 */
func redirect(w http.ResponseWriter, r *http.Request) {

	http.Redirect(w, r, "/registration", 301)
}


/**
 * This was suppose to be a CSS, JS, and Image resource link
 * This still needs to be implemented
 */
func getResources(response http.ResponseWriter, request *http.Request){

	path := request.URL.Path
	fmt.Println(string(path))
	data, err := ioutil.ReadFile(string(path))
	if(err == nil){
		response.Write(data)

	} else{
		response.WriteHeader(404)
		response.Write([]byte("404 - " + http.StatusText(404)))
	}
}


/**
 * Error function... not too special
 */
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}


/**
 * This function will eventually get the info from the request and
 * persist it to the mysql database
 */
func createPlayer(response http.ResponseWriter, request *http.Request){
	//extract data and put it in the sql queries.

	//check to see if the user and email is unique
	//

	db, err := sql.Open("mysql", "master:12345678@tcp(mauza.duckdns.org:3306)/AquireGo?charset=utf8")//dsn info here.
	checkErr(err)

	// insert
	stmt, err := db.Prepare("INSERT players SET email=?,username=?,games_played=?, password=?")
	checkErr(err)

	res, err := stmt.Exec("example@example.com", "SeymourButts", "0", "secret")
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	fmt.Println(id)//this will print the id of the record that wou
	// update
	stmt, err = db.Prepare("update players set updated_at=? where id=?")
	checkErr(err)

	res, err = stmt.Exec("2016-12-10", id) //update on players where uid = id returned from above
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	// query
	rows, err := db.Query("SELECT * FROM players")
	checkErr(err)

	for rows.Next() {
		var uid int
		var username string
		var email string
		var created_at string
		var updated_at string
		var games_played string
		var password string

		err = rows.Scan(&uid, &username, &email, &created_at, &updated_at, &games_played, &password)
		checkErr(err)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(email)
		fmt.Println(created_at)
		fmt.Println(updated_at)
		fmt.Println(games_played)
		fmt.Println(password)
	}

	// delete
	//stmt, err = db.Prepare("delete from players where id=?")
	//checkErr(err)
	//
	//res, err = stmt.Exec(id)
	//checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	db.Close()
}


/**
 * This is the struct that will create the Player objects
 */
type Player struct {
	username, email, created_at, updated_at string
	games_played int
}


/**
 * This will render the registrationPage with the '/registration' route
 */
func registrationPage(responseWriter http.ResponseWriter, request *http.Request){
	path := "views" + request.URL.Path + ".html"
	fmt.Println(string(path))
	data, err := ioutil.ReadFile(string(path))
	if(err == nil){
		responseWriter.Write(data)
	} else{
		responseWriter.WriteHeader(404)
		responseWriter.Write([]byte("404 - " + http.StatusText(404)))
	}

	//on submit either redirect to login with errors or redirect to profile page
}


/**
 * This will eventually be the function for updating the Player profile
 * this is currently on the wrong route
 */
func updateProfile(responseWriter http.ResponseWriter, request *http.Request){
	//get request data
	//update profile
	//display updated profile
}
