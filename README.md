# Replicant

Deploy Mottainai environments in your instance with a control repository.

    Mottainai Replicant
    Copyright (c) 2017-2019 Mottainai

    Command line interface for Mottainai replicant

    Usage:
       [flags]
       [command]

    Examples:
    $> replicant -m http://127.0.0.1:8080 environment apply --revision origin/master

    $> replicant -m http://127.0.0.1:8080 environment deploy --revision origin/master


    Available Commands:
      environment Handle replicant environments
      help        Help about any command

    Flags:
      -k, --apikey string    Mottainai API key (default "fb4h3bhgv4421355")
      -h, --help             help for this command
      -m, --master string    MottainaiCI webUI URL (default "http://localhost:8080")
      -p, --profile string   Use specific profile for call API.
          --version          version for this command

    Use " [command] --help" for more information about a command.
    
## How to

Either run commands in the git folder which defines your environment, or specify one with ```--environment``` flag.

Similarly as mottainai-cli, it reads profiles if you don't want to specify manually ```--master``` and ```--apikey```.

### Deploy environment

    replicant environment deploy --revision origin/master
    
Deploys the environment remotely. Specify a revision or it will default to origin/master.

### Destroy remote environments

    replicant environment destroy 
    
It wipes the remote environment.
    
### Apply environment 

    replicant environment apply --revision origin/master
    
Apply and syncronize state remotely. If you perform changes in the git repository, will keep track of the state remotely.
