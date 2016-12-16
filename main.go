package main

import (
	_ "github.com/go-sql-driver/mysql"
	//"database/sql"
	"fmt"
	"net/http"
	"io/ioutil"
	 "github.com/gorilla/mux"

	"net"
	"bufio"
	"log"
	"encoding/json"
	"time"
	"database/sql"

)

type User struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
}

func main() {

	router := mux.NewRouter()
	rootRouter := mux.NewRouter()
	resourceRouter := mux.NewRouter()

	router.HandleFunc("/user", registrationPage).Methods("GET")
	router.HandleFunc("/user", createPlayer).Methods("POST")
	router.HandleFunc("/user", updateProfile).Methods("UPDATE")//this won't be implemented just yet
	http.Handle("/user", router)

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

	http.Redirect(w, r, "/user", 301)
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


type Player struct {
	Username string
	Email string
	Password string
}


/**
 * This function will eventually get the info from the request and
 * persist it to the mysql database
 */
func createPlayer(response http.ResponseWriter, request *http.Request){
	//extract data and put it in the sql queries.

	//get username, email, password

	var user Player

	user = parseRequest(request)

	log.Println(user.Username)
	//fmt.Println(t.Username)

	moment := time.Now()
	current_time := moment.Format("2006-01-02T15:04:05")

	db, err := sql.Open("mysql", "master:12345678@tcp(mauza.duckdns.org:3306)/AquireGo?charset=utf8")//dsn info here.
	checkErr(err)

	fmt.Println("above is what was printed from the request")

	//check to see if it's valid
	//rows, err := db.Prepare("SELECT * FROM players WHERE username=?")
	//checkErr(err)
	//defer rows.Close()
	//
	//row, errr := rows.Exec(user.Username)
	//if err != nil{
	//	panic(errr)
	//}

	//do unique username and email check


	//if(count > 0 || row == nil){
	//	//if not exit
	//	return "Player already exists"//exit with errror code/message
	//}

	// insert
	stmt, err := db.Prepare("INSERT players SET email=?,username=?,games_played=?, password=?, created_at=?, updated_at=?")
	checkErr(err)

	//encrypt the password

	result , err := stmt.Exec(user.Username, user.Email, "0", "secret",current_time,current_time)
	checkErr(err)

	fmt.Println(result)

	response.WriteHeader(201);

//	return response;//this needs to be revised
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

	var updatePlayer Player
	updatePlayer = parseRequest(request)
	log.Println(updatePlayer.Username)

	db, err := sql.Open("mysql", "master:12345678@tcp(mauza.duckdns.org:3306)/AquireGo?charset=utf8")//dsn info here.
	checkErr(err)

	// update
	stmt, err := db.Prepare("update players set username=?, email=?, password, updated_at=? where email=?")
	checkErr(err)

	moment := time.Now()
	current_time := moment.Format("2006-01-02T15:04:05")

	res, err := stmt.Exec(updatePlayer.Email,updatePlayer.Username,updatePlayer.Password,current_time,updatePlayer.Email)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println(affect)

	//return

}

func connectToJava(response http.ResponseWriter, request *http.Request){
	//this is how we connect to the java instance
	conn, err := net.Dial("tcp", "mauza.duckdns.org:8484")
	if err != nil {
		// handle error
	}
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")//this sends the "GET / HTTP/1.0" to the java server

	status, err := bufio.NewReader(conn).ReadString('\n')
	// ...
	fmt.Println(status)

}

func parseRequest(request *http.Request) (Player){

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic("panic")
	}
	//log.Println(string(body))
	var parsePlayer Player
	err = json.Unmarshal(body, &parsePlayer)
	if err != nil {
		panic("panic")
	}

	return parsePlayer
}