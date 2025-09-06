package go_parking_cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	fmt.Println("BEFORE ALL TESTS")
	m.Run()
	fmt.Println("AFTER ALL TESTS")
}

func TestSkip(t *testing.T) {
	if runtime.GOOS == "windows" { 
		t.Skip("Contoh skip test pada OS Windows")
	}
}

// TestHelloWorld direfactor menjadi table-driven test.
func TestHelloWorld(t *testing.T) {
	strPtr := func(s string) *string { return &s }

	testCases := []struct {
		name     string
		input    *string
		expected string
	}{
		{
			name:     "should return Hello World when name is nil",
			input:    nil,
			expected: "Hello World",
		},
		{
			name:     "should return Hello Mizz when name is provided",
			input:    strPtr("Mizz"),
			expected: "Hello Mizz",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := HelloWorld(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestHandleParkir(t *testing.T) {
	testCases := []struct {
		name          string
		perintah      []string
		initialParkir map[string]time.Time
		expectSuccess bool
		expectError   bool
		checkPlate    string // Plat nomor yang akan diperiksa setelah tes
	}{
		{
			name:          "sukses memarkir kendaraan baru",
			perintah:      []string{"parkir", "B", "1234", "XYZ"},
			initialParkir: make(map[string]time.Time),
			expectSuccess: true,
			expectError:   false,
			checkPlate:    "B 1234 XYZ",
		},
		{
			name:          "gagal memarkir karena plat kosong",
			perintah:      []string{"parkir", ""},
			initialParkir: make(map[string]time.Time),
			expectSuccess: false,
			expectError:   true,
			checkPlate:    "",
		},
		{
			name:     "gagal memarkir karena plat sudah ada",
			perintah: []string{"parkir", "G", "5678", "HIJ"},
			initialParkir: map[string]time.Time{
				"G 5678 HIJ": time.Now(), // Mobil sudah diparkir
			},
			expectSuccess: false,
			expectError:   true,
			checkPlate:    "G 5678 HIJ",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Buat salinan map agar setiap sub-test tidak saling mempengaruhi
			areaParkir := make(map[string]time.Time)
			for k, v := range tc.initialParkir {
				areaParkir[k] = v
			}

			sukses, err := handleParkir(tc.perintah, areaParkir)

			assert.Equal(t, tc.expectSuccess, sukses)

			if tc.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			// Lakukan pengecekan tambahan jika perlu
			if tc.expectSuccess {
				assert.Contains(t, areaParkir, tc.checkPlate, "Plat nomor seharusnya ada di map setelah parkir")
			}
		})
	}
}

// TestHandleKeluar menggabungkan semua kasus keluar ke dalam satu table test.
func TestHandleKeluar(t *testing.T) {
	testCases := []struct {
		name          string
		perintah      []string
		initialParkir map[string]time.Time
		expectSuccess bool
		expectError   bool
		checkPlate    string
	}{
		{
			name:     "sukses mengeluarkan kendaraan yang terparkir",
			perintah: []string{"keluar", "D", "4321", "CBA"},
			initialParkir: map[string]time.Time{
				"D 4321 CBA": time.Now(),
			},
			expectSuccess: true,
			expectError:   false,
			checkPlate:    "D 4321 CBA",
		},
		{
			name:          "gagal mengeluarkan karena kendaraan tidak ditemukan",
			perintah:      []string{"keluar", "Z", "9999", "ZZZ"},
			initialParkir: make(map[string]time.Time),
			expectSuccess: false,
			expectError:   true,
			checkPlate:    "Z 9999 ZZZ",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			areaParkir := tc.initialParkir // Bisa langsung digunakan karena map akan dimodifikasi oleh fungsi
			sukses, err := handleKeluar(tc.perintah, areaParkir)

			assert.Equal(t, tc.expectSuccess, sukses)

			if tc.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.NotContains(t, areaParkir, tc.checkPlate, "Plat nomor seharusnya sudah dihapus dari map")
			}
		})
	}
}

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	_ = r.Close()
	return buf.String()
}

