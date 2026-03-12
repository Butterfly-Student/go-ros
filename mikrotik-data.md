1. PPP
1.A. /ppp/secret/add
Creates new item with specified property values.

caller-id -- Sets IP address for PPTP, L2TP or MAC address for PPPoE
comment -- Short description of the item
copy-from -- Item number
disabled -- Defines whether item is ignored or used
ipv6-routes -- 
limit-bytes-in -- maximum amount of bytes user can transmit
limit-bytes-out -- maximum amount of bytes user can receive
local-address -- Assigns an individual address to the PPP-server
name -- Name of the user
password -- User password
profile -- Profile name for the user
remote-address -- Assigns an individual address to the PPP-client
routes -- Routes that appear on the server when the client is connected
service -- Specifies service that will use this user

1.B /ppp/secret/print 

name
service 
aller-id
password
profile
caller-id
routes
local-address
remote-address
ipv6-routes
limit-bytes-in
limit-bytes-out
last-logged-out
comment
last-caller-id
last-disconnect-reason
disabled

1.C /ppp/active/print
name
service
caller-id
encoding
address
uptime

1.C /ppp/profile/add
Creates new item with specified property values.

address-list -- 
bridge -- 
bridge-horizon -- 
bridge-learning -- 
bridge-path-cost -- 
bridge-port-priority -- 
change-tcp-mss -- Change or not TCP protocol's Maximum Segment Size
comment -- Short description of the item
copy-from -- Item number
dns-server -- DNS server address
idle-timeout -- The time limit when the link will be terminated if there
 is no activity
incoming-filter -- Firewall chain name for incoming packets
insert-queue-before -- 
interface-list -- 
local-address -- Assigns an individual address to the PPP-server
name -- Profile name
on-down -- 
on-up -- 
only-one -- Allow only one connection at a time
outgoing-filter -- Firewall chain name for outgoing packets
parent-queue -- 
queue-type -- 
rate-limit -- Data rate limitations for the client
remote-address -- Assigns an individual address to the PPP
session-timeout -- The maximum time the connection can stay up
use-compression -- Defines whether compress traffic or not
use-encryption -- Defines whether encrypt traffic or not
use-mpls -- 
use-upnp -- 
wins-server -- Windows Internet Naming Service server

1.D /ppp/profile/print

address-list -- 
bridge -- 
bridge-horizon -- 
bridge-learning -- 
bridge-path-cost -- 
bridge-port-priority -- 
change-tcp-mss -- Change or not TCP protocol's Maximum Segment Size
comment -- Short description of the item
copy-from -- Item number
dns-server -- DNS server address
idle-timeout -- The time limit when the link will be terminated if there
 is no activity
incoming-filter -- Firewall chain name for incoming packets
insert-queue-before -- 
interface-list -- 
local-address -- Assigns an individual address to the PPP-server
name -- Profile name
on-down -- 
on-up -- 
only-one -- Allow only one connection at a time
outgoing-filter -- Firewall chain name for outgoing packets
parent-queue -- 
queue-type -- 
rate-limit -- Data rate limitations for the client
remote-address -- Assigns an individual address to the PPP
session-timeout -- The maximum time the connection can stay up
use-compression -- Defines whether compress traffic or not
use-encryption -- Defines whether encrypt traffic or not
use-mpls -- 
use-upnp -- 
wins-server -- Windows Internet Naming Service server

2. Hotspot

2.A ip/hotspot/user/add

Creates new item with specified property values.

address -- Static IP address
comment -- Short description of the item
copy-from -- Item number
disabled -- Defines whether item is ignored or used
email -- 
limit-bytes-in -- Total uploaded byte limit for user
limit-bytes-out -- Total downloaded byte limit for user
limit-bytes-total -- Total transferred byte limit for user
limit-uptime -- Total uptime limit for user
mac-address -- Static MAC address
name -- User name
password -- User password
profile -- List of profiles for local HotSpot users
routes -- User routes
server -- Which server is this user allowed to log in to

2.B ip/hotspot/user/print

address -- Static IP address
comment -- Short description of the item
copy-from -- Item number
disabled -- Defines whether item is ignored or used
email -- 
limit-bytes-in -- Total uploaded byte limit for user
limit-bytes-out -- Total downloaded byte limit for user
limit-bytes-total -- Total transferred byte limit for user
limit-uptime -- Total uptime limit for user
mac-address -- Static MAC address
name -- User name
password -- User password
profile -- List of profiles for local HotSpot users
routes -- User routes
server -- Which server is this user allowed to log in to
uptime
bytes-in
bytes-out
packets-in
packets-out

2.D ip/hotspot/active/print
server
user
domain
address
uptime
idle-time
session-time-left
rx-rate
tx-rate

2.E ip/hotspot/profile/add

[admin@G-Net] /ip hotspot> profile add 
copy-from                mac-auth-mode          radius-mac-format 
dns-name                 mac-auth-password      rate-limit        
hotspot-address          name                   smtp-server       
html-directory           nas-port-type          split-user-domain 
html-directory-override  radius-accounting      ssl-certificate   
http-cookie-lifetime     radius-default-domain  trial-uptime-limit
http-proxy               radius-interim-update  trial-uptime-reset
https-redirect           radius-location-id     trial-user-profile
login-by                 radius-location-name   use-radius   

2.F ip/hotspot/profile/print same ass add

2.G ip/hotspot/hosts/print

mac-address
address
to-address
server
idl-time
rx-rate
tx-rate

2.H ip/hotspot/ip-binding/add
[admin@G-Net] /ip hotspot> ip-binding add 
address  copy-from  mac-address   server      type
comment  disabled   place-before  to-address

2.I ip/hotspot/ip-binding/print
address
mac-address
to-address
server
type
comment
disabled

2.J ip/hotspot/cookie/print
user
dodmain
mac-address
expires-in