from subprocess import run, CalledProcessError
from multiprocessing import Pool, cpu_count
import shlex
import shutil
import os
import json

# build config, would be altered by init_config()
APP_NAME = "iris"
STRIP = True
VERBOSE = False
FORMAT = True
PLATFORMS: list[str] = []


def init_config():
    global APP_NAME, STRIP, VERBOSE, FORMAT, PLATFORMS
    with open("./scripts/release.config.json") as f:
        config = json.load(f)

        APP_NAME = config["app_name"]
        STRIP = config["strip_binaries"]
        VERBOSE = config["verbose"]
        FORMAT = config["format_code"]
        PLATFORMS = config["platforms"]


def init_folders() -> None:
    """
    Makes sure that the `temp` and `dist` folders exist.
    """
    if not os.path.exists("./dist"):
        os.mkdir("./dist")

    if not os.path.exists("./temp"):
        os.mkdir("./temp")


def pack(dir: str, platform: str) -> None:
    """
    Copies README, LICENSE, CHANGELOG and iris logo to the output directory and creates a tarball file for the given platform.
    """
    shutil.copyfile("./README.md", f"{dir}/README.md")
    shutil.copyfile("./LICENSE.txt", f"{dir}/LICENSE.txt")
    # shutil.copyfile("./CHANGELOG.md", f"{dir}/CHANGELOG.md")
    shutil.copyfile("./assets/icon.png", f"{dir}/icon.png")

    splitted = platform.split("/")
    build_os = splitted[0]
    build_arch = splitted[1]

    compression = "zip" if build_os == "windows" else "gztar"

    shutil.make_archive(f"dist/{APP_NAME}_{build_os}_{build_arch}", compression, dir)


def build(platform: str) -> None:
    """
    Calls the go compiler to build the application for the given platform, and the pack function.
    """
    try:
        print(f"==> Building for {platform}.")
        splitted = platform.split("/")
        build_os = splitted[0]
        build_arch = splitted[1]

        output_dir = f"temp/{APP_NAME}_{build_os}_{build_arch}"
        if not os.path.exists(output_dir):
            os.makedirs(output_dir)

        executable_path = f"{output_dir}/{APP_NAME}"
        if build_os == "windows":
            executable_path += ".exe"

        os.environ["GOOS"] = build_os
        os.environ["GOARCH"] = build_arch

        run(
            shlex.split(
                "go build -o {} {} {}".format(
                    executable_path,
                    '-ldflags="-s -w"' if STRIP else "",
                    "-v" if VERBOSE else "",
                )
            ),
            check=True,
        )

        print(f"==> Packing for {platform}.")
        pack(output_dir, platform)

    except CalledProcessError:
        print(f"==> Failed to build for {platform}: The Go compiler returned an error.")

    except Exception as e:
        print(f"==> Failed to build for {platform}.")
        print(e)


def cleanup() -> None:
    """
    Removes the `temp` folder.
    """
    print("==> Cleaning up.")
    shutil.rmtree("./temp")


if __name__ == "__main__":
    print("==> Initialising folders, executing prebuild commands.")
    init_config()
    init_folders()
    if FORMAT:
        run(shlex.split("go fmt ./..."), check=True)

    max_procs = cpu_count()
    print(f"==> Starting builds with {max_procs} parallel processes.")

    with Pool(processes=max_procs) as pool:
        pool.map(build, PLATFORMS)

    cleanup()
