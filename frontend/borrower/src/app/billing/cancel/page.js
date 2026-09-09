'use client';

/**
 * Stripe Checkout 取消回跳页；与 success 共用 BillingCallbackView。
 */

import BillingCallbackView from '../components/BillingCallbackView';

export default function BillingCancelPage() {
	return (
		<BillingCallbackView
			status="cancel"
			title="Checkout cancelled"
			description="Nothing was charged. You can subscribe anytime from Pricing."
		/>
	);
}
