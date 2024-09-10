## CMS
Project where trying to learn and create wordpress styled back-end in go, where
webpages is created in some kind of editor and pages are managed in
server side rendering

### Run project 

#### Setup .env file
```
PORT=1234
ENVIRONMENT=development|production
```
#### Run front end scripts in development mode
```bash
cd ./src 
npm install
npm run dev
```

#### Run back-end scripts in development mode
Note: Vite must be in running to get scripts and styles in development
mode

```bash
go mod tidy
```
generate templ go files (templ package must be installed)
```
templ generate
```
```
go run .
```


### Used Packages
- https://templ.guide/
- https://github.com/go-chi/chi 
- https://github.com/joho/godotenv
- https://github.com/mattn/go-sqlite3
- https://github.com/go-playground/form 
- https://github.com/go-playground/validator 
- https://pkg.go.dev/golang.org/x/crypto/argon2
- https://github.com/alexedwards/scs

### Front-end packages 
- https://vitejs.dev/
- https://htmx.org/

### Backend vite integration
- https://github.com/olivere/vite

### Template parsing
- https://stackoverflow.com/questions/59948906/is-it-possible-to-combine-parse-and-parsefiles-templates-in-go
- https://stackoverflow.com/questions/40714378/how-to-extend-a-template-in-go

### Other resources
- https://www.alexedwards.net/blog/how-to-hash-and-verify-passwords-with-argon2-in-go
