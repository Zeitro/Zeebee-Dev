/**
 * Pricing 页底部卖点与联系入口；静态文案，不依赖套餐接口。
 */

import {
	CustomerServiceOutlined,
	SafetyCertificateOutlined,
	SyncOutlined,
} from '@ant-design/icons';
import { DEFAULT_SALES_EMAIL } from '@repo/shared-methods/componentsUi/Pricing/constants';

const HIGHLIGHTS = [
	{
		icon: SyncOutlined,
		title: 'Flexible billing',
		description:
			'Upgrade, downgrade, or cancel anytime. Your team keeps access through the current billing period.',
	},
	{
		icon: SafetyCertificateOutlined,
		title: 'Secure checkout',
		description:
			'Payments are processed through Stripe with industry-standard encryption and compliance.',
	},
	{
		icon: CustomerServiceOutlined,
		title: 'Dedicated support',
		description: 'Need a custom plan or have billing questions? Our team is ready to help.',
	},
];

export default function PricingHighlights() {
	return (
		<section className="mt-10 md:mt-14">
			<div className="grid grid-cols-1 gap-4 md:grid-cols-3 md:gap-5">
				{HIGHLIGHTS.map(({ icon: Icon, title, description }) => (
					<div
						key={title}
						className="rounded-xl border border-base-border bg-background px-5 py-5 shadow-sm"
					>
						<div className="mb-3 flex h-10 w-10 items-center justify-center rounded-lg bg-base-active text-base-primary">
							<Icon className="text-lg" />
						</div>
						<h3 className="m-0 text-base font-medium leading-6 text-base-text">
							{title}
						</h3>
						<p className="m-0 mt-2 text-sm leading-6 text-base-gray">{description}</p>
					</div>
				))}
			</div>

			<p className="m-0 mt-8 text-center text-sm leading-6 text-base-label">
				Questions about billing?{' '}
				<a
					href={`mailto:${DEFAULT_SALES_EMAIL}`}
					className="font-medium text-base-primary hover:underline"
				>
					{DEFAULT_SALES_EMAIL}
				</a>
			</p>
		</section>
	);
}
