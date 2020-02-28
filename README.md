# Local Mail

Heavy inspired by python implementation [localmail](https://github.com/mistio/localmail) 
(which works on python2) and description has been partially copied.


`go-localmail` is an SMTP and IMAP server that stores all messages into a single in-memory mailbox. 
It is designed to be used to speed up running test suites on systems that send email,
such as new account sign up emails with confirmation codes. 
It can also be used to test SMTP/IMAP client code.

Features:
* Fast and robust IMAP/SMTP implementations, including multipart messages and unicode support.
* Compatible with python's stdlib clients.
* Authentication based on password length, correct only passwords longer then 6 chars. 
* SSL support

WARNING
: not a real SMTP/IMAP server - **not for production usage**.


## Running go-localmail

```bash
go-localmail
```

This will run localmail in the foreground, SMTP on port `2025` and IMAP on `2993`.

You can pass in arguments to control parameters

```bash
go-localmail -imap <port> -smtp <port> -cert <path>
```
