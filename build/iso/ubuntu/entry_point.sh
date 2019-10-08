#!/bin/bash
set -x
set -e


cat <<EOF | tee /boot/grub/grub.cfg
set timeout=1

serial --speed=115200 --unit=0 --word=8 --parity=no --stop=1
terminal_input console serial
terminal_output console serial

menuentry 'ubuntu-live' {
  linux `ls /boot/vmlinuz* | head -1` boot=live console=tty0 console=ttyS0,115200
  initrd `ls /boot/initrd.img* | head -1`
}
EOF
cp -r /boot /isoroot/boot

# wget -O/isoroot/live/rootfs.squashfs http://cdimage.ubuntu.com/ubuntu-server/daily/current/eoan-server-amd64.squashfs
mksquashfs \
  / /isoroot/live/rootfs.squashfs \
  -comp gzip \
  -xattrs \
  -noD \
  -progress \
  -regex \
  -e \
  '^proc$/.*' \
  '^dev$/.*' \
  '^sys$/.*' \
  '^tmp$/.*' \
  '^boot$/.*' \
  '^isoroot$/.*' \
  '^dst$/.*'
  # /boot \
  # /isoroot \
  # /dst \
  # /proc \
  # /dev \
  # /sys \
  # /tmp
  # -ef \
  # -mem 5G \

grub-mkrescue -o /dst/live.iso /isoroot
