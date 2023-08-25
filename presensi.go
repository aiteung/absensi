package absensi

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/aiteung/atmessage"
	"github.com/aiteung/module/model"
	"github.com/aiteung/musik"
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

func tidakhadirHandler(Pesan model.IteungMessage, mongoconn *mongo.Database) string {
	nama := GetNamaFromPhoneNumber(mongoconn, Pesan.Phone_number)
	return MessageTidakMasukKerja(nama, Pesan.Longitude, Pesan.Latitude)
}

func hadirHandler(Pesan model.IteungMessage, lokasi string, mongoconn *mongo.Database) (msg string) {
	// CODE AWAL
	presensihariini := getPresensiTodayFromPhoneNumber(mongoconn, Pesan.Phone_number)
	presensipulanghariini := getPresensiPulangTodayFromPhoneNumber(mongoconn, Pesan.Phone_number)
	durasikerja, persentasekerja := DurasiKerja(time.Now().UTC().Sub(presensihariini.ID.Timestamp()), presensihariini.ID.Timestamp(), time.Now().UTC())
	// durasikerja := DurasiKerja(time.Now().UTC().Sub(presensihariini.ID.Timestamp()), presensihariini.ID.Timestamp(), time.Now().UTC())
	karyawan := getKaryawanFromPhoneNumber(mongoconn, Pesan.Phone_number)
	waktu := GetTimeSekarang(karyawan)
	pulang := GetTimePulang(karyawan)
	selisihpulangcepat := SelisihJamPulangCepat(karyawan)
	selisihpulang := SelisihJamPulang(karyawan)
	masuk := GetTimeKerja(karyawan)
	selisihmasukcepat := SelisihJamMasukCepat(karyawan)
	selisihmasuk := SelisihJamMasuk(karyawan)
	fmt.Println(karyawan.Jam_kerja[0].Durasi)
	if !reflect.ValueOf(presensihariini).IsZero() {
		fmt.Println(presensihariini)
		aktifjamkerja := time.Now().UTC().Sub(presensihariini.ID.Timestamp().UTC())
		fmt.Println(aktifjamkerja)
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
	} else if waktu < masuk {
		keterangan := "Lebih Cepat"
		id := InsertPresensi(Pesan, "masuk", keterangan, mongoconn)
		msg = MessageMasukKerjaCepat(karyawan, id, lokasi, selisihmasukcepat, keterangan)
	} else if waktu > masuk {
		keterangan := "Terlambat"
		id := InsertPresensi(Pesan, "masuk", keterangan, mongoconn)
		msg = MessageTerlambatKerja(karyawan, id, lokasi, selisihmasuk, keterangan)
	} else {
		keterangan := "Tepat Waktu"
		id := InsertPresensi(Pesan, "masuk", keterangan, mongoconn)
		msg = MessageMasukKerjaTepatWaktu(karyawan, id, lokasi, keterangan)
	}
	return
}

func DurasiKerja(durasi time.Duration, start time.Time, end time.Time) (string, string) {
	aktifjamkerja := end.Sub(start)

	loc, _ := time.LoadLocation("Asia/Jakarta")
	endInJakarta := end.In(loc)

	aktifjamkerja = endInJakarta.Sub(start)

	totalMinutes := aktifjamkerja.Minutes()
	percentageOfWork := (totalMinutes / (8*60 + 30)) * 100

	hours := int(aktifjamkerja.Hours())
	minutes := int(aktifjamkerja.Minutes()) % 60
	seconds := int(aktifjamkerja.Seconds()) % 60

	durasiFormatted := fmt.Sprintf("%d Jam %d Menit %d Detik", hours, minutes, seconds)
	percentageFormatted := fmt.Sprintf("%.2f%%", percentageOfWork)
	return durasiFormatted, percentageFormatted
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

func SelisihJamMasukCepat(karyawan Karyawan) (selisihJamFormatted string) {
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
	selisihJam := jamMasuk.Sub(formatjam).String()
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

func SelisihJamPulangCepat(karyawan Karyawan) (selisihJamFormatted string) {
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
	selisihJam := jamKeluar.Sub(formatjam).String()
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
