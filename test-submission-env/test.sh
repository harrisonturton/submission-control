#!/bin/bash

echo "Beginning tests..."
if python student_code.py; then
	echo "Successfully passed all tests."
else
	echo "All tests failed!"
fi
