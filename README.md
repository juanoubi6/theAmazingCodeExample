# Small Microservice application for test purposes. Kind of a Social Network

Tools needed to run this example:
 - RabbitMQ (configure host and port on config.go)
 - NATS
 - Redis (optional, can be run without it)
 - Minio server (local s3 mock)

## Functionality

### Sign up/Login
 - Login with email/password or Google account using oauth2. Unique JWT will be generated for each user. Some server calls needs the JWT in order to work
 - Password recover (via email). All emails are sent from another service (theAmazingEmailSender). The communication between this service and the other one is by NATS.

### User profile
- Change name and last name
- Upload a profile pictures (using amazon S3 and a pool of workers made with go-routines for parallel image uploading). It's meaningless for this functionality because there is only one profile picture but useful when uploading many pictures at once
- Address CRUD
- Email change (via email)
- Phone modification (using SMS message). The API that sends SMS and verifies phone numbers is implemented in another service (theAmazingSmsSender). The communication between this service and the other one is by RabbitMQ queue of messages.

### Role management
- See all existing roles. Each role has many permissions that tells what API routes you can access
- Modify the permissions each role has

### User management
- See all existing users (with pagination)
- Modifying users (name, last name and role)
- Enable/disabling users

### Posts and comments (theAmazingPostManager)
- Post CRUD
- Users can comment and vote posts. Also, users can comment and vote other comments
- See all existing posts with a certain order (with pagination)
- Get last created posts and comments (uses Redis if the application could establish connection with the server)

## How to run each service
First, you will need to create the database and schema. You must use the [migration script](https://github.com/juanoubi6/migrationScript) to set up the initial table structure. You will need [Govendor](https://github.com/kardianos/govendor) for dependencies. Once in the project root:
```sh
$ govendor init
$ govendor add +external
$ go build && migrationScript.exe
```
You can run all services with those 3 commands. The first one creates the vendor folder, the second one gets all dependencies and the third one compiles and executes the service.

### Considerations to run each service
- [theAmazingCodeExample](https://github.com/juanoubi6/theAmazingCodeExample) -  You'll need RabbitMQ, NATS and Minio server running. Also, you'll need a google places api key and a google project oauth key.
- [theAmazingSmsSender](https://github.com/juanoubi6/theAmazingSmsSender) - You will need RabbitMQ server running. Also, you'll have to create an account in Twilio to fill the env params.  
- [theAmazingEmailSender](https://github.com/juanoubi6/theAmazingEmailSender) - You will need NATS server running. Also, you'll have to create an account in Sendgrid and create an API key to fill the env params.  
- [theAmazingPostManager](https://github.com/juanoubi6/theAmazingPostManager) - You can start Redis server before running the service if you want it to use it, but it's optional
