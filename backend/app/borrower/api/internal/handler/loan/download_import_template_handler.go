package loan

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"

	"zeebee-backend/app/borrower/api/internal/svc"
	"zeebee-backend/common/response"

	"zeebee-backend/common/errorx"
)

// DownloadImportTemplateHandler download import template
func DownloadImportTemplateHandler(_ *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			if err != nil {
				logx.WithContext(r.Context()).
					WithFields(logx.Field("stacks", errorx.GetErrorStack(err))).
					Error(err)
			}
		}()

		templatePath := "./statics/template/borrower_loans.xlsx"

		// check file
		if _, err := os.Stat(templatePath); os.IsNotExist(err) {
			response.BadRequest(w, "Template file not found")
			return
		}

		file, err := os.Open(templatePath)
		if err != nil {
			response.BadRequest(w, "Failed to open template file")
			return
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			response.BadRequest(w, "Failed to get template file info")
			return
		}

		// set header
		w.Header().Set("Content-Disposition", "attachment; filename=borrower_loans.xlsx")
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

		_, err = io.Copy(w, file)
		if err != nil {
			response.BadRequest(w, "Failed to send template file")
			return
		}
		return
	}
}
