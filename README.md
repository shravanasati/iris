# iris

<p align="center"> 
	<img src="assets/gopher.png" height="300px">
</p>

iris is an easy to use, cross platform and customizable wallpaper manager.

<br>

## 🌐 Table of Contents

- [Features](https://github.com/Shravan-1908/iris#-features)
- [Installation](https://github.com/Shravan-1908/iris#%EF%B8%8F-installation)
- [Motivation](https://github.com/Shravan-1908/iris#-motivation)
- [Usage](https://github.com/Shravan-1908/iris#-usage)
    * [Root command](https://github.com/Shravan-1908/iris#root-command)
    * [Customization](https://github.com/Shravan-1908/iris#customization)
    * [Update](https://github.com/Shravan-1908/iris#update)
- [Changelog](https://github.com/Shravan-1908/iris#-changelog)
- [Versioning](https://github.com/Shravan-1908/iris#-versioning)
- [Licensing](https://github.com/Shravan-1908/iris#-license)
- [Contribution](https://github.com/Shravan-1908/iris#-contribution)


<br>

## ✨ Features
- Cross platform (Windows and Linux)
- Customizable
- Easy to use
- Low memory overhead and CPU usage
- Support for remote wallpapers as well as local wallpapers
- Free & Open Source

<br>

## ⚡️ Installation

Open powershell **as Admin** and execute the following command:

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; (Invoke-WebRequest -Uri https://raw.githubusercontent.com/Shravan-1908/iris/master/scripts/windows_install.ps1 -UseBasicParsing).Content | powershell -
```

This installation script will automatically add iris to PATH and startup applications, so whenever the PC is booted, iris will be launched.

<br>

## 💫 Motivation
I wanted a wallpaper manager which gave a bing wallpaper + nitrogen like interface, good wallpapers and customizability with a bunch of features.

<br>

## 💡 Usage

### Root command

Simply calling `iris` without any flags or arguments from the terminal would launch iris and it will change the desktop wallpaper according to the set configuration.

### Customization

iris uses [unsplash](https://source.unsplash.com) for fetching remote wallpapers. However, you can use your own collection of wallpapers too.

You can customize iris to work as you wish by editing the `config.json` located in `~/.iris` folder. The first time you run iris, it automatically creates a `config.json` file with default configuration, which looks like this:

```json
{
    "search_terms": [
        "nature"
    ],
    "resolution": "1600x900",
    "change_wallpaper": false,
    "change_wallpaper_duration": -1,
    "wallpaper_directory": "",
    "selection_type": "random",
    "save_wallpaper": false
}
```

All configuration fields are pretty self explanatory, still I'd like to describe them all in brief.

- Search Terms: The search terms for unsplash images, i.e., which kind of wallpaper do you want. You can have multiple search terms, but its recommended to not to have more than 3 since it narrows down the search results.

- Resolution: The desired wallpaper resolution. Can only be one of the following:
    * `1024x768`
	* `1600x900`
	* `1920x1080`
	* `3840x2160`

- Change wallpaper: Boolean value for whether to continuously change wallpapers or not.

- Change wallpaper duration: If to change wallpapers, then after how long. The duration value will be counted in minutes. Defaults to 15 mins if the change wallpaper property is set to true.

- Wallpaper directory: Specify your own wallpaper directory if you don't want iris to use unsplash.

- Selection type: If to use wallpapers from the local system, then what should be the selection type: random or sorted.

- Save wallpaper: Boolean value for whether to save the unsplash wallpapers or delete them after usage. If this is set to true, then the wallpapers will be stored in `~/.iris/wallpapers` directory.

### Update

To update iris, run
```
iris up
```
This will update iris to its latest version.

<br>


## ⏪ Changelog

The changes made in the latest release, v0.1.0 are:

- First release

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
