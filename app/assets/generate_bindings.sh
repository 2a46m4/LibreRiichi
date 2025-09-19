#!/bin/sh

DIR="/riichi-mahjong-tiles/"
TYPES="Black Regular"
SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
OUT="$SCRIPT_DIR/image-imports.tsx"
EXPORTS=""

echo "// Auto-generated file" > "$OUT"

for type in $TYPES; do
for file in ".$DIR""$type"/*; do
  filename=$(basename "$file")
  filename_no_ext=${filename%.*}
  filename_fix=$(echo "$filename_no_ext" | tr '-' '_')
  echo "import tile_""$type""_$filename_fix from '$file';" >> "$OUT"
  EXPORTS=tile_"$type"_"$filename_fix, $EXPORTS"
done
done

echo "export { $EXPORTS }" >> "$OUT"