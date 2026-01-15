import { formatAmount } from '@repo/shared-methods/utils';

const Income = ({ record }) => {
	return (
		<div className="flex justify-between items-center px-6">
			<div>
				<div>
					80% AMI: &nbsp;
					<span className="font-medium">{formatAmount(record.ami * 0.8)}</span>
				</div>
				<div>
					80% MFI: &nbsp;
					<span className="font-medium">{formatAmount(record.mfi * 0.8)}</span>
				</div>
			</div>
			<div>
				<div>
					100% AMI: &nbsp;<span className="font-medium">{formatAmount(record.ami)}</span>
				</div>
				<div>
					100% MFI: &nbsp;<span className="font-medium">{formatAmount(record.mfi)}</span>
				</div>
			</div>
			<div>
				<div>
					150% AMI: &nbsp;
					<span className="font-medium">{formatAmount(record.ami * 1.5)}</span>
				</div>
				<div>
					150% MFI: &nbsp;
					<span className="font-medium">{formatAmount(record.mfi * 1.5)}</span>
				</div>
			</div>
		</div>
	);
};
export default Income;
