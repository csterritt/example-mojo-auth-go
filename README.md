### example-mojo-auth-go

An example (minimal) golang web service, using MojoAuth for authentication.

You'll have to visit [MojoAuth](https://mojoauth.com), create an account, and set up a project there.

Environmental variables which must be set:

* `MOJO_APP_ID`: Your project's API key.
* `MAIN_AUTH_SECRET`: A random string for encoding the session.
* `MAIN_ENC_SECRET`: A random string for encoding the session.
