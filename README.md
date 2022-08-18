# url_shortener
Tiny URL clone

http://shorturl.adadkins.com

TODO's:
- implement logger such as logurus
- implment cloud flare tunneling rather than have an open port 80/443
- can we go embed the html files?
- docker file just needs the binary, implement a 2 stage process to shrink the final image
- logic to check if http checking needs improved
- add logic so we dont save duplicate urls with different hashs
- tests
- can the makefile be cleaned up/modified to be more generic
- can the dockerfile be more generic/not project specific