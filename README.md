<a name="readme-top"></a>
<h2 align="center">Goshh</h3>
  
<!-- LOGO -->
<br />
<div align="center">
  <a href="https://github.com/othneildrew/Best-README-Template">
    <img src="logos/root.png" alt="Logo" width="60%" height="60%">
  </a>

<p align="center">
  A Go message and file sharing service
  <br />
  <a href="https://github.com/othneildrew/Best-README-Template"><strong>Wiki»</strong></a>
</div>

<!-- TOC -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <li><a href="#acknowledgments">Acknowledgments</a></li>
  </ol>
</details>



<!-- ABOUT -->
## About
I was looking for a secret/OPT sharing service in [libhunt](https://selfhosted.libhunt.com/), [awesome-selfhosted](https://github.com/awesome-selfhosted/awesome-selfhosted#communication---custom-communication-systems), and [awesome-privacy](https://github.com/pluja/awesome-privacy#pastebin-and-secret-sharing).  
I saw a few projects that were interesting like; [ots](https://github.com/Luzifer/ots) and [privatebin](https://github.com/PrivateBin/PrivateBin), but I dont like javascropt or php.  
The others I saw were way too complex for me to build and understand, so I opted to make something for myself.

### Built With
[![Go][Go]][Go-url] [![Powershell][Powershell]][Powershell-url] [![Vim][Vim]][Vim-url] [![Exchange][StackExchange]][StackExchange-url] [![Overflow][StackOverflow]][StackOverflow-url] [![Windows][Windows]][Windows-url]

<!-- GETTING STARTED -->
## Getting Started

You can download the pre-built binaries or you can clone the repository and build it yourself.
In my old hardware it takes about 6 seconds to build, so it can only get better if you have any recent hardware.

### Prerequisites

### Go versions above 1.16

This is an example of how to list things you need to use the software and how to install them.
#### Windows
  ```Powershell
  # Download Scoop and make life easier: https://scoop.sh/
  irm get.scoop.sh | iex
  scoop bucket add main
  scoop install main/go
  
  # Or just download it from: https://go.dev/dl/
  ```
  
#### Linux
  ```Shell
  apt install golang
  
  # Or use: https://go.dev/doc/install
  # Download from: https://go.dev/dl/
  rm -rf /usr/local/go && tar -C /usr/local -xzf go*.linux-amd64.tar.gz
  
  # Add to /etc/profile or /etc/bash.bashrc
  export PATH=$PATH:/usr/local/go/bin
  ```

### Installation

#### Server
##### Windows
1. Clone the repo
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
  # Use default gin router (if flase new is created):
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

##### Linux
1. Clone the repository
```Shell
git clone https://github.com/5ur/Goshh-Server.git
cd Goshh-Server/
```
2. Build in the same way
```Shell
# If need be (eg; missing go.mod or you have a diferent go version)
go mod init Gosh-Server
go mod tidy

# Finally:
go build .

# You can run a go install, but it shouldn't be needed.
```

<!-- USAGE EXAMPLES -->
## Usage
### Standalone:
Windows:
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
Linux:
```Bash
./Goshh-Server
```
### As a service
#### Windows
.NET
I'm not explining or giving a template for this.
You can read this doc: https://learn.microsoft.com/en-us/dotnet/framework/windows-services/walkthrough-creating-a-windows-service-application-in-the-component-designer

Powershell
```Powershell
New-Service -Name "Goshh Server" -BinaryPathName "D:\GitHub\Goshh\Goshh-Server\Goshh-Server.exe"
```
Or
```Powershell
sc.exe create "Goshh Server" binpath= <Full path to Goshh-Server>.exe
```

#### Linux
systemd:
Make a new user for the service:
```Bash
useradd -m -d /home/gohh -s /bin/bash gohh

# Lock the user, you won't be using it for anything. Besides you can just su to it.
passwd -l gohh
```
vim /etc/systemd/system/goshh-server.service
and place the following, or something like it in the new file: [the man page](https://www.freedesktop.org/software/systemd/man/systemd.service.html)
```Bash
[Unit]
Description=Goshh Server
After=network.target

[Service]
User=gohh
WorkingDirectory=/home/gohh/
ExecStart=/home/gohh/<path-to-your-binary>
Restart=on-failure
RestartSec=5s
StandardOutput=file:/var/log/goshh-server.log
StandardError=file:/var/log/goshh-server.log

[Install]
WantedBy=default.target
```

Enable and start:
```Shell
systemctl enable goshh-server
systemctl start goshh-server
```

_For more examples, please refer to the [Documentation](https://placeholder)_

<!-- ROADMAP -->
## Roadmap
 - [x] Set the random :id generation to a rune/charset because pluses break it as well
 - [x] Clean up the "\"\n" shit at the end of the message contents.
 - [x] Add the QR code as a flag or config file option.
 - [x] Add file upload
 - [x] Add a manifest and branding to the scripts.
 - [] Add an openapi documentation html page as the root of the server

<!-- CONTRIBUTING -->
## Contributing
Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- LICENSE -->
## License
Distributed under the MIT License. See `LICENSE` for more information.
I chose MIT, because it seems ok for the project.  
In essense, feel free to clone copy edit distribute, temm your friends you made this on your own, and/or whatever you want to. I made this for myself, I have the script, so I don't care if someone else is using it or making "money" or gettinf street cred from it.  
I know for a fact that if it wasnt the go docs, stack*, other go projects, and grep.app I wouldn't have been able to make that, so It's not even mine, I just placed the lines one beneath the other.   
Consider this set of files your property.


<!-- CONTACT -->
## Contact
Petar - [@placeholder](https://example.com) - 5150@penev.xyz
Don't know what you would email me for, but here you go.  
Maybe you are a rich motherfucker looking at random git repos and you'll give me money. I wouldn't mind getting more money.

<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

* [Guy who made the OTS project](https://github.com/Luzifer/ots) - Gave me the idea
* [Guy who made this README.md emplate](https://github.com/othneildrew/Best-README-Template) - Gave me the template which you are reading right now
* [MS Paint](https://apps.microsoft.com/store/detail/paint/9PCFS5B6T72H) - Used it to draw the logos
* [Pixlr](https://pixlr.com/x/) - Used it to make the logos transparent
* [//Grep.app](https://grep.app/) - Found some nice examples in there.
* [ChatGPT](https://openai.com/blog/chatgpt) - Retarded, slow, confusing and wrong, but gave me some **really** good tips when when I just started asking it for questions about logic rather than making it write code.

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

