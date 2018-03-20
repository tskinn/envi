#!/bin/bash

set -e

cd store
go test
printf "\ngo tests passed!\n\n"
cd ..

set +e

fail() {
		echo "Failure"
		exit 1
}

ID="supertest__somegarbageid"

# test setting
echo "Testing 'set' command..."
if ! go run *.go s -i ${ID} -e one=two,three=four; then
		fail
fi
printf "\t'set' succesfully tested.\n"

# test getting
echo "Testing 'get' command..."
results=$(go run *.go g -i ${ID})
success=$?
if [ ${success} -ne 0 ]; then
		fail
fi
if ! echo "${results}" | grep -q "one=two"; then
		fail
fi
if ! echo "${results}" | grep -q "three=four"; then
		fail
fi
printf "\t'get' successfully tested.\n"


# test updating
echo "Testing 'update' command..."
if ! go run *.go u -i ${ID} -e one=one; then
		fail
fi
results=$(go run *.go g -i ${ID})
if ! echo "${results}" | grep -q "one=one"; then
		fail
fi
if ! echo "${results}" | grep -q "three=four"; then
		fail
fi

printf "\t'update' succesfully tested.\n"

# test deleting variable
echo "Testing 'delete' variable command..."
go run *.go d -i ${ID} -e one
results=$(go run *.go g -i ${ID})
success=$?
if [ ${success} -ne 0 ]; then
		fail
fi

lines=$(echo "${results}" | wc -l)
if [ ${lines} -ne 1 ]; then
	fail
fi

printf "\t'delete' variable succesfully tested.\n"

# test deleting config
echo "Testing 'delete' command..."
if ! go run *.go d -i ${ID}; then
		fail
fi
results=$(go run *.go g -i ${ID})
success=$?
if [ ${success} -ne 0 ]; then
		fail
fi

lines=$(echo "${results}" | wc -l)
if [ ${lines} -eq 0 ]; then
		fail
fi
printf "\t'delete' succesfully tested.\n"

echo "Tests successful"
