package absensi

import (
	"context"
	"fmt"

	"github.com/aiteung/module/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetNamaFromPhoneNumber(mongoconn *mongo.Database, phone_number string) (nama string) {
	karyawan := mongoconn.Collection("karyawan")
	filter := bson.M{"phone_number": phone_number}
	var staf Karyawan
	err := karyawan.FindOne(context.TODO(), filter).Decode(&staf)
	if err != nil {
		fmt.Printf("GetNamaFromPhoneNumber: %v\n", err)
	}
	return staf.Nama
}

func GetBiodataFromPhoneNumber(mongoconn *mongo.Database, phone_number string) (staf Karyawan) {
	karyawan := mongoconn.Collection("karyawan")
	filter := bson.M{"phone_number": phone_number}
	err := karyawan.FindOne(context.TODO(), filter).Decode(&staf)
	if err != nil {
		fmt.Printf("GetBiodataFromPhoneNumber: %v\n", err)
	}
	return staf
}

func GetKaryawanFromPhoneNumber(mongoconn *mongo.Database, phone_number string) (staf Karyawan) {
	karyawan := mongoconn.Collection("karyawan")
	filter := bson.M{"phone_number": phone_number}
	err := karyawan.FindOne(context.TODO(), filter).Decode(&staf)
	if err != nil {
		fmt.Printf("getKaryawanFromPhoneNumber: %v\n", err)
	}
	return staf
}

func GetPresensiTodayFromPhoneNumber(mongoconn *mongo.Database, phone_number string) (presensi Presensi) {
	coll := mongoconn.Collection("presensi")
	today := bson.M{
		"$gte": primitive.NewObjectIDFromTimestamp(GetDateSekarang()),
	}
	filter := bson.M{"phone_number": phone_number, "_id": today}
	err := coll.FindOne(context.TODO(), filter).Decode(&presensi)
	if err != nil {
		fmt.Printf("GetPresensiTodayFromPhoneNumber: %v\n", err)
	}
	return presensi
}

func GetPresensiCurrentMonth(mongoconn *mongo.Database) (allpresensi []Presensi) {
	startdate, enddate := GetFirstLastDateCurrentMonth()
	coll := mongoconn.Collection("presensi")
	today := bson.M{
		"$gte": primitive.NewDateTimeFromTime(startdate),
		"$lte": primitive.NewDateTimeFromTime(enddate),
	}
	filter := bson.M{"datetime": today}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		fmt.Printf("getPresensiTodayFromPhoneNumber: %v\n", err)
	}
	err = cursor.All(context.TODO(), &allpresensi)
	if err != nil {
		fmt.Println(err)
	}

	return
}

func GetLokasi(mongoconn *mongo.Database, long float64, lat float64) (namalokasi string) {
	lokasicollection := mongoconn.Collection("lokasi")
	filter := bson.M{
		"batas": bson.M{
			"$geoIntersects": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{long, lat},
				},
			},
		},
	}
	var lokasi Lokasi
	err := lokasicollection.FindOne(context.TODO(), filter).Decode(&lokasi)
	if err != nil {
		fmt.Printf("GetLokasi: %v\n", err)
	}
	return lokasi.Nama

}

func getKaryawanFromPhoneNumber(mongoconn *mongo.Database, phone_number string) (staf Karyawan) {
	karyawan := mongoconn.Collection("karyawan")
	filter := bson.M{"phone_number": phone_number}
	err := karyawan.FindOne(context.TODO(), filter).Decode(&staf)
	if err != nil {
		fmt.Printf("getKaryawanFromPhoneNumber: %v\n", err)
	}
	return staf
}

func getPresensiTodayFromPhoneNumber(mongoconn *mongo.Database, phone_number string) (presensi Presensi) {
	coll := mongoconn.Collection("presensi")
	today := bson.M{
		"$gte": primitive.NewObjectIDFromTimestamp(GetDateSekarang()),
	}
	filter := bson.M{"phone_number": phone_number, "_id": today}
	err := coll.FindOne(context.TODO(), filter).Decode(&presensi)
	if err != nil {
		fmt.Printf("getPresensiTodayFromPhoneNumber: %v\n", err)
	}
	return presensi
}

func getPresensiPulangTodayFromPhoneNumber(mongoconn *mongo.Database, phone_number string) (pulang Pulang) {
	coll := mongoconn.Collection("presensi_pulang")
	today := bson.M{
		"$gte": primitive.NewObjectIDFromTimestamp(GetDateSekarang()),
	}
	filter := bson.M{"phone_number": phone_number, "_id": today}
	err := coll.FindOne(context.TODO(), filter).Decode(&pulang)
	if err != nil {
		fmt.Printf("getPresensiTodayFromPhoneNumber: %v\n", err)
	}
	return pulang
}

func InsertPresensi(Pesan model.IteungMessage, Checkin string, Keterangan string, mongoconn *mongo.Database) (InsertedID interface{}) {
	insertResult, err := mongoconn.Collection("presensi").InsertOne(context.TODO(), fillStructPresensi(Pesan, Checkin, Keterangan, mongoconn))
	if err != nil {
		fmt.Printf("InsertOneDoc: %v\n", err)
	}
	return insertResult.InsertedID
}

func InsertPresensiPulang(Pesan model.IteungMessage, Checkin string, Keterangan string, Durasi string, Persentase string, mongoconn *mongo.Database) (InsertedID interface{}) {
	insertResult, err := mongoconn.Collection("presensi_pulang").InsertOne(context.TODO(), fillStructPresensiPulang(Pesan, Checkin, Keterangan, Durasi, Persentase, mongoconn))
	if err != nil {
		fmt.Printf("InsertOneDoc: %v\n", err)
	}
	return insertResult.InsertedID
}
