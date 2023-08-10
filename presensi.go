package absensi

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/aiteung/atmessage"
	"github.com/aiteung/musik"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
)

const Keyword string = "adorable"

func Handler(Info *types.MessageInfo, Message *waProto.Message, whatsapp *whatsmeow.Client, mongoconn *mongo.Database) {
	if Message.LiveLocationMessage != nil {
		LiveLocationMessage(Info, Message, whatsapp, mongoconn)
	} else if Message.ButtonsResponseMessage != nil {
		ButtonMessage(Info, Message, whatsapp)
	} else {
		MultiKey(mongoconn, Info, Message, whatsapp)
	}
}

func MultiKey(mongoconn *mongo.Database, Info *types.MessageInfo, Message *waProto.Message, whatsapp *whatsmeow.Client) {
	m := musik.NormalizeString(Message.GetConversation())
	complete, match := musik.IsMatch(m, "ini", "rekap", "presen", "absen", "hrd", "sdm", "excel", "data", "bulan")
	fmt.Println(complete)
	if match >= 2 {
		resp, err := GenerateReportCurrentMonth(mongoconn, Info.Chat, whatsapp)
		if err != nil {
			atmessage.SendMessage("error GenerateReportCurrentMonth", Info.Chat, whatsapp)
		}
		fmt.Println(resp)

	}

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

func LiveLocationMessage(Info *types.MessageInfo, Message *waProto.Message, whatsapp *whatsmeow.Client, mongoconn *mongo.Database) {
	lokasi := GetLokasi(mongoconn, *Message.LiveLocationMessage.DegreesLongitude, *Message.LiveLocationMessage.DegreesLatitude)
	if lokasi != "" {
		hadirHandler(Info, Message, lokasi, whatsapp, mongoconn)
	} else {
		tidakhadirHandler(Info, Message, whatsapp, mongoconn)
	}

}

func tidakhadirHandler(Info *types.MessageInfo, Message *waProto.Message, whatsapp *whatsmeow.Client, mongoconn *mongo.Database) {
	lat, long := atmessage.GetLiveLoc(Message)
	nama := GetNamaFromPhoneNumber(mongoconn, Info.Sender.User)
	MessageTidakMasukKerja(nama, long, lat, Info, whatsapp)
}

func hadirHandler(Info *types.MessageInfo, Message *waProto.Message, lokasi string, whatsapp *whatsmeow.Client, mongoconn *mongo.Database) {
	// presensihariini := getPresensiTodayFromPhoneNumber(mongoconn, Info.Sender.User)
	// karyawan := getKaryawanFromPhoneNumber(mongoconn, Info.Sender.User)
	// fmt.Println(karyawan.Jam_kerja[0].Durasi)

	// currentTime := time.Now().UTC()
	// jamMasuk := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 8, 0, 0, 0, currentTime.Location())

	// if !reflect.ValueOf(presensihariini).IsZero() {
	// 	fmt.Println(presensihariini)
	// 	aktifjamkerja := currentTime.Sub(presensihariini.ID.Timestamp().UTC())
	// 	fmt.Println(aktifjamkerja)
	// 	if int(aktifjamkerja.Hours()) >= karyawan.Jam_kerja[0].Durasi {
	// 		id := InsertPresensi(Info, Message, "pulang", mongoconn)
	// 		MessagePulangKerja(karyawan, aktifjamkerja, id, lokasi, Info, whatsapp)
	// 	} else {
	// 		MessageJamKerja(karyawan, aktifjamkerja, presensihariini, Info, whatsapp)
	// 	}
	// } else {
	// 	selisih := SelisihJamMasuk(karyawan)
	// 	if currentTime.Before(jamMasuk) {
	// 		id := InsertPresensi(Info, Message, "masuk", mongoconn)
	// 		MessageMasukKerjaCepat(karyawan, id, lokasi, selisih, Info, whatsapp)
	// 	} else if currentTime.After(jamMasuk) {
	// 		id := InsertPresensi(Info, Message, "masuk", mongoconn)
	// 		MessageTerlambatKerja(karyawan, id, lokasi, selisih, Info, whatsapp)
	// 	} else {
	// 		id := InsertPresensi(Info, Message, "masuk", mongoconn)
	// 		MessageMasukKerja(karyawan, id, lokasi, selisih, Info, whatsapp)
	// 	}
	// }
	presensihariini := getPresensiTodayFromPhoneNumber(mongoconn, Info.Sender.User)
	karyawan := getKaryawanFromPhoneNumber(mongoconn, Info.Sender.User)
	fmt.Println(karyawan.Jam_kerja[0].Durasi)
	if !reflect.ValueOf(presensihariini).IsZero() {
		fmt.Println(presensihariini)
		aktifjamkerja := time.Now().UTC().Sub(presensihariini.ID.Timestamp().UTC())
		fmt.Println(aktifjamkerja)
		waktu := GetTimeSekarang(karyawan)
		pulang := GetTimePulang(karyawan)
		selisih := SelisihJamPulang(karyawan)

		// Ganti kondisi di bawah ini
		if int(aktifjamkerja.Hours()) >= karyawan.Jam_kerja[0].Durasi || !presensihariini.ID.IsZero() {
			id := InsertPresensi(Info, Message, "pulang", mongoconn)
			if waktu <= pulang {
				MessagePulangKerjaCepat(karyawan, aktifjamkerja, id, lokasi, selisih, Info, whatsapp)
			} else if waktu >= pulang {
				MessagePulangLebihLama(karyawan, aktifjamkerja, id, lokasi, selisih, Info, whatsapp)
			} else {
				MessagePulangKerja(karyawan, aktifjamkerja, id, lokasi, Info, whatsapp)
			}
		} else {
			MessageJamKerja(karyawan, aktifjamkerja, presensihariini, Info, whatsapp)
		}
	} else {
		waktu := GetTimeSekarang(karyawan)
		masuk := GetTimeKerja(karyawan)
		if waktu <= masuk {
			id := InsertPresensi(Info, Message, "masuk", mongoconn)
			selisih := SelisihJamMasuk(karyawan)
			MessageMasukKerjaCepat(karyawan, id, lokasi, selisih, Info, whatsapp)
		} else if waktu >= masuk {
			id := InsertPresensi(Info, Message, "masuk", mongoconn)
			selisih := SelisihJamMasuk(karyawan)
			MessageTerlambatKerja(karyawan, id, lokasi, selisih, Info, whatsapp)
		} else {
			id := InsertPresensi(Info, Message, "masuk", mongoconn)
			selisih := SelisihJamMasuk(karyawan)
			MessageMasukKerja(karyawan, id, lokasi, selisih, Info, whatsapp)
		}
	}
}

func SelisihJamMasuk(karyawan Karyawan) (selisihJamFormatted string) {
	// Replace 10.00 ke 10:00
	jam := strings.Replace(karyawan.Jam_kerja[0].Jam_masuk, ".", ":", 1)
	fmt.Println("Jam Masuk :", jam)

	// Definisi lokasi waktu sekarang
	location, _ := time.LoadLocation("Asia/Jakarta")
	jamMasuk := time.Now().In(location)
	fmt.Println("Datetime Now :", jamMasuk)

	// Convert string menjadi time
	jamMasuk, _ = time.Parse("15:04", jam)
	fmt.Println("Datetime Masuk :", jamMasuk)

	// Waktu Sekarang dan Convert Waktu Sekarang menjadi format 15:04 (string)
	waktuSekarang := time.Now().In(location).Format("15:04")
	fmt.Println("Final Waktu Sekarang :", waktuSekarang)

	// Dijadikan datetime agar bisa dihitung selisih nya
	formatjam, _ := time.Parse("15:04", waktuSekarang)
	fmt.Println("Datetime Waktu Sekarang :", formatjam)

	// Hitung selisih waktu
	selisihJam := formatjam.Sub(jamMasuk).String()
	// fmt.Println("Selisih Jam Masuk :", selisihJam)

	// Ubah Hours, Minutes dan Seconds ke Jam, Menit dan Detik
	selisihJam = strings.Replace(selisihJam, "m", " menit ", 1)
	selisihJam = strings.Replace(selisihJam, "h", " jam ", 1)
	selisihJam = strings.Replace(selisihJam, "s", " detik ", 1)
	fmt.Println("Final Selisih Jam Masuk :", selisihJam)
	return selisihJam
}

func SelisihJamPulang(karyawan Karyawan) (selisihJamFormatted string) {
	// Replace 10.00 ke 10:00
	jam := strings.Replace(karyawan.Jam_kerja[0].Jam_keluar, ".", ":", 1)
	fmt.Println("Jam Pulang :", jam)

	// Definisi lokasi waktu sekarang
	location, _ := time.LoadLocation("Asia/Jakarta")
	jamKeluar := time.Now().In(location)
	fmt.Println("Datetime Now :", jamKeluar)

	// Convert string menjadi time
	jamKeluar, _ = time.Parse("15:04", jam)
	fmt.Println("Datetime Pulang :", jamKeluar)

	// Waktu Sekarang dan Convert Waktu Sekarang menjadi format 15:04 (string)
	waktuSekarang := time.Now().In(location).Format("15:04")
	fmt.Println("Final Waktu Sekarang :", waktuSekarang)

	// Dijadikan datetime agar bisa dihitung selisih nya
	formatjam, _ := time.Parse("15:04", waktuSekarang)
	fmt.Println("Datetime Waktu Sekarang :", formatjam)

	// Hitung selisih waktu
	selisihJam := formatjam.Sub(jamKeluar).String()
	// fmt.Println("Selisih Jam Masuk :", selisihJam)

	// Ubah Hours, Minutes dan Seconds ke Jam, Menit dan Detik
	selisihJam = strings.Replace(selisihJam, "m", " menit ", 1)
	selisihJam = strings.Replace(selisihJam, "h", " jam ", 1)
	selisihJam = strings.Replace(selisihJam, "s", " detik ", 1)
	fmt.Println("Final Selisih Jam Pulang :", selisihJam)
	return selisihJam
}

func GetTimeSekarang(karyawan Karyawan) (timeSekarangFormatted string) {
	// Definisi lokasi waktu sekarang
	location, _ := time.LoadLocation("Asia/Jakarta")

	// Waktu Sekarang dan Convert Waktu Sekarang menjadi format 15:04 (string)
	waktuSekarang := time.Now().In(location).Format("15:04")
	return waktuSekarang
}

func GetTimeKerja(karyawan Karyawan) (timeKerjaFormatted string) {
	jam := strings.Replace(karyawan.Jam_kerja[0].Jam_masuk, ".", ":", 1)
	return jam
}

func GetTimePulang(karyawan Karyawan) (timePulangFormatted string) {
	jam := strings.Replace(karyawan.Jam_kerja[0].Jam_keluar, ".", ":", 1)
	return jam
}

func fillStructPresensi(Info *types.MessageInfo, Message *waProto.Message, Checkin string, mongoconn *mongo.Database) (presensi Presensi) {
	presensi.Latitude, presensi.Longitude = atmessage.GetLiveLoc(Message)
	presensi.Location = GetLokasi(mongoconn, *Message.LiveLocationMessage.DegreesLongitude, *Message.LiveLocationMessage.DegreesLatitude)
	presensi.Phone_number = Info.Sender.User
	presensi.Datetime = primitive.NewDateTimeFromTime(time.Now().UTC())
	presensi.Checkin = Checkin
	presensi.Biodata = GetBiodataFromPhoneNumber(mongoconn, Info.Sender.User)
	return presensi
}

func Member(Info *types.MessageInfo, Message *waProto.Message, mongoconn *mongo.Database) (status bool) {
	if GetNamaFromPhoneNumber(mongoconn, Info.Sender.User) != "" && Info.Chat.Server != "g.us" && (Message.LiveLocationMessage != nil || Message.ButtonsResponseMessage != nil) {
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
