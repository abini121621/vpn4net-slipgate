# ƔƤƝ⇃4↾ƝЄƬ Slipgate — Original Flow Secure-Link Release

This release keeps the original Slipgate installer, full interactive menu, Quick Wizard, user management, tunnel actions, and service behavior. The approved ƔƤƝ⇃4↾ƝЄƬ box-style banner and `@vpnfornet` line are included.

The only configuration-generation change is that generated profile links now use the unchanged VPN4Net app format:

```text
vpnfornet://secure/<Base64 version-1 secure envelope>
```

These links are intended for the existing unchanged applications:

| Application | Required existing version |
|---|---:|
| ƔƤƝ⇃4↾ƝЄƬ Client | 3.0.3 |
| Ɣ⇃4↾Ɲ Administrator | 3.0.1 |

No new APK is included or required. No additional VPS signing-key archive is required.

## One-line installation

After the release files are uploaded to the `main` branch and the Linux binary is attached to the new GitHub Release, install with:

```bash
curl -fsSL https://raw.githubusercontent.com/abini121621/vpn4net-slipgate/main/install.sh | sudo bash
```

The installer follows the original Slipgate flow: it downloads `slipgate-linux-amd64` and then runs the normal interactive `slipgate install` action. It is designed for Linux x86_64 KVM VPS systems.

## Normal usage

After installation, start the unchanged full menu:

```bash
sudo slipgate
```

Use the normal Quick Wizard or Tunnel Management → Share Tunnel flow. A generated profile must start with `vpnfornet://secure/`. Import that link into the unchanged Client v3.0.3 or Administrator v3.0.1 app.

## Public repository contents

Only these public files are required:

- `install.sh`
- `slipgate-linux-amd64`
- `SHA256SUMS`
- `README.md`

The source archive is provided for private backup and maintenance. It contains the compatible link generator, but no app APK, Android signing material, VPS credential, or separate private-key archive.
