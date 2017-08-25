# Installing Intel® Clear Containers 3.0 on CentOS and RHEL version 7

## Scope

This document covers installing Intel® Clear Containers on a
[CentOS](https://www.centos.org/) or [Red Hat Enterprise
Linux](https://www.redhat.com/) system.

## Required Setup

The installation requires the current user to run sudo without
specifying a password. Verify this with the following commands:

```
$ su -
# echo "$some_user ALL=(ALL:ALL) NOPASSWD: ALL" | (EDITOR="tee -a" visudo)
$ exit

```

## Installation steps

1. Ensure the system packages are up-to-date by running the following:

```
$ sudo yum -y update

```
2. Install Git:

```
$ sudo yum install -y git

```
3. Create the installation directory and clone the repository with the following commands:

```
$ mkdir -p $HOME/go/src/github/clearcontainers
$ cd $HOME/go/src/github/clearcontainers
$ git clone https://github.com/clearcontainers/runtime
$ cd runtime

```
4. Run the installation script:

```
$ script -efc ./installation/centos+rhel-setup.sh

```

Notes:

- Running the installation script can take a long time as it needs to
  download source packages and compile them.

- Although it is not strictly necessary to run the installation
  script using the `script(1)` command, doing so ensures that a log of the
  installation is written to the file `typescript`. This is useful for
  administrators to see what changes were made and can also be used to
  debug any issues.

## Verify the installation was successful

1. Check the `cc-runtime` version with the following command:

```
$ cc-runtime --version

```

2. Run an example with the following command:

```
$ sudo docker run -ti busybox sh

```
