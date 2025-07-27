<main>
 <h1>Parquet: What it is and why it matters</h1>

  <p>
    <strong>Parquet</strong> is a columnar storage file format optimized for big data processing and analytics. Unlike traditional row-based formats like CSV or JSON, Parquet stores data by column, enabling faster queries, better compression, and efficient use of storage.
  </p>

  <hr />

  <h2>Why use Parquet?</h2>
  <ul>
    <li><strong>Columnar Storage:</strong> Only reads the columns needed for a query, reducing I/O and speeding up data access.</li>
    <li><strong>Efficient Compression:</strong> Columnar data compresses much better, saving storage space and improving performance.</li>
    <li><strong>Schema Awareness:</strong> Parquet files store metadata and schema information, enabling type-safe and consistent reads.</li>
    <li><strong>Compatibility:</strong> Widely supported in data processing engines like Apache Spark, DuckDB, and many cloud platforms.</li>
  </ul>

  <h2>Common DuckDB commands with Parquet</h2>
  <p>DuckDB is a powerful embedded SQL database engine optimized for analytical queries. It can query Parquet files directly without importing them.</p>

  <div class="grid">
  <pre class="w-full bg-card px-6 py-4 rounded overflow-x-auto whitespace-pre"><code class="bg-card">-- Read 10 rows from a parquet file
SELECT * FROM 'myfile.parquet' LIMIT 10;

-- Filter rows where 'email' column contains 'john'
SELECT * FROM 'myfile.parquet'
WHERE lower(email) LIKE '%john%';

-- Select specific columns
SELECT name, email FROM 'myfile.parquet'
WHERE email LIKE '%@example.com%';

-- Count rows in the parquet file
SELECT COUNT(*) FROM 'myfile.parquet';

-- Convert a CSV file to Parquet
duckdb -c "CREATE TABLE thetable AS FROM read_csv_auto('./file.csv', ignore_errors=true, all_varchar=true); COPY thetable TO './file.parquet' (FORMAT 'parquet', COMPRESSION 'ZSTD');"

-- Convert a CSV file and rename a column
duckdb -c "CREATE TABLE thetable AS SELECT * EXCLUDE(name), name AS full_name FROM read_csv_auto('./file.csv', ignore_errors=true, all_varchar=true); COPY thetable TO './file.parquet' (FORMAT 'parquet', COMPRESSION 'ZSTD');"

-- Convert a CSV file and merge two columns
duckdb -c "CREATE TABLE thetable AS SELECT *, first_name || ' ' || last_name AS full_name FROM read_csv_auto('./file.csv', ignore_errors=true, all_varchar=true); COPY thetable TO './file.parquet' (FORMAT 'parquet', COMPRESSION 'ZSTD');"
</code></pre>
  </div>

  <h2>Tips for working with Data wells</h2>
  <ul>
  <li>Use <code>all_varchar=true</code> when importing data because the Eleakxir search engine only works with strings.</li>
  <li>Trim down columns to keep only the essentials; for example, remove ID columns as they are often not useful for searches.</li>
  <li>Remove duplicate entries whenever possible to optimize storage and search efficiency.</li>
  <li>Use <code>ZSTD</code> compression if disk space is limited. --This saves space but uses more CPU during compression and decompression.</li>
  <li>Ensure consistent columns across files; for example, merge <code>first_name</code> and <code>last_name</code> into a single <code>full_name</code> column for better search results.</li>
  </ul>
</main>
