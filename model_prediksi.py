import pyodbc
from sklearn.linear_model import LinearRegression
import pandas as pd


conn_str = r"DRIVER={ODBC Driver 17 for SQL Server};SERVER=LAPTOP-NDII57D7\DAMSSQLSERVER;DATABASE=MachineLearningDB;UID=damar;PWD=123"

conn = pyodbc.connect(conn_str)

query_transaksi = """
SELECT  Jumlah_Klik_Aplikasi,
        Qty_Transaksi,
        Omset
FROM Fact_DataTransaksi a
WHERE a.Tipe_Data = 'AKTUAL'
"""

df_historical_data = pd.read_sql(query_transaksi, conn)

X = df_historical_data[['Jumlah_Klik_Aplikasi', 'Qty_Transaksi']]
y = df_historical_data[['Omset']]

model = LinearRegression()

model.fit(X, y)

print(f"Hasil Pelatihan model dengan R2-Score: {model.score(X, y) : .4f}")

# asumsi perkiraan di juni nantinya akan ada klik sebanyak 500x, dan pembelian 25 produk.
data_BulanJuni_Baru = [[500, 25]]

hasil_prediksi = model.predict(data_BulanJuni_Baru)
hasil_prediksi_omset_BulanJuni = float(hasil_prediksi[0][0])

print(f"Hasil prediksi omset bulan juni 2026: {hasil_prediksi_omset_BulanJuni: .2f}")

cursor = conn.cursor()

insert_predict = """
    INSERT INTO Fact_DataTransaksi
    (PelangganSK, ProdukSK, WaktuKey, Jumlah_Klik_Aplikasi, Qty_Transaksi, Omset, Tipe_Data)
    VALUES(?,?,?,?,?,?,?)
"""

fitur_data_BulanJuni = data_BulanJuni_Baru[0]
value_Inject = (1,1, 202606, fitur_data_BulanJuni[0], fitur_data_BulanJuni[1], hasil_prediksi_omset_BulanJuni, "PREDIKSI")

cursor.execute(insert_predict, value_Inject)
conn.commit()

print(f"Data prediksi juni sukses di insert ke SQL Server.")

cursor.close()
conn.close()



