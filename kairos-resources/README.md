Build the kairos base image:

```bash
./build-image.sh
```

Flash the result to the sdcard and add the config as described here:
https://kairos.io/docs/installation/raspberry/

```bash
 cat build/image.img | sudo dd of=/dev/sda oflag=sync status=progress bs=10MB
```

Add the config:

```bash
./fix-img.sh
```

Plug it on rpi4 and boot power it on
