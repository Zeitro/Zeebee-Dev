/**
 * Pricing 路由页外壳：渐变背景、装饰光斑、页眉引导与底部卖点。
 * 仅包装 {@link Pricing}，不修改共享 Pricing 组件本身。
 */

import { CrownOutlined } from '@ant-design/icons';
import PricingHighlights from '../PricingHighlights';

export default function PricingPageShell({ children }) {
	return (
		<div className="relative flex-1 overflow-y-auto bg-gradient-to-b from-blue-light-1 via-background to-gray-light w-full">
			<div className="relative mx-auto w-full max-w-6xl px-4 py-8 md:px-6 md:py-12">
				<div className="mb-6 text-center md:mb-8">
					<span className="inline-flex items-center gap-1.5 rounded-full border border-base-border bg-background px-3 py-1 text-xs font-medium leading-5 text-base-primary shadow-sm">
						<CrownOutlined className="text-[11px]" />
						Subscription &amp; billing
					</span>
					<p className="m-0 mx-auto mt-3 max-w-xl text-sm leading-6 text-base-gray">
						Pick the plan that matches your pipeline today — scale up as your team
						grows.
					</p>
				</div>

				<div className="rounded-2xl border border-base-border bg-background/90 px-2 py-2 shadow-sm backdrop-blur-sm md:px-4 md:py-4">
					{children}
				</div>

				<PricingHighlights />
			</div>
		</div>
	);
}
