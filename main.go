package main

import (
	"database/sql"
	"net/http"
	"strconv"

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

// CREATE
func createBioskop(c *gin.Context) {
	var bioskop Bioskop

	if err := c.ShouldBindJSON(&bioskop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request body tidak valid"})
		return
	}

	if bioskop.Nama == "" || bioskop.Lokasi == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Nama dan Lokasi tidak boleh kosong"})
		return
	}

	query := `INSERT INTO bioskop (nama, lokasi, rating) VALUES ($1, $2, $3) RETURNING id`
	err := db.QueryRow(query, bioskop.Nama, bioskop.Lokasi, bioskop.Rating).Scan(&bioskop.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Bioskop berhasil ditambahkan", "data": bioskop})
}

// READ ALL
func getAllBioskop(c *gin.Context) {
	rows, err := db.Query("SELECT id, nama, lokasi, rating FROM bioskop")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer rows.Close()

	bioskops := []Bioskop{}
	for rows.Next() {
		var b Bioskop
		if err := rows.Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		bioskops = append(bioskops, b)
	}

	c.JSON(http.StatusOK, gin.H{"data": bioskops})
}

// READ BY ID
func getBioskopByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}

	var b Bioskop
	query := "SELECT id, nama, lokasi, rating FROM bioskop WHERE id = $1"
	err = db.QueryRow(query, id).Scan(&b.ID, &b.Nama, &b.Lokasi, &b.Rating)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "Bioskop tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": b})
}

// UPDATE
func updateBioskop(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}

	var bioskop Bioskop
	if err := c.ShouldBindJSON(&bioskop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Request body tidak valid"})
		return
	}

	if bioskop.Nama == "" || bioskop.Lokasi == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Nama dan Lokasi tidak boleh kosong"})
		return
	}

	query := `
		UPDATE bioskop 
		SET nama = $1, lokasi = $2, rating = $3
		WHERE id = $4
		RETURNING id
	`
	var updatedID int
	err = db.QueryRow(query, bioskop.Nama, bioskop.Lokasi, bioskop.Rating, id).Scan(&updatedID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "Bioskop tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	bioskop.ID = updatedID
	c.JSON(http.StatusOK, gin.H{"message": "Bioskop berhasil diperbarui", "data": bioskop})
}

// DELETE
func deleteBioskop(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}

	result, err := db.Exec("DELETE FROM bioskop WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Bioskop tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bioskop berhasil dihapus"})
}

func main() {
	connectDB()
	defer db.Close()

	router := gin.Default()

	router.POST("/bioskop", createBioskop)
	router.GET("/bioskop", getAllBioskop)
	router.GET("/bioskop/:id", getBioskopByID)
	router.PUT("/bioskop/:id", updateBioskop)
	router.DELETE("/bioskop/:id", deleteBioskop)

	router.Run(":8080")
}