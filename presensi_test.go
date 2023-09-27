package absensi

import (
	"fmt"
	"os"
	"testing"

	"github.com/aiteung/atdb"
	_ "github.com/mattn/go-sqlite3"
)

var MongoInfo = atdb.DBInfo{
	DBString: os.Getenv("MONGOSTRING"),
	DBName:   "hris",
}

var MongoConn = atdb.MongoConnect(MongoInfo)

// var pengirim = types.JID{
// 	User: "6289522910966",
// }

var Long float64 = 107.57606294114771
var Lat float64 = -6.873439789736144

// var client whatsmeow.Client

// func TestGetPresensiThisMonth(t *testing.T) {
// 	// GenerateReportCurrentMonth(MongoConn)
// 	Pesan.Sender.User = "6289522910966"
// 	go wa.Runwa(&client)
// 	Handler(&Pesan, &Msg, &client, MongoConn)
// }

// func TestSelisih(t *testing.T) {
// 	Pesan.Sender.User = "6285722697918"
// 	karyawan := getKaryawanFromPhoneNumber(MongoConn, Pesan.Sender.User)
// 	cekhadir := SelisihJamMasuk(karyawan)
// 	fmt.Println(cekhadir)
// }

// func TestDurasi(t *testing.T) {
// 	start := time.Date(2023, time.August, 22, 9, 0, 0, 0, time.UTC)
// 	end := time.Date(2023, time.August, 22, 17, 30, 0, 0, time.UTC)
// 	durasi := end.Sub(start)

// 	durasiFormatted, percentageFormatted := DurasiKerja(durasi, start, end)

// 	fmt.Println("Durasi Kerja:", durasiFormatted)
// 	fmt.Println("Persentase Kerja:", percentageFormatted)
// }

// func TestTimeStamp(t *testing.T) {
// 	objectIDStr := "64e7f243ca06a39f7e741b9d" // Contoh ObjectID MongoDB
// 	timestamp, err := GetTimestampFromObjectID(objectIDStr)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Timestamp:", timestamp)
// }

// func TestTimeStamp(t *testing.T) {
// 	// karyawan := getKaryawanFromPhoneNumber(MongoConn, Pesan.Sender.User)
// 	//waktu := GetDateSekarang()
// 	//fmt.Println(waktu)
// 	waktu := GetPresensiTodayFromPhoneNumber(MongoConn, "6281312000300")
// 	fmt.Println(waktu)
// }

func TestTimeStamp(t *testing.T) {
	// presensihariini := getPresensiTodayFromPhoneNumber(MongoConn, "6281312000300")
	// fmt.Println(presensihariini)
	test := GetMulaiPresensi()
	fmt.Println(test)
	test2 := GetTimeNow()
	fmt.Println(test2)
}

// func TestTimeStamp(t *testing.T) {
// 	CopyCollectionKaryawanToPresensiBelum(MongoConn)
// }

// func TestPersentase(t *testing.T) {
// 	start := time.Date(2023, time.August, 22, 9, 0, 0, 0, time.UTC)
// 	end := time.Date(2023, time.August, 22, 17, 30, 0, 0, time.UTC)

// 	durasiKerja := DurasiKerja(start, end)
// 	fmt.Println("Durasi Kerja:", durasiKerja)

// 	aktifjamkerja := end.Sub(start)
// 	persentase := PersentaseKerja(aktifjamkerja)
// 	fmt.Printf("Persentase Kerja: %.2f%%\n", persentase)
// }

// func TestGetPresensiThisMonth(t *testing.T) {
// 	// GenerateReportCurrentMonth(MongoConn)
// 	Pesan.Sender.User = "6289522910966"
// 	go wa.Runwa(&client)
// 	Handler(&Pesan, &Msg, &client, MongoConn)
// }
