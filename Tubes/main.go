package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Stock struct {
	Name, Code string
	Price      float64
	Volume     int
}

type User struct {
	Name      string
	Balance   float64
	Portfolio map[string]int
}

func (u *User) BuyStock(stock Stock, quantity int) {
	if quantity <= 0 {
		fmt.Println("Kuantitas untuk membeli harus positif.")
		return
	}
	totalCost := stock.Price * float64(quantity)
	if u.Balance >= totalCost {
		u.Balance -= totalCost
		u.Portfolio[stock.Code] += quantity
		fmt.Printf("Berhasil membeli %d lembar saham %s (%s) seharga $%.2f per lembar. Saldo baru: $%.2f\n", quantity, stock.Name, stock.Code, stock.Price, u.Balance)
	} else {
		fmt.Printf("Saldo tidak mencukupi untuk membeli %d lembar saham %s. Dibutuhkan: $%.2f, Tersedia: $%.2f\n", quantity, stock.Name, totalCost, u.Balance)
	}
}

func (u *User) SellStock(stock Stock, quantity int) {
	if quantity <= 0 {
		fmt.Println("Kuantitas untuk menjual harus positif.")
		return
	}
	if currentQuantity, ok := u.Portfolio[stock.Code]; ok && currentQuantity >= quantity {
		totalSale := stock.Price * float64(quantity)
		u.Balance += totalSale
		u.Portfolio[stock.Code] -= quantity
		if u.Portfolio[stock.Code] == 0 {
			delete(u.Portfolio, stock.Code)
		}
		fmt.Printf("Berhasil menjual %d lembar saham %s (%s) seharga $%.2f per lembar. Saldo baru: $%.2f\n", quantity, stock.Name, stock.Code, stock.Price, u.Balance)
	} else {
		fmt.Printf("Tidak cukup lembar saham %s untuk dijual. Anda memiliki %d, mencoba menjual %d.\n", stock.Name, u.Portfolio[stock.Code], quantity)
	}
}

func (u *User) DisplayPortfolio(marketStocks []Stock) {
	fmt.Printf("\nPortofolio untuk %s:\n", u.Name)
	if len(u.Portfolio) == 0 {
		fmt.Println("  Portofolio Anda kosong.")
	} else {
		fmt.Println("--------------------------------------------------------------------------")
		fmt.Printf("  %-20s | %-10s | %-10s | %-15s | %-15s\n", "Nama Saham", "Kode", "Kuantitas", "Harga Saat Ini", "Nilai Saat Ini")
		fmt.Println("--------------------------------------------------------------------------")
		var totalPortfolioValue float64
		for code, quantity := range u.Portfolio {
			if quantity == 0 {
				continue
			}
			var currentStockPrice float64 = -1
			var stockName string = "N/A"
			foundInMarket := false
			for _, marketStock := range marketStocks {
				if strings.EqualFold(marketStock.Code, code) {
					currentStockPrice = marketStock.Price
					stockName = marketStock.Name
					foundInMarket = true
					break
				}
			}

			if foundInMarket {
				currentValue := currentStockPrice * float64(quantity)
				totalPortfolioValue += currentValue
				fmt.Printf("  %-20s | %-10s | %-10d | $%14.2f | $%14.2f\n", stockName, code, quantity, currentStockPrice, currentValue)
			} else {
				fmt.Printf("  %-20s | %-10s | %-10d | %-15s | %-15s\n", "Tidak Dikenal/Delisting", code, quantity, "Harga T/A", "Nilai T/A")
			}
		}
		fmt.Println("--------------------------------------------------------------------------")
		fmt.Printf("  Total Nilai Saham: $%.2f\n", totalPortfolioValue)
	}
	fmt.Printf("  Saldo Tunai: $%.2f\n", u.Balance)
	fmt.Println("==========================================================================")
}

