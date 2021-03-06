<div align="center">

# box

A container-first command runner 🚀⚡🔥<br>

_Suggestions are always welcome!_

![](https://resources.machinelearning.one/box.png)

</div>

## About

`box` is a command runner that uses a container-first approach to execute commands. It runs on host system as well as inside docker equally well. It is a simple, yet powerful tool that can be used to avoid the hassle of installing heavy or one-off packages and remembering the docker syntax.

It is lightweight (less than 3MB as opposed to docker CLI's 60MB) making it suitable for use in containerized development environments.

`box` executes commands in an ephemeral container, which is destroyed after the command finishes. Results are streamed to the terminal and generated artifacts are mapped to working directory.


## Prerequisites

Download and install:

- [Docker](https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-20-04) 
- [NVIDIA Container Toolkit](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/install-guide.html#docker) (Optional - For GPU Support)

## Getting Started

Box provides pre-compiled binary releases for Linux. 

Simply download the latest release from github and put it on $PATH.

## Usage

### Adding a new command

```
box add -k <key> -i <image> -c <command>
```

Optional Flags:
- `-g` : Enable GPU support
- `-m` : Avoid mounting current directory inside container

### Removing a command

```
box remove -k <key>
```

Optional Flags:
- `-i` : Keep image

### Running a command

```
box run <key> <args...>
```

### Listing local commands

```
box ls
```

### Listing remote commands

```
box ls remote
```

### Fetching remote commands

```
box fetch <key>
```

Optional Flags:
- `-k` : Use alternate key for saving to avoid conflicts

### Adding an alias for frequently used commands

```
box alias -a <alias> <key> <optional-commands>
```

### Example

Run the following command to fetch docker/whalesay image and add it to box:

```
box add -k say -i docker/whalesay:latest -c cowsay
```
After successfully adding the command, you can run it with:

```
box run say "box rocks!"
```
It should output:

```
 ____________ 
< box rocks! >
 ------------ 
    \
     \
      \     
,                    ##        .            
,              ## ## ##       ==            
,           ## ## ## ##      ===            
)       /""""""""""""""""___/ ===        
,  ~~~ {~~ ~~~~ ~~~ ~~~~ ~~ ~ /  ===- ~~~   
*       \______ o          __/            
'        \    \        __/             
____\______/    
```

Let's say we use this command frequently and want to add an alias for it:

```
box alias -a say say
```

Now as long as you are in a directory with proper .box file, you can simply run `say` to execute the command.

```
say "box rocks!"
```

This should still produce the same output.

## Building from Source

Download and install the following software

- [go toolchain](https://go.dev/)
- [upx](https://upx.github.io/) (Optional - For Compression)

1. Clone the repository:

```sh
git clone https://github.com/machinelearning-one/box.git
```

2. Change directories and run the following command to build the binary:

```sh
cd box
go build -ldflags="-s -w"
```

3. (Optional) Compress the binary:

```sh
upx --brute box
```

## Features Roadmap

- [x] `add` new commands to box
- [x] `rm` commands from box
- [x] `run` commands from box
- [x] Works inside a docker container (socket mounting required)
- [x] GPU support
- [x] `ls` local commands and images
- [x] `ls` remote commands and images
- [x] `fetch` for getting remote commands
- [x] `alias` to avoid extra typing
- [ ] Create a library of common commands - (In progress)
- [ ] V2: Use BubbleTea for interactivity

## Acknowledgements

The banner for this repository is a modified version of art by [Ashley McNamara](https://github.com/ashleymcnamara/gophers) released under [CC BY-NC-SA 4.0](https://creativecommons.org/licenses/by-nc-sa/4.0/). Thank you Ashley for making such great art available to everyone!