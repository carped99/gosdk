package entx

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"path"
	"slices"
	"strings"
)

type Executor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func Migrate(ctx context.Context, executor Executor, fsys fs.FS, dir string) error {
	// 디렉토리 내 파일 목록 읽기
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return fmt.Errorf("failed to read directory (%s): %w", dir, err)
	}

	// SQL 파일만 필터링
	var sqlFiles []fs.DirEntry
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".sql") {
			sqlFiles = append(sqlFiles, entry)
		}
	}

	// 파일 이름으로 정렬
	slices.SortFunc(entries, func(a, b fs.DirEntry) int {
		return strings.Compare(a.Name(), b.Name())
	})

	// 각 파일 실행
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := path.Join(dir, entry.Name())
		content, err := fs.ReadFile(fsys, filePath)
		if err != nil {
			return fmt.Errorf("failed to read file (%s): %w", filePath, err)
		}

		if _, err := executor.ExecContext(ctx, string(content)); err != nil {
			return fmt.Errorf("execution failed (%s): %w", filePath, err)
		}
	}

	return nil
}

//
//func RunAtlasMigration(ctx context.Context, cfg Config, fsys fs.FS, dir string) error {
//	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
//	defer cancel()
//
//	// 연결 생성
//	client, err := sqlclient.Open(ctx, cfg.DSN)
//	if err != nil {
//		return fmt.Errorf("failed to open db: %w", err)
//	}
//	defer func() {
//		if err := client.Close(); err != nil {
//			fmt.Printf("failed to close db: %v\n", err)
//		}
//	}()
//
//	tempDir, err := util.CopyToTempDir(fsys, dir)
//	if err != nil {
//		return fmt.Errorf("failed to copy migration files to temp dir: %w", err)
//	}
//
//	defer func(path string) {
//		err := os.RemoveAll(path)
//		if err != nil {
//			slog.Error("failed to remove temp dir", "path", path, "error", err)
//		} else {
//			slog.Info("removed temp dir", "path", path)
//		}
//	}(tempDir) // cleanup
//
//	// 4. 마이그레이션 디렉토리 로드
//	localDir, err := migrate.NewLocalDir(tempDir)
//	if err != nil {
//		return fmt.Errorf("failed to create local migration dir: %w", err)
//	}
//
//	// 5. Executor 생성
//	executor, err := migrate.NewExecutor(client.Driver, localDir, migrate.NopRevisionReadWriter{})
//	if err != nil {
//		return fmt.Errorf("failed to create migration executor: %w", err)
//	}
//
//	// 6. 마이그레이션 실행
//	realm, err := executor.Replay(ctx, localDir)
//	if err != nil {
//		return fmt.Errorf("migration replay failed: %w", err)
//	}
//
//	return nil
//}
//
//func extractToTempDir(fsys embed.FS, root string) (string, error) {
//	tmpDir, err := os.MkdirTemp("", "atlas-migrations-*")
//	if err != nil {
//		return "", err
//	}
//
//	err = fs.WalkDir(fsys, root, func(p string, d fs.DirEntry, err error) error {
//		if err != nil {
//			return err
//		}
//		if d.IsDir() {
//			return nil
//		}
//		data, err := fsys.ReadFile(p)
//		if err != nil {
//			return err
//		}
//		rel := strings.TrimPrefix(p, root+"/")
//		dest := filepath.Join(tmpDir, rel)
//		if err := os.WriteFile(dest, data, 0o644); err != nil {
//			return err
//		}
//		return nil
//	})
//	if err != nil {
//		return "", err
//	}
//	return tmpDir, nil
//}
