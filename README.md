# 🚀 LaunchRDP

[![Go Version](https://img.shields.io/github/go-mod/go-version/chrilep/LaunchRDP)](https://golang.org)
[![License](https://img.shields.io/github/license/chrilep/LaunchRDP)](LICENSE)
[![Release](https://img.shields.io/github/v/release/chrilep/LaunchRDP)](../../releases)
[![Issues](https://img.shields.io/github/issues/chrilep/LaunchRDP)](../../issues)
[![Stars](https://img.shields.io/github/stars/chrilep/LaunchRDP?style=social)](../../stargazers)

**LaunchRDP** is a Windows utility written in **Go** with a **Fyne** GUI.  
It manages credentials (username + password) and host definitions, then dynamically generates `.rdp` files with all supported settings. With one click, LaunchRDP launches sessions directly using **mstsc.exe** (Microsoft Remote Desktop Connection).

## ⚠️ Project Status

> **Note:** LaunchRDP is in the **early stages of development**.  
> Features, UI, and storage formats are subject to change. Expect breaking changes until the first stable release.

## ✨ Planned Features

- Manage multiple **credentials** (user + password) securely in Window's cred store
- Define **hosts** with (hopefully) full `.rdp` configuration options
- Dynamically generate temporary `.rdp` files on demand
- One‑click session launch via `mstsc.exe`
- Remember last position and size
- Use custom resolution on remote host to max efficiency between remote resolution and local resolution
- GUI built with [Fyne](https://fyne.io)  
- Export/import settings for portability  

## 🛠️ Installation

Precompiled binaries will be published on the [Releases](../../releases) page once available.  

## 🚦 Usage
- Start the application.
- Add credentials (username + password).
- Add a host definition (address, port, and optional RDP settings).
- Select a host + credential pair and click Launch.
- LaunchRDP generates a temporary .rdp file and opens it with mstsc.exe.

## 📌 Roadmap
- [ ] Credential storage
- [ ] UI for managing hosts and credentials
- [ ] Advanced .rdp settings (screen size, color depth, etc.)
- [ ] Profiles and grouping
- [ ] Auto‑update mechanism

## 📄 License
This project is licensed under the MIT License — see the LICENSE file for details.

## 🤝 Contributing
Contributions, issues, and feature requests are welcome!
Feel free to open an issue or submit a pull request.

## 🙏 Acknowledgements
- Go
- Fyne
- Microsoft’s mstsc.exe client
