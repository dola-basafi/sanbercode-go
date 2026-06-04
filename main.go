package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Bioskop struct {
	ID     int     `json:"id"`
	Nama   string  `json:"nama"`
	Lokasi string  `json:"lokasi"`
	Rating float64 `json:"rating"`
}

var db *sql.DB

func connectDB() {
	connStr := "host=localhost port=5432 user=dola password=123 dbname=sanbercode sslmode=disable"

	var err error

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	println("Database connected")
}

func createBioskop(c *gin.Context) {
	var bioskop Bioskop

	if err := c.ShouldBindJSON(&bioskop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Request body tidak valid",
		})
		return
	}

	if bioskop.Nama == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Nama tidak boleh kosong",
		})
		return
	}

	if bioskop.Lokasi == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Lokasi tidak boleh kosong",
		})
		return
	}

	query := `
		INSERT INTO bioskop (nama, lokasi, rating)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := db.QueryRow(
		query,
		bioskop.Nama,
		bioskop.Lokasi,
		bioskop.Rating,
	).Scan(&bioskop.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Bioskop berhasil ditambahkan",
		"data":    bioskop,
	})
}

func main() {
	connectDB()
	defer db.Close()

	router := gin.Default()

	router.POST("/bioskop", createBioskop)

	router.Run(":8080")
}