func SequentialSearch(stocks []Stock, searchTerm string) *Stock {
	searchTermLower := strings.ToLower(searchTerm)
	for i, stock := range stocks {
		if strings.ToLower(stock.Code) == searchTermLower || strings.ToLower(stock.Name) == searchTermLower {
			return &stocks[i]
		}
	}
	return nil
}

func BinarySearch(stocks []Stock, searchCode string) *Stock {
	sortedStocks := make([]Stock, len(stocks))
	copy(sortedStocks, stocks)
	sort.Slice(sortedStocks, func(i, j int) bool {
		return strings.ToLower(sortedStocks[i].Code) < strings.ToLower(sortedStocks[j].Code)
	})

	left, right := 0, len(sortedStocks)-1
	targetCodeLower := strings.ToLower(searchCode)
	var foundStockInSortedSlice *Stock

	for left <= right {
		mid := left + (right-left)/2
		midCodeLower := strings.ToLower(sortedStocks[mid].Code)

		if midCodeLower == targetCodeLower {
			foundStockInSortedSlice = &sortedStocks[mid]
			break
		}
		if midCodeLower < targetCodeLower {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	if foundStockInSortedSlice != nil {
		for i := range stocks {
			if strings.EqualFold(stocks[i].Code, foundStockInSortedSlice.Code) {
				return &stocks[i]
			}
		}
	}
	return nil
}

func SortStocksByPrice(stocks []Stock) {
	sort.Slice(stocks, func(i, j int) bool {
		return stocks[i].Price > stocks[j].Price
	})
}

func SortStocksByVolume(stocks []Stock) {
	sort.Slice(stocks, func(i, j int) bool {
		return stocks[i].Volume > stocks[j].Volume
	})
}

func readString(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func readFloat(prompt string) float64 {
	for {
		s := readString(prompt)
		val, err := strconv.ParseFloat(s, 64)
		if err == nil && val >= 0 {
			return val
		}
		fmt.Println("Input tidak valid. Silakan masukkan angka non-negatif.")
	}
}

func readInt(prompt string) int {
	for {
		s := readString(prompt)
		val, err := strconv.Atoi(s)
		if err == nil && val >= 0 {
			return val
		}
		fmt.Println("Input tidak valid. Silakan masukkan bilangan bulat non-negatif.")
	}
}

func printStocks(stockList []Stock, title string) {
	if title != "" {
		fmt.Println(title)
	}
	if len(stockList) == 0 {
		fmt.Println("Tidak ada saham untuk ditampilkan.")
		return
	}
	fmt.Println("------------------------------------------------------------------")
	fmt.Printf("%-20s | %-10s | %-10s | %-10s\n", "Nama", "Kode", "Harga", "Volume")
	fmt.Println("------------------------------------------------------------------")
	for _, stock := range stockList {
		fmt.Printf("%-20s | %-10s | $%8.2f | %-10d\n", stock.Name, stock.Code, stock.Price, stock.Volume)
	}
	fmt.Println("------------------------------------------------------------------")
}

func printStocksVolume(stockList []Stock, title string) {
	if title != "" {
		fmt.Println(title)
	}
	if len(stockList) == 0 {
		fmt.Println("Tidak ada saham untuk ditampilkan.")
		return
	}
	fmt.Println("------------------------------------------------------------------")
	fmt.Printf("%-20s | %-10s | %-10s | %-10s\n", "Nama", "Kode", "Volume", "Harga")
	fmt.Println("------------------------------------------------------------------")
	for _, stock := range stockList {
		fmt.Printf("%-20s | %-10s | %-10d | $%8.2f\n", stock.Name, stock.Code, stock.Volume, stock.Price)
	}
	fmt.Println("------------------------------------------------------------------")
}

func main() {
	fmt.Println("Selamat Datang di Simulator Perdagangan Saham!")

	numStocks := readInt("Masukkan jumlah saham yang tersedia di pasar: ")
	stocks := make([]Stock, numStocks)
	for i := 0; i < numStocks; i++ {
		fmt.Printf("\nMasukkan detail untuk Saham #%d:\n", i+1)
		stocks[i].Name = readString("  Nama: ")

		for {
			codeCandidate := readString("  Kode (mis., AAPL, GOOG - harus unik): ")
			isUnique := true
			for j := 0; j < i; j++ {
				if strings.EqualFold(stocks[j].Code, codeCandidate) {
					isUnique = false
					fmt.Println("  Kode saham sudah ada. Silakan masukkan kode yang unik.")
					break
				}
			}
			if isUnique {
				stocks[i].Code = codeCandidate
				break
			}
		}
		stocks[i].Price = readFloat("  Harga Saat Ini: $")
		stocks[i].Volume = readInt("  Volume: ")
	}
	fmt.Println("\nSaham pasar berhasil diinisialisasi.")
	printStocks(stocks, "Saham Pasar Saat Ini:")

	userName := readString("\nMasukkan nama Anda: ")
	initialBalance := readFloat("Masukkan saldo awal perdagangan Anda: $")
	user := User{Name: userName, Balance: initialBalance, Portfolio: make(map[string]int)}
	fmt.Printf("\nPengguna %s berhasil dibuat dengan saldo $%.2f\n", user.Name, user.Balance)

	for {
		fmt.Println("\n------------------------------------")
		fmt.Println("Pilih tindakan:")
		fmt.Println("1. Beli Saham")
		fmt.Println("2. Jual Saham")
		fmt.Println("3. Tampilkan Portofolio & Saldo")
		fmt.Println("4. Lihat Saham Pasar (Urut berdasarkan Harga)")
		fmt.Println("5. Lihat Saham Pasar (Urut berdasarkan Volume)")
		fmt.Println("6. Cari Saham (berdasarkan Nama atau Kode - Sekuensial)")
		fmt.Println("7. Cari Saham (berdasarkan Kode - Biner)")
		fmt.Println("8. Hitung Total Untung/Rugi (berdasarkan saldo awal Anda)")
		fmt.Println("9. Keluar")
		fmt.Println("------------------------------------")

		choiceStr := readString("Masukkan pilihan Anda (1-9): ")
		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("Pilihan tidak valid. Silakan masukkan angka antara 1 dan 9.")
			continue
		}

		switch choice {
		case 1:
			fmt.Println("\n--- Beli Saham ---")
			if len(stocks) == 0 {
				fmt.Println("Tidak ada saham yang tersedia di pasar untuk dibeli.")
				continue
			}
			printStocks(stocks, "Saham Tersedia untuk Dibeli:")
			stockCode := readString("Masukkan Kode saham yang akan dibeli: ")
			stockToBuy := SequentialSearch(stocks, stockCode)
			if stockToBuy == nil {
				fmt.Println("Saham tidak ditemukan di pasar.")
				continue
			}
			quantity := readInt(fmt.Sprintf("Masukkan kuantitas %s yang akan dibeli: ", stockToBuy.Name))
			user.BuyStock(*stockToBuy, quantity)

		case 2:
			fmt.Println("\n--- Jual Saham ---")
			if len(user.Portfolio) == 0 {
				fmt.Println("Portofolio Anda kosong. Tidak ada yang bisa dijual.")
				continue
			}
			user.DisplayPortfolio(stocks)
			stockCode := readString("Masukkan Kode saham yang akan dijual: ")

			ownedQuantity, ok := user.Portfolio[strings.ToUpper(stockCode)]
			if !ok || ownedQuantity == 0 {
				fmt.Printf("Anda tidak memiliki saham %s atau kodenya salah.\n", stockCode)
				continue
			}

			stockToSell := SequentialSearch(stocks, stockCode)
			if stockToSell == nil {
				fmt.Println("Error: Detail saham tidak ditemukan di pasar untuk saham dalam portofolio Anda. Ini mungkin menunjukkan saham tersebut telah delisting.")
				continue
			}
			quantity := readInt(fmt.Sprintf("Masukkan kuantitas %s yang akan dijual (Anda memiliki %d): ", stockToSell.Name, ownedQuantity))
			user.SellStock(*stockToSell, quantity)

		case 3:
			fmt.Println("\n--- Portofolio Anda ---")
			user.DisplayPortfolio(stocks)

		case 4:
			fmt.Println("\n--- Saham Pasar (Urut berdasarkan Harga: Tertinggi ke Terendah) ---")
			if len(stocks) == 0 {
				fmt.Println("Tidak ada saham di pasar untuk ditampilkan.")
				continue
			}
			stocksToSort := make([]Stock, len(stocks))
			copy(stocksToSort, stocks)
			SortStocksByPrice(stocksToSort)
			printStocks(stocksToSort, "")

		case 5:
			fmt.Println("\n--- Saham Pasar (Urut berdasarkan Volume: Tertinggi ke Terendah) ---")
			if len(stocks) == 0 {
				fmt.Println("Tidak ada saham di pasar untuk ditampilkan.")
				continue
			}
			stocksToSort := make([]Stock, len(stocks))
			copy(stocksToSort, stocks)
			SortStocksByVolume(stocksToSort)
			printStocksVolume(stocksToSort, "")

		case 6:
			fmt.Println("\n--- Pencarian Saham Sekuensial ---")
			if len(stocks) == 0 {
				fmt.Println("Tidak ada saham di pasar untuk dicari.")
				continue
			}
			searchTerm := readString("Masukkan Nama atau Kode saham untuk dicari: ")
			foundStock := SequentialSearch(stocks, searchTerm)
			if foundStock != nil {
				fmt.Printf("Ditemukan: %s (%s), Harga: $%.2f, Volume: %d\n", foundStock.Name, foundStock.Code, foundStock.Price, foundStock.Volume)
			} else {
				fmt.Println("Saham tidak ditemukan menggunakan pencarian sekuensial.")
			}

		case 7:
			fmt.Println("\n--- Pencarian Saham Biner (berdasarkan Kode) ---")
			if len(stocks) == 0 {
				fmt.Println("Tidak ada saham di pasar untuk dicari.")
				continue
			}
			searchTerm := readString("Masukkan Kode saham untuk dicari: ")
			foundStock := BinarySearch(stocks, searchTerm)
			if foundStock != nil {
				fmt.Printf("Ditemukan: %s (%s), Harga: $%.2f, Volume: %d\n", foundStock.Name, foundStock.Code, foundStock.Price, foundStock.Volume)
			} else {
				fmt.Println("Saham tidak ditemukan menggunakan pencarian biner berdasarkan kode.")
			}

		case 8:
			fmt.Println("\n--- Perhitungan Untung/Rugi ---")
			currentTotalValue := user.Balance
			for code, quantity := range user.Portfolio {
				stockFound := false
				for _, marketStock := range stocks {
					if strings.EqualFold(marketStock.Code, code) {
						currentTotalValue += marketStock.Price * float64(quantity)
						stockFound = true
						break
					}
				}
				if !stockFound {
					fmt.Printf("Peringatan: Saham %s dalam portofolio tidak ditemukan di pasar saat ini untuk perhitungan U/R.\n", code)
				}
			}
			profitLoss := currentTotalValue - initialBalance
			fmt.Printf("Saldo Awal: $%.2f\n", initialBalance)
			fmt.Printf("Total Nilai Portofolio Saat Ini (Tunai + Saham): $%.2f\n", currentTotalValue)
			if profitLoss >= 0 {
				fmt.Printf("Total Untung: $%.2f\n", profitLoss)
			} else {
				fmt.Printf("Total Rugi: $%.2f\n", -profitLoss)
			}

		case 9:
			fmt.Println("\nTerima kasih telah menggunakan Simulator Perdagangan Saham. Sampai jumpa!")
			return

		default:
			fmt.Println("Pilihan tidak valid. Silakan masukkan angka antara 1 dan 9.")
		}
	}
}
