# 📍 Pusher

Pusher is an application designed to quickly setup a bare-metal server 
for running your applications using Traefik and Docker. Within a few 
minutes you can have a machine setup and your application running.

> [!WARNING]
> This tool isn't meant to be used in a production setting. In fact,
> it really isn't meant for anyone but myself. I do NOT guarantee
> anything about this app, its performance, or outcomes. And I do not
> guarantee I'll fix or change anything.

## Example

```bash
pusher prepare
pusher deploy
```

## 🚀 Quickstart

To begin there are a few prerequisites. 

#### Your Machine

1. Docker
2. A config file with your host setup in `~/.ssh.config`.
    - This host must havea **HostName**, **User**, and **IdentityFile**
3. SSH tools, like `scp`

#### Remote Machine

- Your remote machine must have your SSH key setup 
- The server must be a Debian-flavored system with tools like `apt`

#### Your Application

1. Must have a **Dockerfile**
2. Must have an **.env** file of some type (doesn't have to be named that)

### Preparing your Server

Before deploying your application you must prepare your server. This should
be a one time task.

```bash
pusher prepare
```

You will be asked the following:

- Choose a host. This list comes from your `~/.ssh/config` file.
- Enter an email for Let's Encrypt SSL certificates.

<details>
<summary>What does this do?</summary>

- Updates and upgrades the server's packages
- Installs OS certificates and tools like wget, htop, neovim, and git
- Creates three directories in your user's home directory: `applications`, `services`, and `traefik`
- Sets up Docker
    - Installs Docker
    - Adds your user to the **docker** group
    - Creates two Docker networks: `applications` and `web`
    - Installs [lazydocker](https://github.com/jesseduffield/lazydocker)
- Sets up [Traefik](https://traefik.io)
    - Creates configurations for the Traefik image and config file in `~/traefik`
    - Starts the Traefik container
</details>

### Deploy your Application

To deploy your application, run the following.

```bash
pusher deploy
```

You will be asked the following:

- The name of your application. This should be a directory/url friendly name, with no spaces.
- The port exposed by the application. This should be unique on the server, as it is forwarded to `localhost` on the same port.
- The domain this will be bound to. For example `testing.mydomain.com`, or `mydomain.net`.
- The name of an environment file to use, such as `.env`. I tend to use `.env.production` (make sure to .gitignore these!!)
- Any dependencies this application container has
- Any volume mounts you wish to setup.

<details>
<summary>What does this do?</summary>

- Creates a directory at `~/applications/<app name>`. This is where your Docker compose and env files are copied to.
- If you specified any mounts, those directories will be created on the server if they do not exist.
- A Docker image is built locally into a TAR file, then uploaded to the server.
- The container is launched on the server.
</details>
