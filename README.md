Code example for test purposes

Tools needed to run this example:
-RabbitMQ (configure host and port on config.go)
-NATS
-Redis (optional, can be run without it)

Functionality:

Sign up/Login
-Login with email/password or Google account using oauth2. Unique JWT will be generated for each user. Some server
calls needs the JWT in order to work
-Password recover (via email). All emails are sent from another service (theAmazingEmailSender). The communication
between this service and the other one is by NATS.

User profile
-Change name and last name
-Upload a profile pictures (using amazon S3 and a pool of workers made with go-routines for parallel image uploading).
It's meaningless for this functionality because there is only one profile picture but useful when uploading many
pictures at once
-Address CRUD
-Email change (via email)
-Phone modification (using SMS message). The API that sends SMS and verifies phone numbers is implemented
 in another service (theAmazingSmsSender). The communication between this service and the other one is
 by RabbitMQ queue of messages.

Role management
-See all existing roles. Each role has many permissions that tells what API routes you can access
-Modify the permissions each role has

User management
-See all existing users (with pagination)
-Modifying users (name, last name and role)
-Enable/disabling users

Posts and comments (theAmazingPostManager)
-Post CRUD
-Users can comment and vote posts. Also, users can comment and vote other comments
-See all existing posts with a certain order (with pagination)
-Get last created posts and comments (uses Redis if the application could establish connection with the server)
