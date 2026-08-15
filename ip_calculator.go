// ip_calculator.go — Go версия

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

type IPv6Info struct {
	IP           string   `json:"ip"`
	Mask         string   `json:"mask"`
	CIDR         string   `json:"cidr"`
	Network      string   `json:"network"`
	FirstHost    string   `json:"first_host"`
	LastHost     string   `json:"last_host"`
	NumAddresses uint64   `json:"num_addresses"`
	Version      string   `json:"version"`
	PrefixLen    int      `json:"prefixlen"`
}

func calculateIPv6(cidr string) (*IPv6Info, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}
	// Для IPv6
	ones, bits := ipnet.Mask.Size()
	if bits != 128 {
		return nil, fmt.Errorf("не IPv6 адрес")
	}
	mask := net.IP(ipnet.Mask)

	// Первый и последний адрес (для IPv6)
	first := ip.Mask(ipnet.Mask)
	last := make(net.IP, len(ip))
	copy(last, ip)
	// Инвертируем маску для получения broadcast
	maskInv := make(net.IP, len(mask))
	for i := range mask {
		maskInv[i] = ^mask[i]
	}
	last = last.Mask(ipnet.Mask)
	// Добавляем маску инвертированную
	for i := range last {
		last[i] |= maskInv[i]
	}
	// Вычисляем количество адресов (2^(128-prefix))
	num := uint64(1) << (128 - ones)
	firstHost := first.String()
	lastHost := last.String()
	if num == 1 {
		firstHost = "N/A"
		lastHost = "N/A"
	} else {
		// Для первого и последнего хоста нужно добавить/убавить 1
		// В IPv6 это сложно, но можно упрощённо показать
		// Для простоты оставим как сетевой и широковещательный
		// В IPv6 нет широковещательного, но мы можем показать первый и последний
	}
	// Упрощённо: показываем как есть
	return &IPv6Info{
		IP:           ip.String(),
		Mask:         mask.String(),
		CIDR:         fmt.Sprintf("/%d", ones),
		Network:      ipnet.IP.String(),
		FirstHost:    firstHost,
		LastHost:     lastHost,
		NumAddresses: num,
		Version:      "IPv6",
		PrefixLen:    ones,
	}, nil
}

func printInfo(info *IPv6Info) {
	fmt.Println("\x1b[36m🌐 IP Calculator (IPv6) (Go)\x1b[0m")
	fmt.Printf("Входные данные: %s\n", os.Args[1])
	fmt.Println()
	fmt.Println("\x1b[32m─────────────────────────────────────────\x1b[0m")
	fmt.Printf("IP-адрес:       %s\n", info.IP)
	fmt.Printf("Маска (CIDR):   %s (%s)\n", info.CIDR, info.Mask)
	fmt.Printf("Сетевой адрес:  %s\n", info.Network)
	fmt.Printf("Количество адресов: %d\n", info.NumAddresses)
	fmt.Printf("Первый хост:    %s\n", info.FirstHost)
	fmt.Printf("Последний хост: %s\n", info.LastHost)
	fmt.Printf("Версия:         %s\n", info.Version)
	fmt.Println("\x1b[32m─────────────────────────────────────────\x1b[0m")
}

func saveJSON(info *IPv6Info) {
	data := map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
		"input":     os.Args[1],
		"result":    info,
	}
	jsonData, _ := json.MarshalIndent(data, "", "  ")
	os.WriteFile("ip_calc_output.json", jsonData, 0644)
	fmt.Printf("\x1b[32m💾 Сохранено JSON: ip_calc_output.json\x1b[0m\n")
}

func saveCSV(info *IPv6Info) {
	file, _ := os.Create("ip_calc_output.csv")
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	writer.Write([]string{"Parameter", "Value"})
	writer.Write([]string{"IP", info.IP})
	writer.Write([]string{"Mask", info.Mask})
	writer.Write([]string{"CIDR", info.CIDR})
	writer.Write([]string{"Network", info.Network})
	writer.Write([]string{"FirstHost", info.FirstHost})
	writer.Write([]string{"LastHost", info.LastHost})
	writer.Write([]string{"NumAddresses", strconv.FormatUint(info.NumAddresses, 10)})
	writer.Write([]string{"Version", info.Version})
	writer.Write([]string{"PrefixLen", strconv.Itoa(info.PrefixLen)})
	fmt.Printf("\x1b[32m💾 Сохранено CSV: ip_calc_output.csv\x1b[0m\n")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run ip_calculator.go <IPv6/CIDR>")
		fmt.Println("Пример: go run ip_calculator.go 2001:db8::/32")
		os.Exit(1)
	}
	cidr := os.Args[1]
	info, err := calculateIPv6(cidr)
	if err != nil {
		fmt.Printf("\x1b[31m❌ Ошибка: %v\x1b[0m\n", err)
		os.Exit(1)
	}
	printInfo(info)
	saveJSON(info)
	saveCSV(info)
}
