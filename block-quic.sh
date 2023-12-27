iptables --table filter --insert OUTPUT --protocol udp --dport 80 --jump DROP
iptables --table filter --insert OUTPUT --protocol udp --dport 443 --jump DROP
