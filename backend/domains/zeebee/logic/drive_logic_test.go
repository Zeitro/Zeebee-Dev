package logic

import (
	"testing"

	"zeebee-backend/domains/zeebee/constants"
	"zeebee-backend/domains/zeebee/drive/client"
	"zeebee-backend/domains/zeebee/model"
)

func TestBuildTransactionFolderName(t *testing.T) {
	tx := &model.BorrowerPropertyTransaction{
		Code:   "PT-001",
		Type:   constants.BorrowerPropertyTransactionType("Purchase"),
		Status: constants.BorrowerPropertyTransactionStatus("Pending"),
	}
	if got, want := buildTransactionFolderName(tx), "Purchase Transaction PT-001 (Pending)"; got != want {
		t.Fatalf("name = %q, want %q", got, want)
	}
	if got, want := transactionFolderNamePrefix(tx), "Purchase Transaction PT-001 ("; got != want {
		t.Fatalf("prefix = %q, want %q", got, want)
	}

	// Code 为空时优雅降级，不产生多余空格
	txNoCode := &model.BorrowerPropertyTransaction{
		Type:   constants.BorrowerPropertyTransactionType("Refinance"),
		Status: constants.BorrowerPropertyTransactionStatus("Approved"),
	}
	if got, want := buildTransactionFolderName(txNoCode), "Refinance Transaction (Approved)"; got != want {
		t.Fatalf("no-code name = %q, want %q", got, want)
	}
}

func TestPickTransactionFolder(t *testing.T) {
	const target = "Purchase Transaction PT-001 (Approved)"
	const prefix = "Purchase Transaction PT-001 ("

	cases := []struct {
		name       string
		folders    []client.DriveFolderInfo
		wantId     string // "" 表示期望返回 nil（需要新建）
		wantRename bool
	}{
		{
			name:    "无匹配则新建",
			folders: []client.DriveFolderInfo{{Id: "a", Name: "Purchase Transaction PT-999 (Pending)"}},
			wantId:  "",
		},
		{
			name:       "旧状态文件夹改名复用",
			folders:    []client.DriveFolderInfo{{Id: "b", Name: "Purchase Transaction PT-001 (Pending)"}},
			wantId:     "b",
			wantRename: true,
		},
		{
			name:       "当前状态文件夹直接复用",
			folders:    []client.DriveFolderInfo{{Id: "c", Name: target}},
			wantId:     "c",
			wantRename: false,
		},
		{
			name: "同时存在当前与旧状态时优先用当前状态",
			folders: []client.DriveFolderInfo{
				{Id: "old", Name: "Purchase Transaction PT-001 (Pending)"},
				{Id: "cur", Name: target},
			},
			wantId:     "cur",
			wantRename: false,
		},
		{
			name:    "其它交易（不同 code）不误匹配",
			folders: []client.DriveFolderInfo{{Id: "x", Name: "Purchase Transaction PT-002 (Approved)"}},
			wantId:  "",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			folder, rename := pickTransactionFolder(c.folders, target, prefix)
			if c.wantId == "" {
				if folder != nil {
					t.Fatalf("期望 nil（新建），实际得到 %+v", folder)
				}
				return
			}
			if folder == nil {
				t.Fatalf("期望复用 id=%s，实际为 nil", c.wantId)
			}
			if folder.Id != c.wantId || rename != c.wantRename {
				t.Fatalf("got (id=%s, rename=%v), want (id=%s, rename=%v)", folder.Id, rename, c.wantId, c.wantRename)
			}
		})
	}
}
