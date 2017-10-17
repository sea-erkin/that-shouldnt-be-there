# That Shouldn't Be There
TSBT is a tool which can be used to automate external asset identification based off a domain name. 
TSBT tracks it's own results and alerts on any changes. Examples of an alert would be a new subdomain being
identified for a target domain or a system with 8080 open externally which did not have this port open in the past.

# How it works
The tool can be broken down into the following steps:

1. Subdomain identification via Sublist3r and AltDNS. Output: Subdomains
2. DNS resolving identified subdomains. Output: IP addresses
3. Nmap scanning IP addresses for common web ports. Output: Ports
4. Screenshot identified open ports for provided IP addresses and subdomains. Output: Images
5. Tracking to identify changes over time. Output: Data stored in sqlite database
6. Alerting on identified new hosts or open ports. Output: Email

# Configuration
Currently the alerting module is email based. In order to receive alerts, you must configure an authenticated SMTP account in the ./state/config.json file.

