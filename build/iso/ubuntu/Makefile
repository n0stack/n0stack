build-iso:
	docker run -it --rm \
		-v $(PWD):/dst \
		n0stack/base-ubuntu \
		bash /dst/entry_point.sh

up: build-iso
	sudo -E qemu-system-x86_64 -enable-kvm -cpu host -m 4G -serial stdio -cdrom live.iso
