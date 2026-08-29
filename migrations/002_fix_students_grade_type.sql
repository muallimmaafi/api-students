-- Kolom grade seharusnya bertipe numerik (IPK/nilai), bukan teks.
-- Tabel masih kosong sehingga konversi tipe ini aman dilakukan.
ALTER TABLE students
    ALTER COLUMN grade TYPE DOUBLE PRECISION USING grade::double precision;