/**
 * 中英文切换：写 Cookie（setLocaleAction）并触发服务端重渲染。
 * 后续可挪到 Header；选项来自 i18n/config，勿硬编码语言列表。
 */
'use client';

import { useTransition } from 'react';
import { Segmented } from 'antd';
import { useLocale, useTranslations } from 'next-intl';
import { setLocaleAction } from '@/i18n/actions';
import { localeLabels, locales } from '@/i18n/config';

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
					setLocaleAction(value);
				});
			}}
		/>
	);
}
