package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/microsoft/go-mssqldb"
)

type OmsetData struct {
	Labels []string  `json:"labels"`
	Data   []float64 `json:"data"`
}

func getOmset(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		query := `
			SELECT 
				b.Nama_Bulan, 
				SUM(a.Omset) AS Total_Omset
			FROM Fact_DataTransaksi a
			JOIN WAKTU b ON a.WaktuKey = b.WaktuKey
			GROUP BY b.Bulan, b.Nama_Bulan
			ORDER BY b.Bulan
		`
		rows, err := db.Query(query)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var response OmsetData

		response.Labels = []string{}
		response.Data = []float64{}

		for rows.Next() {
			var bulan string
			var omset float64

			if err := rows.Scan(&bulan, &omset); err != nil {
				log.Println(err)
				continue
			}

			response.Labels = append(response.Labels, bulan)

			response.Data = append(response.Data, omset)
		}

		json.NewEncoder(w).Encode(response)

	}
}

func main() {

	connectData := "server=localhost;port=1433;user id=damar;password=123;database=MachineLearningDB;"

	db, err := sql.Open("sqlserver", connectData)

	if err != nil {
		log.Fatal("Koneksi ke SQL Server tidak berhasil ERR: ", err)
	}

	defer db.Close()

	http.HandleFunc("/omsetPenjualan", getOmset(db))

	fmt.Println("Backend Steady, test di postman http://localhost:8080/omsetPenjualan")

	log.Fatal(http.ListenAndServe(":8080", nil))

}
