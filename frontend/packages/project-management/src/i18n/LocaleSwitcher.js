/**
 * 项目管理中英文切换：写 project_locale Cookie 并触发服务端重渲染。
 */
'use client';

import { useTransition } from 'react';
import { Segmented } from 'antd';
import { useLocale, useTranslations } from 'next-intl';
import { setProjectLocaleAction } from './actions';
import { localeLabels, locales } from './config';

export default function LocaleSwitcher() {
	const locale = useLocale();
	const t = useTranslations('common');
	const [isPending, startTransition] = useTransition();

	return (
		<Segmented
			size="small"
			disabled={isPending}
			value={locale}
			aria-label={t('language')}
			options={locales.map((value) => ({
				label: localeLabels[value],
				value,
			}))}
			onChange={(value) => {
				startTransition(() => {
					setProjectLocaleAction(value);
				});
			}}
		/>
	);
}
