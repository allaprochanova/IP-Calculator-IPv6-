// ip_calculator.rs — Rust версия

use std::env;
use std::net::{Ipv6Addr, Ipv6Network};
use std::fs;
use std::collections::HashMap;
use serde::{Deserialize, Serialize};
use serde_json;

#[derive(Debug, Serialize, Deserialize)]
struct IPv6Info {
    ip: String,
    mask: String,
    cidr: String,
    network: String,
    first_host: String,
    last_host: String,
    num_addresses: String,
    version: String,
    prefixlen: u8,
}

fn calculate_ipv6(cidr: &str) -> Result<IPv6Info, Box<dyn std::error::Error>> {
    let network: Ipv6Network = cidr.parse()?;
    let ip = network.network();
    let mask = network.mask();
    let prefix = network.prefix();
    let num = 1u128.checked_shl(128 - prefix as u32).unwrap_or(0);
    let num_str = if num == 0 { "0".to_string() } else { num.to_string() };

    // Первый и последний (упрощённо)
    let first_host = ip.to_string();
    let last_host = ip.to_string();

    Ok(IPv6Info {
        ip: ip.to_string(),
        mask: mask.to_string(),
        cidr: format!("/{}", prefix),
        network: ip.to_string(),
        first_host,
        last_host,
        num_addresses: num_str,
        version: "IPv6".to_string(),
        prefixlen: prefix,
    })
}

fn print_info(info: &IPv6Info, input: &str) {
    println!("\x1b[36m🌐 IP Calculator (IPv6) (Rust)\x1b[0m");
    println!("Входные данные: {}", input);
    println!();
    println!("\x1b[32m─────────────────────────────────────────\x1b[0m");
    println!("IP-адрес:       {}", info.ip);
    println!("Маска (CIDR):   {} ({})", info.cidr, info.mask);
    println!("Сетевой адрес:  {}", info.network);
    println!("Количество адресов: {}", info.num_addresses);
    println!("Первый хост:    {}", info.first_host);
    println!("Последний хост: {}", info.last_host);
    println!("Версия:         {}", info.version);
    println!("\x1b[32m─────────────────────────────────────────\x1b[0m");
}

fn save_json(info: &IPv6Info, input: &str) -> Result<(), Box<dyn std::error::Error>> {
    let data = serde_json::json!({
        "timestamp": chrono::Utc::now().to_rfc3339(),
        "input": input,
        "result": info
    });
    let json = serde_json::to_string_pretty(&data)?;
    fs::write("ip_calc_output.json", json)?;
    println!("\x1b[32m💾 Сохранено JSON: ip_calc_output.json\x1b[0m");
    Ok(())
}

fn save_csv(info: &IPv6Info) -> Result<(), Box<dyn std::error::Error>> {
    let mut csv = String::from("Parameter,Value\n");
    csv.push_str(&format!("ip,{}\n", info.ip));
    csv.push_str(&format!("mask,{}\n", info.mask));
    csv.push_str(&format!("cidr,{}\n", info.cidr));
    csv.push_str(&format!("network,{}\n", info.network));
    csv.push_str(&format!("first_host,{}\n", info.first_host));
    csv.push_str(&format!("last_host,{}\n", info.last_host));
    csv.push_str(&format!("num_addresses,{}\n", info.num_addresses));
    csv.push_str(&format!("version,{}\n", info.version));
    csv.push_str(&format!("prefixlen,{}\n", info.prefixlen));
    fs::write("ip_calc_output.csv", csv)?;
    println!("\x1b[32m💾 Сохранено CSV: ip_calc_output.csv\x1b[0m");
    Ok(())
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    let args: Vec<String> = env::args().collect();
    if args.len() < 2 {
        println!("Usage: cargo run -- <IPv6/CIDR>");
        println!("Пример: cargo run -- 2001:db8::/32");
        std::process::exit(1);
    }
    let input = &args[1];
    let info = calculate_ipv6(input)?;
    print_info(&info, input);
    save_json(&info, input)?;
    save_csv(&info)?;
    Ok(())
}
