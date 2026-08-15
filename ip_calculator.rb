# ip_calculator.rb — Ruby версия

require 'ipaddr'
require 'json'
require 'time'
require 'csv'

def calculate_ipv6(cidr)
  ip = IPAddr.new(cidr)
  network = ip
  mask = ip.inspect.split('/')[0] # упрощённо
  prefix = ip.prefix
  num_addresses = 2 ** (128 - prefix)
  first_host = network.to_s
  last_host = network.to_s
  {
    ip: ip.to_s,
    mask: mask,
    cidr: "/#{prefix}",
    network: network.to_s,
    first_host: first_host,
    last_host: last_host,
    num_addresses: num_addresses.to_s,
    version: "IPv6",
    prefixlen: prefix
  }
end

def print_info(info, input)
  puts "\e[36m🌐 IP Calculator (IPv6) (Ruby)\e[0m"
  puts "Входные данные: #{input}"
  puts
  puts "\e[32m─────────────────────────────────────────\e[0m"
  puts "IP-адрес:       #{info[:ip]}"
  puts "Маска (CIDR):   #{info[:cidr]} (#{info[:mask]})"
  puts "Сетевой адрес:  #{info[:network]}"
  puts "Количество адресов: #{info[:num_addresses]}"
  puts "Первый хост:    #{info[:first_host]}"
  puts "Последний хост: #{info[:last_host]}"
  puts "Версия:         #{info[:version]}"
  puts "\e[32m─────────────────────────────────────────\e[0m"
end

def save_json(info, input)
  data = {
    timestamp: Time.now.iso8601,
    input: input,
    result: info
  }
  File.write('ip_calc_output.json', JSON.pretty_generate(data))
  puts "\e[32m💾 Сохранено JSON: ip_calc_output.json\e[0m"
end

def save_csv(info)
  CSV.open('ip_calc_output.csv', 'w') do |csv|
    csv << ['Parameter', 'Value']
    info.each { |k, v| csv << [k.to_s, v] }
  end
  puts "\e[32m💾 Сохранено CSV: ip_calc_output.csv\e[0m"
end

def main
  if ARGV.length < 1
    puts "Usage: ruby ip_calculator.rb <IPv6/CIDR>"
    puts "Пример: ruby ip_calculator.rb 2001:db8::/32"
    exit 1
  end
  input = ARGV[0]
  begin
    info = calculate_ipv6(input)
    print_info(info, input)
    save_json(info, input)
    save_csv(info)
  rescue => e
    puts "\e[31m❌ Ошибка: #{e.message}\e[0m"
    exit 1
  end
end

main if __FILE__ == $0
