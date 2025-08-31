package go_parking_cli

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandleParkir(t *testing.T) {
	t.Run("sukses memarkir kendaraan baru", func(t *testing.T) {
		areaParkir := make(map[string]time.Time)
		perintah := []string{"parkir", "B", "1234", "XYZ"}

		sukses, err := handleParkir(perintah, areaParkir)

		assert.True(t, sukses, "Seharusnya berhasil parkir kendaraan")
		assert.Nil(t, err, "Seharusnya tidak ada error saat parkir berhasil")
		assert.Contains(t, areaParkir, "B 1234 XYZ", "Plat nomor seharusnya ada di map setelah parkir")
	})

	t.Run("gagal memarkir karena plat kosong", func(t *testing.T) {
		areaParkir := make(map[string]time.Time)
		perintah := strings.Fields("perintah ")
		perintah = append(perintah, " ")
		
		sukses, err := handleParkir(perintah, areaParkir)

		assert.False(t, sukses, "Seharusnya mengembalikan false saat gagal parkir")
		assert.NotNil(t, err, "Seharusnya mengembalikan error saat plat kosong")
	})

	t.Run("gagal memarkir karena plat sudah ada", func(t *testing.T) {
		platNomor := "G 5678 HIJ"
		areaParkir := map[string]time.Time{
			platNomor: time.Now(), // Mobil sudah diparkir sebelumnya
		}
		perintah := strings.Fields("parkir " + platNomor)

		sukses, err := handleParkir(perintah, areaParkir)

		assert.False(t, sukses, "Seharusnya mengembalikan false saat gagal parkir")
		assert.NotNil(t, err, "Seharusnya mengembalikan error saat plat sudah ada")
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

		assert.True(t, sukses, "Seharusnya mengembalikan true saat sukses keluar")
		assert.Nil(t, err, "Seharusnya tidak ada error saat mengeluarkan kendaraan")
		assert.NotContains(t, areaParkir, platNomor, "Plat nomor seharusnya sudah dihapus dari map")
	})

	t.Run("gagal mengeluarkan karena kendaraan tidak ditemukan", func(t *testing.T) {
		areaParkir := make(map[string]time.Time)
		perintah := []string{"keluar", "Z", "9999", "ZZZ"}

		sukses, err := handleKeluar(perintah, areaParkir)

		assert.False(t, sukses, "Seharusnya mengembalikan false saat gagal keluar")
		assert.NotNil(t, err, "Seharusnya mengembalikan error saat kendaraan tidak ditemukan")
	})
}