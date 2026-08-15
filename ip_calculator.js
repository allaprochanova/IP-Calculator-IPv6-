// ip_calculator.js — JavaScript версия

const fs = require('fs');
const os = require('os');

class IPv6Calculator {
    constructor(cidr) {
        this.cidr = cidr;
        this.parse();
    }

    parse() {
        const parts = this.cidr.split('/');
        if (parts.length !== 2) throw new Error('Неверный формат CIDR');
        this.address = parts[0];
        this.prefix = parseInt(parts[1]);
        if (this.prefix < 0 || this.prefix > 128) throw new Error('Префикс должен быть 0-128');
        // Неполная реализация, используем простые вычисления
        // В реальности нужно использовать библиотеку, но для демонстрации упростим
        this.network = this.address; // Упрощённо
        this.mask = this.calculateMask();
        this.numAddresses = BigInt(1) << BigInt(128 - this.prefix);
        this.firstHost = this.address; // упрощённо
        this.lastHost = this.address; // упрощённо
    }

    calculateMask() {
        // Возвращает строку маски
        return `/${this.prefix}`;
    }

    getInfo() {
        return {
            ip: this.address,
            mask: this.mask,
            cidr: `/${this.prefix}`,
            network: this.network,
            first_host: this.firstHost,
            last_host: this.lastHost,
            num_addresses: this.numAddresses.toString(),
            version: 'IPv6',
            prefixlen: this.prefix
        };
    }

    printInfo() {
        const info = this.getInfo();
        console.log('\x1b[36m🌐 IP Calculator (IPv6) (JavaScript)\x1b[0m');
        console.log(`Входные данные: ${this.cidr}`);
        console.log();
        console.log('\x1b[32m─────────────────────────────────────────\x1b[0m');
        console.log(`IP-адрес:       ${info.ip}`);
        console.log(`Маска (CIDR):   ${info.cidr} (${info.mask})`);
        console.log(`Сетевой адрес:  ${info.network}`);
        console.log(`Количество адресов: ${info.num_addresses}`);
        console.log(`Первый хост:    ${info.first_host}`);
        console.log(`Последний хост: ${info.last_host}`);
        console.log(`Версия:         ${info.version}`);
        console.log('\x1b[32m─────────────────────────────────────────\x1b[0m');
    }

    saveJSON(filename = 'ip_calc_output.json') {
        const data = {
            timestamp: new Date().toISOString(),
            input: this.cidr,
            result: this.getInfo()
        };
        fs.writeFileSync(filename, JSON.stringify(data, null, 2));
        console.log(`\x1b[32m💾 Сохранено JSON: ${filename}\x1b[0m`);
    }

    saveCSV(filename = 'ip_calc_output.csv') {
        const info = this.getInfo();
        let csv = 'Parameter,Value\n';
        for (const [key, val] of Object.entries(info)) {
            csv += `${key},${val}\n`;
        }
        fs.writeFileSync(filename, csv);
        console.log(`\x1b[32m💾 Сохранено CSV: ${filename}\x1b[0m`);
    }
}

function main() {
    const args = process.argv.slice(2);
    if (args.length < 1) {
        console.log('Usage: node ip_calculator.js <IPv6/CIDR>');
        console.log('Пример: node ip_calculator.js 2001:db8::/32');
        process.exit(1);
    }
    const cidr = args[0];
    try {
        const calc = new IPv6Calculator(cidr);
        calc.printInfo();
        calc.saveJSON();
        calc.saveCSV();
    } catch (err) {
        console.error(`\x1b[31m❌ Ошибка: ${err.message}\x1b[0m`);
        process.exit(1);
    }
}

if (require.main === module) main();
