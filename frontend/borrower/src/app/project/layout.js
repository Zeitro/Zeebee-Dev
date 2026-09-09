import { cookies } from 'next/headers';
import {
	ProjectIntlProvider,
	getProjectMessages,
	readProjectLocaleFromCookies,
	LocaleSwitcher,
} from '@repo/project-management/i18n';
import { NewHeader } from '@repo/shared-methods';

/**
 * 项目管理路由壳：仅 /project 子树启用 PM i18n，不污染 borrower 其余页面。
 */
export default async function ProjectLayout({ children }) {
	const cookieStore = await cookies();
	const locale = readProjectLocaleFromCookies(cookieStore);
	const messages = getProjectMessages(locale);

	return (
		<ProjectIntlProvider locale={locale} messages={messages}>
			<NewHeader isShowCrumbs>
				<LocaleSwitcher />
			</NewHeader>
			{children}
		</ProjectIntlProvider>
	);
}
