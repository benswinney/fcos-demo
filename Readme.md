# Simple Fedora CoreOS Demo

## Pre-reqs:

LibVirt / KVM

### Create Ignition config with fcct

Install fcct (Fedora CoreOS Configuration Transpiler) via brew (https://brew.sh) on MacOS
```shell
$ brew install fcct
```

I've included a simple yaml file that we will convert into an Ignition file via fcct. The simple yaml file takes a basic http server built in go and serves a Welcome message.

```yaml
variant: fcos
version: 1.0.0
passwd:
  users:
    - name: core
      ssh_authorized_keys:
        - ecdsa-sha2-nistp521 AAA.........

systemd:
  units:
    - name: fcosdemo.service
      enabled: true
      contents: |
        [Unit]
        Description=A Webserver to use for a Fedora CoreOS Demo
        After=network-online.target
        Wants=network-online.target
        [Service]
        Type=forking
        KillMode=none
        Restart=on-failure
        RemainAfterExit=yes
        ExecStartPre=podman pull docker.io/benswinney/fcos-demo
        ExecStart=podman run -d --name fcos-demo-server -p 8080:8080 docker.io/benswinney/fcos-demo
        ExecStop=podman stop -t 10 fcos-demo-server
        ExecStopPost=podman rm fcos-demo-server
        [Install]
        WantedBy=multi-user.target
```

This ignition file will start the http server via a service.

```shell
$ fcct --strict fcos-demo.yaml --output fcos-demo.ign
```

### Boot Fedora CoreOS with Ignition file

```shell
$ wget https://builds.coreos.fedoraproject.org/prod/streams/stable/builds/32.20200809.3.0/x86_64/fedora-coreos-32.20200809.3.0-qemu.x86_64.qcow2.xz
$ xz -d fedora-coreos-32.20200809.3.0-qemu.x86_64.qcow2.xz

IGNITION_CONFIG="/vmstore/fcos-demo.ign"
IMAGE="/vmstore/fedora-coreos-32.20200809.3.0-qemu.x86_64.qcow2"
VM_NAME="fcos-demo"
RAM_MB="2048"
DISK_GB="10"

virt-install --connect qemu:///system -n "${VM_NAME}" -r "${RAM_MB}" --os-variant=fedora32 \
        --import --graphics=vnc --disk "size=${DISK_GB},backing_store=${IMAGE}" \
        --qemu-commandline="-fw_cfg name=opt/com.coreos/config,file=${IGNITION_CONFIG}"
```

TIP: If running on a host with SELinux enabled (use the sestatus command to check SELinux status), make sure your OS image and Ignition file are labeled as svirt_home_t. You can do this by placing them under ~/.local/share/libvirt/images/ or running chcon -t svirt_home_t /path/to/file.

NOTE: Depending on your version of virt-install, you may not be able to use --os-variant=fedora32 and will get an error. In this case, you should pick an older Fedora variant (--os-variant=fedora31 or --os-variant=fedora30 for example). You can find the variants that are supported by you current version of virt-install with osinfo-query os | grep "fedora[2-3][0-9]".

NOTE: virt-install requires both the OS image and Ignition file to be specified as absolute paths.

Once the install completes, you will be able to connect to the virtual machine you've just created via the Virtual Machine IP address and the 8080 port. E.g. http://<VM-IP>:8080/

You should be greeted with Welcome Page


### What actually happened?

I'll quickly explain what's happened:
    1. The VM is booted
    2. The Ignition File you created earlier is loaded into Fedora CoreOS 
    3. The Service is started
    4. The service pulls down a docker image (benswinney/fcos-demo)
    5. The service runs the docker image via podman run
    6. The httpd server serves the simple go application we built
    7. You're smiling at the Welcome message :)

There are multiple ways to load boot Fedora CoreOS and load an ignition file. The above is one such way. The user guide has many more: https://docs.fedoraproject.org/en-US/fedora-coreos/

As you can see, it's very simple to use Fedora CoreOS, and I encourage you to experiment further.