#!/bin/bash

# Set the output directory
output_dir="release_assets"

# Create the output directory if it doesn't exist
mkdir -p "${output_dir}"

# Generate 200 text files
for i in {1..200}; do
  filename=$(printf "%s/asset_%03d.txt" "${output_dir}" "${i}")
  echo "This is asset file number ${i}" >"${filename}"
  echo "Created ${filename}"
done

echo "Generated 200 text files in the '${output_dir}' directory."
