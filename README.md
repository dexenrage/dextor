# Dextor
Redirects your internet traffic through the [Tor](https://www.torproject.org/) network.

# Installing
* Go to [releases](https://github.com/dexenrage/dextor/releases).
* Choose the latest version.
* Download `PKGBUILD`.
* Go to your `downloads` folder and run `makepkg -si` in your terminal.
* Installed :white_check_mark:

> Type `dextor` to run.


# Usage of Dextor
`dextor --help`:
```
  -c, --connect      Connect to the Tor network.
  -d, --disconnect   Disconnect from the Tor network.
  -g, --fixcfg       Use this for restore configs if something went wrong.
  -s, --fixdns       Use this if the website address can't be resolved.
  -h, --help         Show this message.
  -r, --reconnect    Reconnect to the Tor netork (changes IP address).
  -i, --showip       Show your current IP address.
  -v, --version      Show the current version of the program.
```

# Notes before you use Tor
Tor can't help you completely anonymous, just almost:
* [Is Tor Broken? How the NSA Is Working to De-Anonymize You When Browsing the Deep Web](https://null-byte.wonderhowto.com/how-to/is-tor-broken-nsa-is-working-de-anonymize-you-when-browsing-deep-web-0148933/)
* [Use Traffic Analysis to Defeat TOR](https://null-byte.wonderhowto.com/how-to/use-traffic-analysis-defeat-tor-0149100/)

It's recommended that you should use [NoScript](https://noscript.net) before before surfing the web with Tor.
> The NoScript extension is free, open source add-on allows JavaScript, Java, Flash and other plugins to be executed only by trusted web sites of your choice (e.g. your online bank).

# And please
* **Don't spam or perform DoS attacks with Tor.** It's not effective, you will only make Tor get hated and waste Tor's money.
* **Don't torrent over Tor.** If you want to keep anonymous while torrenting, use a no-logs VPN please.

#
*Rewritten from [TorghostNG](https://github.com/GitHackTools/TorghostNG) (currently deleted :disappointed_relieved:) using [Go](https://golang.org/).*

*Was tested on [Arch Linux](https://archlinux.org/download/).*
