'use client';

/**
 * Stripe Checkout 回跳提示壳：按状态强化成功 / 取消 / 异常视觉，无跳转按钮。
 * success / cancel 共用。
 */

import { Spin } from 'antd';
import { CheckCircleFilled, CloseCircleFilled, ExclamationCircleFilled } from '@ant-design/icons';

/** @typedef {'processing' | 'success' | 'cancel' | 'error'} BillingCallbackStatus */

const STATUS_THEME = {
	processing: {
		panelClass: 'border-base-active bg-base-active',
		iconWrapClass: 'bg-white text-base-primary',
		titleClass: 'text-base-text',
		icon: null,
	},
	success: {
		panelClass: 'bg-green-light',
		panelStyle: { borderColor: 'var(--green-light-100)' },
		iconWrapClass: 'bg-white text-base-success',
		titleClass: 'text-base-success',
		icon: <CheckCircleFilled className="text-[36px]" />,
	},
	cancel: {
		panelClass: 'border-base-border bg-gray-light',
		iconWrapClass: 'bg-white text-base-label',
		titleClass: 'text-base-text',
		icon: <CloseCircleFilled className="text-[36px]" />,
	},
	error: {
		panelClass: '',
		panelStyle: {
			backgroundColor: 'var(--orange-light)',
			borderColor: 'var(--orange-medium)',
		},
		iconWrapClass: 'bg-white',
		iconWrapStyle: { color: 'var(--orange-dark)' },
		titleClass: '',
		titleStyle: { color: 'var(--orange-dark)' },
		icon: <ExclamationCircleFilled className="text-[36px]" />,
	},
};

/**
 * @param {object} props
 * @param {BillingCallbackStatus} props.status
 * @param {string} props.title
 * @param {string} [props.description]
 */
export default function BillingCallbackView({ status, title, description }) {
	const theme = STATUS_THEME[status] ?? STATUS_THEME.error;

	return (
		<div className="flex h-full w-full items-center justify-center bg-background px-6 py-12">
			<div
				className={[
					'w-full max-w-[480px] rounded-2xl border px-8 py-10 text-center',
					theme.panelClass,
				]
					.filter(Boolean)
					.join(' ')}
				style={theme.panelStyle}
			>
				<div className="mb-5 flex justify-center">
					{status === 'processing' ? (
						<div
							className={[
								'flex h-16 w-16 items-center justify-center rounded-full',
								theme.iconWrapClass,
							].join(' ')}
						>
							<Spin size="large" />
						</div>
					) : (
						<div
							className={[
								'flex h-16 w-16 items-center justify-center rounded-full shadow-sm',
								theme.iconWrapClass,
							]
								.filter(Boolean)
								.join(' ')}
							style={theme.iconWrapStyle}
						>
							{theme.icon}
						</div>
					)}
				</div>

				<h1
					className={['m-0 text-[22px] font-medium leading-8', theme.titleClass]
						.filter(Boolean)
						.join(' ')}
					style={theme.titleStyle}
				>
					{title}
				</h1>

				{description ? (
					<p className="mt-3 mb-0 text-sm leading-6 text-base-gray">{description}</p>
				) : null}
			</div>
		</div>
	);
}
