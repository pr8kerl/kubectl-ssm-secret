# kubectl-ssm-secret

A kubectl plugin to allow import/export of kubernetes secrets to/from AWS SSM Parameter Store path.

The plugin is opinionated. It will look for parameters under a single path. It will not recursively search more than one level under a given path. All parameters found under the given parameter store path can be imported into a single kubernetes secret as StringData.

Useful if you are reprovisioning clusters or namespaces and need to provision the same secrets over and over.
Or perhaps useful to backup/restore your LetsEncrypt or other certificates.
 
## examples

Given a couple of parameters stored in param store under the path `/foo`, these can easily be imported into kubernetes into a single secret.

If an AWS parameter at path `/foo/bar` contains a secret value, and the parameter `/foo/passwd` contains a secure password, we can view the keys and values in parameter store using the `kubectl ssm-secret list` subcommand:

```
% kubectl ssm-secret list --ssm-path /foo
ssm:/foo/bar: foobar
ssm:/foo/passwd: SuperSecretSquirrelPassword
```

These params can then be imported with the following import command:
```
% kubectl ssm-secret import foo --ssm-path /foo
imported secret: foo
```

The resulting kubernetes secret created will look like this:
```
% kubectl get secret foo -o yaml
apiVersion: v1
data:
  bar: Zm9vYmFy
  passwd: U3VwZXJTZWNyZXRTcXVpcnJlbFBhc3N3b3Jk
kind: Secret
metadata:
  creationTimestamp: "2019-08-10T00:42:35Z"
  name: foo
  namespace: default
  resourceVersion: "5565641"
  selfLink: /api/v1/namespaces/default/secrets/foo
  uid: bf0fe887-bb07-11e9-9531-02946becbcee
type: Opaque
```

ssm-secret can also be used to then view the plain-text contents of the kubernetes secret using list subcommand:
```
% kubectl ssm-secret list foo
k8s:default/foo/bar: foobar
k8s:default/foo/passwd: SuperSecretSquirrelPassword
```

Additionally, we can export a secret from kubernetes into a parameter store path:
```
% kubectl ssm-secret export foo --ssm-path /bar
created parameter: /bar/bar, version: 1
created parameter: /bar/passwd, version: 1
exported secret: foo
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
% cd kubectl-ssm-secret
% GOOS=darwin docker-compose run --rm make
```

## Use

* Authenticate to AWS
* Authenticate to your kubernetes cluster
* Use the `list` subcommand to list keys and decoded values from a kubernetes secret or from a ssm parameter store path
* Use the `import` subcommand to create a kubernetes secret from key/values stored under a parameter store path
* Use the `export` subcommand to copy from a kubernetes secret to a parameter store path
* Use the `--overwrite` flag to overwrite an existing kubernetes secret or existing parameter store keys.
* Use the `--tls` flag with the import subcommand to create a kubernetes tls secret instead of the default opaque type
* Use the `--namespace` flag to to override the kubernetes namespace in the current context

```
% kubectl ssm-secret --help
view or import/export k8s secrets from/to aws ssm param store

Usage:
  ssm-secret list|import|export secret [flags]
  ssm-secret [command]

Examples:

        # view the parameter store keys and values located in parameter store path /param/path/foo
        kubectl ssm-secret list --ssm-path /param/path/foo

        # view the kubernetes secret called foo
        kubectl ssm-secret list foo

        # import to a kubernetes secret called foo from key/values stored at parameter store path /param/path/foo
        kubectl ssm-secret import foo --ssm-path /param/path/foo

        # export a kubernetes secret called foo to aws ssm parameter store path /param/path/foo
        kubectl ssm-secret export foo --ssm-path /param/path/foo

        # display the plugin version
        kubectl ssm-secret version


Available Commands:
  export      export a kubernetes secret to aws ssm param store
  help        Help about any command
  import      import a kubernetes secret from aws ssm param store
  list        list ssm parameters by path 
  version     print the ssm-secret version

Flags:
  -h, --help               help for ssm-secret
  -n, --namespace string   kubernetes namespace (default "default")

Use "ssm-secret [command] --help" for more information about a command.
```

```
% kubectl ssm-secret export --help
export a kubernetes secret to aws ssm param store

Usage:
  ssm-secret export [flags]

Flags:
  -e, --encode            gzip, base64 encode values in parameter store
  -h, --help              help for export
  -o, --overwrite         if parameter store key exists, overwite its values with those from k8s secret
  -s, --ssm-path string   ssm parameter store path to write data to

Global Flags:
  -n, --namespace string   kubernetes namespace (default "default")
```

```
% kubectl ssm-secret import --help
import a kubernetes secret from aws ssm param store

Usage:
  ssm-secret import [flags]

Flags:
  -d, --decode            treat store values in param store as gzipped, base64 encoded strings
  -h, --help              help for import
  -o, --overwrite         if k8s secret exists, overwite its values with those from param store
  -s, --ssm-path string   ssm parameter store path to read data from
  -t, --tls               import ssm param store values to k8s tls secret

Global Flags:
  -n, --namespace string   kubernetes namespace (default "default")
```

```
% kubectl ssm-secret list --help
Flags:
  -e, --env               output as environment variable key pairs
  -h, --help              help for list
  -s, --ssm-path string   ssm parameter store path to list parameters from

Global Flags:
  -n, --namespace string   kubernetes namespace (default "svcs")
```
