package absensi

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/aiteung/atdb"
	_ "github.com/mattn/go-sqlite3"
)

var MongoInfo = atdb.DBInfo{
	DBString: os.Getenv("MONGOSTRING"),
	DBName:   "hris",
}

var MongoConn = atdb.MongoConnect(MongoInfo)

func TestRekapBulanan(t *testing.T) {
	bulan := time.August // Ganti dengan bulan yang sesuai
	tahun := 2023        // Ganti dengan tahun yang sesuai
	test, _ := GetDataPresensiMasukBulanan(bulan, tahun, MongoConn)
	fmt.Println(test)
}

func TestJamRamadhan(t *testing.T) {
	fmt.Println("Selisih Waktu Masuk Ramadhan : ", SelisihJamMasukRamadhan())
}

func TestJamMasuk(t *testing.T) {
	karyawan := GetKaryawanFromPhoneNumber(MongoConn, "6289522910966")
	test := GetTimeKerja(karyawan)
	fmt.Println(test)
}
