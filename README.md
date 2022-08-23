# Dashwood
Dashwood is a scalable, clusterable, high-availability, and most importantly **fast** DNS server. Its the main backend for VelocityDNS, the DNS provider for VelocityHost.

# How it works?
![proposal-dns drawio](https://user-images.githubusercontent.com/56168307/186231491-b03d8b84-a8a7-48e0-8516-0f5040da7b89.png)
Above is a diagram basing out the structure. It is comprised of two main components (which you will find in this repo) the Control Surface (API Provider) and the DNS Agent which runs on the servers to serve DNS requests.

Name|Contribution
|-------------|-------------------|
|[Nathan Hedge](https://github.com/10nates)| DNS Agent and Control Surface |
|[Owen Rummage](https://github.com/ash-quinn)| API Docs, README, [Frontend Development and Design](#)|
|[Logan Shaw](https://github.com/themoddedchicken)| [Frontend Development and Design](#)|
