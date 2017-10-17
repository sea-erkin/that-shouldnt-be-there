# That Shouldn't Be There
How many times have you performed recon for a client on a pentest engangement upon which the client is shocked with the externally facing systems you've identified. This tool aims to help identify those types of systems through an automated means.

TSBT is a tool which can be used to automate external asset identification based off a domain name. 
TSBT tracks it's own results and alerts on any changes. Examples of an alert would be a new subdomain being
identified for a target domain or a system with 8080 open externally which did not have this port open in the past.

# How it works
The tool can be broken down into the following steps:

1. User provides new line separated target domains via domains.txt file.
2. Subdomain identification via Sublist3r and AltDNS. Output: Subdomains
3. DNS resolving identified subdomains. Output: IP addresses
4. Nmap scanning IP addresses for common web ports. Output: Ports
5. Screenshot identified open ports for provided IP addresses and subdomains. Output: Images
6. Tracking to identify changes over time. Output: Data stored in sqlite database
7. Alerting on identified new hosts or open ports. Output: Email

# Configuration
Currently the alerting module is email based. In order to receive alerts, you must configure an authenticated SMTP account in the ./state/config.json file.

# Use Cases
The main benefit of TSBT is that changes are tracked over time. If you are an organization and are not confident about your external presence, you can configure TSBT to run on your identified external assets and alert you if any changes have happened. Let's say a lousy developer like myself spun up a webserver on a host they did not know was externally facing. TSBT would create an alert and send you a screenshot of the webpage so you can determine for yourself whether the identified web enbled system is legitimate.

If you are a pentester and would like to automate a portion of your recon phase, this is certainly a great tool to run. If you are a pentester on a longer term external engagement, you could configure TSBT to run and alert you if any additional hosts were identified or ports were opened since you first performed your recon.
