<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Terraform Provider for Amixr
=============================

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) >= 1.13 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/alertmixer/terraform-provider-amixr`

```sh
$ mkdir -p $GOPATH/src/github.com/alertmixer; cd $GOPATH/src/github.com/alertmixer
$ git clone git@github.com:alertmixer/terraform-provider-amixr
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-amixr
$ make build
```

Using the provider
----------------------
## Fill in for each provider

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-amixr
...
```

### Running tests
In order to run the full suite of acceptance tests, you need to export your email, registered in Amixr,
as the environment variable AMIXR_TEST_USER_EMAIL and then run `make acctest`.
```sh
$ make acctest
```