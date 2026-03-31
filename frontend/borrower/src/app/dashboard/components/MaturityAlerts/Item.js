import { Button } from 'antd';

const MaturityAlertsItem = ({ item }) => {
	return (
		<div className="border rounded-lg border-base-border">
			<div className="flex items-center justify-between">
				<div className="flex">
					LN-100
					<div>3rd party</div>
				</div>

				<span>12/20/2026</span>
			</div>
			<div className="flex items-center justify-between">
				<span>883 cn xxxxxxxxxxxxx</span>
				<div>Due in 9days</div>
			</div>
		</div>
	);
};
export default MaturityAlertsItem;
