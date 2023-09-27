package absensi

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/aiteung/atmessage"
	"github.com/aiteung/module/model"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

const Keyword string = "ulbi"

func Handler(Pesan model.IteungMessage, mongoconn *mongo.Database) (reply string) {
	lokasi := GetLokasi(mongoconn, Pesan.Longitude, Pesan.Latitude)
	if lokasi != "" {
		reply = hadirHandler(Pesan, lokasi, mongoconn)
	} else {
		reply = tidakhadirHandler(Pesan, mongoconn)
	}
	return
}

func ButtonMessage(Info *types.MessageInfo, Message *waProto.Message, whatsapp *whatsmeow.Client) {
	var msg string
	switch *Message.ButtonsResponseMessage.SelectedButtonId {
	case "adorable|ijin|wekwek":
		msg = "mau ijin kemana kak? yuk ingetin c obos buat set nomornya dulu biar bisa ijin ini jalan."
	case "adorable|sakit|lalala":
		msg = "Semoga lekas sembuh kak. yuk ingetin c obos buat set nomornya dulu biar bisa ijin sakit"
	case "adorable|dinas|kopkop":
		msg = "Ciee cieee yang lagi dinas diluar. yuk ingetin c obos buat set nomornya dulu biar bisa ijin sakit"
	case "adorable|lembur|wekwek":
		msg = "semangat 45 kak kejar setoran. kasih tau c obos belum set nomor nya biar bisa approve lembur"
	case "adorable|pulang|wekwek":
		msg = "hati hati di jalan ya kak, lihat lurus kedepan jangan ke lain hati kak nanti kesasar, sakit rasanya."
	default:
		msg = "Selamat datang di modul presensi, saat ini anda mengakses modul presensi."
	}

	atmessage.SendMessage(msg, Info.Sender, whatsapp)
}

func tidakhadirHandler(Pesan model.IteungMessage, mongoconn *mongo.Database) string {
	nama := GetNamaFromPhoneNumber(mongoconn, Pesan.Phone_number)
	return MessageTidakMasukKerja(nama, Pesan.Longitude, Pesan.Latitude)
}

func hadirHandler(Pesan model.IteungMessage, lokasi string, mongoconn *mongo.Database) (msg string) {
	// CODE AWAL
	presensihariini := getPresensiTodayFromPhoneNumber(mongoconn, Pesan.Phone_number)
	presensipulanghariini := getPresensiPulangTodayFromPhoneNumber(mongoconn, Pesan.Phone_number)
	durasikerja, persentasekerja := DurasiKerja(time.Now().UTC().Sub(presensihariini.Id.Timestamp()), presensihariini.Id.Timestamp(), time.Now().UTC())
	// durasikerja := DurasiKerja(time.Now().UTC().Sub(presensihariini.ID.Timestamp()), presensihariini.ID.Timestamp(), time.Now().UTC())
	karyawan := getKaryawanFromPhoneNumber(mongoconn, Pesan.Phone_number)
	waktu := GetTimeSekarang()
	pulang := GetTimePulang(karyawan)
	selisihpulangcepat := SelisihJamPulangCepat(karyawan)
	selisihpulang := SelisihJamPulang(karyawan)
	masuk := GetTimeKerja(karyawan)
	selisihmasukcepat := SelisihJamMasukCepat(karyawan)
	selisihmasuk := SelisihJamMasuk(karyawan)
	tutup := GetBatasPresensi()
	aktifjamkerja := time.Now().UTC().Sub(presensihariini.Id.Timestamp().UTC())
	mulaipresensi := GetMulaiPresensi()
	timenow := GetTimeNow()
	fmt.Println(karyawan.Jam_kerja[0].Durasi)

	if timenow.Before(mulaipresensi) {
		msg = MessageBelumBisaPresensiMasuk(karyawan)
	}

	if !reflect.ValueOf(presensihariini).IsZero() {
		fmt.Println(presensihariini)
		fmt.Println(aktifjamkerja)
		// Tambahkan kondisi ini untuk memeriksa durasi kerja sebelum mengizinkan presensi pulang
		if aktifjamkerja >= (2 * time.Hour) {
			if waktu < pulang && reflect.ValueOf(presensipulanghariini).IsZero() {
				keterangan := "Lebih Cepat"
				id := InsertPresensiPulang(Pesan, "pulang", keterangan, durasikerja, persentasekerja, mongoconn)
				msg = MessagePulangKerjaCepat(karyawan, durasikerja, persentasekerja, keterangan, id, lokasi, selisihpulangcepat)
			} else if waktu > pulang && reflect.ValueOf(presensipulanghariini).IsZero() {
				keterangan := "Lebih Lama"
				id := InsertPresensiPulang(Pesan, "pulang", keterangan, durasikerja, persentasekerja, mongoconn)
				msg = MessagePulangLebihLama(karyawan, durasikerja, persentasekerja, keterangan, id, lokasi, selisihpulang)
			} else if waktu == pulang && reflect.ValueOf(presensipulanghariini).IsZero() {
				keterangan := "Tepat Waktu"
				id := InsertPresensiPulang(Pesan, "pulang", keterangan, durasikerja, persentasekerja, mongoconn)
				msg = MessagePulangKerja(karyawan, durasikerja, persentasekerja, keterangan, id, lokasi)
			} else if !reflect.ValueOf(presensipulanghariini).IsZero() {
				msg = MessagePresensiSudahPulang(karyawan)
			} else {
				msg = MessageJamKerja(karyawan, aktifjamkerja, presensihariini)
			}
		} else {
			msg = MessageBelumBisaPresensiPulang(karyawan)
		}
	} else if waktu < masuk && waktu < tutup {
		keterangan := "Lebih Cepat"
		id := InsertPresensi(Pesan, "masuk", keterangan, mongoconn)
		msg = MessageMasukKerjaCepat(karyawan, id, lokasi, selisihmasukcepat, keterangan)
	} else if waktu > masuk && waktu < tutup {
		keterangan := "Terlambat"
		id := InsertPresensi(Pesan, "masuk", keterangan, mongoconn)
		msg = MessageTerlambatKerja(karyawan, id, lokasi, selisihmasuk, keterangan)
	} else if waktu == masuk && waktu < tutup {
		keterangan := "Tepat Waktu"
		id := InsertPresensi(Pesan, "masuk", keterangan, mongoconn)
		msg = MessageMasukKerjaTepatWaktu(karyawan, id, lokasi, keterangan)
	} else {
		msg = MessagePresensiDitutup(karyawan)
	}
	return
}

