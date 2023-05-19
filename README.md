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

This is an example of how to list things you need to use the software and how to install them.
* Windows (Scoop)
  ```Powershell
  # Download Scoop and make life easier: https://scoop.sh/
  irm get.scoop.sh | iex
  scoop bucket add main
  scoop install main/go
  ```

### Installation

_Below is an example of how you can instruct your audience on installing and setting up your app. This template doesn't rely on any external dependencies or services._

1. Clone the repo
   ```Powershell
   git clone Placeholder
   ```
2. Download the external go packages
   ```Powershell
   cd Placeholder
   go mod tidy
   ```
3. Build the binary
   ```Powershell
   go build
   ```
4. Create your config file
  ```Powershell
  mv config.yaml.example config.yaml; code/vim/notepad/ed config.yaml
  ```
5. That's it

<!-- USAGE EXAMPLES -->
## Usage
Placeholder
_For more examples, please refer to the [Documentation](https://example.com)_

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

