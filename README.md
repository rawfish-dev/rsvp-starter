[![Go Report Card](https://goreportcard.com/badge/github.com/rawfish-dev/rsvp-starter)](https://goreportcard.com/report/github.com/rawfish-dev/rsvp-starter)

## Overview

This is a template for developers to use as a starting point for their own RSVP websites. Packaged functionality includes:

1) Displaying mobile friendly event details
2) Creating / Editing / Deleting invitations
3) Creating / Editing / Deleting guest RSVPs

If you already know [React](https://facebook.github.io/react/), [Redux](http://redux.js.org/) and [Golang](https://golang.org/), you'll find it pretty simple to extend the template!

## Pages

##### Details of your event

This is a sample of how details are displayed on the landing page, using dummy values and images from my own event.

![Details Page](https://image.ibb.co/hf757k/details_screengrab.png)

##### Login

The login page protects the control panel section from unauthorized access.

![Login Page](https://image.ibb.co/gjvp05/login_screengrab.png)

##### Control Panel

The control panel is to go to page for managing categories, invitations and RSVPs. 

![Control Panel Page](https://image.ibb.co/ch1wf5/control_panel_screengrab.png)


## Requirements

- Install NodeJS for npm (built with v8.1.2)
- Install Golang (built with 1.8)
- Install Postgres (built with 9.6.2)
	-  Create an empty DB called `rsvp_starter_development` e.g `createdb rsvp_starter_development`

## How to use

The project consists of 2 main parts bundled into one, the client and the server. The front end is built using React+Redux and requires us to bundle the CSS and JS/JSX files for the backend to serve.

1) Make sure you've installed all the requirements above.
2) Clone this repository into your local workstation.
3) `cd` into the root directory of the repository and run `npm install`
4) Run `node_modules/.bin/webpack --progress --colors` to watch JS and CSS files in the `client` folder. You should see similar output:
```
➜  rsvp-starter git:(master) ✗ node_modules/.bin/webpack --progress --colors
Hash: b3ca840d2868a20ef74d
Version: webpack 1.15.0
Time: 8136ms
     Asset     Size  Chunks             Chunk Names
 bundle.js  2.69 MB       0  [emitted]  bundle
bundle.css  7.61 kB       0  [emitted]  bundle
    + 950 hidden modules
Child extract-text-webpack-plugin:
        + 2 hidden modules
Child extract-text-webpack-plugin:
        + 2 hidden modules
Child extract-text-webpack-plugin:
        + 2 hidden modules
Child extract-text-webpack-plugin:
        + 2 hidden modules
```

The steps above bundle the CS and JS/JSX files into the `server/static/build` folder. It also automatically rebuilds them each time we save some changes.

We are now ready to start the backend server!

1) In a separate terminal tab, `cd` to the root directory of the repository and `cd` into the `server` folder.
2) Run `go get bitbucket.org/liamstask/goose/cmd/goose` and `goose up` to run the DB migrations on your newly created DB.
3) Run `POSTGRES_URL=postgres://postgres@localhost:5432/rsvp_starter_development?sslmode=disable HMAC_SECRET="some_secret" TOKEN_ISSUER="some_issuer" gin --appPort 6001`. You should see the following output:
```
[gin] listening on port 3000
```
This uses the [gin](https://github.com/codegangsta/gin) code utility to watch our backend Go files and reloads the server if we save changes.

If you visit `http://localhost:3000`, you should see the default landing page.
