# Parking System CLI

A simple command-line interface (CLI) parking system implemented in Go. This application allows users to manage parking slots by parking vehicles, removing vehicles, and checking the status of the parking area.

## Features

- Park vehicles with license plate numbers
- Remove vehicles from parking slots
- Check parking status and available slots
- Input validation for all operations
- Capacity management to prevent overbooking

## Installation

### Prerequisites

- Go 1.25 or higher

### Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/misbahkun/go_parking_cli.git
   cd parking_system_cli
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Usage

To run the parking system CLI:

```bash
go run cmd/parking-cli/main.go [capacity]
```

Where `[capacity]` is an optional parameter to set the parking capacity (default is 10).

Examples:
```bash
# Run with default capacity (10)
go run cmd/parking-cli/main.go

# Run with capacity of 20
go run cmd/parking-cli/main.go 20
```

Upon starting, you'll be prompted to enter commands. The system supports the following commands:

### Commands

1. **parkir** - Park a vehicle
   ```
   parkir [license_plate]
   ```
   Example:
   ```
   parkir B 1234 XYZ
   ```

2. **keluar** - Remove a vehicle from parking
   ```
   keluar [license_plate]
   ```
   Example:
   ```
   keluar B 1234 XYZ
   ```

3. **status** - Check parking status and available slots
   ```
   status
   ```

4. **exit** - Exit the application
   ```
   exit
   ```

### Example Session

```
======================================
Selamat Datang di Sistem Parkir CLI
Kapasitas Parkir: 5 kendaraan
======================================

Masukkan perintah (parkir/keluar/status/exit): parkir B 1234 XYZ
... Akan memproses parkir untuk plat: B 1234 XYZ ...
Berhasil parkir!

Masukkan perintah (parkir/keluar/status/exit): status
... Akan menampilkan status parkir ...
================= STATUS PARKIR ================
No 1: Plat B 1234 XYZ | Waktu masuk: 05/09/2025 16:55:00 WIB
================ SISA SLOT PARKIR ================
Sisa slot: 4 dari 5
---------------------

Masukkan perintah (parkir/keluar/status/exit): keluar B 1234 XYZ
... Akan memproses keluar untuk plat: B 1234 XYZ ...
Berhasil keluar dari parkir!

Masukkan perintah (parkir/keluar/status/exit): exit
Terima kasih telah menggunakan sistem ini. Program berhenti.
```

## Testing

To run the tests:

```bash
go test -v
```

To generate a coverage report:

```bash
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Project Structure

```
parking_system_cli/
├── go.mod
├── go.sum
├── parking_cli.go
├── parking_cli_test.go
├── coverage.out
├── README.md
└── cmd/
    └── parking-cli/
        └── main.go
```

## Functions

### Main Functions

- `ParkingCLI(kapasitasParkir int)` - Main function that starts the CLI with a given parking capacity
- `handleParkir(perintah []string, areaParkir map[string]time.Time)` - Handles parking a vehicle
- `handleKeluar(perintah []string, areaParkir map[string]time.Time)` - Handles removing a vehicle
- `handleStatus(kapasitasParkir int, areaParkir map[string]time.Time)` - Displays parking status

### Utility Functions

- `HelloWorld(name *string) string` - Simple greeting function (example/test)

## Error Handling

The application handles various error conditions:

1. Attempting to park a vehicle with an empty license plate
2. Attempting to park a vehicle that is already parked
3. Attempting to remove a vehicle that isn't parked
4. Attempting to park when the parking area is full

## Dependencies

- [github.com/stretchr/testify](https://github.com/stretchr/testify) - For assertions in tests