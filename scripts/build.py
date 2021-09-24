from subprocess import run
from typing import List
from multiprocessing import Pool, cpu_count
import shlex
import shutil
import os

app_name = "iris"


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
    Creates a tarball file for the given directory.
    """
    shutil.copyfile("./README.md", f"{dir}/README.md")
    shutil.copyfile("./LICENSE", f"{dir}/LICENSE")
    shutil.copyfile("./CHANGELOG.md", f"{dir}/CHANGELOG.md")
    shutil.copyfile("./assets/gopher.png", f"{dir}/iris.png")

    build_os = platform.split("/")[0]
    build_arch = platform.split("/")[1]

    shutil.make_archive(f"dist/{app_name}_{build_os}_{build_arch}", "gztar", dir)


def build(platform: str) -> None:
    print(f"==> Building for {platform}.")
    splitted = platform.split("/")
    build_os = splitted[0]
    build_arch = splitted[1]

    output_dir = f"temp/{app_name}_{build_os}_{build_arch}"
    if not os.path.exists(output_dir):
        os.makedirs(output_dir)

    executable_path = f"{output_dir}/{app_name}"
    if build_os == "windows":
        executable_path += ".exe"

    os.environ["GOOS"] = build_os
    os.environ["GOARCH"] = build_arch

    run(shlex.split(f"go build -o {executable_path}"), check=True)

    print(f"==> Packing for {platform}.")
    pack(output_dir, platform)


if __name__ == "__main__":
    platforms: List[str] = ["linux/amd64", "windows/amd64", "darwin/amd64", "darwin/arm64", "linux/arm64"]
    init_folders()

    with Pool(processes=cpu_count()) as pool:
        pool.map(build, platforms)

    print("==> Cleaning up.")
    shutil.rmtree("./temp")
