# PepsPoints
> A discord bot written in GoLang that handles a servers points

[![Build Status](https://travis-ci.com/Druue/PepsPointBot_Go.svg?branch=master)](https://travis-ci.com/Druue/PepsPointBot_Go)

## Release History
* 0.0.1
    * Work in Progress

## Meta
Sophie 
https://github.com/Druue

## Testing
1. have a postgres database ready
2. the database structure is described in db.go the top of https://github.com/Druue/PepsPointBot_Go/blob/master/db.go
2. have a discord bot
3. clone it
4. create a file in the root folder called SECRET.go (actually you can call it whatever you want, but with SECRET.go it's already gitignored)
5. write stuff like postgres data and discord token in there, it should look something like this 
```go
package main

var SECRET = &Secret{
	DISCORD_TOKEN: "<token>",
	DB_HOST: "localhost", //or something else
	DB_NAME: "PepsPointsBot", //or whatever your database name is called
	DB_PORT: "5432", //this is just the standard postgres port
	DB_PASSWORD: "super_secure_password",
	DB_USER: "super_secure_user",
}
``` 

## TODO
- ?whois command, that will parse the user given and return the user
- "Maybe make it so that if there's two Emilies, it defaults to the one with a set nickname, and if someone has the same nick, then inform the user that the must specify which person they meanLike .givepoints Emilie(0) or .givepoints Emilie(1)", from Emilie

## Contributing
1. Fork it (<https://github.com/{usr}/{proj_name}/{fork}>)
2. Create your feature branch (`git checkout -b feature/fooBar`)
3. Commit your changes (`git add {files} && git commit`)
4. Push to the branch (`git push origin feature/fooBar`)
5. Create a new Pull Request