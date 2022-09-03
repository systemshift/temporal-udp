
# Temporal-UDP: A time aware UDP socket

Currently the best way to mitigate denial of service attacks (DDoS) is through installing bigger network capacity. Either by upgrading hardware capacity in house, or by using comercial proxy services.

By using time sensitive sockets to filter unwanted incoming packets, a network operator can augment existing capacity against DDoS flood attacks and increasing the cost and complexity on the spammer.

This will work with the following trade offs:
* A decrease in the possible number of connections
* Low latency variance from median connection



```
                         |
                         |
Spam  -----------------> |
                         | (Close socket)
Spam  -----------------> |
                              SERVER
Client-------------------+->
                           (Open socket)
Spam  -----------------> |
                         | (Close socket)
                         |
```

## Use Cases

### Privileged user access

Under DDoS attack, a network can trade off access for most users in exchange for Privileged users, such as sysadmins, devops teams, account administrators to access the server.


### P2P networks
While p2p networks are hard to attack with DDoS, individual nodes in the network are vulnerable to knock out. Networks that rely on gossip protocols for transactions such as Bitcoin and Ethereum can see DDoS attacks on public facing nodes of miners to delay nodes from propogating a transaction across the networks for a few minutes until another miner finds the block and earns the reward.

```
      ALL CONNECTIONS HAVE TIME SOCKET FILTERS BETWEEN THEM


                      New Node
             +-----------  -------+
Spame Node   v                    v      Spam Node
           +----+              +----+
           |Node|<-------------+Node|
           ++---+              +-^-++
            | ^                  | |
            | |                  | |
            v |                  | v
           +--+-+              +-+--+
           |Node+------------->|Node|
           +----+              +----+


        Spame Node                   Spame Node
```
