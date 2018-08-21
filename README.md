Code example for test purposes

Tools needed to run this example:
-RabbitMQ (configure host and port on config.go)
-NATS

Functionalities:

Signup/Login
-Login with email/password or Google account using oauth2. JWT will be generated and exchanged in each request
-Password recover (via email)

User profile
-Change name and last name
-Upload a profile pictures (using amazon S3 and pool of workers with gorutines for image uploading)
-Address CRUD
-Email change (via email)
-Phone modification (using SMS message). The API that sends SMS and verifies phone numbers is implemented
 in another service (theAmazingSmsSender). The communication between the main service and the other one is
 by RabbitMQ queue of messages.

Role management
-See all existing roles. Each role has many permissions that tells what API routes you can access
-Modify the permissions each role has

User managment
-See all existing users
-Modifying users (name, lastname and role)
-Enable/disabling users


