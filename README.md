# iris

<p align="center"> 
	<img src="assets/gopher.png" height="300px">
</p>

iris is an easy to use, cross platform and customizable wallpaper manager.

<br>

## 🌐 Table of Contents

- [Features](https://github.com/Shravan-1908/iris#-features)

- [Installation](https://github.com/Shravan-1908/iris#%EF%B8%8F-installation)
    * [Installation Scripts](https://github.com/Shravan-1908/iris/#installation-scripts)
    * [Package Managers](https://github.com/Shravan-1908/iris/#package-managers)
    * [Using Go compiler](https://github.com/Shravan-1908/iris/#using-go-compiler)
    * [Build from source](https://github.com/Shravan-1908/iris/#build-from-source)

- [Motivation](https://github.com/Shravan-1908/iris#-motivation)

- [Usage](https://github.com/Shravan-1908/iris#-usage)
    * [Root command](https://github.com/Shravan-1908/iris#root-command)
    * [Customization](https://github.com/Shravan-1908/iris#customization)

- [Changelog](https://github.com/Shravan-1908/iris#-changelog)

- [Versioning](https://github.com/Shravan-1908/iris#-versioning)

- [Licensing](https://github.com/Shravan-1908/iris#-license)

- [Contribution](https://github.com/Shravan-1908/iris#-contribution)


<br>

## ✨ Features
- Cross platform
- Customizable
- Easy to use
- Low memory overhead and CPU usage
- Support for remote wallpapers as well as local wallpapers
- Free & Open Source

<br>

## ⚡️ Installation

### Installation Scripts

Open powershell **as Admin** and execute the following command:

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; (Invoke-WebRequest -Uri https://raw.githubusercontent.com/Shravan-1908/iris/master/scripts/windows_install.ps1 -UseBasicParsing).Content | powershell -
```

This installation script will automatically add iris to PATH and startup applications, so whenever the PC is booted, iris will be launched.

### Package Managers

iris is available on various package managers across different operating systems.

### GitHub Releases

iris binaries for all operating systems are available on the [GitHub Releases](https://github.com/Shravan-1908/iris/releases/latest) tab. You can download them manually and place them on `PATH` in order to use them.

### Using Go compiler

If you've Go compiler (v1.16 or above) installed on your system, you can install iris via the following command. 

For Go v1.17:

```
go install github.com/Shravan-1908/iris@latest
```

For Go v1.16:

```
go get github.com/Shravan-1908/iris@latest
```


### Build from source

You can alternatively build iris from source via the following commands (again, requires go1.16 or above):

```
git clone https://github.com/Shravan-1908/iris.git
cd ./iris
go build
```

If you want to build iris in release mode (stripped binaries, compressed dsitribution and cross compilation), execute the following command. You can also control the release builds behavior using the [`release.config.json`](./scripts/release.config.json) file.

```
python ./scripts/build.py
```

<br>

## 💫 Motivation
I wanted a wallpaper manager which gave a bing wallpaper + nitrogen like interface, good wallpapers and customizability with a bunch of features.

<br>

## 💡 Usage

### Root command

Simply calling `iris` without any flags and arguments from the terminal would launch iris and it will change the desktop wallpaper according to the set configuration.

### Customization

iris uses [unsplash](https://source.unsplash.com) for fetching remote wallpapers. However, you can use your own collection of wallpapers too.

When iris is ran for the first time, it automatically configures itself with sensible defaults.

You can customize iris to work as you wish by using the `config` command.


```
$ iris config --help

iris v0.2.0
The config command is used to customize iris according to your needs. All configuration options are exposed as flags.

Examples:

$ iris config --save-wallpaper
$ iris config --wallpaper-directory /home/user/Pictures/Wallpapers
$ iris config --search-terms landscape,nature
$ iris config --change-wallpaper=false  
$ iris config --resolution 1920x1080
$ iris config list

Usage:
  iris config [flags]
  iris config [command]

Available Commands:
  list        List the iris config.

Flags:
  -c, --change-wallpaper                   Whether to change wallpapers continuosly in the background.
  -h, --help                               help for config
  -r, --resolution string                  The image resolution to use for unsplash wallpapers. (default "1920x1080")
  -s, --save-wallpaper                     Whether to save the wallpaper to the local directory. (default true)
  -u, --save-wallpaper-directory string    The local directory to save wallpapers in. (default "C:\\Users\\LENOVO\\.iris\\wallpapers")
  -q, --search-terms strings               The search terms for unsplash wallpapers. (default [landscape])
  -t, --selection-type random              The selection type for choosing wallpapers from the local directory, either random or `sorted`. (default "random")  
  -d, --wallpaper-change-duration string   The duration between wallpaper changes, if to change them continuosly. (default "5m")
  -w, --wallpaper-directory string         The local directory to get wallpapers from.

Use "iris config [command] --help" for more information about a command.
```

All configuration fields are pretty self explanatory, still I'd like to describe them all in brief.

- Search Terms: The search terms for unsplash images, i.e., which kind of wallpaper do you want. You can have multiple search terms, but its recommended to not to have more than 3 since it narrows down the search results.

- Resolution: The desired wallpaper resolution. Can only be one of the following:
    * `1024x768`
	* `1600x900`
	* `1920x1080`
	* `3840x2160`

- Change wallpaper: Boolean value for whether to continuously change wallpapers or not.

- Change wallpaper duration: If to change wallpapers, then after how long. The duration value can be anything like `30s` `4m5s` `1h` `2h30m8s`.

- Wallpaper directory: Specify your own wallpaper directory if you don't want iris to use unsplash.

- Selection type: If to use wallpapers from the local system, then what should be the selection type: random or sorted.

- Save wallpaper: Boolean value for whether to save the unsplash wallpapers or delete them after usage. If this is set to true, then the wallpapers will be stored in `~/.iris/wallpapers` directory by default, unless the following option is not altered.

- Save wallpaper directory: Choose a directory to save wallpapers in. Defaults to `~/.iris/wallpapers`.


You can also view your iris configuration using `iris config list` command.

```
$ iris config list

iris v0.2.1
+---------------------------+----------------------------------+
|          OPTION           |              VALUE               |
+---------------------------+----------------------------------+
| Search Terms              | nature                           |
| Resolution                | 1920x1080                        |
| Change Wallpaper          | false                            |
| Change Wallpaper Duration | 5m                               |
| Wallpaper Directory       |                                  |
| Selection Type            | random                           |
| Save Wallpaper            | true                             |
| Save Wallpaper Directory  |                                  |
+---------------------------+----------------------------------+
```


## ⏪ Changelog

The changes made in the latest release, v0.2.0 are:

- Linux and Mac support
- Shell completion scripts generation
- Config command
- Minor bug fixes

<br>

## 🔖 Versioning
*iris* releases follow semantic versioning, every release is in the *x.y.z* form, where:
- *x* is the MAJOR version and is incremented when a backwards incompatible change to iris is made.
- *y* is the MINOR version and is incremented when a backwards compatible change to iris is made, like changing dependencies or adding a new function, method, struct field, or type.
- *z* is the PATCH version and is incremented after making minor changes that don't affect iris's public API or dependencies, like fixing a bug.

<br>

## 📄 License
License
© 2021-Present Shravan Asati

This repository is licensed under the MIT license. See [LICENSE](LICENSE) for details.

<br>

## 👥 Contribution
Pull requests are more than welcome. For more information on how to contribute to *iris*, refer [CONTRIBUTING.md](CONTRIBUTING.md).
