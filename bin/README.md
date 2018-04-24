# How To Use This Package

> This ```Sentinel SOCKS5 Client``` is only for ```Ubuntu 16.04``` and up, if you have an Ubuntu Operating System installed on your computer, please test and report back if you were able to crack the firewall with this package or not.
> 

1. Download this binary file [sentinel-socks-client](https://github.com/sentinel-official/sentinel/raw/sentinel-socks-node/bin/sentinel-socks-client-alpha-v0.0.1)

2. Open the Terminal and go to the folder in which the downloaded file exists (usually in the 'Downloads' folder), using ```cd ~/Downloads```
  
3. Make it executable by either ```chmod +x sentinel-socks-client-alpha-v0.0.1``` or from the UI, right click on the package and select ```Properties``` from the bottom of the list, then click on ```Permissions``` and check the box where it says ```Allow executing file as program```

4. Then run the package as ```sudo ./sentinel-socks-client-alpha-v0.0.1```

5. Please be patient as it takes time to install dependencies but after a while, the program will finish executing and you should have made connection to Sentinel's SOCKS5 node.

6. One last step is to change your proxy setting in ```Settings > Networks > Network Proxy > Manual```
 and enter ```127.0.0.1``` into Socks Host and ```1080``` into second field
 
7. Check if it is running by: 

8. Now you should be able to crack firewalls. (Please get back if this method doesn't work.)

> NOTE: this all will be programatic in future updates, so that user will be able to connect with a single click without any extra work. Since its in very early stages, there's a bit of learning curve to it.
> In case you want to remove all the proxy settings or you're unable to connect to internet, Remove all the settings by following method:

1. Open terminal and enter ```sudo killall ss-local```
2. From network proxy settings, change it back to ```None```