#!/bin/sh


test_match(){
  address=$( echo "$1" | grep -A 1 "ANSWER SECTION" |tail -n 1| awk '{print $5}')
  if [[ "$address" == "$2" ]] ; then 
    echo "test passed ✓"
  else 
    echo "test failed x"
  fi
}

# test 1 - 3 ipv4 response
echo "test coshi.com"
response=$(dig A @127.0.0.1 -p 8080 coshi.com)
expected_address="192.168.1.4"
test_match "$response" "$expected_address"

echo "test younes.com"
response=$(dig A @127.0.0.1 -p 8080 younes.com)
expected_address="192.168.1.10"
test_match "$response" "$expected_address"

echo "test google.com"
response=$(dig A @127.0.0.1 -p 8080 google.com)
expected_address="145.255.100.14"
test_match "$response" "$expected_address"

# test 2-5 ipv6 response
echo "test potato.org"
response=$(dig AAAA @127.0.0.1 -p 8080 potato.org)
expected_address="3244:0:301:0:f3:93c:aa00:4"
test_match "$response" "$expected_address"

echo "test historyfacts.org"
response=$(dig AAAA @127.0.0.1 -p 8080 historyfacts.org)
expected_address="1600:f73f:1de1:dd00:f3:4ab:ca00:9"
test_match "$response" "$expected_address"

# test 6 error response
echo "test hello.com"
response=$(dig A @127.0.0.1 -p 8080 hello.com)
expected_address=""
test_match "$response" "$expected_address"

