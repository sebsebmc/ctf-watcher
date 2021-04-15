# CTF Watcher
Use a git-like CLI to `pull` files for CTF challenges or use in `watch` mode to automatically download new challenges.

# WIP
This is a work in progress, so far only init is implemented.

## CTF Support
Currently only supports CTFd based CTFs, if CTF frameworks want to standardize on an API that would be appreciated :wink:

## Usage
`ctf-watcher init ctf-url username password directory` to initialize a folder for a CTF. 

Once the folder is initialized you may `ctf-watcher list` to get a list of challenges and then `ctf-watcher pull challenge-name` to download the challenge files into a directory named for the challenge.