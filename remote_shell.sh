#!/bin/sh

set -eu

# override with env var, else use default
maxruntime=${C87_CONSOLE_MAX_RUNTIME-43200}

githubuser=${C87_CONSOLE_GH_USER-webdestroya}

mkdir -p $HOME/.ssh
chmod 0600 $HOME/.ssh


# to ensure that the environment is properly loaded
printenv | sort | grep -v -E '^(_|DISPLAY|MAIL|USER|TERM|HOME|LOGNAME|SHELL|SHLVL|PWD|SSH_.+)=' > $HOME/.ssh/environment

echo "" >> /etc/environment
cat $HOME/.ssh/environment >> /etc/environment
echo "export \$(cat ${HOME}/.ssh/environment | sed 's/#.*//g' | xargs)" >> $HOME/.bashrc


# setup the authorized keys
echo "" >> $HOME/.ssh/authorized_keys

curl -sSL https://api.github.com/users/${githubuser}/keys | jq -rcM ".[] | .key" >> $HOME/.ssh/authorized_keys
echo "" >> $HOME/.ssh/authorized_keys
chmod 0600 $HOME/.ssh/authorized_keys

echo "" >> $HOME/.bashrc
echo "cd $(pwd)" >> $HOME/.bashrc

timeout \
  --preserve-status \
  --kill-after=10s ${maxruntime} \
    /cloud87/dropbear -F -E -s -g -j -k -p 6000