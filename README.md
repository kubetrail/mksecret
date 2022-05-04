# mksecret
CLI tool for storing and managing versions of secrets on Google secret manager.

## disclaimer
>The use of this tool does not guarantee security or suitability
for any particular use. Please review the code and use at your own risk.

## installation
This step assumes you have [Go compiler toolchain](https://go.dev/dl/)
installed on your system.

```bash
go install github.com/kubetrail/mksecret@latest
```
Install shell completion. For instance `bash` completion can be installed
by adding following line to your `.bashrc`:
```bash
source <(mksecret completion bash)
```

Create a Google cloud project and activate Secrets Manager API. Also
create a service account key and then export following two variables after
replacing the values to your setup:
```bash
GOOGLE_PROJECT_ID=your-project-id
GOOGLE_APPLICATION_CREDENTIALS=service-account-file-path.json
```

## store a secret
A secret can be stored as a string.
```bash
mksecret set --name=foo bar
```

> Please note that the secret name, once created, cannot be changed later.
> Furthermore, secrets are best entered via STDIN to avoid getting
> them captured in command history or files on disk

```bash
mksecret set --name=foo
```
```text
Enter secret as a string: bar 2
bar 2
```

## retrieve the secret value
Secret value can be retrieved formatted as `table`, `json` or `native`
```bash
mksecret get foo --output-format=table
```
```text
  NAME   VERSION   PHRASE  
-------+---------+---------
  foo          2   bar 2   
```

As you can see the version is set at 2 since we created `foo` named secret
twice. We can retrieve a particular version
```bash
mksecret get foo --output-format=table --version=1
```
```text
  NAME   VERSION   PHRASE  
-------+---------+---------
  foo          1   bar     
```

## encrypt secrets before storing
Secrets can be encrypted by using `--encrypt` flag:
```bash
mksecret set --name=encrypted-foo --encrypt my super secret string
```
```text
This input will be encrypted using your password
Enter encryption password (min 8 char): 
Enter encryption password again: 
my super secret string
```

Behind the scenes the code generates an AES key deterministically using your
password and then encrypts the input phrase using that AES key before storing.

## retrieve phrases
Stored phrases can be listed:
```bash
mksecret list --output-format=table
```
```text
                  NAME                  
----------------------------------------
  foo  
  encrypted-foo                                                       
```

## delete phrase
When a named phrase is deleted, all versions of secret material are 
deleted forever.
> Please use caution when using this command
```bash
mksecret delete foo
```
```text
Type secret name to delete: foo
```
The above command will ask for confirmation, however, `--force` option
can be used to skip the confirmation and delete the secret without any
confirmation.

```bash
mksecret delete encrypted-foo --force
```
