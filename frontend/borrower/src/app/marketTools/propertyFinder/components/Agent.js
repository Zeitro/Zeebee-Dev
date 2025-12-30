import { Popover, Typography } from 'antd';
import { useState } from 'react';
import { CheckCircleOutlined, CloseOutlined, CopyOutlined } from '@ant-design/icons';
import { formatPhone } from '@repo/shared-methods/utils';
const { Paragraph } = Typography;

export const getCopyContent = (content, options = {}) => {
	return (
		<Paragraph
			copyable={{
				icon: [<CopyOutlined key="copy-icon" />, <CheckCircleOutlined key="copied-icon" />],
				tooltips: ['click here', 'Copied!'],
				text: content,
			}}
			ellipsis={{
				rows: 1,
				tooltip: {
					title: content,
					placement: 'bottom',
				},
			}}
			className={['!mb-0', options.className].join(' ')}
		>
			{content}
		</Paragraph>
	);
};
const Agent = ({ record }) => {
	const [open, setOpen] = useState(false);

	return (
		<Popover
			open={open}
			onOpenChange={(open) => setOpen(open)}
			content={
				<div className="w-[280px]">
					<div className="flex justify-between items-center pb-2 border-b border-base-border">
						<span className="text-base font-medium">SELLER AGENT</span>
						<CloseOutlined className="cursor-pointer" onClick={() => setOpen(false)} />
					</div>
					{record.agentname &&
						getCopyContent(record.agentname, {
							className: 'text-base pt-1 font-medium',
						})}
					{record.agentphone &&
						getCopyContent(formatPhone(record.agentphone), {
							className: 'text-base pt-1',
						})}
					{record.agentemail &&
						getCopyContent(record.agentemail, {
							className: 'text-base pt-1',
						})}
				</div>
			}
			trigger="click"
			placement="right"
		>
			<span className="text-base-primary">{record.agentname}</span>
		</Popover>
	);
};
export default Agent;
