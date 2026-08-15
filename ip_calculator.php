<?php
// ip_calculator.php — PHP версия

function calculateIPv6($cidr) {
    $parts = explode('/', $cidr);
    if (count($parts) != 2) throw new Exception("Неверный формат CIDR");
    $addr = $parts[0];
    $prefix = intval($parts[1]);
    if ($prefix < 0 || $prefix > 128) throw new Exception("Префикс должен быть 0-128");

    // Проверка IPv6
    if (!filter_var($addr, FILTER_VALIDATE_IP, FILTER_FLAG_IPV6)) {
        throw new Exception("Не IPv6 адрес");
    }

    // Создаём inet_pton и битовые операции
    $packed = inet_pton($addr);
    if ($packed === false) throw new Exception("Неверный адрес");

    // Маска
    $mask = str_repeat('f', intdiv($prefix, 4));
    if ($prefix % 4 != 0) {
        $mask .= dechex((hexdec('f') << (4 - ($prefix % 4))) & 0xf);
    }
    $mask = str_pad($mask, 32, '0');
    $maskBin = hex2bin($mask);

    // Сетевой адрес
    $netBin = $packed & $maskBin;
    $network = inet_ntop($netBin);

    // Количество адресов
    $num = pow(2, 128 - $prefix);
    if ($num > 1e18) $num = "∞"; // упрощённо

    $firstHost = $network;
    $lastHost = $network;

    return [
        'ip' => $addr,
        'mask' => inet_ntop($maskBin),
        'cidr' => "/$prefix",
        'network' => $network,
        'first_host' => $firstHost,
        'last_host' => $lastHost,
        'num_addresses' => (string)$num,
        'version' => 'IPv6',
        'prefixlen' => $prefix
    ];
}

function printInfo($info, $input) {
    echo "\033[36m🌐 IP Calculator (IPv6) (PHP)\033[0m\n";
    echo "Входные данные: $input\n\n";
    echo "\033[32m─────────────────────────────────────────\033[0m\n";
    echo "IP-адрес:       {$info['ip']}\n";
    echo "Маска (CIDR):   {$info['cidr']} ({$info['mask']})\n";
    echo "Сетевой адрес:  {$info['network']}\n";
    echo "Количество адресов: {$info['num_addresses']}\n";
    echo "Первый хост:    {$info['first_host']}\n";
    echo "Последний хост: {$info['last_host']}\n";
    echo "Версия:         {$info['version']}\n";
    echo "\033[32m─────────────────────────────────────────\033[0m\n";
}

function saveJSON($info, $input) {
    $data = [
        'timestamp' => date('c'),
        'input' => $input,
        'result' => $info
    ];
    file_put_contents('ip_calc_output.json', json_encode($data, JSON_PRETTY_PRINT));
    echo "\033[32m💾 Сохранено JSON: ip_calc_output.json\033[0m\n";
}

function saveCSV($info) {
    $fp = fopen('ip_calc_output.csv', 'w');
    fputcsv($fp, ['Parameter', 'Value']);
    foreach ($info as $key => $val) {
        fputcsv($fp, [$key, $val]);
    }
    fclose($fp);
    echo "\033[32m💾 Сохранено CSV: ip_calc_output.csv\033[0m\n";
}

function main($argv) {
    if ($argc < 2) {
        echo "Usage: php ip_calculator.php <IPv6/CIDR>\n";
        echo "Пример: php ip_calculator.php 2001:db8::/32\n";
        exit(1);
    }
    $input = $argv[1];
    try {
        $info = calculateIPv6($input);
        printInfo($info, $input);
        saveJSON($info, $input);
        saveCSV($info);
    } catch (Exception $e) {
        echo "\033[31m❌ Ошибка: " . $e->getMessage() . "\033[0m\n";
        exit(1);
    }
}

$argc = $_SERVER['argc'] ?? 0;
$argv = $_SERVER['argv'] ?? [];
main($argv);
?>