func fillStructPresensi(Pesan model.IteungMessage, Checkin string, Keterangan string, mongoconn *mongo.Database) (presensi Presensi) {
	presensi.Latitude = Pesan.Latitude
	presensi.Longitude = Pesan.Longitude
	presensi.Location = GetLokasi(mongoconn, Pesan.Longitude, Pesan.Latitude)
	presensi.Phone_number = Pesan.Phone_number
	presensi.Checkin = Checkin
	presensi.Datetime = ConvertTimestampToJkt(time.Now())
	presensi.Keterangan = Keterangan
	presensi.Karyawan = GetBiodataFromPhoneNumber(mongoconn, Pesan.Phone_number)
	return presensi
}

func fillStructPresensiPulang(Pesan model.IteungMessage, Checkin string, Keterangan string, Durasi string, Persentase string, mongoconn *mongo.Database) (pulang Pulang) {
	pulang.Latitude = Pesan.Latitude
	pulang.Longitude = Pesan.Longitude
	pulang.Location = GetLokasi(mongoconn, Pesan.Longitude, Pesan.Latitude)
	pulang.Phone_number = Pesan.Phone_number
	pulang.Checkin = Checkin
	pulang.Datetime = ConvertTimestampToJkt(time.Now())
	pulang.Keterangan = Keterangan
	pulang.Durasi = Durasi
	pulang.Persentase = Persentase
	pulang.Karyawan = GetBiodataFromPhoneNumber(mongoconn, Pesan.Phone_number)
	return pulang
}

func Member(Info *types.MessageInfo, Message *waProto.Message, mongoconn *mongo.Database, Pesan model.IteungMessage) (status bool, responseMessage string) {
	if GetNamaFromPhoneNumber(mongoconn, Info.Sender.User) != "" && Info.Chat.Server != "g.us" && (Message.LiveLocationMessage != nil || Message.ButtonsResponseMessage != nil) {
		// Jika pesan tidak berisi LiveLocationMessage
		if Message.LiveLocationMessage == nil {
			karyawan := getKaryawanFromPhoneNumber(mongoconn, Pesan.Phone_number)
			responseMessage = MessageSalahShareLoc(karyawan)
			status = false
			return status, responseMessage
		}
		if Message.ButtonsResponseMessage != nil {
			if strings.Contains(*Message.ButtonsResponseMessage.SelectedButtonId, Keyword) {
				status = true
			}
		} else {
			status = true
		}
	} else if GetNamaFromPhoneNumber(mongoconn, Info.Sender.User) != "" && strings.Contains(Message.GetConversation(), Keyword) {
		status = true
	}
	return

}
