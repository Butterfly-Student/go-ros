# Dokumentasi Teknis Lengkap Mikhmon v4

## Daftar Isi
1. [Overview](#1-overview)
2. [Struktur Project](#2-struktur-project)
3. [Komunikasi dengan MikroTik](#3-komunikasi-dengan-mikrotik)
4. [Command MikroTik API Lengkap](#4-command-mikrotik-api-lengkap)
5. [Script RouterOS (MikroTik)](#5-script-routeros-mikrotik)
6. [JavaScript Functions](#6-javascript-functions)
7. [Data Flow](#7-data-flow)
8. [Security](#8-security)
9. [Utility Functions](#9-utility-functions)
10. [API Endpoints](#10-api-endpoints)

---

## 1. Overview

Mikhmon v4 adalah aplikasi web-based management system untuk MikroTik Hotspot yang dibangun menggunakan PHP dan JavaScript. Aplikasi ini memungkinkan administrator untuk mengelola user hotspot, generate voucher, monitoring traffic, dan mengelola report penjualan.

### Fitur Utama:
- **User Management**: Add, edit, remove hotspot users
- **Voucher Generator**: Generate voucher dalam batch dengan berbagai pattern
- **Profile Management**: Konfigurasi profile dengan on-login script
- **Traffic Monitoring**: Real-time traffic monitoring per interface
- **Reporting**: Sales report yang tersimpan di MikroTik
- **Expire Monitoring**: Auto-disable/remove user yang expired
- **Multi-Router**: Support multiple MikroTik routers

---

## 2. Struktur Project

```
mikhmon_v4/
├── index.php                    # Entry point aplikasi
├── config/                      # Konfigurasi
│   ├── config.php              # Konfigurasi database/session
│   ├── connection.php          # Koneksi ke MikroTik
│   ├── page.php                # Routing page
│   ├── readcfg.php             # Read configuration
│   ├── settheme.php            # Theme settings
│   └── theme.php               # Theme loader
├── core/                        # Core functionality
│   ├── routeros_api.class.php  # RouterOS API Class
│   ├── page_route.php          # Page routing logic
│   ├── route.php               # Route handler
│   ├── jsencode.class.php      # JavaScript encoding
│   ├── generator_functions.php # User generator functions
│   └── no_cache.php            # Cache control
├── get/                         # API GET endpoints
│   ├── get_users.php           # Get hotspot users
│   ├── get_user.php            # Get single user
│   ├── get_profiles.php        # Get user profiles
│   ├── get_profile.php         # Get single profile
│   ├── get_hotspot_active.php  # Get active users
│   ├── get_hotspot_server.php  # Get hotspot servers
│   ├── get_traffic.php         # Get interface traffic
│   ├── get_interface.php       # Get interfaces
│   ├── get_dashboard.php       # Dashboard data
│   ├── get_hosts.php           # Get hotspot hosts
│   ├── get_addr_pool.php       # Get address pools
│   ├── get_parent_queue.php    # Get parent queues
│   ├── get_nat.php             # Get NAT rules
│   ├── get_report.php          # Get reports
│   ├── get_tot_users.php       # Get total users
│   ├── get_expire_mon.php      # Get expire monitor
│   ├── get_connect.php         # Test connection
│   └── get_sys_resource.php    # Get system resource
├── post/                        # API POST endpoints
│   ├── post_add_user.php       # Add hotspot user
│   ├── post_update_user.php    # Update user
│   ├── post_add_userprofile.php    # Add profile
│   ├── post_update_userprofile.php # Update profile
│   ├── post_hotspot_remove.php     # Remove user/profile/host/active
│   ├── post_generate_voucher.php   # Generate vouchers
│   ├── post_cache_voucher.php      # Cache voucher data
│   ├── post_expire_monitor.php     # Setup expire monitor
│   ├── post_a_router.php           # Router management
│   ├── post_logout.php             # Logout handler
│   └── post_template.php           # Template management
├── view/                        # View templates
│   └── print_voucher.php       # Print voucher template
├── assets/                      # Static assets
│   ├── js/
│   │   ├── mikhmon.js          # Main JavaScript
│   │   └── func.js             # Utility functions
│   ├── css/                    # Stylesheets
│   └── img/                    # Images
└── template/                    # Voucher templates
    ├── header.default.txt      # Default header template
    ├── row.default.txt         # Default row template
    ├── footer.default.txt      # Default footer template
    ├── header.small.txt        # Small header template
    ├── row.small.txt           # Small row template
    ├── footer.small.txt        # Small footer template
    ├── header.thermal.txt      # Thermal header template
    ├── row.thermal.txt         # Thermal row template
    └── footer.thermal.txt      # Thermal footer template
```

---

## 3. Komunikasi dengan MikroTik

### 3.1 RouterOS API Class

File: `core/routeros_api.class.php`

Class `RouterosAPI` adalah wrapper untuk komunikasi dengan RouterOS API MikroTik menggunakan protokol binary API.

#### Konfigurasi Koneksi:
```php
$API = new RouterosAPI();
$API->debug = false;     // Debug mode
$API->port = 8728;       // Default API port (8729 untuk SSL)
$API->ssl = false;       // SSL connection
$API->timeout = 3;       // Connection timeout
$API->attempts = 5;      // Connection attempts
$API->delay = 3;         // Delay between attempts
```

#### Metode Utama:
```php
// Connect to MikroTik
$API->connect($ip, $username, $password);

// Disconnect
$API->disconnect();

// Send command with parameters
$API->comm($command, $params);

// Write command
$API->write($command, $param2);

// Read response
$API->read($parse = true);
```

### 3.2 Enkripsi Password

File: `core/routeros_api.class.php`

```php
// Encrypt - XOR dengan key 128 + Base64
function enc_rypt($string, $key = 128) {
    $result = '';
    for ($i = 0, $k = strlen($string); $i < $k; $i++) {
        $char = substr($string, $i, 1);
        $keychar = substr($key, ($i % strlen($key)) - 1, 1);
        $char = chr(ord($char) + ord($keychar));
        $result .= $char;
    }
    return base64_encode($result);
}

// Decrypt
function dec_rypt($string, $key = 128) {
    $result = '';
    $string = base64_decode($string);
    for ($i = 0, $k = strlen($string); $i < $k; $i++) {
        $char = substr($string, $i, 1);
        $keychar = substr($key, ($i % strlen($key)) - 1, 1);
        $char = chr(ord($char) - ord($keychar));
        $result .= $char;
    }
    return $result;
}
```

### 3.3 JavaScript Encoding

File: `assets/js/func.js`

```javascript
// XOR decoder untuk response dari server
var jesD = {
  dec: function (e, t = 25) {
    var a = "";
    for (var s = 0; s < e.length; s++) {
      var o = e.charCodeAt(s) ^ t;
      a += String.fromCharCode(o);
    }
    return a;
  }
};

// Password encoding untuk dikirim ke server
function blah(e) {
  var t = "";
  for (e = btoa(e), e = btoa(e), i = 0; i < e.length; i++) {
    var a = 10 ^ e.charCodeAt(i);
    t += String.fromCharCode(a);
  }
  return (t = btoa(t));
}

// Password decoding
function unblah(encoded) {
  try {
    let decoded = atob(encoded);
    let xorResult = "";
    for (let i = 0; i < decoded.length; i++) {
      xorResult += String.fromCharCode(decoded.charCodeAt(i) ^ 10);
    }
    let original = atob(atob(xorResult));
    return original;
  } catch (e) {
    console.error("Gagal decode:", e);
    return null;
  }
}
```

---

## 4. Command MikroTik API Lengkap (59 Commands)

### 4.1 Hotspot User Management (15 commands)

#### Get Users by Profile
```php
// File: get/get_users.php
$get_users = $API->comm("/ip/hotspot/user/print", array(
    "?profile" => "$uprof"
));
```
**Parameter:**
- `?profile` - Filter berdasarkan profile name

**Response:** Array of user objects dengan properties:
- `.id` - Unique identifier (e.g., "*1F")
- `name` - Username
- `password` - Password
- `profile` - Profile name
- `mac-address` - MAC address
- `limit-uptime` - Time limit (e.g., "1d2h3m")
- `limit-bytes-total` - Data limit dalam bytes
- `comment` - Comment
- `uptime` - Current uptime
- `bytes-in` - Bytes received
- `bytes-out` - Bytes sent
- `disabled` - Status disabled (yes/no)
- `server` - Hotspot server

#### Get Single User by ID
```php
// File: get/get_user.php, view/print_voucher.php
$get_users = $API->comm("/ip/hotspot/user/print", array(
    "?.id" => "$uid"
));
```

#### Get Single User by Name
```php
// File: get/get_user.php, get/get_tot_users.php
$get_users = $API->comm("/ip/hotspot/user/print", array(
    "?name" => "$uname"
));
```

#### Get Users by Comment (untuk voucher)
```php
// File: post/post_generate_voucher.php, post/post_cache_voucher.php
$get_hotspotusers = $API->comm("/ip/hotspot/user/print", array(
    "?comment" => "$commt"
));

// File: view/print_voucher.php
$get_users = $API->comm('/ip/hotspot/user/print', array(
    "?comment" => "$comment", 
    "?uptime" => "0s"
));
```

#### Get Total Users Count
```php
// File: get/get_dashboard.php, get/get_tot_users.php
$get_hotspotusers = $API->comm("/ip/hotspot/user/print", array(
    "count-only" => ""
));
```

#### Add New User
```php
// File: post/post_add_user.php, post/post_generate_voucher.php
$add = $API->comm("/ip/hotspot/user/add", array(
    "server" => "$server",
    "name" => "$name",
    "password" => "$password",
    "profile" => "$profile",
    "mac-address" => "$mac_addr",
    "disabled" => "no",
    "limit-uptime" => "$timelimit",
    "limit-bytes-total" => "$datalimit",
    "comment" => "$comment"
));
```
**Parameters:**
- `server` - Hotspot server name (string)
- `name` - Username (string)
- `password` - Password (string)
- `profile` - User profile name (string)
- `mac-address` - MAC address format "00:00:00:00:00:00" (string)
- `disabled` - "yes" atau "no" (string)
- `limit-uptime` - Time limit format "1d2h3m4s" (string)
- `limit-bytes-total` - Data limit dalam bytes (integer)
- `comment` - Comment (string)

**Response:**
- Success: Returns `.id` (e.g., "*1F")
- Error: Returns `!trap` dengan message error

#### Update User
```php
// File: post/post_update_user.php
$add = $API->comm("/ip/hotspot/user/set", array(
    ".id" => "$uid",
    "server" => "$server",
    "name" => "$name",
    "password" => "$password",
    "profile" => "$profile",
    "mac-address" => "$mac_addr",
    "disabled" => "no",
    "limit-uptime" => "$timelimit",
    "limit-bytes-total" => "$datalimit",
    "comment" => "$comment"
));
```

#### Reset User Counters
```php
// File: post/post_update_user.php
$API->comm("/ip/hotspot/user/reset-counters", array(
    ".id" => "$uid"
));
```

#### Remove User
```php
// File: post/post_hotspot_remove.php
$API->comm("/ip/hotspot/user/remove", array(
    ".id" => $id
));
```

### 4.2 Hotspot Profile Management (7 commands)

#### Get All Profiles
```php
// File: get/get_profiles.php
$get_user_profiles = $API->comm("/ip/hotspot/user/profile/print");
```

**Response Properties:**
- `.id` - Unique identifier
- `name` - Profile name
- `address-pool` - Address pool name
- `rate-limit` - Rate limit (e.g., "1M/2M")
- `shared-users` - Max shared users (integer)
- `status-autorefresh` - Autorefresh interval
- `on-login` - On-login script
- `parent-queue` - Parent queue name

#### Get Single Profile by ID
```php
// File: get/get_profile.php, post/post_add_userprofile.php, post/post_update_userprofile.php
$get_uprofile = $API->comm("/ip/hotspot/user/profile/print", array(
    "?.id" => "$upid"
));
```

#### Get Single Profile by Name
```php
// File: get/get_profile.php, view/print_voucher.php
$get_uprofile = $API->comm("/ip/hotspot/user/profile/print", array(
    "?name" => "$upname"
));
```

#### Add Profile
```php
// File: post/post_add_userprofile.php
$add = $API->comm("/ip/hotspot/user/profile/add", array(
    "name" => "$name",
    "address-pool" => "$addrpool",
    "rate-limit" => "$ratelimit",
    "shared-users" => "$sharedusers",
    "status-autorefresh" => "1m",
    "on-login" => "$onlogin",
    "parent-queue" => "$parent"
));
```
**Parameters:**
- `name` - Profile name (string)
- `address-pool` - IP pool name (string)
- `rate-limit` - Rate limit format "rx/tx" (string)
- `shared-users` - Number of shared users (string)
- `status-autorefresh` - Autorefresh interval (string)
- `on-login` - RouterOS script yang dijalankan saat login (string)
- `parent-queue` - Parent queue name (string)

#### Update Profile
```php
// File: post/post_update_userprofile.php
$add = $API->comm("/ip/hotspot/user/profile/set", array(
    ".id" => "$upid",
    "name" => "$name",
    "address-pool" => "$addrpool",
    "rate-limit" => "$ratelimit",
    "shared-users" => "$sharedusers",
    "status-autorefresh" => "1m",
    "on-login" => "$onlogin",
    "parent-queue" => "$parent"
));
```

#### Remove Profile
```php
// File: post/post_hotspot_remove.php
$API->comm("/ip/hotspot/user/profile/remove", array(
    ".id" => $id
));
```

### 4.3 Hotspot Active Users (3 commands)

#### Get All Active Users
```php
// File: get/get_hotspot_active.php
$get_hotspot_active = $API->comm("/ip/hotspot/active/print");
```

**Response Properties:**
- `.id` - Unique identifier
- `user` - Username
- `address` - IP address
- `mac-address` - MAC address
- `login-by` - Login method
- `uptime` - Session uptime
- `bytes-in` - Bytes received
- `bytes-out` - Bytes sent
- `server` - Hotspot server

#### Get Active Users Count
```php
// File: get/get_dashboard.php
$get_hotspotactive = $API->comm("/ip/hotspot/active/print", array(
    "count-only" => ""
));
```

#### Remove Active User (Kick)
```php
// File: post/post_hotspot_remove.php
$API->comm("/ip/hotspot/active/remove", array(
    ".id" => $id
));
```

### 4.4 Hotspot Hosts (2 commands)

#### Get All Hosts
```php
// File: get/get_hosts.php
$get_hots_active = $API->comm("/ip/hotspot/host/print");
```

**Response Properties:**
- `.id` - Unique identifier
- `mac-address` - MAC address
- `address` - IP address
- `to-address` - To address
- `server` - Hotspot server
- `authorized` - Authorization status
- `bypassed` - Bypass status
- `found-by` - Discovery method

#### Remove Host
```php
// File: post/post_hotspot_remove.php
$API->comm("/ip/hotspot/host/remove", array(
    ".id" => $id
));
```

### 4.5 Hotspot Server (1 command)

#### Get All Hotspot Servers
```php
// File: get/get_hotspot_server.php
$get_hotspot_server = $API->comm("/ip/hotspot/print");
```

**Response Properties:**
- `.id` - Unique identifier
- `name` - Server name
- `interface` - Interface
- `address-pool` - Address pool
- `profile` - Default profile
- `addresses-per-mac` - Max addresses per MAC
- `idle-timeout` - Idle timeout
- `keepalive-timeout` - Keepalive timeout

### 4.6 System Information (5 commands)

#### Get System Resource
```php
// File: get/get_dashboard.php
$get_resource = $API->comm("/system/resource/print")[0];
```

**Response Properties:**
- `uptime` - System uptime
- `version` - RouterOS version
- `build-time` - Build time
- `factory-software` - Factory software version
- `free-memory` - Free memory
- `total-memory` - Total memory
- `cpu` - CPU model
- `cpu-count` - Number of CPUs
- `cpu-frequency` - CPU frequency
- `cpu-load` - Current CPU load
- `free-hdd-space` - Free HDD space
- `total-hdd-space` - Total HDD space
- `architecture-name` - Architecture
- `board-name` - Board name

#### Get System Clock
```php
// File: get/get_dashboard.php
$get_systime = $API->comm("/system/clock/print")[0];
```

**Response Properties:**
- `time` - Current time
- `date` - Current date
- `time-zone-name` - Timezone name
- `gmt-offset` - GMT offset

#### Get Routerboard Info
```php
// File: get/get_dashboard.php
$get_routerboard = $API->comm("/system/routerboard/print")[0];
```

**Response Properties:**
- `model` - Router model
- `serial-number` - Serial number
- `firmware-type` - Firmware type
- `factory-firmware` - Factory firmware
- `current-firmware` - Current firmware
- `upgrade-firmware` - Upgrade firmware

#### Get System Identity
```php
// File: get/get_dashboard.php
$get_sysidentity = $API->comm("/system/identity/print")[0];
```

**Response Properties:**
- `name` - System identity name

#### Get System Health
```php
// File: get/get_dashboard.php
$get_syshealth = $API->comm("/system/health/print")[0];
```

**Response Properties:**
- `voltage` - Voltage
- `temperature` - Temperature

### 4.7 Interface & Traffic (2 commands)

#### Get All Interfaces
```php
// File: get/get_interface.php
$get_interface = $API->comm("/interface/print");
```

**Response Properties:**
- `.id` - Unique identifier
- `name` - Interface name
- `type` - Interface type
- `mtu` - MTU
- `actual-mtu` - Actual MTU
- `l2mtu` - L2MTU
- `mac-address` - MAC address
- `running` - Running status (true/false)
- `disabled` - Disabled status (true/false)

#### Monitor Traffic
```php
// File: get/get_traffic.php, get/get_dashboard.php
$get_interfacetraffic = $API->comm("/interface/monitor-traffic", array(
    "interface" => "$iface",
    "once" => ""
));
```

**Parameters:**
- `interface` - Interface name (string)
- `once` - Monitor once (empty string)

**Response Properties:**
- `name` - Interface name
- `rx-bits-per-second` - RX bits per second
- `tx-bits-per-second` - TX bits per second
- `rx-packets-per-second` - RX packets per second
- `tx-packets-per-second` - TX packets per second
- `rx-drops-per-second` - RX drops per second
- `tx-drops-per-second` - TX drops per second
- `rx-errors-per-second` - RX errors per second
- `tx-errors-per-second` - TX errors per second

### 4.8 Address Pool (1 command)

#### Get All Address Pools
```php
// File: get/get_addr_pool.php
$get_pool = $API->comm("/ip/pool/print");
```

**Response Properties:**
- `.id` - Unique identifier
- `name` - Pool name
- `ranges` - IP ranges
- `next-pool` - Next pool

### 4.9 NAT Rules (1 command)

#### Get All NAT Rules
```php
// File: get/get_nat.php
$get_hots_active = $API->comm("/ip/firewall/nat/print");
```

**Response Properties:**
- `.id` - Unique identifier
- `chain` - Chain (srcnat/dstnat)
- `action` - Action (masquerade/accept/drop/etc)
- `src-address` - Source address
- `dst-address` - Destination address
- `protocol` - Protocol
- `src-port` - Source port
- `dst-port` - Destination port
- `to-addresses` - To addresses
- `to-ports` - To ports
- `comment` - Comment
- `disabled` - Disabled status

### 4.10 Parent Queue (1 command)

#### Get All Simple Queues
```php
// File: get/get_parent_queue.php
$get_allqueue = $API->comm("/queue/simple/print", array(
    "?dynamic" => "false"
));
```

**Parameters:**
- `?dynamic` - Filter dynamic queues ("false" untuk static queues)

**Response Properties:**
- `.id` - Unique identifier
- `name` - Queue name
- `target` - Target addresses
- `parent` - Parent queue
- `packet-marks` - Packet marks
- `priority` - Priority
- `queue` - Queue type
- `limit-at` - Limit at
- `max-limit` - Max limit
- `burst-limit` - Burst limit
- `burst-threshold` - Burst threshold
- `burst-time` - Burst time
- `comment` - Comment
- `disabled` - Disabled status

### 4.11 System Logging (3 commands)

#### Get Logging Configuration
```php
// File: get/get_dashboard.php
$getlogging = $API->comm("/system/logging/print", array(
    "?prefix" => "->"
));
```

**Response Properties:**
- `.id` - Unique identifier
- `topics` - Log topics
- `action` - Log action
- `prefix` - Log prefix

#### Add Logging Rule
```php
// File: get/get_dashboard.php
$API->comm("/system/logging/add", array(
    "action" => "disk",
    "prefix" => "->",
    "topics" => "hotspot,info,debug"
));
```

#### Get Logs
```php
// File: get/get_dashboard.php
$get_log = $API->comm("/log/print", array(
    "?topics" => "hotspot, info, debug"
));
```

**Response Properties:**
- `.id` - Unique identifier
- `time` - Log time
- `topics` - Log topics
- `message` - Log message

### 4.12 System Scheduler (4 commands)

#### Get Scheduler
```php
// File: post/post_expire_monitor.php, get/get_expire_mon.php
$get_expire_mon = $API->comm("/system/scheduler/print", array(
    "?name" => "Mikhmon-Expire-Monitor"
));

// File: get/get_expire_mon.php (check if enabled)
$get_expire_mon = $API->comm("/system/scheduler/print", array(
    "?name" => "Mikhmon-Expire-Monitor", 
    "?disabled" => "false"
));
```

**Response Properties:**
- `.id` - Unique identifier
- `name` - Scheduler name
- `start-time` - Start time
- `interval` - Interval
- `on-event` - Script to execute
- `disabled` - Disabled status
- `comment` - Comment

#### Add Scheduler
```php
// File: post/post_expire_monitor.php
$expmon = $API->comm("/system/scheduler/add", array(
    "name" => "Mikhmon-Expire-Monitor",
    "start-time" => "00:00:00",
    "interval" => "00:01:00",
    "on-event" => "$expire_monitor_src",
    "disabled" => "no",
    "comment" => "Mikhmon Expire Monitor"
));
```

#### Update Scheduler
```php
// File: post/post_expire_monitor.php
$expmon = $API->comm("/system/scheduler/set", array(
    ".id" => "$id",
    "interval" => "00:01:00",
    "on-event" => "$expire_monitor_src",
    "disabled" => "no"
));
```

### 4.13 System Scripts (4 commands)

#### Get Scripts by Owner
```php
// File: get/get_dashboard.php
$get_report = $API->comm("/system/script/print", array(
    "?owner" => "$month"
));
```

#### Get Scripts by Source (Date)
```php
// File: get/get_report.php
$get_report = $API->comm("/system/script/print", array(
    "?source" => "$day"
));
```

#### Count Scripts by Owner
```php
// File: get/get_dashboard.php
$get_tot_report = $API->comm("/system/script/print", array(
    "?owner" => "$month",
    "count-only" => ""
));
```

#### Count Scripts by Source
```php
// File: get/get_report.php
$get_tot_report = $API->comm("/system/script/print", array(
    "?source" => "$day",
    "count-only" => ""
));
```

**Response Properties:**
- `.id` - Unique identifier
- `name` - Script name (format: date-|-time-|-user-|-price-|-ip-|-mac-|-validity-|-profile-|-comment)
- `owner` - Owner (format: monthyear, e.g., "jan2024")
- `source` - Source (date)
- `comment` - Comment ("mikhmon")

---

## 5. Script RouterOS (MikroTik)

### 5.1 On-Login Script (User Profile)

Script ini dijalankan ketika user login ke hotspot. Script ini otomatis di-generate oleh PHP saat membuat/mengupdate profile.

**Lokasi:** `post/post_add_userprofile.php` dan `post/post_update_userprofile.php`

#### Struktur Script:

```routeros
# Output debug info
:put (",rem,1000,30d,1500,,Enable,Disable,");

# Set mode (N=Notify, X=Remove)
:local mode "X";

{
    # Get current date
    :local date [ /system clock get date ];
    :local year [ :pick $date 7 11 ];
    :local month [ :pick $date 0 3 ];
    
    # Get user comment
    :local comment [ /ip hotspot user get [/ip hotspot user find where name="$user"] comment];
    :local ucode [:pic $comment 0 2];
    
    # Check if user has code prefix (vc- atau up-)
    :if ($ucode = "vc" or $ucode = "up" or $comment = "") do={
        # Create temporary scheduler untuk menghitung expire date
        /sys sch add name="$user" disable=no start-date=$date interval="30d";
        :delay 2s;
        
        # Get next-run (expire date)
        :local exp [ /sys sch get [ /sys sch find where name="$user" ] next-run];
        :local getxp [len $exp];
        
        # Format expire date berdasarkan length
        :if ($getxp = 15) do={
            # Format: jan/01/2024 12:00:00
            :local d [:pic $exp 0 6];
            :local t [:pic $exp 7 16];
            :local s ("/");
            :local exp ("$d$s$year $t");
            /ip hotspot user set comment="$exp $mode" [find where name="$user"];
        };
        
        :if ($getxp = 8) do={
            # Format: 12:00:00
            /ip hotspot user set comment="$date $exp $mode" [find where name="$user"];
        };
        
        :if ($getxp > 15) do={
            # Format lainnya
            /ip hotspot user set comment="$exp $mode" [find where name="$user"];
        };
        
        # Remove temporary scheduler
        /sys sch remove [find where name="$user"];
    };
};

# MAC Address Locking (jika Enable)
[:local mac "$mac-address"; /ip hotspot user set mac-address=$mac [find where name=$user]];

# Server Locking (jika Enable)
[:local mac "$mac-address"; :local srv [/ip hotspot host get [find where mac-address="$mac"] server]; /ip hotspot user set server=$srv [find where name=$user]]
```

#### Mode Expire:

1. **rem** (Remove): Hapus user saat expired
   ```routeros
   :local mode "X";
   ```

2. **ntf** (Notify): Disable user saat expired (set limit-uptime=1s)
   ```routeros
   :local mode "N";
   ```

3. **remc** (Remove + Comment): Hapus user + record ke report
   ```routeros
   :local mode "X";
   # + recording script
   ```

4. **ntfc** (Notify + Comment): Disable user + record ke report
   ```routeros
   :local mode "N";
   # + recording script
   ```

5. **0** (No Expire): Tidak ada expire
   ```routeros
   :put (",,1000,,1500,noexp,Enable,Disable,");
   ```

#### Recording Script (untuk report):
```routeros
:local mac "$mac-address";
:local time [/system clock get time ];
/system script add name="$date-|-$time-|-$user-|-1000-|-$address-|-$mac-|-30d-|-ProfileName-|-$comment" owner="$month$year" source=$date comment=mikhmon
```

### 5.2 Expire Monitor Script

Script ini berjalan setiap 1 menit via System Scheduler untuk memonitor user yang expired.

**Lokasi:** `assets/js/func.js` (function `setExpMon()`) dan `post/post_expire_monitor.php`

#### Full Script:

```routeros
# Function untuk convert date ke integer format YYYYMMDD
:local dateint do={
    :local montharray ( "jan","feb","mar","apr","may","jun","jul","aug","sep","oct","nov","dec" );
    :local days [ :pick $d 4 6 ];
    :local month [ :pick $d 0 3 ];
    :local year [ :pick $d 7 11 ];
    :local monthint ([ :find $montharray $month]);
    :local month ($monthint + 1);
    :if ( [len $month] = 1) do={
        :local zero ("0");
        :return [:tonum ("$year$zero$month$days")];
    } else={
        :return [:tonum ("$year$month$days")];
    }
};

# Function untuk convert time ke menit
:local timeint do={
    :local hours [ :pick $t 0 2 ];
    :local minutes [ :pick $t 3 5 ];
    :return ($hours * 60 + $minutes);
};

# Get current date dan time
:local date [ /system clock get date ];
:local time [ /system clock get time ];
:local today [$dateint d=$date];
:local curtime [$timeint t=$time];

# Get tahun sekarang dan tahun lalu
:local tyear [ :pick $date 7 11 ];
:local lyear ($tyear-1);

# Loop semua user dengan comment yang mengandung tahun
:foreach i in [ /ip hotspot user find where comment~"/$tyear" || comment~"/$lyear" ] do={
    :local comment [ /ip hotspot user get $i comment];
    :local limit [ /ip hotspot user get $i limit-uptime];
    :local name [ /ip hotspot user get $i name];
    :local gettime [:pic $comment 12 20];
    
    # Check format comment (harus ada / di posisi 3 dan 6)
    :if ([:pic $comment 3] = "/" and [:pic $comment 6] = "/") do={
        :local expd [$dateint d=$comment];
        :local expt [$timeint t=$gettime];
        
        # Check kondisi expired
        :if (($expd < $today and $expt < $curtime) or 
              ($expd < $today and $expt > $curtime) or 
              ($expd = $today and $expt < $curtime) and $limit != "00:00:01") do={
            
            # Mode N = Notify (disable user)
            :if ([:pic $comment 21] = "N") do={
                [ /ip hotspot user set limit-uptime=1s $i ];
                [ /ip hotspot active remove [find where user=$name] ];
            } else={
                # Mode X = Remove (hapus user)
                [ /ip hotspot user remove $i ];
                [ /ip hotspot active remove [find where user=$name] ];
            }
        }
    }
}
```

#### Alur Kerja Expire Monitor:

1. **Scheduler** berjalan setiap 1 menit (interval "00:01:00")
2. **Convert date** dari format "jan/01/2024" ke integer "20240101"
3. **Convert time** dari format "12:00:00" ke menit "720"
4. **Loop** semua user dengan comment yang mengandung tahun sekarang atau tahun lalu
5. **Parse comment** untuk mendapatkan expire date
6. **Compare** expire date dengan current date
7. **Action** berdasarkan mode:
   - Mode "N" (Notify): Set limit-uptime=1s dan kick user
   - Mode "X" (Remove): Hapus user dan kick active

#### RouterOS Commands dalam Expire Monitor:

| Command | Fungsi |
|---------|--------|
| `:local` | Deklarasi variabel lokal |
| `:pick` | Mengambil substring |
| `:find` | Mencari index dalam array |
| `:len` | Mendapatkan panjang string |
| `:tonum` | Convert ke number |
| `:if ... do={...}` | Conditional statement |
| `:foreach ... in ... do={...}` | Loop melalui array |
| `:return` | Return value dari function |
| `:pic` | Alias untuk `:pick` |
| `:put` | Output debug |
| `:delay` | Delay execution |
| `\|\|` | OR operator |
| `\|`, `\|` | String concatenation |

### 5.3 MAC Address Lock Script

**Lokasi:** `post/post_add_userprofile.php`

```routeros
# Lock MAC Address (jika Enable)
[:local mac "$mac-address"; /ip hotspot user set mac-address=$mac [find where name=$user]]
```

Fungsi: Mengunci user ke MAC address device yang pertama kali login.

### 5.4 Server Lock Script

**Lokasi:** `post/post_add_userprofile.php`

```routeros
# Lock Server (jika Enable)
[:local mac "$mac-address"; 
 :local srv [/ip hotspot host get [find where mac-address="$mac"] server]; 
 /ip hotspot user set server=$srv [find where name=$user]]
```

Fungsi: Mengunci user ke hotspot server tertentu berdasarkan MAC address.

### 5.5 Daftar Lengkap RouterOS Script Commands

#### Variabel dan Data:
```routeros
:local nama_variabel "nilai"           # Deklarasi variabel lokal
:global nama_variabel "nilai"          # Deklarasi variabel global
:set nama_variabel "nilai_baru"        # Set nilai variabel
```

#### String Manipulation:
```routeros
:pick $string 0 5                      # Ambil substring (start, end)
:len $string                           # Panjang string
:find $array "value"                   # Cari index dalam array
:tonum "123"                           # Convert ke number
:tostr 123                             # Convert ke string
```

#### Control Flow:
```routeros
:if (kondisi) do={
    # code
}

:if (kondisi) do={
    # code
} else={
    # code
}

:foreach i in=$array do={
    # code menggunakan $i
}

:while (kondisi) do={
    # code
}

:for i from=0 to=10 do={
    # code
}
```

#### Operators:
```routeros
=    # Assignment
==   # Equal
!=   # Not equal
<    # Less than
>    # Greater than
<=   # Less than or equal
>=   # Greater than or equal
~>   # Regex match
~>*  # Regex match (case insensitive)
\|    # OR
\|\|  # OR (shortcut)
&    # AND
&&   # AND (shortcut)
!    # NOT
```

#### System Commands:
```routeros
/system clock get date                 # Get current date
/system clock get time                 # Get current time
/system clock get time-zone-name       # Get timezone

/system resource get uptime            # Get uptime
/system resource get version           # Get RouterOS version
/system resource get cpu-load          # Get CPU load

/system identity get name              # Get system identity

/system scheduler add ...              # Add scheduler
/system scheduler set ...              # Update scheduler
/system scheduler remove ...           # Remove scheduler
/system scheduler print                # List schedulers

/system script add ...                 # Add script
/system script set ...                 # Update script
/system script remove ...              # Remove script
/system script print                   # List scripts

/system logging add ...                # Add logging rule
/system logging print                  # List logging rules

/log print                             # View logs
```

#### Hotspot Commands:
```routeros
/ip hotspot user print                 # List users
/ip hotspot user add ...               # Add user
/ip hotspot user set ...               # Update user
/ip hotspot user remove ...            # Remove user
/ip hotspot user reset-counters ...    # Reset counters
/ip hotspot user find where ...        # Find users
/ip hotspot user get ...               # Get user property

/ip hotspot user profile print         # List profiles
/ip hotspot user profile add ...       # Add profile
/ip hotspot user profile set ...       # Update profile
/ip hotspot user profile remove ...    # Remove profile

/ip hotspot active print               # List active users
/ip hotspot active remove ...          # Kick user

/ip hotspot host print                 # List hosts
/ip hotspot host remove ...            # Remove host

/ip hotspot print                      # List hotspot servers
```

#### Network Commands:
```routeros
/interface print                       # List interfaces
/interface monitor-traffic ...         # Monitor traffic

/ip pool print                         # List address pools
/ip firewall nat print                 # List NAT rules

/queue/simple print                    # List simple queues
/queue/simple add ...                  # Add queue
/queue/simple set ...                  # Update queue
/queue/simple remove ...               # Remove queue
```

---

## 6. JavaScript Functions

### 6.1 Enkripsi/Decoding

#### XOR Decoder (jesD)
File: `assets/js/func.js`

```javascript
var jesD = {
  dec: function (e, t = 25) {
    var a = "";
    for (var s = 0; s < e.length; s++) {
      var o = e.charCodeAt(s) ^ t;
      a += String.fromCharCode(o);
    }
    return a;
  }
};
```

Digunakan untuk decode response dari server PHP yang di-encode dengan XOR key 25.

#### Password Encoding (blah)
File: `assets/js/func.js`

```javascript
function blah(e) {
  var t = "";
  for (e = btoa(e), e = btoa(e), i = 0; i < e.length; i++) {
    var a = 10 ^ e.charCodeAt(i);
    t += String.fromCharCode(a);
  }
  return (t = btoa(t));
}
```

Proses encoding:
1. Base64 encode
2. Base64 encode lagi
3. XOR setiap karakter dengan 10
4. Base64 encode hasil XOR

#### Password Decoding (unblah)
File: `assets/js/func.js`

```javascript
function unblah(encoded) {
  try {
    // Step 1: Decode base64 pertama
    let decoded = atob(encoded);
    
    // Step 2: XOR setiap karakter dengan 10
    let xorResult = "";
    for (let i = 0; i < decoded.length; i++) {
      xorResult += String.fromCharCode(decoded.charCodeAt(i) ^ 10);
    }
    
    // Step 3: Decode base64 dua kali
    let original = atob(atob(xorResult));
    
    return original;
  } catch (e) {
    console.error("Gagal decode:", e);
    return null;
  }
}
```

### 6.2 User Generator Functions

File: `core/generator_functions.php`

```php
// Random Lowercase - huruf kecil saja
function randLC($length) {
    $chars = "abcdefghijklmnopqrstuvwxyz";
    return substr(str_shuffle($chars), 0, $length);
}

// Random Uppercase - huruf besar saja
function randUC($length) {
    $chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ";
    return substr(str_shuffle($chars), 0, $length);
}

// Random Uppercase + Lowercase - huruf besar dan kecil
function randULC($length) {
    $chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ";
    return substr(str_shuffle($chars), 0, $length);
}

// Random Numeric + Lowercase - angka dan huruf kecil
function randNLC($length) {
    $chars = "abcdefghijklmnopqrstuvwxyz0123456789";
    return substr(str_shuffle($chars), 0, $length);
}

// Random Numeric + Uppercase - angka dan huruf besar
function randNUC($length) {
    $chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
    return substr(str_shuffle($chars), 0, $length);
}

// Random Numeric + Uppercase + Lowercase - semua karakter
function randNULC($length) {
    $chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
    return substr(str_shuffle($chars), 0, $length);
}

// Random Numeric - angka saja
function randN($length) {
    $chars = "0123456789";
    return substr(str_shuffle($chars), 0, $length);
}
```

### 6.3 Expire Monitor Setup (setExpMon)

File: `assets/js/func.js`

```javascript
function setExpMon() {
  var e = session;
  "?admin" !== session &&
    $.ajax({
      type: "POST",
      url: "./post/post_expire_monitor.php",
      data: {
        sessname: e,
        expmon: ':local dateint do={...}'  // Script RouterOS lengkap
      },
      dataType: "json",
      success: function (e) {
        "success" == e.message &&
          $("#exp_mon").html(
            '<i class="fa fa-ci fa-circle text-green" title="Expire users monitor is activated."></i>'
          ),
          $("#btn-exp-mon").fadeOut(200));
      }
    });
}
```

Fungsi: Mengirim script expire monitor ke MikroTik untuk dijadikan System Scheduler.

### 6.4 User Management Functions

#### Open User Detail
```javascript
function openUser(e, t) {
  // e = user ID atau name
  // t = type ('uname' untuk name, lainnya untuk ID)
  $("#usract-btn").show();
  upath = "uname" == t ? "/user&name=" : "/user&id=";
  
  $.get(session + upath + e, function (e) {
    var t = jesD.dec(e);  // Decode response
    // Parse dan tampilkan data user
  });
}
```

#### Add/Update User
```javascript
function addUser(e = "no") {
  // e = "yes" jika reset, "no" jika normal update
  
  // Convert data limit
  "M" == s.substr(-1, 1)
    ? (adlimit = 1048576 * Number(s.substr(0, s.length - 1)))
    : "G" == s.substr(-1, 1)
    ? (adlimit = 1073741824 * Number(s.substr(0, s.length - 1)))
    : (adlimit = "" == s ? 0 : s);
  
  // Set MAC address default
  l = "" == $("#add_usrmac").val()
    ? "00:00:00:00:00:00"
    : $("#add_usrmac").val();
  
  // Kirim ke server
  $.ajax({
    type: "POST",
    url: t,  // post_add_user.php atau post_update_user.php
    data: {
      sessname: a,
      uid: $("#add_usrid").val(),
      server: $("#add_hserver").val(),
      name: $("#add_usrname").val(),
      password: $("#add_usrpass").val(),
      profile: $("#add_usrprofile").val(),
      macaddr: l,
      timelimit: atlimit,
      datalimit: adlimit,
      comment: n,
      expdate: i,
      reset: e,
      ucode: $("#ucode").val()
    },
    // ...
  });
}
```

#### Remove User
```javascript
function removeUser(e, t, a) {
  // e = user ID
  // t = username
  // a = profile name
  
  if (confirm("Are you sure want to remove user [" + t + "] ?")) {
    $.ajax({
      type: "POST",
      url: "./post/post_hotspot_remove.php",
      data: { sessname: s, where: "user_", id: e },
      // ...
    });
  }
}
```

### 6.5 Voucher Generator Functions

#### Generate Voucher
```javascript
function generateV() {
  $("#btn-genV").attr("onclick", "");
  localStorage.setItem(session + "_gencode", randomN(101, 999));
  var e = $("#gen_usrqty").val();
  localStorage.setItem(session + "_totgenv", e);
  localStorage.setItem(session + "_genleft", e);
  genVoucher();
}
```

#### Process Voucher Generation
```javascript
function genVoucher() {
  var e = localStorage.getItem(session + "_genleft");
  qty = e > 50 ? 50 : e;  // Max 50 per batch
  
  $.ajax({
    type: "POST",
    url: "./post/post_generate_voucher.php",
    data: {
      sessname: t,
      qty: qty,
      server: $("#gen_hserver").val(),
      user: $("#gen_usrmode").val(),  // "up" atau "vc"
      userl: $("#gen_namelength").val(),
      prefix: $("#gen_usrprefix").val(),
      char: $("#gen_char").val(),
      profile: $("#gen_usrprofile").val(),
      timelimit: $("#gen_tlimit").val(),
      datalimit: $("#gen_dlimit").val(),
      gcomment: $("#gen_usrcomm").val(),
      gencode: localStorage.getItem(session + "_gencode")
    },
    // ...
  });
}
```

### 6.6 Traffic Monitoring

```javascript
$("#iface-name").change(function () {
  var e = $("#iface-name").val();
  localStorage.setItem(session + "_iface", e);
  
  // Clear existing intervals
  clearInterval(window.ifaceI);
  clearInterval(window.ifaceII);
  
  // Get traffic data
  $.getJSON(session + "/get_traffic/&iface=" + e, function (t) {
    txrx = '[{"name":"Tx","data":["' + t.tx + '"]},{"name":"Rx","data":["' + t.rx + '"]}]';
    localStorage.setItem(session + "_traffic_data", txrx);
  }).always(function () {
    trafficMonitor(localStorage.getItem(session + "_theme"));
  });
});
```

---

## 7. Data Flow

### 7.1 Login Flow

1. User mengisi form login (IP, Username, Password)
2. JavaScript `blah()` meng-encode password
3. AJAX ke `post/post_a_router.php` dengan action `save`
4. PHP menyimpan config ke `config/config.php`
5. Test koneksi ke MikroTik via `get/get_connect.php`
6. Jika sukses, redirect ke dashboard

### 7.2 Add User Flow

1. User mengisi form add user
2. JavaScript mengkonversi data limit (M/G ke bytes)
3. AJAX ke `post/post_add_user.php`
4. PHP decode password dengan `dec_rypt()`
5. PHP connect ke MikroTik via RouterOS API
6. PHP kirim command `/ip/hotspot/user/add`
7. Jika sukses, get user data dengan `/ip/hotspot/user/print`
8. Return JSON ke JavaScript
9. JavaScript decode dengan `jesD.dec()`
10. Update UI dengan data user

### 7.3 Generate Voucher Flow

1. User setting voucher parameter (qty, profile, pattern, dll)
2. JavaScript `generateV()` set kode unik
3. `genVoucher()` kirim AJAX ke `post/post_generate_voucher.php`
4. PHP generate username/password sesuai pattern
5. PHP kirim command `/ip/hotspot/user/add` untuk setiap voucher
6. Return count dan comment identifier
7. JavaScript update progress
8. Jika masih ada sisa qty, repeat step 3-7
9. `genCacheVoucher()` cache data voucher
10. `cacheUser()` reload user list

### 7.4 Expire Monitor Flow

1. User klik "Enable Expire Monitor"
2. JavaScript `setExpMon()` kirim script ke `post/post_expire_monitor.php`
3. PHP check existing scheduler dengan `/system/scheduler/print`
4. Jika belum ada, create dengan `/system/scheduler/add`
5. Jika sudah ada tapi disabled, enable dengan `/system/scheduler/set`
6. Scheduler berjalan setiap 1 menit
7. Script check user expired
8. Action: Remove atau Disable user
9. Kick active user jika ada

### 7.5 Report Flow

1. On-login script di profile menyimpan data ke `/system/script/add`
2. Script name format: `date-|-time-|-user-|-price-|-ip-|-mac-|-validity-|-profile-|-comment`
3. Script owner format: `monthyear` (e.g., "jan2024")
4. JavaScript request ke `get/get_dashboard.php` dengan `page=get_livereport`
5. PHP get scripts dengan `/system/script/print` filter by owner
6. PHP encode dan return ke JavaScript
7. JavaScript decode dan tampilkan di table

### 7.6 Print Voucher Flow

1. User klik print voucher
2. JavaScript `pVcr()` buka window baru ke `view/print_voucher.php`
3. PHP get voucher data dengan `/ip/hotspot/user/print` filter by comment
4. PHP get profile info dengan `/ip/hotspot/user/profile/print`
5. PHP parse template file dari `template/` folder
6. PHP generate QR code dengan qrious.js
7. PHP render HTML dengan data voucher
8. Browser auto-print saat window load

---

## 8. Security

### 8.1 PHP Protection

```php
// Protect direct access to PHP files
$get_self = explode("/",$_SERVER['PHP_SELF']);
$self[] = $get_self[count($get_self)-1];

if($self[0] !== "index.php"  && $self[0] !==""){
    include_once("../core/route.php");
}
```

### 8.2 Session Validation

```php
if(isset($_POST['sessname']) && isset($_SESSION["mikhmon"])){
    // Process request
}
```

### 8.3 Data Encoding

1. **Password Storage**: XOR + Base64 double encoding
2. **API Response**: XOR encoding dengan key 25
3. **JavaScript**: `jesD.dec()` untuk decode response

### 8.4 Input Sanitization

```php
$name = (preg_replace('/\s+/', '-',$_POST['name']));
$validity = strtolower($_POST['validity']);
```

---

## 9. Utility Functions

### 9.1 Format Bytes

```php
function formatBytes($size, $decimals = 0) {
    $unit = array(
        '0' => 'Byte',
        '1' => 'KiB',
        '2' => 'MiB',
        '3' => 'GiB',
        '4' => 'TiB',
        '5' => 'PiB',
        '6' => 'EiB',
        '7' => 'ZiB',
        '8' => 'YiB'
    );
    
    for ($i = 0; $size >= 1024 && $i <= count($unit); $i++) {
        $size = $size / 1024;
    }
    
    return round($size, $decimals) . ' ' . $unit[$i];
}
```

### 9.2 Format Bytes2 (Decimal)

```php
function formatBytes2($size, $decimals = 0) {
    $unit = array(
        '0' => 'Byte',
        '1' => 'KB',
        '2' => 'MB',
        '3' => 'GB',
        '4' => 'TB',
        '5' => 'PB',
        '6' => 'EB',
        '7' => 'ZB',
        '8' => 'YB'
    );
    
    for ($i = 0; $size >= 1000 && $i <= count($unit); $i++) {
        $size = $size / 1000;
    }
    
    return round($size, $decimals) . '' . $unit[$i];
}
```

### 9.3 Format Bits (Bandwidth)

```php
function formatBites($size, $decimals = 0) {
    $unit = array(
        '0' => 'bps',
        '1' => 'kbps',
        '2' => 'Mbps',
        '3' => 'Gbps',
        '4' => 'Tbps',
        '5' => 'Pbps',
        '6' => 'Ebps',
        '7' => 'Zbps',
        '8' => 'Ybps'
    );
    
    for ($i = 0; $size >= 1000 && $i <= count($unit); $i++) {
        $size = $size / 1000;
    }
    
    return round($size, $decimals) . ' ' . $unit[$i];
}
```

### 9.4 Get Config

```php
function get_config($string, $start, $end) {
    $string = ' ' . $string;
    $ini = strpos($string, $start);
    if ($ini == 0) return '';
    $ini += strlen($start);
    $len = strpos($string, $end, $ini) - $ini;
    return substr($string, $ini, $len);
}
```

---

## 10. API Endpoints

### 10.1 GET Endpoints

| Endpoint | File | Parameter | Deskripsi |
|----------|------|-----------|-----------|
| `/get_users` | `get/get_users.php` | `prof`, `f`, `c` | Get users by profile |
| `/get_user` | `get/get_user.php` | `id` atau `name` | Get single user |
| `/get_profiles` | `get/get_profiles.php` | `f` | Get all profiles |
| `/get_profile` | `get/get_profile.php` | `id` atau `name` | Get single profile |
| `/get_hotspot_active` | `get/get_hotspot_active.php` | - | Get active users |
| `/get_hotspot_server` | `get/get_hotspot_server.php` | `f` | Get hotspot servers |
| `/get_traffic` | `get/get_traffic.php` | `iface` | Get interface traffic |
| `/get_interface` | `get/get_interface.php` | - | Get interfaces |
| `/get_dashboard` | `get/get_dashboard.php` | `page` | Dashboard data |
| `/get_hosts` | `get/get_hosts.php` | - | Get hotspot hosts |
| `/get_addr_pool` | `get/get_addr_pool.php` | `f` | Get address pools |
| `/get_parent_queue` | `get/get_parent_queue.php` | `f` | Get parent queues |
| `/get_nat` | `get/get_nat.php` | - | Get NAT rules |
| `/get_report` | `get/get_report.php` | `day`, `f` | Get reports by date |
| `/get_tot_users` | `get/get_tot_users.php` | `name` | Get total users |
| `/get_connect` | `get/get_connect.php` | - | Test connection |
| `/get_sys_resource` | `get/get_dashboard.php` | - | Get system resource |
| `/get_hotspotinfo` | `get/get_dashboard.php` | - | Get hotspot info |
| `/get_log` | `get/get_dashboard.php` | `f` | Get logs |
| `/get_livereport` | `get/get_dashboard.php` | `day`, `month`, `f` | Get live report |
| `/get_expire_mon` | `get/get_expire_mon.php` | - | Check expire monitor status |

### 10.2 POST Endpoints

| Endpoint | File | Parameter | Deskripsi |
|----------|------|-----------|-----------|
| `/post_add_user` | `post/post_add_user.php` | `sessname`, `server`, `name`, `password`, `profile`, `macaddr`, `timelimit`, `datalimit`, `comment` | Add user |
| `/post_update_user` | `post/post_update_user.php` | `sessname`, `uid`, `server`, `name`, `password`, `profile`, `macaddr`, `timelimit`, `datalimit`, `comment`, `expdate`, `reset`, `ucode` | Update user |
| `/post_add_userprofile` | `post/post_add_userprofile.php` | `sessname`, `name`, `addresspool`, `sharedusers`, `ratelimit`, `parentqueue`, `expmode`, `validity`, `price`, `sellingprice`, `lockuser`, `lockserver` | Add profile |
| `/post_update_userprofile` | `post/post_update_userprofile.php` | Similar to add | Update profile |
| `/post_hotspot_remove` | `post/post_hotspot_remove.php` | `sessname`, `where`, `id` | Remove user/profile/host/active |
| `/post_generate_voucher` | `post/post_generate_voucher.php` | `sessname`, `qty`, `server`, `user`, `userl`, `prefix`, `char`, `profile`, `timelimit`, `datalimit`, `gcomment`, `gencode` | Generate vouchers |
| `/post_cache_voucher` | `post/post_cache_voucher.php` | `sessname`, `qty`, `user`, `gcomment`, `gencode` | Cache voucher |
| `/post_expire_monitor` | `post/post_expire_monitor.php` | `sessname`, `expmon` | Setup expire monitor |
| `/post_a_router` | `post/post_a_router.php` | `do`, `router_`, ... | Router management |
| `/post_logout` | `post/post_logout.php` | `logout` | Logout |
| `/post_template` | `post/post_template.php` | - | Template management |

---

## 11. Format Data

### 11.1 Comment Format

Format comment untuk user hotspot:

```
[DATE] [MODE]
Contoh: apr/01/2024 12:00:00 N

Mode:
- N = Notify (notifikasi saja, disable user)
- X = Remove (hapus user)
```

### 11.2 User Code Format

```
vc-[CODE]-[DATE]  -> Voucher (username = password)
up-[CODE]-[DATE]  -> User/Password (username != password)
```

Contoh:
- `vc-123-01.15.24-Paket1` - Voucher dengan kode 123
- `up-456-01.15.24-Paket2` - User/Password dengan kode 456

### 11.3 Report Script Name Format

```
[DATE]-|-[TIME]-|-[USER]-|-[PRICE]-|-[IP]-|-[MAC]-|-[VALIDITY]-|-[PROFILE]-|-[COMMENT]
```

Contoh:
```
jan/01/2024-|-12:00:00-|-user123-|-10000-|-192.168.1.100-|-00:11:22:33:44:55-|-30d-|-Paket1-|-vc-123-01.01.24
```

### 11.4 Profile On-Login Format

```
:put (",[MODE],[PRICE],[VALIDITY],[SELLING_PRICE],[NOEXP],[LOCK_USER],[LOCK_SERVER],");
```

Contoh:
```
:put (",rem,10000,30d,15000,,Enable,Disable,");
```

---

## 12. Session Management

### 12.1 PHP Session Variables

```php
$_SESSION['mikhmon']        // Main session flag (boolean)
$_SESSION['m_user']         // Current user/session name (string)
$_SESSION['admin']          // Admin flag (string: "admin")
$_SESSION["$ses"]           // Router sessions (array)
$_SESSION["timezone"]       // MikroTik timezone (string)
$_SESSION["$m_user.$uprof"] // Cached users by profile (string - encoded)
$_SESSION["$m_user.'profiles'"] // Cached profiles (string - encoded)
$_SESSION["$m_user.'addr_pool'"] // Cached address pools (string)
$_SESSION["$m_user.'parentq'"]   // Cached parent queues (string)
$_SESSION["$m_user.'hotspot-server'"] // Cached hotspot servers (string)
$_SESSION["$sessR"]         // Cached reports (string - encoded)
$_SESSION["$m_user.$commt"] // Cached voucher users (array)
```

### 12.2 LocalStorage Keys (JavaScript)

```javascript
// Session-specific
session + "_curr"                    // Currency (string)
session + "_theme"                   // Theme (string)
session + "_iface"                   // Selected interface (string)
session + "_idleto"                  // Idle timeout (string)
session + "_profn"                   // Profile number (string)
session + "_traffic_data"            // Traffic data JSON (string)
session + "_gencode"                 // Generate code (string)
session + "_totgenv"                 // Total generate voucher (string)
session + "_genleft"                 // Generate left (string)
session + "_gencomment"              // Generate comment (string)
session + "_auto_traffic"            // Auto traffic interval (string)
session + "_cache_user_profiles"     // Cached profiles (string)
session + "_temp_hotspot_server"     // Temp servers (string)
session + "_temp_user_profiles"      // Temp profiles (string)
session + "_resume_report"           // Resume report data (string)
session + "_NetInfo"                 // Network info (string)
session + "_force"                   // Force reload flag (string)
session + "AutoReload"               // Auto reload flag (string)

// Global
"typeTemplate"                       // Template type (string)
"nameTemplate"                       // Template name (string)
"tmplPath"                           // Template path (string)
"noR"                                // No reload flag (string)
```

---

## 13. Catatan Penting

### 13.1 Port API
- **Default**: Port 8728 (non-SSL)
- **SSL**: Port 8729 (harus enable api-ssl di IP/Services)

### 13.2 Enkripsi
- Password dienkripsi dengan XOR + Base64 double encoding
- Response API di-encode dengan XOR key 25
- JavaScript decode dengan `jesD.dec()`

### 13.3 Caching
- Data di-cache di session PHP untuk mengurangi beban API
- Cache di-reset dengan parameter `f=true`
- Expire monitor cache: 1 menit

### 13.4 Expire Monitor
- Berjalan setiap 1 menit via System Scheduler
- Script di-generate oleh JavaScript dan dikirim ke PHP
- PHP install script ke MikroTik
- Mode N = Disable user (set limit-uptime=1s)
- Mode X = Remove user (delete)

### 13.5 Report
- Disimpan sebagai System Script di MikroTik
- Format nama script: `date-|-time-|-user-|-price-|-ip-|-mac-|-validity-|-profile-|-comment`
- Owner script: `monthyear` (e.g., "jan2024")
- Source script: tanggal transaksi

### 13.6 Voucher
- Dapat digenerate dalam batch (max 50 per request)
- Pattern: lowercase, uppercase, mixed, numeric, dll
- Prefix dapat ditambahkan
- Comment format: `mode-code-date-comment`

### 13.7 Profile
- Mendukung on-login script untuk auto-expire
- MAC Address Locking (optional)
- Server Locking (optional)
- Rate limiting
- Parent queue

### 13.8 Limitations
- Max 50 voucher per batch generation
- Scheduler expire monitor interval minimal 1 menit
- Report disimpan di MikroTik (bisa penuh jika tidak dihapus)

### 13.9 Total Commands
**59 API Commands** digunakan dalam Mikhmon v4:
- 15 Hotspot User commands
- 7 Hotspot Profile commands
- 3 Hotspot Active commands
- 2 Hotspot Host commands
- 1 Hotspot Server command
- 5 System Information commands
- 2 Interface commands
- 1 Address Pool command
- 1 NAT command
- 1 Queue command
- 3 Logging commands
- 4 Scheduler commands
- 4 Script commands

---

*Dokumentasi ini dibuat berdasarkan analisis kode lengkap Mikhmon v4*
*Total: 59 API Commands, 4 RouterOS Scripts, 7 Generator Functions*
*Terakhir diupdate: 13 Maret 2026*
