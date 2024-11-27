# FSP Codes and their meanings and specs

## General Code Format

`code param1 param2 param3 ...;`<br>
*each code must be terminated with a semicolon ";"*

## Codes

### join
``join <node_ip> <node_mac> <node_client_version>``<br>
when a new node wants to join the network.
1. node sends request to orchistrator
    1. join node_ip node_mac node_client_version;"
3. 