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

## store a passphrase
Passphrases are named entities that have versions. A new phrase version
can be stored as follows:
```bash
mksecret set
Enter passphrase: this is a test phrase
                  NAME                   VERSION          PHRASE          
---------------------------------------+---------+------------------------
  5095fa73-7f24-4952-a258-968ed1d29e31         1   this is a test phrase  
```
Alternatively, a specific name can be provided
```bash
mksecret set --name=my-phrase
Enter passphrase: clarity upper basket money wheel
    NAME      VERSION               PHRASE              
------------+---------+---------------------------------
  my-phrase         1   clarity upper basket money      
                        wheel                           

```
> Please note that the phrase name, once created, cannot be changed later.
> Furthermore, passphrases can only be entered via STDIN to avoid getting
> them captured in command history or files on disk

As you can see, if the named phrase does not already exist, it's version
will be at `1`. Issuing the same command again using an existing named phrase
will generate a new version:
```bash
mksecret set --name=my-phrase
Enter passphrase: tennis water wing code window leaf
    NAME      VERSION               PHRASE              
------------+---------+---------------------------------
  my-phrase         2   tennis water wing code window   
                        leaf                            
```

## encrypt phrases before storing
Additionally encrypt phrases before storing by using `--encrypt` flag:
```bash
mksecret set --encrypt
Enter passphrase: bat country screen puzzle paper ice grain
This input will be encrypted using your password
Enter encryption password (min 8 char): 
Enter encryption password again: 
                  NAME                   VERSION               PHRASE              
---------------------------------------+---------+---------------------------------
  b9fc6cd0-dfc2-4297-b504-08f3da0a9773         1   bat country screen puzzle       
                                                   paper ice grain                 
```
Behind the scenes the code generates an AES key deterministically using your
password and then encrypts the input phrase using that AES key before storing.

## retrieve phrases
Stored phrases can be listed:
```bash
mksecret list
                  NAME                  
----------------------------------------
  5095fa73-7f24-4952-a258-968ed1d29e31  
  my-phrase                                                       
```
And any particular phrase value can be fetched, which always fetches the
`latest` version of the named phrase
```bash
mksecret get 5095fa73-7f24-4952-a258-968ed1d29e31
                  NAME                   VERSION          PHRASE          
---------------------------------------+---------+------------------------
  5095fa73-7f24-4952-a258-968ed1d29e31         1   this is a test phrase  
```
A specific version can be fetched using `--version` flag:
```bash
mksecret get my-phrase --version=2
    NAME      VERSION               PHRASE              
------------+---------+---------------------------------
  my-phrase         2   tennis water wing code window   
                        leaf                            
```

If the phrase was encrypted, you will be asked to provide the password:
```bash
mksecret get b9fc6cd0-dfc2-4297-b504-08f3da0a9773
Enter encryption password: 
                  NAME                   VERSION               PHRASE              
---------------------------------------+---------+---------------------------------
  b9fc6cd0-dfc2-4297-b504-08f3da0a9773         1   bat country screen puzzle       
                                                   paper ice grain                 
```
## delete phrase
When a named phrase is deleted, all versions of secret material are 
deleted forever.
> Please use caution when using this command
```bash
mksecret delete 5095fa73-7f24-4952-a258-968ed1d29e31
Type secret name to delete: 5095fa73-7f24-4952-a258-968ed1d29e31
```
The above command will ask for confirmation, however, `--force` option
can be used to skip the confirmation and delete the secret without any
confirmation.

```bash
mksecret delete 5095fa73-7f24-4952-a258-968ed1d29e31 --force
```