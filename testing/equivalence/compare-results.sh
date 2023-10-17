#!/bin/bash

source $(realpath $(dirname $0))/lib.sh

for i in $(find $(realpath $(dirname $0))/results/terraform -mindepth 1 -not -path '**/*/.*' -type d); do
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
        diff \
            $(realpath $(dirname $0))/results/terraform/$(basename ${i})/$(basename ${f}) \
            $(realpath $(dirname $0))/results/opentofu/$(basename ${i})/$(basename ${f}) \
            || print_style "FAILED\n" "danger" \
            && print_style "PASS\n" "success"
    done
done