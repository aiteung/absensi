package absensi

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Karyawan struct { //data karwayan unik
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Nama         string             `bson:"nama" json:"nama"`
	Phone_number string             `bson:"phone_number" json:"phone_number"`
	Jabatan      string             `bson:"jabatan" json:"jabatan"`
	Jam_kerja    []JamKerja         `bson:"jam_kerja" json:"jam_kerja"`
	Hari_kerja   []string           `bson:"hari_kerja" json:"hari_kerja"`
}

type JamKerja struct { //info tambahan dari karyawan
	Durasi     int      `bson:"durasi,omitempty"`
	Jam_masuk  string   `bson:"jam_masuk,omitempty"`
	Jam_keluar string   `bson:"jam_keluar,omitempty"`
	Gmt        int      `bson:"gmt,omitempty"`
	Hari       []string `bson:"hari,omitempty"`
	Shift      int      `bson:"shift,omitempty"`
	Piket_tim  string   `bson:"piket_tim,omitempty"`
}

type Presensi struct {
	Id                     primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Longitude              float64            `bson:"longitude" json:"longitude"`
	Latitude               float64            `bson:"latitude" json:"latitude"`
	Location               string             `bson:"location" json:"location"`
	Phone_number           string             `bson:"phone_number" json:"phone_number"`
	Checkin                string             `bson:"checkin" json:"checkin"`
	Datetime               time.Time          `bson:"datetime,omitempty"`
	Karyawan               Karyawan           `bson:"biodata" json:"biodata"`
	Keterangan             string             `bson:"ket" json:"ket"`
	Lampiran               string             `bson:"lampiran" json:"lampiran"`
	TanggalMulaiTdkMasuk   string             `bson:"tgl_mulai_tdk_masuk" json:"tgl_mulai_tdk_masuk"`
	TanggalSelesaiTdkMasuk string             `bson:"tgl_selesai_tdk_masuk" json:"tgl_selesai_tdk_masuk"`
}

type Pulang struct { // input presensi, dimana pulang adalaha kewajiban 8 jam
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Longitude    float64            `bson:"longitude" json:"longitude"`
	Latitude     float64            `bson:"latitude" json:"latitude"`
	Location     string             `bson:"location" json:"location"`
	Phone_number string             `bson:"phone_number" json:"phone_number"`
	Checkin      string             `bson:"checkin" json:"checkin"`
	Datetime     time.Time          `bson:"datetime" json:"datetime"`
	Durasi       string             `bson:"durasi" json:"durasi"`
	Persentase   string             `bson:"persentase" json:"persentase"`
	Keterangan   string             `bson:"ket" json:"ket"`
	Karyawan     Karyawan           `bson:"biodata" json:"biodata"`
}

type RekapPresensi struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	In            Presensi           `bson:"in,omitempty"`
	Out           Presensi           `bson:"out,omitempty"`
	Lembur        Presensi           `bson:"lembur,omitempty"`
	Keterangan    string             `bson:"keterangan,omitempty"`
	TotalJamKerja primitive.DateTime `bson:"totaljamkerja,omitempty"`
	Late          primitive.DateTime `bson:"late,omitempty"`
}

type Lokasi struct { //lokasi yang bisa melakukan presensi
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Nama     string             `bson:"nama,omitempty"`
	Batas    Geometry           `bson:"batas,omitempty"`
	Kategori string             `bson:"kategori,omitempty"`
}

type Geometry struct { //data geometry untuk lokasi presensi
	Type        string      `json:"type" bson:"type"`
	Coordinates interface{} `json:"coordinates" bson:"coordinates"`
}

type TidakMasuk struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Nama       string             `bson:"nama" json:"nama"`
	Jabatan    string             `bson:"jabatan" json:"jabatan"`
	Keterangan string             `bson:"ket" json:"ket"`
	Lampiran   string             `bson:"lampiran" json:"lampiran"`
}

type User struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Nama        string             `bson:"nama" json:"nama"`
	PhoneNumber string             `bson:"phone_number" json:"phone_number"`
	Role        string             `bson:"role" json:"role"`
}
