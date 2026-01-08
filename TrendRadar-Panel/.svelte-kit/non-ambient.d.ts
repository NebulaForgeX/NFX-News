
// this file is generated — do not edit it


declare module "svelte/elements" {
	export interface HTMLAttributes<T> {
		'data-sveltekit-keepfocus'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-noscroll'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-preload-code'?:
			| true
			| ''
			| 'eager'
			| 'viewport'
			| 'hover'
			| 'tap'
			| 'off'
			| undefined
			| null;
		'data-sveltekit-preload-data'?: true | '' | 'hover' | 'tap' | 'off' | undefined | null;
		'data-sveltekit-reload'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-replacestate'?: true | '' | 'off' | undefined | null;
	}
}

export {};


declare module "$app/types" {
	export interface AppTypes {
		RouteId(): "/" | "/focus" | "/hottest" | "/realtime";
		RouteParams(): {
			
		};
		LayoutParams(): {
			"/": Record<string, never>;
			"/focus": Record<string, never>;
			"/hottest": Record<string, never>;
			"/realtime": Record<string, never>
		};
		Pathname(): "/" | "/focus" | "/focus/" | "/hottest" | "/hottest/" | "/realtime" | "/realtime/";
		ResolvedPathname(): `${"" | `/${string}`}${ReturnType<AppTypes['Pathname']>}`;
		Asset(): "/Baloo2-Bold.subset.ttf" | "/apple-touch-icon.png" | "/icon.svg" | "/icons/36kr.png" | "/icons/acfun.png" | "/icons/aljazeeracn.png" | "/icons/baidu.png" | "/icons/bilibili.png" | "/icons/cankaoxiaoxi.png" | "/icons/chongbuluo.png" | "/icons/cls.png" | "/icons/coolapk.png" | "/icons/default.png" | "/icons/douban.png" | "/icons/douyin.png" | "/icons/fastbull.png" | "/icons/gelonghui.png" | "/icons/genshin.png" | "/icons/ghxi.png" | "/icons/github.png" | "/icons/hackernews.png" | "/icons/hellogithub.png" | "/icons/honkai.png" | "/icons/hupu.png" | "/icons/ifeng.png" | "/icons/ithome.png" | "/icons/jianshu.png" | "/icons/jin10.png" | "/icons/juejin.png" | "/icons/kaopu.png" | "/icons/kuaishou.png" | "/icons/linuxdo.png" | "/icons/mktnews.png" | "/icons/nowcoder.png" | "/icons/pcbeta.png" | "/icons/peopledaily.png" | "/icons/producthunt.png" | "/icons/smzdm.png" | "/icons/solidot.png" | "/icons/sputniknewscn.png" | "/icons/sspai.png" | "/icons/starrail.png" | "/icons/steam.png" | "/icons/tencent.png" | "/icons/thepaper.png" | "/icons/tieba.png" | "/icons/toutiao.png" | "/icons/v2ex.png" | "/icons/wallstreetcn.png" | "/icons/weibo.png" | "/icons/weread.png" | "/icons/xueqiu.png" | "/icons/zaobao.png" | "/icons/zhihu.png" | "/robots.txt" | "/vite.svg" | string & {};
	}
}