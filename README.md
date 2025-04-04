# Govertly

## Dependencies

- https://github.com/FiloSottile/age
  - https://htmlpreview.github.io/?https://github.com/FiloSottile/age/blob/main/doc/age.1.html
- https://github.com/jedisct1/libsodium ?

## Why yet another password manager?

- The simplicity of pass without pgp [problems](https://www.latacora.com/blog/2019/07/16/the-pgp-problem/)

## COMMANDS
### init [ --path=sub-folder, -p sub-folder ] [ --key=key-file, -k key-file ] [ --git, -g ] [ --remote=repo, -r repo ]

Initialize new password storage and use key-file for encryption. Multiple keys may be specified in the key-file, in order to encrypt each password with multiple ids. This command must be run first before a password store can be used. If the specified key-file is different from the key used in any existing files, these files will be reencrypted to use the new id. If --path or -p is specified, along with an argument, a specific key-file or set of key-files is assigned for that specific sub folder of the password store. If only one key-file is given, and it is an empty string, then the current .key-file file for the specified sub-folder (or root if unspecified) is removed. If --git or -g is specified, creates local git repo, unless --remote or -r is informed with a remote repo, which implies the git option.

### ls subfolder

List names of passwords inside the tree at subfolder. This command is alternatively named list.

### grep [GREPOPTIONS] search-string

Searches inside each decrypted password file for search-string, and displays line containing matched string along with filename. GREPOPTIONS are passed to grep(1) as-is. (Note: the GREP_OPTIONS environment variable functions as well.)

### find govertly-names...

List names of passwords inside the tree that match govertly-names by using the tree(1) program. This command is alternatively named search.

### show [ --clip[=line-number], -c[line-number] ] [--qrcode[=line-number], -q[line-number] ] govertly-name

Decrypt and print a password named govertly-name. If --clip or -c is specified, do not print the password but instead copy the first (or otherwise specified) line to the clipboard using xclip(1) or wl-clipboard(1) and then restore the clipboard after 45 (or GOVERTLY_CLIP_TIME) seconds. If --qrcode or -q is specified, do not print the password but instead display a QR code using qrencode(1) either to the terminal or graphically if supported.

### insert [ --echo, -e | --multiline, -m ] [ --force, -f ] govertly-name

Insert a new password into the password store called govertly-name. This will read the new password from standard in. If --echo or -e is not specified, disable keyboard echo when the password is entered and confirm the password by asking for it twice. If --multiline or -m is specified, lines will be read until EOF or Ctrl+D is reached. Otherwise, only a single line from standard in is read. Prompt before overwriting an existing password, unless --force or -f is specified. This command is alternatively named add.

### edit govertly-name

Insert a new password or edit an existing password using the default text editor specified by the environment variable EDITOR or using vi(1) as a fallback. This mode makes use of temporary files for editing, but care is taken to ensure that temporary files are created in /dev/shm in order to avoid writing to difficult-to-erase disk sectors. If /dev/shm is not accessible, fallback to the ordinary TMPDIR location, and print a warning.

### generate [ --no-symbols, -n ] [ --clip, -c ] [ --in-place, -i | --force, -f ] govertly-name [govertly-length]

Generate a new password of length govertly-length (or GOVERTLY_GENERATED_LENGTH if unspecified) and insert into govertly-name. If --no-symbols or -n is specified, do not use any non-alphanumeric characters in the generated password. The character sets used in generating passwords can be changed with the GOVERTLY_CHARACTER_SET and GOVERTLY_CHARACTER_SET_NO_SYMBOLS environment variables, described below. If --clip or -c is specified, do not print the password but instead copy it to the clipboard using xclip(1) or wl-clipboard(1) and then restore the clipboard after 45 (or GOVERTLY_CLIP_TIME) seconds. If --qrcode or -q is specified, do not print the password but instead display a QR code using qrencode(1) either to the terminal or graphically if supported. Prompt before overwriting an existing password, unless --force or -f is specified. If --in-place or -i is specified, do not interactively prompt, and only replace the first line of the password file with the new generated password, keeping the remainder of the file intact.

### rm [ --recursive, -r ] [ --force, -f ] govertly-name

Remove the password named govertly-name from the password store. This command is alternatively named remove or delete. If --recursive or -r is specified, delete govertly-name recursively if it is a directory. If --force or -f is specified, do not interactively prompt before removal.

### mv [ --force, -f ] old-path new-path

Renames the password or directory named old-path to new-path. This command is alternatively named rename. If --force is specified, silently overwrite new-path if it exists. If new-path ends in a trailing /, it is always treated as a directory. Passwords are selectively reencrypted to the corresponding keys of their new destination.

### cp [ --force, -f ] old-path new-path

Copies the password or directory named old-path to new-path. This command is alternatively named copy. If --force is specified, silently overwrite new-path if it exists. If new-path ends in a trailing /, it is always treated as a directory. Passwords are selectively reencrypted to the corresponding keys of their new destination.

### git git-command-args...

If the password store is a git repository, govertly git-command-args as arguments to git(1) using the password store as the git repository. If git-command-args is init, in addition to initializing the git repository, add the current contents of the password store to the repository in an initial commit. If the git config key govertly.signcommits is set to true, then all commits will be signed using user.signingkey or the default git signing key. This config key may be turned on using: ‘govertly git config --bool --add govertly.signcommits true‘

### help

Show usage message.

### version

Show version information.

## Configuration

FILES
˜/.local/share/govertly

The default password storage directory.

˜/.config/govertly/.key

Contains the default key identification used for encryption and decryption. Multiple age keys may be specified in this file, one per line. If this file exists in any sub directories, passwords inside those sub directories are encrypted using those keys. This should be set using the init command.

˜/.local/share/govertly/.extensions

The directory containing extension files.

ENVIRONMENT VARIABLES
GOVERTLY_DIR

Overrides the default password storage directory.

GOVERTLY_KEY

Overrides the default age key identification set by init. Keys must not contain spaces and thus use of the hexadecimal key signature is recommended. Multiple keys may be specified separated by spaces.

GOVERTLY_age_OPTS

Additional options to be passed to all invocations of age.

GOVERTLY_X_SELECTION

Overrides the selection passed to xclip, by default clipboard. See xclip(1) for more info.

GOVERTLY_CLIP_TIME

Specifies the number of seconds to wait before restoring the clipboard, by default 45 seconds.

GOVERTLY_UMASK

Sets the umask of all files modified by govertly, by default 077.

GOVERTLY_GENERATED_LENGTH

The default password length if the govertly-length parameter to generate is unspecified.

GOVERTLY_CHARACTER_SET

The character set to be used in password generation for generate. This value is to be interpreted by tr. See tr(1) for more info.

GOVERTLY_CHARACTER_SET_NO_SYMBOLS

The character set to be used in no-symbol password generation for generate, when --no-symbols, -n is specified. This value is to be interpreted by tr. See tr(1) for more info.

GOVERTLY_ENABLE_EXTENSIONS

This environment variable must be set to "true" for extensions to be enabled.

GOVERTLY_EXTENSIONS_DIR

The location to look for executable extension files, by default GOVERTLY_DIR/.extensions.

GOVERTLY_SIGNING_KEY

If this environment variable is set, then all .age-id files and non-system extension files must be signed using a detached signature using the age key specified by the full 40 character upper-case fingerprint in this variable. If multiple fingerprints are specified, each separated by a whitespace character, then signatures must match at least one. The init command will keep signatures of .age-id files up to date.

EDITOR

The location of the text editor used by edit.

