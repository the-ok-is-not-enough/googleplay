# Research

- https://github.com/NoMore201/googleplay-api/issues/155
- https://github.com/NoMore201/googleplay-api/issues/156

## /fdfe/browse

preFetch:

~~~
GET /fdfe/browse?c=3&scat=FINANCE&stcid=apps_topselling_free HTTP/1.1
Host: android.clients.google.com
Authorization: Bearer ya29.A0ARrdaM_rOO6nnrwek1N5zRvc8JzUxa9nmS4-MkWWOxRvxrM9...
X-Dfe-Device-Id: 3588cd1e2b1...
~~~

## /fdfe/getHomeStream

Recommended for you:

~~~
GET /fdfe/getHomeStream?c=3&cat=FINANCE&ctr=apps_topselling_free HTTP/1.1
Host: android.clients.google.com
Authorization: Bearer ya29.A0ARrdaM_rOO6nnrwek1N5zRvc8JzUxa9nmS4-MkWWOxRvxrM9...
X-Dfe-Device-Id: 3588cd1e2b1...
~~~

## /fdfe/homeV2

preFetch:

~~~
GET /fdfe/homeV2?c=3&cat=FINANCE&ctr=apps_topselling_free HTTP/1.1
Host: android.clients.google.com
Authorization: Bearer ya29.A0ARrdaM_rOO6nnrwek1N5zRvc8JzUxa9nmS4-MkWWOxRvxrM9...
X-Dfe-Device-Id: 3588cd1e2b1...
~~~

preFetch:

~~~
GET /fdfe/homeV2?c=3&scat=FINANCE&stcid=apps_topselling_free HTTP/1.1
Host: android.clients.google.com
Authorization: Bearer ya29.A0ARrdaM_rOO6nnrwek1N5zRvc8JzUxa9nmS4-MkWWOxRvxrM9...
X-Dfe-Device-Id: 3588cd1e2b1...
~~~

## /fdfe/listTopChartItems

same as list
