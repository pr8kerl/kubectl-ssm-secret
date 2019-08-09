# kubectl-ssm-secret

A kubectl plugin to allow import/export of kubernetes secrets key/value pairs to/from AWS SSM Parameter Store under a common path.
The plugin is opinionated. It will look for parameters under a single path. It will not recursively search more than one level under a given path.
Useful if you are reprovisioning clusters or namespaces and need to provision the same secrets over and over.
Or perhaps useful to backup/restore your LetsEncrypt or other certificates.
 
Import example.

Given a couple of parameters stored in param store under the path `/foo`, these can easily be imported into kubernetes into a single secret.

If an AWS parameter at path `/foo/bar` contains a secret value, and the parameter `/foo/passwd` contains a secure password, we can view the keys and vaules in parameter store using the list subcommand:

```
% kubectl ssm-secret list /foo
bar: foobar
passwd: SuperSecretSquirrelPassword
```

These params can then be imported with the following import command:
```
% kubectl ssm-secret import test-secret --ssm-path /foo
imported secret: test-secret
```

And we can then check the secret using the beaut `view-secret` kubectl plugin:
```
% kubectl view-secret test-secret
Multiple sub keys found. Specify another argument, one of:
-> bar
-> passwd

% kubectl view-secret test-secret bar
foobar%

% kubectl view-secret test-secret passwd
SuperSecretSquirrelPassword%
```


## Install

Use [krew](https://github.com/kubernetes-sigs/krew) to install.

```
% curl -LO https://raw.githubusercontent.com/pr8kerl/kubectl-ssm-secret/master/ssm-secret.yaml
% kubectl krew install --manifest=ssm-secret.yaml
```

## Build 

Requires docker and docker-compose installed locally.

* clone the repository
* set your `GOOS` environment variable to match your platform

```
% git clone git@github.com:pr8kerl/kubectl-ssm-secret.git
% cd kubectl-ssm-secret                                                                                                                                                               !11252
% GOOS=darwin docker-compose run --rm make
```

## Use

* Authenticate to AWS
* Authenticate to your kubernetes cluster
* Use the `import` subcommand to create a kubernetes secret from key/valus stored under a parameter store path
* Use the `export` subcommand to copy from a kubernetes secret to a parameter store path
* Use the `--overwrite` flag to overwrite an existing kubernetes secret or existing parameter store keys.
* Use the `--tls` flag with the import subcommand to create a kubernetes tls secret instead of the default opaque type

Examples

```
± % kubectl ssm-secret --help
view or import/export k8s secrets from/to aws ssm param store

Usage:
  ssm-secret list|import|export secret [flags]
  ssm-secret [command]

Examples:

	# view the parameter store keys and values located in parameter store path /param/path/foo
	kubectl ssm-secret list /param/path/foo

	# import to a kubernetes secret called foo from key/values stored at parameter store path /param/path/foo
	kubectl ssm-secret import foo --ssm-path /param/path/foo

	# export a kubernetes secret called foo to aws ssm parameter store path /param/path/foo
	kubectl ssm-secret export foo --ssm-path /param/path/foo

	# display the plugin version
	kubectl ssm-secret version

```

