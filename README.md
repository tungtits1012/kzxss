# kzxss

`kzxss` takes a list of URLs and checks for **parameter reflection** (potential XSS points).  
It modifies each query parameter with a payload (`kzxss`) and tests both **GET** and **POST** requests to detect whether the payload is reflected in the HTTP response.

This tool is useful for quickly triaging large sets of URLs to identify potential XSS injection points, especially when combined with tools like [hakrawler](https://github.com/hakluke/hakrawler) or [gau](https://github.com/lc/gau).

## Features
- **Dynamic Parameter Handling**: Automatically parses query parameters and injects payloads.  
- **GET & POST Support**: Tests both request methods for reflection.  
- **Concurrent Processing**: Speed up testing with multiple workers.  
- **Configurable Timeout & Retries**: Handle unstable connections gracefully.  
- **STDIN Friendly**: Easily chain with other reconnaissance tools.  

---

## Installation
```bash
go install github.com/xkmikze/kzxss@latest
```

# Usage
```
cat urls.txt | kzxss
echo "https://example.com/?parameter=value" | kzxss
```
# Sample Usage
```
echo "https://example.com/showroom/?parameter1=value" | kzxss
[REFLECTION:GET] https://example.com/?parameter1=value (param: parameter1)
[REFLECTION:POST] https://example.com/?parameter2=value (param: parameter2)

echo "https://example.com/showroom/?parameter3=value&parameter4=value" | kzxss
[REFLECTION:GET] https://example.com/?parameter3=value&parameter4=value (param: parameter3)
[REFLECTION:GET] https://example.com/?parameter3=value&parameter4=value (param: parameter4)
[REFLECTION:POST] https://example.com/?parameter3=value&parameter4=value (param: parameter3)
[REFLECTION:POST] https://example.com/?parameter3=value&parameter4=value (param: parameter4)

cat urls.txt | kzxss
...
```
# Combine with other tools:
```
echo example.com | hakrawler | kzxss
```
