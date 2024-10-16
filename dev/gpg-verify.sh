# ./gpg-verify.sh <filename-to-verify>
# its signature should be in the same folder with name <filename-to-verify>.sig

gpg --verify $1.sig $1
