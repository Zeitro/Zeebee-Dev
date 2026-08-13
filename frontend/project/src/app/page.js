import { getTranslations } from 'next-intl/server';

/**
 * 登录后落地页（非 Dashboard）：问候 + 当前组织 + 简短说明。
 * 导航交给侧边栏；套餐不在此展示。
 */
const Home = async () => {
	const t = await getTranslations('home');

	return (
		<div className="flex h-full w-full items-center justify-center overflow-y-auto bg-base-text2 px-6 py-12">
			{t('title')}
		</div>
	);
};

export default Home;
