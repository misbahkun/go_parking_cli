package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

const kapasitasParkir = 5

func handleStatus(areaParkir map[string]time.Time) {
	fmt.Println("================= STATUS PARKIR ================")

	no := 1

	for plat, time := range areaParkir {
		fmt.Printf("No %d: Plat %s | Waktu masuk: %s %s\n", no, plat, time.Format("02/01/2006 15:04:05"), "WIB")
		no++
	}

	fmt.Println("================ SISA SLOT PARKIR ================")

	sisaSlot := kapasitasParkir - len(areaParkir)

	fmt.Printf("Sisa slot: %d dari %d\n", sisaSlot, kapasitasParkir)
	fmt.Println("---------------------")
}

func handleParkir(perintah []string, areaParkir map[string]time.Time) (bool, error) {
	commands := perintah[1:]
	noPlat := strings.Join(commands, " ")

	_, exist := areaParkir[noPlat]

	if exist {
    return false, errors.New("Kendaraan dengan plat " + noPlat + " sudah terparkir")
	}

	fmt.Printf("... Akan memproses parkir untuk plat: %s ...\n", noPlat)
	areaParkir[noPlat] = time.Now()

	return true, nil
}

func handleKeluar(perintah []string,areaParkir map[string]time.Time) (bool, error) {
	commands := perintah[1:]
	noPlat := strings.Join(commands, " ")

	_, exist := areaParkir[noPlat]

	if exist {
		fmt.Printf("... Akan memproses keluar untuk plat: %s ...\n", noPlat)

		delete(areaParkir, noPlat)
    return true, nil
	}

	return false, errors.New("kendaraan anda sudah keluar dari parkir")
}

func main() {
	areaParkir := make(map[string]time.Time)

	fmt.Println("======================================")
	fmt.Println("Selamat Datang di Sistem Parkir CLI")
	fmt.Println("Kapasitas Parkir:", kapasitasParkir, "kendaraan")
	fmt.Println("======================================")

	// 'bufio.NewScanner' adalah cara yang baik untuk membaca input baris per baris
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\nMasukkan perintah (parkir/keluar/status/exit): ")

		scanner.Scan()
		input := scanner.Text()

		// Memecah input menjadi beberapa kata (perintah + argumen)
		// strings.Fields lebih baik dari strings.Split karena bisa menangani spasi berlebih
		perintah := strings.Fields(input)

		if len(perintah) == 0 {
			fmt.Println("Perintah tidak boleh kosong!")
			continue
		}

		command := strings.ToLower(perintah[0])

		switch command {
		case "parkir":
			if len(perintah) < 2 {
				fmt.Println("Mohon masukkan nomor plat. Contoh: parkir G 12345 XYZ")
			} else if len(areaParkir) >= kapasitasParkir {
				fmt.Println("============ MAAF AREA PARKIR SUDAH PENUH ===============")
				continue
			}

			result, err := handleParkir(perintah, areaParkir)

			if err != nil {
				fmt.Println(err)
			}

			if result {
				fmt.Println("Berhasil parkir!")
			}

		case "keluar":
			if len(perintah) < 2 {
				fmt.Println("Mohon masukkan nomor plat. Contoh: keluar B1234XYZ")
			} else {
				result, err := handleKeluar(perintah, areaParkir)

				if err != nil {
					fmt.Println(err)
				}

				if result {
					fmt.Println("Berhasil keluar dari parkir!")
				}
			}

		case "status":
			fmt.Println("... Akan menampilkan status parkir ...")
			handleStatus(areaParkir)

		case "exit":
			fmt.Println("Terima kasih telah menggunakan sistem ini. Program berhenti.")
			return

		default:
			fmt.Println("Perintah tidak dikenali. Gunakan: parkir, keluar, status, atau exit.")
		}
	}
}