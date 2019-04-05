# identity
WebTech identity microservice 🚐

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
 "last_name": "Lyon"
}
```

*employee*

`/valid/apgart`: 

```json
{
 "error": false,
 "user_type": "Employee",
 "first_name": "Travis",
 "last_name": "Apgar"
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

