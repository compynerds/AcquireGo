package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"net/http"
	"io/ioutil"
	"github.com/gorilla/mux"
	"net"
	//"bufio"
	"log"
	"encoding/json"
	"time"
	"crypto/sha256"
	"io"
	//"strings"
	"os"
)

var c net.Conn = connectToJava()

func main() {
	//c.Write([]byte("Connection Established"))
	router := mux.NewRouter()
	rootRouter := mux.NewRouter()
	resourceRouter := mux.NewRouter()
	gameRouter := mux.NewRouter()
	gameResponse := mux.NewRouter()

	router.HandleFunc("/user", registrationPage).Methods("GET")
	router.HandleFunc("/user", createPlayer).Methods("POST")
	router.HandleFunc("/user", updateProfile).Methods("UPDATE")
	router.HandleFunc("/user", deleteUser).Methods("DELETE")
	http.Handle("/user", router)

	rootRouter.HandleFunc("/", loginFunc)
	http.Handle("/",rootRouter)

	gameRouter.HandleFunc("/game", displayGame).Methods("GET")
	gameRouter.HandleFunc("/game", communicate).Methods("POST")
	http.Handle("/game", gameRouter)

	gameRouter.HandleFunc("/initgame", initGame).Methods("POST")
	http.Handle("/initgame", gameRouter)

	gameResponse.HandleFunc("/messages", getGameMessages).Methods("GET")
	http.Handle("/messages", gameResponse)

	resourceRouter.HandleFunc("/resources", getResources).Methods("GET")
	http.Handle("/resources",resourceRouter)

	http.ListenAndServe(":8000", nil)

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
func createPlayer(response http.ResponseWriter, request *http.Request) {
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

	//make a call to validateUnique
	exists := validateUnique(user)

	if exists != true {
		http.Redirect(response, request, "/invalidCreate", http.StatusBadRequest)
	}
	// insert
	stmt, err := db.Prepare("INSERT players SET email=?,username=?,games_played=?, password=?, created_at=?, updated_at=?")
	checkErr(err)

	//encrypt the password
	cryptPass := encryptPassword(user.Password)//this needs to be salted like crazy

	result, err := stmt.Exec(user.Username, user.Email, "0", cryptPass, current_time, current_time)
	checkErr(err)

	fmt.Println(result)

	http.Redirect(response, request, "/", http.StatusCreated)

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
func updateProfile(response http.ResponseWriter, request *http.Request){

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

	http.Redirect(response, request, "/", http.StatusAccepted)

}


func deleteUser(response http.ResponseWriter, request *http.Request){
	var deletePlayer Player

	deletePlayer = parseRequest(request)

	db, err := sql.Open("mysql", "master:12345678@tcp(mauza.duckdns.org:3306)/AquireGo?charset=utf8")//dsn info here.
	checkErr(err)

	// delete
	deleteStm, err := db.Prepare("delete from players where email=?")
	checkErr(err)

	res, err := deleteStm.Exec(deletePlayer.Email)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)
	log.Println(affect)

	db.Close()

	http.Redirect(response, request, "/", http.StatusAccepted)
}


func connectToJava()(net.Conn){
	//this is how we connect to the java instance
	conn, err := net.Dial("tcp", "localhost:8484")
	if err != nil {
		// handle error
	}
	//fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")//this sends the "GET / HTTP/1.0" to the java server

	//status, err := bufio.NewReader(conn).ReadString('\n')
	// ...
	//fmt.Println(status)
	return conn

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


func checkConnectionDB() (bool){
	//this will ping the DB to make sure it's reachable
	db, err := sql.Open("mysql", "master:12345678@tcp(mauza.duckdns.org:3306)/AquireGo?charset=utf8")//dsn info here.
	checkErr(err)
	defer db.Close()
	err = db.Ping();
	if err == nil{
		return false;
	}
	return true
}


func displayGame(response http.ResponseWriter, request *http.Request){
	path := "views" + request.URL.Path + ".html"
	fmt.Println(string(path))
	data, err := ioutil.ReadFile(string(path))
	if(err == nil){
		response.Write(data)
	} else{
		response.WriteHeader(404)
		response.Write([]byte("404 - " + http.StatusText(404)))
	}

}

func communicate(response http.ResponseWriter, request *http.Request){

	request.ParseForm()
	message := request.Form["message"][0]
	fmt.Println(message)
	//Write message to the game engine
	c.Write([]byte(message + "\n"))
	log.Println("Message sent")
	//Read back response from engine.
	m := make([]byte, 1024);
	_, error := c.Read(m);
	if error != nil {
		fmt.Printf("Cannot read: %s\n", error);
		os.Exit(1);
	}

	response.Write(m)
}


func initGame(response http.ResponseWriter, request *http.Request){
	m := make([]byte, 1024);
	_, error := c.Read(m);
	if error != nil {
		fmt.Printf("Cannot read: %s\n", error);
		os.Exit(1);
	}

	response.Write(m)
}

func loginFunc(response http.ResponseWriter, request *http.Request){
	path := "views/index.html"
	fmt.Println(string(path))
	data, err := ioutil.ReadFile(string(path))
	if(err == nil){
		response.Write(data)
	} else{
		response.WriteHeader(404)
		response.Write([]byte("404 - " + http.StatusText(404)))
	}

}


func validateUnique(player Player) (bool){
	db, err := sql.Open("mysql", "master:12345678@tcp(mauza.duckdns.org:3306)/AquireGo?charset=utf8")//dsn info here.
	checkErr(err)
	defer db.Close()
	//check to see if it's valid
	rows, err := db.Prepare("SELECT * FROM players WHERE username=?")
	checkErr(err)
	defer rows.Close()

	row, errr := rows.Exec(player.Username)
	if err != nil{
		panic(errr)
	}
	log.Println(row)
	//call the count method or something

	//do unique username and email check

	//if(count > 0 || row == nil){
	//	//if not exit
	//	return false
	//}
	return true
}


func getGameMessages(response http.ResponseWriter, request *http.Request){
	//pass the json to the endpoint for JS to grab and parse.
	fmt.Println(encryptPassword("encrypt this bitch!"))//this is for testing not production
}


func encryptPassword(passwd string) ([]byte){
	h := sha256.New()
	io.WriteString(h, "His money is twice tainted: 'taint yours and 'taint mine.")
	crypt := h.Sum(nil)
	log.Println("Encrypted pass below")
	log.Println( crypt)

	return crypt
}



