#!/bin/bash

source $(realpath $(dirname $0))/lib.sh

for i in $(find $(realpath $(dirname $0))/results/terraform -mindepth 1 -maxdepth 1 -type d); do
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

        terraform_file=$(realpath $(dirname $0))/results/terraform/$(basename ${i})/$(basename ${f})
        opentofu_file=$(realpath $(dirname $0))/results/opentofu/$(basename ${i})/$(basename ${f})

        # We normalize the JSON output to make them comparable regardless of key order
        diff \
            <(cat $terraform_file | jq -reM '""' 2>/dev/null || cat $terraform_file) \
            <(cat $terraform_file | jq -reM '""' 2>/dev/null || cat $opentofu_file) \
            || print_style "FAILED\n" "danger" \
            && print_style "PASS\n" "success"
    done
done