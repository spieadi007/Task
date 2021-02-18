package service

import (
	"auth/models"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strconv"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// db connection
func createConnection() *sql.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}

//create user
func CreateUser(c *gin.Context) {

	var user models.User

	err := c.BindJSON(&user)

	if err != nil {
		log.Fatalf("request body error.  %v", err)
	}

	insertID := insertUser(user)

	res := response{
		ID:      insertID,
		Message: "User created successfully",
	}

	c.JSON(http.StatusCreated, res)
}

//get user
func GetUser(c *gin.Context) {
	param := c.Params.ByName("id")

	id, err := strconv.Atoi(param)

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	user, err := getUser(int64(id))

	if err != nil {
		log.Fatalf("user not found. %v", err)
	}

	c.JSON(http.StatusOK, user)
}

//get all users
func GetAllUser(c *gin.Context) {

	users, err := getAllUsers()

	if err != nil {
		log.Fatalf("Users not found. %v", err)
	}

	c.JSON(http.StatusOK, users)
}

//update user
func UpdateUser(c *gin.Context) {

	param := c.Params.ByName("id")

	id, err := strconv.Atoi(param)

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	var user models.User

	errDecode := c.BindJSON(&user)

	if errDecode != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	updatedRows := updateUser(int64(id), user)

	msg := fmt.Sprintf("User updated successfully")

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	c.JSON(http.StatusOK, res)

}

// delete user
func DeleteUser(c *gin.Context) {

	param := c.Params.ByName("id")

	id, err := strconv.Atoi(param)

	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}

	deletedRows := deleteUser(int64(id))

	msg := fmt.Sprintf("User updated successfully.")

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	c.JSON(http.StatusOK, res)
}

// insert one user in the DB
func insertUser(user models.User) int64 {

	db := createConnection()

	defer db.Close()

	sqlStatement := `INSERT INTO users (name, location, age) VALUES ($1, $2, $3) RETURNING userid`

	var id int64

	err := db.QueryRow(sqlStatement, user.Name, user.Location, user.Age).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)

	return id
}

// get one user from the DB by its userid
func getUser(id int64) (models.User, error) {
	db := createConnection()

	defer db.Close()

	var user models.User

	sqlStatement := `SELECT * FROM users WHERE userid=$1`

	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&user.ID, &user.Name, &user.Age, &user.Location)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return user, nil
	case nil:
		return user, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return user, err
}

// get one user from the DB by its userid
func getAllUsers() ([]models.User, error) {

	db := createConnection()

	defer db.Close()

	var users []models.User

	sqlStatement := `SELECT * FROM users`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User

		err = rows.Scan(&user.ID, &user.Name, &user.Age, &user.Location)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		users = append(users, user)

	}

	return users, err
}

// update user in the DB
func updateUser(id int64, user models.User) int64 {

	db := createConnection()

	defer db.Close()

	sqlStatement := `UPDATE users SET name=$2, location=$3, age=$4 WHERE userid=$1`

	res, err := db.Exec(sqlStatement, id, user.Name, user.Location, user.Age)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}

// delete user in the DB
func deleteUser(id int64) int64 {

	db := createConnection()

	defer db.Close()

	sqlStatement := `DELETE FROM users WHERE userid=$1`

	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
