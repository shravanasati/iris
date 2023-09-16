#!/bin/sh

# This script installs iris.
#
# Quick install: `curl https://raw.githubusercontent.com/Shravan-1908/iris/main/scripts/install.sh | bash`
#
# Acknowledgments:
#   - getmic.ro: https://github.com/zyedidia/eget

set -e -u

githubLatestTag() {
  finalUrl=$(curl "https://github.com/$1/releases/latest" -s -L -I -o /dev/null -w '%{url_effective}')
  printf "%s\n" "${finalUrl##*v}"
}

ensure() {
    if ! "$@"; then err "command failed: $*"; fi
}

platform=''
machine=$(uname -m)

if [ "${GETIRIS_PLATFORM:-x}" != "x" ]; then
  platform="$GETIRIS_PLATFORM"
else
  case "$(uname -s | tr '[:upper:]' '[:lower:]')" in
    "linux")
      case "$machine" in
        "arm64"* | "aarch64"* ) platform='linux_arm64' ;;
        *"86") platform='linux_386' ;;
        *"64") platform='linux_amd64' ;;
      esac
      ;;
    "darwin")
      case "$machine" in
        "arm64"* | "aarch64"* ) platform='darwin_arm64' ;;
        *"64") platform='darwin_amd64' ;;
      esac
      ;;
    "msys"*|"cygwin"*|"mingw"*|*"_nt"*|"win"*)
      case "$machine" in
        *"86") platform='windows_386' ;;
        *"64") platform='windows_amd64' ;;
        "arm64"* | "aarch64"* ) platform='windows_arm64' ;;
      esac
      ;;
  esac
fi

if [ "x$platform" = "x" ]; then
  cat << 'EOM'
/=====================================\\
|      COULD NOT DETECT PLATFORM      |
\\=====================================/
Uh oh! We couldn't automatically detect your operating system.
To continue with installation, please choose from one of the following values:
- linux_arm64
- linux_386
- linux_amd64
- darwin_amd64
- darwin_arm64
- windows_386
- windows_arm64
- windows_amd64
Export your selection as the GETIRIS_PLATFORM environment variable, and then
re-run this script.
For example:
  $ export GETIRIS_PLATFORM=linux_amd64
  $ curl https://raw.githubusercontent.com/Shravan-1908/iris/main/scripts/install.sh | bash
EOM
  exit 1
else
  printf "Detected platform: %s\n" "$platform"
fi

TAG=$(githubLatestTag Shravan-1908/iris)

if [ "x$platform" = "xwindows_amd64" ] || [ "x$platform" = "xwindows_386" ] || [ "x$platform" = "xwindows_arm64" ]; then
  extension='zip'
else
  extension='tar.gz'
fi

printf "Latest Version: %s\n" "$TAG"
printf "Downloading https://github.com/Shravan-1908/iris/releases/download/v%s/iris_%s.%s\n" "$TAG" "$platform" "$extension"

ensure curl -L "https://github.com/Shravan-1908/iris/releases/download/v$TAG/iris_$platform.$extension" > "iris.$extension"

case "$extension" in
  "zip") ensure unzip -j "iris.$extension" -d "./iris" ;;
  "tar.gz") ensure tar -xvzf "iris.$extension" "./iris" ;;
esac

bin_dir="${HOME}/.local/bin"
ensure mkdir -p "${bin_dir}"
ensure mv "./iris" "${bin_dir}"
ensure chmod +x "${bin_dir}/iris"

ensure rm "iris.$extension"
ensure rm -rf "$platform"

echo 'iris has been downloaded to' ${bin_dir}
echo "You can run it with:"
echo "iris"

if ! echo ":${PATH}:" | grep -Fq ":${bin_dir}:"; then
  echo "NOTE: ${bin_dir} is not on your \$PATH. iris will not work unless it is added to \$PATH."
fi
