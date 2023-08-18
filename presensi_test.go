package absensi

import (
	"fmt"
	"os"
	"testing"

	// "time"

	"github.com/aiteung/atdb"
	_ "github.com/mattn/go-sqlite3"

	// "go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

var MongoInfo = atdb.DBInfo{
	DBString: os.Getenv("MONGOSTRING"),
	DBName:   "hris",
}

var MongoConn = atdb.MongoConnect(MongoInfo)

// var pengirim = types.JID{
// 	User: "6289522910966",
// }

var Pesan = types.MessageInfo{
	PushName: "Testing",
}

var Long float64 = 107.57606294114771
var Lat float64 = -6.873439789736144

var LokConn = waProto.LiveLocationMessage{
	DegreesLatitude:  &Lat,
	DegreesLongitude: &Long,
}

var Msg = waProto.Message{
	LiveLocationMessage: &LokConn,
}

// var client whatsmeow.Client

// func TestGetPresensiThisMonth(t *testing.T) {
// 	// GenerateReportCurrentMonth(MongoConn)
// 	Pesan.Sender.User = "6289522910966"
// 	go wa.Runwa(&client)
// 	Handler(&Pesan, &Msg, &client, MongoConn)
// }

func TestSelisih(t *testing.T) {
	Pesan.Sender.User = "6285722697918"
	karyawan := getKaryawanFromPhoneNumber(MongoConn, Pesan.Sender.User)
	cekhadir := NewSelisihPulangCepat(karyawan)
	fmt.Println(cekhadir)
}

// func TestDurasi(t *testing.T) {
// 	start := time.Date(2023, time.August, 18, 9, 0, 0, 0, time.UTC)
// 	end := time.Date(2023, time.August, 18, 17, 30, 0, 0, time.UTC)

// 	durasiFormatted := DurasiKerja(start, end)
// 	fmt.Println("Durasi Kerja:", durasiFormatted)
// }

// func TestGetPresensiThisMonth(t *testing.T) {
// 	// GenerateReportCurrentMonth(MongoConn)
// 	Pesan.Sender.User = "6289522910966"
// 	go wa.Runwa(&client)
// 	Handler(&Pesan, &Msg, &client, MongoConn)
// }
