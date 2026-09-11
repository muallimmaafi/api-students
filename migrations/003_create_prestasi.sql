CREATE TABLE IF NOT EXISTS prestasi (
    id_prestasi SERIAL PRIMARY KEY,
    nama_prestasi VARCHAR(150) NOT NULL,
    id_student INTEGER NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    juara VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS prestasi_id_student_idx ON prestasi (id_student);