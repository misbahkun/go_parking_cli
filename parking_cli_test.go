package go_parking_cli

import (
	"strings"
	"testing"
	"time"
)

func TestHandleParkir(t *testing.T) {
	t.Run("sukses memarkir kendaraan baru", func(t *testing.T) {
		areaParkir := make(map[string]time.Time)
		perintah := []string{"parkir", "B", "1234", "XYZ"}
		platNomor := "B 1234 XYZ"

		sukses, err := handleParkir(perintah, areaParkir)

		if err != nil {
			t.Errorf("Seharusnya tidak ada error, tapi malah dapat: %v", err)
		}
		if !sukses {
			t.Error("Seharusnya mengembalikan true saat berhasil parkir")
		}
		if _, ada := areaParkir[platNomor]; !ada {
			t.Errorf("Kendaraan dengan plat '%s' seharusnya ada di map, tapi tidak ditemukan", platNomor)
		}
	})

	t.Run("gagal memarkir karena plat kosong", func(t *testing.T) {
		areaParkir := make(map[string]time.Time)

		platNomor := " "
		perintah := strings.Fields("parkir " + platNomor)

		sukses, err := handleParkir(perintah, areaParkir)

		if err == nil {
			t.Error("Seharusnya mengembalikan error saat plat kosong, tapi malah nill")
		}

		if sukses {
			t.Error("Seharusnya mengembalikan false saat gagal parkir")
		}
	})								

	t.Run("gagal memarkir karena plat sudah ada", func(t *testing.T) {
		platNomor := "G 5678 HIJ"
		areaParkir := map[string]time.Time{
			platNomor: time.Now(), // Mobil sudah diparkir sebelumnya
		}
		perintah := strings.Fields("parkir " + platNomor) // ["parkir", "G", "5678", "HIJ"]

		sukses, err := handleParkir(perintah, areaParkir)

		if err == nil {
			t.Error("Seharusnya mengembalikan error saat plat sudah ada, tapi malah nil")
		}
		if sukses {
			t.Error("Seharusnya mengembalikan false saat gagal parkir karena plat kosong")
		}
	})
}

func TestHandleKeluar(t *testing.T) {
	t.Run("sukses mengeluarkan kendaraan yang terparkir", func(t *testing.T) {
		platNomor := "D 4321 CBA"
		areaParkir := map[string]time.Time{
			platNomor: time.Now(),
		}
		perintah := strings.Fields("keluar " + platNomor)

		sukses, err := handleKeluar(perintah, areaParkir)

		if err != nil {
			t.Errorf("Seharusnya tidak ada error, tapi malah dapat: %v", err)
		}
		if !sukses {
			t.Error("Seharusnya mengembalikan true saat sukses keluar")
		}
		if _, ada := areaParkir[platNomor]; ada {
			t.Errorf("Kendaraan dengan plat '%s' seharusnya sudah dihapus dari map", platNomor)
		}
	})

	t.Run("gagal mengeluarkan karena kendaraan tidak ditemukan", func(t *testing.T) {
		areaParkir := make(map[string]time.Time)
		perintah := []string{"keluar", "Z", "9999", "ZZZ"}

		sukses, err := handleKeluar(perintah, areaParkir)

		if err == nil {
			t.Error("Seharusnya mengembalikan error saat kendaraan tidak ditemukan, tapi malah nil")
		}
		if sukses {
			t.Error("Seharusnya mengembalikan false saat gagal keluar")
		}
	})
}
