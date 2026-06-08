# simple_dns

simple_dns is a very simple DNS server written in Go. It was made to learn the basics of the DNS protocol.

## Features

* written in Go using only the standard library
* single-threaded, simple event loop
* hand-rolled DNS message encoding and decoding (no third-party DNS library)
* supports A (IPv4) and AAAA (IPv6) queries
* hardcoded namespace — no master zone file (yet)
* no full protection against malformed requests :|

## Build

```bash
git clone https://github.com/coshi-muhammad/simple_dns.git
cd simple_dns
make build 
```

## Run

Start the server:

```bash
$ ./simple_dns
DNS server listening on :8080
```
or just
```bash
make 
```

## Test

With the server running, use [dig](https://linux.die.net/man/1/dig) in another terminal to send a DNS query.

**IPv4 (A record):**

```bash
$ dig A @127.0.0.1 -p 8080 coshi.com

; <<>> DiG 9.18.x <<>> A @127.0.0.1 -p 8080 coshi.com
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 12345
;; flags: qr aa; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 0

;; QUESTION SECTION:
;coshi.com.                     IN      A

;; ANSWER SECTION:
coshi.com.              18000   IN      A       192.168.1.4

;; Query time: 1 msec
;; SERVER: 127.0.0.1#8080(127.0.0.1)
```

**IPv6 (AAAA record):**

```
$ dig AAAA @127.0.0.1 -p 8080 potato.org

;; ANSWER SECTION:
potato.org.             18000   IN      AAAA    3244:0:301:0:f3:93c:aa00:4
```

You can also run the included test script, which fires several queries and checks the responses automatically:

```
$ ./test.sh
test coshi.com
test passed ✓
test younes.com
test passed ✓
test google.com
test passed ✓
test potato.org
test passed ✓
test historyfacts.org
test passed ✓
test hello.com
test passed ✓
```

Note: `dig` is part of the `dnsutils` package on Debian/Ubuntu (`sudo apt install dnsutils`).

## Hardcoded DNS entries

The namespace is loaded from `loadNameSpace()` in `main.go`. The current entries are:

| Domain            | Type | Address          |
|-------------------|------|------------------|
| coshi.com         | A    | 192.168.1.4      |
| younes.com        | A    | 192.168.1.10     |
| google.com        | A    | 145.255.100.14   |
| potato.org        | AAAA | 3244:0:301:0:f3:93c:aa00:4   |
| historyfacts.org  | AAAA | 1600:f73f:1de1:dd00:f3:4ab:ca00:9 |

To add or change entries, edit `loadNameSpace()` directly. Full master zone file support would be added at some point in the future

## Recommended Reading

The DNS section of the [TCP/IP Guide](http://www.tcpipguide.com/free/t_TCPIPDomainNameSystemDNS.htm) is a great resource for understanding the protocol
