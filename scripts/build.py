import hashlib
from subprocess import run, CalledProcessError
from multiprocessing import Pool, cpu_count
import shlex
import shutil
import os
import json
from pathlib import Path

# build config, would be altered by init_config()
APP_NAME = "iris"
STRIP = True
VERBOSE = False
FORMAT = True
PLATFORMS: list[str] = []


def hash_file(filename: str):
    h = hashlib.sha256()

    with open(filename, "rb") as file:
        chunk = 0
        while chunk != b"":
            chunk = file.read(1024)
            h.update(chunk)

    return h.hexdigest()


def init_config():
    try:
        global APP_NAME, STRIP, VERBOSE, FORMAT, PLATFORMS
        release_config_file = Path(__file__).parent.resolve() / 'release.config.json'
        with open(str(release_config_file)) as f:
            config = json.load(f)

            APP_NAME = config["app_name"]
            STRIP = config["strip_binaries"]
            VERBOSE = config["verbose"]
            FORMAT = config["format_code"]
            PLATFORMS = config["platforms"]

    except Exception as e:
        print(f"==> ❌ Some error occured while reading the release config:\n{e}")
        exit(1)


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
    Copies README, LICENSE, CHANGELOG and iris logo to the output directory and creates an archive for the given platform.
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
        print(f"==> 🚧 Building for {platform}.")
        splitted = platform.split("/")
        build_os = splitted[0]
        build_arch = splitted[1]

        output_dir = f"temp/{APP_NAME}_{build_os}_{build_arch}"
        if not os.path.exists(output_dir):
            os.makedirs(output_dir, exist_ok=True)

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

        print(f"==> ✅ Packing for {platform}.")
        pack(output_dir, platform)

    except CalledProcessError:
        print(f"==> ❌ Failed to build for {platform}: The Go compiler returned an error.")

    except Exception as e:
        print(f"==> ❌ Failed to build for {platform}.")
        print(e)


def generate_checksums() -> None:
    project_base = Path(__file__).parent.parent
    dist_folder = project_base / "dist"
    checksum = ""

    for item in dist_folder.iterdir():
        checksum += f"{hash_file(item.absolute())}  {item.name}\n"

    checksum_file = dist_folder / "checksums.txt"
    with open(str(checksum_file), 'w') as f:
        f.write(checksum)


def cleanup() -> None:
    """
    Removes the `temp` folder.
    """
    print("==> 👍 Cleaning up.")
    shutil.rmtree("./temp")


if __name__ == "__main__":
    print("==> ⌛ Initialising folders, executing prebuild commands.")
    init_config()
    init_folders()
    if FORMAT:
        run(shlex.split("go fmt ./..."), check=True)

    max_procs = cpu_count()
    print(f"==> 🔥 Starting builds with {max_procs} parallel processes.")

    with Pool(processes=max_procs) as pool:
        pool.map(build, PLATFORMS)

    print("==> #️⃣  Generating checksums.")
    generate_checksums()

    cleanup()
