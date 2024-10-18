# identity
WebTech identity microservice 🚐

## Developing & Deployment

### Local Development
Install Docker/Docker Desktop

Run: 
```code
docker build . -t identity
```
Check whether or not you have the newly created image
```code
docker images
```
It should look something like below:
```
REPOSITORY   TAG       IMAGE ID       CREATED             SIZE
identity     latest    9d2fd23ff250   21 seconds ago      845MB
```
Now we need to run the image:
```
docker run -p 127.0.0.1:8080:8080 identity
```

### Production Development
See [the Production Server Administration article](https://github.com/wtg/Shuttle-Tracker-Server/wiki/Production-Server-Administration) on [our wiki](https://github.com/wtg/Shuttle-Tracker-Server/wiki) for important information on code deployment, certificate renewal, and other common administration tasks. There are also many other useful articles about various relevant topics on our wiki.

## Usage

Perform a get request on `https://webtech.union.rpi.edu/identity/valid/[rcsid]` with an authorization header containing the identity token.

```
Authorization: Token [key]
```

Example request: 

```bash
curl http://localhost:8080/valid/lyonj4 -H "Authorization: Token one"
```

**Responses**

**Bad auth token:** 

`403` `invalid key`

**valid rcs id:**

*student*

`/valid/lyonj4`: 

```json
{
 "error": false,
 "user_type": "Student",
 "first_name": "Joseph",
 "last_name": "Lyon",
 "entry_date": "2016-09-01",
 "class_by_credit": "Senior"
}
```

*employee*

`/valid/apgart`: 

```json
{
 "error": false,
 "user_type": "Employee",
 "first_name": "Travis",
 "last_name": "Apgar",
 "entry_date": "",
 "class_by_credit": ""
}
```

*invalid rcs*

`/valid/notlyonj4`: 

```json
{
 "error": true,
 "message": "Invalid RCS ID"
}
```

