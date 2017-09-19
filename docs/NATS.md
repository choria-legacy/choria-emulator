# NATS Connection and Subscription rates

It's important to see that NATS handles incoming connections at a good rate, other middleware are particularly slow in creating new queues.  We found for example that it can take 10s of minutes to start > 50 000 nodes against RabbitMQ.

Measuring this is particularly difficult because on one hand initial connections are done in a loop so not concurrently and future reconnects are subject to 2 second sleeps.  So take these measurements with some salt but they should show you real world worst case scenarios anyway:

The idea is to start the whole network up and wait for it to be stable.  Then run this:

```
./nats-monitor.rb --out subs.csv
Waiting for num_subscriptions to go below current threshold of 4950
.......................
```

While it does this press ^C on the NATS broker, wait a few seconds and start it again, the display continues like this:

```
!!!!!!!!!!!!!!!!!!!
Waiting for /subsz num_subscriptions, starting at 0
>....................................<
Stable on 4950 num_subscriptions @ 2017-09-19 11:32:30 +0200 elapsed 9.64s
```

You'll have CSV records like this in the file:

```
0.0,0.1,0.2,0.3,0.4,0.5,0.6,0.7,0.8,0.9,1.0,1.1,1.2,1.3,1.4,1.5,1.6,1.7,1.8,1.9,2.0,2.1
3,99,207,241,242,267,258,237,240,270,264,240,234,309,258,270,297,327,174,201,177,135
```

The first row is 0.1 second buckets and the 2nd row is increase in subscriptions during that 0.1 second

You can also monitor the incoming connection rate using:

```
./nats-monitor.rb --out /dev/stdout --source varz --variable connections
```
