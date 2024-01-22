# Contributing to iris

👍🎉 First off, thanks for taking the time to contribute! 🎉👍

The following is a set of guidelines for contributing to *iris*, which is hosted on GitHub. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.


## Project Structure
```
├── .github
|  └── ISSUE_TEMPLATE
|     ├── bug_report.md
|     ├── custom.md
|     └── feature_request.md
├── .gitignore
├── CHANGELOG.md
├── CONTRIBUTING.md
├── LICENSE
├── README.md
├── Taskfile.yml
├── assets
|  └── gopher.png
├── cmd
|  └── iris
|     └── main.go
├── go.mod
├── go.sum
├── internal
|  ├── config.go
|  ├── init.go
|  ├── utils.go
|  └── wallpapers.go
└── scripts
   └── windows_install.ps11
```

## Setup Development Environment
This section shows how you can setup your development environment to contribute to iris.

- Fork the repository.
- Clone it using Git (`git clone https://github.com/<YOUR USERNAME>iris.git`).
- Create a new git branch (`git checkout -b "BRANCH NAME"`).
- Install the project dependencies. (`go get ./...`)
- Additionally, you may want to install [task](https://taskfile.dev), an easy to use alternative to make.
- Make changes.
- Stage and commit (`git add .` and `git commit -m "COMMIT MESSAGE"`).
- Push it to your remote repository (`git push`).
- Open a pull request by clicking [here](https://github.com/shravanasati/iris/compare).


## Reporting Issues
If you know a bug in the code or you want to file a feature request, open an issue.
Choose the correct issue template from [here](https://github.com/shravanasati/iris/issues/new/choose).