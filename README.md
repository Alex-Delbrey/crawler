# Web Crawler
Basic Web Crawler with Go using goroutines\
Guided by ![Boot.dev](https://www.boot.dev)
## FYI
Project is set to stop the second it reaches a url with a different host than the one provided. This is so that our laptops won't go berserk in the crawl.
## Usage
To run project, use the following command:
```git clone git@github.com:Alex-Delbrey/crawler.git```
```./crawler <website.com> <maxConcurrency> <maxPages>```
Where ```<maxConcurrency>``` is the buffer for go's channel\
and ```<maxPages>``` is the maximum amount of pages you would reach.
## Ideas for extending project by Boot.dev
- Make the script run on a timer and deploy it to a server. Have it email you every so often with a report.
- Save the report as a CSV spreadsheet rather than printing it to the console
- Make requests concurrently to speed up the crawling process
