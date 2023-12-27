<a name="readme-top"></a>
<h2 align="center">Goshh-Server</h3>

<!-- LOGO -->
<br />
<div align="center">
  <a href="Placeholder">
    <img src="https://github.com/5ur/Goshh/blob/main/logos/server_logo.png" alt="Logo" width="35%" height="35%">
  </a>

<p align="center">
  A Go message and file sharing service
  <br />
  <a href="Placeholder"><strong>Wiki»</strong></a>
</div>


<!-- TOC -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#prerequisites">Prerequisites</a></li>
    <li><a href="#Installation">Installation</a></li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#examples">Examples</a></li>
    <li><a href="#recommendation">Recommendation</a></li>
  </ol>
</details>

### Built With
[![Go][Go]][Go-url] [![Powershell][Powershell]][Powershell-url] [![Vim][Vim]][Vim-url] [![Exchange][StackExchange]][StackExchange-url] [![Overflow][StackOverflow]][StackOverflow-url] [![Windows][Windows]][Windows-url]

# Prerequisites
**Minimum version of go required is 1.16  **
## Windows
Start by installing Go:  
You can download and install Go from: https://go.dev/dl/

Or you and use a package manager:  
**Scoop**: https://scoop.sh/
```Powershell
irm get.scoop.sh | iex
scoop bucket add main
scoop install main/go
```
**Winget**
```Powershell
winget install -e --id GoLang.Go
```

## Linux
**apt**
```Shell
apt install golang -y
```

You can also manually install the bin (eg; apt ins installing an older version of go lesser than 1.16): https://go.dev/dl/  

```Shell
wget https://go.dev/dl/go1.20.4.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.20.4.linux-amd64.tar.gz

```

Depending on your distro you might have one of these two files:  
`vim /etc/profile` or `/etc/bash.bashrc`  
You can add the the Go bins to path in one of those or both to make the binary available for all users:  
```Shell
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
```
You can also just add the above to a specif user's bashrc or profile.  

# Installation
## Windows
1. Clone the repository:
   ```Powershell
   git clone https://github.com/5ur/Goshh-Server.git
   cd .\Goshh-Server\   
   ```
2. Download the external go packages
   ```Powershell
   # If need be (eg; you use a newer/older version of go, delete the go.mod file and create your own):
   go mod init Goshh-Server
   go get .
   go mod tidy
   ```
3. Build the binary
   ```Powershell
   go build .
   ```
4. Create your config file
  ```Powershell
  mv config.yaml.example config.yaml
  vim config.yaml
  ```
  Example config file:
  ```YAML
  # Gin debug mode:
  debugMode: true
  # Port for the engine:
  serverPort: 5150
  # Use default gin router (if false new is created):
  useDefault: true
  # Trusted proxy slice/array
  trustedProxies:
    # Local host only:
    - 127.0.0.1
    # Some RFC 1918 range:
    - 10.0.0.0/8
  # Time in which the slice of messages will be dumped
  cleanupInterval: "10s"
  # This achieves the same as adding all 1918 ranges in trustedProxies, just more convenient
  allowLocalNetworkAccess: true
  # You understand what this does with no context needed
  allowedFileTypes:
    - txt
    - md
    - jpg
  # You understand what this does with no context needed
  fileSavePath: "tmp/"
  # Time after which a stale file will be deleted (ie; not downloaded at all, not downloaded enough times to reach the allowedFileDownloadCount limit)
  staleFileTTL: "30s"
  # The amount of times a file is allowed to be downloaded (kept in check by the file struct values)
  allowedFileDownloadCount: 1 
  ```
5. That's it
Mind the follwoing:  
If the server is started with no configuration file or a missing value/s, it will just start with the default ones hard-coded into the server binary.  
eg:
```Powershell
❯ .\Goshh-Server.exe
2023/05/20 12:03:18 Error reading config file: open config.yaml: The system cannot find the file specified.
2023/05/20 12:03:18 Loading configuration values:
 fileSavePath="/path/to/save/files"
 staleFileTTL=30s
 debugMode=false
 allowLocalNetworkAccess=false
 allowedFileTypes=["txt" "md" "jpg"]
 cleanupInterval=30s
 allowedFileDownloadCount=1
 serverPort=5150
 useDefault=false
 trustedProxies=["127.0.0.1"]
```

