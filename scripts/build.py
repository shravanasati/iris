from subprocess import run
from typing import List
from multiprocessing import Pool, cpu_count
import shlex
import shutil
import os

# build config
APP_NAME = "iris"
STRIP = True
VERBOSE = False
FORMAT = True


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
    shutil.copyfile("./CHANGELOG.md", f"{dir}/CHANGELOG.md")
    shutil.copyfile("./assets/gopher.png", f"{dir}/iris.png")

    splitted = platform.split("/")
    build_os = splitted[0]
    build_arch = splitted[1]

    shutil.make_archive(f"dist/{APP_NAME}_{build_os}_{build_arch}", "gztar", dir)


def build(platform: str) -> None:
    """
    Calls the go compiler to build the application for the given platform, and the pack function.
    """
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


def cleanup() -> None:
    """
    Removes the `temp` folder.
    """
    print("==> Cleaning up.")
    shutil.rmtree("./temp")


if __name__ == "__main__":
    platforms: List[str] = [
        "linux/amd64",
        "windows/amd64",
        "darwin/amd64",
        "darwin/arm64",
        "linux/arm64",
    ]

    init_folders()
    if FORMAT:
        run(shlex.split("go fmt ./..."), check=True)

    max_procs = cpu_count()
    print(f"==> Starting builds with {max_procs} processes.")

    with Pool(processes=max_procs) as pool:
        pool.map(build, platforms)

    cleanup()
