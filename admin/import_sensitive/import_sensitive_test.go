package import_sensitive

import "testing"

func TestImportSensitiveWordFromExcel(t *testing.T) {
	type args struct {
		excelPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				excelPath: "./OnBio_sensitive.xlsx",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ImportSensitiveWordFromExcel(tt.args.excelPath); (err != nil) != tt.wantErr {
				t.Errorf("ImportSensitiveWordFromExcel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
