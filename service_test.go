package main

import (
	"reflect"
	"testing"
)

func Test_service_search(t *testing.T) {
	type fields struct {
		domain string
	}
	type args struct {
		query searchQuery
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []series
		wantErr bool
	}{
		{
			name: "Normal",
			fields: fields{
				domain: "www.rrys2020.com",
			},
			args: args{
				query: searchQuery{
					Keyword: "权力的游戏",
					Details: false,
				},
			},
			want: []series{
				{
					ID:     "10733",
					CnName: "《权力的游戏》(Game of Thrones)[冰与火之歌 / 权力的游戏下载 / 权利的游戏 / 冰火]",
					Poster: "http://tu.jstucdn.com/ftp/2019/0322/d2b4282fe50dffaad4c73b6f3d6176ff.jpg",
				},
				{
					ID:     "35844",
					CnName: "《权力的游戏：征服与反抗》(Game of Thrones: Conquest and Rebellion)",
					Poster: "http://tu.jstucdn.com/ftp/2017/1214/754eb87fb49adbadcbbe46348370ff73.jpg",
				},
			},
			wantErr: false,
		},
		{
			name: "Domaind Deprecated",
			fields: fields{
				domain: "www.rrys2019.com",
			},
			args: args{
				query: searchQuery{
					Keyword: "权力的游戏",
					Details: false,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &service{
				domain: tt.fields.domain,
			}
			got, err := svc.search(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_detail(t *testing.T) {
	type fields struct {
		domain string
	}
	type args struct {
		seriesID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    series
		wantErr bool
	}{
		{
			name: "Normal",
			fields: fields{
				domain: "www.rrys2020.com",
			},
			args: args{
				seriesID: "10733",
			},
			want: series{
				ID:       "10733",
				CnName:   "权力的游戏",
				Poster:   "http://tu.jstucdn.com/ftp/2019/0322/d2b4282fe50dffaad4c73b6f3d6176ff.jpg",
				EnName:   "Game of Thrones",
				Link:     "http://www.rrys2020.com/resource/10733",
				RssLink:  "http://rss.rrys.tv/rss/feed/10733",
				Area:     "美国",
				Category: "战争/剧情/魔幻/历史/古装/史诗",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &service{
				domain: tt.fields.domain,
			}
			got, err := svc.detail(tt.args.seriesID)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.detail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.detail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_episodes(t *testing.T) {
	type fields struct {
		domain string
	}
	type args struct {
		seriesID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []episode
		wantErr bool
	}{
		{
			name: "Normal",
			fields: fields{
				domain: "www.rrys2020.com",
			},
			args: args{
				seriesID: "40595",
			},
			want: []episode{
				{
					SeriesID: "40595",
					Name:     "异星灾变.Raised.by.Wolves.2020.S01E01.WEBrip.720P-人人影视.V2.mp4",
					Season:   1,
					Episode:  1,
					Magnet:   "magnet:?xt=urn:btih:28f95d483107488fdf39dbe1a802649eed4b0549&tr=http://tr.cili001.com:8070/announce&tr=udp://p4p.arenabg.com:1337&tr=udp://tracker.opentrackr.org:1337/announce&tr=udp://open.demonii.com:1337",
				},
				{
					SeriesID: "40595",
					Name:     "异星灾变.Raised.by.Wolves.2020.S01E02.WEBrip.720P-人人影视.V2.mp4",
					Season:   1,
					Episode:  2,
					Magnet:   "magnet:?xt=urn:btih:8d19a84f5769ebafdc5981abb320aa9baab669e5&tr=http://tr.cili001.com:8070/announce&tr=udp://p4p.arenabg.com:1337&tr=udp://tracker.opentrackr.org:1337/announce&tr=udp://open.demonii.com:1337",
				},
				{
					SeriesID: "40595",
					Name:     "异星灾变.Raised.by.Wolves.2020.S01E03.WEBrip.720P-人人影视.mp4",
					Season:   1,
					Episode:  3,
					Magnet:   "magnet:?xt=urn:btih:86cca1c3910c3f233fe17d3ed6fc3e33fc9e6bb5&tr=http://tr.cili001.com:8070/announce&tr=udp://p4p.arenabg.com:1337&tr=udp://tracker.opentrackr.org:1337/announce&tr=udp://open.demonii.com:1337",
				},
				{
					SeriesID: "40595",
					Name:     "异星灾变.Raised.by.Wolves.2020.S01E04.WEBrip.720P-人人影视.V2.mp4",
					Season:   1,
					Episode:  4,
					Magnet:   "magnet:?xt=urn:btih:b9b881e459d3989f7cfa042f04d1ba5720ef758f&tr=http://tr.cili001.com:8070/announce&tr=udp://p4p.arenabg.com:1337&tr=udp://tracker.opentrackr.org:1337/announce&tr=udp://open.demonii.com:1337",
				},
				{
					SeriesID: "40595",
					Name:     "异星灾变.Raised.by.Wolves.2020.S01E05.WEBrip.720P-人人影视.mp4",
					Season:   1,
					Episode:  5,
					Magnet:   "magnet:?xt=urn:btih:ab78ca780b45bcc9e6ee696d8f0a10a04b72e727&tr=http://tr.cili001.com:8070/announce&tr=udp://p4p.arenabg.com:1337&tr=udp://tracker.opentrackr.org:1337/announce&tr=udp://open.demonii.com:1337",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &service{
				domain: tt.fields.domain,
			}
			got, err := svc.episodes(tt.args.seriesID)
			if (err != nil) != tt.wantErr {
				t.Errorf("service.episodes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("service.episodes() = %v, want %v", got, tt.want)
			}
		})
	}
}
