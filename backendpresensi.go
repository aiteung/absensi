package absensi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetDataPresensi(db *mongo.Database) (data []Presensi, err error) {
	presensi := db.Collection("presensi")
	filter := bson.M{} // Empty filter to get all data
	cur, err := presensi.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	if err := cur.All(context.Background(), &data); err != nil {
		return nil, err
	}

	if len(data) < 1 {
		return nil, errors.New("Data Tidak Ada")
	}
	return data, nil
}

func GetDataPresensiMasukHarian(db *mongo.Database) (data []Presensi, err error) {
	presensi := db.Collection("presensi")

	// Get current date
	now := time.Now()
	startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	endOfToday := startOfToday.Add(24 * time.Hour)

	// Create filter to query data for today
	filter := bson.M{
		"datetime": bson.M{
			"$gte": startOfToday,
			"$lt":  endOfToday,
		},
	}

	// Query the database
	cur, err := presensi.Find(context.Background(), filter)
	defer cur.Close(context.Background())

	if err := cur.All(context.Background(), &data); err != nil {
		return nil, err
	}

	if len(data) < 1 {
		return nil, errors.New("Data Tidak Ada")
	}
	return data, nil
}

func GetDataPresensiPulang(db *mongo.Database) (data []Pulang, err error) {
	pulang := db.Collection("presensi_pulang")
	filter := bson.M{} // Empty filter to get all data
	cur, err := pulang.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	if err := cur.All(context.Background(), &data); err != nil {
		return nil, err
	}

	if len(data) < 1 {
		return nil, errors.New("Data Tidak Ada")
	}
	return data, nil
}

func GetDataPresensiPulangHarian(db *mongo.Database) (data []Pulang, err error) {
	presensi := db.Collection("presensi_pulang")

	// Ambil tanggal hari ini
	today := time.Now().UTC()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Buat filter berdasarkan rentang waktu hari ini
	filter := bson.M{
		"datetime": bson.M{
			"$gte": startOfDay,
			"$lt":  endOfDay,
		},
	}

	cur, err := presensi.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	if err := cur.All(context.Background(), &data); err != nil {
		return nil, err
	}

	if len(data) < 1 {
		return nil, errors.New("Data Tidak Ada")
	}
	return data, nil
}

func GetOnePresensi(Id primitive.ObjectID, db *mongo.Database) (data Presensi, err error) {
	presensi := db.Collection("presensi")
	filter := bson.M{"_id": Id}
	err = presensi.FindOne(context.TODO(), filter).Decode(&data)
	if err != nil {
		fmt.Printf("Data Tidak Ada : %v\n", err)
	}
	return data, err
}

func GetOnePresensiPulang(Id primitive.ObjectID, db *mongo.Database) (data Pulang, err error) {
	presensi := db.Collection("presensi_pulang")
	filter := bson.M{"_id": Id}
	err = presensi.FindOne(context.TODO(), filter).Decode(&data)
	if err != nil {
		fmt.Printf("Data Tidak Ada : %v\n", err)
	}
	return data, err
}

func UpdatePresensi(db *mongo.Database, Id primitive.ObjectID, update bson.M) error {
	presensi := db.Collection("presensi")

	filter := bson.M{"_id": Id}

	updateResult, err := presensi.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	if updateResult.ModifiedCount == 0 {
		return errors.New("Data Tidak Ditemukan")
	}

	return nil
}

func GetDataKaryawan(db *mongo.Database) (data []Karyawan, err error) {
	karyawan := db.Collection("karyawan")
	filter := bson.M{} // Empty filter to get all data
	cur, err := karyawan.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	if err := cur.All(context.Background(), &data); err != nil {
		return nil, err
	}

	if len(data) < 1 {
		return nil, errors.New("Data Tidak Ada")
	}
	return data, nil
}

func GetOneKaryawan(Id primitive.ObjectID, db *mongo.Database) (data Karyawan, err error) {
	karyawan := db.Collection("karyawan")
	filter := bson.M{"_id": Id}
	err = karyawan.FindOne(context.TODO(), filter).Decode(&data)
	if err != nil {
		fmt.Printf("Data Tidak Ada : %v\n", err)
	}
	return data, err
}

func GetBiodataFromId(mongoconn *mongo.Database, Id primitive.ObjectID) (staf Karyawan) {
	karyawan := mongoconn.Collection("karyawan")
	filter := bson.M{"_id": Id}
	err := karyawan.FindOne(context.TODO(), filter).Decode(&staf)
	if err != nil {
		fmt.Printf("GetBiodataFromId: %v\n", err)
	}
	return staf
}

func GetDataUser(PhoneNumber string, db *mongo.Database) (data User, err error) {
	user := db.Collection("user")
	filter := bson.M{"phone_number": PhoneNumber}
	err = user.FindOne(context.TODO(), filter).Decode(&data)
	if err != nil {
		fmt.Printf("Data Tidak Ada : %v\n", err)
	}
	return data, err
}

func InsertOneDoc(db *mongo.Database, collection string, doc interface{}) (insertedID interface{}) {
	insertResult, err := db.Collection(collection).InsertOne(context.TODO(), doc)
	if err != nil {
		fmt.Printf("InsertOneDoc: %v\n", err)
	}
	return insertResult.InsertedID
}

func InsertDokumenTidakMasuk(db *mongo.Database, Id primitive.ObjectID, keterangan string, lampiran string, TanggalMulaiTdkMasuk string, TanggalSelesaiTdkMasuk string) (InsertedID interface{}, err error) {
	var presensi Presensi
	presensi.Latitude, presensi.Longitude = 0, 0
	presensi.Location = "Tidak Ada"
	presensi.Karyawan = GetBiodataFromId(db, Id)
	presensi.Keterangan = keterangan
	presensi.Lampiran = lampiran
	presensi.TanggalMulaiTdkMasuk = TanggalMulaiTdkMasuk
	presensi.TanggalSelesaiTdkMasuk = TanggalSelesaiTdkMasuk
	return InsertOneDoc(db, "presensi", presensi), err
}

func ExportToExcel(data []Presensi, filename string) error {
	f := excelize.NewFile()
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		fmt.Println(err)
		return err
	}

	f.SetCellValue("Sheet1", "A1", "Tanggal")
	f.SetCellValue("Sheet1", "B1", "Nama")
	f.SetCellValue("Sheet1", "C1", "Jabatan")
	f.SetCellValue("Sheet1", "D1", "Jam Masuk")
	f.SetCellValue("Sheet1", "E1", "Checkin")
	f.SetCellValue("Sheet1", "F1", "Keterangan")

	for i, presensi := range data {
		rowNum := i + 2

		// waktuPresensi := presensi.ID.Timestamp().(loc)
		waktuPresensi := ConvertTimestampToJkt(GetTimestampFromObjectID(presensi.Id))

		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", rowNum), waktuPresensi.Format("2006-01-02"))
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", rowNum), presensi.Karyawan.Nama)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", rowNum), presensi.Karyawan.Jabatan)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", rowNum), waktuPresensi.Format("15:04"))
		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", rowNum), presensi.Checkin)
		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", rowNum), presensi.Keterangan)
	}

	f.SetActiveSheet(index)
	return f.SaveAs(filename)
}
