# Eleakxir

## Prompt

## Cleaning the csv files

```bash
nix-shell -p csvkit csvtk
# First, check that the header is correct:
input="myfile.csv" &&\
echo "File: $input" &&\
echo "Number of line: $(wc -l $input)" &&\
echo "Header:" &&\
cat "$input" | head -n1

# CSV Correction:
csvclean -d "," --fill-short-rows --omit-error-rows $input > clean-$input
csvtk uniq clean-$input -o uniq-$input -j 16

# Check que la première colonne c'est que des emails:
tail -n +2 fichier.csv | cut -d',' -f1 | grep -vE '^[^@]+@[^@]+\.[^@]+$' | wc -l

# Changer la première colonne:
{ echo "nouvelle_premiere_ligne"; tail -n +2 fichier.csv; } > tmp 
mv tmp fichier.csv
```

## Todolist

- Responsive
- README
- Nix developpe launch bun server
- ingest new dataleaks

- Ghunt, gh-recon, sherlock & holehe implementation
