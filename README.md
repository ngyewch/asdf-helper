# asdf-helper

[![Build](https://github.com/ngyewch/asdf-helper/actions/workflows/build.yml/badge.svg)](https://github.com/ngyewch/asdf-helper/actions/workflows/build.yml)

Helper tool for [asdf](https://asdf-vm.com/).

## Installation

asdf-helper can be installed via [asdf](https://asdf-vm.com/). See https://github.com/ngyewch/asdf-helper-plugin

## Features

### Install (recursive)

```
asdf-helper install
```

asdf plugins are also installed automatically. Custom/unlisted plugins may be specified via a `.plugin-versions` file in the same directory as `.tool-versions`. 

Sample `.plugin-versions`
```
ansible-core https://github.com/amrox/asdf-pyapp.git
earthly https://github.com/YR-ZR0/asdf-earthly
```

### Latest version (recursive)

```
asdf-helper latest
```

Version constraints can be specified in `.tool-versions`. See https://github.com/Masterminds/semver for more details.

```
nodejs 16.17.1 # (constraint ^16)
java openjdk-11.0.2 # (constraint ^11)
golang 1.17.13 # (constraint >= 1.17, < 1.19)
```
