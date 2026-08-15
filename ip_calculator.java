// ip_calculator.java — Java версия

import java.io.*;
import java.net.*;
import java.nio.file.*;
import java.time.*;
import java.util.*;
import com.google.gson.Gson;
import com.google.gson.GsonBuilder;

public class ip_calculator {
    public static void main(String[] args) throws Exception {
        if (args.length < 1) {
            System.out.println("Usage: java ip_calculator <IPv6/CIDR>");
            System.out.println("Пример: java ip_calculator 2001:db8::/32");
            System.exit(1);
        }
        String cidr = args[0];
        try {
            IPv6Info info = calculateIPv6(cidr);
            printInfo(info, cidr);
            saveJSON(info, cidr);
            saveCSV(info);
        } catch (Exception e) {
            System.out.println("\u001B[31m❌ Ошибка: " + e.getMessage() + "\u001B[0m");
            System.exit(1);
        }
    }

    static IPv6Info calculateIPv6(String cidr) throws Exception {
        // Разбор CIDR
        String[] parts = cidr.split("/");
        if (parts.length != 2) throw new Exception("Неверный формат CIDR");
        String addr = parts[0];
        int prefix = Integer.parseInt(parts[1]);
        if (prefix < 0 || prefix > 128) throw new Exception("Префикс должен быть 0-128");

        InetAddress inet = InetAddress.getByName(addr);
        if (!(inet instanceof Inet6Address)) throw new Exception("Не IPv6 адрес");

        // Маска
        byte[] mask = new byte[16];
        int bits = prefix;
        for (int i = 0; i < 16; i++) {
            if (bits >= 8) {
                mask[i] = (byte)0xff;
                bits -= 8;
            } else if (bits > 0) {
                mask[i] = (byte)(0xff << (8 - bits));
                bits = 0;
            } else {
                mask[i] = 0;
            }
        }
        InetAddress maskAddr = InetAddress.getByAddress(mask);

        // Сетевой адрес
        byte[] netBytes = new byte[16];
        byte[] addrBytes = inet.getAddress();
        for (int i = 0; i < 16; i++) {
            netBytes[i] = (byte)(addrBytes[i] & mask[i]);
        }
        InetAddress network = InetAddress.getByAddress(netBytes);

        // Количество адресов
        BigInteger num = BigInteger.ONE.shiftLeft(128 - prefix);

        // Первый и последний (упрощённо)
        String firstHost = network.getHostAddress();
        String lastHost = network.getHostAddress();

        IPv6Info info = new IPv6Info();
        info.ip = inet.getHostAddress();
        info.mask = maskAddr.getHostAddress();
        info.cidr = "/" + prefix;
        info.network = network.getHostAddress();
        info.firstHost = firstHost;
        info.lastHost = lastHost;
        info.numAddresses = num.toString();
        info.version = "IPv6";
        info.prefixlen = prefix;
        return info;
    }

    static class IPv6Info {
        String ip, mask, cidr, network, firstHost, lastHost, numAddresses, version;
        int prefixlen;
    }

    static void printInfo(IPv6Info info, String input) {
        System.out.println("\u001B[36m🌐 IP Calculator (IPv6) (Java)\u001B[0m");
        System.out.println("Входные данные: " + input);
        System.out.println();
        System.out.println("\u001B[32m─────────────────────────────────────────\u001B[0m");
        System.out.printf("IP-адрес:       %s\n", info.ip);
        System.out.printf("Маска (CIDR):   %s (%s)\n", info.cidr, info.mask);
        System.out.printf("Сетевой адрес:  %s\n", info.network);
        System.out.printf("Количество адресов: %s\n", info.numAddresses);
        System.out.printf("Первый хост:    %s\n", info.firstHost);
        System.out.printf("Последний хост: %s\n", info.lastHost);
        System.out.printf("Версия:         %s\n", info.version);
        System.out.println("\u001B[32m─────────────────────────────────────────\u001B[0m");
    }

    static void saveJSON(IPv6Info info, String input) throws IOException {
        Map<String, Object> data = new LinkedHashMap<>();
        data.put("timestamp", Instant.now().toString());
        data.put("input", input);
        data.put("result", info);
        Gson gson = new GsonBuilder().setPrettyPrinting().create();
        String json = gson.toJson(data);
        Files.write(Paths.get("ip_calc_output.json"), json.getBytes());
        System.out.println("\u001B[32m💾 Сохранено JSON: ip_calc_output.json\u001B[0m");
    }

    static void saveCSV(IPv6Info info) throws IOException {
        try (PrintWriter pw = new PrintWriter(new File("ip_calc_output.csv"))) {
            pw.println("Parameter,Value");
            pw.printf("ip,%s\n", info.ip);
            pw.printf("mask,%s\n", info.mask);
            pw.printf("cidr,%s\n", info.cidr);
            pw.printf("network,%s\n", info.network);
            pw.printf("first_host,%s\n", info.firstHost);
            pw.printf("last_host,%s\n", info.lastHost);
            pw.printf("num_addresses,%s\n", info.numAddresses);
            pw.printf("version,%s\n", info.version);
            pw.printf("prefixlen,%d\n", info.prefixlen);
        }
        System.out.println("\u001B[32m💾 Сохранено CSV: ip_calc_output.csv\u001B[0m");
    }
}
