


**1. Python (ip_calculator.py)**

```python
# ip_calculator.py — Python версия

import sys
import json
import ipaddress
import argparse
from datetime import datetime

def calculate_ipv6(cidr):
    network = ipaddress.ip_network(cidr, strict=False)
    info = {
        "ip": str(network.network_address),
        "mask": str(network.netmask),
        "cidr": f"/{network.prefixlen}",
        "network": str(network.network_address),
        "first_host": str(network.network_address + 1) if network.num_addresses > 1 else "N/A",
        "last_host": str(network.broadcast_address - 1) if network.num_addresses > 1 else "N/A",
        "num_addresses": network.num_addresses,
        "version": f"IPv{network.version}",
        "prefixlen": network.prefixlen,
    }
    return info

def print_info(info):
    print("\n\u001B[36m🌐 IP Calculator (IPv6) (Python)\u001B[0m")
    print(f"Входные данные: {sys.argv[1] if len(sys.argv) > 1 else 'N/A'}")
    print()
    print("\u001B[32m─────────────────────────────────────────\u001B[0m")
    print(f"IP-адрес:       {info['ip']}")
    print(f"Маска (CIDR):   {info['cidr']} ({info['mask']})")
    print(f"Сетевой адрес:  {info['network']}")
    print(f"Количество адресов: {info['num_addresses']}")
    print(f"Первый хост:    {info['first_host']}")
    print(f"Последний хост: {info['last_host']}")
    print(f"Версия:         {info['version']}")
    print("\u001B[32m─────────────────────────────────────────\u001B[0m")

def save_json(info, filename="ip_calc_output.json"):
    data = {
        "timestamp": datetime.now().isoformat(),
        "input": sys.argv[1] if len(sys.argv) > 1 else "",
        "result": info
    }
    with open(filename, 'w') as f:
        json.dump(data, f, indent=2)
    print(f"\u001B[32m💾 Сохранено JSON: {filename}\u001B[0m")

def save_csv(info, filename="ip_calc_output.csv"):
    import csv
    with open(filename, 'w', newline='') as f:
        writer = csv.writer(f)
        writer.writerow(["Parameter", "Value"])
        for key, val in info.items():
            writer.writerow([key, val])
    print(f"\u001B[32m💾 Сохранено CSV: {filename}\u001B[0m")

def main():
    if len(sys.argv) < 2:
        print("Usage: python ip_calculator.py <IPv6/CIDR>")
        print("Пример: python ip_calculator.py 2001:db8::/32")
        sys.exit(1)

    cidr = sys.argv[1]
    try:
        info = calculate_ipv6(cidr)
        print_info(info)
        save_json(info)
        save_csv(info)
    except ValueError as e:
        print(f"\u001B[31m❌ Ошибка: {e}\u001B[0m")
        sys.exit(1)

if __name__ == "__main__":
    main()
