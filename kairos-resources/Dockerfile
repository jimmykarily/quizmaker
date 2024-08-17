# https://www.spectrocloud.com/blog/kiosk-mode-for-edge-computing-with-kairos-and-kubernetes
#FROM busybox AS downloader

#ADD chrome-linux.zip .
#RUN unzip chrome-linux.zip
#RUN mv chrome-linux /

FROM quay.io/kairos/opensuse:tumbleweed-standard-arm64-rpi4-v3.1.1-k3sv1.29.4-k3s1
RUN zypper ref && zypper --non-interactive in -y cage wpa_supplicant firefox chromium xorg-x11-fonts Mesa-dri-vc4 chrony-pool-openSUSE chrony
RUN zypper --non-interactive in -y fetchmsttfonts || : # This could fail

#COPY --from=downloader /chrome-linux /etc/chrome-linux
#RUN ln -s /etc/chrome-linux/chrome /usr/bin/chrome

COPY <<'EOF' /etc/wpa_supplicant/wpa_supplicant.conf
# This file should be generated using wpa_passphrase
# passphrase here is 12345678
ctrl_interface=/var/run/wpa_supplicant
ctrl_interface_group=admin
network={
  ssid="Yakari"
  psk=09ac611bbc74e4b970529cfc1d2ca4f1d997e9feade0a80ac45d57e13fa876e4
}
EOF
RUN chmod 600 /etc/wpa_supplicant/wpa_supplicant.conf

COPY <<'EOF' /var/lib/rancher/k3s/server/manifests/quizmaker.yaml
---
apiVersion: v1
kind: Secret
metadata:
  name: cookie-secret
type: Opaque
data:
  secret: c2VjcmV0cGFzc3dvcmQ=  # This is the base64 encoded value of "secretpassword"

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: quizmaker
spec:
  replicas: 1
  selector:
    matchLabels:
      app: quizmaker
  template:
    metadata:
      labels:
        app: quizmaker
    spec:
      containers:
        - name: quizmaker
          image: quay.io/jimmykarily/quizmaker:latest
          command: ["/bin/quizmaker"]
          args:
            - "-question-pool"
            - "/data/questions.yaml"
          env:
            - name: GIN_MODE
              value: release
            - name: PORT
              value: "8080"
            - name: QUIZMAKER_COOKIE_SECRET
              valueFrom:
                secretKeyRef:
                  name: cookie-secret
                  key: secret
          volumeMounts:
            - name: data
              mountPath: /data
      volumes:
        - name: data
          hostPath:
            path: /home/kairos/quizmaker-data
            type: Directory 

---
apiVersion: v1
kind: Service
metadata:
  name: quizmaker
spec:
  type: NodePort
  selector:
    app: quizmaker
  ports:
  - protocol: TCP
    port: 80          # Port that the service will expose inside the cluster
    targetPort: 8080  # Port on the container that the service will forward traffic to
    nodePort: 30000   # Port on each node that will forward traffic to the service
EOF


COPY <<'EOF' /etc/pam.d/cage
auth           required        pam_unix.so nullok
account        required        pam_unix.so
session        required        pam_unix.so
session        required        pam_systemd.so
EOF

COPY <<'EOF' /usr/bin/start_cage.sh
#!/bin/bash

# Get the IPv4 address of wlan0
export KAIROS_IP=$(ip -4 addr show wlan0 | grep -oP '(?<=inet\s)\d+(\.\d+){3}')
#/usr/bin/cage chromium-browser http://${KAIROS_IP}:30000 -- --kiosk --enable-features=UseOzonePlatform --ozone-platform=wayland
/usr/bin/cage firefox http://${KAIROS_IP}:30000 -- --kiosk
EOF
RUN chmod 755 /usr/bin/start_cage.sh


COPY <<'EOF' /etc/systemd/system/cage@.service
# This is a system unit for launching Cage with auto-login as the
# user configured here. For this to work, wlroots must be built
# with systemd logind support.

[Unit]
Description=Cage Wayland compositor on %I
# Make sure we are started after logins are permitted. If Plymouth is
# used, we want to start when it is on its way out.
After=systemd-user-sessions.service plymouth-quit-wait.service
# Since we are part of the graphical session, make sure we are started
# before it is complete.
Before=graphical.target
# On systems without virtual consoles, do not start.
ConditionPathExists=/dev/tty0
# D-Bus is necessary for contacting logind, which is required.
Wants=dbus.socket systemd-logind.service
After=dbus.socket systemd-logind.service
# Replace any (a)getty that may have spawned, since we log in
# automatically.
Conflicts=getty@%i.service
After=getty@%i.service

[Service]
Type=simple
ExecStart=/bin/bash /usr/bin/start_cage.sh
Restart=always
User=kairos
# Log this user with utmp, letting it show up with commands 'w' and
# 'who'. This is needed since we replace (a)getty.
UtmpIdentifier=%I
UtmpMode=user
# A virtual terminal is needed.
TTYPath=/dev/%I
TTYReset=yes
TTYVHangup=yes
TTYVTDisallocate=yes
# Fail to start if not controlling the virtual terminal.
StandardInput=tty-fail

StandardOutput=file:/home/kairos/cage.log
StandardError=file:/home/kairos/cage.err.log

# Set up a full (custom) user session for the user, required by Cage.
PAMName=cage

[Install]
WantedBy=graphical.target
#Alias=display-manager.service
DefaultInstance=tty7
EOF

COPY <<'EOF' /etc/systemd/resolved.conf.d/dns.conf
[Resolve]
DNS=8.8.8.8
FallbackDNS=1.1.1.1
EOF


RUN chmod 644 /etc/systemd/system/cage@.service
RUN chown root /etc/systemd/system/cage@.service