## Linux
1. Clone the repository:
```Shell
git clone https://github.com/5ur/Goshh-Server.git
cd Goshh-Server/
```
2. Build in the same way
```Shell
# If need be (eg; missing go.mod or you have a different go version)
go mod init Gosh-Server
go mod tidy

# Finally:
go build .
```
3. Create your config file
```Shell
mv config.yaml.example config.yaml
vim config.yaml
```
Example config file:
```YAML
# Gin debug mode:
debugMode: true
# Port for the engine:
serverPort: 5150
# Use default gin router (if false new is created):
useDefault: true
# Trusted proxy slice/array
trustedProxies:
  # Local host only:
  - 127.0.0.1
  # Some RFC 1918 range:
  - 10.0.0.0/8
# Time in which the slice of messages will be dumped
cleanupInterval: "10s"
# This achieves the same as adding all 1918 ranges in trustedProxies, just more convenient
allowLocalNetworkAccess: true
# You understand what this does with no context needed
allowedFileTypes:
  - txt
  - md
  - jpg
# You understand what this does with no context needed
fileSavePath: "tmp/"
# Time after which a stale file will be deleted (ie; not downloaded at all, not downloaded enough times to reach the allowedFileDownloadCount limit)
staleFileTTL: "30s"
# The amount of times a file is allowed to be downloaded (kept in check by the file struct values)
allowedFileDownloadCount: 1 
```
4. That's it.
Mind the follwoing:  
If the server is started with no configuration file or a missing value/s, it will just start with the default ones hard-coded into the server binary.  
eg:
```Shell
> ./Goshh-Server
2023/05/20 12:01:23 Error reading config file: open config.yaml: no such file or directory
2023/05/20 12:01:23 Loading configuration values:
 debugMode=false
 serverPort=5150
 allowedFileTypes=["txt" "md" "jpg"]
 fileSavePath="/path/to/save/files"
 staleFileTTL=30s
 allowedFileDownloadCount=1
 useDefault=false
 trustedProxies=["127.0.0.1"]
 cleanupInterval=30s
 allowLocalNetworkAccess=false
```

## Nginx + SSL(LE)

Install Certbot:  
```Shell
apt install certbot -y
```
Generate Let's Encrypt Certificate:  
```shell
certbot certonly --standalone -d your_domain_name_here.com
```

Create a new nginx site:  
`vim /etc/nginx/sites-available/goshh-server`
```nginx
server {
    listen 80;
    server_name your_domain_name_here.com;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl;
    server_name your_domain_name_here.com;

    ssl_certificate /etc/letsencrypt/live/your_domain_name_here.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your_domain_name_here.com/privkey.pem;

    location / {
        proxy_pass http://serverhost.local:5150;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_redirect off;
    }
}
```
Add a link:  
`ln -s /etc/nginx/sites-available/goshh-server /etc/nginx/sites-enables/goshh-server`

Reload nginx, so it will take the new config:  
```shell
nginx -t && systemctl reload nginx
```

# Usage
## Standalone:
### Windows:
```Powerhsell
PS ❯ .\Goshh-Server.exe
2023/05/19 14:21:05 Loading configuration values:
 debugMode=true
 serverPort=5150
 allowLocalNetworkAccess=true
 fileSavePath="tmp/"
 staleFileTTL=30s
 useDefault=true
 trustedProxies=["127.0.0.1"]
 cleanupInterval=10s
 allowedFileTypes=["txt" "md" "jpg"]
 allowedFileDownloadCount=1
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /                         --> main.main.func2 (4 handlers)
[GIN-debug] POST   /message                  --> main.main.func3 (4 handlers)
[GIN-debug] GET    /message/:id              --> main.main.func4 (4 handlers)
[GIN-debug] POST   /upload                   --> main.main.func6 (4 handlers)
[GIN-debug] GET    /download/:filename       --> main.main.func7 (4 handlers)
[GIN-debug] Listening and serving HTTP on :5150
```
### Linux:
```Shell
> ./Goshh-Server
2023/05/20 08:21:44 Loading configuration values:
 useDefault=true
 trustedProxies=["127.0.0.1"]
 cleanupInterval=10s
 debugMode=true
 serverPort=5150
 fileSavePath="tmp/"
 staleFileTTL=30s
 allowedFileDownloadCount=1
 allowLocalNetworkAccess=true
 allowedFileTypes=["zip" "txt" "md" "jpg"]
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /                         --> main.main.func2 (4 handlers)
[GIN-debug] POST   /message                  --> main.main.func3 (4 handlers)
[GIN-debug] GET    /message/:id              --> main.main.func4 (4 handlers)
[GIN-debug] POST   /upload                   --> main.main.func6 (4 handlers)
[GIN-debug] GET    /download/:filename       --> main.main.func7 (4 handlers)
[GIN-debug] Listening and serving HTTP on :5150
```
## Service:
### Windows
**.NET**
I'm not explaining or giving a template for this, but you cab use this doc: [.NET service](https://learn.microsoft.com/en-us/dotnet/framework/windows-services/walkthrough-creating-a-windows-service-application-in-the-component-designer)

