#!/bin/sh

# Using :+
DISABLE_SSL="no"
echo ${DISABLE_SSL:+-k} # Outputs: -k
DISABLE_SSL=""
echo ${DISABLE_SSL:+-k} # Outputs: nothing

# Using :-
DISABLE_SSL="no"
echo "${DISABLE_SSL:--k}" # Outputs: yes
DISABLE_SSL=""
echo "${DISABLE_SSL:--k}" # Outputs: -k
