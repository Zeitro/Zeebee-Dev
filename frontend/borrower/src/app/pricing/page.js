import { Pricing } from '@repo/shared-methods';
import PricingPageShell from './components/PricingPageShell';

export default function PricingPage() {
	return (
		<PricingPageShell>
			<Pricing variant="page" />
		</PricingPageShell>
	);
}
