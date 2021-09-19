# Records

The following snippet from the RECORDS.csv contains the Vivosport records for 1km, 1mi, and 5km.

```csv
Type,Local Number,Message,Field 1,Value 1,Units 1,Field 2,Value 2,Units 2,Field 3,Value 3,Units 3,Field 4,Value 4,Units 4,Field 5,Value 5,Units 5,Field 6,Value 6,Units 6,Field 7,Value 7,Units 7,Field 8,Value 8,Units 8,Field 9,Value 9,Units 9,Field 10,Value 10,Units 10,
Data,2,unknown,unknown,"1000886458",,unknown,"100000",,unknown,"99500",,unknown,"101600",,unknown,"273036",,unknown,"2",,unknown,"0",,unknown,"0",,unknown,"1",,,,,
Data,2,unknown,unknown,"1000886426",,unknown,"160900",,unknown,"160100",,unknown,"163500",,unknown,"451672",,unknown,"2",,unknown,"1",,unknown,"0",,unknown,"1",,,,,
Data,2,unknown,unknown,"1000886425",,unknown,"500000",,unknown,"498000",,unknown,"507900",,unknown,"1528180",,unknown,"2",,unknown,"2",,unknown,"0",,unknown,"1",,,,,
```

The distance is represented by the `Value 2` field, above these are "10000", "160900", and "50000". Although the units are omitted, it appears they are in centimeters.

The time is represented by the `Value 5` field, above these are "270036", "451672", and "1528180". Again, units are omitted but they are in milliseconds.

Values in `Value 3` and `Value 4` fields appear to be upper and lower bounds for the distance, also in centib
meters.
