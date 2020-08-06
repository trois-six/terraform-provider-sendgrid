<a href="https://terraform.io">
    <img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" alt="Terraform logo" title="Terraform" align="right" height="50" />
</a>

# Terraform Provider for Sendgrid

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

```sh
$ make testacc
```

## [Documentation](docs/index.md)

The documentation is created thank's to a fork of https://github.com/terraform-providers/terraform-provider-baiducloud/tree/master/gendocs.

## Known issues

The API KEY API is not completely documented: when you don't set scopes, you get all scopes. This is managed by the provider.

When you set one or multiple scopes, even if you don't set the scope `sender_verification_eligible`, you will get it in the end. So if you want to manage your API keys with this provider, add `sender_verification_eligible` to your list of scopes, if you don't do it, you will always have a difference between what you have in your .tf file vs what you really have in Sendgrid.