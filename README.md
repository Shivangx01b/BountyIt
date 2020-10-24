<h1 align="center">
  <br>
  <a href=""><img src="https://github.com/Shivangx01b/LfiMe/blob/main/logo.png" alt="" width="200px;"></a>
  <br>
  <img src="https://img.shields.io/github/languages/top/Shivangx01b/CorsMe?style=flat-square">
  <a href="https://goreportcard.com/report/github.com/Shivangx01b/CorsMe"><img src="https://goreportcard.com/badge/github.com/Shivangx01b/CorsMe"></a>
  <a href="https://twitter.com/intent/follow?screen_name=shivangx01b"><img src="https://img.shields.io/twitter/follow/shivangx01b?style=flat-square"></a>
</h1>

## What is LfiMe ?
A local file inclusion fuzzer made in golang that's it !

## Help
```
 -grep string
        Specify custom grepping singantures. Ex -grep singantures.txt
  -method string
        Add method name if required. Ex: -method PUT. Default "GET" (default "GET")
  -p string
        Feed the list of payloads to fuzz. Ex: -p ~/wordlists/lfi.txt
  -t int
        Number of workers to use..default 40. Ex: -t 50 (default 40)
  -verify
        Only prints confirmed results. Ex -verify

```

## How to Install

```
$ go get -u -v github.com/shivangx01b/BountyIt
```
## Usage

Single Url
```plain
echo "https://example.com" | BountyIt
```
Multiple Url
```plain
cat http_https.txt | BountyIt -t 70 -p payloads.txt -verify
```
Add another method if required
```plain
cat http_https.txt | LfiMe -t 70  -method "POST" -p payloads.txt -grep signatures.txt
```

## Screenshot
![1414](https://github.com/Shivangx01b/LfiMe/blob/main/run.PNG)

## Note:

- Scanner stores the error results as "error_requests.txt"... which contains urls which cannot be requested


