#!/bin/bash
this_file=$(basename $0)
base_dir_len=`expr ${#this_file} - 3`
test_dir=${this_file: 0: `expr $base_dir_len`}
cd $test_dir
./$this_file
cd -
