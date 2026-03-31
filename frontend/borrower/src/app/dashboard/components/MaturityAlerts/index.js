import { Button } from 'antd';
import MaturityAlertsItem from './Item';
import ListModal from './ListModal';
import { useState } from 'react';

const MaturityAlerts = ({ list }) => {
	const [open, setOpen] = useState(false);
	return (
		<div className="border rounded-lg border-base-border w-full lg:w-1/3">
			<div className="flex items-center justify-between">
				<div className="flex">
					Maturity alerts<div className="text-base-error">5 loans</div>
				</div>
				<Button type="text" onClick={() => setOpen(true)}>
					View all
				</Button>
			</div>
			<MaturityAlertsItem />
			{open && <ListModal setOpen={setOpen} list={list} />}
		</div>
	);
};
export default MaturityAlerts;
