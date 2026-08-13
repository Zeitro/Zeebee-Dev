import type { Metadata } from 'next';
import { Roboto } from 'next/font/google';
import Script from 'next/script';
import { NextIntlClientProvider } from 'next-intl';
import { getLocale, getMessages, getTranslations } from 'next-intl/server';
import AppProviders from '@/components/AppProviders';
import LocaleSwitcher from '@/components/LocaleSwitcher';
import './globals.css';
import '@ant-design/v5-patch-for-react-19';

export const roboto = Roboto({
	weight: ['400', '500', '700'], // 选择需要的字重
	subsets: ['latin', 'latin-ext'], // 字符集
	display: 'swap', // 字体加载策略
	variable: '--font-roboto', // CSS 变量名
});

export async function generateMetadata(): Promise<Metadata> {
	const t = await getTranslations('meta');
	return {
		title: t('title'),
		description: t('description'),
	};
}

const RootLayout = async ({
	children,
}: Readonly<{
	children: React.ReactNode;
}>) => {
	const locale = await getLocale();
	const messages = await getMessages();

	return (
		<html lang={locale} suppressHydrationWarning>
			<head>
				<Script
					src={`https://maps.googleapis.com/maps/api/js?key=${process.env.NEXT_PUBLIC_GOOGLE_MAPS_API_KEY}&libraries=places&loading=async`}
					strategy="lazyOnload" // 页面交互前加载
				/>
			</head>
			<body className={`${roboto.variable}`} suppressHydrationWarning>
				<NextIntlClientProvider locale={locale} messages={messages}>
					<AppProviders locale={locale}>
						{/* 临时入口：有 Header 后挪到导航区即可 */}
						<div className="fixed right-4 top-4 z-50">
							<LocaleSwitcher />
						</div>
						{children}
					</AppProviders>
				</NextIntlClientProvider>
			</body>
		</html>
	);
};

export default RootLayout;
