#!/bin/bash

if [[ "$OSTYPE" == "darwin"* ]];
then
    if ! command -v grealpath &> /dev/null;
    then
        print_style "ERROR: This script requires grealpath when run on Mac. Please install it with 'brew install coreutils'\n" "danger"
        exit 1
    else
        alias realpath=grealpath
    fi
fi

source $(grealpath $(dirname $0))/lib.sh

for i in $(find $(grealpath $(dirname $0))/results/terraform -mindepth 1 -maxdepth 1 -type d); do
    test_case_name=$(basename ${i})
    print_style "Test case "
    print_style "${test_case_name}" "info"
    print_style ":\n"
    for f in $(find ${i} -type f); do
        print_style "  * comparing "
        print_style "results/terraform/$(basename ${f})" "info"
        print_style " and "
        print_style "results/opentofu/$(basename ${f})" "info"
        print_style ": "

        terraform_file=$(grealpath $(dirname $0))/results/terraform/$(basename ${i})/$(basename ${f})
        opentofu_file=$(grealpath $(dirname $0))/results/opentofu/$(basename ${i})/$(basename ${f})

        # We normalize the JSON output to make them comparable regardless of key order
        diff \
            <(cat $terraform_file | jq -reM '""' 2>/dev/null || cat $terraform_file) \
            <(cat $terraform_file | jq -reM '""' 2>/dev/null || cat $opentofu_file) \
            || print_style "FAILED\n" "danger" \
            && print_style "PASS\n" "success"
    done
done