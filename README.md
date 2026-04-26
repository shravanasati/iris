# iris

<p align="center"> 
	<img src="assets/icon.png" height="300px">
</p>

iris is an easy to use, cross platform and customizable wallpaper manager.

<br>

## 🌐 Table of Contents

- [Features](https://github.com/shravanasati/iris#-features)

- [Installation](https://github.com/shravanasati/iris#%EF%B8%8F-installation)
    * [Package Managers](https://github.com/shravanasati/iris/#package-managers)
    * [Using Go compiler](https://github.com/shravanasati/iris/#using-go-compiler)
    * [Build from source](https://github.com/shravanasati/iris/#build-from-source)

- [Motivation](https://github.com/shravanasati/iris#-motivation)

- [Usage](https://github.com/shravanasati/iris#-usage)
    * [Root command](https://github.com/shravanasati/iris#root-command)
    * [Cache](https://github.com/shravanasati/iris#cache)
    * [Customization](https://github.com/shravanasati/iris#customization)
    * [Shell Completions](https://github.com/shravanasati/iris#shell-completions)

- [Changelog](https://github.com/shravanasati/iris#-changelog)

- [Versioning](https://github.com/shravanasati/iris#-versioning)

- [Licensing](https://github.com/shravanasati/iris#-license)

- [Contribution](https://github.com/shravanasati/iris#-contribution)


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

#### Linux and macOS

```bash
curl https://raw.githubusercontent.com/shravanasati/iris/main/scripts/install.sh | bash
```

### Package Managers

#### Windows
```powershell
scoop install https://github.com/shravanasati/iris/raw/main/scripts/iris.json
```

<!-- iris is available on various package managers across different operating systems.

iris is present in the AUR. If you're on an Arch based Linux distro,
execute:

```
git clone https://aur.archlinux.org/iris-bin.git
cd iris-bin
makepkg -si
```

Or use any AUR helper like yay:
```
yay -S iris-bin
``` -->


### GitHub Releases

Use [`eget`](https://github.com/zyedidia/eget) to automatically download and extract the binaries:

```
eget shravanasati/iris
```

iris binaries for all operating systems are available on the [GitHub Releases](https://github.com/shravanasati/iris/releases/latest) tab. You can download them manually and place them on `PATH` in order to use them.



### Using Go compiler

If you've Go compiler (v1.18 or above) installed on your system, you can install iris via the following command. 


```
go install github.com/shravanasati/iris@latest
```


### Build from source

You can alternatively build iris from source via the following commands (again, requires go1.18 or above):

```
git clone https://github.com/shravanasati/iris.git
cd ./iris
go build
```

If you want to build iris in release mode (stripped binaries, compressed distribution and cross compilation), execute the following command. You can also control the release builds behavior using the [`release.config.json`](./scripts/release.config.json) file.

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

### Get

The get command prints path to the currently set wallpaper.

```
iris get
```

### Set

The set command accepts a filepath as an argument and sets it as the desktop wallpaper.

Example:
```
iris set ~/Pictures/my-fav-image.jpg
```

### Video

The video wallpaper support is currently experimental and requires ffmpeg to work.

Example:
```
iris video ~/Videos/path-to-video.mp4
```

Only mp4, mkv and gif files are supported at the moment.

The first time running this command on a single video might take some time since iris first converts the video into frames using ffmpeg and then iterates through the frames and sets each one of them as wallpaper every few milliseconds.

### Cache

The `cache` command provides access to manage iris storage. This includes video frame caches and the GitHub repository result cache.

#### Video Cache
Since video wallpapers are implemented by converting videos into frames, iris caches these frames to avoid reconverting them every time you set a video wallpaper.

```bash
# Print the total cache size used by iris (videos + remote sources)
$ iris cache size

# List out the paths of all videos iris has cached
$ iris cache video list

# Remove a specific video from the cache
$ iris cache video rm "/path/to/video.mp4"

# Clear all video caches
$ iris cache video clear
```

#### GitHub Cache
iris caches the results of GitHub repository listings to support offline mode and reduce API traffic by checking commit SHAs before fetching updates.

```bash
# List all GitHub repositories currently cached
$ iris cache github list

# Force a refresh of the cache for a specific repository on next run
$ iris cache github sync "github.com/user/repo/tree/branch"

# Clear all cached GitHub repository results
$ iris cache github clear
```

#### Global Cache Helpers
```bash
# Clear the entire cache (deletes everything: video frames, github data, etc.)
$ iris cache clear

# Total aggregate size of all caches
$ iris cache size
```

### Customization

iris supports multiple remote sources for fetching wallpapers, including **Windows Spotlight**, **GitHub**, and **Reddit**. You can also use your own local collection of wallpapers.

#### Remote Sources

- **Windows Spotlight**: Fetches images from [windows10spotlight.com](https://windows10spotlight.com).
  - Use with: `iris config --remote-source spotlight`
  - Filter with: `iris config --search-terms "landscape,ocean"`

- **GitHub**: Fetches wallpapers from a specific folder in a GitHub repository.
  - Use with: `iris config --remote-source "https://github.com/owner/repo/tree/branch/path/to/wallpapers"`
  - **TIP**: If you reach rate limits, set your token with `iris config --github-token <your-token>`.

- **Reddit**: Fetches top images from specified subreddits.
  - Basic: `iris config --remote-source "r/wallpapers"`
  - Specialized: `iris config --remote-source "r/wallpapers+earthporn"` (combine subreddits)
  - Advanced: Use sort and time filters:
    - `r/wallpapers/top?t=day` (Top of the day)
    - `r/wallpapers/new` (Newest)
    - `r/wallpapers/hot` (Hot)

When iris is ran for the first time, it automatically configures itself with sensible defaults.

You can customize iris to work as you wish by using the `config` command.


```
$ iris config --help

iris v0.4.0
The config command is used to customize iris according to your needs. All configuration options are exposed as flags.
	
Examples:

$ iris config --remote-source spotlight
$ iris config --search-terms landscape,nature
$ iris config --save-wallpaper[=false]
$ iris config --wallpaper-directory /home/user/Pictures/Wallpapers
$ iris config --change-wallpaper[=false]
$ iris config list

Usage:
  iris config [flags]
  iris config [command]

Available Commands:
  list        List the iris config.

Flags:
  -c, --change-wallpaper                   Whether to change wallpapers continuosly in the background.
      --check-for-updates                  Whether to check for updates of iris from github. (default true)
      --github-token string                The GitHub Personal Access Token (PAT), used to perform authorized requests to fetch wallpapers from GitHub repositories.
  -h, --help                               help for config
  --remote-source string               Remote source to select wallpapers from. Valid options are: spotlight, github, reddit.
  -s, --save-wallpaper                     Whether to save the wallpaper to the local directory.
  -u, --save-wallpaper-directory string    The local directory to save wallpapers in. (default "C:\\Users\\devsh\\.iris\\wallpapers")
  -q, --search-terms strings               The search terms for spotlight remote wallpapers. (default [nature])
  -t, --selection-type random              The selection type for choosing wallpapers from the local directory, either random or `sorted`. (default "random")
  -d, --wallpaper-change-duration string   The duration between wallpaper changes, if to change them continuosly. (default "5m")
  -w, --wallpaper-directory string         The local directory to get wallpapers from. (default "C:\\Users\\devsh\\OneDrive\\Pictures\\favorites")
  -f, --wallpaper-file string              Path to the wallpaper file.

Use "iris config [command] --help" for more information about a command.

```

All configuration fields are pretty self explanatory, still I'd like to describe them all in brief.

- Remote Source: Specify where to fetch wallpapers from. Supported sources are `spotlight` (Windows Spotlight) and `github` (GitHub repository).

- Search Terms: The search terms for remote wallpapers. You can have multiple search terms, but its recommended to not to have more than 3 since it narrows down the search results. The search terms are used only when the remote source is set to spotlight.

- Change wallpaper: Boolean value for whether to continuously change wallpapers or not.

- Change wallpaper duration: If to change wallpapers, then after how long. The duration value can be anything in format `30s` `4m5s` `1h` `2h30m8s`.

- Wallpaper file: Specify path to a single wallpaper file.

- Wallpaper directory: Specify your own wallpaper directory if you don't want iris to use a remote source.

- Selection type: If to use wallpapers from the local system, then what should be the selection type: random or sorted.

- Save wallpaper: Boolean value for whether to save the remote wallpapers or delete them after usage. If this is set to true, then the wallpapers will be stored in `~/.iris/wallpapers` directory by default, unless the following option is not altered.

- Save wallpaper directory: Choose a directory to save wallpapers in. Defaults to `~/.iris/wallpapers`.

- GitHub Token: The GitHub Personal Access Token (PAT) used to perform authorized requests when fetching wallpapers from GitHub repositories. You can create one from [here](https://github.com/settings/tokens). You need to grant "Contents" repository permission for a fine-grained token or the "repo" scope for the classic PAT.


You can also view your iris configuration using `iris config list` command.

```
$ iris config list

iris v0.2.1
+---------------------------+----------------------------------+
|          OPTION           |              VALUE               |
+---------------------------+----------------------------------+
| Search Terms              | nature                           |
| Change Wallpaper          | false                            |
| Change Wallpaper Duration | 5m                               |
| Wallpaper Directory       |                                  |
| Selection Type            | random                           |
| Save Wallpaper            | true                             |
| Save Wallpaper Directory  |                                  |
+---------------------------+----------------------------------+
```

### Shell-Completions

iris can generate shell completions for powershell, fish, bash and zsh.

`iris completion shell_name`

It will output a completion script for your shell. Copy and paste it on your shell profile.


## ⏪ Changelog

Entire changelog can be viewed in the [`CHANGELOG.md`](CHANGELOG.md)

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
