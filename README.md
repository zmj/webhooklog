# webhooklog

A tool to explore and debug webhook APIs.

Begin logging by browsing to `/log/myLogId`.

Subscribe to a webhook with a URL that includes the query string parameter `log=myLogId`. Incoming requests with that parameter will be printed to the log page.