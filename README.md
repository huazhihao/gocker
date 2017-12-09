# Gocker

A flyweight docker implemented in Go and inspired by [Bocker](https://github.com/p8952/bocker), using BTRFS subvolumes as union file system layers, cgroup as virtualization method and linux bridge as virtual networking.

# How to use

    $ GOOS=linux go build

    $ vagrant up

    $ vagrant ssh

    $ gocker pull centos 7

    $ gocker images
    centos:7

    $ gocker run centos 7 cat /etc/centos-release
    CentOS Linux release 7.4.1708 (Core)

    $ gocker ps
    CONTAINER	COMMAND
    42031	cat /etc/centos-release

    $ gocker logs 42031
    CentOS Linux release 7.4.1708 (Core)

# License

Copyright (C) 2017 Hua Zhihao

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.