package utils

import "testing"

func TestDownloadPiscToTmp(t *testing.T) {
	type args struct {
		imgUrl   string
		userName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test download",
			args: args{
				imgUrl:   "https://gimg2.baidu.com/image_search/src=http%3A%2F%2Fcdn.duitang.com%2Fuploads%2Fitem%2F201302%2F04%2F20130204092701_FLrZs.thumb.700_0.jpeg&refer=http%3A%2F%2Fcdn.duitang.com&app=2002&size=f9999,10000&q=a80&n=0&g=0n&fmt=jpeg?sec=1613456524&t=47ae722b9c11bc773a8e38d8b111d3b2",
				userName: "kriswu",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err, _ := DownloadPiscToTmp(tt.args.imgUrl, tt.args.userName); (err != nil) != tt.wantErr {
				t.Errorf("DownloadPiscToTmp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
