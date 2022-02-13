# mkphrase
CLI tool for storing and managing versions of crypto passphrases on Google secret manager.

## disclaimer
>The use of this tool does not guarantee security or suitability
for any particular use. Please review the code and use at your own risk.

## installation
Download the code to a folder and cd to the folder, then run
```bash
go install
```
Install shell completion. For instance `bash` completion can be installed
by adding following line to your `.bashrc`:
```bash
source <(mkphrse completion bash)
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
mkphrase set
Enter passphrase: this is a test phrase
                  NAME                   VERSION          PHRASE          
---------------------------------------+---------+------------------------
  5095fa73-7f24-4952-a258-968ed1d29e31         1   this is a test phrase  
```
Alternatively, a specific name can be provided
```bash
mkphrase set --name=my-phrase
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
mkphrase set --name=my-phrase
Enter passphrase: tennis water wing code window leaf
    NAME      VERSION               PHRASE              
------------+---------+---------------------------------
  my-phrase         2   tennis water wing code window   
                        leaf                            
```

## retrieve phrases
Stored phrases can be listed:
```bash
mkphrase list
                  NAME                  
----------------------------------------
  5095fa73-7f24-4952-a258-968ed1d29e31  
  my-phrase                                                       
```
And any particular phrase value can be fetched, which always fetches the
`latest` version of the named phrase
```bash
mkphrase get 5095fa73-7f24-4952-a258-968ed1d29e31
                  NAME                   VERSION          PHRASE          
---------------------------------------+---------+------------------------
  5095fa73-7f24-4952-a258-968ed1d29e31         1   this is a test phrase  
```
A specific version can be fetched using `--version` flag:
```bash
mkphrase get my-phrase --version=2
    NAME      VERSION               PHRASE              
------------+---------+---------------------------------
  my-phrase         2   tennis water wing code window   
                        leaf                            
```

## delete phrase
When a named phrase is deleted, all versions of secret material are 
deleted forever.
> Please use caution when using this command
```bash
mkphrase delete 5095fa73-7f24-4952-a258-968ed1d29e31
Type secret name to delete: 5095fa73-7f24-4952-a258-968ed1d29e31
```
The above command will ask for confirmation, however, `--force` option
can be used to skip the confirmation and delete the secret without any
confirmation.

```bash
mkphrase delete 5095fa73-7f24-4952-a258-968ed1d29e31 --force
```