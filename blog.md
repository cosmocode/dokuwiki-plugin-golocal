# DokuWiki Fileserver Linking

One of the questions our customers using DokuWiki ask again and again is how to link to files on "the fileserver".

Of course the first goal of a wiki should be to replace many of the files that used live on the file server with proper wiki pages. But that doesn't make sense for everything. And uploading those files is not always the right solution either. Having a centrally stored file (like an Excel Spreadsheet for example) that can be opened and edited right there is often the best (and well tested) solution.

So how do you link to this file? The obvious answer is to use DokuWiki's builtin mechanism to link to UNC paths (aka. Windows share paths) using a link like this ``[[\\fileserver\share\path\file.xlsx|My Excel Sheet]]``. Seems simple enough, but when clicking on this link you will most likely be greeted with this popup:

FIXME screenshot

So what's going on? DokuWiki is a web application. It may run on a server in your local intranet, but you still use a browser to access it and that browser was built to be used on the Internet. And as you know, the Internet is a potentially dagerous place with potentially malicious websites trying to trick you into all kinds of bad stuff. For that reason, web browsers limit what websites (like DokuWiki) can do. And one of the things website may not do is opening links to local resources (like your fileserver).

Wouldn't it be nice if you could tell your browser to trust the DokuWiki site and let it open local links from there? That's exactly what you can do. In Internet Explorer this is the "Local Zone" and it defaults to your local LAN, so that's why that link may work out of the box on IE. On [Firefox](http://kb.mozillazine.org/Links_to_local_pages_do_not_work) and [Chrome](https://github.com/tksugimoto/chrome-extension_open-local-file-link) there are ways to disable the check through hidden settings or extensions.

Once you enabled those security overrides, you can disable the warning in DokuWiki by creating a file called ``conf/lang/en/lang.php`` (adjust the `en` for your language) with the following content:

```
<?php
$lang['js']['nosmblinks'] = '';  
```

So that's your answer. It's possible but requires tweaking at every user's browser in your company (unless you happen to use IE anyway). And the tweak is somwhat cumbersome to apply.

## custom protocol handler as an alternative?

Is there another way possible? Actually there might. One thing browser and operating systems support is to install "protocol handlers". Your browser is the default handler for the "http" protocol - whenever you click on something starting with "http://" your default browser will open and load the resource.

This mechanism could be used to link local files from DokuWiki. A special application would be installed on the users computer handling "locallink://" links. Whenever a user clicks that kind of link, the browser will ask if that link should be opened with the default handler (our program) and would then pass the link on. Our program then could open Excel or the File Explorer on the file on the fileserver.

So how is this better than the solution above you may wonder. After all this handler still has to be installd on each computer. That's true, but most users are familiar with downloading and executing a program from the internet - the download could be provided on the wiki and new users could simply install it themselves. Alternatively, it is relatively easy for the IT department to centrally distribute a single binary to all computers in a large corporate network - easier thant to distribute Browser hacks at least.

However, creating this custom protocol handler would open the opportunity to do much more than just open the link. Since the programm has control over opening the file, it could react on files that aren't found, eg, by displaying a message and/or opening the parent folder instead. It would also be possible to pass alternative access methods. Eg. try a common drive letter first and fall back to an UNC path if that drive isn't mounted. It could also handle the access to the files differently based on the Users' operation system.   

I spent about a day evaluating the feasability of this approach by writing a simple protocol handler in the language go, which allows the same program to be compiled for different operating systems. My prototype works on Linux and Windows already. It shouldn't be too hard to add Mac support.

What's missing now, is turing this prototype into a working product. That includes creating the needed DokuWiki plugin, polishing the code, defining the interaction between the wiki and the application, finishing the Mac version etc. I estimate that there are about 3 to 4 days work left until it's in a state that could be used and released to the public.

If you're interested in having this functionality for your company and want to sponsor the development, please contact us at [dokuwiki@cosmocode.de](mailto:dokuwiki@cosmocode.de).
