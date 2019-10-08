FROM ubuntu:bionic

# for Japan
RUN sed -i.bak -e "s%http://archive.ubuntu.com/ubuntu/%http://ftp.jaist.ac.jp/pub/Linux/ubuntu/%g" /etc/apt/sources.list

RUN apt update \
 && apt install -y \
      curl

# ubuntu vanilla
ENV DEBIAN_FRONTEND=noninteractive
RUN apt update \
 && apt upgrade -y \
 && apt install -y `curl https://cloud-images.ubuntu.com/bionic/current/bionic-server-cloudimg-amd64-root.manifest | awk '{print $1}'`

# live
RUN apt update \
 && apt install -y \
      live-boot \
      live-boot-initramfs-tools \
      xorriso \
      squashfs-tools \
 && mkdir -p /isoroot/live

# some custom
RUN apt install -y \
      linux-generic-hwe-18.04 \
      openssh-server