**Powershell** (core is recommended, since the *-Service commandlets are more developed there)
```Powershell
New-Service -Name "Goshh Server" -BinaryPathName "Full_path_to_Goshh-Server_here.exe"
```
**Or**
```Powershell
sc.exe create "Goshh Server" binpath= "Full_path_to_Goshh-Server_here.exe"
```
Unless it's changed in the future, the service will be created with no additional prompts, and will be set to Automatic by default.  

### Linux
**Systemd:**  
Make a new user for the service:
```Bash
useradd -m -d /home/goshh -s /bin/bash goshh

# Lock the user, you won't be using it for anything. Besides you can just su to it.
passwd -l goshh
```
Create a service file:  
`vim /etc/systemd/system/goshh-server.service`  
and place the following, or something like it in the new service file:  
[Have a look here for more arguments and options.](https://www.freedesktop.org/software/systemd/man/systemd.service.html)
```Bash
[Unit]
Description=Goshh Server
After=network.target

[Service]
User=goshh
WorkingDirectory=/home/goshh/
ExecStart=/home/goshh/Goshh-Server
Restart=on-failure
RestartSec=5s
StandardOutput=file:/var/log/goshh-server.log
StandardError=file:/var/log/goshh-server.log

[Install]
WantedBy=default.target
```

Reload the daemon:  
`systemctl daemon-reload`

Enable and start:
```Shell
systemctl enable goshh-server
systemctl start goshh-server
```
**SysV init script:**  
Create a new user for the service:
```Bash
useradd -m -d /home/goshh -s /bin/bash goshh

# Lock the user, you won't be using it for anything. Besides you can just su to it.
passwd -l goshh
```
Create the init script:
`vim /etc/init.d/goshh-server`
And place in something like this:  
See:  
https://manpages.debian.org/testing/sysvinit-utils/init-d-script.5.en.html  
https://www.cyberciti.biz/tips/linux-write-sys-v-init-script-to-start-stop-service.html  
```Shell
#!/bin/sh
### BEGIN INIT INFO
# Provides:          Goshh-Server
# Required-Start:    $remote_fs $syslog
# Required-Stop:     $remote_fs $syslog
# Default-Start:     2 3 4 5
# Default-Stop:      0 1 6
# Short-Description: Start/stop Goshh-Server
### END INIT INFO

NAME="Goshh-Server"
USER="goshh"
PATH="/home/goshh/"
CMD="/home/goshh/Goshh-Server"
LOG_FILE="/var/log/goshh-server.log"
LOG_LINES=11

case "$1" in
  start)
    echo "Starting $NAME"
    start-stop-daemon --start --chuid "$USER" --chdir "$PATH" --background --make-pidfile --pidfile /var/run/$NAME.pid --startas /bin/bash -- -c "exec $CMD >> $LOG_FILE 2>&1"
    ;;
  stop)
    echo "Stopping $NAME"
    start-stop-daemon --stop --pidfile /var/run/$NAME.pid --retry=TERM/5/KILL/10 >/dev/null
    ;;
  restart)
    echo "Restarting $NAME"
    $0 stop
    sleep 1
    $0 start
    ;;
  status)
    echo "Checking $NAME status"
    if [ -e /var/run/$NAME.pid ]; then
      echo "$NAME is running"
      echo "Last $LOG_LINES lines of the log file:"
      tail -$LOG_LINES $LOG_FILE
      exit 0
    else
      echo "$NAME is not running"
      exit 1
    fi
    ;;
  *)
    echo "Usage: $0 {start|stop|restart|status}"
    exit 1
    ;;
esac

exit 0
```
Make the script executable:  
`chmod +x /etc/init.d/goshh-server`  
Make the service startup automatically:  
`update-rc.d Goshh-Server defaults`  
Reload SysV:  
`init q`  

# Examples
## curl
Most basic:  
```Shell
❯ curl -X POST http://serverhost.local:5150/message -H 'Content-Type: application/json' -d '{"message": "A message."}'
http://serverhost.local:5150/message/20230520094532

❯ curl http://serverhost.local:5150/message/20230520094532
A message.

❯
```

With a custom rune:  
```Shell
❯ curl -X POST http://serverhost.local:5150/message -H 'Content-Type: application/json' -d '{"message": "Another message","rune": "goshh-server"}'
http://serverhost.local:5150/message/goshh-server

❯ curl http://serverhost.local:5150/message/goshh-server
Another message

❯ 

# There is user input verification for the rune as well:
❯ curl -X POST http://serverhost.local:5150/message -H 'Content-Type: application/json' -d '{"message": "Test","rune": "@@#_(*^&@# *@#_%)*(@&#%)@#& @#N%V@#N &%*@#& %*)@#&%_@)#(*&%*(@#&%_)(@*#&^%@#%)(@*#%1"}'
http://local:5150/message/__NVN__1

❯ curl http://serverhost.local:5150/message/__NVN__1
Test

❯
```
Sending a file:  
```Shell
❯ touch testfile.md

❯ curl -X POST -F "file=@testfile.md" -H "Content-Type: multipart/form-data" http://serverhost.local:5150/upload
http://serverhost.local:5150/download/testfile.md

❯ wget http://serverhost.local:5150/download/testfile.md
--2023-05-20 10:22:03--  http://serverhost.local:5150/download/testfile.md
Resolving serverhost.local (serverhost.local)... fe80::f4ac:8f:b906:8250, 192.168.100.7
Connecting to serverhost.local (serverhost.local)|fe80::f4ac:8f:b906:8250|:5150... connected.
HTTP request sent, awaiting response... 200 OK
Length: 0 [application/octet-stream]
Saving to: 'testfile.md.1'

testfile.md.1                                  [ <=>                                                                                     ]       0  --.-KB/s    in 0s      

2023-05-20 10:22:03 (0.00 B/s) - 'testfile.md.1' saved [0/0]

❯ 

```
## iwr/irm
iwr:  
Basic message:  
```Powershell
$uri = "http://serverhost.local:5150/message"
$headers = @{
    "Content-Type" = "application/json"
}
$body = @{
    message = "A message."
} | ConvertTo-Json

$response = Invoke-WebRequest -Uri $uri -Method POST -Headers $headers -Body $body
```

irm:  
```Powershell
$uri = "http://serverhost.local:5150/message"
$headers = @{
    "Content-Type" = "application/json"
}
$body = @{
    message = "A message."
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri $uri -Method POST -Headers $headers -Body $body
```

With rune:  
```Powerhsell
$uri = 'http://serverhost.local:5150/message'
$headers = @{
    'Content-Type' = 'application/json'
}
$message = 'Another message'
$rune = 'goshh-server'
$body = @{
    'message' = $message
    'rune' = $rune
} | ConvertTo-Json

# Using Invoke-WebRequest (iwr)
$response = Invoke-WebRequest -Uri $uri -Method POST -Headers $headers -Body $body

# Using Invoke-RestMethod (irm)
$response = Invoke-RestMethod -Uri $uri -Method POST -Headers $headers -Body $body
```

Sending a file:  
iwr:  
```Powershell
$file = Get-Item -Path "C:\path\to\testfile.md"
$url = "http://serverhost.local:5150/upload"

Invoke-WebRequest -Uri $url -Method POST -InFile $file -ContentType "multipart/form-data"
```
irm:  
```Powershell
$file = Get-Item -Path "C:\path\to\testfile.md"
$url = "http://serverhost.local:5150/upload"

Invoke-RestMethod -Uri $url -Method POST -InFile $file -ContentType "multipart/form-data"
```

# Recommendation
I spent a shit ton of time adding comments to the source-code. You can completely rebuild anything just by reading the comments, and you should, because this is the way that I made it for myself.  

I urge you to have a look at the supplementary tool I made - [Goshh-Client](https://github.com/5ur/Goshh-Client)  
It's purpose it to make this a bit better, since it features options like QR code generation, off/online usage, pipeline input translation, etc..
<img src="https://github.com/5ur/Goshh/blob/main/logos/client_logo.png" alt="Logo" width="20%" height="20%">
<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
[product-screenshot]: logo/logo.png
[Go]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[Go-url]: https://go.dev/doc/

[Powershell]: https://img.shields.io/badge/powershell-5391FE?style=for-the-badge&logo=powershell&logoColor=white
[Powershell-url]: https://github.com/PowerShell/PowerShell

[Vim]: https://img.shields.io/badge/NeoVim-%2357A143.svg?&style=for-the-badge&logo=neovim&logoColor=white
[Vim-url]: https://github.com/AstroNvim/AstroNvim

[StackExchange]: https://img.shields.io/badge/StackExchange-%23ffffff.svg?&style=for-the-badge&logo=StackExchange&logoColor=white
[StackExchange-url]: https://stackexchange.com/

[StackOverflow]: https://img.shields.io/badge/Stack_Overflow-FE7A16?style=for-the-badge&logo=stack-overflow&logoColor=white
[StackOverflow-url]: https://stackoverflow.com/

[Windows]: https://img.shields.io/badge/Windows-0078D6?style=for-the-badge&logo=windows&logoColor=white
[Windows-url]: https://www.microsoft.com/en-us/windows?r=1

[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/othneildrew
