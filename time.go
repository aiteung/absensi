package absensi

import (
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetFirstLastDateCurrentMonth() (firstOfMonth, lastOfMonth time.Time) {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth = time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth = firstOfMonth.AddDate(0, 1, -1)

	fmt.Println(firstOfMonth)
	fmt.Println(lastOfMonth)
	return

}

func GetTimestampFromObjectID(objectID primitive.ObjectID) time.Time {
	timestamp := objectID.Timestamp()
	return timestamp
}

func ConvertTimestampToJkt(waktu time.Time) time.Time {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	return waktu.In(loc)
}

func DurasiKerja(durasi time.Duration, start time.Time, end time.Time) (string, string) {
	//aktifjamkerja := end.Sub(start)

	loc, _ := time.LoadLocation("Asia/Jakarta")
	endInJakarta := end.In(loc)

	aktifjamkerja := endInJakarta.Sub(start)

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

func SelisihJamMasukRamadhan() (selisihJamFormatted string) {
	// Replace 10.00 ke 10:00
	jam := GetJamRamadhan()
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

func GetTimeSekarang() (timeSekarangFormatted string) {
	// Definisi lokasi waktu sekarang
	location, _ := time.LoadLocation("Asia/Jakarta")

	// Waktu Sekarang dan Convert Waktu Sekarang menjadi format 15:04 (string)
	waktuSekarang := time.Now().In(location).Format("15:04")
	return waktuSekarang
}

func GetDateSekarang() (datesekarang time.Time) {
	// Definisi lokasi waktu sekarang
	location, _ := time.LoadLocation("Asia/Jakarta")

	t := time.Now().In(location) //.Truncate(24 * time.Hour)
	datesekarang = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	return
}

func GetDateKemarin() (datekemarin time.Time) {
	// Definisi lokasi waktu sekarang
	location, _ := time.LoadLocation("Asia/Jakarta")

	t := time.Now().AddDate(0, 0, -1).In(location) //.Truncate(24 * time.Hour)
	datekemarin = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	return
}

func GetTimeKerja(karyawan Karyawan) (timeKerjaFormatted string) {
	jam := strings.Replace(karyawan.Jam_kerja[0].Jam_masuk, ".", ":", 1)
	return jam
}

func GetTimePulang(karyawan Karyawan) (timePulangFormatted string) {
	jam := strings.Replace(karyawan.Jam_kerja[0].Jam_keluar, ".", ":", 1)
	return jam
}

func GetMulaiPresensi() (waktumulai string) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	waktuSekarang := time.Now().In(loc)
	waktuMulai := time.Date(waktuSekarang.Year(), waktuSekarang.Month(), waktuSekarang.Day(), 6, 0, 0, 0, loc)
	return waktuMulai.Format("15:04")
}

func GetJamRamadhan() (waktumasuk string) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	waktuSekarang := time.Now().In(loc)
	waktuMulai := time.Date(waktuSekarang.Year(), waktuSekarang.Month(), waktuSekarang.Day(), 8, 30, 0, 0, loc)
	return waktuMulai.Format("15:04")
}

func GetJamPulangRamadhan() (waktumasuk string) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	waktuSekarang := time.Now().In(loc)
	waktuMulai := time.Date(waktuSekarang.Year(), waktuSekarang.Month(), waktuSekarang.Day(), 16, 0, 0, 0, loc)
	return waktuMulai.Format("15:04")
}
