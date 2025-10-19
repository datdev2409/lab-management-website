package sheets

import (
	"context"
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"
)

type ReportFile struct {
	File *excelize.File
}

func (rf *ReportFile) GetIOReader(ctx context.Context) (io.Reader, error) {
	f := rf.File
	r, w := io.Pipe()
	go func() {
		defer func() {
			_ = f.Close()
			_ = w.Close()
		}()
		if err := f.Write(w); err != nil {
			// Error during write (e.g., I/O error): close writer with error
			_ = w.CloseWithError(fmt.Errorf("excelize write error: %w", err))
			return
		}
	}()

	return r, nil
}
