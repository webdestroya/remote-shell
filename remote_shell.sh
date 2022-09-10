#!/bin/bash

set -e

script_path=$(realpath $0)
script_dir=$(dirname $script_path)

echo "Assuming Cloud87 Remote shell folder of: ${script_dir}"
echo ""

while getopts u:e:p:i:k:m: flag
do
  case "${flag}" in
    u) opt_github=${OPTARG};;
    e) opt_envname=${OPTARG};;
    p) opt_port=${OPTARG};;
    i) opt_idletimeout=${OPTARG};;
    k) opt_keepalive=${OPTARG};;
    m) opt_maxruntime=${OPTARG};;
  esac
done

env_name=${opt_envname-unknown}
db_idle_timeout=${opt_idletimeout-${C87_RSHELL_IDLE_TIMEOUT-0}}
db_keepalive=${opt_keepalive-${C87_RSHELL_KEEPALIVE-300}}
db_port=${opt_port-${C87_RSHELL_PORT-8722}}
maxruntime=${opt_maxruntime-${C87_RSHELL_MAX_RUNTIME-43200}}
githubuser=${opt_github-$C87_RSHELL_GH_USER}


if [ -z "$githubuser" ]; then
  echo "ERROR: GitHub username not provided!"
  exit 1;
fi

echo "GitHub User:  ${githubuser}"
echo "Port:         ${db_port}"
echo "Idle Timeout: ${db_idle_timeout}"
echo "Keepalive:    ${db_keepalive}"
echo "Max Runtime:  ${maxruntime}"

set -u
set -o pipefail

mkdir -p $HOME/.ssh
chmod 0600 $HOME/.ssh
mkdir -p /etc/dropbear


# to ensure that the environment is properly loaded
printenv | sort | grep -v -E '^(_|DISPLAY|MAIL|USER|TERM|HOME|LOGNAME|SHELL|SHLVL|PWD|SSH_.+)=' > $HOME/.ssh/environment

# echo "" >> /etc/environment
# cat $HOME/.ssh/environment >> /etc/environment
# echo "export \$(cat ${HOME}/.ssh/environment | sed 's/#.*//g' | xargs)" >> $HOME/.bashrc


# setup the authorized keys
echo "" >> $HOME/.ssh/authorized_keys

curl -sSL https://api.github.com/users/${githubuser}/keys | ${script_dir}/jq -rcM ".[] | .key" >> $HOME/.ssh/authorized_keys
echo "" >> $HOME/.ssh/authorized_keys
chmod 0600 $HOME/.ssh/authorized_keys

echo "Authorized Keys:"
cat $HOME/.ssh/authorized_keys


echo "" >> $HOME/.bashrc
echo "cd $(pwd)" >> $HOME/.bashrc

# timeout \
#   --preserve-status \
#   --kill-after=10s ${maxruntime} \
#     ${script_dir}/dropbear -v -F -R -s -g \
#       -I ${db_idle_timeout} \
#       -K ${db_keepalive} \
#       -p ${db_port}

strace ${script_dir}/dropbear -vvvv -F -R -s -g \
      -I ${db_idle_timeout} \
      -K ${db_keepalive} \
      -p ${db_port}