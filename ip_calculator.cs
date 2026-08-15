// ip_calculator.cs — C# версия

using System;
using System.Collections.Generic;
using System.Net;
using System.IO;
using System.Text.Json;
using System.Linq;

class IPv6Info
{
    public string IP { get; set; }
    public string Mask { get; set; }
    public string CIDR { get; set; }
    public string Network { get; set; }
    public string FirstHost { get; set; }
    public string LastHost { get; set; }
    public string NumAddresses { get; set; }
    public string Version { get; set; }
    public int PrefixLen { get; set; }
}

class Program
{
    static void Main(string[] args)
    {
        if (args.Length < 1)
        {
            Console.WriteLine("Usage: dotnet run <IPv6/CIDR>");
            Console.WriteLine("Пример: dotnet run 2001:db8::/32");
            return;
        }
        string cidr = args[0];
        try
        {
            var info = CalculateIPv6(cidr);
            PrintInfo(info, cidr);
            SaveJSON(info, cidr);
            SaveCSV(info);
        }
        catch (Exception e)
        {
            Console.WriteLine($"\u001B[31m❌ Ошибка: {e.Message}\u001B[0m");
        }
    }

    static IPv6Info CalculateIPv6(string cidr)
    {
        var parts = cidr.Split('/');
        if (parts.Length != 2) throw new Exception("Неверный формат CIDR");
        string addr = parts[0];
        int prefix = int.Parse(parts[1]);
        if (prefix < 0 || prefix > 128) throw new Exception("Префикс должен быть 0-128");

        IPAddress ip = IPAddress.Parse(addr);
        if (ip.AddressFamily != System.Net.Sockets.AddressFamily.InterNetworkV6)
            throw new Exception("Не IPv6 адрес");

        // Маска
        byte[] maskBytes = new byte[16];
        int bits = prefix;
        for (int i = 0; i < 16; i++)
        {
            if (bits >= 8) { maskBytes[i] = 0xff; bits -= 8; }
            else if (bits > 0) { maskBytes[i] = (byte)(0xff << (8 - bits)); bits = 0; }
            else maskBytes[i] = 0;
        }
        IPAddress mask = new IPAddress(maskBytes);

        // Сетевой адрес
        byte[] ipBytes = ip.GetAddressBytes();
        byte[] netBytes = new byte[16];
        for (int i = 0; i < 16; i++) netBytes[i] = (byte)(ipBytes[i] & maskBytes[i]);
        IPAddress network = new IPAddress(netBytes);

        // Количество адресов
        BigInteger num = BigInteger.One << (128 - prefix);

        // Первый и последний (упрощённо)
        string firstHost = network.ToString();
        string lastHost = network.ToString();

        return new IPv6Info
        {
            IP = ip.ToString(),
            Mask = mask.ToString(),
            CIDR = "/" + prefix,
            Network = network.ToString(),
            FirstHost = firstHost,
            LastHost = lastHost,
            NumAddresses = num.ToString(),
            Version = "IPv6",
            PrefixLen = prefix
        };
    }

    static void PrintInfo(IPv6Info info, string input)
    {
        Console.WriteLine("\u001B[36m🌐 IP Calculator (IPv6) (C#)\u001B[0m");
        Console.WriteLine($"Входные данные: {input}");
        Console.WriteLine();
        Console.WriteLine("\u001B[32m─────────────────────────────────────────\u001B[0m");
        Console.WriteLine($"IP-адрес:       {info.IP}");
        Console.WriteLine($"Маска (CIDR):   {info.CIDR} ({info.Mask})");
        Console.WriteLine($"Сетевой адрес:  {info.Network}");
        Console.WriteLine($"Количество адресов: {info.NumAddresses}");
        Console.WriteLine($"Первый хост:    {info.FirstHost}");
        Console.WriteLine($"Последний хост: {info.LastHost}");
        Console.WriteLine($"Версия:         {info.Version}");
        Console.WriteLine("\u001B[32m─────────────────────────────────────────\u001B[0m");
    }

    static void SaveJSON(IPv6Info info, string input)
    {
        var data = new
        {
            timestamp = DateTime.Now.ToString("o"),
            input = input,
            result = info
        };
        string json = JsonSerializer.Serialize(data, new JsonSerializerOptions { WriteIndented = true });
        File.WriteAllText("ip_calc_output.json", json);
        Console.WriteLine("\u001B[32m💾 Сохранено JSON: ip_calc_output.json\u001B[0m");
    }

    static void SaveCSV(IPv6Info info)
    {
        using var writer = new StreamWriter("ip_calc_output.csv");
        writer.WriteLine("Parameter,Value");
        writer.WriteLine($"ip,{info.IP}");
        writer.WriteLine($"mask,{info.Mask}");
        writer.WriteLine($"cidr,{info.CIDR}");
        writer.WriteLine($"network,{info.Network}");
        writer.WriteLine($"first_host,{info.FirstHost}");
        writer.WriteLine($"last_host,{info.LastHost}");
        writer.WriteLine($"num_addresses,{info.NumAddresses}");
        writer.WriteLine($"version,{info.Version}");
        writer.WriteLine($"prefixlen,{info.PrefixLen}");
        Console.WriteLine("\u001B[32m💾 Сохранено CSV: ip_calc_output.csv\u001B[0m");
    }
}
