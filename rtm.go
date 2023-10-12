package absensi

import "database/sql"

func GetKaryawanByPhoneNumberRtm(db *sql.DB, phoneNumber string) TblRtm {
	// Query untuk mengambil data dari tabel tblMHS dengan kondisi WHERE Nomor Telepon
	query := "SELECT id_users, full_name, email, nomor_telepon, id_user_level, id_siap, jabatan_id FROM tbl_user WHERE nomor_telepon = ?"

	var result TblRtm

	// Eksekusi query dan ambil data
	err := db.QueryRow(query, phoneNumber).Scan(&result.IdUsers, &result.FullName, &result.Email, &result.NomorTelepon, &result.IdUserLevel, &result.IdSiap, &result.JabatanId)
	if err != nil {
		return TblRtm{}
	}

	return result
}

func GetNamaFromPhoneNumberRtm(db *sql.DB, phoneNumber string) (TblRtm, error) {
	// Query untuk mengambil data dari tabel tblMHS dengan kondisi WHERE Nomor Telepon
	query := "SELECT full_name FROM tbl_user WHERE nomor_telepon = ?"

	var result TblRtm

	// Eksekusi query dan ambil data
	err := db.QueryRow(query, phoneNumber).Scan(&result.FullName)
	if err != nil {
		return TblRtm{}, err
	}

	return result, nil
}