func TestHandleStatus_CaptureStdout(t *testing.T) {
	fixedTime, _ := time.Parse("2006-01-02 15:04:05", "2025-09-05 16:55:00")
	initialParkir := map[string]time.Time{
		"G 1227 XYZ": fixedTime,
		"Y 9999 XYZ": fixedTime,
	}

	out := captureOutput(t, func() {
		handleStatus(10, initialParkir)
	})

	fmt.Println(out)

	assert.Contains(t, out, "================= STATUS PARKIR ================")
	assert.Contains(t, out, "================ SISA SLOT PARKIR ================")
	assert.Contains(t, out, "Sisa slot: 8 dari 10")
	assert.Contains(t, out, "---------------------")

	assert.Contains(t, out, "G 1227 XYZ")
	assert.Contains(t, out, "Y 9999 XYZ")
}

func TestParkingCLI_AllBranches(t *testing.T) {
	// Urutan perintah memicu SEMUA cabang:
	// - status (awal kosong)
	// - parkir  (kurang arg -> pesan contoh)
	// - parkir B 1111 XYZ (sukses)
	// - parkir B 1111 XYZ (duplikat -> err != nil tercetak)
	// - keluar (kurang arg -> pesan contoh)
	// - keluar Z 9999 ZZZ (tidak ada -> err != nil tercetak)
	// - parkir C 2222 CCC (sukses)
	// - parkir D 3333 DDD (sukses, kapasitas penuh)
	// - parkir E 4444 EEE (penuh -> cabang "MAAF AREA PARKIR SUDAH PENUH")
	// - status
	// - abc (default: perintah tidak dikenali)
	// - exit (keluar loop)
	inputScript := strings.Join([]string{
		"status",
		"parkir",
		"parkir B 1111 XYZ",
		"parkir B 1111 XYZ",
		"keluar",
		"keluar Z 9999 ZZZ",
		"parkir C 2222 CCC",
		"parkir D 3333 DDD",
		"parkir E 4444 EEE",
		"status",
		"abc",
		"exit",
	}, "\n") + "\n"

	// Backup FDs
	oldIn, oldOut := os.Stdin, os.Stdout
	defer func() { os.Stdin, os.Stdout = oldIn, oldOut }()

	// Pipe untuk stdin & stdout (butuh *os.File)
	rIn, wIn, _ := os.Pipe()
	rOut, wOut, _ := os.Pipe()
	os.Stdin, os.Stdout = rIn, wOut

	// Feed input lalu tutup writer agar EOF terbaca
	go func() {
		_, _ = io.WriteString(wIn, inputScript)
		_ = wIn.Close()
	}()

	ParkingCLI(2)

	// Ambil output
	_ = wOut.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, rOut)
	_ = rIn.Close()
	_ = rOut.Close()
	out := buf.String()

	assert.Contains(t, out, "Selamat Datang di Sistem Parkir CLI")
	assert.Contains(t, out, "Kapasitas Parkir: 2 kendaraan")

	// argumen kurang
	assert.Contains(t, out, "Mohon masukkan nomor plat. Contoh: parkir G 12345 XYZ")
	assert.Contains(t, out, "Mohon masukkan nomor plat. Contoh: keluar B1234XYZ")

	// sukses parkir & keluar
	assert.Contains(t, out, "Berhasil parkir!")

	// duplikat -> err != nil diprint
	assert.Contains(t, out, "sudah terparkir")

	// keluar non-existent -> err != nil diprint
	assert.Contains(t, out, "tidak ada nomor plat Z 9999 ZZZ di tempat parkir")

	// penuh
	assert.Contains(t, out, "MAAF AREA PARKIR SUDAH PENUH")

	// default
	assert.Contains(t, out, "Perintah tidak dikenali")

	// status muncul
	assert.Contains(t, out, "================ SISA SLOT PARKIR ================")

	// exit
	assert.Contains(t, out, "Program berhenti.")
}