![Untitled drawing (4)](https://github.com/hartsfield/bolt/assets/30379836/b551f0d4-53e5-4309-b7d7-9fc26b7eaa76)

# BOLT

Bolt is a Go program that generates the scaffolding for a web app using a 
component based architecture, and generates server code and other files that
are useful for rapid prototyping and development. Bolt also has other features
useful for rapid development. 

NOTE: Bolt is under heavy development and is still considered beta. Currently 
bolt can be used to build single-page static websites, but in the future, bolt 
will be able to aid in the development of large dynamic websites, and also help
generate code for forms and database models. 

Goals:
 - Generate basic scaffolding to begin building a webapp
 - Tooling to speed up the process of creating generic web components
 - Convert data modeled in CSV into code for forms and basic database procedures
 - Install components from third party git repos
 - deploy to testing server

Create a directory called `boltApp` and `cd boltApp` and run the command:

    bolt --init

This command will generate the scaffolding for a new web app.

The structure of a bolt app looks like this:

    > .
         ├── internal
         │   ├── components
         │   │   ├── footer
         │   │   │   ├── footer.css
         │   │   │   ├── footer.js
         │   │   │   └── footer.tmpl
         │   │   └── head
         │   │       └── head.tmpl
         │   ├── pages
         │   │   └── main
         │   │       └── main.tmpl
         │   └── shared
         │       ├── css
         │       │   └── main.css
         │       └── js
         │           └── main.js
         ├── public
         │   └── media
         ├── autoload.sh
         ├── bolt_conf.json
         ├── Dockerfile
         ├── globals.go
         ├── handlers.go
         ├── helpers.go
         ├── logging.go
         ├── main.go
         ├── restart-service.sh
         ├── server.go
         └── viewdata.go
         
         12 directories, 18 files

A bolt app is composed of html templates called `pages` and `components`. `pages`
are typically composed of components, and components can also be composed of
other components. 

Bolt will also generate some basic server code and build files to begin the 
process of rapid prototyping and development. 

By running the following, we can visit the website generated by the previous 
command:

    ./autoload boltapp 9001
    Server started @ http://localhost:9001

Visit `localhost:9001` in your web browser and you should see the following:

![Screenshot from 2023-07-21 19-07-26](https://github.com/hartsfield/bolt/assets/30379836/832f4789-9212-4af9-9d00-594043bfaa41)

To automatically create a navigation bar using bolt, run the following:

    bolt --autonav=section1,section2,section3

Now visit `localhost:9001` and you should see the following:

![Screenshot from 2023-07-21 19-28-08](https://github.com/hartsfield/bolt/assets/30379836/51c8d948-e086-4d4c-bb90-67f1590b8030)

This command generated a navigation bar, and three components. These 
components are also added to the main `page`. The sections have no content yet, so 
can't be seen, but the boiler plate code to begin their creation has been 
generated.

Components are stored in `internal/components`, and are composed of three files,
one for javascript, another for css, and a .tmpl file which is html with golang 
template directives.

    > .
         ├── internal
         │   ├── components
         │   │   ├── footer
         │   │   │   ├── footer.css
         │   │   │   ├── footer.js
         │   │   │   └── footer.tmpl


Once a component is created, it must then be added to a `page` by specifying it 
in a template directive.
