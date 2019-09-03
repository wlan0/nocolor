Nocolor
------

Nocolor is a command line tool which strips color output from the output of shell commands.

![nocolor output](nocolor.png)

### Usage
-----------

Install nocolor

#### Linux

```bash
$> curl -s https://api.github.com/repos/wlan0/nocolor/releases/latest | jq -r ".assets[0].browser_download_url" | xargs wget -qO- | tar --strip-components=1 -xzf - releases/nocolor-linux-amd64; echo "nocolor installing to /usr/local/bin/nocolor"; sudo mv ./nocolor-linux-amd64 /usr/local/bin/nocolor
```

#### OSX

```bash
$> curl -s https://api.github.com/repos/wlan0/nocolor/releases/latest | jq -r ".assets[0].browser_download_url" | xargs wget -qO- | tar -xzf - releases/nocolor-darwin-amd64 --strip-components=1; echo "nocolor installing to /usr/local/bin/nocolor"; sudo mv ./nocolor-darwin-amd64 /usr/local/bin/nocolor
```

Use it

```bash
$> echo "\x1b[30;34malert\x1b[0m"
alert # blue color
$> echo "\x1b[30;34malert\x1b[0m" | nocolor
alert # no color
```
