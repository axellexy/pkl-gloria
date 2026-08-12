package handler

import (
	"database/sql"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Jurnal struct {
	ID         int    `json:"id"`
	Tanggal    string `json:"tanggal"`
	Judul      string `json:"judul"`
	Kegiatan   string `json:"kegiatan"`
	Keterangan string `json:"keterangan"`
	Nama       string `json:"nama"`
	CreatedAt  string `json:"created_at,omitempty"`
}

var (
	db     *sql.DB
	router *gin.Engine
	once   sync.Once
)

func initApp() {
	once.Do(func() {
		connectionString := "postgresql://postgres.rkarsceacmamrzwryczj:Manasayatau23@aws-0-ap-southeast-1.pooler.supabase.com:5432/postgres"

		var err error
		db, err = sql.Open("postgres", connectionString)
		if err != nil {
			log.Println("Gagal konek ke database:", err)
			return
		}

		err = db.Ping()
		if err != nil {
			log.Println("Database tidak merespon:", err)
			return
		}

		gin.SetMode(gin.ReleaseMode)
		router = gin.New()
		router.Use(gin.Recovery())

		// CORS Middleware
		router.Use(func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(204)
				return
			}
			c.Next()
		})

		// GET: Hanya ambil data milik user yang mengirimkan parameter ?nama=email_user
		router.GET("/api/jurnal", func(c *gin.Context) {
			emailParam := c.Query("nama")

			var rows *sql.Rows
			var err error

			if emailParam != "" {
				query := "SELECT id, tanggal, judul, kegiatan, COALESCE(keterangan, ''), COALESCE(nama, '') FROM jurnal WHERE nama = $1 ORDER BY id DESC"
				rows, err = db.Query(query, emailParam)
			} else {
				query := "SELECT id, tanggal, judul, kegiatan, COALESCE(keterangan, ''), COALESCE(nama, '') FROM jurnal ORDER BY id DESC"
				rows, err = db.Query(query)
			}

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data"})
				return
			}
			defer rows.Close()

			listJurnal := []Jurnal{}
			for rows.Next() {
				var j Jurnal
				if err := rows.Scan(&j.ID, &j.Tanggal, &j.Judul, &j.Kegiatan, &j.Keterangan, &j.Nama); err != nil {
					continue
				}
				listJurnal = append(listJurnal, j)
			}

			c.JSON(http.StatusOK, listJurnal)
		})

		// POST: Simpan jurnal dengan identifier email user
		router.POST("/api/jurnal", func(c *gin.Context) {
			var input Jurnal

			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Format data salah"})
				return
			}

			query := "INSERT INTO jurnal (tanggal, judul, kegiatan, keterangan, nama) VALUES ($1, $2, $3, $4, $5)"
			_, err := db.Exec(query, input.Tanggal, input.Judul, input.Kegiatan, input.Keterangan, input.Nama)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan ke database"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Jurnal berhasil disimpan!"})
		})
	})
}

func Handler(w http.ResponseWriter, r *http.Request) {
	initApp()
	router.ServeHTTP(w, r)
}